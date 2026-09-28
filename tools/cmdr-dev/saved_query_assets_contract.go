package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	savedQueryAssetsContractPath      = "engineering/implementation/saved-query-assets/executable-contract.json"
	savedQueryAssetsFixturesPath      = "engineering/implementation/saved-query-assets/fixtures.json"
	savedQueryAssetsManifestPath      = "work/lots/E10-INV-006A-CONTRACT/manifest.json"
	savedQueryAssetsCapabilityPath    = "cmdr-product-spec/07-investigate/modules/signals-and-hunt/capabilities/saved-searches-and-query-assets.md"
	savedQueryAssetsRuntimeBoundaryID = "saved-query-assets-runtime"
	savedQueryAssetsRuntimeRoot       = "product-runtime/saved-query-assets/**"
	savedQueryAssetsRuntimeMarker     = "product-runtime/saved-query-assets/README.md"
)

type QueryAssetRuntimeBoundary struct {
	ID                          string   `json:"id"`
	Root                        string   `json:"root"`
	ExpectedState               string   `json:"expected_state"`
	ExternalRuntimeDependencies []string `json:"external_runtime_dependencies"`
}

type QueryAssetScopePolicy struct {
	TenantRequired          bool `json:"tenant_required"`
	EnvironmentRequired     bool `json:"environment_required"`
	WildcardTenantForbidden bool `json:"wildcard_tenant_forbidden"`
	CrossTenantForbidden    bool `json:"cross_tenant_forbidden"`
}

type QueryAssetPolicy struct {
	AllowedKinds                 []string `json:"allowed_kinds"`
	StableIdentifierRequired     bool     `json:"stable_identifier_required"`
	CallerAccessDecisionRequired bool     `json:"caller_access_decision_required"`
	DeniedProjectionForbidden    bool     `json:"denied_projection_forbidden"`
	AuthorRequired               bool     `json:"author_required"`
	ValidationStateRequired      bool     `json:"validation_state_required"`
	ParametersPreserved          bool     `json:"parameters_preserved"`
	LineagePreservedWhenPresent  bool     `json:"lineage_preserved_when_present"`
	DeprecationMetadataReadonly  bool     `json:"deprecation_metadata_readonly"`
}

type QueryAssetQueryPolicy struct {
	SharedOwned                  bool `json:"shared_owned"`
	TenantScoped                 bool `json:"tenant_scoped"`
	StableIdentifierRequired     bool `json:"stable_identifier_required"`
	ExactVersionRequired         bool `json:"exact_version_required"`
	CallerAccessDecisionRequired bool `json:"caller_access_decision_required"`
	DeniedProjectionForbidden    bool `json:"denied_projection_forbidden"`
	OwnershipTransferForbidden   bool `json:"ownership_transfer_forbidden"`
}

type QueryAssetCompatibilityPolicy struct {
	SourcePrerequisitesRequired         bool `json:"source_prerequisites_required"`
	SourceTenantMustMatchAsset          bool `json:"source_tenant_must_match_asset"`
	MissingSourceIncompatible           bool `json:"missing_source_incompatible"`
	MissingOrRemovedFieldIncompatible   bool `json:"missing_or_removed_field_incompatible"`
	StalePrerequisiteVisible            bool `json:"stale_prerequisite_visible"`
	StaleValidationVisible              bool `json:"stale_validation_visible"`
	DeterministicSortedDiagnostics      bool `json:"deterministic_sorted_diagnostics"`
	SynthesisOfMissingFactsForbidden    bool `json:"synthesis_of_missing_facts_forbidden"`
}

type QueryAssetProjectionPolicy struct {
	Immutable                        bool `json:"immutable"`
	DeepCopyIsolation                bool `json:"deep_copy_isolation"`
	QueryIdentityVisible             bool `json:"query_identity_visible"`
	ParameterMetadataVisible         bool `json:"parameter_metadata_visible"`
	SourceFieldPrerequisitesVisible  bool `json:"source_field_prerequisites_visible"`
	AuthorValidationProvenanceVisible bool `json:"author_validation_provenance_visible"`
	CompatibilityVisible             bool `json:"compatibility_visible"`
	HandoffOnlyWhenCompatible        bool `json:"handoff_only_when_compatible"`
}

