package main

import "testing"

func validCaseQueueContractForTest() CaseQueueExecutableContract {
	return CaseQueueExecutableContract{
		SchemaVersion:1, ContractID:"CASE-QUEUE-READONLY-CONTRACT-V1",
		ContractKind:"provider-neutral-readonly-local-executable-contract",
		Capability:"CAP-INV-101",
		RuntimeBoundary:CaseQueueRuntimeBoundary{ID:caseQueueRuntimeBoundaryID,Root:caseQueueRuntimeRoot,ExpectedState:"preimplementation"},
		ScopePolicy:CaseQueueScopePolicy{TenantRequired:true,EnvironmentRequired:true,WildcardTenantForbidden:true,CrossTenantForbidden:true,CaseAccessDecisionRequired:true},
		SearchPolicy:CaseQueueSearchPolicy{AllowedFields:[]string{"id","human_id","title","owner","next_action"},CaseInsensitive:true,TrimSpace:true,InaccessibleDataForbidden:true},
		FilterPolicy:CaseQueueFilterPolicy{AllowedFields:[]string{"status","owner","incident_priority","finding_status","freshness"},UnknownFilterRejected:true,EmptyValueRejected:true,CaseInsensitiveExactMatch:true},
		SortPolicy:CaseQueueSortPolicy{
			AllowedKeys:[]string{"updated_at","human_id","title","owner","status"},
			DefaultSort:[]CaseQueueSortSpec{{Field:"updated_at",Direction:"desc"},{Field:"human_id",Direction:"asc"},{Field:"id",Direction:"asc"}},
			StableTieBreakers:[]string{"human_id","id"},UnknownSortRejected:true,
		},
		SavedViewPolicy:CaseQueueSavedViewPolicy{SharedOwned:true,Optional:true,StableIdentifierRequired:true,VersionRequired:true,TenantScoped:true,CallerAccessDecisionRequired:true,FieldPermissionReevaluationRequired:true,FilterPermissionReevaluationRequired:true,SortPermissionReevaluationRequired:true,DeniedElementsOmitted:true,MutationForbidden:true},
		FreshnessPolicy:CaseQueueFreshnessPolicy{NamedStates:[]string{"available","empty","partial","stale","permission-filtered","view-dirty"},MissingProjectionNamed:true,StaleProjectionNamed:true,UnavailableProjectionNamed:true,DeterministicSortedDiagnostics:true,SynthesisOfMissingFactsForbidden:true},
		ProjectionPolicy:CaseQueueProjectionPolicy{Immutable:true,DeepCopyIsolation:true,CaseIsOnlyOwnedWorkItem:true,IncidentContextOptional:true,FindingContextOptional:true,ActivityContextOptional:true},
		Exclusions:[]string{
			"case-create-update-assignment-lifecycle:CAP-INV-102:OPEN-013","saved-view-create-update-share-archive-manage",
			"export-job-creation-or-export-backend","task-decision-response-run-ownership-or-aggregation","command-work-queue-behavior",
			"canonical-case-schema-change","final-case-state-machine","final-queue-columns","storage-provider-retention","collaboration-backend","final-cap-inv-101-ui",
		},
		SemanticRules:[]string{"1","2","3","4","5","6","7","8","9"},
		SecurityInvariants:[]string{"1","2","3","4","5","6","7","8"},
	}
}

func baseCaseQueueInputForTest() CaseQueueFixtureInput {
	return CaseQueueFixtureInput{
		TenantRef:"tenant-a",AuthorizedTenants:[]string{"tenant-a"},EnvironmentRef:"env",
		Cases:[]CaseQueueCaseInput{
			{ID:"case-b",TenantRef:"tenant-a",HumanID:"CASE-002",Title:"Beta",Status:"review",Owner:"alice",UpdatedAt:"2026-09-29T11:00:00Z",AccessDecision:"allow",Freshness:"available",ActivityState:"available"},
			{ID:"case-a",TenantRef:"tenant-a",HumanID:"CASE-001",Title:"Alpha",Status:"investigating",Owner:"alice",UpdatedAt:"2026-09-29T10:00:00Z",AccessDecision:"allow",Freshness:"available",ActivityState:"available"},
		},
	}
}

