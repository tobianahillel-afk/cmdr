# E10-INV-120A-SNAPSHOT — operator principal and complete offline Case flow

## Working integration

The existing loopback demo previously accepted only `analyst-a`/`analyst-b`, even when an operator-selected snapshot defined its own local investigator principals. This increment **reuses** `loadDemoSnapshot`, `firstDemoAuthorizedCase`, `NewMemoryCaseTimelineHTTPHandlers` and all existing Case/Timeline/typed-source authorization. No new provider, DB, framework, object schema or permission algorithm was added.

The bundled `testdata/example-owners.json` adapts the reusable two-tenant synthetic owner snapshot from draft PR #40 and adds harmless source preview strings. After cloning the repository, from the `product-runtime/investigation-timeline` directory:

```sh
go run ./cmd/case-timeline-demo -snapshot="$(pwd)/cmd/case-timeline-demo/testdata/example-owners.json" -principal=investigator-42
```

Open the printed loopback analyst URL, enter the ephemeral bearer printed in **that** terminal, discover `case-incident-42`, view the chronological Timeline and open the permitted source preview. `investigator-99` sees only the other tenant's Case when selected in a **separate instance**. This is synthetic offline analysis, **not real customer evidence**.

The operator-selected snapshot determines the host's principal set and per-object reader ACLs; the principal is never taken from an HTTP parameter. Before any listener opens, the selected principal must be in the owner snapshot and have a readable Case. Each HTTP call independently checks the bearer, Case, Timeline Entry and typed source, including source preview. Unknown principals and foreign Cases fail closed. Without `-snapshot`, the original embedded demo accepts only `analyst-a` / `analyst-b`. The startup Case link URL-encodes the owner Case ID.

## Verification and boundaries

The new end-to-end test drives the **actual ephemeral loopback listener**, checking:
- synthetic imported investigator identity and first Case selection;
- authorized Case discovery, ordered Timeline entries, source navigation and source preview;
- denial of anonymous requests, foreign Case access, and foreign source metadata;
- independent other-tenant Case discovery and denial of cross-tenant Case;
- shutdown and default embedded-principal isolation.

The existing tests additionally reject malformed/oversized/symlinked snapshots, invalid owner memberships and forged `SourcePermission`. Tests and CI must be checked on the exact commit: publishing this file is not a verification result.

**Product exclusions:** this is a local demo with operator-controlled ACL snapshot and temporary bearer; it does NOT implement production authentication, durable persistence, live SIEM ingestion, multi-user transactional ACL updates, real audit storage, network-exposed API, or completion of CAP-INV-110. It must not be given real customer secrets or publicly exposed. `OPEN-013` remains open.

## PR dependency and choice

Build a distinct scoped functional PR against `engineering/autonomous-development-foundation` (PR #36), instead of appending further micro-fixes to a large draft branch. PRs #37/#38/#39/#40 were reviewed for reuse, not merged or overwritten. In particular, the synthetic example schema is reused from PR #40 without importing its separate loader/authorization implementation. The new PR should not be merged before its base and required gates are accepted.
