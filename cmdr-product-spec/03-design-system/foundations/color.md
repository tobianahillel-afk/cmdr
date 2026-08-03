---
id: foundation-color
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-003
  - REQ-BRAND-004
  - REQ-BRAND-005
  - REQ-BRAND-006
---
# Color foundation — source boundaries

## Responsibility

This foundation translates approved brand palettes into semantic and component tokens during Phase 3. It does not own brand values and must not copy product palette definitions into a second canonical table.

## Canonical sources

- CMDR parent palette: [`../../02-brand/cmdr/palette.md`](../../02-brand/cmdr/palette.md)
- Command palette: [`../../02-brand/command/palette.md`](../../02-brand/command/palette.md)
- Investigate proposals: [`../../02-brand/investigate/palette-proposals.md`](../../02-brand/investigate/palette-proposals.md)
- Govern proposals: [`../../02-brand/govern/palette-proposals.md`](../../02-brand/govern/palette-proposals.md)
- Studio proposals: [`../../02-brand/studio/palette-proposals.md`](../../02-brand/studio/palette-proposals.md)

Only approved palettes may become product tokens. Proposed values remain evaluation material and must not be shipped as defaults.

## Phase 3 responsibility

Phase 3 will define:

- brand-to-product-to-semantic token mapping;
- light and dark theme aliases;
- focus, selection and interaction tokens;
- status, severity, confidence and data-visualization semantics;
- contrast validation and fallback behavior.

## Non-negotiable constraints

- Color never carries meaning alone.
- CMDR Moss and Ember remain brand accents, not automatic success or critical colors.
- Product identity and status semantics are separate layers.
- Absolute black is not the automatic dark-theme canvas.
- A coherent interface around an external engine does not authorize vendor colors to replace CMDR hierarchy.

## Acceptance criterion

**Given** an Investigate proposal marked `proposed`,  
**When** Phase 3 creates product tokens,  
**Then** no token uses that proposal as the default until `OPEN-001` is explicitly resolved and the approved decision is recorded.
