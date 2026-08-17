---
id: phase-4b2-investigate-capability-closure
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-05
source-of-truth: quality-report
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-INV-001
  - REQ-INV-002
  - REQ-INV-003
  - REQ-INV-004
  - REQ-INV-005
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-011
  - OPEN-012
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Phase 4B.2 — Investigate capability closure

## Verdict
**PASS after post-publication verification.** Phase 4B.2A, 4B.2B.1, 4B.2B.2A, 4B.2B.2B, 4B.2B.3A, 4B.2B.3B.1 and 4B.2B.3B.2 provide a continuous functional path from authorized collection to static, dynamic, reverse/debugger, memory, disk/filesystem and network analysis. Phase 4B.3 remains not started.

## Sub-phase inventory
| Sub-phase | Capability range | Capabilities | Sections | Mandatory tables | Report | Publication SHA | Status |
|---|---|---:|---:|---:|---|---|---|
| 4B.2A — Collection and Live Response | CAP-INV-201..215 | 15 | 405 | 90 | `phase-4b2a-collection-live-response-capability-conformance.md` | `2bc2d37e2d726d15e1ee6af1ce3ae633703b4cc3` corrective closure | PASS |
| 4B.2B.1 — Static Analysis | CAP-INV-301..313 | 13 | 351 | 78 | `phase-4b2b1-analysis-workbench-static-capability-conformance.md` | `df06eeaa94a7a3bc982efc9f5e5c6b4571b279c0` restoration closure | PASS |
| 4B.2B.2A — Dynamic Sandbox | CAP-INV-314..328 | 15 | 405 | 90 | `phase-4b2b2a-dynamic-sandbox-capability-conformance.md` | `eb60be74cdd4d3b08a37d5ac0a22d781b51aef86` | PASS |
| 4B.2B.2B — Reverse and Debugger | CAP-INV-329..346 | 18 | 486 | 108 | `phase-4b2b2b-reverse-debugger-capability-conformance.md` | `ce26d6237202572c88f7428b4a532cd809a9c74c` | PASS |
| 4B.2B.3A — Memory Forensics | CAP-INV-347..362 | 16 | 432 | 96 | `phase-4b2b3a-memory-forensics-capability-conformance.md` | `b6a16d66d66117273f0490afb362038f63e7934e` | PASS |
| 4B.2B.3B.1 — Disk and Filesystem | CAP-INV-363..379 | 17 | 459 | 102 | `phase-4b2b3b1-disk-filesystem-forensics-capability-conformance.md` | `516226778850a5d0ef4e6bbec6afc3175af90aba` | PASS |
| 4B.2B.3B.2 — Network Forensics and Closure | CAP-INV-380..397 | 18 | 486 | 108 | `phase-4b2b3b2-network-forensics-capability-conformance.md` | containing fifth commit; exact immutable SHA recorded in PR #2 and final execution report | PASS |
| **Phase 4B.2 total** | **CAP-INV-201..397 programme families** | **112** | **3024** | **672** | **8 reports including this closure report** | **final branch SHA recorded after fast-forward** | **PASS** |

## Publication chain for Network Forensics
| Commit | Message | Scope |
|---|---|---|
| `b9b8f7942f40e520cf0243bf117c1353883b443d` | `docs: establish Investigate network forensics boundaries` | module foundation and CAP-INV-380..383 |
| `cdf865293fcf19db5cedf76555b6c7560e865e82` | `docs: define capture inspection flow and conversation capabilities` | CAP-INV-384..389 |
| `4739bcf61a4529b3ff024afacad5462137ed620c` | `docs: specify network transactions entities and artifact extraction` | CAP-INV-390..392 and CAP-INV-395 plus permissions/payload handling |
| `90e01601ee9ae3cecbcae6a70dddcba0e2d50f0c` | `docs: document network correlation provenance and future handoffs` | CAP-INV-393,394,396,397 and module maps |
| containing commit | `docs: close Investigate analysis workbench and phase four two quality gates` | registers, traceability, reports, STATUS and CHANGELOG |

The exact fifth commit SHA cannot be embedded in the commit that creates this report without an additional commit. PR #2 and the final execution report are the immutable post-publication evidence.

## Totals
| Scope | Capabilities | Defined | Proposed | Planned | Sections | Mandatory tables |
|---|---:|---:|---:|---:|---:|---:|
| Phase 4B.1 | 22 | 21 | 1 | 22 | 594 | 132 |
| Phase 4B.2 | 112 | 112 | 0 | 112 | 3024 | 672 |
| Investigate | 134 | 133 | 1 | 134 | 3618 | 804 |
| Command | 27 | 26 | 1 | 27 | 729 | 162 |
| Command + Investigate | 161 | 159 | 2 | 161 | 4347 | 966 |

The totals are calculated from the active Capability Register shards and the conformance reports. They match the theoretical programme totals.

