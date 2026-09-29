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
	caseQueueContractPath      = "engineering/implementation/case-queue/executable-contract.json"
	caseQueueFixturesPath      = "engineering/implementation/case-queue/fixtures.json"
	caseQueueManifestPath      = "work/lots/E10-INV-101A-CONTRACT/manifest.json"
	caseQueueCapabilityPath    = "cmdr-product-spec/07-investigate/modules/cases-and-evidence/capabilities/case-queue.md"
	caseLifecycleCapabilityPath = "cmdr-product-spec/07-investigate/modules/cases-and-evidence/capabilities/case-lifecycle-and-coordination.md"
	caseQueueScreenPath        = "cmdr-product-spec/07-investigate/modules/case-workspace/screens/case-workspace.md"
	caseQueueRuntimeBoundaryID = "case-queue-runtime"
	caseQueueRuntimeRoot       = "product-runtime/case-queue/**"
	caseQueueRuntimeMarker     = "product-runtime/case-queue/README.md"
)

type CaseQueueRuntimeBoundary struct {
	ID                          string   `json:"id"`
	Root                        string   `json:"root"`
	ExpectedState               string   `json:"expected_state"`
	ExternalRuntimeDependencies []string `json:"external_runtime_dependencies"`
}

type CaseQueueScopePolicy struct {
	TenantRequired              bool `json:"tenant_required"`
	EnvironmentRequired         bool `json:"environment_required"`
	WildcardTenantForbidden     bool `json:"wildcard_tenant_forbidden"`
	CrossTenantForbidden        bool `json:"cross_tenant_forbidden"`
	CaseAccessDecisionRequired  bool `json:"case_access_decision_required"`
}

type CaseQueueSearchPolicy struct {
	AllowedFields              []string `json:"allowed_fields"`
	CaseInsensitive            bool     `json:"case_insensitive"`
	TrimSpace                  bool     `json:"trim_space"`
	InaccessibleDataForbidden  bool     `json:"inaccessible_data_forbidden"`
}

type CaseQueueFilterPolicy struct {
	AllowedFields             []string `json:"allowed_fields"`
	UnknownFilterRejected     bool     `json:"unknown_filter_rejected"`
	EmptyValueRejected        bool     `json:"empty_value_rejected"`
	CaseInsensitiveExactMatch bool     `json:"case_insensitive_exact_match"`
}

type CaseQueueSortSpec struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type CaseQueueSortPolicy struct {
	AllowedKeys          []string            `json:"allowed_keys"`
	DefaultSort          []CaseQueueSortSpec `json:"default_sort"`
	StableTieBreakers    []string            `json:"stable_tie_breakers"`
	UnknownSortRejected  bool                `json:"unknown_sort_rejected"`
}

type CaseQueueSavedViewPolicy struct {
	SharedOwned                        bool `json:"shared_owned"`
	Optional                           bool `json:"optional"`
	StableIdentifierRequired           bool `json:"stable_identifier_required"`
	VersionRequired                    bool `json:"version_required"`
	TenantScoped                       bool `json:"tenant_scoped"`
	CallerAccessDecisionRequired       bool `json:"caller_access_decision_required"`
	FieldPermissionReevaluationRequired bool `json:"field_permission_reevaluation_required"`
	FilterPermissionReevaluationRequired bool `json:"filter_permission_reevaluation_required"`
	SortPermissionReevaluationRequired bool `json:"sort_permission_reevaluation_required"`
	DeniedElementsOmitted              bool `json:"denied_elements_omitted"`
	MutationForbidden                  bool `json:"mutation_forbidden"`
}

type CaseQueueFreshnessPolicy struct {
	NamedStates                    []string `json:"named_states"`
	MissingProjectionNamed         bool     `json:"missing_projection_named"`
	StaleProjectionNamed           bool     `json:"stale_projection_named"`
	UnavailableProjectionNamed     bool     `json:"unavailable_projection_named"`
	DeterministicSortedDiagnostics bool     `json:"deterministic_sorted_diagnostics"`
	SynthesisOfMissingFactsForbidden bool   `json:"synthesis_of_missing_facts_forbidden"`
}

type CaseQueueProjectionPolicy struct {
	Immutable                        bool `json:"immutable"`
	DeepCopyIsolation                bool `json:"deep_copy_isolation"`
	CaseIsOnlyOwnedWorkItem          bool `json:"case_is_only_owned_work_item"`
	IncidentContextOptional          bool `json:"incident_context_optional"`
	FindingContextOptional           bool `json:"finding_context_optional"`
	ActivityContextOptional          bool `json:"activity_context_optional"`
	TaskDecisionResponseRunAggregation bool `json:"task_decision_response_run_aggregation"`
	FinalColumnsClaimed              bool `json:"final_columns_claimed"`
	FinalCaseStateMachineClaimed     bool `json:"final_case_state_machine_claimed"`
}

type CaseQueueMutationPolicy struct {
	CaseCreateExecutable       bool `json:"case_create_executable"`
	CaseUpdateExecutable       bool `json:"case_update_executable"`
	AssignmentExecutable       bool `json:"assignment_executable"`
	StatusTransitionExecutable bool `json:"status_transition_executable"`
	LifecycleExecutable        bool `json:"lifecycle_executable"`
	SavedViewMutationExecutable bool `json:"saved_view_mutation_executable"`
	ExportJobCreationExecutable bool `json:"export_job_creation_executable"`
	TaskCreationExecutable     bool `json:"task_creation_executable"`
	DecisionCreationExecutable bool `json:"decision_creation_executable"`
	ResponseRunCreationExecutable bool `json:"response_run_creation_executable"`
}

