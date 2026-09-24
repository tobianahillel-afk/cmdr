package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	pilotScopePath        = "engineering/pilot/scope.json"
	pilotManifestPath     = "work/lots/E9-PILOT-001A/manifest.json"
	permissionCatalogPath = "cmdr-product-spec/14-security-permissions-and-trust/permission-catalog.md"
)

type PilotRejectedCandidate struct {
	Capability string `json:"capability"`
	Reason     string `json:"reason"`
}

type PilotRuntimeIntent struct {
	Component                   string   `json:"component"`
	BoundaryKind                string   `json:"boundary_kind"`
	ExternalRuntimeDependencies []string `json:"external_runtime_dependencies"`
	AuthorizationModel          string   `json:"authorization_model"`
}

type PilotScope struct {
	SchemaVersion          int                      `json:"schema_version"`
	SliceID                string                   `json:"slice_id"`
	Title                  string                   `json:"title"`
	Capability             string                   `json:"capability"`
	CapabilitySource       string                   `json:"capability_source"`
	Requirements           []string                 `json:"requirements"`
	Permissions            []string                 `json:"permissions"`
	OpenDecisions          []string                 `json:"open_decisions"`
	ImplementationContracts []string                `json:"implementation_contracts"`
	CanonicalObjects       []string                 `json:"canonical_objects"`
	DependencySources      []string                 `json:"dependency_sources"`
	MVPSource              string                   `json:"mvp_source"`
	PilotPlanSource        string                   `json:"pilot_plan_source"`
	DependencyRoadmapSource string                  `json:"dependency_roadmap_source"`
	ActionClasses          []string                 `json:"action_classes"`
	IncludedBehaviors      []string                 `json:"included_behaviors"`
	ExcludedBehaviors      []string                 `json:"excluded_behaviors"`
	RuntimeIntent          PilotRuntimeIntent       `json:"runtime_intent"`
	SelectionRationale     []string                 `json:"selection_rationale"`
	RejectedCandidates     []PilotRejectedCandidate `json:"rejected_candidates"`
}

type PilotScopeAuditSummary struct {
	SliceID             string `json:"slice_id"`
	Capability          string `json:"capability"`
	Requirements        int    `json:"requirements"`
	Permissions         int    `json:"permissions"`
	OpenDecisions       int    `json:"open_decisions"`
	Contracts           int    `json:"implementation_contracts"`
	CanonicalObjects    int    `json:"canonical_objects"`
	DependencySources   int    `json:"dependency_sources"`
	RejectedCandidates  int    `json:"rejected_candidates"`
	ProductGraphDigest  string `json:"product_graph_digest"`
	Status              string `json:"status"`
}

var pilotSliceIDPattern = regexp.MustCompile(`^PILOT-[A-Z0-9-]+$`)

func runPilotScopeAudit(root string, state CurrentState) (PilotScopeAuditSummary, error) {
	var scope PilotScope
	if err := decodeStrict(root, pilotScopePath, &scope); err != nil {
		return PilotScopeAuditSummary{}, err
	}
	if err := validatePilotScopeShape(scope); err != nil {
		return PilotScopeAuditSummary{}, err
	}

	inventory, err := buildSpecInventory(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return PilotScopeAuditSummary{}, err
	}
	activeScreens, err := loadActiveScreenIDs(root, state.ProductSpec.CanonicalPath)
	if err != nil {
		return PilotScopeAuditSummary{}, err
	}
	graph, err := buildProductGraph(inventory, activeScreens)
	if err != nil {
		return PilotScopeAuditSummary{}, err
	}

	var manifest WorkManifestV2
	manifest, err = decodeWorkManifestV2(root, pilotManifestPath)
	if err != nil {
		return PilotScopeAuditSummary{}, err
	}
	if err := validatePilotScopeAgainstProduct(root, scope, manifest, inventory, graph); err != nil {
		return PilotScopeAuditSummary{}, err
	}

	return PilotScopeAuditSummary{
		SliceID:            scope.SliceID,
		Capability:         scope.Capability,
		Requirements:       len(scope.Requirements),
		Permissions:        len(scope.Permissions),
		OpenDecisions:      len(scope.OpenDecisions),
		Contracts:          len(scope.ImplementationContracts),
		CanonicalObjects:   len(scope.CanonicalObjects),
		DependencySources:  len(scope.DependencySources),
		RejectedCandidates: len(scope.RejectedCandidates),
		ProductGraphDigest: graph.SpecTreeDigest,
		Status:             "PASS",
	}, nil
}

