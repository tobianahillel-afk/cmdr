package main

import "testing"

func waveRecord(id, state, canonical string) ImplementationReadinessRecord {
	return ImplementationReadinessRecord{
		Capability:     id,
		Name:           id,
		State:          state,
		DeliveryStatus: "defined",
		DeliveryMode:   "planned",
		RegisterPath:   "register.md",
		CanonicalFile:  canonical,
	}
}

func TestSelectImplementationWavePrefersImplementedFamily(t *testing.T) {
	program := ImplementationReadinessProgram{
		ProductSpecBaseline: "baseline",
		SpecTreeDigest:      "digest",
		Status:              "PASS",
		Records: []ImplementationReadinessRecord{
			waveRecord("CAP-CMD-001", "READY", "cmd.md"),
			waveRecord("CAP-SET-004", "IMPLEMENTED", "context.md"),
			waveRecord("CAP-SET-011", "READY", "set-11.md"),
			waveRecord("CAP-SET-005", "READY", "set-05.md"),
			waveRecord("CAP-INV-001", "BLOCKED", "inv.md"),
		},
	}
	got, err := selectImplementationWave(program)
	if err != nil {
		t.Fatal(err)
	}
	if got.Capability != "CAP-SET-005" || got.PreferredFamily != "CAP-SET" {
		t.Fatalf("unexpected reuse-first selection: %#v", got)
	}
	if got.ReadyCandidates != 3 {
		t.Fatalf("unexpected READY denominator: %#v", got)
	}
}

func TestSelectImplementationWaveFallsBackToStableCapabilityOrder(t *testing.T) {
	program := ImplementationReadinessProgram{
		ProductSpecBaseline: "baseline",
		SpecTreeDigest:      "digest",
		Status:              "PASS",
		Records: []ImplementationReadinessRecord{
			waveRecord("CAP-SET-004", "IMPLEMENTED", "context.md"),
			waveRecord("CAP-INV-010", "READY", "inv-10.md"),
			waveRecord("CAP-CMD-002", "READY", "cmd-02.md"),
		},
	}
	got, err := selectImplementationWave(program)
	if err != nil {
		t.Fatal(err)
	}
	if got.Capability != "CAP-CMD-002" || got.PreferredFamily != "" {
		t.Fatalf("unexpected fallback selection: %#v", got)
	}
}

func TestSelectImplementationWaveExcludesNonReadyStates(t *testing.T) {
	program := ImplementationReadinessProgram{
		ProductSpecBaseline: "baseline",
		SpecTreeDigest:      "digest",
		Status:              "PASS",
		Records: []ImplementationReadinessRecord{
			waveRecord("CAP-SET-004", "IMPLEMENTED", "context.md"),
			waveRecord("CAP-SET-005", "BLOCKED", "blocked.md"),
			waveRecord("CAP-SET-006", "PROPOSED", "proposed.md"),
			waveRecord("CAP-SET-007", "READY", "ready.md"),
		},
	}
	got, err := selectImplementationWave(program)
	if err != nil {
		t.Fatal(err)
	}
	if got.Capability != "CAP-SET-007" || got.ReadyCandidates != 1 {
		t.Fatalf("non-READY capability leaked into wave selection: %#v", got)
	}
}

func TestResolveWaveCanonicalFileFallsBackToCanonicalInventoryID(t *testing.T) {
	selection := ImplementationWaveSelection{Capability: "CAP-EPT-065"}
	inventory := SpecInventory{
		SpecRoot: "cmdr-product-spec",
		Documents: []SpecDocument{
			{
				Path: "cmdr-product-spec/11-endpoint-agent/capabilities/cap-ept-065-example.md",
				ID:   "CAP-EPT-065", Active: true, Canonical: true,
			},
			{
				Path: "cmdr-product-spec/99-archive/cap-ept-065-old.md",
				ID:   "CAP-EPT-065", Active: false, Canonical: true,
			},
		},
	}
	got, err := resolveWaveCanonicalFile(selection, inventory)
	if err != nil {
		t.Fatal(err)
	}
	if got != "11-endpoint-agent/capabilities/cap-ept-065-example.md" {
		t.Fatalf("unexpected canonical fallback %q", got)
	}
}

func TestResolveWaveCanonicalFileFailsOnAmbiguousOrMissingCanonicalDocument(t *testing.T) {
	for name, docs := range map[string][]SpecDocument{
		"missing": nil,
		"ambiguous": {
			{Path: "cmdr-product-spec/a.md", ID: "CAP-EPT-065", Active: true, Canonical: true},
			{Path: "cmdr-product-spec/b.md", ID: "CAP-EPT-065", Active: true, Canonical: true},
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := resolveWaveCanonicalFile(
				ImplementationWaveSelection{Capability: "CAP-EPT-065"},
				SpecInventory{SpecRoot: "cmdr-product-spec", Documents: docs},
			)
			if err == nil {
				t.Fatal("expected canonical resolution failure")
			}
		})
	}
}

func TestResolveWaveCanonicalFilePreservesExplicitRegisterPath(t *testing.T) {
	got, err := resolveWaveCanonicalFile(
		ImplementationWaveSelection{Capability: "CAP-SET-005", CanonicalFile: "10-platform-settings/capabilities/cap-set-005.md"},
		SpecInventory{SpecRoot: "cmdr-product-spec"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if got != "10-platform-settings/capabilities/cap-set-005.md" {
		t.Fatalf("unexpected explicit canonical path %q", got)
	}
}

func TestNormalizeCanonicalCapabilityPath(t *testing.T) {
	got, err := normalizeCanonicalCapabilityPath("cmdr-product-spec", "10-platform-settings/capabilities/example.md")
	if err != nil {
		t.Fatal(err)
	}
	if got != "cmdr-product-spec/10-platform-settings/capabilities/example.md" {
		t.Fatalf("unexpected canonical path %q", got)
	}
	if _, err := normalizeCanonicalCapabilityPath("cmdr-product-spec", "../outside.md"); err == nil {
		t.Fatal("expected Product Spec escape rejection")
	}
}