type CaseQueueTechnologyPolicy struct {
	StorageEngineSelected         bool `json:"storage_engine_selected"`
	ProviderSelected              bool `json:"provider_selected"`
	RetentionPolicySelected       bool `json:"retention_policy_selected"`
	CollaborationBackendSelected  bool `json:"collaboration_backend_selected"`
	ExportBackendSelected         bool `json:"export_backend_selected"`
	FinalColumnsSelected          bool `json:"final_columns_selected"`
	FinalCaseStateMachineSelected bool `json:"final_case_state_machine_selected"`
	FinalUISelected               bool `json:"final_ui_selected"`
}

type CaseQueueExecutableContract struct {
	SchemaVersion        int                       `json:"schema_version"`
	ContractID           string                    `json:"contract_id"`
	ContractKind         string                    `json:"contract_kind"`
	CanonicalSchemaClaim bool                      `json:"canonical_schema_claim"`
	Capability           string                    `json:"capability"`
	RuntimeBoundary      CaseQueueRuntimeBoundary  `json:"runtime_boundary"`
	ProductRefs          ProductRefsV2             `json:"product_refs"`
	ScopePolicy          CaseQueueScopePolicy      `json:"scope_policy"`
	SearchPolicy         CaseQueueSearchPolicy     `json:"search_policy"`
	FilterPolicy         CaseQueueFilterPolicy     `json:"filter_policy"`
	SortPolicy           CaseQueueSortPolicy       `json:"sort_policy"`
	SavedViewPolicy      CaseQueueSavedViewPolicy  `json:"saved_view_policy"`
	FreshnessPolicy      CaseQueueFreshnessPolicy  `json:"freshness_policy"`
	ProjectionPolicy     CaseQueueProjectionPolicy `json:"projection_policy"`
	MutationPolicy       CaseQueueMutationPolicy   `json:"mutation_policy"`
	TechnologyPolicy     CaseQueueTechnologyPolicy `json:"technology_policy"`
	Exclusions           []string                  `json:"exclusions"`
	SemanticRules        []string                  `json:"semantic_rules"`
	SecurityInvariants   []string                  `json:"security_invariants"`
}

type CaseQueueFilter struct {
	Field  string   `json:"field"`
	Values []string `json:"values"`
}

type CaseQueueSavedViewField struct {
	Name           string `json:"name"`
	AccessDecision string `json:"access_decision"`
}

type CaseQueueSavedViewFilter struct {
	Field          string   `json:"field"`
	Values         []string `json:"values"`
	AccessDecision string   `json:"access_decision"`
}

type CaseQueueSavedViewSort struct {
	Field          string `json:"field"`
	Direction      string `json:"direction"`
	AccessDecision string `json:"access_decision"`
}

type CaseQueueSavedViewInput struct {
	ID             string                     `json:"id"`
	TenantRef      string                     `json:"tenant_ref"`
	Version        string                     `json:"version"`
	AccessDecision string                     `json:"access_decision"`
	Fields         []CaseQueueSavedViewField  `json:"fields"`
	Filters        []CaseQueueSavedViewFilter `json:"filters"`
	Sort           []CaseQueueSavedViewSort   `json:"sort"`
}

type CaseQueueIncidentContext struct {
	ID             string `json:"id"`
	TenantRef      string `json:"tenant_ref"`
	AccessDecision string `json:"access_decision"`
	State          string `json:"state"`
	Priority       string `json:"priority"`
	Impact         string `json:"impact"`
	Owner          string `json:"owner"`
	UpdatedAt      string `json:"updated_at"`
}

type CaseQueueFindingContext struct {
	ID             string `json:"id"`
	TenantRef      string `json:"tenant_ref"`
	CaseID         string `json:"case_id"`
	AccessDecision string `json:"access_decision"`
	State          string `json:"state"`
	Status         string `json:"status"`
	UpdatedAt      string `json:"updated_at"`
}

type CaseQueueCaseInput struct {
	ID             string                    `json:"id"`
	TenantRef      string                    `json:"tenant_ref"`
	HumanID        string                    `json:"human_id"`
	Title          string                    `json:"title"`
	Status         string                    `json:"status"`
	Owner          string                    `json:"owner"`
	NextAction     string                    `json:"next_action"`
	UpdatedAt      string                    `json:"updated_at"`
	AccessDecision string                    `json:"access_decision"`
	Freshness      string                    `json:"freshness"`
	Incident       *CaseQueueIncidentContext `json:"incident"`
	Findings       []CaseQueueFindingContext `json:"findings"`
	ActivityState  string                    `json:"activity_state"`
}

type CaseQueueFixtureInput struct {
	TenantRef                    string                   `json:"tenant_ref"`
	AuthorizedTenants            []string                 `json:"authorized_tenants"`
	EnvironmentRef               string                   `json:"environment_ref"`
	SearchQuery                  string                   `json:"search_query"`
	Filters                      []CaseQueueFilter        `json:"filters"`
	Sort                         []CaseQueueSortSpec      `json:"sort"`
	ViewDirty                    bool                     `json:"view_dirty"`
	SavedView                    *CaseQueueSavedViewInput `json:"saved_view"`
	Cases                        []CaseQueueCaseInput     `json:"cases"`
	MutationRequested            string                   `json:"mutation_requested"`
	ExportRequested              bool                     `json:"export_requested"`
	WorkQueueAggregationRequested bool                    `json:"work_queue_aggregation_requested"`
}

