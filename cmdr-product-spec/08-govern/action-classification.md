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
# Action Classification — Govern GOV-1

## Principle

Action class describes functional effect and authority, not button appearance. GOV-1 governs whether an action may proceed but executes no target mutation and starts no Response Run.

| Class | GOV-1 meaning | Representative GOV-1 actions | Execution boundary |
|---:|---|---|---|
| 0 | observation | read, search, filter, inspect, compare, view Policy, authority, Approval, Decision, history | local read under permission |
| 1 | bounded no-effect assessment | completeness check, Policy Evaluation, risk/reversibility assessment, comparison, bounded report, no-effect simulation | deterministic/attributed assessment only |
| 2 | reversible governance mutation/preparation | request-context edit, Govern assignment, information request, Exception Candidate, Approval Request, comment, permitted delegation, Decision Draft, pre-final conditions, withdrawal/supersession, Execution Handoff preparation | versioned and reversible; OPEN-013 remains open |
| 3 | authority-bearing governance act | final Approval of a high-risk action, Decision authorizing production action, emergency authorization, active exception allowing high-impact action | may authorize later effect but still performs no target mutation in GOV-1 |
| 4 | destructive/irreversible target effect | none executed in GOV-1 | prohibited in this lot; future Govern/runtime authority required |

## Mandatory distinctions

- a class-3 Approval or Decision is **authority**, not execution;
- `approved` does not mean `executed`;
- preparing a class-3/4 target action is a local class-2 governance preparation until an authorized authority-bearing record is finalized;
- a Policy outcome, risk assessment or AI recommendation cannot silently promote an action class;
- OPEN-013 remains open for default governance/step-up of class-2 mutations.

## Prohibitions

GOV-1 never starts Response Run, invokes rollback, mutates a target, executes Endpoint commands, changes Cloud permissions, revokes credentials, blocks networks, deploys Detection content or deletes provenance.
