package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	workLeasePolicyPath   = "engineering/recovery/lease-policy.json"
	workLeaseRegistryPath = "engineering/recovery/claims.json"
)

type WorkLeasePolicy struct {
	SchemaVersion            int      `json:"schema_version"`
	RegistryKind             string   `json:"registry_kind"`
	DefaultPolicy            string   `json:"default_policy"`
	DefaultLeaseMinutes      int      `json:"default_lease_minutes"`
	MinLeaseMinutes          int      `json:"min_lease_minutes"`
	MaxLeaseMinutes          int      `json:"max_lease_minutes"`
	MaxActiveClaimsPerAgent  int      `json:"max_active_claims_per_agent"`
	AcquireStatuses          []string `json:"acquire_statuses"`
}

type WorkLease struct {
	ID          string `json:"id"`
	WorkUnit    string `json:"work_unit"`
	AgentID     string `json:"agent_id"`
	BaseHeadSHA string `json:"base_head_sha"`
	AcquiredAt  string `json:"acquired_at"`
	ExpiresAt   string `json:"expires_at"`
	ReleasedAt  string `json:"released_at,omitempty"`
}

type WorkLeaseRegistry struct {
	SchemaVersion int         `json:"schema_version"`
	RegistryKind  string      `json:"registry_kind"`
	Claims        []WorkLease `json:"claims"`
}

type WorkLeaseStatus struct {
	ID       string `json:"id"`
	WorkUnit string `json:"work_unit"`
	AgentID  string `json:"agent_id"`
	Status   string `json:"status"`
}

type WorkLeaseAuditSummary struct {
	AsOf     string            `json:"as_of"`
	Claims   int               `json:"claims"`
	Active   int               `json:"active"`
	Expired  int               `json:"expired"`
	Released int               `json:"released"`
	Entries  []WorkLeaseStatus `json:"entries,omitempty"`
}

type WorkLeaseActionRequest struct {
	Action          string
	LeaseID         string
	AgentID         string
	WorkUnit        string
	BaseHeadSHA     string
	AsOf            string
	DurationMinutes int
}

type WorkLeaseEvaluation struct {
	Action        string     `json:"action"`
	Allowed       bool       `json:"allowed"`
	Reasons       []string   `json:"reasons,omitempty"`
	ProposedClaim *WorkLease `json:"proposed_claim,omitempty"`
}

var (
	workLeaseIDPattern = regexp.MustCompile(`^LEASE-[A-Z0-9][A-Z0-9._-]{7,63}$`)
	agentIDPattern     = regexp.MustCompile(`^AGENT-[A-Za-z0-9][A-Za-z0-9._-]{2,95}$`)
)

func loadWorkLeaseConfiguration(root string) (WorkLeasePolicy, WorkLeaseRegistry, error) {
	var policy WorkLeasePolicy
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(workLeasePolicyPath)), &policy); err != nil {
		return policy, WorkLeaseRegistry{}, err
	}
	if err := validateWorkLeasePolicy(policy); err != nil {
		return policy, WorkLeaseRegistry{}, err
	}
	var registry WorkLeaseRegistry
	if err := decodeStrict(root, filepath.Join(root, filepath.FromSlash(workLeaseRegistryPath)), &registry); err != nil {
		return policy, registry, err
	}
	if registry.SchemaVersion != 1 || registry.RegistryKind != "engineering-work-leases" {
		return policy, registry, fmt.Errorf("invalid work lease registry header")
	}
	return policy, registry, nil
}

