package investigationtimeline

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// This is a trusted server-side integration stub, not a request-provided ACL.
// It exposes the ordering of independent Case authorization decisions.
type discoveryCheckerFunc func(context.Context, AccessRequest) (string, error)

func (f discoveryCheckerFunc) CheckAccess(ctx context.Context, req AccessRequest) (string, error) {
	return f(ctx, req)
}

// The real read-only Case picker must never publish an earlier allowed
// Case when an *authority* fails on a later candidate. Ordinary independent
// denials may filter a list; operational ACL failure or an unknown decision
// must fail the whole response, not return an incomplete success receipt.
func TestCaseDiscoveryDoesNotLeakPartialCatalogAfterLateACLFailure(t *testing.T) {
	for _, mode := range []string{"provider-failure", "unknown-decision"} {
		t.Run(mode, func(t *testing.T) {
			seen := []string{}
			catalog := caseDiscoveryCatalogFunc(func(ctx context.Context, scope Scope) ([]string, error) {
				if scope.Tenant != "tenant-a" || scope.Environment != "prod" {
					t.Fatalf("server scope was not derived from authenticated session: %+v", scope)
				}
				return []string{"case-one", "case-two"}, nil
			})
			handler, err := NewCaseDiscoveryHTTPHandler(timelineHTTPSession{},
				discoveryCheckerFunc(func(ctx context.Context, req AccessRequest) (string, error) {
					if req.Resource != "case" || req.Permission != caseReadPermission ||
						req.CaseID != req.ObjectID || req.Tenant != "tenant-a" ||
						req.Environment != "prod" {
						t.Fatalf("Case authorization used untrusted selector or scope: %+v", req)
					}
					seen = append(seen, req.CaseID)
					if req.CaseID == "case-one" {
						return "allow", nil
					}
					switch mode {
					case "provider-failure":
						return "", errors.New("SENSITIVE-ACL-INTERNAL-FAILURE")
					case "unknown-decision":
						return "unrecognized-policy-outcome", nil
					}
				}), catalog)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(http.MethodGet, "/cases", nil)
			req.Header.Set("Authorization", "Bearer analyst-a")
			resp := httptest.NewRecorder()
			handler.ServeHTTP(resp, req)
			if resp.Code != http.StatusNotFound ||
				resp.Body.String() != "{\"error\":\"request-unavailable\"}\n" ||
				len(seen) != 2 || seen[0] != "case-one" || seen[1] != "case-two" {
				t.Fatalf("late %s did not fail atomically: HTTP %d body=%q checks=%v",
					mode, resp.Code, resp.Body.String(), seen)
			}
			for _, private := range []string{
				"case-one", "case-two", "SENSITIVE-ACL-INTERNAL-FAILURE",
			} {
				if strings.Contains(resp.Body.String(), private) {
					t.Fatalf("partial Case identity or internal provider failure leaked: %q", private)
				}
			}
		})
	}
}

// A session-derived cancellation during the second Case authorization is a
// realistic race when a browser closes the picker. No partial Case list may
// be returned, even though an earlier Case had already been authorized.
func TestCaseDiscoveryRejectsLateCancelledSessionWithoutPartialResponse(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	checks := 0
	handler, err := NewCaseDiscoveryHTTPHandler(
		cancelledCaseDiscoverySession{ctx},
		discoveryCheckerFunc(func(_ context.Context, req AccessRequest) (string, error) {
			checks++
			if checks == 2 {
				cancel()
			}
			return "allow", nil
		}),
		caseDiscoveryCatalogFunc(func(context.Context, Scope) ([]string, error) {
			return []string{"case-one", "case-two"}, nil
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	resp := caseDiscoveryGet(handler, "Bearer analyst-a")
	if resp.Code != http.StatusNotFound || checks != 2 ||
		resp.Body.String() != "{\"error\":\"request-unavailable\"}\n" {
		t.Fatalf("revocation after first Case leaked a partial list: HTTP %d checks %d %q",
			resp.Code, checks, resp.Body.String())
	}
}

type cancelledCaseDiscoverySession struct{ ctx context.Context }

func (s cancelledCaseDiscoverySession) ResolveCaseTimelineSession(r *http.Request) (context.Context, Scope, error) {
	if r.Header.Get("Authorization") != "Bearer analyst-a" {
		return nil, Scope{}, errors.New("unauthorized")
	}
	return s.ctx, Scope{Tenant: "tenant-a", Environment: "prod"}, nil
}
