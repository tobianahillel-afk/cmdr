package eventsearch

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type JobState string

const (
	JobQueued    JobState = "queued"
	JobRunning   JobState = "running"
	JobCompleted JobState = "completed"
	JobPartial   JobState = "partial"
	JobFailed    JobState = "failed"
	JobCancelled JobState = "cancelled"
)

var (
	ErrIllegalTransition    = errors.New("illegal search job transition")
	ErrInvalidAuditContext  = errors.New("invalid search job audit context")
	ErrInvalidIdentity      = errors.New("invalid search job identity")
	ErrBackendExecution     = errors.New("search backend execution failed")
	ErrBackendResultInvalid = errors.New("search backend result rejected")
)

type AuditContext struct {
	ActorRef  string
	Rationale string
}

type AuditEvent struct {
	At            time.Time
	ActorRef      string
	TenantRef     string
	Action        string
	TargetRef     string
	Outcome       string
	Rationale     string
	CorrelationID string
	QueryVersion  string
	RunID         string
	ErrorCodes    []string
}

type JobProvenance struct {
	TenantRef      string
	EnvironmentRef string
	TimeStart      time.Time
	TimeEnd        time.Time
	Sources        []string
	CorrelationID  string
	QueryVersion   string
	ParentRunID    string
}

type ResultReference struct {
	Source string
	Ref    string
}

type SourceFailure struct {
	Source string
	Code   string
}

type SearchJob struct {
	ID            string
	RunID         string
	State         JobState
	Version       uint64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Provenance    JobProvenance
	ResultRefs    []ResultReference
	FailedSources []SourceFailure
	AuditTrail    []AuditEvent
}

type BackendRequest struct {
	TenantRef      string
	EnvironmentRef string
	TimeStart      time.Time
	TimeEnd        time.Time
	Sources        []string
	CorrelationID  string
	QueryVersion   string
	RunID          string
}

type BackendResult struct {
	TenantRef        string
	CompletedSources []string
	ResultRefs       []ResultReference
	FailedSources    []SourceFailure
}

type Backend interface {
	Execute(context.Context, BackendRequest) (BackendResult, error)
}

type Orchestrator struct {
	backend Backend
	now     func() time.Time
	newID   func(string) (string, error)
}

func NewOrchestrator(backend Backend) (*Orchestrator, error) {
	return newOrchestrator(backend, time.Now, randomOpaqueID)
}

func newOrchestrator(backend Backend, now func() time.Time, newID func(string) (string, error)) (*Orchestrator, error) {
	if backend == nil || now == nil || newID == nil {
		return nil, fmt.Errorf("orchestrator dependencies are required")
	}
	return &Orchestrator{backend: backend, now: now, newID: newID}, nil
}

func (o *Orchestrator) Create(in Input, audit AuditContext) (SearchJob, error) {
	if !validAuditContext(audit) {
		return SearchJob{}, ErrInvalidAuditContext
	}
	validated := Validate(in)
	if !validated.Allowed {
		return SearchJob{}, fmt.Errorf("event search validation rejected: %s", validated.ErrorCode)
	}

	jobID, err := o.nextIdentity("job")
	if err != nil {
		return SearchJob{}, err
	}
	runID, err := o.nextIdentity("run")
	if err != nil {
		return SearchJob{}, err
	}
	now, err := o.currentTime()
	if err != nil {
		return SearchJob{}, err
	}
	job := SearchJob{
		ID:        jobID,
		RunID:     runID,
		State:     JobQueued,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
		Provenance: JobProvenance{
			TenantRef:      validated.Envelope.TenantRef,
			EnvironmentRef: validated.Envelope.EnvironmentRef,
			TimeStart:      validated.Envelope.TimeStart,
			TimeEnd:        validated.Envelope.TimeEnd,
			Sources:        append([]string(nil), validated.Envelope.Sources...),
			CorrelationID:  validated.Envelope.CorrelationID,
			QueryVersion:   validated.Envelope.QueryVersion,
		},
	}
	job.AuditTrail = append(job.AuditTrail, job.auditEvent(now, audit, "search-job.create", "queued"))
	return job, nil
}