func validateWorkLeasePolicy(policy WorkLeasePolicy) error {
	if policy.SchemaVersion != 1 || policy.RegistryKind != "work-lease-policy" ||
		policy.DefaultPolicy != "deny-conflicting-active-lease" {
		return fmt.Errorf("invalid work lease policy header")
	}
	if policy.MinLeaseMinutes < 1 || policy.DefaultLeaseMinutes < policy.MinLeaseMinutes ||
		policy.DefaultLeaseMinutes > policy.MaxLeaseMinutes || policy.MaxLeaseMinutes > 24*60 {
		return fmt.Errorf("invalid lease duration policy")
	}
	if policy.MaxActiveClaimsPerAgent != 1 {
		return fmt.Errorf("max_active_claims_per_agent must be 1")
	}
	if len(policy.AcquireStatuses) != 1 || policy.AcquireStatuses[0] != "READY" {
		return fmt.Errorf("acquire_statuses must contain only READY")
	}
	return nil
}

func runWorkLeaseAudit(root, asOfValue string, graph WorkGraph) (WorkLeaseAuditSummary, error) {
	policy, registry, err := loadWorkLeaseConfiguration(root)
	if err != nil {
		return WorkLeaseAuditSummary{}, err
	}
	asOf, err := parseLeaseTimeOrNow(asOfValue)
	if err != nil {
		return WorkLeaseAuditSummary{}, err
	}
	return auditWorkLeases(policy, registry, graph, asOf)
}

func auditWorkLeases(policy WorkLeasePolicy, registry WorkLeaseRegistry, graph WorkGraph, asOf time.Time) (WorkLeaseAuditSummary, error) {
	summary := WorkLeaseAuditSummary{AsOf: asOf.UTC().Format(time.RFC3339), Claims: len(registry.Claims)}
	nodes := map[string]WorkNode{}
	for _, node := range graph.Nodes {
		nodes[node.ID] = node
	}
	ids := map[string]bool{}
	activeWork := map[string]string{}
	activeAgent := map[string]int{}
	for _, claim := range registry.Claims {
		acquired, expires, released, err := validateWorkLeaseRecord(policy, claim, nodes, asOf)
		if err != nil {
			return summary, err
		}
		if ids[claim.ID] {
			return summary, fmt.Errorf("duplicate work lease id %s", claim.ID)
		}
		ids[claim.ID] = true
		status := derivedLeaseStatus(acquired, expires, released, asOf)
		summary.Entries = append(summary.Entries, WorkLeaseStatus{ID: claim.ID, WorkUnit: claim.WorkUnit, AgentID: claim.AgentID, Status: status})
		switch status {
		case "active":
			node := nodes[claim.WorkUnit]
			if leaseTerminalStatus(node.Status) {
				return summary, fmt.Errorf("active lease %s references terminal work unit %s (%s)", claim.ID, claim.WorkUnit, node.Status)
			}
			if owner := activeWork[claim.WorkUnit]; owner != "" {
				return summary, fmt.Errorf("work unit %s has conflicting active leases %s and %s", claim.WorkUnit, owner, claim.ID)
			}
			activeWork[claim.WorkUnit] = claim.ID
			activeAgent[claim.AgentID]++
			if activeAgent[claim.AgentID] > policy.MaxActiveClaimsPerAgent {
				return summary, fmt.Errorf("agent %s exceeds max active claims", claim.AgentID)
			}
			summary.Active++
		case "expired":
			summary.Expired++
		case "released":
			summary.Released++
		}
	}
	sort.Slice(summary.Entries, func(i, j int) bool { return summary.Entries[i].ID < summary.Entries[j].ID })
	return summary, nil
}