type CaseQueueFixtureExpected struct {
	Allowed                    bool     `json:"allowed"`
	ErrorCode                  string   `json:"error_code"`
	States                     []string `json:"states"`
	Diagnostics                []string `json:"diagnostics"`
	ProjectedCaseIDs           []string `json:"projected_case_ids"`
	ProjectedCount             int      `json:"projected_count"`
	PermissionFilteredCases    int      `json:"permission_filtered_cases"`
	PermissionFilteredElements int      `json:"permission_filtered_elements"`
	AppliedSavedViewFields     []string `json:"applied_saved_view_fields"`
	AppliedSort                []string `json:"applied_sort"`
	ImmutableProjection        bool     `json:"immutable_projection"`
	DeepCopyIsolated           bool     `json:"deep_copy_isolated"`
	MutationExecuted           bool     `json:"mutation_executed"`
	ExportJobCreated           bool     `json:"export_job_created"`
	WorkQueueAggregated        bool     `json:"work_queue_aggregated"`
	AuditRequired              bool     `json:"audit_required"`
}

type CaseQueueFixtureCase struct {
	ID       string                   `json:"id"`
	Input    CaseQueueFixtureInput    `json:"input"`
	Expected CaseQueueFixtureExpected `json:"expected"`
}

type CaseQueueFixtureSet struct {
	SchemaVersion int                    `json:"schema_version"`
	ContractID    string                 `json:"contract_id"`
	Cases         []CaseQueueFixtureCase `json:"cases"`
}

type CaseQueueContractAuditSummary struct {
	ContractID          string `json:"contract_id"`
	Capability          string `json:"capability"`
	RuntimeBoundary     string `json:"runtime_boundary"`
	RuntimeState        string `json:"runtime_state"`
	Fixtures            int    `json:"fixtures"`
	PositiveFixtures    int    `json:"positive_fixtures"`
	NegativeFixtures    int    `json:"negative_fixtures"`
	PartialFixtures     int    `json:"partial_fixtures"`
	StaleFixtures       int    `json:"stale_fixtures"`
	PermissionFilteredFixtures int `json:"permission_filtered_fixtures"`
	ViewDirtyFixtures   int    `json:"view_dirty_fixtures"`
	RuntimeDependencies int    `json:"runtime_dependencies"`
	Status              string `json:"status"`
}

type caseQueueEvaluatedRow struct {
	Case             CaseQueueCaseInput
	IncidentPriority string
	FindingStatuses  []string
	Updated          time.Time
	Partial          bool
	Stale            bool
	Diagnostics      []string
	PermissionFilteredElements int
}

func runCaseQueueContractAudit(root string) (CaseQueueContractAuditSummary, error) {
	var contract CaseQueueExecutableContract
	if err := decodeStrict(root, caseQueueContractPath, &contract); err != nil {
		return CaseQueueContractAuditSummary{}, err
	}
	var fixtures CaseQueueFixtureSet
	if err := decodeStrict(root, caseQueueFixturesPath, &fixtures); err != nil {
		return CaseQueueContractAuditSummary{}, err
	}
	manifest, err := decodeWorkManifestV2(root, caseQueueManifestPath)
	if err != nil {
		return CaseQueueContractAuditSummary{}, err
	}
	if err := validateCaseQueueContractCore(root, manifest, contract, fixtures, true); err != nil {
		return CaseQueueContractAuditSummary{}, err
	}

	var registry ArchitectureRegistry
	if err := decodeStrict(root, architectureRegistryPath, &registry); err != nil {
		return CaseQueueContractAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return CaseQueueContractAuditSummary{}, err
	}
	boundary, ok := architectureBoundaryByID(registry, caseQueueRuntimeBoundaryID)
	if !ok || boundary.Kind != "product-runtime" || !sameStringSet(boundary.Roots, []string{caseQueueRuntimeRoot}) {
		return CaseQueueContractAuditSummary{}, fmt.Errorf("Case Queue runtime boundary is missing or inconsistent")
	}
	runtimeState, err := validateCaseQueueRuntimeLayout(root, boundary, contract.RuntimeBoundary.ExpectedState)
	if err != nil {
		return CaseQueueContractAuditSummary{}, err
	}
	deps, unsupported, err := scanRuntimeBoundary(root, boundary.ID, caseQueueRuntimeRoot)
	if err != nil {
		return CaseQueueContractAuditSummary{}, err
	}
	if len(deps) != 0 || len(unsupported) != 0 {
		return CaseQueueContractAuditSummary{}, fmt.Errorf("Case Queue contract phase introduced runtime dependencies: deps=%v unsupported=%v", deps, unsupported)
	}
	negative, partial, stale, permissionFiltered, viewDirty := 0, 0, 0, 0, 0
	for _, fixture := range fixtures.Cases {
		if !fixture.Expected.Allowed {
			negative++
		}
		if containsString(fixture.Expected.States, "partial") {
			partial++
		}
		if containsString(fixture.Expected.States, "stale") {
			stale++
		}
		if containsString(fixture.Expected.States, "permission-filtered") {
			permissionFiltered++
		}
		if containsString(fixture.Expected.States, "view-dirty") {
			viewDirty++
		}
	}
	return CaseQueueContractAuditSummary{
		ContractID: contract.ContractID, Capability: contract.Capability,
		RuntimeBoundary: boundary.ID, RuntimeState: runtimeState,
		Fixtures: len(fixtures.Cases), PositiveFixtures: len(fixtures.Cases)-negative,
		NegativeFixtures: negative, PartialFixtures: partial, StaleFixtures: stale,
		PermissionFilteredFixtures: permissionFiltered, ViewDirtyFixtures: viewDirty,
		RuntimeDependencies: len(deps), Status: "PASS",
	}, nil
}

