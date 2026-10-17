# E10-INV-120A — integrate authorized offline Case dossier on current foundation

## Functional gap and selected approach

PR #45 implements useful owner-authorized Case -> Timeline -> typed-source dossier behavior but its historical base diverges by 10 commits from engineering foundation and its adaptive kernel CI fails (security-critical changed coverage 88.83%). Rebase by building a clean, narrow functional integration PR from the exact current engineering HEAD instead of force-pushing or silently overwriting changes in PR #45. Reuse the exact first-party dossier implementation and behavioral on-disk tests, plus the shared source-inspection parser. The current engineering CLI contains working JSONL import/triage and authorized JSONL export that the historical PR did not know about; preserve those modes, and reject combining them with dossier flags.

The dossier collects the existing complete authorized Case pages, selected Timeline Entry and same-source chronology, rechecks the original typed-source owner endpoint, and emits a *single* sensitive JSON report only after success. User-visible capability is an offline, explicit command: `case-timeline-demo -snapshot=/absolute/path/owners.json -principal=analyst-a -dossier-case=case-a -dossier-entry=event-old`.

Alternatives: (1) **selected** reuse existing first-party Go handlers and `encoding/json`; (2) merge an old conflicting PR wholesale (would overwrite importer/exporter), rejected; (3) build a new DB/authorization engine or independent source parser, rejected as unnecessary and higher security risk. No new dependencies or paid API.

## Acceptance and security limits

- On-disk snapshot → authenticated Case pages → Timeline → independently authorized typed source preview and same-source chronology through existing bounded first-party handlers.
- Denied/cross-tenant or missing source, malformed input, unsupported mode combinations and cancellation produce no JSON report.
- Preserve current `-import-jsonl`, `-triage-import`, `-export-case` flows and check mutual exclusion with dossier mode.
- Run exact-SHA `investigation-timeline-runtime` Go race/vet and `engineering-kernel` PR CI. Do not attribute historical green CI to the new branch.
- This remains an **offline operator-controlled synthetic identity/snapshot**, not production IAM, live SIEM ingestion, durable custody or a signed evidence artifact. OPEN-013 remains unresolved.
- Historical PR #45 is preserved until this clean integration is tested. Avoid broad merge or product-spec changes, no change to `main`.

## Evidence

This note alone is not CI evidence. Link exact workflow runs and independent L2 findings here only after inspection. The substantial functional integration merits ONE independent read-only L2 on a stable immutable published snapshot, not an L2 for every formatting fix.

## New functional analyst context

The integrated dossier additionally presents `nearby_entries`: eligible **owner-authorized** Case Timeline entries from different or same sources in the five-minute **half-open** interval `[selected-5m, selected+5m)`, excluding the selected entry, malformed timestamps and entries marked temporally uncertain. `context_time_certain=false` means no temporal inference is made; it does NOT override source ACLs or claim causality. The already-authorized pagination order is retained. A separate disk-backed behavioral test checks a nearby telemetry entry accompanying selected evidence, a distant entry excluded, and denied/foreign entries absent from the sensitive JSON. No new source reads, background store, DB, query service, or data authority.
