---
id: govern-functional-dependency-map
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-006, REQ-PROD-015, REQ-PROD-019, REQ-PROD-020]
open_decisions: [OPEN-007, OPEN-013, OPEN-015, OPEN-019]
---
# Functional Dependency Map — Govern GOV-1

| Capability | Dependency | Type | Failure behavior |
|---|---|---|---|
| CAP-GOV-001 | Action Request source/version/return origin | source | remain incomplete/blocked; never invent context |
| CAP-GOV-001 | tenant/environment and permission context | Settings/Security projection | deny or partial without cross-scope leakage |
| CAP-GOV-002 | CAP-GOV-001/003 | queue/lifecycle | preserve received request; expose unavailable dependency |
| CAP-GOV-002 | Shared Search/Notifications/Assignments/Inspector | shared | degraded queue remains readable where source allows |
| CAP-GOV-003 | Action Request canonical object | object | no lifecycle mutation if current version unavailable |
| CAP-GOV-003 | Shared Versioning/Linking/Trace | shared | no silent overwrite; conflict/retry explicit |
| CAP-GOV-004 | Incident/Case/Finding/Evidence/target projections | cross-product | ambiguous/restricted/stale remains visible and blocks readiness as required |
| CAP-GOV-004 | Settings tenant/environment/identity context | projection | unknown or restricted; never infer authority |
| CAP-GOV-005 | CAP-GOV-004 | context | risk/reversibility cannot claim completeness without target/scope |
| CAP-GOV-005 | future GOV-2 rollback capability context | future boundary | record requirement/unknown; never invent rollback availability |
| CAP-GOV-006 | Action Request + source Evidence/Finding projections | completeness | request information or mark partial; never requalify Evidence/Finding |
| CAP-GOV-007 | Policy/version/scope | Govern/Security | outcome unknown/not-applicable if required source unavailable |
| CAP-GOV-007 | deterministic Policy evaluator future | technique | manual review/checklist remains valid; no implementation assumed |
| CAP-GOV-008 | CAP-GOV-007 | conflict/exception | conflict stays unresolved; no automatic rejection or bypass |
| CAP-GOV-008 | authority/approver context | Govern/Security | Exception Candidate remains non-active |
| CAP-GOV-009 | Security Decision Authority rules | security | missing authority => not decision-ready |
| CAP-GOV-009 | Settings Principal/Role/tenant/environment projections | admin source | configured role alone never proves contextual authority |
| CAP-GOV-010 | CAP-GOV-009 | authority | no eligible approver selected when authority incomplete |
| CAP-GOV-010 | Security SoD/step-up | security | mark ineligible/step-up-required with reason |
| CAP-GOV-011 | CAP-GOV-010 | approval routing | Approval Request blocked if no eligible approver |
| CAP-GOV-011 | Shared Notifications/Collaboration/Trace | shared | preserve request state and allow explicit retry/manual routing |
| CAP-GOV-012 | CAP-GOV-009/010/011 | delegation/escalation | no implicit delegation; retain original authority context |
| CAP-GOV-013 | Security emergency-access/step-up | security | emergency path blocked if authority/expiry/justification unavailable |
| CAP-GOV-014 | CAP-GOV-004..013 | decision inputs | remain not-ready; unresolved questions visible |
| CAP-GOV-014 | Studio AI/Workflow provenance | optional automation | manual/deterministic Decision Draft remains available |
| CAP-GOV-015 | CAP-GOV-014 | Decision preparation | no Decision finalization unless authorized review is complete |
| CAP-GOV-015 | Action Request/Approval/Policy/authority/versioning | governance | preserve exact source snapshots/references and dissent |
| CAP-GOV-016 | CAP-GOV-015 | authoritative Decision | no handoff if Decision rejected/deferred/expired or scope unresolved |
| CAP-GOV-016 | Studio Tool/Tool Call/Automation Run provenance | provenance | package marks missing automation provenance; no execution inferred |
| CAP-GOV-016 | future GOV-2 Response Run | future handoff | package remains prepared/queued-for-future only; no Response Run created |
| CAP-GOV-001..016 | Shared Linking/Trace/Activity/Reporting/Export/Collaboration | shared | owner data preserved; degraded shared service exposed |
| CAP-GOV-001..016 | OPEN-013 | authority policy | no silent default for class-2 step-up/governance |
| CAP-GOV-010..016 | OPEN-007 | Human Gate relation | Human Gate never treated as Approval/Decision |
| CAP-GOV-016 | OPEN-015 | Automation Run/Response Run bridge | package preserves Run refs without converting run type |

## Dependency rules

Dependencies never transfer ownership. GOV-1 may remain functionally defined while implementation dependencies are planned or unresolved. Missing data, stale sources, permission denial or unavailable Shared/Studio/Settings services must produce explicit partial/blocked/unknown behavior rather than a silent default.