## End-to-end functional coverage
1. Case, Hypothesis and scope establish the investigation context.
2. Collection and Live Response prepares and tracks authorized acquisition without owning downstream analysis.
3. Artifact Management, integrity, custody and provenance preserve source identity and limitations.
4. Static Analysis inspects content without execution and creates bounded Derived Artifacts.
5. Dynamic Sandbox observes isolated behavior with explicit environment and network safety boundaries.
6. Reverse Engineering and Debugger provide static and isolated runtime interpretation without real Endpoint execution.
7. Memory Forensics reconstructs memory-resident state with platform/profile and partiality controls.
8. Disk and Filesystem Forensics reconstructs persistent state, deleted/residual content and Disk Timeline.
9. Network Forensics analyzes captures, packets, flows, sessions, conversations, protocols, DNS, transactions, encrypted metadata, transfers, Entities, behaviors and Network Timeline.
10. All families preserve provenance, no-AI operation and explicit Evidence/Finding handoffs.

## Ownership audit
- Investigate owns Case, Hypothesis, Artifact, analytical contexts and sessions, observations, interpretations, Derived Artifacts, Evidence candidates and Finding Drafts.
- Command retains Detection, Signal, Alert, Incident and operational coordination.
- Endpoint Agent and Collection execute or report authorized acquisition and do not own analysis.
- Platform Settings owns Fleet, sensors, policies, providers, storage, retention, health, time synchronization and secrets.
- CMDR Studio owns Tool, Tool Call, Workflow, Automation Run, versions and evaluation context.
- Govern owns Decision, Approval, Response Run, Result and all real-target authority.
- Shared owns Entity, Graph, Timeline, Jobs, Notifications, Trace, Activity, Linking, Search, Export, Reporting, Versioning, Collaboration, Audit Hooks and Recovery.

Concurrent owners introduced: **0**.

## Object and concept audit
Canonical objects are consumed without silent mutation. New Network, Memory, Disk and analysis session concepts remain functional gaps for the Objects phase; no complete schema, JSON Schema, final cardinality or object state machine is created. Artifact, Derived Artifact, Runtime Artifact, Evidence, Finding, Capture Artifact, Memory Image and Disk Image remain distinct. Tool Call, Automation Run and Response Run remain distinct.

## Transition audit
Collection→Artifact→Analysis transitions preserve source, tenant, environment, Case, permissions, restrictions, errors, partiality and return origin. Cross-source correlation never transfers ownership. Handoffs to Evidence and Finding retain qualification in CAP-INV-107/108/109. Future Detection Engineering and Intelligence receive packages only; Phase 4B.3 remains not started.

## Requirements and OPEN decisions
- Requirement IDs: **122** — 99 conform, 20 partial, 3 absent, 0 contradictory.
- Open decisions: **15**; no decision closed by this phase.
- OPEN-005, OPEN-008, OPEN-013, OPEN-014 and OPEN-015 remain directly relevant.
- OPEN-011 Mobile and OPEN-012 Cloud remain open and outside scope.

## Migrations and screens
- Static legacy functional documents migrated previously: preserved.
- Reverse/Debugger generic functional documents migrated previously: preserved.
- Memory generic functional documents migrated previously: preserved.
- Disk generic functional documents migrated previously: preserved.
- Full Network Forensics competing functional source: none; migrations in 4B.2B.3B.2: **0**.
- Detailed screen rewrites across closure: **0** for the Network phase; existing active screens remain active.
- New Network Screen IDs: **0**.

## Quality and boundaries
- empty capability files: 0;
- placeholders: 0;
- duplicate IDs: 0;
- concurrent owners: 0;
- generic mandatory tables: 0;
- active contradictions: 0;
- APIs, protocols, engines, imposed products, real commands and code: 0;
- packet generation, injection, active scanning/interception, replay and target interaction: 0;
- unauthorized decryption, secret extraction/use, exploit and evasion guidance: 0;
- Detection rules or canonical Intelligence objects: 0;
- Cloud and Mobile analysis content: 0.

## Status closure
- Phase 4B.2A: PASS.
- Phase 4B.2B.1: PASS.
- Phase 4B.2B.2A: PASS.
- Phase 4B.2B.2B: PASS.
- Phase 4B.2B.2 global: PASS.
- Phase 4B.2B.3A: PASS.
- Phase 4B.2B.3B.1: PASS.
- Phase 4B.2B.3B.2: PASS after remote verification.
- Phase 4B.2B.3B global: PASS.
- Phase 4B.2B.3 global: PASS.
- Phase 4B.2B global: PASS.
- Phase 4B.2 global: PASS.
- Phase 4B global: PARTIAL.
- Phase 4 global: PARTIAL.
- Repository global maturity: PARTIAL.
- Phase 4B.3 Detection Engineering and Intelligence: NOT STARTED.
