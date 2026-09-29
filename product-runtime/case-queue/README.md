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

The dedicated Saved View task `E10-INV-101C-VIEW` remains required to finalize the isolated adapter/fallback behavior; this runtime does not claim Shared Saved View persistence or management.

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
