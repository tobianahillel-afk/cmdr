---
id: phase-4a-capability-template-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-04
source-of-truth: validation-report
requirements:
  - REQ-PROD-006
  - REQ-PROD-013
  - REQ-PROD-019
  - REQ-OBJ-001
  - REQ-OBJ-012
---

# Phase 4A capability template conformance

## Scope

This report closes the documentation-only corrective gate for the 27 published `CAP-CMD-*` specifications. It does not validate implementation, APIs, protocols, detailed objects, permissions, screens or Phase 4B.

## Initial audit

- capabilities audited: 27;
- required sections: 729 (`27 × 27`);
- mandatory canonical tables expected: 162 (`27 × 6`);
- exact canonical tables initially present: 39;
- table sections requiring normalization: 123;
- initial prose-only or abbreviated-header failures: 123;
- Capability IDs changed: 0;
- Requirement IDs changed: 0;
- OPEN decisions changed: 0;
- owners or delivery classifications changed: 0.

## Initial defects and corrections

| Capability ID | Canonical file | Initial failure | Correction | Final proof |
|---|---|---|---|---|
| CAP-CMD-001 | `06-command/modules/mission-control/capabilities/situation-overview.md` | S13 used the extended authority matrix only | added the mandatory six-column S13 matrix | six canonical tables and 27 sections |
| CAP-CMD-002 | `06-command/modules/mission-control/capabilities/priority-management.md` | S13 used the extended authority matrix only | added the mandatory six-column S13 matrix | six canonical tables and 27 sections |
| CAP-CMD-003 | `06-command/modules/mission-control/capabilities/situation-timeline.md` | S10 was prose; S13 non-canonical | added explicit no-mutation/job rows and canonical S13 | six canonical tables and 27 sections |
| CAP-CMD-004 | `06-command/modules/mission-control/capabilities/handover.md` | S13 non-canonical | normalized manual/deterministic/automation/AI alternatives | six canonical tables and 27 sections |
| CAP-CMD-005 | `06-command/modules/mission-control/capabilities/operational-blockers.md` | S13 non-canonical | normalized blocker-specific automation rows | six canonical tables and 27 sections |
| CAP-CMD-006 | `06-command/modules/mission-control/capabilities/recent-results-and-outcomes.md` | S13 non-canonical | normalized Result-consumption automation rows | six canonical tables and 27 sections |
| CAP-CMD-101 | `06-command/modules/incidents-and-work-queue/capabilities/unified-work-queue.md` | S13 non-canonical | normalized view/filter/next-action rows | six canonical tables and 27 sections |
| CAP-CMD-102 | `06-command/modules/incidents-and-work-queue/capabilities/work-assignment.md` | S13 non-canonical | normalized candidate/proposal/mutation/conflict rows | six canonical tables and 27 sections |
| CAP-CMD-103 | `06-command/modules/incidents-and-work-queue/capabilities/operational-ownership.md` | S8/S9/S10/S13/S16/S17 abbreviated headers | replaced all six with exact canonical headers and ownership-specific rows | six canonical tables and 27 sections |
| CAP-CMD-104 | `06-command/modules/incidents-and-work-queue/capabilities/priority-and-severity-coordination.md` | six abbreviated/non-canonical tables | normalized dimensions, objects, mutation, automation, output and transitions | six canonical tables and 27 sections |
| CAP-CMD-105 | `06-command/modules/incidents-and-work-queue/capabilities/sla-tracking.md` | six abbreviated/non-canonical tables | normalized policy, calculation, pause, output and transition rows | six canonical tables and 27 sections |
| CAP-CMD-106 | `06-command/modules/incidents-and-work-queue/capabilities/incident-coordination.md` | six abbreviated/non-canonical tables | normalized Incident, Case, Govern and Result boundaries | six canonical tables and 27 sections |
| CAP-CMD-107 | `06-command/modules/incidents-and-work-queue/capabilities/task-coordination.md` | six abbreviated/non-canonical tables | normalized Task/source/dependency behavior | six canonical tables and 27 sections |
| CAP-CMD-108 | `06-command/modules/incidents-and-work-queue/capabilities/bulk-coordination.md` | six abbreviated/non-canonical tables | normalized preview, per-item mutation, outputs and retries | six canonical tables and 27 sections |
| CAP-CMD-109 | `06-command/modules/incidents-and-work-queue/capabilities/work-freshness-and-staleness.md` | six abbreviated/non-canonical tables | normalized metadata, health, follow-up and refresh transitions | six canonical tables and 27 sections |
| CAP-CMD-110 | `06-command/modules/incidents-and-work-queue/capabilities/escalation.md` | six abbreviated/non-canonical tables | normalized escalation package and destination ownership | six canonical tables and 27 sections |
| CAP-CMD-201 | `06-command/modules/risk-and-coverage/capabilities/service-context.md` | S8/S9/S10/S13/S16/S17 prose or abbreviated | normalized Service projection without creating `service.md` | six canonical tables and 27 sections |
| CAP-CMD-202 | `06-command/modules/risk-and-coverage/capabilities/exposure-overview.md` | six non-canonical tables | normalized external Exposure projection without creating `exposure.md` | six canonical tables and 27 sections |
| CAP-CMD-203 | `06-command/modules/risk-and-coverage/capabilities/coverage-overview.md` | six non-canonical tables | normalized five coverage families and source ownership | six canonical tables and 27 sections |
| CAP-CMD-204 | `06-command/modules/risk-and-coverage/capabilities/business-impact-context.md` | six sections were prose only | added specific impact/source/certainty tables | six canonical tables and 27 sections |
| CAP-CMD-205 | `06-command/modules/risk-and-coverage/capabilities/risk-prioritization-context.md` | six sections were prose only | added factor/proposal/disposition tables | six canonical tables and 27 sections |
| CAP-CMD-301 | `06-command/modules/readiness-and-operations/capabilities/readiness-overview.md` | S8 abbreviated; S9/S10/S13/S16/S17 prose | normalized assessment/evidence/gap behavior | six canonical tables and 27 sections |
| CAP-CMD-302 | `06-command/modules/readiness-and-operations/capabilities/exercise-coordination.md` | six sections were prose only | added exercise/observation/task/Studio transition tables | six canonical tables and 27 sections |
| CAP-CMD-303 | `06-command/modules/readiness-and-operations/capabilities/improvement-actions.md` | six sections were prose only | added source/outcome/Task/verification tables | six canonical tables and 27 sections |
| CAP-CMD-304 | `06-command/modules/readiness-and-operations/capabilities/operational-plans.md` | six sections were prose only | added plan/reference/review/execution-boundary tables | six canonical tables and 27 sections |
| CAP-CMD-305 | `06-command/modules/readiness-and-operations/capabilities/capability-readiness.md` | six sections were prose only | added Register/health/evidence/assessment tables | six canonical tables and 27 sections |
| CAP-CMD-401 | `06-command/modules/customers-and-delivery/capabilities/customers-and-delivery-context.md` | S8 abbreviated; S9/S10/S13/S16/S17 prose | normalized deployment/report/engagement tables while retaining proposal status | six canonical tables and 27 sections |

