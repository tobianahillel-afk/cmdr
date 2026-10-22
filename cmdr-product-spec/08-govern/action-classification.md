---
id: govern-action-classification
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical
requirements: [REQ-PROD-004, REQ-PROD-015, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
---
# Action Classification — Govern GOV-1 + GOV-2

## Principle

Action class describes functional effect and authority, not UI appearance. GOV-1 creates authority-bearing records but no target effect. GOV-2 may coordinate authorized execution while preserving exact Decision, target, scope, conditions and executor ownership.

| Class | Govern meaning | Representative actions | Execution boundary |
|---:|---|---|---|
| 0 | observation | read, search, filter, inspect, compare, view Policy/authority/Approval/Decision/Run/Result/history | local read under permission |
| 1 | bounded no-effect assessment | completeness, Policy Evaluation, compatibility/readiness/authorization reconciliation, verification query/comparison, outcome reconciliation | deterministic/attributed no-effect assessment |
| 2 | reversible governance preparation/mutation | request-context edit, assignment, information request, Exception Candidate, Approval Request, Decision Draft, Execution Handoff, Playbook selection, Execution Plan/bindings, Run draft/scheduling, pause/resume request, retry proposal, Verification Plan, Rollback Plan, annotations | versioned/reversible; OPEN-013 remains open |
| 3 | authority-bearing or reversible production effect | final high-risk Approval/Decision; start/stop/cancel of authorized reversible Response Run; approved production effect; governed retry/rollback/recovery when reversible | requires exact authority/reconciliation and technical-owner execution |
| 4 | potentially irreversible/destructive production effect | action or rollback/recovery whose effect is destructive/irreversible or cannot be safely reversed | explicit heightened governance; no command/bypass defined by documentation |

## GOV-1 invariants preserved

- class-3 Approval/Decision is authority, not execution;
- `approved` != `executed`;
- Policy outcome, risk score or AI recommendation does not create authority;
- GOV-1 still stops at Execution Handoff Package.

## GOV-2 execution rules

- Plan/readiness/reconciliation remain no-effect until an effectful control is authorized;
- start requested != started; stop requested != stopped;
- retry never expands target/scope, bypasses Decision expiry or silently substitutes Playbook/target;
- compensation != rollback;
- rollback planning is class 2; effectful rollback/recovery is class 3/4 according to underlying effect and authority;
- a technical executor result cannot downgrade action class or become a canonical Result automatically;
- class 4 handling can be documented functionally but GOV-2 defines no destructive command, exploit, bypass or implementation.

## AI and automation

AI may propose plans, candidates, retry/rollback recommendations and Result drafts. AI never authorizes, starts, retries or rolls back silently. Deterministic automation may perform an effect only under an explicitly applicable canonical governance contract and remains interruptible/auditable.

## OPEN

OPEN-013 remains open for default class-2 governance and effectful step-up policy. OPEN-007 and OPEN-015 remain open for Human Gate and Automation Run/Response Run semantics.