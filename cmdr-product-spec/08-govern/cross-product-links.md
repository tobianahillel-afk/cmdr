---
id: govern-cross-product-links
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-008, REQ-PROD-013, REQ-PROD-014, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015, OPEN-019]
---
# Cross-product Links and Transitions — Govern GOV-1 + GOV-2

## Transition contract

Every transition preserves tenant, environment where applicable, stable source references, source owner, permission boundaries, restrictions, freshness/version, return origin, errors, authority context and provenance. A link never transfers ownership or permission.

## GOV-1 transitions preserved

| Source | Trigger | Destination | Destination owner | Context transmitted | Return / failure behavior |
|---|---|---|---|---|---|
| Finding | governed action proposed | Action Request | Govern | Case/Finding/Evidence refs, proposed action, target, impact, uncertainty, requester | return to Finding/Case; incomplete remains explicit |
| Incident | governed action proposed | Action Request | Govern | Incident, Service/impact, urgency, targets, alternatives, requester | return to Incident Detail |
| Detection Engineering | production change package submitted | Action Request | Govern | change draft, target, risks, rollout/rollback context, provenance | return to Detection Project |
| Threat Intelligence | governed operational/external action needed | Action Request | Govern | sourced package, restrictions, target/destination, expiry, provenance | return to Intelligence context |
| Cloud/Mobile analysis | response need prepared | Action Request | Govern | candidate context, target, restrictions, source lineage, no-effect handoff | return to analysis session |
| Action Request | received | CAP-GOV-001 Intake | Govern | submitted version + return origin | out-of-scope/incomplete retained without execution |
| Intake | accepted for review | Response Inbox | Govern | request state, blockers, deadlines, assignment | source request retained |
| Response Inbox | reviewer opens request | Action Center | Govern | exact request/version/review context | return to same queue selection |
| Action Request | completeness/context/risk/policy/authority review | CAP-GOV-004..014 | Govern | request, source refs, target/scope, assessments, Policy/authority/Approval context | missing/blocked context returns to owner/request lineage |
| Decision Preparation | decision-ready and authorized | CAP-GOV-015 Decision | Govern | request, Policy, authority, Approvals, risk, conditions, rationale | disposition returned to source product |
| Decision | approve / approve-with-conditions | CAP-GOV-016 Execution Handoff | Govern | exact action/scope/targets/conditions/expiry/rollback/verification requirements | no-effect package for GOV-2 |
| information-required / rejected / deferred | review or Decision disposition | source product | source owner | exact request/Decision refs, rationale/questions/restrictions | no deletion or execution |
| superseded Decision | replacement recorded | Decision Register history | Govern | old/new versions and effective status | old record remains resolvable |

## GOV-2 transitions