func (o *Orchestrator) Start(job SearchJob, audit AuditContext) (SearchJob, error) {
	job = cloneJob(job)
	if !validAuditContext(audit) {
		return job, ErrInvalidAuditContext
	}
	if err := validateJobIdentity(job); err != nil {
		return job, err
	}
	if job.State != JobQueued {
		return job, ErrIllegalTransition
	}
	now, err := o.currentTime()
	if err != nil {
		return job, err
	}
	job.State = JobRunning
	job.Version++
	job.UpdatedAt = now
	job.AuditTrail = append(job.AuditTrail, job.auditEvent(now, audit, "search-job.start", "running"))
	return job, nil
}

func (o *Orchestrator) Cancel(job SearchJob, audit AuditContext) (SearchJob, error) {
	job = cloneJob(job)
	if !validAuditContext(audit) {
		return job, ErrInvalidAuditContext
	}
	if err := validateJobIdentity(job); err != nil {
		return job, err
	}
	if job.State != JobQueued && job.State != JobRunning {
		return job, ErrIllegalTransition
	}
	now, err := o.currentTime()
	if err != nil {
		return job, err
	}
	job.State = JobCancelled
	job.Version++
	job.UpdatedAt = now
	job.AuditTrail = append(job.AuditTrail, job.auditEvent(now, audit, "search-job.cancel", "cancelled"))
	return job, nil
}

func (o *Orchestrator) Execute(ctx context.Context, job SearchJob, audit AuditContext) (SearchJob, error) {
	job = cloneJob(job)
	if ctx == nil {
		return job, fmt.Errorf("%w: nil context", ErrBackendExecution)
	}
	if !validAuditContext(audit) {
		return job, ErrInvalidAuditContext
	}
	if err := validateJobIdentity(job); err != nil {
		return job, err
	}
	if job.State == JobQueued {
		var err error
		job, err = o.Start(job, audit)
		if err != nil {
			return job, err
		}
	} else if job.State != JobRunning {
		return job, ErrIllegalTransition
	}

	if errors.Is(ctx.Err(), context.Canceled) {
		return o.cancelFromExecution(job, audit)
	}
	if ctx.Err() != nil {
		return o.failExecution(job, audit, "context-deadline"), fmt.Errorf("%w: context deadline", ErrBackendExecution)
	}

	result, err := o.backend.Execute(ctx, backendRequest(job))
	if errors.Is(err, context.Canceled) || errors.Is(ctx.Err(), context.Canceled) {
		return o.cancelFromExecution(job, audit)
	}
	if err != nil {
		return o.failExecution(job, audit, "backend-execution"), ErrBackendExecution
	}

	state, refs, failures, err := validateBackendResult(job, result)
	if err != nil {
		return o.failExecution(job, audit, "backend-result-invalid"), err
	}
	now, err := o.currentTime()
	if err != nil {
		return job, err
	}
	job.State = state
	job.ResultRefs = refs
	job.FailedSources = failures
	job.Version++
	job.UpdatedAt = now
	job.AuditTrail = append(job.AuditTrail, job.auditEvent(now, audit, "search-job.finish", string(state)))
	return job, nil
}

