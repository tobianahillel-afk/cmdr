# E10-INV-110B — canonical source-object authority and alias safety

## Risk and evidence

A prior independent L2 reviewed an **older SHA** and reported a major object-alias ACL concern. Independently re-examining the current implementation identified the same root cause in `memory_provider.go`: the source owner lookup resolved references by `Reference` while `CheckAccess` searched all source records by `CaseID/Kind/ObjectID` and returned `allow` if **any** alias listed the principal as a reader. Two separate references to one canonical object, with different reader lists, could therefore authorize a restricted reference via a permissive alias. The concern is real in the in-memory reference provider, not evidence of a deployed exploit or a production provider.

## Fix and reuse-first decision

Keep existing owner-verified `SourceOwnerIdentity`, `AccessChecker` and canonical permission names. Introduce a small first-party Go map keyed by **Case ID, source kind and canonical object ID**, built once at immutable provider construction. Reject duplicate canonical-object aliases at construction, even when their distinct reference strings would otherwise have passed. `CheckAccess` uses that exact-key record and its reader ACL in O(1) average time instead of scanning every reference O(number of sources) per authorization. A distinct Case may separately own the same kind/object ID; the Case boundary is part of the key.

This is stricter than a production store that explicitly supports aliases: such a store would require an owner-defined alias policy and reference-aware authority contract rather than assuming a permissive alias. Do NOT silently loosen the policy. No new source schema, new provider, migration, third-party search/index or paid service is introduced. This correction applies to the static memory-backed reference only; a future persistent provider must independently enforce its authoritative access and referential constraints.

## Regression

`memory_source_identity_test.go` constructs one restricted evidence object and a second permissive reference to the **same** canonical object; construction must fail. It also verifies that distinct Case-scoped objects remain allowed and that authorized vs denied evidence stays isolated in the integrated HTTP Case timeline. Existing in-process journey tests must still pass. Run exact-SHA GitHub Actions and independent L2 before claiming this risk remediated in a verified release.

Work graph B remains READY; no production SLA, deployed ACL system or full CAP-INV-110 claim.
