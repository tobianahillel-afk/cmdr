# Case Queue runtime boundary

This directory contains the first bounded executable runtime of **CAP-INV-101 Case Queue**.

## Implemented scope

The runtime is a Go standard-library-only, provider-neutral, persistence-neutral, read-only projection over caller-owned Case snapshots.

It implements the bounded semantics already compiled in `CASE-QUEUE-READONLY-CONTRACT-V1`:

- mandatory single-selected-Tenant and environment scope with no wildcard/default fallback;
- fail-closed cross-tenant Case, Incident, Finding and Saved View inputs;
- Case access filtering without leaking denied Case identity;
- immutable, deep-copy-isolated queue rows;
- deterministic local search, filters and stable sort;
- optional Incident/Finding context with explicit partial/stale/unavailable diagnostics;
- contract-level application of caller-owned Saved View snapshots with per-element permission re-evaluation and no Shared mutation;
- explicit `available`, `empty`, `partial`, `stale`, `permission-filtered` and `view-dirty` states.

The Saved View adapter applies caller-owned Shared snapshots immutably, omits permission-denied elements, uses a generic deterministic default-view fallback for missing or access-denied views, and rejects cross-tenant views. It does not claim Shared Saved View persistence, sharing, archival or management.

## Explicit exclusions

This runtime does not select or implement:

- Case create/update/assignment/status/lifecycle mutation, which remains CAP-INV-102 territory and is governed by OPEN-013;
- Export Job creation or export backend selection;
- Saved View create/update/share/archive/manage mutation;
- Task, Decision or Response Run ownership, aggregation or Command Work Queue behavior;
- canonical Case schema or final Case state-machine changes;
- final Case Queue columns or a dedicated final CAP-INV-101 UI;
- storage/provider/retention/collaboration backend selection;
- a production Case Queue SLO or full-capability completion claim.

Case remains Investigate-owned, Incident remains Command-owned and Saved View remains Shared-owned. External product-runtime dependencies remain deny-by-default.

## Bounded Case Workspace handoff (E10-INV-101D-HANDOFF)

`OpenCaseWorkspace` derives a deterministic, read-only transition to `INV-CAS-001` from the **already permission-filtered queue projection**. It refuses cases hidden by access control or active filters, cross-tenant inputs, stale selected cases, unsupported origins and unsafe scroll/search state.

The return token carries only the tenant/environment, safe origin marker, selected Case reference, scroll position, caller search/filter/sort state and **accessible** Shared Saved View ID/version. Forbidden Saved View elements, foreign Incident data and Case snapshots are not serialized. `RestoreCaseQueue` validates the token structure and rechecks the *current* caller's tenant, environment, Saved View version/access and selected Case visibility before restoring state. Changed permissions or outdated Case data block restoration rather than reuse previous authority.

**The base64 token is not signed, encrypted, persisted or an authorization credential.** It must be treated as untrusted, transient presentation state; any consumer must re-check authorization and avoid exposing the token in logs or durable/public URLs. The final routing/storage/UI transport is not selected here. No Case/Incident/Saved View mutation, assignment, export, Task, Decision or Response Run is introduced.

## Native authenticated read-only HTTP boundary (E10-INV-101H-HTTP)

`NewAuthorizedQueueHTTPHandler` adapts the same `ProjectWithSavedViewFallback` logic to a host-mounted read-only `GET` endpoint. The host injects a verified request/session scope, Case candidate reader and first-party Case/Incident/Finding object-read authorizer; none of these authorities are accepted from client query params or the candidate's `AccessDecision` strings. Denied candidate/nested records are omitted before local search/sort, with no denied Case counts or content in the response. One optional bounded literal `q` parameter provides native server-side search and does not mutate Case state or store a Saved View. Generic fail-closed errors, no-store JSON and a bounded output are enforced.

This is an integration **contract and runtime**, not a production deployment: accepted real IAM, authoritative Case/Incident/Finding providers, TLS, durable audit, throttling, full Case Queue UI, Saved View storage and OPEN-013 mutation remain separate product obligations. The unchanged Go projection is still the only native queue search/sort implementation. Refer to `work/lots/E10-INV-101H-HTTP/HANDOFF.md` for the exact acceptance and Git boundaries.

## Authorized Class-0 HTTP filters and ordering (E10-INV-101I-FILTERS)

The native host-mounted GET adapter now exposes the **existing** Case Queue projection's literal `status`, `owner`, `freshness` filters and deterministic read-only sorting, without a second search/sort engine. Examples: `?q=response&status=review&owner=analyst-a&freshness=available&sort=updated_at:desc` or `?sort=title:asc`. Up to one value per key, valid UTF-8, no controls, bounded values, strict existing field and direction allowlists. The Case owner provider and independent server-side Case/Incident/Finding checks run **before** query matching and sorting; forbidden Case identities and counts are not returned. Client-supplied tenant, environment, permissions, membership, authorization, Saved View and unknown parameters fail closed.

This remains an injected **native Class-0 read**; no production identity service or candidate provider is provided, no live Case Queue UI or Saved View storage, no cursor pagination, mutation, TLS deployment or signed permission attestation. The original Go Project function owns the field/filter/sort semantics and tie-breakers. See `work/lots/E10-INV-101I-FILTERS/HANDOFF.md` for exact evidence, exclusions, ownership and Git/CI status. Parent #139 remains L2-frozen and is not automatically accepted by this child.

