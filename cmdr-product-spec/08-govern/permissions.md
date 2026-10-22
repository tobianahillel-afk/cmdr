---
id: 08-govern-permissions
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-SEC-001, REQ-SEC-002, REQ-PROD-015]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015]
---
# Functional Permission Needs — Govern GOV-1 + GOV-2

## Scope

This document identifies **functional permission needs only**. The canonical security source remains `../14-security-permissions-and-trust/permission-model.md` and its catalog. GOV-1/GOV-2 do not finalize an RBAC/ABAC namespace, atomic permission matrix, step-up policy or SoD algorithm.

Authority is distinct from technical permission. A Principal can have CRUD/execute permission and still lack contextual Decision or Run authority. A configured role is not itself approval authority.

## GOV-1 permission needs — preserved

| Functional need | Capability | Risk | Class | Masking / scope | Possible step-up | SoD | Owner | Future owner phase |
|---|---|---|---:|---|---|---|---|---|
| Govern Inbox read | 001,002 | restricted request visibility | 0 | tenant/env/restricted refs | low/conditional | no | Govern/Security | Permissions |
| Action Request read/review | 001,003,004,006 | sensitive evidence/target context | 0 | permission-aware projections | conditional | reviewer/requester distinction | Govern | Permissions |
| assign/reassign Govern review | 002 | queue ownership change | 2 | same tenant/env | possible | assignee != authority | Govern | Permissions |
| request information | 003,006,014 | disclosure/loop risk | 2 | authorized source context only | possible | reviewer attributed | Govern | Permissions |
| request-context update | 003,004 | scope/target drift | 2 | versioned and scoped | possible | author/reviewer rules future | Govern | Permissions |
| Impact/Risk/Reversibility Assessment create/update | 005 | Decision influence | 1/2 | source citations/uncertainty | conditional | reviewer independence possible | Govern | Permissions |
| Policy read/evaluate | 007,008 | policy confidentiality/outcome influence | 0/1/2 | request/version/scope | possible | evaluator != Decision Maker if policy requires | Govern/Security | Permissions/Technique |
| Exception propose/review/revoke | 008 | bypass candidate/high impact | 2/3 | exact scope/time bound | likely | independent authority where required | Govern/Security | Permissions |
| Authority Context / approver candidate read | 009,010 | privileged identity/authority data | 0 | minimal tenant/action scope | possible | requester relation visible | Govern/Settings/Security | Permissions |
| eligibility review | 010 | incorrect approver selection | 1/2 | authority/scope/expiry | likely | explicit SoD evaluation | Govern/Security | Permissions |
| Approval Request manage | 011 | authority solicitation | 2 | request/version/requirement | likely | requester cannot satisfy own requirement under SoD | Govern | Permissions |
| Approval read / approve/reject/abstain | 011,014,015 | authority-bearing act | 0/3 | exact action/scope/authority | likely/mandatory high risk | self-approval prohibited where applicable | Govern/Security | Permissions |
| delegation create/revoke / escalation | 012 | authority transfer/routing | 2/3 | scope/duration/restrictions | likely | delegator authority required | Govern/Security + Settings identity | Permissions |
| emergency request/review | 013 | shortened path/high impact | 2/3 | limited target/scope/duration | mandatory candidate | emergency SoD rules | Govern/Security | Permissions |
| Decision Draft/review/record/disposition | 014,015 | authoritative Decision | 0/2/3 | exact request/target/scope/time | likely/mandatory | authority + Approval/SoD preserved | Govern/Security | Permissions |
| conditions update / Decision supersede | 014,015 | scope/authority change | 2/3 | diff + old Decision retained | likely | may invalidate Approval | Govern | Permissions |
| Execution Handoff prepare | 016 | future execution scope | 2 | exact approved bounds | possible | preparation != execution | Govern | GOV-2/Permissions |

## GOV-2 permission needs

