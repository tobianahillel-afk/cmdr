---
id: investigate-threat-intelligence-states
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Functional states

| Family | States |
|---|---|
| Intake | draft, incomplete, ready, blocked, insufficient-context, source-required, permission-blocked, out-of-scope, superseded |
| Intelligence Requirement | draft, active, paused, partially-satisfied, satisfied, blocked, obsolete, cancelled, superseded, archived |
| Knowledge Project | draft, ready, active, paused, blocked, partial, completed-foundation, failed, archived, superseded |
| Source Assessment | not-assessed, assessing, reliable-candidate, mixed, unreliable-candidate, credibility-unknown, corroborated, contradicted, stale, disputed, superseded |
| Candidate Knowledge | proposed, under-review, supported, weakly-supported, contradicted, inconclusive, disputed, superseded, withdrawn |
| Sighting | observed, imported, inferred-candidate, partial, contradicted, disputed, superseded, withdrawn |
| Relationship | proposed, under-review, supported, contradicted, ambiguous, disputed, superseded, withdrawn |
| Lifecycle | active-candidate, stale, expiring, expired, revocation-proposed, revoked, superseded, withdrawn, archived, disputed |
| Analysis Handoff | draft, collecting, partial, complete-for-declared-scope, reviewed, handoff-ready, returned, withdrawn, superseded |

These are functional projections. No final object state machine, physical event model or exchange contract is defined.
