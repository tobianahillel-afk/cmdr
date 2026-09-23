package main

import (
	"fmt"
	"math"
	"path/filepath"
	"regexp"
	"strings"
)

const decisionValidationRegistryPath = "engineering/decisions/decision-validation.json"

type CandidateTradeoff struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	FamilyRef    string   `json:"family_ref"`
	Strengths    []string `json:"strengths"`
	Weaknesses   []string `json:"weaknesses"`
	EvidenceRefs []string `json:"evidence_refs"`
}

type BenchmarkResult struct {
	CandidateID string  `json:"candidate_id"`
	Metric      string  `json:"metric"`
	Unit        string  `json:"unit"`
	Value       float64 `json:"value"`
}

type DecisionBenchmark struct {
	ID                   string            `json:"id"`
	RunAt                string            `json:"run_at"`
	Workload             string            `json:"workload"`
	Environment          string            `json:"environment"`
	DatasetScale         string            `json:"dataset_scale"`
	Representative       bool              `json:"representative"`
	ComparedCandidateIDs []string          `json:"compared_candidate_ids"`
	Results              []BenchmarkResult `json:"results"`
	ArtifactRef          string            `json:"artifact_ref"`
}

type DecisionPrototype struct {
	ID                     string   `json:"id"`
	CandidateID            string   `json:"candidate_id"`
	Question               string   `json:"question"`
	RunAt                  string   `json:"run_at"`
	ArtifactRef            string   `json:"artifact_ref"`
	Outcome                string   `json:"outcome"`
	UnresolvedUncertainties []string `json:"unresolved_uncertainties"`
}

type AdversarialReview struct {
	ID                    string   `json:"id"`
	RunAt                 string   `json:"run_at"`
	TargetCandidateID     string   `json:"target_candidate_id"`
	Hypothesis            string   `json:"hypothesis"`
	FalsificationAttempts []string `json:"falsification_attempts"`
	CounterEvidence       []string `json:"counter_evidence"`
	BlockingFindings      []string `json:"blocking_findings"`
	Outcome               string   `json:"outcome"`
	EvidenceRefs          []string `json:"evidence_refs"`
}

type DecisionValidation struct {
	ID                         string              `json:"id"`
	DecisionID                 string              `json:"decision_id"`
	PreferredCandidateID       string              `json:"preferred_candidate_id"`
	ExternalEvidenceSufficient bool                `json:"external_evidence_sufficient"`
	CandidateTradeoffs         []CandidateTradeoff `json:"candidate_tradeoffs"`
	Benchmarks                 []DecisionBenchmark `json:"benchmarks"`
	Prototypes                 []DecisionPrototype `json:"prototypes"`
	AdversarialReviews         []AdversarialReview `json:"adversarial_reviews"`
	BlockingCounterEvidence    []string            `json:"blocking_counter_evidence"`
	UnresolvedMaterialUncertainty []string         `json:"unresolved_material_uncertainty"`
}

type DecisionValidationRegistry struct {
	SchemaVersion int                  `json:"schema_version"`
	RegistryKind  string               `json:"registry_kind"`
	Validations   []DecisionValidation `json:"validations"`
}

type DecisionGateAuditSummary struct {
	Validations          int `json:"validations"`
	AcceptedEvaluated    int `json:"accepted_evaluated"`
	Candidates           int `json:"candidates"`
	RepresentativeBenchmarks int `json:"representative_benchmarks"`
	Prototypes           int `json:"prototypes"`
	AdversarialReviews   int `json:"adversarial_reviews"`
}

var decisionValidationIDPattern = regexp.MustCompile(`^DECVAL-[0-9]{4,}$`)
var candidateIDPattern = regexp.MustCompile(`^CAND-[A-Z0-9][A-Z0-9-]{1,63}$`)
var benchmarkIDPattern = regexp.MustCompile(`^BENCH-[0-9]{4,}$`)
var prototypeIDPattern = regexp.MustCompile(`^PROTO-[0-9]{4,}$`)
var adversarialIDPattern = regexp.MustCompile(`^ADV-[0-9]{4,}$`)

var knownPrototypeOutcomes = map[string]bool{"supports": true, "rejects": true, "inconclusive": true}
var knownAdversarialOutcomes = map[string]bool{"survived": true, "rejected": true, "inconclusive": true}

