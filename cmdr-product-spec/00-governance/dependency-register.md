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

This register records functional and documentary dependencies only. It chooses no protocol, API, engine, hypervisor, provider, command or storage design.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN | Review |
|---|---|---|---|---|---|---|---|---|---|
| DEP-001 | source material | product specifications | decision | apply mission, principles and boundaries | active | Product Architecture | yes | REQ-PROD-001..020 | continuous |
| DEP-002 | product boundaries / ownership | all modules | boundary | prevent concurrent ownership | active | Product Architecture | yes | REQ-PROD-013..018 | continuous |
| DEP-003 | object/permission phases | capabilities | model/security | final schemas and atomic permissions are future | partial | Architecture/Security | before implementation | REQ-OBJ/SEC | Phase 7 |
| DEP-004 | Shared mechanisms | all products | shared | Jobs, Trace, Activity, Linking, Export, Versioning, Notifications | partial | Shared | yes | REQ-PROD-019,020 | Technique |
| DEP-005 | OPEN-005 | Analysis Workbench/Dynamic Sandbox | future engine | no engine or hypervisor selected | open | Investigate | before delivery | OPEN-005 | Phase 8 |
| DEP-006 | OPEN-013 | class-2 actions | authority | default step-up/Govern policy unresolved | open | Security | yes | OPEN-013 | Phase 7 |
| DEP-007 | OPEN-014 | Artifact relations | object | Artifact versus Attachment unresolved | open | Investigate | model blocking | OPEN-014 | Phase 7 |
| DEP-008 | OPEN-015 | Studio/Govern/Investigate runs | run model | Automation Run / Response Run bridge unresolved | open | Studio + Govern | yes | OPEN-015 | Phase 7 |
| DEP-INV-301 | CAP-INV-105/213/214 | CAP-INV-301..313 | object/provenance | Artifact, custody and provenance sources | active | Investigate | yes | REQ-PROD-014,020 | 4B.2B.1 |
| DEP-INV-302 | CAP-INV-107/108/109 | CAP-INV-313 | handoff | explicit Evidence/Finding qualification | active | Investigate | yes | REQ-OBJ-004 | 4B.2B.1 |
| DEP-INV-303 | Studio Tool/Tool Call/Automation Run | CAP-INV-303/312 | tool/provenance | selection, execution and attribution | partial | Studio | yes | REQ-OBJ-009; OPEN-015 | Phase 7 |
| DEP-INV-304 | Settings environments/providers/secrets | CAP-INV-303/312 | administration | Investigate consumes only projections | partial | Platform Settings | yes | REQ-PROD-014 | Settings/Technique |
| DEP-INV-314 | CAP-INV-301..313 | CAP-INV-314/324/327/328 | static-to-dynamic handoff | static results, Derived Artifacts and provenance | active | Investigate | yes | REQ-INV-002,005 | 4B.2B.2A |
| DEP-INV-315 | Settings Sandbox Environment/Health/Policies | CAP-INV-315/317/326 | environment | availability, isolation declaration, network policy and reset | partial | Platform Settings | yes | REQ-PROD-017; REQ-SEC-001 | Settings/Technique |
| DEP-INV-316 | Studio Tool/Tool Call/Workflow/Automation Run | CAP-INV-317/323/327 | execution/provenance | Tools and orchestration remain Studio-owned | partial | Studio | yes | REQ-OBJ-009; OPEN-015 | Phase 7 |
| DEP-INV-317 | Shared Background Jobs/Notifications | CAP-INV-317 | shared mechanism | queue, progress, cancel and partial status | partial | Shared | yes | REQ-PROD-019 | Technique |
| DEP-INV-318 | Shared Timeline/Trace/Activity/Linking | CAP-INV-318..328 | shared mechanism | sourced behavioral events and navigation | partial | Shared | yes | REQ-PROD-020 | Trust/Technique |
| DEP-INV-319 | Govern authority | dynamic module | boundary | real-target action leaves sandbox and routes to governed products | active | Govern | class 3/4 | REQ-SEC-002 | 4C/5 |
| DEP-INV-320 | CAP-INV-107/108/109 | CAP-INV-328 | handoff | results remain candidates/drafts | active | Investigate | yes | REQ-OBJ-004 | 4B.1 |
| DEP-INV-321 | Phase 4B.2B.2B | future Reverse/Debugger handoff | future boundary | only context is prepared now | planned | Investigate | no | REQ-INV-003,004 | future |
| DEP-INV-322 | Phase 4B.2B.3 | future forensics handoff | future boundary | no forensic capability in this phase | planned | Investigate | no | REQ-INV-001 | future |

OPEN-005, OPEN-013, OPEN-014 and OPEN-015 remain open. No low-level dependency is created.
