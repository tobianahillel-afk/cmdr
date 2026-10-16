# E10-INV-120A — actionable cross-Case investigation pivot

## Functional increment and technology selection

The existing operator-selected offline snapshot already supports Case discovery, Timeline pagination, server-validated typed-source previews, temporal triage and an authenticated **read-only** comparison of source-reference strings across two authorized Cases. However, correlation results were text-only, requiring investigators to manually re-enter the Case and search to investigate a shared reference.

This increment adds bounded UI actions for up to five shared references (both Case A and Case B), reusing the existing authenticated `/case-timeline` search **without introducing a new backend, cross-Case index, database, provider, storage or frontend dependency**. A tap selects the Case and exact visible source-reference search text, clears old results/preview, resets the page, and issues a **fresh** server-authorized HTTP request. This is a triage pivot, not automatic owner-object opening or proof of causal linkage. Long refs beyond the 96-character query contract remain text-only; no silently truncated filter is supplied.

## Security and failure semantics

The complete comparison must first have succeeded for **both** Cases. The action is invalidated by any current Case, principal bearer, comparison Case or filter change, and by aborted/incomplete requests. Clicking does not reuse an authorization decision: new Case+Timeline+source ACLs are evaluated server-side, and revoked/foreign Cases return a generic empty/denied state. No tokens are stored or included in URLs. Every displayed reference uses `textContent`, not HTML injection. No new route is exposed.

The snapshot is operator-controlled synthetic data in a loopback-only demonstration, with per-request in-memory ACLs rather than live production IAM. This is NOT live SIEM ingestion, transactional access consistency, durable evidence storage or a cryptographic forensic report; `OPEN-013` remains unresolved.

## Testing and review

Behavioral Node tests use the real `analyst.js` VM and HTTP mocks to exercise two authorized Case reads followed by a newly authenticated pivot search, revoked destination Case, invalidated stale action and long-reference fallback. Existing Go tests cover actual operator-file CLI to Case/Timeline/source HTTP paths. Verify GitHub Actions CI **at the exact published SHA** before treating tests as successful; request a single independent risk-based read-only L2 on the functional milestone when stable. This PR remains a draft based on PR #36; do not merge into main prematurely.
