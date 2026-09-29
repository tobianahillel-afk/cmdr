# Hunt Management runtime boundary

This directory contains the bounded executable read-only runtime of **CAP-INV-005 Hunt Management**.

## Implemented scope

The runtime validates and immutably projects caller-owned Hunt workspace context only:

- one explicit authorized Tenant and environment;
- explicit question, scope, ordered RFC3339 period and owner;
- copied contributor references;
- typed tenant-scoped accessible references to Query, Search Job, Hypothesis, Case and Incident;
- canonical ownership labels without ownership transfer;
- copied version/provenance metadata when supplied by the caller;
- deterministic stale/partial/unavailable limitations;
- Search Job visibility constrained to a projected Query;
- deep-copy isolation from caller-owned slices and references.

All rejection paths return no projection. The module uses the Go standard library only and performs no network, filesystem, persistence or query-execution work.

## Explicit exclusions

This runtime does **not** create a canonical Hunt object or final Hunt state machine. While OPEN-013 remains unresolved it performs no Hunt create/update/activate/pause/close/archive transition, contributor mutation, Case promotion/link mutation, Hypothesis mutation, Query/Saved Search mutation, or Search Job execution.

It also selects no storage engine, provider, retention policy, collaboration backend, final query dialect or final Hunt UI.

External product-runtime dependencies remain deny-by-default.
