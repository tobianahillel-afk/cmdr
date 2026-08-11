---
id: endpoint-ept5-automation-and-ai-model
domain: 11-endpoint-agent
status: draft
owner: Endpoint Agent Product Lead
updated: 2026-08-11
source-of-truth: canonical
---
# Endpoint EPT-5 — Automation and AI Model

AI is optional. Deterministic/manual paths are mandatory for every essential EPT-5 function.

AI may explain readiness, suggest a source-backed primitive candidate, summarize technical output, explain a verification mismatch or drift, propose escalation and summarize reversal state.

AI never grants permission or authority, creates Approval/Decision, declares an action authorized, executes an effectful primitive autonomously, bypasses Govern, invents target state or verification, converts a technical outcome into Result, hides partial/failure/unknown state, removes provenance or silently chooses rollback/release.

| Function | Deterministic/manual alternative | AI role | Authority |
|---|---|---|---|
| primitive selection | declared-capability catalog + operator | suggestion | none |
| eligibility/readiness | capability/policy/state checks | explanation | none |
| technical verification | criteria + observations | summary | none |
| drift/partial handling | explicit state rules | explanation/escalation suggestion | none |
| reversal candidate | capability/precondition matrix | suggestion | Govern remains required |
| reconciliation handoff | typed refs/checklist | sourced summary | no Result authority |