package investigationtimeline

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func memoryDemoFixture() MemoryTimelineData {
	return MemoryTimelineData{
		Principals: []MemoryTimelinePrincipal{
			{ID: "analyst-a", Tenant: "tenant-a", Environment: "prod"},
			{ID: "analyst-secret", Tenant: "tenant-a", Environment: "prod"},
			{ID: "analyst-b", Tenant: "tenant-b", Environment: "prod"},
		},
		Cases: []MemoryTimelineCase{
			{ID: "case-a", Tenant: "tenant-a", Environment: "prod", Readers: []string{"analyst-a", "analyst-secret"},
				EntryIDs: []string{"new", "old", "hidden"}},
			{ID: "case-b", Tenant: "tenant-b", Environment: "prod", Readers: []string{"analyst-b"}, EntryIDs: []string{"foreign"}},
		},
		Entries: []MemoryTimelineEntry{
			{CaseID: "case-a", Value: Entry{ID: "new", Tenant: "tenant-a", Kind: "observed",
				SourceTime: "2026-10-01T11:00:00Z", Timezone: "+00:00", Version: 2,
				SourceRef: "src-new", Provenance: "telemetry", Supersedes: "old", SourcePermission: "allow"},
				Readers: []string{"analyst-a", "analyst-secret"}},
			{CaseID: "case-a", Value: Entry{ID: "old", Tenant: "tenant-a", Kind: "inferred",
				SourceTime: "2026-10-01T10:00:00Z", Timezone: "+00:00", Version: 1,
				SourceRef: "src-old", Provenance: "evidence", SourcePermission: "allow"},
				Readers: []string{"analyst-a", "analyst-secret"}},
			{CaseID: "case-a", Value: Entry{ID: "hidden", Tenant: "tenant-a", Kind: "observed",
				SourceTime: "2026-10-01T09:00:00Z", Timezone: "+00:00", Version: 1,
				SourceRef: "src-secret", Provenance: "restricted", SourcePermission: "allow"},
				Readers: []string{"analyst-a", "analyst-secret"}},
			{CaseID: "case-b", Value: Entry{ID: "foreign", Tenant: "tenant-b", Kind: "observed",
				SourceTime: "2026-10-01T09:00:00Z", Timezone: "+00:00", Version: 1,
				SourceRef: "src-b", Provenance: "foreign"},
				Readers: []string{"analyst-b"}},
		},
		Sources: []MemoryTimelineSource{
			{Identity: SourceOwnerIdentity{Reference: "src-new", ObjectID: "event-a", Kind: "telemetry-event",
				Tenant: "tenant-a", Environment: "prod", CaseID: "case-a"},
				Readers: []string{"analyst-a", "analyst-secret"}},
			{Identity: SourceOwnerIdentity{Reference: "src-old", ObjectID: "evidence-a", Kind: "evidence",
				Tenant: "tenant-a", Environment: "prod", CaseID: "case-a"},
				Readers: []string{"analyst-a", "analyst-secret"}},
			{Identity: SourceOwnerIdentity{Reference: "src-secret", ObjectID: "evidence-secret", Kind: "evidence",
				Tenant: "tenant-a", Environment: "prod", CaseID: "case-a"},
				Readers: []string{"analyst-secret"}},
			{Identity: SourceOwnerIdentity{Reference: "src-b", ObjectID: "incident-b", Kind: "incident",
				Tenant: "tenant-b", Environment: "prod", CaseID: "case-b"},
				Readers: []string{"analyst-b"}},
		},
	}
}

func demoAuthenticator(r *http.Request) (string, error) {
	switch r.Header.Get("Authorization") {
	case "Bearer aa":
		return "analyst-a", nil
	case "Bearer ab":
		return "analyst-secret", nil
	case "Bearer bb":
		return "analyst-b", nil
	default:
		return "", errors.New("unauthenticated")
	}
}

