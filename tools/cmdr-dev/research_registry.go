package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	researchPolicyPath   = "engineering/research/research-policy.json"
	researchRegistryPath = "engineering/research/evidence-registry.json"
)

type ResearchLimits struct {
	MaxPackets          int `json:"max_packets"`
	MaxFamiliesPerPacket int `json:"max_families_per_packet"`
	MaxSourcesPerPacket int `json:"max_sources_per_packet"`
	MaxClaimsPerPacket  int `json:"max_claims_per_packet"`
	MaxLimitationsPerPacket int `json:"max_limitations_per_packet"`
}

type ResearchPolicy struct {
	SchemaVersion      int      `json:"schema_version"`
	DefaultPolicy      string   `json:"default_policy"`
	PacketStatuses     []string `json:"packet_statuses"`
	SourceKinds        []string `json:"source_kinds"`
	ClaimCriticalities []string `json:"claim_criticalities"`
	ClaimStatuses      []string `json:"claim_statuses"`
	FamilyDispositions []string `json:"family_dispositions"`
	Limits             ResearchLimits `json:"limits"`
}

type ResearchSolutionFamily struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Disposition string `json:"disposition"`
	Rationale   string `json:"rationale"`
}

type ResearchSource struct {
	ID                   string `json:"id"`
	Kind                 string `json:"kind"`
	Title                string `json:"title"`
	Publisher            string `json:"publisher"`
	Locator              string `json:"locator"`
	PublishedOrUpdatedAt string `json:"published_or_updated_at"`
	AccessedAt           string `json:"accessed_at"`
}

type ResearchClaim struct {
	ID             string   `json:"id"`
	Statement      string   `json:"statement"`
	Criticality    string   `json:"criticality"`
	Status         string   `json:"status"`
	SourceRefs     []string `json:"source_refs"`
	FamilyRefs     []string `json:"family_refs"`
	LimitationRefs []string `json:"limitation_refs"`
}

type ResearchLimitation struct {
	ID         string   `json:"id"`
	Statement  string   `json:"statement"`
	SourceRefs []string `json:"source_refs"`
}

type ResearchSaturation struct {
	MajorFamilyIDs                       []string `json:"major_family_ids"`
	MajorSolutionFamiliesCovered         bool     `json:"major_solution_families_covered"`
	NewSourcesMateriallyChangeCandidates bool     `json:"new_sources_materially_change_candidates"`
	CriticalLimitationsKnown             bool     `json:"critical_limitations_known"`
	UnresolvedCriticalClaimIDs           []string `json:"unresolved_critical_claim_ids"`
}

type ResearchPacket struct {
	ID               string                   `json:"id"`
	DecisionID       string                   `json:"decision_id"`
	Status           string                   `json:"status"`
	Topic            string                   `json:"topic"`
	Question         string                   `json:"question"`
	AsOf             string                   `json:"as_of"`
	Constraints      []string                 `json:"constraints"`
	SolutionFamilies []ResearchSolutionFamily `json:"solution_families"`
	Sources          []ResearchSource         `json:"sources"`
	Claims           []ResearchClaim          `json:"claims"`
	Limitations      []ResearchLimitation     `json:"limitations"`
	Saturation       ResearchSaturation       `json:"saturation"`
	SupersededBy     string                   `json:"superseded_by"`
}

type ResearchRegistry struct {
	SchemaVersion int              `json:"schema_version"`
	RegistryKind  string           `json:"registry_kind"`
	Packets       []ResearchPacket `json:"packets"`
}

type ResearchAuditSummary struct {
	Packets              int            `json:"packets"`
	Saturated            int            `json:"saturated"`
	Collecting           int            `json:"collecting"`
	Superseded           int            `json:"superseded"`
	Sources              int            `json:"sources"`
	Claims               int            `json:"claims"`
	CriticalClaims       int            `json:"critical_claims"`
	UnresolvedCritical   int            `json:"unresolved_critical_claims"`
	BySourceKind         map[string]int `json:"by_source_kind"`
}

type ResearchContextBundle struct {
	PacketID          string                   `json:"packet_id"`
	DecisionID        string                   `json:"decision_id"`
	Topic             string                   `json:"topic"`
	Question          string                   `json:"question"`
	AsOf              string                   `json:"as_of"`
	Constraints       []string                 `json:"constraints"`
	SolutionFamilies  []ResearchSolutionFamily `json:"solution_families"`
	Sources           []ResearchSource         `json:"sources"`
	Claims            []ResearchClaim          `json:"claims"`
	Limitations       []ResearchLimitation     `json:"limitations"`
	Saturation        ResearchSaturation       `json:"saturation"`
}

