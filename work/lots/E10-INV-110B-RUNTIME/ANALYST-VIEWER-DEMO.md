# E10-INV-110B — Interactive read-only Case investigation demo

## Delivered functional path

One Go standard-library executable now connects a loopback browser page to the already implemented Case → Timeline HTTP handler, immutable owner-backed sample Case/Timeline/source data, server-derived principal and tenant, Case catalog, Timeline Entry ACL, typed source ACL, stable projection, filters and authorized pagination.

Manual quickstart (no external dependency or paid API):

```sh
cd product-runtime/investigation-timeline
go run ./cmd/case-timeline-demo -principal=analyst-a
```

Open the **Local analyst viewer** URL shown in the terminal, enter the separately displayed ephemeral bearer (the UI does not receive it automatically), and choose `case-a`. A second instance using `-principal=analyst-b` allows `case-b` but cannot retrieve `case-a`. Filter `observed` vs `inferred`, page authorized results and inspect provenance and accessible supersession. There are no persistence writes, external calls, database changes or browser credentials in URLs.

## Reuse and alternatives

- **Selected:** `net/http`, `embed`, the existing `NewMemoryCaseTimelineHTTPHandler`, `OpenCaseTimeline`, and its authenticated ACL pipeline. A handful of static HTML/JS/CSS files, with no new server abstraction or framework.
- **Deferred:** React/SPA, a separate API gateway, remote IdP, third-party database/index and a new Timeline Engine. None is needed to deliver the bounded interactive reference flow; those are separate product integrations requiring ownership and production decisions.
- **Boundary:** UI markup is untrusted presentation. Every refresh, filter and page hits the same server-owned Case/Timeline/source permission chain. No client-supplied tenant or role flags. Never use this demo memory provider or ephemeral console token as production identity.

## Acceptance and limitations

Tests check HTTP delivery/security headers, no anonymous data access, tenant/source isolation on the existing API, denied HTTP methods, and static-client practices (textContent, no local storage, same-origin credentials omit). The existing live loopback tests verify startup and shutdown. This addition does not constitute a full browser automation suite, production authorization backend, durable storage or independent product navigation lot C; graph B stays READY until acceptance gates are met.

The official performance evidence from older SHA dae5fe2c remains a synthetic historic measurement; no production p95 or RSS is inferred. No merge, Product Spec modification or new architecture decision is authorized by this step.

## Completed next-step integration: analyst time windows

The local synthetic analyst viewer now exposes **From/To** local-time controls for the existing server-side `[from,to)` selection, converts those browser-local instants to RFC3339 UTC, and combines time + kind + Case + authorized pagination. Uncertain chronology remains visible and marked uncertain according to the existing server contract. The browser rejects reversed intervals before sending a request; the server independently rejects malformed/unauthorized selections. Cancellation plus monotonically increasing request revisions prevents an older authorized response from overwriting a changed Case, token or filter selection.

Functional integration tests exercise inclusive lower/exclusive upper boundaries, offset equivalence, kind conjunction, empty windows, malformed intervals and cross-tenant/anonymous denial through the **real in-memory HTTP handler**. No client-supplied permission flag, additional provider, external API or database is introduced. Browser component automation and production identity/provider connections remain separate milestones. An existing GitHub engineering-kernel push check was failing at 82.22% security-critical coverage on the previous demo SHA; this increment does not claim to have repaired that coverage threshold without exact-SHA CI evidence.

## Local authenticated Case chooser

An analyst may now enter the ephemeral terminal bearer, click **Afficher mes Cases**, choose an authorized synthetic Case and then submit the same timeline query with time/kind filters. The dedicated `/demo-cases` endpoint authenticates on every request, derives the principal on the server, lists only fixture Case IDs with matching tenant/environment and reader permission, and returns no entry/source metadata. It rejects callers attempting to supply a tenant, principal or Case selector. The browser clears stale catalogue results on bearer change and uses inert DOM text, same-origin, no credential persistence, and abortable fetches.

Integration tests use the **real loopback HTTP handler** from authenticated catalogue discovery through to ACL-protected Timeline opening for both synthetic tenants; invalid authentication and forbidden Case paths are negative scenarios. Unlike a remote connector or durable Case repository, this only makes the synthetic analyst demo more navigable. The host's actual owner-backed catalogue/authentication/persistence remains an outstanding functional integration and security gate; do not treat the demo as a production end-to-end CMDR deployment.
