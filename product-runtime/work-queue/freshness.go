package workqueue

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// SourceFreshness is supplied only by the trusted Command owner for a work
// item. Health and explicit tenant/module policy thresholds must come from
// their authoritative sources; the browser never selects them.
type SourceFreshness struct {
	Source string
	ObservedAt string
	Health string // available, partial, unavailable
	PolicyVersion string
	AgingAfterSeconds int64
	StaleAfterSeconds int64
}

// FreshnessProjection is a read-only, explainable source projection.
// No global stale threshold is invented when policy data is missing.
type FreshnessProjection struct {
	State string `json:"state"`
	Source string `json:"source"`
	ObservedAt string `json:"observed_at,omitempty"`
	PolicyVersion string `json:"policy_version,omitempty"`
}

var errFreshness = errors.New("invalid authoritative source freshness metadata")

// sourceFreshness runs only AFTER the Case/Task owner authorizes the item.
// A missing policy is UNKNOWN, not a guessed healthy/current result.
func sourceFreshness(value *SourceFreshness, now time.Time) (*FreshnessProjection, error) {
	if value == nil {
		return nil, nil // old clients get the same JSON field set
	}
	if len(value.Source) > 128 || invalidScope(value.Source) || !utf8.ValidString(value.Source) ||
		len(value.PolicyVersion) > 64 || (value.PolicyVersion != "" && (invalidScope(value.PolicyVersion) || !utf8.ValidString(value.PolicyVersion))) {
		return nil, errFreshness
	}
	result := &FreshnessProjection{State:"unknown", Source:value.Source, PolicyVersion:value.PolicyVersion}
	switch value.Health {
	case "unavailable":
		result.State = "source-unavailable"
	case "partial":
		result.State = "partial"
	case "available", "":
		// Source metadata is separate from the Command item's UpdatedAt.
		// Do not infer freshness from the last Incident/Task mutation.
	default:
		return nil, errFreshness
	}
	if value.ObservedAt == "" {
		return result, nil
	}
	observed, err := time.Parse(time.RFC3339Nano, value.ObservedAt)
	if err != nil || observed.IsZero() || strings.TrimSpace(value.ObservedAt) != value.ObservedAt {
		return nil, errFreshness
	}
	result.ObservedAt = observed.UTC().Format(time.RFC3339Nano)
	if result.State != "unknown" {
		return result, nil // explicit owner health takes precedence
	}
	// No valid, versioned policy: show source/time but never claim current.
	if value.PolicyVersion == "" || value.AgingAfterSeconds <= 0 ||
		value.StaleAfterSeconds <= value.AgingAfterSeconds {
		return result, nil
	}
	if observed.After(now) {
		return result, nil // future or inconsistent source clock is unknown
	}
	age := now.Sub(observed)
	if age.Seconds() >= float64(value.StaleAfterSeconds) {
		result.State = "stale"
	} else if age.Seconds() >= float64(value.AgingAfterSeconds) {
		result.State = "aging"
	} else {
		result.State = "current"
	}
	return result, nil
}

// ViewFreshness is an aggregate of ONLY the current authorized system-view
// rows after the selected status filter, including authorized later pages.
// It never includes denied/cross-tenant objects or claims a global threshold.
type ViewFreshness struct {
	Total int `json:"total"`
	Current int `json:"current"`
	Aging int `json:"aging"`
	Stale int `json:"stale"`
	Unknown int `json:"unknown"`
	Partial int `json:"partial"`
	Unavailable int `json:"source_unavailable"`
}
