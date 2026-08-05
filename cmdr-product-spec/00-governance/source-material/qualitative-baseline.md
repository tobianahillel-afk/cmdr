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
Remote start: `516226778850a5d0ef4e6bbec6afc3175af90aba`. It contained PASS through Disk and Filesystem Forensics, 143 registered capabilities, 116 Investigate capabilities and 79 CAP-INV-3xx.

## Source audit
- Governance, template, registers, ADRs, ownership, status and terminology were reread.
- Foundation, Signals/Hunt, Cases/Evidence, Collection, Static, Dynamic, Reverse/Debugger, Memory, Disk, Command/SIEM, Endpoint Agent, Settings, Studio, Govern, Shared and Technical Workbench boundaries were audited.
- Network-specific active or generic sources were distributed across Collection Network Capture Request, Dynamic Sandbox Network Behavior, Memory Network State, Endpoint Agent network telemetry/connections, Event Search, Entity Graph and Shared mechanisms.
- No active standalone Network Forensics functional module or autonomous screen existed before the phase.
- Competing functional sources migrated to deprecated pointers: **0**; none was equivalent to full Network Forensics.
- Existing screens read: **10**; screen specifications modified: **0**; detailed rewrites: **0**; new Screen IDs: **0**.

## Measures
| Measure | Before | After 4B.2 closure |
|---|---:|---:|
| Network Forensics targeted module files | 0 | 30 |
| Active Network functional documents | 0 | 30 |
| Deprecated Network functional pointers | 0 | 0 |
| Generic active Network Forensics documents / placeholders | 0 / 0 | 0 / 0 |
| CAP-INV-3xx | 79 | 97 |
| Investigate capabilities | 116 | 134 |
| Registered capabilities | 143 | 161 |
| Defined / proposed | 141 / 2 | 159 / 2 |
| Delivery mode planned | 143 | 161 |
| New Network capability documents | 0 | 18 |
| Network sections expected/present | 0 | 486/486 |
| Network mandatory tables expected/present | 0 | 108/108 |
| Empty / prose-only / generic mandatory tables | 0 | 0 |
| Capability files missing owner/user/input/output/object/action/no-AI/GWT | 0 | 0 |
| Duplicate IDs / active duplicates / concurrent owners | 0 | 0 |
| Legacy functional sources migrated | 0 | 0 |
| Screens read / specs modified / rewritten / new IDs | 10 / 0 / 0 / 0 | 10 / 0 / 0 / 0 |
| Object map updated / canonical object files created | 0 / 0 | 1 / 0 |
| Functional permission documents / atomic permissions | 0 / 0 | 1 / 0 |
| APIs / protocols / engines / tools imposed / commands | 0 | 0 |
| Product code / detection rules / Intelligence objects | 0 | 0 |
| Broken targeted links / empty targeted files | 0 | 0 |
| Requirement IDs / OPEN | 122 / 15 | 122 / 15 |
| Phase 4B.3 content | 0 | 0 |
| Phase 4B.2 capabilities / sections / tables | 94 / 2538 / 564 | 112 / 3024 / 672 |
| Investigate capabilities / sections / tables | 116 / 3132 / 696 | 134 / 3618 / 804 |
| Command + Investigate capabilities / sections / tables | 143 / 3861 / 858 | 161 / 4347 / 966 |
| Phase 4B.2 sub-phase reports present | 6 | 8 including Network and closure |
| Incoherent reports / statuses after closure | 0 | 0 |

## Maturity
Phase 4B.2B.3B.2, Phase 4B.2B.3B, Phase 4B.2B.3, Phase 4B.2B and Phase 4B.2 are PASS after publication verification. Phase 4B, Phase 4 and global maturity remain PARTIAL because Phase 4B.3 Detection Engineering and Intelligence is not started. No implementation is established.
