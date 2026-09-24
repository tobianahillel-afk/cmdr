package main

import "testing"

func readinessRecord(id, delivery string, opens ...string) CapabilityRegistryRecord {
	return CapabilityRegistryRecord{
		ID: id, Name: id, Status: "draft", DeliveryStatus: delivery, DeliveryMode: "planned",
		CanonicalFile: "capability.md", OpenDecisions: opens, SourcePath: "register.md",
	}
}

func TestParseCapabilityRegistryContent(t *testing.T) {
	content := `# Register

| ID | Name | Owner | Status | Delivery status | Delivery mode | Canonical file | OPEN |
|---|---|---|---|---|---|---|---|
| CAP-SET-001 | Tenant | Owner | draft | defined | planned | a.md | OPEN-013 |
| CAP-SET-004 | Context | Owner | draft | defined | planned | b.md | none |
`
	records, err := parseCapabilityRegistryContent("register.md", content)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 || records[0].ID != "CAP-SET-001" || len(records[0].OpenDecisions) != 1 {
		t.Fatalf("unexpected capability records: %#v", records)
	}
	if records[1].CanonicalFile != "b.md" {
		t.Fatalf("unexpected canonical file: %#v", records[1])
	}
}

func TestParseCapabilityRegistryCombinedDeliveryVariant(t *testing.T) {
	content := `# Endpoint register

| ID | Capability | Primary role | Primary concepts | Key consumers | OPEN | Delivery |
|---|---|---|---|---|---|---|
| CAP-EPT-031 | Local Detection | Analyst | match | Investigate | OPEN-008, OPEN-017 | defined / planned |
`
	records, err := parseCapabilityRegistryContent("endpoint.md", content)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected one record, got %#v", records)
	}
	record := records[0]
	if record.ID != "CAP-EPT-031" || record.Name != "Local Detection" {
		t.Fatalf("unexpected identity: %#v", record)
	}
	if record.DeliveryStatus != "defined" || record.DeliveryMode != "planned" {
		t.Fatalf("unexpected delivery fields: %#v", record)
	}
	if record.Status != "" || record.CanonicalFile != "" {
		t.Fatalf("row-absent optional metadata must remain empty: %#v", record)
	}
	if len(record.OpenDecisions) != 2 || record.OpenDecisions[0] != "OPEN-008" || record.OpenDecisions[1] != "OPEN-017" {
		t.Fatalf("unexpected OPEN refs: %#v", record.OpenDecisions)
	}
}

func TestParseCapabilityRegistryStatusDeliveryVariant(t *testing.T) {
	content := `# Cloud register

| ID | Title | Primary role | Canonical path | Open decisions | Status | Delivery |
|---|---|---|---|---|---|---|
| CAP-INV-601 | Cloud Analysis Intake | Analyst | cloud-intake.md | OPEN-008/OPEN-012 | defined | planned |
| CAP-INV-699 | Future Cloud Capability | Analyst | future-cloud.md | OPEN-014 | proposed | planned |
`
	records, err := parseCapabilityRegistryContent("cloud.md", content)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatalf("expected two records, got %#v", records)
	}
	if records[0].Name != "Cloud Analysis Intake" || records[0].CanonicalFile != "cloud-intake.md" {
		t.Fatalf("unexpected cloud identity metadata: %#v", records[0])
	}
	if len(records[0].OpenDecisions) != 2 || records[0].OpenDecisions[0] != "OPEN-008" || records[0].OpenDecisions[1] != "OPEN-012" {
		t.Fatalf("unexpected cloud OPEN refs: %#v", records[0].OpenDecisions)
	}
	if records[0].DeliveryStatus != "defined" || records[0].DeliveryMode != "planned" {
		t.Fatalf("unexpected defined delivery mapping: %#v", records[0])
	}
	if records[1].DeliveryStatus != "proposed" || records[1].DeliveryMode != "planned" {
		t.Fatalf("unexpected proposed delivery mapping: %#v", records[1])
	}
}

func TestParseCapabilityRegistryRejectsDocumentStatusAsDeliveryStatus(t *testing.T) {
	content := `# Ambiguous register

| ID | Name | Status | Delivery |
|---|---|---|---|
| CAP-INV-601 | Cloud Analysis Intake | draft | planned |
`
	if _, err := parseCapabilityRegistryContent("ambiguous.md", content); err == nil {
		t.Fatal("expected draft/planned schema ambiguity rejection")
	}
}

func TestParseCapabilityRegistryUnderscoreDeliveryColumns(t *testing.T) {
	content := `# Settings register

| Capability ID | Name | Status | delivery_status | delivery_mode | Canonical file | OPEN |
|---|---|---|---|---|---|---|
| CAP-SET-005 | Principal Administrative Lifecycle | draft | defined | planned | settings.md | OPEN-013 |
`
	records, err := parseCapabilityRegistryContent("settings.md", content)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected one record, got %#v", records)
	}
	record := records[0]
	if record.DeliveryStatus != "defined" || record.DeliveryMode != "planned" || record.Status != "draft" {
		t.Fatalf("unexpected settings delivery mapping: %#v", record)
	}
	if record.ID != "CAP-SET-005" || record.CanonicalFile != "settings.md" {
		t.Fatalf("unexpected settings identity metadata: %#v", record)
	}
}

