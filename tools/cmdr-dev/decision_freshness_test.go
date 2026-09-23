package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func freshnessFixture(t *testing.T) (string, CurrentState, WorkGraph, DecisionRegistry, ResearchRegistry, DecisionValidationRegistry, DecisionFreshnessRecord) {
	t.Helper()
	root := t.TempDir()
	state := CurrentState{}
	state.ProductSpec.BaselineCommit = "baseline-abc"
	graph := decisionTestGraph()
	decisions, research, validations := validDecisionValidationFixture("C", true)
	decision := decisions.Decisions[0]
	snapshot, err := buildFreshnessSnapshot(root, decision, research, validations, state.ProductSpec.BaselineCommit, nil)
	if err != nil {
		t.Fatal(err)
	}
	record := DecisionFreshnessRecord{
		ID: "FRESH-0001", DecisionID: decision.ID, ReviewedAt: "2026-09-23", FreshUntil: "2026-12-31",
		FreshnessRationale: "critical choice reviewed against current research and representative validation",
		Snapshot:           snapshot,
		RevisitTriggers: []string{
			"date-expiry", "decision-change", "research-evidence-change", "validation-evidence-change",
			"repository-input-change", "product-spec-baseline-change", "material-new-research", "benchmark-regression",
		},
		Signals: []RevisitSignal{
			{Kind: "material-new-research", State: "clear", ObservedAt: "2026-09-23", EvidenceRefs: []string{"RES-PKT-0001"}, Detail: "no material newer evidence identified during review"},
			{Kind: "benchmark-regression", State: "clear", ObservedAt: "2026-09-23", EvidenceRefs: []string{"RES-PKT-0001"}, Detail: "representative benchmark remains within accepted constraints"},
		},
	}
	return root, state, graph, decisions, research, validations, record
}

func canonicalFreshnessPolicy() FreshnessPolicy {
	return FreshnessPolicy{
		SchemaVersion: 1, DefaultPolicy: "deny-stale-accepted-decision",
		AutomaticTriggers: []string{"date-expiry", "decision-change", "research-evidence-change", "validation-evidence-change", "repository-input-change", "product-spec-baseline-change"},
		ExternalTriggers:  []string{"upstream-version-change", "security-advisory", "threat-model-change", "benchmark-regression", "assumption-change", "material-new-research"},
		TrackedInputKinds: []string{"assumption", "threat-model", "benchmark-fixture", "constraint", "dependency-lock", "architecture-contract", "security-policy"},
		SignalStates:      []string{"clear", "fired"},
	}
}

func TestFreshAcceptedDecisionIsReusable(t *testing.T) {
	root, state, _, decisions, research, validations, record := freshnessFixture(t)
	registry := FreshnessRegistry{SchemaVersion: 1, RegistryKind: "engineering-decision-freshness", Records: []DecisionFreshnessRecord{record}}
	if err := validateFreshnessRegistry(canonicalFreshnessPolicy(), registry, decisions); err != nil {
		t.Fatal(err)
	}
	fresh, reasons, current, err := evaluateFreshnessRecord(root, mustResearchDate(t, "2026-10-01"), canonicalFreshnessPolicy(), record, decisions.Decisions[0], research, validations, state)
	if err != nil {
		t.Fatal(err)
	}
	if !fresh || len(reasons) != 0 || digestCanonical(current) != digestCanonical(record.Snapshot) {
		t.Fatalf("expected fresh reusable evidence, reasons=%v", reasons)
	}
}

func TestExpiredAcceptedDecisionBecomesStale(t *testing.T) {
	root, state, _, decisions, research, validations, record := freshnessFixture(t)
	fresh, reasons, _, err := evaluateFreshnessRecord(root, mustResearchDate(t, "2027-01-01"), canonicalFreshnessPolicy(), record, decisions.Decisions[0], research, validations, state)
	if err != nil {
		t.Fatal(err)
	}
	if fresh || !containsString(reasons, "date-expiry") {
		t.Fatalf("expected date expiry, got %v", reasons)
	}
}

