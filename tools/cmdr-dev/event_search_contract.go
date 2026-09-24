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
	eventSearchContractPath      = "engineering/implementation/event-search/executable-contract.json"
	eventSearchFixturesPath      = "engineering/implementation/event-search/fixtures.json"
	eventSearchManifestPath      = "work/lots/E10-INV-002A-CONTRACT/manifest.json"
	eventSearchCapabilityPath    = "cmdr-product-spec/07-investigate/modules/signals-and-hunt/capabilities/event-search.md"
	eventSearchScreenPath        = "cmdr-product-spec/07-investigate/modules/event-search/screens/event-search.md"
	eventSearchDecisionLogPath   = "cmdr-product-spec/00-governance/decision-log.md"
	eventSearchRuntimeBoundaryID = "event-search-runtime"
	eventSearchRuntimeRoot       = "product-runtime/event-search/**"
	eventSearchRuntimeMarker     = "product-runtime/event-search/README.md"
)

type EventSearchRuntimeBoundary struct {
	ID                          string   `json:"id"`
	Root                        string   `json:"root"`
	ExpectedState               string   `json:"expected_state"`
	ExternalRuntimeDependencies []string `json:"external_runtime_dependencies"`
}

type EventSearchScopePolicy struct {
	TenantRequired          bool `json:"tenant_required"`
	EnvironmentRequired     bool `json:"environment_required"`
	TimeRangeRequired       bool `json:"time_range_required"`
	SourcesRequired         bool `json:"sources_required"`
	WildcardTenantForbidden bool `json:"wildcard_tenant_forbidden"`
	CrossTenantForbidden    bool `json:"cross_tenant_forbidden"`
	EmptySourcesForbidden   bool `json:"empty_sources_forbidden"`
}

type EventSearchPermissionPolicy struct {
	ServerEvaluationRequired     bool     `json:"server_evaluation_required"`
	ClientTrustForbidden         bool     `json:"client_trust_forbidden"`
	RequiredExecutionPermissions []string `json:"required_execution_permissions"`
	InspectionPermission         string   `json:"inspection_permission"`
	QueryManagementPermission    string   `json:"query_management_permission"`
	DenyMasksProtectedData       bool     `json:"deny_masks_protected_data"`
	RawAccessIsSeparate          bool     `json:"raw_access_is_separate"`
}

type EventSearchQueryPolicy struct {
	QueryRequired                              bool `json:"query_required"`
	ValidationBeforeExecution                  bool `json:"validation_before_execution"`
	DeterministicManualPathRequired            bool `json:"deterministic_manual_path_required"`
	AIExecutionWithoutExplicitTriggerForbidden bool `json:"ai_execution_without_explicit_trigger_forbidden"`
	FinalQueryLanguageSelected                 bool `json:"final_query_language_selected"`
	FinalIndexSelected                         bool `json:"final_index_selected"`
	FinalStorageEngineSelected                 bool `json:"final_storage_engine_selected"`
	FinalProviderSelected                      bool `json:"final_provider_selected"`
}

type EventSearchJobPolicy struct {
	States                       []string `json:"states"`
	TerminalStates               []string `json:"terminal_states"`
	AllowedTransitions           []string `json:"allowed_transitions"`
	RelaunchCreatesNewJob        bool     `json:"relaunch_creates_new_job"`
	CancellationAudited          bool     `json:"cancellation_audited"`
	PartialPreservesValidResults bool     `json:"partial_preserves_valid_results"`
	FailedSourcesVisible         bool     `json:"failed_sources_visible"`
	NoCompleteResultSimulation   bool     `json:"no_complete_result_simulation"`
}

type EventSearchAuditPolicy struct {
	ActorRequired                         bool `json:"actor_required"`
	TenantRequired                        bool `json:"tenant_required"`
	ActionRequired                        bool `json:"action_required"`
	TargetRequired                        bool `json:"target_required"`
	OutcomeRequired                       bool `json:"outcome_required"`
	RationaleRequired                     bool `json:"rationale_required"`
	CorrelationIDRequired                 bool `json:"correlation_id_required"`
	QueryVersionRequired                  bool `json:"query_version_required"`
	PeriodRequired                        bool `json:"period_required"`
	SourcesRequired                       bool `json:"sources_required"`
	RunIDRequired                         bool `json:"run_id_required"`
	ErrorsRecorded                        bool `json:"errors_recorded"`
	SensitiveDataForbiddenInCorrelationID bool `json:"sensitive_data_forbidden_in_correlation_id"`
}

