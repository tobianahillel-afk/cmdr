package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	freshnessPolicyPath   = "engineering/research/freshness-policy.json"
	freshnessRegistryPath = "engineering/research/freshness-registry.json"
)

type FreshnessPolicy struct {
	SchemaVersion     int      `json:"schema_version"`
	DefaultPolicy     string   `json:"default_policy"`
	AutomaticTriggers []string `json:"automatic_triggers"`
	ExternalTriggers  []string `json:"external_triggers"`
	TrackedInputKinds []string `json:"tracked_input_kinds"`
	SignalStates      []string `json:"signal_states"`
}

type EvidenceBinding struct {
	ID     string `json:"id"`
	Digest string `json:"sha256"`
}

type TrackedInputBinding struct {
	Kind   string `json:"kind"`
	Path   string `json:"path"`
	Digest string `json:"sha256"`
}

type FreshnessSnapshot struct {
	ProductSpecBaseline string                `json:"product_spec_baseline"`
	Decision             EvidenceBinding       `json:"decision"`
	ResearchPackets      []EvidenceBinding     `json:"research_packets"`
	Validation           EvidenceBinding       `json:"validation"`
	TrackedInputs        []TrackedInputBinding `json:"tracked_inputs"`
}

type RevisitSignal struct {
	Kind         string   `json:"kind"`
	State        string   `json:"state"`
	ObservedAt   string   `json:"observed_at"`
	EvidenceRefs []string `json:"evidence_refs"`
	Detail       string   `json:"detail"`
}

type DecisionFreshnessRecord struct {
	ID                 string            `json:"id"`
	DecisionID         string            `json:"decision_id"`
	ReviewedAt         string            `json:"reviewed_at"`
	FreshUntil         string            `json:"fresh_until"`
	FreshnessRationale string            `json:"freshness_rationale"`
	Snapshot           FreshnessSnapshot `json:"snapshot"`
	RevisitTriggers    []string          `json:"revisit_triggers"`
	Signals            []RevisitSignal   `json:"signals"`
}

type FreshnessRegistry struct {
	SchemaVersion int                       `json:"schema_version"`
	RegistryKind  string                    `json:"registry_kind"`
	Records       []DecisionFreshnessRecord `json:"records"`
}

type DecisionFreshnessAuditSummary struct {
	AsOf              string         `json:"as_of"`
	Records           int            `json:"records"`
	AcceptedEvaluated int            `json:"accepted_evaluated"`
	Fresh             int            `json:"fresh"`
	Stale             int            `json:"stale"`
	Reusable          int            `json:"reusable"`
	ByTrigger         map[string]int `json:"by_trigger"`
}

type ReusableEvidenceEntry struct {
	DecisionID        string   `json:"decision_id"`
	DecisionKey       string   `json:"decision_key"`
	Class             string   `json:"class"`
	FreshUntil        string   `json:"fresh_until"`
	BasisDigest       string   `json:"basis_sha256"`
	ResearchPacketIDs []string `json:"research_packet_ids"`
	ValidationID      string   `json:"validation_id"`
}

type ReusableEvidenceCache struct {
	SchemaVersion int                     `json:"schema_version"`
	AsOf          string                  `json:"as_of"`
	Entries       []ReusableEvidenceEntry `json:"entries"`
	Digest        string                  `json:"sha256"`
}

var freshnessRecordIDPattern = regexp.MustCompile(`^FRESH-[0-9]{4,}$`)
var sha256Pattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

var knownFreshnessAutomaticTriggers = map[string]bool{
	"date-expiry": true,
	"decision-change": true,
	"research-evidence-change": true,
	"validation-evidence-change": true,
	"repository-input-change": true,
	"product-spec-baseline-change": true,
}

var knownFreshnessExternalTriggers = map[string]bool{
	"upstream-version-change": true,
	"security-advisory": true,
	"threat-model-change": true,
	"benchmark-regression": true,
	"assumption-change": true,
	"material-new-research": true,
}

