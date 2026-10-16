# E10-INV-110B — bounded empty/opaque projection scope envelope

## Finding and reused technology

The first-party `Project` path used the 8 MiB limit for accessible entry strings but did **not** reject a tenant/environment scope that exceeded the budget if there were no visible entries. This differed from the existing trusted entrypoints `ProjectWithAccessChecker` and `ProjectFromOwnerSnapshots`, which already have envelope preflight budgets. Reuse the same `MaxInputBytes` engineering constant and a subtraction-based length check (avoids combined integer overflow); no new dependency, authority framework, database or migration.

The correction enforces the envelope cap unconditionally, before iterating or filtering entries. A denied/foreign entry still cannot consume an otherwise valid scope's **visible-entry** budget. Regression tests cover tenant, environment, combined overflow, empty and denied/missing entry sets, and an exactly-at-boundary empty scope. Untrusted resource strings do not trigger a provider call because this is a pure in-process projection.

## Boundaries

This 8 MiB limit is an engineering safeguard, not a canonical Product Spec limit or validated production SLO. `Project` still accepts caller-supplied permission flags; production authority and Shared ownership remain opt-in adapter responsibilities. `OPEN-013` stays open, Navigation/C remains BLOCKED. No historical benchmark is reassigned to this snapshot. Verify SHA-specific Go race/vet/format/security gates and obtain independent L2 before considering B complete.
