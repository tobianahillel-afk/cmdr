---
id: dependency-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-019
---
# Dependency Register

This register records product and documentary dependencies only. It chooses no protocol, API, storage, engine, provider, port or command.

| ID | Source | Dependent | Type | Reason | Status | Impact | Owner | Blocking | Requirement IDs | Review |
|---|---|---|---|---|---|---|---|---|---|---|
| DEP-001 | source-material/ | 01-product-vision/ | decision | apply mission, principles and boundaries | active | vision inconsistency | Head of Product | yes | REQ-PROD-001..020 | Phase 1 |
| DEP-002 | product-boundaries.md | product READMEs | boundary | prevent concurrent ownership | active | duplicate responsibility | Product Architecture | yes | REQ-PROD-013..018 | continuous |
| DEP-003 | ownership-register.md | object model | ownership | guide objects and states | active | concurrent objects | Product Architecture | yes | REQ-OBJ-001..012 | Phase 7 |
| DEP-004 | terminology-rules.md | UX/content | vocabulary | distinguish Page/View/Mode/Filter/Run/Agent | active | ambiguous navigation | Content Design | yes | REQ-UX-001; REQ-OBJ-009 | continuous |
| DEP-005 | ADR-0005 | Work Queue | UX | consolidate views | resolved | screen clones | UX Architecture | yes | REQ-UX-008,009 | Phase 3 |
| DEP-006 | OPEN-005 | Analysis Workbench | future engine | initial forensic engines | open | 4B.2B engine boundary unknown | Investigate Product Lead | before delivery | REQ-PROD-052 | Phase 4B.2B/8 |
| DEP-007 | Capability map | modules | classification | do not present planned as delivered | active | misleading promise | Product Architecture | yes | REQ-PROD-012,019 | Phase 4 |
| DEP-008 | OPEN-007 | Studio/Govern | authority | Human Gate versus Decision | open | double approval/bypass | Security Architecture | yes | REQ-PROD-054 | Phase 4C/7 |
| DEP-009 | OPEN-008 | Endpoint Agent | platform | initial platform support | open | capabilities/support unknown | Endpoint Product Lead | before delivery | REQ-PROD-055 | Phase 4D/8 |
| DEP-010 | OPEN-014 | Investigate | object model | Artifact versus Attachment | open | ambiguous promotion/custody | Investigate Product Lead | model blocking | REQ-PROD-061 | Phase 7 |
| DEP-011 | OPEN-015 | Studio/Govern | run model | Automation Run versus Response Run | open | audit confusion | Studio + Govern | yes | REQ-PROD-062; REQ-OBJ-009 | Phase 7 |
| DEP-012 | Phase 3 | screen pilots | UX | quality level before screen rewrite | planned | template screens | UX Architecture | yes | REQ-UX-010 | Phase 6 |
| DEP-013 | Phase 5 | journeys | transition | end-to-end context | planned | broken navigation | Product Architecture | yes | REQ-JRN-001..008 | Phase 5 |
| DEP-014 | Phase 7 | permissions | security | atomize namespaces after boundaries | planned | inconsistent access control | Security Architecture | yes | REQ-SEC-001..005 | Phase 7 |
| DEP-015 | OPEN-013 | class 2 actions | governance | default governance/step-up | open | reversible action authority unknown | Security Architecture | yes | REQ-PROD-060; REQ-SEC-002 | Phase 4C/7 |
| DEP-CMD-001 | Capability template | 27 Command capabilities | documentary | 27 sections and immutable IDs | active | incomplete specifications | Product Architecture | yes | REQ-PROD-006,013 | Phase 4A |
| DEP-CMD-002 | Incident/Task | Work Queue | object | coordination, assignment, priority and SLA | partial | Draft schemas/states | Command Product Lead | later | REQ-OBJ-001,012 | Phase 7 |
| DEP-CMD-003 | Business Service Catalog | Risk and Coverage | shared | Service projection | partial | criticité/dependencies draft | Shared Capabilities Lead | no | REQ-PROD-021 | Phase 4D/7 |
| DEP-CMD-004 | Exposure sources | Exposure Overview | integration | Exposure projection | open | source/freshness variable | Command + source owner | no | REQ-PROD-037..042 | Phase 4D/7 |
| DEP-CMD-005 | Saved Views | Unified Work Queue | shared | permission-aware views | active | no second source | Shared Capabilities Lead | yes | REQ-OBJ-011,012 | Phase 4A/4D |
| DEP-CMD-006 | OPEN-006 | Customers and Delivery | scope | deployment applicability | open | module not universal | Command Product Lead | before activation | REQ-PROD-053 | review |
| DEP-CMD-007 | OPEN-010 | density defaults | UX | role/activity defaults | open | preferences not final | UX Architecture | no | REQ-PROD-057 | Phase 3 review |
| DEP-CMD-008 | OPEN-013 | Command class 2 | authority | default governance | open | step-up/Govern context | Security Architecture | yes | REQ-PROD-060; REQ-SEC-002 | Phase 4C/7 |
| DEP-CMD-009 | Investigate Case | Incident Coordination | transition | idempotent create/link | partial | detailed contract future | Investigate Product Lead | no | REQ-PROD-008; REQ-OBJ-002 | Phase 4B/5 |
| DEP-CMD-010 | Govern Action Request | Escalation | transition | context/impact/action package | partial | authority remains Govern | Govern Product Lead | class 3/4 | REQ-PROD-015; REQ-SEC-002 | Phase 4C/5 |
| DEP-CMD-011 | Govern Result | Situation/Incident | projection | verification and residual risk | partial | Result remains Govern | Govern Product Lead | no | REQ-PROD-008; REQ-OBJ-007 | Phase 4C/5 |
| DEP-CMD-012 | Permission catalog | Command actions | permission | assignment/SLA/bulk/handover families | partial | no final atomic permissions | Security Architecture | before implementation | REQ-SEC-001..005 | Phase 7 |
| DEP-INV-001 | Command Signal/Alert/Incident | Signal Triage and Case | projection/transition | source and operational context | active | no ownership transfer | Command + Investigate | yes | REQ-PROD-008,014 | Phase 4B.1 |
| DEP-INV-002 | Shared Query/Search Job | Event Search/Hunt | shared | query execution and results | partial | engine/language future | Shared Capabilities Lead | no | REQ-PROD-014,019 | Technique |
| DEP-INV-003 | Shared Trace/Activity | Search/Hunt provenance | shared | reconstruct runs and sources | partial | record model future | Shared/Trust | yes | REQ-AI-002; REQ-PROD-020 | Phase 7 |
| DEP-INV-004 | Case object | Cases and Evidence | object | workspace and relations | partial | states/cardinalities future | Investigate Product Lead | later | REQ-OBJ-002; REQ-PROD-014 | Phase 7 |
| DEP-INV-005 | Artifact/Evidence/Finding | reasoning chain | object | explicit qualification and review | partial | trust/states future | Investigate Product Lead | yes | REQ-OBJ-003..005 | Phase 7 |
| DEP-INV-006 | Shared Entity | Entity management | shared | relations and resolution | partial | merge/identity future | Shared Capabilities Lead | yes | REQ-OBJ-006 | Phase 7/Technique |
| DEP-INV-007 | OPEN-014 | Attachment Handling | object decision | Artifact/Attachment boundary | open | promotion/retention ambiguous | Investigate Product Lead | yes model | REQ-PROD-061 | Phase 7 |
| DEP-INV-008 | Govern Action Request | Action Request Preparation | transition | authority handoff | partial | Govern lifecycle future | Govern Product Lead | class 3/4 | REQ-PROD-016 | Phase 4C/5 |
| DEP-INV-009 | Shared Reporting Engine | Case Reporting | shared | draft/citations/export | partial | Report object/export future | Shared Capabilities Lead | later | REQ-PROD-018 | Shared/Technique |
| DEP-INV-010 | Studio Workflow/Automation Run | optional assistance | automation | proposals and provenance | partial | OPEN-015 | Studio Product Lead | no essential workflow | REQ-AI-002; REQ-PROD-062 | Phase 4C/7 |
| DEP-INV-011 | OPEN-013 | Investigate class 2 | authority | mutations/review/links | open | step-up defaults unknown | Security Architecture | yes | REQ-PROD-060 | Phase 4C/7 |
| DEP-INV-012 | OPEN-007 | Human Gate/Govern | authority | submission and approval boundary | open | double gate risk | Security Architecture | yes | REQ-PROD-054 | Phase 4C/7 |
| DEP-INV-013 | Phase 4B.2A | Evidence collection gaps | future capability | collection requests | resolved by 4B.2A definitions | implementation still planned | Investigate Product Lead | no | REQ-PROD-014,018 | Phase 4B.2A |
| DEP-INV-014 | Phase 4B.2B/4B.3 | Workbench/Detection handoffs | future boundary | Artifacts and review outputs | planned | no capabilities yet | Investigate Product Lead | no | REQ-PROD-020,052 | future |
| DEP-INV-201 | Platform Settings Fleet projection | CAP-INV-201 | projection | Endpoint identity, posture and freshness | partial | read only; admin remains Settings | Platform Settings Product Lead | yes | REQ-PROD-014,055 | Phase 4B.2A/4D |
| DEP-INV-202 | Endpoint Agent health/capabilities | CAP-INV-201,202 | product boundary | availability and supported operations | partial | OPEN-008 | Endpoint Agent Product Lead | before execution | REQ-PROD-055 | Phase 4D/8 |
| DEP-INV-203 | Endpoint Policy | CAP-INV-202,204..211 | policy projection | scope limits and allowed actions | partial | no local assignment | Platform Settings Product Lead | yes | REQ-SEC-001,002 | Phase 4D/7 |
| DEP-INV-204 | Shared Background Jobs | CAP-INV-203 | shared mechanism | queue/progress/cancel/partial | partial | Collection Job remains business context | Shared Capabilities Lead | yes | REQ-PROD-019 | Technique |
| DEP-INV-205 | Endpoint Agent collection | CAP-INV-204..208 | execution | local acquisition and result reporting | planned | no engine/protocol/platform claim | Endpoint Agent Product Lead | before delivery | REQ-PROD-014,055 | Phase 4D/8 |
| DEP-INV-206 | CAP-INV-105 Artifact Management | collection outputs | object handoff | results become Artifacts explicitly | active | Artifact never Evidence automatically | Investigate Product Lead | yes | REQ-OBJ-003,004 | Phase 4B.2A |
| DEP-INV-207 | CAP-INV-107 Evidence Creation | Artifact/result review | object handoff | human qualification | active | no automatic conversion | Investigate Product Lead | yes | REQ-OBJ-004 | Phase 4B.1 |
| DEP-INV-208 | Endpoint Agent availability | CAP-INV-209,210 | execution/session | open/reconnect/close and local operations | planned | offline/conflict explicit | Endpoint Agent Product Lead | before delivery | REQ-PROD-055 | Phase 4D/8 |
| DEP-INV-209 | Govern Action Request/Decision/Run/Result | CAP-INV-215/113 | authority | containment and risky actions | partial | no local Decision or Run | Govern Product Lead | class 3/4 | REQ-PROD-016; REQ-SEC-002 | Phase 4C/5 |
| DEP-INV-210 | Studio Workflow/Automation Run/Tool Calls | CAP-INV-202,209,210,214 | optional automation | orchestration and provenance | partial | OPEN-015; no essential dependency | Studio Product Lead | no | REQ-AI-002; REQ-PROD-062 | Phase 4C/7 |
| DEP-INV-211 | OPEN-013 | CAP-INV-202,203,209..213 | authority | class 2 defaults/step-up | open | no default invented | Security Architecture | yes | REQ-PROD-060; REQ-SEC-002 | Phase 4C/7 |
| DEP-INV-212 | OPEN-007 | CAP-INV-209,210,215 | authority | Human Gate versus Govern | open | no gate equivalence assumed | Security Architecture | yes | REQ-PROD-054 | Phase 4C/7 |
| DEP-INV-213 | OPEN-008 | CAP-INV-201..211,215 | platform | platform and capability support | open | unsupported explicit | Endpoint Agent Product Lead | yes before delivery | REQ-PROD-055 | Phase 4D/8 |
| DEP-INV-214 | OPEN-015 | CAP-INV-209,210,212,214,215 | run model | Automation Run versus Response Run | open | results remain distinct | Studio + Govern | yes | REQ-PROD-062; REQ-OBJ-009 | Phase 7 |

## Phase 4B.2A interpretation

- Investigate owns collection/session business context and Case relations, not Fleet, Policies, local execution or Govern authority.
- Collection Job remains distinct from Shared Background Job; Operation Result remains distinct from Govern Result.
- OPEN-005 is a future Phase 4B.2B dependency only and no forensic engine is selected.
- OPEN-007, OPEN-008, OPEN-013 and OPEN-015 remain open.
- Shared mechanisms and product projections never transfer ownership.
