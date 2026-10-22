---
id: requirements-traceability-matrix-endpoint-ept4
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-11
source-of-truth: supporting-traceability
---
# Requirements Traceability — Endpoint EPT-4

This is an additive EPT-4 supplement. It does **not** create Requirement IDs or change the canonical distribution **122 = 99 conform / 20 partial / 3 absent / 0 contradictory**.

| Requirement family | EPT-4 evidence | Notes |
|---|---|---|
| REQ-PROD-006 | CAP-EPT-047/048/056/058/064 | bounded cross-product functional contracts |
| REQ-PROD-014 | CAP-EPT-047..064 | Collection and Live Response technical execution coverage |
| REQ-PROD-018 | CAP-EPT-047..064 | Endpoint Agent local technical capability semantics |
| REQ-PROD-019 | 047/048/050/053/056..064 | state, orchestration boundaries and target context |
| REQ-PROD-020 | 047/049/051/052/054/055/056/058/060..064 | provenance, output and custody-compatible metadata |
| REQ-PROD-052 | CAP-EPT-051 | memory acquisition boundary only; analysis external |
| REQ-PROD-055 | CAP-EPT-051/052 | bounded memory/network acquisition where supported |
| REQ-OBJ-004 | CAP-EPT-054/064 | integrity/custody metadata without final object identity |
| REQ-OBJ-007 | CAP-EPT-064 | technical output != Govern Result |
| REQ-SEC-001 | all effect/sensitive families | least privilege, tenant isolation, authority refs |
| REQ-SEC-002 | session/execution/control families | permission/authority and sensitive runtime state |
| REQ-SEC-004 | file/memory/network/package/output/provenance | classification, masking and provenance |

OPEN-008/014/015/017 remain open. No canonical Requirement status is promoted solely because EPT-4 is documented.