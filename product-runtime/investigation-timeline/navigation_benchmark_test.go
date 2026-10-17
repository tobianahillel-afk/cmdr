package investigationtimeline

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// This measures the complete in-process HTTP handler and synthetic Case,
// Timeline Entry and typed-source permission chain, not production IAM,
// transport, persistence, customer load or request-level p95.
func BenchmarkAuthorizedCaseSourceNavigationHTTP(b *testing.B) {
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(
		func(*http.Request) (string, error) { return "analyst-a", nil },
		memorySourcePreviewFixture(),
	)
	if err != nil {
		b.Fatal(err)
	}
	const endpoint = "/case-source-navigation?case_id=case-a&entry_id=entry-a&kind=observed&limit=20&offset=0"
	req := httptest.NewRequest(http.MethodGet, endpoint, nil)
	// Preflight asserts the expected owner identity, not merely HTTP status.
	preflight := httptest.NewRecorder()
	handlers.Navigate.ServeHTTP(preflight, req)
	if preflight.Code != http.StatusOK || !strings.Contains(preflight.Body.String(), `"object_id":"object-a"`) ||
		!strings.Contains(preflight.Body.String(), `"return_is_authority":false`) {
		b.Fatalf("navigation fixture did not authorize canonical source: %d %s", preflight.Code, preflight.Body.String())
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := httptest.NewRecorder()
		handlers.Navigate.ServeHTTP(out, req)
		if out.Code != http.StatusOK {
			b.Fatalf("navigation owner chain denied authorized fixture: %d", out.Code)
		}
	}
	b.StopTimer()
	postflight := httptest.NewRecorder()
	handlers.Navigate.ServeHTTP(postflight, req)
	if postflight.Code != http.StatusOK || !strings.Contains(postflight.Body.String(), `"object_id":"object-a"`) {
		b.Fatal("navigation postflight owner identity failed")
	}
}
