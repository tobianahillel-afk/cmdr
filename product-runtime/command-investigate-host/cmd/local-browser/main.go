// Command local-browser launches a synthetic, loopback-only demonstration of
// the existing independently authorized Command -> Investigate HTTP journey.
// It is NOT a production IAM, storage adapter, Wazuh connector or release.
package main

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	host "github.com/tobianahillel-afk/cmdr/product-runtime/command-investigate-host"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
	queue "github.com/tobianahillel-afk/cmdr/product-runtime/work-queue"
)

const demoTenant = "synthetic-tenant"
const demoEnvironment = "lab"

type demoPrincipalKey struct{}

type demoSecrets struct {
	analyst   string
	commander string
}

func newDemoSecrets() (demoSecrets, error) {
	return demoSecretsFrom(rand.Reader)
}

// A short/unavailable system entropy source must refuse startup, never
// silently produce fixed or low-entropy demonstration credentials.
func demoSecretsFrom(source io.Reader) (demoSecrets, error) {
	nonce := make([]byte, 64)
	if _, err := io.ReadFull(source, nonce); err != nil {
		return demoSecrets{}, err
	}
	return demoSecrets{analyst: hex.EncodeToString(nonce[:32]),
		commander: hex.EncodeToString(nonce[32:])}, nil
}

func (s demoSecrets) valid() bool {
	return len(s.analyst) == 64 && len(s.commander) == 64 && s.analyst != s.commander
}

// The front-door verifies Basic credentials once for each HTTP request. Both
// product handlers still consult independent owner-scoped permissions.
func (s demoSecrets) principal(r *http.Request) string {
	username, password, ok := r.BasicAuth()
	if !ok || !s.valid() {
		return ""
	}
	switch username {
	case "analyst":
		if subtle.ConstantTimeCompare([]byte(password), []byte(s.analyst)) == 1 {
			return "analyst"
		}
	case "commander":
		if subtle.ConstantTimeCompare([]byte(password), []byte(s.commander)) == 1 {
			return "commander"
		}
	}
	return ""
}

func principalFromContext(ctx context.Context) string {
	p, _ := ctx.Value(demoPrincipalKey{}).(string)
	return p
}

type demoCommandSessions struct{}

func (demoCommandSessions) ResolveSession(ctx context.Context, _ *http.Request) (queue.SessionScope, error) {
	principal := principalFromContext(ctx)
	if principal == "" {
		return queue.SessionScope{}, errors.New("unverified demo session")
	}
	return queue.SessionScope{Principal: principal, Tenant: demoTenant, Environment: demoEnvironment}, nil
}

type demoCommandOwner struct {
	linked queue.LinkedCaseOwnerJoin
}

func (demoCommandOwner) List(_ context.Context, tenant, environment string) ([]queue.Item, error) {
	if tenant != demoTenant || environment != demoEnvironment {
		return nil, errors.New("scope unavailable")
	}
	return []queue.Item{{Kind: queue.Incident, ID: "incident-a",
		Tenant: demoTenant, Environment: demoEnvironment,
		Owner: "analyst", Status: "active", UpdatedAt: "2026-10-21T02:00:00Z"}}, nil
}

func (demoCommandOwner) CanRead(_ context.Context, principal, tenant, environment string, kind queue.Kind, id string) (bool, error) {
	return (principal == "analyst" || principal == "commander") &&
		tenant == demoTenant && environment == demoEnvironment &&
		kind == queue.Incident && id == "incident-a", nil
}

func (o demoCommandOwner) ReadAuthorizedLinkedCases(ctx context.Context, principal, tenant, environment, incidentID string) (*queue.AuthorizedLinkedCases, error) {
	return o.linked.ReadAuthorizedLinkedCases(ctx, principal, tenant, environment, incidentID)
}

type demoSharedLinks struct{}

func (demoSharedLinks) ReadAuthorizedIncidentCaseReferences(_ context.Context, principal, tenant, environment, incidentID string) (*queue.AuthorizedIncidentCaseReferences, error) {
	if tenant != demoTenant || environment != demoEnvironment || incidentID != "incident-a" ||
		(principal != "analyst" && principal != "commander") {
		return nil, errors.New("link owner unavailable")
	}
	return &queue.AuthorizedIncidentCaseReferences{References: []queue.LinkedCaseReference{{
		CaseID: "case-a", LinkSource: "synthetic-shared-link", LinkVersion: "v1",
	}}}, nil
}

type demoInvestigateCaseOwner struct{}

func (demoInvestigateCaseOwner) ReadAuthorizedCaseIdentity(_ context.Context, principal, tenant, environment, id string) (*queue.AuthorizedCaseIdentity, error) {
	if principal != "analyst" || tenant != demoTenant || environment != demoEnvironment || id != "case-a" {
		return nil, nil // denied and absent intentionally indistinguishable
	}
	return &queue.AuthorizedCaseIdentity{
		ID: id, Tenant: tenant, Environment: environment,
		Status: "investigating", StatusSource: "synthetic-investigate",
		StatusVersion: "v1", StatusObservedAt: "2026-10-21T02:00:00Z",
		StatusFreshness: "unknown",
	}, nil
}

