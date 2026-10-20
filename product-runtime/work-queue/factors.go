package workqueue

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

// Dimension is a single owner-authorized source claim, not a numeric rank.
// A factor never confers authority to change Command priority or a source's
// technical severity, confidence, or contractual SLA policy.
type Dimension struct {
	Value string `json:"value"`
	Source string `json:"source"`
	ObservedAt string `json:"observed_at,omitempty"`
}

// AuthorizedFactors retains the six distinct canonical CAP-CMD-104 factors.
// A nil factor means unknown/not disclosed, NEVER zero, healthy, or low.
type AuthorizedFactors struct {
	Priority *Dimension `json:"priority,omitempty"`
	Severity *Dimension `json:"severity,omitempty"`
	Urgency *Dimension `json:"urgency,omitempty"`
	Impact *Dimension `json:"impact,omitempty"`
	Confidence *Dimension `json:"confidence,omitempty"`
	SLARisk *Dimension `json:"sla_risk,omitempty"`
}

// FactorsOwner is an optional first-party *host* capability. In addition to
// Owner.CanRead on the selected Incident/Task, the host MUST independently
// authorize every returned factor against its actual owner/source permission
// (including severity/confidence and SLA). Inaccessible factors MUST be nil.
// Never resolve source permissions from an HTTP query or a caller-provided
// factor. Absence of this adapter yields six unknown dimensions in the UI.
type FactorsOwner interface {
	ReadAuthorizedFactors(ctx context.Context, principal, tenant, environment string, kind Kind, id string) (AuthorizedFactors, error)
}

var errFactors = errors.New("invalid or unavailable authorized coordination factors")

func validateAuthorizedFactors(f AuthorizedFactors) (AuthorizedFactors, error) {
	for _,dim := range []*Dimension{
		f.Priority, f.Severity, f.Urgency, f.Impact, f.Confidence, f.SLARisk,
	} {
		if dim == nil { continue }
		if len(dim.Value) > 96 || len(dim.Source) > 128 ||
			invalidScope(dim.Value) || invalidScope(dim.Source) ||
			!utf8.ValidString(dim.Value) || !utf8.ValidString(dim.Source) ||
			strings.TrimSpace(dim.Value) != dim.Value ||
			strings.TrimSpace(dim.Source) != dim.Source ||
			len(dim.ObservedAt) > 48 {
			return AuthorizedFactors{}, errFactors
		}
		if dim.ObservedAt != "" {
			t,err := time.Parse(time.RFC3339Nano, dim.ObservedAt)
			if err != nil || t.IsZero() || strings.TrimSpace(dim.ObservedAt) != dim.ObservedAt {
				return AuthorizedFactors{}, errFactors
			}
		}
	}
	return f,nil
}
