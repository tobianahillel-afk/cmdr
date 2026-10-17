package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// localImportedComparison is a deliberately aggregate-only local analysis.
// The imported owner data and source previews never enter the output.
type localImportedComparison struct {
	Format          string              `json:"format"`
	Synthetic       bool                `json:"synthetic"`
	ImportedEntries int                 `json:"imported_entries"`
	Comparison      localCaseComparison `json:"comparison"`
}

// writeDemoImportedJSONLComparison connects normalized offline telemetry to
// the existing owner-authorized two-Case comparison in one process. Nothing is
// emitted unless importing and BOTH independent Case reads have succeeded.
func writeDemoImportedJSONLComparison(ctx context.Context, principal, ownerPath, intoCase, withCase, eventsPath string, out io.Writer) error {
	merged, imported, err := prepareImportedAuthorizedComparison(ctx, principal, ownerPath, intoCase, withCase, eventsPath)
	if err != nil {
		return err
	}
	comparison, err := collectAuthorizedCaseComparison(ctx, principal, merged, intoCase, withCase)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(localImportedComparison{
		Format: "cmdr-local-imported-case-comparison-v1",
		Synthetic: true, ImportedEntries: imported, Comparison: comparison,
	})
}

 // prepareImportedAuthorizedComparison performs the bounded local import once
 // for both aggregate comparison and temporal lead modes. It never emits the
 // merged, potentially cross-tenant owner graph.
func prepareImportedAuthorizedComparison(ctx context.Context, principal, ownerPath, intoCase, withCase, eventsPath string) (timeline.MemoryTimelineData, int, error) {
	if err := ctx.Err(); err != nil {
		return timeline.MemoryTimelineData{}, 0, err
	}
	if principal == "" || intoCase == "" || withCase == "" || intoCase == withCase ||
		len(intoCase) > 1024 || len(withCase) > 1024 ||
		ownerPath == "" || eventsPath == "" || ownerPath == eventsPath {
		return timeline.MemoryTimelineData{}, 0, errors.New("select two distinct Cases and separate local owner/event files")
	}
	original, err := loadDemoSnapshot(ownerPath)
	if err != nil {
		return timeline.MemoryTimelineData{}, 0, err
	}
	content, err := readLocalEventFile(eventsPath)
	if err != nil {
		return timeline.MemoryTimelineData{}, 0, err
	}
	merged, err := importLocalEvents(ctx, original, principal, intoCase, content)
	if err != nil {
		return timeline.MemoryTimelineData{}, 0, err
	}
	imported := len(merged.Entries) - len(original.Entries)
	if imported < 1 || imported > timeline.MaxEntries {
		return timeline.MemoryTimelineData{}, 0, errors.New("inconsistent imported event count")
	}

	return merged, imported, nil
}

// localImportedProximity is an explicitly requested sensitive, unsigned
// temporal hint; it is not a causal conclusion or chain-of-custody record.
type localImportedProximity struct {
	Format          string             `json:"format"`
	Synthetic       bool               `json:"synthetic"`
	ImportedEntries int                `json:"imported_entries"`
	Proximity       localCaseProximity `json:"proximity"`
}

// writeDemoImportedJSONLProximity composes the already implemented import,
// independent HTTP Case/typed-source ACL reads and certain-time proximity.
// Nothing reaches output until both Case traversals and selection succeed.
func writeDemoImportedJSONLProximity(ctx context.Context, principal, ownerPath, intoCase, withCase, eventsPath string, seconds int, out io.Writer) error {
	if seconds < 1 || seconds > 300 {
		return errors.New("temporal proximity window must be 1..300 seconds")
	}
	merged, imported, err := prepareImportedAuthorizedComparison(ctx, principal, ownerPath, intoCase, withCase, eventsPath)
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
	sequence, err := collectAuthorizedCrossCaseSequence(ctx, handler, token, intoCase, withCase)
	if err != nil {
		return err
	}
	leads, err := selectCaseProximity(ctx, sequence, seconds)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(localImportedProximity{
		Format: "cmdr-local-imported-case-proximity-v1",
		Synthetic: true, ImportedEntries: imported, Proximity: leads,
	})
}
