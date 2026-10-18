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

	result, err := compareLoadedAuthorizedSavedCases(ctx, principal, first, second, firstCase, secondCase)
	if err != nil { return err }

	// Stage the complete receipt before writing so a denied second Case can
	// never cause a partial report about the first Case.
	var encoded bytes.Buffer
	if err := json.NewEncoder(&encoded).Encode(result); err != nil { return err }
	if err := ctx.Err(); err != nil { return err }
	// Do not report a successful comparison when the destination accepted
	// fewer bytes than the fully authorized result.
	_, err = io.Copy(output, &encoded)
	return err
}


 // Compare already verified and *independently authorized* owners. This
 // single shared path powers both the metadata-only CLI comparison and the
 // private two-source drilldown, avoiding a second disk read/owner traversal.
 // The caller must have loaded each Case through loadAuthorizedImportedCase.
func compareLoadedAuthorizedSavedCases(ctx context.Context, principal string,
	first, second timeline.MemoryTimelineData, firstCase, secondCase string) (localCaseComparison, error) {
	if ctx == nil || ctx.Err() != nil {
		return localCaseComparison{}, errors.New("comparison cancelled")
	}
	if principal == "" || firstCase == "" || secondCase == "" || firstCase == secondCase ||
		!sameComparisonEnvironment(first, second) {
		return localCaseComparison{}, errors.New("saved Cases require distinct IDs and a common authorized tenant/environment")
	}
	// Each owner is independently paginated through existing Case/Timeline
	// and typed-source authorization before any Wazuh metadata is correlated.
	a, err := collectDemoAuthorizedCaseTriage(ctx, principal, first, firstCase)
	if err != nil { return localCaseComparison{}, err }
	b, err := collectDemoAuthorizedCaseTriage(ctx, principal, second, secondCase)
	if err != nil { return localCaseComparison{}, err }
	if err := ctx.Err(); err != nil { return localCaseComparison{}, err }
	result, err := buildCaseComparison(a, b)
	if err != nil { return localCaseComparison{}, err }
	result.Format = "cmdr-local-verified-saved-case-comparison-v1"
	// Both temporal views must use exactly the SAME imported/authorized
	// Wazuh owner observations. Previously each collected and sorted both
	// Cases separately, quadrupling the source walk for one comparison.
	firstWazuh, err := collectAuthorizedWazuhObservations(first, firstCase)
	if err != nil { return localCaseComparison{}, err }
	secondWazuh, err := collectAuthorizedWazuhObservations(second, secondCase)
	if err != nil { return localCaseComparison{}, err }
	// Pure derivations cannot fail after both independently authorized
	// source cohorts have been collected and validated. No phantom error
	// branches or duplicate security decisions at this presentation layer.
	result.SharedWazuhWindows = sharedWazuhWindowsFromObservations(firstWazuh, secondWazuh)
	result.SharedWazuhEventPairs = sharedWazuhPairsFromObservations(firstWazuh, secondWazuh)
	if err := ctx.Err(); err != nil { return localCaseComparison{}, err }
	return result, nil
}

func sameComparisonEnvironment(first, second timeline.MemoryTimelineData) bool {
	return len(first.Principals) == 1 && len(second.Principals) == 1 &&
		first.Principals[0].Tenant != "" &&
		first.Principals[0].Environment != "" &&
		first.Principals[0].Tenant == second.Principals[0].Tenant &&
		first.Principals[0].Environment == second.Principals[0].Environment
}
