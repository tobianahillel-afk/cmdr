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
	huntManagementContractPath      = "engineering/implementation/hunt-management/executable-contract.json"
	huntManagementFixturesPath      = "engineering/implementation/hunt-management/fixtures.json"
	huntManagementManifestPath      = "work/lots/E10-INV-005A-CONTRACT/manifest.json"
	huntManagementCapabilityPath    = "cmdr-product-spec/07-investigate/modules/signals-and-hunt/capabilities/hunt-management.md"
	huntManagementRuntimeBoundaryID = "hunt-management-runtime"
	huntManagementRuntimeRoot       = "product-runtime/hunt-management/**"
	huntManagementRuntimeMarker     = "product-runtime/hunt-management/README.md"
)

type HuntRuntimeBoundary struct {
	ID                          string   `json:"id"`
	Root                        string   `json:"root"`
	ExpectedState               string   `json:"expected_state"`
	ExternalRuntimeDependencies []string `json:"external_runtime_dependencies"`
}

type HuntScopePolicy struct {
	TenantRequired          bool `json:"tenant_required"`
	EnvironmentRequired     bool `json:"environment_required"`
	QuestionRequired        bool `json:"question_required"`
	ScopeRequired           bool `json:"scope_required"`
	TimeRangeRequired       bool `json:"time_range_required"`
	OwnerRequired           bool `json:"owner_required"`
	WildcardTenantForbidden bool `json:"wildcard_tenant_forbidden"`
	CrossTenantForbidden    bool `json:"cross_tenant_forbidden"`
}

type HuntReferencePolicy struct {
	AllowedTypes                       []string `json:"allowed_types"`
	TenantScoped                       bool     `json:"tenant_scoped"`
	StableIdentifierRequired           bool     `json:"stable_identifier_required"`
	CallerAccessDecisionRequired       bool     `json:"caller_access_decision_required"`
	DeniedReferenceProjectionForbidden bool     `json:"denied_reference_projection_forbidden"`
	OwnershipTransferForbidden         bool     `json:"ownership_transfer_forbidden"`
	ProvenanceRequired                 bool     `json:"provenance_required"`
	SearchJobRequiresVisibleQuery      bool     `json:"search_job_requires_visible_query"`
	UnavailableStateBecomesLimitation  bool     `json:"unavailable_state_becomes_limitation"`
	StaleStateBecomesLimitation        bool     `json:"stale_state_becomes_limitation"`
	PartialStateBecomesLimitation      bool     `json:"partial_state_becomes_limitation"`
}

type HuntProjectionPolicy struct {
	Immutable                        bool `json:"immutable"`
	DeepCopyIsolation                bool `json:"deep_copy_isolation"`
	QuestionVisible                  bool `json:"question_visible"`
	ScopeVisible                     bool `json:"scope_visible"`
	TimeRangeVisible                 bool `json:"time_range_visible"`
	OwnerVisible                     bool `json:"owner_visible"`
	ContributorsVisible              bool `json:"contributors_visible"`
	ProvenanceVisible                bool `json:"provenance_visible"`
	LimitationsVisible               bool `json:"limitations_visible"`
	SynthesisOfMissingFactsForbidden bool `json:"synthesis_of_missing_facts_forbidden"`
}

type HuntMutationPolicy struct {
	CanonicalHuntObjectCreated      bool `json:"canonical_hunt_object_created"`
	HuntStateMutationExecutable     bool `json:"hunt_state_mutation_executable"`
	ContributorMutationExecutable   bool `json:"contributor_mutation_executable"`
	CasePromotionMutationExecutable bool `json:"case_promotion_mutation_executable"`
	CaseLinkMutationExecutable      bool `json:"case_link_mutation_executable"`
	HypothesisMutationExecutable    bool `json:"hypothesis_mutation_executable"`
	QueryMutationExecutable         bool `json:"query_mutation_executable"`
	QueryExecutionExecutable        bool `json:"query_execution_executable"`
	SavedSearchMutationExecutable   bool `json:"saved_search_mutation_executable"`
	SearchJobExecutionExecutable    bool `json:"search_job_execution_executable"`
}

