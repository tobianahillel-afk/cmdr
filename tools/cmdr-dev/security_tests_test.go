package main

import (
	"os"
	"path/filepath"
	"testing"
)

func securityTestArchitecture(runtime ...ArchitectureBoundary) ArchitectureRegistry {
	boundaries := []ArchitectureBoundary{{
		ID: "product-spec", Kind: "product-documentation", Roots: []string{"cmdr-product-spec/**"},
	}}
	boundaries = append(boundaries, runtime...)
	return ArchitectureRegistry{SchemaVersion: 1, DefaultPolicy: "deny-unregistered-path", Boundaries: boundaries}
}

func securityTestGates(authActive, tenantActive bool) SecurityGateRegistry {
	authReadiness, authKey := "deferred-runtime", ""
	if authActive {
		authReadiness, authKey = "active", "auth-test-v1"
	}
	tenantReadiness, tenantKey := "deferred-runtime", ""
	if tenantActive {
		tenantReadiness, tenantKey = "active", "tenant-test-v1"
	}
	return SecurityGateRegistry{Gates: []SecurityGate{
		{ID: "SEC-AUTH-NEG-001", Readiness: authReadiness, ImplementationKey: authKey},
		{ID: "SEC-TENANT-ISO-001", Readiness: tenantReadiness, ImplementationKey: tenantKey},
	}}
}

func baseRuntimePolicy(scopes ...RuntimeSecurityScope) RuntimeSecurityPolicy {
	return RuntimeSecurityPolicy{
		SchemaVersion: 1, DefaultPolicy: "deny-unregistered-runtime-boundary",
		GlobalCoverageFloorPercent: 80, ChangedSecurityCriticalCoverageFloor: 90,
		AuthorizationGateID: "SEC-AUTH-NEG-001", TenantIsolationGateID: "SEC-TENANT-ISO-001",
		Scopes: scopes,
	}
}

func ptrFloat(v float64) *float64 { return &v }
func ptrBool(v bool) *bool        { return &v }

func TestSecurityTestAuditReportsNotApplicableWithoutRuntime(t *testing.T) {
	evidence := RuntimeSecurityEvidence{
		SchemaVersion: 1, Status: "not-applicable",
		Reason: "no product runtime",
	}
	summary, err := evaluateSecurityTestPolicy(
		securityTestArchitecture(), securityTestGates(false, false),
		baseRuntimePolicy(), evidence, nil, "0123456789abcdef0123456789abcdef01234567",
	)
	if err != nil {
		t.Fatal(err)
	}
	if summary.CoverageStatus != "not-applicable" || summary.RuntimeBoundaries != 0 || summary.RegisteredScopes != 0 {
		t.Fatalf("unexpected not-applicable summary: %#v", summary)
	}
}

func TestSecurityTestAuditRejectsUnregisteredRuntimeBoundary(t *testing.T) {
	runtime := ArchitectureBoundary{ID: "api", Kind: "product-runtime", Roots: []string{"services/api/**"}}
	evidence := RuntimeSecurityEvidence{SchemaVersion: 1, Status: "measured", SourceCommit: "0123456789abcdef0123456789abcdef01234567", GlobalCoveragePercent: ptrFloat(90)}
	if _, err := evaluateSecurityTestPolicy(
		securityTestArchitecture(runtime), securityTestGates(false, false),
		baseRuntimePolicy(), evidence, []string{"services/api/a.go"}, evidence.SourceCommit,
	); err == nil {
		t.Fatal("expected unregistered runtime boundary rejection")
	}
}

