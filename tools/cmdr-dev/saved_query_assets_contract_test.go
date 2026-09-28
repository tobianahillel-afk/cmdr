package main

import "testing"

func queryAssetTestContract() QueryAssetExecutableContract {
	return QueryAssetExecutableContract{
		AssetPolicy: QueryAssetPolicy{AllowedKinds: []string{"saved-search", "query-asset"}},
	}
}

func queryAssetTestInput() QueryAssetFixtureInput {
	return QueryAssetFixtureInput{
		TenantRef: "tenant-a", AuthorizedTenants: []string{"tenant-a"}, EnvironmentRef: "env-prod",
		AssetID: "asset-1", AssetKind: "saved-search", AssetAccessDecision: "allow",
		QueryRef: "query-1", QueryTenantRef: "tenant-a", QueryVersion: "v3", QueryAccessDecision: "allow",
		AuthorRef: "author-1", ValidationState: "validated",
		Sources: []QueryAssetFixtureSource{{
			ID: "source-1", TenantRef: "tenant-a", State: "available",
			Fields: []QueryAssetFixtureField{{Name: "host.name", State: "available"}},
		}},
	}
}

func TestEvaluateSavedQueryAssetsFixtureAcceptsCompatibleSnapshot(t *testing.T) {
	got := evaluateSavedQueryAssetsFixture(queryAssetTestContract(), queryAssetTestInput())
	if !got.Allowed || got.Compatibility != "compatible" || !got.HandoffEligible {
		t.Fatalf("unexpected compatible projection: %#v", got)
	}
	if !got.ImmutableProjection || !got.DeepCopyIsolated || got.MutationExecuted || got.SearchJobCreated || got.SearchJobExecuted {
		t.Fatalf("projection invariants missing: %#v", got)
	}
}

func TestEvaluateSavedQueryAssetsFixtureMarksMissingAndRemovedPrerequisitesIncompatible(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*QueryAssetFixtureInput)
		want string
	}{
		{"missing-source", func(in *QueryAssetFixtureInput) {
			in.Sources[0].State = "missing"
			in.Sources[0].Fields = nil
		}, "source:source-1:missing"},
		{"removed-field", func(in *QueryAssetFixtureInput) {
			in.Sources[0].Fields[0].State = "removed"
		}, "field:source-1/host.name:removed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := queryAssetTestInput()
			tc.edit(&in)
			got := evaluateSavedQueryAssetsFixture(queryAssetTestContract(), in)
			if !got.Allowed || got.Compatibility != "incompatible" || got.HandoffEligible ||
				!containsString(got.Diagnostics, tc.want) {
				t.Fatalf("unexpected incompatible projection: %#v", got)
			}
		})
	}
}

func TestEvaluateSavedQueryAssetsFixtureMakesStaleSnapshotVisibleButNotExecutable(t *testing.T) {
	in := queryAssetTestInput()
	in.ValidationState = "stale"
	in.Sources[0].Fields[0].State = "stale"
	got := evaluateSavedQueryAssetsFixture(queryAssetTestContract(), in)
	if !got.Allowed || got.Compatibility != "stale" || got.HandoffEligible {
		t.Fatalf("unexpected stale projection: %#v", got)
	}
	if !containsString(got.Diagnostics, "validation:stale") ||
		!containsString(got.Diagnostics, "field:source-1/host.name:stale") {
		t.Fatalf("stale diagnostics missing: %#v", got.Diagnostics)
	}
}

func TestEvaluateSavedQueryAssetsFixtureRejectsAccessAndCrossTenantInputs(t *testing.T) {
	for _, tc := range []struct {
		name string
		edit func(*QueryAssetFixtureInput)
		want string
	}{
		{"asset-denied", func(in *QueryAssetFixtureInput) { in.AssetAccessDecision = "deny" }, "asset-access-denied"},
		{"query-denied", func(in *QueryAssetFixtureInput) { in.QueryAccessDecision = "deny" }, "query-access-denied"},
		{"query-cross-tenant", func(in *QueryAssetFixtureInput) { in.QueryTenantRef = "tenant-b" }, "cross-tenant-query"},
		{"source-cross-tenant", func(in *QueryAssetFixtureInput) { in.Sources[0].TenantRef = "tenant-b" }, "cross-tenant-source"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			in := queryAssetTestInput()
			tc.edit(&in)
			got := evaluateSavedQueryAssetsFixture(queryAssetTestContract(), in)
			if got.Allowed || got.ErrorCode != tc.want || got.HandoffEligible {
				t.Fatalf("unexpected rejection: %#v", got)
			}
		})
	}
}

func TestEvaluateSavedQueryAssetsFixtureRejectsMutationAndSearchJobExecution(t *testing.T) {
	in := queryAssetTestInput()
	in.MutationRequested = "share"
	got := evaluateSavedQueryAssetsFixture(queryAssetTestContract(), in)
	if got.Allowed || got.ErrorCode != "mutation-open-decision" || got.MutationExecuted {
		t.Fatalf("unexpected mutation result: %#v", got)
	}

	in = queryAssetTestInput()
	in.SearchJobExecutionRequested = true
	got = evaluateSavedQueryAssetsFixture(queryAssetTestContract(), in)
	if got.Allowed || got.ErrorCode != "search-job-execution-forbidden" || got.SearchJobCreated || got.SearchJobExecuted {
		t.Fatalf("unexpected Search Job result: %#v", got)
	}
}
