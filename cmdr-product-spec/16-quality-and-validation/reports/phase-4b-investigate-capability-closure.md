---
id: phase-4b-investigate-capability-closure
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-07
source-of-truth: quality-report
requirements: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-014, REQ-PROD-019, REQ-PROD-020, REQ-INV-001, REQ-INV-002, REQ-INV-003, REQ-INV-004, REQ-INV-005, REQ-INV-006]
open_decisions: [OPEN-005, OPEN-008, OPEN-011, OPEN-012, OPEN-013, OPEN-014, OPEN-015, OPEN-017, OPEN-018, OPEN-019]
---
# Phase 4B — Investigate Capability Closure

## Verdict before Mobile post-publication verification
**PARTIAL.** Phase 4B.1, 4B.2 and 4B.3 are PASS. Phase 4B.4A Cloud Analysis is verified PASS. Phase 4B.4B Mobile Forensics has complete prepared functional content but remains `PENDING POST-PUBLICATION VERIFICATION` until the fifth Mobile functional commit and remote checks pass. No other mandatory functional category is missing from the audited Phase 4B roadmap.

## Inventory
| Phase | Capability range | Capabilities | Sections | Tables | Status |
|---|---|---:|---:|---:|---|
| 4B.1 — Signals/Hunt and Cases/Evidence | CAP-INV-001..114 families | 22 | 594 | 132 | PASS |
| 4B.2 — Collection and Analysis Workbench | CAP-INV-201..397 families | 112 | 3024 | 672 | PASS |
| 4B.3A — Detection Engineering | CAP-INV-401..435 | 35 | 945 | 210 | PASS |
| 4B.3B — Threat Intelligence | CAP-INV-501..537 | 37 | 999 | 222 | PASS |
| 4B.4A — Cloud Analysis | CAP-INV-601..618 | 18 | 486 | 108 | PASS AFTER POST-PUBLICATION VERIFICATION |
| 4B.4B — Mobile Forensics | CAP-INV-701..719 | 19 | 513 | 114 | PENDING POST-PUBLICATION VERIFICATION |
| **Investigate total** | **all active CAP-INV families** | **243** | **6561** | **1458** | **PARTIAL pending Mobile remote verification** |

## Functional coverage audited
- Signals, search, Hunt, Cases, Evidence, Findings and investigation coordination.
- Authorized Collection and Live Response preparation/execution ownership boundaries.
- Static, dynamic, reverse/debugger, memory, disk/filesystem and network analysis.
- Complete Detection Engineering functional lifecycle.
- Complete Threat Intelligence foundations, analysis, product, dissemination-preparation, operationalization-handoff, monitoring, feedback and correction lifecycle.
- Provider-neutral Cloud Analysis from intake through provenance.
- Provider-neutral Mobile Forensics from intake/acquisition-context review through storage/apps/private data/connectivity/backups/recovery/timeline/hypotheses/Derived Artifacts/handoffs/provenance.

## Ownership and transitions
Investigate retains Case/Hypothesis/Artifact/Evidence/Finding and its analytical concepts. Collection owns acquisition requests/jobs/execution/results/custody. Command, Govern, Settings, Studio, Endpoint and Shared retain their canonical objects. Mobile does not assume Endpoint Agent presence and executes no real-device action. Transitions preserve tenant, environment, device/source scope, source owner, permissions, classifications, versions, errors, provenance and return origin. Concurrent owners introduced: **0**.

## Open decisions after functional coverage
OPEN-011 remains open for Mobile platform/version/tool/acquisition delivery strategy; OPEN-012 remains open for Cloud provider/service/delivery strategy. These decisions block implementation choices but do not inherently block provider-neutral documentary functional completeness once the corresponding subphase passes its verification gates. OPEN-005/008/013/014/015/017/018/019 likewise remain assigned to their future owner phases.

## Remaining global object, permission, screen and implementation work
- no complete schemas, JSON Schemas, final cardinalities, physical graph or final state machines;
- no atomic permissions, final RBAC/ABAC, step-up or separation-of-duties matrix;
- no API, protocol, runtime, language, provider/platform, engine, connector, command or product code;
- no detailed Mobile/Cloud/TI screen rewrite or new Screen ID;
- no actual unauthorized collection, unlock/bypass/root/jailbreak, secret use, active scan, watchlist, Indicator, rule, block, response, external sharing or device mutation.

These later global gaps keep Phase 4/global maturity PARTIAL but do not represent missing Phase 4B capability families.

## Totals and requirements
- Registered capabilities: **270 global — 27 Command and 243 Investigate**.
- Defined/proposed/planned: **268 / 2 / 270**.
- Investigate: **6561 sections / 1458 tables**.
- Command + Investigate: **7290 sections / 1620 tables**.
- Requirements: **122 — 99 conform, 20 partial, 3 absent, 0 contradictory**.
- Open decisions: **18**; none closed by Mobile; OPEN-009 remains the only historically resolved item.

## Phase 4B closure rule
If final Mobile remote verification passes all 250 gates and Cloud remains verified PASS, then every audited Phase 4B child is functionally documented and Phase 4B may be promoted to **PASS** without closing OPEN-011/012 delivery-strategy decisions. If any Mobile remote gate fails, Phase 4B remains PARTIAL. No status is forced before that verification.

## Parent status before final Mobile verification
- Phase 4B.3: PASS.
- Phase 4B.4A: PASS after publication verification.
- Phase 4B.4B: PENDING POST-PUBLICATION VERIFICATION.
- Phase 4B.4: PARTIAL.
- Phase 4B: PARTIAL.
- Phase 4 and global maturity: PARTIAL.
