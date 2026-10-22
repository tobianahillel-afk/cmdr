# E10-INV-110D-CLOSURE — real HTTP journey and navigation performance scope

## What now runs, and what does not

The product module and analyst demo already implement an authenticated, read-only Case picker, tenant-scoped Timeline with kind/time/search filtering, source navigation draft, separately reauthorized typed-source preview, and return to the Case. The D acceptance test now executes these as ONE sequence through a genuine Go HTTP test server for analyst A and analyst B, including negative anonymous/cross-tenant/denied-source requests. It exercises actual handlers with the existing first-party immutable synthetic owner data, not just independent handler unit tests or a mock that fabricates HTTP bodies.

The navigation Go benchmark measures the **in-process HTTP handler plus Case/Timeline/typed source checks** against the same first-party synthetic owner fixture. It is intentionally separate from the historical projection-only benchmark. The workflow uses `go test -run '^$' -bench '^BenchmarkAuthorizedCaseSourceNavigationHTTP$' -benchmem -count=3 -benchtime=100x .` and identifies ns/op and allocation statistics. No conversion of this result into request-level p95, production SLO or a real remote IAM/provider benchmark is permitted. The on-demand navigation handler registry remains preimplementation until a correctly bound probe/target is measured and reviewed; this command must not be substituted for that policy check.

## Reuse-before-invention

Chosen: `net/http/httptest`, existing authenticated memory-backed owner chain, existing `NewMemoryCaseTimelineHTTPHandlers` and existing GitHub Actions module workflow. Rejected for this bounded closure: new database, application gateway, unauthorized JWT parser, duplicate navigation framework, storage/index engine or testing service.

## Blocking gaps outside the bounded D task

No real identity/ACL provider, tenant-aware persistent owner repository, source ingestion, durable evidence and audit log, global streaming timeline, deployed browser/API or production performance evidence. `OPEN-013` remains open; Shared Timeline Entry mutation is excluded. This is an executable integrated **reference journey**, not delivery of CAP-INV-110 to production.

This change does not promote D to VERIFIED by itself: read exact GitHub CI for its published SHA, assess platform/security checks, obtain and qualify independent L2 findings where available, and record deterministic closure evidence before status promotion.
