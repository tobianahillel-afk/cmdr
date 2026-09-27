# Event Inspection runtime boundary

This directory contains the first bounded executable runtime of **CAP-INV-004 Event Inspection and Pivot**.

## Current state

**Implemented read-only inspection core.** The runtime is a Go standard-library module with zero external runtime dependencies. It performs immutable, single-Tenant Event inspection projection only.

Implemented in E10-INV-004B-RUNTIME:

- explicit server-side telemetry-event.read enforcement;
- independent raw/rendered access decisions;
- zero raw leakage when raw is denied, including raw-derived sensitive values;
- copied source and normalized fields with source/parser/missing-field provenance kept distinct;
- copied derived enrichments with producer/version/freshness and stale markers;
- stable source-unavailable, tombstone, partial, and restricted state projection;
- correlation and Event identity preservation;
- fail-closed mutation rejection while OPEN-013 and OPEN-014 remain unresolved.

## Deferred to later bounded tasks

E10-INV-004C-PIVOT adds deterministic dialect-neutral PivotDraft creation and return-context handling. Pivot preparation validates that the authorized field/value is actually present in the visible projection, preserves source provenance and normalizes the time window without executing Event Search. This runtime does not select or implement a final query dialect or provider.

## Explicit exclusions

This boundary does not select or implement:

- Case-link or annotation mutation while OPEN-013 remains open;
- Artifact/Evidence-candidate mutation or retention semantics while OPEN-013/OPEN-014 remain open;
- parser execution;
- Entity Resolution;
- a final query language or dialect;
- search/index/storage/provider technology;
- automatic Evidence or Finding creation;
- a dedicated final Event Inspector UI;
- any production retrieval/search SLO.

External product-runtime dependencies remain deny-by-default.