func syntheticTimeline() timeline.MemoryTimelineData {
	return timeline.MemoryTimelineData{
		Principals: []timeline.MemoryTimelinePrincipal{
			{ID: "analyst", Tenant: demoTenant, Environment: demoEnvironment},
			{ID: "commander", Tenant: demoTenant, Environment: demoEnvironment},
		},
		Cases: []timeline.MemoryTimelineCase{{
			ID: "case-a", Tenant: demoTenant, Environment: demoEnvironment,
			Readers: []string{"analyst"}, EntryIDs: []string{"event-v1", "event-v2"},
		}},
		Entries: []timeline.MemoryTimelineEntry{
			{CaseID: "case-a", Readers: []string{"analyst"}, Value: timeline.Entry{
				ID: "event-v1", Kind: "observed", Tenant: demoTenant,
				SourceTime: "2026-10-21T02:05:00Z", Timezone: "+00:00",
				SourceRef: "source-v1", Provenance: "synthetic-telemetry",
				Version: 1,
			}},
			{CaseID: "case-a", Readers: []string{"analyst"}, Value: timeline.Entry{
				ID: "event-v2", Kind: "observed", Tenant: demoTenant,
				SourceTime: "2026-10-21T02:06:00Z", Timezone: "+00:00",
				SourceRef: "source-v2", Provenance: "synthetic-telemetry",
				Version: 2, Supersedes: "event-v1",
			}},
		},
		Sources: []timeline.MemoryTimelineSource{
			{Identity: timeline.SourceOwnerIdentity{Reference: "source-v1",
				ObjectID: "object-v1", Kind: "telemetry-event",
				Tenant: demoTenant, Environment: demoEnvironment, CaseID: "case-a"},
				Readers: []string{"analyst"}, Preview: "Synthetic source: original observed event, version 1 (not production evidence)"}},
			{Identity: timeline.SourceOwnerIdentity{Reference: "source-v2",
				ObjectID: "object-v2", Kind: "telemetry-event",
				Tenant: demoTenant, Environment: demoEnvironment, CaseID: "case-a"},
				Readers: []string{"analyst"}, Preview: "Synthetic source: corrected observed event, version 2 (not production evidence)"}},
		},
	}
}

func isLoopbackHost(value string) bool {
	name, _, err := net.SplitHostPort(value)
	if err != nil {
		name = value
	}
	return name == "127.0.0.1"
}

func newLocalDemoHandler(secrets demoSecrets) (http.Handler, error) {
	if !secrets.valid() {
		return nil, errors.New("two separate generated demo credentials are required")
	}
	owner := demoCommandOwner{linked: queue.LinkedCaseOwnerJoin{
		Links: demoSharedLinks{}, Cases: demoInvestigateCaseOwner{},
	}}
	command := queue.NewHTTPHandler(demoCommandSessions{}, owner)
	investigate, err := timeline.NewMemoryCaseTimelineHTTPHandlers(
		func(r *http.Request) (string, error) {
			principal := principalFromContext(r.Context())
			if principal == "" {
				return "", errors.New("no verified Investigate principal")
			}
			return principal, nil
		}, syntheticTimeline())
	if err != nil {
		return nil, err
	}
	routes, err := host.NewHandler(command, timeline.WithLocalSourceLinks(investigate.Timeline))
	if err != nil {
		return nil, err
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		if r == nil || r.URL == nil || !isLoopbackHost(r.Host) {
			http.Error(w, "local demonstration only", http.StatusForbidden)
			return
		}
		principal := secrets.principal(r)
		if principal == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="CMDR local synthetic demo", charset="UTF-8"`)
			http.Error(w, "authentication required", http.StatusUnauthorized)
			return
		}
		verified := r.WithContext(context.WithValue(r.Context(), demoPrincipalKey{}, principal))
		if r.URL.Path == "/demo-source-open" {
			// Existing memory owner rechecks Case, Entry, typed source and
			// destination permission; no browser-provided sourceRef is trusted.
			investigate.OpenSource.ServeHTTP(w, verified)
			return
		}
		routes.ServeHTTP(w, verified)
	}), nil
}

func runLocalBrowser(ctx context.Context, out io.Writer) error {
	return runLocalBrowserFromEntropy(ctx, out, rand.Reader)
}

// Production always supplies the OS cryptographic source; a test can supply
// a failing reader to prove that no listener starts without safe credentials.
func runLocalBrowserFromEntropy(ctx context.Context, out io.Writer, entropy io.Reader) error {
	secrets, err := demoSecretsFrom(entropy)
	if err != nil {
		return fmt.Errorf("generate local credentials: %w", err)
	}
	return serveDemo(ctx, out, secrets, net.Listen)
}

// Keep startup failures observable and prevent serving when the operator
// cannot receive the freshly generated, ephemeral authentication secrets.
func serveDemo(ctx context.Context, out io.Writer, secrets demoSecrets,
	listen func(string, string) (net.Listener, error)) error {
	app, err := newLocalDemoHandler(secrets)
	if err != nil {
		return err
	}
	listener, err := listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return err
	}
	defer listener.Close()
	if _, err := fmt.Fprintf(out, "CMDR SYNTHETIC LOCAL DEMO — NOT PRODUCTION\nURL: http://%s/command/work-queue?view=incidents&selected=Incident:incident-a\nAnalyst: analyst / %s\nCommander (no Case permission): commander / %s\nStop: Ctrl+C\n", listener.Addr(), secrets.analyst, secrets.commander); err != nil {
		return fmt.Errorf("cannot announce temporary local credentials: %w", err)
	}
	server := &http.Server{Handler: app, ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout: 20 * time.Second, WriteTimeout: 20 * time.Second,
		IdleTimeout: 30 * time.Second, MaxHeaderBytes: 8192}
	stopped := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			_ = server.Close()
		case <-stopped:
		}
	}()
	err = server.Serve(listener)
	close(stopped)
	if errors.Is(err, http.ErrServerClosed) && ctx.Err() != nil {
		return nil
	}
	return err
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := runLocalBrowser(ctx, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "CMDR local browser unavailable:", err)
		os.Exit(1)
	}
}
