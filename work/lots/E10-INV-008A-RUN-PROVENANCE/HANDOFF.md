# CMDR — E10-INV-008A-RUN-PROVENANCE

## Product / engineering / code / Git ownership

- Canonical capability: `CAP-INV-008`, [full fiche](../../../cmdr-product-spec/07-investigate/modules/signals-and-hunt/capabilities/search-and-hunt-provenance.md), Investigate / Signals and Hunt. Dependencies: Shared-owned `CAP-INV-002` Search Job and `CAP-INV-004` Telemetry source. This implementation only displays owner-authorized completed Search Job provenance, not a new Trace object.
- Engineering: separately scoped `E10-INV-008A-RUN-PROVENANCE` under `E10-INV-008` epic; parent is **not** CAP-INV-007, E10-INV-120A or a previous DFIR PR. `work/graph.json` status **IMPLEMENTED only**, not `VERIFIED` or canonically accepted.
- Temporary UI code owner: preexisting `product-runtime/investigation-timeline/cmd/case-timeline-demo/analyst.js` and `analyst.html`, because the same first-party completed Search Job is already consumed in its local synthetic investigation workspace. This is **not** a final production architecture.
- Git parent: PR #130 exact frozen SHA `188c0c5970fccb1725d34f0351d1e2575515be93` (which consolidated PRs #128/#129/#130 above PR #127), not `main`. Preserve every old reviewed SHA and native approval gate.

## Canonical obligation → behavior delivered → explicit gap

| CAP-INV-008 requirement | Actual behavior | Outstanding |
|---|---|---|
| Consult Query/run parameters and provenance | From an already completed first-party owner-authorized Search Job, show bounded RunID, QueryVersion, UTC start/end, source **count** and failed-source **count**, and audit-event count | Full Query author identity, transformation events, Data Source inventory, actual Trace links |
| Partial/missing provenance | Every absent or malformed field becomes **indisponible**, with an explicit partial label; never fabricate an actor, Tool Call, Case relation or Trace event | Shared Trace redaction/retrieval, completeness statuses |
| Read-only permission-safe access | Uses the exact Search Job response after existing Case/typed-source server authorization and full result membership checks; no new endpoint, no source dereference, no bearer or source labels printed, no additional HTTP calls | Production IAM, audit-read permission/retention, redacted Trace projection |
| UI expiry/invalidation | Any new Case, bearer, query, field, family, date selection or rejected Job clears the provenance panel and transient result workspace | Persisted versioned history, expired historical run access |

`OPEN-015` remains unresolved regarding audit retention/storage and redaction. `OPEN-013` stays open for writes. This Class-0 stage cannot be called full CAP-INV-008, Shared Trace persistence, signed provenance, Case link or production SIEM.

## Evidence / testing / limitations

Browser test `analyst_summary.test.cjs` exercises an authorized completed Search Job with actual provenance, confirming visible run/version/UTC/source/error/audit counters **without raw source labels, result refs or bearer**; changing query or bearer invalidates the panel. A second test presents missing metadata (explicit partial) then an invalid/foreign result (no partial provenance). Existing CAP-INV-007 result grouping, inspection, return-navigation, authorization and stale selection tests must all continue to pass.

Required exact SHA: `investigation-timeline-runtime` Linux Go + Node browser, macOS Go, Windows Go; `engineering-kernel` mandatory `CHK-SECURITY-TESTS` unchanged **90%** changed security-critical floor, manifest and architecture gates. Parent PR #130 previously showed **89.30% < 90%**, so this candidate must remain DRAFT/unmerged if inherited security control stays red. A browser-only edit does **not** exempt the Go security gate. Do not manufacture tests to manipulate the score or mark the canonical capability accepted.

L2: one independent read-only risk-based review of this material new CAP-INV-008 read-only run-provenance integration, covering the exact Git parent/head diff, permission/metadata leakage, stale invalidation, ownership and missing fields. If the scheduled reviewer/report is inaccessible, mark `L2_EN_ATTENTE`. L2 is distinct from native GitHub approval.

## Next safe work

Reconcile the current failing security coverage scope/root cause and upstream massive draft PR stack; do not stack unbounded UI patches just because the code is testable. Future CAP-INV-008 stages require independently authorized Shared Trace/Audit Hooks, settled OPEN-015 retention/redaction and versioned provenance, with a new bounded architecture ownership decision.
