package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestGo12610PinnedVerifiedRuntimeTransitions(t *testing.T) {
	for _, tc := range []struct {
		name, dir string
	}{
		{"event-search", eventSearchRuntimeDirectory},
		{"hunt-management", huntManagementRuntimeDirectory},
		{"saved-query-assets", savedQueryAssetsRuntimeDirectory},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			gitInspectionToolchainFixture(t, root, "init", "-q")
			gitInspectionToolchainFixture(t, root, "config", "user.email", "fixture@example.invalid")
			gitInspectionToolchainFixture(t, root, "config", "user.name", "Toolchain regression fixture")
			module := fmt.Sprintf("module github.com/tobianahillel-afk/cmdr/%s\n\n", tc.dir)
			path := tc.dir + "/go.mod"
			writeInspectionToolchainFixture(t, root, path, module+"go 1.26.8\n")
			gitInspectionToolchainFixture(t, root, "add", ".")
			gitInspectionToolchainFixture(t, root, "commit", "-qm", "immutable old verifier evidence")
			old := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")
			if unchanged, err := checkPinnedGoToolchainDelta(root, old, old, tc.dir); err != nil || !unchanged {
				t.Fatalf("unchanged source rejected: %v %v", unchanged, err)
			}
			writeInspectionToolchainFixture(t, root, path, module+"go 1.26.10\n")
			gitInspectionToolchainFixture(t, root, "add", ".")
			gitInspectionToolchainFixture(t, root, "commit", "-qm", "single reviewed stdlib security fix")
			patched := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")
			if unchanged, err := checkPinnedGoToolchainDelta(root, old, patched, tc.dir); err != nil || unchanged {
				t.Fatalf("approved module byte transition rejected: %v %v", unchanged, err)
			}
			writeInspectionToolchainFixture(t, root, tc.dir+"/unreviewed.go", "package runtime\n")
			gitInspectionToolchainFixture(t, root, "add", ".")
			gitInspectionToolchainFixture(t, root, "commit", "-qm", "unapproved source")
			badSource := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")
			if _, err := checkPinnedGoToolchainDelta(root, old, badSource, tc.dir); err == nil {
				t.Fatal("runtime code changed but historical attestation was falsely retained")
			}
			writeInspectionToolchainFixture(t, root, path, module+"go 1.26.11\n")
			gitInspectionToolchainFixture(t, root, "add", ".")
			gitInspectionToolchainFixture(t, root, "commit", "-qm", "unapproved later compiler")
			other := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")
			if _, err := checkPinnedGoToolchainDelta(root, old, other, tc.dir); err == nil {
				t.Fatal("unapproved compiler and executable source were accepted")
			}
		})
	}
	if _, err := checkPinnedGoToolchainDelta(t.TempDir(), "0000000000000000000000000000000000000000", "0000000000000000000000000000000000000000", "nonexistent/runtime"); err == nil {
		t.Fatal("unregistered runtime scope must fail closed before inspecting Git")
	}
}

func TestPinnedGoToolchainTransitionRejectsDeniedExecution(t *testing.T) {
	root := t.TempDir()
	gitInspectionToolchainFixture(t, root, "init", "-q")
	gitInspectionToolchainFixture(t, root, "config", "user.email", "fixture@example.invalid")
	gitInspectionToolchainFixture(t, root, "config", "user.name", "Toolchain validation denial regression")
	module := "module github.com/tobianahillel-afk/cmdr/product-runtime/event-search\n\n"
	path := eventSearchRuntimeDirectory + "/go.mod"
	writeInspectionToolchainFixture(t, root, path, module+"go 1.26.8\n")
	gitInspectionToolchainFixture(t, root, "add", ".")
	gitInspectionToolchainFixture(t, root, "commit", "-qm", "original reviewed source")
	old := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")
	writeInspectionToolchainFixture(t, root, path, module+"go 1.26.10\n")
	gitInspectionToolchainFixture(t, root, "add", ".")
	gitInspectionToolchainFixture(t, root, "commit", "-qm", "approved stdlib transition")
	patched := gitInspectionToolchainFixture(t, root, "rev-parse", "HEAD")

	bin := filepath.Join(root, "fakebin")
	if err := os.MkdirAll(bin, 0755); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(bin, "go")
	payload := "#!/bin/sh\necho 'Access is denied' >&2\nexit 1\n"
	if runtime.GOOS == "windows" {
		executable += ".cmd"
		payload = "@echo off\r\necho Access is denied\r\nexit /b 1\r\n"
	}
	if err := os.WriteFile(executable, []byte(payload), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	unchanged, err := reverifyPinnedGoToolchainTransition(root, old, patched, eventSearchRuntimeDirectory)
	if unchanged || err == nil {
		t.Fatalf("a failed mandatory Go test, including Windows Access is denied, must fail closed: unchanged=%v err=%v", unchanged, err)
	}
}
