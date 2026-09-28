# Saved Query Assets runtime boundary

This directory is reserved for the first bounded executable runtime of **CAP-INV-006 Saved Searches and Query Assets**.

## Current state

**Preimplementation only.** This README is the sole file intentionally present in the runtime boundary at this stage. It carries no executable product behavior and must not be treated as runtime coverage, authorization-negative evidence, tenant-isolation evidence, performance proof, persistence, a canonical Saved Search schema, or a Search Job implementation.

## Bounded future scope

Later verified tasks may implement only the read-only compatibility and handoff core already materialized in `WAVE-CAP-INV-006-CORE-001`:

- validation of one caller-owned, explicitly tenant-scoped Saved Search / Query Asset snapshot;
- immutable projection of the stable Shared Query reference and exact Query version;
- copied parameters/variables, source and field prerequisites, author, validation and existing lineage/deprecation metadata;
- deterministic compatibility, stale, partial and unavailable-prerequisite diagnostics without synthesizing missing facts;
- fail-closed access decisions and cross-tenant rejection;
- construction of an immutable backend-neutral Event Search handoff draft that preserves the original Query/version and requires Event Search to re-evaluate source and permission context.

## Explicit exclusions

This boundary does not select or implement:

- Saved Search or Query Asset create/save/share/duplicate/deprecate/archive mutation while OPEN-013 remains unresolved;
- Shared Query mutation or ownership transfer;
- Saved View mutation or conceptual merging with a Saved Search;
- Detection Rule creation, conversion or deployment;
- Search Job creation or execution inside this runtime;
- cross-tenant copy;
- storage/version-store/provider/retention implementation;
- collaboration backend, approval or revalidation policy;
- a final query dialect;
- a dedicated final CAP-INV-006 UI;
- a production Saved Search/Event Search SLO or full-capability completion claim.

External product-runtime dependencies remain deny-by-default. The initial implementation direction is Go standard library first unless a later evidence-backed technical decision proves otherwise.