func TestParseCapabilityRegistryFileLevelDeliveryDefaults(t *testing.T) {
	content := "# Register\n\nAll capabilities are documentary draft, delivery status `defined`, delivery mode `planned`.\n\n" +
		"| ID | Capability | Primary role | Canonical file | OPEN |\n" +
		"|---|---|---|---|---|\n" +
		"| CAP-INV-519 | Intelligence Analysis Intake | Analyst | ti.md | OPEN-013/014 |\n"
	records, err := parseCapabilityRegistryContent("threat-intel.md", content)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("expected one record, got %#v", records)
	}
	record := records[0]
	if record.DeliveryStatus != "defined" || record.DeliveryMode != "planned" {
		t.Fatalf("unexpected file-level delivery defaults: %#v", record)
	}
	if record.Name != "Intelligence Analysis Intake" || record.CanonicalFile != "ti.md" {
		t.Fatalf("unexpected file-level record metadata: %#v", record)
	}
}

func TestParseCapabilityRegistryWithoutDeliveryEvidenceFails(t *testing.T) {
	content := "# Register\n\n| ID | Capability | OPEN |\n|---|---|---|\n| CAP-INV-519 | Intake | none |\n"
	if _, err := parseCapabilityRegistryContent("missing-delivery.md", content); err == nil {
		t.Fatal("expected missing delivery evidence rejection")
	}
}

func TestExpandCapabilitySelectorsSupportsExactRangeAndSlash(t *testing.T) {
	known := map[string]CapabilityRegistryRecord{
		"CAP-INV-001": readinessRecord("CAP-INV-001", "defined"),
		"CAP-INV-002": readinessRecord("CAP-INV-002", "defined"),
		"CAP-INV-003": readinessRecord("CAP-INV-003", "defined"),
		"CAP-INV-007": readinessRecord("CAP-INV-007", "defined"),
	}
	got := expandCapabilitySelectors("CAP-INV-001..003 plus CAP-INV-007/002", known)
	want := []string{"CAP-INV-001", "CAP-INV-002", "CAP-INV-003", "CAP-INV-007"}
	if len(got) != len(want) {
		t.Fatalf("unexpected selector expansion: %#v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("unexpected selector expansion: %#v", got)
		}
	}
}

func TestCompileImplementationReadinessStates(t *testing.T) {
	records := map[string]CapabilityRegistryRecord{
		"CAP-SET-001": readinessRecord("CAP-SET-001", "defined", "OPEN-013"),
		"CAP-SET-002": readinessRecord("CAP-SET-002", "defined"),
		"CAP-SET-004": readinessRecord("CAP-SET-004", "defined"),
		"CAP-INV-002": readinessRecord("CAP-INV-002", "defined"),
		"CAP-INV-106": readinessRecord("CAP-INV-106", "proposed", "OPEN-014"),
	}
	activeOpen := map[string][]string{
		"OPEN-013": {"unresolved-decisions.md"},
		"OPEN-014": {"unresolved-decisions.md"},
	}
	dependencies := []DependencyReadinessEvidence{{
		ID: "DEP-X-001", Dependent: "CAP-SET-002", Status: "partial", Blocking: "before implementation",
		AffectedCapabilities: []string{"CAP-SET-002"}, SourcePath: "dependency-register.md",
	}}
	implemented := map[string]string{"CAP-SET-004": "E9-PILOT-001C-RUNTIME"}
	program, err := compileImplementationReadiness(records, activeOpen, dependencies, implemented, "baseline", "digest")
	if err != nil {
		t.Fatal(err)
	}
	states := map[string]string{}
	for _, record := range program.Records {
		states[record.Capability] = record.State
	}
	if states["CAP-SET-001"] != "BLOCKED" || states["CAP-SET-002"] != "BLOCKED" ||
		states["CAP-SET-004"] != "IMPLEMENTED" || states["CAP-INV-002"] != "READY" ||
		states["CAP-INV-106"] != "PROPOSED" {
		t.Fatalf("unexpected readiness states: %#v", states)
	}
	if program.Summary.Total != 5 || program.Summary.Implemented != 1 || program.Summary.Ready != 1 ||
		program.Summary.Blocked != 2 || program.Summary.Proposed != 1 {
		t.Fatalf("unexpected readiness summary: %#v", program.Summary)
	}
}

func TestCompileReadinessRejectsImplementedProposedCapability(t *testing.T) {
	records := map[string]CapabilityRegistryRecord{
		"CAP-INV-106": readinessRecord("CAP-INV-106", "proposed"),
	}
	_, err := compileImplementationReadiness(records, nil, nil, map[string]string{"CAP-INV-106": "WORK"}, "baseline", "digest")
	if err == nil {
		t.Fatal("expected proposed implementation claim rejection")
	}
}

func TestDependencyBlockingPolicy(t *testing.T) {
	cases := []struct {
		status, blocking string
		want             bool
	}{
		{"active", "yes", false},
		{"partial", "yes where required", true},
		{"open", "before implementation", true},
		{"planned", "model blocking", true},
		{"partial", "no essential AI dependency", false},
	}
	for _, tc := range cases {
		got := dependencyBlocksImplementation(DependencyReadinessEvidence{Status: tc.status, Blocking: tc.blocking})
		if got != tc.want {
			t.Fatalf("status=%s blocking=%s: got %t want %t", tc.status, tc.blocking, got, tc.want)
		}
	}
}
