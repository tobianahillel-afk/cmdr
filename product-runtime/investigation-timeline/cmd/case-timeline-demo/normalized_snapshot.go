package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Save only a single Case and records readable by this synthetic principal.
// The full source snapshot may contain unrelated or denied tenant data.
func scopedNormalizedCaseSnapshot(data timeline.MemoryTimelineData, principal, caseID string) (timeline.MemoryTimelineData, error) {
	if _, err := newDemoHandlerWithData(strings.Repeat("0", 64), principal, data); err != nil {
		return timeline.MemoryTimelineData{}, errors.New("input owner graph invalid")
	}
	readable := func(ids []string) bool {
		for _, id := range ids {
			if id == principal { return true }
		}
		return false
	}
	var current timeline.MemoryTimelinePrincipal
	for _, p := range data.Principals {
		if p.ID == principal { current = p; break }
	}
	var chosen timeline.MemoryTimelineCase
	found := false
	for _, c := range data.Cases {
		if c.ID != caseID { continue }
		if current.ID == "" || c.Tenant != current.Tenant ||
			c.Environment != current.Environment || !readable(c.Readers) {
			return timeline.MemoryTimelineData{}, errors.New("Case unavailable to selected principal")
		}
		chosen = c
		found = true
		break
	}
	if !found { return timeline.MemoryTimelineData{}, errors.New("selected Case absent") }
	entryByID := make(map[string]timeline.MemoryTimelineEntry)
	sourceByRef := make(map[string]timeline.MemoryTimelineSource)
	for _, entry := range data.Entries {
		if entry.CaseID == caseID { entryByID[entry.Value.ID] = entry }
	}
	for _, src := range data.Sources {
		if src.Identity.CaseID == caseID { sourceByRef[src.Identity.Reference] = src }
	}
	limited := timeline.MemoryTimelineData{
		Principals: []timeline.MemoryTimelinePrincipal{current},
		Cases: []timeline.MemoryTimelineCase{chosen},
	}
	limited.Cases[0].Readers = []string{principal}
	limited.Cases[0].EntryIDs = make([]string, 0, len(chosen.EntryIDs))
	addedSource := make(map[string]bool)
	for _, id := range chosen.EntryIDs {
		entry, ok := entryByID[id]
		if !ok { return timeline.MemoryTimelineData{}, errors.New("Case entry missing") }
		src, found := sourceByRef[entry.Value.SourceRef]
		if !readable(entry.Readers) || !found || !readable(src.Readers) {
			continue // denied metadata and source previews are never saved
		}
		entry.Readers = []string{principal}
		limited.Entries = append(limited.Entries, entry)
		limited.Cases[0].EntryIDs = append(limited.Cases[0].EntryIDs, id)
		if !addedSource[src.Identity.Reference] {
			src.Readers = []string{principal}
			limited.Sources = append(limited.Sources, src)
			addedSource[src.Identity.Reference] = true
		}
	}
	if _, err := newDemoHandlerWithData(strings.Repeat("0", 64), principal, limited); err != nil {
		return timeline.MemoryTimelineData{}, errors.New("scoped Case owner invalid")
	}
	return limited, nil
}

// Explicit private persistence of only the selected authorized Case owner.
func writeDemoPersistedNormalizedCase(ctx context.Context, principal, snapshotPath, caseID string,
	paths []string, destination string, output io.Writer) error {
	if ctx == nil || ctx.Err() != nil || !filepath.IsAbs(destination) ||
		len(paths) < 1 || len(paths) > 4 || snapshotPath == "" || caseID == "" {
		return errors.New("select an absolute private output path, a Case and 1..4 normalized sources")
	}
	var data timeline.MemoryTimelineData
	var err error
	if len(paths) == 1 {
		if paths[0] == snapshotPath {
			return errors.New("normalized source cannot be the owner snapshot")
		}
		var initial timeline.MemoryTimelineData
		initial, err = loadDemoSnapshot(snapshotPath)
		if err == nil {
			var file []byte
			file, err = readLocalEventFile(paths[0])
			if err == nil {
				data, err = importLocalEvents(ctx, initial, principal, caseID, file)
			}
		}
	} else {
		data, _, err = loadMultipleNormalizedJSONLFeeds(ctx, principal, snapshotPath, caseID, paths)
	}
	if err != nil { return err }
	data, err = scopedNormalizedCaseSnapshot(data, principal, caseID)
	if err != nil { return err }
	if _, err := importedNormalizedIDs(data, caseID); err != nil { return err }
	// Marshal the minimal authorized Case graph before any output is opened.
	payload, err := json.Marshal(data)
	if err != nil || len(payload)+1 > maxDemoSnapshotBytes {
		return errors.New("persisted normalized Case exceeds bounded snapshot size")
	}
	payload = append(payload, '\n')
	if err := writePrivateSnapshotAtomic(ctx, destination, payload); err != nil {
		return err
	}
	// Digest is useful for accidental corruption checks, NOT a signature or
	// a trustworthy custody assertion; output never contains source payload.
	sum := sha256.Sum256(payload)
	_, err = fmt.Fprintf(output,
		"Sensitive local Case owner snapshot saved to %q (sha256 %x); no bearer stored; this is NOT signed forensic custody.\n",
		destination, sum)
	return err
}

