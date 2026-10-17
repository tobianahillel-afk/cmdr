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
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || intoCase == "" || withCase == "" || intoCase == withCase ||
		len(intoCase) > 1024 || len(withCase) > 1024 ||
		ownerPath == "" || eventsPath == "" || ownerPath == eventsPath {
		return errors.New("select two distinct Cases and separate local owner/event files")
	}
	original, err := loadDemoSnapshot(ownerPath)
	if err != nil {
		return err
	}
	content, err := readLocalEventFile(eventsPath)
	if err != nil {
		return err
	}
	merged, err := importLocalEvents(ctx, original, principal, intoCase, content)
	if err != nil {
		return err
	}
	imported := len(merged.Entries) - len(original.Entries)
	if imported < 1 || imported > timeline.MaxEntries {
		return errors.New("inconsistent imported event count")
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
