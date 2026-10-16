# E10-INV-110B — reusable authenticated Case discovery to Timeline journey

## Functional delivery

Prior state: the local analyst page discovered Cases through a *demo-only* fixture reader, while the reusable `NewCaseTimelineHTTPHandler` supported only opening a known Case ID. A future host would otherwise have to rebuild Case listing, independent reader checks and response policies.

This increment introduces `CaseListCatalog` and `NewCaseDiscoveryHTTPHandler` as a **server-authorized, provider-neutral, read-only** Case selection boundary. The trusted owner service supplies candidate IDs, not permission verdicts. The handler obtains tenant, environment and principal from `CaseTimelineSessionResolver` (exact same identity boundary as `OpenCaseTimeline`) and checks `perm.investigate.case.read` on **each candidate** with the existing `AccessChecker`, prior to release. Unsupported/failed owner checks fail closed without returning a partial list. Results have deterministic order and contain only authorized Case IDs, with no Timeline Entry/source metadata.

`NewMemoryCaseTimelineHTTPHandlers` assembles the Cases and Timeline handlers on **one shared, immutable in-process owner snapshot**; the loopback-only browser demonstration uses both. Existing `NewMemoryCaseTimelineHTTPHandler` remains backwards compatible. The host may inject real authenticated Case-owner and ACL providers later; none are silently bundled.

## Reuse-first design choice

- Chosen: existing Go `net/http`, `CaseTimelineSessionResolver`, `AccessChecker`, `AccessRequest`, in-memory owner fixture, and `sort.Strings`.
- Rejected for this scope: a second identity service, ad-hoc browser authority, new permissions schema, database/index, external case catalogue, generic policy framework, separate SPA.
- Limits: bounded to 5,000 candidate IDs, 1,024 bytes per ID and 8 MiB aggregate candidate strings. Duplicate/malformed catalogs are rejected before any owner ACL calls.
- The Case-only listing does **not** expose Timeline metadata and cannot substitute for per-entry or per-source authorization when opening a Case.

## Security and operational boundary

Server-side auth and current owner ACL must both pass; unknown/denied Case IDs are omitted. Any checker error prevents all output. No HTTP body or query selectors accepted. Responses are no-store with JSON/nosniff/referrer/Vary headers. Memory fixtures stay loopback-only, not a deployed source or real production auth/ACL. Host deployment still owns TLS, session verification, revocation semantics, audit, rate limits and durable Case/Shared/source repositories.

Checks: `go test -race ./...` for the Timeline module; full Go engineering-kernel as required; integration tests validate both routes with the same owner, cross-tenant and same-tenant restricted Case confidentiality, invalid credentials, source authorization and invalid candidate behavior. These are **process-local** end-to-end integration proofs only; not claims that external providers or persistence exist.

Governance: B stays READY pending its actual acceptance and L2; C/D/E stay dependency-blocked. Product Spec and OPEN-013 remain unchanged. No claim of production readiness, new migration or permanent storage.
