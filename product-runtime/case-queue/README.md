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

Optional `GET /case-queue?limit=1..100` pages the SAME existing server-authorized, Shared-view-filtered, deterministic Queue projection. Responses keep the original projection JSON shape with only the current `rows` and an optional `next_cursor`. Supply the original filters/sort/limit plus `cursor` for the following page. Each page **reauthorizes all candidate Case/Incident/Finding objects and re-resolves the current Shared Saved View** before checking the cursor against a SHA-256 fingerprint of the complete current projection and view. Changed ACL, Case metadata, tenant/environment or view/version invalidates old cursors rather than disclosing partial data. Cursors contain no raw records or bearer and are not credentials or signed evidence. No paging flags: old full GET unchanged. The paged Case-open and return integration in E10-INV-101N-RETURN extends this initial listing-only behavior; current Case/Incident/Finding ACL and cursor validation remain mandatory. Limit 100, candidates <=5000, no datastore/IAM/TLS/CSRF/production provider added. See `work/lots/E10-INV-101M-PAGE/HANDOFF.md`.

## Native paged Case selection and return (E10-INV-101N-RETURN)

The former Queue pagination has now been integrated with the **same existing** Case Workspace handoff and return engine. A selected Case from a bounded currently authorized page may use `GET /case-queue?limit=2&cursor=<opaque>&sort=title:asc&open_case=CASE_ID&scroll=91`. The selected Case must be on exactly that page; its placement or page cursor grants **no permission**. The normal first-party Case/Incident/Finding owner checks and current host Shared Saved View apply before both the open and the independent `POST /case-queue/return`.

The transient `ReturnState` includes optional `page_limit` and `page_cursor`, alongside existing Case, filters, sort, scroll and saved-view version. On return the SAME current owner-authorized projection and cursor fingerprint are checked again. Revocation, metadata or Saved View changes that invalidate the page cause generic fail-closed instead of restoring a stale selection. Old unpaged clients and return tokens remain compatible. No token is authorization or a durable snapshot, and no production IAM/provider/UI was introduced. See `work/lots/E10-INV-101N-RETURN/HANDOFF.md`.

## Native Command Incident → authorized Case Queue selection (E10-INV-101O-INCIDENT)

The canonical CAP-INV-101 §17 Command Incident transition can now select the
already projected, **currently authorized** linked parent with
`GET /case-queue?incident_id=INCIDENT_ID`. The `incident_id` value is an
opaque, **exact** identity: uppercase/lowercase and whitespace are meaningful.
It does not grant a Case or Command Incident permission, choose a tenant or
modify either object. The current injected host authenticates the session,
independently reauthorizes each candidate Case and its Command-owned Incident,
then the SAME `ProjectWithSavedViewFallback` applies this extra selector
only to the Incident context that survived the owner ACL and optional source
availability check. A denied or unavailable Incident never matches; the
projection does not reveal its original identifier.

This optional selector composes with the existing `q/status/owner/freshness/sort`,
host-selected Shared Saved View, `limit/cursor`, and native
`open_case/scroll` navigation. The existing short-lived return token retains
the opaque Incident selector alongside the selected page. An independent
`POST /case-queue/return` rechecks current Case **and Incident** permissions,
Shared View and the original paginated projection before restoring context.
Changing or revoking rights fails closed without partial Case or Incident
source metadata.

**Boundary:** the current synthetic CaseSnapshot contains **one optional
IncidentContext**. This implements a bounded linkage projection only, *not*
the complete canonical Case `linked-incidents` many-to-many relation, a
Command relationship resolver, production host-side object provider, or a
deployed Incident Detail link. The host must independently authenticate the
Command Incident and authorize its current projection; merely supplying an
`incident_id` string is not proof of any relationship or read grant.
No Case/Incident mutation, Export Job, assignment, provider storage, new
Saved View engine, third-party dependency or paid API.


## Read-only Case Queue → Inspector selection (E10-INV-101P-INSPECT)

The canonical CAP-INV-101 §17 includes `Case Queue → Inspector` before/alongside
full `Case Workspace` navigation. Optional native
`GET /case-queue?inspect_case=EXACT_CASE_ID` returns **one existing
source-authorized Case Queue `Row`**, including currently authorized
Incident/Finding context and the `partial`/`stale` indicators. It never
creates another Case owner, Inspector search engine, source data provider,
mutation path or UI implementation.

The selection composes with the SAME existing Queue query, status/owner/
freshness/Incident filters, host-selected Shared Saved View, deterministic
sorting and optional `limit/cursor`. A selected Case must appear in the
currently authorized and filtered page. Each request rechecks the verified
session and independent Case/Incident/Finding object read authorizations
*before* selecting the visible row. If a source is revoked, its content is
omitted, or an obsolete page cursor is rejected; denied/filtered/out-of-page
Case IDs produce only the generic not-found response. A stale yet authorized
Case is still inspectable with a stale indicator; the existing full Workspace
handoff keeps its stricter stale-selection refusal. `inspect_case` cannot be
combined with `open_case` or the latter's `scroll`.

Only the read-only *data projection* needed for a future Shared Inspector is
provided. The canonical Shared Inspector's keyboard, focus, responsive,
accessible and full-screen behaviors and the final native Case Queue UI are
NOT implemented here. The host still owns IAM, Case/Incident/Finding
providers, TLS, CSRF, audit, rate limiting, and a real deployed frontend.
Existing `open_case`, saved view, paging, and POST return interfaces are
not modified. See `work/lots/E10-INV-101P-INSPECT/HANDOFF.md`.

## Opt-in native read-only Case Queue browser (E10-INV-101Q-BROWSER)

A host already supplying real authenticated Queue sessions, current Case candidates,
independent Case/Incident/Finding object-read authorization, and optionally a
trusted Shared Saved View provider can mount
`NewAuthorizedQueueBrowserHandler(sessions, cases, access, views)`. It serves
the SAME first-party `GET /case-queue` and `POST /case-queue/return` handler
plus public static `GET /case-queue/ui`, `ui.js`, and `ui.css` assets. Assets
contain **no Case data, source IDs, tenant, bearer or authorization claims**.

The bounded browser displays currently server-authorized Cases, optional literal
search and revocation-safe 25-row cursor paging; a revoked/failed request clears
any previously displayed private rows. All dynamic data are assigned as inert
`textContent`, never executable HTML. Assets are no-store and subject to a
restrictive same-origin CSP, without paid API or external UI libraries.
The host must mount it on the same origin and provide a trusted session mechanism
usable by browser requests (for example a host-managed HttpOnly session cookie);
**no bearer-token injection or client-selected tenant permission** is offered.

This is **only a minimal read-only queue presentation**, not the canonical
Shared Inspector visual component/focus workflow or final Case Queue product UI.
No mutation, assignment, export, Saved View management, production IAM/TLS/audit,
or live Incident integration is claimed. A future Shared Inspector consumer may
use PR #147's `inspect_case` read projection without cloning its component.
The separate engineering manifest and acceptance are at
`work/lots/E10-INV-101Q-BROWSER/manifest.json`. Parent #147 and all earlier
ancestor PRs remain draft/unmerged, and this child must pass its own SHA-matched
CI plus L2/native approvals before merging.
