---
id: govern-gov2-execution-boundaries
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-008, REQ-PROD-009, REQ-PROD-015, REQ-PROD-016, REQ-PROD-019, REQ-PROD-020, REQ-AI-002, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015]
---
# GOV-2 Execution Boundaries

## Scope

GOV-2 begins with a current **Decision** and **Execution Handoff Package** from GOV-1 and ends with a Govern-owned canonical **Result** plus permission-aware downstream handoffs. It owns response governance and coordination, not every technical executor.

Canonical functional chain:

`Decision → Execution Handoff Package → Playbook Selection → Compatibility Review → Execution Plan → Target Resolution/Readiness → Authorization Reconciliation → Response Run → Execution Coordination → Runtime Reconciliation → Verification → Rollback/Recovery when required → Result → downstream handoff`.

## Govern owns

Response Playbook semantics, Playbook Selection, Execution Plan, Decision-to-execution reconciliation, Response Run and its functional lifecycle, execution coordination, runtime reconciliation, verification governance, rollback/recovery governance, canonical Result, outcome classification and GOV-2 provenance.

## Studio boundary

CMDR Studio retains Skill, Tool, Tool Call, Workflow, Workflow Version, Automation Agent, Human Gate and Automation Run. Govern may route an authorized execution request to a deployed Workflow/Tool and consume its attributed technical output.

- Response Playbook ≠ Workflow.
- Response Run ≠ Automation Run.
- Tool Call ≠ Response Run.
- Human Gate ≠ Govern Approval.
- Studio compensation semantics do not redefine Govern rollback.

`OPEN-007` and `OPEN-015` remain open.

## Endpoint/runtime boundary

Endpoint Agent and other technical owners retain execution primitives, target-side execution, runtime/device health, technical response, technical retry mechanics and technical rollback primitives when supported. Govern supplies exact authorized intent/bounds and reconciles returned status.

- runtime owner ≠ Govern owner;
- executor accepted ≠ Run started;
- executor success ≠ intended outcome achieved;
- technical output ≠ canonical Result.

## Platform Settings boundary

Platform Settings retains providers, integrations, credentials, secrets, connections, tenant/environment configuration, target administration and runtime health. Govern consumes only authorized metadata and Secret References.

Secret handling remains separated into existence, reference, metadata, masked preview, executor resolution, copy, export and use. GOV-2 normally consumes only references/metadata; raw secret values never belong in Decision, Playbook, Execution Plan, Response Run, Result, log or report.

## Command and Investigate boundaries

Command retains Incident, Work Queue, prioritization and operational coordination. Investigate retains Case, Evidence, Finding and analytical/forensic content. Result may be handed back to those products under their own transition rules but cannot silently mutate or requalify their objects.

- Result ≠ Finding.
- Result ≠ Evidence.
- Result does not rewrite Decision or Evidence.

## Shared boundary

Shared Capabilities retains generic Jobs, Notifications, Trace, Activity, Versioning, Reporting, Search, Linking, Collaboration and Recovery mechanisms. A Shared Job may carry background work but is not a Response Run.

## Execution safety invariants

Before an effectful start, GOV-2 preserves and rechecks exact Decision version, exact Playbook version, exact targets and scope, prohibited scope, conditions, time bounds/expiry, Approval/Exception validity, target drift, dependency/readiness state, rollback requirement and verification requirement.

Target drift never expands scope. Retry never expands scope or bypasses expiry. Playbook changes are never silently substituted. Partial success remains explicit. Verification is separate from runtime completion. Rollback requires its own eligibility/authority/preconditions.

## AI boundary

AI may suggest candidates, plans, mappings, anomaly explanations, retry/rollback recommendations and Result drafts. AI never authorizes, starts, expands scope, selects secret values, bypasses expiry/policy/SoD, retries silently, rolls back silently, declares success without evidence or deletes trace.

Essential GOV-2 paths remain available through deterministic checks, forms, matrices, tables, diffs, status views and human-controlled workflows without AI.

## Delivery boundary

These are provider/runtime-neutral functional contracts. GOV-2 selects no API, protocol, command, provider, exploit/bypass, execution engine, storage schema or product code. Final objects, atomic RBAC/ABAC and detailed screen controls remain later-phase work.