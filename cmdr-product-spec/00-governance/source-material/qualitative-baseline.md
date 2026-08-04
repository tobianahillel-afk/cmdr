---
id: qualitative-baseline
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-04
source-of-truth: source-material
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-UX-010
---
# Qualitative Baseline

## Corrective starting state
The exact remote start was `2bc2d37e2d726d15e1ee6af1ce3ae633703b4cc3`. It contained PASS evidence through Phase 4B.2A only: 64 registered capabilities, 37 Investigate capabilities and no CAP-INV-3xx. A previous conversational report claiming published Phase 4B.2B.1 work was non-authoritative because no remote commit or file existed; the sub-phase was reconstructed before 4B.2B.2A.

## Source audit
- Combined governance, Investigate, Studio, Settings, Govern, Shared, Design System, Experience Architecture and screen manifest: **114 distinct paths**.
- Relevant technical screens read: **9**; detailed rewrites: **0**.
- Missing canonical concepts recorded: Analysis Session, Analysis Result, Derived Artifact, Dynamic Analysis Session, Sandbox Run, Runtime Artifact, Behavioral/Process/Network/System Observation and Provenance Record.

## Measures
| Measure | Before | After 4B.2B.2A |
|---|---:|---:|
| CAP-INV-3xx | 0 | 28 |
| Investigate capabilities | 37 | 65 |
| Registered capabilities | 64 | 92 |
| Defined / proposed | 62 / 2 | 90 / 2 |
| Delivery mode planned | 64 | 92 |
| Static / Dynamic capability documents | 0 / 0 | 13 / 15 |
| Sections expected/present | 0 | 756/756 |
| Mandatory tables expected/present | 0 | 168/168 |
| Empty / prose-only / generic mandatory tables | 0 | 0 |
| Duplicate IDs / concurrent owners | 0 | 0 |
| Active module placeholders | 2 generic static docs | 0; two deprecated pointers |
| Screens read / modified / rewritten | 0 | 9 / 0 / 0 |
| Object / atomic permission sources modified | 0 | 0 |
| APIs / protocols / engines / hypervisors / commands | 0 | 0 |
| Product code / fonts | 0 | 0 |
| Requirement IDs / OPEN | 122 / 15 | 122 / 15 |
| Reverse/Debugger/Forensics capabilities | 0 | 0 |

Phase 4B.2B remains PARTIAL because Reverse Engineering, Debugger and forensic work are not complete. No implementation is established.
