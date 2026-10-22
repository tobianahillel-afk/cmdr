package casequeue

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func handoffFixture() Input {
	return Input{
		TenantRef: "tenant-a", AuthorizedTenants: []string{"tenant-a"},
		EnvironmentRef: "prod", SearchQuery: "malware",
		Filters: []Filter{{Field: "status", Values: []string{"open"}}},
		Sort: []SortSpec{{Field: "updated_at", Direction: "desc"}},
		ViewDirty: true,
		SavedView: &SavedView{
			ID: "view-1", TenantRef: "tenant-a", Version: "v1", AccessDecision: "allow",
			Fields: []SavedViewField{
				{Name: "title", AccessDecision: "allow"},
				{Name: "secret", AccessDecision: "deny"},
			},
			Filters: []SavedViewFilter{{Field: "owner", Values: []string{"alice"}, AccessDecision: "allow"}},
		},
		Cases: []CaseSnapshot{
			{
				ID: "case-1", TenantRef: "tenant-a", HumanID: "CASE-1",
				Title: "Malware triage", Status: "open", Owner: "alice",
				UpdatedAt: "2026-09-29T10:00:00Z", AccessDecision: "allow", Freshness: "available",
				Incident: &IncidentContext{ID: "inc-1", TenantRef: "tenant-a", AccessDecision: "deny"},
			},
			{
				ID: "case-hidden", TenantRef: "tenant-a", HumanID: "CASE-H",
				Title: "Malware confidential", Status: "open", Owner: "alice",
				UpdatedAt: "2026-09-29T10:00:00Z", AccessDecision: "deny", Freshness: "available",
			},
		},
	}
}

func TestCaseWorkspaceHandoffRoundTripWithoutMutation(t *testing.T) {
	in := handoffFixture()
	original, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	result := OpenCaseWorkspace(in, "case-1", 432, OriginCaseQueue)
	if !result.Allowed || result.ErrorCode != ErrorNone || result.Handoff == nil {
		t.Fatalf("handoff rejected: %+v", result)
	}
	h := result.Handoff
	if h.ScreenID != CaseWorkspaceScreen || h.CaseID != "case-1" ||
		h.TenantRef != "tenant-a" || h.CaseUpdatedAt != "2026-09-29T10:00:00Z" ||
		!result.AuditRequired {
		t.Fatalf("incorrect safe envelope: %+v", h)
	}
	payload, err := base64.RawURLEncoding.DecodeString(h.ReturnToken)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(payload), "case-hidden") ||
		strings.Contains(string(payload), "secret") || strings.Contains(string(payload), "inc-1") {
		t.Fatalf("return token leaked denied metadata: %q", payload)
	}
	restored, code := RestoreCaseQueue(h.ReturnToken, in)
	if code != ErrorNone || restored.SelectedCaseID != "case-1" ||
		restored.ScrollOffset != 432 || restored.SearchQuery != in.SearchQuery ||
		restored.SavedViewID != "view-1" || restored.SavedViewVersion != "v1" ||
		restored.Origin != OriginCaseQueue || !restored.ViewDirty ||
		!reflect.DeepEqual(restored.Filters, in.Filters) ||
		!reflect.DeepEqual(restored.Sort, in.Sort) {
		t.Fatalf("inexact return context: %+v, %s", restored, code)
	}
	restored.Filters[0].Values[0] = "tampered"
	if in.Filters[0].Values[0] != "open" {
		t.Fatal("restored filter aliased caller input")
	}
	after, err := json.Marshal(in)
	if err != nil || !reflect.DeepEqual(original, after) {
		t.Fatal("handoff or restoration mutated input")
	}
	if result.Handoff.ReturnToken != OpenCaseWorkspace(in, "case-1", 432, OriginCaseQueue).Handoff.ReturnToken {
		t.Fatal("identical inputs produced nondeterministic token")
	}
}

