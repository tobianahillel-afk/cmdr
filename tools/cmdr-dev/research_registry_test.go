package main

import "testing"

func canonicalResearchPolicy() ResearchPolicy {
	return ResearchPolicy{
		SchemaVersion:  1,
		DefaultPolicy:  "deny-unsourced-material-claim",
		PacketStatuses: []string{"collecting", "saturated", "superseded"},
		SourceKinds: []string{
			"official-doc", "scientific-paper", "standard", "security-advisory",
			"source-code", "benchmark-report", "independent-analysis", "community-discussion",
		},
		ClaimCriticalities: []string{"supporting", "material", "critical"},
		ClaimStatuses:      []string{"supported", "contested", "unresolved"},
		FamilyDispositions: []string{"candidate", "excluded"},
		Limits:             ResearchLimits{MaxPackets: 256, MaxFamiliesPerPacket: 16, MaxSourcesPerPacket: 64, MaxClaimsPerPacket: 96, MaxLimitationsPerPacket: 32},
	}
}

func researchDecision(class, status string) EngineeringDecision {
	d := baseDecision()
	d.Class = class
	d.Status = status
	d.RiskTags = []string{"architecture"}
	if class != "A" {
		d.EvidenceRefs = nil
	}
	return d
}

func validSaturatedPacket() ResearchPacket {
	return ResearchPacket{
		ID: "RES-PKT-0001", DecisionID: "ENG-DEC-0001", Status: "saturated",
		Topic: "Index structure", Question: "Which indexing family fits the workload?", AsOf: "2026-09-23",
		Constraints: []string{"large forensic datasets", "local-first execution"},
		SolutionFamilies: []ResearchSolutionFamily{
			{ID: "FAMILY-LSM", Name: "LSM family", Description: "Log structured merge approaches", Disposition: "candidate", Rationale: "write-heavy candidate"},
			{ID: "FAMILY-BTREE", Name: "B-tree family", Description: "B-tree derived approaches", Disposition: "candidate", Rationale: "read-oriented candidate"},
		},
		Sources: []ResearchSource{
			{ID: "SRC-001", Kind: "scientific-paper", Title: "Paper A", Publisher: "Conference", Locator: "https://example.org/paper-a", PublishedOrUpdatedAt: "2026-09-01", AccessedAt: "2026-09-23"},
			{ID: "SRC-002", Kind: "official-doc", Title: "Reference B", Publisher: "Project B", Locator: "https://example.org/ref-b", PublishedOrUpdatedAt: "2026-09-10", AccessedAt: "2026-09-23"},
		},
		Limitations: []ResearchLimitation{
			{ID: "LIMIT-001", Statement: "Benchmarks are not CMDR-specific", SourceRefs: []string{"SRC-001"}},
		},
		Claims: []ResearchClaim{
			{ID: "CLAIM-001", Statement: "LSM is a credible family", Criticality: "critical", Status: "supported", SourceRefs: []string{"SRC-001"}, FamilyRefs: []string{"FAMILY-LSM"}, LimitationRefs: []string{"LIMIT-001"}},
			{ID: "CLAIM-002", Statement: "B-tree is a credible family", Criticality: "material", Status: "supported", SourceRefs: []string{"SRC-002"}, FamilyRefs: []string{"FAMILY-BTREE"}},
		},
		Saturation: ResearchSaturation{
			MajorFamilyIDs:                       []string{"FAMILY-LSM", "FAMILY-BTREE"},
			MajorSolutionFamiliesCovered:         true,
			NewSourcesMateriallyChangeCandidates: false,
			CriticalLimitationsKnown:             true,
		},
	}
}

func TestResearchRegistryAllowsEmptyRegistry(t *testing.T) {
	registry := ResearchRegistry{SchemaVersion: 1, RegistryKind: "engineering-research-evidence"}
	if err := validateResearchRegistry(canonicalResearchPolicy(), registry, DecisionRegistry{}); err != nil {
		t.Fatal(err)
	}
}

