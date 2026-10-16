# E10-INV-120A — reuse existing authenticated analyst Case selector

The branch briefly introduced a CLI `-case` startup selector, but the existing `/analyst` Case picker and server-verified `/case-timeline?case_id=...` route already select any readable Case, with independently enforced per-request Case, Timeline Entry and typed source ACLs. Keep ONE selection path and restore the proven base `main.go` and `snapshot.go`, rather than duplicate authorization checks.

The synthetic operator-owned snapshot now provides TWO Cases for tenant-A investigator-42 and one restricted tenant-B Case for investigator-99. The real loopback integration test opens the non-default tenant-A Case `case-incident-43`, entry `alert-3`, typed source `event-3` and preview; anonymous and foreign-tenant views are rejected. The CLI starts at the first readable Case, but the analyst UI can switch Cases. This is a functional integration of existing reusable modules and a representative example dataset, not a production IAM or durable Case store.

Reuse-first: first-party Go owner provider, existing `net/http` handlers and analyst UI were preferred over parallel `-case` CLI parsing/authorization or a third-party framework/index. The prior `-case` code is retained only in Git history; `CASE-DEEP-LINK.md` marks it superseded, and this document is normative for this PR. No new dependencies, APIs, product entity ownership or changes to OPEN-013.

Required checks: `go test -race ./...` and `go vet` on Timeline module, Linux/macOS/Windows runner tests, mandatory engineering-kernel security/coverage policy, independent risk-based L2 at immutable material snapshot. Do not claim a passing security gate before the exact-SHA GitHub workflow finishes.
