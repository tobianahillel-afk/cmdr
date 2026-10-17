# E10-INV-120A — one-command JSONL import → authenticated analyst viewer

## Functional journey

Existing separate paths required exporting a complete merged owner snapshot (which may include other tenants and restricted source previews) before starting the analyst viewer. This integration avoids writing or emitting that sensitive intermediate owner graph:

```sh
go run ./cmd/case-timeline-demo -principal=analyst-a \
  -snapshot=/absolute/path/owner.json \
  -import-jsonl=/absolute/path/normalized-events.jsonl \
  -into-case=case-a -serve-import
```

The operator supplies two distinct absolute local files, and an already authorized Case. The command loads the existing trusted local owner snapshot, parses strictly bounded normalized JSONL, validates the imported owner graph, then attaches the already available Case picker, timeline pages, source preview, source navigation and analyst JavaScript UI to **127.0.0.1 on an ephemeral port**. The terminal displays a new ephemeral bearer; the browser must send it explicitly. No merged owner JSON is written to disk or sent to stdout. Original input files remain unchanged.

## Reuse-first architecture and risk boundary

**Selected:** existing `loadDemoSnapshot`, `readLocalEventFile`, `importLocalEvents`, `serveDemoWithData` (thin refactor of `serveDemoWithSnapshot`), and `newDemoHandlerWithData`. Go standard library only.

**Deferred:** new listener, import API endpoint, importer framework, database, filesystem staging area, or ACL implementation. Each would duplicate trusted Case/Timeline/typed source boundaries and increase the disclosure surface. Import is explicit CLI-only and cannot be combined with triage/export/comparison modes.

Security properties: reject invalid/oversized/symlink input, malformed or access-forging JSONL, unauthorized Case, cross-tenant sources, and cancelled setup before the server advertises a bearer or binds. The existing HTTP handler reauthorizes every Case and source read; public browser routes never accept raw files or client-provided authority. The bootstrap screen and Case picker remain synthetic.

Limitations: operator-controlled synthetic owner ACLs, local-only single-process in-memory data, no real SIEM/EDR connector, no production IAM/DB/audit or durable cryptographic chain of custody. The command does not save an evidence copy, and the normalized input is NOT authenticated forensic evidence. `OPEN-013` remains unresolved. No claim of production readiness.

Acceptance: a disk snapshot plus two JSONL rows starts a real loopback server, returns an imported entry through bearer-authenticated Case paging and typed-source preview, hides the foreign Case, denies anonymous access, preserves both input files byte-for-byte, shuts down on cancellation, and returns no token for malformed/unauthorized input. Confirm CI by exact GitHub SHA.