var knownFreshnessTrackedInputKinds = map[string]bool{
	"assumption": true,
	"threat-model": true,
	"benchmark-fixture": true,
	"constraint": true,
	"dependency-lock": true,
	"architecture-contract": true,
	"security-policy": true,
}

var knownFreshnessSignalStates = map[string]bool{"clear": true, "fired": true}

func runDecisionFreshnessAudit(root, asOf string, state CurrentState, graph WorkGraph) (DecisionFreshnessAuditSummary, error) {
	summary, _, err := evaluateDecisionFreshness(root, asOf, state, graph)
	return summary, err
}

func runReusableEvidenceCache(root, asOf string, state CurrentState, graph WorkGraph) (ReusableEvidenceCache, error) {
	_, cache, err := evaluateDecisionFreshness(root, asOf, state, graph)
	return cache, err
}

func runDecisionFreshnessSnapshot(root, decisionID string, state CurrentState, graph WorkGraph) (FreshnessSnapshot, error) {
	if !decisionIDPattern.MatchString(decisionID) {
		return FreshnessSnapshot{}, fmt.Errorf("invalid decision id %q", decisionID)
	}
	decisions, research, validations, err := loadValidatedDecisionEvidence(root, graph)
	if err != nil {
		return FreshnessSnapshot{}, err
	}
	for _, decision := range decisions.Decisions {
		if decision.ID == decisionID {
			if decision.Class == "A" {
				return FreshnessSnapshot{}, fmt.Errorf("Class A decision %s does not require a freshness snapshot", decisionID)
			}
			return buildFreshnessSnapshot(root, decision, research, validations, state.ProductSpec.BaselineCommit, nil)
		}
	}
	return FreshnessSnapshot{}, fmt.Errorf("decision %s not found", decisionID)
}

func evaluateDecisionFreshness(root, asOf string, state CurrentState, graph WorkGraph) (DecisionFreshnessAuditSummary, ReusableEvidenceCache, error) {
	asOfDate, err := parseFreshnessAsOf(asOf)
	if err != nil {
		return DecisionFreshnessAuditSummary{}, ReusableEvidenceCache{}, err
	}
	policy, registry, err := loadFreshnessRegistry(root)
	if err != nil {
		return DecisionFreshnessAuditSummary{}, ReusableEvidenceCache{}, err
	}
	decisions, research, validations, err := loadValidatedDecisionEvidence(root, graph)
	if err != nil {
		return DecisionFreshnessAuditSummary{}, ReusableEvidenceCache{}, err
	}
	if err := validateFreshnessRegistry(policy, registry, decisions); err != nil {
		return DecisionFreshnessAuditSummary{}, ReusableEvidenceCache{}, err
	}

	decisionByID := map[string]EngineeringDecision{}
	for _, decision := range decisions.Decisions {
		decisionByID[decision.ID] = decision
	}
	recordByDecision := map[string]DecisionFreshnessRecord{}
	for _, record := range registry.Records {
		recordByDecision[record.DecisionID] = record
	}

	summary := DecisionFreshnessAuditSummary{AsOf: asOfDate.Format("2006-01-02"), Records: len(registry.Records), ByTrigger: map[string]int{}}
	cache := ReusableEvidenceCache{SchemaVersion: 1, AsOf: summary.AsOf}
	for _, decision := range decisions.Decisions {
		if decision.Class == "A" || decision.Status != "accepted" {
			continue
		}
		summary.AcceptedEvaluated++
		record, ok := recordByDecision[decision.ID]
		if !ok {
			return summary, cache, fmt.Errorf("accepted Class %s decision %s requires a freshness record", decision.Class, decision.ID)
		}
		fresh, reasons, currentSnapshot, err := evaluateFreshnessRecord(root, asOfDate, policy, record, decision, research, validations, state)
		if err != nil {
			return summary, cache, err
		}
		if !fresh {
			summary.Stale++
			for _, reason := range reasons {
				summary.ByTrigger[reason]++
			}
			return summary, cache, fmt.Errorf("accepted decision %s is stale (%s); move it to revisit-required or refresh its evidence before reuse", decision.ID, strings.Join(reasons, ", "))
		}
		summary.Fresh++
		summary.Reusable++
		var packetIDs []string
		for _, binding := range currentSnapshot.ResearchPackets {
			packetIDs = append(packetIDs, binding.ID)
		}
		cache.Entries = append(cache.Entries, ReusableEvidenceEntry{
			DecisionID: decision.ID,
			DecisionKey: decision.DecisionKey,
			Class: decision.Class,
			FreshUntil: record.FreshUntil,
			BasisDigest: digestCanonical(currentSnapshot),
			ResearchPacketIDs: packetIDs,
			ValidationID: currentSnapshot.Validation.ID,
		})
	}
	sort.Slice(cache.Entries, func(i, j int) bool { return cache.Entries[i].DecisionID < cache.Entries[j].DecisionID })
	cache.Digest = digestCanonical(struct {
		SchemaVersion int                     `json:"schema_version"`
		AsOf          string                  `json:"as_of"`
		Entries       []ReusableEvidenceEntry `json:"entries"`
	}{cache.SchemaVersion, cache.AsOf, cache.Entries})
	return summary, cache, nil
}

