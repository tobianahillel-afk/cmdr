---
id: unresolved-decisions-govern-gov2-audit
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: decision-audit-addendum
---
# Open Decision Audit — Govern GOV-2

The canonical `unresolved-decisions.md` remains unchanged because GOV-2 discovers no independently resolved decision and requires no new OPEN. **18 decisions remain OPEN; OPEN-009 remains historically resolved.**

## GOV-2 relevant decisions

| OPEN | GOV-2 relevance | GOV-2 treatment | Status after GOV-2 |
|---|---|---|---|
| OPEN-007 | Human Gate vs Govern Approval/Decision and effectful human control | preserve Human Gate as Studio-owned; never treat it as Approval/Decision | OPEN |
| OPEN-008 | actual platform/source/executor support | keep all Playbook/Endpoint/provider/runtime semantics provider-neutral and availability-aware | OPEN |
| OPEN-013 | default governance/step-up of class-2 and effectful actions | classify actions and expose possible step-up/SoD without inventing final policy | OPEN |
| OPEN-015 | Tool Call/Automation Run / Response Run bridge | keep Automation Run/Tool Call distinct and source-owned; use explicit correlation refs | OPEN |
| OPEN-019 | dissemination/releasability/sharing/access | keep external/cross-tenant handoff permission-aware; no external release is implemented | OPEN |
| OPEN-010 | final role/activity density | no detailed screen rewrite or final control density in GOV-2 | OPEN |
| OPEN-014 | final attachment/material relations | provenance uses references only and does not finalize object relationships | OPEN |

Other OPEN decisions remain unaffected and OPEN. GOV-2 creates **0** new OPEN and closes **0** OPEN.

## Not decided by GOV-2

GOV-2 deliberately does not select a Playbook engine, Workflow bridge implementation, provider/runtime, executor protocol, retry algorithm, rollback primitive, verification engine, physical Result schema, event model, final state machine, RBAC/ABAC or GOV-3 audit/metrics architecture.