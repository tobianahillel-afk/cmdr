package main

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Save Wazuh into the existing private, scoped Case snapshot through the
// same preparation routine used by the one-step mixed-source import.
func writePersistedWazuhCase(ctx context.Context, principal, snapshotPath, caseID string,
	alertFiles []string, destination string, output io.Writer) error {
	if !filepath.IsAbs(destination) {
		return errors.New("choose a new absolute private destination")
	}
	if destination == snapshotPath {
		return errors.New("cannot overwrite the Case owner")
	}
	for _, path := range alertFiles {
		if path == destination {
			return errors.New("cannot overwrite a Wazuh evidence source")
		}
	}
	merged, err := preparePersistedWazuhCase(ctx, principal, snapshotPath, caseID, alertFiles)
	if err != nil { return err }
	return persistScopedCaseOwner(ctx, merged, principal, caseID, destination, output)
}

// Prepare the validated Wazuh Case in memory without creating a file.
// Wazuh agent/rule/severity metadata is linked only to the canonical typed
// sources created by the already existing first-party normalizer.
func preparePersistedWazuhCase(ctx context.Context, principal, snapshotPath, caseID string,
	alertFiles []string) (timeline.MemoryTimelineData, error) {
	if ctx == nil || ctx.Err() != nil || principal == "" || caseID == "" ||
		snapshotPath == "" || len(alertFiles) < 1 || len(alertFiles) > maxWazuhAlertFiles {
		return timeline.MemoryTimelineData{}, errors.New("select a Case and 1..4 Wazuh source files")
	}
	files := append([]string{snapshotPath}, alertFiles...)
	for _, path := range files {
		if !filepath.IsAbs(path) {
			return timeline.MemoryTimelineData{}, errors.New("Wazuh owner and evidence paths must be absolute")
		}
	}
	// All selected files must be distinct regular files, not symlinks or
	// hardlink aliases of the owner or of another Wazuh source.
	seen := make([]os.FileInfo, 0, len(files))
	for _, path := range files {
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() {
			return timeline.MemoryTimelineData{}, errors.New("Wazuh evidence and owner must be regular files")
		}
		for _, earlier := range seen {
			if os.SameFile(info, earlier) {
				return timeline.MemoryTimelineData{}, errors.New("duplicate Wazuh evidence or owner alias")
			}
		}
		seen = append(seen, info)
	}
	var merged timeline.MemoryTimelineData
	var levels map[string]wazuhRuleLevel
	if len(alertFiles) == 1 {
		var raw []byte
		var err error
		merged, _, raw, err = prepareWazuh4CaseWithRaw(ctx, principal, snapshotPath, caseID, alertFiles[0])
		if err != nil { return timeline.MemoryTimelineData{}, err }
		levels, err = wazuhRuleLevels(raw, caseID)
		if err != nil { return timeline.MemoryTimelineData{}, err }
	} else {
		var err error
		merged, _, levels, err = loadMultipleWazuhFeedsWithLevels(ctx, principal, snapshotPath, caseID, alertFiles)
		if err != nil { return timeline.MemoryTimelineData{}, err }
	}
	linked := 0
	for i := range merged.Sources {
		source := &merged.Sources[i]
		if !strings.HasPrefix(source.Identity.Reference, "local-jsonl-source-") {
			continue
		}
		entryID := "local-jsonl-entry-" + strings.TrimPrefix(source.Identity.Reference, "local-jsonl-source-")
		meta, present := levels[entryID]
		if !present { continue }
		info := &timeline.WazuhSourceMetadata{
			RuleID: meta.RuleID, AgentID: meta.AgentID, Feed: meta.Feed,
		}
		if meta.Known {
			value := meta.Level
			info.RuleLevel = &value
		}
		source.Wazuh = info
		linked++
	}
	if linked != len(levels) || linked == 0 {
		return timeline.MemoryTimelineData{}, errors.New("Wazuh metadata does not match imported typed sources")
	}
	return merged, nil
}
