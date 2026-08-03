---
id: foundation-borders-and-radius
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-001
  - REQ-BRAND-008
---
# Borders and radius foundation — source boundaries

## Responsibility

Phase 3 owns exact radius and border tokens. Brand defines the intended geometry: edited, stable and restrained rather than a collection of floating rounded cards.

## Canonical source

- [`../../02-brand/shape-and-surfaces.md`](../../02-brand/shape-and-surfaces.md)

## Phase 3 responsibility

Define:

- radius scale;
- border widths and contrast;
- separators;
- selected and focused outlines;
- overlay boundaries;
- high-contrast-mode behavior.

## Constraints

- Small or moderate radii are the default direction.
- Fully pill-shaped containers are reserved for genuinely compact labels, filters or segmented controls.
- Large decorative capsules, excessive rounding and glass layers are prohibited.
- Borders and tonal surfaces communicate structure before shadows.

## Acceptance criterion

**Given** a dense evidence table and an adjacent Inspector,  
**When** their boundaries are rendered,  
**Then** the hierarchy remains clear through alignment, separators and tonal contrast without wrapping each region in a large floating card.
