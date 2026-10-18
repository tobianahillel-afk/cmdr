# E10-INV-360A — Integrate PR #52 owner reference/provenance search with actual Case Search Jobs

## Functional objective

Independent draft [PR #52](https://github.com/tobianahillel-afk/cmdr/pull/52) provides an offline Search Workbench over **already authorized** Case Timeline event IDs, source references, and provenance. Current canonical stack [PR #75](https://github.com/tobianahillel-afk/cmdr/pull/75) already has a richer first-party Event Search Job, Case/Timeline/typed-source owner, authenticated analyst UI, normalized JSONL imports, retained Wazuh source metadata, rule/agent pivots, source-linked chronological Search Job results and unified 60-second burst leads. Maintaining a **second** search orchestration, listener, CLI or IAM would increase overlap and review debt without product value.

This increment incorporates **the useful search semantics of PR #52** in the **single existing** `POST /demo-local-search` / `event-search` orchestrator and existing analyst UI. Select one of three extra **literal metadata fields**:
- `entry_id`: bounded case-insensitive substring of the Case-owner-authorized Timeline Entry identity;
- `source_ref`: bounded case-insensitive substring of its typed-source reference;
- `provenance`: bounded case-insensitive substring of its normalized import origin slot (`operator-local-jsonl` / `operator-local-jsonl-feed-N`).

The actual source and case ACL are **still re-evaluated first** for every time-eligible imported event through the canonical typed-source owner; only then compare its **already authorized** projected metadata. Search remains bounded to at most 500 time-scoped candidates and 200 completed result references, with all-or-nothing failure on denied owner, malformed input, invalid time or overflow. It shares existing source-open navigation/revocation, import scope, principal/tenant/Case restrictions, Search Job audit and the exact same browser result list. Query selector and literal are sent only in authenticated, max-4KiB **POST JSON**, never in URLs. No second Search Job engine, external dependency, API, provider, database or production IAM assertion. Historical `preview`, `rule_id`, `agent_id`, Wazuh burst and cross-Case aggregation paths remain unchanged.

### Acceptance evidence

- Two actual normalized JSONL input feeds have the **same original event ID** but independent canonical Case owner Entry/Source IDs. Query `provenance=feed-2` produces only the second in-window imported event, while `feed-1` produces the first; `entry_id` substring and `source_ref` separately identify the second source in a completed first-party Event Search Job.
- Every returned reference remains independently owner-openable; raw source preview, secret messages and bearer never appear in search receipts. Foreign principal/bearer or Case and unrecognized arbitrary field names return no source identities. An absent feed completes with zero results. Events outside the selected UTC interval do not match.
- Browser uses existing `local-search-field` selection and POST; switching fields clears old results and source-opening callbacks. Existing Wazuh agent/rule/search burst UI and tests remain.
- Tests use actual Go source owner + Event Search orchestrator and browser State tests on Linux; exact-head GitHub Actions on macOS/Windows and the **unchanged mandatory 90% security-critical coverage gate**. Failures block merge.
- `entry_id`, `source_ref`, and `provenance` describe authorized synthetic **source metadata**, not payload search; metadata visibility can be sensitive, which is why authorization happens BEFORE even comparing a filter.

### Reconciliation / integration risk

This is a child of immutable PR #75 HEAD `e99c881d872098a220b23b27cd1b29c80ab59dad`; #75 itself consolidated PR #51 burst algorithm with PR #58. Independent PR #52 remains draft and **must not be merged wholesale into the canonical chain** in a way that duplicates the search frontend/backends or changes security policy. Only after authentic GitHub review and approved integration should it be reconciled/retired. No parent L2 audited SHA changes and no premature merge. Main is not a production SIEM. This is operator-provided offline synthetic Case data with no real IAM, continuous ingestion, durable Search Job persistence, signed custody or verified attack classification.

One read-only independent risk-based L2 is appropriate for this materially useful PR #52/#75 integration, on its exact published base..head. Do not request another L2 for minor fixes or repeat full parent reviews.
