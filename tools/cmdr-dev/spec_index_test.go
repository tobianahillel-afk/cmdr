package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseFrontMatter(t *testing.T) {
	content := `---
id: sample-id
domain: 00-governance
status: validated
owner: Example Owner
source-of-truth: canonical
requirements:
  - REQ-PROD-009
  - REQ-PROD-012
---
# Body
CAP-INV-001 uses perm.investigate.read and OPEN-007.
`
	doc := parseSpecDocument("cmdr-product-spec/a.md", []byte(content))
	if doc.ID != "sample-id" || doc.Domain != "00-governance" || !doc.Active || !doc.Canonical {
		t.Fatalf("unexpected document metadata: %#v", doc)
	}
	wantReqs := []string{"REQ-PROD-009", "REQ-PROD-012"}
	if !reflect.DeepEqual(doc.Requirements, wantReqs) {
		t.Fatalf("requirements: want %v got %v", wantReqs, doc.Requirements)
	}
	for _, want := range []string{"CAP-INV-001", "OPEN-007", "REQ-PROD-009", "REQ-PROD-012", "perm.investigate.read"} {
		if !contains(doc.References, want) {
			t.Fatalf("missing reference %s in %v", want, doc.References)
		}
	}
}

func TestArchiveIsNotActive(t *testing.T) {
	doc := parseSpecDocument("cmdr-product-spec/99-archive/old.md", []byte("---\nid: old\nsource-of-truth: canonical\n---\n"))
	if doc.Active {
		t.Fatal("archive document must not be active")
	}
}

func TestRejectDuplicateCanonicalIDs(t *testing.T) {
	docs := []SpecDocument{
		{Path: "a.md", ID: "same", Active: true, Canonical: true},
		{Path: "b.md", ID: "same", Active: true, Canonical: true},
	}
	if err := rejectDuplicateCanonicalIDs(docs); err == nil {
		t.Fatal("expected duplicate id error")
	}
}

func TestDuplicateArchivedIDIsAllowed(t *testing.T) {
	docs := []SpecDocument{
		{Path: "a.md", ID: "same", Active: true, Canonical: true},
		{Path: "99-archive/a.md", ID: "same", Active: false, Canonical: true},
	}
	if err := rejectDuplicateCanonicalIDs(docs); err != nil {
		t.Fatal(err)
	}
}

func TestBuildSpecInventoryDeterministicOrdering(t *testing.T) {
	root := t.TempDir()
	specRoot := filepath.Join(root, "cmdr-product-spec")
	if err := os.MkdirAll(specRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, id := range map[string]string{"z.md": "z-id", "a.md": "a-id"} {
		content := "---\nid: " + id + "\nsource-of-truth: canonical\n---\n"
		if err := os.WriteFile(filepath.Join(specRoot, name), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	inv, err := buildSpecInventory(root, "cmdr-product-spec")
	if err != nil {
		t.Fatal(err)
	}
	if len(inv.Documents) != 2 || !strings.HasSuffix(inv.Documents[0].Path, "a.md") || !strings.HasSuffix(inv.Documents[1].Path, "z.md") {
		t.Fatalf("documents not sorted: %#v", inv.Documents)
	}
	firstDigest := inv.TreeDigest
	inv2, err := buildSpecInventory(root, "cmdr-product-spec")
	if err != nil {
		t.Fatal(err)
	}
	if inv2.TreeDigest != firstDigest {
		t.Fatalf("tree digest changed: %s != %s", inv2.TreeDigest, firstDigest)
	}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
