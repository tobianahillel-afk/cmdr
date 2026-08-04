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

## Repository baseline
Phase 0 recorded 781 Markdown files, 559 generic skeletons, 61 Template-level screens, 57 insufficiently formalized objects, 461 `À compléter` occurrences and 289 repeated placeholders.

## Completed evidence before Phase 4B.2A
- Phases 0–3: PASS.
- Phase 4A Command: PASS.
- Phase 4B.1 Investigate: PASS.
- Requirement coverage: 99 conform, 20 partial, 3 absent, 0 contradictory.
- Registered capabilities: 49 (27 Command + 22 Investigate).
- Open decisions: 15.

## Phase 4B.2A source scope
- source documents read integrally: **140**;
- active Investigate screens: **14**;
- Platform Settings screens: **2**;
- Govern screens: **2**;
- total screen specifications read: **18**;
- initial canonical Collection and Live Response module: absent;
- absent object files: Endpoint, Collection Job, Live Session, Endpoint Operation, Operation Result, Custody Record and Provenance Record;
- absent canonical screens: Endpoint Search, Collection, Live Session, Jobs and Terminal;
- renamed source: Technical Workbench Shell is under `03-design-system/layouts/`;
- absent requested screen path: Platform Settings System Health screen.

## Corrective substantive audit

The initial publication at `cd93be822a06a721884a243ec34085270b323e4b` satisfied structural counts but failed the substantive generic-table gate for `CAP-INV-204` through `CAP-INV-212`. The affected S9, S12, S13, S17, permissions and acceptance criteria were replaced in two corrective commits. Final counts are unchanged, but **generic mandatory tables are now verified at 0 rather than inferred from structure**.

## Phase 4B.2A measures

| Measure | Before | Final |
|---|---:|---:|
| Collection and Live Response files | 0 | 30 |
| Active module files | 0 | 30 |
| Deprecated module files | 0 | 0 |
| Generic active module files | 0 | 0 |
| Active placeholders | 0 | 0 |
| CAP-INV-2xx | 0 | 15 |
| Investigate capabilities | 22 | 37 |
| Total registered capabilities | 49 | 64 |
| CAP-INV-3xx added | 0 | 0 |
| Capabilities without owner | N/A | 0 |
| Without users | N/A | 0 |
| Without inputs | N/A | 0 |
| Without outputs | N/A | 0 |
| Without objects | N/A | 0 |
| Without classified action | N/A | 0 |
| Without no-AI alternative | N/A | 0 |
| Without GWT criteria | N/A | 0 |
| Numbered sections expected/present | 0 | 405/405 |
| Mandatory tables expected/present | 0 | 90/90 |
| Empty mandatory tables | 0 | 0 |
| Prose-only mandatory sections | 0 | 0 |
| Generic mandatory tables | 0 | 0 |
| Delivery status `defined` in 4B.2A | 0 | 15 |
| Delivery mode `planned` in 4B.2A | 0 | 15 |
| Current native/integrated/temporary claims | 0 | 0 |
| Old competing Investigate modules migrated | 0 | 0; none existed |
| Active duplicate architectures | 0 | 0 |
| Concurrent owners | 0 | 0 |
| Screens read | 0 | 18 |
| Screens modified minimally | 0 | 0 |
| Detailed screens rewritten | 0 | 0 |
| Object source files modified | 0 | 0 |
| Atomic permission sources modified | 0 | 0 |
| APIs / protocols / engines / exact commands | 0 | 0 |
| Product code / fonts | 0 | 0 |
| Requirement IDs | 122 | 122 |
| Open decisions | 15 | 15 |
| Phase 4B.2B content | 0 | 0 |

## Functional result
The module distinguishes Endpoint, Agent, Fleet and Policy; Request, Job and Background Job; Live Session, Console, Automation Run and Response Run; Operation Result and Govern Result; Artifact and Evidence. Offline, partial, reconnect, cancel, retry, custody and Govern handoffs are explicit.

## Limits
Phase 4B.2 remains PARTIAL because Analysis Workbench is not started. Phase 4B, Phase 4 and repository maturity remain PARTIAL. No implementation is established.
