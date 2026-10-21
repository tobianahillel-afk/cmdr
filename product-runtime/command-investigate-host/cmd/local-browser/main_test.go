package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"html"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	queue "github.com/tobianahillel-afk/cmdr/product-runtime/work-queue"
)

func demoRequest(t *testing.T, client *http.Client, base, path, role, secret, accept string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, base+path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if role != "" {
		req.SetBasicAuth(role, secret)
	}
	if accept != "" {
		req.Header.Set("Accept", accept)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func demoBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	contents, err := io.ReadAll(io.LimitReader(resp.Body, 50000))
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}

func TestLocalBrowserDemonstratesRealIncidentToCaseHTMLAndFilters(t *testing.T) {
	secrets, err := newDemoSecrets()
	if err != nil || !secrets.valid() {
		t.Fatalf("secure secrets failed: %v", err)
	}
	if _, err := hex.DecodeString(secrets.analyst); err != nil {
		t.Fatal(err)
	}
	if secrets.analyst == secrets.commander || len(secrets.commander) != 64 {
		t.Fatal("secrets not separately randomized")
	}
	app, err := newLocalDemoHandler(secrets)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(app)
	defer server.Close()
	client := server.Client()
	incidentPath := "/command/work-queue?view=incidents&selected=Incident:incident-a"
	start := demoRequest(t, client, server.URL, incidentPath, "analyst", secrets.analyst, "text/html")
	page := demoBody(t, start)
	if start.StatusCode != http.StatusOK || !strings.Contains(page, "incident-a") ||
		!strings.Contains(html.UnescapeString(page), "case_id=case-a") ||
		!strings.Contains(page, "investigating") {
		t.Fatalf("selected Incident and independent Case link absent: %d %s", start.StatusCode, page)
	}
	if strings.Contains(page, secrets.analyst) || strings.Contains(page, secrets.commander) {
		t.Fatal("secret exposed in authorized HTML")
	}
	casePath := "/case-timeline?case_id=case-a&limit=1&offset=0"
	first := demoRequest(t, client, server.URL, casePath, "analyst", secrets.analyst, "text/html")
	firstPage := demoBody(t, first)
	if first.StatusCode != http.StatusOK || !strings.Contains(firstPage, "event-v1") ||
		!strings.Contains(firstPage, "Remplacée") || strings.Contains(firstPage, "event-v2</th>") ||
		!strings.Contains(first.Header.Get("Content-Security-Policy"), "form-action 'self'") {
		t.Fatalf("real Case HTML page failed: %d %s", first.StatusCode, firstPage)
	}
	next := regexp.MustCompile(`rel="next" href="([^"]+)"`).FindStringSubmatch(firstPage)
	if len(next) != 2 {
		t.Fatalf("missing next authorized page: %s", firstPage)
	}
	second := demoRequest(t, client, server.URL, html.UnescapeString(next[1]), "analyst", secrets.analyst, "text/html")
	secondPage := demoBody(t, second)
	if second.StatusCode != http.StatusOK || !strings.Contains(secondPage, "event-v2") ||
		!strings.Contains(secondPage, "event-v1") || !strings.Contains(secondPage, "Non déterminé") {
		t.Fatalf("followed authorized version navigation failed: %d %s", second.StatusCode, secondPage)
	}
	filter := demoRequest(t, client, server.URL, "/case-timeline?case_id=case-a&kind=observed&q=event-v1&from=&to=&limit=1", "analyst", secrets.analyst, "text/html")
	filterPage := demoBody(t, filter)
	if filter.StatusCode != http.StatusOK || !strings.Contains(filterPage, "event-v1") ||
		strings.Contains(filterPage, "event-v2</th>") || !strings.Contains(filterPage, "Réinitialiser les filtres") {
		t.Fatalf("native browser filters did not run: %d %s", filter.StatusCode, filterPage)
	}
	jsonRead := demoRequest(t, client, server.URL, "/case-timeline?case_id=case-a", "analyst", secrets.analyst, "application/json")
	jsonBody := demoBody(t, jsonRead)
	if jsonRead.StatusCode != http.StatusOK || !strings.Contains(jsonBody, "\"superseded\":true") {
		t.Fatalf("original JSON read regressed: %d %s", jsonRead.StatusCode, jsonBody)
	}
}

func TestLocalBrowserRejectsCommanderCaseAndHostSpoofing(t *testing.T) {
	secrets, err := newDemoSecrets()
	if err != nil {
		t.Fatal(err)
	}
	app, err := newLocalDemoHandler(secrets)
	if err != nil {
		t.Fatal(err)
	}
	s := httptest.NewServer(app)
	defer s.Close()
	for _, tc := range []struct {
		path, role, secret string
		want int
	}{
		{"/command/work-queue?selected=Incident:incident-a", "", "", http.StatusUnauthorized},
		{"/command/work-queue?selected=Incident:incident-a", "analyst", "wrong", http.StatusUnauthorized},
		{"/command/work-queue?selected=Incident:incident-a", "commander", secrets.commander, http.StatusOK},
		{"/case-timeline?case_id=case-a", "commander", secrets.commander, http.StatusNotFound},
		{"/case-timeline?case_id=case-a", "analyst", secrets.analyst, http.StatusOK},
		{"/case-timeline?case_id=case-a&tenant=synthetic-tenant", "analyst", secrets.analyst, http.StatusBadRequest},
		{"/unknown", "analyst", secrets.analyst, http.StatusNotFound},
	} {
		resp := demoRequest(t, s.Client(), s.URL, tc.path, tc.role, tc.secret, "text/html")
		body := demoBody(t, resp)
		if resp.StatusCode != tc.want {
			t.Fatalf("access not correctly isolated: %s %s -> %d want %d: %s", tc.role, tc.path, resp.StatusCode, tc.want, body)
		}
		if resp.Header.Get("Cache-Control") != "no-store" {
			t.Fatalf("response cacheable on %s", tc.path)
		}
		if resp.StatusCode == http.StatusUnauthorized &&
			!strings.Contains(resp.Header.Get("WWW-Authenticate"), "Basic ") {
			t.Fatal("native browser cannot authenticate via challenge")
		}
		if resp.StatusCode != http.StatusOK && (strings.Contains(body, "event-v1") ||
			strings.Contains(body, "source-v2")) {
			t.Fatalf("denied response exposed source: %s", body)
		}
	}
	req, err := http.NewRequest(http.MethodGet, s.URL+"/case-timeline?case_id=case-a", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("analyst", secrets.analyst)
	req.Host = "attacker.example:8080"
	resp, err := s.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = demoBody(t, resp)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("host header rebinding was accepted: %d", resp.StatusCode)
	}
	for _, value := range []string{"attacker.example:8080", "[::1]:8080", "localhost:8080", "127.0.0.2:8080", ""} {
		if isLoopbackHost(value) {
			t.Fatalf("unexpected non-canonical host accepted: %q", value)
		}
	}
	if !isLoopbackHost("127.0.0.1:8080") {
		t.Fatal("loopback listener Host rejected")
	}
	for _, bad := range []demoSecrets{{}, {analyst: secrets.analyst, commander: secrets.analyst}} {
		if _, err := newLocalDemoHandler(bad); err == nil {
			t.Fatalf("invalid credentials accepted: %+v", bad)
		}
	}
}

type testServerIntro struct {
	values chan string
}

func (w testServerIntro) Write(p []byte) (int, error) {
	w.values <- string(p)
	return len(p), nil
}

func TestLocalBrowserCanActuallyStartAndStopOnLoopback(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	output := testServerIntro{values: make(chan string, 1)}
	finished := make(chan error, 1)
	go func() { finished <- runLocalBrowser(ctx, output) }()
	var notice string
	select {
	case notice = <-output.values:
	case err := <-finished:
		t.Fatalf("local browser failed before listening: %v", err)
	case <-time.After(15 * time.Second):
		t.Fatal("loopback listener did not start")
	}
	matches := regexp.MustCompile(`URL: (http://127[.]0[.]0[.]1:[0-9]+/[^\n]+)`).FindStringSubmatch(notice)
	credentials := regexp.MustCompile(`Analyst: analyst / ([a-f0-9]{64})`).FindStringSubmatch(notice)
	if len(matches) != 2 || len(credentials) != 2 ||
		strings.Contains(matches[1], credentials[1]) {
		t.Fatalf("listener did not announce a scoped URL and separate local secret: %q", notice)
	}
	req, err := http.NewRequest(http.MethodGet, matches[1], nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("analyst", credentials[1])
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body := demoBody(t, resp)
	if resp.StatusCode != http.StatusOK || !strings.Contains(body, "incident-a") {
		t.Fatalf("running listener did not render real Incident: %d %s", resp.StatusCode, body)
	}
	cancel()
	select {
	case err = <-finished:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("clean signal stop failed: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("local listener did not stop after cancellation")
	}
}


func TestLocalBrowserRefusesShortOrReusedEntropy(t *testing.T) {
	for _, source := range []string{"", "short", strings.Repeat("x", 63)} {
		if _, err := demoSecretsFrom(strings.NewReader(source)); err == nil {
			t.Fatal("short entropy incorrectly yielded credentials")
		}
	}
	// An accidental repeated 32-byte secret is denied before any session
	// handler can be constructed; these bytes are never used in production.
	repeated, err := demoSecretsFrom(bytes.NewReader(bytes.Repeat([]byte{0x5a}, 64)))
	if err != nil || repeated.valid() {
		t.Fatalf("weak pair accepted: err=%v", err)
	}
	if _, err := newLocalDemoHandler(repeated); err == nil {
		t.Fatal("reused analyst and commander credentials were accepted")
	}
	safeInput := append(bytes.Repeat([]byte{0xa1}, 32), bytes.Repeat([]byte{0xb2}, 32)...)
	separate, err := demoSecretsFrom(bytes.NewReader(safeInput))
	if err != nil || !separate.valid() {
		t.Fatalf("separately provided entropy unexpectedly refused: %v", err)
	}
}

func TestLocalBrowserSyntheticOwnerScopeAndPrincipalFailClosed(t *testing.T) {
	ctx := context.Background()
	s := demoCommandSessions{}
	if _, err := s.ResolveSession(ctx, nil); err == nil {
		t.Fatal("Command accepted context without host-authenticated identity")
	}
	scope, err := s.ResolveSession(context.WithValue(ctx, demoPrincipalKey{}, "analyst"), nil)
	if err != nil || scope.Principal != "analyst" || scope.Tenant != demoTenant {
		t.Fatalf("host principal not bound to Command source scope: %+v %v", scope, err)
	}
	owner := demoCommandOwner{linked: queue.LinkedCaseOwnerJoin{
		Links: demoSharedLinks{}, Cases: demoInvestigateCaseOwner{},
	}}
	if _, err := owner.List(ctx, "foreign-tenant", demoEnvironment); err == nil {
		t.Fatal("foreign Command snapshot available")
	}
	for _, who := range []string{"", "not-analyst"} {
		granted, err := owner.CanRead(ctx, who, demoTenant, demoEnvironment, queue.Incident, "incident-a")
		if err != nil || granted {
			t.Fatalf("unknown Command principal authorized: %q", who)
		}
	}
	for _, bad := range []struct{ principal, tenant, environment, incident string }{
		{"untrusted", demoTenant, demoEnvironment, "incident-a"},
		{"analyst", "foreign", demoEnvironment, "incident-a"},
		{"analyst", demoTenant, "prod", "incident-a"},
		{"analyst", demoTenant, demoEnvironment, "other-incident"},
	} {
		_, err := (demoSharedLinks{}).ReadAuthorizedIncidentCaseReferences(ctx,
			bad.principal, bad.tenant, bad.environment, bad.incident)
		if err == nil {
			t.Fatalf("Shared link authority accepted an invalid scope %+v", bad)
		}
	}
	linked, err := owner.ReadAuthorizedLinkedCases(ctx, "commander", demoTenant,
		demoEnvironment, "incident-a")
	if err != nil || linked == nil || len(linked.Cases) != 0 {
		t.Fatalf("Commander obtained Investigate Case owner data: %+v, %v", linked, err)
	}
	linked, err = owner.ReadAuthorizedLinkedCases(ctx, "analyst", demoTenant,
		demoEnvironment, "incident-a")
	if err != nil || linked == nil || len(linked.Cases) != 1 ||
		linked.Cases[0].CaseID != "case-a" {
		t.Fatalf("Case owner didn't independently grant analyst: %+v, %v", linked, err)
	}
	if c, err := (demoInvestigateCaseOwner{}).ReadAuthorizedCaseIdentity(ctx,
		"commander", demoTenant, demoEnvironment, "case-a"); err != nil || c != nil {
		t.Fatal("Command principal obtained Investigate Case identity")
	}
}

func TestLocalBrowserCredentialAndHostInputsCannotBecomeAuthority(t *testing.T) {
	s, err := newDemoSecrets()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ username, password, expected string }{
		{"analyst", s.analyst, "analyst"},
		{"commander", s.commander, "commander"},
		{"analyst", s.commander, ""},
		{"commander", s.analyst, ""},
		{"administrator", s.analyst, ""},
		{"analyst", "", ""},
	} {
		r := httptest.NewRequest(http.MethodGet, "/command/work-queue", nil)
		r.SetBasicAuth(tc.username, tc.password)
		if got := s.principal(r); got != tc.expected {
			t.Fatalf("spoofed authentication %q resolved to %q", tc.username, got)
		}
	}
	if s.principal(httptest.NewRequest(http.MethodGet, "/", nil)) != "" {
		t.Fatal("unprovided credentials resolved to identity")
	}
}

type unavailableConsole struct{}

func (unavailableConsole) Write(_ []byte) (int, error) {
	return 0, errors.New("console unavailable")
}

func TestLocalBrowserStartupRejectsBadSecretsListenerAndConsole(t *testing.T) {
	secrets, err := newDemoSecrets()
	if err != nil {
		t.Fatal(err)
	}
	listenCalls := 0
	for _, invalid := range []demoSecrets{{}, {analyst: secrets.analyst, commander: secrets.analyst}} {
		err := serveDemo(context.Background(), io.Discard, invalid,
			func(_, _ string) (net.Listener, error) {
				listenCalls++
				return nil, errors.New("listener must not start")
			})
		if err == nil || listenCalls != 0 {
			t.Fatalf("invalid secrets opened listener: %v, calls=%d", err, listenCalls)
		}
	}
	err = serveDemo(context.Background(), io.Discard, secrets,
		func(network, address string) (net.Listener, error) {
			listenCalls++
			if network != "tcp4" || address != "127.0.0.1:0" {
				t.Fatalf("unexpected binding %q %q", network, address)
			}
			return nil, errors.New("loopback listener unavailable")
		})
	if err == nil || listenCalls != 1 {
		t.Fatalf("failed TCP4 binding did not prevent startup: %v, calls=%d", err, listenCalls)
	}

	var listener net.Listener
	err = serveDemo(context.Background(), unavailableConsole{}, secrets,
		func(network, address string) (net.Listener, error) {
			var listenErr error
			listener, listenErr = net.Listen(network, address)
			return listener, listenErr
		})
	if err == nil || !strings.Contains(err.Error(), "cannot announce") || listener == nil {
		t.Fatalf("unable to print secret did not prevent serving: %v", err)
	}
	// A failed console write closes the allocated listener immediately.
	conn, dialErr := net.DialTimeout("tcp4", listener.Addr().String(), 100*time.Millisecond)
	if dialErr == nil {
		_ = conn.Close()
		t.Fatal("secret-unavailable startup left a reachable HTTP listener")
	}
}


func TestLocalBrowserRefusesFailedOperatingSystemEntropyBeforeListening(t *testing.T) {
	var output bytes.Buffer
	err := runLocalBrowserFromEntropy(context.Background(), &output,
		strings.NewReader("too-short"))
	if err == nil || !strings.Contains(err.Error(), "generate local credentials") ||
		output.Len() != 0 {
		t.Fatalf("failed crypto source started/announced server: %v output=%q", err, output.String())
	}
	err = runLocalBrowserFromEntropy(context.Background(), &output,
		bytes.NewReader(bytes.Repeat([]byte{42}, 64)))
	if err == nil || output.Len() != 0 {
		t.Fatalf("identical role credentials started local server: %v output=%q", err, output.String())
	}
}


func TestLocalBrowserCaseTimelineOpensOnlyFreshlyAuthorizedSource(t *testing.T) {
	secrets, err := newDemoSecrets()
	if err != nil {
		t.Fatal(err)
	}
	app, err := newLocalDemoHandler(secrets)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(app)
	defer server.Close()
	client := server.Client()
	timelinePath := "/case-timeline?case_id=case-a&kind=observed&limit=1&offset=0"
	first := demoRequest(t, client, server.URL, timelinePath, "analyst", secrets.analyst, "text/html")
	page := demoBody(t, first)
	if first.StatusCode != http.StatusOK {
		t.Fatalf("authorized Case page unavailable: %d %s", first.StatusCode, page)
	}
	urls := regexp.MustCompile(`href="(/demo-source-open\\?[^"]+)"`).FindStringSubmatch(page)
	if len(urls) != 2 {
		t.Fatalf("trusted host did not expose existing source handler: %s", page)
	}
	sourcePath := html.UnescapeString(urls[1])
	if sourcePath != "/demo-source-open?case_id=case-a&entry_id=event-v1" ||
		strings.Contains(sourcePath, "source-v1") || strings.Contains(sourcePath, secrets.analyst) {
		t.Fatalf("client acquired source selector or credential: %q", sourcePath)
	}
	opened := demoRequest(t, client, server.URL, sourcePath, "analyst", secrets.analyst, "text/html")
	preview := demoBody(t, opened)
	if opened.StatusCode != http.StatusOK || !strings.Contains(preview, "original observed event") ||
		!strings.Contains(preview, "\"object_id\":\"object-v1\"") ||
		!strings.Contains(preview, "\"synthetic\":true") ||
		opened.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("same-origin source owner failed to re-authorize: status=%d %s", opened.StatusCode, preview)
	}
	// A browser's Back action restores the original validated Case URL,
	// including the specific authorized filter and pagination selection.
	back := demoRequest(t, client, server.URL, timelinePath, "analyst", secrets.analyst, "text/html")
	if back.StatusCode != http.StatusOK || !strings.Contains(demoBody(t, back), "event-v1") {
		t.Fatal("returning to the Case's same filtered page failed")
	}
	for _, tc := range []struct{ role, password, path string }{
		{"", "", sourcePath},
		{"commander", secrets.commander, sourcePath},
		{"analyst", secrets.analyst, "/demo-source-open?case_id=case-b&entry_id=event-v1"},
		{"analyst", secrets.analyst, "/demo-source-open?case_id=case-a&entry_id=not-visible"},
		{"analyst", secrets.analyst, sourcePath + "&source_ref=source-v1"},
		{"analyst", "wrong-credential", sourcePath},
	} {
		result := demoRequest(t, client, server.URL, tc.path, tc.role, tc.password, "application/json")
		body := demoBody(t, result)
		if result.StatusCode == http.StatusOK ||
			strings.Contains(body, "object-v1") || strings.Contains(body, "original observed event") {
			t.Fatalf("source URL bypassed owner checks: role=%q status=%d body=%q", tc.role, result.StatusCode, body)
		}
	}
}
