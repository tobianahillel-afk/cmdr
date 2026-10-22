# E10-INV-110B — Legacy adapter three-phase authorization preflight

## Finding and reuse-first approach

The trusted `ProjectFromOwnerSnapshots` path already enforces three stages: check Case and **all** Timeline Entry ACL decisions, resolve/validate all permitted snapshots, then check linked source ACLs. The opt-in **legacy** `ProjectWithAccessChecker` path was still calling a source authority immediately after each Timeline check. When a later Timeline provider returned an error or a later authorized snapshot was invalid, the overall output became `deny` but earlier source identifiers had already been sent to the source authority.

Rather than importing an authorization framework, writing a new provider or duplicating the projection algorithm, reuse the existing `AccessChecker`, `validTimelineEntryBeforeSourceLookup`, `boundedSourceSelector` and Go slices to apply the owner path's ordering to the legacy adapter. The separate Case, Timeline and source rights remain independently checked; all caller-supplied permissions are discarded; foreign or malformed selectors remain opaque. No Product Spec changes, no new dependencies, storage or production authority provider.

## Regression proofs

- A late Timeline ACL **error** yields `deny` with three provider calls (Case, Timeline x, Timeline later), and **zero** source ACL calls.
- Each invalid **later allowed** owner metadata case (unknown kind, empty provenance, wildcard SourceRef, version 0) yields the same zero-source result.
- Valid two-entry requests retain two authorized entries and all five provider obligations, in the order Case, Timeline x, Timeline later, Source x, Source later.
- Earlier fixture tests are updated to reflect the safer ordering; the existing negative permission and tenant isolation tests must still pass.
- This is not an atomic external ACL snapshot, a canonical owner resolver, a real provider integration, or a production SLO. The legacy adapter remains a transitional caller-snapshot path; the owner-resolved adapter is the preferred trusted boundary.

This document describes the intended exact code diff, **not** automatic proof of CI/L2 success. Match commit SHA and GitHub Actions jobs before marking verified. `OPEN-013` remains open and B stays READY pending exit criteria/review.
