package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

)

// Persist one authorized Case with both source families in ONE atomic write.
// Reuse both proven parsers and the existing Case/Timeline/typed-source owner;
// there is no intermediate private Case archive and no new source permission.
func writePersistedMixedCase(ctx context.Context, principal, snapshotPath, caseID string,
	wazuhFiles, normalizedFiles []string, destination string, output io.Writer) error {
	if ctx == nil || ctx.Err() != nil || principal == "" || caseID == "" ||
		!filepath.IsAbs(snapshotPath) || !filepath.IsAbs(destination) ||
		len(wazuhFiles) < 1 || len(wazuhFiles) > maxWazuhAlertFiles ||
		len(normalizedFiles) < 1 || len(normalizedFiles) > 4 {
		return errors.New("select one Case, 1..4 Wazuh files, 1..4 normalized JSONL files and new private output")
	}

	// The inputs are the only selected evidence. Distinct real files are
	// required so an alias cannot manufacture multiple independent feeds.
	all := append([]string{snapshotPath}, wazuhFiles...)
	all = append(all, normalizedFiles...)
	seen := make([]os.FileInfo, 0, len(all))
	var total int64
	for _, path := range all {
		if !filepath.IsAbs(path) || path == destination {
			return errors.New("all input evidence paths must be separate absolute paths")
		}
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() {
			return errors.New("Case and evidence inputs must be independent regular files")
		}
		for _, old := range seen {
			if os.SameFile(old, info) {
				return errors.New("duplicated or hardlinked mixed Case evidence")
			}
		}
		seen = append(seen, info)
		if len(seen) > 1 {
			total += info.Size()
			if total > maxDemoSnapshotBytes {
				return errors.New("combined mixed Case input exceeds 8MiB")
			}
		}
	}

	merged, err := preparePersistedWazuhCase(ctx, principal, snapshotPath, caseID, wazuhFiles)
	if err != nil { return err }

	// One generic feed is intentionally given the same deterministic feed-1
	// namespace as multiple generic feeds. A collision with an existing
	// Wazuh-derived Entry is rejected by the existing Case owner importer;
	// source metadata or alert IDs never grant import or read permission.
	for i, path := range normalizedFiles {
		if err := ctx.Err(); err != nil { return err }
		raw, err := readLocalEventFile(path)
		if err != nil { return err }
		merged, err = importLocalEventsFromFeed(ctx, merged, principal, caseID, raw, i+1)
		if err != nil {
			return fmt.Errorf("generic feed %d: %w", i+1, err)
		}
	}
	// Same single trust boundary as existing individual Wazuh and JSONL
	// saves; failed normalization, ACL validation or disk publication cannot
	// create a partly committed merged owner.
	return persistScopedCaseOwner(ctx, merged, principal, caseID, destination, output)
}

