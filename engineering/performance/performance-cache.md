# Performance evidence cache

E6-PERF-001D reuses benchmark evidence only when every identity dimension is still current.

## Source identity

A reusable benchmark result does not use the repository HEAD commit as its cache source identity. The cache file itself is versioned, so using HEAD would make a newly committed cache record stale immediately.

Each performance target derives a SHA-256 source digest from the content and repository-relative paths of Git-tracked files matched by that target's `trigger_paths`.

Therefore:

- benchmark-relevant source changes invalidate the record;
- unrelated work-state or cache metadata does not invalidate it;
- changing `trigger_paths` still invalidates through the target digest;
- a target whose trigger paths match no tracked files fails closed.

The benchmark audit summary keeps the exact HEAD commit separately for run traceability.

## Cache key

A reusable record binds the target source digest, full target definition, metric budgets, environment, workload, Go toolchain and current reusable research/decision evidence basis. The result has a separate evidence digest.

Any mismatch becomes a cache miss.

## Revisit behavior

A cached regression is never reusable. If it is bound to an accepted performance-sensitive engineering decision whose research evidence is still reusable, the cache audit emits a deterministic `benchmark-regression` E5 revisit signal.

Regression evidence without a valid decision binding fails closed.

## Adaptive CI

`CHK-PERFORMANCE-CACHE` runs after deep-performance scheduling and decision freshness. Empty cache plus no performance targets is `not-applicable`; fresh records are `reusable`; stale records are `cache-miss`; regression records require revisit.
