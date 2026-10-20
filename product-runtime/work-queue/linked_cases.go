package workqueue

import (
 "context"
 "errors"
 "strings"
 "time"
 "unicode/utf8"
)

// AuthorizedLinkedCases is an independently source-authorized, read-only
// summary of Investigate-owned Cases ALREADY linked to ONE Command Incident.
// Empty Cases means the source authorized an empty set; a nil projection
// means the linking/Investigate source is unavailable (unknown, not empty).
type AuthorizedLinkedCases struct {
 Cases []LinkedCaseSummary `json:"cases"`
}

type LinkedCaseSummary struct {
 CaseID string `json:"case_id"`
 Status string `json:"status"`
 LinkSource string `json:"link_source"`
 LinkVersion string `json:"link_version,omitempty"`
 // These fields come ONLY from the independently authorized Investigate
 // Case source, not from Shared link metadata or the Command Incident.
 StatusSource string `json:"status_source,omitempty"`
 StatusVersion string `json:"status_version,omitempty"`
 StatusObservedAt string `json:"status_observed_at,omitempty"`
 StatusFreshness string `json:"status_freshness,omitempty"`
}

// LinkedCasesOwner is a trusted HOST boundary. The host MUST check BOTH
// the Shared Object Linking relation and Investigate Case read permission,
// individually for every returned Case, in the supplied session scope.
// Command Incident CanRead alone NEVER grants Case or Finding access.
type LinkedCasesOwner interface {
 ReadAuthorizedLinkedCases(ctx context.Context, principal, tenant, environment, incidentID string) (*AuthorizedLinkedCases, error)
}

var errLinkedCases = errors.New("linked Case source unavailable")

func validateLinkedCases(value *AuthorizedLinkedCases) (*AuthorizedLinkedCases, error) {
 if value == nil { return nil,nil }
 if len(value.Cases)>16 { return nil,errLinkedCases }
 seen:=make(map[string]bool,len(value.Cases))
 for _,item:=range value.Cases {
  for _,field:=range []struct{ s string; max int; required bool }{
   {item.CaseID,256,true},
   {item.LinkSource,128,true},
   {item.LinkVersion,96,false},
   {item.StatusSource,128,false},
   {item.StatusVersion,96,false},
  }{
   if (field.required && field.s=="") || len(field.s)>field.max ||
    !utf8.ValidString(field.s) ||
    (field.s!="" && (invalidScope(field.s) || strings.TrimSpace(field.s)!=field.s)) {
    return nil,errLinkedCases
   }
  }
  switch item.Status {
  case "intake","triage","investigating","awaiting-input","review","completed","archived","unknown":
  default: return nil,errLinkedCases
  }
  if item.StatusObservedAt!="" {
   if len(item.StatusObservedAt)>64 || !utf8.ValidString(item.StatusObservedAt) ||
      strings.TrimSpace(item.StatusObservedAt)!=item.StatusObservedAt {
     return nil,errLinkedCases
   }
   if _,err:=time.Parse(time.RFC3339Nano,item.StatusObservedAt);err!=nil {
     return nil,errLinkedCases
   }
  }
  switch item.StatusFreshness {
  case "", "unknown":
  case "current", "stale", "partial":
   // A positive source-state assertion without source identity and
   // observation time would mislead an Incident Commander.
   if item.StatusSource=="" || item.StatusObservedAt=="" {
    return nil,errLinkedCases
   }
  default: return nil,errLinkedCases
  }
  if seen[item.CaseID] { return nil,errLinkedCases }
  seen[item.CaseID]=true
 }
 return value,nil
}
