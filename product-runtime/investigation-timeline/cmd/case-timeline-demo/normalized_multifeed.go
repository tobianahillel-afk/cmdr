package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Import an explicit operator-selected set of independent normalized feeds.
// Every file is checked by the same bounded regular-file reader, and every
// appended event still goes through the existing authoritative Case owner.
func loadMultipleNormalizedJSONLFeeds(ctx context.Context, principal, snapshotPath, caseID string, paths []string) (
	timeline.MemoryTimelineData, map[string]struct{}, error,
) {
	if ctx == nil || len(paths) < 2 || len(paths) > 4 || snapshotPath == "" || caseID == "" {
		return timeline.MemoryTimelineData{}, nil, errors.New("select 2 to 4 normalized JSONL feeds and a Case snapshot")
	}
	original, err := loadDemoSnapshot(snapshotPath)
	if err != nil {
		return timeline.MemoryTimelineData{}, nil, err
	}
	merged := original
	seenFiles := make([]os.FileInfo, 0, len(paths))
	totalBytes := 0
	for i, path := range paths {
		if err := ctx.Err(); err != nil {
			return timeline.MemoryTimelineData{}, nil, err
		}
		if path == snapshotPath {
			return timeline.MemoryTimelineData{}, nil, errors.New("normalized feed cannot reuse the owner snapshot")
		}
		// Reject aliases and hard links, not just repeated path strings,
		// to avoid manufacturing a second feed out of the same evidence.
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() {
			return timeline.MemoryTimelineData{}, nil, errors.New("normalized feed must be a regular local file")
		}
		for _, previous := range seenFiles {
			if os.SameFile(info, previous) {
				return timeline.MemoryTimelineData{}, nil, errors.New("normalized source file selected more than once")
			}
		}
		seenFiles = append(seenFiles, info)
		content, err := readLocalEventFile(path)
		if err != nil {
			return timeline.MemoryTimelineData{}, nil, err
		}
		totalBytes += len(content)
		if totalBytes > maxDemoSnapshotBytes {
			return timeline.MemoryTimelineData{}, nil, errors.New("combined normalized feeds exceed 8MiB input")
		}
		// Feed ordinal partitions original event IDs without rewriting,
		// mutating or trusting the original forensic record.
		merged, err = importLocalEventsFromFeed(ctx, merged, principal, caseID, content, i+1)
		if err != nil {
			return timeline.MemoryTimelineData{}, nil, fmt.Errorf("normalized feed %d: %w", i+1, err)
		}
	}
	imported := make(map[string]struct{}, len(merged.Entries)-len(original.Entries))
	for _, entry := range merged.Entries[len(original.Entries):] {
		if _, exists := imported[entry.Value.ID]; exists {
			return timeline.MemoryTimelineData{}, nil, errors.New("duplicate derived normalized source identity")
		}
		imported[entry.Value.ID] = struct{}{}
	}
	if len(imported) == 0 {
		return timeline.MemoryTimelineData{}, nil, errors.New("no normalized events imported")
	}
	return merged, imported, nil
}

// Keeps all feeds on the same existing loopback Case/Timeline/source owner.
// No persistent file system writing, external listener or ingestion service.
func serveDemoWithMultipleNormalizedJSONL(ctx context.Context, principal, snapshotPath, caseID string,
	paths []string, output io.Writer,
) error {
	data, imported, err := loadMultipleNormalizedJSONLFeeds(ctx, principal, snapshotPath, caseID, paths)
	if err != nil {
		return err
	}
	return serveDemoWithDataRoute(ctx, principal, data, output, func(token string, owner http.Handler) http.Handler {
		return newDemoLocalJSONLSearchHandler(token, principal, caseID, data, imported, owner)
	})
}
