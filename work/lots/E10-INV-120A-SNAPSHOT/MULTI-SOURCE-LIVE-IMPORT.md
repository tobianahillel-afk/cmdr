# E10-INV-120A — multi-file offline telemetry to one authorized Case

## Functional increment and acceptance

The local analyst can open up to **four independent normalized JSONL telemetry files** in an existing Case, with no intermediate owner-snapshot export or network ingestion endpoint:

```sh
go run ./cmd/case-timeline-demo \
  -principal=analyst-a \
  -snapshot=/absolute/local/owner.json \
  -into-case=forensic-case-on-disk \
  -import-jsonl=/absolute/local/sysmon.jsonl \
  -import-jsonl=/absolute/local/windows-security.jsonl \
  -serve-import
```

Only `-serve-import` allows repeated `-import-jsonl`. Offline `-import-jsonl` transformation and `-triage-import` remain **exactly one input** for compatibility. The imported Case is selected as the analyst viewer's initial Case. The existing authenticated `/demo-cases`, `/case-timeline`, `/case-source-navigation`, `/demo-source-open` handlers apply per-Case, per-entry and per-source synthetic owner ACL. Neither the browser nor a HTTP parameter supplies the event file paths, principal role, tenant, or authorization.

Inputs are checked at process startup before bearer generation and listener binding. Reuse `readLocalEventFile` for the absolute, regular, no-symlink, 8 MiB **per-file** boundary and `importLocalEvents` for normalized fields, immutable snapshot, canonical owner graph validation, unique source identities and 5,000 total Case entries. Repeated paths and hardlinked files are rejected. If any input is invalid, or its IDs collide with earlier files, the whole load fails with no partial listener or sensitive stdout. Operator's original owner snapshot and JSONL files are not written.

## Reuse-first decision

- **Selected:** sequentially compose the existing bounded Go JSONL importer and the same memory owner/ACL HTTP route. New code is one bounded path-list flag and one multi-file orchestration function; no generic ETL adapter, parser rewrite, storage abstraction, event bus or watcher.
- **Deferred:** on-demand browser uploads, HTTP ingestion, persistent multi-source index, production authorization, custody signing, background file monitoring and live Sysmon/Defender connector. Those require ownership and trust decisions not settled by the Case snapshot task or `OPEN-013`.
- **Limit:** maximum four files, 8 MiB each and 5,000 Case entries. This is still operator-controlled offline, synthetic principal/ACL and read-only loopback. A process restart is required to load changes; no production IAM, real retention/correlation store, custody certification or production SLO.

## Evidence expectations

`TestMultipleLocalJSONLFilesReachAuthorizedAnalystViewer` checks real localhost startup, bearer, Case page, two independent imported entry identities and typed source previews, denied anonymous access, foreign Case exclusion and unchanged owner file. `TestMultipleLocalJSONLFilesFailClosedBeforeBearerOrListener` checks duplicate IDs, bad later input, repeated path, owner path, hardlink, more than four files, unauthorized Case and rejection of multi-file offline mode. Exact SHA CI must be verified after publication; green CI alone is not an independent L2 verdict. One risk-based read-only L2 is appropriate for this **new multi-source functional path**, not for subsequent formatting/test-only commits.
