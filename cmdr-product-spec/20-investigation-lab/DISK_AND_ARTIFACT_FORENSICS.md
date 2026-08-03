# Disk and Artifact Forensics

## Objective

Analyse disk images and endpoint artefacts while preserving forensic integrity.

## Scope

This specification owns the page-local behaviour of **Disk and Artifact Forensics** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Evidence
- Disk image
- Filesystem artefact
- Registry artefact
- Browser artefact
- Timeline entry
- Finding

## Features

- Filesystem and partition browsing
- Deleted-file and metadata analysis
- Registry, event log, browser and execution artefacts
- Hash sets and known-file filtering
- Super-timeline generation
- Carving and derived artefact extraction
- Cross-source correlation

## UX and interactions

- Original images are mounted read-only
- Timezone and filesystem semantics are explicit
- Recovered artefacts show confidence and method
- Large images support indexed navigation
- Preview and export are independently permissioned

## Permissions

`analysis.execute`; artefact export requires `evidence.export`; sensitive user data is ABAC-scoped.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Analysis jobs and Evidence use canonical states; filesystem allocation state is an artefact attribute.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Disk forensic engine, evidence storage, hash-set service, timeline.

## Acceptance criteria

- Image integrity is checked before analysis
- Mount mode is demonstrably read-only
- Derived artefacts preserve offsets and method
- Super-timeline entries link to source artefacts
- Partial parsing is explicit

## Open questions

- Which filesystem and artefact parsers are mandatory?
- How are privacy minimisation rules applied?