## Final conformance matrix

| Capability ID | Sections 1–27 | S8 | S9 | S10 | S13 | S16 | S17 | Front matter | Verdict |
|---|---:|---:|---:|---:|---:|---:|---:|---|---|
| CAP-CMD-001 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-002 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-003 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-004 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-005 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-006 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-101 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-102 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-103 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-104 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-105 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-106 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-107 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-108 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-109 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-110 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-201 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-202 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-203 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-204 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-205 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-301 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-302 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-303 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-304 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-305 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-CMD-401 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |

## Final totals

| Measure | Result |
|---|---:|
| Capabilities conforming | 27 / 27 |
| Sections 1–27 present | 729 / 729 |
| Mandatory tables expected | 162 |
| Mandatory tables present | 162 |
| Table sections corrected | 123 |
| Mandatory sections represented only by prose | 0 |
| Empty mandatory tables | 0 |
| Generic copied tables without capability-specific adaptation | 0 |
| Duplicate Capability IDs | 0 |
| Concurrent active owners | 0 |
| New functional contradiction | 0 |

## Module result

| Module | Capabilities | Corrected | Remaining defect |
|---|---:|---:|---:|
| Mission Control | 6 | 6 | 0 |
| Incidents and Work Queue | 10 | 10 | 0 |
| Risk and Coverage | 5 | 5 | 0 |
| Readiness and Operations | 5 | 5 | 0 |
| Customers and Delivery | 1 | 1 | 0 |

## Closure gates — 75/75

The placeholder gate applies to the Phase 4A functional and corrective scope. Ten legacy active Command screen specifications remain intentionally deferred and unchanged, because this corrective mission expressly prohibits detailed screen rewriting.