var researchPacketIDPattern = regexp.MustCompile(`^RES-PKT-[0-9]{4,}$`)
var researchFamilyIDPattern = regexp.MustCompile(`^FAMILY-[A-Z0-9][A-Z0-9-]{1,63}$`)
var researchSourceIDPattern = regexp.MustCompile(`^SRC-[0-9]{3,}$`)
var researchClaimIDPattern = regexp.MustCompile(`^CLAIM-[0-9]{3,}$`)
var researchLimitationIDPattern = regexp.MustCompile(`^LIMIT-[0-9]{3,}$`)

var knownResearchPacketStatuses = map[string]bool{"collecting": true, "saturated": true, "superseded": true}
var knownResearchSourceKinds = map[string]bool{
	"official-doc": true, "scientific-paper": true, "standard": true, "security-advisory": true,
	"source-code": true, "benchmark-report": true, "independent-analysis": true, "community-discussion": true,
}
var knownResearchClaimCriticalities = map[string]bool{"supporting": true, "material": true, "critical": true}
var knownResearchClaimStatuses = map[string]bool{"supported": true, "contested": true, "unresolved": true}
var knownResearchFamilyDispositions = map[string]bool{"candidate": true, "excluded": true}

func runResearchPacketAudit(root string) (ResearchAuditSummary, error) {
	policy, registry, decisions, err := loadResearchState(root)
	if err != nil {
		return ResearchAuditSummary{}, err
	}
	if err := validateResearchRegistry(policy, registry, decisions); err != nil {
		return ResearchAuditSummary{}, err
	}
	summary := ResearchAuditSummary{BySourceKind: map[string]int{}}
	summary.Packets = len(registry.Packets)
	for _, packet := range registry.Packets {
		switch packet.Status {
		case "collecting":
			summary.Collecting++
		case "saturated":
			summary.Saturated++
		case "superseded":
			summary.Superseded++
		}
		summary.Sources += len(packet.Sources)
		summary.Claims += len(packet.Claims)
		for _, source := range packet.Sources {
			summary.BySourceKind[source.Kind]++
		}
		for _, claim := range packet.Claims {
			if claim.Criticality == "critical" {
				summary.CriticalClaims++
				if claim.Status == "unresolved" {
					summary.UnresolvedCritical++
				}
			}
		}
	}
	return summary, nil
}

func compileResearchContext(root, packetID string) (ResearchContextBundle, error) {
	if !researchPacketIDPattern.MatchString(packetID) {
		return ResearchContextBundle{}, fmt.Errorf("invalid research packet id %q", packetID)
	}
	policy, registry, decisions, err := loadResearchState(root)
	if err != nil {
		return ResearchContextBundle{}, err
	}
	if err := validateResearchRegistry(policy, registry, decisions); err != nil {
		return ResearchContextBundle{}, err
	}
	for _, packet := range registry.Packets {
		if packet.ID == packetID {
			return ResearchContextBundle{
				PacketID: packet.ID, DecisionID: packet.DecisionID, Topic: packet.Topic, Question: packet.Question,
				AsOf: packet.AsOf, Constraints: append([]string(nil), packet.Constraints...),
				SolutionFamilies: append([]ResearchSolutionFamily(nil), packet.SolutionFamilies...),
				Sources: append([]ResearchSource(nil), packet.Sources...),
				Claims: append([]ResearchClaim(nil), packet.Claims...),
				Limitations: append([]ResearchLimitation(nil), packet.Limitations...),
				Saturation: packet.Saturation,
			}, nil
		}
	}
	return ResearchContextBundle{}, fmt.Errorf("research packet %s not found", packetID)
}

func loadResearchState(root string) (ResearchPolicy, ResearchRegistry, DecisionRegistry, error) {
	var policy ResearchPolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(researchPolicyPath)), &policy); err != nil {
		return policy, ResearchRegistry{}, DecisionRegistry{}, err
	}
	if err := validateResearchPolicy(policy); err != nil {
		return policy, ResearchRegistry{}, DecisionRegistry{}, err
	}
	var registry ResearchRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(researchRegistryPath)), &registry); err != nil {
		return policy, registry, DecisionRegistry{}, err
	}
	var decisions DecisionRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(decisionRegistryPath)), &decisions); err != nil {
		return policy, registry, decisions, err
	}
	return policy, registry, decisions, nil
}

