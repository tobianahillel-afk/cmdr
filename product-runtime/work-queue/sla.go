package workqueue

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// AuthorizedSLA is a source-owned, independently authorized read-only
// projection. Command does not invent an SLA policy, threshold, pause rule,
// clock, or contract. A missing source means unknown, never healthy.
type AuthorizedSLA struct {
	State string `json:"state"`
	PolicySource string `json:"policy_source,omitempty"`
	PolicyVersion string `json:"policy_version,omitempty"`
	ObservedAt string `json:"observed_at,omitempty"`
	DueAt string `json:"due_at,omitempty"`
	RemainingSeconds *int64 `json:"remaining_seconds,omitempty"`
	// These details are source-authorized explanations, not pause permissions
	// or a locally implemented SLA policy/calendar engine.
	CalendarSource string `json:"calendar_source,omitempty"`
	PausedAt string `json:"paused_at,omitempty"`
	PauseReason string `json:"pause_reason,omitempty"`
}

// SLAOwner is an optional trusted host adapter distinct from the Incident /
// Task owner. The host MUST authorize the policy/engagement separately for the
// current principal, tenant, environment and selected object before returning
// a claim; inaccessible policy or engagement details must not be returned.
// Item CanRead alone does NOT grant access to SLA policy or contract metadata.
type SLAOwner interface {
	ReadAuthorizedSLA(ctx context.Context, principal, tenant, environment string, kind Kind, id string) (*AuthorizedSLA,error)
}

var errSLA = errors.New("invalid or unavailable authorized SLA projection")

func validateAuthorizedSLA(value *AuthorizedSLA) (*AuthorizedSLA,error) {
	if value == nil { return nil,nil }
	switch value.State {
	case "healthy", "approaching", "at-risk", "breached", "paused", "not-applicable", "unknown":
	default: return nil,errSLA
	}
	for _, field := range []struct{v string;max int}{
		{value.PolicySource,128},{value.PolicyVersion,96},
		{value.CalendarSource,128},{value.PauseReason,256},
	} {
		if len(field.v)>field.max || !utf8.ValidString(field.v) ||
			(field.v!="" && (invalidScope(field.v) || strings.TrimSpace(field.v)!=field.v)) {
			return nil,errSLA
		}
	}
	for _, stamp := range []string{value.ObservedAt,value.DueAt,value.PausedAt} {
		if stamp=="" { continue }
		if len(stamp)>48 || strings.TrimSpace(stamp)!=stamp {
			return nil,errSLA
		}
		parsed,err:=time.Parse(time.RFC3339Nano,stamp)
		if err!=nil || parsed.IsZero() { return nil,errSLA }
	}
	if value.RemainingSeconds != nil {
		// A remaining-time claim is authoritative source metadata: never
		// infer it from due-at without the source's calendar/pause policy.
		if value.DueAt=="" || value.ObservedAt=="" ||
			*value.RemainingSeconds < -315360000 || *value.RemainingSeconds > 315360000 {
			return nil,errSLA
		}
	}
	if value.State!="unknown" && value.State!="not-applicable" &&
		(value.PolicySource=="" || value.PolicyVersion=="" || value.ObservedAt=="") {
		return nil,errSLA
	}
	if value.State=="unknown" || value.State=="not-applicable" {
		if value.RemainingSeconds != nil { return nil,errSLA }
	}
	// The host alone owns pause authorization and source history; never
	// fabricate a reason or silently expose paused metadata for a state
	// that is not paused. Incomplete paused state is explicitly unknown.
	if value.State!="paused" && (value.PausedAt!="" || value.PauseReason!="") {
		return nil,errSLA
	}
	if value.PausedAt!="" {
		start,err:=time.Parse(time.RFC3339Nano,value.PausedAt)
		if err!=nil { return nil,errSLA }
		observed,err:=time.Parse(time.RFC3339Nano,value.ObservedAt)
		if err!=nil || start.After(observed) { return nil,errSLA }
	}
	return value,nil
}
