---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: canonical
---
# Requirements Traceability Matrix — through Phase 4B.3A

The 122 source Requirement IDs remain unchanged. `conform` records documentary evidence only, not implementation.

| State | Before 4B.3A.2 | After 4B.3A.2 |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| **Total** | **122** | **122** |

| Evidence range | Scope | Requirements strengthened | Global state change | Reason |
|---|---|---|---|---|
| CAP-INV-001..397 | Investigate foundation, collection and analysis | REQ-INV-001..005 and supporting requirements | none | final engines/models/implementation remain future |
| CAP-INV-401..417 | Detection authoring, validation, historical replay, candidate review and coverage | REQ-INV-006 and supporting product/AI/security/UX requirements | none | runtime lifecycle and Intelligence were future at that point |
| CAP-INV-418..435 | formal review, readiness, governed promotion/deployment coordination, runtime health/quality, tuning, bounded proposals, drift, performance, rollback, retirement and continuous improvement | REQ-INV-006; REQ-PROD-013,014,015,019,020,055,060,061,062; REQ-OBJ-005,006,007,010; REQ-AI-002,010,011; REQ-SEC-001,002; REQ-UX-006,010 | none | Threat Intelligence, final objects/permissions/screens/runtime strategy and implementation remain future |

## Disposition
- REQ-INV-006 now has complete functional evidence for Detection Engineering from intake through lifecycle feedback.
- REQ-INV-006 remains globally partial because Threat Intelligence and implementation remain future.
- OPEN-017 records the unresolved runtime, target-language and portability strategy; no option is selected.
- REQ-PROD-055 remains partial and dependent on OPEN-008; target/runtime support is not invented.
- REQ-PROD-020 and REQ-OBJ-009 remain partial because final trace/run contracts are not defined.
- REQ-PROD-060/061/062 remain tied to OPEN-013/014/015 and final permissions/object relations.
- REQ-UX-010 remains partial because no detailed Detection Engineering screen is created or rewritten.
- New Requirement IDs: 0; removed IDs: 0; active contradictions: 0.
