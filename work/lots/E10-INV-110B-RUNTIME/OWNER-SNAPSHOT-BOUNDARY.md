# E10-INV-110B — trusted snapshot owner seam

## Source-backed engineering choice
Reuse existing Go standard-library interfaces, `AccessChecker`, `AccessRequest` and deterministic `Project`. Do **not** build a database, vendor-specific timeline service, custom graph engine or a second permission registry for an in-process 5,000-entry projection.

`ProjectWithAccessChecker` previously checked resource IDs, but its `Input.Entries` metadata (timestamps, observed/inferred classification, provenance and source links) was still provided by the calling layer. A caller can therefore attempt metadata substitution even after successful ID authorization if that layer is not trustworthy.

The new optional `ProjectFromOwnerSnapshots` accepts only caller-selected IDs. It first obtains Case read authorization, then Timeline Entry authorization per ID, and **only after an allow** calls the injected trusted `EntryResolver` to obtain canonical snapshot fields. The returned ID and tenant must match, and a separate linked-source owner check is performed before a projection. A denied Timeline Entry is never resolved or queried for its source; provider errors/cancellation/mismatched ID/tenant fail atomically. Permissions embedded in resolved snapshots are ignored. Negative tests assert the exact sequence and that metadata comes from the resolver, not from caller-side flags.

## Scope limits
No concrete authenticated resolver, ACL backend, source-specific permission mapping, transactional snapshot/authorization consistency, production transport or persistence is delivered. A malicious or misconfigured resolver/checker can still lie; both must be trusted server integrations. `ProjectWithAccessChecker` remains an older caller-snapshot integration seam, not a safe path for arbitrary untrusted metadata. This new interface does not complete real provider integration. CAP-INV-110 / OPEN-013 and the canonical Product Spec remain unchanged.

Do not claim CI or L2 approval until exact-SHA GitHub Actions and independent report are checked. Existing performance benchmarks do **not** measure the new resolver/provider path.
