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
---
# Requirements Traceability Matrix

## Interpretation

Source for all rows: Master Execution Prompt as captured in `source-material/`. `conform` means substantive Phase 1 product/governance coverage and traceability, not completion of later UX, functional, security, technical or implementation work.

## Coverage

| State | Phase 0 published | Reconciled before | After Phase 1 |
|---|---:|---:|---:|
| conform | 0 | 0 | 81 |
| partial | 68 | 54 | 32 |
| absent | 47 | 61 | 6 |
| contradictory | 7 | 7 | 3 |
| total | 122 | 122 | 122 |

The published Phase 0 summary did not arithmetically match its 122 catalog rows. The published values and reconciled recount are both retained.

## Coverage rules

- `conform`: Phase 1 intent is substantive, traced and has no active Phase 1 contradiction.
- `partial`: product principle or boundary exists; its owning later phase must define detailed behavior.
- `absent`: deliberately deferred with no delivery claim.
- `contradictory`: an active competing structure remains.

## Dependents by family

- `REQ-PROD-*`: all product owners and later product/UX phases.
- `REQ-AI-*`: CMDR Studio, Command, Investigate, Govern, Settings and Security.
- `REQ-OBJ-*`: ownership register, domain model, product projections and permissions.
- `REQ-UX-*`: Experience Architecture, Design System, journeys and screens.
- `REQ-BRAND-*`: Brand and Design System Phase 2.
- `REQ-SEC-*`: Security/Trust, Govern, Endpoint Agent and dangerous-action UX.
- `REQ-JRN-*`: journeys, transitions and pilot screens.
- `REQ-INV-*`: Investigate modules and capability classification.

## Requirement rows

