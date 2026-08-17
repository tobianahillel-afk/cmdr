---
id: phase-4b2a-collection-live-response-capability-conformance
domain: 16-quality-and-validation
status: draft
owner: QA and Traceability Lead
updated: 2026-08-04
source-of-truth: quality-report
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-AI-002
  - REQ-UX-010
---
# Phase 4B.2A — Collection and Live Response capability conformance

## Verdict

**PASS — 156/156 gates after corrective revalidation.** Phase 4B.2 global, Phase 4B global, Phase 4 global and repository maturity remain PARTIAL. Phase 4B.2B and Phase 4B.3 are not started.

## Corrective validation history

The first published report at `cd93be822a06a721884a243ec34085270b323e4b` was a false positive for gate 58. A manual substantive audit found generic copied scaffolding in `CAP-INV-204` through `CAP-INV-212`, especially in S9, S12, S13, S17, permissions and acceptance criteria. The sub-phase was therefore treated as **PARTIAL** until correction.

Corrective commits:

- `33aa8337e654e526873c7079b7bf4273956489ab` — `docs: enforce collection capability contracts`;
- `17709e80c498bb8d8f845f98251433a30debb934` — `docs: enforce live response capability contracts`.

The nine affected contracts were replaced with capability-specific inputs, objects, actions, automation alternatives, outputs, transitions, permissions, limits and Given/When/Then scenarios. No ID, owner, status, delivery mode, screen, object schema or product boundary changed.

## Scope and evidence

- Initial Phase 4B.2A base: `71b45ef99feb568a2ef0cbabfa0fbe008fcbaea7`.
- Source documents read integrally: **140**.
- Screen specifications read: **18**; modified: **0**; detailed rewrites: **0**.
- Canonical module files: **30**; active: **30**; deprecated: **0**; generic active: **0**; placeholders: **0**.
- Canonical capabilities: **15** (`CAP-INV-201..215`).
- Numbered sections expected/present: **405/405**.
- Mandatory tables expected/present: **90/90**.
- Empty mandatory tables: **0**.
- Prose-only required table sections: **0**.
- Generic unadapted tables after correction: **0**.
- Duplicate IDs: **0**.
- Concurrent owners: **0**.
- Active contradictions introduced: **0**.
- Object files modified: **0**.
- Atomic permission sources modified: **0**.
- CAP-INV-3xx created: **0**.
- APIs, protocols, engines, exact commands, product code and fonts: **0**.

## Source domains

| Domain | Files read |
|---|---:|
| Governance, source material, registers and ADRs | 27 |
| Phase 4B.1 contracts | 17 |
| Endpoint Agent | 40 |
| Platform Settings including screens | 10 |
| Govern including screens | 9 |
| CMDR Studio | 4 |
| Shared Capabilities | 8 |
| Design System and canonical template | 6 |
| Domain objects | 5 |
| Investigate screens | 14 |
| **Total** | **140** |

## Absent or renamed paths

- `07-investigate/modules/collection-and-live-response/` was absent before this sub-phase.
- Canonical object files were absent for Endpoint, Collection Job, Live Session, Endpoint Operation, Operation Result, Custody Record and Provenance Record.
- Canonical Investigate screens were absent for Endpoint Search, Collection, Live Session, Jobs and Terminal.
- Technical Workbench Shell exists under `03-design-system/layouts/`, not `03-design-system/shells/`.
- A Platform Settings System Health screen was not present under the requested path.
- Absence was recorded; no object, screen or API was invented.

## Capability count decision

The target list of fifteen capabilities is retained. `CAP-INV-215` is not a duplicate of `CAP-INV-113`: it captures Endpoint-specific target, policy, availability, impact, alternatives and rollback, then delegates the Action Request lifecycle to `CAP-INV-113` and Govern.

## Capability-template matrix

| Capability ID | Sections 1–27 | S8 | S9 | S10 | S13 | S16 | S17 | Front matter | Verdict |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---|
| CAP-INV-201 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-202 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-203 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-204 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-205 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-206 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-207 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-208 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-209 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-210 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-211 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-212 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-213 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-214 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |
| CAP-INV-215 | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS | PASS |

## Corrected capability specificity

