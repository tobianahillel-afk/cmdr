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
Remote start: `b6a16d66d66117273f0490afb362038f63e7934e`. It contained PASS through Memory Forensics, 126 registered capabilities, 99 Investigate capabilities and 62 CAP-INV-3xx.

## Source audit
- Inherited relevant source manifest: 120 references from the previous phase.
- Phase-specific Disk/Filesystem sources directly audited: 27 documentary sources before construction, supplemented by active screen and cross-analysis checks.
- Legacy Disk sources: two generic functional documents and one active screen.
- Disk Image is canonical and owned by Investigate; Disk Forensics Session and observation concepts are gaps, not new canonical schemas.
- Relevant screens read: 9; screen specifications modified: 0; detailed rewrites: 0.
- Network Forensics is referenced only as a future handoff and was not specified.

## Measures
| Measure | Before | After 4B.2B.3B.1 |
|---|---:|---:|
| Disk/Filesystem targeted files | 3 | 31 |
| Active functional/screen files | 3 | 29 |
| Deprecated functional pointers | 0 | 2 |
| Generic active functional files / placeholders | 2 / 2 | 0 / 0 |
| CAP-INV-3xx | 62 | 79 |
| Investigate capabilities | 99 | 116 |
| Registered capabilities | 126 | 143 |
| Defined / proposed | 124 / 2 | 141 / 2 |
| Delivery mode planned | 126 | 143 |
| New Disk/Filesystem capability documents | 0 | 17 |
| Sections expected/present | 0 | 459/459 |
| Mandatory tables expected/present | 0 | 102/102 |
| Empty / prose-only / generic mandatory tables | 0 | 0 |
| Capability files missing owner/user/input/output/object/action/no-AI/GWT | 0 | 0 |
| Duplicate IDs / active duplicates / concurrent owners | 0 | 0 |
| Legacy functional sources migrated | 0 | 2 |
| Screens read / specs modified / rewritten / new IDs | 9 / 0 / 0 / 0 | 9 / 0 / 0 / 0 |
| Object maps modified / canonical object files created | 0 / 0 | 1 / 0 |
| Functional permission docs / atomic permission sources | 0 / 0 | 1 / 0 |
| APIs / protocols / engines / tools imposed / commands | 0 | 0 |
| Product code / fonts | 0 | 0 |
| Broken targeted links / empty targeted files | 0 | 0 |
| Requirement IDs / OPEN | 122 / 15 | 122 / 15 |
| Full Network Forensics capabilities / 4B.2B.3B.2 content | 0 / 0 | 0 / 0 |

Phase 4B.2B.3B.1 is PASS. Phase 4B.2B.3B, Phase 4B.2B.3 and Phase 4B.2B remain PARTIAL because Network Forensics and closure are not started. No implementation is established.
