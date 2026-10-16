package investigationtimeline

import (
	"fmt"
	"net/http"
)

// CaseInvestigationHTTPHandlers pairs Case discovery and authorized Timeline
// reads for a host that supplies its own authenticated owner services.
// The host retains URL routing, identity, persistence, audit, rate limits
// and deployment ownership; neither handler accepts caller ACL flags.
type CaseInvestigationHTTPHandlers struct {
	Cases    http.Handler
	Timeline http.Handler
	Navigate http.Handler
	Compare http.Handler
}

// NewCaseInvestigationHTTPHandlers composes existing providers rather than
// replacing or silently upgrading them. Both routes use the SAME session
// resolver and Case authority. Discovery checks Case read on every candidate;
// Timeline independently rechecks Case, Timeline Entry and linked-source ACLs.
func NewCaseInvestigationHTTPHandlers(
	sessions CaseTimelineSessionResolver, checker AccessChecker,
	candidateCases CaseListCatalog, timelineCases CaseTimelineCatalog,
	entries EntryResolver, sources SourceOwnerResolver,
) (CaseInvestigationHTTPHandlers, error) {
	if missingAuthority(sessions) || missingAuthority(checker) ||
		missingAuthority(candidateCases) || missingAuthority(timelineCases) ||
		missingAuthority(entries) || missingAuthority(sources) {
		return CaseInvestigationHTTPHandlers{}, fmt.Errorf("Case investigation requires authenticated Case and Timeline owner services")
	}
	cases, err := NewCaseDiscoveryHTTPHandler(sessions, checker, candidateCases)
	if err != nil {
		return CaseInvestigationHTTPHandlers{}, err
	}
	timeline, err := NewCaseTimelineHTTPHandler(sessions, checker, timelineCases, entries, sources)
	if err != nil {
		return CaseInvestigationHTTPHandlers{}, err
	}
	navigation, err := NewCaseSourceNavigationHTTPHandler(sessions, checker, timelineCases, entries, sources)
	if err != nil {
		return CaseInvestigationHTTPHandlers{}, err
	}
	comparison, err := NewCaseVersionComparisonHTTPHandler(sessions, checker, timelineCases, entries, sources)
	if err != nil {
		return CaseInvestigationHTTPHandlers{}, err
	}
	return CaseInvestigationHTTPHandlers{Cases: cases, Timeline: timeline, Navigate: navigation, Compare: comparison}, nil
}
