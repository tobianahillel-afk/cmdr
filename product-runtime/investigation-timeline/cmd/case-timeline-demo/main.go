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
	principal := flags.String("principal", "analyst-a", "synthetic analyst-a or analyst-b")
	snapshot := flags.String("snapshot", "", "absolute path to operator-owned read-only Case snapshot JSON (loopback demo only)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("unexpected positional arguments")
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
	if principal != "analyst-a" && principal != "analyst-b" {
		return nil, errors.New("unknown synthetic principal")
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
	if principal != "analyst-a" && principal != "analyst-b" {
		return errors.New("select analyst-a or analyst-b")
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
	fmt.Fprintf(output, "URL: http://%s/case-timeline?case_id=%s&limit=2\n", listener.Addr(), caseID)
	fmt.Fprintf(output, "Local analyst viewer: http://%s/analyst\n", listener.Addr())
	fmt.Fprintf(output, "Bearer token (this run only): %s\n", token)
	fmt.Fprintf(output, "Principal: %s; available synthetic Case: %s\n", principal, caseID)
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
