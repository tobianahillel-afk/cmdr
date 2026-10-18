package main

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Save Wazuh alerts into the SAME private, scoped Case owner snapshot that
// already supports verified offline reopen, investigation, and comparison.
// The Wazuh normalizer is the existing trusted parser: raw Wazuh JSON,
// agent IP/hostname and original full_log are not copied as owner metadata.
// The retained preview is the rule description, or the bounded full_log
// fallback only when no rule description was present.
func writePersistedWazuhCase(ctx context.Context, principal, snapshotPath, caseID string,
	alertFiles []string, destination string, output io.Writer) error {
	if ctx == nil || ctx.Err() != nil || principal == "" || caseID == "" ||
		snapshotPath == "" || !filepath.IsAbs(destination) ||
		len(alertFiles) < 1 || len(alertFiles) > maxWazuhAlertFiles {
		return errors.New("select a Case, 1..4 Wazuh files and a new absolute private destination")
	}
	// The output cannot alias original evidence or the owner source.
	// Publication must never overwrite even if an operator accidentally
	// selects a different hard-link path for an existing evidence file.
	var files = append([]string{snapshotPath}, alertFiles...)
	for _, file := range files {
		if !filepath.IsAbs(file) {
			return errors.New("Wazuh owner and evidence paths must be absolute")
		}
		if file == destination {
			return errors.New("cannot replace the original Wazuh or Case source")
		}
	}
	var merged timeline.MemoryTimelineData
	if len(alertFiles) == 1 {
		source, err := os.Lstat(snapshotPath)
		if err != nil || !source.Mode().IsRegular() {
			return errors.New("Wazuh Case owner must be a regular file")
		}
		alert, err := os.Lstat(alertFiles[0])
		if err != nil || !alert.Mode().IsRegular() ||
			os.SameFile(source, alert) {
			return errors.New("Wazuh alerts must be independent of the Case owner")
		}
		merged, _, err = prepareWazuh4Case(ctx, principal, snapshotPath, caseID, alertFiles[0])
		if err != nil {
			return err
		}
	} else {
		var err error
		merged, _, err = loadMultipleWazuhFeeds(ctx, principal, snapshotPath, caseID, alertFiles)
		if err != nil {
			return err
		}
	}
	// No other file output or stdout before the entire Wazuh import and
	// scoped owner validation succeeds. Existing no-replace atomic writer
	// enforces the private destination filesystem boundary.
	return persistScopedCaseOwner(ctx, merged, principal, caseID, destination, output)
}
