---
id: foundation-typography
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-007
  - REQ-PROD-051
---
# Typography foundation — source boundaries

## Responsibility

This foundation will implement the typographic system after the brand study and `OPEN-004` are reviewed. The brand domain defines the desired editorial character; the Design System defines exact sizes, line heights, weights, responsive behavior and tokens.

## Canonical source

- Brand typography study: [`../../02-brand/cmdr/typography.md`](../../02-brand/cmdr/typography.md)

## Pending decision

`OPEN-004` remains open. Inter, Inter Tight, IBM Plex Sans and equivalent editorial sans combinations are candidates, not approved defaults. Monospace selection is also subject to readability, licensing, language and platform testing.

## Phase 3 responsibility

Phase 3 will define role tokens for Display, Page title, Section title, Body, UI label, Metadata, Table, Badge, Code, Query and Identifier, together with:

- exact metrics;
- tabular-number behavior;
- truncation and wrapping;
- dense-interface legibility;
- fallback stacks;
- localization expansion;
- Windows, macOS and Linux rendering tests.

## Constraints

- One coherent family system serves all products.
- Product differentiation does not use unrelated typefaces.
- Monospace is reserved for technical material and identifiers.
- Typography carries hierarchy before color or decorative containers.
- No font binaries are stored in this documentation repository.

## Acceptance criterion

**Given** `OPEN-004` is still open,  
**When** a component specification references typography,  
**Then** it uses semantic role names rather than asserting a final font family or unsupported metric.
