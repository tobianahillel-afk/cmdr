package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	pilotContractPath      = "engineering/pilot/context-envelope-contract.json"
	pilotFixturesPath      = "engineering/pilot/context-envelope-fixtures.json"
	pilotRuntimeBoundaryID = "pilot-context-envelope-runtime"
	pilotRuntimeRoot       = "product-runtime/context-envelope/**"
	pilotRuntimeMarker     = "product-runtime/context-envelope/README.md"
)

type PilotContextRuntimeBoundary struct {
	ID                          string   `json:"id"`
	Root                        string   `json:"root"`
	PackageIdentity             string   `json:"package_identity"`
	ExternalRuntimeDependencies []string `json:"external_runtime_dependencies"`
}

type PilotReferencePolicy struct {
	Kind                       string `json:"kind"`
	MinBytes                   int    `json:"min_bytes"`
	MaxBytes                   int    `json:"max_bytes"`
	WildcardForbidden          bool   `json:"wildcard_forbidden"`
	WhitespaceForbidden        bool   `json:"whitespace_forbidden"`
	ControlCharactersForbidden bool   `json:"control_characters_forbidden"`
}

type PilotReturnOriginPolicy struct {
	RelativeRouteOnly bool `json:"relative_route_only"`
	MaxBytes          int  `json:"max_bytes"`
	QueryForbidden    bool `json:"query_forbidden"`
	FragmentForbidden bool `json:"fragment_forbidden"`
	SchemeForbidden   bool `json:"scheme_forbidden"`
}

type PilotExecutableContract struct {
	SchemaVersion         int                         `json:"schema_version"`
	ContractID            string                      `json:"contract_id"`
	ContractKind          string                      `json:"contract_kind"`
	CanonicalSchemaClaim  bool                        `json:"canonical_schema_claim"`
	SliceID               string                      `json:"slice_id"`
	Capability            string                      `json:"capability"`
	RuntimeBoundary       PilotContextRuntimeBoundary `json:"runtime_boundary"`
	ProductRefs           ProductRefsV2               `json:"product_refs"`
	ActionClasses         []string                    `json:"action_classes"`
	AuthorizationOutcomes []string                    `json:"authorization_outcomes"`
	FunctionalStates      []string                    `json:"functional_states"`
	ErrorCodes            []string                    `json:"error_codes"`
	ReferencePolicy       PilotReferencePolicy        `json:"reference_policy"`
	ReturnOriginPolicy    PilotReturnOriginPolicy     `json:"return_origin_policy"`
	SemanticRules         []string                    `json:"semantic_rules"`
	SecurityInvariants    []string                    `json:"security_invariants"`
}

type PilotFixtureInput struct {
	SourceProduct                  string `json:"source_product"`
	DestinationProduct             string `json:"destination_product"`
	SourceTenantRef                string `json:"source_tenant_ref"`
	DestinationTenantRef           string `json:"destination_tenant_ref"`
	EnvironmentRef                 string `json:"environment_ref"`
	EnvironmentTenantRef           string `json:"environment_tenant_ref"`
	DestinationRequiresEnvironment bool   `json:"destination_requires_environment"`
	EnvironmentCompatible          bool   `json:"environment_compatible"`
	DestinationAuthorization       string `json:"destination_authorization"`
	ContextFresh                   bool   `json:"context_fresh"`
	ReturnOrigin                   string `json:"return_origin"`
}

type PilotFixtureExpected struct {
	State                     string `json:"state"`
	TenantRef                 string `json:"tenant_ref"`
	EnvironmentRef            string `json:"environment_ref"`
	ReturnOrigin              string `json:"return_origin"`
	EnvironmentChoiceRequired bool   `json:"environment_choice_required"`
	Cleared                   bool   `json:"cleared"`
	ClearedReason             string `json:"cleared_reason"`
	ProtectedRefsVisible      bool   `json:"protected_refs_visible"`
	TraceMinimalOnly          bool   `json:"trace_minimal_only"`
	Error                     string `json:"error"`
}

type PilotFixtureCase struct {
	ID       string               `json:"id"`
	Input    PilotFixtureInput    `json:"input"`
	Expected PilotFixtureExpected `json:"expected"`
}

type PilotFixtureSet struct {
	SchemaVersion int                `json:"schema_version"`
	ContractID    string             `json:"contract_id"`
	Cases         []PilotFixtureCase `json:"cases"`
}

type PilotContractAuditSummary struct {
	ContractID          string `json:"contract_id"`
	SliceID             string `json:"slice_id"`
	RuntimeBoundary     string `json:"runtime_boundary"`
	RuntimeState        string `json:"runtime_state"`
	Fixtures            int    `json:"fixtures"`
	NegativeFixtures    int    `json:"negative_fixtures"`
	RuntimeDependencies int    `json:"runtime_dependencies"`
	RuntimeFiles        int    `json:"runtime_files"`
	Status              string `json:"status"`
}

