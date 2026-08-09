---
id: 08-govern-permissions
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-SEC-001, REQ-SEC-002, REQ-PROD-015]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# Functional Permission Needs — Govern GOV-1

## Scope

This document identifies **functional permission needs only**. The canonical security source remains `../14-security-permissions-and-trust/permission-model.md` and its catalog. GOV-1 does not finalize an RBAC/ABAC namespace, atomic permission matrix, step-up policy or SoD algorithm.

Authority is distinct from technical permission. A Principal can have CRUD permission and still lack contextual Decision Authority. A configured role is not itself approval authority.

| Functional need | Capability | Risk | Class | Masking / scope | Possible step-up | SoD | Owner | Future owner phase |
|---|---|---|---:|---|---|---|---|---|
| Govern Inbox read | 001,002 | restricted request visibility | 0 | tenant/env/restricted refs | low/conditional | no | Govern/Security | Permissions |
| Action Request read/review | 001,003,004,006 | sensitive evidence/target context | 0 | permission-aware projections | conditional | reviewer/requester distinction | Govern | Permissions |
| assign/reassign Govern review | 002 | queue ownership change | 2 | same tenant/env | possible | assignee ≠ authority | Govern | Permissions |
| request information | 003,006,014 | disclosure/loop risk | 2 | only authorized source context | possible | reviewer attributed | Govern | Permissions |
| request-context update | 003,004 | scope/target drift | 2 | versioned and scoped | possible | author/reviewer rules future | Govern | Permissions |
| Impact/Risk/Reversibility Assessment create/update | 005 | decision influence | 1/2 | source citations, uncertainty | conditional | reviewer independence possible | Govern | Permissions |
| Policy read | 007,008 | policy confidentiality | 0 | version/scope restricted | conditional | no | Govern/Security | Permissions |
| Policy Evaluation run/read | 007 | governance outcome influence | 1 | request/version scoped | possible | evaluator ≠ decision maker if policy requires | Govern | Permissions/Technique |
| Exception propose | 008 | bypass candidate | 2 | scope/time bound | likely | requester ≠ approving authority | Govern | Permissions |
| Exception review/revoke | 008 | high-impact governance | 2/3 | exact target/scope/duration | likely | independent authority required by policy | Govern/Security | Permissions |
| Authority Context read | 009 | privileged authority metadata | 0 | tenant/env/action scoped | possible | no | Govern/Security | Permissions |
| approver candidates read | 010 | identity/authority exposure | 0 | minimal identity projection | possible | requester relationship visible | Govern/Settings projection | Permissions |
| eligibility review | 010 | incorrect approver selection | 1/2 | authority/scope/expiry visible | likely | explicit SoD evaluation | Govern/Security | Permissions |
| Approval Request create/update/cancel | 011 | authority solicitation | 2 | request/version/requirement scoped | likely | requester cannot satisfy own requirement under SoD | Govern | Permissions |
| Approval read | 011,014,015 | sensitive rationale | 0 | tenant/request scoped | possible | no | Govern | Permissions |
| approve/reject/abstain | 011 | authority-bearing act | 3 where policy requires | exact action/scope/authority | likely/mandatory for high risk | requester self-approval prohibited where SoD applies | Govern/Security | Permissions |
| delegation create/revoke | 012 | authority transfer misuse | 2/3 | scope/duration/restrictions | likely | delegator authority required | Govern/Security + Settings identity | Permissions |
| escalation | 012 | authority routing | 2 | no privilege creation | possible | escalation ≠ Approval | Govern | Permissions |
| emergency request/review | 013 | shortened path/high impact | 2/3 | limited target/scope/duration | mandatory candidate | emergency requester/approver rules | Govern/Security | Permissions |
| Decision Draft create/update | 014 | decision influence | 2 | versioned; not authoritative | possible | author/decision-maker distinction visible | Govern | Permissions |
| Decision review | 014 | authority context exposure | 0/2 | all inputs permission-aware | likely | challenge/reviewer attribution | Govern | Permissions |
| Decision record | 015 | authoritative disposition | 3 when authorizing high-risk effect | exact scope/target/time bound | likely/mandatory | authority + SoD re-evaluated | Govern/Security | Permissions |
| approve/reject/defer Decision | 015 | production authority | 3 | Decision context only | likely/mandatory | Decision Maker cannot bypass required Approvals | Govern | Permissions |
| conditions update before finalization | 014,015 | scope expansion risk | 2 | diff + re-evaluation | likely if material change | may invalidate prior Approval | Govern | Permissions |
| Decision supersede | 015 | authority history | 2/3 | old Decision retained | likely | new authority evaluation | Govern | Permissions |
| Execution Handoff prepare | 016 | future execution scope | 2 | exact approved scope/target | possible | preparation ≠ execution | Govern | GOV-2/Permissions |
| provenance export | 016 | sensitive cross-object disclosure | 0/1 | redact/minimize; export ≠ read | possible | export authority separate | Shared/Security/Govern | Permissions |
| automated recommendation request | 004..014 | model/tool disclosure/influence | 1 | source minimization and provenance | conditional | automation cannot approve | Studio/Settings/Security, consumed by Govern | Studio/Permissions |
| cross-tenant review | 001..016 | tenant isolation breach | 0–3 depending action | denied unless explicit authorized scope | mandatory candidate | explicit authority + SoD | Security/Settings/Govern | Permissions |
| restricted context read | 001..016 | Evidence/secret/policy disclosure | 0 | masking/minimization | possible | no inferred authority | source owner/Security | Permissions |

## Existing permission families

GOV-1 may reference existing families such as `perm.govern.action-request.*`, `perm.govern.approval.*`, `perm.govern.authority.*`, `perm.govern.decision.*` and `perm.govern.policy.*`. Their current breadth is **not** treated as the final atomic permission design. Missing functional granularity is recorded for the future Permissions phase rather than invented here.

## Invariants

- read does not imply export, approve or execute;
- platform administration does not grant response authority;
- authority is re-evaluated against action, target, scope, tenant/environment, expiry and applicable policy;
- Decision/Approval can require fresh authentication without defining the final step-up mechanism;
- no GOV-1 permission starts a Response Run or target mutation.
