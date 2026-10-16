# E10-INV-120A — merge conflict reconciliation, functional not cosmetic

## Why the latest PR #41 stopped receiving usable CI

At GitHub read time, base `engineering/autonomous-development-foundation` was `a02fc9233772e9c45cf9fe99547ddefb9bed01f2`, and the feature branch was `74dbe21ebe5661640ec4c2b107d217fdad6e5cb0`, with shared ancestor `3035c6ad0ab48010a5bad862939c48efc2614456`. GitHub flagged PR #41 `mergeable_state=dirty`: both branches had modified `analyst.html`, `analyst.js` and `analyst_summary.test.cjs`. This disabled useful pull_request checks on the new feature SHA. No CI pass claim is made for the merge until actual results are verified.

## Explicit integration rather than dropping either feature

The base's two commits introduce Case-wide observed/inferred counts, certain-time envelopes and an HTTP test proving independent Case-scoped source ACLs across two imported Cases. The feature branch adds authorized source-reference intersections, cross-Case source pivots, explicit local JSON exports, and per-common-source chronology classification.

We preserve **both**: one `sourcesFor` pass already authorized by existing Go handlers computes Case-wide counts and UTC envelopes alongside per-source correlation. The rendered comparison and explicit JSON export include both sets. The base's independent Case/ACL integration test and provenance document are retained. No new API, library, provider or persistent state is introduced. The merge commit has the feature tip and current foundation tip as two parents, retaining immutable history and avoiding a force push.

## Review boundaries

CI and L2 must inspect the exact merged SHA. The previously scheduled L2 for `74dbe21` audits that immutable subfeature only; it does not automatically approve the merged integration. The snapshot and source owner ACLs are still operator-provided synthetic data in a loopback demo, not IAM, live SIEM ingestion or forensic chain-of-custody. `OPEN-013` remains open.
