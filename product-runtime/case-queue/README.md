# Case Queue runtime boundary

This directory is reserved for the first bounded executable runtime of **CAP-INV-101 Case Queue**.

## Current state

**Preimplementation only.** This README is the sole file intentionally present in the runtime boundary at this stage. It carries no executable product behavior and must not be treated as runtime coverage, authorization-negative evidence, tenant-isolation evidence, performance proof, persistence, a final Case Queue UI, or a Case lifecycle implementation.

## Bounded future scope

Later verified tasks may implement only the read-only queue core already materialized in `WAVE-CAP-INV-101-CORE-001`:

- validation of caller-owned, explicitly tenant-scoped Case snapshots and access decisions;
- immutable projection of accessible Case queue rows without redefining the canonical Case object;
- deterministic local search, filter and stable sort over the bounded projection;
- optional Incident and Finding context with explicit freshness, partial and unavailable diagnostics;
- permission-re-evaluated application of a caller-owned Shared Saved View without mutating Shared state or exposing forbidden fields/filters;
- deterministic available/empty/partial/stale/permission-filtered/view-dirty states;
- immutable Case Workspace handoff and safe return-context restoration for view, filters, scroll and selection.

## Explicit exclusions

This boundary does not select or implement:

- Case create/update/assignment/status/lifecycle mutation, which remains CAP-INV-102 territory and is governed by OPEN-013;
- Export Job creation or export backend selection;
- Saved View create/update/share/archive/manage mutation;
- Task, Decision or Response Run ownership, aggregation or Command Work Queue behavior;
- canonical Case schema or final Case state-machine changes;
- final Case Queue columns or a dedicated final CAP-INV-101 UI;
- storage/provider/retention/collaboration backend selection;
- a production Case Queue SLO or full-capability completion claim.

Case remains Investigate-owned, Incident remains Command-owned and Saved View remains Shared-owned. External product-runtime dependencies remain deny-by-default. The initial implementation direction is Go standard library first unless a later evidence-backed technical decision proves otherwise.
