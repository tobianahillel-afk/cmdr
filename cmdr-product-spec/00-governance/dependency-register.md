---
id: dependency-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: registry
requirements: [REQ-PROD-006, REQ-PROD-009, REQ-PROD-019, REQ-PROD-020]
---
# Dependency Register

Functional/documentary dependencies only; no protocol, API, engine, language, provider, exchange format, storage architecture or production pipeline is chosen.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN | Review |
|---|---|---|---|---|---|---|---|---|---|
| DEP-001 | source material | all specifications | decision | mission, principles and boundaries | active | Product Architecture | yes | REQ-PROD-001..020 | continuous |
| DEP-002 | ownership/boundaries | all modules | boundary | prevent concurrent ownership | active | Product Architecture | yes | REQ-PROD-013..018 | continuous |
| DEP-003 | object/permission phases | all capabilities | model/security | final schemas and atomic permissions future | partial | Architecture/Security | before implementation | REQ-OBJ/SEC | Phase 7 |
| DEP-004 | Shared mechanisms | all products | shared | Search, Graph, Timeline, Linking, Versioning, Reporting, Export, Notifications, Collaboration, Trace, Activity and Recovery | partial | Shared | yes | REQ-PROD-019/020 | Technique |
| DEP-005 | OPEN-005 | Forensics | future engine | forensic engines only | open | Investigate | before delivery | OPEN-005 | Phase 8 |
| DEP-006 | OPEN-008 | sources/platform/runtime | platform | support and configured targets unresolved | open | Settings/Endpoint | yes | OPEN-008 | Settings/Technique |
| DEP-007 | OPEN-013 | class-2 actions | authority | step-up/Govern policy unresolved | open | Security | yes | OPEN-013 | Phase 7 |
| DEP-008 | OPEN-014 | Artifact/Material/dataset relations | object | relation/retention unresolved | open | Investigate | model blocking | OPEN-014 | Phase 7 |
| DEP-009 | OPEN-015 | Studio/Govern/Investigate runs | run model | Automation Run/Response Run bridge unresolved | open | Studio/Govern | yes | OPEN-015 | Phase 7 |
| DEP-010 | OPEN-017 | Detection runtime/language/portability | architecture | Detection implementation strategy unresolved | open | Product Architecture | before implementation | OPEN-017 | Phase 7/8 |
| DEP-011 | OPEN-018 | Intelligence ontology/interoperability/exchange | architecture | representation unresolved | open | Product Architecture | before object/exchange implementation | OPEN-018 | Objects/Technique |
| DEP-012 | OPEN-019 | dissemination/releasability/sharing/access | policy | audience, tenant, client, external and recall policy unresolved | open | Product Architecture/Security/Govern | before external sharing implementation | OPEN-019 | Permissions/Govern/Technique |
| DEP-INV-201 | Collection/Artifact | Analysis Workbench | source/provenance | authorized acquisition feeds analysis | active | Investigate/Endpoint | yes | REQ-INV-001 | 4B.2 |
| DEP-INV-301 | Artifact Management | CAP-INV-301..397 | analysis source | custody, versions and lineage | active | Investigate | yes | REQ-OBJ-003/004 | 4B.2 |
| DEP-INV-401 | Cases/Hunts/Findings/technical handoffs | CAP-INV-401..435 | Detection lifecycle | authoring through improvement | active | Investigate/source owners | yes | REQ-INV-006; OPEN-017 | 4B.3A |
| DEP-INV-432 | CAP-INV-435 | CAP-INV-501 | Intelligence intake | sourced candidate context | active | Investigate | no | REQ-INV-006; OPEN-018 | 4B.3B.1 |
| DEP-INV-501 | Case/Finding/Hunt/Incident/technical analysis | CAP-INV-501..518 | Intelligence foundations | need, sources, candidates, restrictions, provenance | active | Investigate/source owners | yes | OPEN-013/014/015/018 | 4B.3B.1 |
| DEP-INV-504 | Settings sources/providers/connectors/access/health | CAP-INV-504..506/519..535 | administrative projection | Settings retains source administration | partial | Settings | yes | OPEN-008/018/019 | Settings/Technique |
| DEP-INV-508 | Shared Entity/Graph/Linking | CAP-INV-508/514/516/519..537 | shared projection | canonical Entity and generic Graph remain Shared-owned | partial | Shared | yes | OPEN-018 | Shared/Objects |
| DEP-INV-509 | Analysis Workbench | CAP-INV-509/510/512/518/524 | technical knowledge | technical outputs feed assessments | active | Investigate | no | REQ-INV-001..005 | 4B.2/3B |
| DEP-INV-510 | Command Detection/Signal/Alert/Incident | CAP-INV-501/513/518/531/532/534/537 | operational projection | Sighting/change notification ≠ Signal/Alert | active | Command | yes | REQ-PROD-013/015 | Command/Objects |
| DEP-INV-511 | Detection Engineering | CAP-INV-507/512/518/531/534/535/537 | handoff | candidate/package informs a new Detection cycle; no rule automatic | active | Investigate Detection Engineering | no | REQ-INV-006; OPEN-017 | 4B.3A/3B |
| DEP-INV-512 | Studio Tool/Tool Call/Automation Run | CAP-INV-506/515/519..537 | execution/provenance | optional attributed automation | partial | Studio | no essential AI dependency | OPEN-015 | Studio/Objects |
| DEP-INV-513 | Shared Versioning/Reporting/Export/Notifications/Trace | CAP-INV-501..537 | shared | generic rendering, delivery mechanisms and lineage | partial | Shared | yes | REQ-PROD-019/020 | Shared/Technique |
| DEP-INV-516 | CAP-INV-518 | CAP-INV-519 | analysis intake | Analysis Handoff Package becomes Session input, not Report | active | Investigate | yes | OPEN-018 | 4B.3B.2 |
| DEP-INV-519 | CAP-INV-519..524 | CAP-INV-525..527 | analysis-to-product | assessments feed Product Plan/Draft/Review | active | Investigate | yes | OPEN-018 | 4B.3B.2 |
| DEP-INV-525 | CAP-INV-525..527 | CAP-INV-528/529 | product-to-dissemination | quality and release recommendation precede internal publication | active | Investigate/Security projections | yes | OPEN-019 | 4B.3B.2 |
| DEP-INV-528 | CAP-INV-528 | CAP-INV-529/533 | dissemination | internal publication or external preparation only | active | Investigate/Govern | yes | OPEN-019 | 4B.3B.2 |
| DEP-INV-530 | Indicator candidates | CAP-INV-530/531 | operationalization | definition/package only; runtime owner activates | active | Investigate/runtime owner | yes | OPEN-008/017/018/019 | 4B.3B.2/Technique |
| DEP-INV-532 | sources/Sightings/publications | CAP-INV-532..537 | monitoring/feedback/lifecycle | changes and feedback feed correction/improvement | active | Investigate/source owners | yes | OPEN-008/015/019 | 4B.3B.2 |
| DEP-INV-533 | CAP-INV-533 | Govern/Settings future execution | external sharing | Action Request/package only; no transmission | planned | Govern/destination owner | yes | OPEN-018/019 | Govern/Technique |
| DEP-INV-537 | CAP-INV-501..536 | Continuous Improvement Package | closure | complete lineage returns to Intake/Requirement/Case/Hunt/Detection/Settings/Studio | active | Investigate/destination owners | no mutation | REQ-PROD-020 | 4B.3B.2 |
| DEP-INV-CLOUD | OPEN-012 | Phase 4B closure | scope | Cloud Analysis remains open/not implemented | blocking | Product Architecture | yes for Phase 4B PASS | OPEN-012 | roadmap decision |
| DEP-INV-MOBILE | OPEN-011 | Phase 4B closure | scope | Mobile Forensics remains open/not implemented | blocking | Product Architecture | yes for Phase 4B PASS | OPEN-011 | roadmap decision |

Eighteen decisions are open. Phase 4B.3 closes; Phase 4B remains PARTIAL because Cloud/Mobile are neither completed nor explicitly deferred outside its approved scope.
