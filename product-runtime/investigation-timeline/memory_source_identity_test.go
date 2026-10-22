package investigationtimeline

import (
	"net/http"
	"strings"
	"testing"
)

func TestMemorySourceAliasCannotBorrowAnotherReadersPermission(t *testing.T) {
	data := memoryDemoFixture()
	// The evidence-secret canonical object is not readable by analyst-a.
	// Previously, an alias with a different SourceRef and an analyst-a
	// reader could make the CheckAccess kind/object lookup return allow
	// for the restricted original reference.
	alias := data.Sources[2]
	alias.Identity.Reference = "alias-permitted"
	alias.Readers = []string{"analyst-a"}
	data.Sources = append(data.Sources, alias)
	if _, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, data); err == nil ||
		!strings.Contains(err.Error(), "ambiguous canonical source") {
		t.Fatalf("aliased canonical object ACL was accepted: %v", err)
	}
	// A distinct canonical source in a separate Case is allowed, even if
	// its object ID and source type coincide: Case binding is part of ACL.
	data = memoryDemoFixture()
	foreign := data.Sources[2]
	foreign.Identity.Reference = "case-b-other-reference"
	foreign.Identity.Tenant = "tenant-b"
	foreign.Identity.CaseID = "case-b"
	foreign.Readers = []string{"analyst-b"}
	data.Sources = append(data.Sources, foreign)
	if _, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, data); err != nil {
		t.Fatalf("different Case scoped source rejected: %v", err)
	}
}

func TestMemorySourceCanonicalACLStillProtectsDeniedEvidence(t *testing.T) {
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(demoAuthenticator, memoryDemoFixture())
	if err != nil {
		t.Fatal(err)
	}
	regular := memoryRequest(handlers.Timeline, "case-a", "Bearer aa")
	if regular.Code != http.StatusOK ||
		strings.Contains(regular.Body.String(), "src-secret") ||
		strings.Contains(regular.Body.String(), "evidence-secret") {
		t.Fatalf("restricted canonical source leaked: %d %s", regular.Code, regular.Body.String())
	}
	authorized := memoryRequest(handlers.Timeline, "case-a", "Bearer ab")
	if authorized.Code != http.StatusOK ||
		!strings.Contains(authorized.Body.String(), "src-secret") {
		t.Fatalf("independently authorized evidence source missing: %d %s", authorized.Code, authorized.Body.String())
	}
}
