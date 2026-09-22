package main

import "testing"

func TestBuildObligationSetExcludesReferenceOnlyHistoricalEntities(t *testing.T) {
	graph := ProductGraph{
		SpecTreeDigest: "digest",
		Entities: []ProductEntity{
			{ID: "CAP-CMD-001", Kind: "capability", SourcePath: "cap.md", SourceSHA256: "a"},
			{ID: "CAP-CMD-999", Kind: "capability", ReferenceOnly: true},
			{ID: "CMD-MC-001", Kind: "screen", SourcePath: "screen.md", SourceSHA256: "b", RegistryStatus: "registered-active"},
			{ID: "contract-a", Kind: "implementation-contract", SourcePath: "contract.md", SourceSHA256: "c"},
		},
		Edges: []ProductEdge{
			{From: "CAP-CMD-001", To: "CAP-CMD-999", Kind: "reference", SourcePath: "cap.md"},
		},
	}
	set, err := buildObligationSet(
		graph,
		map[string][]string{"CAP-CMD-001": {"capability-register.md"}},
		map[string][]string{"REQ-PROD-001": {"source.md"}},
		map[string]struct{}{"CMD-MC-001": {}},
		map[string][]string{"perm.command.read": {"permission-register.md"}},
	)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, obligation := range set.Obligations {
		got[obligation.ID] = true
		if len(obligation.SourcePaths) == 0 {
			t.Fatalf("missing provenance: %#v", obligation)
		}
	}
	for _, want := range []string{
		"OBL:capability:CAP-CMD-001",
		"OBL:requirement:REQ-PROD-001",
		"OBL:screen:CMD-MC-001",
		"OBL:permission:perm.command.read",
		"OBL:implementation-contract:contract-a",
	} {
		if !got[want] {
			t.Fatalf("missing obligation %s", want)
		}
	}
	if got["OBL:capability:CAP-CMD-999"] {
		t.Fatal("reference-only historical capability became an obligation")
	}
	if len(set.UnresolvedReferences) != 1 || set.UnresolvedReferences[0].ID != "CAP-CMD-999" {
		t.Fatalf("unexpected unresolved references: %#v", set.UnresolvedReferences)
	}
}

func TestNewObligationDeduplicatesAndSortsSources(t *testing.T) {
	obligation := newObligation("requirement", "REQ-PROD-001", []string{"z.md", "a.md", "z.md"})
	if len(obligation.SourcePaths) != 2 || obligation.SourcePaths[0] != "a.md" || obligation.SourcePaths[1] != "z.md" {
		t.Fatalf("unexpected source paths: %v", obligation.SourcePaths)
	}
}