type EventSearchPerformanceBudget struct {
	ValidationP95MS          float64 `json:"validation_p95_ms"`
	OrchestrationP95MS       float64 `json:"orchestration_p95_ms"`
	ProductionSearchSLOClaim bool    `json:"production_search_slo_claim"`
}

type EventSearchExecutableContract struct {
	SchemaVersion        int                          `json:"schema_version"`
	ContractID           string                       `json:"contract_id"`
	ContractKind         string                       `json:"contract_kind"`
	CanonicalSchemaClaim bool                         `json:"canonical_schema_claim"`
	Capability           string                       `json:"capability"`
	Screen               string                       `json:"screen"`
	RuntimeBoundary      EventSearchRuntimeBoundary   `json:"runtime_boundary"`
	ProductRefs          ProductRefsV2                `json:"product_refs"`
	ScopePolicy          EventSearchScopePolicy       `json:"scope_policy"`
	PermissionPolicy     EventSearchPermissionPolicy  `json:"permission_policy"`
	QueryPolicy          EventSearchQueryPolicy       `json:"query_policy"`
	SearchJobPolicy      EventSearchJobPolicy         `json:"search_job_policy"`
	AuditPolicy          EventSearchAuditPolicy       `json:"audit_policy"`
	PerformanceBudget    EventSearchPerformanceBudget `json:"performance_budget"`
	Exclusions           []string                     `json:"exclusions"`
	SemanticRules        []string                     `json:"semantic_rules"`
	SecurityInvariants   []string                     `json:"security_invariants"`
}

type EventSearchFixtureInput struct {
	TenantRef            string   `json:"tenant_ref"`
	AuthorizedTenants    []string `json:"authorized_tenants"`
	EnvironmentRef       string   `json:"environment_ref"`
	TimeStart            string   `json:"time_start"`
	TimeEnd              string   `json:"time_end"`
	Sources              []string `json:"sources"`
	AuthorizedSources    []string `json:"authorized_sources"`
	Permissions          []string `json:"permissions"`
	QueryPresent         bool     `json:"query_present"`
	QueryValid           bool     `json:"query_valid"`
	CurrentJobState      string   `json:"current_job_state"`
	RequestedJobState    string   `json:"requested_job_state"`
	Relaunch             bool     `json:"relaunch"`
	PartialFailedSources []string `json:"partial_failed_sources"`
	CaseLinkRequested    bool     `json:"case_link_requested"`
}

type EventSearchFixtureExpected struct {
	Allowed                  bool   `json:"allowed"`
	ErrorCode                string `json:"error_code"`
	NextJobState             string `json:"next_job_state"`
	ProtectedDataVisible     bool   `json:"protected_data_visible"`
	ValidResultsPreserved    bool   `json:"valid_results_preserved"`
	FailedSourcesVisible     bool   `json:"failed_sources_visible"`
	NewJobRequired           bool   `json:"new_job_required"`
	AuditRequired            bool   `json:"audit_required"`
	CaseLinkMutationExecuted bool   `json:"case_link_mutation_executed"`
}

type EventSearchFixtureCase struct {
	ID       string                     `json:"id"`
	Input    EventSearchFixtureInput    `json:"input"`
	Expected EventSearchFixtureExpected `json:"expected"`
}

type EventSearchFixtureSet struct {
	SchemaVersion int                      `json:"schema_version"`
	ContractID    string                   `json:"contract_id"`
	Cases         []EventSearchFixtureCase `json:"cases"`
}

type EventSearchContractAuditSummary struct {
	ContractID          string `json:"contract_id"`
	Capability          string `json:"capability"`
	RuntimeBoundary     string `json:"runtime_boundary"`
	RuntimeState        string `json:"runtime_state"`
	Fixtures            int    `json:"fixtures"`
	NegativeFixtures    int    `json:"negative_fixtures"`
	PositiveFixtures    int    `json:"positive_fixtures"`
	RuntimeDependencies int    `json:"runtime_dependencies"`
	Status              string `json:"status"`
}

