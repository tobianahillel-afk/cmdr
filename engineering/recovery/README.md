# Recovery and multi-agent coordination

This directory is the Git-versioned coordination control plane for autonomous CMDR development.

## Work leases

A work lease is **coordination metadata**, not an authentication credential. Agent IDs are opaque non-secret identifiers.

The policy intentionally allows only one active mutating claim per agent and one active claim per work unit. A claim is bounded to 15–240 minutes; the default is 90 minutes. Expiry is derived from canonical UTC RFC3339 timestamps and does not require a cleanup commit.

The registry starts empty. Agents first run `lease-evaluate` for `acquire`, `renew` or `release`; the command never mutates Git. An allowed proposal is then persisted through an ordinary compare-and-swap Git commit. This keeps coordination auditable and avoids any hidden coordinator service.

Stale/expired claims do not block continuation. Conflicting non-expired claims fail closed. Git/CI reality remains authoritative over this registry; exact reconciliation is implemented in later E7 units.


## Recovery journal

The execution journal is append-only by invariant and protected by a SHA-256 hash chain. Every event binds a sequence, work unit, agent, lease, UTC observation time, exact Git SHA, bounded evidence references and the previous event digest.

The journal is initially empty. Events must occur inside the bound lease interval. Sequence gaps, event mutation, chain rewrites, unknown work units or mismatched lease ownership fail closed.

`resume-checkpoint` is derived data: it compiles the last valid state for one work unit and includes at most five recent event IDs. It never overrides Git or CI reality and is intentionally not persisted as a second source of truth.


## Reconciliation

`recovery-reconcile` compares a derived work-unit checkpoint with the local Git HEAD and the current lease registry. It never calls GitHub or another network service.

CI evidence is optional input. When supplied, it must include the exact commit SHA, positive run ID, push/pull_request event, supported conclusion and canonical UTC observation time as one complete observation. A journal event that says `validation-passed` without matching successful CI evidence is conservatively classified `revalidate`.

Outcomes are `resume`, `revalidate`, `stale-claim` or `conflict`. Diverged Git history and missing lease bindings fail closed as conflicts. An expired/released lease requires reacquisition before mutation. A checkpoint head that is merely an ancestor of the current HEAD requires revalidation instead of silently trusting old evidence.
