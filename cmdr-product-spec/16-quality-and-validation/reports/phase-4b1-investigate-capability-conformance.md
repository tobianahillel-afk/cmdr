---
id: phase-4b1-investigate-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-04
source-of-truth: quality-report
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-AI-002
  - REQ-UX-010
---

# Phase 4B.1 — Investigate capability conformance

## Verdict

**PASS — 140/140 gates.** Phase 4B global, Phase 4 global and repository maturity remain PARTIAL. Phase 4B.2 and Phase 4B.3 are not started.

## Scope and evidence

- Initial published head: `8d9aec3ddae2a6e45b0f4dffdc97a2ec0c8ad903`.
- Canonical modules treated: Signals and Hunt; Cases and Evidence.
- Canonical capabilities: 22.
- Numbered sections expected/present: 594/594.
- Mandatory tables expected/present: 132/132.
- Active Investigate screens read: 14; modified: 0; detailed rewrites: 0.
- Object files modified: 0.
- Atomic permission sources modified: 0.
- CAP-INV-2xx/3xx/4xx/5xx created: 0.
- APIs, protocols, engine choices, product code and fonts: 0.

## Detached-work disposition

| Commit or tree | Content | Decision | Final commit | Reason |
|---|---|---|---|---|
| `f8daf01d6a4fd540edbd8e6258ebb058efc11526` | Investigate framework, ownership, maps and migration pointers | reused as-is | same SHA | direct child of published head; scope-clean |
| `de82bfa522b0310292012733100b2ddb53b441f8` | CAP-INV-001..008 | reused as-is | same SHA | eight complete canonical capabilities |
| detached trees for CAP-INV-101..109 | first Cases and Evidence specifications | partially recovered | `b29febaa571e41f0aa9632950b187d6a2e41073e` | valid blobs consolidated; test history excluded |
| temp branch for CAP-INV-110..114 | five final specifications | partially recovered | `b29febaa571e41f0aa9632950b187d6a2e41073e` | content consolidated into one functional commit |
| `tmp/noop-test*`, per-file temp commits and duplicate checkpoint commits | build experiments or no-op commits | abandoned | none | not functional, not referenced by final branch |

## Initial defects and corrections

| Initial defect | Affected scope | Correction | Final evidence |
|---|---|---|---|
| No registered CAP-INV IDs | all Investigate | 22 immutable IDs and canonical files | Capability Register and capability map |
| Old module architecture remained active or generic | eight legacy 4B.1 modules | needs migrated; module sources converted to migration/deprecated pointers | old module README files and canonical module READMEs |
| No complete capability contract | all target functions | 27 sections and six exact mandatory tables per capability | matrix below |
| Artifact/Attachment ambiguity was only briefly noted | OPEN-014 | complete 27-section option analysis without closure | `07-investigate/modules/cases-and-evidence/artifact-versus-attachment.md` |
| Transitions and automation boundaries were too aggregated | cross-product flows | detailed transition, action and automation matrices | cross-product links, action classification, AI model |
| Registers contained Command only | governance | added 22 CAP-INV and 14 functional dependencies | Capability and Dependency Registers |
| Phase 4B.1 had no quality evidence | validation | this report, baseline, status and changelog | quality and governance files |

## Capability-template matrix

| Capability ID | Sections 1–27 | S8 | S9 | S10 | S13 | S16 | S17 | Front matter | Verdict |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| CAP-INV-001 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-002 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-003 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-004 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-005 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-006 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-007 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-008 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-101 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-102 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-103 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-104 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-105 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-106 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-107 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-108 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-109 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-110 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-111 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-112 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-113 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-114 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |

## Template totals

| Control | Result |
|---|---:|
| Canonical files | 22/22 |
| Unique IDs | 22/22 |
| Front matter | 22/22 |
| Sections 1–27 | 594/594 |
| S8 input tables | 22/22 |
| S9 read-object tables | 22/22 |
| S10 create/modify tables | 22/22 |
| S13 automation/AI tables | 22/22 |
| S16 output tables | 22/22 |
| S17 transition tables | 22/22 |
| Mandatory tables total | 132/132 |
| Empty mandatory tables | 0 |
| Prose-only mandatory sections | 0 |
| Generic unadapted tables | 0 |
| Duplicate IDs | 0 |
| Concurrent owners | 0 |
| Active contradictions introduced | 0 |

## Ownership and concept assertions

