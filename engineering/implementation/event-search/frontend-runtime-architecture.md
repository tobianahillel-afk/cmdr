# Event Search frontend runtime architecture

Status: **accepted engineering decision**  
Decision: `ENG-DEC-0001`  
Validation: `DECVAL-0001`  
Research: `RES-PKT-0001`  
Benchmark: `BENCH-0001`

## Selected family

Event Search will use a **first-party browser Web Platform client**. The implementation may use build-time development tooling, but it does not introduce a third-party UI framework/runtime into the CMDR product runtime by default.

The choice is an engineering architecture decision only. It does not change Product Spec behavior.

## Runtime rules

- First-party ES modules and standard DOM APIs.
- Native semantic HTML controls by default.
- Custom Elements only where a stable component boundary materially helps; they are not a replacement for native semantics.
- No third-party component kit.
- No framework-owned global state container.
- No unsafe rendering of untrusted values through `innerHTML`, `outerHTML`, `insertAdjacentHTML` or equivalent string-to-DOM sinks.
- CSP is mandatory. Trusted Types is defense in depth where supported, not the sole XSS control.
- Server-side authorization remains authoritative; frontend permission checks only shape UX.
- Request cancellation uses standards-based cancellation such as `AbortController`; the architecture does not choose a final streaming transport.
- Third-party product-runtime dependencies remain deny-by-default and require a new source-backed engineering decision.

## State model

The frontend must keep one explicit Event Search state model with a strict split:

**URL/deep-link state**
- query text or opaque backend-neutral query representation;
- search scope;
- filters;
- stable navigation/return context that is product-defined.

**transient client state**
- loading and cancellation handles;
- in-flight request identity;
- partial-result assembly state;
- focus/selection state that is not product-defined as deep-linkable;
- view-window/virtualization bookkeeping.

URL/deep-link state is serialized through URL/History primitives and must survive refresh. Hidden framework state must never be required to reconstruct a deep link.

## Rendering model

- Initial rendering may use semantic server-provided shell markup where useful, but Event Search interactive state is owned by the first-party client model.
- Frequent result/state updates use bounded incremental DOM changes rather than unconditional full-fragment replacement.
- Large-result rendering must remain windowed/virtualized once the concrete UI lot defines the result-table geometry; this decision does not select the final virtualization algorithm.
- Loading, Empty, Partial, Error, Offline and Permission denied are explicit render states.

## Benchmark interpretation

`BENCH-0001` compared the selected candidate with a first-party server-progressive alternative on 250 rendered result rows.

Across seven repeated headless Chromium runs:
- selected Native Web median p95 initial render: **12.5 ms**;
- server-progressive median p95 initial fragment: **9.5 ms**;
- selected Native Web median p95 repeated update: **0.7 ms**;
- server-progressive full-fragment update: **9.8 ms**.

The alternative therefore retained an initial-render advantage, while the selected architecture had roughly fourteen-times lower repeated-update p95 in the measured interactive workload. These figures are engineering decision evidence, not production Search SLOs.

## Explicit non-decisions

This decision does **not** choose or imply:
- a final query language or dialect;
- search index technology;
- storage engine;
- search provider;
- final streaming transport;
- persisted result-storage model;
- Saved Search behavior;
- Case-link mutation while `OPEN-013` remains unresolved.

## Revisit conditions

The decision is reusable only while `FRESH-0001` remains fresh. Automatic repository/spec/evidence drift, material browser/runtime changes, security advisories, threat-model changes, benchmark regressions or material new research force revalidation.