func validateCaseQueueContractCore(root string, manifest WorkManifestV2, contract CaseQueueExecutableContract, fixtures CaseQueueFixtureSet, validateSources bool) error {
	if contract.SchemaVersion != 1 || contract.ContractID != "CASE-QUEUE-READONLY-CONTRACT-V1" ||
		contract.ContractKind != "provider-neutral-readonly-local-executable-contract" || contract.CanonicalSchemaClaim {
		return fmt.Errorf("Case Queue contract identity/schema claim is invalid")
	}
	if contract.Capability != "CAP-INV-101" {
		return fmt.Errorf("Case Queue capability mismatch")
	}
	if contract.RuntimeBoundary.ID != caseQueueRuntimeBoundaryID ||
		contract.RuntimeBoundary.Root != caseQueueRuntimeRoot ||
		(contract.RuntimeBoundary.ExpectedState != "preimplementation" && contract.RuntimeBoundary.ExpectedState != "implemented") ||
		len(contract.RuntimeBoundary.ExternalRuntimeDependencies) != 0 {
		return fmt.Errorf("Case Queue runtime-boundary contract is invalid")
	}
	if validateSources {
		if err := validateProductRefs(root, contract.ProductRefs); err != nil {
			return err
		}
	}
	if !productRefsEqual(contract.ProductRefs, manifest.ProductRefs) {
		return fmt.Errorf("Case Queue Product Spec references do not match active work manifest")
	}
	if validateSources {
		if err := validateCaseQueueSourceAnchors(root, contract.RuntimeBoundary.ExpectedState); err != nil {
			return err
		}
	}
	scope := contract.ScopePolicy
	if !scope.TenantRequired || !scope.EnvironmentRequired || !scope.WildcardTenantForbidden ||
		!scope.CrossTenantForbidden || !scope.CaseAccessDecisionRequired {
		return fmt.Errorf("Case Queue scope policy is incomplete or unsafe")
	}
	if !sameStringSet(contract.SearchPolicy.AllowedFields, []string{"id","human_id","title","owner","next_action"}) ||
		!contract.SearchPolicy.CaseInsensitive || !contract.SearchPolicy.TrimSpace || !contract.SearchPolicy.InaccessibleDataForbidden {
		return fmt.Errorf("Case Queue search policy is incomplete or unsafe")
	}
	if !sameStringSet(contract.FilterPolicy.AllowedFields, []string{"status","owner","incident_priority","finding_status","freshness"}) ||
		!contract.FilterPolicy.UnknownFilterRejected || !contract.FilterPolicy.EmptyValueRejected ||
		!contract.FilterPolicy.CaseInsensitiveExactMatch {
		return fmt.Errorf("Case Queue filter policy is incomplete or unsafe")
	}
	if !sameStringSet(contract.SortPolicy.AllowedKeys, []string{"updated_at","human_id","title","owner","status"}) ||
		!sameStringSet(contract.SortPolicy.StableTieBreakers, []string{"human_id","id"}) ||
		!contract.SortPolicy.UnknownSortRejected ||
		!caseQueueSortEqual(contract.SortPolicy.DefaultSort, []CaseQueueSortSpec{{Field:"updated_at",Direction:"desc"},{Field:"human_id",Direction:"asc"},{Field:"id",Direction:"asc"}}) {
		return fmt.Errorf("Case Queue sort policy is incomplete or nondeterministic")
	}
	view := contract.SavedViewPolicy
	if !view.SharedOwned || !view.Optional || !view.StableIdentifierRequired || !view.VersionRequired ||
		!view.TenantScoped || !view.CallerAccessDecisionRequired || !view.FieldPermissionReevaluationRequired ||
		!view.FilterPermissionReevaluationRequired || !view.SortPermissionReevaluationRequired ||
		!view.DeniedElementsOmitted || !view.MutationForbidden {
		return fmt.Errorf("Case Queue Saved View policy is incomplete or unsafe")
	}
	fresh := contract.FreshnessPolicy
	if !sameStringSet(fresh.NamedStates, []string{"available","empty","partial","stale","permission-filtered","view-dirty"}) ||
		!fresh.MissingProjectionNamed || !fresh.StaleProjectionNamed || !fresh.UnavailableProjectionNamed ||
		!fresh.DeterministicSortedDiagnostics || !fresh.SynthesisOfMissingFactsForbidden {
		return fmt.Errorf("Case Queue freshness policy is incomplete")
	}
	proj := contract.ProjectionPolicy
	if !proj.Immutable || !proj.DeepCopyIsolation || !proj.CaseIsOnlyOwnedWorkItem ||
		!proj.IncidentContextOptional || !proj.FindingContextOptional || !proj.ActivityContextOptional ||
		proj.TaskDecisionResponseRunAggregation || proj.FinalColumnsClaimed || proj.FinalCaseStateMachineClaimed {
		return fmt.Errorf("Case Queue projection policy violates bounded read-only scope")
	}
	mut := contract.MutationPolicy
	if mut.CaseCreateExecutable || mut.CaseUpdateExecutable || mut.AssignmentExecutable ||
		mut.StatusTransitionExecutable || mut.LifecycleExecutable || mut.SavedViewMutationExecutable ||
		mut.ExportJobCreationExecutable || mut.TaskCreationExecutable || mut.DecisionCreationExecutable ||
		mut.ResponseRunCreationExecutable {
		return fmt.Errorf("Case Queue mutation policy violates bounded read-only scope")
	}
	tech := contract.TechnologyPolicy
	if tech.StorageEngineSelected || tech.ProviderSelected || tech.RetentionPolicySelected ||
		tech.CollaborationBackendSelected || tech.ExportBackendSelected || tech.FinalColumnsSelected ||
		tech.FinalCaseStateMachineSelected || tech.FinalUISelected {
		return fmt.Errorf("Case Queue contract selects unresolved technology or final UX")
	}
	for _, exclusion := range []string{
		"case-create-update-assignment-lifecycle:CAP-INV-102:OPEN-013",
		"saved-view-create-update-share-archive-manage",
		"export-job-creation-or-export-backend",
		"task-decision-response-run-ownership-or-aggregation",
		"command-work-queue-behavior",
		"canonical-case-schema-change",
		"final-case-state-machine",
		"final-queue-columns",
		"storage-provider-retention",
		"collaboration-backend",
		"final-cap-inv-101-ui",
	} {
		if !containsString(contract.Exclusions, exclusion) {
			return fmt.Errorf("Case Queue exclusion %s is missing", exclusion)
		}
	}
	if len(contract.SemanticRules) < 9 || len(contract.SecurityInvariants) < 8 {
		return fmt.Errorf("Case Queue semantic/security contract is incomplete")
	}
	if fixtures.SchemaVersion != 1 || fixtures.ContractID != contract.ContractID {
		return fmt.Errorf("Case Queue fixtures do not bind the executable contract")
	}
	return validateCaseQueueFixtures(contract, fixtures)
}

