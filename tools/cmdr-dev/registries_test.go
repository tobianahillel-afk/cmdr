package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadActiveScreenIDsStopsBeforeAliases(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "cmdr-product-spec", "00-governance", "registers")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "# Screen Register\n\n## Écrans actifs — 1\n\n| ID | Product |\n|---|---|\n| \\x60CMD-MC-001\\x60 | command |\n\n## Aliases de migration\n\n| \\x60CMD-IWQ-001\\x60 | deprecated |\n"
	if err := os.WriteFile(filepath.Join(dir, "screen-register.md"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	ids, err := loadActiveScreenIDs(root, "cmdr-product-spec")
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) != 1 {
		t.Fatalf("expected 1 active screen, got %d: %v", len(ids), ids)
	}
	if _, ok := ids["CMD-MC-001"]; !ok {
		t.Fatal("missing active screen")
	}
	if _, ok := ids["CMD-IWQ-001"]; ok {
		t.Fatal("deprecated alias must not be active")
	}
}
