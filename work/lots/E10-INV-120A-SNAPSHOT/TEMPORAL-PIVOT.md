# E10-INV-120A — reconciliation of independent temporal Case pivots

This milestone merges upstream engineering base commit `36e11e1c968d8e39666af76feb0c58c3d244abd3` into PR #41 **without overwriting its existing operator-selected offline multi-Case / source-navigation journey**. The base added five-minute UTC co-visible buckets in two authorized Cases. PR #41 independently contains source-reference and individual correlated-event pivots. All three actions are preserved in distinct, inert UI surfaces, routed through one existing `reopenAuthorizedCase` helper that reauthorizes Case, Timeline Entry and source on every new HTTP request.

- `comparison-actions`: common source-reference pivots for Case A or B.
- `comparison-windows`: common 5-minute UTC window pivots, including when source references are disjoint.
- `comparison-chronology`: individual correlated event drilldown and ±30-second contextual timeline reopening.
- All actions require the exact original bearer and selected Case/revision; none transfers authorization, and time-uncertain events do not create synthetic windows. Window timestamps are preserved exactly across DST folds. New requests start from offset zero.
- Workflow uses the upstream single Node test invocation for both `analyst_return.test.cjs` and `analyst_summary.test.cjs`, removing the duplicated setup/run. Linux + macOS + Windows Go testing and existing security gates remain intact.

Reuse JavaScript Map/Date and the existing server-authorized pagination. No new index, correlation service, IAM, persistent store, provider, external dependency or paid API. Synthetic offline data only, no production chain of custody, no transactional snapshot/causality claim. OPEN-013 remains open. Check the **two-parent merge SHA**, exact CI and risk-based L2 before qualifying this milestone. This is a scoped draft PR on the existing engineering base and is not merged into main.
