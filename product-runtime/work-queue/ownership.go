package workqueue

import (
 "context"
 "errors"
 "strings"
 "unicode/utf8"
)

// AuthorizedCanonicalOwnership is the Ownership Register's independent
// Class-0, source-authorized projection for ONE selected Command work item.
// Product owner and operational assignee are different ownership concepts.
// A nil projection means the register is unavailable or does not publish
// the selected object's canonical owner; the UI must say unknown.
type AuthorizedCanonicalOwnership struct {
 ProductOwner string `json:"product_owner,omitempty"`
 RegisterSource string `json:"register_source,omitempty"`
 RegisterVersion string `json:"register_version,omitempty"`
 IdentityStatus string `json:"identity_status,omitempty"`
}

// CanonicalOwnershipOwner is a separate trusted host adapter. The host MUST
// check access to the Ownership Register and identity source independently;
// the Incident/Task CanRead grant is insufficient to disclose those records.
type CanonicalOwnershipOwner interface {
 ReadAuthorizedCanonicalOwnership(ctx context.Context, principal, tenant, environment string, kind Kind, id string) (*AuthorizedCanonicalOwnership, error)
}

var errCanonicalOwnership = errors.New("canonical ownership source unavailable")

func validateAuthorizedCanonicalOwnership(v *AuthorizedCanonicalOwnership) (*AuthorizedCanonicalOwnership, error) {
 if v == nil { return nil,nil }
 for _,field := range []struct{value string; max int}{
  {v.ProductOwner,128},{v.RegisterSource,128},{v.RegisterVersion,96},
 } {
  if len(field.value)>field.max || !utf8.ValidString(field.value) ||
   (field.value!="" && (invalidScope(field.value) || strings.TrimSpace(field.value)!=field.value)) {
   return nil,errCanonicalOwnership
  }
 }
 switch v.IdentityStatus {
 case "", "resolved","identity-unresolved","unknown":
 default: return nil,errCanonicalOwnership
 }
 if v.ProductOwner!="" && v.RegisterSource=="" {
  return nil,errCanonicalOwnership
 }
 if v.RegisterVersion!="" && v.RegisterSource=="" {
  return nil,errCanonicalOwnership
 }
 // Source ownership is supplied by the host, not re-derived from kind
 // strings, a Command row, or any client query parameter.
 return v,nil
}
