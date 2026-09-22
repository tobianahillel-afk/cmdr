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