func runEventSearchContractAudit(root string) (EventSearchContractAuditSummary, error) {
	var contract EventSearchExecutableContract
	if err := decodeStrict(root, eventSearchContractPath, &contract); err != nil {
		return EventSearchContractAuditSummary{}, err
	}
	var fixtures EventSearchFixtureSet
	if err := decodeStrict(root, eventSearchFixturesPath, &fixtures); err != nil {
		return EventSearchContractAuditSummary{}, err
	}
	manifest, err := decodeWorkManifestV2(root, eventSearchManifestPath)
	if err != nil {
		return EventSearchContractAuditSummary{}, err
	}
	if err := validateEventSearchContractCore(root, manifest, contract, fixtures, true); err != nil {
		return EventSearchContractAuditSummary{}, err
	}

	var registry ArchitectureRegistry
	if err := decodeStrict(root, architectureRegistryPath, &registry); err != nil {
		return EventSearchContractAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return EventSearchContractAuditSummary{}, err
	}
	boundary, ok := architectureBoundaryByID(registry, eventSearchRuntimeBoundaryID)
	if !ok || boundary.Kind != "product-runtime" || !sameStringSet(boundary.Roots, []string{eventSearchRuntimeRoot}) {
		return EventSearchContractAuditSummary{}, fmt.Errorf("Event Search runtime boundary is missing or inconsistent")
	}
	runtimeState, err := validateEventSearchPreimplementationLayout(root)
	if err != nil {
		return EventSearchContractAuditSummary{}, err
	}
	deps, unsupported, err := scanRuntimeBoundary(root, boundary.ID, eventSearchRuntimeRoot)
	if err != nil {
		return EventSearchContractAuditSummary{}, err
	}
	if len(deps) != 0 || len(unsupported) != 0 {
		return EventSearchContractAuditSummary{}, fmt.Errorf("Event Search contract phase introduced runtime dependencies: deps=%v unsupported=%v", deps, unsupported)
	}

	negative := 0
	for _, fixture := range fixtures.Cases {
		if !fixture.Expected.Allowed {
			negative++
		}
	}
	return EventSearchContractAuditSummary{
		ContractID: contract.ContractID, Capability: contract.Capability,
		RuntimeBoundary: boundary.ID, RuntimeState: runtimeState,
		Fixtures: len(fixtures.Cases), NegativeFixtures: negative,
		PositiveFixtures:    len(fixtures.Cases) - negative,
		RuntimeDependencies: len(deps), Status: "PASS",
	}, nil
}

