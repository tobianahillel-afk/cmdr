package eventsearch

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"
	"time"
)

type fixtureBackend struct {
	result BackendResult
	err    error
	calls  int
	got    BackendRequest
}

func (b *fixtureBackend) Execute(_ context.Context, req BackendRequest) (BackendResult, error) {
	b.calls++
	b.got = req
	return b.result, b.err
}

func TestSearchJobPartialPreservesValidReferencesAndFailedSources(t *testing.T) {
	fixture := lifecycleFixture(t, "partial-preserves-valid-results")
	backend := &fixtureBackend{result: BackendResult{
		TenantRef:        fixture.Input.TenantRef,
		CompletedSources: []string{"source-a"},
		ResultRefs:       []ResultReference{{Source: "source-a", Ref: "result:source-a:1"}},
		FailedSources:    []SourceFailure{{Source: "source-b", Code: "source-timeout"}},
	}}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(runtimeInput(fixture), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	got, err := orchestrator.Execute(context.Background(), job, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	if got.State != JobPartial || len(got.ResultRefs) != 1 || len(got.FailedSources) != 1 {
		t.Fatalf("partial result lost evidence: %#v", got)
	}
	if got.ResultRefs[0].Ref != "result:source-a:1" || got.FailedSources[0].Source != "source-b" {
		t.Fatalf("partial result provenance mismatch: %#v", got)
	}
	if backend.got.TenantRef != fixture.Input.TenantRef ||
		!reflect.DeepEqual(backend.got.Sources, fixture.Input.Sources) ||
		backend.got.QueryVersion != "query-v1" || backend.got.CorrelationID != "corr-job-001" {
		t.Fatalf("backend request lost validated provenance: %#v", backend.got)
	}
}

func TestSearchJobRejectsCrossTenantBackendResultWithoutLeakingReferences(t *testing.T) {
	backend := &fixtureBackend{result: BackendResult{
		TenantRef:        "tenant-b",
		CompletedSources: []string{"source-a"},
		ResultRefs:       []ResultReference{{Source: "source-a", Ref: "foreign-result"}},
	}}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(validInput(), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	got, err := orchestrator.Execute(context.Background(), job, testAudit())
	if !errors.Is(err, ErrBackendResultInvalid) {
		t.Fatalf("expected backend result rejection, got %v", err)
	}
	if got.State != JobFailed || len(got.ResultRefs) != 0 || len(got.FailedSources) != 0 {
		t.Fatalf("cross-tenant backend data leaked into job: %#v", got)
	}
}

func TestSearchJobBackendMustAccountForEveryRequestedSource(t *testing.T) {
	backend := &fixtureBackend{result: BackendResult{
		TenantRef:        "tenant-a",
		CompletedSources: []string{},
	}}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(validInput(), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	got, err := orchestrator.Execute(context.Background(), job, testAudit())
	if !errors.Is(err, ErrBackendResultInvalid) || got.State != JobFailed {
		t.Fatalf("unaccounted source must fail closed: job=%#v err=%v", got, err)
	}
}

func TestSearchJobCancellationDoesNotFabricateCompletion(t *testing.T) {
	fixture := lifecycleFixture(t, "running-job-can-be-cancelled")
	backend := &fixtureBackend{}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(runtimeInput(fixture), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	job, err = orchestrator.Start(job, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	got, err := orchestrator.Cancel(job, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	if got.State != JobCancelled || len(got.ResultRefs) != 0 || len(got.FailedSources) != 0 {
		t.Fatalf("cancel fabricated completion evidence: %#v", got)
	}
	if backend.calls != 0 {
		t.Fatalf("explicit cancellation should not call backend, calls=%d", backend.calls)
	}
}

func TestSearchJobContextCancellationProducesCancelledState(t *testing.T) {
	backend := &fixtureBackend{}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(validInput(), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	got, err := orchestrator.Execute(ctx, job, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	if got.State != JobCancelled || backend.calls != 0 {
		t.Fatalf("cancelled context executed backend or wrong state: %#v calls=%d", got, backend.calls)
	}
}

func TestSearchJobTerminalStateCannotRestartInPlace(t *testing.T) {
	fixture := lifecycleFixture(t, "illegal-terminal-transition-is-rejected")
	backend := &fixtureBackend{}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(runtimeInput(fixture), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	job.State = JobCompleted
	job.Version++
	got, err := orchestrator.Execute(context.Background(), job, testAudit())
	if !errors.Is(err, ErrIllegalTransition) || got.State != JobCompleted || backend.calls != 0 {
		t.Fatalf("terminal job restarted in place: job=%#v err=%v", got, err)
	}
}

func TestSearchJobRetryCreatesNewIdentityAndPreservesProvenance(t *testing.T) {
	backend := &fixtureBackend{result: BackendResult{
		TenantRef:        "tenant-a",
		CompletedSources: []string{"source-a"},
	}}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(validInput(), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	completed, err := orchestrator.Execute(context.Background(), job, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	retry, err := orchestrator.Retry(completed, AuditContext{ActorRef: "analyst-1", Rationale: "retry after reviewing prior run"})
	if err != nil {
		t.Fatal(err)
	}
	if retry.State != JobQueued || retry.ID == completed.ID || retry.RunID == completed.RunID ||
		retry.Provenance.ParentRunID != completed.RunID {
		t.Fatalf("retry did not create a new run identity: old=%#v retry=%#v", completed, retry)
	}
	if retry.Provenance.TenantRef != completed.Provenance.TenantRef ||
		retry.Provenance.EnvironmentRef != completed.Provenance.EnvironmentRef ||
		!retry.Provenance.TimeStart.Equal(completed.Provenance.TimeStart) ||
		!retry.Provenance.TimeEnd.Equal(completed.Provenance.TimeEnd) ||
		!reflect.DeepEqual(retry.Provenance.Sources, completed.Provenance.Sources) ||
		retry.Provenance.CorrelationID != completed.Provenance.CorrelationID ||
		retry.Provenance.QueryVersion != completed.Provenance.QueryVersion {
		t.Fatalf("retry lost query/scope provenance: old=%#v retry=%#v", completed.Provenance, retry.Provenance)
	}
	if len(retry.ResultRefs) != 0 || len(retry.FailedSources) != 0 || len(retry.AuditTrail) != 1 {
		t.Fatalf("retry inherited previous execution evidence: %#v", retry)
	}
}

func TestSearchJobBackendErrorFailsWithoutPersistingErrorText(t *testing.T) {
	backend := &fixtureBackend{err: errors.New("provider secret diagnostic must not persist")}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(validInput(), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	got, err := orchestrator.Execute(context.Background(), job, testAudit())
	if !errors.Is(err, ErrBackendExecution) {
		t.Fatalf("expected backend execution error, got %v", err)
	}
	if got.State != JobFailed || len(got.AuditTrail) == 0 {
		t.Fatalf("backend failure not represented safely: %#v", got)
	}
	last := got.AuditTrail[len(got.AuditTrail)-1]
	if !reflect.DeepEqual(last.ErrorCodes, []string{"backend-execution"}) {
		t.Fatalf("backend error text leaked or stable code missing: %#v", last)
	}
}

func TestSearchJobRejectsInvalidAuditContext(t *testing.T) {
	backend := &fixtureBackend{}
	orchestrator := testOrchestrator(t, backend)
	for _, audit := range []AuditContext{
		{},
		{ActorRef: "actor with spaces", Rationale: "valid"},
		{ActorRef: "analyst-1", Rationale: ""},
		{ActorRef: "analyst-1", Rationale: "line\nbreak"},
	} {
		if _, err := orchestrator.Create(validInput(), audit); !errors.Is(err, ErrInvalidAuditContext) {
			t.Fatalf("expected invalid audit rejection for %#v, got %v", audit, err)
		}
	}
}

func TestSearchJobCopiesCallerOwnedSlices(t *testing.T) {
	backend := &fixtureBackend{result: BackendResult{
		TenantRef:        "tenant-a",
		CompletedSources: []string{"source-a"},
		ResultRefs:       []ResultReference{{Source: "source-a", Ref: "result-1"}},
	}}
	orchestrator := testOrchestrator(t, backend)
	in := validInput()
	job, err := orchestrator.Create(in, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	in.Sources[0] = "mutated-source"
	if job.Provenance.Sources[0] != "source-a" {
		t.Fatal("job provenance aliases caller input")
	}
	got, err := orchestrator.Execute(context.Background(), job, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	backend.result.ResultRefs[0].Ref = "mutated-result"
	if got.ResultRefs[0].Ref != "result-1" {
		t.Fatal("job results alias backend-owned result slice")
	}
}

func lifecycleFixture(t *testing.T, id string) struct {
	ID    string
	Input struct {
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
	} `json:"input"`
	Expected map[string]any `json:"expected"`
} {
	t.Helper()
	data, err := os.ReadFile("../../engineering/implementation/event-search/fixtures.json")
	if err != nil {
		t.Fatal(err)
	}
	var set struct {
		Cases []struct {
			ID    string
			Input struct {
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
			} `json:"input"`
			Expected map[string]any `json:"expected"`
		} `json:"cases"`
	}
	if err := json.Unmarshal(data, &set); err != nil {
		t.Fatal(err)
	}
	for _, fixture := range set.Cases {
		if fixture.ID == id {
			return fixture
		}
	}
	t.Fatalf("fixture %s not found", id)
	panic("unreachable")
}

func runtimeInput(fixture struct {
	ID    string
	Input struct {
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
	} `json:"input"`
	Expected map[string]any `json:"expected"`
}) Input {
	return Input{
		TenantRef: fixture.Input.TenantRef, AuthorizedTenants: append([]string(nil), fixture.Input.AuthorizedTenants...),
		EnvironmentRef: fixture.Input.EnvironmentRef, TimeStart: fixture.Input.TimeStart, TimeEnd: fixture.Input.TimeEnd,
		Sources: append([]string(nil), fixture.Input.Sources...), AuthorizedSources: append([]string(nil), fixture.Input.AuthorizedSources...),
		Permissions: append([]string(nil), fixture.Input.Permissions...), QueryPresent: fixture.Input.QueryPresent,
		QueryValid: fixture.Input.QueryValid, CorrelationID: "corr-job-001", QueryVersion: "query-v1",
		CaseLinkRequested: fixture.Input.CaseLinkRequested,
	}
}

func testAudit() AuditContext {
	return AuditContext{ActorRef: "analyst-1", Rationale: "bounded Event Search test"}
}

func testOrchestrator(t *testing.T, backend Backend) *Orchestrator {
	t.Helper()
	clockValue := time.Date(2026, 9, 24, 12, 0, 0, 0, time.UTC)
	counter := 0
	ids := func(prefix string) (string, error) {
		counter++
		return prefix + "-test-" + string(rune('0'+counter)), nil
	}
	orchestrator, err := newOrchestrator(backend, func() time.Time {
		clockValue = clockValue.Add(time.Millisecond)
		return clockValue
	}, ids)
	if err != nil {
		t.Fatal(err)
	}
	return orchestrator
}
