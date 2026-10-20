# E10-CMD-109A-FRESHNESS — canonical Command Work Freshness Class-0 slice

## Product authority and ownership

CAP-CMD-109 §§1–27 (Command / Incidents and Work Queue) requires visible source, last-observed time, explainable freshness, and explicit dependency degradation. CAP-CMD-101 consumes the projection in its existing six-view queue. The canonical Incident and Task objects remain Command-owned. Shared Data Quality/Settings owns any external health source; tenant/module policy remains authoritative to the host and may be absent.

Engineering parent is **E10-CMD-109** (newly materialized epic), child E10-CMD-109A-FRESHNESS; dependency **E10-SCALE-005B VERIFIED**. Git integration base is the exact immutable **PR #163** head `b5e71280e2af3960cf63f45fea92a71cff557a68`. **Git ancestry is not product or work parent.** CAP-CMD-101 selection PR #163 remains L2-frozen and unchanged.

## Requirement → code → evidence → remaining

| Canonical requirement | Behavior here | Evidence | Remaining |
|---|---|---|---|
| CAP-CMD-109 §§3,8,11 source/time/status | Optional first-party owner `SourceFreshness` is mapped *after* tenant/environment + Owner.CanRead to `FreshnessProjection`; supports `current/aging/stale/partial/source-unavailable/unknown` and source/time/version | `product-runtime/work-queue/freshness.go`, Go behavior tests and CI | Real host provider data and health pipeline integration |
| §§4,19,26 no invented thresholds | `current/aging/stale` only when owner supplies explicit versioned `AgingAfterSeconds/StaleAfterSeconds`; absent/invalid/inconsistent/future metadata remains `unknown`, not current | direct threshold tests and absence of policy tests | Tenant/module policy ownership, policy version registry and audit |
| §§15,22 partial source without losing valid data | Only authorized items in the selected view/filter mark the projection `partial-sources` and `work-source-degraded` when host says `partial` or `unavailable`; rows preserved | filtered-authorized Go tests | Linked Case/Run source health, safe read of remote owner and transitions |
| CAP-CMD-101 §§3,8/11 existing route | Exactly six system views, single GET, existing ACL and 5000/50 bounds, same browser representation, full native escaping; missing owner metadata keeps original JSON field set | multi-platform work-queue regression suite | Shared Saved Views, Inspector/Incident Detail and scroll restoration, no final cursor design |

### Scope and decisions

This is a **read-only Class-0 host-owned projection** of authorized Command source metadata, NOT owner-independent SLA policy, historical evidence, health telemetry ingestion, a new authorization service, a durable record or a right to mutate. Owner provenance/source labels are not external URLs and never grant access. No client-provided thresholds, new route, query DSL, paid API or third-party product engine. Data remains a synthetic host adapter until production IAM and trusted source metadata are supplied. Source clock skew yields `unknown` rather than a false healthy status. A local source-health failure does not delete otherwise authorized Incident/Task rows.

The time thresholds in tests are **test fixtures only** (60/120 seconds), not normative CMDR defaults. Source freshness is distinct from Incident/Task `UpdatedAt` used for sorting. No change to canonical specification or OPEN decision. No cursor schema or final SLA configuration invented.

## Acceptance and handoff

1. One authorized owner row with explicit versioned policy produces deterministic freshness `current/aging/stale`, with source, UTC observed time and policy in JSON/browser; no policy yields `unknown`.
2. Explicit owner `source-unavailable/partial` preserves otherwise authorized items and marks only matching current view as partial; missing future health is never guessed.
3. Denied/foreign source identifiers are absent from JSON/HTML, malformed owner metadata fails closed without partial rows, HTML source text is contextually escaped.
4. Previous six-view filtering, sorting, selection, first-50 bound, and JSON without source metadata remain compatible, no extra Owner.List/CanRead.
5. **Exact SHA** command-work-queue-runtime Linux/macOS/Windows Go and engineering-kernel mandatory security floor/architecture/manifest checks MUST pass. A successful child does not approve parent #163, foundational draft #36/#2 or production release.

**GitHub writes** occur only in this new feature branch, with no changes to the frozen PR #163 SHA. Once final exact HEAD is stable: update PR evidence, schedule ONE independent risk-based read-only L2 for CAP-CMD-109 new source-metadata trust boundary; no duplicate existing CAP-CMD-101 L2. **Status IMPLEMENTED only**, even if workflows succeed; no native review, merge, accepted CAP or prod claim. A future agent should re-check HEAD and all gates before touching files.
