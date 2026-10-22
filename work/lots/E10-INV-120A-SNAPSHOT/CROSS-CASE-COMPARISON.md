# E10-INV-120A — comparison of two authorized Cases in the local analyst journey

## Functional path and reuse-first decision

An analyst may compare **two different Case IDs** belonging to the current synthetic principal's accessible Case set. The browser reuses the existing authenticated, server-authorized `/case-timeline` paginated endpoint for each Case. It shows the number of **exact source-reference strings shared** by the two already authorized sets, their respective visible-entry counts, and at most five leading source groups. The same event-kind, text and date filters apply. A common source reference is only a triage lead, not causal or identity proof.

Chosen solution: native browser `fetch`, `Map` and deterministic sort, bound to existing read-only Case/Timeline/source owner ACLs. Rejected: a new server authorization shortcut, UI-supplied permissions, cross-tenant correlation index, Elasticsearch, new provider, database, third-party library, or backend write. This is useful cross-Case navigation without inventing an additional security authority.

## Security and bounded semantics

- Each read requests `/case-timeline` with `Authorization: Bearer` **in headers only**, `cache: no-store`, `credentials: omit`, `redirect: error` and `mode: same-origin`.
- The API already checks current Case membership, Timeline Entry visibility and typed-source permission. **Both** Cases must be readable for a result. A failure, incomplete page, inconsistent total, duplicate ID, malformed authorized entry, revoked token or changed form selection discards the entire comparison.
- Max 5,000 authorized entries/Case, 200 entries/page, at most 50 sequential page reads across the pair; no raw entries, source contents or bearer exported or stored. All output uses `textContent` and truncates displayed source labels to 120 code units. No new server route or persisted state.
- The two page traversals are **not transactionally consistent**, so counts are a best-effort observation, and hidden/denied/uncertain sources may exist. Never interpret overlapping source reference strings as a causal graph or proof of same real-world object.
- Operator-controlled disk snapshot, bearer and principal remain synthetic loopback fixtures; this is **not** production IAM, persistence, live SIEM ingestion, audited export, or complete CAP-INV-110. `OPEN-013` remains unresolved.

## Acceptance

Node browser behavioral tests check actual two-Case HTTP pagination (205 events across two pages), shared-source counts, no bearer in URLs, denied second Case failing closed with no partial disclosure, cancellation on token change, and same-Case/empty-bearer rejection. Existing Go HTTP integration tests cover owner snapshot to Case and typed-source access. CI and an independent risk-based L2 must be verified at exact SHA before marking this functional milestone complete.
