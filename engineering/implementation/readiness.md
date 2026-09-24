# Capability implementation readiness

E10 compiles implementation readiness from the canonical Product Spec without modifying it.

The compiler emits exactly one record for every active registered capability and uses four states:

- `IMPLEMENTED`: present in the verified runtime implementation ledger.
- `PROPOSED`: Product Spec delivery status is `proposed`; never silently scheduled.
- `BLOCKED`: at least one active OPEN reference on the capability or one explicitly CAP-scoped unsatisfied blocking dependency exists.
- `READY`: defined/planned, not implemented, not proposed and no source-backed blocker is attributable to that capability.

## Dependency attribution

The canonical `00-governance/dependency-register.md` is parsed directly. A dependency becomes a capability blocker only when:

1. its `Dependent` field explicitly names a CAP ID, numeric CAP range, or CAP slash-list;
2. its status is not exactly `active`; and
3. its `Blocking` field explicitly says yes, before-*, blocking, or required.

Blocking dependencies that do not identify a CAP selector are preserved in `unscoped_blocking_dependencies`; they are not guessed onto every capability. This prevents a documentary/global dependency from silently blocking or unblocking hundreds of capabilities without capability-level evidence.

## Open decisions

The capability-register shard `OPEN` column is intersected with the canonical active OPEN decision set. Resolved OPEN references therefore do not remain blockers, while active OPEN references fail closed.

The verified implementation ledger has precedence over READY: a ledger-backed capability is `IMPLEMENTED`, never scheduled again.
