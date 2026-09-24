package main

import "testing"

func waveRecord(id, state, canonical string) ImplementationReadinessRecord {
	return ImplementationReadinessRecord{
		Capability: id,
		Name: id,
		State: state,
		DeliveryStatus: "defined",
		DeliveryMode: "planned",
		RegisterPath: "register.md",
		CanonicalFile: canonical,
	}
}

func TestSelectImplementationWavePrefersImplementedFamily(t *testing.T) {
	program := ImplementationReadinessProgram{
		ProductSpecBaseline: "baseline",
		SpecTreeDigest: "digest",
		Status: "PASS",
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
		SpecTreeDigest: "digest",
		Status: "PASS",
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
		SpecTreeDigest: "digest",
		Status: "PASS",
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

func TestSelectImplementationWaveRequiresCanonicalFile(t *testing.T) {
	program := ImplementationReadinessProgram{
		ProductSpecBaseline: "baseline",
		SpecTreeDigest: "digest",
		Status: "PASS",
		Records: []ImplementationReadinessRecord{
			waveRecord("CAP-SET-004", "IMPLEMENTED", "context.md"),
			waveRecord("CAP-SET-005", "READY", ""),
		},
	}
	if _, err := selectImplementationWave(program); err == nil {
		t.Fatal("expected missing canonical-file rejection")
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
