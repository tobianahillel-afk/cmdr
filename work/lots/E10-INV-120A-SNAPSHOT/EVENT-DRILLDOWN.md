# E10-INV-120A — source-aware correlated event drill-down

## Functional increment

From the already reauthorized two-Case temporal comparison, each displayed correlated Timeline row now offers **“Rechercher cet événement dans Case A/B”**, constrained by the existing 96-character search bound. The click reopens the correct Case and performs a **fresh** request to the first-party `/case-timeline` handler. This lets an investigator proceed from a cross-Case chronological hint to its owning Case and, separately, to the typed-source preview through existing authorization paths. It does not silently open an owner resource or trust the cross-Case correlation as permission.

The shared-source group pivots now reuse the same short in-view callback, avoiding two separate authorization/pivot implementations. Stale comparison state, switched Case, replaced bearer, aborted requests or invalid IDs cause a no-op. The original comparison and owner previews are cleared before querying.

## Reuse-first and remaining limits

Reuse JavaScript DOM buttons and the existing `caseID` change event + `fetchPage` flow, and the Go Case/Timeline/source authorization pipeline. A new endpoint, provider, framework, index, persistence model or additional frontend permission cache would introduce duplication. The search is a **candidate lookup** (not a unique signed deep link); the backend rechecks authorization and can return no result or several matches.

Only the first 20 already authorized correlated Timeline rows are displayed. IDs above the current 96-character search limit remain text-only. No production IAM, live SIEM, durable evidence, forensic chain of custody or transactionally consistent two-Case snapshot is claimed. `OPEN-013` remains unresolved.

Behavioral Node tests demonstrate the event->Case HTTP flow, bearer/no-credentials/no-store, Case revocation/denial, stale revision and oversized ID. Existing tests still verify the shared-source pivot, temporal comparison, report export and Go offline operator journey. Verify exact SHA CI before considering this increment accepted. The larger 120A functional milestone is already under an earlier risk-based independent read-only L2; do not create redundant L2 tasks for this small UX extension.
