# E10-INV-110B — canonical source-reference consistency

## Bounded security issue

The opt-in typed-source boundary already checks Case and Timeline Entry access, validates owner snapshots, resolves every linked source owner identity, then checks source ACLs. It previously accepted two distinct canonical `SourceOwnerIdentity` values for the *same* trusted `SourceRef` in one request if both were individually well-formed. That permits an internally contradictory mapping of one displayed reference to different owner objects/kinds, despite the immutable projection showing a single source string.

## Minimal reuse-first correction

Reuse Go's ordinary comparable `SourceOwnerIdentity` struct and a request-local `map[string]SourceOwnerIdentity`; no new provider, graph engine, cache, database, parsing of client prefixes or external dependency. Each occurrence still invokes the source owner resolver; no stale memoized shortcut is added. Two resolutions of one reference must be identical. A mismatch fails closed before any linked-source ACL call, preserving the existing three-phase validation order.

The tests use a deliberately changing trusted mock resolver and check stable references remain allowed, while changed object ID, canonical kind or Case deny with zero linked-source ACL calls. Note that the Case mismatch also exercises the existing Case-binding gate. Provider identity itself, revocation races after this comparison, external ACL consistency and real deployment remain integration obligations: this per-request check is **not** a transactional provider snapshot.

`OPEN-013` remains unresolved, `E10-INV-110B-RUNTIME` remains READY; do not promote C-NAV. Exact-SHA CI and L2 must be checked before any pass/approval claim.
