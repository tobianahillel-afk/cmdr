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

This register records functional and documentary dependencies only. It chooses no protocol, API, engine, language, parser implementation, model, provider, command, exchange format, storage architecture or production pipeline.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN | Review |
|---|---|---|---|---|---|---|---|---|---|
| DEP-001 | source material | product specifications | decision | apply mission, principles and boundaries | active | Product Architecture | yes | REQ-PROD-001..020 | continuous |
| DEP-002 | product boundaries / ownership | all modules | boundary | prevent concurrent ownership | active | Product Architecture | yes | REQ-PROD-013..018 | continuous |
| DEP-003 | object/permission phases | capabilities | model/security | final schemas and atomic permissions are future | partial | Architecture/Security | before implementation | REQ-OBJ/SEC | Phase 7 |
| DEP-004 | Shared mechanisms | all products | shared | Jobs, Trace, Activity, Linking, Search, Export, Reporting, Versioning, Collaboration and Recovery | partial | Shared | yes | REQ-PROD-019,020 | Technique |
| DEP-005 | OPEN-005 | Analysis Workbench and Forensics | future engine | forensic engines only; not Detection or Intelligence | open | Investigate | before forensic delivery | OPEN-005 | Phase 8 |
| DEP-006 | OPEN-008 | platform/source support | platform | source, sensor and endpoint platform support unresolved | open | Endpoint Agent / Settings | before delivery | OPEN-008 | Phase 4D/8 |
| DEP-007 | OPEN-013 | class-2 actions | authority | default step-up/Govern policy unresolved | open | Security | yes | OPEN-013 | Phase 7 |
| DEP-008 | OPEN-014 | Artifact/Attachment/dataset/material relations | object | relation and retention semantics unresolved | open | Investigate | model blocking | OPEN-014 | Phase 7 |
| DEP-009 | OPEN-015 | Studio/Govern/Investigate runs | run model | Automation Run / Response Run bridge unresolved | open | Studio + Govern | yes | OPEN-015 | Phase 7 |
| DEP-010 | OPEN-017 | Detection runtime, target language and portability | product architecture | Detection-only strategy undecided | open | Product Architecture | before Detection implementation | OPEN-017 | Phase 7/8 |
| DEP-011 | OPEN-018 | Threat Intelligence ontology, interoperability and exchange | product architecture | canonical ontology, standards, projections and external representations remain undecided | open | Product Architecture | before object/exchange implementation | OPEN-018 | Objects/Technique/4B.3B.2 |
| DEP-INV-201 | Collection and Live Response | Analysis Workbench families | source/provenance | authorized acquisition outputs feed analysis | active | Investigate / Endpoint | yes | REQ-INV-001 | 4B.2 |
| DEP-INV-301 | Artifact Management | CAP-INV-301..397 | analysis source | Artifact, versions, custody and lineage | active | Investigate | yes | REQ-OBJ-003,004 | 4B.2 |
| DEP-INV-314 | Static Analysis | Dynamic Sandbox | handoff | Derived Artifacts and static context | active | Investigate | no | REQ-INV-002,005 | 4B.2 |
| DEP-INV-329 | Static/Dynamic | Reverse and Debugger | handoff | code and runtime analysis context | active | Investigate | no | REQ-INV-003,004 | 4B.2 |
| DEP-INV-347 | Collection/Artifact | Memory Forensics | source/custody | Memory Image and acquisition context | active | Investigate / Endpoint | yes | REQ-INV-001; OPEN-008 | 4B.2 |
| DEP-INV-363 | Collection/Artifact | Disk and Filesystem Forensics | source/custody | Disk Image and acquisition context | active | Investigate / Endpoint | yes | REQ-INV-001; OPEN-008 | 4B.2 |
| DEP-INV-380 | Collection/Artifact | Network Forensics | source/custody | Capture Artifact and acquisition context | active | Investigate / Endpoint | yes | REQ-INV-001; OPEN-008 | 4B.2 |
| DEP-INV-401 | Finding/Hypothesis/Case/Incident/Hunt and technical handoffs | CAP-INV-401..417 | Detection authoring | need, sources, authoring, validation, replay, review and coverage | active | Investigate | yes | REQ-INV-006; OPEN-008/013/014/015 | 4B.3A.1 |
| DEP-INV-418 | CAP-INV-417 | CAP-INV-418..435 | Detection lifecycle | review, readiness, governed runtime coordination and feedback | active | Investigate / source owners | yes | REQ-INV-006; OPEN-017 | 4B.3A.2 |
| DEP-INV-432 | CAP-INV-435 | CAP-INV-501 | Intelligence intake | candidate context, sources, Sightings, limitations and return origin | active | Investigate | no | REQ-INV-006; OPEN-018 | 4B.3B.1 |
| DEP-INV-501 | Case/Finding/Hunt/Incident/technical analysis | CAP-INV-501..503 | Intelligence intake/project | need, scope, questions, sources, restrictions and return origin | active | Investigate / Command projection | yes | REQ-PROD-014; OPEN-013/014/015/018 | 4B.3B.1 |
| DEP-INV-502 | CAP-INV-501 | CAP-INV-502 | requirement | qualified need becomes versioned Intelligence Requirement, not collection execution | active | Investigate | yes | REQ-PROD-014; OPEN-018 | 4B.3B.1 |
| DEP-INV-503 | CAP-INV-502 | CAP-INV-503 | workspace | Requirement, origins, owners and future consumers form a Knowledge Project | active | Investigate | yes | REQ-PROD-020 | 4B.3B.1 |
| DEP-INV-504 | Platform Settings providers/sources/connectors/access/health | CAP-INV-504..506 | source projection | source administration and secrets remain Settings-owned | partial | Platform Settings | yes | REQ-PROD-055; OPEN-008/018 | Settings/Technique |
| DEP-INV-505 | CAP-INV-504/506/513/514 | CAP-INV-505/515 | assessment | reliability, credibility, corroboration, contradiction and confidence remain distinct | active | Investigate | yes | REQ-PROD-020; OPEN-018 | 4B.3B.1 |
| DEP-INV-506 | Artifact/source material / Studio Tool | CAP-INV-506 | material intake | preserve original, restrictions, extraction and ambiguity | active/partial | source owner / Studio | yes | REQ-OBJ-003; OPEN-014/015/018 | Objects/Studio |
| DEP-INV-507 | CAP-INV-506 and source observations | CAP-INV-507..513 | candidate knowledge | Observables, Indicators, entities, malware, infrastructure, campaigns, TTP and Sightings remain candidates | active | Investigate | yes | REQ-PROD-014; OPEN-018 | 4B.3B.1 |
| DEP-INV-508 | Shared Entity/Graph/Linking | CAP-INV-508/514/516 | shared projection | canonical Entity and generic Graph remain Shared-owned | partial | Shared | yes | REQ-PROD-019,020; OPEN-018 | Shared/Objects |
| DEP-INV-509 | Analysis Workbench handoffs | CAP-INV-509/510/512/518 | technical knowledge | Artifact, behavior, infrastructure observations and provenance feed candidate knowledge | active | Investigate | no | REQ-INV-001..005 | 4B.2/3B.1 |
| DEP-INV-510 | Command Detection/Signal/Alert/Incident | CAP-INV-501/513/518 | operational projection | Sighting and Intelligence knowledge do not replace Command objects | active | Command | yes | REQ-PROD-013,015 | Command/Objects |
| DEP-INV-511 | Detection Engineering | CAP-INV-501/507/512/518 | cross-domain handoff | Indicator/TTP candidates may inform a new Detection cycle but create no rule automatically | active | Investigate Detection Engineering | no | REQ-INV-006 | 4B.3A/3B.1 |
| DEP-INV-512 | Studio Tool/Tool Call/Automation Run | CAP-INV-506/515/516/518 | execution/provenance | optional extraction, comparison and attribution of automation | partial | Studio | no essential AI dependency | OPEN-015 | Studio/Objects |
| DEP-INV-513 | Shared Versioning/Comparison/Timeline/Trace/Activity | CAP-INV-501..518 | shared | generic lineage, comparison and time navigation | partial | Shared | yes | REQ-PROD-019,020 | Shared/Technique |
| DEP-INV-514 | CAP-INV-507..515 | CAP-INV-516/517 | lifecycle | dedup, versions, supersession, expiry and revocation preserve source restrictions | active | Investigate | yes | OPEN-013/014/018 | 4B.3B.1 |
| DEP-INV-515 | CAP-INV-501..517 | CAP-INV-518 | provenance/handoff | complete foundation lineage and unresolved issues | active | Investigate / source owners | yes | REQ-PROD-020; OPEN-014/015/018 | 4B.3B.1 |
| DEP-INV-516 | CAP-INV-518 | future 4B.3B.2 | future analysis | candidate package only; no Report, attribution, dissemination or operationalization | planned | future Investigate owner | no | OPEN-018 | 4B.3B.2 |
| DEP-INV-517 | Govern external sharing authority | future 4B.3B.2 | future governance | information release and external sharing are not performed now | planned | Govern | future | REQ-SEC-001,002 | 4B.3B.2/Govern |

Seventeen decisions are open. OPEN-018 is the canonical open decision for Threat Intelligence ontology, interoperability and future exchange. No option is selected; OPEN-017 remains Detection-only.
