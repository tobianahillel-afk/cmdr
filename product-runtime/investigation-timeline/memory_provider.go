package investigationtimeline

import (
	"context"
	"errors"
	"net/http"
	"unicode/utf8"
)

// MemoryTimelineData is an immutable, process-local owner snapshot for a
// runnable Case -> Timeline HTTP integration. The host authenticates users;
// this module never treats an HTTP header as proof of identity.
type MemoryTimelineData struct {
	Principals []MemoryTimelinePrincipal
	Cases      []MemoryTimelineCase
	Entries    []MemoryTimelineEntry
	Sources    []MemoryTimelineSource
}

type MemoryTimelinePrincipal struct{ ID, Tenant, Environment string }
type MemoryTimelineCase struct {
	ID, Tenant, Environment string
	Readers, EntryIDs       []string
}
type MemoryTimelineEntry struct {
	CaseID  string
	Value   Entry
	Readers []string
}

// WazuhSourceMetadata is narrowly bounded descriptive information retained
// with an operator-approved source; it is NEVER an access-control decision.
// Original alert JSON, hostnames, IPs and full_log are not stored here.
type WazuhSourceMetadata struct {
	RuleID    string `json:"rule_id,omitempty"`
	AgentID   string `json:"agent_id,omitempty"`
	RuleLevel *int   `json:"rule_level,omitempty"`
	Feed      int    `json:"feed,omitempty"`
}

type MemoryTimelineSource struct {
	Identity SourceOwnerIdentity
	Readers  []string
	// Preview belongs to the synthetic source owner, never the Timeline Entry.
	// It is available only via a fresh, authenticated source-open operation.
	Preview string
	// This optional, source-owned metadata is not returned by source-open.
	// Only already-authorized investigation jobs may use it for filtering.
	Wazuh *WazuhSourceMetadata `json:"Wazuh,omitempty"`
}

type memoryTimelinePrincipalKey struct{}

// ACL is scoped to a canonical source object, not a caller-chosen reference.
// Only one reference per Case/kind/object is accepted in the memory provider.
type memorySourceObjectKey struct {
	caseID, kind, objectID string
}

type memoryTimelineProvider struct {
	authenticate  func(*http.Request) (string, error)
	principals    map[string]MemoryTimelinePrincipal
	cases         map[string]MemoryTimelineCase
	entries       map[string]MemoryTimelineEntry
	sources       map[string]MemoryTimelineSource
	sourceObjects map[memorySourceObjectKey]MemoryTimelineSource
}

// NewMemoryCaseTimelineHTTPHandler creates a fully wired read-only HTTP flow
// using a static, trusted owner dataset and a HOST-VERIFIED authentication
// callback. No persistence, remote ACL, mutation API or production auth
// provider is supplied. Do not use this adapter as a production data store.
// MemoryTimelineHTTPHandlers shares ONE immutable owner snapshot for Case
// discovery and Case -> Timeline reads, preventing a split demo authority.
type MemoryTimelineHTTPHandlers struct {
	Timeline   http.Handler
	Cases      http.Handler
	Navigate   http.Handler
	OpenSource http.Handler // synthetic owner preview only; not a production source backend
}

func NewMemoryCaseTimelineHTTPHandler(authenticate func(*http.Request) (string, error), data MemoryTimelineData) (http.Handler, error) {
	handlers, err := NewMemoryCaseTimelineHTTPHandlers(authenticate, data)
	if err != nil {
		return nil, err
	}
	return handlers.Timeline, nil
}

