package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	eventInspectionContractPath      = "engineering/implementation/event-inspection/executable-contract.json"
	eventInspectionFixturesPath      = "engineering/implementation/event-inspection/fixtures.json"
	eventInspectionManifestPath      = "work/lots/E10-INV-004A-CONTRACT/manifest.json"
	eventInspectionCapabilityPath    = "cmdr-product-spec/07-investigate/modules/signals-and-hunt/capabilities/event-inspection-and-pivot.md"
	eventInspectionRuntimeBoundaryID = "event-inspection-runtime"
	eventInspectionRuntimeRoot       = "product-runtime/event-inspection/**"
	eventInspectionRuntimeMarker     = "product-runtime/event-inspection/README.md"
)

type EventInspectionRuntimeBoundary struct {
	ID                          string   `json:"id"`
	Root                        string   `json:"root"`
	ExpectedState               string   `json:"expected_state"`
	ExternalRuntimeDependencies []string `json:"external_runtime_dependencies"`
}

type EventInspectionScopePolicy struct {
	TenantRequired                    bool `json:"tenant_required"`
	EnvironmentRequired               bool `json:"environment_required"`
	EventRefRequired                  bool `json:"event_ref_required"`
	WildcardTenantForbidden           bool `json:"wildcard_tenant_forbidden"`
	CrossTenantForbidden              bool `json:"cross_tenant_forbidden"`
	EventTenantMustMatchSelectedTenant bool `json:"event_tenant_must_match_selected_tenant"`
}

type EventInspectionPermissionPolicy struct {
	ServerEvaluationRequired                 bool     `json:"server_evaluation_required"`
	ClientTrustForbidden                     bool     `json:"client_trust_forbidden"`
	RequiredProjectionPermissions            []string `json:"required_projection_permissions"`
	RawAccessDecisionRequired                bool     `json:"raw_access_decision_required"`
	RenderedAccessDecisionRequired           bool     `json:"rendered_access_decision_required"`
	RawDenialMasksRawPayload                 bool     `json:"raw_denial_masks_raw_payload"`
	RawDenialMasksRawDerivedSensitiveValues  bool     `json:"raw_denial_masks_raw_derived_sensitive_values"`
	RenderedMayRemainWhenRawDenied           bool     `json:"rendered_may_remain_when_raw_denied"`
	NoNewPermissionNamespace                 bool     `json:"no_new_permission_namespace"`
}

type EventInspectionProjectionPolicy struct {
	TelemetryEventImmutable       bool `json:"telemetry_event_immutable"`
	SourceFieldsDistinguished     bool `json:"source_fields_distinguished"`
	NormalizedFieldsDistinguished bool `json:"normalized_fields_distinguished"`
	EnrichmentsDistinguished      bool `json:"enrichments_distinguished"`
	ParserStatusVisible           bool `json:"parser_status_visible"`
	MissingFieldsVisible          bool `json:"missing_fields_visible"`
	ProvenanceVisible             bool `json:"provenance_visible"`
	SourceUnavailableStateVisible bool `json:"source_unavailable_state_visible"`
	TombstoneStateVisible         bool `json:"tombstone_state_visible"`
	PartialStateVisible           bool `json:"partial_state_visible"`
	RestrictedStateVisible        bool `json:"restricted_state_visible"`
}

type EventInspectionEnrichmentPolicy struct {
	ProducerRequired             bool `json:"producer_required"`
	VersionRequired              bool `json:"version_required"`
	FreshnessRequired            bool `json:"freshness_required"`
	MayNotMasqueradeAsSource     bool `json:"may_not_masquerade_as_source"`
	AbsentEnrichmentNotInferred  bool `json:"absent_enrichment_not_inferred"`
}

type EventInspectionPivotPolicy struct {
	DraftOnly                 bool     `json:"draft_only"`
	DialectNeutral            bool     `json:"dialect_neutral"`
	FinalQueryDialectSelected bool     `json:"final_query_dialect_selected"`
	FinalProviderSelected     bool     `json:"final_provider_selected"`
	PreservesContext          []string `json:"preserves_context"`
}

type EventInspectionMutationPolicy struct {
	SourceEventMutationForbidden          bool `json:"source_event_mutation_forbidden"`
	CaseLinkMutationExecutable            bool `json:"case_link_mutation_executable"`
	ArtifactProposalMutationExecutable    bool `json:"artifact_proposal_mutation_executable"`
	EvidenceCandidateMutationExecutable   bool `json:"evidence_candidate_mutation_executable"`
	AutomaticEvidenceCreationForbidden    bool `json:"automatic_evidence_creation_forbidden"`
	AutomaticFindingCreationForbidden     bool `json:"automatic_finding_creation_forbidden"`
}

