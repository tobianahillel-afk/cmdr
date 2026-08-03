# Case Timeline

## Objective

Reconstruct events, assertions and actions into a defensible chronology.

## Scope

This specification owns the page-local behaviour of **Case Timeline** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Case
- Timeline entry
- Event
- Evidence
- Entity
- Finding
- Decision
- Run

## Features

- Multi-source chronological reconstruction
- Observed, inferred and user-asserted entry types
- Gap and clock-skew handling
- Entity and evidence filters
- Timeline annotations
- Export and snapshot

## UX and interactions

- Drag/reorder is prohibited for observed events
- Analysts may place inferred events with explicit uncertainty
- Every entry expands to provenance
- Clock corrections are visible and reversible

## Permissions

`case.read`; annotations require `case.manage`; export requires `evidence.export` or `report.manage` based on content.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Timeline entries do not own source-object states; inferred entries have draft/review/accepted status local to the timeline feature.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Event search, evidence registry, time normalisation, audit.

## Acceptance criteria

- Observed ordering is reproducible
- Clock-skew adjustments preserve originals
- Inferences are visually distinct
- Export includes provenance and timezone

## Open questions

- How are conflicting timestamps adjudicated?
- Should collaborative timeline branches be supported?
