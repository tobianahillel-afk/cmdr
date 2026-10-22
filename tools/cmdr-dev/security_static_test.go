package main

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func TestGoSCAReportsOnlyValidatedPublicAdvisoryIDsAndStillBlocks(t *testing.T) {
	stream := strings.Join([]string{
		`{"config":{"protocol_version":"v1.0.0","scanner_name":"govulncheck"}}`,
		`{"SBOM":{"modules":[{"path":"private.example/credentials","version":"v1.2.3"}]}}`,
		`{"finding":{"osv":"GO-2099-0001","trace":[{"module":"private.example/credentials","package":"secret/source","function":"secretAccessToken"}]}}`,
	}, "\n")
	summary, err := parseGovulncheckStream([]byte(stream))
	if err == nil || summary.ActionableFindings != 1 || !strings.Contains(err.Error(), "GO-2099-0001") {
		t.Fatalf("reachable OSV not safely actionable: %#v %v", summary, err)
	}
	for _, secret := range []string{"credentials", "secretAccessToken", "secret/source", "private.example"} {
		if strings.Contains(err.Error(), secret) {
			t.Fatalf("unsafe scanner identity disclosed: %q", err)
		}
	}
}

func TestGoSCARejectsMalformedAdvisoryOutputWithoutWaiver(t *testing.T) {
	got := safeGoVulnAdvisoryIDs([]SCAFinding{
		{OSV: "GO-2099-0002"},
		{OSV: "GO-2099-0002"},
		{OSV: "GO-2099-0001"},
		{OSV: "GO-2099-0003\\npassword=123"},
		{OSV: "CVE-2099-1234 private-source"},
		{OSV: "GO-2099-9999<script>"},
	})
	if got != "GO-2099-0001,GO-2099-0002" {
		t.Fatalf("unvalidated advisory data made it to CI output: %q", got)
	}
	if safeGoVulnAdvisoryIDs([]SCAFinding{{OSV: "GO-2099-A"}, {OSV: "GHSA-secret"}}) != "identifier-unavailable" {
		t.Fatal("missing public advisory must remain safely redacted")
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

func TestGoSecurityModuleRootsFollowReviewedGoRuntimeAdapters(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, root, "tools/cmdr-dev/go.mod", "module example/tools\n\ngo 1.26.8\n")

	roots, err := goSecurityModuleRoots(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(roots) != 1 {
		t.Fatalf("expected only engineering module before runtime implementation, got %d", len(roots))
	}

	for _, module := range []struct {
		path string
		want int
	}{
		{"product-runtime/context-envelope/go.mod", 2},
		{"product-runtime/event-search/go.mod", 3},
		{"product-runtime/event-inspection/go.mod", 4},
	} {
		writeTestFile(t, root, module.path, "module example/runtime\n\ngo 1.26.8\n")
		roots, err = goSecurityModuleRoots(root)
		if err != nil {
			t.Fatal(err)
		}
		if len(roots) != module.want {
			t.Fatalf("after %s expected %d reviewed modules, got %d: %v", module.path, module.want, len(roots), roots)
		}
	}

	writeTestFile(t, root, "product-runtime/unregistered/go.mod", "module example/unregistered\n\ngo 1.26.8\n")
	roots, err = goSecurityModuleRoots(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(roots) != 4 {
		t.Fatalf("unregistered runtime must not silently enter reviewed SAST/SCA scope: %v", roots)
	}
}

func TestGosecProcessingErrorsIdentifyOnlySafeExistingRepoPaths(t *testing.T) {
	root := t.TempDir()
	module := filepath.Join(root, "tools", "cmdr-dev")
	if err := os.MkdirAll(module, 0700); err != nil {
		t.Fatal(err)
	}
	allowed := filepath.Join(module, "current.go")
	if err := os.WriteFile(allowed, []byte("package main\n"), 0600); err != nil {
		t.Fatal(err)
	}
	report := map[string]any{
		"Golang errors": map[string]any{
			allowed: []map[string]any{{"error": "secret compiler diagnostic with credentials"}},
			filepath.Join(root, "..", "foreign-private.go"): []map[string]any{{"error": "outside repo"}},
			"../../leak.go":     []map[string]any{{"error": "outside relative"}},
			"missing.go":        []map[string]any{{"error": "nonexistent source"}},
			"arbitrary-package": []map[string]any{{"error": "not a Go file"}},
		},
		"Issues": []any{},
	}
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	summary, err := parseGosecReport(root, module, raw)
	if err == nil || summary.FindingCount != 0 {
		t.Fatalf("processing error must still block: %v %+v", err, summary)
	}
	out := err.Error()
	if !strings.Contains(out, "processing error(s)") ||
		!strings.Contains(out, "tools/cmdr-dev/current.go") {
		t.Fatalf("processing error failed to identify repository-local Go path: %s", out)
	}
	for _, forbidden := range []string{
		"credentials", "secret compiler", "foreign-private.go", "leak.go",
		"missing.go", "arbitrary-package", root,
	} {
		if strings.Contains(out, forbidden) {
			t.Fatalf("SAST processing error leaked raw path or compiler details: %s", out)
		}
	}
	// An analyzer that returns only package-load errors without a safe local
	// .go path MUST continue blocking, but report no attacker-supplied key.
	report["Golang errors"] = map[string]any{
		"/outside/private/secret-package": []map[string]any{{"error": "secrets"}},
	}
	raw, err = json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	_, err = parseGosecReport(root, module, raw)
	if err == nil || !strings.Contains(err.Error(), "source paths unavailable") ||
		strings.Contains(err.Error(), "secret-package") {
		t.Fatalf("unresolved gosec processing error not safely fail-closed: %v", err)
	}
}

func TestGosecProcessingKindsNeverEchoUntrustedErrorMessage(t *testing.T) {
	cases := []struct {
		diagnostic string
		want       string
	}{
		{"undefined: PersonalSecretToken", "unresolved-symbol"},
		{"could not import private.example/vendor-token", "dependency-import"},
		{"syntax error: expected semicolon near SECRET", "syntax-or-parser"},
		{"panic in analyzer when visiting user data", "analyzer-failure"},
		{"failed to load sensitive package", "package-load"},
		{"credentials=super-secret-password unknown", "unclassified"},
	}
	for _, tc := range cases {
		raw, err := json.Marshal(map[string]string{"error": tc.diagnostic})
		if err != nil {
			t.Fatal(err)
		}
		kinds := safeGosecProcessingKinds(map[string][]json.RawMessage{
			"something.go": {json.RawMessage(raw)},
		})
		if len(kinds) != 1 || kinds[0] != tc.want {
			t.Fatalf("unexpected safe classification for %q: %v", tc.want, kinds)
		}
		if strings.Contains(strings.Join(kinds, ","), "secret") {
			t.Fatal("raw gosec diagnostic escaped through safe kind")
		}
	}
}

// Match the actual pinned kernel module toolchain before attributing SAST,
// SBOM, or benchmark evidence to a Go patch release.
func TestGoSecurityEvidenceVersionMatchesModule(t *testing.T) {
	contents, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(contents), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "go ") {
			continue
		}
		actual := strings.TrimSpace(strings.TrimPrefix(line, "go "))
		if actual != engineeringGoVersion {
			t.Fatalf("engineering security evidence Go version %q does not match go.mod directive %q", engineeringGoVersion, actual)
		}
		return
	}
	t.Fatal("tools/cmdr-dev/go.mod must declare its Go version")
}
