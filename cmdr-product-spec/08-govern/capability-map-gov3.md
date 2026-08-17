---
id: govern-capability-map-gov3
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-09
source-of-truth: canonical-addendum
---
# Govern Capability Map — GOV-3 and Closure Addendum

Parent: **Delivery Roadmap Phase 4 — Govern**, canonical id `roadmap-phase-4-govern`. GOV-1/GOV-2/GOV-3 are execution lots, not roadmap phases. `Phase 4C/4D/4E Govern` do not exist.

## Historical families retained
- GOV-1 `CAP-GOV-001..016`: Action Request, Policy, Authority, Approval, Decision and Execution Handoff — 16 / 432 / 96, historical 180/180 PASS.
- GOV-2 `CAP-GOV-017..033`: Playbooks, planning/readiness, Response Run, runtime, verification, rollback/recovery, Result and provenance — 17 / 459 / 102, historical 190/190 PASS.

## GOV-3 family
| ID | Capability | Module | Primary role | Delivery |
|---|---|---|---|---|
| CAP-GOV-034 | Govern Audit Trail Intake and Event Semantics | Audit Trail | Govern Auditor | defined / planned |
| CAP-GOV-035 | Decision, Approval and Authority Audit Reconstruction | Audit Trail | Govern Auditor | defined / planned |
| CAP-GOV-036 | Response Run, Verification, Rollback and Result Audit Reconstruction | Audit Trail | Govern Auditor | defined / planned |
| CAP-GOV-037 | Audit Completeness, Integrity, Gap and Contradiction Assessment | Audit Trail | Govern Auditor | defined / planned |
| CAP-GOV-038 | Govern Audit Review, Search and Evidence Package Preparation | Audit Trail | Govern Auditor | defined / planned |
| CAP-GOV-039 | Policy, Exception and Emergency Governance Metrics | Response Metrics | Govern Control Reviewer | defined / planned |
| CAP-GOV-040 | Approval, Authority and Separation-of-Duties Metrics | Response Metrics | Govern Control Reviewer | defined / planned |
| CAP-GOV-041 | Decision Flow, Disposition and Timeliness Metrics | Response Metrics | Govern Control Reviewer | defined / planned |
| CAP-GOV-042 | Response Run Execution and Reliability Metrics | Response Metrics | Govern Control Reviewer | defined / planned |
| CAP-GOV-043 | Verification, Rollback and Recovery Metrics | Response Metrics | Govern Control Reviewer | defined / planned |
| CAP-GOV-044 | Response Outcome, Residual Risk and Effectiveness Metrics | Response Metrics | Govern Control Reviewer | defined / planned |
| CAP-GOV-045 | Govern Queue, Ageing and Lifecycle Flow Metrics | Response Metrics | Govern Coordinator / Control Reviewer | defined / planned |
| CAP-GOV-046 | Govern Trend, Comparison and Control Health Assessment | Response Metrics | Govern Control Reviewer | defined / planned |
| CAP-GOV-047 | Govern Continuous Improvement, Closure and Provenance | Response Metrics | Govern Product Lead / QA | defined / planned |

## Full Govern functional chain
`Action Request → Policy/Authority/Approval → Decision → Execution Handoff → Playbook/Execution Plan/Readiness → Response Run/Steps → technical execution refs → Verification → Rollback/Recovery when required → Result → Audit Reconstruction → Govern Metrics/Trends/Control Health → Continuous Improvement Package / documentary closure`.

Every derived audit/metric artifact preserves source ownership and provenance. Audit does not replay execution. Metrics do not become objectives/policies/SLOs automatically. Improvement packages apply no change.

## Counts after GOV-3 functional files
- GOV-3: **14 capabilities / 378 sections / 84 mandatory tables**.
- Govern cumulative: **47 capabilities / 1269 sections / 282 mandatory tables**.
- Global registered capabilities: **317** — 27 Command / 243 Investigate / 47 Govern.
- Delivery status: **315 defined / 2 proposed / 317 planned**.
- Total sections/tables: **8559 / 1902**.

Final GOV-3/parent PASS is not claimed in this construction map until the fifth functional commit is published and the 200 remote-inclusive gates are verified.