- Investigate owns Case, Hypothesis, Artifact, Evidence and Finding.
- Command owns Detection, Signal, Alert, Incident and operational Task.
- Govern owns Action Request lifecycle, Decision, Approval, Response Run and Result.
- Shared Capabilities owns Telemetry Event, Entity, Query, Search Job, Timeline mechanisms, Saved Views and Reporting Engine according to current registers.
- CMDR Studio owns Skill, Tool, Tool Call, Workflow, Automation Agent, Agent Team, Human Gate and Automation Run.
- Platform Settings owns Endpoint Agent Fleet and Endpoint Policies.
- Event is not Evidence; Detection is not Signal; Signal is not Alert; Alert is not Incident; Incident is not Case.
- Hypothesis is not Finding; Entity is not Artifact; Artifact is not Evidence.
- Attachment is neither automatically Artifact nor Evidence; OPEN-014 remains open.
- Finding is neither Decision nor Result; Action Request is not Decision; Response Run is not Automation Run; Report is not Finding; Note is not Task.

## Quality gates — 140/140

| # | Gate | Verdict |
|---:|---|---|
| 1 | Repository correct | PASS |
| 2 | Working branch correct | PASS |
| 3 | PR #2 correct | PASS |
| 4 | Base is main | PASS |
| 5 | PR open | PASS |
| 6 | PR Draft | PASS |
| 7 | PR unmerged | PASS |
| 8 | No force-push | PASS |
| 9 | History not rewritten | PASS |
| 10 | Root README unchanged | PASS |
| 11 | main unchanged | PASS |
| 12 | Final Phase 4B.1 commits reachable from final branch | PASS |
| 13 | All active Investigate files read | PASS |
| 14 | All legacy 4B.1 modules read | PASS |
| 15 | Fourteen active screens read | PASS |
| 16 | Required objects read or absences recorded | PASS |
| 17 | Command/Govern/Studio boundaries verified | PASS |
| 18 | Detached framework and capability work audited | PASS |
| 19 | No detached content represented as published before ref update | PASS |
| 20 | Twenty-two Capability IDs | PASS |
| 21 | Twenty-two unique IDs | PASS |
| 22 | Twenty-two canonical files | PASS |
| 23 | Twenty-two owners | PASS |
| 24 | Twenty-two user sets | PASS |
| 25 | Twenty-two specific user problems | PASS |
| 26 | Twenty-two objective sets | PASS |
| 27 | Twenty-two non-objective sets | PASS |
| 28 | Twenty-two input sets | PASS |
| 29 | Twenty-two read-object definitions | PASS |
| 30 | Twenty-two create/modify-object definitions | PASS |
| 31 | Twenty-two classified action sets | PASS |
| 32 | Twenty-two S13 matrices | PASS |
| 33 | Twenty-two specific state sets | PASS |
| 34 | Twenty-two output sets | PASS |
| 35 | Twenty-two transition tables | PASS |
| 36 | Twenty-two sources of truth | PASS |
| 37 | Twenty-two provenance sections | PASS |
| 38 | Twenty-two functional permission sections | PASS |
| 39 | Twenty-two limits/error sections | PASS |
| 40 | Twenty-two conceptual metric sets | PASS |
| 41 | Twenty-two delivery classifications | PASS |
| 42 | Twenty-two Given/When/Then sets | PASS |
| 43 | Twenty-two no-AI alternatives | PASS |
| 44 | Twenty-two documentary-consumer sections | PASS |
| 45 | 22/22 S8 tables | PASS |
| 46 | 22/22 S9 tables | PASS |
| 47 | 22/22 S10 tables | PASS |
| 48 | 22/22 S13 tables | PASS |
| 49 | 22/22 S16 tables | PASS |
| 50 | 22/22 S17 tables | PASS |
| 51 | 132/132 mandatory tables | PASS |
| 52 | No empty mandatory table | PASS |
| 53 | No prose-only mandatory section | PASS |
| 54 | No generic unadapted table | PASS |
| 55 | 594/594 numbered sections | PASS |
| 56 | Investigate owns Case | PASS |
| 57 | Investigate owns Hypothesis | PASS |
| 58 | Investigate owns Artifact | PASS |
| 59 | Investigate owns Evidence | PASS |
| 60 | Investigate owns Finding | PASS |
| 61 | Investigate does not own Detection | PASS |
| 62 | Investigate does not own Signal | PASS |
| 63 | Investigate does not own Alert | PASS |
| 64 | Investigate does not own Incident | PASS |
| 65 | Investigate does not own Decision | PASS |
| 66 | Investigate does not own Response Run | PASS |
| 67 | Investigate does not own Result | PASS |
| 68 | Investigate does not own Automation Run | PASS |
| 69 | Investigate does not own Endpoint Agent Fleet | PASS |
| 70 | Action Request lifecycle remains Govern | PASS |
| 71 | Event distinct from Evidence | PASS |
| 72 | Detection distinct from Signal | PASS |
| 73 | Signal distinct from Alert | PASS |
| 74 | Alert distinct from Incident | PASS |
| 75 | Incident distinct from Case | PASS |
| 76 | Hypothesis distinct from Finding | PASS |
| 77 | Artifact distinct from Evidence | PASS |
| 78 | Attachment not arbitrarily merged; OPEN-014 open | PASS |
| 79 | Finding distinct from Decision | PASS |
| 80 | Finding distinct from Result | PASS |
| 81 | Action Request distinct from Decision | PASS |
| 82 | OPEN-014 retained | PASS |
| 83 | Signals and Hunt has eight conforming capabilities | PASS |
| 84 | Cases and Evidence has fourteen conforming capabilities | PASS |
| 85 | Case Queue distinct from Command Work Queue | PASS |
| 86 | Saved Search distinct from Saved View | PASS |
| 87 | Query Asset distinct from Detection Rule | PASS |
| 88 | Reporting Engine not duplicated | PASS |
| 89 | Trace not duplicated | PASS |
| 90 | Timeline component not duplicated | PASS |
| 91 | AI optional | PASS |
| 92 | No mandatory chatbot | PASS |
| 93 | No automatic Evidence qualification | PASS |
| 94 | No automatic Finding confirmation | PASS |
| 95 | No Decision created by Investigate | PASS |
| 96 | No silent Entity merge | PASS |
| 97 | No silent sensitive Query execution | PASS |
| 98 | Automated provenance visible | PASS |
| 99 | No-AI alternative for essential capabilities | PASS |
| 100 | No Govern bypass | PASS |
| 101 | No detailed screen rewrite | PASS |
| 102 | No complete object schema created | PASS |
| 103 | No final cardinality created | PASS |
| 104 | No atomic permission matrix finalized | PASS |
| 105 | No global namespace normalization | PASS |
| 106 | No API created | PASS |
| 107 | No protocol created | PASS |
| 108 | No search engine selected | PASS |
| 109 | No final query language selected | PASS |
| 110 | No forensic engine selected | PASS |
| 111 | No product code added | PASS |
| 112 | No font file added | PASS |
| 113 | No CAP-INV-2xx capability | PASS |
| 114 | No CAP-INV-3xx/4xx/5xx capability | PASS |
| 115 | Phase 4B.2 not started | PASS |
| 116 | Phase 4B.3 not started | PASS |
| 117 | Capability Register updated | PASS |
| 118 | Dependency Register updated | PASS |
| 119 | Object Consumption Map created | PASS |
| 120 | Screen Capability Map created | PASS |
| 121 | Automation and AI Model created | PASS |
| 122 | Action Classification created | PASS |
| 123 | Cross-product transitions created | PASS |
| 124 | Artifact versus Attachment substantive | PASS |
| 125 | Requirements Matrix updated | PASS |
| 126 | Baseline updated | PASS |
| 127 | STATUS updated | PASS |
| 128 | CHANGELOG updated | PASS |
| 129 | PR description coherent after publication | PASS |
| 130 | No generic placeholder in targeted active scope | PASS |
| 131 | No targeted empty file | PASS |
| 132 | No broken local link introduced | PASS |
| 133 | No concurrent active owner | PASS |
| 134 | No duplicated Shared Capability | PASS |
| 135 | No artificially Validated document | PASS |
| 136 | Fifteen OPEN decisions coherent | PASS |
| 137 | No unnecessary new OPEN | PASS |
| 138 | Requirement scores justified and unchanged | PASS |
| 139 | Conformance report published | PASS |
| 140 | Final remote SHA verified after publication | PASS |

## Final limits

Phase 4B.1 defines functions, ownership, transitions and quality evidence. It does not prove implementation. Collection and Live Response, Analysis Workbench, Detection Engineering, Intelligence, detailed screens, complete objects, permissions, technical architecture and software remain future work.
