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
	tick := string(rune(96))
	content := "# Screen Register\n\n## Écrans actifs — 1\n\n| ID | Product |\n|---|---|\n| " + tick + "CMD-MC-001" + tick + " | command |\n\n## Aliases de migration\n\n| " + tick + "CMD-IWQ-001" + tick + " | deprecated |\n"
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

func TestFirstTableCellAcceptsPlainAndBacktickValues(t *testing.T) {
	for line, want := range map[string]string{
		"| CAP-CMD-001 | Name |":            "CAP-CMD-001",
		"| `CAP-SET-014` | Name |":          "CAP-SET-014",
		"| `perm.command.read` | Command |": "perm.command.read",
	} {
		got, ok := firstTableCell(line)
		if !ok || got != want {
			t.Fatalf("%q: got %q ok=%t want %q", line, got, ok, want)
		}
	}
}
