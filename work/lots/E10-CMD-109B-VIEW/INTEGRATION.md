# E10-CMD-109B-VIEW — Authorized whole-view source freshness

## Canonical ownership and actual gap

Product **CAP-CMD-109** (Work Freshness and Staleness) under Command / Incidents and Work Queue; **CAP-CMD-101** consumes its existing single Work Queue projection. Engineering epic **E10-CMD-109**, new bounded task **E10-CMD-109B-VIEW**, depends on already-implemented **109A owner freshness** and **101G authorized cursor pages**. Code owner `product-runtime/work-queue`. Git base PR #168 is **NOT** the business parent.

Canonical CAP-CMD-109 §§3,7–8,11,14–16,19,22,25 requires a meaningful **view-level** source freshness presentation, not only per-item states. Existing PR #168 and prior ancestors already exposed policy-versioned `SourceFreshness` for individual authorized Incident/Task rows, but no aggregate covering the whole selected filtered cohort. On multi-page queues the analyst could not tell that stale items existed beyond the current page. This new Class-0 projection uses the SAME first-party `Project`, host-only source policy, existing single `GET /command/work-queue` route, six canonical views and server-side permissions; adds **no new host adapter, API, threshold, persistence or routing surface**.

### Real end-to-end behavior

After the established first-party `Owner.List` and per-object `Owner.CanRead` authorizations, view choice and manual exact status filtering, a single O(n) pass over the resulting bounded (<=5000) authorized matching cohort counts seven figures: total, current, aging, stale, partial, source-unavailable and unknown. Unknown includes matching items with absent source metadata *when at least one source reports metadata*; a fully metadata-free owner has **no new `view_freshness` field**, preserving legacy JSON. The existing source projector alone defines aging/stale using explicitly versioned policy thresholds and observed_at; item UpdatedAt is NOT used for inferred freshness. No source names, tenant metadata or raw item IDs are added to the aggregate.

The same non-sensitive counts are available in selected scope's JSON and native `html/template` single-route HTML. Unlike the 50 displayed rows/page the view summary covers **all** authorized matches even on later pages. The existing cursor owner-validation must pass before anything can be returned; unauthorized, foreign, filtered-out or newly revoked objects never contribute. Degraded sources still drive existing `partial-sources` diagnostics; an unknown timestamp/policy is NOT a source outage. This is an explainable read-only view, not a health inference or a refresh job.

### Acceptance / risk checks

1. Current, stale and unavailable items, denied other-tenant and source-hidden items: report only current authorized View+Status cohort counts; no source metadata leakage.
2. Page 1 with >50 items and page 2 must report the SAME full authorized summary. Revoking an item after cursor creation returns no rows, no counts, and the existing stale-cursor error.
3. Null metadata is explicitly unknown only in an otherwise reporting cohort; all-null legacy provider does not grow the JSON schema.
4. Partial source, unavailable source and policy-unknown have distinct semantics; original Work Queue state/diagnostics remain unmodified. Native HTML contains counts without leaking underlying source identifiers.
5. Mandatory `command-work-queue-runtime` (Linux, macOS, Windows Go tests) and `engineering-kernel` checks including CHK-MANIFESTS, CHK-COMPLEXITY, CHK-ARCHITECTURE, CHK-RUNTIME-DEPS, CHK-SECURITY-TESTS >=90% changed-critical coverage and Go SAST on EXACT SHA. None can be waived.

### Exclusions and integration blockers

The epic and capability are NOT canonically accepted; no production source adapter, Shared health engine, refresh/job action, long-term event storage, Signed forensic custody, or Command Class-2 mutation. This child of frozen parent **PR #168 SHA 40ee320eeaac44847f9a5c088f91ff451030a21e** remains a DRAFT until both CI and required upstream approvals/reviews are actually satisfied. Long parent chain (#2/#36 and many Command/Investigate branches) must still be consolidated before merging to main. One independent read-only L2 for this material view-level integration, not for each commit or test.