func runPilotContractAudit(root string) (PilotContractAuditSummary, error) {
	var scope PilotScope
	if err := decodeStrict(root, pilotScopePath, &scope); err != nil {
		return PilotContractAuditSummary{}, err
	}
	var contract PilotExecutableContract
	if err := decodeStrict(root, pilotContractPath, &contract); err != nil {
		return PilotContractAuditSummary{}, err
	}
	var fixtures PilotFixtureSet
	if err := decodeStrict(root, pilotFixturesPath, &fixtures); err != nil {
		return PilotContractAuditSummary{}, err
	}
	if err := validatePilotExecutableContract(scope, contract, fixtures); err != nil {
		return PilotContractAuditSummary{}, err
	}

	var registry ArchitectureRegistry
	if err := decodeStrict(root, architectureRegistryPath, &registry); err != nil {
		return PilotContractAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return PilotContractAuditSummary{}, err
	}
	boundary, ok := architectureBoundaryByID(registry, pilotRuntimeBoundaryID)
	if !ok {
		return PilotContractAuditSummary{}, fmt.Errorf("pilot runtime boundary %s is missing", pilotRuntimeBoundaryID)
	}
	if boundary.Kind != "product-runtime" || !boundary.MutableByImplementation ||
		!sameStringSet(boundary.Roots, []string{pilotRuntimeRoot}) || len(boundary.MayDependOn) != 0 {
		return PilotContractAuditSummary{}, fmt.Errorf("pilot runtime boundary is not isolated as required")
	}

	runtimeFiles, runtimeState, err := validatePilotRuntimeLayout(root, contract.RuntimeBoundary.PackageIdentity)
	if err != nil {
		return PilotContractAuditSummary{}, err
	}
	deps, err := runDependencyAudit(root)
	if err != nil {
		return PilotContractAuditSummary{}, err
	}
	if deps.RuntimeBoundaries != 1 || deps.RuntimeDependencies != 0 || deps.Unapproved != 0 || deps.UnsupportedManifests != 0 {
		return PilotContractAuditSummary{}, fmt.Errorf("pilot runtime dependency floor violated: %#v", deps)
	}

	for _, manifestPath := range []string{
		"work/lots/E9-PILOT-001B/manifest.json",
		"work/lots/E9-PILOT-001C/manifest.json",
	} {
		manifest, err := decodeWorkManifestV2(root, manifestPath)
		if err != nil {
			return PilotContractAuditSummary{}, err
		}
		if err := validatePilotManifestAgainstContract(manifest, scope); err != nil {
			return PilotContractAuditSummary{}, fmt.Errorf("%s: %w", manifest.ID, err)
		}
	}

	negative := 0
	for _, fixture := range fixtures.Cases {
		if fixture.Expected.Error != "" || fixture.Expected.State == "permission-denied" ||
			fixture.Expected.State == "expired" || fixture.Expected.State == "source-missing" ||
			fixture.Expected.State == "context-incompatible" {
			negative++
		}
	}
	return PilotContractAuditSummary{
		ContractID:          contract.ContractID,
		SliceID:             contract.SliceID,
		RuntimeBoundary:     boundary.ID,
		RuntimeState:        runtimeState,
		Fixtures:            len(fixtures.Cases),
		NegativeFixtures:    negative,
		RuntimeDependencies: deps.RuntimeDependencies,
		RuntimeFiles:        runtimeFiles,
		Status:              "PASS",
	}, nil
}

