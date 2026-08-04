---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-04
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

| State | Before | After 4B.2B.2A |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| total | 122 | 122 |

| Evidence range | Scope | Requirements strengthened | Global state change | Reason |
|---|---|---|---|---|
| CAP-INV-301..313 | Static Analysis Workbench | REQ-INV-002; REQ-PROD-014,020; REQ-OBJ-003,004,009; REQ-AI-002; REQ-SEC-001,002; REQ-UX-002,006 | none | final objects, permissions, engines, journeys, screens and implementation remain future |
| CAP-INV-314..328 | Dynamic Sandbox and Behavioral Analysis | REQ-INV-005; REQ-PROD-014,017,020; REQ-OBJ-003,004,009; REQ-AI-002; REQ-SEC-001,002; REQ-UX-007 | none | environment administration, final run/observation objects, technical execution and implementation remain future |

## Disposition
- REQ-INV-002 and REQ-INV-005 gain specific capability evidence while delivery remains planned.
- REQ-PROD-020 remains partial because provenance objects and technical contracts are not final.
- REQ-OBJ-009 remains partial because Tool Call, Automation Run and analysis concepts are not fully modeled.
- REQ-PROD-052 remains dependent on OPEN-005; no engine or hypervisor is selected.
- REQ-PROD-060/061/062 remain tied to OPEN-013/014/015.
- REQ-UX-010 remains partial because no detailed screen is rewritten.
- New Requirement IDs: 0; removed IDs: 0; active contradictions: 0.

The complete conform/partial/absent inventory established in Phase 4B.2A remains unchanged.