func validatePilotScopeShape(scope PilotScope) error {
	if scope.SchemaVersion != 1 {
		return fmt.Errorf("unsupported pilot scope schema_version %d", scope.SchemaVersion)
	}
	if !pilotSliceIDPattern.MatchString(scope.SliceID) {
		return fmt.Errorf("invalid pilot slice id %q", scope.SliceID)
	}
	if strings.TrimSpace(scope.Title) == "" {
		return fmt.Errorf("pilot title is required")
	}
	if !capabilityIDPattern.MatchString(scope.Capability) {
		return fmt.Errorf("pilot capability %q is invalid", scope.Capability)
	}
	if len(scope.OpenDecisions) != 0 {
		return fmt.Errorf("pilot capability must not select unresolved open decisions: %v", scope.OpenDecisions)
	}
	if len(scope.Requirements) == 0 || len(scope.Permissions) == 0 ||
		len(scope.ImplementationContracts) == 0 || len(scope.CanonicalObjects) == 0 ||
		len(scope.DependencySources) == 0 {
		return fmt.Errorf("pilot scope requires requirements, permissions, contracts, canonical objects, and dependency sources")
	}
	for _, permission := range scope.Permissions {
		if !strings.HasPrefix(permission, "perm.") || !strings.HasSuffix(permission, ".read") {
			return fmt.Errorf("pilot permission must be an existing read permission, got %q", permission)
		}
	}
	for _, class := range scope.ActionClasses {
		if class != "0" && class != "1" {
			return fmt.Errorf("pilot action class %q exceeds read/navigation/validation scope", class)
		}
	}
	if scope.RuntimeIntent.BoundaryKind != "product-runtime" {
		return fmt.Errorf("pilot runtime intent must target product-runtime")
	}
	if strings.TrimSpace(scope.RuntimeIntent.Component) == "" || strings.TrimSpace(scope.RuntimeIntent.AuthorizationModel) == "" {
		return fmt.Errorf("pilot runtime intent requires component and authorization model")
	}
	if len(scope.RuntimeIntent.ExternalRuntimeDependencies) != 0 {
		return fmt.Errorf("pilot runtime intent must not introduce external runtime dependencies")
	}
	if len(scope.IncludedBehaviors) == 0 || len(scope.ExcludedBehaviors) == 0 || len(scope.SelectionRationale) == 0 {
		return fmt.Errorf("pilot behavior and selection rationale must be explicit")
	}
	return nil
}

