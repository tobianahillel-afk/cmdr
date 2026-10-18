# E10-INV-750A — One bidirectional, source-authorized Wazuh investigation journey

## What is integrated

**Integration, not reinvention.** PR #117 combined the existing independently authorized Wazuh cross-Case pair pagination and first-side exact-Entry discovery. Both sibling PRs #118 and #119 then passed all four GitHub Actions jobs independently but were based on the SAME frozen PR #117 HEAD and therefore their features could not safely be assumed to coexist.

This deliverable creates a native **two-parent Git merge/integration commit** whose ordered parents are:
- `44c2f89145d7bef0e590ac4a9928e4966158f618` — frozen PR #119, source-authorized Wazuh pair pagination filtered by EXACT `rule_id` AND `agent_id`;
- `bbe6360fd66ced154716cc0d483a256e47533c6d` — frozen PR #118, reverse navigation from a SECOND-Case Wazuh Entry ID to all authorized paired FIRST-Case leads in bounded ten-item pages.

The integrated source tree takes PR #119 (including its `normalized_wazuh_pair_pages.go` and pivot tests) and reconciles PR #118's CLI arguments and mutually exclusive dispatch into `main.go`; it adds the unmodified reverse-search implementation and its original tests. Both parent histories are preserved; neither prior reviewed SHA or branch was changed. A NEW end-to-end integration test crosses the actual feature boundary. No new engine, owner model or generalized ACL layer was written.

## Actual synthetic analyst workflow

1. Select TWO previously operator-imported and independently SHA-256-verified Wazuh Case archives. The Case/Timeline/typed-source owner enforces the principal, tenant/environment and source READ scopes independently for each.
2. Use the existing metadata-only `-list-saved-wazuh-pair-offset 0 -pair-rule-id 5710 -pair-agent-id 001` to select one authorized cohort. The 60-second nearest-peer matching and ten-result paging engine remains first-party; the pair's two Entry IDs are **navigation selectors only**.
3. From one returned second-side Entry ID, use `-find-saved-pair-second-entry EXACT_ID -find-saved-pair-second-offset 0`. The engine returns ALL correlated first-side leads, retaining the full authorized second cohort during nearest-peer determination and the existing ten-result/has_more semantics.
4. Choose one returned pair of exact canonical Entry IDs and separately call `-inspect-saved-pair-first-entry ... -inspect-saved-pair-second-entry ... -private-report /absolute/new-owner-dossier.json`. Both source owners are reauthorized before the previous 0600 no-overwrite private dossier is created. Neither discovery command emits Wazuh source previews, bearer tokens or URL listeners; private evidence never goes to stdout.

The pivot selects authorized source-derived cohort events and the reverse lookup selects exact authorized source Entry IDs. Neither expands Case access; a rule/agent selector never authorizes a source. CLI prohibits simultaneous pivot+reverse discovery, reverse+first-ID discovery, pagination+reverse, or pivot+private report. No false co-occurrence causality claim; this is a lab-local Wazuh temporal correlation.

## Integration acceptance

- The two-parent merge commit is a real published tree reconciled from both sibling branches; GitHub SHA parent relationships, base PR #117 and changed-file tree can be audited.
- A single Go integration test creates two independent authorized saved Wazuh owner snapshots with the SAME two-source rule/agent cohort. It proves filtered page count two, reverse second-side ID returns **the same two first-side canonical IDs**, and exact IDs reopen both actual private typed-source descriptions while stdout remains metadata-only.
- Conflicting selectors fail with **no partial receipt or artifact**.
- Both original sibling tests continue to run. CI must pass first-party `investigation-timeline-runtime` Go/Node/Linux, macOS and Windows and unmodified `engineering-kernel` security checks, particularly the 90% changed critical-security coverage gate. Do not waive, lower or merge over failures.

## Release disposition

This is a synthetic offline lab integration, not live Wazuh collection, production IAM, verified forensic custody or a general-purpose SIEM. All ancestor PRs remain DRAFT/UNMERGED; `main` contains only specs/documentation. This integration requires one independent risk-based READ-ONLY L2 on the actual immutable HEAD and canonical parent/merge ancestry, plus the repository's native branch/approval gates before any merge. No paid API, local PC, Codex or Work used.
