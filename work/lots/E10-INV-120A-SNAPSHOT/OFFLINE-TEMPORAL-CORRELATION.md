# E10-INV-120A — authorized offline five-minute temporal pivot

## User-facing increment
The existing offline `case-timeline-demo -compare-case=... -with-case=...` now emits up to **five shared certain five-minute UTC windows**, with each Case's visible event count and observed/inferred breakdown. This is the first machine-readable precise temporal co-occurrence pivot in the offline comparison, rather than merely intersecting two enclosing time intervals.

## Reuse-first decision
Reuse existing disk snapshot loader, one synthetic principal, authenticated Case discovery and source-authorized `/case-timeline` pagination for both Cases, plus Go standard library `time.Truncate` and `sort.Slice`. There is **no new permission engine, HTTP endpoint, persistence backend, index, schema, vendor tool or paid service**. Window aggregates are private, unexported Go receipt fields and are not duplicated into standalone Case receipts. Two Cases must complete independently before one receipt is emitted.

**Important ownership limit:** the current memory provider requires globally unique source references with a single Case owner. Matching raw source strings across Cases would be misleading or impossible here; therefore no shared-source identity, event identity or causal correlation is asserted.

## Security and measurement
Only server-authorized, known-timestamp entries enter the window counts. Uncertain clocks never become certain. Window count max five, stable chronological order, local stdout only. No bearer, raw event rows, source previews or denied IDs added to the receipt. This is nontransactional, unsigned, operator-controlled synthetic data; **not production IAM, live SIEM correlation, cryptographic custody or proof of causality**. No production SLO.

## Evidence
Tests use real disk JSON and two actual owner-authorized paginated Cases, with >200 rows in Case A, exact observed/inferred counts, denied sources excluded and a second Case within the 10:00–10:05 UTC bucket. A second test proves unknown source time yields no asserted shared bucket. Existing tests cover failure without emitting a partial comparison. Verify exact-SHA GitHub Actions and independent risk-based L2 before closing this lot; do not equate green CI with readiness for production.