func runDecisionGateAudit(root string) (DecisionGateAuditSummary, error) {
	policy, research, decisions, err := loadResearchState(root)
	if err != nil {
		return DecisionGateAuditSummary{}, err
	}
	if err := validateResearchRegistry(policy, research, decisions); err != nil {
		return DecisionGateAuditSummary{}, err
	}
	var registry DecisionValidationRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(decisionValidationRegistryPath)), &registry); err != nil {
		return DecisionGateAuditSummary{}, err
	}
	if err := validateDecisionValidationRegistry(registry, decisions, research); err != nil {
		return DecisionGateAuditSummary{}, err
	}
	summary := DecisionGateAuditSummary{Validations: len(registry.Validations)}
	for _, validation := range registry.Validations {
		summary.Candidates += len(validation.CandidateTradeoffs)
		summary.Prototypes += len(validation.Prototypes)
		summary.AdversarialReviews += len(validation.AdversarialReviews)
		for _, benchmark := range validation.Benchmarks {
			if benchmark.Representative {
				summary.RepresentativeBenchmarks++
			}
		}
	}
	for _, decision := range decisions.Decisions {
		if decision.Status == "accepted" && decision.Class != "A" {
			summary.AcceptedEvaluated++
		}
	}
	return summary, nil
}

func validateDecisionValidationRegistry(registry DecisionValidationRegistry, decisions DecisionRegistry, research ResearchRegistry) error {
	if registry.SchemaVersion != 1 || registry.RegistryKind != "engineering-decision-validation" {
		return fmt.Errorf("invalid decision validation registry header")
	}
	decisionByID := map[string]EngineeringDecision{}
	for _, decision := range decisions.Decisions {
		decisionByID[decision.ID] = decision
	}
	familiesByDecision := map[string]map[string]bool{}
	packetsByDecision := map[string]map[string]bool{}
	for _, packet := range research.Packets {
		if familiesByDecision[packet.DecisionID] == nil {
			familiesByDecision[packet.DecisionID] = map[string]bool{}
			packetsByDecision[packet.DecisionID] = map[string]bool{}
		}
		packetsByDecision[packet.DecisionID][packet.ID] = true
		for _, family := range packet.SolutionFamilies {
			familiesByDecision[packet.DecisionID][family.ID] = true
		}
	}

	byDecision := map[string]DecisionValidation{}
	seenIDs := map[string]bool{}
	for _, validation := range registry.Validations {
		if !decisionValidationIDPattern.MatchString(validation.ID) || seenIDs[validation.ID] {
			return fmt.Errorf("invalid or duplicate decision validation id %q", validation.ID)
		}
		seenIDs[validation.ID] = true
		decision, ok := decisionByID[validation.DecisionID]
		if !ok {
			return fmt.Errorf("decision validation %s references unknown decision %s", validation.ID, validation.DecisionID)
		}
		if decision.Class == "A" {
			return fmt.Errorf("decision validation %s cannot attach to Class A decision", validation.ID)
		}
		if _, duplicate := byDecision[decision.ID]; duplicate {
			return fmt.Errorf("decision %s has more than one validation record", decision.ID)
		}
		if err := validateDecisionValidation(validation, decision, familiesByDecision[decision.ID], packetsByDecision[decision.ID]); err != nil {
			return err
		}
		byDecision[decision.ID] = validation
	}
	for _, decision := range decisions.Decisions {
		if decision.Status != "accepted" || decision.Class == "A" {
			continue
		}
		if _, ok := byDecision[decision.ID]; !ok {
			return fmt.Errorf("accepted Class %s decision %s requires a decision validation record", decision.Class, decision.ID)
		}
	}
	return nil
}

