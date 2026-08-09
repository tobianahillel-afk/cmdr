---
id: govern-capability-map
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# Govern Capability Map — GOV-1

## Canonical identity

- Parent roadmap: **Delivery Roadmap Phase 4 — Govern**.
- Canonical roadmap id: `roadmap-phase-4-govern`.
- Execution lot: **GOV-1 — Action Requests, Policy, Authorities and Decisions**.
- `GOV-1` is an execution-lot identifier, not a roadmap phase.
- `Phase 4C Govern`: **DOES NOT EXIST**.
- GOV-2 and GOV-3 remain NOT STARTED.

## Allocation

The `CAP-GOV` namespace was audited before allocation and contained no published capability or reserved/recycled identifier. GOV-1 therefore allocates exactly `CAP-GOV-001..016`.

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

## Functional chain

`Action Recommendation / source context → Action Request → Govern Intake → Response Inbox → Action Center → completeness/context/risk → Policy Evaluation → conflict/exception review → Authority Assessment → Approver Eligibility → Approval Request → Approval → Decision Preparation → Decision → Execution Handoff Package → future GOV-2`.

This chain is nominal, not a guarantee that every request uses every step. Policy, authority, emergency and information-return branches remain explicit and provenance-preserving.

## Ownership

Govern owns Action Request processing, Policy Evaluation, policy conflict/exception governance, contextual authority assessment, Approval lifecycle, Decision and the GOV-1 execution-handoff package. It consumes but does not own Command Incident/Work Queue, Investigate Case/Evidence/Finding, Studio Workflow/Human Gate/Automation Run, Settings users/roles/tenants/secrets/administrative authority configuration, Endpoint execution primitives or Shared engines.

## GOV-2/GOV-3 boundary

No GOV-1 capability creates or executes a Response Run, executes a target mutation, performs rollback, qualifies a Result, owns the Audit Trail product surface or defines Response Metrics. Those remain future GOV-2/GOV-3 concerns.

## Expected conformance

If all 16 canonical files conform to the capability template, GOV-1 contributes **432 numbered sections** and **96 mandatory tables**. These counts are verified from the final published files; this map is not itself proof of completion.
