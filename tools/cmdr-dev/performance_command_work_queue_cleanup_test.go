package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommandWorkQueueProbeCleanupOnlyRemovesOwnedArtifact(t *testing.T) {
	tempDir := t.TempDir()
	privateDir, err := os.MkdirTemp(tempDir, ".cmdr-command-work-queue-perf-")
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(privateDir, "command-work-queue-perf")
	if err := os.WriteFile(binary, []byte("test probe binary"), 0600); err != nil {
		t.Fatal(err)
	}

	cleanupCommandWorkQueueProbe(tempDir, binary, privateDir)
	if _, err := os.Lstat(binary); !os.IsNotExist(err) {
		t.Fatalf("expected owned probe binary removed, got %v", err)
	}
	if _, err := os.Lstat(privateDir); !os.IsNotExist(err) {
		t.Fatalf("expected empty temporary directory removed, got %v", err)
	}
	// Cleanup is safe if the build failed before producing a binary.
	cleanupCommandWorkQueueProbe(tempDir, binary, privateDir)
}

func TestCommandWorkQueueProbeCleanupPreservesUnexpectedFile(t *testing.T) {
	workspace := t.TempDir()
	privateDir, err := os.MkdirTemp(workspace, ".cmdr-command-work-queue-perf-")
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(privateDir, "command-work-queue-perf")
	other := filepath.Join(privateDir, "unrelated-artifact")
	if err := os.WriteFile(binary, []byte("probe binary"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(other, []byte("unrelated data must remain"), 0600); err != nil {
		t.Fatal(err)
	}
	cleanupCommandWorkQueueProbe(workspace, binary, privateDir)
	if _, err := os.Lstat(binary); !os.IsNotExist(err) {
		t.Fatalf("owned binary not removed: %v", err)
	}
	got, err := os.ReadFile(other)
	if err != nil || string(got) != "unrelated data must remain" {
		t.Fatalf("unrelated artifact was removed or modified: %v", err)
	}
	if st, err := os.Stat(privateDir); err != nil || !st.IsDir() {
		t.Fatalf("nonempty private directory was removed: %v", err)
	}
}

func TestCommandWorkQueueProbeCleanupCannotFollowWorkspaceSymlink(t *testing.T) {
	workspace := t.TempDir()
	outside := t.TempDir()
	privateFile := filepath.Join(outside, "command-work-queue-perf")
	if err := os.WriteFile(privateFile, []byte("unrelated external evidence"), 0600); err != nil {
		t.Fatal(err)
	}
	tempName := filepath.Join(workspace, ".cmdr-command-work-queue-perf-replaced")
	if err := os.Symlink(outside, tempName); err != nil {
		t.Skipf("directory symlink creation unavailable: %v", err)
	}
	cleanupCommandWorkQueueProbe(workspace, filepath.Join(tempName, "command-work-queue-perf"), tempName)
	data, err := os.ReadFile(privateFile)
	if err != nil || string(data) != "unrelated external evidence" {
		t.Fatalf("cleanup escaped workspace root through substituted symlink: %v", err)
	}
}
