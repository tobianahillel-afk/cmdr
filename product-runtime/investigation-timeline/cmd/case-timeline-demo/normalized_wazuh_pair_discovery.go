package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
)

// writeVerifiedSelectedWazuhPairLead finds a nearest-in-time authorized
// second-Case event for an exact first-Case Entry ID, even when the first event
// lies beyond the comparison receipt's ten-lead display limit.
// Returns metadata only. Source previews still require separate private
// inspection with both exact IDs and all Case receipts.
func writeVerifiedSelectedWazuhPairLead(ctx context.Context, principal,
	firstPath, firstCase, firstDigest, secondPath, secondCase, secondDigest,
	firstEntryID string, output io.Writer) error {
	if ctx == nil || ctx.Err() != nil || principal == "" || firstEntryID == "" ||
		firstCase == secondCase || output == nil {
		return errors.New("select two authorized Cases and an exact first Wazuh Entry ID")
	}
	first, _, err := loadAuthorizedImportedCase(ctx, principal, firstPath, firstCase, firstDigest)
	if err != nil { return err }
	second, _, err := loadAuthorizedImportedCase(ctx, principal, secondPath, secondCase, secondDigest)
	if err != nil { return err }
	verified, err := compareLoadedAuthorizedSavedCases(ctx, principal, first, second, firstCase, secondCase)
	if err != nil { return err }
	if verified.First.CaseID != firstCase || verified.Second.CaseID != secondCase {
		return errors.New("compared Case identity is not authorized")
	}
	left, err := collectAuthorizedWazuhObservations(first, firstCase)
	if err != nil { return err }
	right, err := collectAuthorizedWazuhObservations(second, secondCase)
	if err != nil { return err }
	var selected map[authorizedWazuhCohort][]authorizedWazuhObservation
	for cohort, events := range left {
		for _, event := range events {
			if event.id != firstEntryID { continue }
			if selected != nil { return errors.New("ambiguous Wazuh Entry ID") }
			selected = map[authorizedWazuhCohort][]authorizedWazuhObservation{
				cohort: {event},
			}
		}
	}
	if selected == nil { return errors.New("first Wazuh Entry ID not found in authorized Case") }
	matches := sharedWazuhPairsFromObservations(selected, right)
	if len(matches) != 1 {
		return errors.New("no authorized matching second-Case Wazuh event within 60 seconds")
	}
	if err := ctx.Err(); err != nil { return err }
	return json.NewEncoder(output).Encode(struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		Pair sharedAuthorizedWazuhEventPair `json:"pair"`
	}{
		Format: "cmdr-authorized-wazuh-pair-lead-v1",
		Synthetic: true, Pair: matches[0],
	})
}
