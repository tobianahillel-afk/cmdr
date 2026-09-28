# Saved Query Assets runtime boundary

This directory contains the first bounded executable read-only runtime of **CAP-INV-006 Saved Searches and Query Assets**.

## Implemented scope

The runtime validates and immutably projects one caller-owned Saved Search / Query Asset snapshot only:

- one explicit authorized Tenant and environment;
- stable asset identity and kind (`saved-search` or `query-asset`);
- fail-closed asset and Shared Query access decisions;
- the original Shared Query reference, tenant and exact Query version;
- copied parameter names;
- copied source and field prerequisites with their caller-resolved state;
- copied author and validation state;
- existing lineage, deprecation reason and replacement-asset references when supplied;
- deterministic sorted compatibility diagnostics;
- handoff eligibility only when all prerequisites are compatible;
- deep-copy isolation from caller-owned slices and nested source/field data.

Missing or removed prerequisites are reported as incompatible. Stale validation/source/field facts remain visible and make the asset ineligible for execution handoff. Missing facts are never synthesized.

The module uses the Go standard library only and performs no network, filesystem, persistence or query-execution work.

## Explicit exclusions

This runtime does **not** implement or select:

- Saved Search / Query Asset create, save, update, share, duplicate, deprecate or archive mutation while OPEN-013 remains unresolved;
- Shared Query mutation or ownership transfer;
- Saved View mutation or conceptual merging with Saved Search / Query Asset;
- Detection Rule creation, conversion or deployment;
- Search Job creation or execution;
- Event Search execution handoff construction (reserved for E10-INV-006C-HANDOFF);
- cross-tenant copy;
- storage/version-store/provider/retention;
- collaboration backend, approval or revalidation policy;
- a final query dialect;
- a dedicated final CAP-INV-006 UI;
- production Saved Search/Event Search SLOs or full-capability completion.

External product-runtime dependencies remain deny-by-default.