func loadValidatedDecisionEvidence(root string, graph WorkGraph) (DecisionRegistry, ResearchRegistry, DecisionValidationRegistry, error) {
	if _, err := runDecisionRegistryAudit(root, graph); err != nil {
		return DecisionRegistry{}, ResearchRegistry{}, DecisionValidationRegistry{}, err
	}
	if _, err := runResearchPacketAudit(root); err != nil {
		return DecisionRegistry{}, ResearchRegistry{}, DecisionValidationRegistry{}, err
	}
	if _, err := runDecisionGateAudit(root); err != nil {
		return DecisionRegistry{}, ResearchRegistry{}, DecisionValidationRegistry{}, err
	}
	var decisions DecisionRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(decisionRegistryPath)), &decisions); err != nil {
		return decisions, ResearchRegistry{}, DecisionValidationRegistry{}, err
	}
	var research ResearchRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(researchRegistryPath)), &research); err != nil {
		return decisions, research, DecisionValidationRegistry{}, err
	}
	var validations DecisionValidationRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(decisionValidationRegistryPath)), &validations); err != nil {
		return decisions, research, validations, err
	}
	return decisions, research, validations, nil
}

func loadFreshnessRegistry(root string) (FreshnessPolicy, FreshnessRegistry, error) {
	var policy FreshnessPolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(freshnessPolicyPath)), &policy); err != nil {
		return policy, FreshnessRegistry{}, err
	}
	if err := validateFreshnessPolicy(policy); err != nil {
		return policy, FreshnessRegistry{}, err
	}
	var registry FreshnessRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(freshnessRegistryPath)), &registry); err != nil {
		return policy, registry, err
	}
	return policy, registry, nil
}

func validateFreshnessPolicy(policy FreshnessPolicy) error {
	if policy.SchemaVersion != 1 || policy.DefaultPolicy != "deny-stale-accepted-decision" {
		return fmt.Errorf("invalid freshness policy header")
	}
	if err := exactStringSet("automatic freshness trigger", policy.AutomaticTriggers, knownFreshnessAutomaticTriggers); err != nil {
		return err
	}
	if err := exactStringSet("external freshness trigger", policy.ExternalTriggers, knownFreshnessExternalTriggers); err != nil {
		return err
	}
	if err := exactStringSet("freshness tracked input kind", policy.TrackedInputKinds, knownFreshnessTrackedInputKinds); err != nil {
		return err
	}
	return exactStringSet("freshness signal state", policy.SignalStates, knownFreshnessSignalStates)
}

