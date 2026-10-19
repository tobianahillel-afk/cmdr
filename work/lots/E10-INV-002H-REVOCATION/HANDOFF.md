# E10-INV-002H-REVOCATION — bounded GitHub recovery handoff

## Exact responsibility and scope

- Canonical behavior: **CAP-INV-002 — Event Search**, Investigate / Signals and Hunt, consuming **CAP-INV-110 Investigation Timeline** / Shared typed-source owner. Canonical read-only permission enforcement; never broaden raw telemetry, source ACL, Case scope or ownership.
- Engineering unit: pre-existing `E10-INV-002H-REVOCATION`, actual parent `E10-INV-002`, dependency `E10-INV-002G` + `E10-INV-110E`, recorded **IMPLEMENTED, not VERIFIED**. Git parent PR #135 is not its engineering parent.
- Code placement: first-party Go synthetic operator-local `product-runtime/investigation-timeline/cmd/case-timeline-demo`. Not production IAM, real ingestion, persistent Shared SearchJob or canonical CAP acceptance.
- Canonical sources: `cmdr-product-spec/07-investigate/modules/signals-and-hunt/capabilities/event-search.md` and `cmdr-product-spec/07-investigate/modules/cases-and-evidence/capabilities/investigation-timeline.md`, permission-enforcement contract, canonical Case/Telemetry Event/Search Job objects. `OPEN-013` remains unresolved for Case linking or Timeline mutations.

## Implemented bounded behavior and proof

The actual integration test `product-runtime/investigation-timeline/cmd/case-timeline-demo/normalized_search_revocation_integration_test.go` imports two distinct normalized JSONL events within the requested search window into an authorized Case. A completed Search Job returns the **first matching** source reference; the second is an **authorized nonmatching** candidate. When the second typed-source owner re-open is revoked or the request cancelled, the **entire** subsequent search must fail. Neither a partial prior hit, a Search Job, source preview nor bearer may be returned. The test checks a real HTTP `POST /demo-local-search` against existing first-party Case/Timeline/source authorization; it adds **no new runtime behavior**, dependency or exception to a required check.

Git publication was verified for existing PR [#136](https://github.com/tobianahillel-afk/cmdr/pull/136), based on frozen PR [#135](https://github.com/tobianahillel-afk/cmdr/pull/135) SHA `4f27c41c4554b21fe835fc2a366b2da9b46a1423`, not `main`. On the *original published behavior/test SHA* `7feae9601f03682fbbf7dd8acbb46410004a54d9`:

- [investigation-timeline-runtime run 39216688574](https://github.com/tobianahillel-afk/cmdr/actions/runs/39216688574): COMPLETED **SUCCESS**, three jobs Linux Go/Node browser, macOS Go, Windows Go.
- [engineering-kernel run 39216688631](https://github.com/tobianahillel-afk/cmdr/actions/runs/39216688631): COMPLETED **FAILURE**, mandatory `CHK-SECURITY-TESTS`: investigation-timeline-runtime Go changed security-critical coverage **89.33% < 90.00%**. Reported diagnostics include `case_http.go`, `case_discovery.go`, `normalized_search.go`, `normalized_snapshot.go` and CLI. A green functional workflow does **not** waive this real release blocker.
- PR #136 draft/open, no GitHub-native APPROVED reviews identified during this reconciliation. Parent PR #135 is draft with a separately failing required kernel; the root engineering PR #36 remains draft (763 commits, 871 changed files), root product-spec PR #2 also draft. The full product is **not** merged or production-ready.

**Reconciliation rules:** this handoff and execution-cache updates do not retroactively change the test's source SHA above or transfer its CI to a new documentation SHA. Any new commit triggers new exact-HEAD checks; report their actual conclusions separately. Existing canonical acceptance, security floor and approved review requirements remain unchanged. No independent L2 is triggered for this **test-only / recovery-documentation** change; do not create a duplicate L2 for parents.

## Actual remaining work and next safe selection

- `work/graph.json` currently has **149 units: 124 VERIFIED, 17 DECOMPOSED, 8 IMPLEMENTED, 0 READY**. Do **not** infer a READY product task from the highest PR number or relabel E10-INV-002H VERIFIED.
- Product spec CAP-INV-002 expects multi-source partial/error presentation, query provenance, cancellation/replay and Case return context. This test only proves local **fail-closed owner revocation** for a synthetic single SearchJob source; it does **not** implement full partial results or source freshness. A denied permission must never be treated as a permissible partial result.
- Re-evaluate the *common* security coverage blocker with a risk-based, design-level inspection of uncovered critical Go paths and the actual canonical permissions model. After repeated test-only attempts without a pass, do **not** add redundant near-identical regressions. If a concrete material source bug is demonstrated, define a bounded new work unit with real CAP/engineering parent, allowed paths and reproducible acceptance before changing code; otherwise pursue an independent authorized ready unit.
- Address draft stack integration and upstream rules on PR #36/#2 and the selected parents. No retroactive SHA mutations, waivers, native approvals or premature merges.
- `engineering/state/current.json` is **execution-cache only**; the work graph, product spec, exact Git SHAs, CI and acceptance each answer different questions. The next agent must re-fetch `main`, PR heads, status checks, graph and L2 evidence.

This handoff adds no executable functionality and does not claim a new product delivery. The actual behavior belongs to the original source commit and CI above; any new documentation HEAD is just a recoverable checkpoint.
