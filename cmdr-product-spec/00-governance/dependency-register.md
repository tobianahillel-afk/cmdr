---
id: dependency-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-07
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
| DEP-005 | OPEN-005 | Analysis Workbench and Forensics | future engine | forensic engines only; not Detection, Intelligence or Cloud provider selection | open | Investigate | before forensic delivery | OPEN-005 | Phase 8 |
| DEP-006 | OPEN-008 | platform/source support | platform | source, sensor, endpoint and Cloud provider support unresolved | open | Endpoint Agent / Settings | before delivery | OPEN-008 | Settings/Technique |
| DEP-007 | OPEN-013 | class-2 actions | authority | default step-up/Govern policy unresolved | open | Security | yes | OPEN-013 | Phase 7 |
| DEP-008 | OPEN-014 | Artifact/Attachment/dataset/material relations | object | relation and retention semantics unresolved | open | Investigate | model blocking | OPEN-014 | Phase 7 |
| DEP-009 | OPEN-015 | Studio/Govern/Investigate runs | run model | Automation Run / Response Run bridge unresolved | open | Studio + Govern | yes | OPEN-015 | Phase 7 |
| DEP-010 | OPEN-017 | Detection runtime, target language and portability | product architecture | Detection-only strategy undecided | open | Product Architecture | before Detection implementation | OPEN-017 | Phase 7/8 |
| DEP-011 | OPEN-018 | Threat Intelligence ontology, interoperability and exchange | product architecture | canonical ontology, standards, projections and external representations remain undecided | open | Product Architecture | before object/exchange implementation | OPEN-018 | Objects/Technique |
| DEP-012 | OPEN-019 | dissemination, releasability, sharing and consumer access | policy | internal audiences, cross-tenant/client/external sharing, publication authority and recall remain unresolved | open | Product Architecture / Security / Govern | before external sharing implementation | OPEN-019 | Permissions/Govern/Technique |
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
| DEP-INV-504 | Platform Settings providers/sources/connectors/access/health | CAP-INV-504..506 and CAP-INV-519..535 | source projection | source administration and secrets remain Settings-owned | partial | Platform Settings | yes | REQ-PROD-055; OPEN-008/018/019 | Settings/Technique |
| DEP-INV-505 | CAP-INV-504/506/513/514 | CAP-INV-505/515/521 | assessment | reliability, credibility, corroboration, contradiction and confidence remain distinct | active | Investigate | yes | REQ-PROD-020; OPEN-018 | 4B.3B |
| DEP-INV-506 | Artifact/source material / Studio Tool | CAP-INV-506 | material intake | preserve original, restrictions, extraction and ambiguity | active/partial | source owner / Studio | yes | REQ-OBJ-003; OPEN-014/015/018 | Objects/Studio |
| DEP-INV-507 | CAP-INV-506 and source observations | CAP-INV-507..513 | candidate knowledge | Observables, Indicators, entities, malware, infrastructure, campaigns, TTP and Sightings remain candidates | active | Investigate | yes | REQ-PROD-014; OPEN-018 | 4B.3B.1 |
| DEP-INV-508 | Shared Entity/Graph/Linking | CAP-INV-508/514/516 and CAP-INV-519..537 | shared projection | canonical Entity and generic Graph remain Shared-owned | partial | Shared | yes | REQ-PROD-019,020; OPEN-018 | Shared/Objects |
| DEP-INV-509 | Analysis Workbench handoffs | CAP-INV-509/510/512/518/524 | technical knowledge | Artifact, behavior, infrastructure observations and provenance feed knowledge and assessments | active | Investigate | no | REQ-INV-001..005 | 4B.2/3B |
| DEP-INV-510 | Command Detection/Signal/Alert/Incident | CAP-INV-501/513/518/531/532/534/537 | operational projection | Sighting, Watchlist Match and Change Notification do not replace Command objects | active | Command | yes | REQ-PROD-013,015 | Command/Objects |
| DEP-INV-511 | Detection Engineering | CAP-INV-501/507/512/518/531/534/535/537 | cross-domain handoff | Indicator/TTP candidates and packages may inform a new Detection cycle but create no rule automatically | active | Investigate Detection Engineering | no | REQ-INV-006; OPEN-017 | 4B.3A/3B |
| DEP-INV-512 | Studio Tool/Tool Call/Automation Run | CAP-INV-506/515/516/518 and CAP-INV-519..537 | execution/provenance | optional extraction, comparison, drafting, monitoring and attributed automation | partial | Studio | no essential AI dependency | OPEN-015 | Studio/Objects |
| DEP-INV-513 | Shared Versioning/Comparison/Timeline/Trace/Activity/Reporting/Export/Notifications | CAP-INV-501..537 | shared | generic lineage, comparison, rendering, delivery mechanisms and time navigation | partial | Shared | yes | REQ-PROD-019,020 | Shared/Technique |
| DEP-INV-514 | CAP-INV-507..515 | CAP-INV-516/517 | lifecycle | dedup, versions, supersession, expiry and revocation preserve source restrictions | active | Investigate | yes | OPEN-013/014/018 | 4B.3B.1 |
| DEP-INV-515 | CAP-INV-501..517 | CAP-INV-518 | provenance/handoff | complete foundation lineage and unresolved issues | active | Investigate / source owners | yes | REQ-PROD-020; OPEN-014/015/018 | 4B.3B.1 |
| DEP-INV-516 | CAP-INV-518 | CAP-INV-519 | analysis intake | Analysis Handoff Package becomes Session input, not Report or publication | active | Investigate | yes | OPEN-018 | 4B.3B.2 |
| DEP-INV-517 | Govern external sharing authority | CAP-INV-533 and future execution | governance | external release and sharing are prepared but never executed by Investigate | active/planned | Govern / destination owner | yes for execution | REQ-SEC-001,002; OPEN-019 | Govern/Technique |
| DEP-INV-519 | CAP-INV-519..524 | CAP-INV-525..527 | analysis-to-product | sourced assessments feed Product Plan, Draft and Review | active | Investigate | yes | OPEN-018 | 4B.3B.2 |
| DEP-INV-525 | CAP-INV-525..527 | CAP-INV-528/529 | product-to-dissemination | quality and Release Recommendation precede internal publication | active | Investigate / Security projections | yes | OPEN-019 | 4B.3B.2 |
| DEP-INV-528 | CAP-INV-528 | CAP-INV-529/533 | dissemination | internal publication or external preparation only | active | Investigate / Govern | yes | OPEN-019 | 4B.3B.2 |
| DEP-INV-530 | Indicator candidates | CAP-INV-530/531 | operationalization | definition/package only; runtime owner activates or deploys | active | Investigate / runtime owner | yes | OPEN-008/017/018/019 | 4B.3B.2/Technique |
| DEP-INV-532 | sources/Sightings/publications | CAP-INV-532..537 | monitoring/feedback/lifecycle | changes and feedback feed correction and improvement | active | Investigate / source owners | yes | OPEN-008/015/019 | 4B.3B.2 |
| DEP-INV-533 | CAP-INV-533 | Govern/Settings future execution | external sharing | Action Request/package only; no transmission | planned | Govern / destination owner | yes | OPEN-018/019 | Govern/Technique |
| DEP-INV-537 | CAP-INV-501..536 | Continuous Improvement Package | closure | complete lineage returns to Intake/Requirement/Case/Hunt/Detection/Settings/Studio | active | Investigate / destination owners | no active mutation | REQ-PROD-020 | 4B.3B.2 |
| DEP-INV-601 | Case/Incident/Finding/Hunt/Signal/Cloud source/Artifact | CAP-INV-601..603 | Cloud intake/scope | origin, scope, periods, sources, access, restrictions and return origin | active | Investigate / source owners | yes | REQ-INV-001/006; OPEN-008/012/013/014/015 | 4B.4A |
| DEP-INV-604 | Platform Settings providers/sources/connectors/health/schemas | CAP-INV-604..614 | Cloud source projection | Settings retains all Cloud administration; Investigate consumes authorized observations only | partial | Platform Settings | yes | REQ-PROD-055; OPEN-008/012 | Settings/Technique |
| DEP-INV-605 | Cloud identities/roles/policies/resources | CAP-INV-605/606 | identity/IAM analysis | identity, assignment, declared policy and effective-permission candidate remain distinct | active | Investigate / source owners | yes | OPEN-012/013/014 | 4B.4A |
| DEP-INV-607 | Cloud audit/activity materials | CAP-INV-607/616 | event/timeline | actor, event time, gaps, duplicates and lateness preserve source uncertainty | active/partial | Investigate / Settings | yes | OPEN-008/012 | 4B.4A |
| DEP-INV-608 | inventory/configuration observations | CAP-INV-608..614 | workload/data analysis | observed configuration, workloads, network, storage and sensitive material remain candidates | active | Investigate | yes | OPEN-012/014 | 4B.4A |
| DEP-INV-609 | Cloud compute/workload observations | Endpoint/Memory/Disk/Network owners | specialist handoff | Cloud Analysis prepares scope and evidence context but executes no acquisition or command | active | destination owners | no active mutation | REQ-INV-001..005; OPEN-008/013/014 | 4B.4A/4B.2 |
| DEP-INV-615 | CAP-INV-604..614 | CAP-INV-615/616 | interpretation | anomalies, Hypotheses and correlations preserve supporting, contradicting and missing elements | active | Investigate | yes | REQ-PROD-020; OPEN-012/014/015 | 4B.4A |
| DEP-INV-617 | CAP-INV-601..616 | Evidence/Finding/Detection/TI/Collection/Govern owners | handoff | packages are candidates, drafts or requests and create no destination object or action automatically | active | Investigate / destination owners | yes | OPEN-013/014/015/017 | 4B.4A |
| DEP-INV-618 | CAP-INV-601..617 plus Studio/Shared lineage | provenance | preserve sources, permissions, Tools, Runs, errors, human decisions and versions | active/partial | Investigate / Studio / Shared | yes | REQ-PROD-020; OPEN-014/015 | 4B.4A |
| DEP-INV-CLOUD | OPEN-012 | Cloud provider/service delivery | scope/implementation | provider priorities, service support, SaaS, Cloud evidence and cross-tenant strategy remain unresolved; provider-neutral functional documentation is verified | open | Product Architecture | before Cloud implementation, not before provider-neutral documentation | OPEN-012 | roadmap/Technique |
| DEP-INV-MOBILE | OPEN-011 | Phase 4B.4 and Phase 4B closure | scope | Mobile Forensics remains open and not implemented | blocking | Product Architecture | yes for Phase 4B PASS | OPEN-011 | next subphase |

Eighteen decisions are open. OPEN-017 remains Detection-only; OPEN-018 governs Threat Intelligence ontology and interoperability; OPEN-019 governs Intelligence dissemination, releasability, sharing and access; OPEN-012 governs unresolved Cloud provider, service and delivery strategy. Phase 4B.4A is **PASS AFTER POST-PUBLICATION VERIFICATION**. Phase 4B.4 and Phase 4B remain PARTIAL because Mobile Forensics is not started.
