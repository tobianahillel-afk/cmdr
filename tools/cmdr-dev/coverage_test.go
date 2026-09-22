package main

import (
	"reflect"
	"testing"
)

func TestClassifyOwnedDocument(t *testing.T) {
	cases := []struct {
		doc  SpecDocument
		want string
	}{
		{SpecDocument{ID: "CAP-CMD-001"}, "capability"},
		{SpecDocument{ID: "CMD-MC-001", Type: "screen"}, "screen"},
		{SpecDocument{ID: "journey-alert-to-incident", Domain: "13-user-journeys"}, "journey"},
		{SpecDocument{ID: "contract-x", Path: "cmdr-product-spec/17-implementation-contracts/x.md"}, "implementation-contract"},
		{SpecDocument{ID: "incident", Path: "cmdr-product-spec/05-domain-model/objects/incident.md"}, "canonical-object"},
	}
	for _, tc := range cases {
		if got := classifyOwnedDocument(tc.doc); got != tc.want {
			t.Fatalf("%s: want %s got %s", tc.doc.ID, tc.want, got)
		}
	}
}

func TestBuildProductGraphCreatesSourceBackedAndReferenceOnlyEntities(t *testing.T) {
	inventory := SpecInventory{
		TreeDigest: "digest",
		Documents: []SpecDocument{
			{
				Path:          "cmdr-product-spec/cap.md",
				SHA256:        "abc",
				ID:            "CAP-CMD-001",
				Active:        true,
				Canonical:     true,
				Requirements:  []string{"REQ-PROD-001"},
				OpenDecisions: []string{"OPEN-007"},
				Permissions:   []string{"perm.command.read"},
				References:    []string{"CAP-CMD-002", "REQ-PROD-001"},
			},
			{
				Path:      "cmdr-product-spec/cap2.md",
				SHA256:    "def",
				ID:        "CAP-CMD-002",
				Active:    true,
				Canonical: true,
			},
		},
	}
	graph, err := buildProductGraph(inventory)
	if err != nil {
		t.Fatal(err)
	}
	entities := map[string]ProductEntity{}
	for _, entity := range graph.Entities {
		entities[entity.ID] = entity
	}
	if entities["CAP-CMD-001"].ReferenceOnly {
		t.Fatal("owned capability unexpectedly reference-only")
	}
	if !entities["REQ-PROD-001"].ReferenceOnly || entities["REQ-PROD-001"].Kind != "requirement" {
		t.Fatalf("unexpected requirement entity: %#v", entities["REQ-PROD-001"])
	}
	if !entities["OPEN-007"].ReferenceOnly || entities["OPEN-007"].Kind != "open-decision" {
		t.Fatalf("unexpected open decision entity: %#v", entities["OPEN-007"])
	}
	if !entities["perm.command.read"].ReferenceOnly || entities["perm.command.read"].Kind != "permission" {
		t.Fatalf("unexpected permission entity: %#v", entities["perm.command.read"])
	}
}

func TestBuildProductGraphOrderingIsDeterministic(t *testing.T) {
	inventory := SpecInventory{
		TreeDigest: "digest",
		Documents: []SpecDocument{
			{Path: "z.md", SHA256: "z", ID: "CAP-CMD-002", Active: true},
			{Path: "a.md", SHA256: "a", ID: "CAP-CMD-001", Active: true, References: []string{"REQ-PROD-002", "REQ-PROD-001"}},
		},
	}
	first, err := buildProductGraph(inventory)
	if err != nil {
		t.Fatal(err)
	}
	second, err := buildProductGraph(inventory)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatal("product graph is not deterministic")
	}
	for i := 1; i < len(first.Entities); i++ {
		if first.Entities[i-1].ID > first.Entities[i].ID {
			t.Fatalf("entities not sorted: %#v", first.Entities)
		}
	}
}

func TestAddEntityRejectsMultipleOwnedSources(t *testing.T) {
	entities := map[string]ProductEntity{}
	if err := addEntity(entities, ProductEntity{ID: "X", SourcePath: "a.md", SourceSHA256: "a"}); err != nil {
		t.Fatal(err)
	}
	if err := addEntity(entities, ProductEntity{ID: "X", SourcePath: "b.md", SourceSHA256: "b"}); err == nil {
		t.Fatal("expected duplicate owned source error")
	}
}
