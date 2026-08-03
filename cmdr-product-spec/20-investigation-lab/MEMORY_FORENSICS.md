# Memory Forensics

## Objective

Analyse memory captures for processes, credentials exposure, injections, network state and kernel artefacts.

## Scope

This specification owns the page-local behaviour of **Memory Forensics** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Evidence
- Memory image
- Process/thread entity
- Module
- Network entity
- Derived artefact
- Finding

## Features

- Profile/symbol validation
- Process and tree analysis
- Handles, modules and injected regions
- Network and socket reconstruction
- Credential-exposure indicators without displaying secrets unnecessarily
- Kernel and driver analysis
- Plugin execution and result comparison

## UX and interactions

- Analysts can compare plugin outputs
- Sensitive secret-like values are masked by default
- Every result links to image offset or source plugin
- Long jobs expose progress and cancellation
- Derived dumps require explicit gated export

## Permissions

`analysis.execute`; sensitive result access and dumps are ABAC-scoped; export requires `evidence.export`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Analysis jobs use canonical job states; memory Evidence follows canonical Evidence states.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Memory analysis engine, symbol services, isolated workers, evidence storage.

## Acceptance criteria

- Image hash and acquisition metadata are verified
- Plugin versions and parameters are stored
- Masked values stay protected in exports unless authorised
- Derived dumps retain offsets and provenance
- Failed plugins do not invalidate successful results

## Open questions

- Which memory frameworks are first-class?
- What secret-redaction policy is required?
