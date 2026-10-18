package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Compare two independently saved, verified Cases without ever combining
// their source owners, raw source bytes, or ACLs into a shared graph.
// This is a retrospective time-window comparison, not evidence correlation
// or attribution of event identities across different Case owners.
func compareVerifiedSavedCases(ctx context.Context, principal, firstPath, firstCase, firstDigest,
	secondPath, secondCase, secondDigest string, output io.Writer) error {
	if ctx == nil || ctx.Err() != nil || principal == "" || firstCase == "" ||
		secondCase == "" || firstCase == secondCase || firstDigest == "" ||
		secondDigest == "" || firstPath == "" || secondPath == "" {
		return errors.New("two distinct saved Cases and independent SHA-256 receipts are required")
	}
	// The SAME first-party integrity and Case/source ACL boundary now
	// guards restore, offline inventory and BOTH comparison operands.
	// Neither owner can inherit access from the other.
	first, _, err := loadAuthorizedImportedCase(ctx, principal, firstPath, firstCase, firstDigest)
	if err != nil { return err }
	second, _, err := loadAuthorizedImportedCase(ctx, principal, secondPath, secondCase, secondDigest)
	if err != nil { return err }

	if !sameComparisonEnvironment(first, second) {
		return errors.New("saved Cases belong to different tenant or environment scopes")
	}
	// Existing first-party triage traverses every owner-authorized Timeline
	// page and canonical typed source independently before any comparison.
	a, err := collectDemoAuthorizedCaseTriage(ctx, principal, first, firstCase)
	if err != nil { return err }
	b, err := collectDemoAuthorizedCaseTriage(ctx, principal, second, secondCase)
	if err != nil { return err }
	if err := ctx.Err(); err != nil { return err }
	result, err := buildCaseComparison(a, b)
	if err != nil { return err }
	result.Format = "cmdr-local-verified-saved-case-comparison-v1"
	// Both owners have already passed independent fingerprint, Case,
	// Timeline and typed-source authorization before any Wazuh metadata
	// is correlated. Generic JSONL snapshots simply produce no cohorts.
	result.SharedWazuhWindows, err = sharedAuthorizedWazuhWindows(first, second, firstCase, secondCase)
	if err != nil { return err }

	// Stage the complete receipt before writing so a denied second Case can
	// never cause a partial report about the first Case.
	var encoded bytes.Buffer
	if err := json.NewEncoder(&encoded).Encode(result); err != nil { return err }
	if err := ctx.Err(); err != nil { return err }
	_, err = output.Write(encoded.Bytes())
	return err
}

func sameComparisonEnvironment(first, second timeline.MemoryTimelineData) bool {
	return len(first.Principals) == 1 && len(second.Principals) == 1 &&
		first.Principals[0].Tenant != "" &&
		first.Principals[0].Environment != "" &&
		first.Principals[0].Tenant == second.Principals[0].Tenant &&
		first.Principals[0].Environment == second.Principals[0].Environment
}
