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

This register records functional and documentary dependencies only. It chooses no protocol, API, engine, plugin, debugger, hypervisor, provider, command, format or storage design.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN | Review |
|---|---|---|---|---|---|---|---|---|---|
| DEP-001 | source material | product specifications | decision | apply mission, principles and boundaries | active | Product Architecture | yes | REQ-PROD-001..020 | continuous |
| DEP-002 | product boundaries / ownership | all modules | boundary | prevent concurrent ownership | active | Product Architecture | yes | REQ-PROD-013..018 | continuous |
| DEP-003 | object/permission phases | capabilities | model/security | final schemas and atomic permissions are future | partial | Architecture/Security | before implementation | REQ-OBJ/SEC | Phase 7 |
| DEP-004 | Shared mechanisms | all products | shared | Jobs, Trace, Activity, Linking, Export, Versioning, Notifications | partial | Shared | yes | REQ-PROD-019,020 | Technique |
| DEP-005 | OPEN-005 | Analysis Workbench/Dynamic/Reverse/Debugger/Forensics | future engine | no engine, plugin or framework selected | open | Investigate | before delivery | OPEN-005 | Phase 8 |
| DEP-006 | OPEN-008 | Endpoint platform support | platform | initial supported platforms unresolved | open | Endpoint Agent / Settings | before delivery | OPEN-008 | Phase 4D/8 |
| DEP-007 | OPEN-013 | class-2 actions | authority | default step-up/Govern policy unresolved | open | Security | yes | OPEN-013 | Phase 7 |
| DEP-008 | OPEN-014 | Artifact relations | object | Artifact versus Attachment unresolved | open | Investigate | model blocking | OPEN-014 | Phase 7 |
| DEP-009 | OPEN-015 | Studio/Govern/Investigate runs | run model | Automation Run / Response Run bridge unresolved | open | Studio + Govern | yes | OPEN-015 | Phase 7 |
| DEP-INV-301 | CAP-INV-105/213/214 | CAP-INV-301..313 | object/provenance | Artifact, custody and provenance sources | active | Investigate | yes | REQ-PROD-014,020 | 4B.2B.1 |
| DEP-INV-314 | CAP-INV-301..313 | CAP-INV-314..328 | static-to-dynamic handoff | static results, Derived Artifacts and provenance | active | Investigate | yes | REQ-INV-002,005 | 4B.2B.2A |
| DEP-INV-329 | CAP-INV-301/307/311/312/313 and CAP-INV-314/324/327/328 | CAP-INV-329..346 | reverse/debug context | static/dynamic results and provenance | active | Investigate | yes | REQ-INV-003,004 | 4B.2B.2B |
| DEP-INV-347 | CAP-INV-207 | CAP-INV-347 | acquisition handoff | consume Memory Image and declared context; no acquisition | active | Investigate / Endpoint Agent | yes | REQ-INV-001; OPEN-008 | 4B.2B.3A |
| DEP-INV-348 | CAP-INV-105 | CAP-INV-347..362 | object source | Memory Image and Artifact ownership/lineage | active | Investigate | yes | REQ-OBJ-003 | Objects |
| DEP-INV-349 | CAP-INV-213/214 | CAP-INV-349/361 | custody/provenance | acquisition integrity, transfers and trace | active | Investigate / Shared | yes | REQ-PROD-020 | 4B.2A |
| DEP-INV-350 | CMDR Studio Tool/Tool Call | CAP-INV-347..362 | tool/provenance | selection, version, execution and attribution | partial | Studio | yes | REQ-OBJ-009; OPEN-015 | Studio/Objects |
| DEP-INV-351 | CMDR Studio Automation Run | automated memory analysis | automation/provenance | orchestration remains Studio-owned | partial | Studio | no essential dependency | REQ-AI-002; OPEN-015 | Studio/Objects |
| DEP-INV-352 | Platform Settings Fleet/Policy/storage/retention/health | CAP-INV-347..361 | projection/admin | support, restrictions and administrative state | partial | Platform Settings | yes | REQ-PROD-055; OPEN-008 | Settings/Technique |
| DEP-INV-353 | Shared Jobs/Trace/Activity/Timeline/Linking/Export/Recovery | CAP-INV-347..362 | shared mechanism | progress, provenance, correlation and export | partial | Shared | yes | REQ-PROD-019,020 | Shared/Technique |
| DEP-INV-354 | CAP-INV-311 Derived Artifact Management | CAP-INV-360 | artifact handoff | extracted content and lineage | active | Investigate | yes | REQ-OBJ-003 | 4B.2B.1 |
| DEP-INV-355 | CAP-INV-329..346 | CAP-INV-357/360/362 | reverse/debug handoff | extracted content and anomaly context | active | Investigate | no | REQ-INV-003,004 | 4B.2B.2B |
| DEP-INV-356 | CAP-INV-107/108 | CAP-INV-362 | Evidence handoff | qualification remains Evidence owner | active | Investigate | yes | REQ-OBJ-004 | 4B.1 |
| DEP-INV-357 | CAP-INV-109 | CAP-INV-362 | Finding handoff | draft remains unconfirmed | active | Investigate | yes | REQ-PROD-016 | 4B.1 |
| DEP-INV-358 | Future Detection Engineering | CAP-INV-362 | future handoff | package only; no rule creation/deployment | planned | future owner | no | REQ-INV-006 | 4B.3 |
| DEP-INV-359 | Memory network observations | future Phase 4B.2B.3B | future handoff | no full Network Forensics now | planned | Investigate | no | REQ-INV-001 | 4B.2B.3B |
| DEP-INV-360 | Disk and Filesystem Forensics | future Phase 4B.2B.3B | future boundary | no disk/filesystem capability now | planned | Investigate | no | REQ-INV-001 | 4B.2B.3B |
| DEP-INV-361 | Govern authority | any real-target or sensitive action | boundary | block or route; no local Endpoint action | active | Govern | class 3/4 | REQ-SEC-002 | 4C/5 |

OPEN-005, OPEN-008, OPEN-013, OPEN-014 and OPEN-015 remain open. OPEN-011/012 remain future. No low-level dependency is created.
