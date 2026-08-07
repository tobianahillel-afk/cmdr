---
id: investigate-mobile-forensics-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-07
source-of-truth: canonical
requirements: [REQ-INV-001, REQ-PROD-020]
open_decisions: [OPEN-011, OPEN-014]
---
# Mobile Forensics Functional States

These states are functional vocabulary only; they do not create final object state machines.

| Concept | States |
|---|---|
| Intake | `draft`, `incomplete`, `ready`, `partial-extraction`, `encrypted`, `locked`, `restricted`, `unsupported`, `source-required`, `permission-blocked`, `out-of-scope`, `superseded` |
| Session | `draft`, `ready`, `active`, `paused`, `blocked`, `partial`, `completed`, `failed`, `archived`, `superseded` |
| Acquisition Context | `declared`, `under-review`, `accepted-with-limitations`, `incomplete`, `inconsistent`, `disputed`, `unsupported`, `superseded` |
| Integrity / Accessibility | `unverified`, `verifying`, `verified`, `partially-verified`, `incomplete`, `corrupted`, `encrypted`, `locked`, `restricted`, `inaccessible`, `disputed` |
| Observation / Candidate | `proposed`, `under-review`, `supported`, `weakly-supported`, `contradicted`, `inconclusive`, `disputed`, `superseded`, `withdrawn` |
| Extraction / Recovery | `proposed`, `queued`, `processing`, `available`, `partial`, `invalid`, `restricted`, `failed`, `cancelled`, `superseded`, `withdrawn-from-use` |

## Invariants
- `locked`, `encrypted` or `inaccessible` does not mean data is absent.
- `partial` never becomes complete by inference.
- `recovered`/`available` does not establish complete original content or attribution.
- `supported` for an observation means sourced support, not qualified Evidence or confirmed Finding.
- `completed` Session means the bounded analytical work closed; it does not assert the source was complete.
- supersession preserves prior versions and provenance.

## Interface behavior
Loading preserves current context; Empty explains which representation/scope is empty; Partial lists missing areas and impact; Error preserves valid observations; Offline is read-only where safe; Permission denied exposes no protected value; Stale shows source/version/time and requires re-evaluation before mutation.