func TestCaseWorkspaceHandoffFailClosedSelection(t *testing.T) {
	testCases := []struct {
		name string
		edit func(*Input)
		id   string
		want ErrorCode
	}{
		{"denied", func(*Input) {}, "case-hidden", ErrorCaseNotVisible},
		{"not-in-filter", func(i *Input) { i.Filters[0].Values = []string{"closed"} }, "case-1", ErrorCaseNotVisible},
		{"wrong-tenant", func(i *Input) { i.Cases[0].TenantRef = "tenant-b" }, "case-1", ErrorCrossTenantCase},
		{"stale", func(i *Input) { i.Cases[0].Freshness = "stale" }, "case-1", ErrorCaseStale},
		{"revoked-saved-view", func(i *Input) { i.SavedView.AccessDecision = "deny" }, "case-1", ErrorNone},
		{"cross-tenant-view", func(i *Input) { i.SavedView.TenantRef = "tenant-b" }, "case-1", ErrorCrossTenantSavedView},
		{"unauthorized-tenant", func(i *Input) { i.AuthorizedTenants = nil }, "case-1", ErrorCrossTenant},
		{"unknown-selection", func(*Input) {}, "other", ErrorCaseNotVisible},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in := handoffFixture()
			tc.edit(&in)
			got := OpenCaseWorkspace(in, tc.id, 10, OriginCaseQueue)
			if tc.want == ErrorNone {
				if !got.Allowed || got.Handoff == nil { t.Fatalf("expected safe fallback: %+v", got) }
				return
			}
			if got.Allowed || got.Handoff != nil || got.ErrorCode != tc.want {
				t.Fatalf("should fail closed with %s: %+v", tc.want, got)
			}
		})
	}
}

func TestCaseWorkspaceHandoffRejectsUnsafeNavigation(t *testing.T) {
	for _, tc := range []struct{ origin string; scroll int; query string }{
		{"javascript:alert(1)", 0, "malware"},
		{"https://evil.invalid", 0, "malware"},
		{OriginCaseQueue, -1, "malware"},
		{OriginCaseQueue, maxReturnScroll+1, "malware"},
		{OriginCaseQueue, 0, strings.Repeat("a", 257)},
		{OriginCaseQueue, 0, "malware\nsecret"},
	} {
		in := handoffFixture()
		in.SearchQuery = tc.query
		got := OpenCaseWorkspace(in, "case-1", tc.scroll, tc.origin)
		if got.Allowed || got.Handoff != nil || got.ErrorCode != ErrorInvalidHandoffState {
			t.Fatalf("unsafe navigation accepted: %+v", tc)
		}
	}
}

func TestCaseWorkspaceRestoreRevalidatesAuthorityAndVersion(t *testing.T) {
	in := handoffFixture()
	h := OpenCaseWorkspace(in, "case-1", 7, OriginCommandIncident).Handoff
	if h == nil { t.Fatal("fixture failed") }
	for _, tc := range []struct{ name string; edit func(*Input); want ErrorCode }{
		{"case-access-revoked", func(i *Input) { i.Cases[0].AccessDecision = "deny" }, ErrorCaseNotVisible},
		{"case-stale", func(i *Input) { i.Cases[0].Freshness = "stale" }, ErrorCaseStale},
		{"tenant-changed", func(i *Input) { i.TenantRef = "tenant-b" }, ErrorInvalidHandoffState},
		{"environment-changed", func(i *Input) { i.EnvironmentRef = "staging" }, ErrorInvalidHandoffState},
		{"view-version-changed", func(i *Input) { i.SavedView.Version = "v2" }, ErrorReturnViewChanged},
		{"view-denied", func(i *Input) { i.SavedView.AccessDecision = "deny" }, ErrorReturnViewChanged},
		{"tenant-auth-revoked", func(i *Input) { i.AuthorizedTenants = nil }, ErrorCrossTenant},
	} {
		t.Run(tc.name, func(t *testing.T) {
			latest := handoffFixture()
			tc.edit(&latest)
			restored, code := RestoreCaseQueue(h.ReturnToken, latest)
			if code != tc.want || restored.SelectedCaseID != "" {
				t.Fatalf("restoration should fail closed with %s, got %s: %+v", tc.want, code, restored)
			}
		})
	}
}

