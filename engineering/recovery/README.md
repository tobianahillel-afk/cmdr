# Recovery and multi-agent coordination

This directory is the Git-versioned coordination control plane for autonomous CMDR development.

## Work leases

A work lease is **coordination metadata**, not an authentication credential. Agent IDs are opaque non-secret identifiers.

The policy intentionally allows only one active mutating claim per agent and one active claim per work unit. A claim is bounded to 15–240 minutes; the default is 90 minutes. Expiry is derived from canonical UTC RFC3339 timestamps and does not require a cleanup commit.

The registry starts empty. Agents first run `lease-evaluate` for `acquire`, `renew` or `release`; the command never mutates Git. An allowed proposal is then persisted through an ordinary compare-and-swap Git commit. This keeps coordination auditable and avoids any hidden coordinator service.

Stale/expired claims do not block continuation. Conflicting non-expired claims fail closed. Git/CI reality remains authoritative over this registry; exact reconciliation is implemented in later E7 units.