func validateCaseQueueSourceAnchors(root, expectedState string) error {
	runtimeStateAnchor := ""
	switch expectedState {
	case "preimplementation":
		runtimeStateAnchor = "Preimplementation only"
	case "implemented":
		runtimeStateAnchor = "Implemented scope"
	default:
		return fmt.Errorf("unsupported Case Queue runtime state %q", expectedState)
	}
	sources := []struct {
		path string
		anchors []string
	}{
		{caseQueueCapabilityPath, []string{
			"id: CAP-INV-101", "Case Queue ≠ Work Queue", "Saved View avec champ interdit",
			"permission-filtered", "OPEN-013",
		}},
		{caseLifecycleCapabilityPath, []string{
			"id: CAP-INV-102", "Machine finale reportée", "OPEN-013",
		}},
		{caseQueueScreenPath, []string{
			"id: INV-CAS-001", "perm.investigate.case.read", "État Partial", "État Permission denied",
		}},
		{caseQueueRuntimeMarker, []string{
			"Case Queue runtime boundary", runtimeStateAnchor, "OPEN-013", "Task, Decision or Response Run ownership",
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
				return fmt.Errorf("Case Queue source anchor %q missing from %s", anchor, source.path)
			}
		}
	}
	return nil
}

func validateCaseQueueFixtures(contract CaseQueueExecutableContract, fixtures CaseQueueFixtureSet) error {
	required := map[string]bool{
		"available-stable-default-order": false,
		"empty-queue-is-explicit": false,
		"partial-incident-context-is-visible": false,
		"stale-context-is-visible": false,
		"permission-filtered-case-is-hidden": false,
		"saved-view-denied-elements-are-omitted": false,
		"view-dirty-state-is-visible": false,
		"search-filter-sort-is-deterministic": false,
		"missing-tenant-is-rejected": false,
		"wildcard-tenant-is-rejected": false,
		"missing-environment-is-rejected": false,
		"cross-tenant-case-is-rejected": false,
		"missing-case-id-is-rejected": false,
		"invalid-case-updated-at-is-rejected": false,
		"unknown-filter-is-rejected": false,
		"invalid-sort-is-rejected": false,
		"saved-view-access-denied-is-rejected": false,
		"cross-tenant-saved-view-is-rejected": false,
		"missing-saved-view-version-is-rejected": false,
		"cross-tenant-incident-is-rejected": false,
		"cross-tenant-finding-is-rejected": false,
		"case-mutation-is-rejected": false,
		"export-request-is-rejected": false,
		"work-queue-aggregation-is-rejected": false,
	}
	if len(fixtures.Cases) < len(required) {
		return fmt.Errorf("Case Queue fixture set is too small: %d", len(fixtures.Cases))
	}
	ids := map[string]bool{}
	for _, fixture := range fixtures.Cases {
		if strings.TrimSpace(fixture.ID) == "" || ids[fixture.ID] {
			return fmt.Errorf("Case Queue fixture id is empty or duplicate: %q", fixture.ID)
		}
		ids[fixture.ID] = true
		if _, ok := required[fixture.ID]; ok {
			required[fixture.ID] = true
		}
		want := evaluateCaseQueueFixture(contract, fixture.Input)
		if !caseQueueExpectedEqual(fixture.Expected, want) {
			return fmt.Errorf("Case Queue fixture %s expectation mismatch: want %#v got %#v", fixture.ID, want, fixture.Expected)
		}
	}
	for id, present := range required {
		if !present {
			return fmt.Errorf("required Case Queue fixture %s is missing", id)
		}
	}
	return nil
}

