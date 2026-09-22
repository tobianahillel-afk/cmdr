# CMDR product coverage graph

This directory contains engineering artifacts derived from the canonical Product Spec.

The Product Spec remains authoritative. The coverage graph is generated from source identities and references and must never be used to redefine product semantics.

## E0-COVERAGE-001A

`cmdr-dev coverage-graph` builds the typed product entity/reference graph on demand.

Owned entities retain:
- source path;
- source SHA-256;
- document status;
- domain/product/module metadata when present;
- source-of-truth classification.

References that do not own a document in the current inventory remain explicit `reference_only: true` entities. This is deliberate: the graph records that a source referenced the identifier without inventing ownership or content.

Edge types in this first sublot are restricted to facts directly visible in source metadata or references:
- `requirement`;
- `open-decision`;
- `permission`;
- `reference`.

Implementation obligations are **not** inferred in this sublot; that is E0-COVERAGE-001B.


## E0-COVERAGE-001B

`cmdr-dev obligations` derives the implementation obligation set from active normative evidence.

An identifier becomes an active obligation only when its family-specific source of authority says it is active:
- capabilities: active Capability Register shards + canonical capability source;
- requirements: the ten core sponsor/source requirement documents listed by Source Material governance;
- screens: the active section of the Screen Register + canonical screen source;
- permissions: Permission Register rows;
- implementation contracts and canonical objects: active canonical source documents.

Reference-only or historical identifiers remain traceable in `unresolved_references` but do not silently become delivery obligations.

The detailed obligation file is generated on demand while corpus compilation remains cheap.


### Source requirement identity baseline

The initial Phase-0 Requirement Catalog is the identity authority for the 122 source Requirement IDs. The compact engineering snapshot is `source-requirements-baseline.json`.

Its provenance is pinned to Product Spec commit `0d4c6183daa5d7a99b51718a99250a2dfcb4ecc6`, source file `00-governance/source-material/requirements-traceability-matrix.md`, and source blob `652c7e2300664a70a9257d4fbb6f971b210a5a8e`.

The snapshot contains identities only, never requirement semantics. Current Product Spec documents remain authoritative for meaning and evidence. This prevents later documentary references, historical aliases or extra `REQ-*`-shaped identifiers from silently changing the implementation obligation universe.

The obligations command hard-fails when registered capability, source-requirement or active-screen counts diverge from the repository state baseline.


## E0-COVERAGE-001C

`cmdr-dev coverage-audit` compares authoritative active sets with owned canonical entities and observed source references.

The audit is report-only. A non-zero gap count does not by itself fail CI because the purpose of this stage is to expose documentary/reference inconsistencies without silently rewriting Product Spec. CI does fail if the authority sets drift from known repository invariants: 498 capabilities, 122 source Requirements, 56 active screens and 17 open decisions.

Gap records are deterministic and source-backed. Historical/reference-only identifiers are never promoted into active obligations merely to make the report green.
