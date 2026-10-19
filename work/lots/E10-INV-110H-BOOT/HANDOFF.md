# E10-INV-110H-BOOT — Recoverable product integration handoff

## Exact product/engineering authority

- **Canonical behavior**: `CAP-INV-110` Investigation Timeline — Case selection and read-only navigation; consumer `CAP-INV-101` Case Queue, `INV-TIM-001` and `INV-CAS-001`. Case / Timeline Entry objects remain canonical; Shared owns the Timeline mechanism.
- **Engineering**: task `E10-INV-110H-BOOT`, actual parent `E10-INV-110`. Execution prerequisite `E10-INV-110G-EXACT` is `IMPLEMENTED` only, **not VERIFIED**. `E10-INV-110E` remains `VERIFIED`. No new CAP, no canonical contract edits, no OPEN decision resolved.
- **Temporary code owners**: first-party Go Investigation Timeline Case selector validation (`product-runtime/investigation-timeline/case_http.go`) and operator-local synthetic viewer bootstrap (`cmd/case-timeline-demo/snapshot.go`). This is NOT production IAM or canonical persistent Case storage.
- **Git relation**: branch `feature/e10-inv-110h-safe-demo-case-startup` directly from frozen PR #136 exact HEAD `a0be908a5712ae24ec6c0ee6e791648058f9995b`, **not main**. The PR's Git base is not the engineering parent. No changes to the parent HEAD or previous reviewer snapshots.

## Reproducible material mismatch

The server's first-party picker rejects the entire in-scope provider Case catalog if *any* ID fails `validHTTPCaseSelector` (NUL, CR/LF, invalid UTF-8, blank, over 1024 bytes), **before** applying individual Case reader ACLs. The local CLI's `firstDemoAuthorizedCase` previously selected the lexically smallest Case based only on tenant/environment and `Readers` without checking this same identity invariant. An operator-controlled snapshot with an invalid in-scope identifier could therefore cause the CLI to print an unusable Case URL and ephemeral credential despite the picker refusing its provider catalog. An invalid unreadable in-scope Case was enough to break the selector consistency. Valid leading/trailing spaces are **distinct opaque IDs** and must NOT be trimmed.

## Minimal correction

- The same unmodified selector predicate is exposed as `timeline.ValidHTTPCaseSelector` in its owner package; existing internal `validHTTPCaseSelector` routes to it for compatibility.
- `firstDemoAuthorizedCase` checks every in-scope catalog ID **before** selecting a reader-visible Case. A malformed candidate aborts before `randomBearer` or local loopback listener. Foreign tenant/environment candidates never become selected.
- Existing Go tests prove invalid IDs including NUL, CR/LF, invalid UTF-8, overlong and whitespace-only; an in-scope unreadable malformed candidate is still rejected; a valid opaque ID with surrounding spaces is passed through byte-for-byte and accepted by the authenticated first-party owner; foreign tenant candidate remains isolated.
- No new ACL engine, Case mutation, actual telemetry provider, endpoint daemon, HTTP route, API, paid dependency, or storage. Distinguish a runtime test in a synthetic operator-local viewer from production readiness.

## Correct acceptance

- Exact SHA runtime Go+Node test workflow must pass on Linux, macOS and Windows. Required engineering `CHK-SECURITY-TESTS` still requires at least 90% changed-critical Go coverage. No waiver; parent #136 has required engineering failure at 89.33%. This task remains `IMPLEMENTED` until independent review, all required CI, dependencies and merge approvals permit stronger claims.
- This is a read-only **boundary correction** to Case bootstrap and therefore qualifies for ONE distinct risk-based L2 on the frozen child SHA once stable; no duplicate L2 of #136 test-only revocation unit.
- The Case picker, Timeline and source HTTP handlers remain the authorization authority. The snapshot is untrusted operator-local data. No Evidence status, signed custody or production IAM.

## Next actionable work

1. Inspect exact-HEAD GitHub Actions and address any genuine compile or behavior break; distinguish functional success from required security gate.
2. If mandatory Go coverage remains below 90%, do not promote/fuse this child. Avoid duplicative micro-tests; review uncovered owner boundaries and long stacked PR chain instead.
3. Reconcile PR #36, product-spec PR #2, and new standalone feature PRs before any merge; all predecessor approvals/CI remain independent obligations.

No indispensable truth should depend on ChatGPT conversation context; the manifest, graph, code, tests, PR and checks form the restartable handoff.