type QueryAssetMutationPolicy struct {
	SavedSearchCreateExecutable       bool `json:"saved_search_create_executable"`
	SavedSearchUpdateExecutable       bool `json:"saved_search_update_executable"`
	ShareExecutable                   bool `json:"share_executable"`
	DuplicateExecutable               bool `json:"duplicate_executable"`
	DeprecateExecutable               bool `json:"deprecate_executable"`
	ArchiveExecutable                 bool `json:"archive_executable"`
	QueryMutationExecutable           bool `json:"query_mutation_executable"`
	SavedViewMutationExecutable       bool `json:"saved_view_mutation_executable"`
	DetectionRuleConversionExecutable bool `json:"detection_rule_conversion_executable"`
	SearchJobCreationExecutable       bool `json:"search_job_creation_executable"`
	SearchJobExecutionExecutable      bool `json:"search_job_execution_executable"`
}

type QueryAssetTechnologyPolicy struct {
	StorageEngineSelected         bool `json:"storage_engine_selected"`
	VersionStoreSelected          bool `json:"version_store_selected"`
	ProviderSelected              bool `json:"provider_selected"`
	RetentionPolicySelected       bool `json:"retention_policy_selected"`
	CollaborationBackendSelected  bool `json:"collaboration_backend_selected"`
	ApprovalPolicySelected        bool `json:"approval_policy_selected"`
	RevalidationPolicySelected    bool `json:"revalidation_policy_selected"`
	FinalQueryDialectSelected     bool `json:"final_query_dialect_selected"`
	FinalUISelected               bool `json:"final_ui_selected"`
}

type QueryAssetExecutableContract struct {
	SchemaVersion        int                           `json:"schema_version"`
	ContractID           string                        `json:"contract_id"`
	ContractKind         string                        `json:"contract_kind"`
	CanonicalSchemaClaim bool                          `json:"canonical_schema_claim"`
	Capability           string                        `json:"capability"`
	RuntimeBoundary      QueryAssetRuntimeBoundary     `json:"runtime_boundary"`
	ProductRefs          ProductRefsV2                 `json:"product_refs"`
	ScopePolicy          QueryAssetScopePolicy         `json:"scope_policy"`
	AssetPolicy          QueryAssetPolicy              `json:"asset_policy"`
	QueryPolicy          QueryAssetQueryPolicy         `json:"query_policy"`
	CompatibilityPolicy  QueryAssetCompatibilityPolicy `json:"compatibility_policy"`
	ProjectionPolicy     QueryAssetProjectionPolicy    `json:"projection_policy"`
	MutationPolicy       QueryAssetMutationPolicy      `json:"mutation_policy"`
	TechnologyPolicy     QueryAssetTechnologyPolicy    `json:"technology_policy"`
	Exclusions           []string                      `json:"exclusions"`
	SemanticRules        []string                      `json:"semantic_rules"`
	SecurityInvariants   []string                      `json:"security_invariants"`
}

type QueryAssetFixtureField struct {
	Name  string `json:"name"`
	State string `json:"state"`
}

type QueryAssetFixtureSource struct {
	ID        string                   `json:"id"`
	TenantRef string                   `json:"tenant_ref"`
	State     string                   `json:"state"`
	Fields    []QueryAssetFixtureField `json:"fields"`
}

