package main

import "testing"

func validDecisionValidationFixture(class string, performanceSensitive bool) (DecisionRegistry, ResearchRegistry, DecisionValidationRegistry) {
	d := researchDecision(class, "accepted")
	d.EvidenceRefs = []string{"RES-PKT-0001"}
	d.PerformanceSensitive = performanceSensitive
	if performanceSensitive {
		d.BenchmarkRefs = []string{"BENCH-0001"}
	}
	if class == "C" {
		d.CriticalFactors = []string{"architecture-critical"}
		d.AdversarialRefs = []string{"ADV-0001"}
	}
	packet := validSaturatedPacket()
	validation := DecisionValidation{
		ID: "DECVAL-0001", DecisionID: d.ID, PreferredCandidateID: "CAND-LSM", ExternalEvidenceSufficient: true,
		CandidateTradeoffs: []CandidateTradeoff{
			{ID: "CAND-LSM", Name: "LSM implementation", FamilyRef: "FAMILY-LSM", Strengths: []string{"write throughput"}, Weaknesses: []string{"compaction"}, EvidenceRefs: []string{"RES-PKT-0001"}},
			{ID: "CAND-BTREE", Name: "B-tree implementation", FamilyRef: "FAMILY-BTREE", Strengths: []string{"point reads"}, Weaknesses: []string{"write amplification"}, EvidenceRefs: []string{"RES-PKT-0001"}},
		},
		Benchmarks: []DecisionBenchmark{{
			ID: "BENCH-0001", RunAt: "2026-09-23", Workload: "representative ingest/query", Environment: "fixed CI host",
			DatasetScale: "representative synthetic forensic corpus", Representative: true,
			ComparedCandidateIDs: []string{"CAND-LSM", "CAND-BTREE"},
			Results: []BenchmarkResult{
				{CandidateID: "CAND-LSM", Metric: "latency", Unit: "ms", Value: 8.2},
				{CandidateID: "CAND-BTREE", Metric: "latency", Unit: "ms", Value: 11.4},
			},
			ArtifactRef: "artifact:bench-0001",
		}},
		AdversarialReviews: []AdversarialReview{{
			ID: "ADV-0001", RunAt: "2026-09-23", TargetCandidateID: "CAND-LSM",
			Hypothesis: "Preferred candidate fails under adversarial failure/load assumptions",
			FalsificationAttempts: []string{"stress failure recovery", "challenge durability assumptions"},
			CounterEvidence: []string{"compaction pressure remains a known trade-off"},
			Outcome: "survived", EvidenceRefs: []string{"RES-PKT-0001"},
		}},
	}
	return DecisionRegistry{Decisions: []EngineeringDecision{d}}, ResearchRegistry{SchemaVersion: 1, RegistryKind: "engineering-research-evidence", Packets: []ResearchPacket{packet}}, DecisionValidationRegistry{SchemaVersion: 1, RegistryKind: "engineering-decision-validation", Validations: []DecisionValidation{validation}}
}

func TestAcceptedDecisionGatePassesWithRequiredEvidence(t *testing.T) {
	decisions, research, registry := validDecisionValidationFixture("C", true)
	if err := validateDecisionValidationRegistry(registry, decisions, research); err != nil {
		t.Fatal(err)
	}
}

func TestAcceptedDecisionRequiresComparisonAndNoBlockingCounterEvidence(t *testing.T) {
	decisions, research, registry := validDecisionValidationFixture("B", false)
	registry.Validations[0].CandidateTradeoffs = registry.Validations[0].CandidateTradeoffs[:1]
	if err := validateDecisionValidationRegistry(registry, decisions, research); err == nil {
		t.Fatal("expected comparison requirement")
	}
	_, _, registry = validDecisionValidationFixture("B", false)
	decisions, research, _ = validDecisionValidationFixture("B", false)
	registry.Validations[0].BlockingCounterEvidence = []string{"unresolved durability failure"}
	if err := validateDecisionValidationRegistry(registry, decisions, research); err == nil {
		t.Fatal("expected blocking counter-evidence rejection")
	}
}

func TestPerformanceDecisionRequiresRepresentativeComparison(t *testing.T) {
	decisions, research, registry := validDecisionValidationFixture("B", true)
	registry.Validations[0].Benchmarks[0].Representative = false
	if err := validateDecisionValidationRegistry(registry, decisions, research); err == nil {
		t.Fatal("expected representative benchmark rejection")
	}
}

func TestClassCRequiresCleanAdversarialFalsification(t *testing.T) {
	decisions, research, registry := validDecisionValidationFixture("C", false)
	registry.Validations[0].AdversarialReviews[0].BlockingFindings = []string{"recovery model falsified"}
	if err := validateDecisionValidationRegistry(registry, decisions, research); err == nil {
		t.Fatal("expected adversarial blocking finding rejection")
	}
}

func TestPrototypeRequiredWhenExternalEvidenceInsufficient(t *testing.T) {
	decisions, research, registry := validDecisionValidationFixture("B", false)
	registry.Validations[0].ExternalEvidenceSufficient = false
	if err := validateDecisionValidationRegistry(registry, decisions, research); err == nil {
		t.Fatal("expected missing prototype rejection")
	}
	registry.Validations[0].Prototypes = []DecisionPrototype{{
		ID: "PROTO-0001", CandidateID: "CAND-LSM", Question: "Does the candidate satisfy local workload constraints?",
		RunAt: "2026-09-23", ArtifactRef: "artifact:proto-0001", Outcome: "supports",
	}}
	if err := validateDecisionValidationRegistry(registry, decisions, research); err != nil {
		t.Fatal(err)
	}
}

func TestValidationRejectsCandidateFamilyOutsideResearchEvidence(t *testing.T) {
	decisions, research, registry := validDecisionValidationFixture("B", false)
	registry.Validations[0].CandidateTradeoffs[0].FamilyRef = "FAMILY-NOT-RESEARCHED"
	if err := validateDecisionValidationRegistry(registry, decisions, research); err == nil {
		t.Fatal("expected unknown research family rejection")
	}
}