func TestSecurityTestAuditEnforcesCoverageAndNegativeTests(t *testing.T) {
	runtime := ArchitectureBoundary{ID: "api", Kind: "product-runtime", Roots: []string{"services/api/**"}}
	scope := RuntimeSecurityScope{
		BoundaryID: "api", Owner: "security-platform",
		SecurityCriticalPaths: []string{"services/api/auth/**"},
		AuthorizationRequired: true, AuthorizationRationale: "API authorizes subject access",
		TenantIsolationRequired: true, TenantIsolationRationale: "API is tenant scoped",
	}
	commit := "0123456789abcdef0123456789abcdef01234567"
	evidence := RuntimeSecurityEvidence{
		SchemaVersion: 1, Status: "measured", SourceCommit: commit, GlobalCoveragePercent: ptrFloat(82),
		Scopes: []RuntimeSecurityEvidenceScope{{
			BoundaryID: "api", ImplementationState: "implemented", ChangedSecurityCritical: true,
			ChangedSecurityCriticalPercent: ptrFloat(91),
			AuthorizationNegativePassed:    ptrBool(true), TenantIsolationNegativePassed: ptrBool(true),
		}},
	}
	summary, err := evaluateSecurityTestPolicy(
		securityTestArchitecture(runtime), securityTestGates(true, true),
		baseRuntimePolicy(scope), evidence, []string{"services/api/auth/policy.go"}, commit,
	)
	if err != nil {
		t.Fatal(err)
	}
	if summary.CoverageStatus != "measured" || summary.ChangedSecurityCriticalScopes != 1 {
		t.Fatalf("unexpected measured summary: %#v", summary)
	}
}

func TestSecurityTestAuditRejectsCoverageBelowFloors(t *testing.T) {
	runtime := ArchitectureBoundary{ID: "api", Kind: "product-runtime", Roots: []string{"services/api/**"}}
	scope := RuntimeSecurityScope{
		BoundaryID: "api", Owner: "security-platform",
		SecurityCriticalPaths:    []string{"services/api/auth/**"},
		AuthorizationRationale:   "authorization is outside this synthetic scope",
		TenantIsolationRationale: "tenant isolation is outside this synthetic scope",
	}
	commit := "0123456789abcdef0123456789abcdef01234567"
	for _, tc := range []RuntimeSecurityEvidence{
		{
			SchemaVersion: 1, Status: "measured", SourceCommit: commit, GlobalCoveragePercent: ptrFloat(79.99),
			Scopes: []RuntimeSecurityEvidenceScope{{BoundaryID: "api", ImplementationState: "implemented", ChangedSecurityCritical: false}},
		},
		{
			SchemaVersion: 1, Status: "measured", SourceCommit: commit, GlobalCoveragePercent: ptrFloat(80),
			Scopes: []RuntimeSecurityEvidenceScope{{
				BoundaryID: "api", ImplementationState: "implemented", ChangedSecurityCritical: true, ChangedSecurityCriticalPercent: ptrFloat(89.99),
			}},
		},
	} {
		if _, err := evaluateSecurityTestPolicy(
			securityTestArchitecture(runtime), securityTestGates(false, false),
			baseRuntimePolicy(scope), tc, []string{"services/api/auth/policy.go"}, commit,
		); err == nil {
			t.Fatal("expected coverage floor rejection")
		}
	}
}

func TestSecurityTestAuditRejectsStaleEvidence(t *testing.T) {
	runtime := ArchitectureBoundary{ID: "api", Kind: "product-runtime", Roots: []string{"services/api/**"}}
	scope := RuntimeSecurityScope{
		BoundaryID: "api", Owner: "security-platform",
		SecurityCriticalPaths:  []string{"services/api/**"},
		AuthorizationRationale: "not required", TenantIsolationRationale: "not required",
	}
	evidence := RuntimeSecurityEvidence{
		SchemaVersion: 1, Status: "measured",
		SourceCommit:          "0123456789abcdef0123456789abcdef01234567",
		GlobalCoveragePercent: ptrFloat(90),
		Scopes:                []RuntimeSecurityEvidenceScope{{BoundaryID: "api", ImplementationState: "implemented", ChangedSecurityCritical: true, ChangedSecurityCriticalPercent: ptrFloat(95)}},
	}
	if _, err := evaluateSecurityTestPolicy(
		securityTestArchitecture(runtime), securityTestGates(false, false),
		baseRuntimePolicy(scope), evidence, []string{"services/api/x.go"},
		"abcdef0123456789abcdef0123456789abcdef01",
	); err == nil {
		t.Fatal("expected stale evidence rejection")
	}
}

