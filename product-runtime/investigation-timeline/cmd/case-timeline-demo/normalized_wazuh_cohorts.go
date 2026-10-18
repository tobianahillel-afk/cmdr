package main

import (
	"errors"
	"sort"
	"time"

	timeline "github.com/tobianahillel-afk/cmdr/product-runtime/investigation-timeline"
)

// A cohort is scoped to ONE already authorized Case owner. It is not a
// global alert identity and cannot grant read access to another Case.
type authorizedWazuhCohort struct { rule, agent string }
type authorizedWazuhObservation struct { id string; at time.Time }

// One bounded, canonical source-time traversal powers both the existing
// five-minute counts and the precise 60-second event-linked leads.
// All source authorization happens before invoking this function in the
// first-party dual-Case comparison; neither data owner is joined or copied.
func collectAuthorizedWazuhObservations(data timeline.MemoryTimelineData, caseID string) (map[authorizedWazuhCohort][]authorizedWazuhObservation, error) {
	imported, err := importedNormalizedIDs(data, caseID)
	if err != nil { return nil, err }
	sources := make(map[string]timeline.MemoryTimelineSource, len(data.Sources))
	for _, source := range data.Sources {
		if source.Identity.CaseID == caseID && source.Identity.Kind == "telemetry-event" {
			sources[source.Identity.Reference] = source
		}
	}
	groups := make(map[authorizedWazuhCohort][]authorizedWazuhObservation)
	for _, entry := range data.Entries {
		if entry.CaseID != caseID { continue }
		if _, authorized := imported[entry.Value.ID]; !authorized { continue }
		source, found := sources[entry.Value.SourceRef]
		if !found || source.Wazuh == nil { continue }
		if source.Wazuh.RuleID == "" || source.Wazuh.AgentID == "" { continue }
		if entry.Value.Kind != "observed" || entry.Value.ClockUncertaintyMS != 0 {
			return nil, errors.New("Wazuh Case entry must have an observed certain timestamp")
		}
		at, err := time.Parse(time.RFC3339Nano, entry.Value.SourceTime)
		if err != nil { return nil, errors.New("Wazuh Case source timestamp invalid") }
		key := authorizedWazuhCohort{rule:source.Wazuh.RuleID,agent:source.Wazuh.AgentID}
		groups[key] = append(groups[key],authorizedWazuhObservation{id:entry.Value.ID,at:at.UTC()})
	}
	for key := range groups {
		sort.Slice(groups[key],func(i,j int)bool{
			a,b:=groups[key][i],groups[key][j]
			if a.at.Equal(b.at) { return a.id<b.id }
			return a.at.Before(b.at)
		})
	}
	return groups,nil
}
