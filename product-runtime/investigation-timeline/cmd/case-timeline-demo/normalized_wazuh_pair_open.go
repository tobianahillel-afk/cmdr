package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
)

// Drill down from a concrete first_entry_id / second_entry_id in an
// independently authorized cross-Case correlation receipt. The operator
// chooses the corresponding Case snapshot explicitly; the entry ID from a
// comparison is only a selector, NEVER an access grant.
//
// Source previews are sensitive. The existing private receipt publisher
// writes ONLY to a new absolute file, atomically without replacing evidence,
// and never falls back to stdout, a browser listener, or an HTTP redirect.
func writeVerifiedSavedCaseEntryPrivate(ctx context.Context, principal, snapshotPath,
	caseID, expectedSHA256, entryID, destination string) error {
	if ctx == nil || ctx.Err() != nil || principal == "" ||
		snapshotPath == "" || caseID == "" || expectedSHA256 == "" ||
		entryID == "" || destination == "" {
		return errors.New("explicit authorized Case, fingerprint, entry and private destination required")
	}
	return writeAuthorizedReceiptPrivate(ctx, destination, func(output io.Writer) error {
		// EXACT same bounded single-read digest + principal/Case/tenant/
		// environment + typed-source owner checks as the resumed Case
		// viewer, verification, and cross-Case comparison.
		data, imported, err := loadAuthorizedImportedCase(ctx, principal,
			snapshotPath, caseID, expectedSHA256)
		if err != nil {
			return err
		}
		// A caller-supplied Entry ID must belong to the selected Case's
		// canonical normalized import before any preview may be opened.
		if _, ok := imported[entryID]; !ok {
			return errors.New("selected imported event not available in authorized Case")
		}
		// Do not reimplement source authorization. The first-party owner
		// independently validates Case + Entry + typed Source each time.
		inspection, err := inspectDemoEntryWithData(ctx, principal, data, caseID, entryID)
		if err != nil {
			return err
		}
		if inspection.CaseID != caseID || inspection.EntryID != entryID || !inspection.Synthetic {
			return errors.New("source inspection did not return selected authorized identity")
		}
		return json.NewEncoder(output).Encode(inspection)
	})
}