func validateWorkLeaseRecord(policy WorkLeasePolicy, claim WorkLease, nodes map[string]WorkNode, asOf time.Time) (time.Time, time.Time, *time.Time, error) {
	if !workLeaseIDPattern.MatchString(claim.ID) || !agentIDPattern.MatchString(claim.AgentID) {
		return time.Time{}, time.Time{}, nil, fmt.Errorf("lease %q has invalid lease or agent id", claim.ID)
	}
	if _, ok := nodes[claim.WorkUnit]; !ok {
		return time.Time{}, time.Time{}, nil, fmt.Errorf("lease %s references unknown work unit %s", claim.ID, claim.WorkUnit)
	}
	if _, err := validateFullCommitID(claim.BaseHeadSHA); err != nil {
		return time.Time{}, time.Time{}, nil, fmt.Errorf("lease %s base_head_sha: %w", claim.ID, err)
	}
	acquired, err := parseCanonicalLeaseTime(claim.AcquiredAt)
	if err != nil {
		return time.Time{}, time.Time{}, nil, fmt.Errorf("lease %s acquired_at: %w", claim.ID, err)
	}
	expires, err := parseCanonicalLeaseTime(claim.ExpiresAt)
	if err != nil {
		return time.Time{}, time.Time{}, nil, fmt.Errorf("lease %s expires_at: %w", claim.ID, err)
	}
	if acquired.After(asOf) {
		return time.Time{}, time.Time{}, nil, fmt.Errorf("lease %s acquired_at is in the future", claim.ID)
	}
	duration := expires.Sub(acquired)
	if duration%time.Minute != 0 || duration < time.Duration(policy.MinLeaseMinutes)*time.Minute ||
		duration > time.Duration(policy.MaxLeaseMinutes)*time.Minute {
		return time.Time{}, time.Time{}, nil, fmt.Errorf("lease %s duration violates policy", claim.ID)
	}
	var released *time.Time
	if claim.ReleasedAt != "" {
		value, err := parseCanonicalLeaseTime(claim.ReleasedAt)
		if err != nil || value.Before(acquired) {
			return time.Time{}, time.Time{}, nil, fmt.Errorf("lease %s has invalid released_at", claim.ID)
		}
		released = &value
	}
	return acquired, expires, released, nil
}

func derivedLeaseStatus(acquired, expires time.Time, released *time.Time, asOf time.Time) string {
	if released != nil && !released.After(asOf) {
		return "released"
	}
	if !expires.After(asOf) {
		return "expired"
	}
	return "active"
}

func runWorkLeaseEvaluation(root string, request WorkLeaseActionRequest, graph WorkGraph) (WorkLeaseEvaluation, error) {
	policy, registry, err := loadWorkLeaseConfiguration(root)
	if err != nil {
		return WorkLeaseEvaluation{}, err
	}
	asOf, err := parseLeaseTimeOrNow(request.AsOf)
	if err != nil {
		return WorkLeaseEvaluation{}, err
	}
	if _, err := auditWorkLeases(policy, registry, graph, asOf); err != nil {
		return WorkLeaseEvaluation{}, err
	}
	return evaluateWorkLeaseAction(policy, registry, graph, request, asOf), nil
}

