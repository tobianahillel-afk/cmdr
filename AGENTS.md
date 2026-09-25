# CMDR autonomous engineering rules

This file contains repository-wide rules for human and AI contributors. It does not replace the canonical product specification.

## 1. Canonical product authority

`cmdr-product-spec/**` is the canonical source of truth for CMDR product behavior, capabilities, requirements, journeys, screens, permissions, trust expectations, implementation contracts and product acceptance.

Engineering documentation MUST NOT silently redefine product behavior. Generated indexes, summaries and work manifests are navigation and execution aids only. If generated state conflicts with the canonical product specification, the product specification wins and generated state must be rebuilt.

Never convert an OPEN decision, proposal or planned capability into an implemented fact by assumption.

## 2. "Continue development" protocol

When asked to continue development, do not ask the user what to work on while a safe READY path exists.

1. Read `AI_START_HERE.md`.
2. Reconcile recorded state with Git/GitHub reality.
3. Select the next READY work unit from the canonical work graph.
4. Compile only the context required by that unit.
5. Verify Definition of Ready and dependencies.
6. Split the unit before implementation if it exceeds its complexity budget.
7. Reuse valid technical decisions; research only when a material decision is missing or invalidated.
8. Implement only within the unit's authorized scope.
9. Run the minimum safe validation selected from risk and impact.
10. Record evidence and leave a deterministic handoff.
11. Recalculate the graph and leave the repository resumable.

No indispensable state may exist only in an AI conversation.

## 3. Product dependency policy

Third-party product functionality is DENY-BY-DEFAULT.

Do not introduce third-party UI/component libraries, authorization/RBAC/ABAC engines, workflow engines, rule engines, indexing/search engines, orchestration frameworks, business-logic frameworks, analytics SDKs or similar product functionality merely for convenience.

CMDR-owned functionality is implemented in CMDR-owned code unless an explicit Trusted Base exception exists.

Security primitives that are dangerous to reimplement (for example cryptographic primitives and TLS) require an explicit, audited Trusted Base decision. Development-only analysis/build/test tools are governed separately and are not product runtime dependencies.

## 4. Evidence-driven technical choices

For non-trivial technical choices, do not choose solely from model memory or familiarity.

- Reuse an accepted CMDR decision if it remains valid.
- Otherwise classify the decision and perform the required current research.
- Performance-critical choices require representative CMDR benchmarks.
- Critical choices require an adversarial attempt to falsify the preferred solution.
- Claims of superiority must be supported by current evidence or reproducible measurement.
- Prefer simplicity when two solutions satisfy the same verified constraints.

## 5. Scope and safety

- No direct changes to `main`.
- No force-push or history rewrite without explicit authorization.
- Do not silently mutate product contracts to simplify implementation.
- Do not widen allowed paths or scope without updating the work unit.
- A locally blocked unit does not stop the project: continue with the next independent READY unit when safe.
- `IMPLEMENTED` is not `VERIFIED`.
- Dependents unlock only from the state defined by their dependency contract, normally `VERIFIED`.

## 6. Required completion properties

A completed work unit must be traceable from product/engineering obligation through work unit, code, tests, evidence, commit and PR. The repository must remain recoverable by a fresh agent with no conversation history.