| Functional need | Capability | Risk | Class | Masking / scope | Possible step-up | SoD | Owner | Execution owner / future phase |
|---|---|---|---:|---|---|---|---|---|
| Playbook read/select | 017 | wrong procedure selection | 0/2 | tenant/env/action/target metadata | possible | selector != execution authority where required | Govern | Permissions |
| Playbook/version compare | 017,018 | hidden incompatible version | 0/1 | exact version/dependency metadata | low/conditional | no | Govern | Permissions |
| Execution Plan create/update | 019 | execution-intent/scope risk | 2 | exact Decision/Playbook/target scope | likely for material changes | author != final execution authority where required | Govern | Permissions |
| parameter binding | 019 | wrong/sensitive input | 2 | definitions + source refs only | likely for sensitive actions | binding cannot grant authority | Govern | Settings/Studio/Endpoint source contracts |
| restricted parameter / Secret Reference metadata read | 019,025,030,031 | secret disclosure | 0 | metadata/masked only; never value | likely | no | Settings/Security consumed by Govern | Permissions/Technique |
| target resolve | 020 | wrong-target execution | 0/1 | tenant/env/exact approved refs | possible | no | Govern + source owner | Objects/Permissions |
| readiness run/read | 020 | stale/false readiness | 1 | source freshness/health restricted | possible | reviewer independence possible | Govern | Settings/Endpoint/Studio projections |
| authorization reconcile | 021,023,027,030,031 | expired/mismatched authority | 1 | exact Decision/Approval/Exception versions | likely | SoD/authority source preserved | Govern/Security | Permissions |
| Response Run create | 022 | creates execution envelope | 2 | exact pinned Decision/Plan/targets | likely | creator != effect authorizer where required | Govern | Permissions |
| schedule | 023 | delayed effect may outlive authority | 2 | start window/expiry/bounds | likely | schedule != approval | Govern | Shared scheduling mechanism future |
| start | 023 | production effect | 3/4 | exact target/scope/Decision/Playbook | likely/mandatory by risk | required authority/SoD | Govern | Studio/Endpoint/provider technical owner |
| pause/resume | 023 | runtime control/effect continuation | 2/3 | Run/target/current conditions | possible/likely | resume requires revalidation | Govern | technical executor |
| stop/cancel | 023 | effect/control interruption | 3/4 by underlying action | exact Run/affected scope | likely | request != confirmation | Govern | technical executor |
| per-step inspect | 024 | technical/sensitive state | 0 | target/output masking | conditional | no | Govern | source executor data |
| effectful step dispatch | 024,025 | production effect | 3/4 | exact Run/Step/action/target | mandatory candidate high risk | inherited Decision/SoD | Govern | Studio/Endpoint/provider |
| technical output inspect | 025,026,027,029,031,032 | sensitive raw output | 0 | source permission/masking | conditional | no | source owner + Govern projection | Permissions |
| retry request/propose | 027 | duplicate/repeated effect | 2 | same scope/target/attempt bound | likely | proposer != authority where required | Govern | Permissions |
| effectful retry | 027 | repeated production effect | 3/4 | revalidated Decision/readiness/expiry | likely/mandatory | authority/SoD rechecked | Govern | technical executor |
| compensation prepare | 027 | secondary effect | 2 | affected step/target only | possible | compensation != rollback | Govern | Permissions |
| Verification Plan create/update | 028 | biased/missing verification | 2 | source/criterion visibility | possible | reviewer independence possible | Govern | Permissions |
| verification execute/read | 028,029,031 | observation access | 0/1/2 if no-effect query | source-specific masking | conditional | no inferred authority | Govern + source owner | Permissions/Technique |
| Rollback Plan create | 030 | reverse-effect risk | 2 | original Run/exact affected targets | likely | planner != effect authority where required | Govern | Permissions |
| rollback authorize/start | 030,031 | production reverse/destructive effect | 3/4 | exact rollback Plan/target/scope/authority | likely/mandatory | explicit authority + SoD | Govern | Studio/Endpoint/provider |
| recovery coordinate/effect | 031 | recovery may alter target | 2/3/4 | exact fallback scope | likely | separate authority where required | Govern | technical owner |
| Result create/review/close | 032 | operational truth/closure | 2/3 by policy | source refs/verification/uncertainty | possible | reviewer/finalizer separation possible | Govern | Permissions |
| cross-product handoff | 033 | disclosure/downstream mutation risk | 2 | destination/minimum context | likely for sensitive/cross-tenant | destination owner retains write authority | Govern | Command/Investigate/Studio/Settings/etc. |
| provenance export | 016,033 | sensitive cross-object disclosure | 0/1/2 | redact/minimize; export != read | possible | export authority separate | Shared/Security/Govern | Permissions |
| automated recommendation request | 004..014,017..033 | model/tool disclosure/influence | 1 | minimization/provenance | conditional | automation cannot authorize/execute | Studio/Settings/Security consumed by Govern | Studio/Permissions |
| cross-tenant review/execution context | 001..033 | tenant isolation breach | 0–4 | denied unless explicit authorized scope | mandatory candidate | explicit authority + SoD | Security/Settings/Govern | Permissions |
| restricted context read | 001..033 | Evidence/secret/policy/technical disclosure | 0 | masking/minimization | possible | no inferred authority | source owner/Security | Permissions |

## Existing permission families

Govern may reference existing `perm.govern.*` families including action-request, approval, authority, decision, policy, playbook, response-run, response-step, response-rollback and result. Their current breadth is **not** the final atomic permission design. Functional needs above expose future decomposition rather than inventing RBAC/ABAC.

## Invariants

- read != export != approve != execute;
- platform administration does not grant response authority;
- Secret Reference metadata access != reveal/copy/export/use of secret value;
- Decision/Approval validity does not by itself prove current target readiness;
- Response Run create/schedule != start;
- effectful start/retry/rollback/recovery preserves exact Decision, target, scope, conditions and expiry;
- technical executor permission does not replace Govern authority, and Govern authority does not grant technical executor permission;
- Result finalization does not grant permission to mutate Incident, Case, Finding, Evidence, Workflow or Settings;
- final RBAC/ABAC, step-up and SoD algorithms remain future Permissions work.