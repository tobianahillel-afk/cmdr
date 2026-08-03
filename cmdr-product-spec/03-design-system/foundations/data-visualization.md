---
id: foundation-data-visualization
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-001
  - REQ-BRAND-002
  - REQ-BRAND-008
---
# Data visualization foundation — source boundaries

## Responsibility

Phase 3 will define chart tokens, series behavior and accessible interaction. Brand defines the visual language: evidence-led, quiet, annotated and never decorative.

## Canonical source

- [`../../02-brand/data-visualization-language.md`](../../02-brand/data-visualization-language.md)

## Phase 3 responsibility

Define:

- qualitative and sequential series tokens;
- axes, grids, labels and annotations;
- confidence and uncertainty treatment;
- selection, comparison and hover;
- print and export behavior;
- table alternatives;
- color-vision-deficiency validation.

## Constraints

- A chart must answer a stated operational question.
- Source, time range, freshness and units remain visible.
- 3D charts, decorative maps, gauges resembling game counters and saturated rainbow scales are prohibited.
- Product accents do not replace semantic or series systems.
- Every essential visualization has an accessible textual or tabular alternative.

## Acceptance criterion

**Given** a trend that supports a Decision,  
**When** Govern displays it,  
**Then** the chart identifies source, interval, units and uncertainty, and the same values are available without relying on color or pointer interaction.
