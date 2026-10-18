package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

// This exercises the operator's actual CLI entry into the already-tested
// authenticated multi-feed Case/Search viewer. Earlier tests called importer
// and HTTP handler directly, leaving this complete delivery route unchecked.
func TestRealCLIMultipleNormalizedFeedsStartsOnlyAuthorizedViewer(t *testing.T) {
	const caseID = "forensic-case-on-disk"
	owner := snapshotFixture(t)
	first := writeJSONLEvents(t,
		"{\"id\":\"shared-id\",\"timestamp\":\"2026-10-17T10:00:00Z\",\"message\":\"CONFIDENTIAL-ENDPOINT-INDICATOR\"}\n")
	second := writeJSONLEvents(t,
		"{\"id\":\"shared-id\",\"timestamp\":\"2026-10-17T10:00:20Z\",\"message\":\"CONFIDENTIAL-FIREWALL-INDICATOR\"}\n")
	ctx, cancel := context.WithTimeout(context.Background(), 650*time.Millisecond)
	defer cancel()
	var output bytes.Buffer
	err := runDemoCLI(ctx, []string{
		"-principal", "analyst-a", "-snapshot", owner,
		"-into-case", caseID, "-import-jsonl", first,
		"-import-jsonl-extra", second, "-serve-import",
	}, &output)
	if err != nil {
		t.Fatalf("real multi-feed analyst CLI failed: %v", err)
	}
	got := output.String()
	if !strings.Contains(got, "Local analyst viewer: http://127.0.0.1:") ||
		!strings.Contains(got, "Bearer token (this run only):") {
		t.Fatalf("authorized multi-feed analyst viewer did not start: %q",got)
	}
	// The viewer intentionally prints a fresh one-run bearer to its local
	// operator, but NEVER the confidential source previews or raw files.
	for _, secret := range []string{
		"CONFIDENTIAL-ENDPOINT-INDICATOR",
		"CONFIDENTIAL-FIREWALL-INDICATOR",
		"shared-id", first, second,
	} {
		if strings.Contains(got, secret) {
			t.Fatalf("loopback viewer startup disclosed private source bytes: %q",secret)
		}
	}
	// A source selected twice is not two distinct forensic feeds. Refuse
	// before opening a viewer or issuing a bearer for invalid evidence.
	var denied bytes.Buffer
	if err := runDemoCLI(context.Background(), []string{
		"-principal", "analyst-a", "-snapshot", owner,
		"-into-case", caseID, "-import-jsonl", first,
		"-import-jsonl-extra", first, "-serve-import",
	}, &denied); err == nil || denied.Len() != 0 {
		t.Fatalf("duplicate physical source started a viewer: err=%v output=%q",err,denied.String())
	}
}
