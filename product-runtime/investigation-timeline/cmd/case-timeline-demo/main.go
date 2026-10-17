// Command case-timeline-demo serves a loopback-only, synthetic Case -> Timeline
// integration. It is deliberately not a production identity or owner service.
package main

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := runDemoCLI(ctx, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// CLI parsing is separate from process exit so real requests and rejected
// arguments can be exercised in tests without a shell or subprocess runner.
func runDemoCLI(ctx context.Context, args []string, output io.Writer) error {
	flags := flag.NewFlagSet("case-timeline-demo", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	principal := flags.String("principal", "analyst-a", "synthetic analyst ID from the selected local owner snapshot")
	snapshot := flags.String("snapshot", "", "absolute path to operator-owned read-only Case snapshot JSON (loopback demo only)")
	importJSONL := flags.String("import-jsonl", "", "absolute local normalized event JSONL for offline owner snapshot preparation (no listener)")
	intoCase := flags.String("into-case", "", "existing Case for imported local events (requires -import-jsonl)")
	triageImport := flags.Bool("triage-import", false, "analyze locally imported events; use -with-case for aggregate-only two-Case comparison (no merged owner export)")
	reportCase := flags.String("report-case", "", "produce one authorized offline Case triage JSON receipt and exit (no listener)")
	exportCase := flags.String("export-case", "", "export authorized Case Timeline events as sensitive local JSONL, without listener")
	compareCase := flags.String("compare-case", "", "first authorized Case for an offline two-Case comparison receipt")
	withCase := flags.String("with-case", "", "second authorized Case for the offline comparison (requires -compare-case)")
	compareEvents := flags.Bool("compare-events", false, "explicit sensitive chronological event-level cross-Case JSON receipt")
	compareProximity := flags.Bool("compare-proximity", false, "explicit nearest certain-time event leads across two authorized Cases")
	proximitySeconds := flags.Int("proximity-seconds", 60, "bounded nearest-event window in seconds, 1..300 (requires -compare-proximity)")
	sourceCase := flags.String("source-case", "", "authorized Case for offline source-specific event timeline")
	sourceRef := flags.String("source-ref", "", "exact authorized source reference (requires -source-case)")
	inspectCase := flags.String("inspect-case", "", "Case containing an explicitly selected source Timeline Entry")
	inspectEntry := flags.String("inspect-entry", "", "Timeline Entry ID for an authorized offline typed-source preview (requires -inspect-case)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("unexpected positional arguments")
	}
	windowSet := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "proximity-seconds" {
			windowSet = true
		}
	})
	if windowSet && !*compareProximity {
		return errors.New("proximity window requires comparison proximity mode")
	}
	if *compareProximity && *compareEvents {
		return errors.New("choose exactly one cross-Case comparison mode")
	}
	if *importJSONL != "" || *intoCase != "" || *triageImport {
		if *importJSONL == "" || *intoCase == "" || *reportCase != "" || *exportCase != "" ||
			*compareCase != "" || (*withCase != "" && !*triageImport) || *compareEvents || *compareProximity ||
			*sourceCase != "" || *sourceRef != "" || *inspectCase != "" || *inspectEntry != "" {
			return errors.New("local JSONL import requires a Case and cannot be combined with investigation modes")
		}
		if *triageImport {
			if *withCase != "" {
				return writeDemoImportedJSONLComparison(ctx, *principal, *snapshot, *intoCase, *withCase, *importJSONL, output)
			}
			return writeDemoImportedJSONLTriage(ctx, *principal, *snapshot, *intoCase, *importJSONL, output)
		}
		return writeDemoImportedJSONL(ctx, *principal, *snapshot, *intoCase, *importJSONL, output)
	}
	if *exportCase != "" {
		if *reportCase != "" || *compareCase != "" || *withCase != "" ||
			*compareEvents || *compareProximity || *sourceCase != "" ||
			*sourceRef != "" || *inspectCase != "" || *inspectEntry != "" {
			return errors.New("JSONL export cannot be combined with another investigation mode")
		}
		return writeDemoCaseJSONL(ctx, *principal, *snapshot, *exportCase, output)
	}
	if *inspectCase != "" || *inspectEntry != "" {
		if *inspectCase == "" || *inspectEntry == "" || *reportCase != "" ||
			*compareCase != "" || *withCase != "" || *compareEvents || *compareProximity ||
			*sourceCase != "" || *sourceRef != "" {
			return errors.New("select one complete offline inspection mode")
		}
		return writeDemoEntryInspection(ctx, *principal, *snapshot, *inspectCase, *inspectEntry, output)
	}
	if *sourceCase != "" || *sourceRef != "" {
		if *sourceCase == "" || *sourceRef == "" || *reportCase != "" || *compareCase != "" || *withCase != "" || *compareEvents || *compareProximity {
			return errors.New("select one complete offline report mode")
		}
		return writeDemoSourceTimeline(ctx, *principal, *snapshot, *sourceCase, *sourceRef, output)
	}
	if *reportCase != "" && (*compareCase != "" || *withCase != "") {
		return errors.New("select one offline report mode")
	}
	if *compareCase != "" || *withCase != "" || *compareEvents || *compareProximity {
		if *compareProximity {
			return writeDemoCaseProximity(ctx, *principal, *snapshot, *compareCase, *withCase, *proximitySeconds, output)
		}
		if *compareEvents {
			return writeDemoCaseEventSequence(ctx, *principal, *snapshot, *compareCase, *withCase, output)
		}
		return writeDemoCaseComparison(ctx, *principal, *snapshot, *compareCase, *withCase, output)
	}
	if *reportCase != "" {
		return writeDemoCaseTriage(ctx, *principal, *snapshot, *reportCase, output)
	}
	return serveDemoWithSnapshot(ctx, *principal, *snapshot, output)
}

