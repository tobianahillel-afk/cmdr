---
id: dependency-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-019
---
# Dependency Register

This register records functional and documentary dependencies only. It chooses no protocol, API, engine, debugger, hypervisor, provider, command or storage design.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN | Review |
|---|---|---|---|---|---|---|---|---|---|
| DEP-001 | source material | product specifications | decision | apply mission, principles and boundaries | active | Product Architecture | yes | REQ-PROD-001..020 | continuous |
| DEP-002 | product boundaries / ownership | all modules | boundary | prevent concurrent ownership | active | Product Architecture | yes | REQ-PROD-013..018 | continuous |
| DEP-003 | object/permission phases | capabilities | model/security | final schemas and atomic permissions are future | partial | Architecture/Security | before implementation | REQ-OBJ/SEC | Phase 7 |
| DEP-004 | Shared mechanisms | all products | shared | Jobs, Trace, Activity, Linking, Export, Versioning, Notifications | partial | Shared | yes | REQ-PROD-019,020 | Technique |
| DEP-005 | OPEN-005 | Analysis Workbench/Dynamic/Reverse/Debugger/Forensics | future engine | no engine, debugger or hypervisor selected | open | Investigate | before delivery | OPEN-005 | Phase 8 |
| DEP-006 | OPEN-013 | class-2 actions | authority | default step-up/Govern policy unresolved | open | Security | yes | OPEN-013 | Phase 7 |
| DEP-007 | OPEN-014 | Artifact relations | object | Artifact versus Attachment unresolved | open | Investigate | model blocking | OPEN-014 | Phase 7 |
| DEP-008 | OPEN-015 | Studio/Govern/Investigate runs | run model | Automation Run / Response Run bridge unresolved | open | Studio + Govern | yes | OPEN-015 | Phase 7 |
| DEP-INV-301 | CAP-INV-105/213/214 | CAP-INV-301..313 | object/provenance | Artifact, custody and provenance sources | active | Investigate | yes | REQ-PROD-014,020 | 4B.2B.1 |
| DEP-INV-314 | CAP-INV-301..313 | CAP-INV-314..328 | static-to-dynamic handoff | static results, Derived Artifacts and provenance | active | Investigate | yes | REQ-INV-002,005 | 4B.2B.2A |
| DEP-INV-329 | CAP-INV-301/307/311/312/313 | CAP-INV-329 | static-to-reverse intake | Artifact, static results, Derived Artifacts and provenance | active | Investigate | yes | REQ-INV-002,003 | 4B.2B.2B |
| DEP-INV-330 | CAP-INV-314/324/327/328 | CAP-INV-329 | dynamic-to-reverse intake | Runtime Artifacts, Runs, observations and provenance | active | Investigate | yes | REQ-INV-003,005 | 4B.2B.2B |
| DEP-INV-331 | CAP-INV-105 Artifact Management | CAP-INV-329..346 | object source | source Artifact, versions, restrictions and lineage | active | Investigate | yes | REQ-OBJ-003 | Objects |
| DEP-INV-332 | CMDR Studio Tool/Tool Call | CAP-INV-329..346 | tool/provenance | selection, execution, version and attribution | partial | Studio | yes | REQ-OBJ-009; OPEN-015 | Studio/Objects |
| DEP-INV-333 | CMDR Studio Automation Run | automated reverse assistance | automation/provenance | agentic orchestration remains Studio-owned | partial | Studio | no essential dependency | REQ-AI-002; OPEN-015 | Studio/Objects |
| DEP-INV-334 | Platform Settings environment/health/policy | CAP-INV-329/339/344 | environment | availability, isolation and restrictions | partial | Platform Settings | yes before execution | REQ-PROD-017; REQ-SEC-001 | Settings/Technique |
| DEP-INV-335 | Shared Workbench/Inspector/Graph/Trace/Versioning/Recovery | CAP-INV-331..346 | shared mechanism | navigation, graph, provenance, collaboration and recovery | partial | Shared | yes | REQ-PROD-019,020; REQ-UX-002 | Shared/Technique |
| DEP-INV-336 | CAP-INV-311 Derived Artifact | CAP-INV-342/344 | artifact handoff | extraction and isolated experiment outputs | active | Investigate | yes | REQ-OBJ-003 | 4B.2B.1 |
| DEP-INV-337 | CAP-INV-107/108 | CAP-INV-346 | Evidence handoff | qualification remains Evidence owner | active | Investigate | yes | REQ-OBJ-004 | 4B.1 |
| DEP-INV-338 | CAP-INV-109 | CAP-INV-346 | Finding handoff | draft remains unconfirmed until owner review | active | Investigate | yes | REQ-PROD-016 | 4B.1 |
| DEP-INV-339 | Future Detection Engineering | CAP-INV-346 | future handoff | knowledge package only; no rule creation/deployment | planned | Investigate/future owner | no | REQ-INV-006 | 4B.3 |
| DEP-INV-340 | OPEN-013 | CAP-INV-330/337/339/340/344/345/346 | authority | reversible class-2 defaults unresolved | open | Security | yes | OPEN-013 | Phase 7 |
| DEP-INV-341 | OPEN-014 | CAP-INV-329/342/346 | object ambiguity | Artifact versus Attachment when exporting/citing | open | Investigate | model blocking | OPEN-014 | Phase 7 |
| DEP-INV-342 | OPEN-015 | CAP-INV-329..346 | run/provenance | Automation Run distinct from Debugger/Reverse/Response sessions | open | Studio + Govern | yes | OPEN-015 | Phase 7 |
| DEP-INV-343 | Govern authority | real-target requests | boundary | no direct debugger or patch on real target | active | Govern | class 3/4 | REQ-SEC-002 | 4C/5 |
| DEP-INV-344 | Phase 4B.2B.3 | Forensics | future boundary | Memory/Disk/Filesystem/Network Forensics not started | planned | Investigate | no | REQ-INV-001 | next phase |
| DEP-INV-345 | Phase 4B.3 | Detection Engineering | future boundary | no rule engineering or deployment in this phase | planned | Investigate | no | REQ-INV-006 | future |

OPEN-005, OPEN-013, OPEN-014 and OPEN-015 remain open. OPEN-011/012 remain future. No low-level dependency is created.
