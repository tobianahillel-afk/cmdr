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
