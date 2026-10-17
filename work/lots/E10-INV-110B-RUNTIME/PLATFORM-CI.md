# E10-INV-110B — Windows/macOS process-boundary verification

The Timeline subprocess runner has separate Go implementations: `performance_timeline_process_unix.go` (Linux/macOS process-group SIGKILL) and `performance_timeline_process_windows.go` (best-effort Windows taskkill /T /F /PID). Earlier successful engineering-kernel CI ran on Linux only; it did not test the OS-specific paths.

Reuse the existing, first-party Go tests and authorized, SHA-pinned GitHub Actions workflow `.github/workflows/investigation-timeline-evidence.yml`. Add a two-OS runner matrix for Windows and macOS with actual Go build, go vet, cancellation, descendant, output-bound and projection-probe tests. Cross-compiling alone would not test child-process cancellation. No paid API, product dependency, store/index or new workflow is introduced.

This is synthetic GitHub-runner verification only, not support certification for every OS version, every subprocess session, provider authorization or production SLO. Linux on-demand benchmark evidence remains separate and unchanged. Do not assert any Windows/macOS result until the exact-SHA matrix job conclusions and test logs have been read. Work unit B remains READY pending complete scope and L2, navigation C stays BLOCKED.

## GitHub evidence verified at exact SHA 939aca2f55d5f453bfb9360c6ea7434ddb3ed772

- Workflow [investigation-timeline-evidence #38847405796](https://github.com/tobianahillel-afk/cmdr/actions/runs/38847405796): **SUCCESS** with 3/3 jobs.
  - Windows AMD64 job **116800954371**: **SUCCESS**, Go 1.26.8, Go build, vet, real probe and subprocess tests including `TestTimelineProcessCancellationKillsDescendants`, `TestTimelineProcessBoundsStdoutAndStderr`, `TestTimelineBoundedWriterDoesNotBlockOnOverflow`.
  - macOS ARM64 job **116800954606**: **SUCCESS**, Go 1.26.8, corresponding real-process tests including descendant cancellation.
  - Ubuntu evidence job **116800955072**: **SUCCESS**, explicit, non-reused synthetic projection latency and resource target runs, validated against immutable outer Git commit SHA.
- Workflow [engineering-kernel push #38847405704](https://github.com/tobianahillel-afk/cmdr/actions/runs/38847405704): **SUCCESS**.
- Workflow [engineering-kernel PR #38847410659](https://github.com/tobianahillel-afk/cmdr/actions/runs/38847410659): **SUCCESS**.
- Workflow [investigation-timeline-runtime PR #38847410648](https://github.com/tobianahillel-afk/cmdr/actions/runs/38847410648): **SUCCESS**.

### Immutable synthetic measurements (Ubuntu on this exact SHA)

- `PERF-TGT-INV-TIMELINE-PROJECTION`: observed **3.382318 ms**, budget 5 ms, PASS; this is p95 **of five batch-mean measurements**, 200 projections per batch, **not request-level p95**.
- `PERF-DEEP-INV-TIMELINE-PROJECTION`: 1,000 operations, duration 3,262 ms, periodically sampled Go heap **7,837,544 bytes**, total allocations **1,887,095,184 bytes**, PASS; this is not maximum RSS.
- Artifact [#12023636993](https://github.com/tobianahillel-afk/cmdr/actions/runs/38847405796/artifacts/12023636993) contains the two JSON reports, retention 30 days, archive digest SHA256 `7387943e5f2f15f3354264a0dc8b96ab7e8e928a49cb8192db157ec69de530e8`.
- Important field semantics: the **top-level** `source_sha` in the benchmark summary is Git HEAD `939aca2f...`; the per-target benchmark `results[0].source_sha` is the target's **content/source identity digest** (`performanceTargetSourceDigest`), here `1726a979...`, NOT an alternative Git commit. Do not interpret them as mismatched commits or conflate with a benchmark baseline.

### Explicit qualification limits

These checks verify only the specified hosted runner versions (Windows amd64, macOS arm64, Ubuntu amd64). They do **not** establish compatibility with every Windows/macOS/Linux version, provider authorization enforcement, malicious process escape from a separate session, production SLO, backend latency or actual investigation navigation. `OPEN-013` stays unresolved and `E10-INV-110B-RUNTIME` remains READY. The independent read-only L2 review is scheduled on the **code SHA 939aca2f**; no report or native GitHub approval has been observed. This documentation-only subsequent commit is not a substitute L2 snapshot.
