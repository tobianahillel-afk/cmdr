package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"sort"
	"time"
)

// Cohort counts are navigation leads only. Both Case source owners must be
// authorized and SHA-256-verified before any cohort metadata is disclosed.
type verifiedWazuhCohortCount struct {
	RuleID       string `json:"rule_id"`
	AgentID      string `json:"agent_id"`
	FirstEvents  int    `json:"first_events"`
	SecondEvents int    `json:"second_events"`
	MatchedPairs int    `json:"matched_pairs"`
	// First source-authorized pair, for a direct exact Entry-ID navigation.
	// No report contents or access authority are included.
	LeadPair *sharedAuthorizedWazuhEventPair `json:"lead_pair,omitempty"`
}

// Count precisely the first-Case events with an eligible second-Case peer
// within the same 60-second window as the existing pair-pagination engine.
// The nearest-peer tie-break determines WHICH ID is selected but not WHETHER
// the first observation has a match. Input cohorts are already time-sorted.
func countAuthorizedCohortPairs(first, second []authorizedWazuhObservation) int {
	const radius = 60 * time.Second
	count := 0
	for _, event := range first {
		at := event.at
		idx := sort.Search(len(second), func(i int) bool {
			return !second[i].at.Before(at)
		})
		for _, candidate := range []int{idx - 1, idx} {
			if candidate < 0 || candidate >= len(second) {
				continue
			}
			delta := second[candidate].at.Sub(at)
			if delta >= -radius && delta <= radius {
				count++
				break
			}
		}
	}
	return count
}

// Return only matching cohorts, ranked by useful correlated evidence volume.
// Never include unpaired rule/agent metadata from either Case.
func authorizedWazuhCohortCatalog(
	left, right map[authorizedWazuhCohort][]authorizedWazuhObservation,
	offset int, filters ...string,
) ([]verifiedWazuhCohortCount, bool, error) {
	if len(filters) != 0 && len(filters) != 2 {
		return nil, false, errors.New("cohort catalog requires zero or two exact filters")
	}
	var ruleID, agentID string
	if len(filters) == 2 { ruleID, agentID = filters[0], filters[1] }
	if offset < 0 || offset > 5000 {
		return nil, false, errors.New("Wazuh cohort catalog offset out of bounds")
	}
	found := make([]verifiedWazuhCohortCount, 0)
	for key, first := range left {
		if (ruleID != "" && key.rule != ruleID) || (agentID != "" && key.agent != agentID) {
			continue
		}
		second := right[key]
		if len(first) == 0 || len(second) == 0 {
			continue
		}
		count := countAuthorizedCohortPairs(first, second)
		if count == 0 {
			continue
		}
		found = append(found, verifiedWazuhCohortCount{
			RuleID: key.rule, AgentID: key.agent,
			FirstEvents: len(first), SecondEvents: len(second), MatchedPairs: count,
		})
	}
	sort.Slice(found, func(i, j int) bool {
		a, b := found[i], found[j]
		if a.MatchedPairs != b.MatchedPairs {
			return a.MatchedPairs > b.MatchedPairs
		}
		if a.RuleID != b.RuleID {
			return a.RuleID < b.RuleID
		}
		return a.AgentID < b.AgentID
	})
	if offset >= len(found) {
		return []verifiedWazuhCohortCount{}, false, nil
	}
	end := offset + 10
	if end > len(found) {
		end = len(found)
	}
	// Resolve only the page's (at most ten) authorized cohorts. Reuse the
	// exact pair navigator so IDs and time tie-breaks cannot diverge.
	page := found[offset:end]
	for i := range page {
		key := authorizedWazuhCohort{rule: page[i].RuleID, agent: page[i].AgentID}
		first, _ := sharedWazuhPairsPageFromObservations(
			map[authorizedWazuhCohort][]authorizedWazuhObservation{key: left[key]},
			map[authorizedWazuhCohort][]authorizedWazuhObservation{key: right[key]}, 0)
		if len(first) == 0 {
			return nil, false, errors.New("authorized Wazuh cohort has no navigable pair")
		}
		page[i].LeadPair = &first[0]
	}
	return page, end < len(found), nil
}

func listVerifiedSavedWazuhCohortCatalog(ctx context.Context, principal,
	firstPath, firstCase, firstSHA, secondPath, secondCase, secondSHA string,
	offset int, filterRuleID, filterAgentID string, output io.Writer,
) error {
	if ctx == nil || ctx.Err() != nil || output == nil || principal == "" ||
		firstCase == "" || secondCase == "" || firstCase == secondCase ||
		firstSHA == "" || secondSHA == "" || offset < 0 || offset > 5000 {
		return errors.New("two distinct verified Cases and a bounded cohort offset are required")
	}
	first, _, err := loadAuthorizedImportedCase(ctx, principal, firstPath, firstCase, firstSHA)
	if err != nil {
		return err
	}
	second, _, err := loadAuthorizedImportedCase(ctx, principal, secondPath, secondCase, secondSHA)
	if err != nil {
		return err
	}
	if _, err := compareLoadedAuthorizedSavedCases(ctx, principal, first, second, firstCase, secondCase); err != nil {
		return err
	}
	left, err := collectAuthorizedWazuhObservations(first, firstCase)
	if err != nil {
		return err
	}
	right, err := collectAuthorizedWazuhObservations(second, secondCase)
	if err != nil {
		return err
	}
	cohorts, hasMore, err := authorizedWazuhCohortCatalog(left, right, offset, filterRuleID, filterAgentID)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	var encoded bytes.Buffer
	err = json.NewEncoder(&encoded).Encode(struct {
		Format       string                     `json:"format"`
		Synthetic    bool                       `json:"synthetic"`
		FirstCaseID  string                     `json:"first_case_id"`
		SecondCaseID string                     `json:"second_case_id"`
		Offset       int                        `json:"offset"`
		FilterRuleID string                     `json:"filter_rule_id,omitempty"`
		FilterAgentID string                    `json:"filter_agent_id,omitempty"`
		HasMore      bool                       `json:"has_more"`
		Cohorts      []verifiedWazuhCohortCount `json:"cohorts"`
	}{
		Format: "cmdr-authorized-wazuh-cohort-catalog-v1",
		Synthetic: true, FirstCaseID: firstCase, SecondCaseID: secondCase,
		Offset: offset, FilterRuleID: filterRuleID, FilterAgentID: filterAgentID,
		HasMore: hasMore, Cohorts: cohorts,
	})
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	_, err = io.Copy(output, &encoded)
	return err
}