func TestSecurityScopeRejectsPathOutsideBoundary(t *testing.T) {
	boundary := ArchitectureBoundary{ID: "api", Kind: "product-runtime", Roots: []string{"services/api/**"}}
	scope := RuntimeSecurityScope{
		BoundaryID: "api", Owner: "security-platform",
		SecurityCriticalPaths:  []string{"services/other/**"},
		AuthorizationRationale: "not required", TenantIsolationRationale: "not required",
	}
	if err := validateRuntimeSecurityScope(scope, boundary); err == nil {
		t.Fatal("expected out-of-bound security path rejection")
	}
}

func TestSecurityTestAuditAcceptsReservedPreimplementationRuntime(t *testing.T) {
	runtime := ArchitectureBoundary{ID: "api", Kind: "product-runtime", Roots: []string{"services/api/**"}}
	scope := RuntimeSecurityScope{
		BoundaryID: "api", Owner: "E9-PILOT-001C",
		SecurityCriticalPaths: []string{"services/api/**"},
		AuthorizationRequired: true, AuthorizationRationale: "authorization is required before executable runtime",
		TenantIsolationRequired: true, TenantIsolationRationale: "tenant isolation is required before executable runtime",
	}
	evidence := RuntimeSecurityEvidence{
		SchemaVersion: 1, Status: "preimplementation",
		Reason: "runtime boundary is reserved but contains no executable source",
		Scopes: []RuntimeSecurityEvidenceScope{{
			BoundaryID: "api", ImplementationState: "preimplementation",
		}},
	}
	summary, err := evaluateSecurityTestPolicyWithImplementationStates(
		securityTestArchitecture(runtime), securityTestGates(false, false),
		baseRuntimePolicy(scope), evidence, []string{"services/api/README.md"},
		"0123456789abcdef0123456789abcdef01234567",
		map[string]string{"api": "preimplementation"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if summary.CoverageStatus != "preimplementation" || summary.PreimplementationScopes != 1 ||
		summary.ImplementedScopes != 0 || summary.AuthorizationRequired != 1 || summary.TenantIsolationRequired != 1 {
		t.Fatalf("unexpected preimplementation summary: %#v", summary)
	}
}

func TestSecurityTestAuditRejectsPreimplementationClaimAfterRuntimeImplementation(t *testing.T) {
	runtime := ArchitectureBoundary{ID: "api", Kind: "product-runtime", Roots: []string{"services/api/**"}}
	scope := RuntimeSecurityScope{
		BoundaryID: "api", Owner: "security-platform",
		SecurityCriticalPaths:  []string{"services/api/**"},
		AuthorizationRationale: "not required yet", TenantIsolationRationale: "not required yet",
	}
	evidence := RuntimeSecurityEvidence{
		SchemaVersion: 1, Status: "preimplementation", Reason: "stale reservation",
		Scopes: []RuntimeSecurityEvidenceScope{{BoundaryID: "api", ImplementationState: "preimplementation"}},
	}
	if _, err := evaluateSecurityTestPolicyWithImplementationStates(
		securityTestArchitecture(runtime), securityTestGates(false, false),
		baseRuntimePolicy(scope), evidence, []string{"services/api/runtime.go"},
		"0123456789abcdef0123456789abcdef01234567",
		map[string]string{"api": "implemented"},
	); err == nil {
		t.Fatal("expected stale preimplementation evidence rejection")
	}
}

func TestDetectRuntimeBoundaryImplementationState(t *testing.T) {
	root := t.TempDir()
	runtimeRoot := filepath.Join(root, "services", "api")
	if err := os.MkdirAll(runtimeRoot, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(runtimeRoot, "README.md"), []byte("reserved\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	boundary := ArchitectureBoundary{ID: "api", Kind: "product-runtime", Roots: []string{"services/api/**"}}
	state, err := detectRuntimeBoundaryImplementationState(root, boundary)
	if err != nil {
		t.Fatal(err)
	}
	if state != "preimplementation" {
		t.Fatalf("expected preimplementation, got %s", state)
	}
	if err := os.WriteFile(filepath.Join(runtimeRoot, "runtime.go"), []byte("package api\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	state, err = detectRuntimeBoundaryImplementationState(root, boundary)
	if err != nil {
		t.Fatal(err)
	}
	if state != "implemented" {
		t.Fatalf("expected implemented, got %s", state)
	}
}

func TestParsePilotCoverageProfileMeasuresGlobalAndChangedStatements(t *testing.T) {
	profile := []byte("mode: atomic\n" +
		pilotSecurityModuleIdentity + "/envelope.go:10.1,12.2 2 1\n" +
		pilotSecurityModuleIdentity + "/envelope.go:14.1,18.2 3 0\n" +
		pilotSecurityModuleIdentity + "/other.go:1.1,2.2 1 1\n")
	got, err := parsePilotCoverageProfile(profile, []string{"product-runtime/context-envelope/envelope.go"})
	if err != nil {
		t.Fatal(err)
	}
	if got.GlobalPercent != 50 {
		t.Fatalf("global coverage: got %.2f want 50", got.GlobalPercent)
	}
	if got.ChangedExecutableStmts != 5 || got.ChangedExecutablePercent != 40 {
		t.Fatalf("changed coverage mismatch: %#v", got)
	}
}

func TestParsePilotCoverageProfileRejectsForeignModule(t *testing.T) {
	profile := []byte("mode: atomic\nexample.invalid/other/file.go:1.1,2.2 1 1\n")
	if _, err := parsePilotCoverageProfile(profile, nil); err == nil {
		t.Fatal("expected foreign coverage source rejection")
	}
}

func TestRuntimeSecurityAdaptersAreExplicitForKnownRuntimes(t *testing.T) {
	adapters := runtimeSecurityAdapters()
	for _, id := range []string{pilotRuntimeBoundaryID, eventSearchRuntimeBoundaryID} {
		adapter, ok := adapters[id]
		if !ok {
			t.Fatalf("missing explicit runtime security adapter %s", id)
		}
		if adapter.RuntimeKind != "go" || adapter.RuntimeRoot == "" || adapter.ModuleIdentity == "" ||
			adapter.AuthorizationTestRegex == "" || adapter.TenantTestRegex == "" {
			t.Fatalf("incomplete Go runtime security adapter %s: %#v", id, adapter)
		}
	}
	frontend, ok := adapters["event-search-frontend-runtime"]
	if !ok {
		t.Fatal("missing explicit Event Search frontend runtime security adapter")
	}
	if frontend.RuntimeKind != "node" || frontend.RuntimeRoot != eventSearchFrontendSecurityRuntimeRoot ||
		len(frontend.TestFiles) != 5 || len(frontend.CoverageIncludes) != 4 ||
		frontend.AuthorizationTestRegex == "" || frontend.TenantTestRegex == "" ||
		!containsString(frontend.TestFiles, "e2e.test.mjs") {
		t.Fatalf("incomplete frontend runtime security adapter: %#v", frontend)
	}
	if len(adapters) != 3 {
		t.Fatalf("unexpected implicit runtime security adapters: %#v", adapters)
	}
}

func TestParseEventSearchCoverageProfileMeasuresExactModule(t *testing.T) {
	adapter := runtimeSecurityAdapters()[eventSearchRuntimeBoundaryID]
	profile := []byte("mode: atomic\n" +
		eventSearchSecurityModuleIdentity + "/validator.go:10.1,12.2 2 1\n" +
		eventSearchSecurityModuleIdentity + "/validator.go:14.1,18.2 3 0\n")
	got, err := parseRuntimeCoverageProfile(profile, []string{"product-runtime/event-search/validator.go"}, adapter)
	if err != nil {
		t.Fatal(err)
	}
	if got.TotalStatements != 5 || got.CoveredStatements != 2 || got.GlobalPercent != 40 {
		t.Fatalf("unexpected Event Search coverage: %#v", got)
	}
	if got.ChangedExecutableStmts != 5 || got.ChangedExecutablePercent != 40 {
		t.Fatalf("unexpected changed Event Search coverage: %#v", got)
	}
}

func TestParseNodeLCOVMeasuresExactFrontendSource(t *testing.T) {
	root := t.TempDir()
	moduleRoot := filepath.Join(root, "product-runtime", "event-search-frontend")
	if err := os.MkdirAll(moduleRoot, 0o750); err != nil {
		t.Fatal(err)
	}
	lcov := []byte("TN:\nSF:state.mjs\nDA:1,1\nDA:2,1\nDA:3,0\nLF:3\nLH:2\nend_of_record\n")
	adapter := runtimeSecurityAdapter{
		BoundaryID:       "event-search-frontend-runtime",
		RuntimeRoot:      "product-runtime/event-search-frontend",
		RuntimeKind:      "node",
		CoverageIncludes: []string{"state.mjs"},
	}
	measurement, err := parseNodeLCOV(root, moduleRoot, lcov, []string{"product-runtime/event-search-frontend/state.mjs"}, adapter)
	if err != nil {
		t.Fatal(err)
	}
	if measurement.TotalStatements != 3 || measurement.CoveredStatements != 2 {
		t.Fatalf("unexpected node coverage counts: %#v", measurement)
	}
	if measurement.GlobalPercent < 66.6 || measurement.ChangedExecutablePercent < 66.6 || measurement.ChangedExecutableStmts != 3 {
		t.Fatalf("unexpected node coverage percentages: %#v", measurement)
	}
}

func TestParseNodeLCOVMeasuresMultipleExplicitFrontendSources(t *testing.T) {
	root := t.TempDir()
	moduleRoot := filepath.Join(root, "product-runtime", "event-search-frontend")
	if err := os.MkdirAll(moduleRoot, 0o750); err != nil {
		t.Fatal(err)
	}
	lcov := []byte("SF:state.mjs\nLF:2\nLH:2\nend_of_record\nSF:shell.mjs\nLF:3\nLH:3\nend_of_record\n")
	adapter := runtimeSecurityAdapter{
		BoundaryID:       "event-search-frontend-runtime",
		RuntimeRoot:      "product-runtime/event-search-frontend",
		RuntimeKind:      "node",
		CoverageIncludes: []string{"state.mjs", "shell.mjs"},
	}
	measurement, err := parseNodeLCOV(
		root,
		moduleRoot,
		lcov,
		[]string{"product-runtime/event-search-frontend/shell.mjs"},
		adapter,
	)
	if err != nil {
		t.Fatal(err)
	}
	if measurement.TotalStatements != 5 || measurement.CoveredStatements != 5 || measurement.GlobalPercent != 100 {
		t.Fatalf("unexpected multi-source node coverage: %#v", measurement)
	}
	if measurement.ChangedExecutableStmts != 3 || measurement.ChangedExecutablePercent != 100 {
		t.Fatalf("unexpected changed shell coverage: %#v", measurement)
	}
}

func TestParseNodeLCOVRejectsCoverageOutsideExpectedSource(t *testing.T) {
	root := t.TempDir()
	moduleRoot := filepath.Join(root, "product-runtime", "event-search-frontend")
	if err := os.MkdirAll(moduleRoot, 0o750); err != nil {
		t.Fatal(err)
	}
	lcov := []byte("SF:other.mjs\nLF:1\nLH:1\nend_of_record\n")
	adapter := runtimeSecurityAdapter{
		BoundaryID:       "event-search-frontend-runtime",
		RuntimeRoot:      "product-runtime/event-search-frontend",
		RuntimeKind:      "node",
		CoverageIncludes: []string{"state.mjs"},
	}
	if _, err := parseNodeLCOV(root, moduleRoot, lcov, nil, adapter); err == nil {
		t.Fatal("expected unexpected Node coverage source rejection")
	}
}