func evaluateCaseQueueFixture(contract CaseQueueExecutableContract, in CaseQueueFixtureInput) CaseQueueFixtureExpected {
	reject := func(code string) CaseQueueFixtureExpected {
		return CaseQueueFixtureExpected{
			ErrorCode: code, AuditRequired: true,
			States: []string{}, Diagnostics: []string{}, ProjectedCaseIDs: []string{},
			AppliedSavedViewFields: []string{}, AppliedSort: []string{},
		}
	}
	if strings.TrimSpace(in.MutationRequested) != "" {
		return reject("mutation-open-decision")
	}
	if in.ExportRequested {
		return reject("export-forbidden")
	}
	if in.WorkQueueAggregationRequested {
		return reject("work-queue-aggregation-forbidden")
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

	requestFilters, errCode := validateCaseQueueFilters(contract, in.Filters)
	if errCode != "" {
		return reject(errCode)
	}
	requestSort, errCode := validateCaseQueueSort(contract, in.Sort)
	if errCode != "" {
		return reject(errCode)
	}

	diagnostics := map[string]bool{}
	permissionFilteredElements := 0
	appliedFields := []string{}
	viewFilters := []CaseQueueFilter{}
	viewSort := []CaseQueueSortSpec{}
	if in.SavedView != nil {
		view := in.SavedView
		if strings.TrimSpace(view.ID) == "" {
			return reject("missing-saved-view-id")
		}
		if view.TenantRef != in.TenantRef {
			return reject("cross-tenant-saved-view")
		}
		if strings.TrimSpace(view.Version) == "" {
			return reject("missing-saved-view-version")
		}
		switch view.AccessDecision {
		case "allow":
		case "deny":
			return reject("saved-view-access-denied")
		default:
			return reject("invalid-saved-view-access-decision")
		}
		fieldSeen := map[string]bool{}
		for _, field := range view.Fields {
			if field.AccessDecision == "deny" {
				permissionFilteredElements++
				diagnostics["saved-view:permission-filtered"] = true
				continue
			}
			if field.AccessDecision != "allow" {
				return reject("invalid-saved-view-element-access-decision")
			}
			name := strings.TrimSpace(field.Name)
			if name == "" {
				return reject("invalid-saved-view-field")
			}
			if fieldSeen[name] {
				return reject("duplicate-saved-view-field")
			}
			fieldSeen[name] = true
			appliedFields = append(appliedFields, name)
		}
		for _, filter := range view.Filters {
			if filter.AccessDecision == "deny" {
				permissionFilteredElements++
				diagnostics["saved-view:permission-filtered"] = true
				continue
			}
			if filter.AccessDecision != "allow" {
				return reject("invalid-saved-view-element-access-decision")
			}
			validated, code := validateCaseQueueFilters(contract, []CaseQueueFilter{{Field:filter.Field, Values:filter.Values}})
			if code != "" {
				return reject(code)
			}
			viewFilters = append(viewFilters, validated...)
		}
		for _, spec := range view.Sort {
			if spec.AccessDecision == "deny" {
				permissionFilteredElements++
				diagnostics["saved-view:permission-filtered"] = true
				continue
			}
			if spec.AccessDecision != "allow" {
				return reject("invalid-saved-view-element-access-decision")
			}
			validated, code := validateCaseQueueSort(contract, []CaseQueueSortSpec{{Field:spec.Field, Direction:spec.Direction}})
			if code != "" {
				return reject(code)
			}
			viewSort = append(viewSort, validated...)
		}
	}
	sort.Strings(appliedFields)

	effectiveFilters := append([]CaseQueueFilter{}, viewFilters...)
	effectiveFilters = append(effectiveFilters, requestFilters...)
	effectiveSort := requestSort
	if len(effectiveSort) == 0 {
		effectiveSort = viewSort
	}
	if len(effectiveSort) == 0 {
		effectiveSort = append([]CaseQueueSortSpec{}, contract.SortPolicy.DefaultSort...)
	}
	effectiveSort = appendCaseQueueTieBreakers(effectiveSort, contract.SortPolicy.StableTieBreakers)

	search := strings.ToLower(strings.TrimSpace(in.SearchQuery))
	permissionFilteredCases := 0
	caseIDs := map[string]bool{}
	rows := []caseQueueEvaluatedRow{}
	for _, item := range in.Cases {
		if item.TenantRef != in.TenantRef {
			return reject("cross-tenant-case")
		}
		if strings.TrimSpace(item.ID) == "" {
			return reject("missing-case-id")
		}
		if caseIDs[item.ID] {
			return reject("duplicate-case-id")
		}
		caseIDs[item.ID] = true
		switch item.AccessDecision {
		case "deny":
			permissionFilteredCases++
			diagnostics["case:permission-filtered"] = true
			continue
		case "allow":
		default:
			return reject("invalid-case-access-decision")
		}
		if strings.TrimSpace(item.HumanID) == "" || strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.Status) == "" {
			return reject("incomplete-case-snapshot")
		}
		updated, err := time.Parse(time.RFC3339, item.UpdatedAt)
		if err != nil {
			return reject("invalid-case-updated-at")
		}
		row := caseQueueEvaluatedRow{Case:item, Updated:updated, FindingStatuses:[]string{}, Diagnostics:[]string{}}
		switch item.Freshness {
		case "available":
		case "partial":
			row.Partial = true
			row.Diagnostics = append(row.Diagnostics, "case:"+item.ID+":partial")
		case "stale":
			row.Stale = true
			row.Diagnostics = append(row.Diagnostics, "case:"+item.ID+":stale")
		default:
			return reject("invalid-case-freshness")
		}
		if item.Incident != nil {
			incident := item.Incident
			if incident.TenantRef != in.TenantRef {
				return reject("cross-tenant-incident")
			}
			switch incident.AccessDecision {
			case "deny":
				row.PermissionFilteredElements++
				row.Diagnostics = append(row.Diagnostics, "incident-context:permission-filtered")
			case "allow":
				switch incident.State {
				case "available":
					row.IncidentPriority = incident.Priority
				case "partial":
					row.Partial = true
					row.IncidentPriority = incident.Priority
					row.Diagnostics = append(row.Diagnostics, "incident-context:"+item.ID+":partial")
				case "stale":
					row.Stale = true
					row.IncidentPriority = incident.Priority
					row.Diagnostics = append(row.Diagnostics, "incident-context:"+item.ID+":stale")
				case "unavailable":
					row.Partial = true
					row.Diagnostics = append(row.Diagnostics, "incident-context:"+item.ID+":unavailable")
				default:
					return reject("invalid-incident-state")
				}
			default:
				return reject("invalid-incident-access-decision")
			}
		}
		for _, finding := range item.Findings {
			if finding.TenantRef != in.TenantRef {
				return reject("cross-tenant-finding")
			}
			if finding.CaseID != item.ID {
				return reject("finding-case-mismatch")
			}
			switch finding.AccessDecision {
			case "deny":
				row.PermissionFilteredElements++
				row.Diagnostics = append(row.Diagnostics, "finding-context:permission-filtered")
				continue
			case "allow":
			default:
				return reject("invalid-finding-access-decision")
			}
			switch finding.State {
			case "available":
				row.FindingStatuses = append(row.FindingStatuses, finding.Status)
			case "partial":
				row.Partial = true
				row.FindingStatuses = append(row.FindingStatuses, finding.Status)
				row.Diagnostics = append(row.Diagnostics, "finding-context:"+item.ID+":partial")
			case "stale":
				row.Stale = true
				row.FindingStatuses = append(row.FindingStatuses, finding.Status)
				row.Diagnostics = append(row.Diagnostics, "finding-context:"+item.ID+":stale")
			case "unavailable":
				row.Partial = true
				row.Diagnostics = append(row.Diagnostics, "finding-context:"+item.ID+":unavailable")
			default:
				return reject("invalid-finding-state")
			}
		}
		switch item.ActivityState {
		case "", "available":
		case "partial", "unavailable":
			row.Partial = true
			row.Diagnostics = append(row.Diagnostics, "activity:"+item.ID+":"+item.ActivityState)
		case "stale":
			row.Stale = true
			row.Diagnostics = append(row.Diagnostics, "activity:"+item.ID+":stale")
		default:
			return reject("invalid-activity-state")
		}
		if !caseQueueMatchesSearch(row, search) || !caseQueueMatchesFilters(row, effectiveFilters) {
			continue
		}
		rows = append(rows, row)
	}
	sort.SliceStable(rows, func(i, j int) bool {
		return caseQueueRowLess(rows[i], rows[j], effectiveSort)
	})

	states := map[string]bool{}
	if len(rows) == 0 {
		states["empty"] = true
	} else {
		states["available"] = true
	}
	if permissionFilteredCases > 0 || permissionFilteredElements > 0 {
		states["permission-filtered"] = true
	}
	if in.ViewDirty {
		states["view-dirty"] = true
	}
	for _, row := range rows {
		if row.Partial {
			states["partial"] = true
		}
		if row.Stale {
			states["stale"] = true
		}
		permissionFilteredElements += row.PermissionFilteredElements
		for _, diagnostic := range row.Diagnostics {
			diagnostics[diagnostic] = true
		}
	}
	if permissionFilteredElements > 0 {
		states["permission-filtered"] = true
	}
	projected := make([]string, 0, len(rows))
	for _, row := range rows {
		projected = append(projected, row.Case.ID)
	}
	return CaseQueueFixtureExpected{
		Allowed:true, States:sortedBoolKeys(states), Diagnostics:sortedBoolKeys(diagnostics),
		ProjectedCaseIDs:projected, ProjectedCount:len(projected),
		PermissionFilteredCases:permissionFilteredCases, PermissionFilteredElements:permissionFilteredElements,
		AppliedSavedViewFields:appliedFields, AppliedSort:caseQueueSortStrings(effectiveSort),
		ImmutableProjection:true, DeepCopyIsolated:true, AuditRequired:true,
	}
}