type EventInspectionTechnologyPolicy struct {
	ParserEngineSelected          bool `json:"parser_engine_selected"`
	EntityResolutionEngineSelected bool `json:"entity_resolution_engine_selected"`
	StorageEngineSelected         bool `json:"storage_engine_selected"`
	ProviderSelected              bool `json:"provider_selected"`
	FinalQueryDialectSelected     bool `json:"final_query_dialect_selected"`
}

type EventInspectionExecutableContract struct {
	SchemaVersion        int                              `json:"schema_version"`
	ContractID           string                           `json:"contract_id"`
	ContractKind         string                           `json:"contract_kind"`
	CanonicalSchemaClaim bool                             `json:"canonical_schema_claim"`
	Capability           string                           `json:"capability"`
	RuntimeBoundary      EventInspectionRuntimeBoundary   `json:"runtime_boundary"`
	ProductRefs          ProductRefsV2                    `json:"product_refs"`
	ScopePolicy          EventInspectionScopePolicy       `json:"scope_policy"`
	PermissionPolicy     EventInspectionPermissionPolicy  `json:"permission_policy"`
	ProjectionPolicy     EventInspectionProjectionPolicy  `json:"projection_policy"`
	EnrichmentPolicy     EventInspectionEnrichmentPolicy  `json:"enrichment_policy"`
	PivotPolicy          EventInspectionPivotPolicy       `json:"pivot_policy"`
	MutationPolicy       EventInspectionMutationPolicy    `json:"mutation_policy"`
	TechnologyPolicy     EventInspectionTechnologyPolicy  `json:"technology_policy"`
	Exclusions           []string                         `json:"exclusions"`
	SemanticRules        []string                         `json:"semantic_rules"`
	SecurityInvariants   []string                         `json:"security_invariants"`
}

type EventInspectionFixtureInput struct {
	TenantRef                     string   `json:"tenant_ref"`
	AuthorizedTenants             []string `json:"authorized_tenants"`
	EnvironmentRef                string   `json:"environment_ref"`
	EventRef                      string   `json:"event_ref"`
	EventTenantRef                string   `json:"event_tenant_ref"`
	Permissions                   []string `json:"permissions"`
	RawAccessDecision             string   `json:"raw_access_decision"`
	RenderedAccessDecision        string   `json:"rendered_access_decision"`
	RawPayloadPresent             bool     `json:"raw_payload_present"`
	RawDerivedSensitivePresent    bool     `json:"raw_derived_sensitive_present"`
	RenderedPayloadPresent        bool     `json:"rendered_payload_present"`
	EnrichmentsPresent            bool     `json:"enrichments_present"`
	EnrichmentProducer            string   `json:"enrichment_producer"`
	EnrichmentVersion             string   `json:"enrichment_version"`
	EnrichmentFreshness           string   `json:"enrichment_freshness"`
	EnrichmentPresentedAsSource   bool     `json:"enrichment_presented_as_source"`
	PivotRequested                bool     `json:"pivot_requested"`
	PivotField                    string   `json:"pivot_field"`
	PivotValuePresent             bool     `json:"pivot_value_present"`
	PivotTimeStart                string   `json:"pivot_time_start"`
	PivotTimeEnd                  string   `json:"pivot_time_end"`
	ReturnOrigin                  string   `json:"return_origin"`
	MutationRequested             string   `json:"mutation_requested"`
}

type EventInspectionFixtureExpected struct {
	Allowed                    bool   `json:"allowed"`
	ErrorCode                  string `json:"error_code"`
	RawVisible                 bool   `json:"raw_visible"`
	RenderedVisible            bool   `json:"rendered_visible"`
	RawDerivedSensitiveVisible bool   `json:"raw_derived_sensitive_visible"`
	SourceFieldsDistinguished  bool   `json:"source_fields_distinguished"`
	MissingFieldsVisible       bool   `json:"missing_fields_visible"`
	EnrichmentsDistinguished   bool   `json:"enrichments_distinguished"`
	ParserStatusVisible        bool   `json:"parser_status_visible"`
	ProvenanceVisible          bool   `json:"provenance_visible"`
	PivotDraftCreated          bool   `json:"pivot_draft_created"`
	PivotDialectSelected       bool   `json:"pivot_dialect_selected"`
	ReturnContextPreserved     bool   `json:"return_context_preserved"`
	MutationExecuted           bool   `json:"mutation_executed"`
	AuditRequired              bool   `json:"audit_required"`
}

