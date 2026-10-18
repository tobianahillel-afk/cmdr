# E10-INV-500A — Make bounded mixed Case searches correctly source-family scoped

## Functional fix and exact first-party boundary

A saved and owner-authorized mixed Wazuh 4.x + generic normalized JSONL Case can contain thousands of imported events (up to the existing Case cap), but an actual first-party Event Search Job is limited to **500 candidates and 200 matching results**. An analyst may explicitly filter by `source_family=generic`, `wazuh`, or `all`.

Previously the `500 candidates` counter was incremented **before** the source-family filter. In a perfectly valid Case with 501 Wazuh alerts and one generic event in the same time window, an otherwise selective `generic` search failed with a candidate-limit error even though it requested just one relevant entry. This blocked source-specific hunts at realistic mixed feed volumes.

**Minimal fix:** preserve the exact same Case/Timeline traversal and canonical first-party typed-source reopening **for EVERY in-window imported entry**, including entries of the unselected family. Only after the source is reauthorized, the already trusted owner-attached Wazuh metadata identifies its family. The existing counter is moved immediately **after the source-family filter**, before evaluating any query predicate. Therefore:
- `generic` with 501 Wazuh + one generic in scope succeeds with one authorized result;
- `wazuh` with 501 Wazuh still denies (500 ceiling unchanged);
- `all` with 502 selected entries still denies (500 ceiling unchanged);
- if a Wazuh source owner is revoked, even `generic` query still fails closed before returning a partial result, because ignored-family entries are NOT skipped from first-party authorization;
- source preview and bearer never appear in the completed Search Job, no new permissions/metadata trust, no changes to 200 result cap or event search orchestrator.

The workload scans bounded owner data exactly as before; it does not add extra parsers, auth wrappers, APIs, state, database or new correlation logic. The Case cap still bounds traversal even for family-specific searches. This is a concrete functional fix and a least-privilege security-preserving placement change, not coverage gaming.

## GitHub PR topology and release integration

This branch descends directly from the immutable #90 HEAD `7a5f9441fe17b6d50d001c1ce5f1209a813f610e`, preserving all reviewed code. The new PR intentionally targets the immutable **GREEN exact-head #89** `250c0e27a71dc92adfe2117d07e7ece7b4474d61` instead of the larger #81 integration base, so reviewers and mandatory changed-lines security coverage can separately inspect the six #90 integration delta commits and THIS functional fix while retaining all #89 code. This is NOT a waiver of red #88/#90 PRs; no source SHA is modified, native reviews/branch rules/parent dependencies still apply. No forced merge.

Exact-head CI required: Linux Go + Node/browser, macOS, Windows and unchanged `engineering-kernel` security-critical changed-line 90% coverage floor. One risk-based independent read-only L2 on complete #89..HEAD if stable; do not duplicate an earlier #90 L2. This remains an offline synthetic authorized owner, not production SIEM IAM/live ingestion or signed custody.

## Test

The behavior integration test creates an actual saved mixed Case from **501 real Wazuh JSONL alerts and one real generic normalized JSONL event**, uses the same authenticated `POST /demo-local-search`, checks one completed generic-family Job ref, rejects Wazuh/all family queries above 500, and verifies revoking a Wazuh typed source still denies the whole generic query without leaked refs. No arbitrary test-only backend or mocked permission grant.
