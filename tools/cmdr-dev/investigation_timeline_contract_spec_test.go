package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type timelineContractFixtureFile struct {
	SchemaVersion int    `json:"schema_version"`
	ContractID   string `json:"contract_id"`
	Fixtures     []struct {
		ID string `json:"id"`
		Input struct {
			Scope struct {
				Tenant string `json:"tenant"`
				Environment string `json:"environment"`
				CasePermission string `json:"case_permission"`
				TimelinePermission string `json:"timeline_permission"`
			} `json:"scope"`
			Entries []struct {
				ID string `json:"id"`
				Kind string `json:"kind"`
				Tenant string `json:"tenant"`
				SourcePermission string `json:"source_permission"`
				Provenance string `json:"provenance"`
			} `json:"entries"`
		} `json:"input"`
		Expected struct {
			Status string `json:"status"`
			VisibleIDs []string `json:"visible_ids"`
			VisibleKinds []string `json:"visible_kinds"`
		} `json:"expected"`
	} `json:"fixtures"`
}

func verifyTimelineFixtureSafety(id string, tenant, environment, casePermission, timelinePermission string, entryIDs, entryKinds, entryTenants, entryDecisions, entryProvenance, visibleIDs, visibleKinds []string, outcome string) error {
	if len(visibleIDs) != len(visibleKinds) {
		return fmt.Errorf("%s: visible IDs/types mismatch", id)
	}
	if tenant == "" || tenant == "*" || environment == "" || casePermission != "allow" || timelinePermission != "allow" {
		if outcome != "deny" || len(visibleIDs) != 0 {
			return fmt.Errorf("%s: unsafe tenant/access expectation", id)
		}
		return nil
	}
	if len(entryIDs) != len(entryKinds) || len(entryIDs) != len(entryTenants) || len(entryIDs) != len(entryDecisions) || len(entryIDs) != len(entryProvenance) {
		return fmt.Errorf("%s: malformed fixture columns", id)
	}
	visible := make(map[string]string, len(visibleIDs))
	for i, key := range visibleIDs {
		if key == "" || visible[key] != "" {
			return fmt.Errorf("%s: missing or duplicate visible id", id)
		}
		visible[key] = visibleKinds[i]
	}
	for i, key := range entryIDs {
		if _, exposed := visible[key]; exposed {
			if entryTenants[i] != tenant || entryDecisions[i] != "allow" || entryProvenance[i] == "" || visible[key] != entryKinds[i] {
				return fmt.Errorf("%s: unsafe disclosure or type promotion of %s", id, key)
			}
		}
	}
	return nil
}

func TestTimelineContractFixturesAreIndependentSafeOracles(t *testing.T) {
	path := filepath.Join("..", "..", "engineering", "implementation", "investigation-timeline", "fixtures.json")
	blob, err := os.ReadFile(path)
	if err != nil { t.Fatal(err) }
	var data timelineContractFixtureFile
	if err := json.Unmarshal(blob, &data); err != nil { t.Fatal(err) }
	if data.SchemaVersion != 1 || data.ContractID != "TIMELINE-READONLY-CONTRACT-V1" || len(data.Fixtures) < 18 {
		t.Fatalf("missing bounded timeline contract matrix: %#v", data)
	}
	statuses := map[string]bool{}
	seen := map[string]bool{}
	for _, fixture := range data.Fixtures {
		if fixture.ID == "" || seen[fixture.ID] { t.Fatalf("duplicate fixture id: %q", fixture.ID) }
		seen[fixture.ID] = true
		var ids, kinds, tenants, decisions, provenance []string
		for _, e := range fixture.Input.Entries {
			ids = append(ids, e.ID)
			kinds = append(kinds, e.Kind)
			tenants = append(tenants, e.Tenant)
			decisions = append(decisions, e.SourcePermission)
			provenance = append(provenance, e.Provenance)
		}
		v := fixture.Input.Scope
		if err := verifyTimelineFixtureSafety(fixture.ID, v.Tenant, v.Environment, v.CasePermission, v.TimelinePermission, ids, kinds, tenants, decisions, provenance, fixture.Expected.VisibleIDs, fixture.Expected.VisibleKinds, fixture.Expected.Status); err != nil { t.Fatal(err) }
		statuses[fixture.Expected.Status] = true
	}
	for _, required := range []string{"allow", "deny", "filtered", "partial", "uncertain", "empty"} {
		if !statuses[required] { t.Errorf("missing %s test oracle", required) }
	}
}

func TestTimelineFixtureOracleRejectsMaskedObjectDisclosure(t *testing.T) {
	err := verifyTimelineFixtureSafety("adversarial", "tenant-a", "prod", "allow", "allow", []string{"secret"}, []string{"observed"}, []string{"tenant-b"}, []string{"deny"}, []string{"sensitive"}, []string{"secret"}, []string{"observed"}, "filtered")
	if err == nil || !strings.Contains(err.Error(), "unsafe disclosure") { t.Fatalf("expected rejected exfiltration oracle, got %v", err) }
}
