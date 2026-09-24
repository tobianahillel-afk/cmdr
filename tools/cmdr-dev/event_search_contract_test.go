package main

import "testing"

func TestEventSearchExecutableContractValidates(t *testing.T) {
	manifest, contract, fixtures := validEventSearchContractTestData()
	if err := validateEventSearchContractCore("", manifest, contract, fixtures, false); err != nil {
		t.Fatal(err)
	}
}

func TestEventSearchFixtureRejectsWildcardTenant(t *testing.T) {
	_, contract, fixtures := validEventSearchContractTestData()
	for i := range fixtures.Cases {
		if fixtures.Cases[i].ID == "wildcard-tenant-is-rejected" {
			fixtures.Cases[i].Expected.Allowed = true
			if err := validateEventSearchFixtures(contract, fixtures); err == nil {
				t.Fatal("expected wildcard tenant fixture mismatch rejection")
			}
			return
		}
	}
	t.Fatal("wildcard fixture missing")
}

func TestEventSearchFixtureRejectsProtectedDataOnPermissionDenied(t *testing.T) {
	_, contract, fixtures := validEventSearchContractTestData()
	for i := range fixtures.Cases {
		if fixtures.Cases[i].ID == "permission-denied-masks-protected-data" {
			fixtures.Cases[i].Expected.ProtectedDataVisible = true
			if err := validateEventSearchFixtures(contract, fixtures); err == nil {
				t.Fatal("expected permission-denied disclosure rejection")
			}
			return
		}
	}
	t.Fatal("permission-denied fixture missing")
}

func TestEventSearchContractRejectsBackendSelection(t *testing.T) {
	manifest, contract, fixtures := validEventSearchContractTestData()
	contract.QueryPolicy.FinalIndexSelected = true
	if err := validateEventSearchContractCore("", manifest, contract, fixtures, false); err == nil {
		t.Fatal("expected final index selection rejection")
	}
}

func TestEventSearchContractRequiresOpenDecisionExclusion(t *testing.T) {
	manifest, contract, fixtures := validEventSearchContractTestData()
	contract.Exclusions = []string{"final-ui", "final-query-language", "final-index", "final-storage-engine", "final-provider"}
	if err := validateEventSearchContractCore("", manifest, contract, fixtures, false); err == nil {
		t.Fatal("expected OPEN-013 exclusion rejection")
	}
}

func TestEventSearchTerminalTransitionIsRejected(t *testing.T) {
	_, contract, _ := validEventSearchContractTestData()
	in := validEventSearchInput()
	in.CurrentJobState = "completed"
	in.RequestedJobState = "running"
	got := evaluateEventSearchFixture(contract, in)
	if got.Allowed || got.ErrorCode != "illegal-transition" {
		t.Fatalf("unexpected terminal transition result: %#v", got)
	}
}

