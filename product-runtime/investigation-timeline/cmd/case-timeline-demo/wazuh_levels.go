package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Source metadata, not an authorization claim or a derived severity score.
// Wazuh 4.x rule level is an optional numeric field (0..16).
type wazuhRuleLevel struct {
	Level int
	Known bool
}

type wazuhLevelCount struct {
	Level int `json:"level"`
	Count int `json:"count"`
}

type wazuhRuleTriageReceipt struct {
	Format string `json:"format"`
	Synthetic bool `json:"synthetic"`
	ImportedEntries int `json:"imported_entries"`
	Case caseTriageReceipt `json:"case_triage"`
	RuleLevels []wazuhLevelCount `json:"rule_levels"`
	UnknownRuleLevels int `json:"unknown_rule_levels"`
}

// This mapping is computed from the SAME bytes that were normalized and
// imported into the Case; no second file open or independent ACL is introduced.
func wazuhRuleLevels(raw []byte, caseID string) (map[string]wazuhRuleLevel, error) {
	result := make(map[string]wazuhRuleLevel)
	for _, line := range bytes.Split(raw, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		var alert struct {
			ID string `json:"id"`
			Rule struct {
				Level *int `json:"level"`
			} `json:"rule"`
		}
		if err := json.Unmarshal(line, &alert); err != nil {
			return nil, errors.New("invalid Wazuh rule level metadata")
		}
		digest := sha256.Sum256([]byte(caseID + "\x00wazuh-" + alert.ID))
		entryID := fmt.Sprintf("local-jsonl-entry-%x", digest)
		if _, exists := result[entryID]; exists {
			return nil, errors.New("duplicate Wazuh rule identity")
		}
		level := wazuhRuleLevel{}
		if alert.Rule.Level != nil {
			if *alert.Rule.Level < 0 || *alert.Rule.Level > 16 {
				return nil, errors.New("Wazuh rule level outside 0..16")
			}
			level = wazuhRuleLevel{Level: *alert.Rule.Level, Known: true}
		}
		result[entryID] = level
	}
	return result, nil
}

// writeDemoWazuh4RuleTriage adds genuine source rule-level analysis to the
// existing full authenticated Case/Timeline/typed-source pagination path.
// Metadata is counted ONLY when its imported Timeline Entry is authorized.
func writeDemoWazuh4RuleTriage(ctx context.Context, principal, snapshotPath, intoCase, alertsPath string, output io.Writer) error {
	merged, imported, raw, err := prepareWazuh4CaseWithRaw(ctx, principal, snapshotPath, intoCase, alertsPath)
	if err != nil {
		return err
	}
	levels, err := wazuhRuleLevels(raw, intoCase)
	if err != nil {
		return err
	}
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, merged)
	if err != nil {
		return err
	}
	counts := make(map[int]int)
	visibleImported, unknown := 0, 0
	receipt, err := collectDemoCaseTriageWithObserver(ctx, handler, token, intoCase, func(entry timeline.ProjectedEntry) error {
		level, exists := levels[entry.ID]
		if !exists {
			return nil
		}
		visibleImported++
		if level.Known {
			counts[level.Level]++
		} else {
			unknown++
		}
		return nil
	})
	if err != nil {
		return err
	}
	if visibleImported != imported || visibleImported != len(levels) {
		return errors.New("not every imported Wazuh alert is visible to the authorized Case reader")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	keys := make([]int, 0, len(counts))
	for level := range counts {
		keys = append(keys, level)
	}
	sort.Ints(keys)
	histogram := make([]wazuhLevelCount, 0, len(keys))
	for _, level := range keys {
		histogram = append(histogram, wazuhLevelCount{Level: level, Count: counts[level]})
	}
	return json.NewEncoder(output).Encode(wazuhRuleTriageReceipt{
		Format: "cmdr-local-wazuh4-case-triage-v1",
		Synthetic: true, ImportedEntries: imported, Case: receipt,
		RuleLevels: histogram, UnknownRuleLevels: unknown,
	})
}
