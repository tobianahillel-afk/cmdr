# E10-INV-110B — authenticated HTTP delivery adapter for Case Timeline

## Functional gap and reuse-first selection

Current first-party code already supports `OpenCaseTimeline`, trusted Case ID discovery, authoritative Case and Timeline Entry checks, typed source identity resolution, owner snapshots and immutable chronology. The missing bounded user-facing seam was a server HTTP request -> authenticated Case selection -> JSON read-only projection.

Alternatives:
1. **Selected**: Go standard-library `net/http`, `net/url`, `encoding/json` and already implemented owner/ACL pipeline. The application/gateway injects existing trusted providers. Zero third-party runtime dependencies, zero schema/persistence migrations, and no new permission engine.
2. A new web framework, generalized GraphQL engine, OAuth provider, event store or timeline index: **not selected**; no supported deployed Case/Shared stores or trusted migration contract is available.
3. Treat Case Queue `return_token`, browser `X-Tenant`, `scope` flags, or an unsigned navigation token as a server permission: **rejected**; the canonical security contract requires server-side authorization.

## Implemented behavior

`NewCaseTimelineHTTPHandler(sessions, checker, catalog, entries, sources)` returns a Go `http.Handler`. The trusted server session resolver accepts the request and returns an authenticated principal context and scoped tenant/environment. The handler forwards only a single bounded `case_id` selector (GET only, without body). All object and linked-source access decisions are obtained from the trusted AccessChecker and owner adapters via `OpenCaseTimeline`.

- Authorization fails before listing any Case Timeline IDs; a second Case authorization check guards against revocation after discovery.
- Canonical Timeline and source entry decisions are independent; absent/denied linked entries are not serialized.
- Responses disable caching and MIME sniffing. Denial/provider failures expose one opaque 404; invalid request 400; no session 401; non-GET 405.
- Existing ordering and provenance, observed/inferred kind, visible supersession and partial/filtered status remain intact.
- Tests use `httptest` and a two-tenant in-memory first-party owner fixture. This is **real executable request-to-projection integration in-process**, not a deployed HTTP service or production auth/persistence integration.

## Remaining production/security gates

- The injected `CaseTimelineSessionResolver`, `CaseTimelineCatalog`, `AccessChecker`, `EntryResolver`, `SourceOwnerResolver` remain interfaces. Production implementations must authenticate principals from the actual server session, enforce tenant/Case/entry/source membership and preserve scope+permission consistency.
- The host gateway must supply a canonical route, TLS, CSRF protection as applicable to its chosen auth mechanism, audit events for sensitive reads/denials, rate limits, traffic logging/redaction and an accepted error/correlation-ID contract.
- Real Case/Shared owner repositories, transactional authorization snapshot or fail-closed consistency checks, read-only UI path from Case Queue, and p95 of **server/API** response (not projection-only sample means) remain absent.
- Product Spec, Shared Timeline ownership and open decision `OPEN-013` are untouched. The existing `E10-INV-110B-RUNTIME` remains READY until acceptance gates and independent L2 findings are reconciled. `E10-INV-110C-NAV` remains BLOCKED; the HTTP adapter is not final navigation.

## Verification expectations and handoff

CI must run `investigation-timeline-runtime` with race, Go vet, coverage, behavioral HTTP tests, Go formatting; the engineering-kernel must independently pass for the exact committed SHA. Evidence is recorded in GitHub Actions, not inferred from the implementation. Schedule a **single** new independent L2 for this material HTTP trust boundary, pinned to its published exact SHA/diff; no duplicate for correction-only commits. Never call a successful in-memory integration a deployed production workflow.
