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

## Exact-SHA CI verified (2026-10-16)

Code snapshot: `a1d89fb9691b8e0a81175f4cb6752644a87efc7d`, parent `6841d883d8a1e64cf0c4f1cc8382615ffdf309ac`, [PR #36](https://github.com/tobianahillel-afk/cmdr/pull/36).

- [engineering-kernel push #38894589373](https://github.com/tobianahillel-afk/cmdr/actions/runs/38894589373): **SUCCESS**.
- [engineering-kernel PR #38894595290](https://github.com/tobianahillel-afk/cmdr/actions/runs/38894595290): **SUCCESS**.
- [investigation-timeline-runtime push #38894589322](https://github.com/tobianahillel-afk/cmdr/actions/runs/38894589322): **SUCCESS**, with Linux, native macOS and native Windows jobs all successful.
- [investigation-timeline-runtime PR #38894595288](https://github.com/tobianahillel-afk/cmdr/actions/runs/38894595288): **SUCCESS**, with Linux, native macOS and native Windows jobs all successful.

This is **CI evidence only**, not a production approval. A one-shot, second-context, read-only L2 was created for precisely this code SHA and diff `6841d883..a1d89fb`; its report was **not directly available** when writing this document (**L2_EN_ATTENTE**). Any new commit must not overwrite the reviewed SHA.

## Next gating actions

1. Retrieve the independent L2 findings with source provenance, reconcile any reproducible R1/R2 defect and retest rather than assuming approval.
2. Verify an actual principal-bound Case / Shared Timeline Entry / source owner integration before asserting real authorization. Existing mocked checker and resolver are **not** a deployed authority.
3. Measure request-level latency and end-to-end resource behavior on a representative future environment before claiming a production SLO; the current batch-mean synthetic 5,000-entry evidence belongs to another SHA.
4. Only promote `E10-INV-110B-RUNTIME` after its full acceptance evidence is satisfied; do not unlock `E10-INV-110C-NAV` prematurely. `OPEN-013` stays unresolved.

This separate documentation-only checkpoint does **not** change the reviewed code or imply that a new L2 has been conducted.
