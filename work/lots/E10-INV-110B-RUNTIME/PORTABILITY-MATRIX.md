# E10-INV-110B-RUNTIME — native OS portability gate

## Technical choice

Reuse the existing GitHub Actions `investigation-timeline-runtime` workflow, the first-party Go toolchain and the current `TestTimelineProcess*` suite on **native Windows and macOS** hosted runners. Linux keeps its existing race/coverage/vet/benchmark job. This is more representative than merely cross-compiling Windows/macOS binaries and avoids adding a vendor testing framework, a custom runner or any paid external API.

## Expected checks / tradeoffs

- On Windows/macOS: compile, run all Timeline projection tests and `go vet`, then execute bounded process-start, cancellation, descendant-termination and stdout/stderr overflow tests for the Go development probe.
- Run both platforms independently (`fail-fast: false`) so a failure on one does not erase evidence from the other.
- Uses immutable pinned `actions/checkout` and `actions/setup-go`, read-only repository permissions, no persisted checkout credentials, no secrets, no production infrastructure.
- Extra OS runners consume GitHub Actions resources; their scope is restricted to the existing Timeline runtime workflow paths, with no unrelated project jobs.
- No claim of Windows/macOS success until the **exact GitHub SHA's matrix jobs** reach SUCCESS. A green Linux job alone is not cross-platform verification.
- Cross-platform tests do not establish full OS-version population coverage, crash isolation, production authorization, provider networking, source ownership, actual p95 per request, or the completion of CAP-INV-110. `OPEN-013` remains open. B stays READY until reconciliation and independent review.

## Review focus

The Unix path uses Go `syscall.Setpgid` to terminate its process group; the Windows path uses native `taskkill /T /F /PID` with direct-process fallback. The test must show that a real grandchild cannot outlive cancellation, not merely that the parent returns a timeout. If a platform job fails, preserve logs, fix the smallest underlying incompatibility and rerun before treating the job as covered.
