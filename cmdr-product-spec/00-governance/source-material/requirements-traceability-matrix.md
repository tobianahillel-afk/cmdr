---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-03
source-of-truth: source-material
requirements:
  - REQ-PROD-001
  - REQ-PROD-062
  - REQ-BRAND-001
  - REQ-BRAND-008
---
# Requirements Traceability Matrix

## Interpretation

This matrix preserves all 122 source Requirement IDs. Phase 2 adds detailed traceability for brand, product-identity and AI-expression requirements. `conform` means substantive coverage in the current owning phase, not completion of later UX, component, technical or implementation work.

## Coverage

| State | Phase 0 published | Reconciled before Phase 1 | After Phase 1 | After Phase 2 |
|---|---:|---:|---:|---:|
| conform | 0 | 0 | 81 | 88 |
| partial | 68 | 54 | 32 | 25 |
| absent | 47 | 61 | 6 | 6 |
| contradictory | 7 | 7 | 3 | 3 |
| total | 122 | 122 | 122 | 122 |

The Phase 0 published categories total 122 but differ from the row-level recount. Both baselines remain visible.

## Phase 2 detailed traceability

| Requirement | Source document | Canonical document / section | Before | After | Coverage | Contradiction resolved | Dependents | Proof | Open decision |
|---|---|---|---|---|---|---|---|---|---|
| REQ-BRAND-001 | `00-governance/source-material/brand-and-visual-decisions.md` | `02-brand/operational-editorial-modernism.md; visual-principles.md; brand-personality.md` — direction, principles and personality | partial | conform | Phase 2 substantive; specific rules, examples and acceptance criteria | — | Brand, Design System, product shells, marketing, documentation | direction, personality, principles and acceptance tests | — |
| REQ-BRAND-002 | `00-governance/source-material/brand-and-visual-decisions.md` | `02-brand/forbidden-directions.md` — anti-examples and CMDR alternatives | partial | conform | Phase 2 substantive; specific rules, examples and acceptance criteria | — | Brand, Design System, product shells, marketing, documentation | anti-examples with reason, risk and CMDR alternative | — |
| REQ-BRAND-003 | `00-governance/source-material/brand-and-visual-decisions.md` | `02-brand/cmdr/palette.md` — exact values, roles, contrast and theme behavior | partial | conform | Phase 2 substantive; specific rules, examples and acceptance criteria | — | Brand, Design System, product shells, marketing, documentation | exact CMDR values, roles, contrast and theme behavior | — |
| REQ-BRAND-004 | `00-governance/source-material/brand-and-visual-decisions.md` | `02-brand/command/palette.md` — exact values, surface hierarchy and restrictions | partial | conform | Phase 2 substantive; specific rules, examples and acceptance criteria | — | Brand, Design System, product shells, marketing, documentation | exact Command values, surface hierarchy and restrictions | — |
| REQ-BRAND-005 | `00-governance/source-material/brand-and-visual-decisions.md` | `02-brand/brand-essence-and-signature.md` — Moss + Ember signature and semantic separation | partial | conform | Phase 2 substantive; specific rules, examples and acceptance criteria | — | Brand, Design System, product shells, marketing, documentation | Moss + Ember signature separated from semantics | — |
| REQ-BRAND-006 | `00-governance/source-material/brand-and-visual-decisions.md` | `02-brand/investigate/palette-proposals.md; govern/palette-proposals.md; studio/palette-proposals.md` — three proposed directions per product; OPEN-001..003 remain open | partial | conform | Phase 2 substantive; specific rules, examples and acceptance criteria | resolved: placeholder/overclaim removed | Brand, Design System, product shells, marketing, documentation | three argued proposed directions per open product palette | OPEN-001, OPEN-002, OPEN-003 |
| REQ-BRAND-007 | `00-governance/source-material/brand-and-visual-decisions.md` | `02-brand/cmdr/typography.md; 03-design-system/foundations/typography.md` — candidate study; final stack remains OPEN-004 | partial | partial | Phase 2 study; decision open; specific rules, examples and acceptance criteria | — | Brand, Design System, product shells, marketing, documentation | candidate study, licensing evidence and validation matrix | OPEN-004 |
| REQ-BRAND-008 | `00-governance/source-material/brand-and-visual-decisions.md` | `02-brand/brand-architecture.md; command/README.md; investigate/README.md; govern/README.md; studio/README.md; logo-and-wordmark.md` — shared/adaptable matrix and product identities; final symbol remains OPEN-016 | partial | conform | Phase 2 substantive; specific rules, examples and acceptance criteria | resolved: placeholder/overclaim removed | Brand, Design System, product shells, marketing, documentation | common/adaptable matrix and four substantive product identities | OPEN-016 |
| REQ-PROD-003 | `00-governance/source-material/cmdr-master-product-brief.md` | `01-product-vision/product-vision.md; product-principles.md` — vision and product principles | conform | conform | unchanged Phase 1 evidence | — | Product owners and owning later phases | brand application references existing product principle or boundary | — |
| REQ-PROD-005 | `00-governance/source-material/cmdr-master-product-brief.md` | `01-product-vision/product-vision.md; product-principles.md` — vision and product principles | conform | conform | unchanged Phase 1 evidence | — | Product owners and owning later phases | brand application references existing product principle or boundary | — |
| REQ-PROD-007 | `00-governance/source-material/cmdr-master-product-brief.md` | `01-product-vision/product-vision.md; product-principles.md` — vision and product principles | partial | partial | unchanged; later owning phase | — | Product owners and owning later phases | brand application references existing product principle or boundary | — |
| REQ-PROD-008 | `00-governance/source-material/cmdr-master-product-brief.md` | `01-product-vision/product-vision.md; product-principles.md` — vision and product principles | partial | partial | unchanged; later owning phase | — | Product owners and owning later phases | brand application references existing product principle or boundary | — |
| REQ-PROD-013 | `00-governance/source-material/cmdr-master-product-brief.md` | `01-product-vision/product-boundaries.md` — product ownership and exclusions | conform | conform | unchanged Phase 1 evidence | — | Product owners and owning later phases | brand application references existing product principle or boundary | — |
| REQ-PROD-014 | `00-governance/source-material/cmdr-master-product-brief.md` | `01-product-vision/product-boundaries.md` — product ownership and exclusions | conform | conform | unchanged Phase 1 evidence | — | Product owners and owning later phases | brand application references existing product principle or boundary | — |
| REQ-PROD-015 | `00-governance/source-material/cmdr-master-product-brief.md` | `01-product-vision/product-boundaries.md` — product ownership and exclusions | conform | conform | unchanged Phase 1 evidence | — | Product owners and owning later phases | brand application references existing product principle or boundary | — |
| REQ-PROD-016 | `00-governance/source-material/cmdr-master-product-brief.md` | `01-product-vision/product-boundaries.md` — product ownership and exclusions | conform | conform | unchanged Phase 1 evidence | — | Product owners and owning later phases | brand application references existing product principle or boundary | — |
| REQ-PROD-048 | `00-governance/source-material/unresolved-decisions.md` | `00-governance/source-material/unresolved-decisions.md` — structured open or resolved decision | conform | conform | unchanged Phase 1 evidence | — | Product owners and owning later phases | named proposal set and structured open-decision evidence | OPEN-001 |
| REQ-PROD-049 | `00-governance/source-material/unresolved-decisions.md` | `00-governance/source-material/unresolved-decisions.md` — structured open or resolved decision | conform | conform | unchanged Phase 1 evidence | — | Product owners and owning later phases | named proposal set and structured open-decision evidence | OPEN-002 |
| REQ-PROD-050 | `00-governance/source-material/unresolved-decisions.md` | `00-governance/source-material/unresolved-decisions.md` — structured open or resolved decision | conform | conform | unchanged Phase 1 evidence | — | Product owners and owning later phases | named proposal set and structured open-decision evidence | OPEN-003 |
| REQ-PROD-051 | `00-governance/source-material/unresolved-decisions.md` | `00-governance/source-material/unresolved-decisions.md` — structured open or resolved decision | conform | conform | unchanged Phase 1 evidence | — | Product owners and owning later phases | typography candidate study; final decision intentionally open | OPEN-004 |
| REQ-AI-001 | `00-governance/source-material/ai-and-automation-constraints.md` | `01-product-vision/product-principles.md; product-boundaries.md` — AI and automation boundaries | conform | conform | unchanged Phase 1 evidence | — | Studio, Command, Investigate, Govern, Security | AI visual language distinguishes provenance, proposal and authority | — |
| REQ-AI-002 | `00-governance/source-material/ai-and-automation-constraints.md` | `01-product-vision/product-principles.md; product-boundaries.md` — AI and automation boundaries | conform | conform | unchanged Phase 1 evidence | — | Studio, Command, Investigate, Govern, Security | AI visual language distinguishes provenance, proposal and authority | — |
| REQ-AI-003 | `00-governance/source-material/ai-and-automation-constraints.md` | `01-product-vision/product-principles.md; product-boundaries.md` — AI and automation boundaries | conform | conform | unchanged Phase 1 evidence | — | Studio, Command, Investigate, Govern, Security | AI visual language distinguishes provenance, proposal and authority | — |
| REQ-AI-007 | `00-governance/source-material/ai-and-automation-constraints.md` | `01-product-vision/product-principles.md; product-boundaries.md` — AI and automation boundaries | partial | partial | unchanged; later owning phase | — | Studio, Command, Investigate, Govern, Security | AI visual language distinguishes provenance, proposal and authority | — |
| REQ-AI-010 | `00-governance/source-material/ai-and-automation-constraints.md` | `01-product-vision/product-principles.md; product-boundaries.md` — AI and automation boundaries | partial | partial | unchanged; later owning phase | — | Studio, Command, Investigate, Govern, Security | AI visual language distinguishes provenance, proposal and authority | — |