type QueryAssetFixtureInput struct {
	TenantRef                   string                    `json:"tenant_ref"`
	AuthorizedTenants           []string                  `json:"authorized_tenants"`
	EnvironmentRef              string                    `json:"environment_ref"`
	AssetID                     string                    `json:"asset_id"`
	AssetKind                   string                    `json:"asset_kind"`
	AssetAccessDecision         string                    `json:"asset_access_decision"`
	QueryRef                    string                    `json:"query_ref"`
	QueryTenantRef              string                    `json:"query_tenant_ref"`
	QueryVersion                string                    `json:"query_version"`
	QueryAccessDecision         string                    `json:"query_access_decision"`
	AuthorRef                   string                    `json:"author_ref"`
	ValidationState             string                    `json:"validation_state"`
	ParameterNames              []string                  `json:"parameter_names"`
	Sources                     []QueryAssetFixtureSource `json:"sources"`
	LineageRef                  string                    `json:"lineage_ref"`
	DeprecationReason           string                    `json:"deprecation_reason"`
	ReplacementAssetRef         string                    `json:"replacement_asset_ref"`
	MutationRequested           string                    `json:"mutation_requested"`
	SearchJobExecutionRequested bool                      `json:"search_job_execution_requested"`
}

type QueryAssetFixtureExpected struct {
	Allowed             bool     `json:"allowed"`
	ErrorCode           string   `json:"error_code"`
	Compatibility       string   `json:"compatibility"`
	Diagnostics         []string `json:"diagnostics"`
	HandoffEligible     bool     `json:"handoff_eligible"`
	ImmutableProjection bool     `json:"immutable_projection"`
	DeepCopyIsolated    bool     `json:"deep_copy_isolated"`
	MutationExecuted    bool     `json:"mutation_executed"`
	SearchJobCreated    bool     `json:"search_job_created"`
	SearchJobExecuted   bool     `json:"search_job_executed"`
	AuditRequired       bool     `json:"audit_required"`
}

type QueryAssetFixtureCase struct {
	ID       string                    `json:"id"`
	Input    QueryAssetFixtureInput    `json:"input"`
	Expected QueryAssetFixtureExpected `json:"expected"`
}

type QueryAssetFixtureSet struct {
	SchemaVersion int                     `json:"schema_version"`
	ContractID    string                  `json:"contract_id"`
	Cases         []QueryAssetFixtureCase `json:"cases"`
}

type QueryAssetContractAuditSummary struct {
	ContractID          string `json:"contract_id"`
	Capability          string `json:"capability"`
	RuntimeBoundary     string `json:"runtime_boundary"`
	RuntimeState        string `json:"runtime_state"`
	Fixtures            int    `json:"fixtures"`
	PositiveFixtures    int    `json:"positive_fixtures"`
	NegativeFixtures    int    `json:"negative_fixtures"`
	IncompatibleCases   int    `json:"incompatible_cases"`
	StaleCases          int    `json:"stale_cases"`
	RuntimeDependencies int    `json:"runtime_dependencies"`
	Status              string `json:"status"`
}

func runSavedQueryAssetsContractAudit(root string) (QueryAssetContractAuditSummary, error) {
	var contract QueryAssetExecutableContract
	if err := decodeStrict(root, savedQueryAssetsContractPath, &contract); err != nil {
		return QueryAssetContractAuditSummary{}, err
	}
	var fixtures QueryAssetFixtureSet
	if err := decodeStrict(root, savedQueryAssetsFixturesPath, &fixtures); err != nil {
		return QueryAssetContractAuditSummary{}, err
	}
	manifest, err := decodeWorkManifestV2(root, savedQueryAssetsManifestPath)
	if err != nil {
		return QueryAssetContractAuditSummary{}, err
	}
	if err := validateSavedQueryAssetsContractCore(root, manifest, contract, fixtures, true); err != nil {
		return QueryAssetContractAuditSummary{}, err
	}

	var registry ArchitectureRegistry
	if err := decodeStrict(root, architectureRegistryPath, &registry); err != nil {
		return QueryAssetContractAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return QueryAssetContractAuditSummary{}, err
	}
	boundary, ok := architectureBoundaryByID(registry, savedQueryAssetsRuntimeBoundaryID)
	if !ok || boundary.Kind != "product-runtime" || !sameStringSet(boundary.Roots, []string{savedQueryAssetsRuntimeRoot}) {
		return QueryAssetContractAuditSummary{}, fmt.Errorf("Saved Query Assets runtime boundary is missing or inconsistent")
	}
	runtimeState, err := validateSavedQueryAssetsRuntimeLayout(root, boundary, contract.RuntimeBoundary.ExpectedState)
	if err != nil {
		return QueryAssetContractAuditSummary{}, err
	}
	deps, unsupported, err := scanRuntimeBoundary(root, boundary.ID, savedQueryAssetsRuntimeRoot)
	if err != nil {
		return QueryAssetContractAuditSummary{}, err
	}
	if len(deps) != 0 || len(unsupported) != 0 {
		return QueryAssetContractAuditSummary{}, fmt.Errorf("Saved Query Assets contract phase introduced runtime dependencies: deps=%v unsupported=%v", deps, unsupported)
	}
	negative, incompatible, stale := 0, 0, 0
	for _, fixture := range fixtures.Cases {
		if !fixture.Expected.Allowed {
			negative++
		}
		switch fixture.Expected.Compatibility {
		case "incompatible":
			incompatible++
		case "stale":
			stale++
		}
	}
	return QueryAssetContractAuditSummary{
		ContractID: contract.ContractID, Capability: contract.Capability,
		RuntimeBoundary: boundary.ID, RuntimeState: runtimeState,
		Fixtures: len(fixtures.Cases), PositiveFixtures: len(fixtures.Cases) - negative,
		NegativeFixtures: negative, IncompatibleCases: incompatible, StaleCases: stale,
		RuntimeDependencies: len(deps), Status: "PASS",
	}, nil
}

