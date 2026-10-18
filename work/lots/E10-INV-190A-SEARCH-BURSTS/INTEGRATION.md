# E10-INV-190A — Authorized Wazuh Search Job temporal correlations

## Functional integration

This work integrates the **existing 60-second sliding-window alert burst algorithm from draft PR #51** with the first-party, authenticated **Event Search Job** investigation flow at PR #57. It is more than a standalone cohort counter:

1. Analyst searches offline operator-selected Wazuh 4.x feeds by normalized agent or rule and precise time range using the existing bearer-guarded Case owner and Event Search Orchestrator.
2. All matched source refs and their normalized rule/agent IDs, UTC timestamps and server-assigned feed provenance must pass the same complete owner-authorized Case/Timeline pagination and per-hit typed-source re-open as the existing Search Job.
3. Only AFTER completed Search Job are **bounded 60-second burst leads** derived from its already authorized, maximum 200 matched events. For each rule/agent combination, pick its strongest inclusive 60-second window of at least three alerts, with counts and distinct feed count; return at most ten ranked leads.
4. Every lead includes the exact subset of stable Entry IDs in its strongest window. The analyst UI can browse each supporting authorized event using the existing Case/source navigation (with fresh authorization on click). These are leads, **not** compromise indicators, alert classifications, incident verdicts or causal findings.

The first-party two-pointer/sort algorithm in `wazuh_search_bursts.go` is adapted from the separately implemented PR #51 `wazuh_bursts.go`; it is re-scoped to completed Search Jobs and adds exact source identities for drilldown. The overlapping algorithm should be consolidated when the independent PR #51 is reconciled into the delivery train; do not duplicate it again. No extra source scan, new external package, backend, new search DSL, product mutation or paid API.

## Security and honest limits

- Synthetic local lab only. No real production IAM, Wazuh live ingest, signed custody, persisted incidents, correlation index or backend query store.
- No rule/agent fields used as permissions. Only already-owner-authorized SearchJob hit metadata feed the burst algorithm; excluded or revoked/foreign sources never appear. The server sends a complete receipt only on completed Search Job and authorized hits.
- The HTTP result `search_bursts` is bounded 10 groups, 3–200 members each, and is returned alongside the existing `job` and `search_hits`. Each member is one reference in the completed Job, never a raw log, hostname, IP address or bearer.
- Browser validates correlation members against the same already verified SearchJob hit map, exact rule/agent, UTC interval and feed count; entire screen's new results AND query replay export fail closed if one burst group is incoherent. Clicking any supporting event requires a fresh owner source-open, never trusts a serialized result as authority.
- The grouping is limited to **the selected Search Job**, not the whole Case or a continuous Wazuh stream; changes to filters can change leads. Source clocks have not been independently corroborated. Do not interpret concentration as intrusion verification.

## Tests and acceptance

- Real operator-file two-feed authorized Search Job with four matching alert events, three in 40 seconds and a fourth outside that burst; validates one lead, 2 feeds, correct supporting Entry IDs, source privacy, denied bearer and no raw leaks.
- Inclusive 60-second boundary and >60-second exclusion, minimum-three threshold, different rule/agent separation.
- Browser accepts a valid source-linked lead and rejects foreign member IDs atomically with no stale query replay.
- Exact-head GitHub CI (Linux/Windows/macOS Go+Node and engineering-kernel) must be verified; mandatory 90% changed security coverage **never waived**. Existing PR #57 `engineering-kernel` failure 80.00% remains upstream merge-blocking until repaired; no amount of new functionality is permission to merge.
- Branch intentionally stacks on parent PR #57 audited SHA `01011d157e598fd7e61cad9282735f24d02f30c5`, leaving PR #57/#56/#55/#54/#53/#50/#51/#36/#2 frozen. Coordinate reconciliation of the independent PR #51 without force-pushing or bypassing reviews.
