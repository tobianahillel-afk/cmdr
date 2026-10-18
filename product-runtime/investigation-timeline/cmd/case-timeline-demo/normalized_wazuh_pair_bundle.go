package main

import (
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
	return writeVerifiedCorrelatedPairPrivateWithGenericEvidence(ctx, principal,
		firstPath, firstCase, firstDigest, secondPath, secondCase, secondDigest,
		pairOrdinal, destination, false)
}

// Explicitly expand source-owned generic leads to their sensitive previews
// only inside an operator-selected NEW private report. Ordinary pair inspection
// retains its metadata-only generic leads and the v2 report format.
func writeVerifiedCorrelatedPairPrivateWithGenericEvidence(ctx context.Context, principal,
	firstPath, firstCase, firstDigest, secondPath, secondCase, secondDigest string,
	pairOrdinal int, destination string, includeGenericEvidence bool) error {
	return writeVerifiedCorrelatedPairPrivateSelected(ctx, principal,
		firstPath, firstCase, firstDigest, secondPath, secondCase, secondDigest,
		pairOrdinal, "", "", destination, includeGenericEvidence)
}

// Entry IDs identify the same currently authorized correlation lead, not
// permission to read either source. Both owners are always reauthorized.
func writeVerifiedCorrelatedPairPrivateByEntryIDs(ctx context.Context, principal,
	firstPath, firstCase, firstDigest, secondPath, secondCase, secondDigest string,
	firstEntryID, secondEntryID, destination string, includeGenericEvidence bool) error {
	return writeVerifiedCorrelatedPairPrivateSelected(ctx, principal,
		firstPath, firstCase, firstDigest, secondPath, secondCase, secondDigest,
		0, firstEntryID, secondEntryID, destination, includeGenericEvidence)
}

func writeVerifiedCorrelatedPairPrivateSelected(ctx context.Context, principal,
	firstPath, firstCase, firstDigest, secondPath, secondCase, secondDigest string,
	pairOrdinal int, firstEntryID, secondEntryID, destination string,
	includeGenericEvidence bool) error {
	byID := pairOrdinal == 0 && firstEntryID != "" && secondEntryID != ""
	byOrdinal := pairOrdinal >= 1 && pairOrdinal <= 10 &&
		firstEntryID == "" && secondEntryID == ""
	if ctx == nil || ctx.Err() != nil || principal == "" ||
		(!byID && !byOrdinal) || destination == "" || firstCase == secondCase {
		return errors.New("select an authorized cross-Case Wazuh pair and a new private destination")
	}
	return writeAuthorizedReceiptPrivate(ctx, destination, func(output io.Writer) error {
		// Independently verify/load each snapshot once. The very same owner
		// bytes and typed sources will be used for BOTH correlation and the
		// final drilldown: no repeated disk reads or owner pagination.
		first, firstIDs, err := loadAuthorizedImportedCase(ctx, principal, firstPath, firstCase, firstDigest)
		if err != nil { return err }
		second, secondIDs, err := loadAuthorizedImportedCase(ctx, principal, secondPath, secondCase, secondDigest)
		if err != nil { return err }
		verified, err := compareLoadedAuthorizedSavedCases(ctx, principal, first, second, firstCase, secondCase)
		if err != nil { return err }
		if verified.First.CaseID != firstCase || verified.Second.CaseID != secondCase {
			return errors.New("correlated Case identity not authorized")
		}
		var pair sharedAuthorizedWazuhEventPair
		if byID {
			found := false
			for _, candidate := range verified.SharedWazuhEventPairs {
				if candidate.FirstEntryID == firstEntryID &&
					candidate.SecondEntryID == secondEntryID {
					pair, found = candidate, true
					break
				}
			}
			if !found {
				return errors.New("selected Wazuh Entry pair is not available in authorized Cases")
			}
		} else {
			if pairOrdinal > len(verified.SharedWazuhEventPairs) {
				return errors.New("selected Wazuh pair is unavailable in both authorized Cases")
			}
			pair = verified.SharedWazuhEventPairs[pairOrdinal-1]
		}

		// Selected correlation IDs are selectors, never access grants.
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
		// One first-party source reauthorization per selected generic lead.
		// Default v2 retains only metadata; explicit private v3 reuses the
		// very same already authorized source inspection, without reopening.
		var firstGenericEvidence, secondGenericEvidence []localEntryInspection
		var firstPrivate, secondPrivate *[]localEntryInspection
		if includeGenericEvidence {
			firstPrivate, secondPrivate = &firstGenericEvidence, &secondGenericEvidence
		}
		firstContext, err := nearbyAuthorizedGenericPairContext(
			ctx, principal, first, firstCase, firstIDs, pair.FirstTimeUTC, firstPrivate)
		if err != nil { return err }
		secondContext, err := nearbyAuthorizedGenericPairContext(
			ctx, principal, second, secondCase, secondIDs, pair.SecondTimeUTC, secondPrivate)
		if err != nil { return err }
		format := "cmdr-private-verified-wazuh-pair-v2"
		if includeGenericEvidence {
			format = "cmdr-private-verified-wazuh-pair-v3"
		}
		if err := ctx.Err(); err != nil { return err }
		return json.NewEncoder(output).Encode(struct {
			Format string `json:"format"`
			Synthetic bool `json:"synthetic"`
			Pair sharedAuthorizedWazuhEventPair `json:"pair"`
			First localEntryInspection `json:"first"`
			Second localEntryInspection `json:"second"`
			FirstGenericContext authorizedGenericPairContextSummary `json:"first_generic_context"`
			SecondGenericContext authorizedGenericPairContextSummary `json:"second_generic_context"`
			FirstGenericEvidence []localEntryInspection `json:"first_generic_evidence,omitempty"`
			SecondGenericEvidence []localEntryInspection `json:"second_generic_evidence,omitempty"`
		}{
			Format: format,
			Synthetic: true, Pair: pair,
			First: firstInspection, Second: secondInspection,
			FirstGenericContext: firstContext, SecondGenericContext: secondContext,
			FirstGenericEvidence: firstGenericEvidence, SecondGenericEvidence: secondGenericEvidence,
		})
	})
}

