# E10-INV-120A — live local JSONL to authenticated analyst journey

## Functional gap closed

The existing offline `-import-jsonl/-into-case` pipeline produced a full owner snapshot JSON document that had to be saved and then reloaded into the analyst viewer. This introduced an unnecessary **sensitive intermediate file** and prevented an operator from examining normalized Sysmon/Windows telemetry through the already working analyst UI in one process.

New explicit CLI mode (loopback only):

```sh
go run ./cmd/case-timeline-demo -principal=analyst-a \
  -snapshot=/absolute/path/owner.json \
  -import-jsonl=/absolute/path/events.jsonl \
  -into-case=forensic-case-on-disk -serve-import
```

Existing functions `loadDemoSnapshot`, `readLocalEventFile`, `importLocalEvents` and the same `serveDemoWithData` HTTP handler are reused. Files are validated and the immutable synthetic owner graph revalidated BEFORE token creation, output or listener binding. The browser then uses the **existing** authenticated `/demo-cases`, `/case-timeline`, `/case-source-navigation` and `/demo-source-open` routes with source ACL rechecks. The operator-selected owner file is not modified; no intermediate owner JSON is emitted. Existing offline import/triage/export modes remain unchanged.

## Technical selection

- **Chosen:** compose first-party Go JSONL parser, local memory owner, source ACL and loopback HTML viewer; use the same in-memory dataset for the complete session.
- **Deferred:** new HTTP ingestion route, a filesystem watcher, database, SIEM connector, third-party event bus, external login, user-triggered browser uploads, new persistence abstraction. These would introduce authorities/migrations outside the current verified snapshot boundary.
- **Trade-off:** a local start-up dataset is static. Re-import requires restart; concurrent Case updates, real producer ingestion and forensic immutable custody are not implemented.

## Security and limits

8 MiB per operator-selected local input file, maximum 5,000 Case Timeline entries, normalized JSONL IDs, untrusted event fields narrowed to ID/timestamp/message, no JSONL-supplied tenant, ACL or source role. Existing source permissions originate from the operator-owned demo snapshot and synthetic principal registry; the local bearer is ephemeral. No write API. Denied/foreign owner data remains server-side hidden. This is **not** real production IAM, persisted Case/Shared Timeline service, durable chain of custody, ingestion from a live endpoint, browser-upload support, nor general SIEM connector.

`OPEN-013` remains unresolved; no Shared Timeline Entry mutation occurs. Existing product spec and provider ownership are unchanged.

## Evidence

The new test exercises a **real loopback HTTP listener** from disk import through Case discovery, Timeline, navigation and typed source preview plus anonymous/cross-tenant rejection and cancellation. Invalid inputs, conflicting CLI modes and foreign principal fail with zero stdout (no bearer/listener). Exact GitHub SHA CI and one independent read-only L2 are required before calling this functional increment reviewed.
