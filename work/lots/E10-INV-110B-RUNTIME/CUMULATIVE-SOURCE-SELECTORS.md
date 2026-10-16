# E10-INV-110B — cumulative source authority selector budget

## Confirmed boundary issue
At the prior exact GitHub code SHA `f216001bb3f826b1e0c44bbd1929f694c1891401`, both the owner-resolved `ProjectFromOwnerSnapshots` and the legacy `ProjectWithAccessChecker` rejected a single SourceRef larger than the engineering 8 MiB limit before `CheckAccess(source,...)`. However, multiple individually legal SourceRef strings could cumulatively exceed 8 MiB while still being passed to the external access provider. The pure `Project` function rejected oversized authorized visible input only *after* provider calls.

## Minimal reuse-first correction
A first-party `boundedSourceSelector` helper uses the existing `MaxInputBytes` engineering budget, a subtract-before-add length check and `invalidScope` validation; no new framework, persistence, third-party indexing or authorization provider. Both entrypoints count source selectors cumulatively with tenant, environment and Case ID bytes before each source-authorization call. Denied Timeline Entry and cross-tenant selectors are not charged because they are never sent to the source authority. Existing per-entry behavior remains unchanged for ordinary cases.

Two multi-entry tests verify that the second cumulative oversize source selector **never reaches** the source access provider and that result is atomic `deny`; separate tests ensure a denied Timeline Entry with a large hidden SourceRef is not resolved or budget-charged.

## Limits and evidence
This bound controls what the library submits to the *source authority*, not network bandwidth, policy correctness, or the budget of the first-party owner resolution itself. An explicit trusted, authenticated production implementation of `AccessChecker`/`EntryResolver` remains absent; passing mocks do not establish real Case membership or source ownership. Navigation and `OPEN-013` stay blocked/open. Existing synthetic on-demand benchmarks belong to their original SHA, not automatically to this update.

**Validation pending:** Do not claim CI pass or independent L2 approval unless the exact new commit's GitHub Actions and review context are independently verified. Keep B READY until all its acceptance gates close.
