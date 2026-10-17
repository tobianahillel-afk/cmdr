# E10-INV-120A — actionable cross-Case source-reference pivot

## Implementation and choice

The loopback analyst comparison already retrieves full authorized pages of two Cases, counts common exact source-reference strings and identifies overlapping time buckets. Now each top-five common reference offers an **open representative event from Case 1 or Case 2** action, taking the analyst directly from comparison to the existing typed-source owner preview. The representative entry ID comes from the corresponding server-authorized Case page. No untrusted owner object ID or `source_ref` is used as an authorization token. Both existing source navigation and owner preview handlers independently reauthorize Case/Entry/typed-source access on click.

Reuse-first comparison: existing Go Case/Timeline/typed-source handlers + first-party `openSource(entryID)`, JS `Map` grouping and same-origin `fetch` meet this bounded functional goal. No new index/database/search service, ACL engine, third-party library, API or migration. UI pivots are invalidated on Case, bearer or filter changes; Case switching discards old views. Source labels remain escaped text and length-capped.

Limits: shared source-reference strings are triage leads, not proof of same real-world object, causal link, transactionally consistent view, or production authorization. Operator-selected JSON dataset and principal tokens remain synthetic, loopback-only. `OPEN-013` stays open; production persistence, identity provider and forensic chain of custody are absent.

## Acceptance

Node browser tests cover two Case HTTP page reads, selected representative Entry ID, reauthorization through both routes, absence of bearer or source_ref in URL, revoked destination, stale actions and text-safe labels. Exact-SHA CI and one independent risk-based L2 are needed before declaring the milestone complete.
