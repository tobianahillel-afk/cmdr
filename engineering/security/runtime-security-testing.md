# Runtime security test obligations and coverage floors

This control is deliberately strict without pretending that CMDR runtime code already exists.

## Current state

The architecture registry currently contains **zero** `product-runtime` boundaries. Therefore runtime authorization tests, tenant-isolation tests and runtime coverage are reported as **not applicable**, not passed.

`runtime-security-evidence.json` records that state explicitly. It is rejected as soon as a product-runtime boundary exists.

## Runtime activation contract

Every future `product-runtime` boundary must have exactly one entry in `runtime-security-policy.json` with:

- boundary ID;
- non-empty owner;
- one or more security-critical path patterns contained by that boundary;
- an explicit authorization-required boolean and rationale;
- an explicit tenant-isolation-required boolean and rationale.

The default is fail closed: a runtime boundary without a registered security scope is invalid.

If authorization is required, `SEC-AUTH-NEG-001` must already be active and runtime evidence must show its negative tests passed. The same applies to tenant isolation through `SEC-TENANT-ISO-001`.

## Coverage floors

When product runtime exists, evidence is `measured` and must be generated for the exact checked-out Git commit.

- global runtime coverage: **at least 80%**;
- if any security-critical path in a scope changed: changed security-critical coverage for that scope: **at least 90%**.

The audit receives the real normalized Git change set from the adaptive validation engine. A measured evidence file whose `source_commit` differs from `git rev-parse HEAD` is rejected as stale.

Before runtime exists, absence of measurable runtime coverage is an explicit `not-applicable` state. It is never represented as 100%, PASS, or zero coverage.

## Evidence lifecycle

The committed evidence file is currently the documentary N/A marker. Once runtime tests exist, the runtime test harness must generate/replace this file in the CI workspace before `security-test-audit` executes. That later harness is not invented by this lot because there is no executable product runtime yet.
