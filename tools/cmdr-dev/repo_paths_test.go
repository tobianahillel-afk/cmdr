package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestResolveRepoPathRejectsEscapes(t *testing.T) {
	root := t.TempDir()
	if _, err := resolveRepoPath(root, "../outside.txt", true); err == nil {
		t.Fatal("expected relative path escape rejection")
	}
	outside := filepath.Join(filepath.Dir(root), "outside.txt")
	if _, err := resolveRepoPath(root, outside, true); err == nil {
		t.Fatal("expected absolute path escape rejection")
	}
}

func TestResolveRepoPathRejectsSymlinkTraversal(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "link")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if _, err := resolveRepoPath(root, filepath.Join("link", "file.txt"), true); err == nil {
		t.Fatal("expected symlink traversal rejection")
	}
}

func TestWriteAndReadRepoFileStayInsideRoot(t *testing.T) {
	root := t.TempDir()
	path, err := writeRepoFile(root, "nested/evidence.json", []byte("ok\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "nested", "evidence.json")
	if path != want {
		t.Fatalf("unexpected output path: got %s want %s", path, want)
	}
	data, err := readRepoFile(root, "nested/evidence.json")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "ok\n" {
		t.Fatalf("unexpected data %q", data)
	}
	if runtime.GOOS != "windows" {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm()&0o077 != 0 {
			t.Fatalf("generated evidence is too permissive: %o", info.Mode().Perm())
		}
	}
}
