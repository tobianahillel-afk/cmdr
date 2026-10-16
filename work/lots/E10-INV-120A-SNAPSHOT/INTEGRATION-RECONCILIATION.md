# E10-INV-120A — functional PR41 base reconciliation

## Verified branch reconciliation

Integrate the existing offline operator Case → Timeline → typed source-preview end-to-end scenario from draft PR #41 with the current engineering foundation. Preserve *all* files from the latest canonical integration branch `engineering/autonomous-development-foundation` and reapply the PR41 synthetic two-tenant fixture, actual-loopback integration tests and operator documentation. In particular preserve the newer base `cmd/case-timeline-demo/main.go` with quoted terminal identifiers and broader positive synthetic-principal selection; do NOT restore the older, unquoted PR41 main.go. This is an explicit two-parent merge commit, not a force push or a merge into `main`.

## Chosen architecture

Reuse first-party `loadDemoSnapshot`, `firstDemoAuthorizedCase`, `NewMemoryCaseTimelineHTTPHandlers`, and the existing local analyst viewer. The operator supplies an absolute JSON snapshot path *only at process startup*. The server authenticates loopback requests with a per-run bearer token and rechecks Case/Timeline/source membership via existing in-memory owner handlers. No API, third-party data service, new persistence layer or new RBAC implementation is introduced. A database/index/production provider would add scope and complexity without enabling this offline proof of integration.

## Acceptance and limitations

- Launch `go run ./cmd/case-timeline-demo -snapshot=<absolute path to example-owners.json> -principal=investigator-42`; the terminal prints a localhost Case URL, local analyst viewer URL and ephemeral bearer.
- Actual running HTTP server checks Case discovery, ordered Timeline, typed-source navigation, source preview, anonymous denial, different tenant denial, and clean shutdown.
- Fail closed for malformed owner snapshots, unknown principal, foreign Case references and invalid source memberships, before a listener prints an access token.
- Source fixture is SYNTHETIC. Authorization is from operator-provided in-memory snapshot, **not authoritative IAM**; no live SIEM ingestion, durable data persistence, trusted evidence-chain signing, production network service, or platform-grade audit trail. CAP-INV-110 is **not** end-to-end production-complete; OPEN-013 remains unresolved.
- This document does not claim CI PASS before actual Github Actions logs/results are checked on the resulting exact SHA. Security-critical coverage checks remain mandatory; do not reduce thresholds merely to mark green.
