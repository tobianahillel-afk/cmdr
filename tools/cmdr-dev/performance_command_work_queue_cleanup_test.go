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

	cleanupCommandWorkQueueProbe(binary, privateDir)
	if _, err := os.Lstat(binary); !os.IsNotExist(err) {
		t.Fatalf("expected owned probe binary removed, got %v", err)
	}
	if _, err := os.Lstat(privateDir); !os.IsNotExist(err) {
		t.Fatalf("expected empty temporary directory removed, got %v", err)
	}
	// Cleanup is safe if the build failed before producing a binary.
	cleanupCommandWorkQueueProbe(binary, privateDir)
}

func TestCommandWorkQueueProbeCleanupPreservesUnexpectedFile(t *testing.T) {
	privateDir, err := os.MkdirTemp(t.TempDir(), ".cmdr-command-work-queue-perf-")
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
	cleanupCommandWorkQueueProbe(binary, privateDir)
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
