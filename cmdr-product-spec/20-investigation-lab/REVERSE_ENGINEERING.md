# Reverse Engineering

## Objective

Provide disassembly, decompilation and binary-navigation workflows for expert analysis.

## Scope

This specification owns the page-local behaviour of **Reverse Engineering** in **Investigation Lab**. Shared object definitions, states, permissions, navigation and design rules remain canonical in the linked shared documents and are not redefined here.

## Functional owner

DFIR Lead

## Affected objects

- Evidence
- Binary image
- Function
- Symbol
- Cross-reference
- Annotation
- Finding

## Features

- Disassembly and decompiler views
- Function and symbol browser
- Strings and cross-references
- Control-flow graph
- Hex and memory layout
- Type and symbol annotation
- Patch-diff or version comparison
- Export of analyst notes and derived artefacts

## UX and interactions

- Views remain synchronised on address/function
- Analyst annotations are versioned
- Navigation history supports back/forward
- Unsafe binary content is never rendered as active content
- Large binaries load progressively

## Permissions

`reverse.execute`; collaboration requires `case.manage`; exporting binaries or patches requires `evidence.export`.

The permission semantics are governed by [`PERMISSION_MODEL.md`](../02-domain-model/PERMISSION_MODEL.md).

## States

Analysis session has loading/ready/failed states; Findings remain canonical.

Canonical state names and transition constraints are governed by [`STATE_MODELS.md`](../02-domain-model/STATE_MODELS.md).

## Dependencies

Reverse engine, symbol services, evidence, case notes.

## Acceptance criteria

- Addresses resolve consistently across views
- Tool version and loader settings are stored
- Annotations survive session reopening
- Derived outputs retain sample provenance
- No edit mutates original evidence

## Open questions

- Which architectures and binary formats are mandatory?
- Will collaborative annotation be real-time?