| Requirement | Before | After | Canonical evidence / section | Contradiction | Open decision |
|---|---|---|---|---|---|
| REQ-PROD-001 | partial | conform | 01-product-vision/product-principles.md + product-vision.md | — | — |
| REQ-PROD-002 | partial | conform | 01-product-vision/product-principles.md + product-vision.md | — | — |
| REQ-PROD-003 | partial | conform | 01-product-vision/product-principles.md + product-vision.md | — | — |
| REQ-PROD-004 | partial | conform | 01-product-vision/product-principles.md + product-vision.md | — | — |
| REQ-PROD-005 | partial | conform | 01-product-vision/product-principles.md + product-vision.md | — | — |
| REQ-PROD-006 | contradictory | partial | 01-product-vision/product-principles.md + product-vision.md | resolved | — |
| REQ-PROD-007 | partial | partial | 01-product-vision/product-principles.md + product-vision.md | — | — |
| REQ-PROD-008 | partial | partial | 01-product-vision/product-principles.md + product-vision.md | — | — |
| REQ-PROD-009 | partial | conform | 01-product-vision/product-principles.md + product-vision.md | — | — |
| REQ-PROD-010 | absent | partial | 01-product-vision/product-principles.md + product-vision.md | — | — |
| REQ-PROD-011 | partial | conform | 01-product-vision/product-principles.md + product-vision.md | — | — |
| REQ-PROD-012 | contradictory | conform | 01-product-vision/product-principles.md + product-vision.md | resolved | — |
| REQ-AI-001 | partial | conform | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-AI-002 | partial | conform | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-AI-003 | partial | conform | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-AI-004 | partial | conform | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-AI-005 | partial | conform | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-AI-006 | partial | conform | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-AI-007 | partial | partial | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-AI-008 | absent | partial | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-AI-009 | absent | partial | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-AI-010 | partial | partial | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-AI-011 | partial | conform | 01-product-vision/product-principles.md + product-boundaries.md + product README boundaries | — | — |
| REQ-OBJ-001 | partial | conform | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | — | — |
| REQ-OBJ-002 | partial | conform | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | — | — |
| REQ-OBJ-003 | partial | conform | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | — | — |
| REQ-OBJ-004 | partial | conform | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | — | — |
| REQ-OBJ-005 | partial | conform | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | — | — |
| REQ-OBJ-006 | partial | conform | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | — | — |
| REQ-OBJ-007 | partial | conform | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | — | — |
| REQ-OBJ-008 | partial | conform | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | — | — |
| REQ-OBJ-009 | contradictory | partial | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | resolved | — |
| REQ-OBJ-010 | partial | conform | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | — | — |
| REQ-OBJ-011 | contradictory | contradictory | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | active | — |
| REQ-OBJ-012 | contradictory | contradictory | 00-governance/ownership-register.md + 01-product-vision/product-boundaries.md | active | — |
| REQ-UX-001 | partial | partial | 00-governance/adr/ADR-0005-page-view-mode-filter-rules.md + Phase 3 specifications | — | — |
| REQ-UX-002 | partial | partial | 00-governance/adr/ADR-0005-page-view-mode-filter-rules.md + Phase 3 specifications | — | — |
| REQ-UX-003 | partial | partial | 00-governance/adr/ADR-0005-page-view-mode-filter-rules.md + Phase 3 specifications | — | — |
| REQ-UX-004 | absent | absent | 00-governance/adr/ADR-0005-page-view-mode-filter-rules.md + Phase 3 specifications | — | — |
| REQ-UX-005 | absent | absent | 00-governance/adr/ADR-0005-page-view-mode-filter-rules.md + Phase 3 specifications | — | — |
| REQ-UX-006 | partial | partial | 00-governance/adr/ADR-0005-page-view-mode-filter-rules.md + Phase 3 specifications | — | — |
| REQ-UX-007 | absent | partial | 00-governance/adr/ADR-0005-page-view-mode-filter-rules.md + Phase 3 specifications | — | — |
| REQ-UX-008 | contradictory | contradictory | 00-governance/adr/ADR-0005-page-view-mode-filter-rules.md + Phase 3 specifications | active | — |
| REQ-UX-009 | partial | partial | 00-governance/adr/ADR-0005-page-view-mode-filter-rules.md + Phase 3 specifications | — | — |
| REQ-UX-010 | absent | absent | 00-governance/adr/ADR-0005-page-view-mode-filter-rules.md + Phase 3 specifications | — | — |
| REQ-BRAND-001 | partial | partial | 00-governance/source-material/brand-and-visual-decisions.md + Phase 2 | — | — |
| REQ-BRAND-002 | partial | partial | 00-governance/source-material/brand-and-visual-decisions.md + Phase 2 | — | — |
| REQ-BRAND-003 | partial | partial | 00-governance/source-material/brand-and-visual-decisions.md + Phase 2 | — | — |
| REQ-BRAND-004 | partial | partial | 00-governance/source-material/brand-and-visual-decisions.md + Phase 2 | — | — |
| REQ-BRAND-005 | partial | partial | 00-governance/source-material/brand-and-visual-decisions.md + Phase 2 | — | — |
| REQ-BRAND-006 | partial | partial | 00-governance/source-material/brand-and-visual-decisions.md + Phase 2 | — | — |
| REQ-BRAND-007 | partial | partial | 00-governance/source-material/brand-and-visual-decisions.md + Phase 2 | — | — |
| REQ-BRAND-008 | partial | partial | 00-governance/source-material/brand-and-visual-decisions.md + Phase 2 | — | — |
| REQ-PROD-013 | partial | conform | 01-product-vision/product-boundaries.md | — | — |
| REQ-PROD-014 | partial | conform | 01-product-vision/product-boundaries.md | — | — |
| REQ-PROD-015 | partial | conform | 01-product-vision/product-boundaries.md | — | — |
| REQ-PROD-016 | partial | conform | 01-product-vision/product-boundaries.md | — | — |
| REQ-PROD-017 | partial | conform | 01-product-vision/product-boundaries.md | — | — |
| REQ-PROD-018 | contradictory | conform | 01-product-vision/product-boundaries.md | resolved | — |
| REQ-PROD-019 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-020 | partial | partial | product-principles.md + capability-map.md | — | — |
| REQ-SEC-001 | absent | conform | 01-product-vision/operating-model.md + later trust specifications | — | — |
| REQ-SEC-002 | partial | conform | 01-product-vision/operating-model.md + later trust specifications | — | — |
| REQ-SEC-003 | absent | partial | 01-product-vision/operating-model.md + later trust specifications | — | — |
| REQ-SEC-004 | partial | partial | 01-product-vision/operating-model.md + later trust specifications | — | — |
| REQ-SEC-005 | partial | partial | 01-product-vision/operating-model.md + later trust specifications | — | — |
| REQ-PROD-021 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-022 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-023 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-024 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-025 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-026 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-027 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-028 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-029 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-030 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-031 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-032 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-033 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-034 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-035 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-PROD-036 | absent | conform | 01-product-vision/target-users.md | — | — |
| REQ-JRN-001 | partial | partial | 01-product-vision/operating-model.md + Phase 5 journeys | — | — |
| REQ-JRN-002 | absent | absent | 01-product-vision/operating-model.md + Phase 5 journeys | — | — |
| REQ-JRN-003 | partial | partial | 01-product-vision/operating-model.md + Phase 5 journeys | — | — |
| REQ-JRN-004 | partial | partial | 01-product-vision/operating-model.md + Phase 5 journeys | — | — |
| REQ-JRN-005 | absent | absent | 01-product-vision/operating-model.md + Phase 5 journeys | — | — |
| REQ-JRN-006 | partial | partial | 01-product-vision/operating-model.md + Phase 5 journeys | — | — |
| REQ-JRN-007 | partial | partial | 01-product-vision/operating-model.md + Phase 5 journeys | — | — |
| REQ-JRN-008 | absent | absent | 01-product-vision/operating-model.md + Phase 5 journeys | — | — |
| REQ-PROD-037 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-038 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-039 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-040 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-041 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-042 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-INV-001 | absent | conform | 01-product-vision/capability-map.md + product-boundaries.md | — | — |
| REQ-INV-002 | absent | conform | 01-product-vision/capability-map.md + product-boundaries.md | — | — |
| REQ-INV-003 | absent | conform | 01-product-vision/capability-map.md + product-boundaries.md | — | — |
| REQ-INV-004 | absent | conform | 01-product-vision/capability-map.md + product-boundaries.md | — | — |
| REQ-INV-005 | absent | conform | 01-product-vision/capability-map.md + product-boundaries.md | — | — |
| REQ-INV-006 | absent | conform | 01-product-vision/capability-map.md + product-boundaries.md | — | — |
| REQ-PROD-043 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-044 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-045 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-046 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-047 | absent | conform | 01-product-vision/capability-map.md | — | — |
| REQ-PROD-048 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-001 |
| REQ-PROD-049 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-002 |
| REQ-PROD-050 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-003 |
| REQ-PROD-051 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-004 |
| REQ-PROD-052 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-005 |
| REQ-PROD-053 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-006 |
| REQ-PROD-054 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-007 |
| REQ-PROD-055 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-008 |
| REQ-PROD-056 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-009 (resolved Phase 1) |
| REQ-PROD-057 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-010 |
| REQ-PROD-058 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-011 |
| REQ-PROD-059 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-012 |
| REQ-PROD-060 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-013 |
| REQ-PROD-061 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-014 |
| REQ-PROD-062 | absent | conform | 00-governance/source-material/unresolved-decisions.md | — | OPEN-015 |

## Resolved contradictions

Source precedence, product ownership, optional AI, Studio placement, Endpoint target-versus-delivery wording, delivery taxonomy and capability review authority are resolved at Phase 1 level. OPEN-009 is recorded as resolved.

## Remaining contradictions

- `REQ-OBJ-011` and `REQ-OBJ-012`: generic versus Work Queue Saved Views in existing consumers.
- `REQ-UX-008`: Work Queue variants remain separate screen files until Phase 3/6.
- Permission namespace conflicts remain assigned to Phase 7.

## Proof and update rule

The canonical evidence column is the Phase 1 proof. Later changes require the owning canonical section, dependent files, validation evidence, and an open-decision record when applicable.
