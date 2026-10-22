# E10-INV-110B — recorded source-backed synthetic microbenchmarks

## Immutable source and verification

- Code SHA: `2c5efaa9ac819aacd2a6854108c852eebc1cc999` on `engineering/autonomous-development-foundation`, PR [#36](https://github.com/tobianahillel-afk/cmdr/pull/36).
- [Dedicated Timeline workflow push #38839381575](https://github.com/tobianahillel-afk/cmdr/actions/runs/38839381575) **SUCCESS**, job `116775702265`.
- [Dedicated Timeline workflow PR #38839385911](https://github.com/tobianahillel-afk/cmdr/actions/runs/38839385911) **SUCCESS**, job `116775715824`.
- [Engineering kernel push #38839381424](https://github.com/tobianahillel-afk/cmdr/actions/runs/38839381424) **SUCCESS**; [engineering kernel PR #38839385868](https://github.com/tobianahillel-afk/cmdr/actions/runs/38839385868) **SUCCESS**.
- The initial benchmark commit `31f0dfdfbb2955ce5feb240342cd1eca63090ed4` failed only the dedicated `gofmt` check; the one-space formatting mismatch was corrected at the verified SHA above. **Do not relabel that earlier failed run as green.**

## Reproducible raw outputs (GitHub shared runner, `go test` standard microbenchmark)

Command: `go test -run '^$' -bench '^BenchmarkProject(500|5000)Entries$' -benchmem -count=3 -benchtime=100x .`

| Run/event | Workload | ns/op repetitions (three) | B/op | allocs/op |
|---|---|---|---|---|
| 38839381575 push | 500 synthetic entries | 136852, 121425, 129617 | 202296–202297 | 11 |
| 38839381575 push | 5000 synthetic entries | 1604302, 1521162, 1521531 | 1887081–1887085 | 39 |
| 38839385911 PR | 500 synthetic entries | 136511, 171670, 152423 | 202296–202349 | 11 |
| 38839385911 PR | 5000 synthetic entries | 1720142, 1829354, 1810478 | 1887080–1887084 | 39 |

The GitHub dedicated workflow also reported `go test -race` **PASS**, `go vet` **PASS**, `gofmt` **PASS**, and Go runtime statement coverage **100.0%** for the tiny current module. This does **not** establish 100% specification coverage, production security correctness or full capability completion.

## Honest scope and decision

Only a caller-supplied, synthetically authorized, local immutable 500/5000-entry projection is measured. Go `testing.B` reports mean `ns/op`, allocated bytes and allocation count, **not p95 latency**, query/storage network time, peak RSS, deep resource under load, source permission-provider latency, UI response or production SLO. The 5ms engineering p95 registry target and both registered Timeline/deep handlers remain intentionally **PREIMPLEMENTATION / fail closed**. No regression baseline accepted and no storage/indexing/provider technology selected. Results need calibration on representative sources and independent L2 review; `E10-INV-110B-RUNTIME` remains open.

Reuse decision: Go standard-library `testing.B` + `sort.SliceStable` is proportionate for bounded local projection; a custom benchmarking framework or third-party event store/index would introduce unnecessary surface and cannot be justified without source-backed comparative benchmarks for a persistent future architecture.
