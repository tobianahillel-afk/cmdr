# E10-INV-110B — observed on-demand synthetic performance evidence

## Exact GitHub Actions sources

- Source repository: `tobianahillel-afk/cmdr`, PR #36, branch `engineering/autonomous-development-foundation`.
- Immutable source commit **`dae5fe2c64e1eb1befb0cc0130cc42040ba5bfdf`**, checked out in GitHub Actions run **[38842748892](https://github.com/tobianahillel-afk/cmdr/actions/runs/38842748892)**, job **116786325537**.
- `investigation-timeline-evidence` run **completed SUCCESS**; executable validates exact source SHA, one explicitly selected projection target, `executed_targets = 1`, `reused_targets = 0`, `status = pass`. Its job log is the source of the following measured figures.
- `engineering-kernel` push run **[38842748868](https://github.com/tobianahillel-afk/cmdr/actions/runs/38842748868)**, PR run **[38842752024](https://github.com/tobianahillel-afk/cmdr/actions/runs/38842752024)**, and dedicated Timeline PR run **[38842752002](https://github.com/tobianahillel-afk/cmdr/actions/runs/38842752002)** were also SUCCESS on that source commit.
- The benchmark result's internal `source_sha = eea6241cbfce93405077a1bc71ca90f89f7873f6777faa6143cdfd60a14f3227` is a **source fingerprint** rather than Git commit SHA. The top-level audit `source_sha` identifies the actual Git commit.

## Measurements, exact scope

| Recorded field | Observation |
| --- | ---: |
| `PERF-TGT-INV-TIMELINE-PROJECTION` p95 across **five means of 200 local synthetic projections** | **4.70007566 ms** |
| Provisional local engineering upper-bound | 5 ms |
| `absolute_pass`, benchmark cache `reused_targets` | true, 0 |
| `PERF-DEEP-INV-TIMELINE-PROJECTION` stage/resource | on-demand / resource |
| Deep resource `operations` | 1000 |
| `duration_ms` | 4572 |
| `peak_heap_bytes` — periodically sampled Go heap | **5,147,088 bytes** |
| `total_allocation_bytes` | 1,887,094,288 bytes |
| `throughput_ops_per_second` | 218.70193321070846 |

This is **measured GitHub CI evidence for the exact older SHA above**, not a benchmark rerun on later security/docs commits. The p95 metric is between *batch averages*, not a latency distribution of individual requests. The sampled heap is **not peak RSS**; total allocations are cumulative, not retained heap. All input sources were synthetic pre-authorized entries; no ACL provider, external network, storage, UI, export or production traffic was exercised. Therefore these numbers neither prove production SLO nor full CAP-INV-110 delivery. `OPEN-013` remains unresolved. The navigation handler is still fail-closed and not benchmarked.

Subsequent code commits must be verified by their own CI; performance results cannot be rebased or marked current merely because the old run passed. Do not suppress previously failed runs or substitute a source digest for the Git commit identity.
