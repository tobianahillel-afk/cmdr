# E10-INV-110B — bounded Timeline diagnostic streams

## Issue and correction

The prior independent L2 of SHA `2421f8be` identified that the Timeline subprocess runner attached unbounded `bytes.Buffer` values to child stdout and stderr. An untrusted/broken probe or compiler could fill memory faster than the child deadline. This independent increment changes **only the Timeline measurement runner**.

Use an existing Go `io.Writer` interface and standard-library `bytes.Buffer` wrapped by `timelineBoundedOutput` with a 64 KiB maximum retained per stream (128 KiB combined). Writes beyond the cap are reported as fully consumed to avoid deadlock, marked `exceeded`, and cause a fail-closed result; truncated diagnostic bodies are never accepted as a valid performance sample. A context timeout takes precedence over overflow to preserve error provenance.

Tests exercise exact-bound and overflow, stdout and stderr oversized child streams, no blocking, no leakage of bytes past the cap. No new third-party package, storage, background process, authorization backend or performance claim is introduced. This is the second, independent remediation after the separate process-tree cancellation SHA `728ac4934c97d3fad8bdaaa3376bbdcbed7e7d63`. The first L2 is pinned to that SHA and must not be treated as reviewing this new output-handling code.

The synthetic performance history (4.70007566 ms p95 of five 200-operation batch means, sampled Go heap 5,147,088 bytes) belongs to GitHub run `38842748892` at older SHA `dae5fe2c`; it is not automatically revalidated by this patch. Further L2, exact-SHA CI and platform-specific validation are still necessary. Work graph B remains READY.
