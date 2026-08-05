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

This register records functional and documentary dependencies only. It chooses no protocol, API, engine, plugin, provider, command, disk-image format, carving algorithm or storage design.

| ID | Source | Dependent | Type | Reason | Status | Owner | Blocking | Requirement / OPEN | Review |
|---|---|---|---|---|---|---|---|---|---|
| DEP-001 | source material | product specifications | decision | apply mission, principles and boundaries | active | Product Architecture | yes | REQ-PROD-001..020 | continuous |
| DEP-002 | product boundaries / ownership | all modules | boundary | prevent concurrent ownership | active | Product Architecture | yes | REQ-PROD-013..018 | continuous |
| DEP-003 | object/permission phases | capabilities | model/security | final schemas and atomic permissions are future | partial | Architecture/Security | before implementation | REQ-OBJ/SEC | Phase 7 |
| DEP-004 | Shared mechanisms | all products | shared | Jobs, Trace, Activity, Linking, Export, Versioning, Notifications and Recovery | partial | Shared | yes | REQ-PROD-019,020 | Technique |
| DEP-005 | OPEN-005 | Analysis Workbench and Forensics | future engine | no forensic engine, plugin or framework selected | open | Investigate | before delivery | OPEN-005 | Phase 8 |
| DEP-006 | OPEN-008 | Endpoint platform support | platform | supported acquisition and filesystem platforms unresolved | open | Endpoint Agent / Settings | before delivery | OPEN-008 | Phase 4D/8 |
| DEP-007 | OPEN-013 | class-2 actions | authority | default step-up/Govern policy unresolved | open | Security | yes | OPEN-013 | Phase 7 |
| DEP-008 | OPEN-014 | Artifact relations | object | Artifact versus Attachment unresolved | open | Investigate | model blocking | OPEN-014 | Phase 7 |
| DEP-009 | OPEN-015 | Studio/Govern/Investigate runs | run model | Automation Run / Response Run bridge unresolved | open | Studio + Govern | yes | OPEN-015 | Phase 7 |
| DEP-INV-301 | CAP-INV-105/213/214 | CAP-INV-301..313 | object/provenance | Artifact, custody and provenance sources | active | Investigate | yes | REQ-PROD-014,020 | 4B.2B.1 |
| DEP-INV-314 | CAP-INV-301..313 | CAP-INV-314..328 | static-to-dynamic handoff | static results, Derived Artifacts and provenance | active | Investigate | yes | REQ-INV-002,005 | 4B.2B.2A |
| DEP-INV-329 | Static/Dynamic foundations | CAP-INV-329..346 | reverse/debug context | static/dynamic results and provenance | active | Investigate | yes | REQ-INV-003,004 | 4B.2B.2B |
| DEP-INV-347 | CAP-INV-207/213/214 | CAP-INV-347..362 | memory acquisition/custody | consume Memory Image and declared context without acquisition | active | Investigate / Endpoint Agent | yes | REQ-INV-001; OPEN-008 | 4B.2B.3A |
| DEP-INV-363 | Collection Job / Artifact Management | CAP-INV-363 | disk intake | consume Disk Image and declared acquisition context; no acquisition | active | Investigate / Endpoint Agent | yes | REQ-INV-001; OPEN-008 | 4B.2B.3B.1 |
| DEP-INV-364 | CAP-INV-105 | CAP-INV-363..379 | object source | Disk Image/Artifact ownership, versions, restrictions and lineage | active | Investigate | yes | REQ-OBJ-003; OPEN-014 | Objects |
| DEP-INV-365 | CAP-INV-203/205/212/213/214 | CAP-INV-365/378 | custody/provenance | collection result, errors, integrity context, transfers and trace | active | Investigate / Shared | yes | REQ-PROD-020 | 4B.2A |
| DEP-INV-366 | CMDR Studio Tool/Tool Call | CAP-INV-363..379 | tool/provenance | selection, version, execution and attribution | partial | Studio | yes | REQ-OBJ-009; OPEN-015 | Studio/Objects |
| DEP-INV-367 | CMDR Studio Automation Run | automated disk analysis | automation/provenance | orchestration remains Studio-owned | partial | Studio | no essential dependency | REQ-AI-002; OPEN-015 | Studio/Objects |
| DEP-INV-368 | Platform Settings Fleet/Policy/storage/retention/health | CAP-INV-363..378 | projection/admin | support, restrictions, storage and administrative state | partial | Platform Settings | yes | REQ-PROD-055; OPEN-008 | Settings/Technique |
| DEP-INV-369 | Shared Jobs/Trace/Activity/Timeline/Linking/Export/Recovery/Search | CAP-INV-363..379 | shared mechanism | progress, navigation, provenance, correlation and export | partial | Shared | yes | REQ-PROD-019,020 | Shared/Technique |
| DEP-INV-370 | CAP-INV-311 Derived Artifact Management | CAP-INV-368/369/375/376/379 | artifact handoff | extracted/recovered content and lineage | active | Investigate | yes | REQ-OBJ-003; OPEN-014 | 4B.2B.1 |
| DEP-INV-371 | CAP-INV-301..313 / CAP-INV-329..346 | CAP-INV-368/375/379 | static/reverse handoff | file or Derived Artifact analysis | active | Investigate | no | REQ-INV-002,003 | 4B.2B.1/2B |
| DEP-INV-372 | CAP-INV-347..362 | CAP-INV-371/373/374/377/379 | memory correlation | compare persistent and memory-resident observations | active | Investigate | no | REQ-INV-001 | 4B.2B.3A |
| DEP-INV-373 | CAP-INV-107/108 | CAP-INV-379 | Evidence handoff | qualification remains Evidence owner | active | Investigate | yes | REQ-OBJ-004 | 4B.1 |
| DEP-INV-374 | CAP-INV-109 | CAP-INV-379 | Finding handoff | draft remains unconfirmed | active | Investigate | yes | REQ-PROD-016 | 4B.1 |
| DEP-INV-375 | Future Detection Engineering | CAP-INV-373/379 | future handoff | package only; no rule creation or deployment | planned | future owner | no | REQ-INV-006 | 4B.3 |
| DEP-INV-376 | Persistent network configuration/artifacts | future Phase 4B.2B.3B.2 | future handoff | no capture, flow, protocol or full Network Forensics analysis | planned | Investigate | no | REQ-INV-001 | 4B.2B.3B.2 |
| DEP-INV-377 | Govern authority | any real-target or sensitive action | boundary | block or route; no local Endpoint action | active | Govern | class 3/4 | REQ-SEC-002 | 4C/5 |

OPEN-005, OPEN-008, OPEN-013, OPEN-014 and OPEN-015 remain open. OPEN-011/012 remain future. No low-level dependency is created.
