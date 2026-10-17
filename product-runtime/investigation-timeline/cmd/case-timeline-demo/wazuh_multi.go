package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Wazuh 4.x is parsed only by the existing validated normalizer. This small
// adapter adds a trusted operator-side feed label to the normalized event ID,
// before passing it to the unchanged owner-authorized JSONL importer.
func namespaceWazuhEvents(ctx context.Context, normalized []byte, feed int) ([]byte, error) {
	if feed < 1 || feed > maxWazuhAlertFiles {
		return nil, errors.New("invalid Wazuh feed number")
	}
	var result bytes.Buffer
	lines := bytes.Split(normalized, []byte("\n"))
	for i, line := range lines {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if i == len(lines)-1 && len(line) == 0 {
			continue
		}
		var event normalizedLocalEvent
		if err := json.Unmarshal(line, &event); err != nil {
			return nil, errors.New("invalid normalized Wazuh event")
		}
		event.ID = fmt.Sprintf("feed-%d-%s", feed, event.ID)
		if !validLocalEventID(event.ID) {
			return nil, errors.New("Wazuh feed identity exceeds supported length")
		}
		if err := json.NewEncoder(&result).Encode(event); err != nil {
			return nil, err
		}
		if result.Len() > maxDemoSnapshotBytes {
			return nil, errors.New("multi-feed normalized Wazuh data too large")
		}
	}
	return result.Bytes(), nil
}

// All input files and source decisions are validated before a bearer is
// issued, an HTTP listener starts, or a cross-Case report is emitted.
func loadMultipleWazuhFeeds(ctx context.Context, principal, snapshotPath, caseID string, paths []string) (timeline.MemoryTimelineData, int, error) {
	if err := ctx.Err(); err != nil {
		return timeline.MemoryTimelineData{}, 0, err
	}
	if snapshotPath == "" || principal == "" || caseID == "" || len(caseID) > 1024 ||
		len(paths) < 2 || len(paths) > maxWazuhAlertFiles {
		return timeline.MemoryTimelineData{}, 0, errors.New("select one owner snapshot, an authorized Case and 2..4 Wazuh alert files")
	}
	owner, err := loadDemoSnapshot(snapshotPath)
	if err != nil {
		return timeline.MemoryTimelineData{}, 0, err
	}
	info, err := os.Lstat(snapshotPath)
	if err != nil {
		return timeline.MemoryTimelineData{}, 0, errors.New("cannot identify owner snapshot")
	}
	before := len(owner.Entries)
	seen := []os.FileInfo{info}
	for i, path := range paths {
		if err := ctx.Err(); err != nil {
			return timeline.MemoryTimelineData{}, 0, err
		}
		data, err := readLocalEventFile(path)
		if err != nil {
			return timeline.MemoryTimelineData{}, 0, err
		}
		current, err := os.Lstat(path)
		if err != nil {
			return timeline.MemoryTimelineData{}, 0, errors.New("cannot identify Wazuh alert file")
		}
		for _, previous := range seen {
			if os.SameFile(current, previous) {
				return timeline.MemoryTimelineData{}, 0, errors.New("Wazuh alert inputs must be distinct files")
			}
		}
		seen = append(seen, current)
		events, err := normalizeWazuh4Alerts(ctx, data)
		if err != nil {
			return timeline.MemoryTimelineData{}, 0, err
		}
		events, err = namespaceWazuhEvents(ctx, events, i+1)
		if err != nil {
			return timeline.MemoryTimelineData{}, 0, err
		}
		owner, err = importLocalEvents(ctx, owner, principal, caseID, events)
		if err != nil {
			return timeline.MemoryTimelineData{}, 0, err
		}
	}
	count := len(owner.Entries) - before
	if count < len(paths) || count > timeline.MaxEntries {
		return timeline.MemoryTimelineData{}, 0, errors.New("inconsistent Wazuh feed event count")
	}
	return owner, count, nil
}

// Multi-feed viewer and two-Case time leads share precisely the same
// bounded owner import; no second parser, permission resolver or storage.
func serveDemoWithWazuhAlertsFiles(ctx context.Context, principal, snapshotPath, caseID string, paths []string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(paths) == 1 {
		return serveDemoWithWazuhAlerts(ctx, principal, snapshotPath, caseID, paths[0], output)
	}
	owner, _, err := loadMultipleWazuhFeeds(ctx, principal, snapshotPath, caseID, paths)
	if err != nil {
		return err
	}
	return serveDemoWithData(ctx, principal, owner, output)
}

// Multi-feed proximity is a bounded unsigned temporal lead, not a causal
// correlation, production finding, provider-integrated Case or custody record.
func writeDemoMultiWazuhProximity(ctx context.Context, principal, snapshotPath, intoCase, withCase string, paths []string, seconds int, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if intoCase == withCase || withCase == "" || len(withCase) > 1024 ||
		seconds < 1 || seconds > 300 {
		return errors.New("choose two distinct Cases and a 1..300-second interval")
	}
	owner, imported, err := loadMultipleWazuhFeeds(ctx, principal, snapshotPath, intoCase, paths)
	if err != nil {
		return err
	}
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, owner)
	if err != nil {
		return err
	}
	both, err := collectAuthorizedCrossCaseSequence(ctx, handler, token, intoCase, withCase)
	if err != nil {
		return err
	}
	leads, err := selectCaseProximity(ctx, both, seconds)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(localImportedProximity{
		Format: "cmdr-local-multi-wazuh4-case-proximity-v1",
		Synthetic: true, ImportedEntries: imported, Proximity: leads,
	})
}

// writeDemoMultiWazuhTriage makes an actionable per-Case receipt from multiple
// Wazuh 4.x alert feeds. Every aggregate is derived from the same authorized
// Case -> Timeline -> typed source pagination as the analyst UI, not raw input.
func writeDemoMultiWazuhTriage(ctx context.Context, principal, snapshotPath, caseID string, paths []string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	owner, imported, err := loadMultipleWazuhFeeds(ctx, principal, snapshotPath, caseID, paths)
	if err != nil {
		return err
	}
	receipt, err := collectDemoAuthorizedCaseTriage(ctx, principal, owner, caseID)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(localImportedTriage{
		Format: "cmdr-local-multi-wazuh4-case-triage-v1",
		Synthetic: true, ImportedEntries: imported, Case: receipt,
	})
}

// writeDemoMultiWazuhComparison reuses the already-authorized two-Case
// aggregate evaluator after atomic import of all selected local feeds.
// Both Cases must pass their owner checks before any output is emitted.
func writeDemoMultiWazuhComparison(ctx context.Context, principal, snapshotPath, intoCase, withCase string, paths []string, output io.Writer) error {
	if intoCase == "" || withCase == "" || intoCase == withCase || len(withCase) > 1024 {
		return errors.New("select two distinct authorized Cases for multi-feed comparison")
	}
	owner, imported, err := loadMultipleWazuhFeeds(ctx, principal, snapshotPath, intoCase, paths)
	if err != nil {
		return err
	}
	comparison, err := collectAuthorizedCaseComparison(ctx, principal, owner, intoCase, withCase)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(output).Encode(localImportedComparison{
		Format: "cmdr-local-multi-wazuh4-case-comparison-v1",
		Synthetic: true, ImportedEntries: imported, Comparison: comparison,
	})
}
