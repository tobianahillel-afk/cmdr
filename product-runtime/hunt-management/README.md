# Hunt Management runtime boundary

This directory is reserved for the first bounded executable runtime of **CAP-INV-005 Hunt Management**.

## Current state

**Preimplementation only.** This README is the sole file intentionally present in the runtime boundary at this stage. It carries no executable product behavior and must not be treated as runtime coverage, authorization-negative evidence, tenant-isolation evidence, performance proof, persistence, or a canonical Hunt domain object.

## Bounded future scope

Later verified tasks may implement only the read-only workspace core already materialized in `WAVE-CAP-INV-005-CORE-001`:

- validation of one caller-owned tenant/environment Hunt workspace envelope;
- immutable projection of question, scope, period, owner and contributors;
- provenance-preserving references to accessible Query and Search Job records;
- references to accessible Hypothesis, Case and Incident context without ownership transfer;
- deterministic completeness diagnostics;
- explicit stale, partial and unavailable-reference limitations.

## Explicit exclusions

This boundary does not select or implement:

- a canonical Hunt object/schema or final Hunt state machine;
- Hunt create/update/activate/pause/close/archive mutation while OPEN-013 remains open;
- Case promotion/create/link mutation;
- Hypothesis create/link/status mutation;
- Query mutation, Saved Search/Query Asset mutation or Search Job execution;
- CAP-INV-003, CAP-INV-006, CAP-INV-008 or CAP-INV-103 runtime assumptions;
- collaboration/presence backend selection;
- storage/provider/retention implementation;
- a final query dialect;
- a dedicated final Hunt UI;
- a production Hunt workflow SLO or full-capability completion claim.

External product-runtime dependencies remain deny-by-default. The initial implementation direction is Go standard library first unless a later evidence-backed technical decision proves otherwise.
