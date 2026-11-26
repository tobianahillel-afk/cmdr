package harness

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// SetupMockSecurityRepo creates an isolated CMDR mock repository under .cmdr-tmp/
// with its own git baseline commit, required bootstrap files, cmdr-product-spec/README.md,
// and runtime boundary stubs so security-test-audit runs without mutating the real repo.
// When syncSourceCommit is true, source_commit in the mock runtime-security-evidence.json
// is updated to the mock repository's HEAD commit SHA.
func SetupMockSecurityRepo(t *testing.T, syncSourceCommit bool) string {
	t.Helper()
	root, err := RepoRoot()
	if err != nil {
		t.Fatalf("failed to locate repo root: %v", err)
	}

	mockRoot := TempDir(t, "sec-evidence-mock")

	copyFiles := []string{
		"AGENTS.md",
		"AI_START_HERE.md",
		"engineering/state/current.json",
		"work/graph.json",
		"engineering/architecture/boundaries.json",
		"engineering/security/security-gates.json",
		"engineering/security/runtime-security-policy.json",
		"engineering/testing/runtime-security-evidence.json",
		"cmdr-product-spec/README.md",
	}

	for _, rel := range copyFiles {
		src := filepath.Join(root, filepath.FromSlash(rel))
		dst := filepath.Join(mockRoot, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			t.Fatalf("failed to create dir for %s: %v", dst, err)
		}
		data, err := os.ReadFile(src)
		if err != nil {
			t.Fatalf("failed to read %s: %v", src, err)
		}
		if err := os.WriteFile(dst, data, 0644); err != nil {
			t.Fatalf("failed to copy %s to mock: %v", rel, err)
		}
	}

	// Create directories and stub files for all 12 Go runtime modules
	for _, mod := range RuntimeGoModules {
		dir := filepath.Join(mockRoot, filepath.FromSlash(mod))
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create runtime stub dir %s: %v", dir, err)
		}
		if err := os.WriteFile(filepath.Join(dir, "stub.go"), []byte("package stub\n"), 0644); err != nil {
			t.Fatalf("failed to write stub.go in %s: %v", dir, err)
		}
	}

	// Create event-search-frontend runtime boundary directory and stub file
	frontendDir := filepath.Join(mockRoot, "product-runtime", "event-search-frontend")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		t.Fatalf("failed to create event-search-frontend dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(frontendDir, "stub.js"), []byte("// stub\n"), 0644); err != nil {
		t.Fatalf("failed to write stub.js in %s: %v", frontendDir, err)
	}

	// Initialize an isolated git repository so currentRepositoryCommit(mockRoot) never escapes to parent .git
	gitSteps := [][]string{
		{"init", "-b", "main"},
		{"config", "user.name", "Test Runner"},
		{"config", "user.email", "test@cmdr.local"},
		{"add", "."},
		{"commit", "-m", "mock security baseline"},
	}
	for _, step := range gitSteps {
		args := append([]string{"-C", mockRoot}, step...)
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("git %v in mockRoot failed: %v (%s)", step, err, string(out))
		}
	}

	if syncSourceCommit {
		headOut, err := exec.Command("git", "-C", mockRoot, "rev-parse", "HEAD").CombinedOutput()
		if err != nil {
			t.Fatalf("failed to read mockRoot HEAD commit: %v (%s)", err, string(headOut))
		}
		mockHead := strings.TrimSpace(string(headOut))
		if mockHead == "" {
			t.Fatal("mockRoot HEAD commit is empty")
		}

		evidencePath := filepath.Join(mockRoot, "engineering", "testing", "runtime-security-evidence.json")
		orig, err := os.ReadFile(evidencePath)
		if err != nil {
			t.Fatalf("failed to read mock runtime-security-evidence.json: %v", err)
		}
		var doc map[string]interface{}
		if err := json.Unmarshal(orig, &doc); err != nil {
			t.Fatalf("failed to unmarshal mock runtime-security-evidence.json: %v", err)
		}
		doc["source_commit"] = mockHead
		mutated, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			t.Fatalf("failed to marshal mock runtime-security-evidence.json: %v", err)
		}
		if err := os.WriteFile(evidencePath, append(mutated, '\n'), 0644); err != nil {
			t.Fatalf("failed to write synced mock runtime-security-evidence.json: %v", err)
		}
	}

	return mockRoot
}