type HuntTechnologyPolicy struct {
	StorageEngineSelected       bool `json:"storage_engine_selected"`
	ProviderSelected            bool `json:"provider_selected"`
	RetentionPolicySelected     bool `json:"retention_policy_selected"`
	CollaborationBackendSelected bool `json:"collaboration_backend_selected"`
	FinalQueryDialectSelected   bool `json:"final_query_dialect_selected"`
	FinalHuntUISelected         bool `json:"final_hunt_ui_selected"`
}

type HuntExecutableContract struct {
	SchemaVersion        int                 `json:"schema_version"`
	ContractID           string              `json:"contract_id"`
	ContractKind         string              `json:"contract_kind"`
	CanonicalSchemaClaim bool                `json:"canonical_schema_claim"`
	Capability           string              `json:"capability"`
	RuntimeBoundary      HuntRuntimeBoundary `json:"runtime_boundary"`
	ProductRefs          ProductRefsV2       `json:"product_refs"`
	ScopePolicy          HuntScopePolicy     `json:"scope_policy"`
	ReferencePolicy      HuntReferencePolicy `json:"reference_policy"`
	ProjectionPolicy     HuntProjectionPolicy `json:"projection_policy"`
	MutationPolicy       HuntMutationPolicy  `json:"mutation_policy"`
	TechnologyPolicy     HuntTechnologyPolicy `json:"technology_policy"`
	Exclusions           []string            `json:"exclusions"`
	SemanticRules        []string            `json:"semantic_rules"`
	SecurityInvariants   []string            `json:"security_invariants"`
}

type HuntFixtureReference struct {
	Kind           string `json:"kind"`
	ID             string `json:"id"`
	TenantRef      string `json:"tenant_ref"`
	AccessDecision string `json:"access_decision"`
	QueryRef       string `json:"query_ref"`
	State          string `json:"state"`
}

type HuntFixtureInput struct {
	TenantRef         string                 `json:"tenant_ref"`
	AuthorizedTenants []string               `json:"authorized_tenants"`
	EnvironmentRef    string                 `json:"environment_ref"`
	Question          string                 `json:"question"`
	Scope             string                 `json:"scope"`
	TimeStart         string                 `json:"time_start"`
	TimeEnd           string                 `json:"time_end"`
	OwnerRef          string                 `json:"owner_ref"`
	ContributorRefs   []string               `json:"contributor_refs"`
	References        []HuntFixtureReference `json:"references"`
	MutationRequested string                 `json:"mutation_requested"`
}

type HuntFixtureExpected struct {
	Allowed              bool     `json:"allowed"`
	ErrorCode            string   `json:"error_code"`
	ProjectedReferences  int      `json:"projected_references"`
	Limitations          []string `json:"limitations"`
	ImmutableProjection  bool     `json:"immutable_projection"`
	DeepCopyIsolated     bool     `json:"deep_copy_isolated"`
	MutationExecuted     bool     `json:"mutation_executed"`
	AuditRequired        bool     `json:"audit_required"`
}

type HuntFixtureCase struct {
	ID       string              `json:"id"`
	Input    HuntFixtureInput    `json:"input"`
	Expected HuntFixtureExpected `json:"expected"`
}

type HuntFixtureSet struct {
	SchemaVersion int               `json:"schema_version"`
	ContractID    string            `json:"contract_id"`
	Cases         []HuntFixtureCase `json:"cases"`
}