func TestResearchOrValidationChangeInvalidatesSnapshot(t *testing.T) {
	root, state, _, decisions, research, validations, record := freshnessFixture(t)
	research.Packets[0].Claims[0].Statement = "materially revised claim"
	fresh, reasons, _, err := evaluateFreshnessRecord(root, mustResearchDate(t, "2026-10-01"), canonicalFreshnessPolicy(), record, decisions.Decisions[0], research, validations, state)
	if err != nil {
		t.Fatal(err)
	}
	if fresh || !containsString(reasons, "research-evidence-change") {
		t.Fatalf("expected research drift, got %v", reasons)
	}

	_, state, _, decisions, research, validations, record = freshnessFixture(t)
	validations.Validations[0].PreferredCandidateID = "CAND-BTREE"
	fresh, reasons, _, err = evaluateFreshnessRecord(root, mustResearchDate(t, "2026-10-01"), canonicalFreshnessPolicy(), record, decisions.Decisions[0], research, validations, state)
	if err != nil {
		t.Fatal(err)
	}
	if fresh || !containsString(reasons, "validation-evidence-change") {
		t.Fatalf("expected validation drift, got %v", reasons)
	}
}

func TestTrackedInputChangeInvalidatesSnapshot(t *testing.T) {
	root, state, _, decisions, research, validations, record := freshnessFixture(t)
	path := filepath.Join(root, "engineering", "research", "assumption.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("assumption v1\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	template := []TrackedInputBinding{{Kind: "assumption", Path: "engineering/research/assumption.md"}}
	snapshot, err := buildFreshnessSnapshot(root, decisions.Decisions[0], research, validations, state.ProductSpec.BaselineCommit, template)
	if err != nil {
		t.Fatal(err)
	}
	record.Snapshot = snapshot
	record.RevisitTriggers = append(record.RevisitTriggers, "assumption-change")
	record.Signals = append(record.Signals, RevisitSignal{Kind: "assumption-change", State: "clear", ObservedAt: "2026-09-23", EvidenceRefs: []string{"RES-PKT-0001"}, Detail: "assumption reviewed"})
	if err := os.WriteFile(path, []byte("assumption v2\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	fresh, reasons, _, err := evaluateFreshnessRecord(root, mustResearchDate(t, "2026-10-01"), canonicalFreshnessPolicy(), record, decisions.Decisions[0], research, validations, state)
	if err != nil {
		t.Fatal(err)
	}
	if fresh || !containsString(reasons, "repository-input-change") {
		t.Fatalf("expected tracked input drift, got %v", reasons)
	}
}

func TestFiredExternalSignalInvalidatesDecision(t *testing.T) {
	root, state, _, decisions, research, validations, record := freshnessFixture(t)
	record.Signals[0].State = "fired"
	fresh, reasons, _, err := evaluateFreshnessRecord(root, mustResearchDate(t, "2026-10-01"), canonicalFreshnessPolicy(), record, decisions.Decisions[0], research, validations, state)
	if err != nil {
		t.Fatal(err)
	}
	if fresh || !containsString(reasons, "material-new-research") {
		t.Fatalf("expected fired signal, got %v", reasons)
	}
}

func TestFreshnessRegistryRequiresClassSpecificTriggers(t *testing.T) {
	_, _, _, decisions, _, _, record := freshnessFixture(t)
	record.RevisitTriggers = []string{"date-expiry", "decision-change", "research-evidence-change", "validation-evidence-change", "repository-input-change", "product-spec-baseline-change"}
	registry := FreshnessRegistry{SchemaVersion: 1, RegistryKind: "engineering-decision-freshness", Records: []DecisionFreshnessRecord{record}}
	if err := validateFreshnessRegistry(canonicalFreshnessPolicy(), registry, decisions); err == nil {
		t.Fatal("expected missing Class C/performance revisit trigger rejection")
	}
}

func mustResearchDate(t *testing.T, value string) time.Time {
	t.Helper()
	d, err := parseResearchDate(value)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
