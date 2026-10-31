package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func gitInspectionToolchainFixture(t *testing.T, root string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git fixture %v failed: %v: %s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func writeInspectionToolchainFixture(t *testing.T, root, path, content string) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(path))
	if err := os.MkdirAll(filepath.Dir(full), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestEventInspectionToolchainDeltaRequiresExactApprovedGoModules(t *testing.T) {
	root := t.TempDir()
	gitInspectionToolchainFixture(t, root, "init", "-q")
	gitInspectionToolchainFixture(t, root, "config", "user.email", "fixture@example.invalid")
	gitInspectionToolchainFixture(t, root, "config", "user.name", "Security fixture")
	original := "module github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection\n\ngo 1.26.8\n"
	patched := "module github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection\n\ngo 1.26.10\n"
	writeInspectionToolchainFixture(t, root, eventInspectionToolchainModPath, original)
	gitInspectionToolchainFixture(t, root, "add", ".")
	gitInspectionToolchainFixture(t, root, "commit", "-qm", "verified predecessor")
	evidenceHead := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")
	if unchanged, err := validateEventInspectionRuntimeSourceDelta(root, evidenceHead, evidenceHead); err != nil || !unchanged {
		t.Fatalf("identical historical source rejected: unchanged=%v err=%v", unchanged, err)
	}

	writeInspectionToolchainFixture(t, root, eventInspectionToolchainModPath, patched)
	gitInspectionToolchainFixture(t, root, "add", ".")
	gitInspectionToolchainFixture(t, root, "commit", "-qm", "only pinned stdlib security upgrade")
	patchHead := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")
	if unchanged, err := validateEventInspectionRuntimeSourceDelta(root, evidenceHead, patchHead); err != nil || unchanged {
		t.Fatalf("exact security toolchain patch rejected or mislabeled unchanged: unchanged=%v err=%v", unchanged, err)
	}

	writeInspectionToolchainFixture(t, root, eventInspectionToolchainModPath,
		"module github.com/tobianahillel-afk/cmdr/product-runtime/event-inspection\n\ngo 1.26.11\n")
	gitInspectionToolchainFixture(t, root, "add", ".")
	gitInspectionToolchainFixture(t, root, "commit", "-qm", "different compiler declaration")
	different := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")
	if _, err := validateEventInspectionRuntimeSourceDelta(root, evidenceHead, different); err == nil {
		t.Fatal("unreviewed compiler declaration accepted")
	}

	writeInspectionToolchainFixture(t, root, "product-runtime/event-inspection/unreviewed.go", "package eventinspection\n")
	gitInspectionToolchainFixture(t, root, "add", ".")
	gitInspectionToolchainFixture(t, root, "commit", "-qm", "additional runtime file")
	extra := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")
	if _, err := validateEventInspectionRuntimeSourceDelta(root, evidenceHead, extra); err == nil {
		t.Fatal("additional executable source accepted as toolchain-only change")
	}

	writeInspectionToolchainFixture(t, root, eventInspectionToolchainModPath, patched)
	gitInspectionToolchainFixture(t, root, "add", ".")
	gitInspectionToolchainFixture(t, root, "commit", "-qm", "restore exact go version but keep extra source")
	restoredWithExtra := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")
	if _, err := validateEventInspectionRuntimeSourceDelta(root, evidenceHead, restoredWithExtra); err == nil {
		t.Fatal("unreviewed executable source hidden by correct go.mod accepted")
	}
}
