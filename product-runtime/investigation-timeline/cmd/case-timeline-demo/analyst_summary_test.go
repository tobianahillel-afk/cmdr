package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestAnalystBrowserSummaryViaNode(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node is unavailable; JS behavior is covered when Node is installed in CI")
	}
	cmd := exec.Command(node, "--test", "analyst_summary.test.cjs", "analyst_return.test.cjs")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("analyst browser behavior failed: %v\n%s", err, output)
	}
}

func TestAnalystSummaryIsAnExplicitUserAction(t *testing.T) {
	html, err := analystAssets.ReadFile("analyst.html")
	if err != nil {
		t.Fatal(err)
	}
	script, err := analystAssets.ReadFile("analyst.js")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(html), `id="analyze"`) ||
		!strings.Contains(string(html), `id="analysis"`) ||
		!strings.Contains(string(script), `analyze.addEventListener("click", analyzeAuthorizedCase)`) ||
		!strings.Contains(string(script), `url.searchParams.set("limit", String(pageSize))`) ||
		!strings.Contains(string(script), `headers: { Authorization: "Bearer " + bearer }`) {
		t.Fatal("explicit, authenticated multi-page investigation summary not wired")
	}
}
