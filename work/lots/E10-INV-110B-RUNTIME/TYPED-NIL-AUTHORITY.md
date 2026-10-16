# E10-INV-110B — fail closed on typed-nil authority dependencies

## Reproduction and scope

Go interface values can be non-nil while holding a typed-nil pointer. Both read-only access entry points previously checked `checker == nil` / `resolver == nil` and then invoked their methods. With `(*AccessCheckerImplementation)(nil)` or `(*EntryResolverImplementation)(nil)`, this could panic if server dependency wiring is uninitialized; a typed-nil `context.Context` could likewise panic on `Err()`. A recovered panic at an API boundary would still cause a denial-of-service/error rather than a predictable fail-closed response.

An early `missingAuthority` guard using the standard Go `reflect.Value.IsNil` for nil-able interface dynamic types now returns the existing `Result{Status:"deny"}`, before a provider call or an owner read. The cost is paid only at the integration entry boundary, not per Timeline Entry. Six adversarial tests cover the three typed-nil dependencies across the two independent adapters.

## Reuse and limitations

Reuse Go standard library and CMDR's existing `AccessChecker`, `EntryResolver`, `ProjectWithAccessChecker` and `ProjectFromOwnerSnapshots`. No dependencies, products, schema migrations or new runtime authority/provider. This does **not** guard against an intentionally broken non-nil provider implementation, does **not** provide authenticated ACL checks, does **not** make provider decisions atomic, and is **not** proof of Windows/macOS behavior without CI.

Do not promote B to VERIFIED until exact-SHA CI, material L2 findings and external integration obligations are addressed. Previous L2 for this same B lot is not an approval of this distinct code snapshot. Canonical Product Spec and OPEN-013 remain untouched.
