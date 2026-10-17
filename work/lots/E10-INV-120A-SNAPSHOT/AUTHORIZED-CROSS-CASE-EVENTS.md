# E10-INV-120A — bounded offline cross-Case event chronology

## Functional capability (risk-based milestone)

New explicit CLI mode:
```sh
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/path/owner.json -compare-case=case-a -with-case=case-c -compare-events
```

This emits ONE operator-local, **sensitive** unsigned JSON receipt combining the read-only Event timelines of two independently authorized Cases. Each event retains its owning Case ID, entry ID, observed/inferred classification, provenance, source reference, source timestamp, uncertainty, version and accessible supersession; the receipt also contains each Case's existing authorized triage summary. Certain timestamps are UTC-ordered; uncertain/unknown entries follow, ordered deterministically without inventing their position in actual time. No Case/source identity correlation or causality is inferred. Up to 5,000 events per Case (10,000 total) and startup file 8 MiB are inherited from existing boundaries.

## Choice of technology

Reused `collectDemoCaseTriageWithObserver` (authenticated Case discovery + paginated Timeline + typed-source owner ACL) twice, once per selected Case, plus `time.Parse`, `sort.Slice`, and `encoding/json` in the existing Go module. This is one bounded chronological join in memory, not a new generic correlation/index engine, query language, database, IAM connector, backend endpoint, plugin or paid dependency. No production object contract or cross-module owner is modified.

## Security/DFIR constraints and acceptance

- `-compare-events` requires two different Case IDs and rejects combinations with single-Case report or source timeline modes; the previous aggregate-only comparison remains unchanged.
- The operator file is selected on the local CLI, never by browser. Both Cases must independently pass the existing authenticated Case picker and all Timeline and source permissions, including across multiple pages, **before writing a single byte**.
- The export deliberately includes authorized Event IDs, source references and provenance; unlike the earlier aggregate-only report, it is sensitive. No bearer, source preview, raw event payload, denied-source row or foreign-tenant Case is included.
- Tests exercise **real operator JSON on disk**, 207 visible rows across two pages in Case A and one from Case C, mixed observed/inferred kinds, deterministic chronology, denied/foreign source suppression, unknown timestamps sorted after certain ones, unauthorized second Case and option conflicts, cancellation and failing output writer.
- Do not treat this as production evidence export, cryptographic chain of custody, signed artifact, immutable/transactional multi-Case snapshot or a live SIEM correlation service. The selected principal/owner dataset and bearer remain synthetic. OPEN-013, production IAM, durable owner/store integration and audit backends remain unresolved.

The code is a materially new investigator-visible functional feature in already READY `E10-INV-120A-SNAPSHOT`. Keep this lot READY until exact-SHA Go/Node CI succeeds and one independent risk-based, read-only L2 evaluates the integration. Do not schedule duplicate L2s for trivial follow-up formatting or tests.
