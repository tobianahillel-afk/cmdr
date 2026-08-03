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
  - REQ-BRAND-001
  - REQ-BRAND-008
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

## Phase 1 scoped result

| Metric | Before | After |
|---|---:|---:|
| Generic skeleton in targeted governance and vision narratives | 19 | 0 |
| `À compléter` in those targeted domains | 6 | 0 |
| Targeted narratives without Requirement IDs | 34 | 0 |
| Substantive role definitions | 0 | 16 |
| Structured product risks | 0 | 19 |
| Decision-specific ADRs | 0 / 4 | 7 / 7 |

## Phase 2 scoped measurements

The Phase 2 measurement scope includes all 34 brand files that existed at the start of the phase, the new canonical brand documents and product entry points, and the minimal Design System foundation references changed to preserve ownership.

| Metric | Before Phase 2 | After Phase 2 |
|---|---:|---:|
| Existing active brand documents using the generic skeleton | 25 / 34 | 0 |
| Active placeholder instructions or unresolved `À compléter` fields in `02-brand/` | approximately 28 | 0 |
| Existing brand documents without explicit Requirement IDs | 34 / 34 | 0 |
| Canonical decided palettes with substantive usage rules | 0 / 2 | 2 / 2 |
| Undecided product palette material | 3 placeholder palette files | 9 argued directions across 3 proposal documents |
| Product identity README and examples documents | 0 / 8 | 8 / 8 |
| Comparative typography studies | 0 | 1 brand study plus a Design System boundary document |
| Argued forbidden-direction examples | 0 | 27 |
| Product identity defined by color alone | present as a risk | prohibited by the shared/adaptable identity matrix |
| Final-logo claims unsupported by approval | 1 | 0 |
| New open decisions | 0 | 1 (`OPEN-016`) |

The three legacy product `palette.md` files remain as `deprecated` navigation records pointing to the proposal documents. They are not canonical palettes.

## Requirement coverage

| State | After Phase 1 | After Phase 2 |
|---|---:|---:|
| conform | 81 | 88 |
| partial | 32 | 25 |
| absent | 6 | 6 |
| contradictory | 3 | 3 |
| total | 122 | 122 |

## Phase 2 decisions applied

- `Operational Editorial Modernism` is the fixed central direction.
- CMDR and Command palette values remain exact.
- Moss plus Ember is a shared brand signature, not a status system.
- Investigate, Govern and Studio each have three argued `proposed` directions.
- The final typography stack remains open.
- The textual wordmark `CMDR` is stable; final optical construction and any symbol remain open.
- Product identities share typography, geometry, component semantics, provenance and accessibility.
- Product differentiation may use rhythm, density, metadata treatment, visualization signatures and approved accent systems, but not color alone.
- AI is represented through provenance, run identity, uncertainty and control rather than magical or anthropomorphic decoration.

## Contradictions resolved in Phase 2

- generic or cloned brand bodies;
- placeholder palettes presented under canonical palette names;
- exact unapproved logo construction;
- semantic status colors mixed into the Command brand palette;
- product-specific typography documents implying separate type systems;
- ambiguous ownership between Brand and Design System foundations.

## Contradictions and work remaining

Outside the Phase 2 scope:

- generic versus Work Queue Saved Views;
- Work Queue screen-file consolidation;
- permission namespace normalization;
- detailed Design System tokens and component behavior;
- human approval of `OPEN-001` through `OPEN-004` and `OPEN-016`;
- final visual assets and visual-regression baselines.

## Outcome

Phase 2 is **PASS** for brand strategy, visual universe, product identities and decision-ready proposals. The repository remains **PARTIAL globally** because the Design System, UX, functional modules, journeys, screens, detailed objects, permissions, technical architecture and implementation are not complete.