## Complete Requirement ID inventory

Each ID appears exactly once below with its Phase 2 state.

### conform

`REQ-PROD-001`, `REQ-PROD-002`, `REQ-PROD-003`, `REQ-PROD-004`, `REQ-PROD-005`, `REQ-PROD-009`, `REQ-PROD-011`, `REQ-PROD-012`, `REQ-PROD-013`, `REQ-PROD-014`, `REQ-PROD-015`, `REQ-PROD-016`, `REQ-PROD-017`, `REQ-PROD-018`, `REQ-PROD-019`, `REQ-PROD-021`, `REQ-PROD-022`, `REQ-PROD-023`, `REQ-PROD-024`, `REQ-PROD-025`, `REQ-PROD-026`, `REQ-PROD-027`, `REQ-PROD-028`, `REQ-PROD-029`, `REQ-PROD-030`, `REQ-PROD-031`, `REQ-PROD-032`, `REQ-PROD-033`, `REQ-PROD-034`, `REQ-PROD-035`, `REQ-PROD-036`, `REQ-PROD-037`, `REQ-PROD-038`, `REQ-PROD-039`, `REQ-PROD-040`, `REQ-PROD-041`, `REQ-PROD-042`, `REQ-PROD-043`, `REQ-PROD-044`, `REQ-PROD-045`, `REQ-PROD-046`, `REQ-PROD-047`, `REQ-PROD-048`, `REQ-PROD-049`, `REQ-PROD-050`, `REQ-PROD-051`, `REQ-PROD-052`, `REQ-PROD-053`, `REQ-PROD-054`, `REQ-PROD-055`, `REQ-PROD-056`, `REQ-PROD-057`, `REQ-PROD-058`, `REQ-PROD-059`, `REQ-PROD-060`, `REQ-PROD-061`, `REQ-PROD-062`, `REQ-AI-001`, `REQ-AI-002`, `REQ-AI-003`, `REQ-AI-004`, `REQ-AI-005`, `REQ-AI-006`, `REQ-AI-011`, `REQ-OBJ-001`, `REQ-OBJ-002`, `REQ-OBJ-003`, `REQ-OBJ-004`, `REQ-OBJ-005`, `REQ-OBJ-006`, `REQ-OBJ-007`, `REQ-OBJ-008`, `REQ-OBJ-010`, `REQ-BRAND-001`, `REQ-BRAND-002`, `REQ-BRAND-003`, `REQ-BRAND-004`, `REQ-BRAND-005`, `REQ-BRAND-006`, `REQ-BRAND-008`, `REQ-SEC-001`, `REQ-SEC-002`, `REQ-INV-001`, `REQ-INV-002`, `REQ-INV-003`, `REQ-INV-004`, `REQ-INV-005`, `REQ-INV-006`