type EventInspectionFixture struct {
	ID       string                         `json:"id"`
	Input    EventInspectionFixtureInput    `json:"input"`
	Expected EventInspectionFixtureExpected `json:"expected"`
}

type EventInspectionFixtureSet struct {
	SchemaVersion int                      `json:"schema_version"`
	ContractID    string                   `json:"contract_id"`
	Cases         []EventInspectionFixture `json:"cases"`
}

type EventInspectionContractAuditSummary struct {
	ContractID          string `json:"contract_id"`
	Capability          string `json:"capability"`
	RuntimeBoundary     string `json:"runtime_boundary"`
	RuntimeState        string `json:"runtime_state"`
	Fixtures            int    `json:"fixtures"`
	PositiveFixtures    int    `json:"positive_fixtures"`
	NegativeFixtures    int    `json:"negative_fixtures"`
	RuntimeDependencies int    `json:"runtime_dependencies"`
	Status              string `json:"status"`
}

func runEventInspectionContractAudit(root string) (EventInspectionContractAuditSummary, error) {
	var contract EventInspectionExecutableContract
	if err := decodeStrict(root, eventInspectionContractPath, &contract); err != nil {
		return EventInspectionContractAuditSummary{}, err
	}
	var fixtures EventInspectionFixtureSet
	if err := decodeStrict(root, eventInspectionFixturesPath, &fixtures); err != nil {
		return EventInspectionContractAuditSummary{}, err
	}
	manifest, err := decodeWorkManifestV2(root, eventInspectionManifestPath)
	if err != nil {
		return EventInspectionContractAuditSummary{}, err
	}
	if err := validateEventInspectionContractCore(root, manifest, contract, fixtures, true); err != nil {
		return EventInspectionContractAuditSummary{}, err
	}

	var registry ArchitectureRegistry
	if err := decodeStrict(root, architectureRegistryPath, &registry); err != nil {
		return EventInspectionContractAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return EventInspectionContractAuditSummary{}, err
	}
	boundary, ok := architectureBoundaryByID(registry, eventInspectionRuntimeBoundaryID)
	if !ok || boundary.Kind != "product-runtime" || !sameStringSet(boundary.Roots, []string{eventInspectionRuntimeRoot}) {
		return EventInspectionContractAuditSummary{}, fmt.Errorf("Event Inspection runtime boundary is missing or inconsistent")
	}
	runtimeState, err := validateEventInspectionRuntimeLayout(root, boundary, contract.RuntimeBoundary.ExpectedState)
	if err != nil {
		return EventInspectionContractAuditSummary{}, err
	}
	deps, unsupported, err := scanRuntimeBoundary(root, boundary.ID, eventInspectionRuntimeRoot)
	if err != nil {
		return EventInspectionContractAuditSummary{}, err
	}
	if len(deps) != 0 || len(unsupported) != 0 {
		return EventInspectionContractAuditSummary{}, fmt.Errorf("Event Inspection contract phase introduced runtime dependencies: deps=%v unsupported=%v", deps, unsupported)
	}

	negative := 0
	for _, fixture := range fixtures.Cases {
		if !fixture.Expected.Allowed {
			negative++
		}
	}
	return EventInspectionContractAuditSummary{
		ContractID: contract.ContractID, Capability: contract.Capability,
		RuntimeBoundary: boundary.ID, RuntimeState: runtimeState,
		Fixtures: len(fixtures.Cases), PositiveFixtures: len(fixtures.Cases) - negative,
		NegativeFixtures: negative, RuntimeDependencies: len(deps), Status: "PASS",
	}, nil
}

