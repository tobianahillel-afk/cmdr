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
