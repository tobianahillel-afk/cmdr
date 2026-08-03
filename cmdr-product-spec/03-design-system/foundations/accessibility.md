---
id: foundation-accessibility
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-003
  - REQ-BRAND-004
  - REQ-BRAND-006
  - REQ-BRAND-007
---
# Accessibility foundation — brand dependencies

## Responsibility

The Design System owns measurable interaction and component accessibility. Brand owns the rule that recognition, hierarchy and identity must survive without color, motion, imagery or a specific font.

## Canonical source

- [`../../02-brand/accessibility-and-brand.md`](../../02-brand/accessibility-and-brand.md)

## Phase 3 responsibility

Validate:

- WCAG 2.2 AA contrast for text and interactive states;
- keyboard navigation and focus;
- zoom, reflow and localization;
- reduced motion;
- forced-colors and high-contrast modes;
- chart and canvas alternatives;
- screen-reader labels and reading order;
- light and dark themes.

## Constraints

- Proposed palettes are not approved merely because selected pairs pass contrast.
- Color contrast is checked in context, including text size and adjacent states.
- Product identity remains recognizable through labels, structure and motifs.
- Critical information is never encoded by color, icon or motion alone.

## Acceptance criterion

**Given** a user in forced-colors mode,  
**When** they move from Command to Investigate,  
**Then** product name, context, selected state and ownership remain understandable even when brand colors are replaced by the operating system.
