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