func validateCaseQueueFilters(contract CaseQueueExecutableContract, filters []CaseQueueFilter) ([]CaseQueueFilter, string) {
	out := make([]CaseQueueFilter, 0, len(filters))
	for _, filter := range filters {
		if !containsString(contract.FilterPolicy.AllowedFields, filter.Field) {
			return nil, "invalid-filter-field"
		}
		if len(filter.Values) == 0 {
			return nil, "empty-filter-value"
		}
		copyFilter := CaseQueueFilter{Field:filter.Field, Values:make([]string,0,len(filter.Values))}
		for _, value := range filter.Values {
			value = strings.TrimSpace(value)
			if value == "" {
				return nil, "empty-filter-value"
			}
			copyFilter.Values = append(copyFilter.Values, strings.ToLower(value))
		}
		out = append(out, copyFilter)
	}
	return out, ""
}

func validateCaseQueueSort(contract CaseQueueExecutableContract, specs []CaseQueueSortSpec) ([]CaseQueueSortSpec, string) {
	out := make([]CaseQueueSortSpec, 0, len(specs))
	seen := map[string]bool{}
	for _, spec := range specs {
		if !containsString(contract.SortPolicy.AllowedKeys, spec.Field) {
			return nil, "invalid-sort-field"
		}
		if spec.Direction != "asc" && spec.Direction != "desc" {
			return nil, "invalid-sort-direction"
		}
		if seen[spec.Field] {
			return nil, "duplicate-sort-field"
		}
		seen[spec.Field] = true
		out = append(out, spec)
	}
	return out, ""
}

func appendCaseQueueTieBreakers(specs []CaseQueueSortSpec, tie []string) []CaseQueueSortSpec {
	out := append([]CaseQueueSortSpec{}, specs...)
	seen := map[string]bool{}
	for _, spec := range out {
		seen[spec.Field] = true
	}
	for _, field := range tie {
		if !seen[field] {
			out = append(out, CaseQueueSortSpec{Field:field, Direction:"asc"})
			seen[field] = true
		}
	}
	return out
}

func caseQueueMatchesSearch(row caseQueueEvaluatedRow, query string) bool {
	if query == "" {
		return true
	}
	values := []string{row.Case.ID,row.Case.HumanID,row.Case.Title,row.Case.Owner,row.Case.NextAction}
	for _, value := range values {
		if strings.Contains(strings.ToLower(value), query) {
			return true
		}
	}
	return false
}

func caseQueueMatchesFilters(row caseQueueEvaluatedRow, filters []CaseQueueFilter) bool {
	for _, filter := range filters {
		matched := false
		switch filter.Field {
		case "status":
			matched = containsNormalized(filter.Values, row.Case.Status)
		case "owner":
			matched = containsNormalized(filter.Values, row.Case.Owner)
		case "incident_priority":
			matched = containsNormalized(filter.Values, row.IncidentPriority)
		case "finding_status":
			for _, status := range row.FindingStatuses {
				if containsNormalized(filter.Values, status) {
					matched = true
					break
				}
			}
		case "freshness":
			matched = containsNormalized(filter.Values, row.Case.Freshness)
		}
		if !matched {
			return false
		}
	}
	return true
}

