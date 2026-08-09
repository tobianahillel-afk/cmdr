---
id: govern-cross-product-links
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-008, REQ-PROD-013, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020]
open_decisions: [OPEN-007, OPEN-013, OPEN-015, OPEN-019]
---
# Cross-product Links and Transitions — Govern GOV-1

## Transition contract

Every transition preserves tenant, environment where applicable, stable source references, source owner, permission boundaries, restrictions, freshness/version, return origin, errors, authority context and provenance. A link never transfers ownership or permission.

| Source | Trigger | Destination | Destination owner | Context transmitted | Return / failure behavior |
|---|---|---|---|---|---|
| Finding | governed action proposed | Action Request | Govern | Case/Finding/Evidence refs, proposed action, target, impact, uncertainty, requester | return to Finding/Case; incomplete remains explicit |
| Incident | governed action proposed | Action Request | Govern | Incident, Service/impact, urgency, targets, alternatives, requester | return to Incident Detail |
| Detection Engineering | production change package submitted | Action Request | Govern | change draft, target, risks, rollout/rollback context, provenance | return to Detection Project |
| Threat Intelligence | governed operational/external action needed | Action Request | Govern | sourced package, restrictions, target/destination, expiry, provenance | return to Intelligence context |
| Cloud/Mobile analysis | future response need prepared | Action Request | Govern | candidate context, target, restrictions, source lineage, no-effect handoff | return to analysis session |
| Action Request | received | CAP-GOV-001 Intake | Govern | complete submitted version plus return origin | reject as out-of-scope or retain incomplete without execution |
| Intake | accepted for Govern review | CAP-GOV-002 Response Inbox | Govern | request state, blockers, deadlines, assignment context | source request retained |
| Response Inbox | reviewer opens request | Action Center | Govern | exact request/version and review context | return to same queue selection |
| Action Request | completeness review | CAP-GOV-006 | Govern | required context and source refs | information request to source product if missing |
| Action Request | context/target review | CAP-GOV-004 | Govern | target/scope/tenant/environment/time bounds | ambiguity/blocker returned to request |
| Context Review | risk review required | CAP-GOV-005 | Govern | target/scope plus impact/reversibility sources | assessment retained as non-decision |
| Action Request | policy review | CAP-GOV-007 | Govern | action, target/scope, request version, Policy candidates | outcome pass/warn/block/unknown/not-applicable |
| Policy Evaluation | conflict detected | CAP-GOV-008 | Govern | Policies/versions/outcomes/conflict reason | no automatic rejection |
| Policy Conflict | exception considered | CAP-GOV-008 | Govern | Exception Candidate context, scope, duration, compensating controls | candidate can be withdrawn/rejected/superseded |
| Action Request | authority review | CAP-GOV-009 | Govern | action class, target, impact, tenant/env, requester | missing authority blocks decision readiness |
| Authority Assessment | approver resolution | CAP-GOV-010 | Govern | required authority, scope, requester relation, SoD | candidates marked eligible/ineligible with reason |
| Eligible approver set | approval needed | CAP-GOV-011 | Govern | Approval Requirement, Action Request/version, authority context, deadline | Approval Request remains distinct from Approval |
| Approval Request | authorized human disposition | Approval | Govern | approver, authority, version, justification, conditions | Approval record; no execution |
| Approval / Policy / Exception | review inputs sufficient | CAP-GOV-014 Decision Preparation | Govern | all source records plus unresolved questions | may return for changes/information |
| Decision Preparation | decision-ready and authorized | CAP-GOV-015 Decision | Govern | request, Policy, authority, Approvals, risk, conditions, rationale | disposition returned to source product |
| Decision | approved / approved-with-conditions | CAP-GOV-016 Execution Handoff | Govern | exact scope, targets, conditions, expiry, rollback/verification requirements | package only; future GOV-2 consumes |
| information-required | reviewer requests source detail | source product | source owner | questions, request/version, due context, restricted links | response returns to same Action Request lineage |
| rejected / deferred | Decision disposition | source product | source owner | Decision ref, rationale, conditions/next review where allowed | no source deletion, no execution |
| superseded Decision | replacement recorded | Decision Register history | Govern | old/new version links, rationale, effective status | old record remains resolvable |

## Explicit boundaries

- Finding/Evidence qualification remains Investigate.
- Incident/Work Queue coordination remains Command.
- Human Gate/Workflow/Automation Run remain Studio.
- users/roles/tenants/secrets and administrative configuration remain Platform Settings.
- technical execution primitives remain Endpoint Agent or target owner.
- generic Search/Linking/Notifications/Jobs/Activity/Trace/Reporting/Collaboration remain Shared.
- GOV-1 ends at Execution Handoff Package; future GOV-2 owns Response Run/execution/verification/rollback specification.
