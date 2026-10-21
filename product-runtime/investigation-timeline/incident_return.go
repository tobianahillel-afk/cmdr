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

// The only return destination is Command's existing owner-authorized queue
// selection route. The Command handler must recheck the Incident on every GET.
func commandIncidentReturnURL(id string) string {
	query := url.Values{
		"view":     {"incidents"},
		"selected": {"Incident:" + id},
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