func validateEventSearchContractCore(root string, manifest WorkManifestV2, contract EventSearchExecutableContract, fixtures EventSearchFixtureSet, validateSources bool) error {
	if contract.SchemaVersion != 1 || contract.ContractID != "EVENT-SEARCH-EXECUTABLE-CONTRACT-V1" ||
		contract.ContractKind != "backend-neutral-local-executable-contract" || contract.CanonicalSchemaClaim {
		return fmt.Errorf("Event Search contract identity/schema claim is invalid")
	}
	if contract.Capability != "CAP-INV-002" || contract.Screen != "INV-EVS-001" {
		return fmt.Errorf("Event Search contract capability/screen mismatch")
	}
	if contract.RuntimeBoundary.ID != eventSearchRuntimeBoundaryID || contract.RuntimeBoundary.Root != eventSearchRuntimeRoot ||
		contract.RuntimeBoundary.ExpectedState != "preimplementation" || len(contract.RuntimeBoundary.ExternalRuntimeDependencies) != 0 {
		return fmt.Errorf("Event Search runtime-boundary contract is invalid")
	}
	if validateSources {
		if err := validateProductRefs(root, contract.ProductRefs); err != nil {
			return err
		}
	}
	if !productRefsEqual(contract.ProductRefs, manifest.ProductRefs) {
		return fmt.Errorf("Event Search contract Product Spec references do not match active work manifest")
	}
	if validateSources {
		if err := validateEventSearchSourceAnchors(root); err != nil {
			return err
		}
	}

	scope := contract.ScopePolicy
	if !scope.TenantRequired || !scope.EnvironmentRequired || !scope.TimeRangeRequired || !scope.SourcesRequired ||
		!scope.WildcardTenantForbidden || !scope.CrossTenantForbidden || !scope.EmptySourcesForbidden {
		return fmt.Errorf("Event Search scope policy is incomplete or unsafe")
	}
	perm := contract.PermissionPolicy
	wantExecution := []string{
		"perm.investigate.search.execute",
		"perm.shared-capabilities.query.read",
		"perm.shared-capabilities.search-job.manage",
		"perm.shared-capabilities.telemetry-event.read",
	}
	if !perm.ServerEvaluationRequired || !perm.ClientTrustForbidden || !perm.DenyMasksProtectedData || !perm.RawAccessIsSeparate ||
		!sameStringSet(perm.RequiredExecutionPermissions, wantExecution) ||
		perm.InspectionPermission != "perm.shared-capabilities.search-job.read" ||
		perm.QueryManagementPermission != "perm.shared-capabilities.query.manage" {
		return fmt.Errorf("Event Search permission policy is incomplete or unsafe")
	}
	query := contract.QueryPolicy
	if !query.QueryRequired || !query.ValidationBeforeExecution || !query.DeterministicManualPathRequired ||
		!query.AIExecutionWithoutExplicitTriggerForbidden || query.FinalQueryLanguageSelected ||
		query.FinalIndexSelected || query.FinalStorageEngineSelected || query.FinalProviderSelected {
		return fmt.Errorf("Event Search query policy selects unresolved technology or weakens deterministic execution")
	}
	job := contract.SearchJobPolicy
	if !sameStringSet(job.States, []string{"queued", "running", "completed", "partial", "failed", "cancelled"}) ||
		!sameStringSet(job.TerminalStates, []string{"completed", "partial", "failed", "cancelled"}) ||
		!sameStringSet(job.AllowedTransitions, []string{
			"queued->running", "queued->cancelled", "running->completed",
			"running->partial", "running->failed", "running->cancelled",
		}) ||
		!job.RelaunchCreatesNewJob || !job.CancellationAudited || !job.PartialPreservesValidResults ||
		!job.FailedSourcesVisible || !job.NoCompleteResultSimulation {
		return fmt.Errorf("Event Search Search Job policy is invalid")
	}
	audit := contract.AuditPolicy
	if !audit.ActorRequired || !audit.TenantRequired || !audit.ActionRequired || !audit.TargetRequired ||
		!audit.OutcomeRequired || !audit.RationaleRequired || !audit.CorrelationIDRequired ||
		!audit.QueryVersionRequired || !audit.PeriodRequired || !audit.SourcesRequired ||
		!audit.RunIDRequired || !audit.ErrorsRecorded || !audit.SensitiveDataForbiddenInCorrelationID {
		return fmt.Errorf("Event Search audit/provenance policy is incomplete")
	}
	if contract.PerformanceBudget.ValidationP95MS <= 0 || contract.PerformanceBudget.ValidationP95MS > 1 ||
		contract.PerformanceBudget.OrchestrationP95MS <= 0 || contract.PerformanceBudget.OrchestrationP95MS > 1 ||
		contract.PerformanceBudget.ProductionSearchSLOClaim {
		return fmt.Errorf("Event Search performance budget is invalid or claims a production Search SLO")
	}
	for _, exclusion := range []string{
		"case-link-mutation:OPEN-013", "final-ui", "final-query-language",
		"final-index", "final-storage-engine", "final-provider",
	} {
		if !containsString(contract.Exclusions, exclusion) {
			return fmt.Errorf("Event Search exclusion %s is missing", exclusion)
		}
	}
	if len(contract.SemanticRules) < 8 || len(contract.SecurityInvariants) < 6 {
		return fmt.Errorf("Event Search semantic/security contract is incomplete")
	}
	if fixtures.SchemaVersion != 1 || fixtures.ContractID != contract.ContractID {
		return fmt.Errorf("Event Search fixtures do not bind the executable contract")
	}
	return validateEventSearchFixtures(contract, fixtures)
}

func validateEventSearchSourceAnchors(root string) error {
	sources := []struct {
		path    string
		anchors []string
	}{
		{eventSearchCapabilityPath, []string{"id: CAP-INV-002", "Lier une sélection", "OPEN-013", "partial"}},
		{eventSearchScreenPath, []string{"id: INV-EVS-001", "perm.investigate.search.execute", "Permission denied"}},
		{eventSearchDecisionLogPath, []string{"OPEN-013 — default governance/authority for reversible class-2 mutations", "OPEN-013", "remains open"}},
	}
	for _, source := range sources {
		data, err := readRepoFile(root, source.path)
		if err != nil {
			return err
		}
		content := string(data)
		for _, anchor := range source.anchors {
			if !strings.Contains(content, anchor) {
				return fmt.Errorf("Event Search source anchor %q missing from %s", anchor, source.path)
			}
		}
	}
	return nil
}