| # | Gate | Verdict | Evidence |
|---:|---|---|---|
| 1 | Correct branch | PASS | `docs/cmdr-product-spec-foundation` |
| 2 | Correct PR | PASS | PR #2 |
| 3 | Base main | PASS | PR base unchanged |
| 4 | PR open | PASS | state open |
| 5 | PR Draft | PASS | draft true |
| 6 | PR unmerged | PASS | merged false |
| 7 | No force-push | PASS | fast-forward commit chain |
| 8 | History not rewritten | PASS | initial SHA remains ancestor |
| 9 | Root README unchanged | PASS | exact `# cmdr` |
| 10 | Main unchanged | PASS | same README blob and no main ref update |
| 11 | 27 Capability IDs | PASS | register and files contain 27 IDs |
| 12 | 27 unique IDs | PASS | no duplicate |
| 13 | 27 owners | PASS | all Command Product Lead |
| 14 | 27 user definitions | PASS | section 6 in every file |
| 15 | 27 input sets | PASS | 27 non-empty S8 tables |
| 16 | 27 output sets | PASS | 27 non-empty S16 tables |
| 17 | 27 read-object definitions | PASS | 27 non-empty S9 tables |
| 18 | 27 create/modify definitions | PASS | 27 non-empty S10 tables |
| 19 | 27 automation matrices | PASS | 27 canonical S13 tables |
| 20 | 27 transition tables | PASS | 27 canonical S17 tables |
| 21 | 27 action classifications | PASS | section 12 in every file |
| 22 | 27 no-AI alternatives | PASS | explicit alternative in each S13 |
| 23 | 27 Given/When/Then sets | PASS | section 25 in every file |
| 24 | 27 delivery classifications | PASS | section 24 and front matter |
| 25 | 27 source-of-truth sections | PASS | section 19 |
| 26 | 27 provenance sections | PASS | section 20 |
| 27 | 27 limits/error sections | PASS | section 22 |
| 28 | 27/27 S8 tables | PASS | exact header |
| 29 | 27/27 S9 tables | PASS | exact header |
| 30 | 27/27 S10 tables | PASS | exact header |
| 31 | 27/27 S13 tables | PASS | exact header |
| 32 | 27/27 S16 tables | PASS | exact header |
| 33 | 27/27 S17 tables | PASS | exact header |
| 34 | 162/162 tables present | PASS | final matrix |
| 35 | No empty table | PASS | at least one specific row per table |
| 36 | No prose-only replacement | PASS | all six tables present per file |
| 37 | No unadapted generic table | PASS | rows are capability-specific |
| 38 | Command owns Incident | PASS | Incident boundaries unchanged |
| 39 | Command owns operational Task | PASS | Task ownership unchanged |
| 40 | Command does not own Case | PASS | Investigate projection only |
| 41 | Command does not own Evidence | PASS | no Evidence mutation |
| 42 | Command does not own Finding | PASS | Investigate projection only |
| 43 | Command does not own Decision | PASS | Govern projection only |
| 44 | Command does not own Response Run | PASS | Govern projection only |
| 45 | Command does not own Automation Run | PASS | Studio projection only |
| 46 | Command does not own Endpoint Agent Fleet | PASS | Settings projection only |
| 47 | One Work Queue workspace | PASS | CAP-CMD-101 unchanged |
| 48 | One Work Queue route | PASS | views remain route parameters |
| 49 | Six exact system views | PASS | All, Incidents, Tasks, Unassigned, SLA Risk, My Work |
| 50 | Team Load absent | PASS | not a system view |
| 51 | Five legacy screens deprecated | PASS | no screen file changed |
| 52 | No new Work Queue screen | PASS | zero screen additions |
| 53 | OPEN-006 open | PASS | CAP-CMD-401 remains proposed |
| 54 | OPEN-010 open | PASS | Readiness files preserve it |
| 55 | OPEN-013 open | PASS | class-2 governance unresolved |
| 56 | Fifteen OPEN decisions coherent | PASS | unresolved-decision set unchanged |
| 57 | OPEN-009 only historical resolved | PASS | decision history unchanged |
| 58 | No detailed screen rewrite | PASS | zero `screens/` changes |
| 59 | No massive object rewrite | PASS | zero object changes in closure |
| 60 | No atomic permission finalization | PASS | permission files unchanged |
| 61 | No global namespace normalization | PASS | permission catalog unchanged |
| 62 | No API created | PASS | documentation-only capability diff |
| 63 | No protocol created | PASS | no technical contract added |
| 64 | No product code | PASS | Markdown-only changes |
| 65 | No font file | PASS | no asset/font change |
| 66 | Phase 4B not started | PASS | no Investigate capability change |
| 67 | No active generic placeholder in corrective scope | PASS | zero placeholder in 27 capabilities/template/report |
| 68 | No targeted empty file | PASS | every changed file substantive |
| 69 | No broken local link introduced | PASS | no new unresolved local Markdown link |
| 70 | No artificial Validated status | PASS | all documents remain draft |
| 71 | Capability Register coherent | PASS | 27 IDs/owners/statuses unchanged |
| 72 | Traceability coherent | PASS | scores and Requirement IDs unchanged |
| 73 | Baseline coherent | PASS | 39→162 and 123 corrections recorded |
| 74 | STATUS coherent | PASS | Phase 4A PASS; global/Phase 4 PARTIAL |
| 75 | PR description coherent | PASS | closure, limits and Phase 4B state stated |

## Verdict

**Phase 4A: PASS.** The strict capability-template closure gate is satisfied. Phase 4 global and repository global remain **PARTIAL**. Phase 4B is **not started**.