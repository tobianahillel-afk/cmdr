package investigationtimeline

import (
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"
)

// A Command Incident ID is only a navigation hint; it never changes Case,
// Timeline Entry, typed-source or destination Incident authorization. This is
// deliberately independent of Command's object registry and permissions.
func validIncidentReturnID(id string) bool {
	return len(id) > 0 && len(id) <= 128 && utf8.ValidString(id) &&
		strings.TrimSpace(id) == id && !invalidScope(id) &&
		strings.IndexFunc(id, unicode.IsControl) == -1
}

// Only fixed Command routes are possible. The returned Incident is reauthorized
// at the Command destination; neither hint is a permission or a relation proof.
func commandIncidentReturnURL(id string, detail ...bool) string {
	if len(detail) != 0 && detail[0] {
		return "/command/incident-detail?" + (url.Values{"incident": {id}}).Encode()
	}
	query := url.Values{
		"view":     {"incidents"},
		"selected": {"Incident:" + id},
		"restore":  {"1"},
	}
	return "/command/work-queue?" + query.Encode()
}

// Consume a single optional HTML-only origin hint before strict Timeline
// query parsing. Never use it to resolve Case or source permissions.
func extractIncidentOrigin(values url.Values) (string, bool) {
	origins, exists := values["origin_incident"]
	if !exists {
		return "", true
	}
	if len(origins) != 1 || !validIncidentReturnID(origins[0]) {
		return "", false
	}
	delete(values, "origin_incident")
	return origins[0], true
}

// extractIncidentReturn accepts only the one optional "detail" destination
// when an Incident origin is also present. Old Work Queue return links remain
// supported. Both keys are removed before strict Case selector parsing.
func extractIncidentReturn(values url.Values) (string, bool, bool) {
	origin, valid := extractIncidentOrigin(values)
	if !valid {
		return "", false, false
	}
	view, present := values["origin_view"]
	if !present {
		return origin, false, true
	}
	if origin == "" || len(view) != 1 || view[0] != "detail" {
		return "", false, false
	}
	delete(values, "origin_view")
	return origin, true, true
}