// NewMemoryCaseTimelineHTTPHandlers is a process-local integration reference,
// NOT a durable provider or a deployed authentication/authorization service.
func NewMemoryCaseTimelineHTTPHandlers(authenticate func(*http.Request) (string, error), data MemoryTimelineData) (MemoryTimelineHTTPHandlers, error) {
	if authenticate == nil {
		return MemoryTimelineHTTPHandlers{}, errors.New("verified session authenticator required")
	}
	s := &memoryTimelineProvider{
		authenticate:  authenticate,
		principals:    make(map[string]MemoryTimelinePrincipal, len(data.Principals)),
		cases:         make(map[string]MemoryTimelineCase, len(data.Cases)),
		entries:       make(map[string]MemoryTimelineEntry, len(data.Entries)),
		sources:       make(map[string]MemoryTimelineSource, len(data.Sources)),
		sourceObjects: make(map[memorySourceObjectKey]MemoryTimelineSource, len(data.Sources)),
	}
	for _, principal := range data.Principals {
		if invalidScope(principal.ID) || invalidScope(principal.Tenant) || invalidScope(principal.Environment) {
			return MemoryTimelineHTTPHandlers{}, errors.New("invalid principal")
		}
		if _, exists := s.principals[principal.ID]; exists {
			return MemoryTimelineHTTPHandlers{}, errors.New("duplicate principal")
		}
		s.principals[principal.ID] = principal
	}
	for _, c := range data.Cases {
		if invalidScope(c.ID) || invalidScope(c.Tenant) || invalidScope(c.Environment) ||
			len(c.EntryIDs) > MaxEntries || !s.validReaders(c.Readers, c.Tenant, c.Environment) {
			return MemoryTimelineHTTPHandlers{}, errors.New("invalid Case record")
		}
		if _, exists := s.cases[c.ID]; exists {
			return MemoryTimelineHTTPHandlers{}, errors.New("duplicate Case record")
		}
		c.Readers = append([]string(nil), c.Readers...)
		c.EntryIDs = append([]string(nil), c.EntryIDs...)
		s.cases[c.ID] = c
	}
	for _, e := range data.Entries {
		c, ok := s.cases[e.CaseID]
		if !ok || invalidScope(e.Value.ID) || e.Value.Tenant != c.Tenant ||
			!validTimelineEntryBeforeSourceLookup(e.Value) ||
			!s.validReaders(e.Readers, c.Tenant, c.Environment) {
			return MemoryTimelineHTTPHandlers{}, errors.New("invalid Timeline Entry record")
		}
		if _, exists := s.entries[e.Value.ID]; exists {
			return MemoryTimelineHTTPHandlers{}, errors.New("duplicate Timeline Entry record")
		}
		e.Readers = append([]string(nil), e.Readers...)
		// Caller-provided source-access claims can never become authority.
		e.Value.SourcePermission = "unknown"
		s.entries[e.Value.ID] = e
	}
	for _, src := range data.Sources {
		id := src.Identity
		c, ok := s.cases[id.CaseID]
		_, validKind := canonicalSourceReadPermission(id.Kind)
		if !ok || !validKind || invalidScope(id.Reference) || invalidScope(id.ObjectID) ||
			id.Tenant != c.Tenant || id.Environment != c.Environment ||
			len(src.Preview) > 2048 || !utf8.ValidString(src.Preview) ||
			!s.validReaders(src.Readers, c.Tenant, c.Environment) {
			return MemoryTimelineHTTPHandlers{}, errors.New("invalid source owner record")
		}
		if src.Wazuh != nil {
			validID := func(v string) bool {
				if len(v) > 16 {
					return false
				}
				for _, ch := range v {
					if ch < '0' || ch > '9' {
						return false
					}
				}
				return true
			}
			if id.Kind != "telemetry-event" || !validID(src.Wazuh.RuleID) ||
				!validID(src.Wazuh.AgentID) || src.Wazuh.Feed < 0 ||
				src.Wazuh.Feed > 4 ||
				(src.Wazuh.RuleLevel != nil &&
					(*src.Wazuh.RuleLevel < 0 || *src.Wazuh.RuleLevel > 16)) {
				return MemoryTimelineHTTPHandlers{}, errors.New("invalid bounded Wazuh owner metadata")
			}
		}
		if _, exists := s.sources[id.Reference]; exists {
			return MemoryTimelineHTTPHandlers{}, errors.New("duplicate source reference")
		}
		key := memorySourceObjectKey{caseID: id.CaseID, kind: id.Kind, objectID: id.ObjectID}
		if _, exists := s.sourceObjects[key]; exists {
			return MemoryTimelineHTTPHandlers{}, errors.New("ambiguous canonical source object references")
		}
		src.Readers = append([]string(nil), src.Readers...)
		// Caller-owned Wazuh metadata must not remain mutable after owner
		// validation. The severity pointer requires its own copy too:
		// otherwise callers can change a validated source after startup.
		if src.Wazuh != nil {
			frozen := *src.Wazuh
			if frozen.RuleLevel != nil {
				level := *frozen.RuleLevel
				frozen.RuleLevel = &level
			}
			src.Wazuh = &frozen
		}
		s.sources[id.Reference] = src
		s.sourceObjects[key] = src
	}
	for _, c := range s.cases {
		seen := make(map[string]bool, len(c.EntryIDs))
		for _, id := range c.EntryIDs {
			e, ok := s.entries[id]
			if !ok || seen[id] || e.CaseID != c.ID {
				return MemoryTimelineHTTPHandlers{}, errors.New("Case catalog points to foreign, missing or duplicate entry")
			}
			seen[id] = true
			src, ok := s.sources[e.Value.SourceRef]
			if !ok || src.Identity.CaseID != c.ID {
				return MemoryTimelineHTTPHandlers{}, errors.New("entry references missing or foreign source owner")
			}
		}
	}
	// The reference in-memory backend exercises the very same wiring used by
	// any host-supplied Case/Timeline/source owner providers.
	combined, err := NewCaseInvestigationHTTPHandlers(s, s, s, s, s, s)
	if err != nil {
		return MemoryTimelineHTTPHandlers{}, err
	}
	return MemoryTimelineHTTPHandlers{
		Timeline: combined.Timeline, Cases: combined.Cases,
		Navigate: combined.Navigate, OpenSource: s.sourcePreviewHandler(),
	}, nil
}

