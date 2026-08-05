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

This register records functional and documentary dependencies only. It chooses no protocol, API, engine, plugin, provider, command, capture format, storage design or low-level field model.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN | Review |
|---|---|---|---|---|---|---|---|---|---|
| DEP-001 | source material | product specifications | decision | apply mission, principles and boundaries | active | Product Architecture | yes | REQ-PROD-001..020 | continuous |
| DEP-002 | product boundaries / ownership | all modules | boundary | prevent concurrent ownership | active | Product Architecture | yes | REQ-PROD-013..018 | continuous |
| DEP-003 | object/permission phases | capabilities | model/security | final schemas and atomic permissions are future | partial | Architecture/Security | before implementation | REQ-OBJ/SEC | Phase 7 |
| DEP-004 | Shared mechanisms | all products | shared | Jobs, Trace, Activity, Entity, Graph, Timeline, Linking, Export, Versioning, Notifications and Recovery | partial | Shared | yes | REQ-PROD-019,020 | Technique |
| DEP-005 | OPEN-005 | Analysis Workbench and Forensics | future engine | no forensic engine, plugin or framework selected | open | Investigate | before delivery | OPEN-005 | Phase 8 |
| DEP-006 | OPEN-008 | Endpoint and sensor platform support | platform | supported acquisition, capture and filesystem platforms unresolved | open | Endpoint Agent / Settings | before delivery | OPEN-008 | Phase 4D/8 |
| DEP-007 | OPEN-013 | class-2 actions | authority | default step-up/Govern policy unresolved | open | Security | yes | OPEN-013 | Phase 7 |
| DEP-008 | OPEN-014 | Artifact relations | object | Artifact versus Attachment unresolved | open | Investigate | model blocking | OPEN-014 | Phase 7 |
| DEP-009 | OPEN-015 | Studio/Govern/Investigate runs | run model | Automation Run / Response Run bridge unresolved | open | Studio + Govern | yes | OPEN-015 | Phase 7 |
| DEP-INV-301 | CAP-INV-105/213/214 | CAP-INV-301..313 | object/provenance | Artifact, custody and provenance sources | active | Investigate | yes | REQ-PROD-014,020 | 4B.2B.1 |
| DEP-INV-314 | CAP-INV-301..313 | CAP-INV-314..328 | static-to-dynamic handoff | static results, Derived Artifacts and provenance | active | Investigate | yes | REQ-INV-002,005 | 4B.2B.2A |
| DEP-INV-329 | Static/Dynamic foundations | CAP-INV-329..346 | reverse/debug context | static/dynamic results and provenance | active | Investigate | yes | REQ-INV-003,004 | 4B.2B.2B |
| DEP-INV-347 | CAP-INV-207/213/214 | CAP-INV-347..362 | memory acquisition/custody | consume Memory Image and declared context without acquisition | active | Investigate / Endpoint Agent | yes | REQ-INV-001; OPEN-008 | 4B.2B.3A |
| DEP-INV-363 | Collection Job / Artifact Management | CAP-INV-363..379 | disk intake and analysis | consume Disk Image and acquisition context without acquiring or mutating | active | Investigate / Endpoint Agent | yes | REQ-INV-001; OPEN-008 | 4B.2B.3B.1 |
| DEP-INV-380 | CAP-INV-208 / Collection Job | CAP-INV-380 | network intake | consume the produced Capture Artifact and declared scope; no acquisition | active | Investigate / Endpoint Agent | yes | REQ-INV-001; OPEN-008 | 4B.2B.3B.2 |
| DEP-INV-381 | CAP-INV-105 | CAP-INV-380..397 | Artifact source | capture ownership, versions, restrictions and lineage | active | Investigate | yes | REQ-OBJ-003; OPEN-014 | Objects |
| DEP-INV-382 | CAP-INV-203/208/212/213/214 | CAP-INV-382/383/396 | collection context | job status, errors, integrity, custody and provenance | active | Investigate / Shared | yes | REQ-PROD-020 | 4B.2A |
| DEP-INV-383 | Platform Settings sensors/Fleet/Policies/storage/retention/health/time synchronization | CAP-INV-380..396 | projection/admin | support, restrictions, storage, health and timebase remain Settings-owned | partial | Platform Settings | yes | REQ-PROD-055; OPEN-008 | Settings/Technique |
| DEP-INV-384 | CMDR Studio Tool/Tool Call | CAP-INV-380..397 | tool/provenance | selection, version, execution and attribution remain Studio-owned | partial | Studio | yes | REQ-OBJ-009; OPEN-005/015 | Studio/Objects |
| DEP-INV-385 | CMDR Studio Automation Run | automated network analysis | automation/provenance | orchestration remains Studio-owned and optional | partial | Studio | no essential dependency | REQ-AI-002; OPEN-015 | Studio/Objects |
| DEP-INV-386 | Shared Entity / Entity Resolution / Graph | CAP-INV-387/391/392/393/397 | shared identity and relationship | no silent merge or local graph ownership | partial | Shared | yes | REQ-PROD-019,020 | Shared/Objects |
| DEP-INV-387 | Shared Timeline / Trace / Activity / Jobs / Search / Linking / Export / Reporting / Recovery | CAP-INV-380..397 | shared mechanism | progress, provenance, correlation, navigation, export and recovery | partial | Shared | yes | REQ-PROD-019,020 | Shared/Technique |
| DEP-INV-388 | CAP-INV-311 Derived Artifact Management | CAP-INV-390/395/397 | extracted Artifact | parent/child relation, restrictions, withdrawal and lineage | active | Investigate | yes | REQ-OBJ-003; OPEN-014 | 4B.2B.1 |
| DEP-INV-389 | CAP-INV-301..313 | CAP-INV-390/395/397 | Static handoff | extracted content analysis without execution | active | Investigate | no | REQ-INV-002 | 4B.2B.1 |
| DEP-INV-390 | CAP-INV-329..346 | CAP-INV-390/395/397 | Reverse handoff | binary or object analysis without ownership transfer | active | Investigate | no | REQ-INV-003,004 | 4B.2B.2B |
| DEP-INV-391 | CAP-INV-347..362 | CAP-INV-355/393/397 | Memory correlation | compare capture observations with memory-resident network state | active | Investigate | no | REQ-INV-001 | 4B.2B.3A |
| DEP-INV-392 | CAP-INV-363..379 | CAP-INV-371/374/379/393/397 | Disk correlation | compare persistent configuration/artifacts with capture observations | active | Investigate | no | REQ-INV-001 | 4B.2B.3B.1 |
| DEP-INV-393 | CAP-INV-002/003 | CAP-INV-393/397 | Event Search correlation | pivot to authorized telemetry without becoming a second SIEM search | active | Investigate / Shared | no | REQ-PROD-014 | 4B.1 |
| DEP-INV-394 | CAP-INV-321 | CAP-INV-392/393/397 | Sandbox network correlation | consume Run-scoped simulated/blocked/authorized observations | active | Investigate | no | REQ-INV-005 | 4B.2B.2A |
| DEP-INV-395 | CAP-INV-107/108 | CAP-INV-397 | Evidence handoff | qualification remains Evidence owner | active | Investigate | yes | REQ-OBJ-004 | 4B.1 |
| DEP-INV-396 | CAP-INV-109 | CAP-INV-397 | Finding handoff | draft remains unconfirmed | active | Investigate | yes | REQ-PROD-016 | 4B.1 |
| DEP-INV-397 | Future Detection Engineering and Intelligence | CAP-INV-397 | future handoff | package only; no rule or Intelligence object created | planned | future owner | no | REQ-INV-006; OPEN-011/012 | 4B.3 |
| DEP-INV-398 | Govern authority | any real-target or sensitive action | boundary | block or route; no local active network operation | active | Govern | class 3/4 | REQ-SEC-002 | 4C/5 |

OPEN-005, OPEN-008, OPEN-013, OPEN-014 and OPEN-015 remain open. OPEN-011 and OPEN-012 remain future. No low-level dependency is created.
