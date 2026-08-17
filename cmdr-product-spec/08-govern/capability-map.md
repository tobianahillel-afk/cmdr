---
id: govern-capability-map
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015, OPEN-019]
---
# Govern Capability Map — GOV-1 + GOV-2

## Canonical identity

- Parent roadmap: **Delivery Roadmap Phase 4 — Govern**.
- Canonical roadmap id: `roadmap-phase-4-govern`.
- GOV-1: **Action Requests, Policy, Authorities and Decisions** — execution lot, historical PASS after post-publication verification.
- GOV-2: **Playbooks, Response Runs, Execution, Verification and Rollback** — current execution lot.
- GOV-3: **Audit Trail, Response Metrics and Govern Closure** — NOT STARTED.
- GOV-1/GOV-2/GOV-3 are not roadmap phases or Capability Specification phases.
- `Phase 4C Govern` and `Phase 4D Govern`: **DO NOT EXIST**.

## Namespace allocation

The `CAP-GOV` namespace contains exactly `CAP-GOV-001..033`: GOV-1 owns 001..016 and GOV-2 owns 017..033. No identifier is recycled or duplicated. GOV-3 allocates no capability in this run.

| ID | Capability | Canonical module | Primary owner role | Delivery status | Delivery mode |
|---|---|---|---|---|---|
| CAP-GOV-001 | Govern Intake and Preconditions | Response Inbox | Govern Reviewer | defined | planned |
| CAP-GOV-002 | Govern Response Inbox and Request Queue Management | Response Inbox | Govern Coordinator | defined | planned |
| CAP-GOV-003 | Action Request Lifecycle Management | Response Inbox | Govern Reviewer | defined | planned |
| CAP-GOV-004 | Action Request Context, Scope and Target Review | Action Center | Govern Reviewer | defined | planned |
| CAP-GOV-005 | Action Impact, Risk and Reversibility Assessment | Action Center | Risk Reviewer | defined | planned |
| CAP-GOV-006 | Action Request Completeness and Evidence Context Review | Action Center | Govern Reviewer | defined | planned |
| CAP-GOV-007 | Policy Applicability and Evaluation | Policy Gates | Policy Reviewer | defined | planned |
| CAP-GOV-008 | Policy Conflict, Exception and Waiver Assessment | Policy Gates | Policy Reviewer | defined | planned |
| CAP-GOV-009 | Authority Requirement and Authorization Context Assessment | Approvals & Authorities | Authority Reviewer | defined | planned |
| CAP-GOV-010 | Approver Eligibility and Separation-of-Duties Assessment | Approvals & Authorities | Authority Reviewer | defined | planned |
| CAP-GOV-011 | Approval Request and Approval Record Management | Approvals & Authorities | Approver / Approval Coordinator | defined | planned |
| CAP-GOV-012 | Delegation, Substitution and Escalation Governance | Approvals & Authorities | Authority Reviewer | defined | planned |
| CAP-GOV-013 | Emergency and Time-Bounded Approval Governance | Approvals & Authorities | Emergency Approver / Security Reviewer | defined | planned |
| CAP-GOV-014 | Decision Preparation and Review | Action Center | Decision Reviewer | defined | planned |
| CAP-GOV-015 | Decision Recording, Disposition, Conditions and Expiration | Decision Register | Decision Maker | defined | planned |
| CAP-GOV-016 | Govern Decision Provenance and Execution Handoff | Decision Register | Govern Reviewer / Decision Maker | defined | planned |
| CAP-GOV-017 | Response Playbook Catalog and Selection | Playbooks | Response Operator | defined | planned |
| CAP-GOV-018 | Playbook Version, Preconditions and Decision Compatibility Review | Playbooks | Govern Reviewer | defined | planned |
| CAP-GOV-019 | Response Execution Plan Preparation and Parameter Binding | Playbooks | Response Operator | defined | planned |
| CAP-GOV-020 | Target Resolution and Execution Readiness Assessment | Runs & Rollback | Response Operator / Govern Reviewer | defined | planned |
| CAP-GOV-021 | Authorization, Decision and Condition Reconciliation | Runs & Rollback | Govern Reviewer | defined | planned |
| CAP-GOV-022 | Response Run Creation and Lifecycle Management | Runs & Rollback | Response Operator | defined | planned |
| CAP-GOV-023 | Response Run Scheduling, Start, Pause, Stop and Cancellation Control | Runs & Rollback | authorized Response Operator | defined | planned |
| CAP-GOV-024 | Response Step and Action Execution Coordination | Runs & Rollback | Response Operator | defined | planned |
| CAP-GOV-025 | Studio Workflow, Tool and Endpoint Execution Handoff | Runs & Rollback | Response Operator | defined | planned |
| CAP-GOV-026 | Runtime Status, Progress and Technical Outcome Reconciliation | Runs & Rollback | Response Operator | defined | planned |
| CAP-GOV-027 | Execution Error, Retry, Partial Success and Compensation Management | Runs & Rollback | Response Operator / Govern Reviewer | defined | planned |
| CAP-GOV-028 | Verification Plan and Expected Outcome Management | Runs & Rollback | Verification Reviewer | defined | planned |
| CAP-GOV-029 | Post-Execution Verification and Residual Risk Assessment | Runs & Rollback | Verification Reviewer | defined | planned |
| CAP-GOV-030 | Rollback Eligibility, Planning and Preconditions | Runs & Rollback | Rollback Reviewer | defined | planned |
| CAP-GOV-031 | Rollback Execution and Recovery Coordination | Runs & Rollback | authorized Response/Rollback Operator | defined | planned |
| CAP-GOV-032 | Response Result Recording and Outcome Classification | Runs & Rollback | Response Operator / Govern Reviewer | defined | planned |
| CAP-GOV-033 | Response Run Provenance and Cross-Product Handoff | Runs & Rollback | Response Operator / Govern Reviewer | defined | planned |

## Functional chain

GOV-1:
`source Finding/Incident/context → Action Request → intake/review → Policy/Authority/Approval → Decision → Execution Handoff Package`.

GOV-2:
`Decision + Execution Handoff Package → Playbook Selection → compatibility → Execution Plan → target resolution/readiness → authorization reconciliation → Response Run → control/steps → Studio/Endpoint/provider handoff → runtime/error reconciliation → Verification Plan → Verification Assessment → rollback/recovery when required → Result → cross-product handoff`.

Every arrow is permission-aware and provenance-preserving. Selection, readiness, technical completion and AI recommendation are never execution authority or verified success by themselves.

## Ownership

Govern owns Action Request processing, Approval/Decision, Response Playbook semantics, Execution Plan, Response Run/Step governance, verification/rollback/recovery and canonical Result. Studio retains Workflow/Tool/Tool Call/Human Gate/Automation Run. Endpoint/provider owners retain technical execution primitives/raw results. Settings retains identities, providers, integrations, secrets and runtime/tenant/environment administration. Command retains Incident/Work Queue; Investigate retains Case/Evidence/Finding/analysis; Shared retains generic mechanisms.

## GOV-3 boundary

Audit Trail and Response Metrics remain future GOV-3 capabilities. GOV-2 provides provenance and conceptual metric inputs but creates no GOV-3 capability or detailed screen rewrite.

## Counts

- GOV-1: **16 capabilities / 432 sections / 96 mandatory tables**.
- GOV-2: **17 capabilities / 459 sections / 102 mandatory tables**.
- Govern cumulative: **33 capabilities / 891 sections / 198 mandatory tables**.

These counts describe the canonical files; final GOV-2 PASS is only recorded after remote post-publication verification.