func validateEventInspectionContractCore(root string, manifest WorkManifestV2, contract EventInspectionExecutableContract, fixtures EventInspectionFixtureSet, validateSources bool) error {
	if contract.SchemaVersion != 1 || contract.ContractID != "EVENT-INSPECTION-READONLY-CONTRACT-V1" ||
		contract.ContractKind != "provider-neutral-readonly-local-executable-contract" || contract.CanonicalSchemaClaim {
		return fmt.Errorf("Event Inspection contract identity/schema claim is invalid")
	}
	if contract.Capability != "CAP-INV-004" {
		return fmt.Errorf("Event Inspection capability mismatch")
	}
	if contract.RuntimeBoundary.ID != eventInspectionRuntimeBoundaryID || contract.RuntimeBoundary.Root != eventInspectionRuntimeRoot ||
		(contract.RuntimeBoundary.ExpectedState != "preimplementation" && contract.RuntimeBoundary.ExpectedState != "implemented") ||
		len(contract.RuntimeBoundary.ExternalRuntimeDependencies) != 0 {
		return fmt.Errorf("Event Inspection runtime-boundary contract is invalid")
	}
	if validateSources {
		if err := validateProductRefs(root, contract.ProductRefs); err != nil {
			return err
		}
	}
	if !productRefsEqual(contract.ProductRefs, manifest.ProductRefs) {
		return fmt.Errorf("Event Inspection Product Spec references do not match active work manifest")
	}
	if validateSources {
		if err := validateEventInspectionSourceAnchors(root); err != nil {
			return err
		}
	}

	scope := contract.ScopePolicy
	if !scope.TenantRequired || !scope.EnvironmentRequired || !scope.EventRefRequired ||
		!scope.WildcardTenantForbidden || !scope.CrossTenantForbidden || !scope.EventTenantMustMatchSelectedTenant {
		return fmt.Errorf("Event Inspection scope policy is incomplete or unsafe")
	}
	perm := contract.PermissionPolicy
	if !perm.ServerEvaluationRequired || !perm.ClientTrustForbidden ||
		!sameStringSet(perm.RequiredProjectionPermissions, []string{"perm.shared-capabilities.telemetry-event.read"}) ||
		!perm.RawAccessDecisionRequired || !perm.RenderedAccessDecisionRequired ||
		!perm.RawDenialMasksRawPayload || !perm.RawDenialMasksRawDerivedSensitiveValues ||
		!perm.RenderedMayRemainWhenRawDenied || !perm.NoNewPermissionNamespace {
		return fmt.Errorf("Event Inspection permission policy is incomplete or unsafe")
	}
	projection := contract.ProjectionPolicy
	if !projection.TelemetryEventImmutable || !projection.SourceFieldsDistinguished ||
		!projection.NormalizedFieldsDistinguished || !projection.EnrichmentsDistinguished ||
		!projection.ParserStatusVisible || !projection.MissingFieldsVisible || !projection.ProvenanceVisible ||
		!projection.SourceUnavailableStateVisible || !projection.TombstoneStateVisible ||
		!projection.PartialStateVisible || !projection.RestrictedStateVisible {
		return fmt.Errorf("Event Inspection projection policy is incomplete")
	}
	enrichment := contract.EnrichmentPolicy
	if !enrichment.ProducerRequired || !enrichment.VersionRequired || !enrichment.FreshnessRequired ||
		!enrichment.MayNotMasqueradeAsSource || !enrichment.AbsentEnrichmentNotInferred {
		return fmt.Errorf("Event Inspection enrichment policy is incomplete")
	}
	pivot := contract.PivotPolicy
	if !pivot.DraftOnly || !pivot.DialectNeutral || pivot.FinalQueryDialectSelected || pivot.FinalProviderSelected ||
		!sameStringSet(pivot.PreservesContext, []string{"tenant", "environment", "event_ref", "field", "value", "time_range", "source", "return_origin"}) {
		return fmt.Errorf("Event Inspection pivot policy is invalid or selects unresolved technology")
	}
	mutation := contract.MutationPolicy
	if !mutation.SourceEventMutationForbidden || mutation.CaseLinkMutationExecutable ||
		mutation.ArtifactProposalMutationExecutable || mutation.EvidenceCandidateMutationExecutable ||
		!mutation.AutomaticEvidenceCreationForbidden || !mutation.AutomaticFindingCreationForbidden {
		return fmt.Errorf("Event Inspection mutation policy violates bounded read-only scope")
	}
	tech := contract.TechnologyPolicy
	if tech.ParserEngineSelected || tech.EntityResolutionEngineSelected || tech.StorageEngineSelected ||
		tech.ProviderSelected || tech.FinalQueryDialectSelected {
		return fmt.Errorf("Event Inspection contract selects unresolved runtime technology")
	}
	for _, exclusion := range []string{
		"case-link-mutation:OPEN-013",
		"artifact-proposal-mutation:OPEN-013+OPEN-014",
		"evidence-candidate-mutation:OPEN-013+OPEN-014",
		"final-event-inspector-ui", "parser-engine", "entity-resolution-engine",
		"storage-engine", "provider", "final-query-dialect",
	} {
		if !containsString(contract.Exclusions, exclusion) {
			return fmt.Errorf("Event Inspection exclusion %s is missing", exclusion)
		}
	}
	if len(contract.SemanticRules) < 8 || len(contract.SecurityInvariants) < 6 {
		return fmt.Errorf("Event Inspection semantic/security contract is incomplete")
	}
	if fixtures.SchemaVersion != 1 || fixtures.ContractID != contract.ContractID {
		return fmt.Errorf("Event Inspection fixtures do not bind the executable contract")
	}
	return validateEventInspectionFixtures(contract, fixtures)
}