func validateDecisionValidation(v DecisionValidation, decision EngineeringDecision, families, packets map[string]bool) error {
	candidates := map[string]CandidateTradeoff{}
	for _, candidate := range v.CandidateTradeoffs {
		if !candidateIDPattern.MatchString(candidate.ID) || candidates[candidate.ID].ID != "" {
			return fmt.Errorf("validation %s has invalid or duplicate candidate id %q", v.ID, candidate.ID)
		}
		if strings.TrimSpace(candidate.Name) == "" || !families[candidate.FamilyRef] {
			return fmt.Errorf("validation %s candidate %s has invalid name/family_ref", v.ID, candidate.ID)
		}
		if len(candidate.Strengths)+len(candidate.Weaknesses) == 0 {
			return fmt.Errorf("validation %s candidate %s lacks trade-off evidence", v.ID, candidate.ID)
		}
		if err := validatePacketRefs(v.ID+" candidate "+candidate.ID, candidate.EvidenceRefs, packets); err != nil {
			return err
		}
		candidates[candidate.ID] = candidate
	}
	if v.PreferredCandidateID != "" {
		if _, ok := candidates[v.PreferredCandidateID]; !ok {
			return fmt.Errorf("validation %s preferred candidate %s is not registered", v.ID, v.PreferredCandidateID)
		}
	}

	benchmarks := map[string]DecisionBenchmark{}
	for _, benchmark := range v.Benchmarks {
		if !benchmarkIDPattern.MatchString(benchmark.ID) || benchmarks[benchmark.ID].ID != "" {
			return fmt.Errorf("validation %s has invalid or duplicate benchmark id %q", v.ID, benchmark.ID)
		}
		if _, err := parseResearchDate(benchmark.RunAt); err != nil {
			return fmt.Errorf("validation %s benchmark %s run_at: %w", v.ID, benchmark.ID, err)
		}
		if strings.TrimSpace(benchmark.Workload) == "" || strings.TrimSpace(benchmark.Environment) == "" ||
			strings.TrimSpace(benchmark.DatasetScale) == "" || strings.TrimSpace(benchmark.ArtifactRef) == "" {
			return fmt.Errorf("validation %s benchmark %s has incomplete representative context", v.ID, benchmark.ID)
		}
		if err := validateCandidateRefs(v.ID+" benchmark "+benchmark.ID, benchmark.ComparedCandidateIDs, candidates); err != nil {
			return err
		}
		if len(benchmark.Results) == 0 {
			return fmt.Errorf("validation %s benchmark %s has no results", v.ID, benchmark.ID)
		}
		compared := map[string]bool{}
		for _, id := range benchmark.ComparedCandidateIDs {
			compared[id] = true
		}
		resultCandidates := map[string]bool{}
		for _, result := range benchmark.Results {
			if !compared[result.CandidateID] || strings.TrimSpace(result.Metric) == "" || strings.TrimSpace(result.Unit) == "" ||
				math.IsNaN(result.Value) || math.IsInf(result.Value, 0) {
				return fmt.Errorf("validation %s benchmark %s contains invalid result", v.ID, benchmark.ID)
			}
			resultCandidates[result.CandidateID] = true
		}
		for id := range compared {
			if !resultCandidates[id] {
				return fmt.Errorf("validation %s benchmark %s has no result for candidate %s", v.ID, benchmark.ID, id)
			}
		}
		benchmarks[benchmark.ID] = benchmark
	}

	prototypes := map[string]DecisionPrototype{}
	for _, prototype := range v.Prototypes {
		if !prototypeIDPattern.MatchString(prototype.ID) || prototypes[prototype.ID].ID != "" {
			return fmt.Errorf("validation %s has invalid or duplicate prototype id %q", v.ID, prototype.ID)
		}
		if _, ok := candidates[prototype.CandidateID]; !ok || !knownPrototypeOutcomes[prototype.Outcome] ||
			strings.TrimSpace(prototype.Question) == "" || strings.TrimSpace(prototype.ArtifactRef) == "" {
			return fmt.Errorf("validation %s prototype %s has invalid metadata", v.ID, prototype.ID)
		}
		if _, err := parseResearchDate(prototype.RunAt); err != nil {
			return fmt.Errorf("validation %s prototype %s run_at: %w", v.ID, prototype.ID, err)
		}
		if err := uniqueNonEmptyStrings(v.ID+" prototype unresolved uncertainties", prototype.UnresolvedUncertainties); err != nil {
			return err
		}
		prototypes[prototype.ID] = prototype
	}

	adversarial := map[string]AdversarialReview{}
	for _, review := range v.AdversarialReviews {
		if !adversarialIDPattern.MatchString(review.ID) || adversarial[review.ID].ID != "" {
			return fmt.Errorf("validation %s has invalid or duplicate adversarial id %q", v.ID, review.ID)
		}
		if _, ok := candidates[review.TargetCandidateID]; !ok || !knownAdversarialOutcomes[review.Outcome] ||
			strings.TrimSpace(review.Hypothesis) == "" || len(review.FalsificationAttempts) == 0 {
			return fmt.Errorf("validation %s adversarial review %s has invalid metadata", v.ID, review.ID)
		}
		if _, err := parseResearchDate(review.RunAt); err != nil {
			return fmt.Errorf("validation %s adversarial review %s run_at: %w", v.ID, review.ID, err)
		}
		if err := uniqueNonEmptyStrings(v.ID+" adversarial attempts", review.FalsificationAttempts); err != nil {
			return err
		}
		if err := uniqueNonEmptyStrings(v.ID+" adversarial counter evidence", review.CounterEvidence); err != nil {
			return err
		}
		if err := uniqueNonEmptyStrings(v.ID+" adversarial blocking findings", review.BlockingFindings); err != nil {
			return err
		}
		if err := validatePacketRefs(v.ID+" adversarial "+review.ID, review.EvidenceRefs, packets); err != nil {
			return err
		}
		adversarial[review.ID] = review
	}

	if decision.Status != "accepted" {
		return nil
	}
	if len(candidates) < 2 {
		return fmt.Errorf("accepted Class %s decision %s requires at least two compared candidates", decision.Class, decision.ID)
	}
	if v.PreferredCandidateID == "" {
		return fmt.Errorf("accepted decision %s requires a preferred candidate", decision.ID)
	}
	if len(v.BlockingCounterEvidence) != 0 || len(v.UnresolvedMaterialUncertainty) != 0 {
		return fmt.Errorf("accepted decision %s still has blocking counter-evidence or unresolved material uncertainty", decision.ID)
	}
	if !v.ExternalEvidenceSufficient {
		prototypeOK := false
		for _, prototype := range prototypes {
			if prototype.CandidateID == v.PreferredCandidateID && prototype.Outcome != "inconclusive" && len(prototype.UnresolvedUncertainties) == 0 {
				prototypeOK = true
			}
		}
		if !prototypeOK {
			return fmt.Errorf("accepted decision %s requires a conclusive preferred-candidate prototype because external evidence is insufficient", decision.ID)
		}
	}
	if decision.PerformanceSensitive {
		if len(decision.BenchmarkRefs) == 0 {
			return fmt.Errorf("accepted performance-sensitive decision %s has no benchmark_refs", decision.ID)
		}
		for _, ref := range decision.BenchmarkRefs {
			benchmark, ok := benchmarks[ref]
			if !ok || !benchmark.Representative {
				return fmt.Errorf("accepted performance-sensitive decision %s benchmark_ref %s is missing or non-representative", decision.ID, ref)
			}
			if !containsString(benchmark.ComparedCandidateIDs, v.PreferredCandidateID) || len(benchmark.ComparedCandidateIDs) < 2 {
				return fmt.Errorf("decision %s benchmark %s must compare preferred candidate with at least one alternative", decision.ID, ref)
			}
		}
	}
	if decision.Class == "C" {
		if len(decision.AdversarialRefs) == 0 {
			return fmt.Errorf("accepted Class C decision %s has no adversarial_refs", decision.ID)
		}
		for _, ref := range decision.AdversarialRefs {
			review, ok := adversarial[ref]
			if !ok {
				return fmt.Errorf("accepted Class C decision %s adversarial_ref %s is not registered", decision.ID, ref)
			}
			if review.TargetCandidateID != v.PreferredCandidateID || review.Outcome != "survived" ||
				len(review.BlockingFindings) != 0 || len(review.EvidenceRefs) == 0 {
				return fmt.Errorf("accepted Class C decision %s adversarial review %s did not survive cleanly", decision.ID, ref)
			}
		}
	}
	return nil
}

func validateCandidateRefs(label string, refs []string, candidates map[string]CandidateTradeoff) error {
	if len(refs) == 0 {
		return fmt.Errorf("%s requires candidate refs", label)
	}
	if err := uniqueNonEmptyStrings(label+" candidate refs", refs); err != nil {
		return err
	}
	for _, ref := range refs {
		if _, ok := candidates[ref]; !ok {
			return fmt.Errorf("%s references unknown candidate %s", label, ref)
		}
	}
	return nil
}

func validatePacketRefs(label string, refs []string, packets map[string]bool) error {
	if len(refs) == 0 {
		return fmt.Errorf("%s requires research packet evidence", label)
	}
	if err := uniqueNonEmptyStrings(label+" evidence refs", refs); err != nil {
		return err
	}
	for _, ref := range refs {
		if !packets[ref] {
			return fmt.Errorf("%s references research packet %s not owned by the decision", label, ref)
		}
	}
	return nil
}
