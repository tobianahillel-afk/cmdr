package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// A Sysmon event and a Windows Security event may legitimately share a local
// event ID. The operator-controlled feed index provides separate, searchable
// source lineage while the original Case/source ACLs remain authoritative.
func TestLiveTwoFeedsSameEventIDAuthorizedSearchAndPreview(t *testing.T) {
	owner, err := loadDemoSnapshot(snapshotFixture(t))
	if err != nil {
		t.Fatal(err)
	}
	const caseID = "forensic-case-on-disk"
	const eventID = "common-42"
	first := []byte("{\"id\":\"common-42\",\"timestamp\":\"2026-10-15T10:20:30Z\",\"message\":\"Sysmon process start\"}\n")
	second := []byte("{\"id\":\"common-42\",\"timestamp\":\"2026-10-15T10:20:31Z\",\"message\":\"Windows authentication\"}\n")
	data, err := importLocalEventsFromFeed(context.Background(), owner, "analyst-a", caseID, first, 1)
	if err != nil {
		t.Fatal(err)
	}
	data, err = importLocalEventsFromFeed(context.Background(), data, "analyst-a", caseID, second, 2)
	if err != nil {
		t.Fatalf("second feed with same local event ID was rejected: %v", err)
	}
	token := strings.Repeat("a", 64)
	handler, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil {
		t.Fatal(err)
	}
	request := func(path string, authenticated bool) *httptest.ResponseRecorder {
		t.Helper()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		if authenticated {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		return res
	}
	for index, preview := range []string{"Sysmon process start", "Windows authentication"} {
		feed := "feed-" + string(rune('1'+index))
		digest := sha256.Sum256([]byte(caseID + "\x00" + feed + "\x00" + eventID))
		entryID := "local-jsonl-entry-" + hex.EncodeToString(digest[:])
		query := url.Values{"case_id": {caseID}, "q": {feed}, "limit": {"200"}}
		target := "/case-timeline?" + query.Encode()
		response := request(target, true)
		if response.Code != http.StatusOK {
			t.Fatalf("feed %s search failed: %d %s", feed, response.Code, response.Body.String())
		}
		var result struct {
			Entries []timeline.ProjectedEntry `json:"entries"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		if len(result.Entries) != 1 || result.Entries[0].ID != entryID ||
			result.Entries[0].Provenance != "operator-local-jsonl/"+feed {
			t.Fatalf("feed %s leaked another source or lost lineage: %+v", feed, result.Entries)
		}
		open := "/demo-source-open?" + url.Values{"case_id": {caseID}, "entry_id": {entryID}}.Encode()
		src := request(open, true)
		if src.Code != http.StatusOK || !strings.Contains(src.Body.String(), preview) {
			t.Fatalf("feed %s source preview not reachable: %d %s", feed, src.Code, src.Body.String())
		}
		if unauthorized := request(open, false); unauthorized.Code != http.StatusUnauthorized ||
			strings.Contains(unauthorized.Body.String(), preview) {
			t.Fatalf("anonymous request disclosed %s preview", feed)
		}
		if unauthorized := request(target, false); unauthorized.Code != http.StatusUnauthorized ||
			strings.Contains(unauthorized.Body.String(), entryID) {
			t.Fatalf("anonymous request disclosed %s entry", feed)
		}
	}
	foreign := request("/case-timeline?case_id=case-b&q=feed-2", true)
	if foreign.Code != http.StatusNotFound {
		t.Fatalf("cross-tenant feed search exposed another Case: %d", foreign.Code)
	}
	if _, err := importLocalEventsFromFeed(context.Background(), owner, "analyst-a", caseID, first, 5); err == nil {
		t.Fatal("untrusted feed namespace accepted")
	}
	if _, err := importLocalEventsFromFeed(context.Background(), owner, "analyst-a", caseID,
		append(first, first...), 1); err == nil {
		t.Fatal("same-feed duplicate ID accepted")
	}
}
