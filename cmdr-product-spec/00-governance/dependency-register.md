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
| DEP-INV-201 | Collection and Live Response | Analysis Workbench families | source/provenance | authorized acquisition outputs feed analysis | active | Investigate / Endpoint | yes | REQ-INV-001 | 4B.2 |
| DEP-INV-301 | Artifact Management | CAP-INV-301..397 | analysis source | Artifact, versions, custody and lineage | active | Investigate | yes | REQ-OBJ-003,004 | 4B.2 |
| DEP-INV-314 | Static Analysis | Dynamic Sandbox | handoff | Derived Artifacts and static context | active | Investigate | no | REQ-INV-002,005 | 4B.2 |
| DEP-INV-329 | Static/Dynamic | Reverse and Debugger | handoff | code and runtime analysis context | active | Investigate | no | REQ-INV-003,004 | 4B.2 |
| DEP-INV-347 | Collection/Artifact | Memory Forensics | source/custody | Memory Image and acquisition context | active | Investigate / Endpoint | yes | REQ-INV-001; OPEN-008 | 4B.2 |
| DEP-INV-363 | Collection/Artifact | Disk and Filesystem Forensics | source/custody | Disk Image and acquisition context | active | Investigate / Endpoint | yes | REQ-INV-001; OPEN-008 | 4B.2 |
| DEP-INV-380 | Collection/Artifact | Network Forensics | source/custody | Capture Artifact and acquisition context | active | Investigate / Endpoint | yes | REQ-INV-001; OPEN-008 | 4B.2 |
| DEP-INV-401 | Finding/Hypothesis/Case/Incident/Hunt | CAP-INV-401 | intake | source need, observations, contradictions and return origin | active | Investigate / Command projection | yes | REQ-INV-006 | 4B.3A.1 |
| DEP-INV-402 | CAP-INV-313/328/346/362/379/397 | CAP-INV-401 | technical handoff | behavior, conditions, sources, limitations and provenance | active | Investigate | no | REQ-INV-006 | 4B.3A.1 |
| DEP-INV-403 | CAP-INV-401 | CAP-INV-402/403 | project/hypothesis | qualified objective and owner context | active | Investigate | yes | REQ-INV-006 | 4B.3A.1 |
| DEP-INV-404 | Platform Settings Data Source/Parser/Health/Retention | CAP-INV-404/405/409/416 | projection/admin | readiness and schema projections remain Settings-owned | partial | Platform Settings | yes | REQ-PROD-055; OPEN-008 | Settings/Technique |
| DEP-INV-405 | Endpoint Agent telemetry capabilities | CAP-INV-404/416 | platform projection | declared source/platform support and limitations | partial | Endpoint Agent / Settings | yes | REQ-PROD-049,050; OPEN-008 | Endpoint/Settings |
| DEP-INV-406 | Shared Telemetry Event/Query/Search/Data Quality/Normalization | CAP-INV-404,405,414,415 | shared data | source events, search, quality and schema projections | partial | Shared | yes | REQ-PROD-019,020 | Shared/Technique |
| DEP-INV-407 | Detection Hypothesis / field review | CAP-INV-406..410 | authoring | objective, source and semantic dependencies | active | Investigate | yes | REQ-INV-006 | 4B.3A.1 |
| DEP-INV-408 | Studio Tool / Tool Call | CAP-INV-411/414 and optional assistance | execution/provenance | validation and replay execution identity | partial | Studio | yes | REQ-OBJ-009; OPEN-015 | Studio/Objects |
| DEP-INV-409 | Studio Automation Run | automated authoring/testing | automation/provenance | optional orchestration and attribution | partial | Studio | no essential dependency | REQ-AI-002,010; OPEN-015 | Studio/Objects |
| DEP-INV-410 | Artifact / Studio Dataset / Shared storage projection | CAP-INV-412/413 | test data | versioned test data and restrictions | partial | source owner | yes | REQ-OBJ-003; OPEN-014 | Objects/Studio |
| DEP-INV-411 | CAP-INV-406..410 | CAP-INV-411 | validation | selected draft, metadata and dependencies | active | Investigate | yes | REQ-INV-006 | 4B.3A.1 |
| DEP-INV-412 | CAP-INV-411 | CAP-INV-412..415 | tests/replay/review | validation results and diagnostics | active | Investigate | no | REQ-INV-006 | 4B.3A.1 |
| DEP-INV-413 | CAP-INV-412/413 | CAP-INV-414/415 | expected/actual comparison | scenarios, datasets and Expected Outcomes | active | Investigate | no | REQ-INV-006 | 4B.3A.1 |
| DEP-INV-414 | Event Search and authorized historical access | CAP-INV-414 | replay | period, environment, retention, gaps and source events | active | Investigate / Shared | yes | REQ-PROD-014,019 | 4B.1/3A.1 |
| DEP-INV-415 | Replay/Validation source events | CAP-INV-415 | match review | actual result, source context and uncertainty | active | Investigate / Shared | yes | REQ-INV-006 | 4B.3A.1 |
| DEP-INV-416 | Data Readiness/Tests/Replay/Review | CAP-INV-416 | coverage | evidence level, environments and gaps | active | Investigate | yes | REQ-INV-006 | 4B.3A.1 |
| DEP-INV-417 | CAP-INV-401..416 | CAP-INV-417 | provenance/handoff | versioned authoring chain and unresolved items | active | Investigate | yes | REQ-PROD-020; OPEN-015 | 4B.3A.1 |
| DEP-INV-418 | CAP-INV-417 | Future Phase 4B.3A.2 | future review/promotion | package only; no Approval, deployment or activation | planned | future owner / Govern as applicable | no | REQ-INV-006; OPEN-013 | 4B.3A.2 |
| DEP-INV-419 | Future runtime Detection | Command | runtime boundary | Command owns Detection, Signal, Alert and Incident | planned | Command | future | REQ-PROD-013,015 | 4B.3A.2/Command |
| DEP-INV-420 | Detection engine/language choice | future technical delivery | unresolved decision gap | no existing OPEN covers detection engines/languages; OPEN-005 remains forensic-only | open gap | Product Architecture | before implementation | REQ-INV-006 | decision phase |
| DEP-INV-421 | Threat Intelligence | future 4B.3B | phase boundary | no Intelligence object or CAP-INV-5xx in 4B.3A.1 | planned | future Investigate owner | no | REQ-INV-006 | 4B.3B |

All fifteen OPEN decisions remain unchanged. The engine/language gap is recorded without inventing a new decision or misusing OPEN-005.