func validatePilotExecutableContract(scope PilotScope, contract PilotExecutableContract, fixtures PilotFixtureSet) error {
	if contract.SchemaVersion != 1 || contract.ContractID != "PILOT-CONTEXT-ENVELOPE-V1" ||
		contract.ContractKind != "pilot-local-executable-contract" || contract.CanonicalSchemaClaim {
		return fmt.Errorf("pilot contract identity/schema claim is invalid")
	}
	if contract.SliceID != scope.SliceID || contract.Capability != scope.Capability {
		return fmt.Errorf("pilot contract does not match selected scope")
	}
	if contract.RuntimeBoundary.ID != pilotRuntimeBoundaryID || contract.RuntimeBoundary.Root != pilotRuntimeRoot ||
		strings.TrimSpace(contract.RuntimeBoundary.PackageIdentity) == "" ||
		len(contract.RuntimeBoundary.ExternalRuntimeDependencies) != 0 {
		return fmt.Errorf("pilot runtime boundary contract is invalid")
	}
	if !sameStringSet(contract.ActionClasses, []string{"0", "1"}) ||
		!sameStringSet(contract.AuthorizationOutcomes, []string{"allow", "deny"}) {
		return fmt.Errorf("pilot contract action/authorization enums are invalid")
	}
	wantStates := []string{"context-valid", "context-incompatible", "permission-denied", "expired", "source-missing"}
	if !sameStringSet(contract.FunctionalStates, wantStates) {
		return fmt.Errorf("pilot functional states mismatch")
	}
	if !sameStringSet(contract.ErrorCodes, []string{"invalid-tenant-ref", "unsafe-return-origin"}) {
		return fmt.Errorf("pilot error codes mismatch")
	}
	if contract.ReferencePolicy.Kind != "opaque-stable-reference" || contract.ReferencePolicy.MinBytes != 1 ||
		contract.ReferencePolicy.MaxBytes < 32 || contract.ReferencePolicy.MaxBytes > 256 ||
		!contract.ReferencePolicy.WildcardForbidden || !contract.ReferencePolicy.WhitespaceForbidden ||
		!contract.ReferencePolicy.ControlCharactersForbidden {
		return fmt.Errorf("pilot reference policy is unsafe")
	}
	if !contract.ReturnOriginPolicy.RelativeRouteOnly || contract.ReturnOriginPolicy.MaxBytes < 128 ||
		contract.ReturnOriginPolicy.MaxBytes > 4096 || !contract.ReturnOriginPolicy.QueryForbidden ||
		!contract.ReturnOriginPolicy.FragmentForbidden || !contract.ReturnOriginPolicy.SchemeForbidden {
		return fmt.Errorf("pilot return-origin policy is unsafe")
	}
	if len(contract.SemanticRules) < 8 || len(contract.SecurityInvariants) < 5 {
		return fmt.Errorf("pilot contract rules/invariants are incomplete")
	}
	wantRefs := pilotProductRefs(scope)
	if !productRefsEqual(contract.ProductRefs, wantRefs) {
		return fmt.Errorf("pilot contract Product Spec references do not match scope")
	}
	if fixtures.SchemaVersion != 1 || fixtures.ContractID != contract.ContractID {
		return fmt.Errorf("pilot fixtures do not bind the selected contract")
	}
	return validatePilotFixtures(contract, fixtures)
}

func validatePilotFixtures(contract PilotExecutableContract, fixtures PilotFixtureSet) error {
	if len(fixtures.Cases) < 10 {
		return fmt.Errorf("pilot fixture set is too small: %d", len(fixtures.Cases))
	}
	required := map[string]bool{
		"preserve-authorized-compatible-environment": false,
		"tenant-change-clears-inherited-environment": false,
		"environment-tenant-mismatch-is-cleared":     false,
		"permission-denied-masks-protected-refs":     false,
		"expired-context-masks-inherited-refs":       false,
		"missing-source-tenant-blocks":               false,
		"wildcard-tenant-is-rejected":                false,
		"unsafe-return-origin-query-is-rejected":     false,
	}
	states := stringSet(contract.FunctionalStates)
	errors := stringSet(contract.ErrorCodes)
	ids := map[string]bool{}
	for _, fixture := range fixtures.Cases {
		if strings.TrimSpace(fixture.ID) == "" || ids[fixture.ID] {
			return fmt.Errorf("pilot fixture id is empty or duplicate: %q", fixture.ID)
		}
		ids[fixture.ID] = true
		if _, ok := required[fixture.ID]; ok {
			required[fixture.ID] = true
		}
		if !containsString(contract.AuthorizationOutcomes, fixture.Input.DestinationAuthorization) {
			return fmt.Errorf("fixture %s has invalid authorization outcome", fixture.ID)
		}
		if fixture.Expected.Error != "" {
			if !errors[fixture.Expected.Error] || fixture.Expected.State != "" ||
				fixture.Expected.TenantRef != "" || fixture.Expected.EnvironmentRef != "" ||
				fixture.Expected.ReturnOrigin != "" || !fixture.Expected.TraceMinimalOnly {
				return fmt.Errorf("fixture %s has unsafe error expectation", fixture.ID)
			}
			continue
		}
		if !states[fixture.Expected.State] {
			return fmt.Errorf("fixture %s has invalid expected state %q", fixture.ID, fixture.Expected.State)
		}
		if fixture.Expected.State == "permission-denied" || fixture.Expected.State == "expired" ||
			fixture.Expected.State == "source-missing" {
			if fixture.Expected.TenantRef != "" || fixture.Expected.EnvironmentRef != "" ||
				fixture.Expected.ProtectedRefsVisible || !fixture.Expected.TraceMinimalOnly {
				return fmt.Errorf("fixture %s leaks protected refs for state %s", fixture.ID, fixture.Expected.State)
			}
		}
		if fixture.Input.SourceTenantRef != fixture.Input.DestinationTenantRef && fixture.Input.SourceTenantRef != "" &&
			fixture.Expected.EnvironmentRef != "" {
			return fmt.Errorf("fixture %s preserves environment across tenant change", fixture.ID)
		}
		if fixture.Input.EnvironmentRef != "" && fixture.Input.EnvironmentTenantRef != fixture.Input.DestinationTenantRef &&
			fixture.Expected.EnvironmentRef != "" {
			return fmt.Errorf("fixture %s preserves cross-tenant environment", fixture.ID)
		}
	}
	for id, present := range required {
		if !present {
			return fmt.Errorf("required pilot fixture %s is missing", id)
		}
	}
	return nil
}