func validatePilotScopeAgainstProduct(root string, scope PilotScope, manifest WorkManifestV2, inventory SpecInventory, graph ProductGraph) error {
	docsByPath := map[string]SpecDocument{}
	for _, doc := range inventory.Documents {
		docsByPath[doc.Path] = doc
	}
	entities := map[string]ProductEntity{}
	for _, entity := range graph.Entities {
		entities[entity.ID] = entity
	}

	capability, ok := entities[scope.Capability]
	if !ok || capability.ReferenceOnly || capability.Kind != "capability" {
		return fmt.Errorf("pilot capability %s is not an owned Product Spec capability", scope.Capability)
	}
	if capability.SourcePath != scope.CapabilitySource {
		return fmt.Errorf("pilot capability source mismatch: graph=%s scope=%s", capability.SourcePath, scope.CapabilitySource)
	}
	capDoc, ok := docsByPath[scope.CapabilitySource]
	if !ok || !capDoc.Active || !capDoc.Canonical {
		return fmt.Errorf("pilot capability source is not active canonical Product Spec")
	}
	if len(capDoc.OpenDecisions) != 0 {
		return fmt.Errorf("pilot capability %s has unresolved open decisions: %v", scope.Capability, capDoc.OpenDecisions)
	}
	if !sameStringSet(scope.Requirements, capDoc.Requirements) {
		return fmt.Errorf("pilot requirements do not exactly match capability requirements: scope=%v capability=%v", scope.Requirements, capDoc.Requirements)
	}

	for _, id := range scope.Requirements {
		entity, ok := entities[id]
		if !ok || entity.Kind != "requirement" {
			return fmt.Errorf("pilot requirement %s does not resolve in product graph", id)
		}
	}
	for _, id := range scope.Permissions {
		entity, ok := entities[id]
		if !ok || entity.Kind != "permission" {
			return fmt.Errorf("pilot permission %s does not resolve in product graph", id)
		}
	}
	for _, path := range scope.ImplementationContracts {
		if err := requireOwnedSourcePath(path, "implementation-contract", docsByPath, graph); err != nil {
			return err
		}
	}
	for _, path := range scope.CanonicalObjects {
		if err := requireOwnedSourcePath(path, "canonical-object", docsByPath, graph); err != nil {
			return err
		}
	}
	for _, path := range append(append([]string(nil), scope.DependencySources...), scope.MVPSource, scope.PilotPlanSource, scope.DependencyRoadmapSource) {
		doc, ok := docsByPath[path]
		if !ok || !doc.Active {
			return fmt.Errorf("pilot dependency source %s is not active in Product Spec inventory", path)
		}
	}

	mvp, err := readRepoFile(root, scope.MVPSource)
	if err != nil {
		return err
	}
	if !strings.Contains(string(mvp), "Canonical objects and shell.") {
		return fmt.Errorf("MVP source no longer contains canonical objects and shell scope")
	}
	pilotPlan, err := readRepoFile(root, scope.PilotPlanSource)
	if err != nil {
		return err
	}
	if !strings.Contains(string(pilotPlan), "Tenants/environments.") {
		return fmt.Errorf("pilot plan no longer contains Tenants/environments scope")
	}
	catalog, err := readRepoFile(root, permissionCatalogPath)
	if err != nil {
		return err
	}
	for _, permission := range scope.Permissions {
		if !strings.Contains(string(catalog), "`"+permission+"`") {
			return fmt.Errorf("pilot permission %s is absent from permission catalog", permission)
		}
	}

	want := ProductRefsV2{
		Capabilities:            []string{scope.Capability},
		Requirements:            append([]string(nil), scope.Requirements...),
		Permissions:             append([]string(nil), scope.Permissions...),
		OpenDecisions:           append([]string(nil), scope.OpenDecisions...),
		ImplementationContracts: append([]string(nil), scope.ImplementationContracts...),
		CanonicalObjects:        append([]string(nil), scope.CanonicalObjects...),
	}
	if !sameStringSet(manifest.ProductRefs.Capabilities, want.Capabilities) ||
		!sameStringSet(manifest.ProductRefs.Requirements, want.Requirements) ||
		!sameStringSet(manifest.ProductRefs.Permissions, want.Permissions) ||
		!sameStringSet(manifest.ProductRefs.OpenDecisions, want.OpenDecisions) ||
		!sameStringSet(manifest.ProductRefs.ImplementationContracts, want.ImplementationContracts) ||
		!sameStringSet(manifest.ProductRefs.CanonicalObjects, want.CanonicalObjects) ||
		len(manifest.ProductRefs.Screens) != 0 {
		return fmt.Errorf("E9-PILOT-001A manifest product_refs do not match selected pilot scope")
	}
	return nil
}

func requireOwnedSourcePath(path, kind string, docs map[string]SpecDocument, graph ProductGraph) error {
	doc, ok := docs[path]
	if !ok || !doc.Active || !doc.Canonical {
		return fmt.Errorf("pilot source %s is not active canonical Product Spec", path)
	}
	for _, entity := range graph.Entities {
		if entity.SourcePath == path && !entity.ReferenceOnly && entity.Kind == kind {
			return nil
		}
	}
	return fmt.Errorf("pilot source %s does not resolve as owned %s in product graph", path, kind)
}