func validateSavedQueryAssetsContractCore(root string, manifest WorkManifestV2, contract QueryAssetExecutableContract, fixtures QueryAssetFixtureSet, validateSources bool) error {
	if contract.SchemaVersion != 1 || contract.ContractID != "QUERY-ASSET-READONLY-CONTRACT-V1" ||
		contract.ContractKind != "provider-neutral-readonly-local-executable-contract" || contract.CanonicalSchemaClaim {
		return fmt.Errorf("Saved Query Assets contract identity/schema claim is invalid")
	}
	if contract.Capability != "CAP-INV-006" {
		return fmt.Errorf("Saved Query Assets capability mismatch")
	}
	if contract.RuntimeBoundary.ID != savedQueryAssetsRuntimeBoundaryID ||
		contract.RuntimeBoundary.Root != savedQueryAssetsRuntimeRoot ||
		(contract.RuntimeBoundary.ExpectedState != "preimplementation" && contract.RuntimeBoundary.ExpectedState != "implemented") ||
		len(contract.RuntimeBoundary.ExternalRuntimeDependencies) != 0 {
		return fmt.Errorf("Saved Query Assets runtime-boundary contract is invalid")
	}
	if validateSources {
		if err := validateProductRefs(root, contract.ProductRefs); err != nil {
			return err
		}
	}
	if !productRefsEqual(contract.ProductRefs, manifest.ProductRefs) {
		return fmt.Errorf("Saved Query Assets Product Spec references do not match active work manifest")
	}
	if validateSources {
		if err := validateSavedQueryAssetsSourceAnchors(root); err != nil {
			return err
		}
	}
	scope := contract.ScopePolicy
	if !scope.TenantRequired || !scope.EnvironmentRequired || !scope.WildcardTenantForbidden || !scope.CrossTenantForbidden {
		return fmt.Errorf("Saved Query Assets scope policy is incomplete or unsafe")
	}
	asset := contract.AssetPolicy
	if !sameStringSet(asset.AllowedKinds, []string{"saved-search", "query-asset"}) ||
		!asset.StableIdentifierRequired || !asset.CallerAccessDecisionRequired || !asset.DeniedProjectionForbidden ||
		!asset.AuthorRequired || !asset.ValidationStateRequired || !asset.ParametersPreserved ||
		!asset.LineagePreservedWhenPresent || !asset.DeprecationMetadataReadonly {
		return fmt.Errorf("Saved Query Assets asset policy is incomplete or unsafe")
	}
	query := contract.QueryPolicy
	if !query.SharedOwned || !query.TenantScoped || !query.StableIdentifierRequired || !query.ExactVersionRequired ||
		!query.CallerAccessDecisionRequired || !query.DeniedProjectionForbidden || !query.OwnershipTransferForbidden {
		return fmt.Errorf("Saved Query Assets Query policy is incomplete or unsafe")
	}
	compat := contract.CompatibilityPolicy
	if !compat.SourcePrerequisitesRequired || !compat.SourceTenantMustMatchAsset || !compat.MissingSourceIncompatible ||
		!compat.MissingOrRemovedFieldIncompatible || !compat.StalePrerequisiteVisible || !compat.StaleValidationVisible ||
		!compat.DeterministicSortedDiagnostics || !compat.SynthesisOfMissingFactsForbidden {
		return fmt.Errorf("Saved Query Assets compatibility policy is incomplete or unsafe")
	}
	proj := contract.ProjectionPolicy
	if !proj.Immutable || !proj.DeepCopyIsolation || !proj.QueryIdentityVisible || !proj.ParameterMetadataVisible ||
		!proj.SourceFieldPrerequisitesVisible || !proj.AuthorValidationProvenanceVisible ||
		!proj.CompatibilityVisible || !proj.HandoffOnlyWhenCompatible {
		return fmt.Errorf("Saved Query Assets projection policy is incomplete")
	}
	mut := contract.MutationPolicy
	if mut.SavedSearchCreateExecutable || mut.SavedSearchUpdateExecutable || mut.ShareExecutable ||
		mut.DuplicateExecutable || mut.DeprecateExecutable || mut.ArchiveExecutable || mut.QueryMutationExecutable ||
		mut.SavedViewMutationExecutable || mut.DetectionRuleConversionExecutable ||
		mut.SearchJobCreationExecutable || mut.SearchJobExecutionExecutable {
		return fmt.Errorf("Saved Query Assets mutation policy violates bounded read-only scope")
	}
	tech := contract.TechnologyPolicy
	if tech.StorageEngineSelected || tech.VersionStoreSelected || tech.ProviderSelected || tech.RetentionPolicySelected ||
		tech.CollaborationBackendSelected || tech.ApprovalPolicySelected || tech.RevalidationPolicySelected ||
		tech.FinalQueryDialectSelected || tech.FinalUISelected {
		return fmt.Errorf("Saved Query Assets contract selects unresolved technology")
	}
	for _, exclusion := range []string{
		"saved-search-mutation:OPEN-013", "query-asset-mutation:OPEN-013",
		"query-mutation-or-ownership-transfer", "saved-view-mutation-or-conceptual-merge",
		"detection-rule-creation-conversion-deployment", "search-job-creation-or-execution",
		"cross-tenant-copy", "storage-version-store-provider-retention",
		"collaboration-approval-revalidation-policy", "final-query-dialect", "final-cap-inv-006-ui",
	} {
		if !containsString(contract.Exclusions, exclusion) {
			return fmt.Errorf("Saved Query Assets exclusion %s is missing", exclusion)
		}
	}
	if len(contract.SemanticRules) < 8 || len(contract.SecurityInvariants) < 7 {
		return fmt.Errorf("Saved Query Assets semantic/security contract is incomplete")
	}
	if fixtures.SchemaVersion != 1 || fixtures.ContractID != contract.ContractID {
		return fmt.Errorf("Saved Query Assets fixtures do not bind the executable contract")
	}
	return validateSavedQueryAssetsFixtures(contract, fixtures)
}

