# E10-INV-770A — Source-authorized cohort discovery in consolidated Case investigation

Integrates the **existing implementation** from PR #122 into the broad, distinct-owner DFIR release candidate PR #123, without changing either immutable parent PR.

## Analyst journey

With two saved Case owner snapshots separately authorized for the requesting principal and verified against independently held SHA-256 hashes:

1. Invoke \`-list-saved-wazuh-cohorts-offset 0\` to discover Wazuh rule+agent cohorts with real 60-second cross-Case matches.
2. Use an exact catalog \`rule_id\`/\`agent_id\` with the existing \`-list-saved-wazuh-pair-offset 0 -pair-rule-id RULE -pair-agent-id AGENT\` mode.
3. Select a returned Wazuh Entry ID to use existing forward/reverse discovery or the explicitly requested private dossier path, **rechecking both owners and source ACL**.

The catalog returns up to 10 **metadata-only** cohorts per offset, without raw source previews or private evidence. No new dependency, product service, cross-tenant join, network API, or merge authority.

## Integration proof

- Reuse the original PR #122 real-two-Case integration tests, permission denials and deterministic pagination.
- New integration regression uses the full PR #123 owner and pair navigation stack. After revoking a first-Case typed source and recomputing its valid snapshot SHA-256, catalog counts and exact cohort pair listings **both** remove the unauthorized lead while preserving other authorized cohorts; metadata-only output must not leak previews.
- Exact candidate HEAD is governed by GitHub Actions (Go and Node on Linux plus platform tests). This branch is stacked on PR #123, **not on main**; none of the parent draft PRs has been merged.
- One context-separated read-only L2 on the stable increment; retain all native GitHub approval and review rules.

## Original source and scope

The catalog code is reused from immutable PR #122 at \`ef0fd7088447f85ed241ef1115645dd3aa8366e0\` and is integrated with PR #123 at \`18fbbeb871813e7fb5ac1cc89dec575303f41805\`. The new integration test is specific to the combined code paths.
