package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDependencyClosureIsDeterministic(t *testing.T) {
	graph := WorkGraph{Nodes: []WorkNode{
		{ID: "E1-X-001A", DependsOn: []string{"E0-X-001"}},
		{ID: "E0-X-001", DependsOn: []string{"E0-X-000"}},
		{ID: "E0-X-000"},
	}}
	got, err := dependencyClosure("E1-X-001A", graph)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"E0-X-000", "E0-X-001"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("unexpected closure: %v", got)
	}
}

func TestContextBundleDigestIsDeterministic(t *testing.T) {
	root := t.TempDir()
	for path, content := range map[string]string{
		"AGENTS.md": "rules\n",
		"engineering/x.md": "engineering\n",
		"work/x.json": "{}\n",
	} {
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	first := contextCollector{root: root, sources: map[string]*ContextSource{}}
	second := contextCollector{root: root, sources: map[string]*ContextSource{}}
	for _, path := range []string{"work/x.json", "AGENTS.md", "engineering/x.md"} {
		if err := first.add(path, "required"); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{"engineering/x.md", "work/x.json", "AGENTS.md"} {
		if err := second.add(path, "required"); err != nil {
			t.Fatal(err)
		}
	}
	a, err := first.bundle("E1-X-001A", "baseline", []string{"E0-X-001"})
	if err != nil {
		t.Fatal(err)
	}
	b, err := second.bundle("E1-X-001A", "baseline", []string{"E0-X-001"})
	if err != nil {
		t.Fatal(err)
	}
	if a.BundleDigest != b.BundleDigest {
		t.Fatalf("bundle digest is not deterministic: %s != %s", a.BundleDigest, b.BundleDigest)
	}
	if a.Sources[0].Authority != "engineering" {
		t.Fatalf("unexpected authority ordering: %#v", a.Sources)
	}
}

func TestCapabilityProductContextRejectsInactiveReference(t *testing.T) {
	root := t.TempDir()
	regDir := filepath.Join(root, "cmdr-product-spec", "00-governance", "registers")
	if err := os.MkdirAll(regDir, 0o755); err != nil {
		t.Fatal(err)
	}
	registry := "| ID | Name | Canonical file |\n|---|---|---|\n| CAP-CMD-001 | One | `06-command/cap.md` |\n"
	if err := os.WriteFile(filepath.Join(regDir, "capability-register-command.md"), []byte(registry), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := resolveProductContextSources(root, "cmdr-product-spec", ProductRefsV2{Capabilities: []string{"CAP-CMD-002"}})
	if err == nil {
		t.Fatal("expected inactive capability rejection")
	}
}

func TestCapabilityCanonicalSourceResolution(t *testing.T) {
	root := t.TempDir()
	regDir := filepath.Join(root, "cmdr-product-spec", "00-governance", "registers")
	if err := os.MkdirAll(regDir, 0o755); err != nil {
		t.Fatal(err)
	}
	registry := "| ID | Name | Canonical file |\n|---|---|---|\n| CAP-CMD-001 | One | `06-command/cap.md` |\n"
	if err := os.WriteFile(filepath.Join(regDir, "capability-register-command.md"), []byte(registry), 0o600); err != nil {
		t.Fatal(err)
	}
	capPath := filepath.Join(root, "cmdr-product-spec", "06-command", "cap.md")
	if err := os.MkdirAll(filepath.Dir(capPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(capPath, []byte("# capability\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := capabilityCanonicalPath(root, "cmdr-product-spec", "CAP-CMD-001")
	if err != nil {
		t.Fatal(err)
	}
	if got != "cmdr-product-spec/06-command/cap.md" {
		t.Fatalf("unexpected capability path %q", got)
	}
}