| Capability | Corrected specificity |
|---|---|
| CAP-INV-204 | versioned triage profiles, per-category partial/error/retry and Artifact outputs |
| CAP-INV-205 | bounded file/directory/pattern scope, preview, missing/locked handling and no deployment |
| CAP-INV-206 | read-only snapshots, process/session/service/connection inspection and containment separation |
| CAP-INV-207 | platform/capability support, impact constraints, Memory Image/Artifact handoff and OPEN-005/008 |
| CAP-INV-208 | bounded duration/volume/filter, explicit start/stop and loss/error reporting |
| CAP-INV-209 | participants, expiry, inactivity, reconnect/conflict and Live Session/Run separation |
| CAP-INV-210 | operation catalogue, classes 0–4, interruption and Govern redirection without commands |
| CAP-INV-211 | upload/download purpose, temporary transfer, Artifact/Attachment distinction and no deployment |
| CAP-INV-212 | output/error/partial review, verification/dispute and Operation Result/Govern Result separation |

## Metrics

| Measure | Before | Final |
|---|---:|---:|
| Collection and Live Response files | 0 | 30 |
| Active module files | 0 | 30 |
| Deprecated module files | 0 | 0 |
| Generic active module files | 0 | 0 |
| Active placeholders | 0 | 0 |
| CAP-INV-2xx | 0 | 15 |
| Investigate capabilities | 22 | 37 |
| Registered capabilities | 49 | 64 |
| Capabilities without owner/user/input/output/object/action/no-AI/GWT | N/A | 0 |
| Sections expected/present | 0 | 405/405 |
| Mandatory tables expected/present | 0 | 90/90 |
| Empty mandatory tables | 0 | 0 |
| Prose-only mandatory sections | 0 | 0 |
| Generic mandatory tables | 0 | 0 |
| Delivery status `defined` | 0 | 15 |
| Delivery mode `planned` | 0 | 15 |
| Native/integrated/current implementation claims | 0 | 0 |
| Competing active module architectures | 0 | 0 |
| Screens read / modified / rewritten | 0 | 18 / 0 / 0 |
| Object or atomic permission sources modified | 0 | 0 |
| APIs/protocols/engines/commands/code/fonts | 0 | 0 |
| Requirement IDs | 122 | 122 |
| Open decisions | 15 | 15 |
| CAP-INV-3xx / Phase 4B.2B content | 0 | 0 |

## 156 closure gates

