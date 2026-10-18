package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"sort"
)

// Unlike a first-Entry pivot, a SECOND-side Wazuh event can be the nearest
// peer of multiple independently reported first-Case events. Keep every
// eligible pair, sorted deterministically and paged, without reinterpreting
// correlation as event identity or authorization.
func verifiedWazuhPairsForSecondEntry(ctx context.Context,
	left, right map[authorizedWazuhCohort][]authorizedWazuhObservation,
	secondID string, offset int) ([]sharedAuthorizedWazuhEventPair, bool, error) {
	if ctx == nil || ctx.Err() != nil || secondID == "" ||
		offset < 0 || offset > 5000 {
		return nil, false, errors.New("an exact second Entry ID and bounded offset are required")
	}
	var target authorizedWazuhCohort
	found := false
	for cohort, events := range right {
		for _, event := range events {
			if event.id != secondID { continue }
			if found { return nil, false, errors.New("ambiguous second Wazuh Entry ID") }
			target, found = cohort, true
		}
	}
	if !found { return nil, false, errors.New("second Wazuh Entry ID not found in authorized Case") }

	first := append([]authorizedWazuhObservation(nil), left[target]...)
	sort.Slice(first, func(i, j int) bool {
		if !first[i].at.Equal(first[j].at) {
			return first[i].at.Before(first[j].at)
		}
		return first[i].id < first[j].id
	})
	pairs := make([]sharedAuthorizedWazuhEventPair, 0, 11)
	eligible := 0
	for _, event := range first {
		if err := ctx.Err(); err != nil { return nil, false, err }
		// Reuse the *full* authorized second cohort, including its competing
		// events. Selecting the sole peer before nearest matching would forge
		// correlations that were not present in the normal Case comparison.
		one := map[authorizedWazuhCohort][]authorizedWazuhObservation{
			target: {event},
		}
		matches := sharedWazuhPairsFromObservations(one, right)
		if len(matches) != 1 || matches[0].SecondEntryID != secondID { continue }
		eligible++
		if eligible <= offset { continue }
		pairs = append(pairs, matches[0])
		if len(pairs) == 11 { return pairs[:10], true, nil }
	}
	return pairs, false, nil
}

// Read-only reverse navigation for two independently SHA-256-verified,
// principal/tenant/Case/typed-source authorized snapshots. No private
// previews on stdout; only Entry IDs and correlation metadata are emitted.
func listVerifiedSavedWazuhPairsForSecondEntry(ctx context.Context, principal,
	firstPath, firstCase, firstSHA, secondPath, secondCase, secondSHA,
	secondID string, offset int, output io.Writer) error {
	if ctx == nil || ctx.Err() != nil || output == nil ||
		principal == "" || secondID == "" ||
		firstCase == "" || secondCase == "" || firstCase == secondCase ||
		firstSHA == "" || secondSHA == "" || offset < 0 || offset > 5000 {
		return errors.New("two distinct verified Cases, an exact second Entry ID and a bounded offset are required")
	}
	first, _, err := loadAuthorizedImportedCase(ctx, principal, firstPath, firstCase, firstSHA)
	if err != nil { return err }
	second, _, err := loadAuthorizedImportedCase(ctx, principal, secondPath, secondCase, secondSHA)
	if err != nil { return err }
	if _, err := compareLoadedAuthorizedSavedCases(ctx, principal,
		first, second, firstCase, secondCase); err != nil { return err }
	left, err := collectAuthorizedWazuhObservations(first, firstCase)
	if err != nil { return err }
	right, err := collectAuthorizedWazuhObservations(second, secondCase)
	if err != nil { return err }
	pairs, more, err := verifiedWazuhPairsForSecondEntry(ctx, left, right, secondID, offset)
	if err != nil { return err }
	if err := ctx.Err(); err != nil { return err }
	var encoded bytes.Buffer
	if err := json.NewEncoder(&encoded).Encode(struct {
		Format string `json:"format"`
		Synthetic bool `json:"synthetic"`
		FirstCase string `json:"first_case_id"`
		SecondCase string `json:"second_case_id"`
		SecondEntryID string `json:"second_entry_id"`
		Offset int `json:"offset"`
		HasMore bool `json:"has_more"`
		Pairs []sharedAuthorizedWazuhEventPair `json:"pairs"`
	}{
		Format: "cmdr-authorized-wazuh-second-entry-pairs-v1",
		Synthetic: true, FirstCase: firstCase, SecondCase: secondCase,
		SecondEntryID: secondID, Offset: offset, HasMore: more, Pairs: pairs,
	}); err != nil { return err }
	if err := ctx.Err(); err != nil { return err }
	_, err = io.Copy(output, &encoded)
	return err
}