### partial

`REQ-PROD-006`, `REQ-PROD-007`, `REQ-PROD-008`, `REQ-PROD-010`, `REQ-PROD-020`, `REQ-AI-007`, `REQ-AI-008`, `REQ-AI-009`, `REQ-AI-010`, `REQ-OBJ-009`, `REQ-UX-001`, `REQ-UX-002`, `REQ-UX-003`, `REQ-UX-006`, `REQ-UX-007`, `REQ-UX-009`, `REQ-BRAND-007`, `REQ-SEC-003`, `REQ-SEC-004`, `REQ-SEC-005`, `REQ-JRN-001`, `REQ-JRN-003`, `REQ-JRN-004`, `REQ-JRN-006`, `REQ-JRN-007`

### absent

`REQ-UX-004`, `REQ-UX-005`, `REQ-UX-010`, `REQ-JRN-002`, `REQ-JRN-005`, `REQ-JRN-008`

### contradictory

`REQ-OBJ-011`, `REQ-OBJ-012`, `REQ-UX-008`

## Phase 2 movement

Seven requirements moved from `partial` to `conform`: `REQ-BRAND-001` through `REQ-BRAND-006`, and `REQ-BRAND-008`.

`REQ-BRAND-007` remains `partial` because the final type stack is still `OPEN-004`. `REQ-PROD-048`, `REQ-PROD-049` and `REQ-PROD-050` remain traceably covered by open decisions and proposal documents; no palette was approved.

## Remaining contradictions

- `REQ-OBJ-011` and `REQ-OBJ-012`: generic Saved Views versus Work Queue Saved Views.
- `REQ-UX-008`: Work Queue variants remain separate screen files.
- Permission namespace conflicts remain assigned to their later owning phase.

## Update rule

A later change updates the source decision, owning canonical document, dependents, proof and open-decision record. A proposal never becomes a decision merely because its individual contrast pairs pass.