func validateEventInspectionSourceAnchors(root string) error {
	sources := []struct {
		path    string
		anchors []string
	}{
		{eventInspectionCapabilityPath, []string{
			"id: CAP-INV-004", "raw et rendered", "ne pas modifier l’Event source",
			"OPEN-014", "Case annotation/link", "Pivot draft",
		}},
		{eventInspectionRuntimeMarker, []string{
			"Preimplementation only", "zero raw leakage when raw is denied", "OPEN-013", "OPEN-014",
		}},
	}
	for _, source := range sources {
		data, err := readRepoFile(root, source.path)
		if err != nil {
			return err
		}
		content := string(data)
		for _, anchor := range source.anchors {
			if !strings.Contains(content, anchor) {
				return fmt.Errorf("Event Inspection source anchor %q missing from %s", anchor, source.path)
			}
		}
	}
	return nil
}

func validateEventInspectionFixtures(contract EventInspectionExecutableContract, fixtures EventInspectionFixtureSet) error {
	required := map[string]bool{
		"valid-raw-and-rendered-projection": false,
		"raw-denied-rendered-remains-visible": false,
		"enrichment-provenance-remains-derived": false,
		"valid-dialect-neutral-pivot": false,
		"missing-tenant-is-rejected": false,
		"wildcard-tenant-is-rejected": false,
		"cross-tenant-is-rejected": false,
		"event-tenant-mismatch-is-rejected": false,
		"missing-event-reference-is-rejected": false,
		"permission-denied-masks-projection": false,
		"both-projections-denied-is-rejected": false,
		"enrichment-missing-provenance-is-rejected": false,
		"enrichment-masquerading-as-source-is-rejected": false,
		"invalid-pivot-context-is-rejected": false,
		"case-link-mutation-remains-excluded": false,
		"artifact-proposal-mutation-remains-excluded": false,
		"evidence-candidate-mutation-remains-excluded": false,
	}
	if len(fixtures.Cases) < len(required) {
		return fmt.Errorf("Event Inspection fixture set is too small: %d", len(fixtures.Cases))
	}
	ids := map[string]bool{}
	for _, fixture := range fixtures.Cases {
		if strings.TrimSpace(fixture.ID) == "" || ids[fixture.ID] {
			return fmt.Errorf("Event Inspection fixture id is empty or duplicate: %q", fixture.ID)
		}
		ids[fixture.ID] = true
		if _, ok := required[fixture.ID]; ok {
			required[fixture.ID] = true
		}
		want := evaluateEventInspectionFixture(contract, fixture.Input)
		if fixture.Expected != want {
			return fmt.Errorf("Event Inspection fixture %s expectation mismatch: want %#v got %#v", fixture.ID, want, fixture.Expected)
		}
	}
	for id, present := range required {
		if !present {
			return fmt.Errorf("required Event Inspection fixture %s is missing", id)
		}
	}
	return nil
}

