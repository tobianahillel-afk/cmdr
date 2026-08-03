---
id: source-material-readme
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: source-material
requirements:
  - REQ-PROD-006
  - REQ-PROD-009
  - REQ-PROD-012
---
# Source Material

## Purpose

This directory preserves sponsor-provided decisions, constraints, preferences and unresolved questions. It records intent and evidence. It does not replace the canonical specification owned by the relevant product or domain.

## Precedence

1. Active explicit decisions and source requirements in this directory.
2. Approved ADRs.
3. Validated canonical documents owned by the affected domain.
4. Canonical registers.
5. Documents in review.
6. Draft documents.
7. Proposals and assumptions.
8. Archived documents, which are never normative.

The detailed conflict-resolution process is defined in [`../source-of-truth-policy.md`](../source-of-truth-policy.md).

## Use

- Canonical documents cite applicable Requirement IDs instead of copying the source text.
- A source requirement states intent; the owning canonical document defines its normative application.
- A proposal remains a proposal until an explicit decision replaces it.
- An unresolved decision remains in `unresolved-decisions.md`.
- Consumers link to the canonical owner and describe local use only.

## Documents and order

1. `cmdr-master-product-brief.md`
2. `product-boundaries.md`
3. `canonical-object-and-ownership-decisions.md`
4. `ai-and-automation-constraints.md`
5. `native-capability-strategy.md`
6. `ux-and-navigation-decisions.md`
7. `brand-and-visual-decisions.md`
8. `product-capability-inventory.md`
9. `user-role-and-journey-inventory.md`
10. `explicit-non-goals.md`
11. `unresolved-decisions.md`
12. `requirements-traceability-matrix.md`
13. `qualitative-baseline.md`

## Phase boundary

Phase 1 propagates source decisions into governance and product vision. Detailed brand, UX, screens, object schemas, permissions, protocols and implementation remain assigned to later phases.