func TestCaseWorkspaceRestoreRejectsMalformedOrForgedTokens(t *testing.T) {
	in := handoffFixture()
	h := OpenCaseWorkspace(in, "case-1", 0, OriginCaseQueue).Handoff
	if h == nil { t.Fatal("fixture failed") }
	for _, token := range []string{"", "!invalid!", strings.Repeat("a", maxReturnTokenBytes*2+1),
		base64.RawURLEncoding.EncodeToString([]byte(`{"version":1,"state":{},"unknown":true}`)),
		base64.RawURLEncoding.EncodeToString([]byte(`{"version":1,"state":{}}junk`)),
	} {
		if got, code := RestoreCaseQueue(token, in); code != ErrorInvalidHandoffState || got.SelectedCaseID != "" {
			t.Fatalf("malformed token accepted: %+v %s", got, code)
		}
	}
	payload, _ := base64.RawURLEncoding.DecodeString(h.ReturnToken)
	var env returnEnvelope
	if err := json.Unmarshal(payload, &env); err != nil { t.Fatal(err) }
	for _, attack := range []struct{ name string; edit func(*returnEnvelope) }{
		{"denied-case", func(e *returnEnvelope) { e.State.SelectedCaseID = "case-hidden" }},
		{"fake-tenant", func(e *returnEnvelope) { e.State.TenantRef = "tenant-other" }},
		{"script-origin", func(e *returnEnvelope) { e.State.Origin = "javascript:alert(1)" }},
		{"invalid-filter", func(e *returnEnvelope) { e.State.Filters = []Filter{{Field:"forbidden",Values:[]string{"x"}}} }},
		{"invalid-sort", func(e *returnEnvelope) { e.State.Sort = []SortSpec{{Field:"forbidden",Direction:"desc"}} }},
		{"invalid-filter-value", func(e *returnEnvelope) { e.State.Filters = []Filter{{Field:"status",Values:[]string{}}} }},
		{"wrong-version", func(e *returnEnvelope) { e.Version = 2 }},
		{"wrong-environment", func(e *returnEnvelope) { e.State.EnvironmentRef = "other" }},
		{"missing-environment", func(e *returnEnvelope) { e.State.EnvironmentRef = "" }},
		{"missing-tenant", func(e *returnEnvelope) { e.State.TenantRef = "" }},
		{"wildcard-tenant", func(e *returnEnvelope) { e.State.TenantRef = "*" }},
		{"missing-case", func(e *returnEnvelope) { e.State.SelectedCaseID = "" }},
		{"negative-scroll", func(e *returnEnvelope) { e.State.ScrollOffset = -1 }},
		{"oversized-scroll", func(e *returnEnvelope) { e.State.ScrollOffset = maxReturnScroll + 1 }},
		{"oversized-query", func(e *returnEnvelope) { e.State.SearchQuery = strings.Repeat("x", 257) }},
		{"control-query", func(e *returnEnvelope) { e.State.SearchQuery = string([]byte{0}) }},
		{"wrong-view-reference", func(e *returnEnvelope) { e.State.SavedViewID = "view-other" }},
		{"wrong-view-version", func(e *returnEnvelope) { e.State.SavedViewVersion = "v2" }},
	} {
		t.Run(attack.name,func(t *testing.T) {
			copy := env
			copy.State.Filters = append([]Filter(nil), env.State.Filters...)
			attack.edit(&copy)
			b, _ := json.Marshal(copy)
			state, code := RestoreCaseQueue(base64.RawURLEncoding.EncodeToString(b), in)
			if code == ErrorNone || state.SelectedCaseID != "" {
				t.Fatalf("tampered payload accepted: %+v", state)
			}
		})
	}
}
