# E10-INV-110B — require real negative-test selection

## Scope

The existing first-party Go security adapter uses `go test -run` to evaluate authorisation and tenant-isolation contract fixtures. Go deliberately exits with status 0 if a pattern selects no subtest. The previous check treated this exit status as a passed negative test even if the named subtest was absent. This increment makes every explicitly selected child alternative **mandatory**, verifying actual `pass` events from the standard `go test -json` stream and rejecting package-level or parent-only PASS.

## Evidence and review limits

- Reuse Go standard library `encoding/json`, `regexp`, `io` and existing security adapter without another runner or dependency.
- Tests cover no match, parent-only match, partial alternative match, skipped/running-but-not-passed tests, malformed JSON, and real Timeline negative fixture selection. The real integration test deliberately substitutes a nonexistent fixture and expects rejection.
- Existing custom Node security adapter remains unchanged; the Go-only selection checks do not attest a real provider permission integration or readiness of the complete capability.
- No canonical Product Spec modifications, no runtime schema or cross-owner writes, no change to `OPEN-013`.
- **Verification status:** publication of this note is not proof that CI passed. Record exact run IDs and findings after querying GitHub. L2 review must independently check that all alternatives actually execute at the immutable commit, not infer success from process exit code.

Technical comparison: Go test's structured JSON events and existing adapter are sufficient for deterministic test-selection evidence. Adding a parallel test framework would increase dependencies and ambiguity without improving the guarantee.
