package main

import "testing"

func TestPilotExecutableContractRejectsCanonicalSchemaClaim(t *testing.T) {
	scope, contract, fixtures := validPilotContractTestData()
	contract.CanonicalSchemaClaim = true
	if err := validatePilotExecutableContract(scope, contract, fixtures); err == nil {
		t.Fatal("expected canonical schema claim rejection")
	}
}

func TestPilotExecutableContractRejectsRuntimeDependency(t *testing.T) {
	scope, contract, fixtures := validPilotContractTestData()
	contract.RuntimeBoundary.ExternalRuntimeDependencies = []string{"example.org/runtime"}
	if err := validatePilotExecutableContract(scope, contract, fixtures); err == nil {
		t.Fatal("expected external runtime dependency rejection")
	}
}

func TestPilotFixturesRejectCrossTenantEnvironmentLeak(t *testing.T) {
	_, contract, fixtures := validPilotContractTestData()
	fixture := fixtures.Cases[0]
	fixture.ID = "environment-tenant-mismatch-is-cleared"
	fixture.Input.EnvironmentTenantRef = "tenant-b"
	fixture.Expected.EnvironmentRef = "env-b"
	fixtures.Cases = append(fixtures.Cases, fixture)
	if err := validatePilotFixtures(contract, fixtures); err == nil {
		t.Fatal("expected cross-tenant environment leak rejection")
	}
}

func TestPilotFixturesRejectProtectedRefsOnPermissionDenied(t *testing.T) {
	_, contract, fixtures := validPilotContractTestData()
	for i := range fixtures.Cases {
		if fixtures.Cases[i].ID == "permission-denied-masks-protected-refs" {
			fixtures.Cases[i].Expected.TenantRef = "tenant-a"
			if err := validatePilotFixtures(contract, fixtures); err == nil {
				t.Fatal("expected permission-denied protected-ref rejection")
			}
			return
		}
	}
	t.Fatal("permission-denied fixture missing")
}

func validPilotContractTestData() (PilotScope, PilotExecutableContract, PilotFixtureSet) {
	scope := validPilotScopeForTest()
	scope.Requirements = []string{"REQ-PROD-008"}
	scope.Permissions = []string{"perm.platform-settings.tenant.read"}
	contract := PilotExecutableContract{
		SchemaVersion: 1,
		ContractID: "PILOT-CONTEXT-ENVELOPE-V1",
		ContractKind: "pilot-local-executable-contract",
		SliceID: scope.SliceID,
		Capability: scope.Capability,
		RuntimeBoundary: PilotContextRuntimeBoundary{
			ID: pilotRuntimeBoundaryID, Root: pilotRuntimeRoot,
			PackageIdentity: "example/context-envelope",
		},
		ProductRefs: pilotProductRefs(scope),
		ActionClasses: []string{"0", "1"},
		AuthorizationOutcomes: []string{"allow", "deny"},
		FunctionalStates: []string{"context-valid", "context-incompatible", "permission-denied", "expired", "source-missing"},
		ErrorCodes: []string{"invalid-tenant-ref", "unsafe-return-origin"},
		ReferencePolicy: PilotReferencePolicy{
			Kind: "opaque-stable-reference", MinBytes: 1, MaxBytes: 128,
			WildcardForbidden: true, WhitespaceForbidden: true, ControlCharactersForbidden: true,
		},
		ReturnOriginPolicy: PilotReturnOriginPolicy{
			RelativeRouteOnly: true, MaxBytes: 1024, QueryForbidden: true,
			FragmentForbidden: true, SchemeForbidden: true,
		},
		SemanticRules: []string{"1","2","3","4","5","6","7","8"},
		SecurityInvariants: []string{"1","2","3","4","5"},
	}
	base := PilotFixtureCase{
		ID: "preserve-authorized-compatible-environment",
		Input: PilotFixtureInput{
			SourceProduct:"platform-settings", DestinationProduct:"command",
			SourceTenantRef:"tenant-a", DestinationTenantRef:"tenant-a",
			EnvironmentRef:"env-a", EnvironmentTenantRef:"tenant-a",
			DestinationRequiresEnvironment:true, EnvironmentCompatible:true,
			DestinationAuthorization:"allow", ContextFresh:true, ReturnOrigin:"/settings",
		},
		Expected: PilotFixtureExpected{
			State:"context-valid", TenantRef:"tenant-a", EnvironmentRef:"env-a",
			ReturnOrigin:"/settings", ProtectedRefsVisible:true,
		},
	}
	fixtures := PilotFixtureSet{SchemaVersion:1,ContractID:contract.ContractID,Cases:[]PilotFixtureCase{base}}
	required := []PilotFixtureCase{
		{ID:"tenant-change-clears-inherited-environment",Input:PilotFixtureInput{SourceTenantRef:"tenant-a",DestinationTenantRef:"tenant-b",DestinationAuthorization:"allow"},Expected:PilotFixtureExpected{State:"context-incompatible",TenantRef:"tenant-b",Cleared:true}},
		{ID:"environment-tenant-mismatch-is-cleared",Input:PilotFixtureInput{SourceTenantRef:"tenant-a",DestinationTenantRef:"tenant-a",EnvironmentRef:"env-b",EnvironmentTenantRef:"tenant-b",DestinationAuthorization:"allow"},Expected:PilotFixtureExpected{State:"context-incompatible",TenantRef:"tenant-a",Cleared:true}},
		{ID:"permission-denied-masks-protected-refs",Input:PilotFixtureInput{SourceTenantRef:"tenant-a",DestinationTenantRef:"tenant-a",DestinationAuthorization:"deny"},Expected:PilotFixtureExpected{State:"permission-denied",TraceMinimalOnly:true}},
		{ID:"expired-context-masks-inherited-refs",Input:PilotFixtureInput{SourceTenantRef:"tenant-a",DestinationTenantRef:"tenant-a",DestinationAuthorization:"allow"},Expected:PilotFixtureExpected{State:"expired",TraceMinimalOnly:true}},
		{ID:"missing-source-tenant-blocks",Input:PilotFixtureInput{DestinationAuthorization:"allow"},Expected:PilotFixtureExpected{State:"source-missing",TraceMinimalOnly:true}},
		{ID:"wildcard-tenant-is-rejected",Input:PilotFixtureInput{SourceTenantRef:"*",DestinationTenantRef:"*",DestinationAuthorization:"allow"},Expected:PilotFixtureExpected{Error:"invalid-tenant-ref",TraceMinimalOnly:true}},
		{ID:"unsafe-return-origin-query-is-rejected",Input:PilotFixtureInput{SourceTenantRef:"tenant-a",DestinationTenantRef:"tenant-a",DestinationAuthorization:"allow"},Expected:PilotFixtureExpected{Error:"unsafe-return-origin",TraceMinimalOnly:true}},
		{ID:"missing-required-environment-needs-choice",Input:PilotFixtureInput{SourceTenantRef:"tenant-a",DestinationTenantRef:"tenant-a",DestinationAuthorization:"allow"},Expected:PilotFixtureExpected{State:"context-incompatible",TenantRef:"tenant-a"}},
		{ID:"tenant-only-optional-environment",Input:PilotFixtureInput{SourceTenantRef:"tenant-a",DestinationTenantRef:"tenant-a",DestinationAuthorization:"allow"},Expected:PilotFixtureExpected{State:"context-valid",TenantRef:"tenant-a"}},
	}
	fixtures.Cases = append(fixtures.Cases, required...)
	return scope, contract, fixtures
}