type HuntContractAuditSummary struct {
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

func runHuntManagementContractAudit(root string) (HuntContractAuditSummary, error) {
	var contract HuntExecutableContract
	if err := decodeStrict(root, huntManagementContractPath, &contract); err != nil {
		return HuntContractAuditSummary{}, err
	}
	var fixtures HuntFixtureSet
	if err := decodeStrict(root, huntManagementFixturesPath, &fixtures); err != nil {
		return HuntContractAuditSummary{}, err
	}
	manifest, err := decodeWorkManifestV2(root, huntManagementManifestPath)
	if err != nil {
		return HuntContractAuditSummary{}, err
	}
	if err := validateHuntManagementContractCore(root, manifest, contract, fixtures, true); err != nil {
		return HuntContractAuditSummary{}, err
	}

	var registry ArchitectureRegistry
	if err := decodeStrict(root, architectureRegistryPath, &registry); err != nil {
		return HuntContractAuditSummary{}, err
	}
	if err := validateArchitectureRegistry(registry); err != nil {
		return HuntContractAuditSummary{}, err
	}
	boundary, ok := architectureBoundaryByID(registry, huntManagementRuntimeBoundaryID)
	if !ok || boundary.Kind != "product-runtime" || !sameStringSet(boundary.Roots, []string{huntManagementRuntimeRoot}) {
		return HuntContractAuditSummary{}, fmt.Errorf("Hunt Management runtime boundary is missing or inconsistent")
	}
	runtimeState, err := validateHuntManagementRuntimeLayout(root, boundary, contract.RuntimeBoundary.ExpectedState)
	if err != nil {
		return HuntContractAuditSummary{}, err
	}
	deps, unsupported, err := scanRuntimeBoundary(root, boundary.ID, huntManagementRuntimeRoot)
	if err != nil {
		return HuntContractAuditSummary{}, err
	}
	if len(deps) != 0 || len(unsupported) != 0 {
		return HuntContractAuditSummary{}, fmt.Errorf("Hunt contract phase introduced runtime dependencies: deps=%v unsupported=%v", deps, unsupported)
	}
	negative := 0
	for _, fixture := range fixtures.Cases {
		if !fixture.Expected.Allowed {
			negative++
		}
	}
	return HuntContractAuditSummary{
		ContractID: contract.ContractID, Capability: contract.Capability,
		RuntimeBoundary: boundary.ID, RuntimeState: runtimeState,
		Fixtures: len(fixtures.Cases), PositiveFixtures: len(fixtures.Cases)-negative,
		NegativeFixtures: negative, RuntimeDependencies: len(deps), Status: "PASS",
	}, nil
}

func validateHuntManagementContractCore(root string, manifest WorkManifestV2, contract HuntExecutableContract, fixtures HuntFixtureSet, validateSources bool) error {
	if contract.SchemaVersion != 1 || contract.ContractID != "HUNT-WORKSPACE-READONLY-CONTRACT-V1" ||
		contract.ContractKind != "provider-neutral-readonly-local-executable-contract" || contract.CanonicalSchemaClaim {
		return fmt.Errorf("Hunt contract identity/schema claim is invalid")
	}
	if contract.Capability != "CAP-INV-005" {
		return fmt.Errorf("Hunt capability mismatch")
	}
	if contract.RuntimeBoundary.ID != huntManagementRuntimeBoundaryID || contract.RuntimeBoundary.Root != huntManagementRuntimeRoot ||
		(contract.RuntimeBoundary.ExpectedState != "preimplementation" && contract.RuntimeBoundary.ExpectedState != "implemented") ||
		len(contract.RuntimeBoundary.ExternalRuntimeDependencies) != 0 {
		return fmt.Errorf("Hunt runtime-boundary contract is invalid")
	}
	if validateSources {
		if err := validateProductRefs(root, contract.ProductRefs); err != nil {
			return err
		}
	}
	if !productRefsEqual(contract.ProductRefs, manifest.ProductRefs) {
		return fmt.Errorf("Hunt Product Spec references do not match active work manifest")
	}
	if validateSources {
		if err := validateHuntManagementSourceAnchors(root); err != nil {
			return err
		}
	}
	scope := contract.ScopePolicy
	if !scope.TenantRequired || !scope.EnvironmentRequired || !scope.QuestionRequired || !scope.ScopeRequired ||
		!scope.TimeRangeRequired || !scope.OwnerRequired || !scope.WildcardTenantForbidden || !scope.CrossTenantForbidden {
		return fmt.Errorf("Hunt scope policy is incomplete or unsafe")
	}
	refPolicy := contract.ReferencePolicy
	if !sameStringSet(refPolicy.AllowedTypes, []string{"query","search-job","hypothesis","case","incident"}) ||
		!refPolicy.TenantScoped || !refPolicy.StableIdentifierRequired || !refPolicy.CallerAccessDecisionRequired ||
		!refPolicy.DeniedReferenceProjectionForbidden || !refPolicy.OwnershipTransferForbidden ||
		!refPolicy.ProvenanceRequired || !refPolicy.SearchJobRequiresVisibleQuery ||
		!refPolicy.UnavailableStateBecomesLimitation || !refPolicy.StaleStateBecomesLimitation ||
		!refPolicy.PartialStateBecomesLimitation {
		return fmt.Errorf("Hunt reference policy is incomplete or unsafe")
	}
	proj := contract.ProjectionPolicy
	if !proj.Immutable || !proj.DeepCopyIsolation || !proj.QuestionVisible || !proj.ScopeVisible ||
		!proj.TimeRangeVisible || !proj.OwnerVisible || !proj.ContributorsVisible || !proj.ProvenanceVisible ||
		!proj.LimitationsVisible || !proj.SynthesisOfMissingFactsForbidden {
		return fmt.Errorf("Hunt projection policy is incomplete")
	}
	mut := contract.MutationPolicy
	if mut.CanonicalHuntObjectCreated || mut.HuntStateMutationExecutable || mut.ContributorMutationExecutable ||
		mut.CasePromotionMutationExecutable || mut.CaseLinkMutationExecutable || mut.HypothesisMutationExecutable ||
		mut.QueryMutationExecutable || mut.QueryExecutionExecutable || mut.SavedSearchMutationExecutable ||
		mut.SearchJobExecutionExecutable {
		return fmt.Errorf("Hunt mutation policy violates bounded read-only scope")
	}
	tech := contract.TechnologyPolicy
	if tech.StorageEngineSelected || tech.ProviderSelected || tech.RetentionPolicySelected ||
		tech.CollaborationBackendSelected || tech.FinalQueryDialectSelected || tech.FinalHuntUISelected {
		return fmt.Errorf("Hunt contract selects unresolved technology")
	}
	for _, exclusion := range []string{
		"canonical-hunt-object","final-hunt-state-machine","hunt-mutation:OPEN-013",
		"case-promotion-or-link-mutation:OPEN-013","hypothesis-mutation:OPEN-013",
		"query-or-query-asset-mutation:OPEN-013","search-job-execution","storage-provider-retention",
		"collaboration-backend","final-query-dialect","final-hunt-ui",
	} {
		if !containsString(contract.Exclusions, exclusion) {
			return fmt.Errorf("Hunt exclusion %s is missing", exclusion)
		}
	}
	if len(contract.SemanticRules) < 8 || len(contract.SecurityInvariants) < 6 {
		return fmt.Errorf("Hunt semantic/security contract is incomplete")
	}
	if fixtures.SchemaVersion != 1 || fixtures.ContractID != contract.ContractID {
		return fmt.Errorf("Hunt fixtures do not bind the executable contract")
	}
	return validateHuntFixtures(contract, fixtures)
}

func validateHuntManagementSourceAnchors(root string) error {
	sources := []struct{
		path string
		anchors []string
	}{
		{huntManagementCapabilityPath, []string{
			"id: CAP-INV-005",
			"sans imposer encore un objet Hunt canonique",
			"OPEN-013",
			"Queries et Search Jobs restent Shared",
			"aucun objet Hunt n’est créé implicitement",
		}},
		{huntManagementRuntimeMarker, []string{
			"Hunt Management runtime boundary",
			"Preimplementation only",
			"OPEN-013",
		}},
	}
	for _, source := range sources {
		data, err := readRepoFile(root, source.path)
		if err != nil { return err }
		content := string(data)
		for _, anchor := range source.anchors {
			if !strings.Contains(content, anchor) {
				return fmt.Errorf("Hunt source anchor %q missing from %s", anchor, source.path)
			}
		}
	}
	return nil
}

func validateHuntFixtures(contract HuntExecutableContract, fixtures HuntFixtureSet) error {
	required := map[string]bool{
		"valid-minimal-workspace":false,
		"valid-referenced-workspace":false,
		"stale-partial-unavailable-become-limitations":false,
		"missing-tenant-is-rejected":false,
		"wildcard-tenant-is-rejected":false,
		"cross-tenant-workspace-is-rejected":false,
		"missing-environment-is-rejected":false,
		"missing-question-is-rejected":false,
		"missing-scope-is-rejected":false,
		"invalid-period-is-rejected":false,
		"missing-owner-is-rejected":false,
		"cross-tenant-reference-is-rejected":false,
		"denied-reference-is-rejected":false,
		"search-job-without-visible-query-is-rejected":false,
		"class-2-mutation-is-rejected":false,
	}
	if len(fixtures.Cases) < len(required) {
		return fmt.Errorf("Hunt fixture set is too small: %d", len(fixtures.Cases))
	}
	ids := map[string]bool{}
	for _, fixture := range fixtures.Cases {
		if strings.TrimSpace(fixture.ID)=="" || ids[fixture.ID] {
			return fmt.Errorf("Hunt fixture id is empty or duplicate: %q", fixture.ID)
		}
		ids[fixture.ID]=true
		if _, ok := required[fixture.ID]; ok { required[fixture.ID]=true }
		want := evaluateHuntFixture(contract, fixture.Input)
		if !huntExpectedEqual(fixture.Expected, want) {
			return fmt.Errorf("Hunt fixture %s expectation mismatch: want %#v got %#v", fixture.ID, want, fixture.Expected)
		}
	}
	for id,present := range required {
		if !present { return fmt.Errorf("required Hunt fixture %s is missing", id) }
	}
	return nil
}

func evaluateHuntFixture(contract HuntExecutableContract, in HuntFixtureInput) HuntFixtureExpected {
	out := HuntFixtureExpected{AuditRequired:true, Limitations:[]string{}}
	reject := func(code string) HuntFixtureExpected {
		out.ErrorCode=code
		out.Limitations=[]string{}
		return out
	}
	if strings.TrimSpace(in.MutationRequested)!="" { return reject("mutation-open-decision") }
	if strings.TrimSpace(in.TenantRef)=="" { return reject("missing-tenant") }
	if in.TenantRef=="*" { return reject("wildcard-tenant") }
	if !containsString(in.AuthorizedTenants,in.TenantRef) { return reject("cross-tenant") }
	if strings.TrimSpace(in.EnvironmentRef)=="" { return reject("missing-environment") }
	if strings.TrimSpace(in.Question)=="" { return reject("missing-question") }
	if strings.TrimSpace(in.Scope)=="" { return reject("missing-scope") }
	start,errStart:=time.Parse(time.RFC3339,in.TimeStart)
	end,errEnd:=time.Parse(time.RFC3339,in.TimeEnd)
	if errStart!=nil || errEnd!=nil || !start.Before(end) { return reject("invalid-time-range") }
	if strings.TrimSpace(in.OwnerRef)=="" { return reject("missing-owner") }
	for _, contributor := range in.ContributorRefs {
		if strings.TrimSpace(contributor)=="" { return reject("invalid-contributor") }
	}
	allowedTypes := map[string]bool{}
	for _, kind := range contract.ReferencePolicy.AllowedTypes { allowedTypes[kind]=true }
	visibleQueries := map[string]bool{}
	for _, ref := range in.References {
		if !allowedTypes[ref.Kind] || strings.TrimSpace(ref.ID)=="" { return reject("invalid-reference") }
		if ref.TenantRef!=in.TenantRef { return reject("cross-tenant-reference") }
		if ref.AccessDecision!="allow" { return reject("reference-access-denied") }
		if ref.Kind=="query" { visibleQueries[ref.ID]=true }
		switch ref.State {
		case "available":
		case "stale","partial","unavailable":
			out.Limitations=append(out.Limitations,ref.Kind+":"+ref.ID+":"+ref.State)
		default:
			return reject("invalid-reference-state")
		}
	}
	for _, ref := range in.References {
		if ref.Kind=="search-job" && (strings.TrimSpace(ref.QueryRef)=="" || !visibleQueries[ref.QueryRef]) {
			return reject("search-job-query-not-visible")
		}
	}
	sort.Strings(out.Limitations)
	out.Allowed=true
	out.ProjectedReferences=len(in.References)
	out.ImmutableProjection=true
	out.DeepCopyIsolated=true
	return out
}

func huntExpectedEqual(a,b HuntFixtureExpected) bool {
	return a.Allowed==b.Allowed && a.ErrorCode==b.ErrorCode &&
		a.ProjectedReferences==b.ProjectedReferences && sameStringSet(a.Limitations,b.Limitations) &&
		a.ImmutableProjection==b.ImmutableProjection && a.DeepCopyIsolated==b.DeepCopyIsolated &&
		a.MutationExecuted==b.MutationExecuted && a.AuditRequired==b.AuditRequired
}

func validateHuntManagementRuntimeLayout(root string, boundary ArchitectureBoundary, expected string) (string,error) {
	scanRoot,err:=runtimeScanRoot(root,huntManagementRuntimeRoot)
	if err!=nil { return "",err }
	files:=[]string{}
	hasGoMod:=false
	hasRuntimeGo:=false
	// #nosec G703 -- scanRoot is repository-confined by runtimeScanRoot; symlink entries are rejected.
	err=filepath.WalkDir(scanRoot,func(path string,entry os.DirEntry,err error) error {
		if err!=nil { return err }
		if entry.Type()&os.ModeSymlink!=0 { return fmt.Errorf("Hunt runtime boundary contains symlink: %s",path) }
		if entry.IsDir(){ return nil }
		rel,err:=filepath.Rel(root,path)
		if err!=nil { return err }
		rel=filepath.ToSlash(rel)
		files=append(files,rel)
		if rel=="product-runtime/hunt-management/go.mod" { hasGoMod=true }
		if strings.HasSuffix(rel,".go") && !strings.HasSuffix(rel,"_test.go") { hasRuntimeGo=true }
		if expected=="preimplementation" && rel!=huntManagementRuntimeMarker && !strings.HasSuffix(rel,"/.gitkeep") {
			return fmt.Errorf("Hunt preimplementation boundary contains unexpected runtime file %s",rel)
		}
		return nil
	})
	if err!=nil { return "",err }
	sort.Strings(files)
	if !containsString(files,huntManagementRuntimeMarker){ return "",fmt.Errorf("Hunt runtime marker is missing") }
	actual,err:=detectRuntimeBoundaryImplementationState(root,boundary)
	if err!=nil { return "",err }
	if actual!=expected { return "",fmt.Errorf("Hunt runtime state mismatch: contract=%s actual=%s",expected,actual) }
	if expected=="implemented" && (!hasGoMod || !hasRuntimeGo) {
		return "",fmt.Errorf("Hunt implemented runtime requires go.mod and executable Go source")
	}
	return actual,nil
}