| Source | Trigger | Destination | Destination owner | Context transmitted | Return / failure behavior |
|---|---|---|---|---|---|
| Decision + Execution Handoff Package | execution planning begins | CAP-GOV-017 Playbook Selection | Govern | exact Decision/action/targets/scope/conditions/expiry/requirements | no compatible candidate blocks planning |
| Playbook Selection | exact candidate selected | CAP-GOV-018 Compatibility Review | Govern | Playbook/version/dependencies/rationale | incompatible/version drift returns to selection or re-decision |
| Compatibility Review | compatible | CAP-GOV-019 Execution Plan | Govern | exact Playbook/version/limitations/Decision lineage | binding/scope gaps remain explicit |
| Execution Plan | plan prepared | CAP-GOV-020 Target Readiness | Govern | target refs/scope/bindings/executor needs/constraints | ambiguous/drifted/unavailable targets block or require review |
| Target Readiness | current target context available | CAP-GOV-021 Authorization Reconciliation | Govern | resolved targets/readiness/timestamps/drift + Plan/Decision | expired/mismatch/invalid authority returns to GOV-1/review |
| Authorization Reconciliation | authorized-for-run-preparation | CAP-GOV-022 Response Run | Govern | exact pinned Decision/Handoff/Plan/Playbook/readiness/authority | creates Run in preparation only |
| Response Run | schedule/start/control requested | CAP-GOV-023 | Govern | Run/version, exact bounds, current authority/readiness | request != confirmation; denial/expiry retained |
| Run control | effectful step eligible | CAP-GOV-024 Response Step Coordination | Govern | Run/Step/target/action/bindings/conditions | dependency/scope/condition failure blocks dispatch |
| Response Step | bounded execution required | CAP-GOV-025 Executor Handoff | Govern → Studio/Endpoint/provider | Run/Step IDs, exact target/action/scope, references, expiry, correlation | source owner accepts/rejects; no ownership transfer |
| Studio Workflow / Tool / Endpoint / provider | technical status/output emitted | CAP-GOV-026 Runtime Reconciliation | source owner → Govern | source object/ref/status/output/error/timestamps/correlation | raw source retained; unknown/contradiction explicit |
| Runtime Reconciliation | error/timeout/partial/unknown | CAP-GOV-027 Error/Retry | Govern | affected steps/targets/raw refs/current authority state | retries bounded; compensation != rollback |
| Run/Step runtime | verification checkpoint/terminal state | CAP-GOV-028 Verification Plan | Govern | intended outcome, targets, criteria/source requirements | unavailable criteria/sources remain explicit |
| Verification Plan + observations | assessment due | CAP-GOV-029 Verification Assessment | Govern | expected vs observed, source refs, technical outcomes, gaps | verification failure does not auto-rollback |
| Verification/Error state | rollback review required | CAP-GOV-030 Rollback Planning | Govern | original Run/Decision/effects/targets/residual risk | unsupported/unsafe/authority gaps require manual/re-decision |
| Rollback Plan | authorized and ready | CAP-GOV-031 Rollback/Recovery | Govern → technical owner | exact reverse/recovery intent, targets, scope, authority, refs | partial/failure/unknown retained and verified |
| Run + verification + rollback/recovery | outcome review sufficient | CAP-GOV-032 Result | Govern | intended/execution/technical/verification/rollback/residual contexts | no raw output promoted wholesale |
| Result / follow-up need | downstream context required | CAP-GOV-033 Provenance/Handoff | Govern | complete Decision→Result lineage, restrictions and return origin | destination owner accepts/acts explicitly |
| Result handoff | operational outcome | Command | Command | Result/version, Run/Decision refs, exact targets/outcome/residual risk | Command owns any Incident transition |
| Result / anomaly / gap | investigation follow-up | Investigate | Investigate | Run/Result/source refs, anomaly/gap, restrictions | Investigate owns Case/Finding/Evidence changes |
| verified response feedback | detection/TI improvement | Detection Engineering / Threat Intelligence | Investigate | relevant response observations/outcome/limitations | no automatic rule/Intel mutation |
| runtime/config issue | source-owner follow-up | Settings / Studio / Endpoint | source owner | exact runtime/config/workflow refs/error/context | no silent config/Workflow mutation |
| GOV-2 provenance | later audit/metrics | GOV-3 future | Govern | immutable/ref-resolvable chain and conceptual metric inputs | GOV-3 capability remains NOT STARTED |

## Execution and ownership boundaries

- Finding/Evidence qualification remains Investigate.
- Incident/Work Queue coordination remains Command.
- Human Gate/Workflow/Tool/Tool Call/Automation Run remain Studio.
- users/roles/tenants/secrets/providers/integrations/runtime administration remain Platform Settings.
- technical execution primitives/raw outcomes remain Endpoint Agent or provider/runtime owner.
- generic Search/Linking/Notifications/Jobs/Activity/Trace/Reporting/Collaboration/Recovery remain Shared.
- Response Playbook, Execution Plan, Response Run/Step governance, verification, rollback/recovery governance and canonical Result are Govern-owned.
- Workflow != Playbook; Automation Run/Tool Call/Job != Response Run; technical output != Result.
- Result never silently rewrites Decision, Incident, Case, Finding or Evidence.
- external/cross-tenant dissemination remains subject to OPEN-019.
- GOV-3 Audit Trail/Response Metrics are not created in GOV-2.