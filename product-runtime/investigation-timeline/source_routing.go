package investigationtimeline

import "context"

// SourceOwnerIdentity is a trusted owner service's resolved identity.
// These are NOT fields to accept from a browser, caller-owned Timeline Entry,
// or a SourceRef prefix. Its owning service must bind the principal, Case,
// tenant, environment, object identity and (where appropriate) version.
type SourceOwnerIdentity struct {
	Reference   string
	Kind        string
	ObjectID    string
	Tenant      string
	Environment string
	CaseID      string
}

// SourceOwnerResolver is an explicit trusted server-only integration seam.
// The incoming AccessRequest is a source-reference lookup scoped to the
// authenticated Case, tenant and environment. Resolving does not grant read:
// the returned owner/type gets a separate canonical AccessChecker decision.
type SourceOwnerResolver interface {
	ResolveSource(context.Context, AccessRequest) (SourceOwnerIdentity, error)
}

func canonicalSourceReadPermission(kind string) (string, bool) {
	switch kind {
	case "evidence":
		return "perm.investigate.evidence.read", true
	case "telemetry-event":
		return "perm.shared-capabilities.telemetry-event.read", true
	case "decision":
		return "perm.govern.decision.read", true
	case "incident":
		return "perm.command.incident.read", true
	default:
		return "", false
	}
}
