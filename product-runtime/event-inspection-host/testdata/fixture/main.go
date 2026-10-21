// A loopback-only, synthetic security fixture for the browser-to-Go integration
// test. NOT a production SessionResolver, IAM engine or telemetry provider.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	inspection "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection"
	host "github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection-host"
)

type sessions struct{ token string }

func (s sessions) ResolveInspectionSession(r *http.Request) (context.Context, host.Session, error) {
	if r.Header.Get("Authorization") != "Bearer "+s.token {
		return nil, host.Session{}, errors.New("unauthenticated synthetic client")
	}
	return r.Context(), host.Session{Principal: "analyst", Tenant: "tenant-a", Environment: "prod"}, nil
}

type searchOwner struct{}

func (searchOwner) CanInspectSearchResult(_ context.Context, s host.Session, r host.Request) (bool, error) {
	return s.Principal == "analyst" && s.Tenant == r.TenantRef &&
		r.JobID == "job-1" && r.RunID == "run-1" && r.CorrelationID == "corr-1" &&
		r.EventRef == "event-1" && r.SourceHint == "edr", nil
}

type telemetryOwner struct{}

func (telemetryOwner) ReadAuthorizedInspectionEvent(_ context.Context, s host.Session, r host.Request) (host.AuthorizedEvent, error) {
	if s.Principal != "analyst" || s.Tenant != r.TenantRef ||
		s.Environment != r.EnvironmentRef || r.EventRef != "event-1" {
		return host.AuthorizedEvent{}, nil
	}
	return host.AuthorizedEvent{
		Allowed: true, RawDecision: inspection.AccessDeny, RenderDecision: inspection.AccessAllow,
		Event: inspection.Event{
			Ref: "event-1", TenantRef: s.Tenant, State: inspection.StatePartial,
			Source: inspection.SourceMetadata{SourceRef: "edr", DataSourceRef: "endpoint-a", Kind: "telemetry"},
			Parser: inspection.ParserMetadata{ParserRef: "parser-a", Version: "v1", Status: "degraded"},
			SourceFields: map[string]string{"source.ip": "192.0.2.4"},
			NormalizedFields: map[string]string{"network.client.ip": "192.0.2.4"},
			MissingFields: []string{"network.server.ip"},
			RawPayload: []byte("RAW-PRIVATE-STRING"), RenderedPayload: []byte("rendered authorized"),
			RawDerivedSensitive: map[string]string{"secret.token": "RAW-PRIVATE-STRING"},
			Enrichments: []inspection.EventEnrichment{{
				Producer: "detector-a", Version: "v3", Freshness: "stale",
				Values: map[string]string{"risk.score": "72"},
			}},
		},
	}, nil
}

// Synthetic action authority is deliberately separate from Event.read.
type pivotOwner struct{}
func (pivotOwner) CanDraftPivot(_ context.Context, s host.Session, q host.Request, field, value string) (bool, error) {
	return s.Principal == "analyst" && s.Tenant == "tenant-a" &&
		q.JobID == "job-1" && q.RunID == "run-1" && q.EventRef == "event-1" &&
		field == "source.ip" && value == "192.0.2.4", nil
}

func main() {
	token := os.Getenv("CMDR_TEST_BEARER")
	if len(token) < 32 {
		log.Fatal("synthetic bearer missing")
	}
	handler, err := host.NewHandlerWithPivotAuthorizer(sessions{token: token}, searchOwner{}, telemetryOwner{}, pivotOwner{})
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	fmt.Printf("http://%s\n", listener.Addr().String())
	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
