# CAP-INV-110 — reuse-first options (E10-SCALE-006B)

Source-backed selection: the `implementation-wave` compiler on `a678f87b0dd828823d93156612f43666c7a8bfef` (push `38833325482`, PR `38833331209`) selected **CAP-INV-110 / Investigation Timeline** among 15 READY capabilities, preferred family `CAP-INV`. Canonical Product Spec: `07-investigate/modules/cases-and-evidence/capabilities/investigation-timeline.md`; screen `INV-TIM-001`. This document records planning alternatives, not benchmark results.

## Alternatives before implementation

| Option | Pros | Risks/cost | Current disposition |
|---|---|---|---|
| **Reuse existing CMDR Case Queue security/snapshot patterns** + caller-owned Shared Timeline Entry envelopes, simple immutable stable timestamp ordering with explicit secondary key | No new infrastructure, preserves tenant/permissions/traceability, smaller attack surface; O(n log n) bounded sort | Not a streaming backend or general timeline store; needs precise time/ordering tests | **Preferred for first bounded read-only projection only** |
| External search/timeline/indexing service or event-sourcing backend | Potential high-scale query and ingest support once requirements quantified | Adds vendor runtime, source ownership, permission, retention and supply-chain concerns; prohibited without Trusted Base review | **Defer** until real query scale, ADR class C, evidence and comparative benchmarks |
| Build new shared Timeline Engine/persistence or incremental ordered-index/tree now | Could support incremental O(log n) updates and complex corrections | Reimplements Shared-owned functionality and guesses unsettled Timeline Entry schema/version semantics | **Reject for first increment**; benchmark bounded sort versus tree if high-scale/streaming requirements emerge |

## Nonnegotiable boundaries

- `Timeline Entry` and shared Timeline Engine remain **Shared**-owned; only Investigate Case selection/projection is new. No runtime-to-runtime import edge unless explicitly owned and allowed.
- Preserve source-declared observed/inferred distinctions, provenance, timezone, original timestamps, versions and corrections; never elevate inference to observed fact or erase superseded evidence.
- Authorization on Case AND every external source is evaluated independently. Denied objects must not be disclosed by title/details even on partial timelines.
- The canonical Product Spec does not settle the final Timeline Entry schema, correction authority, time precision, storage, backend or production SLO. `OPEN-013` prohibits implementing type-2 contribution/correction implicitly.
- Before runtime, execute contract fixtures and negative cases for missing tenant, hidden object, duplicate timestamps, timezone ambiguity, clock drift, version conflicts, stale source and cross-tenant relations; measure deterministic in-process latency then, **not now**.

The planning phase introduces **zero** third-party libraries or new product runtime and does not claim any measured advantage over external alternatives.