func validEventSearchContractTestData() (WorkManifestV2, EventSearchExecutableContract, EventSearchFixtureSet) {
	refs := ProductRefsV2{
		Capabilities: []string{"CAP-INV-002"},
		Requirements: []string{"REQ-PROD-005", "REQ-PROD-010", "REQ-PROD-011", "REQ-PROD-014", "REQ-PROD-045"},
		Screens: []string{"INV-EVS-001"},
		Permissions: []string{
			"perm.investigate.search.execute", "perm.shared-capabilities.query.read",
			"perm.shared-capabilities.query.manage", "perm.shared-capabilities.search-job.read",
			"perm.shared-capabilities.search-job.manage", "perm.shared-capabilities.telemetry-event.read",
		},
		OpenDecisions: []string{"OPEN-013"},
		ImplementationContracts: []string{
			"cmdr-product-spec/17-implementation-contracts/search-query-contract.md",
			"cmdr-product-spec/17-implementation-contracts/background-job-contract.md",
			"cmdr-product-spec/17-implementation-contracts/tenant-scope-contract.md",
			"cmdr-product-spec/17-implementation-contracts/permission-enforcement-contract.md",
			"cmdr-product-spec/17-implementation-contracts/audit-event-contract.md",
			"cmdr-product-spec/17-implementation-contracts/error-contract.md",
			"cmdr-product-spec/17-implementation-contracts/correlation-identifiers.md",
		},
		CanonicalObjects: []string{
			"cmdr-product-spec/05-domain-model/objects/query.md",
			"cmdr-product-spec/05-domain-model/objects/search-job.md",
			"cmdr-product-spec/05-domain-model/objects/telemetry-event.md",
		},
	}
	manifest := WorkManifestV2{ID: "E10-INV-002A-CONTRACT", ProductRefs: refs}
	contract := EventSearchExecutableContract{
		SchemaVersion: 1, ContractID: "EVENT-SEARCH-EXECUTABLE-CONTRACT-V1",
		ContractKind: "backend-neutral-local-executable-contract", Capability: "CAP-INV-002", Screen: "INV-EVS-001",
		RuntimeBoundary: EventSearchRuntimeBoundary{ID: eventSearchRuntimeBoundaryID, Root: eventSearchRuntimeRoot, ExpectedState: "preimplementation"},
		ProductRefs: refs,
		ScopePolicy: EventSearchScopePolicy{
			TenantRequired: true, EnvironmentRequired: true, TimeRangeRequired: true, SourcesRequired: true,
			WildcardTenantForbidden: true, CrossTenantForbidden: true, EmptySourcesForbidden: true,
		},
		PermissionPolicy: EventSearchPermissionPolicy{
			ServerEvaluationRequired: true, ClientTrustForbidden: true,
			RequiredExecutionPermissions: []string{
				"perm.investigate.search.execute", "perm.shared-capabilities.query.read",
				"perm.shared-capabilities.search-job.manage", "perm.shared-capabilities.telemetry-event.read",
			},
			InspectionPermission: "perm.shared-capabilities.search-job.read",
			QueryManagementPermission: "perm.shared-capabilities.query.manage",
			DenyMasksProtectedData: true, RawAccessIsSeparate: true,
		},
		QueryPolicy: EventSearchQueryPolicy{
			QueryRequired: true, ValidationBeforeExecution: true,
			DeterministicManualPathRequired: true, AIExecutionWithoutExplicitTriggerForbidden: true,
		},
		SearchJobPolicy: EventSearchJobPolicy{
			States: []string{"queued", "running", "completed", "partial", "failed", "cancelled"},
			TerminalStates: []string{"completed", "partial", "failed", "cancelled"},
			AllowedTransitions: []string{
				"queued->running", "queued->cancelled", "running->completed",
				"running->partial", "running->failed", "running->cancelled",
			},
			RelaunchCreatesNewJob: true, CancellationAudited: true,
			PartialPreservesValidResults: true, FailedSourcesVisible: true, NoCompleteResultSimulation: true,
		},
		AuditPolicy: EventSearchAuditPolicy{
			ActorRequired: true, TenantRequired: true, ActionRequired: true, TargetRequired: true,
			OutcomeRequired: true, RationaleRequired: true, CorrelationIDRequired: true,
			QueryVersionRequired: true, PeriodRequired: true, SourcesRequired: true,
			RunIDRequired: true, ErrorsRecorded: true, SensitiveDataForbiddenInCorrelationID: true,
		},
		PerformanceBudget: EventSearchPerformanceBudget{ValidationP95MS: 1, OrchestrationP95MS: 1},
		Exclusions: []string{"case-link-mutation:OPEN-013", "final-ui", "final-query-language", "final-index", "final-storage-engine", "final-provider"},
		SemanticRules: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		SecurityInvariants: []string{"1", "2", "3", "4", "5", "6"},
	}
	fixtures := EventSearchFixtureSet{SchemaVersion: 1, ContractID: contract.ContractID}
	add := func(id string, mutate func(*EventSearchFixtureInput)) {
		in := validEventSearchInput()
		if mutate != nil {
			mutate(&in)
		}
		fixtures.Cases = append(fixtures.Cases, EventSearchFixtureCase{ID: id, Input: in, Expected: evaluateEventSearchFixture(contract, in)})
	}
	add("valid-execution-envelope", nil)
	add("partial-preserves-valid-results", func(in *EventSearchFixtureInput) {
		in.Sources = []string{"source-a", "source-b"}
		in.AuthorizedSources = []string{"source-a", "source-b"}
		in.PartialFailedSources = []string{"source-b"}
	})
	add("missing-tenant-is-rejected", func(in *EventSearchFixtureInput) { in.TenantRef = "" })
	add("wildcard-tenant-is-rejected", func(in *EventSearchFixtureInput) {
		in.TenantRef = "*"
		in.AuthorizedTenants = []string{"*"}
	})
	add("cross-tenant-is-rejected", func(in *EventSearchFixtureInput) { in.TenantRef = "tenant-b" })
	add("permission-denied-masks-protected-data", func(in *EventSearchFixtureInput) {
		in.Permissions = []string{
			"perm.shared-capabilities.query.read",
			"perm.shared-capabilities.search-job.manage",
			"perm.shared-capabilities.telemetry-event.read",
		}
	})
	add("invalid-time-range-is-rejected", func(in *EventSearchFixtureInput) {
		in.TimeStart = "2026-09-24T12:00:00Z"
		in.TimeEnd = "2026-09-24T11:00:00Z"
	})
	add("empty-sources-are-rejected", func(in *EventSearchFixtureInput) {
		in.Sources = nil
		in.AuthorizedSources = nil
	})
	add("illegal-terminal-transition-is-rejected", func(in *EventSearchFixtureInput) {
		in.CurrentJobState = "completed"
		in.RequestedJobState = "running"
	})
	add("case-link-remains-excluded", func(in *EventSearchFixtureInput) { in.CaseLinkRequested = true })
	return manifest, contract, fixtures
}

func validEventSearchInput() EventSearchFixtureInput {
	return EventSearchFixtureInput{
		TenantRef: "tenant-a", AuthorizedTenants: []string{"tenant-a"}, EnvironmentRef: "env-prod",
		TimeStart: "2026-09-24T10:00:00Z", TimeEnd: "2026-09-24T11:00:00Z",
		Sources: []string{"source-a"}, AuthorizedSources: []string{"source-a"},
		Permissions: []string{
			"perm.investigate.search.execute", "perm.shared-capabilities.query.read",
			"perm.shared-capabilities.search-job.manage", "perm.shared-capabilities.telemetry-event.read",
		},
		QueryPresent: true, QueryValid: true, RequestedJobState: "queued",
	}
}
