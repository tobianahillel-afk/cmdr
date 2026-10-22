package main

import (
	"os"
	"os/exec"
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

func TestRunSpecBaselineWriteAndCheck(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	specRoot := filepath.Join(root, "cmdr-product-spec")
	if err := os.MkdirAll(specRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "---\nid: sample\nstatus: validated\nsource-of-truth: canonical\n---\n"
	if err := os.WriteFile(filepath.Join(specRoot, "sample.md"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, root, "init")
	runTestGit(t, root, "config", "user.name", "CMDR Test")
	runTestGit(t, root, "config", "user.email", "cmdr-test@example.invalid")
	runTestGit(t, root, "add", "cmdr-product-spec")
	runTestGit(t, root, "commit", "-m", "canonical spec")
	baselineCommit := testGitOutput(t, root, "rev-parse", "HEAD")
	output := filepath.Join(root, "engineering", "spec-index", "baseline.json")
	if _, err := runSpecBaseline(root, "cmdr-product-spec", baselineCommit, output, false); err != nil {
		t.Fatal(err)
	}
	if _, err := runSpecBaseline(root, "cmdr-product-spec", baselineCommit, output, true); err != nil {
		t.Fatalf("expected fresh baseline: %v", err)
	}
	if _, err := runSpecBaseline(root, "cmdr-product-spec", "abc123", output, false); err == nil {
		t.Fatal("an invalid commit must never certify a baseline")
	}

	// Unrelated engineering changes do not invalidate an unchanged Product Spec.
	writeTestFile(t, root, "engineering-note.txt", "not product truth\n")
	runTestGit(t, root, "add", "engineering-note.txt")
	runTestGit(t, root, "commit", "-m", "independent engineering changes")
	if _, err := runSpecBaseline(root, "cmdr-product-spec", baselineCommit, output, true); err != nil {
		t.Fatalf("unchanged spec must remain bound to its canonical source: %v", err)
	}

	// A tracked edit must invalidate the previous commit binding in both
	// write and check modes; regenerating JSON cannot launder a false SHA.
	if err := os.WriteFile(filepath.Join(specRoot, "sample.md"), []byte(content+"changed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := runSpecBaseline(root, "cmdr-product-spec", baselineCommit, output, false); err == nil {
		t.Fatal("expected rejected baseline write for changed spec")
	}
	if _, err := runSpecBaseline(root, "cmdr-product-spec", baselineCommit, output, true); err == nil {
		t.Fatal("expected stale baseline error")
	}
	if err := os.WriteFile(filepath.Join(specRoot, "sample.md"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	// Untracked spec Markdown must not be omitted from the recorded Git tree.
	writeTestFile(t, root, "cmdr-product-spec/untracked.md", "---\nid: untracked\n---\n")
	if _, err := runSpecBaseline(root, "cmdr-product-spec", baselineCommit, output, false); err == nil {
		t.Fatal("an untracked canonical document cannot be attributed to the baseline commit")
	}
}

func TestParseCapabilityFrontMatterAliases(t *testing.T) {
	content := `---
id: CAP-CMD-001
product: command
module: mission-control
requirement_ids:
  - REQ-PROD-013
open_decisions:
  - OPEN-007
permissions:
  - perm.command.read
source-of-truth: canonical
---
`
	doc := parseSpecDocument("cmdr-product-spec/cap.md", []byte(content))
	if doc.Product != "command" || doc.Module != "mission-control" {
		t.Fatalf("unexpected product/module: %#v", doc)
	}
	if !reflect.DeepEqual(doc.Requirements, []string{"REQ-PROD-013"}) {
		t.Fatalf("unexpected requirements: %v", doc.Requirements)
	}
	if !reflect.DeepEqual(doc.OpenDecisions, []string{"OPEN-007"}) {
		t.Fatalf("unexpected open decisions: %v", doc.OpenDecisions)
	}
	if !reflect.DeepEqual(doc.Permissions, []string{"perm.command.read"}) {
		t.Fatalf("unexpected permissions: %v", doc.Permissions)
	}
}

func TestParseInlineCapabilityFrontMatterLists(t *testing.T) {
	content := `---
id: CAP-SET-004
product: platform-settings
module: tenants-and-environments
requirement_ids: [REQ-PROD-008, REQ-UX-006, "REQ-SEC-001"]
open_decisions: []
permissions: [perm.platform-settings.tenant.read, perm.platform-settings.environment.read]
source-of-truth: canonical
---
`
	doc := parseSpecDocument("cmdr-product-spec/cap-inline.md", []byte(content))
	if !reflect.DeepEqual(doc.Requirements, []string{"REQ-PROD-008", "REQ-SEC-001", "REQ-UX-006"}) {
		t.Fatalf("unexpected inline requirements: %v", doc.Requirements)
	}
	if len(doc.OpenDecisions) != 0 {
		t.Fatalf("unexpected inline open decisions: %v", doc.OpenDecisions)
	}
	wantPermissions := []string{"perm.platform-settings.environment.read", "perm.platform-settings.tenant.read"}
	if !reflect.DeepEqual(doc.Permissions, wantPermissions) {
		t.Fatalf("unexpected inline permissions: %v", doc.Permissions)
	}
}
