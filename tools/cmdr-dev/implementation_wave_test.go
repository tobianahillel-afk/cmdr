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

func TestSelectImplementationWavePolicyAndScorePopulated(t *testing.T) {
	program := ImplementationReadinessProgram{
		ProductSpecBaseline: "baseline",
		SpecTreeDigest:      "digest",
		Status:              "PASS",
		Records: []ImplementationReadinessRecord{
			waveRecord("CAP-STD-005", "READY", "studio-05.md"),
		},
	}
	got, err := selectImplementationWave(program)
	if err != nil {
		t.Fatal(err)
	}
	if got.SelectionPolicy != "multi-dimensional-v1" {
		t.Fatalf("expected multi-dimensional-v1 policy, got %q", got.SelectionPolicy)
	}
	if got.CandidateScore == nil {
		t.Fatal("expected candidate score to be populated")
	}
	if got.CandidateScore.TotalScore <= 0 {
		t.Fatalf("expected positive total score, got %f", got.CandidateScore.TotalScore)
	}
}

func TestSelectImplementationWavePrioritizesSecurityDimensions(t *testing.T) {
	implemented := map[string]bool{}

	secCandidate := waveRecord("CAP-EPT-002", "READY", "ept-sec.md")
	secCandidate.Name = "Endpoint Agent Secret Redaction and Tenant Isolation"

	standardCandidate := waveRecord("CAP-EPT-003", "READY", "ept-std.md")
	standardCandidate.Name = "Endpoint Agent Generic Counter Display"

	secScore := evaluateWaveCandidateScore(secCandidate, implemented)
	stdScore := evaluateWaveCandidateScore(standardCandidate, implemented)

	if secScore.SecurityScore <= stdScore.SecurityScore {
		t.Fatalf("expected security-critical capability to have higher security score: %f <= %f",
			secScore.SecurityScore, stdScore.SecurityScore)
	}
	if secScore.TotalScore <= stdScore.TotalScore {
		t.Fatalf("expected security-critical capability to have higher total score: %f <= %f",
			secScore.TotalScore, stdScore.TotalScore)
	}
}

func TestSelectImplementationWavePriorityAndRiskDimensions(t *testing.T) {
	implemented := map[string]bool{}

	// Anchored with canonical file and defined status (low risk)
	lowRisk := waveRecord("CAP-CMD-005", "READY", "cmd-05.md")
	lowRisk.DeliveryStatus = "defined"

	// Missing canonical file (higher risk)
	highRisk := waveRecord("CAP-CMD-006", "READY", "")
	highRisk.DeliveryStatus = "experimental"

	scoreLow := evaluateWaveCandidateScore(lowRisk, implemented)
	scoreHigh := evaluateWaveCandidateScore(highRisk, implemented)

	if scoreLow.RiskScore <= scoreHigh.RiskScore {
		t.Fatalf("expected anchored capability to have higher risk score: %f <= %f",
			scoreLow.RiskScore, scoreHigh.RiskScore)
	}
	if scoreLow.PriorityScore <= scoreHigh.PriorityScore {
		t.Fatalf("expected defined status to have higher priority score: %f <= %f",
			scoreLow.PriorityScore, scoreHigh.PriorityScore)
	}
}

func TestEvaluateWaveCandidateScoreAllDimensions(t *testing.T) {
	implemented := map[string]bool{"CAP-STD": true}

	record := waveRecord("CAP-STD-005", "READY", "studio-tool-input.md")
	record.Name = "Studio Tool Input Security Contract"
	record.DeliveryStatus = "defined"

	score := evaluateWaveCandidateScore(record, implemented)

	if score.ReuseScore != 40.0 {
		t.Errorf("expected ReuseScore 40.0, got %f", score.ReuseScore)
	}
	if score.PriorityScore != 30.0 { // 25 (CAP-STD) + 5 (defined)
		t.Errorf("expected PriorityScore 30.0, got %f", score.PriorityScore)
	}
	if score.SecurityScore != 20.0 { // matches "Security"
		t.Errorf("expected SecurityScore 20.0, got %f", score.SecurityScore)
	}
	if score.RiskScore != 15.0 { // 10 (canonical file) + 5 (no blockers)
		t.Errorf("expected RiskScore 15.0, got %f", score.RiskScore)
	}
	if score.DependencyScore != 15.0 { // 10 (base) + 5 (numeric <= 10)
		t.Errorf("expected DependencyScore 15.0, got %f", score.DependencyScore)
	}
	expectedTotal := 40.0 + 30.0 + 20.0 + 15.0 + 15.0
	if score.TotalScore != expectedTotal {
		t.Errorf("expected TotalScore %f, got %f", expectedTotal, score.TotalScore)
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