func validateSavedQueryAssetsSourceAnchors(root string) error {
	sources := []struct {
		path    string
		anchors []string
	}{
		{savedQueryAssetsCapabilityPath, []string{
			"id: CAP-INV-006",
			"Query reste Shared",
			"Saved Views restent Shared",
			"OPEN-013",
			"Execution handoff",
			"Query inaccessible, source/champ incompatible",
		}},
		{savedQueryAssetsRuntimeMarker, []string{
			"Saved Query Assets runtime boundary",
			"Preimplementation only",
			"OPEN-013",
			"Search Job creation or execution inside this runtime",
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
				return fmt.Errorf("Saved Query Assets source anchor %q missing from %s", anchor, source.path)
			}
		}
	}
	return nil
}

func validateSavedQueryAssetsFixtures(contract QueryAssetExecutableContract, fixtures QueryAssetFixtureSet) error {
	required := map[string]bool{
		"valid-compatible-saved-search": false,
		"valid-compatible-query-asset": false,
		"stale-prerequisites-are-visible": false,
		"missing-source-is-incompatible": false,
		"removed-field-is-incompatible": false,
		"missing-tenant-is-rejected": false,
		"wildcard-tenant-is-rejected": false,
		"cross-tenant-asset-is-rejected": false,
		"missing-environment-is-rejected": false,
		"asset-access-denied-is-rejected": false,
		"missing-query-ref-is-rejected": false,
		"missing-query-version-is-rejected": false,
		"query-access-denied-is-rejected": false,
		"cross-tenant-query-is-rejected": false,
		"cross-tenant-source-is-rejected": false,
		"missing-author-is-rejected": false,
		"class-2-mutation-is-rejected": false,
		"search-job-execution-is-rejected": false,
	}
	if len(fixtures.Cases) < len(required) {
		return fmt.Errorf("Saved Query Assets fixture set is too small: %d", len(fixtures.Cases))
	}
	ids := map[string]bool{}
	for _, fixture := range fixtures.Cases {
		if strings.TrimSpace(fixture.ID) == "" || ids[fixture.ID] {
			return fmt.Errorf("Saved Query Assets fixture id is empty or duplicate: %q", fixture.ID)
		}
		ids[fixture.ID] = true
		if _, ok := required[fixture.ID]; ok {
			required[fixture.ID] = true
		}
		want := evaluateSavedQueryAssetsFixture(contract, fixture.Input)
		if !queryAssetExpectedEqual(fixture.Expected, want) {
			return fmt.Errorf("Saved Query Assets fixture %s expectation mismatch: want %#v got %#v", fixture.ID, want, fixture.Expected)
		}
	}
	for id, present := range required {
		if !present {
			return fmt.Errorf("required Saved Query Assets fixture %s is missing", id)
		}
	}
	return nil
}

