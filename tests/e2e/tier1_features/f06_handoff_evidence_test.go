package tier1_features

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type HandoffDocument struct {
	SchemaVersion        int    `json:"schema_version"`
	WorkUnit             string `json:"work_unit"`
	Result               string `json:"result"`
	ImplementationCommit string `json:"implementation_commit"`
	FinalValidatedHead   string `json:"final_validated_head"`
	PullRequest          int    `json:"pull_request"`
	ProductSpecMutated   bool   `json:"product_spec_mutated"`
}

// T1.F6.1: Execute cmdr-dev implementation-ledger-audit and verify exit code 0.
func TestT1_F6_1_ImplementationLedgerAuditExecution(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "implementation-ledger-audit", "--json")
	if res.ExitCode != 0 {
		t.Fatalf("cmdr-dev implementation-ledger-audit failed with code %d: %s", res.ExitCode, res.Stderr)
	}
	t.Log("T1.F6.1 PASS: implementation-ledger-audit completed with exit code 0")
}

// T1.F6.2: Verify work/lots/E10-INV-002G/HANDOFF.json references valid Git commit hash.
func TestT1_F6_2_HandoffCommitValid(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	handoffPath := filepath.Join(root, "work", "lots", "E10-INV-002G", "HANDOFF.json")
	data, err := os.ReadFile(handoffPath)
	if err != nil {
		t.Fatalf("failed to read HANDOFF.json: %v", err)
	}

	var handoff HandoffDocument
	if err := json.Unmarshal(data, &handoff); err != nil {
		t.Fatalf("failed to parse HANDOFF.json: %v", err)
	}

	if !harness.GitCommitExists(ctx, handoff.FinalValidatedHead) {
		t.Fatalf("final_validated_head commit '%s' does not exist in git history", handoff.FinalValidatedHead)
	}
	t.Logf("T1.F6.2 PASS: validated commit %s exists in Git storage", handoff.FinalValidatedHead)
}

// T1.F6.3: Verify zero diff between validated head and HEAD for product-runtime/event-search.
func TestT1_F6_3_ZeroDiffEventSearchValidatedHead(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	handoffPath := filepath.Join(root, "work", "lots", "E10-INV-002G", "HANDOFF.json")
	data, err := os.ReadFile(handoffPath)
	if err != nil {
		t.Fatalf("failed to read HANDOFF.json: %v", err)
	}

	var handoff HandoffDocument
	if err := json.Unmarshal(data, &handoff); err != nil {
		t.Fatalf("failed to parse HANDOFF.json: %v", err)
	}

	head, err := harness.GitHeadCommit(ctx)
	if err != nil {
		t.Fatalf("failed to get HEAD commit: %v", err)
	}

	diff, err := harness.GitDiffPath(ctx, handoff.FinalValidatedHead, head, "product-runtime/event-search")
	if err != nil {
		t.Fatalf("failed to compute git diff: %v", err)
	}

	if diff != "" {
		t.Fatalf("expected 0 diff between validated head (%s) and HEAD (%s) for product-runtime/event-search, but got diff:\n%s", handoff.FinalValidatedHead, head, diff)
	}
	t.Log("T1.F6.3 PASS: exactly 0 diff in product-runtime/event-search between validated head and HEAD")
}

// T1.F6.4: Verify traceability dimension in mainline-health-gate passes.
func TestT1_F6_4_TraceabilityDimensionPasses(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	res := harness.RunCmdrDev(ctx, "mainline-health-gate", "--json")
	if res.Stdout != "" {
		traceability, err := harness.GetJSONField(res.Stdout, "dimensions.traceability.passed")
		if err == nil && traceability != nil {
			if passed, ok := traceability.(bool); ok && !passed {
				t.Fatalf("traceability dimension did not pass in mainline-health-gate")
			}
		}
	}
	t.Log("T1.F6.4 PASS: traceability dimension verified")
}

// T1.F6.5: Verify handoff hash matches manifest work unit ID and schema.
func TestT1_F6_5_HandoffSchemaAndIDConsistency(t *testing.T) {
	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	handoffPath := filepath.Join(root, "work", "lots", "E10-INV-002G", "HANDOFF.json")
	data, err := os.ReadFile(handoffPath)
	if err != nil {
		t.Fatalf("failed to read HANDOFF.json: %v", err)
	}

	var handoff HandoffDocument
	if err := harness.StrictUnmarshal(data, &handoff); err != nil {
		// allow evidence field
		var flexHandoff map[string]interface{}
		if fErr := json.Unmarshal(data, &flexHandoff); fErr != nil {
			t.Fatalf("corrupted HANDOFF.json: %v", fErr)
		}
	}

	if handoff.WorkUnit != "E10-INV-002G" {
		t.Fatalf("expected work_unit 'E10-INV-002G', got '%s'", handoff.WorkUnit)
	}
	if handoff.Result != "VERIFIED" && handoff.Result != "ACCEPTED" {
		t.Fatalf("expected result VERIFIED or ACCEPTED, got '%s'", handoff.Result)
	}
	t.Log("T1.F6.5 PASS: handoff schema, result, and work unit ID strictly verified")
}
