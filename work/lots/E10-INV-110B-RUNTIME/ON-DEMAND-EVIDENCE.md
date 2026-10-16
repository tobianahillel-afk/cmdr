# E10-INV-110B — exact-SHA CI on-demand projection evidence

This unit extends the existing CMDR compiled performance target selector, rather than implementing a second benchmark evaluator. The new command allows a single `--requested-performance-target PERF-TGT-...` at **on-demand** stage; it fails closed on unknown targets, disallowed stages and configuration errors.

The **investigation-timeline-evidence** workflow is triggered once by the initial push of its own workflow file on `engineering/autonomous-development-foundation`; later runs require an explicit GitHub workflow dispatch. It builds the first-party Go kernel and calls:
- `performance-benchmark-audit --stage on-demand --requested-performance-target PERF-TGT-INV-TIMELINE-PROJECTION --json`
- `deep-performance-audit --stage on-demand --requested-performance-target PERF-DEEP-INV-TIMELINE-PROJECTION --json`

The workflow checks `source_sha`, `status=pass`, exact single target, executed count, and absence of cache reuse. It prints the JSON results and an explicit GitHub Actions job summary for immutable commit/run lookup; this is not an approval or accepted baseline. If either control fails the run fails, with no fabricated positive report.

The registered `investigation-timeline-navigation` operation stays blocked and is not selected. All inputs remain **synthetic, caller-claimed authorized, single process and provider-free**. The latency registry p95 is across five batch-means (200 projections each), **not request-level p95**. Resource evidence reports periodically sampled Go heap, not RSS. No production SLO, backend/index selection, or ACL-provider verification is asserted; Open Decision `OPEN-013` remains open.

Run URLs, measured values and observed outcomes must be appended **only after verifying exact GitHub Actions results**. The historical `go test -bench` results cannot be reused for the new reverse-chronological corpus.

## Durable machine-readable evidence (new workflow version)

The workflow uses the official GitHub `actions/upload-artifact` release **v7.0.2** pinned to verified tag commit `cf430e030ddbb5b0abf93d22962f4752f3646cd9`, with read-only repository permissions. After BOTH exact-SHA audits and strict Python identity assertions pass, it uploads `timeline-latency.json` and `timeline-resource.json` as a named artifact containing source SHA and GitHub run ID; retention is 30 days. Failed checks cannot produce a successful artifact. Older runs have only logs; never claim the new artifact exists for them.

The GitHub job also proves fail-closed CLI integration by invoking an explicit benchmark in forbidden `pr` stage and an unregistered target in `on-demand`, requiring the expected rejection diagnostic in both cases. These are negative integration tests of the existing CMDR selection logic, not new authorization-provider tests. This workflow remains an optional exact-SHA synthetic engineering evidence generator; changing this workflow does not provide a production SLO or automatic independent L2 approval.
