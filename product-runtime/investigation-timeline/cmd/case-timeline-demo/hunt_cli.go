package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// An explicit investigator-driven text hunt on an operator-owned local
// snapshot. Text is a search hint, NOT an IOC verdict or source identity.
const maxHuntCases = 8

type localHuntCase struct {
	CaseID          string            `json:"case_id"`
	Matches         int               `json:"matches"`
	Observed        int               `json:"observed"`
	Inferred        int               `json:"inferred"`
	Uncertain       int               `json:"uncertain"`
	FirstCertainUTC string            `json:"first_certain_utc,omitempty"`
	LastCertainUTC  string            `json:"last_certain_utc,omitempty"`
	TopWindows      []localHuntWindow `json:"top_five_minute_windows,omitempty"`
}

// These UTC buckets contain only entries that passed Case, Entry and source
// permission checks. Their counts are not proof that an IOC is malicious.
type localHuntWindow struct {
	StartUTC string `json:"start_utc"`
	Matches  int    `json:"matches"`
	Observed int    `json:"observed"`
	Inferred int    `json:"inferred"`
}

// Return at most three actionable windows per Case without disclosing raw
// entry IDs, source references or unbounded time series.
func topAuthorizedHuntWindows(windows map[string]certainWindowCount) []localHuntWindow {
	out := make([]localHuntWindow, 0, len(windows))
	for start, count := range windows {
		out = append(out, localHuntWindow{
			StartUTC: start, Matches: count.Total,
			Observed: count.Observed, Inferred: count.Inferred,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Matches != out[j].Matches {
			return out[i].Matches > out[j].Matches
		}
		return out[i].StartUTC < out[j].StartUTC
	})
	if len(out) > 3 {
		out = out[:3]
	}
	return out
}

type localCrossCaseHunt struct {
	Format         string          `json:"format"`
	Synthetic      bool            `json:"synthetic"`
	Query          string          `json:"query"`
	SearchedCases  int             `json:"searched_cases"`
	MatchingCases  int             `json:"matching_cases"`
	VisibleMatches int             `json:"visible_matches"`
	Results        []localHuntCase `json:"results"`
}

func validHuntText(q string) bool {
	return len(q) >= 2 && len(q) <= 96 && strings.TrimSpace(q) == q &&
		utf8.ValidString(q) && strings.IndexFunc(q, unicode.IsControl) < 0
}

// Search authorized cases only; every search page goes back through the
// Case/Timeline Entry/typed-source owner permission checks on the server.
// Never search raw snapshots or infer a shared source identity across Cases.
func collectAuthorizedCrossCaseHunt(
	ctx context.Context, handler http.Handler, token, query string,
) (localCrossCaseHunt, error) {
	if err := ctx.Err(); err != nil {
		return localCrossCaseHunt{}, err
	}
	if !validHuntText(query) {
		return localCrossCaseHunt{}, errors.New("hunt requires a printable 2-96 byte exact text query")
	}
	var discovery struct {
		Cases []string `json:"cases"`
	}
	if err := queryDemoOwner(ctx, handler, token, "/demo-cases", &discovery); err != nil {
		return localCrossCaseHunt{}, err
	}
	if len(discovery.Cases) == 0 || len(discovery.Cases) > maxHuntCases {
		return localCrossCaseHunt{}, errors.New("no bounded authorized Case catalog available")
	}
	result := localCrossCaseHunt{
		Format: "cmdr-local-cross-case-text-hunt-v1", Synthetic: true,
		Query: query, SearchedCases: len(discovery.Cases),
		Results: make([]localHuntCase, 0),
	}
	seen := make(map[string]bool, len(discovery.Cases))
	for _, caseID := range discovery.Cases {
		if caseID == "" || len(caseID) > 1024 || seen[caseID] {
			return localCrossCaseHunt{}, errors.New("invalid authorized Case catalog")
		}
		seen[caseID] = true
		summary, err := collectDemoCaseTriageWithObserverAndQuery(
			ctx, handler, token, caseID, query, nil,
		)
		if err != nil {
			return localCrossCaseHunt{}, err
		}
		if summary.VisibleEntries == 0 {
			continue
		}
		result.MatchingCases++
		result.VisibleMatches += summary.VisibleEntries
		result.Results = append(result.Results, localHuntCase{
			CaseID: caseID, Matches: summary.VisibleEntries,
			Observed: summary.Kinds["observed"],
			Inferred: summary.Kinds["inferred"],
			Uncertain: summary.UncertainEntries,
			FirstCertainUTC: summary.FirstCertainUTC,
			LastCertainUTC: summary.LastCertainUTC,
			TopWindows: topAuthorizedHuntWindows(summary.certainWindows),
		})
	}
	if err := ctx.Err(); err != nil {
		return localCrossCaseHunt{}, err
	}
	return result, nil
}

func writeDemoCrossCaseHunt(
	ctx context.Context, principal, snapshot, query string, output io.Writer,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if principal == "" || !validHuntText(query) {
		return errors.New("explicit synthetic principal and bounded text query required")
	}
	data, err := loadDemoSnapshot(snapshot)
	if err != nil {
		return err
	}
	token, err := randomBearer()
	if err != nil {
		return err
	}
	handler, err := newDemoHandlerWithData(token, principal, data)
	if err != nil {
		return err
	}
	result, err := collectAuthorizedCrossCaseHunt(ctx, handler, token, query)
	if err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	// One machine-readable receipt only after ALL selected Case pages succeed.
	return json.NewEncoder(output).Encode(result)
}