// An O_EXCL hard-link publication is an atomic no-replace transition from a
// private, fully written temporary file. Never overwrite an existing target,
// symlink, or hardlink; never allow HTTP clients to choose this destination.
func writePrivateSnapshotAtomic(ctx context.Context, path string, data []byte) error {
	if ctx == nil || ctx.Err() != nil || !filepath.IsAbs(path) ||
		len(data) == 0 || len(data) > maxDemoSnapshotBytes {
		return errors.New("invalid private owner snapshot output request")
	}
	dir := filepath.Dir(path)
	if _, err := os.Lstat(path); err == nil {
		return errors.New("owner snapshot destination already exists")
	} else if !os.IsNotExist(err) {
		return errors.New("cannot inspect owner snapshot destination")
	}
	// CreateTemp uses mode 0600 on POSIX; Windows filesystem ACL behavior
	// still depends on the operator's selected directory.
	file, err := os.CreateTemp(dir, ".cmdr-private-case-*.partial")
	if err != nil {
		return errors.New("cannot create private snapshot in destination directory")
	}
	tempName := file.Name()
	defer os.Remove(tempName)
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return errors.New("cannot write complete private snapshot")
	}
	if err := file.Sync(); err != nil {
		return errors.New("cannot synchronize private snapshot data")
	}
	if err := file.Close(); err != nil {
		return errors.New("cannot close private snapshot")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	// os.Link fails if path exists, without replacing it. A regular rename
	// could silently overwrite the existing operator evidence file.
	if err := os.Link(tempName, path); err != nil {
		return errors.New("cannot atomically create new private snapshot destination")
	}
	return nil
}

// A reopened snapshot never needs the original raw JSONL files. Eligible
// source refs are recovered from existing first-party normalized import
// provenance and stable owner identities, not an attacker-supplied HTTP list.
// This does NOT attest original JSONL file bytes or make the snapshot tamperproof.
func importedNormalizedIDs(data timeline.MemoryTimelineData, caseID string) (map[string]struct{}, error) {
	ids := make(map[string]struct{})
	for _, entry := range data.Entries {
		if entry.CaseID != caseID {
			continue
		}
		provenance := entry.Value.Provenance
		if provenance != "operator-local-jsonl" &&
			provenance != "operator-local-jsonl-feed-1" &&
			provenance != "operator-local-jsonl-feed-2" &&
			provenance != "operator-local-jsonl-feed-3" &&
			provenance != "operator-local-jsonl-feed-4" {
			continue
		}
		id := strings.TrimPrefix(entry.Value.ID, "local-jsonl-entry-")
		if len(id) != 64 {
			return nil, errors.New("persisted normalized identity is not recognized")
		}
		if _, err := hex.DecodeString(id); err != nil ||
			entry.Value.SourceRef != "local-jsonl-source-"+id ||
			entry.Value.Kind != "observed" {
			return nil, errors.New("persisted normalized source identity is inconsistent")
		}
		if _, duplicate := ids[entry.Value.ID]; duplicate {
			return nil, errors.New("persisted normalized entry identity repeated")
		}
		ids[entry.Value.ID] = struct{}{}
	}
	if len(ids) == 0 {
		return nil, errors.New("no normalized imported entries for the selected Case")
	}
	return ids, nil
}

// One first-party trust boundary for all three saved-Case paths:
// authenticated browser reopen, no-listener verification and offline
// cross-archive investigation comparison. Always validate the exact file
// bytes (when an operator fingerprint is supplied), project only the
// selected principal's Case and sources, and recognize canonical imports.
// This is deliberately limited to our local synthetic owner, not generic IAM.
func loadAuthorizedImportedCase(ctx context.Context, principal, snapshotPath, caseID, expectedDigest string) (
	timeline.MemoryTimelineData, map[string]struct{}, error,
) {
	if ctx == nil || ctx.Err() != nil || principal == "" || snapshotPath == "" || caseID == "" {
		return timeline.MemoryTimelineData{}, nil, errors.New("select a readable saved Case and principal")
	}
	data, err := loadDemoSnapshotWithExpectedDigest(snapshotPath, expectedDigest)
	if err != nil {
		return timeline.MemoryTimelineData{}, nil, err
	}
	// A legacy unscoped operator snapshot is never passed intact to a
	// downstream listener or report, even if its JSON is otherwise valid.
	data, err = scopedNormalizedCaseSnapshot(data, principal, caseID)
	if err != nil {
		return timeline.MemoryTimelineData{}, nil, err
	}
	imported, err := importedNormalizedIDs(data, caseID)
	if err != nil {
		return timeline.MemoryTimelineData{}, nil, err
	}
	if err := ctx.Err(); err != nil {
		return timeline.MemoryTimelineData{}, nil, err
	}
	return data, imported, nil
}

// Reuse exactly the existing Case/Timeline/typed source HTTP and SearchJob
// implementation. Each session generates a NEW short-lived bearer token;
// neither bearer nor SearchJob results are stored in the snapshot.
func serveDemoReopenedNormalizedCase(ctx context.Context, principal, snapshotPath, caseID, expectedDigest string, output io.Writer) error {
	data, ids, err := loadAuthorizedImportedCase(ctx, principal, snapshotPath, caseID, expectedDigest)
	if err != nil {
		return err
	}
	return serveDemoWithDataRoute(ctx, principal, data, output, func(token string, owner http.Handler) http.Handler {
		return newDemoLocalJSONLSearchHandler(token, principal, caseID, data, ids, owner)
	})
}
