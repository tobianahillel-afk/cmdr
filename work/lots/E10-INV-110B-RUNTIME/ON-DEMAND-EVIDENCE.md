# E10-INV-110B — exact-SHA CI on-demand projection evidence

This unit extends the existing CMDR compiled performance target selector, rather than implementing a second benchmark evaluator. The new command allows a single `--requested-performance-target PERF-TGT-...` at **on-demand** stage; it fails closed on unknown targets, disallowed stages and configuration errors.

The **investigation-timeline-evidence** workflow is triggered once by the initial push of its own workflow file on `engineering/autonomous-development-foundation`; later runs require an explicit GitHub workflow dispatch. It builds the first-party Go kernel and calls:
- `performance-benchmark-audit --stage on-demand --requested-performance-target PERF-TGT-INV-TIMELINE-PROJECTION --json`
- `deep-performance-audit --stage on-demand --requested-performance-target PERF-DEEP-INV-TIMELINE-PROJECTION --json`

The workflow checks `source_sha`, `status=pass`, exact single target, executed count, and absence of cache reuse. It prints the JSON results and an explicit GitHub Actions job summary for immutable commit/run lookup; this is not an approval or accepted baseline. If either control fails the run fails, with no fabricated positive report.

The registered `investigation-timeline-navigation` operation stays blocked and is not selected. All inputs remain **synthetic, caller-claimed authorized, single process and provider-free**. The latency registry p95 is across five batch-means (200 projections each), **not request-level p95**. Resource evidence reports periodically sampled Go heap, not RSS. No production SLO, backend/index selection, or ACL-provider verification is asserted; Open Decision `OPEN-013` remains open.

Run URLs, measured values and observed outcomes must be appended **only after verifying exact GitHub Actions results**. The historical `go test -bench` results cannot be reused for the new reverse-chronological corpus.