func containsNormalized(values []string, value string) bool {
	value = strings.ToLower(strings.TrimSpace(value))
	for _, candidate := range values {
		if strings.ToLower(strings.TrimSpace(candidate)) == value {
			return true
		}
	}
	return false
}

func caseQueueRowLess(a,b caseQueueEvaluatedRow, specs []CaseQueueSortSpec) bool {
	for _, spec := range specs {
		cmp := caseQueueCompareField(a,b,spec.Field)
		if cmp == 0 {
			continue
		}
		if spec.Direction == "desc" {
			return cmp > 0
		}
		return cmp < 0
	}
	return false
}

func caseQueueCompareField(a,b caseQueueEvaluatedRow, field string) int {
	switch field {
	case "updated_at":
		if a.Updated.Before(b.Updated) { return -1 }
		if a.Updated.After(b.Updated) { return 1 }
		return 0
	case "id":
		return strings.Compare(strings.ToLower(a.Case.ID), strings.ToLower(b.Case.ID))
	case "human_id":
		return strings.Compare(strings.ToLower(a.Case.HumanID), strings.ToLower(b.Case.HumanID))
	case "title":
		return strings.Compare(strings.ToLower(a.Case.Title), strings.ToLower(b.Case.Title))
	case "owner":
		return strings.Compare(strings.ToLower(a.Case.Owner), strings.ToLower(b.Case.Owner))
	case "status":
		return strings.Compare(strings.ToLower(a.Case.Status), strings.ToLower(b.Case.Status))
	default:
		return 0
	}
}

func caseQueueSortStrings(specs []CaseQueueSortSpec) []string {
	out := make([]string,0,len(specs))
	for _, spec := range specs {
		out = append(out, spec.Field+":"+spec.Direction)
	}
	return out
}

func caseQueueSortEqual(a,b []CaseQueueSortSpec) bool {
	if len(a)!=len(b) { return false }
	for i:=range a {
		if a[i]!=b[i] { return false }
	}
	return true
}

func caseQueueExpectedEqual(a,b CaseQueueFixtureExpected) bool {
	return a.Allowed==b.Allowed && a.ErrorCode==b.ErrorCode &&
		sameStringSet(a.States,b.States) && sameStringSet(a.Diagnostics,b.Diagnostics) &&
		stringSlicesEqual(a.ProjectedCaseIDs,b.ProjectedCaseIDs) &&
		a.ProjectedCount==b.ProjectedCount &&
		a.PermissionFilteredCases==b.PermissionFilteredCases &&
		a.PermissionFilteredElements==b.PermissionFilteredElements &&
		stringSlicesEqual(a.AppliedSavedViewFields,b.AppliedSavedViewFields) &&
		stringSlicesEqual(a.AppliedSort,b.AppliedSort) &&
		a.ImmutableProjection==b.ImmutableProjection && a.DeepCopyIsolated==b.DeepCopyIsolated &&
		a.MutationExecuted==b.MutationExecuted && a.ExportJobCreated==b.ExportJobCreated &&
		a.WorkQueueAggregated==b.WorkQueueAggregated && a.AuditRequired==b.AuditRequired
}

func stringSlicesEqual(a,b []string) bool {
	if len(a)!=len(b) { return false }
	for i:=range a {
		if a[i]!=b[i] { return false }
	}
	return true
}

func sortedBoolKeys(values map[string]bool) []string {
	out:=make([]string,0,len(values))
	for value, present:=range values {
		if present { out=append(out,value) }
	}
	sort.Strings(out)
	return out
}

func validateCaseQueueRuntimeLayout(root string, boundary ArchitectureBoundary, expected string) (string,error) {
	scanRoot, err := runtimeScanRoot(root, caseQueueRuntimeRoot)
	if err != nil { return "",err }
	files:=[]string{}
	hasGoMod:=false
	hasRuntimeGo:=false
	err = filepath.WalkDir(scanRoot, func(path string, entry os.DirEntry, err error) error {
		if err!=nil { return err }
		if entry.Type()&os.ModeSymlink!=0 {
			return fmt.Errorf("Case Queue runtime boundary contains symlink: %s",path)
		}
		if entry.IsDir() { return nil }
		rel,err:=filepath.Rel(root,path)
		if err!=nil { return err }
		rel=filepath.ToSlash(rel)
		files=append(files,rel)
		if rel=="product-runtime/case-queue/go.mod" { hasGoMod=true }
		if strings.HasSuffix(rel,".go") && !strings.HasSuffix(rel,"_test.go") { hasRuntimeGo=true }
		if expected=="preimplementation" && rel!=caseQueueRuntimeMarker && !strings.HasSuffix(rel,"/.gitkeep") {
			return fmt.Errorf("Case Queue preimplementation boundary contains unexpected runtime file %s",rel)
		}
		return nil
	})
	if err!=nil { return "",err }
	sort.Strings(files)
	if !containsString(files,caseQueueRuntimeMarker) {
		return "",fmt.Errorf("Case Queue runtime marker is missing")
	}
	actual,err:=detectRuntimeBoundaryImplementationState(root,boundary)
	if err!=nil { return "",err }
	if actual!=expected {
		return "",fmt.Errorf("Case Queue runtime state mismatch: contract=%s actual=%s",expected,actual)
	}
	if expected=="implemented" && (!hasGoMod || !hasRuntimeGo) {
		return "",fmt.Errorf("Case Queue implemented runtime requires go.mod and executable Go source")
	}
	return actual,nil
}