func validateEventSearchFixtures(contract EventSearchExecutableContract, fixtures EventSearchFixtureSet) error {
	required := map[string]bool{
		"valid-execution-envelope":                false,
		"partial-preserves-valid-results":         false,
		"missing-tenant-is-rejected":              false,
		"wildcard-tenant-is-rejected":             false,
		"cross-tenant-is-rejected":                false,
		"permission-denied-masks-protected-data":  false,
		"invalid-time-range-is-rejected":          false,
		"empty-sources-are-rejected":              false,
		"illegal-terminal-transition-is-rejected": false,
		"case-link-remains-excluded":              false,
	}
	if len(fixtures.Cases) < len(required) {
		return fmt.Errorf("Event Search fixture set is too small: %d", len(fixtures.Cases))
	}
	ids := map[string]bool{}
	for _, fixture := range fixtures.Cases {
		if strings.TrimSpace(fixture.ID) == "" || ids[fixture.ID] {
			return fmt.Errorf("Event Search fixture id is empty or duplicate: %q", fixture.ID)
		}
		ids[fixture.ID] = true
		if _, ok := required[fixture.ID]; ok {
			required[fixture.ID] = true
		}
		want := evaluateEventSearchFixture(contract, fixture.Input)
		if fixture.Expected != want {
			return fmt.Errorf("Event Search fixture %s expectation mismatch: want %#v got %#v", fixture.ID, want, fixture.Expected)
		}
	}
	for id, present := range required {
		if !present {
			return fmt.Errorf("required Event Search fixture %s is missing", id)
		}
	}
	return nil
}

func evaluateEventSearchFixture(contract EventSearchExecutableContract, in EventSearchFixtureInput) EventSearchFixtureExpected {
	out := EventSearchFixtureExpected{AuditRequired: true}
	reject := func(code string) EventSearchFixtureExpected {
		out.ErrorCode = code
		return out
	}
	if in.CaseLinkRequested {
		return reject("case-link-open-decision")
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
	start, errStart := time.Parse(time.RFC3339, in.TimeStart)
	end, errEnd := time.Parse(time.RFC3339, in.TimeEnd)
	if errStart != nil || errEnd != nil || !start.Before(end) {
		return reject("invalid-time-range")
	}
	if len(in.Sources) == 0 {
		return reject("empty-sources")
	}
	for _, source := range in.Sources {
		if strings.TrimSpace(source) == "" || !containsString(in.AuthorizedSources, source) {
			return reject("unauthorized-source")
		}
	}
	if !in.QueryPresent || !in.QueryValid {
		return reject("invalid-query")
	}
	for _, permission := range contract.PermissionPolicy.RequiredExecutionPermissions {
		if !containsString(in.Permissions, permission) {
			return reject("permission-denied")
		}
	}
	if in.Relaunch {
		if !containsString(contract.SearchJobPolicy.TerminalStates, in.CurrentJobState) {
			return reject("illegal-transition")
		}
		out.Allowed = true
		out.NextJobState = "queued"
		out.NewJobRequired = true
		out.ProtectedDataVisible = true
		return out
	}
	if in.CurrentJobState == "" {
		if in.RequestedJobState != "queued" {
			return reject("illegal-transition")
		}
	} else if !containsString(contract.SearchJobPolicy.AllowedTransitions, in.CurrentJobState+"->"+in.RequestedJobState) {
		return reject("illegal-transition")
	}
	out.Allowed = true
	out.ProtectedDataVisible = true
	out.NextJobState = in.RequestedJobState
	if len(in.PartialFailedSources) > 0 {
		for _, source := range in.PartialFailedSources {
			if !containsString(in.Sources, source) {
				return reject("invalid-partial-source")
			}
		}
		out.NextJobState = "partial"
		out.ValidResultsPreserved = true
		out.FailedSourcesVisible = true
	}
	return out
}

func validateEventSearchPreimplementationLayout(root string) (string, error) {
	scanRoot, err := runtimeScanRoot(root, eventSearchRuntimeRoot)
	if err != nil {
		return "", err
	}
	files := []string{}
	// #nosec G703 -- scanRoot is repository-confined by runtimeScanRoot; symlink entries are rejected.
	err = filepath.WalkDir(scanRoot, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("Event Search preimplementation boundary contains symlink: %s", path)
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel != eventSearchRuntimeMarker && !strings.HasSuffix(rel, "/.gitkeep") {
			return fmt.Errorf("Event Search contract phase contains unexpected runtime file %s", rel)
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(files)
	if !containsString(files, eventSearchRuntimeMarker) {
		return "", fmt.Errorf("Event Search runtime marker is missing")
	}
	return "preimplementation", nil
}
