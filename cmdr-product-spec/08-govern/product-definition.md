---
id: 08-govern-product-definition
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# Govern Product Definition — GOV-1 Scope

## Product purpose

Govern converts a **request for governed action** into a traceable authority outcome. GOV-1 determines whether a submitted Action Request is reviewable, what Policy applies, what contextual authority and Approvals are required, what Decision is recorded, under what conditions/time bounds, and what exact package may be handed to future execution.

It does not execute the action.

## User problems

Govern addresses five recurring failures:
1. recommendations being mistaken for authority;
2. incomplete or ambiguous targets reaching response execution;
3. Policies, exceptions and authority being applied invisibly or inconsistently;
4. requester/approver conflicts and emergency paths bypassing SoD;
5. Decisions losing source facts, conditions, expiry or provenance before execution.

## Primary users

- Govern Reviewer / Response Governance Analyst;
- Govern Coordinator;
- Policy Reviewer;
- Risk Reviewer;
- Authority Reviewer;
- authorized Approver;
- Decision Maker;
- Emergency Approver / Security Reviewer;
- Auditor as read-only consumer.

Source-product analysts, Incident Commanders, Detection Engineers, Intelligence analysts, Studio operators, Settings administrators and Endpoint operators are contributors/consumers under their own ownership, not implicit Govern decision makers.

## GOV-1 owned capabilities

`CAP-GOV-001..016` cover intake, request queue/lifecycle, scope/target, impact/risk/reversibility, completeness, Policy Evaluation/conflicts/exceptions, authority/eligibility/SoD, Approvals/delegation/emergency, Decision preparation/recording/expiration and execution-handoff provenance.

## Non-goals

GOV-1 does not:
- investigate Case/Evidence/Finding;
- coordinate the general Command Work Queue;
- administer users/roles/groups/secrets/tenants/providers;
- own Studio Workflow/Human Gate/Automation Run;
- define or execute Endpoint/Cloud/network/identity response primitives;
- create Response Run, Result or rollback execution;
- define final object schemas, cardinalities, Policy engine, RBAC/ABAC matrix, API, protocol or implementation;
- rewrite Govern screens in detail;
- start GOV-2 or GOV-3.

## Source-of-truth model

- Govern: Action Request processing, Policy Evaluation, authority governance, Approval, Decision, execution-handoff package.
- Command: Incident/operational coordination.
- Investigate: Case/Evidence/Finding and analytical conclusions.
- Settings/Security: identity/admin configuration and permission/authority policy sources.
- Studio: Tools/Workflows/Runs/Human Gates.
- Endpoint/runtime owner: technical execution.
- Shared: generic engines and cross-product trace/navigation.

## Decision semantics

A Decision is an authority record based on an explicit request version and review context. It is not a recommendation, Approval, Policy output, risk score, AI output, Response Run or Result. Approval can satisfy an authority requirement but does not itself execute or necessarily finalize a Decision.

Decision dispositions supported functionally in GOV-1 are `approve`, `approve-with-conditions`, `reject`, `defer`, `request-more-information`, `cancel` and `supersede`. Final object-state-machine design remains future.

## Delivery classification

All GOV-1 capabilities are `draft` / `defined` / `planned`. Their target may be native, but no current native/integrated/implemented/deployed claim is made by GOV-1.

## Acceptance

A GOV-1 Decision cannot become decision-ready without explicit source/version, target/scope, completeness disposition, applicable Policy outcomes or explicit unknowns, contextual authority, required Approvals/SoD outcome and provenance. Approved authority still creates no target effect until future GOV-2 execution.
