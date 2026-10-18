package main

import (
	"errors"
	"sort"
	"time"

	search "github.com/tobianahillel-afk/cmdr/product-runtime/event-search"
	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// Enrich ONLY the references of an already completed, source-authorized
// first-party Event Search Job. This is not a second search or source ACL:
// each matched source was re-opened and authorized by the backend before
// its stable reference was included in the completed Search Job.
func authorizedPersistedWazuhHits(data timeline.MemoryTimelineData, caseID string,
	refs []search.ResultReference, metadata map[string]timeline.WazuhSourceMetadata) ([]wazuhSearchHit, error) {
	if len(refs) > maxWazuhSearchMatches || len(metadata) == 0 {
		return nil, errors.New("persisted Wazuh search metadata unavailable")
	}
	entries := make(map[string]timeline.MemoryTimelineEntry, len(data.Entries))
	for _, item := range data.Entries {
		if item.CaseID == caseID {
			entries[item.Value.ID] = item
		}
	}
	hits := make([]wazuhSearchHit, 0, len(refs))
	seen := make(map[string]bool, len(refs))
	for _, ref := range refs {
		if ref.Source != localJSONLSearchSource || seen[ref.Ref] {
			return nil, errors.New("completed persisted Wazuh search has inconsistent result references")
		}
		seen[ref.Ref] = true
		entry, exists := entries[ref.Ref]
		meta, authorized := metadata[ref.Ref]
		if !exists || !authorized || entry.Value.Kind != "observed" || entry.Value.ClockUncertaintyMS != 0 {
			return nil, errors.New("persisted Wazuh result lacks canonical authorized source and time")
		}
		at, err := time.Parse(time.RFC3339Nano, entry.Value.SourceTime)
		if err != nil {
			return nil, errors.New("persisted Wazuh result timestamp invalid")
		}
		hit := wazuhSearchHit{
			EntryID: ref.Ref,
			SourceTimeUTC: at.UTC().Format(time.RFC3339Nano),
			RuleID: meta.RuleID, AgentID: meta.AgentID, Feed: meta.Feed,
			RuleLevel: meta.RuleLevel,
		}
		hits = append(hits, hit)
	}
	sort.Slice(hits, func(i, j int) bool {
		if hits[i].SourceTimeUTC != hits[j].SourceTimeUTC {
			left, _ := time.Parse(time.RFC3339Nano, hits[i].SourceTimeUTC)
			right, _ := time.Parse(time.RFC3339Nano, hits[j].SourceTimeUTC)
			return left.Before(right)
		}
		return hits[i].EntryID < hits[j].EntryID
	})
	return hits, nil
}
