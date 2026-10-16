package investigationtimeline

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

// timelineReadSelection affects presentation only, never provider authority.
type timelineReadSelection struct {
	kind           string
	from, to       time.Time
	hasFrom, hasTo bool
	paged          bool
	limit, offset  int
}

func parseTimelineReadSelection(values url.Values) (timelineReadSelection, bool) {
	var s timelineReadSelection
	for key, entries := range values {
		if len(entries) != 1 {
			return timelineReadSelection{}, false
		}
		switch key {
		case "case_id":
		case "limit":
			value, err := strconv.Atoi(entries[0])
			if err != nil || value < 1 || value > maxTimelinePageSize ||
				strconv.Itoa(value) != entries[0] {
				return timelineReadSelection{}, false
			}
			s.paged, s.limit = true, value
		case "offset":
			value, err := strconv.Atoi(entries[0])
			if err != nil || value < 0 || value > MaxEntries ||
				strconv.Itoa(value) != entries[0] {
				return timelineReadSelection{}, false
			}
			s.offset = value
		case "kind":
			if !validKind(entries[0]) {
				return timelineReadSelection{}, false
			}
			s.kind = entries[0]
		case "from", "to":
			v := entries[0]
			if len(v) < 20 || len(v) > 35 || strings.HasSuffix(v, "-00:00") {
				return timelineReadSelection{}, false
			}
			parsed, err := time.Parse(time.RFC3339Nano, v)
			if err != nil {
				return timelineReadSelection{}, false
			}
			if key == "from" {
				s.from, s.hasFrom = parsed, true
			} else {
				s.to, s.hasTo = parsed, true
			}
		default:
			return timelineReadSelection{}, false
		}
	}
	if _, specified := values["offset"]; specified && !s.paged {
		return timelineReadSelection{}, false
	}
	if s.hasFrom && s.hasTo && !s.from.Before(s.to) {
		return timelineReadSelection{}, false
	}
	return s, true
}

// Select an authorized Case Timeline within [from,to).
// Unknown chronology remains visible and explicitly marked uncertain.
func selectAuthorizedTimeline(in Result, s timelineReadSelection) Result {
	if in.Status == "deny" || s.kind == "" && !s.hasFrom && !s.hasTo {
		return in
	}
	out := Result{Status: in.Status, Entries: make([]ProjectedEntry, 0, len(in.Entries))}
	unknown := false
	for _, entry := range in.Entries {
		if s.kind != "" && entry.Kind != s.kind {
			continue
		}
		if s.hasFrom || s.hasTo {
			instant, err := time.Parse(time.RFC3339Nano, entry.SourceTime)
			if entry.TimeUncertain || err != nil {
				entry.TimeUncertain = true
				unknown = true
			} else if s.hasFrom && instant.Before(s.from) || s.hasTo && !instant.Before(s.to) {
				continue
			}
		}
		out.Entries = append(out.Entries, entry)
	}
	if unknown && out.Status == "allow" {
		out.Status = "partial"
	}
	if len(out.Entries) == 0 && (out.Status == "allow" || out.Status == "uncertain") {
		out.Status = "empty"
	}
	return out
}
