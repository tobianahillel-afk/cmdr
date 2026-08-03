---
id: foundation-elevation
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-001
  - REQ-BRAND-002
---
# Elevation foundation — source boundaries

## Responsibility

Phase 3 owns exact elevation tokens. Brand requires a surface system based primarily on editorial grouping, borders and tonal contrast.

## Canonical source

- [`../../02-brand/shape-and-surfaces.md`](../../02-brand/shape-and-surfaces.md)

## Phase 3 responsibility

Define elevation only for:

- menus and popovers;
- drawers and modals;
- drag previews;
- temporary overlays;
- rare focus-preserving transitions.

## Constraints

- Persistent page sections do not float by default.
- Heavy shadows, glow and glassmorphism are prohibited.
- An overlay must have a functional stacking reason and an accessible focus model.
- Dark themes use tonal separation and borders before luminous shadows.

## Acceptance criterion

**Given** a normal page section and a temporary confirmation dialog,  
**When** elevation is applied,  
**Then** only the dialog receives overlay elevation and the underlying section remains structurally grouped without decorative shadow.
