package main

import "testing"

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
func ptrBool(v bool) *bool { return &v }

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
			BoundaryID: "api", ChangedSecurityCritical: true,
			ChangedSecurityCriticalPercent: ptrFloat(91),
			AuthorizationNegativePassed: ptrBool(true), TenantIsolationNegativePassed: ptrBool(true),
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
		SecurityCriticalPaths: []string{"services/api/auth/**"},
		AuthorizationRationale: "authorization is outside this synthetic scope",
		TenantIsolationRationale: "tenant isolation is outside this synthetic scope",
	}
	commit := "0123456789abcdef0123456789abcdef01234567"
	for _, tc := range []RuntimeSecurityEvidence{
		{
			SchemaVersion: 1, Status: "measured", SourceCommit: commit, GlobalCoveragePercent: ptrFloat(79.99),
			Scopes: []RuntimeSecurityEvidenceScope{{BoundaryID: "api", ChangedSecurityCritical: false}},
		},
		{
			SchemaVersion: 1, Status: "measured", SourceCommit: commit, GlobalCoveragePercent: ptrFloat(80),
			Scopes: []RuntimeSecurityEvidenceScope{{
				BoundaryID: "api", ChangedSecurityCritical: true, ChangedSecurityCriticalPercent: ptrFloat(89.99),
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
		SecurityCriticalPaths: []string{"services/api/**"},
		AuthorizationRationale: "not required", TenantIsolationRationale: "not required",
	}
	evidence := RuntimeSecurityEvidence{
		SchemaVersion: 1, Status: "measured",
		SourceCommit: "0123456789abcdef0123456789abcdef01234567",
		GlobalCoveragePercent: ptrFloat(90),
		Scopes: []RuntimeSecurityEvidenceScope{{BoundaryID: "api", ChangedSecurityCritical: true, ChangedSecurityCriticalPercent: ptrFloat(95)}},
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
		SecurityCriticalPaths: []string{"services/other/**"},
		AuthorizationRationale: "not required", TenantIsolationRationale: "not required",
	}
	if err := validateRuntimeSecurityScope(scope, boundary); err == nil {
		t.Fatal("expected out-of-bound security path rejection")
	}
}
