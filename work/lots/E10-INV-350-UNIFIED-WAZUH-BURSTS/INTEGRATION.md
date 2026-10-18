# E10-INV-350 — Consolidate PR #51 offline Wazuh bursts with SearchJob bursts

## Integration outcome and concrete analyst capability

CMDR has two diverged independent draft paths:
- [PR #51](https://github.com/tobianahillel-afk/cmdr/pull/51): bounded source-owner-authorized **offline** Wazuh 4.x triage detecting concentrations of at least three alerts sharing (agent_id, rule_id) within an inclusive 60-second window, with a source-free JSON lead.
- [PR #58](https://github.com/tobianahillel-afk/cmdr/pull/58), now an ancestor of [PR #74](https://github.com/tobianahillel-afk/cmdr/pull/74): actual first-party Wazuh **Event Search Jobs** returning the strongest 60-second burst within already authorized search hits, including source-linked member EntryIDs.

**Reconcile these two branches in the existing product flow rather than implementing yet another correlation engine.** One `strongestWazuhBurstWindow` is now shared by both offline triage and authorized Search Jobs: sort observed UTC alert times; maintain a two-pointer inclusive window capped at 60 seconds with in-window feed frequency counts; select the highest alert count, break ties by number of distinct imported feeds, choose no lead below three observations. Both modes rank count, distinct feeds and UTC deterministically and cap at 10. A source ID tie-breaker is deterministic for search receipts, but neither mode trusts source data as a permission grant.

The original offline first-party Wazuh rule/agent/level triage now also exposes `temporal_bursts` generated **inside its existing `collectDemoCaseTriageWithObserver`** callback, AFTER the Case/Timeline/typed-source owner authorizes the event. Its summary includes no raw Wazuh IDs, host, IP, description, bearer or EntryIDs. The search result path preserves its current source-linked `search_bursts` with member IDs from the *completed search job*, never arbitrary Case data. The offline summary does **not** gain interactive source references.

This is a genuine consolidation: the independent PR #51 algorithm and test evidence are ported into the current development branch, the pre-existing PR #58 copy of its sliding window is **removed**, and both behaviors are served by one small shared Go function. The old independent PR #51 remains unmerged; do NOT double-merge overlapping files later. GitHub integration sequencing should reconcile/retire its draft only after actual approvals and safe merge of the canonical ancestor chain; no premature close, force-merge or branch rebasing.

## Proof and safeguards

- Retain independent PR #51 tests on real Wazuh JSONL feeds with two original-feed inputs, owner authorization, 60s boundary, multiple rule/agent groups, ranking, and no source-preview leakage.
- Retain PR #58 real source-linked Search Job test and membership validation, browser tests and all cross-Case Wazuh cohorts from PR #74 unchanged.
- New real behavioral equivalence test compares the offline and SearchJob burst outputs derived from the exact same Wazuh observations. Identical UTC/rule/agent/count/feed/window semantics are mandatory; only authorized SearchJob mode may return exact EntryIDs.
- No new Wazuh parser, IAM, source authorization wrapper, database, network connector, index, external Go dependency or paid API. Explicit offline synthetic operator-controlled lab only; a co-occurrence burst is NOT proof of compromise or forensic chain of custody.
- Preserve frozen parent PR #74 SHA `65932273d81b07e9c37d2234ae8aaf1641cb6aee`; independent L2 SHA #51 also untouched. Test this child PR's exact HEAD via `investigation-timeline-runtime` (Go + Node on Linux, Go on macOS/Windows), and unchanged `engineering-kernel` 90% mandatory changed-security-critical coverage. Parent #74 had 80.15% security gate failure; no threshold waiver, no claim of mergeable product while the gate or ancestor approvals fail.
- Once stable, schedule ONE targeted independent read-only L2 for the material PR #51/#58 reconciliation. Do not duplicate preexisting L2 for each commit or test adjustment.
