# E10-INV-110B — legacy access checker selector budget

## Evidence-backed problem

The newer owner-resolved entrypoint `ProjectFromOwnerSnapshots` prevalidates
the cumulative 8 MiB source selector envelope before it calls an authorization
provider. Its legacy caller-snapshot counterpart `ProjectWithAccessChecker`
bounded scope+CaseID and the eventual `SourceRef`, but previously forwarded
arbitrarily large Timeline Entry IDs to `CheckAccess(timeline-entry,...)`.
The eventual `Project` rejection was **too late** to protect the authority
provider from needless allocation, transport and logging.

## Minimal change (reuse first)

Use the existing `MaxInputBytes` engineering limit, Go string `len`, and
the preflight/ID loop from the modern owner-resolved boundary. Reject
oversized same-tenant selectors *before any authority calls*; foreign or
malformed entries remain opaque and do not consume the legitimate request's
selector budget. The subtract-before-add check avoids integer overflow.
No new permission types, packages, backend, index or paid API are introduced.

Tests exercise one oversized ID, two individually bounded IDs whose cumulative
cost exceeds the envelope, and a foreign opaque ID that cannot affect the
visible result or reach the provider.

## Limitations

This is a legacy adapter hardening and does NOT make a caller-provided Timeline
Entry authoritative. There is still no bundled production principal/tenant/
Case/source permission provider; `OPEN-013` remains unresolved. Existing
synthetic benchmarks are not production SLOs. Treat this commit as unqualified
until its exact-SHA CI and independent review evidence are actually obtained.
