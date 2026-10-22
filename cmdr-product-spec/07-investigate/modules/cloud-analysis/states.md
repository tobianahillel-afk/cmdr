---
id: investigate-cloud-analysis-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-UX-006, REQ-UX-010]
open_decisions: [OPEN-012, OPEN-014]
---
# Cloud Analysis states

## Shared functional vocabulary
`draft`, `incomplete`, `ready`, `active`, `paused`, `partial`, `blocked`, `restricted`, `observed`, `stale`, `under-review`, `disputed`, `withdrawn`, `superseded`, `closed`, `reopened`, `archived`.

## Observation qualifiers
`candidate`, `metadata-only`, `control-plane-only`, `runtime-partial`, `source-unavailable`, `unsupported`, `late`, `duplicate-candidate`, `conflicted`, `estimated`, `unknown`.

## State rules
- `observed` always carries source and observation time;
- `partial` names missing periods, fields, services or permissions;
- `restricted` reveals no protected value;
- `stale` never becomes current by display refresh alone;
- `supported` does not mean certain;
- `closed` does not delete or freeze history;
- `superseded` keeps the prior version and reason;
- destination acceptance does not promote a candidate into the destination object without destination review.

These are functional states, not final physical state machines.
