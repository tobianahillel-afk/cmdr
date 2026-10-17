# E10-INV-120A — Functional import: normalized local JSONL → authorized Case snapshot

## The new usable workflow

The previous local analyst journey could *read* a complete operator-owned Case/Timeline/source snapshot but could not turn ordinary event rows into that snapshot. This increment adds an **explicit offline preparation command** for normalized security events. It generates a *new*, immutable operator-local owner JSON snapshot consumed by the already-existing Case picker, Timeline API, typed-source preview, browser cross-Case comparison, temporal leads and JSONL export. It does not open a listener or modify either input file.

Example for **synthetic** lab data only, from inside `product-runtime/investigation-timeline`:

```sh
go run ./cmd/case-timeline-demo -principal=analyst-a \
  -snapshot=/absolute/path/owner.json -into-case=case-a \
  -import-jsonl=/absolute/path/normalized-events.jsonl > /absolute/path/merged-owner.json
go run ./cmd/case-timeline-demo -principal=analyst-a \
  -snapshot=/absolute/path/merged-owner.json
```

Input JSONL: one JSON object per line with EXACT `id`, `timestamp` (RFC3339 with explicit offset) and `message` (bounded 512 bytes) and no other fields. A trailing LF is supported. For example:
```jsonl
{"id":"sysmon-1","timestamp":"2026-10-15T10:20:30+02:00","message":"Normalized offline process event"}
{"id":"sysmon-2","timestamp":"2026-10-15T10:20:31Z","message":"Normalized offline connection event"}
```

The operator-selected owner snapshot provides **the Case, principal, tenant, environment and reader ACLs**. A JSONL row has NO right to select or override access. Import records receive observed kind, synthetic `operator-local-jsonl` provenance, deterministically SHA-256-derived IDs (Case + event ID), typed `telemetry-event` sources and a restricted synthetic preview. Both Timeline Entry and typed source readers are limited to the selected Case-authorized principal. The existing owner provider validates the *entire* graph, Case/Entry/source identities and ACLs **before** any stdout output. This is not authoritative ingestion from an EDR and not cryptographic evidence custody; SHA-256 here is used only for stable names, not signing or integrity proofs.

## Reuse versus alternatives

**Chosen:** Go `encoding/json`, `crypto/sha256`, `time.Parse` and the existing `MemoryTimelineData` plus `NewMemoryCaseTimelineHTTPHandlers` owner validation. This is a deterministic one-time transformation requiring no migration and no external dependency.

**Not chosen now:** a new database/index, vendor-specific Sysmon/Defender/Wazuh SDK, streaming ingestion pipeline, proprietary IAM, or new Case ACL framework. None is necessary to make a normalized local event import exercise the existing read-only Case → Timeline → source workflow; the production connector and durable store choices belong to separate roadmap work, not an implicit change in this lab lot.

## Security and limits

- Host CLI only; **never accepts a browser-selected path**, and creates no new HTTP endpoint.
- Both file paths must be absolute, regular, non-symlink; before/after open identity is checked. 8 MiB per input file, strict one-object-per-line JSON, no unknown fields, up to 5,000 total entries in the selected Case, 4 KiB per input event line. Output must also fit the existing 8 MiB snapshot limit.
- Rejects timestamp ambiguity, duplicate event IDs, unstable identifiers, inconsistent pre-existing owner links or source permission claims; validation fails before emitting bytes.
- Produces a **full owner snapshot that can include other tenants and restricted source previews from the operator's original input file**. The output is SENSITIVE; direct only to a protected local file, never to logs, chat, CI artifacts, repositories or a public stream. Never use real customer data in this demo. File redirection permissions are operator responsibility.
- Real-world authentication, remotely verified per-source ACL, durable storage/transactions, live revocation, chain-of-custody signatures, arbitrary vendor mappings and production-ready multi-user IAM are **still absent**. `OPEN-013` remains unresolved.
- This is a materially new functional input path, not proof that all of CAP-INV-110 is deployed.

## Functional acceptance

Behavioral Go tests build a local JSONL input plus a distinct disk owner JSON, call the CLI, load the emitted snapshot back from disk, then issue real authenticated HTTP requests against the reused in-memory owner handlers. They check visible timeline events, source preview content imported only from JSONL, anonymous rejection, cross-tenant isolation and denial of existing restricted entries. Invalid rows, unauthorized Case selection, forged access fields and malformed filesystem paths must fail before output. Exact-SHA CI and independent L2 review are required; do not claim they passed merely from this documentation.