func randomBearer() (string, error) {
	var secret [32]byte
	if _, err := rand.Read(secret[:]); err != nil {
		return "", fmt.Errorf("generate demo credential: %w", err)
	}
	return hex.EncodeToString(secret[:]), nil
}

func newDemoHandler(token, principal string) (http.Handler, error) {
	return newDemoHandlerWithData(token, principal, demoData())
}

func newDemoHandlerWithData(token, principal string, data timeline.MemoryTimelineData) (http.Handler, error) {
	if len(token) != 64 {
		return nil, errors.New("demo bearer must have 32 random bytes in hex")
	}
	if _, err := hex.DecodeString(token); err != nil {
		return nil, errors.New("invalid demo bearer encoding")
	}
	if _, err := firstDemoAuthorizedCase(data, principal); err != nil {
		return nil, errors.New("principal must be a Case reader in the selected local owner dataset")
	}
	authenticate := func(r *http.Request) (string, error) {
		expected := "Bearer " + token
		if subtle.ConstantTimeCompare([]byte(r.Header.Get("Authorization")), []byte(expected)) != 1 {
			return "", errors.New("not authenticated")
		}
		return principal, nil
	}
	handlers, err := timeline.NewMemoryCaseTimelineHTTPHandlers(authenticate, data)
	if err != nil {
		return nil, fmt.Errorf("assemble synthetic owner chain: %w", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/case-timeline", handlers.Timeline)
	mux.Handle("/case-source-navigation", handlers.Navigate)
	mux.Handle("/demo-source-open", handlers.OpenSource)
	mux.Handle("/demo-cases", handlers.Cases)
	installAnalystDemoRoutes(mux)
	return mux, nil
}

func serveDemo(ctx context.Context, principal string, output io.Writer) error {
	return serveDemoWithSnapshot(ctx, principal, "", output)
}

func serveDemoWithSnapshot(ctx context.Context, principal, snapshot string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || principal == "*" {
		return errors.New("select a synthetic principal from the local owner dataset")
	}
	data, err := loadDemoSnapshot(snapshot)
	if err != nil {
		return err
	}
	caseID, err := firstDemoAuthorizedCase(data, principal)
	if err != nil {
		return err
	}
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil {
		return err
	}
	// Never accept a listener address from HTTP, CLI or an environment flag.
	// The fixture is accessible only on the local IPv4 loopback interface.
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("listen on loopback: %w", err)
	}
	defer listener.Close()
	server := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
		MaxHeaderBytes:    8 << 10,
	}
	fmt.Fprintf(output, "CMDR synthetic Case Timeline demo; NOT a production service\n")
	fmt.Fprintf(output, "URL: http://%s/case-timeline?case_id=%s&limit=2\n", listener.Addr(), url.QueryEscape(caseID))
	fmt.Fprintf(output, "Local analyst viewer: http://%s/analyst\n", listener.Addr())
	fmt.Fprintf(output, "Bearer token (this run only): %s\n", token)
	// Quote operator-controlled identifiers to prevent multiline terminal output.
	fmt.Fprintf(output, "Principal: %q; available synthetic Case: %q\n", principal, caseID)
	fmt.Fprintln(output, "Only authenticated loopback requests are accepted. Ctrl-C exits.")
	ready := make(chan error, 1)
	go func() {
		ready <- server.Serve(listener)
	}()
	select {
	case err := <-ready:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("demo HTTP server: %w", err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("demo HTTP shutdown: %w", err)
		}
		if err := <-ready; err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("demo HTTP server: %w", err)
		}
		return nil
	}
}