func TestSaturatedResearchPacketPasses(t *testing.T) {
	d := researchDecision("B", "researching")
	decisions := DecisionRegistry{SchemaVersion: 1, RegistryKind: "engineering-technical-decisions", Decisions: []EngineeringDecision{d}}
	registry := ResearchRegistry{SchemaVersion: 1, RegistryKind: "engineering-research-evidence", Packets: []ResearchPacket{validSaturatedPacket()}}
	if err := validateResearchRegistry(canonicalResearchPolicy(), registry, decisions); err != nil {
		t.Fatal(err)
	}
}

func TestResearchPacketRejectsClassAAndUnknownDecision(t *testing.T) {
	packet := validSaturatedPacket()
	for _, decisions := range []DecisionRegistry{
		{SchemaVersion: 1, RegistryKind: "engineering-technical-decisions"},
		{SchemaVersion: 1, RegistryKind: "engineering-technical-decisions", Decisions: []EngineeringDecision{researchDecision("A", "accepted")}},
	} {
		if err := validateResearchRegistry(canonicalResearchPolicy(), ResearchRegistry{SchemaVersion: 1, RegistryKind: "engineering-research-evidence", Packets: []ResearchPacket{packet}}, decisions); err == nil {
			t.Fatal("expected packet decision binding rejection")
		}
	}
}

func TestSaturationFailsWhenFamilyCoverageOrCriticalEvidenceIsMissing(t *testing.T) {
	d := researchDecision("B", "researching")
	decisions := DecisionRegistry{SchemaVersion: 1, RegistryKind: "engineering-technical-decisions", Decisions: []EngineeringDecision{d}}
	packet := validSaturatedPacket()
	packet.Claims = packet.Claims[:1]
	if err := validateResearchRegistry(canonicalResearchPolicy(), ResearchRegistry{SchemaVersion: 1, RegistryKind: "engineering-research-evidence", Packets: []ResearchPacket{packet}}, decisions); err == nil {
		t.Fatal("expected missing family coverage rejection")
	}
	packet = validSaturatedPacket()
	packet.Claims[0].LimitationRefs = nil
	if err := validateResearchRegistry(canonicalResearchPolicy(), ResearchRegistry{SchemaVersion: 1, RegistryKind: "engineering-research-evidence", Packets: []ResearchPacket{packet}}, decisions); err == nil {
		t.Fatal("expected critical limitation rejection")
	}
}

func TestAcceptedDecisionRequiresOwnSaturatedPacket(t *testing.T) {
	d := researchDecision("B", "accepted")
	d.EvidenceRefs = []string{"RES-PKT-0001"}
	decisions := DecisionRegistry{SchemaVersion: 1, RegistryKind: "engineering-technical-decisions", Decisions: []EngineeringDecision{d}}
	packet := validSaturatedPacket()
	registry := ResearchRegistry{SchemaVersion: 1, RegistryKind: "engineering-research-evidence", Packets: []ResearchPacket{packet}}
	if err := validateResearchRegistry(canonicalResearchPolicy(), registry, decisions); err != nil {
		t.Fatal(err)
	}
	packet.Status = "collecting"
	if err := validateResearchRegistry(canonicalResearchPolicy(), ResearchRegistry{SchemaVersion: 1, RegistryKind: "engineering-research-evidence", Packets: []ResearchPacket{packet}}, decisions); err == nil {
		t.Fatal("expected non-saturated accepted evidence rejection")
	}
}

func TestResearchContextIsBoundedByValidatedPolicy(t *testing.T) {
	packet := validSaturatedPacket()
	registry := ResearchRegistry{SchemaVersion: 1, RegistryKind: "engineering-research-evidence", Packets: []ResearchPacket{packet}}
	d := researchDecision("B", "researching")
	if err := validateResearchRegistry(canonicalResearchPolicy(), registry, DecisionRegistry{Decisions: []EngineeringDecision{d}}); err != nil {
		t.Fatal(err)
	}
	if len(packet.Sources) > canonicalResearchPolicy().Limits.MaxSourcesPerPacket || len(packet.Claims) > canonicalResearchPolicy().Limits.MaxClaimsPerPacket {
		t.Fatal("fixture exceeds bounded research context")
	}
}
