---
id: govern-automation-and-ai-model
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-005, REQ-AI-007, REQ-AI-009, REQ-AI-010, REQ-PROD-015, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015]
---
# Automation and AI Model — Govern GOV-1 + GOV-2

## Invariant

Every essential Govern function remains usable manually and through explicit deterministic controls without a model provider, Automation Agent or chatbot. Automation is attributed, inspectable, interruptible, permission-aware and cannot acquire authority from its own output.

## GOV-1 functions preserved

| Function | Human | Rule / deterministic | Workflow | Agent / AI | Authority boundary |
|---|---:|---:|---:|---:|---|
| identify missing request fields | yes | yes | yes | suggestion | no authority |
| find candidate Policies | yes | yes | yes | suggestion | Policy applicability remains reviewable |
| summarize impact/risk/reversibility | yes | yes | yes | sourced summary | assessment != Decision |
| propose authority/approver candidates | yes | yes | yes | candidate only | authority/eligibility remain explicit |
| detect Policy conflicts / draft Exception Candidate | yes | yes | yes | candidate/draft | no automatic rejection/exception |
| draft Decision options/conditions/rationale | yes | templates/diffs | yes | draft | no auto-Decision |
| prepare Execution Handoff summary | yes | structured aggregation | yes | draft | package != Response Run |
| record Approval/finalize Decision | accountable authority | validation only | routing only | prohibited as autonomous authority | explicit Govern authority required |

## GOV-2 functions

| Function | Human | Rule / deterministic | Workflow | Agent / AI | Authority boundary |
|---|---:|---:|---:|---:|---|
| propose Playbook candidates | yes | catalog/compatibility rules | yes | recommendation | selection != authorization |
| identify Playbook/version compatibility issues | yes | version/scope/condition diff | yes | explanation | no silent substitution |
| draft Execution Plan / parameter mappings | yes | schema/definition checks | yes | draft/suggestion | Secret Reference only; no value selection |
| flag target drift/readiness concerns | yes | identity/version/health rules | yes | explanation | no scope expansion |
| reconcile Decision/Approval/Exception/expiry | yes | deterministic checks | yes | explanation only | AI cannot create authority |
| summarize runtime/error/partial state | yes | correlation/status rules | yes | sourced summary | technical output != Result |
| propose retry/compensation | yes | eligibility/attempt rules | yes | recommendation | no silent retry |
| draft Verification Plan / residual-risk questions | yes | criteria/checklists | yes | draft | verification evidence remains source-backed |
| propose rollback/recovery path | yes | capability/precondition rules | yes | recommendation | no silent rollback/recovery |
| draft Result / handoff summary | yes | consistency/aggregation rules | yes | draft | no success classification without evidence |
| start/pause/stop/cancel effectful Run | accountable authorized path | governance validation | deterministic execution only if explicitly authorized | prohibited as autonomous authority | exact Decision/Run authority required |
| execute retry/rollback/recovery | accountable authorized path | governance validation | only under explicit canonical contract | prohibited as autonomous authority | target/scope/expiry/step-up preserved |

## Allowed AI assistance

AI may propose missing fields, Policy/Playbook candidates, compatibility concerns, risk summaries, authority/approver candidates, parameter mappings, target anomalies, execution plans, failure summaries, retry recommendations, verification plans, residual-risk questions, rollback/recovery recommendations, Result drafts and cross-product handoff summaries. Every suggestion retains producer/model/run, source context, uncertainty where applicable and human/deterministic disposition.

## Forbidden automation

- automatic Approval or Decision;
- requester self-approval where SoD applies;
- silent Policy bypass, hidden conflict or automatic active exception;
- invented authority, unexplained approver or automatic permission grant;
- silent Playbook/version substitution;
- target or scope expansion;
- selection of raw secret values;
- bypass of Decision/Approval/Exception expiry or target drift;
- autonomous Response Run start, effectful retry, rollback or recovery;
- infinite/unbounded retry;
- success classification without required verification evidence;
- rewriting Decision, Evidence or Finding from Result;
- deletion or hiding of provenance;
- mandatory chatbot/model dependency for essential capability.

## Deterministic execution

A deterministic Workflow/executor may perform an effect only when an explicitly applicable canonical governance contract authorizes the exact Decision version, target, scope, conditions, expiry and action. Technical-owner permission is still required. Govern Response Run and Studio Automation Run remain distinct.

## No-AI alternatives

Forms, deterministic validation, policy/authority/compatibility/readiness matrices, target and version diffs, dependency graphs, parameter forms, Secret Reference pickers, status/error maps, retry counters, verification criteria tables, rollback checklists, Result consistency checks and explicit human workflows provide complete GOV-1/GOV-2 operation without AI.

## Studio and runtime boundary

Studio owns Tool, Tool Call, Workflow, Workflow Version, Automation Agent, Automation Run and Human Gate. Endpoint/provider owners retain technical execution primitives and raw outputs. Govern consumes their attributable projections. Human Gate != Approval or Decision; Automation Run/Tool Call != Response Run; raw technical output != canonical Result. OPEN-007, OPEN-008 and OPEN-015 remain open.