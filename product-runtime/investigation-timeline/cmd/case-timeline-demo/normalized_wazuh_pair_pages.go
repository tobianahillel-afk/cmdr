package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
)

// Discover further source-linked 60-second Wazuh pair leads from two
// INDEPENDENTLY verified Case owners. This reuses the exact first-party
// authorization and correlation rules of the ten-item comparison receipt.
// Returned Entry IDs are navigation selectors, never access grants.
func listVerifiedSavedWazuhPairPage(ctx context.Context, principal,
	firstPath, firstCase, firstSHA, secondPath, secondCase, secondSHA string,
	offset int, output io.Writer) error {
	if ctx == nil || ctx.Err() != nil || offset < 0 || offset > 5000 ||
		firstCase == "" || secondCase == "" || firstCase == secondCase ||
		principal == "" || firstSHA == "" || secondSHA == "" {
		return errors.New("two distinct authorized Cases and a bounded Wazuh pair offset are required")
	}
	first, _, err := loadAuthorizedImportedCase(ctx, principal, firstPath, firstCase, firstSHA)
	if err != nil { return err }
	second, _, err := loadAuthorizedImportedCase(ctx, principal, secondPath, secondCase, secondSHA)
	if err != nil { return err }
	// Both Case/Timeline/typed-source owners must pass their entire existing
	// authenticated pagination before any source-linked correlation is visible.
	if _, err := compareLoadedAuthorizedSavedCases(ctx, principal, first, second, firstCase, secondCase); err != nil {
		return err
	}
	left, err := collectAuthorizedWazuhObservations(first, firstCase)
	if err != nil { return err }
	right, err := collectAuthorizedWazuhObservations(second, secondCase)
	if err != nil { return err }
	pairs, hasMore := sharedWazuhPairsPageFromObservations(left, right, offset)
	if err := ctx.Err(); err != nil { return err }
	var encoded bytes.Buffer
	err = json.NewEncoder(&encoded).Encode(struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		FirstCase string `json:"first_case_id"`
		SecondCase string `json:"second_case_id"`
		Offset int `json:"offset"`
		HasMore bool `json:"has_more"`
		Pairs []sharedAuthorizedWazuhEventPair `json:"pairs"`
	}{
		Format: "cmdr-local-authorized-wazuh-pair-page-v1",
		Synthetic: true, FirstCase: firstCase, SecondCase: secondCase,
		Offset: offset, HasMore: hasMore, Pairs: pairs,
	})
	if err != nil { return err }
	if err := ctx.Err(); err != nil { return err }
	_, err = io.Copy(output, &encoded)
	return err
}