func memoryRequest(handler http.Handler, caseID, bearer string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodGet, "/case-timeline?case_id="+caseID, nil)
	if bearer != "" {
		r.Header.Set("Authorization", bearer)
	}
	return rr(handler, r)
}
func rr(h http.Handler, r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestMemoryOwnerIntegratedCaseToTimelineHTTP(t *testing.T) {
	data := memoryDemoFixture()
	h, err := NewMemoryCaseTimelineHTTPHandler(demoAuthenticator, data)
	if err != nil {
		t.Fatal(err)
	}
	// All constructor inputs are copied: changing the caller's catalog and
	// ACLs afterwards cannot authorize new data or rename a trusted source.
	data.Cases[0].Readers[0] = "attacker"
	data.Cases[0].EntryIDs[0] = "foreign"
	data.Entries[0].Value.Provenance = "attacker"
	data.Sources[0].Readers[0] = "attacker"

	a := memoryRequest(h, "case-a", "Bearer aa")
	if a.Code != http.StatusOK || a.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("authorized Case unavailable: %d %s", a.Code, a.Body.String())
	}
	var result Result
	if err := json.Unmarshal(a.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "filtered" || len(result.Entries) != 2 ||
		result.Entries[0].ID != "old" || result.Entries[0].Kind != "inferred" ||
		result.Entries[1].ID != "new" || result.Entries[1].Kind != "observed" ||
		result.Entries[1].Supersedes != "old" || result.Entries[1].Provenance != "telemetry" {
		t.Fatalf("complete authorized Case chronology incorrect: %+v", result)
	}
	for _, hidden := range []string{"src-secret", "evidence-secret", "hidden", "attacker", "tenant-b"} {
		if strings.Contains(a.Body.String(), hidden) {
			t.Fatalf("sensitive identity leaked: %s", hidden)
		}
	}
	b := memoryRequest(h, "case-b", "Bearer bb")
	if b.Code != http.StatusOK {
		t.Fatalf("tenant B case unavailable: %d", b.Code)
	}
	if err := json.Unmarshal(b.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "allow" || len(result.Entries) != 1 || result.Entries[0].ID != "foreign" {
		t.Fatalf("tenant B wrong chronology %+v", result)
	}
	secret := memoryRequest(h, "case-a", "Bearer ab")
	if secret.Code != http.StatusOK {
		t.Fatalf("granted source access failed: %d", secret.Code)
	}
	if err := json.Unmarshal(secret.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.Status != "allow" || len(result.Entries) != 3 || result.Entries[0].ID != "hidden" {
		t.Fatalf("owner source ACL failed: %+v", result)
	}
}

func TestMemoryOwnerDeniesCrossTenantAndAnonymousAccess(t *testing.T) {
	h, err := NewMemoryCaseTimelineHTTPHandler(demoAuthenticator, memoryDemoFixture())
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		caseID, bearer string
		status         int
	}{
		{"case-a", "", http.StatusUnauthorized},
		{"case-a", "Bearer wrong", http.StatusUnauthorized},
		{"case-b", "Bearer aa", http.StatusNotFound},
		{"case-a", "Bearer bb", http.StatusNotFound},
		{"case-does-not-exist", "Bearer aa", http.StatusNotFound},
		{"case-a&tenant=tenant-b", "Bearer aa", http.StatusBadRequest},
	} {
		got := memoryRequest(h, tc.caseID, tc.bearer)
		if got.Code != tc.status || strings.Contains(got.Body.String(), "src-secret") ||
			strings.Contains(got.Body.String(), "tenant-a") {
			t.Fatalf("leaky denial: %+v => %d %s", tc, got.Code, got.Body.String())
		}
	}
}

func TestMemoryOwnerRejectsInconsistentTrustedSnapshots(t *testing.T) {
	variants := []struct {
		name   string
		change func(*MemoryTimelineData)
	}{
		{"duplicate-principal", func(d *MemoryTimelineData) { d.Principals = append(d.Principals, d.Principals[0]) }},
		{"duplicate-case", func(d *MemoryTimelineData) { d.Cases = append(d.Cases, d.Cases[0]) }},
		{"duplicate-entry", func(d *MemoryTimelineData) { d.Entries = append(d.Entries, d.Entries[0]) }},
		{"duplicate-source", func(d *MemoryTimelineData) { d.Sources = append(d.Sources, d.Sources[0]) }},
		{"unknown-reader", func(d *MemoryTimelineData) { d.Cases[0].Readers = []string{"unknown"} }},
		{"foreign-reader", func(d *MemoryTimelineData) { d.Cases[0].Readers = []string{"analyst-b"} }},
		{"entry-wrong-tenant", func(d *MemoryTimelineData) { d.Entries[0].Value.Tenant = "tenant-b" }},
		{"foreign-entry-catalog", func(d *MemoryTimelineData) { d.Cases[0].EntryIDs[0] = "foreign" }},
		{"duplicate-entry-catalog", func(d *MemoryTimelineData) { d.Cases[0].EntryIDs[0] = "old" }},
		{"missing-source", func(d *MemoryTimelineData) { d.Entries[0].Value.SourceRef = "untrusted" }},
		{"foreign-source-owner", func(d *MemoryTimelineData) { d.Sources[0].Identity.CaseID = "case-b" }},
		{"source-type-unsupported", func(d *MemoryTimelineData) { d.Sources[0].Identity.Kind = "unknown" }},
		{"invalid-entry-chronology", func(d *MemoryTimelineData) { d.Entries[0].Value.Kind = "guessed" }},
	}
	for _, tc := range variants {
		t.Run(tc.name, func(t *testing.T) {
			data := memoryDemoFixture()
			tc.change(&data)
			if _, err := NewMemoryCaseTimelineHTTPHandler(demoAuthenticator, data); err == nil {
				t.Fatal("invalid trusted owner snapshot accepted")
			}
		})
	}
	if _, err := NewMemoryCaseTimelineHTTPHandler(nil, memoryDemoFixture()); err == nil {
		t.Fatal("missing host verifier accepted")
	}
}

func TestMemoryOwnerInputImmutabilityAndConcurrentReads(t *testing.T) {
	data := memoryDemoFixture()
	h, err := NewMemoryCaseTimelineHTTPHandler(demoAuthenticator, data)
	if err != nil {
		t.Fatal(err)
	}
	initial := memoryRequest(h, "case-a", "Bearer aa")
	for i := 0; i < 12; i++ {
		got := memoryRequest(h, "case-a", "Bearer aa")
		if got.Code != initial.Code || !reflect.DeepEqual(got.Body.Bytes(), initial.Body.Bytes()) {
			t.Fatalf("read-only provider output changed: %d", i)
		}
	}
}
