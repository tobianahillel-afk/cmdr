package workqueue

import (
	"net/http"
	"strings"
	"testing"
)

func TestIncidentNextActionShownOnlyFromTrustedCommandOverview(t *testing.T) {
	owner := &overviewFixtureOwner{
		fakeOwner: deepReturnOwner(),
		overview: trustedIncidentOverview(),
	}
	owner.overview.NextAction = "<script>exfiltrate()</script>"
	owner.overview.NextActionDueAt = "2026-10-21T04:00:00+02:00"
	detail := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), owner),
		"/command/incident-detail?incident=inc-100")
	content := detail.Body.String()
	if detail.Code != http.StatusOK || owner.calls != 1 ||
		!strings.Contains(content, "Prochaine action Command") ||
		!strings.Contains(content, "&lt;script&gt;") ||
		!strings.Contains(content, "2026-10-21T04:00:00+02:00") ||
		strings.Contains(content, "<script>") {
		t.Fatalf("authorized next action not safely presented: %d %s", detail.Code, content)
	}
	owner.overview.NextAction = ""
	owner.overview.NextActionDueAt = ""
	partial := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), owner),
		"/command/incident-detail?incident=inc-100")
	if partial.Code != http.StatusOK ||
		!strings.Contains(partial.Body.String(), "Prochaine action Command : inconnue") ||
		strings.Contains(partial.Body.String(), "2026-10-21T04:00:00+02:00") {
		t.Fatalf("missing Command action fabricated: %d %s", partial.Code, partial.Body.String())
	}
	owner.allowed["Incident:inc-100"] = false
	denied := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), owner),
		"/command/incident-detail?incident=inc-100")
	if denied.Code != http.StatusNotFound || owner.calls != 2 ||
		strings.Contains(denied.Body.String(), "Prochaine action") {
		t.Fatalf("revoked Incident released owner action: %d %s", denied.Code, denied.Body.String())
	}
}

func TestIncidentNextActionRejectsMalformedOrCrossScopeTrustedFields(t *testing.T) {
	for _, tc := range []struct {
		name string
		change func(*AuthorizedIncidentOverview)
	}{
		{"oversized", func(o *AuthorizedIncidentOverview) { o.NextAction = strings.Repeat("x", 257) }},
		{"control-byte", func(o *AuthorizedIncidentOverview) { o.NextAction = "run\nunsafe" }},
		{"invalid-utf8", func(o *AuthorizedIncidentOverview) { o.NextAction = string([]byte{0xff}) }},
		{"due-without-action", func(o *AuthorizedIncidentOverview) { o.NextAction = ""; o.NextActionDueAt = "2026-10-21T04:00:00Z" }},
		{"due-malformed", func(o *AuthorizedIncidentOverview) { o.NextActionDueAt = "tomorrow" }},
		{"due-zero", func(o *AuthorizedIncidentOverview) { o.NextActionDueAt = "0001-01-01T00:00:00Z" }},
		{"due-control", func(o *AuthorizedIncidentOverview) { o.NextActionDueAt = "2026-10-21T04:00:00Z\nleak" }},
		{"due-overlong", func(o *AuthorizedIncidentOverview) { o.NextActionDueAt = strings.Repeat("x", 65) }},
		{"wrong-tenant", func(o *AuthorizedIncidentOverview) { o.Tenant = "foreign" }},
		{"wrong-environment", func(o *AuthorizedIncidentOverview) { o.Environment = "foreign" }},
		{"wrong-incident", func(o *AuthorizedIncidentOverview) { o.IncidentID = "hidden" }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			owner := &overviewFixtureOwner{
				fakeOwner: deepReturnOwner(), overview: trustedIncidentOverview(),
			}
			owner.overview.NextAction = "Validate alert"
			owner.overview.NextActionDueAt = "2026-10-21T04:00:00Z"
			tc.change(owner.overview)
			resp := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), owner),
				"/command/incident-detail?incident=inc-100")
			if resp.Code != http.StatusServiceUnavailable ||
				strings.Contains(resp.Body.String(), "Validate alert") ||
				strings.Contains(resp.Body.String(), "Incident Detail — aperçu") ||
				strings.Contains(resp.Body.String(), "hidden") {
				t.Fatalf("invalid Command next action leaked partial HTML: %d %s", resp.Code, resp.Body.String())
			}
		})
	}
}

func TestIncidentNextActionLegacySourceAndUnknownDueRemainPartial(t *testing.T) {
	legacy := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), deepReturnOwner()),
		"/command/incident-detail?incident=inc-100")
	if legacy.Code != http.StatusOK ||
		!strings.Contains(legacy.Body.String(), "prochaine action") ||
		strings.Contains(legacy.Body.String(), "Validate alert") {
		t.Fatalf("older owner requires unknown action, not fake value: %d %s", legacy.Code, legacy.Body.String())
	}
	owner := &overviewFixtureOwner{fakeOwner: deepReturnOwner(), overview: trustedIncidentOverview()}
	owner.overview.NextAction = "Verify alert"
	owner.overview.NextActionDueAt = ""
	resp := detailGET(NewIncidentDetailHTTPHandler(trustedHost(), owner),
		"/command/incident-detail?incident=inc-100")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "échéance : non précisée") {
		t.Fatalf("unknown due time was fabricated: %d %s", resp.Code, resp.Body.String())
	}
}
