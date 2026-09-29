package main

import "testing"

func huntTestContract() HuntExecutableContract {
	return HuntExecutableContract{
		ReferencePolicy: HuntReferencePolicy{
			AllowedTypes: []string{"query", "search-job", "hypothesis", "case", "incident"},
		},
	}
}

func huntTestInput() HuntFixtureInput {
	return HuntFixtureInput{
		TenantRef: "tenant-a", AuthorizedTenants: []string{"tenant-a"}, EnvironmentRef: "env-prod",
		Question: "question", Scope: "scope", TimeStart: "2026-09-28T08:00:00Z", TimeEnd: "2026-09-28T09:00:00Z",
		OwnerRef: "owner-1",
	}
}

func TestEvaluateHuntFixtureAcceptsVisibleReferencesAndLimitations(t *testing.T) {
	in := huntTestInput()
	in.References = []HuntFixtureReference{
		{Kind: "query", ID: "q1", TenantRef: "tenant-a", AccessDecision: "allow", State: "stale"},
		{Kind: "search-job", ID: "j1", TenantRef: "tenant-a", AccessDecision: "allow", QueryRef: "q1", State: "partial"},
	}
	got := evaluateHuntFixture(huntTestContract(), in)
	if !got.Allowed || got.ProjectedReferences != 2 || !sameStringSet(got.Limitations, []string{"query:q1:stale", "search-job:j1:partial"}) {
		t.Fatalf("unexpected projection: %#v", got)
	}
	if !got.ImmutableProjection || !got.DeepCopyIsolated || got.MutationExecuted {
		t.Fatalf("projection invariants missing: %#v", got)
	}
}

func TestEvaluateHuntFixtureRejectsDeniedAndCrossTenantReferences(t *testing.T) {
	for _, tc := range []struct {
		name string
		ref  HuntFixtureReference
		want string
	}{
		{"denied", HuntFixtureReference{Kind: "query", ID: "q1", TenantRef: "tenant-a", AccessDecision: "deny", State: "available"}, "reference-access-denied"},
		{"cross-tenant", HuntFixtureReference{Kind: "query", ID: "q1", TenantRef: "tenant-b", AccessDecision: "allow", State: "available"}, "cross-tenant-reference"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := huntTestInput()
			in.References = []HuntFixtureReference{tc.ref}
			got := evaluateHuntFixture(huntTestContract(), in)
			if got.Allowed || got.ErrorCode != tc.want || got.ProjectedReferences != 0 {
				t.Fatalf("unexpected rejection: %#v", got)
			}
		})
	}
}

func TestEvaluateHuntFixtureRequiresVisibleQueryForSearchJob(t *testing.T) {
	in := huntTestInput()
	in.References = []HuntFixtureReference{
		{Kind: "search-job", ID: "j1", TenantRef: "tenant-a", AccessDecision: "allow", QueryRef: "q1", State: "available"},
	}
	got := evaluateHuntFixture(huntTestContract(), in)
	if got.Allowed || got.ErrorCode != "search-job-query-not-visible" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestEvaluateHuntFixtureRejectsMutationBeforeProjection(t *testing.T) {
	in := huntTestInput()
	in.MutationRequested = "promote-case"
	got := evaluateHuntFixture(huntTestContract(), in)
	if got.Allowed || got.ErrorCode != "mutation-open-decision" || got.MutationExecuted {
		t.Fatalf("unexpected mutation result: %#v", got)
	}
}
