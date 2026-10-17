# E10-INV-120A — Local authorized Case investigation dossier

## Functional gap and implementation

An analyst previously had to run separate commands to obtain Case triage, same-source event chronology and a typed source preview for an explicitly selected entry. This increment composes these EXISTING server-authorized flows into ONE operator-triggered read-only JSON dossier, without a new provider, database, identity engine or product schema.

Use `case-timeline-demo -principal analyst-a -snapshot /absolute/path/owners.json -dossier-case case-a -dossier-entry event-old` (omit `-snapshot` for synthetic built-in fixture).

Output `cmdr-local-case-dossier-v1`: bounded Case triage summary, selected observed/inferred Timeline Entry with provenance and version, all already-authorized events that share its source reference, and the **re-authorized typed owner preview**. The flow first authenticates through existing memory owner handlers, fully traverses Case pages (up to 5000 entries), rechecks the selected entry via the existing `/demo-source-open`, then serializes once all authorization and consistency checks pass. Error, cancellation and denied Case/entry/source produce no dossier.

**Reuse-first comparison:** CMDR already implements `collectDemoCaseTriageWithObserver` and `queryDemoOwner`, `net/url`, `encoding/json`, and immutable `MemoryTimelineData`. Reusing these is simpler and safer than inventing another permission filter, client-side concatenation or introducing a search index/DB. No new product runtime dependency.

## Risk limits and next steps

This is an **offline synthetic loopback demo only**. The local selected principal is NOT production IAM; the JSON is NOT a signed evidence pack or custody record, and timestamps/source matches are observational, not causality. The output is sensitive and must not be shared publicly. A static operator-supplied snapshot can become stale; do not use real customer data. `OPEN-013` still blocks final Shared Timeline Entry mutation/storage decisions. No multiuser session, persistence mutations or production SIEM ingestion has been implemented. The PR targets the existing integration branch; no change is made to `main` or the canonical Product Spec.

Tests cover actual disk-source preview, owner-controlled Case access, matching source across multiple visible entries, denial/cross-tenant, mode conflicts, cancellation and output errors. Risk-based L2 is required on a stable published SHA; do not conflate CI green with L2 approval.
