# Verified on-demand Timeline measurements — immutable source 2071521b

Source code SHA: `2071521b93a73be562e1bcf67e59f8c6c984e297`; branch `engineering/autonomous-development-foundation`, PR #36.
GitHub Actions [run #38842504283](https://github.com/tobianahillel-afk/cmdr/actions/runs/38842504283), job `116785564745`: **SUCCESS**, confirmed from job results and decoded log. Audit date: 2026-10-16 UTC.

## Exact measured outcomes (from JSON printed in the verified job log)

| Measurement | Target | Source | Observed | Engine validation |
|---|---|---|---|---|
| Synthetic local latency | `PERF-TGT-INV-TIMELINE-PROJECTION` | top-level source SHA above; internal **source digest** `eea6241cbfce93405077a1bc71ca90f89f7873f6777faa6143cdfd60a14f3227` | `p95=2.6024978849999996 ms` versus provisional `5 ms` budget | PASS; 1 selected, 1 executed, 0 reused; 5 samples |
| Go heap / allocation / throughput | `PERF-DEEP-INV-TIMELINE-PROJECTION` | top-level same source SHA | 1,000 projections; sampled peak Go heap `4,932,224 bytes`; total allocated `1,887,097,928 bytes`; elapsed `2,368 ms`; throughput `422.15070974189393 ops/s` | PASS; 1 selected and 1 executed |

**Interpretation:** the registry p95 aggregates **five sample-means of 200 iterations each**, and does **not** measure the 95th percentile of individual requests. The deep measurement reports a periodically sampled Go heap, **not** peak process RSS. All 5,000 synthetic entries were caller-marked `allow`; this does **not** establish authoritative Case/Timeline/source permissions, network/database/UI end-to-end performance or a production SLO. The numeric measurements apply **only to this exact SHA**; no accepted regression baseline was created. Navigation remains preimplementation.

**Verified minor workflow defect:** although run #38842504283 concluded SUCCESS and both JSON checks executed, the source-SHA line in the GitHub job summary used shell backtick command substitution. The decoded log shows `line 34: 2071521b93a73be562e1bcf67e59f8c6c984e297: command not found`; Bash `echo` masked the substitution failure and returned success. The workflow's JSON results/validation are genuine, but its summary line was defective. Correct in the **subsequent commit** by writing the source SHA without shell command substitution; re-run on that new SHA before claiming the new summary clean. Preserve the earlier run as-is rather than rewriting evidence.

L2 review for the original CI increment was scheduled on this source SHA, but no reviewer conclusions have been retrieved in this developer conversation: `L2_EN_ATTENTE`.
