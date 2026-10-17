# E10-INV-120A — Wazuh 4.x alerts.json → authorized Case analyst viewer

## Functional operator path

Use a **locally obtained Wazuh 4.x alerts.json** file (newline-separated JSON), an existing operator-controlled read-only Case owner snapshot and an authorized synthetic principal:

```sh
go run ./cmd/case-timeline-demo \
  -principal=analyst-a \
  -snapshot=/absolute/path/owner.json \
  -wazuh-alerts=/absolute/path/alerts.json \
  -into-case=case-a -serve-import
```

One command now reads the real-shaped source alert format, converts it in memory into the already supported bounded normalized JSONL contract, imports those events into the already-authorized Case, and serves the existing bearer-protected loopback Case picker, Timeline, navigation and typed-source previews. No merged owner dataset is written to disk. Neither input is modified. The browser cannot supply a source file, Case owner or permission.

Wazuh 4.x source format is evidenced by its documented `timestamp`, `id`, `rule.description` and `full_log` fields; for example https://documentation.wazuh.com/current/user-manual/ruleset/ruleset-xml-syntax/rules.html. This adapter intentionally does **not** accept Wazuh 5.0 ECS `@timestamp` records; that version uses a different source contract and must be handled separately.

## Normalization and safety

- The canonical Wazuh alert `id` is preserved as `wazuh-<id>`; the importer's existing deterministic per-Case SHA-256 entry/source IDs prevent collisions with unrelated existing owner entries.
- Source timestamp offsets such as `2026-10-17T09:30:15.725+0200` become RFC3339Nano with `+02:00`, without changing the instant.
- The analyst-facing preview uses the rule description, falling back to `full_log` only when absent. This is **a summary**, not a verbatim forensic copy of the raw alert. Unsupported control characters, excessive previews and unknown timestamps are rejected, never silently truncated.
- Optional Wazuh metadata (agent, manager, location, original raw log) is **not** carried into the normalized event. Raw event origin, signature, full payload and complete evidence chain are not retained; do not claim an authenticated forensic evidence record.
- Reuse existing operator-only absolute path/no-symlink/8 MiB file bounds, 5,000 entries, duplicate identity checks, authorized Case selection and server-owned read/preview checks. Conversion is in-memory; all-or-nothing before binding a listener or exposing its bearer.
- The source alert is untrusted. Its fields are **not** a tenant label or access grant. All Case/Timeline/source access still derives solely from the operator's local memory owner snapshot. Real production IAM, external Wazuh API, ingestion queue, database/audit and Shared Timeline Engine remain absent; `OPEN-013` stays unresolved.

## Reuse-first decision and acceptance

Selected: Go `encoding/json`, `time`, existing `readLocalEventFile`, `importLocalEvents`, `serveDemoWithData` and existing HTTP handlers. Rejected for this bounded offline integration: new SIEM SDK, queue, index, parser framework, remote credentials or duplicated access-control engine. This is a functional bridge, **not production SIEM ingestion**.

Tests must exercise real read-only disk input, Wazuh 4.x timestamp normalization, source preview via bearer-authenticated HTTP, cross-tenant isolation, anonymous denial, invalid/oversized inputs, missing permission and cancellation. CI and independent risk-based L2 must be checked at exact published SHA, not inferred from historical green runs.
