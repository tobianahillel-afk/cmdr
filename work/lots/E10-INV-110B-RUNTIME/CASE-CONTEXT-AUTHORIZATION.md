# E10-INV-110B — preserve Case identity in owner authorizations

## Increment decision

The new first-party `ProjectWithAccessChecker` adapter checks a Case, then Shared Timeline Entries, then owner-controlled linked sources. Tenant and object identity alone cannot distinguish whether an otherwise readable Timeline Entry really belongs to the selected Case. This bounded increment adds `CaseID` to each trusted `AccessRequest` so the server-side owner checker can validate **Case membership as well as object permission**. The same selected CaseID is passed to all three resource checks, including source-owner obligations.

Reuse-first: `context.Context` and the existing minimal Go `AccessChecker` abstraction are enough for this contract. No new permission vendor, OPA bundle, database, migrations, external API or shared ownership changes. `Project` and canonical Product Spec are unchanged.

## Boundaries and nonclaims

The interface passes a Case context; it **does not** itself verify Case membership. Only an authoritative backend checker can decide whether a Timeline Entry belongs to the Case, and whether its source is visible under the same principal. No checker implementation is bundled for production, and the caller of the adapter must not supply its own checker. The new tests use explicit provider stubs; they prove propagation and fail-closed denial, **not** production authorization. `OPEN-013` remains open. The B lot remains READY pending end-to-end provider integration, independent review and actual CI validation.

The former `AccessRequest` had no CaseID; the new field is a Go source-level interface change within an unshipped bounded runtime. A future server integration must take this required field into account and fail closed on mismatches.
