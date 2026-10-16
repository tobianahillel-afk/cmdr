# E10-INV-110B — Offline Case owner snapshot, analyst end-to-end flow

## Why this increment

Prior B work already had an authorization-aware Go Case → Shared Timeline Entry → typed source read-only pipeline, paginated/searchable HTTP endpoint, and analyst viewer, but only hard-coded in-memory fixture data was usable from the CLI. This increment makes that end-to-end flow usable with a **local operator-supplied JSON snapshot** without introducing a new database or a second authorization engine. The canonical CAP-INV-110 schema and decisions remain untouched.

## Reuse-first technical choice

- Selected: existing `MemoryTimelineData` owner snapshot + `NewMemoryCaseTimelineHTTPHandler` and existing same-origin Case picker/viewer; bounded Go `encoding/json` from an explicitly supplied regular file, read once. No new Go module dependency.
- Deferred: SQLite/index/event sourcing, remote repositories, canonical Timeline Entry persistence, ingest pipelines and a novel generic authorization layer. They require owner-side transactional design, migrations, external identity/ACL authority and a separate approved lot.
- Alternatives are not claimed inferior in general; the selected solution is the lowest-risk useful increment for **offline, loopback-only** Case investigation, not production ingestion.

## User journey

1. Prepare a trusted local JSON owner snapshot with principals, Cases, Timeline Entries, source references and explicit reader ACLs. Example: `product-runtime/investigation-timeline/cmd/case-timeline-demo/testdata/example-owners.json` (fully synthetic).
2. Run `go run ./cmd/case-timeline-demo -principal=investigator-42 -dataset=cmd/case-timeline-demo/testdata/example-owners.json` from the Timeline module.
3. Open the printed loopback-only `/analyst` UI and supply the one-run bearer token.
4. `/demo-cases` only lists Cases belonging to the authenticated local principal. `/case-timeline` independently rechecks Case, entry and source permissions before search, chronology, filters and pagination.
5. Denied/cross-tenant sources and untrusted fields are absent from visible results; invalid selectors and unauthenticated requests are rejected.

## Trust and non-goals

Input is an **operator-trusted** local owner snapshot, never an HTTP upload or an externally verified source. Its ACLs are part of the trusted fixture, not an independently authorized production identity provider. Dataset <=16 MiB, regular file, a single JSON object with no unknown fields; the existing memory owner constructor validates cross-links, tenant/environment and ACL reader ownership. Data is copied at startup; there is no persistence writer, transactional Case store, multi-user shared service, hot reload, live revocation, production SSO, audit backend or remote listener. The printed bearer is for local demo only.

No changes to product spec, OPEN-013, boundaries, source owner contracts or the main branch. A stacked PR based on `engineering/autonomous-development-foundation` avoids expanding the #36 diff for reviewers; this work is not independently mergeable to main before its parent.

## Acceptance evidence

Automated tests:
- file → JSON → owner constructor → authenticated Case discovery → filtered source-aware Timeline search → generic cross-tenant denials;
- CLI using the imported file starts a real 127.0.0.1 listener with the imported Case instead of a hard-coded Case, then shuts down cleanly;
- malformed JSON, unknown top-level field, empty/null dataset, trailing content, too-large files and directories are rejected;
- existing `go test -race ./...`, Go vet and engineering-kernel checks still apply.

A published commit/CI and independent L2 report must be verified separately before marking this functional increment accepted. Do not count coverage percentages or historical performance numbers as evidence for imported datasets without exact-run verification.
