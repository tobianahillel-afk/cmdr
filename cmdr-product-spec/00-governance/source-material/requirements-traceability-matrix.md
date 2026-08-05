---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: canonical
---
# Requirements Traceability Matrix — through Phase 4B.3A.1

The 122 source Requirement IDs remain unchanged. `conform` records documentary evidence only, not implementation.

| State | Before 4B.3A.1 | After 4B.3A.1 |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| **Total** | **122** | **122** |

| Evidence range | Scope | Requirements strengthened | Global state change | Reason |
|---|---|---|---|---|
| CAP-INV-001..114 | Signals, Hunt, Cases, Evidence and Findings | REQ-PROD-014 and related object/AI requirements | none | implementation/final objects remain future |
| CAP-INV-201..397 | Collection and Analysis Workbench | REQ-INV-001..005 and supporting requirements | none | engines, final models and implementation remain future |
| CAP-INV-401..417 | Detection Engineering foundations, authoring and validation | REQ-INV-006; REQ-PROD-014,019,020,055; REQ-AI-002,010,011; REQ-SEC-001,002; REQ-UX-006,010 | none | 4B.3A.2, runtime bridge, final objects, permissions, screens, engine/language and implementation remain future |

## Disposition
- REQ-INV-006 gains complete functional evidence for foundations, authoring, validation, tests, historical replay, candidate FP/FN review, coverage/gaps and review handoff.
- REQ-INV-006 remains globally partial because promotion, deployment, runtime performance, Threat Intelligence and implementation are future.
- REQ-PROD-052 and OPEN-005 remain forensic-only; no Detection engine or language is selected.
- REQ-PROD-055 remains partial and dependent on OPEN-008; no source/platform support is invented.
- REQ-PROD-020 and REQ-OBJ-009 remain partial because provenance, run and object contracts are not final.
- REQ-PROD-060/061/062 remain tied to OPEN-013/014/015.
- REQ-UX-010 remains partial because no detailed Detection Engineering screen is created or rewritten.
- New Requirement IDs: 0; removed IDs: 0; active contradictions: 0.
