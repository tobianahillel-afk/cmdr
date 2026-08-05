---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-05
source-of-truth: source-material
requirements:
  - REQ-PROD-001
  - REQ-PROD-062
  - REQ-UX-001
  - REQ-UX-010
  - REQ-OBJ-001
  - REQ-OBJ-012
---
# Requirements Traceability Matrix

The 122 source Requirement IDs remain unchanged. `conform` records documentary evidence only, not implementation.

| State | Before 4B.2B.3B.2 | After 4B.2 closure |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| total | 122 | 122 |

| Evidence range | Scope | Requirements strengthened | Global state change | Reason |
|---|---|---|---|---|
| CAP-INV-201..215 | Collection and Live Response | REQ-INV-001 and acquisition, custody, provenance requirements | none | implementation and final objects remain future |
| CAP-INV-301..313 | Static Analysis Workbench | REQ-INV-002 and supporting requirements | none | engine, implementation and final models remain future |
| CAP-INV-314..328 | Dynamic Sandbox | REQ-INV-005 and environment/run/provenance requirements | none | implementation remains future |
| CAP-INV-329..346 | Reverse Engineering and Debugger | REQ-INV-003,004 and supporting requirements | none | engines and final models remain future |
| CAP-INV-347..362 | Memory Forensics | REQ-INV-001 and memory-specific trust, analysis and handoff | none | engine, platforms and implementation remain future |
| CAP-INV-363..379 | Disk and Filesystem Forensics | REQ-INV-001 and persistent-source analysis | none | engine, formats and final objects remain future |
| CAP-INV-380..397 | Network Forensics | REQ-INV-001; REQ-PROD-014,019,020,052,055; REQ-OBJ-003,004,009; REQ-AI-002; REQ-SEC-001,002; REQ-UX-002,006,007,010 | none | engines, sensor support, final objects, permissions, screens, technical contracts and implementation remain future |

## Disposition
- REQ-INV-001 now has complete Phase 4B.2 functional evidence across Collection, Memory, Disk and Network Forensics.
- REQ-PROD-052 remains partial and dependent on OPEN-005; no forensic engine or product is selected.
- REQ-PROD-055 remains partial and dependent on OPEN-008; no final platform, sensor, interface or protocol support list is invented.
- REQ-PROD-020 and REQ-OBJ-009 remain partial because session, observation, provenance and run contracts are not final.
- REQ-PROD-060/061/062 remain tied to OPEN-013/014/015.
- REQ-UX-010 remains partial because no detailed Network Forensics screen is created or rewritten.
- REQ-INV-006 remains future: Detection Engineering and Intelligence packages are handoff-only and Phase 4B.3 is not started.
- New Requirement IDs: 0; removed IDs: 0; active contradictions: 0.
