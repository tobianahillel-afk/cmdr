---
id: requirements-traceability-matrix-govern-gov2
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: traceability-addendum
---
# Requirements Traceability Matrix — Govern GOV-2 Addendum

This additive evidence file preserves the active `requirements-traceability-matrix.md` without deleting or condensing Command/Investigate/GOV-1 evidence. The **122 source Requirement IDs and global states remain unchanged: 99 conform / 20 partial / 3 absent / 0 contradictory**. GOV-2 adds documentary capability evidence only; it does not claim implementation or promote a Requirement state.

Parent: **Delivery Roadmap Phase 4 — Govern** (`roadmap-phase-4-govern`). Execution lot: **GOV-2**. GOV-2 is not a roadmap phase. `Phase 4C Govern` does not exist.

| Evidence range | Scope | Requirements strengthened | Global state change | Reason |
|---|---|---|---|---|
| CAP-GOV-017..018 | Response Playbook catalog/selection and exact-version Decision compatibility | REQ-PROD-004,006,008,015,016,019,020; REQ-SEC-001,002 | none | Playbook != Workflow; selection/compatibility are no-effect and no runtime/provider is selected |
| CAP-GOV-019..021 | Execution Plan/parameter references, target resolution/readiness and execution-time authority reconciliation | REQ-PROD-004,006,008,015,016,017,019,020; REQ-SEC-001,002 | none | raw secrets remain Settings-owned; target drift cannot expand scope; final RBAC/ABAC/runtime remains future |
| CAP-GOV-022..025 | canonical Response Run/lifecycle/control, Response Steps and bounded Studio/Endpoint executor handoff | REQ-PROD-004,006,008,009,015,016,020; REQ-OBJ-007; REQ-SEC-001,002 | none | Run != Automation Run/Job/Tool Call; actual technical execution remains source-owned; no command/API/protocol is selected |
| CAP-GOV-026..027 | runtime reconciliation, errors, retry, partial success and compensation | REQ-PROD-004,005,008,009,015,016,020; REQ-OBJ-007; REQ-SEC-001,002 | none | technical output != Result; retries are bounded and do not reauthorize/expand scope; compensation != rollback |
| CAP-GOV-028..029 | Verification Plan, post-execution verification and residual risk | REQ-PROD-002,004,005,008,015,020; REQ-OBJ-007; REQ-SEC-001,002 | none | runtime success != verified success; source evidence remains owner-controlled; no universal score/verification engine |
| CAP-GOV-030..031 | rollback eligibility/planning and governed rollback/recovery | REQ-PROD-004,005,008,015,016,020; REQ-OBJ-007; REQ-SEC-001,002 | none | rollback support != eligibility/authority/success; technical rollback remains executor-owned; no command/runtime selected |
| CAP-GOV-032..033 | canonical Result, provenance and cross-product handoffs | REQ-PROD-002,004,005,006,008,009,015,019,020; REQ-OBJ-007; REQ-SEC-001,002 | none | Result remains Govern-owned, never rewrites Decision/Evidence/Finding; destination owners retain mutation authority; GOV-3 remains future |

## Preserved upstream evidence

The active matrix still contains the five restored Command evidence ranges:
- `CAP-CMD-001..006`;
- `CAP-CMD-101..110`;
- `CAP-CMD-201..205`;
- `CAP-CMD-301..305`;
- `CAP-CMD-401`.

All Investigate Phase 4B evidence remains active and unmodified. GOV-1 evidence `CAP-GOV-001..016` remains active and its 180/180 historical post-publication PASS is not superseded by this addendum.

## Delivery limitation

All GOV-2 capabilities are `draft` / `defined` / `planned`. This evidence proves functional documentary specification only. Final object schemas, atomic permissions, provider/runtime integrations, detailed screens, APIs/protocols, execution implementation and GOV-3 closure remain future; therefore the global Requirement distribution stays **99 / 20 / 3 / 0**.