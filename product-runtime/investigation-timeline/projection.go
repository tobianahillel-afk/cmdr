// Package investigationtimeline projects authorized Case Timeline Entry snapshots.
// It owns no canonical Timeline Entry or external source data.
package investigationtimeline

import (
	"sort"
	"strings"
	"time"
	"unicode"
)

const (
	MaxEntries    = 5000
	MaxInputBytes = 8 << 20
)

type Scope struct {
	Tenant             string `json:"tenant"`
	Environment        string `json:"environment"`
	CasePermission     string `json:"case_permission"`
	TimelinePermission string `json:"timeline_permission"`
}

type Entry struct {
	ID                 string `json:"id"`
	Kind               string `json:"kind"`
	Tenant             string `json:"tenant"`
	SourcePermission   string `json:"source_permission"`
	Timezone           string `json:"timezone"`
	SourceTime         string `json:"source_time"`
	SourceRef          string `json:"source_ref"`
	Provenance         string `json:"provenance"`
	Version            int    `json:"version"`
	Supersedes         string `json:"supersedes"`
	ClockUncertaintyMS int64  `json:"clock_uncertainty_ms"`
}

type Input struct {
	Scope   Scope   `json:"scope"`
	Entries []Entry `json:"entries"`
}

type ProjectedEntry struct {
	ID            string `json:"id"`
	Kind          string `json:"kind"`
	SourceTime    string `json:"source_time"`
	Timezone      string `json:"timezone"`
	TimeUncertain bool   `json:"time_uncertain"`
	SourceRef     string `json:"source_ref"`
	Provenance    string `json:"provenance"`
	Version       int    `json:"version"`
	Supersedes    string `json:"supersedes"`
}

type Result struct {
	Status  string           `json:"status"`
	Entries []ProjectedEntry `json:"entries"`
}

type sortableEntry struct {
	entry   ProjectedEntry
	instant time.Time
	known   bool
}

// Project produces a new read-only projection from caller-resolved snapshots.
// Access decisions MUST originate from authoritative server-side checks upstream;
// this bounded pure function is not an authorization provider.
func Project(input Input) Result {
	deny := Result{Status: "deny"}
	s := input.Scope
	// The envelope budget must hold even when there are no visible entries.
	// Otherwise an oversized tenant or environment is silently accepted as
	// an empty/filtered projection, unlike the two trusted access adapters.
	if len(s.Tenant) > MaxInputBytes ||
		len(s.Environment) > MaxInputBytes-len(s.Tenant) ||
		invalidScope(s.Tenant) || invalidScope(s.Environment) ||
		s.CasePermission != "allow" || s.TimelinePermission != "allow" {
		return deny
	}
	if len(input.Entries) > MaxEntries {
		return deny
	}
	usedBytes := len(s.Tenant) + len(s.Environment)
	seen := make(map[string]bool, len(input.Entries))
	for _, e := range input.Entries {
		// Only accessible sources may influence ID and byte validation.
		// Hidden data must not turn a permitted Case timeline into a denial
		// or reveal its size through a validation failure. The total number
		// of supplied entries is independently bounded above.
		if e.Tenant != s.Tenant || e.SourcePermission != "allow" {
			continue
		}
		usedBytes += len(e.ID) + len(e.Kind) + len(e.Tenant) + len(e.SourcePermission) +
			len(e.Timezone) + len(e.SourceTime) + len(e.SourceRef) + len(e.Provenance) + len(e.Supersedes)
		if usedBytes > MaxInputBytes || invalidScope(e.ID) || seen[e.ID] {
			return deny
		}
		seen[e.ID] = true
	}
	rows := make([]sortableEntry, 0, len(input.Entries))
	filtered, partial, uncertain := false, false, false
	for _, e := range input.Entries {
		if e.Tenant != s.Tenant || e.SourcePermission == "deny" || e.SourcePermission == "unknown" {
			filtered = true
			continue
		}
		if e.SourcePermission == "missing" {
			partial = true
			continue
		}
		if e.SourcePermission != "allow" || e.Provenance == "" ||
			invalidScope(e.SourceRef) || e.Version < 1 || !validKind(e.Kind) ||
			e.ClockUncertaintyMS < 0 || (e.Supersedes != "" && e.Supersedes == e.ID) {
			return deny
		}
		instant, err := time.Parse(time.RFC3339Nano, e.SourceTime)
		// A numeric negative-zero offset says that the local offset is
		// unknown; time.Parse normalizes it to zero and loses that signal.
		// Do not silently elevate its chronology to a known local offset.
		known := err == nil && e.ClockUncertaintyMS == 0 &&
			!strings.HasSuffix(e.SourceTime, "-00:00") && e.Timezone != "-00:00"
		if known {
			// A timezone label must agree with the source timestamp offset.
			// Unknown or named zones are not silently interpreted as UTC.
			declared, zoneErr := time.Parse("Z07:00", e.Timezone)
			_, actualOffset := instant.Zone()
			_, declaredOffset := declared.Zone()
			known = zoneErr == nil && actualOffset == declaredOffset
		}
		if !known {
			uncertain = true
		}
		rows = append(rows, sortableEntry{
			entry: ProjectedEntry{ID: e.ID, Kind: e.Kind, SourceTime: e.SourceTime,
				Timezone: e.Timezone, TimeUncertain: !known, SourceRef: e.SourceRef,
				Provenance: e.Provenance, Version: e.Version, Supersedes: e.Supersedes},
			instant: instant, known: known,
		})
	}
	// A visible version must not disclose an inaccessible predecessor ID.
	visibleVersions := make(map[string]int, len(rows))
	for _, row := range rows {
		visibleVersions[row.entry.ID] = row.entry.Version
	}
	for i := range rows {
		previous := rows[i].entry.Supersedes
		if previous == "" {
			continue
		}
		version, exists := visibleVersions[previous]
		if !exists {
			rows[i].entry.Supersedes = ""
			partial = true
			continue
		}
		// Version links must be strictly descending. This also makes
		// cycles impossible without an additional traversal/graph engine.
		if version >= rows[i].entry.Version {
			return deny
		}
	}
	// Each visible ID is unique (enforced above) and is the final sort key.
	// This makes the comparison a strict total order even when two entries
	// have the same instant or unknown clocks. Stable sort is therefore
	// unnecessary; reuse the simpler standard-library in-place sorter.
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].known != rows[j].known {
			return rows[i].known
		}
		if rows[i].known && !rows[i].instant.Equal(rows[j].instant) {
			return rows[i].instant.Before(rows[j].instant)
		}
		return rows[i].entry.ID < rows[j].entry.ID
	})
	out := Result{Status: "allow", Entries: make([]ProjectedEntry, 0, len(rows))}
	for _, row := range rows {
		out.Entries = append(out.Entries, row.entry)
	}
	switch {
	case partial:
		out.Status = "partial"
	case filtered:
		out.Status = "filtered"
	case uncertain:
		out.Status = "uncertain"
	case len(rows) == 0:
		out.Status = "empty"
	}
	return out
}

// invalidScope rejects authority selectors with control or invisible
// Unicode format characters as well as empty/wildcard identifiers.
// The actual ACL remains the responsibility of the trusted provider.
func invalidScope(v string) bool {
	trimmed := strings.TrimSpace(v)
	if trimmed == "" || trimmed == "*" {
		return true
	}
	for _, r := range v {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
			return true
		}
	}
	return false
}

func validKind(kind string) bool {
	switch kind {
	case "observed", "inferred", "user-action", "decision", "response":
		return true
	}
	return false
}
