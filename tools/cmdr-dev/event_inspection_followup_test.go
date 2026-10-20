package main

import (
	"strings"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func fixtureGitHead(t *testing.T, root string, args ...string) string {
	t.Helper()
	out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput()
	if err != nil { t.Fatalf("git %v: %s: %v", args, out, err) }
	return string(out)
}

func writePivotFixture(t *testing.T, root, name, value string) {
	t.Helper()
	path := filepath.Join(root, "product-runtime", "event-inspection", name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil { t.Fatal(err) }
	if err := os.WriteFile(path, []byte(value), 0644); err != nil { t.Fatal(err) }
}

func TestPivotRuntimeFollowupGateFailClosedForExtraOrUnreviewedCoreFiles(t *testing.T) {
	root := t.TempDir()
	fixtureGitHead(t, root, "init", "-q")
	fixtureGitHead(t, root, "config", "user.email", "test@example.invalid")
	fixtureGitHead(t, root, "config", "user.name", "CMDR test")
	writePivotFixture(t, root, "pivot.go", "package eventinspection\n")
	fixtureGitHead(t, root, "add", ".")
	fixtureGitHead(t, root, "commit", "-qm", "baseline")
	old := stringTrim(fixtureGitHead(t, root, "rev-parse", "HEAD"))
	// An opportunistic source edit outside the two-file allowlist must
	// fail BEFORE any trusted/expensive compiler execution.
	writePivotFixture(t, root, "pivot.go", "package eventinspection\n// untrusted patch\n")
	writePivotFixture(t, root, "pivot_presence_test.go", "package eventinspection\n")
	writePivotFixture(t, root, "extra.go", "package eventinspection\n")
	fixtureGitHead(t, root, "add", ".")
	fixtureGitHead(t, root, "commit", "-qm", "unverified runtime change")
	head := stringTrim(fixtureGitHead(t, root, "rev-parse", "HEAD"))
	if err := validateEventInspectionPivotFollowupProof(root, old, head); err == nil {
		t.Fatal("extra unreviewed runtime source was permitted")
	}
	// Even with EXACTLY the two named paths, a different blob SHA
	// cannot masquerade as the independently reviewed PR #182 fix.
	fixtureGitHead(t, root, "rm", "product-runtime/event-inspection/extra.go")
	fixtureGitHead(t, root, "commit", "-qm", "two files but wrong source")
	head = stringTrim(fixtureGitHead(t, root, "rev-parse", "HEAD"))
	if err := validateEventInspectionPivotFollowupProof(root, old, head); err == nil {
		t.Fatal("unreviewed pivot.go content was permitted")
	}
	// A source which did not change after the verified commit is accepted
	// without any new snapshot or tampering with historical evidence.
	unchanged, err := validateEventInspectionRuntimeUnchangedOrReverified(root, head, head)
	if err != nil || !unchanged { t.Fatalf("identical runtime snapshot rejected: unchanged=%v err=%v", unchanged, err) }
}

func stringTrim(s string) string {
	return strings.TrimSpace(s)
}
