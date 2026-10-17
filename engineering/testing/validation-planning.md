# Risk-adaptive validation planning

`cmdr-dev validation-plan` combines the E3 impact report with the strict check catalog.

The plan contains only catalogued check IDs and built-in executor keys. It does not contain arbitrary shell commands.

## Selection rules

1. Start with the `always_on_pr` safety floor.
2. Add checks whose trigger paths directly match changed paths.
3. Add checks whose risk domains match product/security scope introduced by the work manifest.
4. If impact is high-risk, unknown, or contains unknown paths, switch to `strict` and include every mandatory check.
5. Close all prerequisites recursively.
6. Order checks topologically so prerequisites run before dependents.

## Tiers

- `fast`: known low-risk control-plane changes;
- `standard`: known product, coverage, architecture, dependency or security domains without explicit high-risk escalation;
- `strict`: unknown/unowned paths, CI/security/architecture high-risk rules, product-runtime changes, or security-sensitive manifest scope.

Cost units are relative planning weights only: fast=1, standard=5, heavy=20. They are not wall-clock promises.

The current GitHub workflow still executes the complete kernel while E3 is being validated. A separate bounded integration lot should switch CI execution to event-derived plans only after the planner itself is verified.
