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

This map links the nine existing Govern screens to functional capabilities. GOV-2 rewrites **0** detailed screen specifications, creates **0** Screen IDs, and defines **0** final buttons, columns, filters, wireframes, animations or shortcuts.

| Existing screen | Status | Primary capabilities | Secondary capabilities / boundary | GOV-2 treatment |
|---|---|---|---|---|
| GOV-INB-001 Response Inbox | active draft | CAP-GOV-001,002,003 | 004,006,009,011 | GOV-1 links preserved; no rewrite |
| GOV-ACT-001 Action Center | active draft | CAP-GOV-004,005,006,014 | 007..013,015,016,021 | GOV-1 links preserved; Run reconciliation may return for review |
| GOV-DEC-001 Decision Register | active draft | CAP-GOV-015,016 | 003,011,014,017..023,032,033 | Decision/Result/Run links only; no control design |
| GOV-POL-001 Policy Gates | active draft | CAP-GOV-007,008 | 004,005,014,021 | policy validity/reconciliation projections only |
| GOV-AUT-001 Approvals & Authorities | active draft | CAP-GOV-009..013 | 007,008,014,015,021,030,031 | authority references/re-review only |
| GOV-PLB-001 Playbooks | active draft | CAP-GOV-017,018,019 | CAP-GOV-020,021; Studio Workflow boundary | capability links only; no detailed rewrite |
| GOV-RUN-001 Runs & Rollback | active draft | CAP-GOV-020..033 | GOV-1 Decision/Handoff; Studio/Endpoint/Settings/Shared projections | capability links only; no detailed rewrite |
| GOV-AUD-001 Audit Trail | active draft | none in GOV-1/GOV-2 | future GOV-3; CAP-GOV-033 emits provenance | boundary/read only |
| GOV-MET-001 Response Metrics | active draft | none in GOV-1/GOV-2 | future GOV-3; GOV-2 emits conceptual metric inputs | boundary/read only |

## Ownership corrections, not screen rewrites

- Response Inbox remains a Govern-only request queue and is not Command Work Queue.
- Playbooks shows Govern Response Playbook semantics; Studio Workflow remains Studio-owned.
- Runs & Rollback shows canonical Response Run/Result and linked technical executor projections; Automation Run, Tool Call, Job and Endpoint technical result remain source-owned and distinct.
- Settings remains owner of secrets, providers, integrations, runtime/tenant/environment configuration.
- Audit Trail and Response Metrics remain future GOV-3 capabilities.

## Counts for GOV-2

- existing Govern screens: **9**;
- primary GOV-2 surfaces: **2** — `GOV-PLB-001`, `GOV-RUN-001`;
- detailed screen specs rewritten: **0**;
- new Screen IDs: **0**;
- wireframes/final controls/columns/filters/animations/shortcuts defined: **0**.