func validateFreshnessRegistry(policy FreshnessPolicy, registry FreshnessRegistry, decisions DecisionRegistry) error {
	if registry.SchemaVersion != 1 || registry.RegistryKind != "engineering-decision-freshness" {
		return fmt.Errorf("invalid freshness registry header")
	}
	decisionByID := map[string]EngineeringDecision{}
	for _, decision := range decisions.Decisions {
		decisionByID[decision.ID] = decision
	}
	seenIDs := map[string]bool{}
	seenDecision := map[string]bool{}
	for _, record := range registry.Records {
		if !freshnessRecordIDPattern.MatchString(record.ID) || seenIDs[record.ID] {
			return fmt.Errorf("invalid or duplicate freshness record id %q", record.ID)
		}
		seenIDs[record.ID] = true
		decision, ok := decisionByID[record.DecisionID]
		if !ok || decision.Class == "A" {
			return fmt.Errorf("freshness record %s references unknown/Class A decision %s", record.ID, record.DecisionID)
		}
		if seenDecision[record.DecisionID] {
			return fmt.Errorf("decision %s has multiple freshness records", record.DecisionID)
		}
		seenDecision[record.DecisionID] = true
		if strings.TrimSpace(record.FreshnessRationale) == "" {
			return fmt.Errorf("freshness record %s requires freshness_rationale", record.ID)
		}
		reviewed, err := parseResearchDate(record.ReviewedAt)
		if err != nil {
			return fmt.Errorf("freshness record %s reviewed_at: %w", record.ID, err)
		}
		freshUntil, err := parseResearchDate(record.FreshUntil)
		if err != nil || freshUntil.Before(reviewed) {
			return fmt.Errorf("freshness record %s has invalid fresh_until", record.ID)
		}
		if !sha256Pattern.MatchString(record.Snapshot.Decision.Digest) {
			return fmt.Errorf("freshness record %s has invalid decision digest", record.ID)
		}
		if strings.TrimSpace(record.Snapshot.ProductSpecBaseline) == "" {
			return fmt.Errorf("freshness record %s requires product_spec_baseline", record.ID)
		}
		if err := validateEvidenceBindings(record.ID+" research", record.Snapshot.ResearchPackets); err != nil {
			return err
		}
		if record.Snapshot.Validation.ID == "" || !sha256Pattern.MatchString(record.Snapshot.Validation.Digest) {
			return fmt.Errorf("freshness record %s requires validation binding", record.ID)
		}
		if err := validateTrackedInputBindings(record.ID, record.Snapshot.TrackedInputs); err != nil {
			return err
		}
		if err := validateFreshnessTriggers(record, decision, policy); err != nil {
			return err
		}
	}
	return nil
}

func validateFreshnessTriggers(record DecisionFreshnessRecord, decision EngineeringDecision, policy FreshnessPolicy) error {
	known := map[string]bool{}
	for trigger := range knownFreshnessAutomaticTriggers { known[trigger] = true }
	for trigger := range knownFreshnessExternalTriggers { known[trigger] = true }
	if err := validateKnownUniqueStrings(record.ID+" revisit trigger", record.RevisitTriggers, known); err != nil {
		return err
	}
	for trigger := range knownFreshnessAutomaticTriggers {
		if !containsString(record.RevisitTriggers, trigger) {
			return fmt.Errorf("freshness record %s must declare automatic trigger %s", record.ID, trigger)
		}
	}
	requiredExternal := map[string]bool{}
	if decision.Class == "C" { requiredExternal["material-new-research"] = true }
	if decision.PerformanceSensitive { requiredExternal["benchmark-regression"] = true }
	if containsString(decision.RiskTags, "security") || containsString(decision.CriticalFactors, "security-critical") {
		requiredExternal["security-advisory"] = true
		requiredExternal["threat-model-change"] = true
	}
	if containsString(decision.RiskTags, "dependency") { requiredExternal["upstream-version-change"] = true }
	for _, input := range record.Snapshot.TrackedInputs {
		switch input.Kind {
		case "assumption":
			requiredExternal["assumption-change"] = true
		case "threat-model":
			requiredExternal["threat-model-change"] = true
		case "benchmark-fixture":
			requiredExternal["benchmark-regression"] = true
		}
	}
	for trigger := range requiredExternal {
		if !containsString(record.RevisitTriggers, trigger) {
			return fmt.Errorf("freshness record %s must declare required external trigger %s", record.ID, trigger)
		}
	}

	signalByKind := map[string]RevisitSignal{}
	for _, signal := range record.Signals {
		if !knownFreshnessExternalTriggers[signal.Kind] || !knownFreshnessSignalStates[signal.State] {
			return fmt.Errorf("freshness record %s has invalid signal %s/%s", record.ID, signal.Kind, signal.State)
		}
		if _, duplicate := signalByKind[signal.Kind]; duplicate {
			return fmt.Errorf("freshness record %s duplicates signal %s", record.ID, signal.Kind)
		}
		if _, err := parseResearchDate(signal.ObservedAt); err != nil {
			return fmt.Errorf("freshness record %s signal %s observed_at: %w", record.ID, signal.Kind, err)
		}
		if strings.TrimSpace(signal.Detail) == "" || len(signal.EvidenceRefs) == 0 {
			return fmt.Errorf("freshness record %s signal %s requires detail and evidence_refs", record.ID, signal.Kind)
		}
		if err := uniqueNonEmptyStrings(record.ID+" signal evidence", signal.EvidenceRefs); err != nil {
			return err
		}
		signalByKind[signal.Kind] = signal
	}
	for _, trigger := range record.RevisitTriggers {
		if knownFreshnessExternalTriggers[trigger] {
			if _, ok := signalByKind[trigger]; !ok {
				return fmt.Errorf("freshness record %s external trigger %s requires a signal observation", record.ID, trigger)
			}
		}
	}
	return nil
}

