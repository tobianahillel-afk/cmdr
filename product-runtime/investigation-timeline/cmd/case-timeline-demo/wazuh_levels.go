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
	"time"

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

// A bounded, descriptive time distribution of visible imported alerts.
// These are alert-volume windows, not inferred incident boundaries.
type wazuhAlertWindow struct {
	StartUTC string `json:"start_utc"`
	Count int `json:"count"`
	HighestKnownRuleLevel *int `json:"highest_known_rule_level,omitempty"`
	UnknownRuleLevels int `json:"unknown_rule_levels"`
}

type wazuhRuleTriageReceipt struct {
	Format string `json:"format"`
	Synthetic bool `json:"synthetic"`
	ImportedEntries int `json:"imported_entries"`
	Case caseTriageReceipt `json:"case_triage"`
	RuleLevels []wazuhLevelCount `json:"rule_levels"`
	UnknownRuleLevels int `json:"unknown_rule_levels"`
	DistinctTimeWindows int `json:"distinct_time_windows"`
	PeakWindows []wazuhAlertWindow `json:"peak_windows"`
	UnplacedAlerts int `json:"unplaced_alerts"`
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
	windows := make(map[string]*wazuhAlertWindow)
	visibleImported, unknown, unplaced := 0, 0, 0
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
		if entry.TimeUncertain {
			unplaced++
			return nil
		}
		when, err := time.Parse(time.RFC3339Nano, entry.SourceTime)
		if err != nil {
			return errors.New("invalid authorized Wazuh Timeline timestamp")
		}
		key := when.UTC().Truncate(5 * time.Minute).Format(time.RFC3339)
		window := windows[key]
		if window == nil {
			window = &wazuhAlertWindow{StartUTC: key}
			windows[key] = window
		}
		window.Count++
		if level.Known {
			if window.HighestKnownRuleLevel == nil || level.Level > *window.HighestKnownRuleLevel {
				highest := level.Level
				window.HighestKnownRuleLevel = &highest
			}
		} else {
			window.UnknownRuleLevels++
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
	peaks := make([]wazuhAlertWindow, 0, len(windows))
	for _, window := range windows {
		peaks = append(peaks, *window)
	}
	sort.Slice(peaks, func(i, j int) bool {
		if peaks[i].Count != peaks[j].Count {
			return peaks[i].Count > peaks[j].Count
		}
		return peaks[i].StartUTC < peaks[j].StartUTC
	})
	// Export at most ten active windows; all windows are still counted.
	if len(peaks) > 10 {
		peaks = peaks[:10]
	}
	histogram := make([]wazuhLevelCount, 0, len(keys))
	for _, level := range keys {
		histogram = append(histogram, wazuhLevelCount{Level: level, Count: counts[level]})
	}
	return json.NewEncoder(output).Encode(wazuhRuleTriageReceipt{
		Format: "cmdr-local-wazuh4-case-triage-v1",
		Synthetic: true, ImportedEntries: imported, Case: receipt,
		RuleLevels: histogram, UnknownRuleLevels: unknown,
		DistinctTimeWindows: len(windows), PeakWindows: peaks, UnplacedAlerts: unplaced,
	})
}
