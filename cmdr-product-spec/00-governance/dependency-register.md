---
id: dependency-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: registry
requirements:
  - REQ-PROD-006
  - REQ-PROD-009
  - REQ-PROD-019
  - REQ-PROD-020
---
# Dependency Register

This register records functional and documentary dependencies only. It chooses no protocol, API, engine, rule language, parser, compiler, model, provider, command, package format, deployment architecture or production pipeline.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN | Review |
|---|---|---|---|---|---|---|---|---|---|
| DEP-001 | source material | product specifications | decision | apply mission, principles and boundaries | active | Product Architecture | yes | REQ-PROD-001..020 | continuous |
| DEP-002 | product boundaries / ownership | all modules | boundary | prevent concurrent ownership | active | Product Architecture | yes | REQ-PROD-013..018 | continuous |
| DEP-003 | object/permission phases | capabilities | model/security | final schemas and atomic permissions are future | partial | Architecture/Security | before implementation | REQ-OBJ/SEC | Phase 7 |
| DEP-004 | Shared mechanisms | all products | shared | Jobs, Trace, Activity, Linking, Search, Export, Reporting, Versioning, Collaboration and Recovery | partial | Shared | yes | REQ-PROD-019,020 | Technique |
| DEP-005 | OPEN-005 | Analysis Workbench and Forensics | future engine | forensic engines only; not Detection Engineering | open | Investigate | before forensic delivery | OPEN-005 | Phase 8 |
| DEP-006 | OPEN-008 | platform/source support | platform | source, sensor and endpoint platform support unresolved | open | Endpoint Agent / Settings | before delivery | OPEN-008 | Phase 4D/8 |
| DEP-007 | OPEN-013 | class-2 actions | authority | default step-up/Govern policy unresolved | open | Security | yes | OPEN-013 | Phase 7 |
| DEP-008 | OPEN-014 | Artifact/Attachment/dataset relations | object | relation and retention semantics unresolved | open | Investigate | model blocking | OPEN-014 | Phase 7 |
| DEP-009 | OPEN-015 | Studio/Govern/Investigate runs | run model | Automation Run / Response Run bridge unresolved | open | Studio + Govern | yes | OPEN-015 | Phase 7 |
| DEP-010 | OPEN-017 | Detection runtime, target language and portability | product architecture | portable model, native content, hybrid, future CMDR runtime or capability-specific combination remain undecided | open | Product Architecture | before implementation/promotion contract | OPEN-017 | Phase 7/8 |
| DEP-INV-201 | Collection and Live Response | Analysis Workbench families | source/provenance | authorized acquisition outputs feed analysis | active | Investigate / Endpoint | yes | REQ-INV-001 | 4B.2 |
| DEP-INV-301 | Artifact Management | CAP-INV-301..397 | analysis source | Artifact, versions, custody and lineage | active | Investigate | yes | REQ-OBJ-003,004 | 4B.2 |
| DEP-INV-314 | Static Analysis | Dynamic Sandbox | handoff | Derived Artifacts and static context | active | Investigate | no | REQ-INV-002,005 | 4B.2 |
| DEP-INV-329 | Static/Dynamic | Reverse and Debugger | handoff | code and runtime analysis context | active | Investigate | no | REQ-INV-003,004 | 4B.2 |
| DEP-INV-347 | Collection/Artifact | Memory Forensics | source/custody | Memory Image and acquisition context | active | Investigate / Endpoint | yes | REQ-INV-001; OPEN-008 | 4B.2 |
| DEP-INV-363 | Collection/Artifact | Disk and Filesystem Forensics | source/custody | Disk Image and acquisition context | active | Investigate / Endpoint | yes | REQ-INV-001; OPEN-008 | 4B.2 |
| DEP-INV-380 | Collection/Artifact | Network Forensics | source/custody | Capture Artifact and acquisition context | active | Investigate / Endpoint | yes | REQ-INV-001; OPEN-008 | 4B.2 |
| DEP-INV-401 | Finding/Hypothesis/Case/Incident/Hunt and technical handoffs | CAP-INV-401..417 | Detection authoring | need, sources, authoring, validation, replay, review, coverage and Review Package | active | Investigate | yes | REQ-INV-006; OPEN-008/013/014/015 | 4B.3A.1 |
| DEP-INV-418 | CAP-INV-417 | CAP-INV-418 | review handoff | immutable candidate version and complete authoring evidence | active | Investigate | yes | REQ-INV-006 | 4B.3A.2 |
| DEP-INV-419 | CAP-INV-418 | CAP-INV-419 | readiness | review disposition, targets, dependencies, rollback and observation plan | active | Investigate / Settings projections | yes | OPEN-008/017 | 4B.3A.2 |
| DEP-INV-420 | CAP-INV-419/421/423/430/433/434 | Govern Action Request | authority handoff | class, targets, risks, rollout, rollback, criteria and evidence | active | Investigate → Govern | yes for class 3 | REQ-SEC-001,002; OPEN-013 | 4B.3A.2 |
| DEP-INV-421 | Settings environments/targets/runtimes | CAP-INV-421/424/425/426 | administrative projection | target selection, capability, state, version and health remain Settings-owned | partial | Platform Settings | yes | OPEN-008/017 | Settings/Technique |
| DEP-INV-422 | Govern Decision/Approval/Response Run/Result | CAP-INV-420/424/433/434 | authority/execution | production change authority and result remain Govern-owned | partial | Govern | yes | REQ-OBJ-006,007,010; OPEN-015 | Govern/Objects |
| DEP-INV-423 | Endpoint Agent/runtime owner | CAP-INV-422..426/431..433 | runtime projection | execution, version, health and target result | partial | Endpoint/Settings/runtime | yes | OPEN-008/017 | Endpoint/Technique |
| DEP-INV-424 | Command Detection/Signal/Alert/Incident | CAP-INV-425..428/434/435 | runtime/feedback | Command retains operational objects and dispositions | active | Command | yes | REQ-PROD-013,015 | Command/Objects |
| DEP-INV-425 | Studio Tool/Tool Call/Automation Run | CAP-INV-422/426/431/432/435 | execution/provenance | optional deterministic assessment and attributed automation | partial | Studio | no essential AI dependency | OPEN-015 | Studio/Objects |
| DEP-INV-426 | Shared Jobs/Trace/Activity/Versioning/Comparison/Metrics/Reporting | CAP-INV-418..435 | shared | generic execution, lineage, metrics and reporting | partial | Shared | yes | REQ-PROD-019,020 | Shared/Technique |
| DEP-INV-427 | CAP-INV-422 | CAP-INV-418/423/427/428 | runtime evaluation | shadow evidence is non-alerting and distinct from replay | active | Investigate | no | REQ-INV-006 | 4B.3A.2 |
| DEP-INV-428 | CAP-INV-423/424 Results | CAP-INV-425/426/427/433 | canary/deployment observation | per-target partiality, versions, errors and stop criteria | active | Investigate projections | yes | OPEN-013/015/017 | 4B.3A.2 |
| DEP-INV-429 | CAP-INV-427/428/431/432 | CAP-INV-429/430/433/434 | continuous improvement | quality, outcomes, drift and performance create proposals only | active | Investigate | no | REQ-INV-006 | 4B.3A.2 |
| DEP-INV-430 | CAP-INV-429/431/432/434 | CAP-INV-406/411..417 | return to authoring | new Draft, revalidation, tests, replay and coverage update | active | Investigate | yes | REQ-INV-006 | 4B.3A |
| DEP-INV-431 | CAP-INV-401..434 | CAP-INV-435 | lifecycle provenance | complete lineage and unresolved links | active | Investigate / source owners | yes | REQ-PROD-020; OPEN-014/015 | 4B.3A.2 |
| DEP-INV-432 | CAP-INV-435 | future 4B.3B | future Intelligence handoff | candidate context only; no Intelligence object or CAP-INV-5xx | planned | future owner | no | REQ-INV-006 | 4B.3B |

OPEN-017 is now the canonical open decision for Detection runtime, target-language and portability strategy. No option is selected. OPEN-005 remains forensic-only. Sixteen decisions are open; OPEN-009 remains the only historically resolved item.
