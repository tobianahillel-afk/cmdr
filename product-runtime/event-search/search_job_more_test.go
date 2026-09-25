package eventsearch

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewOrchestratorRejectsMissingDependencies(t *testing.T) {
	backend := &fixtureBackend{}
	if _, err := newOrchestrator(nil, time.Now, randomOpaqueID); err == nil {
		t.Fatal("expected nil backend rejection")
	}
	if _, err := newOrchestrator(backend, nil, randomOpaqueID); err == nil {
		t.Fatal("expected nil clock rejection")
	}
	if _, err := newOrchestrator(backend, time.Now, nil); err == nil {
		t.Fatal("expected nil id generator rejection")
	}
}

func TestCreateRejectsValidationAndIdentityFailures(t *testing.T) {
	backend := &fixtureBackend{}
	orchestrator := testOrchestrator(t, backend)

	invalid := validInput()
	invalid.TenantRef = ""
	if _, err := orchestrator.Create(invalid, testAudit()); err == nil || !strings.Contains(err.Error(), string(ErrorMissingTenant)) {
		t.Fatalf("expected validation rejection, got %v", err)
	}

	badID, err := newOrchestrator(backend, time.Now, func(string) (string, error) {
		return "id with spaces", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := badID.Create(validInput(), testAudit()); !errors.Is(err, ErrInvalidIdentity) {
		t.Fatalf("expected invalid generated identity, got %v", err)
	}

	genErr := errors.New("entropy unavailable")
	brokenID, err := newOrchestrator(backend, time.Now, func(string) (string, error) {
		return "", genErr
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := brokenID.Create(validInput(), testAudit()); !errors.Is(err, ErrInvalidIdentity) {
		t.Fatalf("expected id generation error, got %v", err)
	}

	zeroClock, err := newOrchestrator(backend, func() time.Time { return time.Time{} }, func(prefix string) (string, error) {
		return prefix + "-valid", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := zeroClock.Create(validInput(), testAudit()); err == nil || !strings.Contains(err.Error(), "zero time") {
		t.Fatalf("expected zero clock rejection, got %v", err)
	}
}

func TestStartAndCancelRejectIllegalStatesAndIdentity(t *testing.T) {
	backend := &fixtureBackend{}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(validInput(), testAudit())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := orchestrator.Start(job, AuditContext{}); !errors.Is(err, ErrInvalidAuditContext) {
		t.Fatalf("expected invalid start audit rejection, got %v", err)
	}
	bad := job
	bad.ID = ""
	if _, err := orchestrator.Start(bad, testAudit()); !errors.Is(err, ErrInvalidIdentity) {
		t.Fatalf("expected invalid start identity rejection, got %v", err)
	}

	running, err := orchestrator.Start(job, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := orchestrator.Start(running, testAudit()); !errors.Is(err, ErrIllegalTransition) {
		t.Fatalf("expected double-start rejection, got %v", err)
	}

	if _, err := orchestrator.Cancel(job, AuditContext{}); !errors.Is(err, ErrInvalidAuditContext) {
		t.Fatalf("expected invalid cancel audit rejection, got %v", err)
	}
	bad = job
	bad.RunID = ""
	if _, err := orchestrator.Cancel(bad, testAudit()); !errors.Is(err, ErrInvalidIdentity) {
		t.Fatalf("expected invalid cancel identity rejection, got %v", err)
	}

	cancelled, err := orchestrator.Cancel(job, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := orchestrator.Cancel(cancelled, testAudit()); !errors.Is(err, ErrIllegalTransition) {
		t.Fatalf("expected terminal cancel rejection, got %v", err)
	}
}

func TestExecuteCoversRunningDeadlineAndBackendCancellation(t *testing.T) {
	backend := &fixtureBackend{result: BackendResult{
		TenantRef:        "tenant-a",
		CompletedSources: []string{"source-a"},
	}}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(validInput(), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	running, err := orchestrator.Start(job, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	completed, err := orchestrator.Execute(context.Background(), running, testAudit())
	if err != nil || completed.State != JobCompleted {
		t.Fatalf("running execution failed: job=%#v err=%v", completed, err)
	}

	if _, err := orchestrator.Execute(nil, job, testAudit()); !errors.Is(err, ErrBackendExecution) {
		t.Fatalf("expected nil context rejection, got %v", err)
	}
	if _, err := orchestrator.Execute(context.Background(), job, AuditContext{}); !errors.Is(err, ErrInvalidAuditContext) {
		t.Fatalf("expected invalid execute audit rejection, got %v", err)
	}
	bad := job
	bad.Provenance.QueryVersion = ""
	if _, err := orchestrator.Execute(context.Background(), bad, testAudit()); !errors.Is(err, ErrInvalidIdentity) {
		t.Fatalf("expected invalid execute identity rejection, got %v", err)
	}

	deadlineCtx, cancelDeadline := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer cancelDeadline()
	failed, err := orchestrator.Execute(deadlineCtx, job, testAudit())
	if !errors.Is(err, ErrBackendExecution) || failed.State != JobFailed {
		t.Fatalf("deadline must fail execution: job=%#v err=%v", failed, err)
	}

	cancelBackend := &fixtureBackend{err: context.Canceled}
	cancelOrchestrator := testOrchestrator(t, cancelBackend)
	cancelJob, err := cancelOrchestrator.Create(validInput(), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	cancelled, err := cancelOrchestrator.Execute(context.Background(), cancelJob, testAudit())
	if err != nil || cancelled.State != JobCancelled {
		t.Fatalf("backend cancellation must cancel job: job=%#v err=%v", cancelled, err)
	}
}

func TestExecuteClassifiesAllFailedSourcesWithoutFabricatedResults(t *testing.T) {
	backend := &fixtureBackend{result: BackendResult{
		TenantRef: "tenant-a",
		FailedSources: []SourceFailure{{
			Source: "source-a",
			Code:   "backend-unavailable",
		}},
	}}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(validInput(), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	got, err := orchestrator.Execute(context.Background(), job, testAudit())
	if err != nil {
		t.Fatal(err)
	}
	if got.State != JobFailed || len(got.ResultRefs) != 0 || len(got.FailedSources) != 1 {
		t.Fatalf("all-source failure misclassified: %#v", got)
	}
}

func TestRetryRejectsInvalidRequestsAndIdentityReuse(t *testing.T) {
	backend := &fixtureBackend{}
	orchestrator := testOrchestrator(t, backend)
	job, err := orchestrator.Create(validInput(), testAudit())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := orchestrator.Retry(job, testAudit()); !errors.Is(err, ErrIllegalTransition) {
		t.Fatalf("expected nonterminal retry rejection, got %v", err)
	}
	job.State = JobFailed

	if _, err := orchestrator.Retry(job, AuditContext{}); !errors.Is(err, ErrInvalidAuditContext) {
		t.Fatalf("expected invalid retry audit rejection, got %v", err)
	}
	bad := job
	bad.ID = ""
	if _, err := orchestrator.Retry(bad, testAudit()); !errors.Is(err, ErrInvalidIdentity) {
		t.Fatalf("expected invalid retry identity rejection, got %v", err)
	}

	reuseCalls := 0
	reuse, err := newOrchestrator(backend, time.Now, func(prefix string) (string, error) {
		reuseCalls++
		if prefix == "job" {
			return job.ID, nil
		}
		return job.RunID, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := reuse.Retry(job, testAudit()); !errors.Is(err, ErrInvalidIdentity) {
		t.Fatalf("expected identity-reuse rejection, got %v (calls=%d)", err, reuseCalls)
	}

	generatorCalls := 0
	broken, err := newOrchestrator(backend, time.Now, func(prefix string) (string, error) {
		generatorCalls++
		if generatorCalls == 1 {
			return "job-new", nil
		}
		return "", errors.New("run id unavailable")
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := broken.Retry(job, testAudit()); !errors.Is(err, ErrInvalidIdentity) {
		t.Fatalf("expected retry id generation failure, got %v", err)
	}

	zeroClock, err := newOrchestrator(backend, func() time.Time { return time.Time{} }, func(prefix string) (string, error) {
		return prefix + "-new", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := zeroClock.Retry(job, testAudit()); err == nil || !strings.Contains(err.Error(), "zero time") {
		t.Fatalf("expected retry zero-clock rejection, got %v", err)
	}
}

func TestValidateBackendResultRejectsMalformedEvidence(t *testing.T) {
	job := validSearchJobForBackendTest()
	tests := []struct {
		name   string
		result BackendResult
	}{
		{
			name: "tenant mismatch",
			result: BackendResult{TenantRef: "tenant-b", CompletedSources: []string{"source-a"}},
		},
		{
			name: "unknown completed source",
			result: BackendResult{TenantRef: "tenant-a", CompletedSources: []string{"source-x"}},
		},
		{
			name: "duplicate completed source",
			result: BackendResult{TenantRef: "tenant-a", CompletedSources: []string{"source-a", "source-a"}},
		},
		{
			name: "unknown failed source",
			result: BackendResult{TenantRef: "tenant-a", CompletedSources: []string{"source-a"}, FailedSources: []SourceFailure{{Source: "source-x", Code: "x"}}},
		},
		{
			name: "source both completed and failed",
			result: BackendResult{TenantRef: "tenant-a", CompletedSources: []string{"source-a"}, FailedSources: []SourceFailure{{Source: "source-a", Code: "x"}}},
		},
		{
			name: "duplicate failed source",
			result: BackendResult{TenantRef: "tenant-a", CompletedSources: []string{"source-a"}, FailedSources: []SourceFailure{{Source: "source-b", Code: "x"}, {Source: "source-b", Code: "y"}}},
		},
		{
			name: "invalid failure code",
			result: BackendResult{TenantRef: "tenant-a", CompletedSources: []string{"source-a"}, FailedSources: []SourceFailure{{Source: "source-b", Code: "bad code"}}},
		},
		{
			name: "unaccounted source",
			result: BackendResult{TenantRef: "tenant-a", CompletedSources: []string{"source-a"}},
		},
		{
			name: "result belongs to failed source",
			result: BackendResult{TenantRef: "tenant-a", CompletedSources: []string{"source-a"}, FailedSources: []SourceFailure{{Source: "source-b", Code: "timeout"}}, ResultRefs: []ResultReference{{Source: "source-b", Ref: "r1"}}},
		},
		{
			name: "invalid result reference",
			result: BackendResult{TenantRef: "tenant-a", CompletedSources: []string{"source-a", "source-b"}, ResultRefs: []ResultReference{{Source: "source-a", Ref: "bad ref"}}},
		},
		{
			name: "duplicate result reference",
			result: BackendResult{TenantRef: "tenant-a", CompletedSources: []string{"source-a", "source-b"}, ResultRefs: []ResultReference{{Source: "source-a", Ref: "r1"}, {Source: "source-a", Ref: "r1"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, refs, failures, err := validateBackendResult(job, tt.result)
			if !errors.Is(err, ErrBackendResultInvalid) || state != JobFailed || refs != nil || failures != nil {
				t.Fatalf("malformed backend evidence accepted: state=%s refs=%#v failures=%#v err=%v", state, refs, failures, err)
			}
		})
	}
}

func TestValidateBackendResultSortsStableEvidence(t *testing.T) {
	job := validSearchJobForBackendTest()
	state, refs, failures, err := validateBackendResult(job, BackendResult{
		TenantRef:        "tenant-a",
		CompletedSources: []string{"source-a"},
		ResultRefs: []ResultReference{
			{Source: "source-a", Ref: "result-z"},
			{Source: "source-a", Ref: "result-a"},
		},
		FailedSources: []SourceFailure{{Source: "source-b", Code: "timeout"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if state != JobPartial || refs[0].Ref != "result-a" || refs[1].Ref != "result-z" ||
		len(failures) != 1 || failures[0].Source != "source-b" {
		t.Fatalf("backend evidence is not deterministic: state=%s refs=%#v failures=%#v", state, refs, failures)
	}
}

func TestValidateJobIdentityRejectsEveryRequiredInvariant(t *testing.T) {
	base := validSearchJobForBackendTest()
	base.State = JobQueued
	base.Version = 1
	base.ID = "job-valid"
	base.RunID = "run-valid"
	base.CreatedAt = time.Now().UTC()
	base.UpdatedAt = base.CreatedAt
	base.Provenance.CorrelationID = "corr-valid"
	base.Provenance.QueryVersion = "query-v1"

	mutations := []func(*SearchJob){
		func(j *SearchJob) { j.ID = "" },
		func(j *SearchJob) { j.RunID = "" },
		func(j *SearchJob) { j.Provenance.CorrelationID = "" },
		func(j *SearchJob) { j.Provenance.QueryVersion = "" },
		func(j *SearchJob) { j.Provenance.TenantRef = "" },
		func(j *SearchJob) { j.Provenance.EnvironmentRef = "" },
		func(j *SearchJob) { j.Provenance.Sources = nil },
		func(j *SearchJob) { j.Provenance.TimeEnd = j.Provenance.TimeStart },
		func(j *SearchJob) { j.Version = 0 },
		func(j *SearchJob) { j.CreatedAt = time.Time{} },
		func(j *SearchJob) { j.UpdatedAt = time.Time{} },
		func(j *SearchJob) { j.Provenance.ParentRunID = "bad parent" },
	}
	for i, mutate := range mutations {
		job := cloneJob(base)
		mutate(&job)
		if err := validateJobIdentity(job); !errors.Is(err, ErrInvalidIdentity) {
			t.Fatalf("mutation %d was not rejected: job=%#v err=%v", i, job, err)
		}
	}
	base.Provenance.ParentRunID = "run-parent"
	if err := validateJobIdentity(base); err != nil {
		t.Fatalf("valid parent run rejected: %v", err)
	}
}

func TestAuditContextAndTerminalStateHelpers(t *testing.T) {
	if validAuditContext(AuditContext{ActorRef: "actor-1", Rationale: strings.Repeat("a", 1025)}) {
		t.Fatal("oversized rationale accepted")
	}
	if validAuditContext(AuditContext{ActorRef: "actor-1", Rationale: "bad	control"}) {
		t.Fatal("control character accepted")
	}
	if !validAuditContext(testAudit()) {
		t.Fatal("valid audit context rejected")
	}
	for _, state := range []JobState{JobCompleted, JobPartial, JobFailed, JobCancelled} {
		if !terminalState(state) {
			t.Fatalf("terminal state not recognized: %s", state)
		}
	}
	for _, state := range []JobState{"", JobQueued, JobRunning} {
		if terminalState(state) {
			t.Fatalf("nonterminal state recognized as terminal: %s", state)
		}
	}
}

func TestCloneJobDeepCopiesNestedSlices(t *testing.T) {
	job := validSearchJobForBackendTest()
	job.ResultRefs = []ResultReference{{Source: "source-a", Ref: "r1"}}
	job.FailedSources = []SourceFailure{{Source: "source-b", Code: "timeout"}}
	job.AuditTrail = []AuditEvent{{ErrorCodes: []string{"e1"}}}
	clone := cloneJob(job)

	clone.Provenance.Sources[0] = "changed"
	clone.ResultRefs[0].Ref = "changed"
	clone.FailedSources[0].Code = "changed"
	clone.AuditTrail[0].ErrorCodes[0] = "changed"

	if job.Provenance.Sources[0] != "source-a" || job.ResultRefs[0].Ref != "r1" ||
		job.FailedSources[0].Code != "timeout" || job.AuditTrail[0].ErrorCodes[0] != "e1" {
		t.Fatalf("clone aliases original slices: original=%#v clone=%#v", job, clone)
	}
}

func TestRandomOpaqueIDProducesValidDistinctIdentity(t *testing.T) {
	first, err := randomOpaqueID("job")
	if err != nil {
		t.Fatal(err)
	}
	second, err := randomOpaqueID("job")
	if err != nil {
		t.Fatal(err)
	}
	if first == second || !validOpaqueReference(first) || !validOpaqueReference(second) {
		t.Fatalf("invalid random identities: %q %q", first, second)
	}
}

func validSearchJobForBackendTest() SearchJob {
	start := time.Date(2026, 9, 24, 10, 0, 0, 0, time.UTC)
	return SearchJob{
		ID: "job-valid", RunID: "run-valid", State: JobRunning, Version: 2,
		CreatedAt: start, UpdatedAt: start.Add(time.Millisecond),
		Provenance: JobProvenance{
			TenantRef: "tenant-a", EnvironmentRef: "env-prod",
			TimeStart: start, TimeEnd: start.Add(time.Hour),
			Sources: []string{"source-a", "source-b"},
			CorrelationID: "corr-valid", QueryVersion: "query-v1",
		},
	}
}
