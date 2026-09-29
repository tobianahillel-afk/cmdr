# Agent context bundles

`cmdr-dev context` compiles the bounded context for exactly one strict-v2 work unit.

The bundle includes:

1. canonical product sources explicitly selected by the work manifest;
2. repository-wide engineering rules and source-authority policy;
3. current execution state and work graph;
4. the selected manifest and parent manifest;
5. the complete dependency closure, including available handoff/progress evidence.

It excludes unrelated repository content by default.

## Product reference resolution

- active capabilities resolve through Capability Register canonical-file columns;
- active screens resolve through the Screen Register canonical source;
- source Requirement IDs are checked against the pinned 122-ID catalog and include current product traceability plus the identity baseline;
- permissions must exist in the Permission Register and resolve to the canonical permission model/catalog;
- OPEN decisions must still be active and resolve to the unresolved-decision authority;
- implementation contracts and canonical objects are explicit product paths in the manifest.

Historical or reference-only identifiers are rejected instead of silently becoming implementation context.

## Content addressing

Every source is read from repository reality, hashed with SHA-256, and embedded with its authority layer and inclusion reason. The bundle digest is deterministic over the selected work unit, product baseline, dependency set, source paths, source hashes and reasons.

Bundles are generated on demand and are not committed by default. This keeps the repository small while allowing a fresh agent to reproduce the exact context used for a work unit.
