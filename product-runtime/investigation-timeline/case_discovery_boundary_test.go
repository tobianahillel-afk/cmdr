package investigationtimeline

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type caseDiscoveryCheckerFunc func(context.Context, AccessRequest) (string, error)

func (f caseDiscoveryCheckerFunc) CheckAccess(ctx context.Context, req AccessRequest) (string, error) {
	return f(ctx, req)
}

func TestCaseDiscoveryAuthorityFailuresNeverEmitPartialList(t *testing.T) {
	for _, tc := range []struct {
		name, decision string
		failure        bool
		want           int
	}{
		{"explicit-denial", "deny", false, http.StatusOK},
		{"missing-access", "missing", false, http.StatusOK},
		{"unknown-access", "unknown", false, http.StatusOK},
		{"invalid-decision", "invalid", false, http.StatusNotFound},
		{"owner-error", "", true, http.StatusNotFound},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := caseDiscoveryCheckerFunc(func(ctx context.Context, req AccessRequest) (string, error) {
				if req.Resource != "case" || req.CaseID != req.ObjectID ||
					req.Permission != caseReadPermission || req.Tenant != "tenant-a" {
					t.Fatalf("Case candidate without bound owner permission: %+v", req)
				}
				if req.CaseID == "case-one" {
					return "allow", nil
				}
				if tc.failure {
					return "", errors.New("restricted owner failure")
				}
				return tc.decision, nil
			})
			h, err := NewCaseDiscoveryHTTPHandler(timelineHTTPSession{}, store,
				caseDiscoveryCatalogFunc(func(context.Context, Scope) ([]string, error) {
					return []string{"case-one", "case-two"}, nil
				}))
			if err != nil {
				t.Fatal(err)
			}
			result := caseDiscoveryGet(h, "Bearer analyst-a")
			if result.Code != tc.want {
				t.Fatalf("unexpected Case discovery status %d: %s", result.Code, result.Body.String())
			}
			if result.Code == http.StatusOK {
				var list struct {
					Cases []string `json:"cases"`
				}
				if err := json.Unmarshal(result.Body.Bytes(), &list); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(list.Cases, []string{"case-one"}) {
					t.Fatalf("restricted Case leaked: %+v", list)
				}
			} else if strings.Contains(result.Body.String(), "case-one") ||
				strings.Contains(result.Body.String(), "case-two") ||
				strings.Contains(result.Body.String(), "owner") {
				t.Fatalf("partial Case list disclosed after failure: %s", result.Body.String())
			}
		})
	}
}

func TestCaseDiscoveryRejectsInvalidSessionAndUnboundedCatalog(t *testing.T) {
	store, _ := workspaceFixture()
	calls := 0
	catalog := caseDiscoveryCatalogFunc(func(context.Context, Scope) ([]string, error) {
		calls++
		return []string{"case-one"}, nil
	})
	for _, mode := range []string{"empty-tenant", "wildcard-tenant", "missing-environment", "nil-context", "cancelled-context"} {
		h, err := NewCaseDiscoveryHTTPHandler(invalidTimelineHTTPSession{mode}, store, catalog)
		if err != nil {
			t.Fatal(err)
		}
		result := caseDiscoveryGet(h, "Bearer analyst-a")
		if result.Code != http.StatusUnauthorized || calls != 0 {
			t.Fatalf("invalid session %q enumerated Cases: HTTP %d, calls=%d", mode, result.Code, calls)
		}
	}
	for _, values := range [][]string{
		make([]string, MaxEntries+1),
		{strings.Repeat("x", maxHTTPCaseIDBytes+1)},
		{string([]byte{0xff})},
	} {
		h, err := NewCaseDiscoveryHTTPHandler(timelineHTTPSession{}, store,
			caseDiscoveryCatalogFunc(func(context.Context, Scope) ([]string, error) {
				return values, nil
			}))
		if err != nil {
			t.Fatal(err)
		}
		result := caseDiscoveryGet(h, "Bearer analyst-a")
		if result.Code != http.StatusNotFound || strings.Contains(result.Body.String(), "case-one") {
			t.Fatalf("unbounded or invalid Case catalog passed: HTTP %d", result.Code)
		}
	}
	if calls != 0 {
		t.Fatalf("invalid sources were enumerated: %d", calls)
	}
	if _, err := NewCaseDiscoveryHTTPHandler(timelineHTTPSession{}, store, nil); err == nil {
		t.Fatal("nil candidate-owner authority accepted")
	}
}

func TestCaseDiscoveryRevocationDuringEnumerationFailsClosed(t *testing.T) {
	catalog := caseDiscoveryCatalogFunc(func(context.Context, Scope) ([]string, error) {
		return []string{"case-one", "case-two"}, nil
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	checker := caseDiscoveryCheckerFunc(func(_ context.Context, req AccessRequest) (string, error) {
		if req.CaseID == "case-one" {
			cancel()
			return "allow", nil
		}
		return "deny", nil
	})
	h, err := NewCaseDiscoveryHTTPHandler(timelineHTTPSession{}, checker, catalog)
	if err != nil {
		t.Fatal(err)
	}
	// Verify a cancelled request cannot serialize a previously authorized Case.
	req := httptest.NewRequest(http.MethodGet, "/cases", nil).WithContext(ctx)
	req.Header.Set("Authorization", "Bearer analyst-a")
	result := rr(h, req)
	if result.Code != http.StatusServiceUnavailable ||
		strings.Contains(result.Body.String(), "case-one") {
		t.Fatalf("revocation leaked partial Case catalogue: %d %s", result.Code, result.Body.String())
	}
}
