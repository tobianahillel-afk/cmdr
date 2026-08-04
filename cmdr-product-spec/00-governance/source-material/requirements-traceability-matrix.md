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

The 122 source Requirement IDs remain unchanged. `conform` means documentary evidence exists; it does not prove implementation or delivery.

## Coverage
| State | After Phase 4B.1 | After Phase 4B.2A |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| total | 122 | 122 |

Phase 4B.2A adds evidence for Investigate collection and live response without promoting requirements dependent on objects, permissions, platforms, engines, journeys, screens or implementation.

## Phase 4B.2A evidence
| Capability | Name | Functional status | Delivery mode | Requirement IDs | Objects/concepts | Classes | OPEN | Canonical evidence |
|---|---|---|---|---|---|---|---|---|
| CAP-INV-201 | Endpoint Investigation Context | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-OBJ-008 | Case, Endpoint, Agent, Fleet/Policy projections | 0,2 | OPEN-008 | `07-investigate/modules/collection-and-live-response/capabilities/endpoint-investigation-context.md` |
| CAP-INV-202 | Collection Scope and Request Preparation | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-SEC-001 | Collection Request, Case, Endpoint, Policy | 1,2 | OPEN-008, OPEN-013 | `07-investigate/modules/collection-and-live-response/capabilities/collection-scope-and-request-preparation.md` |
| CAP-INV-203 | Collection Job Management | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-PROD-019 | Collection Request/Job, Background Job, Artifact | 0,1,2 | OPEN-008 | `07-investigate/modules/collection-and-live-response/capabilities/collection-job-management.md` |
| CAP-INV-204 | Triage Package Collection | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-OBJ-004 | Request, Job, Artifact | 1,2 | OPEN-008, OPEN-013 | `07-investigate/modules/collection-and-live-response/capabilities/triage-package-collection.md` |
| CAP-INV-205 | File and Directory Acquisition | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-OBJ-004 | Request, bounded paths, Artifact | 1,2 | OPEN-008, OPEN-013 | `07-investigate/modules/collection-and-live-response/capabilities/file-and-directory-acquisition.md` |
| CAP-INV-206 | Process and System Inspection | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-SEC-001 | Endpoint snapshot, Agent Command, Operation Result | 0,1,2 | OPEN-008, OPEN-013 | `07-investigate/modules/collection-and-live-response/capabilities/process-and-system-inspection.md` |
| CAP-INV-207 | Memory Acquisition Request | defined | planned | REQ-PROD-014, REQ-PROD-052, REQ-PROD-055 | Request, Job, Memory Image, Artifact | 1,2 | OPEN-005, OPEN-008, OPEN-013 | `07-investigate/modules/collection-and-live-response/capabilities/memory-acquisition-request.md` |
| CAP-INV-208 | Network Capture Request | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-PROD-055 | Request, Job, capture Artifact | 1,2 | OPEN-008, OPEN-013 | `07-investigate/modules/collection-and-live-response/capabilities/network-capture-request.md` |
| CAP-INV-209 | Live Session Management | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-SEC-002 | Case, Endpoint, Live Session, Decision projection | 0,2 | OPEN-007, OPEN-008, OPEN-013, OPEN-015 | `07-investigate/modules/collection-and-live-response/capabilities/live-session-management.md` |
| CAP-INV-210 | Interactive Endpoint Operations | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-SEC-002 | Live Session, Endpoint Operation, Agent Command, Operation Result | 0,1,2,3,4 | OPEN-007, OPEN-008, OPEN-013, OPEN-015 | `07-investigate/modules/collection-and-live-response/capabilities/interactive-endpoint-operations.md` |
| CAP-INV-211 | Session File Transfer | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-SEC-001 | Live Session, Artifact, transfer/result | 1,2 | OPEN-008, OPEN-013 | `07-investigate/modules/collection-and-live-response/capabilities/session-file-transfer.md` |
| CAP-INV-212 | Endpoint Operation Result Handling | defined | planned | REQ-PROD-014, REQ-PROD-018, REQ-OBJ-007 | Operation Result, Artifact, Evidence links | 0,2 | OPEN-013, OPEN-015 | `07-investigate/modules/collection-and-live-response/capabilities/endpoint-operation-result-handling.md` |
| CAP-INV-213 | Collection Integrity and Custody | defined | planned | REQ-PROD-014, REQ-PROD-020, REQ-OBJ-004 | Artifact, custody/provenance events, Evidence | 0,1,2 | OPEN-014 | `07-investigate/modules/collection-and-live-response/capabilities/collection-integrity-and-custody.md` |
| CAP-INV-214 | Collection and Live Response Provenance | defined | planned | REQ-PROD-014, REQ-PROD-020, REQ-AI-002 | collection/session/operation/Studio/Govern references | 0,1,2 | OPEN-007, OPEN-015 | `07-investigate/modules/collection-and-live-response/capabilities/collection-and-live-response-provenance.md` |
| CAP-INV-215 | Containment Request Preparation | defined | planned | REQ-PROD-014, REQ-PROD-016, REQ-SEC-002 | Finding, Evidence, Endpoint, Action Request draft | 0,2,3,4 | OPEN-007, OPEN-008, OPEN-013, OPEN-015 | `07-investigate/modules/collection-and-live-response/capabilities/containment-request-preparation.md` |