func validateResearchPolicy(policy ResearchPolicy) error {
	if policy.SchemaVersion != 1 {
		return fmt.Errorf("unsupported research policy schema_version %d", policy.SchemaVersion)
	}
	if policy.DefaultPolicy != "deny-unsourced-material-claim" {
		return fmt.Errorf("research policy default must be deny-unsourced-material-claim")
	}
	if err := exactStringSet("research packet status", policy.PacketStatuses, knownResearchPacketStatuses); err != nil {
		return err
	}
	if err := exactStringSet("research source kind", policy.SourceKinds, knownResearchSourceKinds); err != nil {
		return err
	}
	if err := exactStringSet("research claim criticality", policy.ClaimCriticalities, knownResearchClaimCriticalities); err != nil {
		return err
	}
	if err := exactStringSet("research claim status", policy.ClaimStatuses, knownResearchClaimStatuses); err != nil {
		return err
	}
	if err := exactStringSet("research family disposition", policy.FamilyDispositions, knownResearchFamilyDispositions); err != nil {
		return err
	}
	if policy.Limits.MaxPackets < 1 || policy.Limits.MaxFamiliesPerPacket < 1 || policy.Limits.MaxSourcesPerPacket < 1 ||
		policy.Limits.MaxClaimsPerPacket < 1 || policy.Limits.MaxLimitationsPerPacket < 1 {
		return fmt.Errorf("research policy limits must all be positive")
	}
	return nil
}

func validateResearchRegistry(policy ResearchPolicy, registry ResearchRegistry, decisions DecisionRegistry) error {
	if registry.SchemaVersion != 1 || registry.RegistryKind != "engineering-research-evidence" {
		return fmt.Errorf("invalid research registry header")
	}
	if len(registry.Packets) > policy.Limits.MaxPackets {
		return fmt.Errorf("research registry exceeds max_packets")
	}
	decisionByID := map[string]EngineeringDecision{}
	for _, d := range decisions.Decisions {
		decisionByID[d.ID] = d
	}
	packetByID := map[string]ResearchPacket{}
	for _, packet := range registry.Packets {
		if !researchPacketIDPattern.MatchString(packet.ID) {
			return fmt.Errorf("invalid research packet id %q", packet.ID)
		}
		if _, exists := packetByID[packet.ID]; exists {
			return fmt.Errorf("duplicate research packet id %s", packet.ID)
		}
		decision, ok := decisionByID[packet.DecisionID]
		if !ok {
			return fmt.Errorf("research packet %s references unknown decision %s", packet.ID, packet.DecisionID)
		}
		if decision.Class == "A" {
			return fmt.Errorf("research packet %s cannot attach to Class A decision %s", packet.ID, decision.ID)
		}
		if !knownResearchPacketStatuses[packet.Status] {
			return fmt.Errorf("research packet %s has unknown status %q", packet.ID, packet.Status)
		}
		if strings.TrimSpace(packet.Topic) == "" || strings.TrimSpace(packet.Question) == "" || len(packet.Constraints) == 0 {
			return fmt.Errorf("research packet %s requires topic question and constraints", packet.ID)
		}
		if _, err := parseResearchDate(packet.AsOf); err != nil {
			return fmt.Errorf("research packet %s as_of: %w", packet.ID, err)
		}
		if len(packet.SolutionFamilies) > policy.Limits.MaxFamiliesPerPacket || len(packet.Sources) > policy.Limits.MaxSourcesPerPacket ||
			len(packet.Claims) > policy.Limits.MaxClaimsPerPacket || len(packet.Limitations) > policy.Limits.MaxLimitationsPerPacket {
			return fmt.Errorf("research packet %s exceeds bounded context limits", packet.ID)
		}
		if err := uniqueNonEmptyStrings(packet.ID+" constraints", packet.Constraints); err != nil {
			return err
		}
		if err := validateResearchPacket(packet); err != nil {
			return err
		}
		packetByID[packet.ID] = packet
	}
	for _, packet := range registry.Packets {
		if packet.Status == "superseded" {
			next, ok := packetByID[packet.SupersededBy]
			if !ok {
				return fmt.Errorf("superseded research packet %s references unknown packet %s", packet.ID, packet.SupersededBy)
			}
			if next.DecisionID != packet.DecisionID || next.Status == "superseded" {
				return fmt.Errorf("invalid research packet supersession %s -> %s", packet.ID, next.ID)
			}
		}
	}
	for _, decision := range decisions.Decisions {
		if decision.Status != "accepted" || decision.Class == "A" {
			continue
		}
		for _, ref := range decision.EvidenceRefs {
			packet, ok := packetByID[ref]
			if !ok {
				return fmt.Errorf("accepted decision %s evidence_ref %s is not a registered research packet", decision.ID, ref)
			}
			if packet.DecisionID != decision.ID || packet.Status != "saturated" {
				return fmt.Errorf("accepted decision %s evidence_ref %s must be its own saturated research packet", decision.ID, ref)
			}
		}
	}
	return nil
}

