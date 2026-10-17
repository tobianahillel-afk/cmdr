# E10-INV-120A — functional two-Case comparison handoff

## Integration implemented

The existing loopback analyst UI now lets the operator voluntarily save an **aggregated JSON comparison** of two independently authorized Cases. It reuses the existing complete paginated Case/Timeline/source ACL HTTP reads, the existing in-browser temporal overlap and source-correlation computations, and the already-proven local Blob download pattern for a single Case. No new backend, source owner, synthetic identity, database, SDK or persistent write was added.

The comparison file includes only two selected Case IDs, current filters, visible row counts, observed/inferred and uncertain counts, bounds of certain UTC instants, boolean/indeterminable range overlap, total exact shared-ref count, up to five truncated 120-character shared reference names with counts, and at most five already-authorized five-minute common UTC buckets. It excludes bearer credentials, raw event arrays, source previews, unfiltered sources, and any claim of actual causal correlation.

The download is disabled until all pages of BOTH Cases completed and were freshly authorized. Denial, malformed/incomplete second Case data, cancellation, token/Case/filter change or running a new comparison invalidates the receipt. Download stays an explicit client-side action without a network upload. The result is **not** transactionally isolated: authorization may change after pages are fetched; the unsignable snapshot cannot be used for evidence custody and must be treated as synthetic analyst notes.

## Reuse-first choice and alternatives

Reuse the JavaScript `Map` and date grouping, page consistency verification, and browser-native Blob URL. A new cross-Case correlation/index service, export endpoint, durable repository or third-party library would introduce new authorization ownership and custody questions that are not justified by an operator-loaded read-only demonstration capped at 5,000 entries per Case. Do not confuse an exported triage summary with a production forensic report.

## Verification and remaining gaps

Node VM browser tests cover complete 205+3 paginated reads with positive local JSON download, bounded fields/no bearer, HTTP header-only authentication, denied peer, token rotation and changed Case selection. Other tests cover cross-Case temporal pivot and invalid entries. Verify exact-SHA Go/Node CI before claiming completion. No production IAM, trusted persistent owner repository, live SIEM feed, audit backend, signed report or resolved Shared Timeline Entry ownership (OPEN-013) exists. This is a functional integration increment of READY E10-INV-120A-SNAPSHOT, not full CAP-INV-110 production readiness.
