# Evidence Board

## Objective

Organise, verify and relate evidence while preserving integrity and provenance.

## Scope

This specification owns the page-local behaviour of **Evidence Board** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Evidence
- Entity
- Case
- Finding
- Hypothesis
- Integrity record

## Features

- Card and table views
- Evidence registration and collection status
- Integrity verification
- Relationship graph
- Annotations and tags
- Parent/derived artefact links
- Evidence package creation

## UX and interactions

- Canvas has a complete list alternative
- Connections are explicit relationships, not decorative lines
- Original artefacts are read-only
- Verification shows hash, collector and method
- Large boards support grouping and focus

## Permissions

`evidence.read`; collection requires `evidence.collect`; verification requires `evidence.verify`; export requires `evidence.export`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Uses canonical Evidence states; visual groupings are not states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Evidence storage, hashing, case service, analysis tools.

## Acceptance criteria

- Every evidence item has provenance
- Verification is independently auditable
- Derived artefacts retain parent links
- Rejected evidence remains visible with reason
- Canvas and list views represent the same objects

## Open questions

- Which evidence formats receive inline preview?
- What chain-of-custody fields are mandatory by jurisdiction?
