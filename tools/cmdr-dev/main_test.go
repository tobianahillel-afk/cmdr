package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSelectNextHonorsDependencyState(t *testing.T) {
	graph := WorkGraph{
		SchemaVersion: 1,
		GraphKind:     "test",
		Nodes: []WorkNode{
			{ID: "A", Status: "IMPLEMENTED"},
			{ID: "B", Status: "READY", DependsOn: []string{"A"}, DependencyTerminalState: "IMPLEMENTED"},
			{ID: "C", Status: "READY", DependsOn: []string{"A"}, DependencyTerminalState: "VERIFIED"},
		},
	}
	got, err := selectNext(graph)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != "B" {
		t.Fatalf("expected B, got %#v", got)
	}
}

func TestSelectNextIsDeterministic(t *testing.T) {
	graph := WorkGraph{
		Nodes: []WorkNode{
			{ID: "Z", Status: "READY"},
			{ID: "A", Status: "READY"},
		},
	}
	got, err := selectNext(graph)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != "A" {
		t.Fatalf("expected A, got %#v", got)
	}
}

func TestValidateAcyclicRejectsCycle(t *testing.T) {
	nodes := []WorkNode{
		{ID: "A", DependsOn: []string{"B"}},
		{ID: "B", DependsOn: []string{"A"}},
	}
	index := map[string]WorkNode{"A": nodes[0], "B": nodes[1]}
	err := validateAcyclic(nodes, index)
	if err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("expected cycle error, got %v", err)
	}
}

func TestDecodeStrictRejectsUnknownField(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "state.json")
	if err := os.WriteFile(path, []byte(`{"schema_version":1,"unknown":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var state CurrentState
	if err := decodeStrict(dir, path, &state); err == nil {
		t.Fatal("expected strict decode error")
	}
}

func TestResolveExplicitRoot(t *testing.T) {
	root := t.TempDir()
	for _, rel := range []string{"AGENTS.md", "AI_START_HERE.md", statePath, graphPath} {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("{}"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	got, err := resolveRoot(root)
	if err != nil {
		t.Fatal(err)
	}
	if got != root {
		t.Fatalf("expected %s, got %s", root, got)
	}
}
