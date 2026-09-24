package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseGosecReportRedactsCodeAndDetails(t *testing.T) {
	report := map[string]any{
		"Golang errors": map[string]any{},
		"Issues": []map[string]any{{
			"severity": "HIGH", "confidence": "MEDIUM", "cwe": map[string]any{"id": "78"},
			"rule_id": "G204", "file": "/repo/tools/cmdr-dev/example.go", "line": "17",
			"code": "super-secret-source-line", "details": "unsafe command",
		}},
	}
	data, _ := json.Marshal(report)
	summary, err := parseGosecReport("/repo", "/repo/tools/cmdr-dev", data)
	if err == nil {
		t.Fatal("expected blocking SAST finding")
	}
	if summary.FindingCount != 1 || summary.Findings[0].RuleID != "G204" || summary.Findings[0].CWE != "78" {
		t.Fatalf("unexpected SAST summary: %#v", summary)
	}
	encoded, _ := json.Marshal(summary)
	if strings.Contains(string(encoded), "super-secret-source-line") || strings.Contains(string(encoded), "unsafe command") {
		t.Fatal("SAST evidence leaked source/details")
	}
}

func TestParseGosecReportFailsClosedOnProcessingErrors(t *testing.T) {
	data := []byte(`{"Golang errors":{"bad.go":[{"line":1}]},"Issues":[]}`)
	if _, err := parseGosecReport("/repo", "/repo/tools/cmdr-dev", data); err == nil {
		t.Fatal("expected processing-error rejection")
	}
}

func TestParseGovulncheckReachableFindingHasExactIdentity(t *testing.T) {
	stream := strings.Join([]string{
		`{"config":{"protocol_version":"v1.0.0","scanner_name":"govulncheck","scanner_version":"v1.8.0","db":"https://vuln.go.dev","go_version":"go1.26.8"}}`,
		`{"SBOM":{"go_version":"go1.26.8","modules":[{"path":"example.org/dep","version":"v1.2.3"}]}}`,
		`{"finding":{"osv":"GO-2099-0001","fixed_version":"v1.2.4","trace":[{"module":"example.org/dep","version":"v1.2.3"}]}}`,
		`{"finding":{"osv":"GO-2099-0001","fixed_version":"v1.2.4","trace":[{"module":"example.org/dep","version":"v1.2.3","package":"example.org/dep/x","function":"Danger"}]}}`,
	}, "\n")
	summary, err := parseGovulncheckStream([]byte(stream))
	if err == nil {
		t.Fatal("expected reachable vulnerability to block")
	}
	if summary.ActionableFindings != 1 || summary.InformationalFindings != 1 {
		t.Fatalf("unexpected SCA summary: %#v", summary)
	}
	f := summary.Findings[0]
	if f.Ecosystem != "Go" || f.Module != "example.org/dep" || f.Version != "v1.2.3" || f.OSV != "GO-2099-0001" {
		t.Fatalf("missing exact dependency identity: %#v", f)
	}
}

func TestParseGovulncheckUsesSBOMVersionAndDeduplicates(t *testing.T) {
	stream := strings.Join([]string{
		`{"config":{"protocol_version":"v1.0.0","scanner_name":"govulncheck"}}`,
		`{"SBOM":{"modules":[{"path":"example.org/dep","version":"v2.0.0"}]}}`,
		`{"finding":{"osv":"GO-2099-0002","trace":[{"module":"example.org/dep","package":"example.org/dep/x","function":"A"}]}}`,
		`{"finding":{"osv":"GO-2099-0002","trace":[{"module":"example.org/dep","package":"example.org/dep/x","function":"B"}]}}`,
	}, "\n")
	summary, err := parseGovulncheckStream([]byte(stream))
	if err == nil {
		t.Fatal("expected reachable vulnerability to block")
	}
	if summary.ActionableFindings != 1 || summary.Findings[0].Version != "v2.0.0" {
		t.Fatalf("unexpected dedupe/version result: %#v", summary)
	}
}

func TestParseGovulncheckFailsClosedWithoutExactVersion(t *testing.T) {
	stream := strings.Join([]string{
		`{"config":{"protocol_version":"v1.0.0","scanner_name":"govulncheck"}}`,
		`{"finding":{"osv":"GO-2099-0003","trace":[{"module":"example.org/dep","package":"example.org/dep/x","function":"Danger"}]}}`,
	}, "\n")
	if _, err := parseGovulncheckStream([]byte(stream)); err == nil {
		t.Fatal("expected missing exact version rejection")
	}
}

func TestGoSecurityModuleRootsIncludeRuntimeOnlyAfterGoMod(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "tools/cmdr-dev/go.mod", "module example/tools\n\ngo 1.26.8\n")

	roots, err := goSecurityModuleRoots(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(roots) != 1 {
		t.Fatalf("expected only engineering module before runtime implementation, got %d", len(roots))
	}

	writeTestFile(t, root, "product-runtime/context-envelope/go.mod", "module example/runtime\n\ngo 1.26.8\n")
	roots, err = goSecurityModuleRoots(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(roots) != 2 {
		t.Fatalf("expected engineering and runtime modules, got %d", len(roots))
	}
}