## Native authorized Case Workspace navigation (E10-INV-101J-HANDOFF)

The existing Class-0 authenticated Case Queue GET endpoint also accepts an optional
`open_case=EXACT_OPAQUE_CASE_ID` and optional `scroll=NONNEGATIVE_INTEGER` (bounded
by the existing handoff limit), alongside the **same** `q/status/owner/freshness/sort`
projection controls. A successful selection returns the existing
`WorkspaceHandoff` JSON pointing to `INV-CAS-001` and carrying a transient,
**unsigned, non-authorizing** navigation return token. The existing
`OpenCaseWorkspace` engine confirms that the Case remains visible after the
current authenticated server-side Case/Incident/Finding authorization, current
filters, and stale checks. The same existing `RestoreCaseQueue` engine must
re-authorize upon return; a serialized token is never a grant. The source
Case ID is preserved byte-for-byte (including distinct leading/trailing spaces)
after standard URL decoding. Request query values and size remain bounded;
an unrecognized tenant/permission/origin key, malformed scroll or unpaired
scroll fails before contacting any authority.

Without `open_case`, the original queue projection response remains unchanged.
This is **navigation state only**, not a launched workspace, production UI, Case
assignment, durable Saved View, or Case lifecycle mutation. The host still
owns IAM/session issuance, HTTP deployment/TLS/rate limiting/audit and Case
providers. Do not store or use the return token as authorization.

## Native reauthorized Case Queue return (E10-INV-101K-RETURN)

The host-mounted read-only Case Queue adapter supports `POST /case-queue/return`
with JSON `{"return_token":"<opaque state from GET open_case>"}` and the SAME
injected trusted session, current Case candidate provider and independent
Case/Incident/Finding reader authorizer. The return token is **never put in a
GET URL** and is not a credential or trusted permission/tenant selector.

After independently reauthorizing every current Case and nested context source,
the adapter invokes the existing first-party `RestoreCaseQueue`, which checks
the caller's tenant/environment, selected Case, stale state, Saved View/version,
current filters, sort and scroll. The response is a read-only ReturnState,
not a raw Case or evidence preview; unauthorized/stale/changed Cases fail closed
without partial Queue data. Accepts only a bounded single `return_token`
JSON key; unknown/duplicate keys, query parameters and wrong route/method are
rejected before host owner access. The host must mount this route and enforce
production authentication, TLS, CSRF protection, request-body log redaction,
rate limits and audit; none are implied by this stand-alone Go handler.

The prior GET listing/filter/sort/open_case behavior and canonical Queue
projector are untouched. No mutation, Saved View persistence, export,
production IAM, external dependencies or API costs.


## Host-selected Shared Saved View in native Queue HTTP (E10-INV-101L-VIEW)

An optional `NewAuthorizedQueueHTTPHandlerWithSavedViews` constructor accepts the host's **trusted Shared Capabilities** `QueueSavedViewProvider`. This adapter resolves the current authenticated principal's selected, versioned, permission-rechecked Saved View **on each request** using the host's trusted tenant/environment; the browser can neither supply a view snapshot nor select a view or tenant by URL. The host is responsible for Shared Saved View read/use ACL and **individual field/filter/sort permission re-evaluation** when producing the `SavedView` and `AccessDecision` values; this is **not** a second Shared storage or IAM implementation.

The same existing first-party `ProjectWithSavedViewFallback` strips `deny` elements without leaking forbidden names, applies allowed view filters/ordering, and reports permission-filtered and fallback states. Existing GET list and GET `open_case` use the selected view, and POST `/case-queue/return` re-resolves it independently before the existing `RestoreCaseQueue` checks version, Case visibility, tenant, environment, scroll, and current access. A version change or lost Shared use right denies the old token; a denied current view produces the existing generic default-view fallback on fresh GET. Cross-tenant and malformed Shared snapshots fail closed with no Case contents. No extra Case permission is obtained from a view.

This is still **host-injected read-only runtime**, not a working persistent Shared provider, real IAM, production UI, CSRF-protected deployed HTTP service, Saved View mutation, Case assignment, export, or full CAP-INV-101 acceptance. Hosts without the new adapter retain the previous `NewAuthorizedQueueHTTPHandler` default-view behavior. The product/Shared owner must decide actual Saved View selection and lifecycle; no new canonically authoritative storage is introduced.

## Authorized native HTTP pagination (E10-INV-101M-PAGE)

Optional `GET /case-queue?limit=1..100` pages the SAME existing server-authorized, Shared-view-filtered, deterministic Queue projection. Responses keep the original projection JSON shape with only the current `rows` and an optional `next_cursor`. Supply the original filters/sort/limit plus `cursor` for the following page. Each page **reauthorizes all candidate Case/Incident/Finding objects and re-resolves the current Shared Saved View** before checking the cursor against a SHA-256 fingerprint of the complete current projection and view. Changed ACL, Case metadata, tenant/environment or view/version invalidates old cursors rather than disclosing partial data. Cursors contain no raw records or bearer and are not credentials or signed evidence. No paging flags: old full GET unchanged. Paging cannot accompany `open_case` and does not change `POST /case-queue/return`. Limit 100, candidates <=5000, no datastore/IAM/TLS/CSRF/production provider added. See `work/lots/E10-INV-101M-PAGE/HANDOFF.md`.