func (s *memoryTimelineProvider) validReaders(ids []string, tenant, environment string) bool {
	for _, id := range ids {
		p, ok := s.principals[id]
		if !ok || p.Tenant != tenant || p.Environment != environment {
			return false
		}
	}
	return true
}

func containsReader(readers []string, id string) bool {
	for _, reader := range readers {
		if reader == id {
			return true
		}
	}
	return false
}

func (s *memoryTimelineProvider) ResolveCaseTimelineSession(r *http.Request) (context.Context, Scope, error) {
	id, err := s.authenticate(r)
	if err != nil || r.Context().Err() != nil {
		return nil, Scope{}, errors.New("unauthenticated")
	}
	p, ok := s.principals[id]
	if !ok {
		return nil, Scope{}, errors.New("unauthenticated")
	}
	return context.WithValue(r.Context(), memoryTimelinePrincipalKey{}, id),
		Scope{Tenant: p.Tenant, Environment: p.Environment}, nil
}

func (s *memoryTimelineProvider) caseAuthorized(ctx context.Context, r AccessRequest) bool {
	if ctx.Err() != nil {
		return false
	}
	id, ok := ctx.Value(memoryTimelinePrincipalKey{}).(string)
	if !ok {
		return false
	}
	p, ok := s.principals[id]
	if !ok || p.Tenant != r.Tenant || p.Environment != r.Environment {
		return false
	}
	c, ok := s.cases[r.CaseID]
	return ok && c.Tenant == r.Tenant && c.Environment == r.Environment &&
		containsReader(c.Readers, id)
}

func (s *memoryTimelineProvider) CheckAccess(ctx context.Context, r AccessRequest) (string, error) {
	if !s.caseAuthorized(ctx, r) {
		return "deny", nil
	}
	principal := ctx.Value(memoryTimelinePrincipalKey{}).(string)
	switch r.Resource {
	case "case":
		if r.ObjectID == r.CaseID && r.Permission == caseReadPermission {
			return "allow", nil
		}
	case "timeline-entry":
		entry, ok := s.entries[r.ObjectID]
		if ok && entry.CaseID == r.CaseID && r.Permission == timelineReadPermission &&
			containsReader(entry.Readers, principal) {
			return "allow", nil
		}
	case "evidence", "telemetry-event", "decision", "incident":
		expected, valid := canonicalSourceReadPermission(r.Resource)
		if valid && r.Permission == expected {
			key := memorySourceObjectKey{caseID: r.CaseID, kind: r.Resource, objectID: r.ObjectID}
			source, exists := s.sourceObjects[key]
			if exists && containsReader(source.Readers, principal) {
				return "allow", nil
			}
		}
	}
	return "deny", nil
}

func (s *memoryTimelineProvider) ListCaseEntryIDs(ctx context.Context, r AccessRequest) ([]string, error) {
	allow, _ := s.CheckAccess(ctx, r)
	if allow != "allow" || r.Resource != "case" || r.ObjectID != r.CaseID || r.Permission != caseReadPermission {
		return nil, errors.New("Case unavailable")
	}
	return append([]string(nil), s.cases[r.CaseID].EntryIDs...), nil
}

func (s *memoryTimelineProvider) ResolveEntry(ctx context.Context, r AccessRequest) (Entry, error) {
	allow, _ := s.CheckAccess(ctx, r)
	entry, exists := s.entries[r.ObjectID]
	if allow != "allow" || !exists || r.Resource != "timeline-entry" {
		return Entry{}, errors.New("Timeline Entry unavailable")
	}
	return entry.Value, nil
}

func (s *memoryTimelineProvider) ResolveSource(ctx context.Context, r AccessRequest) (SourceOwnerIdentity, error) {
	if !s.caseAuthorized(ctx, r) || r.Resource != "source-reference" {
		return SourceOwnerIdentity{}, errors.New("source unavailable")
	}
	source, exists := s.sources[r.ObjectID]
	if !exists || source.Identity.CaseID != r.CaseID {
		return SourceOwnerIdentity{}, errors.New("source unavailable")
	}
	return source.Identity, nil
}

// ListCandidateCaseIDs is a Case-owner listing, not a permission oracle.
// Discovery invokes CheckAccess on every ID before presenting anything.
func (s *memoryTimelineProvider) ListCandidateCaseIDs(ctx context.Context, scope Scope) ([]string, error) {
	principal, ok := ctx.Value(memoryTimelinePrincipalKey{}).(string)
	p, present := s.principals[principal]
	if !ok || !present || ctx.Err() != nil || p.Tenant != scope.Tenant || p.Environment != scope.Environment {
		return nil, errors.New("Case catalogue unavailable")
	}
	ids := make([]string, 0, len(s.cases))
	for id, c := range s.cases {
		if c.Tenant == scope.Tenant && c.Environment == scope.Environment {
			ids = append(ids, id)
		}
	}
	return ids, nil
}