func (o *Orchestrator) Retry(job SearchJob, audit AuditContext) (SearchJob, error) {
	job = cloneJob(job)
	if !validAuditContext(audit) {
		return SearchJob{}, ErrInvalidAuditContext
	}
	if err := validateJobIdentity(job); err != nil {
		return SearchJob{}, err
	}
	if !terminalState(job.State) {
		return SearchJob{}, ErrIllegalTransition
	}

	jobID, err := o.nextIdentity("job")
	if err != nil {
		return SearchJob{}, err
	}
	runID, err := o.nextIdentity("run")
	if err != nil {
		return SearchJob{}, err
	}
	if jobID == job.ID || runID == job.RunID {
		return SearchJob{}, fmt.Errorf("%w: retry must use new job and run identities", ErrInvalidIdentity)
	}
	now, err := o.currentTime()
	if err != nil {
		return SearchJob{}, err
	}
	retry := SearchJob{
		ID:        jobID,
		RunID:     runID,
		State:     JobQueued,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
		Provenance: JobProvenance{
			TenantRef:      job.Provenance.TenantRef,
			EnvironmentRef: job.Provenance.EnvironmentRef,
			TimeStart:      job.Provenance.TimeStart,
			TimeEnd:        job.Provenance.TimeEnd,
			Sources:        append([]string(nil), job.Provenance.Sources...),
			CorrelationID:  job.Provenance.CorrelationID,
			QueryVersion:   job.Provenance.QueryVersion,
			ParentRunID:    job.RunID,
		},
	}
	retry.AuditTrail = append(retry.AuditTrail, retry.auditEvent(now, audit, "search-job.retry", "queued"))
	return retry, nil
}

func (o *Orchestrator) cancelFromExecution(job SearchJob, audit AuditContext) (SearchJob, error) {
	now, err := o.currentTime()
	if err != nil {
		return job, err
	}
	job.State = JobCancelled
	job.ResultRefs = nil
	job.FailedSources = nil
	job.Version++
	job.UpdatedAt = now
	job.AuditTrail = append(job.AuditTrail, job.auditEvent(now, audit, "search-job.cancel", "cancelled"))
	return job, nil
}

func (o *Orchestrator) failExecution(job SearchJob, audit AuditContext, code string) SearchJob {
	now, err := o.currentTime()
	if err != nil {
		now = job.UpdatedAt
	}
	job.State = JobFailed
	job.ResultRefs = nil
	job.FailedSources = nil
	job.Version++
	job.UpdatedAt = now
	event := job.auditEvent(now, audit, "search-job.finish", "failed")
	event.ErrorCodes = []string{code}
	job.AuditTrail = append(job.AuditTrail, event)
	return job
}

func (o *Orchestrator) nextIdentity(prefix string) (string, error) {
	value, err := o.newID(prefix)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidIdentity, err)
	}
	if !validOpaqueReference(value) {
		return "", ErrInvalidIdentity
	}
	return value, nil
}

func (o *Orchestrator) currentTime() (time.Time, error) {
	value := o.now().UTC()
	if value.IsZero() {
		return time.Time{}, fmt.Errorf("orchestrator clock returned zero time")
	}
	return value, nil
}

func (job SearchJob) auditEvent(at time.Time, audit AuditContext, action, outcome string) AuditEvent {
	return AuditEvent{
		At: at, ActorRef: audit.ActorRef, TenantRef: job.Provenance.TenantRef,
		Action: action, TargetRef: job.ID, Outcome: outcome, Rationale: audit.Rationale,
		CorrelationID: job.Provenance.CorrelationID, QueryVersion: job.Provenance.QueryVersion,
		RunID: job.RunID,
	}
}

func backendRequest(job SearchJob) BackendRequest {
	return BackendRequest{
		TenantRef: job.Provenance.TenantRef, EnvironmentRef: job.Provenance.EnvironmentRef,
		TimeStart: job.Provenance.TimeStart, TimeEnd: job.Provenance.TimeEnd,
		Sources: append([]string(nil), job.Provenance.Sources...),
		CorrelationID: job.Provenance.CorrelationID, QueryVersion: job.Provenance.QueryVersion,
		RunID: job.RunID,
	}
}