func evaluateWorkLeaseAction(policy WorkLeasePolicy, registry WorkLeaseRegistry, graph WorkGraph, request WorkLeaseActionRequest, asOf time.Time) WorkLeaseEvaluation {
	result := WorkLeaseEvaluation{Action: request.Action}
	add := func(reason string) {
		if !containsString(result.Reasons, reason) {
			result.Reasons = append(result.Reasons, reason)
		}
	}
	if request.Action != "acquire" && request.Action != "renew" && request.Action != "release" {
		add("invalid-action")
		return result
	}
	if !agentIDPattern.MatchString(request.AgentID) {
		add("invalid-agent-id")
	}
	if !workLeaseIDPattern.MatchString(request.LeaseID) {
		add("invalid-lease-id")
	}
	duration := request.DurationMinutes
	if duration == 0 {
		duration = policy.DefaultLeaseMinutes
	}
	if duration < policy.MinLeaseMinutes || duration > policy.MaxLeaseMinutes {
		add("invalid-duration")
	}

	nodes := map[string]WorkNode{}
	for _, node := range graph.Nodes {
		nodes[node.ID] = node
	}
	claims := map[string]WorkLease{}
	activeByWork := map[string]bool{}
	activeByAgent := map[string]bool{}
	for _, claim := range registry.Claims {
		claims[claim.ID] = claim
		acquired, _ := parseCanonicalLeaseTime(claim.AcquiredAt)
		expires, _ := parseCanonicalLeaseTime(claim.ExpiresAt)
		var released *time.Time
		if claim.ReleasedAt != "" {
			value, _ := parseCanonicalLeaseTime(claim.ReleasedAt)
			released = &value
		}
		if derivedLeaseStatus(acquired, expires, released, asOf) == "active" {
			activeByWork[claim.WorkUnit] = true
			activeByAgent[claim.AgentID] = true
		}
	}

	switch request.Action {
	case "acquire":
		node, ok := nodes[request.WorkUnit]
		if !ok {
			add("unknown-work-unit")
		} else if !containsString(policy.AcquireStatuses, node.Status) {
			add("work-unit-not-ready")
		}
		if _, exists := claims[request.LeaseID]; exists {
			add("lease-id-already-exists")
		}
		if activeByWork[request.WorkUnit] {
			add("work-unit-already-claimed")
		}
		if activeByAgent[request.AgentID] {
			add("agent-already-has-active-claim")
		}
		if _, err := validateFullCommitID(request.BaseHeadSHA); err != nil {
			add("invalid-base-head")
		}
		if len(result.Reasons) == 0 {
			result.ProposedClaim = &WorkLease{
				ID: request.LeaseID, WorkUnit: request.WorkUnit, AgentID: request.AgentID,
				BaseHeadSHA: strings.ToLower(request.BaseHeadSHA), AcquiredAt: asOf.Format(time.RFC3339),
				ExpiresAt: asOf.Add(time.Duration(duration) * time.Minute).Format(time.RFC3339),
			}
		}
	case "renew", "release":
		claim, ok := claims[request.LeaseID]
		if !ok {
			add("lease-not-found")
			break
		}
		if claim.AgentID != request.AgentID {
			add("agent-does-not-own-lease")
		}
		acquired, _ := parseCanonicalLeaseTime(claim.AcquiredAt)
		expires, _ := parseCanonicalLeaseTime(claim.ExpiresAt)
		var released *time.Time
		if claim.ReleasedAt != "" {
			value, _ := parseCanonicalLeaseTime(claim.ReleasedAt)
			released = &value
		}
		if derivedLeaseStatus(acquired, expires, released, asOf) != "active" {
			add("lease-not-active")
		}
		if node, ok := nodes[claim.WorkUnit]; !ok || leaseTerminalStatus(node.Status) {
			add("work-unit-not-continuable")
		}
		if len(result.Reasons) == 0 {
			proposed := claim
			if request.Action == "renew" {
				if request.BaseHeadSHA != "" {
					if _, err := validateFullCommitID(request.BaseHeadSHA); err != nil {
						add("invalid-base-head")
					} else {
						proposed.BaseHeadSHA = strings.ToLower(request.BaseHeadSHA)
					}
				}
				proposed.ExpiresAt = asOf.Add(time.Duration(duration) * time.Minute).Format(time.RFC3339)
			} else {
				proposed.ReleasedAt = asOf.Format(time.RFC3339)
			}
			if len(result.Reasons) == 0 {
				result.ProposedClaim = &proposed
			}
		}
	}
	sort.Strings(result.Reasons)
	result.Allowed = len(result.Reasons) == 0
	return result
}

func leaseTerminalStatus(status string) bool {
	switch status {
	case "DRAFT", "DECOMPOSED", "VERIFIED", "BLOCKED", "BLOCKED_DECISION":
		return true
	default:
		return false
	}
}

func parseLeaseTimeOrNow(value string) (time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return time.Now().UTC().Truncate(time.Second), nil
	}
	return parseCanonicalLeaseTime(value)
}

func parseCanonicalLeaseTime(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil || t.Location() != time.UTC || t.Format(time.RFC3339) != value {
		return time.Time{}, fmt.Errorf("expected canonical UTC RFC3339 timestamp")
	}
	return t, nil
}
