# E10-INV-120A — contextual timeline pivot from two authorized Cases

## Functional step

The existing offline, local analyst journey compares two independently authorized Cases and lists up to twenty correlated visible Timeline Entry references. A certain-time entry now has an additional **“Examine ±30 seconds around this event”** action. It switches to the entry's Case and calls the existing authenticated server `/case-timeline` route with a precise UTC `[event-30s, event+30s)` window and no source/ID query, so the analyst can discover nearby **authorized** events, not just the selected event. Existing kind filter remains active.

Only events returned by both fully authorized paginated Case reads are candidates. The action is absent for uncertain times and long entry IDs, requires the exact originating bearer, Case selection and comparison revision, clears previous correlation and source preview, then reauthorizes on the server. Every navigation is a fresh read; no correlation result is a capability. The UTC window remains exact across local daylight-saving ambiguity using existing `restoredRange`. A changed bearer or a denied destination yields no historical content.

## Technology decision and boundaries

Reuse the existing JavaScript `Date`/URL, Case-pager and Case/Timeline/source ACL provider; do not introduce a graph engine, time-index DB, new routing/authorization wrapper, new endpoint, persistence schema or third-party package for this bounded 5,000-row offline snapshot. Browser-only triage with operator-supplied ACL snapshot; not live SIEM search, deployed IAM, immutable evidence, transactionally consistent results, production performance or causal inference.

## Acceptance

Behavioral Node VM tests on real request URLs: authorized two-Case read then exact 60s bounded re-query with fresh bearer, no query truncation, stale control denied after a credential edit, server-side revocation denied, uncertain timestamp has no time pivot. Existing Go/Node CI, Go race/vet, GitHub engineering-kernel and Timeline workflow should be checked on exact published SHA. A single risk-based L2 of the final functional PR #41 milestone, not each microscopic correction. Source owner `OPEN-013` remains open.
