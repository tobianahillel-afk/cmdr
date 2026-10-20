# Event Search frontend runtime

This directory is reserved for the first-party browser runtime of `CAP-INV-002 / INV-EVS-001`.

Current state: **bounded state/deep-link core, accessible semantic shell, bounded incremental result renderer and backend-neutral transport integration implemented and verified; adversarial end-to-end frontend validation remains next**.

The accepted engineering architecture is `ENG-DEC-0001`: first-party Web Platform ES modules, native semantic HTML by default, bounded incremental DOM updates, URL/History deep-link state and no third-party UI framework/runtime by default.

The boundary-registration, state, shell, result-rendering and transport-integration lots are complete. The executable frontend surface includes first-party deep-link state, semantic shell states, keyed bounded incremental result rendering with position preservation, backend-neutral execution/retry transport, AbortController cancellation, permission-denied fail-closed handling, partial-result preservation and retry-as-new-run semantics. No final streaming transport, Saved Search behavior or Case-link mutation is implemented yet.

Explicitly unresolved/out of scope here:
- Case-link mutation while `OPEN-013` remains open;
- final query dialect;
- search index or storage engine;
- search provider;
- Saved Search implementation;
- persisted result model;
- production Search SLO.

## Investigate Event Search -> Event Inspection read-only journey (E10-INV-004F-JOURNEY)

`inspection.mjs` composes the **existing** Event Search controller with a separate `transport.inspect(request, {signal})` supplied by a host. Its input must be an already server-authorized CAP-INV-004 Event Inspection projection, **not** a browser-computed access decision. The journey opens only a selected result from the currently visible authorized Search Job, binds Tenant/environment/Event/job/run/correlation, rejects a mismatched response and cancels or discards stale inspections. It displays source/normalized/raw-derived/enrichment provenance and raw/rendered separately with native `textContent` sinks, then returns to the exact prior Search state and result window if the source Job remains current.

The browser module **does not** execute the Go inspector core, provide a live Event resolver, authenticate the principal, classify sensitive normalized/enrichment fields, select a search engine, run a PivotDraft, or perform Case/Artifact/Evidence mutations. The Go `product-runtime/event-inspection` core and Shared/Settings authorities are unchanged. Before a real deployment, the backend host must independently reauthorize each Event, tenant, environment, raw/rendered classification, source and requested field, and return a safe projection. This is a tested first-party frontend composition, **not** complete source-to-UI production integration.

Run the existing full frontend suite with `npm test` (Node 24); the isolated exact-HEAD workflow is `.github/workflows/event-search-inspector-journey.yml`. The cross-capability manifest is `work/lots/E10-INV-004F-JOURNEY/manifest.json` and the original engineering foundation remains PR #36 DRAFT.

## Live local browser -> Go HTTP handshake (E10-INV-004H-HTTPBRIDGE)

`createInspectionHTTPTransport()` in `inspection-http.mjs` binds the previously backend-neutral Journey to the existing first-party Go POST `/investigate/event-inspection`, with a fixed same-origin endpoint, POST JSON, explicit cache/redirect controls, AbortSignal and bounded response. `inspection-http.test.mjs` executes the actual HTTP request to a Go host listening only on loopback with independently scoped synthetic Search/Telemetry owner interfaces and a per-test bearer. It verifies authorized Search row -> server permission -> Inspection -> same Search result window as well as wrong bearer/tenant/raw denial/oversize. `npm test` runs it together with all previously published first-party frontend regressions.

The Go host and JS frontend are actually connected **only in this synthetic end-to-end test**. A deployment must supply real authentication, both independent owner adapters, source classification, CSRF protections, TLS, audit sink and actual browser hosting; none are implemented by the fixture or proven in production. Avoid treating this bridge as a generic URL-fetcher or an IAM framework.
