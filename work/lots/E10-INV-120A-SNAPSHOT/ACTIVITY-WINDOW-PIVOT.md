# E10-INV-120A — authorized five-minute activity-window pivot

## Functional milestone

On PR #41's operator-selected local Case → Timeline → source-preview journey, analysts can now identify the top three **five-minute UTC windows containing the most visible, certain-timestamp entries**, then select a window to search again through the existing server-authorized `/case-timeline` route. Filters for the selected Case, entry kind and source text remain in effect. The query uses [from, to) and resets pagination.

Only already-authorized Case pages enter this calculation (up to 5000 entries); unknown/uncertain timestamps are excluded. Source-visible event gaps, dense windows and common references are *triage indicators*, not intrusion findings or proof that unavailable events are absent. No synthetic bearer or raw events are exported; the opt-in aggregate JSON adds only start/end UTC and counts for up to three windows.

## Reuse / choices

Reuse existing Go first-party Case/Timeline/typed-source ACL pipeline, `analyst.js` paging and time filters, `Date`, `Map` and existing Node VM HTTP fetch harness. A search index, rules engine, time-series DB or third-party analytics package would increase scope without value for bounded offline triage. Browser local inputs reflect selected time range; the exact UTC range is retained for daylight-saving ambiguity until edited.

## Security and verification

The click never confers access to an Entry/Case; it triggers a fresh server request. A changed Case/bearer/filters or revoked permission invalidates it. Denial leaves no historical result in the new view. Tests exercise a three-event window, unequal second window, reauthorization, bearer invalidation and unknown-only times. Only synthetic loopback HTTP is represented, not production IAM, SIEM ingestion, durable owner storage, forensic chain of custody or real 95th-percentile SLO. `OPEN-013` remains unresolved. Exact GitHub CI and one risk-based L2 review of the full stable functional lot are required; avoid redundant L2 on this small UI incremental commit.
