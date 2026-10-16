# E10-INV-110B — descendant process cancellation

## Why

A genuinely accessible prior independent L2 review of source SHA `2421f8be0a970a44334efb6babab0fd0cb09cd5c` reported CHANGES_REQUIRED: `exec.CommandContext` terminates the direct Timeline probe, but grandchildren can remain alive (reproduced with a two-second child despite a short context deadline).

## Reuse-first technology comparison

A new supervisor/dependency is unnecessary for this bounded Go benchmark. Reuse existing `exec.CommandContext`; on Linux/macOS, standard-library Unix `syscall.Setpgid` and SIGKILL to the entire process group are used. On Windows, the built-in `taskkill /T /F /PID` is called under a short context, with direct-process kill fallback. No paid service, third-party product dependency, new database, authorization engine or migration.

## Code behavior and evidence limitations

- Before starting an exact-source local Go build or Timeline probe, configure its own Unix process group; cancellation then targets the group, with Go `WaitDelay` limited to 500 ms.
- Windows subtree cancellation is *best effort only* and must not be stated as verified without CI running on Windows.
- A real-process regression test launches a grandchild and checks that it **cannot** write an escape marker after the parent's 350 ms cancellation deadline. The prior single-process timeout tests remain.
- A process creating a new independent Unix process group/session may escape; no operating-system sandbox or total host resource guarantee is claimed.
- This corrects ONLY the benchmark subprocess tree, not the caller-driven authorization provider, open decision `OPEN-013`, production p95, peak RSS or UI/navigation.
- This document is not proof of successful CI or L2 approval. Verify exact-SHA GitHub Actions runs and the independent scheduled review before marking verified.
