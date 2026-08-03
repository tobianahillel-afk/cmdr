---
id: qualitative-baseline
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-03
source-of-truth: source-material
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
---
# Qualitative Baseline

## Repository-wide Phase 0 baseline

| Metric | Baseline |
|---|---:|
| Markdown files before Phase 0 | 781 |
| Draft | 767 |
| In Review | 0 |
| Validated or approved | 0 |
| Implemented | 0 |
| Deprecated | 2 |
| Without status | 12 |
| Generic document skeleton | 559 files |
| Template-level screens | 61 |
| Insufficiently formalized objects | 57 |
| `À compléter` occurrences | 461 |
| Exact repeated placeholder | 289 |
| Files using one of three generic blocks | 677 / 781 |
| Estimated generic or repeated text | approximately 60% |

## Phase 1 scoped measurements

The measurement scope is the active governance and product-vision narrative documents rewritten in Phase 1. Generated registers and later-phase documents are not counted as substantively corrected.

| Metric | Before | After |
|---|---:|---:|
| Repeated generic skeleton in targeted `00-governance/` and `01-product-vision/` documents | 19 | 0 |
| `À compléter` in these two domains | 6 | 0 |
| Targeted narrative documents without explicit Requirement IDs | 34 | 0 |
| Substantive product-boundary documents | 0 | 1 canonical matrix plus supporting portfolio map |
| Substantive role definitions | 0 | 16 |
| Structured product risks | 0 | 19 |
| Structured open decisions | 15 summary rows | 14 open records plus 1 resolved history record |
| Product README boundary contradictions corrected | 0 | 7 |
| ADRs with decision-specific analysis | 0 / 4 | 7 / 7 |

Five generated governance registers still do not carry front matter. They remain indexes, not narrative specifications; their status normalization is a documented follow-up rather than hidden.

## Requirement coverage

### Published Phase 0 summary

| State | Count |
|---|---:|
| conform | 0 |
| partial | 68 |
| absent | 47 |
| contradictory | 7 |

### Reconciled Phase 0 rows and Phase 1 result

| State | Reconciled before | After Phase 1 |
|---|---:|---:|
| conform | 0 | 81 |
| partial | 54 | 32 |
| absent | 61 | 6 |
| contradictory | 7 | 3 |
| total | 122 | 122 |

The original summary was arithmetically inconsistent with its 122 rows. The matrix preserves the published values and records the reconciled recount.

## Contradictions resolved

- source-material precedence;
- product and capability ownership at vision level;
- AI as optional augmentation rather than required interface;
- Studio as owner of agentic concepts, not operational decisions;
- Endpoint Agent as a target native EDR rather than an implemented complete EDR;
- capability delivery taxonomy and review authority;
- product exclusions and consumer restrictions.

## Contradictions remaining

- generic Saved Views versus Work Queue Saved Views in existing screen specifications;
- Work Queue variants represented as separate screens;
- permission namespace duplication;
- Tool Call and Automation Run object files remain absent until the detailed functional model phase.

## Outcome

Phase 1 is assessed independently after link, status, traceability and GitHub-state checks. Even with a Phase 1 PASS, the repository remains globally PARTIAL because brand, UX, product modules, journeys, screens, detailed objects, permissions and implementation have not completed their owning phases.
