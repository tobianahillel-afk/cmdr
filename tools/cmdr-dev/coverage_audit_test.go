package main

import "testing"

func TestBuildCoverageAuditDetectsCoreGapFamilies(t *testing.T) {
	graph := ProductGraph{
		SpecTreeDigest: "digest",
		Entities: []ProductEntity{
			{ID: "CAP-CMD-001", Kind: "capability", SourcePath: "cap1.md", SourceSHA256: "a"},
			{ID: "CAP-CMD-002", Kind: "capability", SourcePath: "cap2.md", SourceSHA256: "b"},
			{ID: "CAP-CMD-999", Kind: "capability", ReferenceOnly: true},
			{ID: "CMD-MC-001", Kind: "screen", SourcePath: "screen.md", SourceSHA256: "c", RegistryStatus: "registered-active"},
			{ID: "REQ-SEC-999", Kind: "requirement", ReferenceOnly: true},
			{ID: "perm.missing.read", Kind: "permission", ReferenceOnly: true},
			{ID: "OPEN-999", Kind: "open-decision", ReferenceOnly: true},
		},
		Edges: []ProductEdge{
			{From: "CAP-CMD-001", To: "CAP-CMD-999", Kind: "reference", SourcePath: "cap1.md"},
			{From: "CAP-CMD-001", To: "REQ-SEC-999", Kind: "reference", SourcePath: "cap1.md"},
			{From: "CAP-CMD-001", To: "perm.missing.read", Kind: "reference", SourcePath: "cap1.md"},
			{From: "CAP-CMD-001", To: "OPEN-999", Kind: "reference", SourcePath: "cap1.md"},
		},
	}
	report := buildCoverageAudit(
		graph,
		map[string][]string{
			"CAP-CMD-001": {"cap-reg.md"},
			"CAP-CMD-003": {"cap-reg.md"},
		},
		map[string][]string{"REQ-PROD-001": {"req-baseline.json"}},
		map[string]struct{}{"CMD-MC-001": {}},
		map[string][]string{"perm.registered.read": {"perm-reg.md"}},
		map[string][]string{"OPEN-001": {"open.md"}},
		"cmdr-product-spec",
	)

	codes := map[string]bool{}
	for _, gap := range report.Gaps {
		codes[gap.Code+":"+gap.SubjectID] = true
		if len(gap.SourcePaths) == 0 {
			t.Fatalf("gap lacks source paths: %#v", gap)
		}
	}
	for _, want := range []string{
		"registered-capability-missing-owner:CAP-CMD-003",
		"owned-capability-unregistered:CAP-CMD-002",
		"capability-reference-outside-register:CAP-CMD-999",
		"active-requirement-unreferenced:REQ-PROD-001",
		"requirement-reference-outside-catalog:REQ-SEC-999",
		"registered-permission-unreferenced:perm.registered.read",
		"permission-reference-outside-register:perm.missing.read",
		"active-open-decision-unreferenced:OPEN-001",
		"open-decision-reference-outside-active-register:OPEN-999",
	} {
		if !codes[want] {
			t.Fatalf("missing expected gap %s", want)
		}
	}
}
