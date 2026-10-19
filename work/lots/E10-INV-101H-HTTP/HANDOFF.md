# E10-INV-101H-HTTP — Native Case Queue read-only server integration

## Canonical scope and ownership

- Product: `CAP-INV-101` Case Queue; requirements `REQ-PROD-014`, `REQ-PROD-020`; canonical `cmdr-product-spec/07-investigate/modules/cases-and-evidence/capabilities/case-queue.md`. Case is Investigate-owned, optional Incident context is Command-owned, Findings are Investigate-owned; Shared owns Saved Views.
- Engineering: bounded new task `E10-INV-101H-HTTP`, parent `E10-INV-101`, depends on verified `E10-INV-101F`. This is **not** an extension of the synthetic `E10-INV-101G-LOOKUP` browser task or a business child of its Git PR. The work graph previously reported no READY node; this new unit is compiled against existing verified canonical/projection sources.
- Code owner: `product-runtime/case-queue/http.go`, within the **native** Case Queue package, reusing `ProjectWithSavedViewFallback`. No modifications to canonical spec, Event Search, Timeline demo, engineering tools or security thresholds.
- Git: branch `feature/e10-inv-101h-native-authorized-queue-http` based on frozen PR #137 SHA `b98fac7334985fb028651ff946b86f170cb048eb`, intentionally parallel to PR #138 (which is L2-frozen and CI security-red). This Git ancestry is not the engineering dependency; merge/reconcile #138 separately.

## Requirement → behavior → evidence → outstanding

- **Native read-only Case Queue**: before this task the read-only Project engine was Go-callable only and the demo exposed a separate Case ID picker. The new `NewAuthorizedQueueHTTPHandler` exposes a host-registered `GET` projection endpoint, reusing exactly the native search and stable sort. An optional literal `q` query is bounded to 256 valid UTF-8 bytes; tenant/environment are never HTTP parameters.
- **Server-side permissions**: the host must inject an authenticated session resolver, trusted candidate provider and independent object permission checker. Every candidate Case is checked. Nested Incident and Finding references are separately checked by type and identifier; denial removes the source data *before* calling the projector. Provider-supplied `AccessDecision` fields are discarded as authority; allowed objects are explicitly tagged after server check. Foreign tenant candidates, invalid owner decisions and provider errors fail closed with generic responses and **no partially authorized projection**. Denied candidate count/identity is not exposed.
- **Transport**: GET-only, no-store JSON, `Vary: Authorization`, `nosniff`, max 5000 candidate Cases, 4 MiB response and explicit generic errors. Host is responsible for TLS, production session/IAM, actual object-provider integrity, rate limits and durable audit logs. Merely defining these interfaces does **not** implement a production provider or deploy an endpoint.
- **Security acceptance**: Go integration tests use an authenticated request and authorized/denied Case and nested source records; verify no denied Case/title/Incident/Finding preview appears, no unauthorized-case count, no extra source request, no query-scoped tenant spoofing or client-provided ACL bypass; reject untrusted requests, malformed fields, scope mismatch, unavailable authorizer, mutations. No external runtime dependencies.

## Scope deliberately excluded

- Full Case Queue UI, native pagination/cursor contracts (still draft); final columns, persisted Saved View management, Incident/Finding providers, production IAM and authentication, Case statuses/assignment via `CAP-INV-102`, export and response workflows. `OPEN-013` remains OPEN.
- This is a **first-party HTTP integration surface with injected providers** plus behavior tests, **not** an actually deployed or end-to-end production Case Queue. Canonical `CAP-INV-101` still defined/planned, not accepted. Do not represent synthetic session fixtures as real production credentials.
- The thin handler's per-resource permission model must be integrated with CMDR's future accepted real IAM and authoritative providers before production. No new technology decision is silently accepted.

## Status and next verified actions

- Code and tests are published in this child PR, but status remains `IMPLEMENTED` (NOT `VERIFIED`) pending exact-head GitHub CI, any applicable static/security checks, independent risk-based L2 and merge approvals.
- Required check floor `CHK-SECURITY-TESTS >=90%` **must not be waived**. Parent #137 had overall CI success; PR #138, which is parallel and frozen, has runtime-success and mandatory security failure 89.34%. This new task must obtain its own exact-head evidence.
- If native Queue and synthetic demo browser are integrated later, preserve authorization semantics and avoid a second Case Queue search/sort engine. Determine real provider/IAM/contracts first; do not merge older draft branches just because later features compile.