| # | Group | Gate | Verdict |
|---:|---|---|---|
| 1 | Git | Repository correct | PASS |
| 2 | Git | Branch correct | PASS |
| 3 | Git | PR correct | PASS |
| 4 | Git | Base `main` | PASS |
| 5 | Git | PR open | PASS |
| 6 | Git | PR Draft | PASS |
| 7 | Git | PR unmerged | PASS |
| 8 | Git | No force-push | PASS |
| 9 | Git | History not rewritten | PASS |
| 10 | Git | Root README unchanged | PASS |
| 11 | Git | `main` unchanged | PASS |
| 12 | Git | Final remote SHA verified | PASS |
| 13 | Sources | All active Collection and Live Response files read | PASS |
| 14 | Sources | All competing legacy modules read | PASS |
| 15 | Sources | Relevant Endpoint Agent documents read | PASS |
| 16 | Sources | Relevant Platform Settings documents read | PASS |
| 17 | Sources | Relevant Govern documents read | PASS |
| 18 | Sources | Relevant screens read | PASS |
| 19 | Sources | Required objects read or absences recorded | PASS |
| 20 | Sources | Phase 4B.1 boundaries reread | PASS |
| 21 | Capabilities | CAP-INV-2xx convention respected | PASS |
| 22 | Capabilities | No duplicate 2xx ID | PASS |
| 23 | Capabilities | No 3xx ID created | PASS |
| 24 | Capabilities | Final capability count justified | PASS |
| 25 | Capabilities | One canonical file per capability | PASS |
| 26 | Capabilities | One owner per capability | PASS |
| 27 | Capabilities | Users per capability | PASS |
| 28 | Capabilities | User problem per capability | PASS |
| 29 | Capabilities | Objectives per capability | PASS |
| 30 | Capabilities | Non-objectives per capability | PASS |
| 31 | Capabilities | Inputs per capability | PASS |
| 32 | Capabilities | Read objects per capability | PASS |
| 33 | Capabilities | Created/modified objects per capability | PASS |
| 34 | Capabilities | Classified actions per capability | PASS |
| 35 | Capabilities | S13 matrix per capability | PASS |
| 36 | Capabilities | Specific states per capability | PASS |
| 37 | Capabilities | Outputs per capability | PASS |
| 38 | Capabilities | Transitions per capability | PASS |
| 39 | Capabilities | Source of truth per capability | PASS |
| 40 | Capabilities | Provenance per capability | PASS |
| 41 | Capabilities | Functional permissions per capability | PASS |
| 42 | Capabilities | Limits and errors per capability | PASS |
| 43 | Capabilities | Metrics per capability | PASS |
| 44 | Capabilities | Delivery classification per capability | PASS |
| 45 | Capabilities | Given/When/Then per capability | PASS |
| 46 | Capabilities | No-AI alternative per capability | PASS |
| 47 | Capabilities | Documentary consumers per capability | PASS |
| 48 | Tables and sections | All S8 present | PASS |
| 49 | Tables and sections | All S9 present | PASS |
| 50 | Tables and sections | All S10 present | PASS |
| 51 | Tables and sections | All S13 present | PASS |
| 52 | Tables and sections | All S16 present | PASS |
| 53 | Tables and sections | All S17 present | PASS |
| 54 | Tables and sections | Expected table count reached | PASS |
| 55 | Tables and sections | Expected section count reached | PASS |
| 56 | Tables and sections | No empty table | PASS |
| 57 | Tables and sections | No prose replacing required table | PASS |
| 58 | Tables and sections | No generic unadapted table | PASS |
| 59 | Ownership | Investigate owns collection context | PASS |
| 60 | Ownership | Investigate owns collection/Case relation | PASS |
| 61 | Ownership | Investigate does not own Fleet | PASS |
| 62 | Ownership | Investigate does not own Endpoint Policy | PASS |
| 63 | Ownership | Platform Settings owns Fleet | PASS |
| 64 | Ownership | Platform Settings owns Endpoint Policies | PASS |
| 65 | Ownership | Endpoint Agent executes locally | PASS |
| 66 | Ownership | Endpoint Agent does not own Case | PASS |
| 67 | Ownership | Govern owns Decision | PASS |
| 68 | Ownership | Govern owns Response Run | PASS |
| 69 | Ownership | Govern owns Result | PASS |
| 70 | Ownership | Studio owns Automation Run | PASS |
| 71 | Ownership | Command owns Incident | PASS |
| 72 | Ownership | Artifact remains Investigate | PASS |
| 73 | Ownership | Evidence remains Investigate | PASS |
| 74 | Ownership | Finding remains Investigate | PASS |
| 75 | Concepts | Endpoint distinct from Fleet | PASS |
| 76 | Concepts | Endpoint distinct from Endpoint Agent | PASS |
| 77 | Concepts | Collection Request distinct from Collection Job | PASS |
| 78 | Concepts | Collection Job distinct from generic Background Job | PASS |
| 79 | Concepts | Collection Result distinct from Artifact | PASS |
| 80 | Concepts | Artifact distinct from Evidence | PASS |
| 81 | Concepts | Live Session distinct from Terminal | PASS |
| 82 | Concepts | Live Session distinct from Automation Run | PASS |
| 83 | Concepts | Live Session distinct from Response Run | PASS |
| 84 | Concepts | Tool Call distinct from Endpoint Operation | PASS |
| 85 | Concepts | Operation Result distinct from Govern Result | PASS |
| 86 | Concepts | Action Request distinct from Decision | PASS |
| 87 | Concepts | Response Run distinct from Automation Run | PASS |
| 88 | Concepts | Containment request distinct from containment execution | PASS |
| 89 | Module | Endpoint Investigation Context defined | PASS |
| 90 | Module | Collection Scope defined | PASS |
| 91 | Module | Collection Job Management defined | PASS |
| 92 | Module | Triage Collection defined | PASS |
| 93 | Module | File Acquisition defined | PASS |
| 94 | Module | Process/System Inspection defined | PASS |
| 95 | Module | Memory Acquisition Request defined | PASS |
| 96 | Module | Network Capture Request defined | PASS |
| 97 | Module | Live Session Management defined | PASS |
| 98 | Module | Interactive Operations defined | PASS |
| 99 | Module | Session File Transfer defined | PASS |
| 100 | Module | Operation Result Handling defined | PASS |
| 101 | Module | Integrity and Custody defined | PASS |
| 102 | Module | Provenance defined | PASS |
| 103 | Module | Containment Request handled without duplication | PASS |
| 104 | AI and product security | AI optional | PASS |
| 105 | AI and product security | No mandatory chatbot | PASS |
| 106 | AI and product security | No silent session opening | PASS |
| 107 | AI and product security | No silent sensitive operation | PASS |
| 108 | AI and product security | No silent scope expansion | PASS |
| 109 | AI and product security | No automatic Artifact-to-Evidence qualification | PASS |
| 110 | AI and product security | No automatic Finding confirmation | PASS |
| 111 | AI and product security | No containment without Govern when required | PASS |
| 112 | AI and product security | Automated provenance visible | PASS |
| 113 | AI and product security | No-AI alternative present | PASS |
| 114 | AI and product security | No self-granted permission | PASS |
| 115 | AI and product security | No Govern bypass | PASS |
| 116 | Limits | No detailed screen rewritten | PASS |
| 117 | Limits | No complete object created | PASS |
| 118 | Limits | No JSON Schema | PASS |
| 119 | Limits | No final cardinality | PASS |
| 120 | Limits | No final object state machine | PASS |
| 121 | Limits | No atomic permission finalized | PASS |
| 122 | Limits | No global namespace normalization | PASS |
| 123 | Limits | No API created | PASS |
| 124 | Limits | No protocol created | PASS |
| 125 | Limits | No collection engine chosen | PASS |
| 126 | Limits | No memory engine chosen | PASS |
| 127 | Limits | No network engine chosen | PASS |
| 128 | Limits | No exact system command | PASS |
| 129 | Limits | No product code | PASS |
| 130 | Limits | No font file | PASS |
| 131 | Limits | No Analysis Workbench capability created | PASS |
| 132 | Limits | Phase 4B.2B not started | PASS |
| 133 | Limits | Phase 4B.3 not started | PASS |
| 134 | Registers and quality | Capability Register updated | PASS |
| 135 | Registers and quality | Dependency Register updated | PASS |
| 136 | Registers and quality | Object Consumption Map updated | PASS |
| 137 | Registers and quality | Action Classification updated | PASS |
| 138 | Registers and quality | Automation and AI Model updated | PASS |
| 139 | Registers and quality | Cross-product links updated | PASS |
| 140 | Registers and quality | Screen Capability Map updated | PASS |
| 141 | Registers and quality | Requirements Matrix updated | PASS |
| 142 | Registers and quality | Baseline updated | PASS |
| 143 | Registers and quality | STATUS updated | PASS |
| 144 | Registers and quality | CHANGELOG updated | PASS |
| 145 | Registers and quality | PR description coherent | PASS |
| 146 | Registers and quality | No active generic placeholder | PASS |
| 147 | Registers and quality | No targeted empty file | PASS |
| 148 | Registers and quality | No introduced broken local link | PASS |
| 149 | Registers and quality | No concurrent owner | PASS |
| 150 | Registers and quality | No Shared Capability duplication | PASS |
| 151 | Registers and quality | No document artificially Validated | PASS |
| 152 | Registers and quality | Fifteen OPEN decisions coherent | PASS |
| 153 | Registers and quality | No unnecessary new OPEN | PASS |
| 154 | Registers and quality | Conformance report published | PASS |
| 155 | Registers and quality | Final commits reachable from branch | PASS |
| 156 | Registers and quality | Local and remote SHA coherent | PASS |

## Publication rule

The publication gates are rechecked after fast-forward. Any remote SHA mismatch, README change, PR state change, unreachable commit or failed gate reverts the sub-phase verdict to PARTIAL.

## Final boundaries

Fleet and Policies remain Platform Settings; local execution remains Endpoint Agent; Case/Artifact/Evidence/Finding remain Investigate; Decision/Response Run/Result remain Govern; Automation Run remains Studio; generic jobs, linking, timeline, trace, notifications and export remain Shared.