## Result of the sub-phase
- `REQ-PROD-014`, `REQ-PROD-018`, `REQ-PROD-020`, `REQ-OBJ-004`, `REQ-OBJ-008`, `REQ-SEC-001` and `REQ-SEC-002` gain new functional evidence without changing their previous global state.
- `REQ-PROD-055` remains documentary evidence only because OPEN-008 keeps platform support unresolved.
- `REQ-PROD-052` remains tied to future Phase 4B.2B engines; no engine is selected.
- `REQ-PROD-060`, `REQ-PROD-062` and `REQ-OBJ-009` remain dependent on OPEN-013/015 and future permissions/object work.
- `REQ-UX-010` remains partial because no detailed screen is rewritten.
- journey requirements remain partial/absent until Phase 5.
- new Requirement IDs: 0; removed IDs: 0; active contradictions: 0.

## Complete inventory
### conform — 99
`REQ-AI-001`, `REQ-AI-002`, `REQ-AI-003`, `REQ-AI-004`, `REQ-AI-005`, `REQ-AI-006`, `REQ-AI-011`, `REQ-BRAND-001`, `REQ-BRAND-002`, `REQ-BRAND-003`, `REQ-BRAND-004`, `REQ-BRAND-005`, `REQ-BRAND-006`, `REQ-BRAND-008`, `REQ-INV-001`, `REQ-INV-002`, `REQ-INV-003`, `REQ-INV-004`, `REQ-INV-005`, `REQ-INV-006`, `REQ-OBJ-001`, `REQ-OBJ-002`, `REQ-OBJ-003`, `REQ-OBJ-004`, `REQ-OBJ-005`, `REQ-OBJ-006`, `REQ-OBJ-007`, `REQ-OBJ-008`, `REQ-OBJ-010`, `REQ-OBJ-011`, `REQ-OBJ-012`, `REQ-PROD-001`, `REQ-PROD-002`, `REQ-PROD-003`, `REQ-PROD-004`, `REQ-PROD-005`, `REQ-PROD-009`, `REQ-PROD-011`, `REQ-PROD-012`, `REQ-PROD-013`, `REQ-PROD-014`, `REQ-PROD-015`, `REQ-PROD-016`, `REQ-PROD-017`, `REQ-PROD-018`, `REQ-PROD-019`, `REQ-PROD-021`, `REQ-PROD-022`, `REQ-PROD-023`, `REQ-PROD-024`, `REQ-PROD-025`, `REQ-PROD-026`, `REQ-PROD-027`, `REQ-PROD-028`, `REQ-PROD-029`, `REQ-PROD-030`, `REQ-PROD-031`, `REQ-PROD-032`, `REQ-PROD-033`, `REQ-PROD-034`, `REQ-PROD-035`, `REQ-PROD-036`, `REQ-PROD-037`, `REQ-PROD-038`, `REQ-PROD-039`, `REQ-PROD-040`, `REQ-PROD-041`, `REQ-PROD-042`, `REQ-PROD-043`, `REQ-PROD-044`, `REQ-PROD-045`, `REQ-PROD-046`, `REQ-PROD-047`, `REQ-PROD-048`, `REQ-PROD-049`, `REQ-PROD-050`, `REQ-PROD-051`, `REQ-PROD-052`, `REQ-PROD-053`, `REQ-PROD-054`, `REQ-PROD-055`, `REQ-PROD-056`, `REQ-PROD-057`, `REQ-PROD-058`, `REQ-PROD-059`, `REQ-PROD-060`, `REQ-PROD-061`, `REQ-PROD-062`, `REQ-SEC-001`, `REQ-SEC-002`, `REQ-UX-001`, `REQ-UX-002`, `REQ-UX-003`, `REQ-UX-004`, `REQ-UX-005`, `REQ-UX-006`, `REQ-UX-007`, `REQ-UX-008`, `REQ-UX-009`

### partial — 20
`REQ-AI-007`, `REQ-AI-008`, `REQ-AI-009`, `REQ-AI-010`, `REQ-BRAND-007`, `REQ-JRN-001`, `REQ-JRN-003`, `REQ-JRN-004`, `REQ-JRN-006`, `REQ-JRN-007`, `REQ-OBJ-009`, `REQ-PROD-006`, `REQ-PROD-007`, `REQ-PROD-008`, `REQ-PROD-010`, `REQ-PROD-020`, `REQ-SEC-003`, `REQ-SEC-004`, `REQ-SEC-005`, `REQ-UX-010`

### absent — 3
`REQ-JRN-002`, `REQ-JRN-005`, `REQ-JRN-008`

### contradictory — 0
None.

## Maintenance rule
A capability may be functionally `defined` while its delivery mode remains `planned`. Only approved implementation and release evidence can change delivery mode.
