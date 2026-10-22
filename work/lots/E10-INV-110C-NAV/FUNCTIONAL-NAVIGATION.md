# E10-INV-110C-NAV — host-integrated source navigation draft (first published increment)

## Functional delivery

A Case analyst can resolve a Timeline Entry selected from the existing authenticated Case -> Timeline HTTP flow into a **typed owner-object navigation draft**. The new `Navigate` handler is returned by `NewCaseInvestigationHTTPHandlers`, by the first-party memory reference wiring, and is mounted by the loopback-only `case-timeline-demo` at `/case-source-navigation`. The host owns any final object route and must repeat destination authorization when the object is actually opened.

GET parameters: required `case_id` and `entry_id`; optional presentation-only `kind`, `q`, `from`, `to`, `limit`, `offset` reused from the Timeline filter parser. Return context is an encoded, allowlisted Case selection **only**; it never grants authority, carries no tenant, owner object ID, token, free-form URL, redirect or external route.

The handler reuses `OpenCaseTimeline`, which already enforces server-injected Case, Timeline Entry and canonical typed-source ACLs before returning permitted entries. Navigation then **rechecks** the Case and selected Entry access, reloads the owner Entry to detect stale time/version/source, re-resolves the typed destination owner, verifies Case/tenant/environment linkage, and rechecks the canonical destination object read. Unknown/denied/foreign/stale destinations yield indistinguishable generic 404 with no sensitive object metadata. HTTP never trusts caller-supplied permissions or source identities.

## Reuse-first decision

Option A: build a new URL resolver, dynamic cross-product router or third-party navigation framework — rejected: canonical route formats are still draft and such a router would conflate UI return hints with server authority.
Option B (selected): reuse first-party Go `OpenCaseTimeline`, `AccessChecker`, `EntryResolver`, `SourceOwnerResolver`, existing HTTP session resolver and Timeline filter parser. A small typed JSON draft is returned; no new package, persistent model, Timeline Entry ownership change, redirect or open decision resolution.

## Tests and limits

Behavioral tests cover authorized Case -> source navigation, denied/cross-tenant entries, revoked destination access after timeline projection, stale owner version, strict HTTP query allowlist, neutral denial and return selection. The reference host is still memory-backed and synthetic. **This is NOT** a deployed identity service, persistent owner repository, a full UI route, transactionally atomic ACL check, a production-ready cross-product deep link or fulfillment of the full CAP-INV-110 capability. `OPEN-013` remains unresolved and product spec unchanged.

CI evidence is pending at initial publication; do not claim success until an exact-SHA GitHub Actions run is checked. Next functional increment: integrate browser initiation and safe in-page return handling using the newly exposed route; run end-to-end tests and an independent read-only L2 at the stable milestone.
