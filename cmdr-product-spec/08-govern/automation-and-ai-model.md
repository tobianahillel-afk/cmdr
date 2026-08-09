---
id: govern-automation-and-ai-model
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-PROD-015, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# Automation and AI Model — Govern GOV-1

## Invariant

Every essential GOV-1 function remains usable manually and through explicit deterministic controls without a model provider, Automation Agent or chatbot. Automation is attributed, inspectable, permission-aware and cannot acquire authority from its own output.

| Function | Human | Rule / deterministic | Workflow | Agent / AI | Authority boundary |
|---|---:|---:|---:|---:|---|
| identify missing request fields | yes | yes | yes | suggestion | no authority |
| find candidate Policies | yes | yes | yes | suggestion | Policy applicability remains reviewable |
| summarize impact/risk/reversibility | yes | yes | yes | sourced summary | assessment ≠ Decision |
| propose authority requirement | yes | yes | yes | candidate only | authority source must be explicit |
| propose approver candidates | yes | yes | yes | candidate only | eligibility and actual approver remain separate |
| detect Policy conflicts | yes | yes | yes | candidate/explanation | conflict ≠ automatic rejection |
| draft Exception Candidate | yes | templates/checklists | yes | draft | candidate ≠ active exception |
| draft Decision options/conditions/rationale | yes | templates/diffs | yes | draft | no auto-Decision |
| prepare Execution Handoff summary | yes | structured aggregation | yes | draft | package ≠ Response Run |
| record Approval | yes | validation only | routing possible | prohibited as autonomous approver | eligible human/authorized authority required |
| finalize Decision | yes | validation only | routing possible | prohibited as autonomous decision maker | Decision Maker authority required |

## Allowed AI assistance

AI may propose missing fields, Policy candidates, risk summaries, authority/approver candidates, conflicts, Exception drafts, Decision options, conditions, rationale, differences and handoff summaries. Every suggestion shows source context, producer/model/run where available, uncertainty and human disposition.

## Forbidden automation

- automatic Approval or Decision;
- requester self-approval when SoD applies;
- silent Policy bypass or hidden conflict;
- automatic activation of an exception;
- invented authority or unexplained approver choice;
- automatic permission grants;
- target mutation, Response Run start or rollback start;
- deletion or rewriting of provenance;
- mandatory chatbot or model dependency for an essential workflow.

## No-AI alternatives

Forms, deterministic validation, policy/authority matrices, checklists, tables, diffs, Policy viewers, authority viewers, Decision templates, explicit human workflows and Shared Search/Linking/Trace provide the complete functional path.

## Studio boundary

Studio owns Tool, Tool Call, Workflow, Automation Agent, Automation Run and Human Gate. Govern may consume their attributed projections. Human Gate is not Approval or Decision; Automation Run is not Response Run. OPEN-007 and OPEN-015 remain open.