func validatePilotRuntimeLayout(root, packageIdentity string) (int, string, error) {
	scanRoot, err := runtimeScanRoot(root, pilotRuntimeRoot)
	if err != nil {
		return 0, "", err
	}
	files := map[string]bool{}
	hasRuntimeSource := false
	count := 0
	// #nosec G703 -- scanRoot is repository-confined by runtimeScanRoot; symlink entries are rejected.
	err = filepath.WalkDir(scanRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("pilot runtime boundary contains symlink: %s", path)
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		files[rel] = true
		count++
		if strings.HasSuffix(rel, ".go") && !strings.HasSuffix(rel, "_test.go") {
			hasRuntimeSource = true
		}
		return nil
	})
	if err != nil {
		return count, "", err
	}
	if !files[pilotRuntimeMarker] {
		return count, "", fmt.Errorf("pilot runtime README marker is missing")
	}
	if count == 1 {
		return count, "preimplementation", nil
	}

	goMod := "product-runtime/context-envelope/go.mod"
	if !files[goMod] {
		return count, "", fmt.Errorf("implemented pilot runtime is missing go.mod")
	}
	if !hasRuntimeSource {
		return count, "", fmt.Errorf("implemented pilot runtime has no non-test Go source")
	}
	data, err := readRepoFile(root, goMod)
	if err != nil {
		return count, "", err
	}
	moduleIdentity := ""
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "module" {
			if moduleIdentity != "" {
				return count, "", fmt.Errorf("pilot runtime go.mod declares module more than once")
			}
			moduleIdentity = fields[1]
		}
	}
	if moduleIdentity != packageIdentity {
		return count, "", fmt.Errorf("pilot runtime module identity mismatch: expected %s, got %s", packageIdentity, moduleIdentity)
	}
	return count, "implemented", nil
}

func validatePilotManifestAgainstContract(manifest WorkManifestV2, scope PilotScope) error {
	if !productRefsEqual(manifest.ProductRefs, pilotProductRefs(scope)) {
		return fmt.Errorf("product_refs do not match selected pilot scope")
	}
	foundRuntime := false
	for _, path := range manifest.AllowedPaths {
		if patternCovers(path, pilotRuntimeRoot) || patternCovers(pilotRuntimeRoot, path) {
			foundRuntime = true
			break
		}
	}
	if !foundRuntime {
		return fmt.Errorf("manifest does not allow owned pilot runtime root")
	}
	return nil
}

func pilotProductRefs(scope PilotScope) ProductRefsV2 {
	return ProductRefsV2{
		Capabilities:            []string{scope.Capability},
		Requirements:            append([]string(nil), scope.Requirements...),
		Permissions:             append([]string(nil), scope.Permissions...),
		OpenDecisions:           append([]string(nil), scope.OpenDecisions...),
		ImplementationContracts: append([]string(nil), scope.ImplementationContracts...),
		CanonicalObjects:        append([]string(nil), scope.CanonicalObjects...),
	}
}

func productRefsEqual(a, b ProductRefsV2) bool {
	return sameStringSet(a.Capabilities, b.Capabilities) &&
		sameStringSet(a.Requirements, b.Requirements) &&
		sameStringSet(a.Screens, b.Screens) &&
		sameStringSet(a.Permissions, b.Permissions) &&
		sameStringSet(a.OpenDecisions, b.OpenDecisions) &&
		sameStringSet(a.ImplementationContracts, b.ImplementationContracts) &&
		sameStringSet(a.CanonicalObjects, b.CanonicalObjects)
}

func architectureBoundaryByID(registry ArchitectureRegistry, id string) (ArchitectureBoundary, bool) {
	for _, boundary := range registry.Boundaries {
		if boundary.ID == id {
			return boundary, true
		}
	}
	return ArchitectureBoundary{}, false
}

func stringSet(values []string) map[string]bool {
	out := map[string]bool{}
	for _, value := range values {
		out[value] = true
	}
	return out
}

func sortedFixtureIDs(fixtures PilotFixtureSet) []string {
	out := make([]string, 0, len(fixtures.Cases))
	for _, fixture := range fixtures.Cases {
		out = append(out, fixture.ID)
	}
	sort.Strings(out)
	return out
}
