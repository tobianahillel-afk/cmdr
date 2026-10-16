package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Unlike runDemoCLI tests, this verifies the actual process entry point:
// an invalid operator-supplied owner snapshot must cause a nonzero exit
// without advertising a listener or printing a one-run credential.
func TestDemoEntrypointRejectsMissingOwnerSnapshot(t *testing.T) {
	if os.Getenv("CMDR_TEST_DEMO_MAIN_HELPER") == "1" {
		os.Args = []string{
			"case-timeline-demo",
			"-principal=investigator-42",
			"-dataset=" + filepath.Join(os.Getenv("CMDR_TEST_DEMO_MISSING_DIR"), "missing.json"),
		}
		main()
		return
	}
	dir := t.TempDir()
	cmd := exec.Command(os.Args[0], "-test.run=^TestDemoEntrypointRejectsMissingOwnerSnapshot$")
	cmd.Env = append(os.Environ(),
		"CMDR_TEST_DEMO_MAIN_HELPER=1",
		"CMDR_TEST_DEMO_MISSING_DIR="+dir,
	)
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("main accepted missing owner file: %s", output)
	}
	text := string(output)
	if !strings.Contains(text, "load local owner snapshot") ||
		strings.Contains(text, "Bearer token") ||
		strings.Contains(text, "URL:") {
		t.Fatalf("entrypoint leaked a local listener or failed to report invalid file: %s", text)
	}
}
