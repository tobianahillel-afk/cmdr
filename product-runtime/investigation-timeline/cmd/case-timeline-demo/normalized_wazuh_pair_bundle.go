package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
)

// Open a concrete cross-Case Wazuh pair as ONE private artifact. The pair
// ordinal selects an existing authorized correlation result; it grants no
// access to either source. The two first-party owners independently recheck
// the operator-held fingerprints, Case ACL and typed-source authorization.
// The existing atomic no-replace publisher emits nothing until BOTH succeed.
func writeVerifiedCorrelatedPairPrivate(ctx context.Context, principal,
	firstPath, firstCase, firstDigest, secondPath, secondCase, secondDigest string,
	pairOrdinal int, destination string) error {
	if ctx == nil || ctx.Err() != nil || principal == "" ||
		pairOrdinal < 1 || pairOrdinal > 10 || destination == "" ||
		firstCase == secondCase {
		return errors.New("select an authorized cross-Case Wazuh pair and a new private destination")
	}
	return writeAuthorizedReceiptPrivate(ctx, destination, func(output io.Writer) error {
		// Reuse the COMPLETE first-party comparison, including independent
		// digests, tenant/environment equality, Timeline/source traversal,
		// and the canonical <=60s Wazuh correlation. No reimplementation.
		var comparison bytes.Buffer
		if err := compareVerifiedSavedCases(ctx, principal, firstPath, firstCase,
			firstDigest, secondPath, secondCase, secondDigest, &comparison); err != nil {
			return err
		}
		var verified localCaseComparison
		if err := json.Unmarshal(comparison.Bytes(), &verified); err != nil {
			return errors.New("cannot decode independently authorized comparison")
		}
		if verified.First.CaseID != firstCase || verified.Second.CaseID != secondCase ||
			pairOrdinal > len(verified.SharedWazuhEventPairs) {
			return errors.New("selected Wazuh pair is unavailable in both authorized Cases")
		}
		pair := verified.SharedWazuhEventPairs[pairOrdinal-1]

		// Entry selectors may be exposed in a correlation receipt but are
		// never an authorization grant. Re-open EACH bounded snapshot
		// independently and require its canonical imported Entry identity.
		first, firstIDs, err := loadAuthorizedImportedCase(ctx, principal, firstPath, firstCase, firstDigest)
		if err != nil { return err }
		second, secondIDs, err := loadAuthorizedImportedCase(ctx, principal, secondPath, secondCase, secondDigest)
		if err != nil { return err }
		if !sameComparisonEnvironment(first, second) {
			return errors.New("correlated source owners are not in the same tenant and environment")
		}
		if _, ok := firstIDs[pair.FirstEntryID]; !ok {
			return errors.New("first paired entry is not available in selected Case")
		}
		if _, ok := secondIDs[pair.SecondEntryID]; !ok {
			return errors.New("second paired entry is not available in selected Case")
		}
		// The existing source handler independently checks selected Case,
		// entry and typed source on both sides, including source read ACL.
		firstInspection, err := inspectDemoEntryWithData(ctx, principal, first, firstCase, pair.FirstEntryID)
		if err != nil { return err }
		secondInspection, err := inspectDemoEntryWithData(ctx, principal, second, secondCase, pair.SecondEntryID)
		if err != nil { return err }
		if firstInspection.CaseID != firstCase || firstInspection.EntryID != pair.FirstEntryID ||
			secondInspection.CaseID != secondCase || secondInspection.EntryID != pair.SecondEntryID ||
			!firstInspection.Synthetic || !secondInspection.Synthetic {
			return errors.New("paired source owner returned inconsistent authorized identity")
		}
		if err := ctx.Err(); err != nil { return err }
		return json.NewEncoder(output).Encode(struct {
			Format string `json:"format"`
			Synthetic bool `json:"synthetic"`
			Pair sharedAuthorizedWazuhEventPair `json:"pair"`
			First localEntryInspection `json:"first"`
			Second localEntryInspection `json:"second"`
		}{
			Format: "cmdr-private-verified-wazuh-pair-v1",
			Synthetic: true, Pair: pair,
			First: firstInspection, Second: secondInspection,
		})
	})
}
