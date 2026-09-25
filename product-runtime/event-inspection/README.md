# Event Inspection runtime boundary

This directory is reserved for the first bounded executable runtime of **CAP-INV-004 Event Inspection and Pivot**.

## Current state

**Preimplementation only.** This README is the sole file intentionally present in the runtime boundary at this stage. It carries no executable product behavior and must not be treated as runtime coverage, authorization-negative evidence, tenant-isolation evidence, or performance proof.

## Bounded future scope

Later verified tasks may implement only the read-only core already materialized in `WAVE-CAP-INV-004-CORE-001`:

- immutable single-Tenant Telemetry Event inspection projection;
- explicit server decisions for event/read/raw/rendered access;
- strict raw/rendered separation, including zero raw leakage when raw is denied;
- source/parser/timestamp/missing-field visibility;
- derived enrichment producer/version/freshness distinction;
- source-unavailable, tombstone, partial and restricted states;
- dialect-neutral PivotDraft preparation with field/value/time/source/return context.

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

External product-runtime dependencies remain deny-by-default. The initial implementation direction is Go standard library first unless a later evidence-backed technical decision proves otherwise.
