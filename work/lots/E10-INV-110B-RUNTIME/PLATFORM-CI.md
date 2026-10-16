# E10-INV-110B — Windows/macOS process-boundary verification

The Timeline subprocess runner has separate Go implementations: `performance_timeline_process_unix.go` (Linux/macOS process-group SIGKILL) and `performance_timeline_process_windows.go` (best-effort Windows taskkill /T /F /PID). Earlier successful engineering-kernel CI ran on Linux only; it did not test the OS-specific paths.

Reuse the existing, first-party Go tests and authorized, SHA-pinned GitHub Actions workflow `.github/workflows/investigation-timeline-evidence.yml`. Add a two-OS runner matrix for Windows and macOS with actual Go build, go vet, cancellation, descendant, output-bound and projection-probe tests. Cross-compiling alone would not test child-process cancellation. No paid API, product dependency, store/index or new workflow is introduced.

This is synthetic GitHub-runner verification only, not support certification for every OS version, every subprocess session, provider authorization or production SLO. Linux on-demand benchmark evidence remains separate and unchanged. Do not assert any Windows/macOS result until the exact-SHA matrix job conclusions and test logs have been read. Work unit B remains READY pending complete scope and L2, navigation C stays BLOCKED.
