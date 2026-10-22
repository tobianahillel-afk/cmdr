# E10-INV-110B — projection microbenchmark evidence

This increment uses the existing Go standard-library `testing.B` harness rather than introducing a benchmarking dependency or a separate indexing/search/storage engine. Alternatives considered: hand-coded `time.Now` loop (less comparable and easier to mis-measure), external benchmarking tool (unnecessary supply-chain cost for this bounded pure function), provider-backed index/event store (outside read-only B-RUNTIME and not selected). The test includes synthetically authorized fixed-scope 500-entry and 5,000-entry workloads, explicit preflight correctness checks and allocs/op.

The dedicated GitHub Actions workflow executes each benchmark three times with 100 iterations. Interpret results as **in-process synthetic ns/op, B/op, allocs/op on GitHub shared runners**; timings are not p95 and cannot establish network, backend, UI, security-provider latency, production SLO or the provisional 5ms p95 registry budget.

Registered Timeline performance targets and deep-resource handlers remain fail-closed preimplementation; this commit does NOT convert them to accepted benchmarks. A future dedicated implementation must bind an exact real target/handler identity and source SHA to CI-produced evidence before activating p95 or deep load/soak gates.

Security: synthetic fixture decisions are hardcoded for the benchmark but MUST NOT be confused with authoritative runtime permission checks. No production dependencies added. This is a separate measurement increment on top of the security hardening snapshot `c47ad463ea38713f187a0e3097e30b7cc71d6d8d`, preserving that SHA for independent L2 review.
