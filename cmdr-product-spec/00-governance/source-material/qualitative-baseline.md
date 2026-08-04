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
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-UX-010
  - REQ-OBJ-001
  - REQ-OBJ-012
---

# Qualitative Baseline

## Baseline repository Phase 0

781 Markdown ; 559 generic skeletons ; 61 Template-level screens ; 57 objects insufficiently formalized ; 461 `À compléter` ; 289 exact repeated placeholders.

## Phase results

- Phase 1: PASS.
- Phase 2: PASS.
- Phase 3: PASS.
- Phase 4A: PASS after corrective closure.
- Phase 4B.1: PASS after capability, register and quality closure.
- Phase 4B global: PARTIAL.
- Requirement coverage: 99 conform, 20 partial, 3 absent, 0 contradictory.

## Phase 4A corrective closure summary

| Measure | Final Phase 4A |
|---|---:|
| Command capabilities | 27 |
| Sections present | 729 / 729 |
| Mandatory tables | 162 / 162 |
| Delivery status | 26 defined, 1 proposed |
| Delivery mode | 27 planned |
| Closure gates | 75 / 75 PASS |
| Screens/objects/code changed in corrective closure | 0 / 0 / 0 |

The ten active Command screen templates remained outside the corrective scope. Customers and Delivery remained proposed and OPEN-006/010/013 stayed open.

## Phase 4B.1 source and scope controls

- all active and deprecated documents under `07-investigate/` were audited before consolidation;
- fourteen active Investigate screens were read without detailed rewrite;
- the two detached functional commits and temporary work were audited before reuse;
- ownership was verified against governance, object registers and dependent product boundaries;
- absent or insufficient object contracts remain documented in `object-consumption-map.md` rather than invented.

## Phase 4B.1 exact change measures

| Measure | Before Phase 4B.1 | After Phase 4B.1 |
|---|---:|---:|
| Investigate files affected by Phase 4B.1 | 0 | 59 |
| Added Investigate files | 0 | 42 |
| Modified Investigate files | 0 | 17 |
| Deleted Investigate files | 0 | 0 |
| CAP-INV IDs | 0 | 22 |
| Signals and Hunt capabilities | 0 | 8 |
| Cases and Evidence capabilities | 0 | 14 |
| Capabilities without owner | N/A | 0 |
| Capabilities without primary users | N/A | 0 |
| Capabilities without inputs | N/A | 0 |
| Capabilities without outputs | N/A | 0 |
| Capabilities without objects | N/A | 0 |
| Capabilities without classified actions | N/A | 0 |
| Capabilities without no-AI alternative | N/A | 0 |
| Capabilities without Given/When/Then | N/A | 0 |
| Required numbered sections | 0 | 594 / 594 |
| Mandatory S8/S9/S10/S13/S16/S17 tables | 0 | 132 / 132 |
| Empty mandatory tables | N/A | 0 |
| Mandatory sections represented only by prose | N/A | 0 |
| Generic tables copied without adaptation | N/A | 0 |
| Duplicate CAP-INV IDs | 0 | 0 |
| Concurrent active owners in Phase 4B.1 | 0 | 0 |
| Functional status `defined` | 0 | 21 |
| Functional status `proposed` | 0 | 1 (CAP-INV-106) |
| Current `planned` delivery mode | 0 | 22 |
| Current `native`/`integrated`/`temporary-integration` claims | 0 | 0 |
| Old functional modules migrated | 0 | 8 |
| Migration/deprecated documents produced | 0 | 12 |
| Active screens read | 0 | 14 |
| Screens modified minimally | 0 | 0 |
| Detailed screens rewritten | 0 | 0 |
| Object files modified | 0 | 0 |
| Investigate permission overview modified | 0 | 1 |
| Atomic permission sources modified | 0 | 0 |
| APIs or protocols created | 0 | 0 |
| Search/query/forensic engines selected | 0 | 0 |
| Product code added | 0 | 0 |
| Font files added | 0 | 0 |
| CAP-INV-2xx/3xx/4xx/5xx created | 0 | 0 |
| Targeted generic placeholders | N/A | 0 |
| Targeted empty files | N/A | 0 |
| Broken local links introduced | N/A | 0 |
| New Requirement IDs | 0 | 0 |
| Open decisions | 15 | 15 |
| New OPEN decisions | 0 | 0 |

## Phase 4B.1 ownership and concept quality

- Investigate owns Case, Hypothesis, Artifact, Evidence and Finding.
- Command retains Detection, Signal, Alert, Incident and operational Task.
- Govern retains Action Request lifecycle, Decision, Approval, Response Run and Result.
- Shared retains Telemetry Event, Entity, Query, Search Job, Timeline mechanisms, Saved Views and Reporting Engine according to the registers.
- CMDR Studio retains Workflow, Agent and Automation Run.
- Event, Artifact and Attachment do not become Evidence automatically.
- Hypothesis does not become Finding automatically.
- Finding is distinct from Decision and Result.
- Case Queue is distinct from the Command Work Queue.

## Phase 4B.1 quality closure

| Measure | Final state |
|---|---:|
| Capability specifications audited | 22 / 22 |
| Front matter valid | 22 / 22 |
| Sections 1–27 | 594 / 594 |
| S8 tables | 22 / 22 |
| S9 tables | 22 / 22 |
| S10 tables | 22 / 22 |
| S13 tables | 22 / 22 |
| S16 tables | 22 / 22 |
| S17 tables | 22 / 22 |
| Mandatory tables total | 132 / 132 |
| OPEN-014 options retained | 4 / 4 |
| Quality gates | 140 / 140 PASS |

## Requirement coverage

| State | After Phase 4A | After Phase 4B.1 |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| total | 122 | 122 |

Capability definition adds evidence but does not prove implementation or complete requirements whose owner phases remain unfinished.

## Phase 4B.1 final result

- one canonical capability template reused;
- 22 immutable CAP-INV IDs and 22 canonical files;
- 8 Signals and Hunt plus 14 Cases and Evidence capabilities;
- 594/594 sections and 132/132 mandatory tables;
- full Artifact versus Attachment analysis with OPEN-014 left open;
- transitions Signal→investigation, Incident→Case, Search/Hunt→Case, Artifact→Evidence, Evidence→Finding, Finding→Action Request, Govern and Result return defined;
- Capability and Dependency Registers updated;
- no screen rewrite, object schema, atomic permission matrix, API, protocol, engine choice, code or font;
- Phase 4B.2 and Phase 4B.3 not started.

## Limits

Phase 4B, Phase 4 global and repository global remain PARTIAL. Collection and Live Response, Analysis Workbench, Detection Engineering, Intelligence, detailed screens, object schemas, permissions, technical architecture and implementation remain later work.
