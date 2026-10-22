# E10-INV-110B — runnable Case-to-Timeline HTTP demonstration

## Functional milestone and technical selection

The existing `CaseTimelineHTTPHandler`, `NewMemoryCaseTimelineHTTPHandler`, `OpenCaseTimeline`, `EntryResolver`, typed source ACL checks and pagination are composed as a runnable localhost-only server in `product-runtime/investigation-timeline/cmd/case-timeline-demo`.

Reuse-first comparison:
- **Selected:** Go standard-library `net/http`, `crypto/rand`, `crypto/subtle`, the existing Case Timeline handler and immutable in-memory source fixtures. No new storage engine, generic framework or authorization implementation.
- **Deferred:** database or append-only event store, durable ACL/session provider, search backend, cross-service transactions, websocket UI and final Timeline navigation. They are not necessary for the read-only process-local functional acceptance; choosing them now would create architecture commitments with `OPEN-013` unresolved.

## Run from the repository branch

```bash
cd product-runtime/investigation-timeline
go run ./cmd/case-timeline-demo -principal analyst-a
```

The program prints a fresh loopback URL and a **random, one-run bearer token** to the local console. Use them in `curl -H 'Authorization: Bearer <printed-token>' '<printed-URL>'`. A second invocation with `-principal analyst-b` demonstrates the isolated `case-b` tenant.

The synthetic analyst-a Case contains a chronological inferred prior version and observed superseding version. Its third source belongs to a restricted evidence owner and is not disclosed. The browser-compatible GET API can filter by kind/time and paginate with `limit` and `offset`. Requests traverse the existing session resolver, Case authorization, catalog, Timeline Entry authorization, typed source-owner authorization and deterministic projection. Integration tests make real HTTP calls using `httptest.NewServer`.

## Trust boundary and non-production limits

This is **only a localhost laboratory demonstration with invented data**, not a production host. The local CLI selects one fixed synthetic principal at process startup, the HTTP listener binds **only** `127.0.0.1` on an ephemeral port, and each run generates its own unpredictable 256-bit bearer credential. No user-selectable listen address, public network exposure, real customer data, persistent secrets or external providers are supplied. Token output is deliberately printed to the console to enable local testing; it is **not** a production secret-handling design. Hostile local users/processes are outside the demo security guarantee.

Real deployments still require independently verified identities, tenant-bound Case and Shared Timeline Entry catalogs, Evidence/Telemetry/Govern/Command provider integrations, durable audit trails, database lifecycle and transaction/revocation semantics, rate limiting, TLS, operational ownership, and a real interface. Production completeness, SLOs and an accepted L2 result are **not** claimed. `OPEN-013` remains unresolved; B stays READY until legitimate delivery gates hold. No canonical spec or owner boundary is rewritten.

Evidence: inspect exact commit and GitHub Actions jobs before claiming test pass. The canonical Product Spec continues to take precedence.