func evaluateFreshnessRecord(root string, asOf time.Time, policy FreshnessPolicy, record DecisionFreshnessRecord, decision EngineeringDecision, research ResearchRegistry, validations DecisionValidationRegistry, state CurrentState) (bool, []string, FreshnessSnapshot, error) {
	_ = policy
	reviewed, _ := parseResearchDate(record.ReviewedAt)
	freshUntil, _ := parseResearchDate(record.FreshUntil)
	if reviewed.After(asOf) {
		return false, nil, FreshnessSnapshot{}, fmt.Errorf("freshness record %s reviewed_at is in the future relative to %s", record.ID, asOf.Format("2006-01-02"))
	}
	current, err := buildFreshnessSnapshot(root, decision, research, validations, state.ProductSpec.BaselineCommit, record.Snapshot.TrackedInputs)
	if err != nil {
		return false, nil, current, err
	}
	var reasons []string
	if asOf.After(freshUntil) { reasons = append(reasons, "date-expiry") }
	if digestCanonical(current) != digestCanonical(record.Snapshot) {
		reasons = append(reasons, classifySnapshotDrift(record.Snapshot, current)...)
	}
	for _, signal := range record.Signals {
		observed, _ := parseResearchDate(signal.ObservedAt)
		if observed.Before(reviewed) {
			return false, nil, current, fmt.Errorf("freshness record %s signal %s predates reviewed_at", record.ID, signal.Kind)
		}
		if observed.After(asOf) {
			return false, nil, current, fmt.Errorf("freshness record %s signal %s is observed in the future", record.ID, signal.Kind)
		}
		if signal.State == "fired" {
			reasons = append(reasons, signal.Kind)
		}
	}
	reasons = uniqueSorted(reasons)
	return len(reasons) == 0, reasons, current, nil
}

