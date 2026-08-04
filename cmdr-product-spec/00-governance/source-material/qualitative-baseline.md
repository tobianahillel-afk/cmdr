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
  - REQ-UX-010
  - REQ-OBJ-001
  - REQ-OBJ-012
---

# Qualitative Baseline

## Baseline repository Phase 0

781 Markdown ; 559 generic skeletons ; 61 Template-level screens ; 57 objects insufficiently formalized ; 461 `À compléter` ; 289 exact repeated placeholders.

## Results before corrective closure

- Phase 1: PASS.
- Phase 2: PASS.
- Phase 3: PASS.
- Phase 4A functional specification: PARTIAL pending strict table normalization.
- Requirement coverage: 99 conform, 20 partial, 3 absent, 0 contradictory.

## Phase 4A source-reading scope

- 142 accessible source files read in the functional pass;
- 27 original Command files read: 22 active and 5 deprecated;
- 10 active Command screens read without rewrite;
- 5 legacy Work Queue files kept deprecated;
- requested object paths confirmed absent: `service.md`, `exposure.md`, `report.md`, `audit-record.md`.

## Functional Phase 4A measures

| Measure | Before Phase 4A | After corrective closure |
|---|---:|---:|
| Command files | 27 | 69 |
| Active Command files | 22 | 61 |
| Deprecated Command files | 5 | 8 |
| Capability IDs | 0 | 27 |
| Registered capabilities | 0 | 27 |
| Capabilities without owner | N/A | 0 |
| Capabilities without users | N/A | 0 |
| Capabilities without inputs | N/A | 0 |
| Capabilities without outputs | N/A | 0 |
| Capabilities without objects | N/A | 0 |
| Capabilities without classified action | N/A | 0 |
| Capabilities without no-AI alternative | N/A | 0 |
| Capabilities without Given/When/Then | N/A | 0 |
| Current `native` delivery mode | 0 | 0 |
| Current `integrated` delivery mode | 0 | 0 |
| Current `temporary-integration` mode | 0 | 0 |
| Current `planned` delivery mode | 0 | 27 |
| Functional status `defined` | 0 | 26 |
| Functional status `proposed` | 0 | 1 |
| Structured Phase 4A dependencies | 0 | 12 |
| Active competing owners in scope | 2 | 0 |
| Detailed screens rewritten | 0 | 0 |
| Objects changed in corrective closure | 0 | 0 |
| Atomic permission sources changed | 0 | 0 |
| Product code added | 0 | 0 |
| APIs or protocols created | 0 | 0 |
| Font files added | 0 | 0 |
| New Requirement IDs | 0 | 0 |
| Open decisions | 15 | 15 |
| New OPEN decisions | 0 | 0 |

## Corrective template-conformance measures

| Measure | Initial published state | Final state |
|---|---:|---:|
| Capabilities audited | 27 | 27 |
| Sections 1–27 present | 729 | 729 |
| Canonical mandatory tables | 39 / 162 | 162 / 162 |
| Table sections requiring correction | 123 | 0 |
| Mandatory sections represented only by prose | 74 | 0 |
| Mandatory tables with abbreviated headers | 49 | 0 |
| Empty mandatory tables | 0 | 0 |
| Capabilities with strict template verdict PASS | 0 / 27 | 27 / 27 |
| Duplicate Capability IDs | 0 | 0 |
| Concurrent active owners | 0 | 0 |
| New functional contradictions | 0 | 0 |
| Targeted capability/template placeholders | 0 | 0 |
| Targeted empty files | 0 | 0 |
| Broken local links introduced | 0 | 0 |
| Closure gates | PARTIAL | 75 / 75 PASS |

The ten active Command screen templates remain outside the corrective scope and unchanged, because the mission prohibits detailed screen rewriting. No capability, module, template or quality report in the Phase 4A corrective scope contains a generic placeholder.

## Requirement coverage

| State | Before closure | After closure |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| total | 122 | 122 |

Formatting normalization adds quality evidence but does not promote a functional Requirement ID.

## Phase 4A final result

- canonical 27-section template with six mandatory tables;
- 27 immutable Command Capability IDs;
- 27/27 capabilities conforming;
- 162/162 mandatory tables present;
- Capability Register and requirement scores unchanged and coherent;
- one Work Queue workspace and six system views unchanged;
- Customers and Delivery remains `proposed`, `planned` and deployment-dependent;
- all fifteen OPEN decisions remain open; OPEN-009 remains the only resolved historical item;
- no screen, object, permission catalog, API, protocol, code or font changed;
- Phase 4B not started.

## Limits

Phase 4 global and repository global remain PARTIAL. Investigate, Govern, Studio, Settings, Endpoint Agent, journeys, detailed screens, detailed objects, permissions, technical architecture and implementation remain later work.