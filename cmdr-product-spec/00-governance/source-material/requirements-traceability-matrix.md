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

| State | Before 4B.2B.2B | After 4B.2B.2B |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| total | 122 | 122 |

| Evidence range | Scope | Requirements strengthened | Global state change | Reason |
|---|---|---|---|---|
| CAP-INV-301..313 | Static Analysis Workbench | REQ-INV-002; REQ-PROD-014,020; REQ-OBJ-003,004,009; REQ-AI-002; REQ-SEC-001,002; REQ-UX-002,006 | none | final objects, permissions, engines, journeys, screens and implementation remain future |
| CAP-INV-314..328 | Dynamic Sandbox | REQ-INV-005; REQ-PROD-014,017,020; REQ-OBJ-003,004,009; REQ-AI-002; REQ-SEC-001,002; REQ-UX-007 | none | environment administration, final run/observation objects and implementation remain future |
| CAP-INV-329..346 | Reverse Engineering and Debugger | REQ-INV-003,004; REQ-PROD-014,017,020,052; REQ-OBJ-003,004,009; REQ-AI-002; REQ-SEC-001,002; REQ-UX-002,006 | none | engines, final concepts, permissions, journeys, detailed screens and implementation remain future |

## Disposition
- REQ-INV-003 and REQ-INV-004 gain complete functional capability evidence while delivery remains planned.
- REQ-PROD-020 remains partial because provenance concepts and technical contracts are not final.
- REQ-OBJ-009 remains partial because Tool Call, Automation Run and analytical session concepts are not fully modeled.
- REQ-PROD-052 remains dependent on OPEN-005; no reverse, decompiler or debugger engine is selected.
- REQ-PROD-060/061/062 remain tied to OPEN-013/014/015.
- REQ-UX-010 remains partial because no detailed screen is rewritten.
- New Requirement IDs: 0; removed IDs: 0; active contradictions: 0.

The complete conform/partial/absent inventory remains unchanged.