func buildFreshnessSnapshot(root string, decision EngineeringDecision, research ResearchRegistry, validations DecisionValidationRegistry, productBaseline string, tracked []TrackedInputBinding) (FreshnessSnapshot, error) {
	snapshot := FreshnessSnapshot{
		ProductSpecBaseline: productBaseline,
		Decision: EvidenceBinding{ID: decision.ID, Digest: digestCanonical(decision)},
	}
	packetByID := map[string]ResearchPacket{}
	for _, packet := range research.Packets { packetByID[packet.ID] = packet }
	for _, ref := range decision.EvidenceRefs {
		packet, ok := packetByID[ref]
		if !ok { return snapshot, fmt.Errorf("decision %s evidence packet %s is missing", decision.ID, ref) }
		snapshot.ResearchPackets = append(snapshot.ResearchPackets, EvidenceBinding{ID: ref, Digest: digestCanonical(packet)})
	}
	sort.Slice(snapshot.ResearchPackets, func(i, j int) bool { return snapshot.ResearchPackets[i].ID < snapshot.ResearchPackets[j].ID })

	for _, validation := range validations.Validations {
		if validation.DecisionID == decision.ID {
			snapshot.Validation = EvidenceBinding{ID: validation.ID, Digest: digestCanonical(validation)}
			break
		}
	}
	if decision.Status == "accepted" && decision.Class != "A" && snapshot.Validation.ID == "" {
		return snapshot, fmt.Errorf("accepted decision %s has no validation record", decision.ID)
	}

	for _, input := range tracked {
		if !knownFreshnessTrackedInputKinds[input.Kind] {
			return snapshot, fmt.Errorf("decision %s has unknown tracked input kind %s", decision.ID, input.Kind)
		}
		data, err := readRepoFile(root, input.Path)
		if err != nil { return snapshot, fmt.Errorf("tracked input %s: %w", input.Path, err) }
		sum := sha256.Sum256(data)
		snapshot.TrackedInputs = append(snapshot.TrackedInputs, TrackedInputBinding{
			Kind: input.Kind, Path: input.Path, Digest: hex.EncodeToString(sum[:]),
		})
	}
	sort.Slice(snapshot.TrackedInputs, func(i, j int) bool {
		if snapshot.TrackedInputs[i].Kind != snapshot.TrackedInputs[j].Kind { return snapshot.TrackedInputs[i].Kind < snapshot.TrackedInputs[j].Kind }
		return snapshot.TrackedInputs[i].Path < snapshot.TrackedInputs[j].Path
	})
	return snapshot, nil
}

func validateEvidenceBindings(label string, bindings []EvidenceBinding) error {
	seen := map[string]bool{}
	for _, binding := range bindings {
		if strings.TrimSpace(binding.ID) == "" || !sha256Pattern.MatchString(binding.Digest) {
			return fmt.Errorf("%s contains invalid binding", label)
		}
		if seen[binding.ID] { return fmt.Errorf("%s duplicates binding %s", label, binding.ID) }
		seen[binding.ID] = true
	}
	return nil
}

func validateTrackedInputBindings(label string, inputs []TrackedInputBinding) error {
	seen := map[string]bool{}
	for _, input := range inputs {
		if !knownFreshnessTrackedInputKinds[input.Kind] || !sha256Pattern.MatchString(input.Digest) {
			return fmt.Errorf("%s has invalid tracked input %s", label, input.Path)
		}
		normalized, err := normalizeChangedPaths([]string{input.Path})
		if err != nil || len(normalized) != 1 || normalized[0] != input.Path {
			return fmt.Errorf("%s has non-canonical tracked input path %q", label, input.Path)
		}
		key := input.Kind + "\x00" + input.Path
		if seen[key] { return fmt.Errorf("%s duplicates tracked input %s", label, input.Path) }
		seen[key] = true
	}
	return nil
}

func classifySnapshotDrift(old, current FreshnessSnapshot) []string {
	var reasons []string
	if old.ProductSpecBaseline != current.ProductSpecBaseline { reasons = append(reasons, "product-spec-baseline-change") }
	if old.Decision != current.Decision { reasons = append(reasons, "decision-change") }
	if digestCanonical(old.ResearchPackets) != digestCanonical(current.ResearchPackets) { reasons = append(reasons, "research-evidence-change") }
	if old.Validation != current.Validation { reasons = append(reasons, "validation-evidence-change") }
	if digestCanonical(old.TrackedInputs) != digestCanonical(current.TrackedInputs) { reasons = append(reasons, "repository-input-change") }
	return reasons
}

func digestCanonical(value any) string {
	data, err := json.Marshal(value)
	if err != nil { panic(err) }
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func parseFreshnessAsOf(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		now := time.Now().UTC()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC), nil
	}
	return parseResearchDate(value)
}

func uniqueSorted(values []string) []string {
	seen := map[string]bool{}
	for _, value := range values { if value != "" { seen[value] = true } }
	out := make([]string, 0, len(seen))
	for value := range seen { out = append(out, value) }
	sort.Strings(out)
	return out
}
