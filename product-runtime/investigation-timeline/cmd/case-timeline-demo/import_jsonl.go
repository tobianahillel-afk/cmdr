package main

import (
	"bytes"
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
	"time"
	"unicode/utf8"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Offline import is an explicit local transformation, not an ingestion
// endpoint. All access labels come from the operator-owned Case snapshot.
type normalizedLocalEvent struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

// Valid IDs are narrow to keep stable records portable across log tools.
func validLocalEventID(id string) bool {
	if len(id) == 0 || len(id) > 96 {
		return false
	}
	for _, c := range []byte(id) {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' ||
			c >= '0' && c <= '9' || c == '.' || c == '_' || c == '-' {
			continue
		}
		return false
	}
	return true
}

// Uses the same local-file boundary as the startup owner snapshot. Never
// accepts HTTP/client paths, symlinks, directories or unbounded input.
func readLocalEventFile(path string) ([]byte, error) {
	if !filepath.IsAbs(path) {
		return nil, errors.New("event JSONL must be an absolute local file")
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxDemoSnapshotBytes {
		return nil, errors.New("event JSONL must be a bounded regular file, not a symlink")
	}
	file, err := os.Open(path) // #nosec G304 -- trusted local CLI file, never HTTP.
	if err != nil {
		return nil, errors.New("cannot open event JSONL")
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(info, opened) {
		return nil, errors.New("event JSONL changed during open")
	}
	data, err := io.ReadAll(io.LimitReader(file, maxDemoSnapshotBytes+1))
	if err != nil || len(data) > maxDemoSnapshotBytes {
		return nil, errors.New("event JSONL exceeds local size limit")
	}
	return data, nil
}

func importLocalEvents(ctx context.Context, snapshot timeline.MemoryTimelineData, principal, caseID string, content []byte) (timeline.MemoryTimelineData, error) {
	return importLocalEventsFromFeed(ctx, snapshot, principal, caseID, content, 0)
}

// Nonzero feed ordinal is used only for multi-feed operator imports.
// Legacy single-file IDs and provenance are preserved unchanged.
func importLocalEventsFromFeed(ctx context.Context, snapshot timeline.MemoryTimelineData, principal, caseID string, content []byte, feed int) (timeline.MemoryTimelineData, error) {
	if feed < 0 || feed > 4 {
		return timeline.MemoryTimelineData{}, errors.New("invalid normalized source feed ordinal")
	}
	if err := ctx.Err(); err != nil {
		return timeline.MemoryTimelineData{}, err
	}
	if principal == "" || caseID == "" || len(caseID) > 1024 || len(content) == 0 || len(content) > maxDemoSnapshotBytes {
		return timeline.MemoryTimelineData{}, errors.New("select an existing authorized Case and a bounded JSONL file")
	}
	// encoding/json replaces invalid UTF-8 with U+FFFD. A forensic import
	// must not silently alter the original identity, message or provenance.
	if !utf8.Valid(content) {
		return timeline.MemoryTimelineData{}, errors.New("event JSONL is not valid UTF-8")
	}
	// Validate the complete existing owner graph before extending it.
	if _, err := newDemoHandlerWithData(strings.Repeat("0", 64), principal, snapshot); err != nil {
		return timeline.MemoryTimelineData{}, errors.New("invalid local owner snapshot")
	}
	selected := -1
	for i, c := range snapshot.Cases {
		if c.ID != caseID {
			continue
		}
		for _, reader := range c.Readers {
			if reader == principal {
				selected = i
				break
			}
		}
	}
	if selected < 0 {
		return timeline.MemoryTimelineData{}, errors.New("Case not available to selected synthetic analyst")
	}
	c := snapshot.Cases[selected]
	rows := bytes.Split(content, []byte("\n"))
	if len(rows) > timeline.MaxEntries+1 || len(c.EntryIDs) >= timeline.MaxEntries {
		return timeline.MemoryTimelineData{}, errors.New("too many Case events")
	}
	seen := make(map[string]bool)
	entryIDs := make(map[string]bool, len(snapshot.Entries))
	sourceIDs := make(map[string]bool, len(snapshot.Sources))
	for _, entry := range snapshot.Entries {
		entryIDs[entry.Value.ID] = true
	}
	for _, src := range snapshot.Sources {
		sourceIDs[src.Identity.Reference] = true
	}
	// Build new slices without mutating the original owner data.
	result := snapshot
	result.Principals = append([]timeline.MemoryTimelinePrincipal(nil), snapshot.Principals...)
	result.Cases = append([]timeline.MemoryTimelineCase(nil), snapshot.Cases...)
	result.Entries = append([]timeline.MemoryTimelineEntry(nil), snapshot.Entries...)
	result.Sources = append([]timeline.MemoryTimelineSource(nil), snapshot.Sources...)
	updatedCase := c
	updatedCase.EntryIDs = append([]string(nil), c.EntryIDs...)
	imported := 0
	for i, line := range rows {
		if err := ctx.Err(); err != nil {
			return timeline.MemoryTimelineData{}, err
		}
		if i == len(rows)-1 && len(line) == 0 { // ordinary trailing newline
			continue
		}
		if len(line) == 0 || len(line) > 4096 {
			return timeline.MemoryTimelineData{}, fmt.Errorf("invalid JSONL line %d", i+1)
		}
		var event normalizedLocalEvent
		dec := json.NewDecoder(bytes.NewReader(line))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&event); err != nil {
			return timeline.MemoryTimelineData{}, fmt.Errorf("invalid normalized event at line %d", i+1)
		}
		var extra any
		if err := dec.Decode(&extra); err != io.EOF {
			return timeline.MemoryTimelineData{}, fmt.Errorf("extra JSON at line %d", i+1)
		}
		if !validLocalEventID(event.ID) || seen[event.ID] ||
			len(event.Message) == 0 || len(event.Message) > 512 ||
			!utf8.ValidString(event.Message) || strings.ContainsAny(event.Message, "\x00\r\n") {
			return timeline.MemoryTimelineData{}, fmt.Errorf("invalid event identity or preview at line %d", i+1)
		}
		at, err := time.Parse(time.RFC3339Nano, event.Timestamp)
		if err != nil {
			return timeline.MemoryTimelineData{}, fmt.Errorf("invalid event timestamp at line %d", i+1)
		}
		seen[event.ID] = true
		if len(updatedCase.EntryIDs) >= timeline.MaxEntries {
			return timeline.MemoryTimelineData{}, errors.New("Case exceeds 5000 entries")
		}
		identity := caseID + "\x00" + event.ID
		if feed > 0 {
			identity = fmt.Sprintf("%s\x00feed-%d\x00%s", caseID, feed, event.ID)
		}
		digest := sha256.Sum256([]byte(identity))
		key := hex.EncodeToString(digest[:])
		entryID := "local-jsonl-entry-" + key
		sourceRef := "local-jsonl-source-" + key
		if entryIDs[entryID] || sourceIDs[sourceRef] {
			return timeline.MemoryTimelineData{}, errors.New("imported identity conflicts with owner snapshot")
		}
		entryIDs[entryID], sourceIDs[sourceRef] = true, true
		provenance := "operator-local-jsonl"
		if feed > 0 {
			provenance = fmt.Sprintf("operator-local-jsonl-feed-%d", feed)
		}
		result.Entries = append(result.Entries, timeline.MemoryTimelineEntry{
			CaseID: caseID, Readers: []string{principal},
			Value: timeline.Entry{
				ID: entryID, Tenant: c.Tenant, Kind: "observed",
				SourceTime: event.Timestamp, Timezone: at.Format("Z07:00"),
				SourceRef: sourceRef, Provenance: provenance,
				Version: 1,
			},
		})
		result.Sources = append(result.Sources, timeline.MemoryTimelineSource{
			Identity: timeline.SourceOwnerIdentity{
				Reference: sourceRef, Kind: "telemetry-event", ObjectID: "local-jsonl-object-" + key,
				CaseID: caseID, Tenant: c.Tenant, Environment: c.Environment,
			},
			Readers: []string{principal}, Preview: event.Message,
		})
		updatedCase.EntryIDs = append(updatedCase.EntryIDs, entryID)
		imported++
	}
	if imported == 0 {
		return timeline.MemoryTimelineData{}, errors.New("JSONL contains no events")
	}
	result.Cases[selected] = updatedCase
	// The existing authoritative-in-this-fixture Case/Entry/source owner
	// constructor must reject any inconsistent derived reference or ACL.
	if _, err := newDemoHandlerWithData(strings.Repeat("0", 64), principal, result); err != nil {
		return timeline.MemoryTimelineData{}, errors.New("imported owner snapshot failed validation")
	}
	return result, nil
}

// Export only the selected authorized Case after full input validation.
// The resulting file remains sensitive because it contains permitted source
// previews, but cannot include other tenants or denied source objects.
func writeDemoImportedJSONL(ctx context.Context, principal, snapshotPath, caseID, eventsPath string, output io.Writer) error {
	if snapshotPath == "" || eventsPath == "" || snapshotPath == eventsPath {
		return errors.New("provide separate absolute owner snapshot and event JSONL files")
	}
	snapshot, err := loadDemoSnapshot(snapshotPath)
	if err != nil {
		return err
	}
	content, err := readLocalEventFile(eventsPath)
	if err != nil {
		return err
	}
	merged, err := importLocalEvents(ctx, snapshot, principal, caseID, content)
	if err != nil {
		return err
	}
	// The legacy stdout JSONL import is ALSO a data export boundary: a
	// caller able to import into one Case may not read other tenants/Cases,
	// or even denied typed sources within their selected Case. Reuse the
	// same first-party minimal owner projection as private disk snapshots.
	// Do not serialize the unfiltered operator-supplied owner graph.
	merged, err = scopedNormalizedCaseSnapshot(merged, principal, caseID)
	if err != nil {
		return err
	}
	payload, err := json.Marshal(merged)
	if err != nil || len(payload)+1 > maxDemoSnapshotBytes {
		return errors.New("resulting local snapshot exceeds 8MiB")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	_, err = output.Write(append(payload, '\n'))
	return err
}


// serveDemoWithImportedJSONL integrates offline event ingestion with the
// existing authenticated analyst UI, without serializing an intermediate
// owner graph containing other tenants or denied source previews.
func serveDemoWithImportedJSONL(ctx context.Context, principal, snapshotPath, caseID, eventsPath string, output io.Writer) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if snapshotPath == "" || eventsPath == "" || snapshotPath == eventsPath {
		return errors.New("select two distinct absolute local owner and event files")
	}
	snapshot, err := loadDemoSnapshot(snapshotPath)
	if err != nil {
		return err
	}
	content, err := readLocalEventFile(eventsPath)
	if err != nil {
		return err
	}
	merged, err := importLocalEvents(ctx, snapshot, principal, caseID, content)
	if err != nil {
		return err
	}
	// The importer's append-only owner graph is the sole authority for which
	// new events belong to this operator import, not an HTTP query or
	// unrelated events in the original Case snapshot.
	imported := make(map[string]struct{}, len(merged.Entries)-len(snapshot.Entries))
	for _, entry := range merged.Entries[len(snapshot.Entries):] {
		imported[entry.Value.ID] = struct{}{}
	}
	// The merged graph remains in memory. Only this import-mode listener
	// offers Event Search; the generic snapshot listener stays read-only.
	return serveDemoWithDataRoute(ctx, principal, merged, output,
		func(token string, base http.Handler) http.Handler {
			return newDemoLocalJSONLSearchHandler(token, principal, caseID, merged, imported, base)
		})
}
