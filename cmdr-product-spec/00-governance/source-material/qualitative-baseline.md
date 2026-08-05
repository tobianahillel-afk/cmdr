---
id: qualitative-baseline
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-05
source-of-truth: source-material
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-UX-010
---
# Qualitative Baseline

## Starting state
Remote start: `eb60be74cdd4d3b08a37d5ac0a22d781b51aef86`. It contained PASS through Dynamic Sandbox, 92 registered capabilities, 65 Investigate capabilities and 28 CAP-INV-3xx.

## Source audit
- Inherited relevant source manifest: 114 distinct paths.
- Phase-specific Reverse/Debugger sources newly audited: 4 functional sources and 2 active screens.
- Reviewed source references for this phase: 120; the six phase-specific paths are listed in the conformance report.
- Relevant technical screens read: 9; screen specifications modified: 0; detailed rewrites: 0.
- Missing concepts recorded, not modeled: Reverse Analysis Session, Code Location, Function, Symbol, Cross Reference, graphs, types, structures, Debugger Session, Breakpoint, Runtime Snapshot, Stack Frame, Memory Region, Module Observation, Exception Event, Debug Trace, Patch Hypothesis and Reproducibility Assessment.

## Measures
| Measure | Before | After 4B.2B.2B |
|---|---:|---:|
| Reverse/Debugger targeted files | 6 | 35 |
| Active functional/screen files | 6 | 31 |
| Deprecated functional pointers | 0 | 4 |
| Generic active functional files | 4 | 0 |
| Active placeholders | 4 | 0 |
| CAP-INV-3xx | 28 | 46 |
| Investigate capabilities | 65 | 83 |
| Registered capabilities | 92 | 110 |
| Defined / proposed | 90 / 2 | 108 / 2 |
| Delivery mode planned | 92 | 110 |
| New Reverse/Debugger capability documents | 0 | 18 |
| Sections expected/present | 0 | 486/486 |
| Mandatory tables expected/present | 0 | 108/108 |
| Empty / prose-only / generic mandatory tables | 0 | 0 |
| Capability files missing owner/user/input/output/object/action/no-AI/GWT | 0 | 0 |
| Duplicate IDs / active duplicates / concurrent owners | 0 | 0 |
| Legacy functional sources migrated | 0 | 4 |
| Screens read / screen specs modified / rewritten / new IDs | 9 / 0 / 0 / 0 | 9 / 0 / 0 / 0 |
| Object maps modified / canonical object files created | 0 / 0 | 1 / 0 |
| Functional permission docs / atomic permission sources | 0 / 0 | 1 / 0 |
| APIs / protocols / engines / debuggers / commands | 0 | 0 |
| Product code / fonts | 0 | 0 |
| Broken local links / empty targeted files | 0 | 0 |
| Requirement IDs / OPEN | 122 / 15 | 122 / 15 |
| Forensics capabilities / 4B.2B.3 content | 0 / 0 | 0 / 0 |

Phase 4B.2B.2 is PASS. Phase 4B.2B remains PARTIAL because Forensics is not started. No implementation is established.