func TestCaseQueueStableDefaultSort(t *testing.T) {
	got:=evaluateCaseQueueFixture(validCaseQueueContractForTest(),baseCaseQueueInputForTest())
	if !got.Allowed || !stringSlicesEqual(got.ProjectedCaseIDs,[]string{"case-b","case-a"}) {
		t.Fatalf("unexpected projection: %#v",got)
	}
	if !stringSlicesEqual(got.AppliedSort,[]string{"updated_at:desc","human_id:asc","id:asc"}) {
		t.Fatalf("unexpected stable sort: %#v",got.AppliedSort)
	}
}

func TestCaseQueueDeniedCaseNeverLeaksIdentity(t *testing.T) {
	in:=baseCaseQueueInputForTest()
	in.Cases[0].AccessDecision="deny"
	got:=evaluateCaseQueueFixture(validCaseQueueContractForTest(),in)
	if !got.Allowed || containsString(got.ProjectedCaseIDs,"case-b") || got.PermissionFilteredCases!=1 ||
		!containsString(got.States,"permission-filtered") {
		t.Fatalf("denied Case leakage/filter mismatch: %#v",got)
	}
	for _,diagnostic:=range got.Diagnostics {
		if diagnostic=="case-b" || diagnostic=="case:case-b:permission-filtered" {
			t.Fatalf("denied Case identity leaked in diagnostic: %q",diagnostic)
		}
	}
}

func TestCaseQueueSavedViewPermissionReevaluation(t *testing.T) {
	in:=baseCaseQueueInputForTest()
	in.SavedView=&CaseQueueSavedViewInput{
		ID:"view-1",TenantRef:"tenant-a",Version:"v1",AccessDecision:"allow",
		Fields:[]CaseQueueSavedViewField{{Name:"title",AccessDecision:"allow"},{Name:"secret",AccessDecision:"deny"}},
		Filters:[]CaseQueueSavedViewFilter{{Field:"owner",Values:[]string{"alice"},AccessDecision:"allow"},{Field:"status",Values:[]string{"closed"},AccessDecision:"deny"}},
	}
	got:=evaluateCaseQueueFixture(validCaseQueueContractForTest(),in)
	if !got.Allowed || got.PermissionFilteredElements!=2 ||
		!stringSlicesEqual(got.AppliedSavedViewFields,[]string{"title"}) ||
		!containsString(got.States,"permission-filtered") {
		t.Fatalf("Saved View permission reevaluation failed: %#v",got)
	}
}

func TestCaseQueuePartialAndStaleAreNamed(t *testing.T) {
	in:=baseCaseQueueInputForTest()
	in.Cases=in.Cases[:1]
	in.Cases[0].Freshness="stale"
	in.Cases[0].Incident=&CaseQueueIncidentContext{ID:"inc",TenantRef:"tenant-a",AccessDecision:"allow",State:"unavailable"}
	got:=evaluateCaseQueueFixture(validCaseQueueContractForTest(),in)
	if !got.Allowed || !containsString(got.States,"stale") || !containsString(got.States,"partial") {
		t.Fatalf("partial/stale state missing: %#v",got)
	}
	if !containsString(got.Diagnostics,"case:case-b:stale") || !containsString(got.Diagnostics,"incident-context:case-b:unavailable") {
		t.Fatalf("partial/stale diagnostics missing: %#v",got.Diagnostics)
	}
}

func TestCaseQueueMutationAndCrossTenantFailClosed(t *testing.T) {
	contract:=validCaseQueueContractForTest()
	in:=baseCaseQueueInputForTest()
	in.MutationRequested="assign"
	if got:=evaluateCaseQueueFixture(contract,in); got.Allowed || got.ErrorCode!="mutation-open-decision" {
		t.Fatalf("mutation not rejected: %#v",got)
	}
	in=baseCaseQueueInputForTest()
	in.Cases[0].TenantRef="tenant-b"
	if got:=evaluateCaseQueueFixture(contract,in); got.Allowed || got.ErrorCode!="cross-tenant-case" {
		t.Fatalf("cross-tenant Case not rejected: %#v",got)
	}
}
