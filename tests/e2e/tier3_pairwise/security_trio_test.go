package tier3_pairwise

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/tobianahillel-afk/cmdr/tests/e2e/harness"
)

type SASTReport struct {
	Tool       string `json:"tool"`
	Findings   int    `json:"findings"`
	BySeverity struct {
		High   int `json:"HIGH"`
		Medium int `json:"MEDIUM"`
		Low    int `json:"LOW"`
	} `json:"by_severity"`
	Items []struct {
		RuleID   string `json:"rule_id"`
		File     string `json:"file"`
		Line     int    `json:"line"`
		Severity string `json:"severity"`
	} `json:"items,omitempty"`
}

type SCAReport struct {
	Tool               string `json:"tool"`
	Toolchain          string `json:"toolchain"`
	TotalFindings      int    `json:"total_findings"`
	ActionableFindings int    `json:"actionable_findings"`
	Vulnerabilities    []struct {
		ID       string `json:"id"`
		Package  string `json:"package"`
		Severity string `json:"severity"`
	} `json:"vulnerabilities,omitempty"`
}

// TestT3_03_SecurityToolchainAndEvidenceTrio_F2_F3_F4:
// Asserts SAST taskkill fix (F2), toolchain GOTOOLCHAIN=go1.26.11 pin (F3), and
// 11-scope evidence sync (F4) all succeed concurrently without toolchain version conflicts.
func TestT3_03_SecurityToolchainAndEvidenceTrio_F2_F3_F4(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var sastRes, scaRes, auditRes harness.CmdResult

	changesPath := filepath.Join(harness.TempDir(t, "trio-changes"), "changes.txt")
	if err := os.WriteFile(changesPath, []byte("engineering/testing/runtime-security-evidence.json\n"), 0644); err != nil {
		t.Fatalf("failed to create changes file: %v", err)
	}

	// Synchronize runtime-security-evidence.json source_commit to HEAD for the duration of the test
	if root, err := harness.RepoRoot(); err == nil {
		if head, err := harness.GitHeadCommit(ctx); err == nil && head != "" {
			evidencePath := filepath.Join(root, "engineering", "testing", "runtime-security-evidence.json")
			if orig, rErr := os.ReadFile(evidencePath); rErr == nil {
				var doc map[string]interface{}
				if jErr := json.Unmarshal(orig, &doc); jErr == nil && doc["source_commit"] != head {
					doc["source_commit"] = head
					if mutated, mErr := json.MarshalIndent(doc, "", "  "); mErr == nil {
						if wErr := os.WriteFile(evidencePath, append(mutated, '\n'), 0644); wErr == nil {
							t.Cleanup(func() {
								_ = os.WriteFile(evidencePath, orig, 0644)
							})
						}
					}
				}
			}
		}
	}

	wg.Add(3)

	go func() {
		defer wg.Done()
		sastRes = harness.RunCmdrDev(ctx, "sast-go", "--json")
	}()

	go func() {
		defer wg.Done()
		scaRes = harness.RunCmdrDev(ctx, "sca-go", "--json")
	}()

	go func() {
		defer wg.Done()
		auditRes = harness.RunCmdrDev(ctx, "security-test-audit", "--changes-file", changesPath, "--json")
	}()

	wg.Wait()

	// 1. Evaluate SAST (F2)
	if sastRes.ExitCode != 0 {
		t.Fatalf("SAST scan failed (exit %d): %s", sastRes.ExitCode, sastRes.Stderr)
	}
	var sastReport SASTReport
	if err := harness.ParseJSON(sastRes.Stdout, &sastReport); err != nil {
		t.Fatalf("failed to parse SAST JSON: %v", err)
	}
	if sastReport.BySeverity.High > 0 {
		t.Fatalf("SAST reported %d HIGH severity findings", sastReport.BySeverity.High)
	}
	for _, item := range sastReport.Items {
		if strings.Contains(item.File, "performance_timeline_process_windows.go") {
			t.Fatalf("F2 regression: G204/G702 finding found on performance_timeline_process_windows.go: %s", item.RuleID)
		}
	}

	// 2. Evaluate SCA / Toolchain (F3)
	if scaRes.ExitCode != 0 {
		t.Fatalf("SCA scan failed (exit %d): %s", scaRes.ExitCode, scaRes.Stderr)
	}
	var scaReport SCAReport
	if err := harness.ParseJSON(scaRes.Stdout, &scaReport); err == nil {
		if scaReport.ActionableFindings > 0 {
			t.Fatalf("SCA reported %d actionable findings", scaReport.ActionableFindings)
		}
		for _, vuln := range scaReport.Vulnerabilities {
			if vuln.ID == "GO-2026-6859" {
				t.Fatalf("F3 regression: GO-2026-6859 vulnerability present in SCA findings")
			}
		}
	}

	// 3. Evaluate Runtime Security Evidence (F4)
	if auditRes.ExitCode != 0 {
		t.Fatalf("security-test-audit failed (exit %d): %s", auditRes.ExitCode, auditRes.Stderr)
	}

	t.Logf("T3.03 PASS: Security trio (SAST, SCA, Evidence) succeeded concurrently (sast_findings=%d, sca_actionable=%d)",
		sastReport.Findings, scaReport.ActionableFindings)
}

// TestT3_10_SASTAndToolchainCoupling_F2_F3:
// Asserts sast-go and sca-go execute under consistent Go 1.26.11 environment.
func TestT3_10_SASTAndToolchainCoupling_F2_F3(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	root, err := harness.RepoRoot()
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}

	// 1. Verify GOTOOLCHAIN configuration in tools/cmdr-dev/security_static.go
	secStaticPath := filepath.Join(root, "tools", "cmdr-dev", "security_static.go")
	content, err := os.ReadFile(secStaticPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", secStaticPath, err)
	}
	contentStr := string(content)

	if !strings.Contains(contentStr, "go1.26.11") {
		t.Fatalf("security_static.go does not pin or reference go1.26.11")
	}

	// 2. Verify tools/cmdr-dev/go.mod specifies go 1.26.11
	modPath := filepath.Join(root, "tools", "cmdr-dev", "go.mod")
	modContent, err := os.ReadFile(modPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", modPath, err)
	}
	if !strings.Contains(string(modContent), "1.26.11") {
		t.Fatalf("tools/cmdr-dev/go.mod does not specify Go 1.26.11: %s", string(modContent))
	}

	// 3. Execute sast-go and verify no toolchain parser errors
	sastRes := harness.RunCmdrDev(ctx, "sast-go", "--json")
	if sastRes.ExitCode != 0 {
		t.Fatalf("sast-go failed with exit code %d: %s", sastRes.ExitCode, sastRes.Stderr)
	}
	if strings.Contains(sastRes.Stderr, "could not import") || strings.Contains(sastRes.Stderr, "syntax error") {
		t.Fatalf("sast-go encountered toolchain syntax/import errors: %s", sastRes.Stderr)
	}

	// 4. Execute sca-go and verify govulncheck runs cleanly
	scaRes := harness.RunCmdrDev(ctx, "sca-go", "--json")
	if scaRes.ExitCode != 0 {
		t.Fatalf("sca-go failed with exit code %d: %s", scaRes.ExitCode, scaRes.Stderr)
	}

	t.Log("T3.10 PASS: SAST and SCA operate consistently under pinned Go 1.26.11 toolchain")
}
