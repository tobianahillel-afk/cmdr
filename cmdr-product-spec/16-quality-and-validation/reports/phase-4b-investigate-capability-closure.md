---
id: phase-4b-investigate-capability-closure
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-06
source-of-truth: quality-report
requirements: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-INV-001, REQ-INV-002, REQ-INV-003, REQ-INV-004, REQ-INV-005, REQ-INV-006]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-012, OPEN-013, OPEN-014, OPEN-015, OPEN-017, OPEN-018, OPEN-019]
---
# Phase 4B — Investigate Capability Closure

## Verdict
**PARTIAL.** Phase 4B.1, 4B.2 and 4B.3 pass their documentary functional gates. Phase 4B cannot become PASS because OPEN-011 Mobile Forensics and OPEN-012 Cloud Analysis remain open, unimplemented and not explicitly deferred outside the approved Phase 4B scope.

## Inventory
| Phase | Capability range | Capabilities | Sections | Tables | Status |
|---|---|---:|---:|---:|---|
| 4B.1 — Signals/Hunt and Cases/Evidence | CAP-INV-001..114 families | 22 | 594 | 132 | PASS |
| 4B.2 — Collection and Analysis Workbench | CAP-INV-201..397 families | 112 | 3024 | 672 | PASS |
| 4B.3A — Detection Engineering | CAP-INV-401..435 | 35 | 945 | 210 | PASS |
| 4B.3B — Threat Intelligence | CAP-INV-501..537 | 37 | 999 | 222 | PASS |
| **Investigate total** | **all active CAP-INV families** | **206** | **5562** | **1236** | **PARTIAL at Phase 4B level** |

## Functional coverage closed
- Signals, search, Hunt, Cases, Evidence, Findings and investigation coordination.
- Authorized collection and Live Response preparation.
- Static, dynamic, reverse/debugger, memory, disk/filesystem and network analysis.
- Complete Detection Engineering functional lifecycle.
- Complete Threat Intelligence foundations, analysis, product, dissemination-preparation, operationalization-handoff, monitoring, feedback and correction lifecycle.

## Ownership and transitions
Investigate retains Case/Hypothesis/Artifact/Evidence/Finding and its analytical concepts. Command, Govern, Settings, Studio, Endpoint and Shared retain their canonical objects. Transitions preserve tenant, environment, source owner, permissions, markings, versions, errors, provenance and return origin. No concurrent owner is introduced.

## Open object, permission and implementation work
- No complete schemas, JSON Schemas, final cardinalities, physical graph or final state machines.
- No atomic permissions, final RBAC/ABAC, step-up or separation-of-duties matrix.
- No API, protocol, runtime, language, provider, engine, connector, command or product code.
- No detailed Threat Intelligence screen rewrite or new Screen ID.
- No active collection, watchlist, Indicator, rule, block, response or external sharing.

## Cloud and Mobile audit
- OPEN-011 Mobile Forensics: open; no capability or coverage added.
- OPEN-012 Cloud Analysis: open; no capability or coverage added.
- The roadmap/product scope does not contain an approved decision explicitly deferring both outside the current Phase 4B delivery.
- Therefore the valid result is **Case B**: 4B.3B.2 PASS, 4B.3B PASS, 4B.3 PASS, but **Phase 4B remains PARTIAL**.

## Totals and requirements
- Registered capabilities: 233 global — 27 Command and 206 Investigate.
- Defined/proposed/planned: 231/2/233.
- Requirements: 122 — 99 conform, 20 partial, 3 absent, 0 contradictory.
- Open decisions: 18; none closed; OPEN-009 remains the only historically resolved item.

## Parent status
- Phase 4B.3B.2: PASS after publication verification.
- Phase 4B.3B: PASS.
- Phase 4B.3: PASS.
- Phase 4B: PARTIAL due to OPEN-011/012.
- Phase 4 and global maturity: PARTIAL.