func evaluateEventInspectionFixture(contract EventInspectionExecutableContract, in EventInspectionFixtureInput) EventInspectionFixtureExpected {
	out := EventInspectionFixtureExpected{AuditRequired: true}
	reject := func(code string) EventInspectionFixtureExpected {
		out.ErrorCode = code
		return out
	}
	if strings.TrimSpace(in.MutationRequested) != "" {
		switch in.MutationRequested {
		case "case-link", "artifact-proposal", "evidence-candidate":
			return reject("mutation-excluded")
		default:
			return reject("unsupported-mutation")
		}
	}
	if strings.TrimSpace(in.TenantRef) == "" {
		return reject("missing-tenant")
	}
	if in.TenantRef == "*" {
		return reject("wildcard-tenant")
	}
	if !containsString(in.AuthorizedTenants, in.TenantRef) {
		return reject("cross-tenant")
	}
	if strings.TrimSpace(in.EnvironmentRef) == "" {
		return reject("missing-environment")
	}
	if strings.TrimSpace(in.EventRef) == "" {
		return reject("missing-event-ref")
	}
	if strings.TrimSpace(in.EventTenantRef) == "" || in.EventTenantRef != in.TenantRef {
		return reject("event-tenant-mismatch")
	}
	for _, permission := range contract.PermissionPolicy.RequiredProjectionPermissions {
		if !containsString(in.Permissions, permission) {
			return reject("permission-denied")
		}
	}
	if !validAccessDecision(in.RawAccessDecision) || !validAccessDecision(in.RenderedAccessDecision) {
		return reject("invalid-access-decision")
	}
	if in.RawAccessDecision == "deny" && in.RenderedAccessDecision == "deny" {
		return reject("projection-denied")
	}
	if in.EnrichmentsPresent {
		if strings.TrimSpace(in.EnrichmentProducer) == "" || strings.TrimSpace(in.EnrichmentVersion) == "" ||
			strings.TrimSpace(in.EnrichmentFreshness) == "" {
			return reject("invalid-enrichment-provenance")
		}
		if in.EnrichmentPresentedAsSource {
			return reject("derived-as-source")
		}
	}
	if in.PivotRequested {
		start, errStart := time.Parse(time.RFC3339, in.PivotTimeStart)
		end, errEnd := time.Parse(time.RFC3339, in.PivotTimeEnd)
		if strings.TrimSpace(in.PivotField) == "" || !in.PivotValuePresent ||
			errStart != nil || errEnd != nil || !start.Before(end) || strings.TrimSpace(in.ReturnOrigin) == "" {
			return reject("invalid-pivot-context")
		}
	}
	out.Allowed = true
	out.RawVisible = in.RawAccessDecision == "allow" && in.RawPayloadPresent
	out.RenderedVisible = in.RenderedAccessDecision == "allow" && in.RenderedPayloadPresent
	out.RawDerivedSensitiveVisible = in.RawAccessDecision == "allow" && in.RawDerivedSensitivePresent
	out.SourceFieldsDistinguished = true
	out.MissingFieldsVisible = true
	out.EnrichmentsDistinguished = true
	out.ParserStatusVisible = true
	out.ProvenanceVisible = true
	if in.PivotRequested {
		out.PivotDraftCreated = true
		out.ReturnContextPreserved = true
	}
	return out
}

func validAccessDecision(value string) bool {
	return value == "allow" || value == "deny"
}

func validateEventInspectionRuntimeLayout(root string, boundary ArchitectureBoundary, expected string) (string, error) {
	scanRoot, err := runtimeScanRoot(root, eventInspectionRuntimeRoot)
	if err != nil {
		return "", err
	}
	files := []string{}
	hasGoMod := false
	hasRuntimeGo := false
	// #nosec G703 -- scanRoot is repository-confined by runtimeScanRoot; symlink entries are rejected.
	err = filepath.WalkDir(scanRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("Event Inspection runtime boundary contains symlink: %s", path)
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		files = append(files, rel)
		if rel == "product-runtime/event-inspection/go.mod" {
			hasGoMod = true
		}
		if strings.HasSuffix(rel, ".go") && !strings.HasSuffix(rel, "_test.go") {
			hasRuntimeGo = true
		}
		if expected == "preimplementation" && rel != eventInspectionRuntimeMarker && !strings.HasSuffix(rel, "/.gitkeep") {
			return fmt.Errorf("Event Inspection preimplementation boundary contains unexpected runtime file %s", rel)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(files)
	if !containsString(files, eventInspectionRuntimeMarker) {
		return "", fmt.Errorf("Event Inspection runtime marker is missing")
	}
	actual, err := detectRuntimeBoundaryImplementationState(root, boundary)
	if err != nil {
		return "", err
	}
	if actual != expected {
		return "", fmt.Errorf("Event Inspection runtime state mismatch: contract=%s actual=%s", expected, actual)
	}
	if expected == "implemented" && (!hasGoMod || !hasRuntimeGo) {
		return "", fmt.Errorf("Event Inspection implemented runtime requires go.mod and executable Go source")
	}
	return actual, nil
}
