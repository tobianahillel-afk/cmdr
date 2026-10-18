package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// This is a real imported, on-disk Wazuh + generic Case. It exceeds the
// search candidate bound in just one family, without exceeding the supported
// Case size. A generic-only hunt must not be denied by 501 unrelated Wazuh
// alerts, while an "all" or Wazuh hunt still fails closed at 500 candidates.
func TestMixedCaseFamilySearchBudgetPreservesCrossFamilyAuthorization(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	var wazuhJSONL strings.Builder
	for i := 0; i < maxWazuhPreviewCandidates+1; i++ {
		fmt.Fprintf(&wazuhJSONL,
			"{\"id\":\"w-%04d\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"agent\":{\"id\":\"001\"},\"rule\":{\"id\":\"5710\",\"description\":\"OTHER-FAMILY-ALERT\",\"level\":7}}\n", i)
	}
	wazuh := writeJSONLEvents(t, wazuhJSONL.String())
	generic := writeJSONLEvents(t,
		"{\"id\":\"generic-single\",\"timestamp\":\"2026-10-17T10:00:10Z\",\"message\":\"AUTHORIZED-GENERIC-ONLY\"}\n")
	snapshot := filepath.Join(t.TempDir(), "mixed-case.json")
	if err := writePersistedMixedCase(context.Background(), "analyst-a",
		snapshotFixture(t), caseID, []string{wazuh}, []string{generic},
		snapshot, &bytes.Buffer{}); err != nil {
		t.Fatalf("real mixed-Case fixture import failed: %v", err)
	}
	data, err := loadDemoSnapshot(snapshot)
	if err != nil { t.Fatal(err) }
	imported, err := importedNormalizedIDs(data, caseID)
	if err != nil || len(imported) != maxWazuhPreviewCandidates+2 {
		t.Fatalf("expected 501 Wazuh and 1 generic authorized event: %d %v", len(imported), err)
	}
	token, err := randomBearer()
	if err != nil { t.Fatal(err) }
	owner, err := newDemoHandlerWithData(token, "analyst-a", data)
	if err != nil { t.Fatal(err) }
	serve := func(handler http.Handler, family, value string) *httptest.ResponseRecorder {
		body, err := json.Marshal(map[string]string{
			"case_id": caseID, "field": "preview", "source_family": family,
			"value": value, "from": "2026-10-17T09:59:00Z",
			"to": "2026-10-17T10:02:00Z",
		})
		if err != nil { t.Fatal(err) }
		req := httptest.NewRequest(http.MethodPost,
			"/demo-local-search", bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		res := httptest.NewRecorder()
		handler.ServeHTTP(res, req)
		return res
	}
	handler := newDemoLocalJSONLSearchHandler(token, "analyst-a", caseID, data, imported, owner)
	// One generic candidate, not 502 Wazuh+generic candidates.
	genericReceipt := serve(handler, "generic", "AUTHORIZED-GENERIC-ONLY")
	var decoded struct { Job search.SearchJob `json:"job"` }
	if genericReceipt.Code != http.StatusOK ||
		json.Unmarshal(genericReceipt.Body.Bytes(), &decoded) != nil ||
		decoded.Job.State != search.JobCompleted ||
		len(decoded.Job.ResultRefs) != 1 {
		t.Fatalf("valid authorized generic hunt denied by unrelated Wazuh load: %d %s",
			genericReceipt.Code, genericReceipt.Body.String())
	}
	for _,s := range []string{"AUTHORIZED-GENERIC-ONLY","OTHER-FAMILY-ALERT",token} {
		if strings.Contains(genericReceipt.Body.String(),s) {
			t.Fatalf("family result disclosed source preview/bearer: %s",s)
		}
	}
	// The original hard ceiling remains effective for selected family
	// searches; a caller cannot increase the 500 cap by choosing a family.
	for _,family := range []string{"wazuh","all"} {
		result := serve(handler, family, "OTHER-FAMILY-ALERT")
		if result.Code == http.StatusOK ||
			strings.Contains(result.Body.String(),"local-jsonl-entry-") {
			t.Fatalf("family %s bypassed candidate ceiling or disclosed partial refs: %d %s",
				family, result.Code, result.Body.String())
		}
	}
	// A revoked source from the otherwise unselected family is NOT
	// silently ignored: it still passes through first-party source-open.
	restricted := data
	restricted.Principals = append([]timeline.MemoryTimelinePrincipal(nil),data.Principals...)
	for _,p := range data.Principals {
		if p.ID == "analyst-a" {
			restricted.Principals=append(restricted.Principals,timeline.MemoryTimelinePrincipal{
				ID:"other-source-reader",Tenant:p.Tenant,Environment:p.Environment,
			})
			break
		}
	}
	restricted.Sources = append([]timeline.MemoryTimelineSource(nil),data.Sources...)
	revoked := false
	for i := range restricted.Sources {
		if restricted.Sources[i].Wazuh != nil {
			restricted.Sources[i].Readers=[]string{"other-source-reader"}
			revoked=true
			break
		}
	}
	if !revoked { t.Fatal("missing Wazuh source to revoke") }
	restrictedOwner, err := newDemoHandlerWithData(token,"analyst-a",restricted)
	if err != nil { t.Fatalf("valid owner with revoked typed source rejected: %v",err) }
	restrictedHandler := newDemoLocalJSONLSearchHandler(token,"analyst-a",caseID,restricted,imported,restrictedOwner)
	denied := serve(restrictedHandler, "generic", "AUTHORIZED-GENERIC-ONLY")
	if denied.Code == http.StatusOK ||
		strings.Contains(denied.Body.String(),"local-jsonl-entry-") ||
		strings.Contains(denied.Body.String(),"AUTHORIZED-GENERIC-ONLY") {
		t.Fatalf("revoked unrelated-family source weakened whole Case authorization: %d %s",
			denied.Code,denied.Body.String())
	}
}
