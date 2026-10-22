# E10-INV-110B — server-side authorization adapter boundary

## Reuse-first decision

Use Go standard-library `context.Context` and a minimal first-party `AccessChecker` interface inside the existing Investigation Timeline runtime. Reuse the Case Queue pattern of explicitly separate Case and linked-object access decisions rather than adopting a new RBAC/ABAC vendor engine, caching permission outcomes or constructing an unowned Shared Timeline service. The bounded `Project` function and its 18 canonical fixtures remain unchanged.

## Trust boundary and deliberate limitations

`ProjectWithAccessChecker` ignores **all caller-supplied** `case_permission`, `timeline_permission` and `source_permission` flags, then consults a server-supplied checker separately for:
1. `perm.investigate.case.read` against a nonempty Case ID in tenant/environment scope.
2. `perm.shared-capabilities.timeline-entry.read` against every same-tenant Timeline Entry ID.
3. An independent `source-owner-read` obligation against each accessible entry's linked source reference.

A foreign-tenant entry is masked without an authorization lookup, unknown/missing/denied decisions are not upgraded, errors/cancellation/invalid decisions fail closed, and the adapter copies entries rather than modifying caller arrays. A denied Timeline Entry never induces a linked-source check. The unmodified low-level `Project` still accepts caller-derived flags; **only trusted server code** may provide `AccessChecker` and call the wrapper.

This increment **does not ship an authoritative provider implementation, authentication middleware, HTTP endpoint, real audit trail, user identity binding, source-owner typed routing, or production-policy deployment**. It cannot be treated as a completed end-to-end authorization integration or a production security approval. `source-owner-read` is an internal adapter obligation: the owner must map it to the canonical conditional permission; default-deny until such an integration is implemented. Never promote this code to a production authorization authority by assuming a test stub or caller-provided checker is trusted. OPEN-013 remains unchanged.

## Evidence and handoff

Scope: CAP-INV-110 bounded read-only Timeline, product spec immutable. The new adapters and adversarial tests are in `product-runtime/investigation-timeline/access.go` and `access_test.go`. Go test/race/coverage/CI and independent L2 must be verified on the exact published commit; no result or provider claim exists merely because the files are published. This is intentionally the minimum safe bridge to future server/provider wiring; no third-party APIs, paid tools, migrations, or new runtime dependencies.