func validateBackendResult(job SearchJob, result BackendResult) (JobState, []ResultReference, []SourceFailure, error) {
	reject := func(reason string) (JobState, []ResultReference, []SourceFailure, error) {
		return JobFailed, nil, nil, fmt.Errorf("%w: %s", ErrBackendResultInvalid, reason)
	}
	if result.TenantRef != job.Provenance.TenantRef {
		return reject("tenant mismatch")
	}

	requested := stringSet(job.Provenance.Sources)
	completed := map[string]bool{}
	for _, source := range result.CompletedSources {
		if !requested[source] || completed[source] {
			return reject("invalid completed source")
		}
		completed[source] = true
	}

	failed := map[string]bool{}
	failures := append([]SourceFailure(nil), result.FailedSources...)
	for _, failure := range failures {
		if !requested[failure.Source] || completed[failure.Source] || failed[failure.Source] ||
			!validOpaqueReference(failure.Code) {
			return reject("invalid failed source")
		}
		failed[failure.Source] = true
	}
	if len(completed)+len(failed) != len(requested) {
		return reject("backend did not account for every requested source")
	}

	refs := append([]ResultReference(nil), result.ResultRefs...)
	seenRefs := map[string]bool{}
	for _, ref := range refs {
		if !completed[ref.Source] || !validOpaqueReference(ref.Ref) {
			return reject("invalid result reference")
		}
		key := ref.Source + "\x00" + ref.Ref
		if seenRefs[key] {
			return reject("duplicate result reference")
		}
		seenRefs[key] = true
	}

	sort.Slice(refs, func(i, j int) bool {
		if refs[i].Source != refs[j].Source {
			return refs[i].Source < refs[j].Source
		}
		return refs[i].Ref < refs[j].Ref
	})
	sort.Slice(failures, func(i, j int) bool {
		if failures[i].Source != failures[j].Source {
			return failures[i].Source < failures[j].Source
		}
		return failures[i].Code < failures[j].Code
	})

	switch {
	case len(failed) == 0:
		return JobCompleted, refs, failures, nil
	case len(completed) == 0:
		return JobFailed, refs, failures, nil
	default:
		return JobPartial, refs, failures, nil
	}
}

func validateJobIdentity(job SearchJob) error {
	if !validOpaqueReference(job.ID) || !validOpaqueReference(job.RunID) ||
		!validOpaqueReference(job.Provenance.CorrelationID) || !validOpaqueReference(job.Provenance.QueryVersion) ||
		strings.TrimSpace(job.Provenance.TenantRef) == "" || strings.TrimSpace(job.Provenance.EnvironmentRef) == "" ||
		len(job.Provenance.Sources) == 0 || !job.Provenance.TimeStart.Before(job.Provenance.TimeEnd) ||
		job.Version == 0 || job.CreatedAt.IsZero() || job.UpdatedAt.IsZero() {
		return ErrInvalidIdentity
	}
	if job.Provenance.ParentRunID != "" && !validOpaqueReference(job.Provenance.ParentRunID) {
		return ErrInvalidIdentity
	}
	return nil
}

func validAuditContext(audit AuditContext) bool {
	if !validOpaqueReference(audit.ActorRef) || !utf8.ValidString(audit.Rationale) ||
		strings.TrimSpace(audit.Rationale) != audit.Rationale || audit.Rationale == "" ||
		len(audit.Rationale) > 1024 {
		return false
	}
	for _, r := range audit.Rationale {
		if unicode.IsControl(r) {
			return false
		}
	}
	return true
}

func terminalState(state JobState) bool {
	switch state {
	case JobCompleted, JobPartial, JobFailed, JobCancelled:
		return true
	default:
		return false
	}
}

func cloneJob(job SearchJob) SearchJob {
	job.Provenance.Sources = append([]string(nil), job.Provenance.Sources...)
	job.ResultRefs = append([]ResultReference(nil), job.ResultRefs...)
	job.FailedSources = append([]SourceFailure(nil), job.FailedSources...)
	job.AuditTrail = append([]AuditEvent(nil), job.AuditTrail...)
	for i := range job.AuditTrail {
		job.AuditTrail[i].ErrorCodes = append([]string(nil), job.AuditTrail[i].ErrorCodes...)
	}
	return job
}

func randomOpaqueID(prefix string) (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	return prefix + "-" + hex.EncodeToString(raw[:]), nil
}
