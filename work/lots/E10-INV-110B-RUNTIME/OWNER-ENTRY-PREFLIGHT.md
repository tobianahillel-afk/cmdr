# E10-INV-110B — structural validation before source-provider calls

## Reuse-first implementation
The two existing read-only access adapters reuse one small Go standard-library helper to validate the structural invariants *of Timeline Entries whose Timeline Entry owner check already returned allow* before making an external linked-source permission request. No new provider interface, dependency, index, storage layer or authorization algorithm is introduced.

Validated prior to source lookup: known observed/inferred/user-action/decision/response kind; nonempty provenance; version >= 1; nonnegative clock uncertainty; and no self-supersession. Existing ID, tenant and cumulative 8 MiB source-selector bounds remain in place. The final pure projection still performs its own checks. Not-yet-authorized and foreign entries are left opaque. Malformed trusted-resolver data causes atomic denial without contacting the source authority. The legacy caller-snapshot path rejects similarly and remains explicitly **not authoritative**.

Intentionally *not* required: parseable timestamps or fully qualified numeric timezone, since legitimate unknown or uncertain source clocks must remain visible as uncertain, not silently discarded. Real ACL providers must still independently verify authenticated principal, tenant, Case membership, Timeline and linked-source rights. This code does not provide them.

## Verification obligations
New tests require exactly two calls (Case + Timeline Entry, no source) on five malformed allowed entry variants in both adapters. Complementary tests prove an unknown clock still flows through all three checks and returns uncertain. Verify the exact published SHA's Go race/tests, vet, formatting, security gates and L2. Do not extrapolate prior benchmark results or prior L2 reports to this change. The work graph remains READY until exit criteria are satisfied.
