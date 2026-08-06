---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-06
source-of-truth: canonical
---
# Requirements Traceability Matrix — through Phase 4B.3B.2

The 122 source Requirement IDs remain unchanged. `conform` records documentary evidence only, not implementation.

| State | Before 4B.3B.2 | After 4B.3B.2 |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| **Total** | **122** | **122** |

| Evidence range | Scope | Requirements strengthened | Global state change | Reason |
|---|---|---|---|---|
| CAP-INV-001..397 | Investigate foundation, collection and analysis | REQ-INV-001..005 and supporting requirements | none | final engines/models/implementation remain future |
| CAP-INV-401..435 | complete Detection Engineering lifecycle | REQ-INV-006 and supporting product/AI/security/UX requirements | none | runtime/language/implementation remain future |
| CAP-INV-501..518 | Threat Intelligence intake, requirements, sources, materials, candidate knowledge, Sightings, relationships, confidence, lifecycle and analysis handoff | REQ-PROD-014,019,020,055,060,061,062; REQ-INV-006; REQ-AI-002,010,011; REQ-SEC-001,002; REQ-UX-006,010 | none | final ontology/objects/permissions/screens/exchange and implementation remain future |
| CAP-INV-519..524 | sessions, competing hypotheses, fusion and actor/campaign/malware/infrastructure assessments | REQ-PROD-014,019,020; REQ-INV-006; REQ-AI-002; REQ-SEC-001,002; REQ-UX-010 | none | final objects, algorithms and implementation remain future |
| CAP-INV-525..529 | product planning, authoring, review, releasability and internal publication | REQ-PROD-014,019,020,055; REQ-INV-006; REQ-AI-002; REQ-SEC-001,002; REQ-UX-010 | none | OPEN-019, final permissions/screens/delivery contracts remain future |
| CAP-INV-530..533 | watchlist definition, operationalization, monitoring and external-sharing preparation | REQ-PROD-014,019,020,055; REQ-INV-006; REQ-SEC-001,002 | none | no activation, deployment, runtime mutation or actual sharing |
| CAP-INV-534..537 | feedback, effectiveness, satisfaction, correction and lifecycle closure | REQ-PROD-014,019,020; REQ-INV-006; REQ-AI-002 | none | final metrics/contracts and implementation remain future |

## Disposition
- Threat Intelligence now has complete functional documentary evidence from intake through lifecycle closure without claiming implementation.
- REQ-INV-006 and the broader Investigate capability remain globally partial because implementation is absent and Cloud/Mobile scope remains unresolved.
- OPEN-018 records the unresolved ontology/interoperability/exchange strategy; OPEN-019 records dissemination/releasability/sharing/access policy.
- OPEN-017 remains Detection-only and unchanged.
- REQ-PROD-055 remains partial and dependent on OPEN-008; source/provider support is not invented.
- REQ-PROD-020 and REQ-OBJ-009 remain partial because final trace/run contracts are not defined.
- REQ-PROD-060/061/062 remain tied to OPEN-013/014/015/019 and final permissions/object relations.
- REQ-UX-010 remains partial because no detailed Threat Intelligence screen is created or rewritten.
- New Requirement IDs: 0; removed IDs: 0; active contradictions: 0.