func evaluateSavedQueryAssetsFixture(contract QueryAssetExecutableContract, in QueryAssetFixtureInput) QueryAssetFixtureExpected {
	out := QueryAssetFixtureExpected{AuditRequired: true, Diagnostics: []string{}}
	reject := func(code string) QueryAssetFixtureExpected {
		out.ErrorCode = code
		out.Diagnostics = []string{}
		return out
	}
	if strings.TrimSpace(in.MutationRequested) != "" {
		return reject("mutation-open-decision")
	}
	if in.SearchJobExecutionRequested {
		return reject("search-job-execution-forbidden")
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
	if strings.TrimSpace(in.AssetID) == "" {
		return reject("missing-asset-id")
	}
	if !containsString(contract.AssetPolicy.AllowedKinds, in.AssetKind) {
		return reject("invalid-asset-kind")
	}
	if in.AssetAccessDecision != "allow" {
		return reject("asset-access-denied")
	}
	if strings.TrimSpace(in.QueryRef) == "" {
		return reject("missing-query-ref")
	}
	if in.QueryTenantRef != in.TenantRef {
		return reject("cross-tenant-query")
	}
	if strings.TrimSpace(in.QueryVersion) == "" {
		return reject("missing-query-version")
	}
	if in.QueryAccessDecision != "allow" {
		return reject("query-access-denied")
	}
	if strings.TrimSpace(in.AuthorRef) == "" {
		return reject("missing-author")
	}
	switch in.ValidationState {
	case "validated":
	case "stale":
		out.Diagnostics = append(out.Diagnostics, "validation:stale")
	case "incompatible":
		out.Diagnostics = append(out.Diagnostics, "validation:incompatible")
	default:
		return reject("invalid-validation-state")
	}
	for _, parameter := range in.ParameterNames {
		if strings.TrimSpace(parameter) == "" {
			return reject("invalid-parameter")
		}
	}
	if len(in.Sources) == 0 {
		return reject("missing-source-prerequisites")
	}

	incompatible := in.ValidationState == "incompatible"
	stale := in.ValidationState == "stale"
	for _, source := range in.Sources {
		if strings.TrimSpace(source.ID) == "" {
			return reject("invalid-source")
		}
		if source.TenantRef != in.TenantRef {
			return reject("cross-tenant-source")
		}
		switch source.State {
		case "available":
		case "stale":
			stale = true
			out.Diagnostics = append(out.Diagnostics, "source:"+source.ID+":stale")
		case "missing":
			incompatible = true
			out.Diagnostics = append(out.Diagnostics, "source:"+source.ID+":missing")
			continue
		default:
			return reject("invalid-source-state")
		}
		for _, field := range source.Fields {
			if strings.TrimSpace(field.Name) == "" {
				return reject("invalid-field")
			}
			switch field.State {
			case "available":
			case "stale":
				stale = true
				out.Diagnostics = append(out.Diagnostics, "field:"+source.ID+"/"+field.Name+":stale")
			case "missing", "removed":
				incompatible = true
				out.Diagnostics = append(out.Diagnostics, "field:"+source.ID+"/"+field.Name+":"+field.State)
			default:
				return reject("invalid-field-state")
			}
		}
	}
	sort.Strings(out.Diagnostics)
	out.Compatibility = "compatible"
	if incompatible {
		out.Compatibility = "incompatible"
	} else if stale {
		out.Compatibility = "stale"
	}
	out.Allowed = true
	out.HandoffEligible = out.Compatibility == "compatible"
	out.ImmutableProjection = true
	out.DeepCopyIsolated = true
	return out
}

func queryAssetExpectedEqual(a, b QueryAssetFixtureExpected) bool {
	return a.Allowed == b.Allowed && a.ErrorCode == b.ErrorCode &&
		a.Compatibility == b.Compatibility && sameStringSet(a.Diagnostics, b.Diagnostics) &&
		a.HandoffEligible == b.HandoffEligible &&
		a.ImmutableProjection == b.ImmutableProjection && a.DeepCopyIsolated == b.DeepCopyIsolated &&
		a.MutationExecuted == b.MutationExecuted && a.SearchJobCreated == b.SearchJobCreated &&
		a.SearchJobExecuted == b.SearchJobExecuted && a.AuditRequired == b.AuditRequired
}

func validateSavedQueryAssetsRuntimeLayout(root string, boundary ArchitectureBoundary, expected string) (string, error) {
	scanRoot, err := runtimeScanRoot(root, savedQueryAssetsRuntimeRoot)
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
			return fmt.Errorf("Saved Query Assets runtime boundary contains symlink: %s", path)
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
		if rel == "product-runtime/saved-query-assets/go.mod" {
			hasGoMod = true
		}
		if strings.HasSuffix(rel, ".go") && !strings.HasSuffix(rel, "_test.go") {
			hasRuntimeGo = true
		}
		if expected == "preimplementation" && rel != savedQueryAssetsRuntimeMarker && !strings.HasSuffix(rel, "/.gitkeep") {
			return fmt.Errorf("Saved Query Assets preimplementation boundary contains unexpected runtime file %s", rel)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(files)
	if !containsString(files, savedQueryAssetsRuntimeMarker) {
		return "", fmt.Errorf("Saved Query Assets runtime marker is missing")
	}
	actual, err := detectRuntimeBoundaryImplementationState(root, boundary)
	if err != nil {
		return "", err
	}
	if actual != expected {
		return "", fmt.Errorf("Saved Query Assets runtime state mismatch: contract=%s actual=%s", expected, actual)
	}
	if expected == "implemented" && (!hasGoMod || !hasRuntimeGo) {
		return "", fmt.Errorf("Saved Query Assets implemented runtime requires go.mod and executable Go source")
	}
	return actual, nil
}
