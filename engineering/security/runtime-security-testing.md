# Runtime security test obligations and coverage floors

This control is strict about runtime security without fabricating evidence before executable product code exists.

## Runtime states

Runtime security is derived from the architecture registry **and the checked-out filesystem**.

- `not-applicable`: there are no `product-runtime` boundaries.
- `preimplementation`: a runtime boundary is reserved but every owned root contains only `README.md` / `.gitkeep` markers.
- `measured`: every registered runtime boundary contains executable/package content.
- `mixed`: implemented and preimplementation runtime boundaries coexist.

Every product-runtime boundary requires an explicit security ownership scope in `runtime-security-policy.json`, including owner, security-critical paths, authorization requirement/rationale and tenant-isolation requirement/rationale.

## Preimplementation safety

A preimplementation scope may declare future authorization and tenant-isolation obligations, but it may **not** claim coverage percentages, authorization-negative PASS or tenant-isolation PASS. Its evidence state must match the filesystem-derived `preimplementation` state.

The moment any non-marker file appears under a product-runtime root, the scope becomes `implemented`; stale preimplementation evidence fails closed automatically.

This permits architecture/performance budgets to be registered before code while preventing false runtime PASS claims.

## Implemented runtime activation contract

For every implemented runtime scope:

- real changed-path evidence is required;
- evidence is bound to the exact checked-out Git commit;
- global runtime coverage must be at least **80%**;
- changed security-critical coverage must be at least **90%** when a security-critical path changes;
- if authorization is required, `SEC-AUTH-NEG-001` must be active and negative tests must pass;
- if tenant isolation is required, `SEC-TENANT-ISO-001` must be active and negative tests must pass.

For mixed repositories, the global measured coverage applies to implemented runtime code only; preimplementation scopes still cannot claim runtime test evidence.

## Evidence lifecycle

The committed pilot evidence is currently a truthful `preimplementation` marker because `product-runtime/context-envelope/` contains only `README.md`.

E9-PILOT-001C must replace/generate measured evidence in the CI workspace as soon as executable runtime content is introduced. It must also activate and satisfy the required authorization and tenant-isolation negative-test gates before the runtime can pass validation.
