---
id: requirements-traceability-matrix-endpoint-ept6
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-12
source-of-truth: traceability
---
# Requirements Traceability Matrix — Endpoint EPT-6

This additive shard records EPT-6 documentary evidence only. It does not replace the canonical requirements matrix or change a requirement state without source-backed evidence.

| Requirement | EPT-6 evidence | Intended capability coverage | Baseline state | EPT-6 effect |
|---|---|---|---|---|
| REQ-PROD-004 | governed effect boundaries preserved for update/reversion/security mutation | CAP-EPT-082/085/088/094/095 | conform | additive evidence |
| REQ-PROD-005 | update/recovery/security provenance and local audit chain | CAP-EPT-087/088/090/093/097/098/099 | conform | additive evidence |
| REQ-PROD-006 | single-owner boundaries across Settings/Studio/Govern/Shared/Security | all EPT-6 | conform | additive evidence |
| REQ-PROD-008 | cross-product handoff with owner/context preservation | CAP-EPT-082/087/093/098/099 | conform | additive evidence |
| REQ-PROD-012 | all capabilities remain defined/planned, never claimed implemented | all EPT-6 | conform | preserved |
| REQ-PROD-017 | Endpoint remains a distinct product component | all EPT-6 | conform | additive evidence |
| REQ-PROD-018 | native Endpoint target semantics without delivery claim | all EPT-6 | partial/conform evidence mix | additive functional evidence only |
| REQ-PROD-019 | capability contracts remain implementation-agnostic | all EPT-6 | conform | preserved |
| REQ-PROD-020 | safe effect/recovery/verification boundaries | CAP-EPT-084..088/091..098 | conform | additive evidence |
| REQ-OBJ-007 | Result/Response Run ownership remains Govern | CAP-EPT-088/098/099 | conform | preserved |
| REQ-OBJ-008 | Fleet remains Settings-owned while local Agent state is Endpoint-owned | CAP-EPT-082/083/099 | conform | additive evidence |
| REQ-SEC-001 | permission/authority/tenant boundaries explicit | all effectful EPT-6 | conform | additive evidence |
| REQ-SEC-002 | SoD/step-up/sensitive-action distinctions retained | CAP-EPT-085/088/094/095/096 | conform | additive evidence |
| REQ-SEC-004 | Endpoint security/privacy/provenance constraints explicit | CAP-EPT-089..099 | conform | additive evidence |
| REQ-SEC-005 | secret/sensitive material remains reference-only and protected | CAP-EPT-096/097/098 | partial/conform evidence mix | additive evidence only |

## OPEN preservation
EPT-6 references but does not resolve OPEN-008, OPEN-013, OPEN-015 or OPEN-017. Other OPEN decisions remain unchanged. The global OPEN baseline stays 18 unless a separate canonical decision proves a real resolution.

## Baseline counts
The global baseline remains 122 requirements: 99 conform / 20 partial / 3 absent / 0 contradictory before EPT-6 build. Any final matrix change must be supported by exact row-level proof; this shard alone does not force count movement.