---
id: 08-govern-readme
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-OBJ-005, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015, OPEN-019]
---
# Govern

## Mission

Govern is CMDR's authority and decision product. It receives governed Action Requests, preserves their source context, evaluates completeness, scope, impact, risk, Policy and contextual authority, obtains required Approvals, records Decisions with conditions and time bounds, and prepares exact execution handoffs. It never converts recommendation or urgency into silent authority.

## Canonical programme identity

- Parent roadmap: **Delivery Roadmap Phase 4 — Govern**.
- Canonical roadmap id: `roadmap-phase-4-govern`.
- Current execution lot: **GOV-1 — Action Requests, Policy, Authorities and Decisions**.
- GOV-1 is not a roadmap phase.
- `Phase 4C Govern`: **DOES NOT EXIST**.
- GOV-2 and GOV-3 are NOT STARTED.

## GOV-1 ownership

Govern owns:
- Action Request processing lifecycle after submission/receipt;
- Govern Response Inbox and Action Center review semantics;
- Policy applicability/evaluation and conflict/exception governance;
- contextual authority assessment and approver eligibility;
- Approval Requests and Approval lifecycle;
- delegation/escalation/emergency governance for authority use;
- Decision Draft preparation, Decision, conditions, expiration and supersession;
- Decision provenance and Execution Handoff Package preparation.

## Consumed, not owned

- Command: Detection, Signal, Alert, Incident, general Work Queue, Task coordination and operational priority.
- Investigate: Case, Hypothesis, Artifact, Evidence, Finding, analyses, forensics, Detection Engineering and Threat Intelligence.
- Studio: Skill, Tool, Tool Call, Workflow, Automation Agent, Automation Run and Human Gate.
- Platform Settings: users, roles, groups, tenant/environment administration, credentials, secrets, integrations/providers and administrative authority configuration.
- Endpoint Agent / runtime owners: technical execution primitives and target state.
- Shared: Search, Linking, Versioning, Notifications, Jobs, Activity, Trace, Reporting, Collaboration, Comments, Inspector and Recovery.

## GOV-1 modules

1. Response Inbox — intake, governed queue and Action Request lifecycle.
2. Action Center — context/scope/target, risk/reversibility, completeness and Decision preparation.
3. Policy Gates — Policy applicability/evaluation, conflict and Exception Candidate assessment.
4. Approvals & Authorities — authority requirements, approver eligibility/SoD, Approvals, delegation and emergency governance.
5. Decision Register — Decision recording, conditions, expiration, supersession and execution handoff provenance.

Modules 6–9 remain active owner surfaces but outside GOV-1 capability creation:
- Playbooks and Runs & Rollback → future GOV-2;
- Audit Trail and Response Metrics → future GOV-3.

## Mandatory distinctions

Action Recommendation ≠ Action Request; Action Request ≠ Decision ≠ Response Run; requester ≠ approver; authority ≠ technical permission; role ≠ authority; Approval Requirement ≠ Approval; Approval Request ≠ Approval; Approval ≠ Decision ≠ execution; Policy Evaluation ≠ Decision; Policy pass ≠ safe; conflict ≠ automatic rejection; Exception Candidate ≠ active exception; Human Gate ≠ Approval/Decision; Automation Run ≠ Response Run; target selected ≠ target verified; AI recommendation/risk score/policy-engine output ≠ Decision; Result ≠ Evidence; audit record ≠ conclusion.

## GOV-2/GOV-3 stop line

GOV-1 creates no Response Run, Result, rollback execution, target mutation, Playbook execution capability, Audit Trail capability or Response Metrics capability. Its final output is an **Execution Handoff Package**, which is not a Response Run and has no effect on the target.

## AI

AI is optional and proposal-only. Manual forms, deterministic checks, matrices, viewers, diffs and human workflows provide all essential functions. No auto-Approval, auto-Decision, self-approval, silent bypass, invented authority or automatic exception is allowed.

## Canonical maps

- `capability-map.md`
- `functional-dependency-map.md`
- `object-consumption-map.md`
- `action-classification.md`
- `automation-and-ai-model.md`
- `cross-product-links.md`
- `permissions.md`
- `screen-capability-map.md`

## Acceptance

GOV-1 may be marked PASS only after all 16 `CAP-GOV-001..016` contracts, 432 sections, 96 mandatory tables, registers, traceability and 180 gates are verified remotely with Command non-regression. Documentary PASS never proves implementation.
