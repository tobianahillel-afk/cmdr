---
id: govern-playbooks
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
---
# Playbooks — GOV-2

## Mission

Own the response-procedure semantics used by Govern between an approved **Execution Handoff Package** and an **Execution Plan**. GOV-2 selects and validates an exact Response Playbook version but does not absorb Studio Workflow ownership or infer execution authority from selection.

## Owned GOV-2 capabilities

- `CAP-GOV-017` — Response Playbook Catalog and Selection.
- `CAP-GOV-018` — Playbook Version, Preconditions and Decision Compatibility Review.

Later GOV-2 capabilities consume the selected compatible version to create an Execution Plan and Response Run.

## Canonical boundary

- **Response Playbook** is Govern-owned response semantics.
- **Studio Workflow** is CMDR Studio-owned orchestration.
- Playbook ≠ Workflow.
- Playbook version ≠ Workflow version.
- Playbook selected ≠ execution authorized.
- compatibility passed ≠ target ready.
- published/deployed Workflow availability is a referenced dependency, not Govern ownership.

A Playbook may reference one or more Studio Workflows, Tools, Endpoint capabilities or other executor capabilities. Those references retain their original owner and permission model.

## Selection contract

Selection must preserve:
- exact Decision and Execution Handoff Package versions;
- exact Playbook candidate/version;
- supported action and target types;
- tenant/environment constraints;
- risk/reversibility metadata;
- rollback and verification support;
- dependencies and declared runtime availability;
- restrictions, deprecation and limitations;
- selection/rejection rationale and provenance.

No candidate/version is silently substituted.

## Version compatibility

A version change after Decision or handoff triggers explicit compatibility re-review. If the change can alter authorized action, target, scope, conditions, rollback or verification requirements, GOV-2 records `re-decision-required` rather than carrying authority forward silently.

Published versions remain resolvable for historical Runs. Deprecation does not delete history.

## Secrets and parameters

Playbook definitions may declare parameter definitions and **Secret Reference requirements**, but never raw secret values. Platform Settings remains owner of secrets, credentials and connections. Parameter binding belongs to CAP-GOV-019.

## Execution boundary

This module performs no target mutation and starts no Response Run. Its output is a selected, compatibility-reviewed Playbook version consumed by CAP-GOV-019. Actual production execution is governed later in GOV-2 and performed by the appropriate technical owner.

## Screen

`GOV-PLB-001` remains the existing Playbooks surface. GOV-2 creates no new Screen ID and does not rewrite its detailed buttons, columns, filters, wireframes, shortcuts or animations. The Screen Capability Map may only add capability links/ownership pointers.

## Dependencies

Decision/Execution Handoff Package, Playbook object, Studio Workflow/Version/Tool, Settings environment/integration/health/Secret Reference metadata, Endpoint capability projections, Shared Search/Versioning/Trace/Linking, Security permissions and OPEN-008/013/015.

## GOV-3 boundary

Audit Trail and Response Metrics consume Playbook/Run provenance later. GOV-2 does not create GOV-3 capability contracts.