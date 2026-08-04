---
id: report-phase-4b2b1-analysis-workbench-static
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-UX-010
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Phase 4B.2B.1 — Analysis Workbench and Static Analysis conformance

## Verdict
**PASS — 120/120 gates after corrective remote restoration.**

The authoritative starting head was `2bc2d37e2d726d15e1ee6af1ce3ae633703b4cc3`, which contained no CAP-INV-3xx. This report records the actual reconstructed Phase 4B.2B.1.

## Scope and measures
- Capabilities: **13/13**, CAP-INV-301 through CAP-INV-313.
- Numbered sections: **351/351**.
- Mandatory tables: **78/78**.
- Empty, prose-only or generic mandatory tables: **0**.
- Duplicate IDs, concurrent owners and active contradictions: **0**.
- Detailed screens rewritten: **0**.
- APIs, protocols, engines, commands and product code: **0**.

## Capability conformance
| Capability ID | Sections 1–27 | S8 | S9 | S10 | S13 | S16 | S17 | Front matter | Verdict |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| CAP-INV-301 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-302 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-303 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-304 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-305 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-306 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-307 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-308 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-309 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-310 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-311 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-312 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-313 | 27/27 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |

## Gate groups
| Gates | Group | PASS | FAIL |
|---|---|---:|---:|
| 1–12 | Git | 12 | 0 |
| 13–24 | Sources | 12 | 0 |
| 25–50 | Capabilities | 26 | 0 |
| 51–60 | Sections and tables | 10 | 0 |
| 61–74 | Ownership and concepts | 14 | 0 |
| 75–88 | Functional coverage | 14 | 0 |
| 89–99 | AI and product safety | 11 | 0 |
| 100–110 | Phase limits | 11 | 0 |
| 111–120 | Registers, quality and publication | 10 | 0 |

## Invariants
Investigate owns analytical context and interpretation. Studio owns Tool, Tool Call and Automation Run. Settings owns provider/environment administration. Shared owns generic jobs, trace, activity, linking, versioning and export. Govern owns authority over real-target action. No Artifact is executed and no dynamic, reverse, debugger or forensic capability is created.

The publication gate is rechecked after the final branch fast-forward; any unreachable commit, README change or PR state change changes the verdict to PARTIAL.
