# E10-INV-110B — independent multi-source Case authorization tests

This independent test-only increment exercises the published Case-bound AccessChecker seam in a four-entry bounded Timeline. Source X is readable; an unrelated same-tenant Timeline Entry is denied at Shared ownership (its source must never be queried); a different Timeline Entry is readable but its owner source is missing; a foreign-tenant entry is masked without querying any provider. Assert exact ordered obligations containing CaseID, tenant, environment and owner-specific permissions. The expected output is exactly X with partial status and no unauthorized source data.

A second test induces a source-owner provider error after an earlier allowed entry; the result must be an atomic denial with zero leaked entries, not an already-computed partial result. Reuses existing Go tests and a short first-party stub, no new vendor permission engine or runtime dependencies.

This validates request/response behavior against a test stub only. It does not provide an actual authenticated principal, production Case-membership authority, real sources or audit logs. Earlier immutable snapshots and L2 reports are not overwritten. CI and independent review must be evaluated on their own exact SHA.
