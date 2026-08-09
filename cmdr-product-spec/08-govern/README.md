---
id: 08-govern-readme
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-OBJ-005, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015, OPEN-019]
---
# Govern

## Mission

Govern is CMDR's authority, execution-governance and response-outcome product. GOV-1 governs Action Requests through Decision and an exact no-effect Execution Handoff Package. GOV-2 governs Playbook selection, execution planning/readiness, Response Runs, technical-executor coordination, verification, rollback/recovery and canonical Result while preserving source ownership and explicit authority.

## Canonical programme identity

- Parent roadmap: **Delivery Roadmap Phase 4 — Govern**.
- Canonical roadmap id: `roadmap-phase-4-govern`.
- GOV-1: **Action Requests, Policy, Authorities and Decisions** — PASS AFTER POST-PUBLICATION VERIFICATION, 16 capabilities / 432 sections / 96 tables / 180 gates.
- Current execution lot: **GOV-2 — Playbooks, Response Runs, Execution, Verification and Rollback**.
- GOV-2 is an execution lot, not a roadmap phase or Capability Specification Phase.
- `Phase 4C Govern` and `Phase 4D Govern`: **DO NOT EXIST**.
- GOV-3 — Audit Trail, Response Metrics and Govern Closure: **NOT STARTED**.

## Govern ownership after GOV-2 functional specification

Govern owns:
- Action Request processing, Policy/authority/Approval/Decision from GOV-1;
- Response Playbook semantics and selection;
- Playbook-to-Decision compatibility review;
- Execution Plan and parameter binding by reference;
- request-specific target resolution/readiness and execution-time authority reconciliation;
- canonical Response Run and Response Step governance;
- scheduling/control intent and execution coordination;
- runtime normalization/reconciliation relative to the Run;
- error/retry/partial-success/compensation governance;
- Verification Plan, Verification Assessment and residual-risk assessment;
- rollback eligibility/plan, Response Rollback and recovery coordination;
- canonical Result/outcome classification;
- GOV-2 provenance and cross-product handoff composition.

## Consumed, not owned

- Command: Detection, Signal, Alert, Incident, general Work Queue, Task coordination and operational priority.
- Investigate: Case, Hypothesis, Artifact, Evidence, Finding, analyses, forensics, Detection Engineering and Threat Intelligence.
- Studio: Skill, Tool, Tool Call, Workflow, Workflow Version, Automation Agent, Automation Run and Human Gate.
- Platform Settings: users, roles, groups, providers, integrations, tenant/environment administration, credentials, secrets, connections and runtime configuration/health.
- Endpoint Agent / provider runtime owners: technical execution primitives, device/runtime health, technical retry/rollback primitives and raw technical responses.
- Shared: Search, Linking, Versioning, Notifications, Jobs, Activity, Trace, Reporting, Collaboration, Comments, Inspector and generic Recovery mechanisms.

## Modules

1. Response Inbox — GOV-1.
2. Action Center — GOV-1, plus return/review links from GOV-2 reconciliation.
3. Decision Register — GOV-1 Decision/Handoff plus Run/Result projections.
4. Policy Gates — GOV-1 plus execution-time validity inputs.
5. Approvals & Authorities — GOV-1 plus execution/rollback authority rechecks.
6. **Playbooks — GOV-2 CAP-GOV-017..019.**
7. **Runs & Rollback — GOV-2 CAP-GOV-020..033.**
8. Audit Trail — future GOV-3 only.
9. Response Metrics — future GOV-3 only.

## Mandatory execution distinctions

Decision != Execution Handoff Package != Execution Plan != Response Run. Playbook != Workflow. Playbook selection/version compatibility != authorization/readiness. Secret Reference != secret value. Target Reference != Resolved Target; target drift != scope expansion. Run created/scheduled/start-requested != started. Step/technical success != Run or verified success. Automation Run/Tool Call/Job != Response Run. Technical output != canonical Result. Retry != reauthorization. Cancel/compensation != rollback. Rollback Plan != rollback execution; rolled-back/recovered != guaranteed exact restoration. Runtime success != verification success. Result != Decision/Evidence/Finding and never rewrites them.

## Execution safety

Every effectful start/resume/retry/rollback/recovery preserves exact Decision version, exact Playbook version, exact target set, allowed/prohibited scope, conditions, expiry, Approval/Exception validity, current readiness and source-executor permission. Unknown/stale/contradictory state is never silently promoted to success or authorization.

## Secrets and AI

Govern normally consumes Secret References/metadata only; raw secret values never belong in Decision, Playbook, Execution Plan, Response Run, Result, logs or reports. AI is optional and proposal-only. It cannot authorize, start, expand scope, select secret values, retry/rollback silently or declare success without evidence. Manual/deterministic paths remain complete.

## Canonical maps

- `capability-map.md`
- `functional-dependency-map.md`
- `object-consumption-map.md`
- `action-classification.md`
- `automation-and-ai-model.md`
- `cross-product-links.md`
- `execution-boundaries.md`
- `permissions.md`
- `screen-capability-map.md`

## Current counts and stop line

- GOV-1: 16 capabilities / 432 sections / 96 mandatory tables — historical verified PASS.
- GOV-2: 17 capabilities / 459 sections / 102 mandatory tables — functional files published; final PASS requires the fifth traceability commit and remote 190-gate verification.
- Govern cumulative: 33 capabilities / 891 sections / 198 mandatory tables.
- GOV-3 capabilities created in GOV-2: **0**.

Govern capability specification and Delivery Roadmap Phase 4 — Govern remain **PARTIAL** until GOV-3 is separately completed. Documentary capability status never proves product implementation.