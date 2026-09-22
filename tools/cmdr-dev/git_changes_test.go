package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestValidateFullCommitID(t *testing.T) {
	valid := "0123456789abcdef0123456789abcdef01234567"
	got, err := validateFullCommitID(valid)
	if err != nil {
		t.Fatal(err)
	}
	if got != valid {
		t.Fatalf("unexpected normalized commit %q", got)
	}
	for _, invalid := range []string{
		"deadbeef",
		"0000000000000000000000000000000000000000",
		"zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz",
	} {
		if _, err := validateFullCommitID(invalid); err == nil {
			t.Fatalf("expected invalid commit rejection for %q", invalid)
		}
	}
}

func TestParseNULPathsNormalizesAndSorts(t *testing.T) {
	raw := []byte("work/z.json\x00engineering/a.md\x00work/z.json\x00")
	got, err := parseNULPaths(raw)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"engineering/a.md", "work/z.json"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("unexpected paths: %#v", got)
	}
}

func TestGitChangesIncludesDeletedAndBothRenameSides(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	runTestGit(t, root, "init")
	runTestGit(t, root, "config", "user.name", "CMDR Test")
	runTestGit(t, root, "config", "user.email", "cmdr-test@example.invalid")

	writeTestFile(t, root, "old.txt", "one\n")
	writeTestFile(t, root, "deleted.txt", "remove me\n")
	runTestGit(t, root, "add", ".")
	runTestGit(t, root, "commit", "-m", "base")
	base := testGitOutput(t, root, "rev-parse", "HEAD")

	runTestGit(t, root, "mv", "old.txt", "new.txt")
	if err := os.Remove(filepath.Join(root, "deleted.txt")); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, root, "added.txt", "new\n")
	runTestGit(t, root, "add", "-A")
	runTestGit(t, root, "commit", "-m", "head")
	head := testGitOutput(t, root, "rev-parse", "HEAD")

	output := filepath.Join(root, "changes.txt")
	summary, err := runGitChanges(root, base, head, output)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"added.txt", "deleted.txt", "new.txt", "old.txt"} {
		if !containsString(summary.ChangedPaths, want) {
			t.Fatalf("change set omitted %s: %#v", want, summary.ChangedPaths)
		}
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("change set output is empty")
	}
}

func writeTestFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func runTestGit(t *testing.T, root string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func testGitOutput(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("git %v failed: %v", args, err)
	}
	return string(bytesTrimSpace(out))
}

func bytesTrimSpace(value []byte) []byte {
	start, end := 0, len(value)
	for start < end && (value[start] == ' ' || value[start] == '\n' || value[start] == '\r' || value[start] == '\t') {
		start++
	}
	for end > start && (value[end-1] == ' ' || value[end-1] == '\n' || value[end-1] == '\r' || value[end-1] == '\t') {
		end--
	}
	return value[start:end]
}
