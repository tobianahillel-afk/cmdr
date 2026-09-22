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