func validateResearchPacket(packet ResearchPacket) error {
	if packet.Status == "superseded" {
		if strings.TrimSpace(packet.SupersededBy) == "" {
			return fmt.Errorf("superseded research packet %s requires superseded_by", packet.ID)
		}
	} else if packet.SupersededBy != "" {
		return fmt.Errorf("research packet %s has superseded_by but status is %s", packet.ID, packet.Status)
	}

	families := map[string]ResearchSolutionFamily{}
	for _, family := range packet.SolutionFamilies {
		if !researchFamilyIDPattern.MatchString(family.ID) || families[family.ID].ID != "" {
			return fmt.Errorf("research packet %s has invalid or duplicate family id %q", packet.ID, family.ID)
		}
		if strings.TrimSpace(family.Name) == "" || strings.TrimSpace(family.Description) == "" ||
			strings.TrimSpace(family.Rationale) == "" || !knownResearchFamilyDispositions[family.Disposition] {
			return fmt.Errorf("research packet %s family %s has incomplete metadata", packet.ID, family.ID)
		}
		families[family.ID] = family
	}

	sources := map[string]ResearchSource{}
	for _, source := range packet.Sources {
		if !researchSourceIDPattern.MatchString(source.ID) || sources[source.ID].ID != "" {
			return fmt.Errorf("research packet %s has invalid or duplicate source id %q", packet.ID, source.ID)
		}
		if !knownResearchSourceKinds[source.Kind] || strings.TrimSpace(source.Title) == "" ||
			strings.TrimSpace(source.Publisher) == "" || !validResearchLocator(source.Locator) {
			return fmt.Errorf("research packet %s source %s has invalid provenance", packet.ID, source.ID)
		}
		published, err := parseResearchDate(source.PublishedOrUpdatedAt)
		if err != nil {
			return fmt.Errorf("research packet %s source %s published_or_updated_at: %w", packet.ID, source.ID, err)
		}
		accessed, err := parseResearchDate(source.AccessedAt)
		if err != nil {
			return fmt.Errorf("research packet %s source %s accessed_at: %w", packet.ID, source.ID, err)
		}
		if accessed.Before(published) {
			return fmt.Errorf("research packet %s source %s accessed before published/updated date", packet.ID, source.ID)
		}
		sources[source.ID] = source
	}

	limitations := map[string]ResearchLimitation{}
	for _, limitation := range packet.Limitations {
		if !researchLimitationIDPattern.MatchString(limitation.ID) || limitations[limitation.ID].ID != "" {
			return fmt.Errorf("research packet %s has invalid or duplicate limitation id %q", packet.ID, limitation.ID)
		}
		if strings.TrimSpace(limitation.Statement) == "" || len(limitation.SourceRefs) == 0 {
			return fmt.Errorf("research packet %s limitation %s requires statement and source_refs", packet.ID, limitation.ID)
		}
		if err := validateResearchRefs(packet.ID+" limitation "+limitation.ID, limitation.SourceRefs, sources); err != nil {
			return err
		}
		limitations[limitation.ID] = limitation
	}

	claims := map[string]ResearchClaim{}
	familyClaimed := map[string]bool{}
	for _, claim := range packet.Claims {
		if !researchClaimIDPattern.MatchString(claim.ID) || claims[claim.ID].ID != "" {
			return fmt.Errorf("research packet %s has invalid or duplicate claim id %q", packet.ID, claim.ID)
		}
		if strings.TrimSpace(claim.Statement) == "" || !knownResearchClaimCriticalities[claim.Criticality] ||
			!knownResearchClaimStatuses[claim.Status] {
			return fmt.Errorf("research packet %s claim %s has invalid metadata", packet.ID, claim.ID)
		}
		if claim.Criticality != "supporting" && len(claim.SourceRefs) == 0 {
			return fmt.Errorf("research packet %s material/critical claim %s requires source_refs", packet.ID, claim.ID)
		}
		if err := validateResearchRefs(packet.ID+" claim "+claim.ID, claim.SourceRefs, sources); err != nil {
			return err
		}
		for _, familyID := range claim.FamilyRefs {
			if _, ok := families[familyID]; !ok {
				return fmt.Errorf("research packet %s claim %s references unknown family %s", packet.ID, claim.ID, familyID)
			}
			familyClaimed[familyID] = true
		}
		for _, limitationID := range claim.LimitationRefs {
			if _, ok := limitations[limitationID]; !ok {
				return fmt.Errorf("research packet %s claim %s references unknown limitation %s", packet.ID, claim.ID, limitationID)
			}
		}
		if claim.Criticality == "critical" {
			if len(claim.SourceRefs) == 0 || len(claim.LimitationRefs) == 0 {
				return fmt.Errorf("research packet %s critical claim %s requires sources and limitations", packet.ID, claim.ID)
			}
			if !hasAuthoritativeResearchSource(claim.SourceRefs, sources) {
				return fmt.Errorf("research packet %s critical claim %s lacks authoritative/primary evidence", packet.ID, claim.ID)
			}
		}
		claims[claim.ID] = claim
	}

	if packet.Status == "saturated" {
		if len(families) == 0 || len(sources) == 0 || len(claims) == 0 {
			return fmt.Errorf("saturated research packet %s requires families sources and claims", packet.ID)
		}
		if !packet.Saturation.MajorSolutionFamiliesCovered || packet.Saturation.NewSourcesMateriallyChangeCandidates ||
			!packet.Saturation.CriticalLimitationsKnown || len(packet.Saturation.UnresolvedCriticalClaimIDs) != 0 {
			return fmt.Errorf("research packet %s does not satisfy saturation stop condition", packet.ID)
		}
		if err := uniqueNonEmptyStrings(packet.ID+" major_family_ids", packet.Saturation.MajorFamilyIDs); err != nil {
			return err
		}
		if len(packet.Saturation.MajorFamilyIDs) != len(families) {
			return fmt.Errorf("research packet %s saturation must name every declared solution family", packet.ID)
		}
		for _, familyID := range packet.Saturation.MajorFamilyIDs {
			if _, ok := families[familyID]; !ok {
				return fmt.Errorf("research packet %s saturation references unknown family %s", packet.ID, familyID)
			}
			if !familyClaimed[familyID] {
				return fmt.Errorf("research packet %s family %s has no sourced claim coverage", packet.ID, familyID)
			}
		}
		for _, claim := range packet.Claims {
			if claim.Criticality == "critical" && claim.Status == "unresolved" {
				return fmt.Errorf("saturated research packet %s still has unresolved critical claim %s", packet.ID, claim.ID)
			}
		}
	}
	return nil
}

func validateResearchRefs[T any](label string, refs []string, index map[string]T) error {
	if err := uniqueNonEmptyStrings(label+" refs", refs); err != nil {
		return err
	}
	for _, ref := range refs {
		if _, ok := index[ref]; !ok {
			return fmt.Errorf("%s references unknown id %s", label, ref)
		}
	}
	return nil
}

func hasAuthoritativeResearchSource(refs []string, sources map[string]ResearchSource) bool {
	for _, ref := range refs {
		switch sources[ref].Kind {
		case "official-doc", "scientific-paper", "standard", "security-advisory", "source-code", "benchmark-report":
			return true
		}
	}
	return false
}

func validResearchLocator(value string) bool {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "doi:") || strings.HasPrefix(value, "git:") {
		return len(value) > 4
	}
	u, err := url.Parse(value)
	return err == nil && u.Scheme == "https" && u.Host != ""
}

func parseResearchDate(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Time{}, fmt.Errorf("date is empty")
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("expected YYYY-MM-DD")
	}
	return t, nil
}
