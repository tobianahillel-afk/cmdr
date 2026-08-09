---
id: govern-screen-capability-map
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-UX-008, REQ-UX-009, REQ-UX-010, REQ-PROD-015]
open_decisions: [OPEN-010]
---
# Screen Capability Map — Govern

This map links existing screens to functional capabilities. GOV-1 does not rewrite detailed screen specifications, create a Screen ID, define final columns/buttons/filters/animations/shortcuts or turn capabilities into screens.

| Existing screen | Status | GOV-1 primary capabilities | Secondary capabilities / boundary | GOV-1 treatment |
|---|---|---|---|---|
| GOV-INB-001 Response Inbox | active draft | CAP-GOV-001,002,003 | 004,006,009,011 | capability links only; no detailed rewrite |
| GOV-ACT-001 Action Center | active draft | CAP-GOV-004,005,006,014 | 007..013,015,016 | capability links only; no detailed rewrite |
| GOV-DEC-001 Decision Register | active draft | CAP-GOV-015,016 | 003,011,014 | capability links only; no detailed rewrite |
| GOV-POL-001 Policy Gates | active draft | CAP-GOV-007,008 | 004,005,014 | capability links only; no detailed rewrite |
| GOV-AUT-001 Approvals & Authorities | active draft | CAP-GOV-009..013 | 007,008,014,015 | capability links only; no detailed rewrite |
| GOV-PLB-001 Playbooks | active draft | none in GOV-1 | future GOV-2; Workflow remains Studio | read for boundary only |
| GOV-RUN-001 Runs & Rollback | active draft | none in GOV-1 | future GOV-2 Response Run/rollback | read for boundary only |
| GOV-AUD-001 Audit Trail | active draft | none in GOV-1 | future GOV-3; GOV-1 emits provenance | read for boundary only |
| GOV-MET-001 Response Metrics | active draft | none in GOV-1 | future GOV-3; GOV-1 defines conceptual metrics only | read for boundary only |

## Corrections by ownership, not screen rewrite

- Response Inbox is a Govern-only request queue and is not Command Work Queue.
- Policy Gates consumes/represents Policy Evaluation; GOV-1 does not establish a final Policy authoring/activation UI or engine.
- Approvals & Authorities evaluates contextual authority and Approval; Platform Settings remains owner of user/role/group administration and administrative authority configuration.
- Human Gates remain Studio-owned and are not Approval or Decision.
- Playbooks/Runs/Audit/Metrics remain future GOV-2/GOV-3 functional specification.

## Counts for this run

- Govern screens read: **9**.
- Detailed screen specifications rewritten by GOV-1: **0**.
- New Screen IDs: **0**.
