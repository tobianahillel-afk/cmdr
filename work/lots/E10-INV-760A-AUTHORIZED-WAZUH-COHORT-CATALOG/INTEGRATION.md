# E10-INV-760A — Explore correlated Wazuh cohorts before exact pair navigation

## Product capability
From two independently authorized and SHA-256-verified saved Wazuh Cases, the existing offline investigator CLI now discovers matching rule+agent cohorts **without requiring an analyst to guess rule or agent IDs**:

```
case-timeline-demo -principal analyst-a -snapshot FIRST.json -compare-saved-case FIRST \
  -expect-snapshot-sha256 FIRST_SHA256 -compare-saved-with-snapshot SECOND.json \
  -compare-saved-with-case SECOND -expect-other-snapshot-sha256 SECOND_SHA256 \
  -list-saved-wazuh-cohorts-offset 0
```

The JSON receipt contains `rule_id`, `agent_id`, `matched_pairs`, `first_events`, `second_events` for **correlated cohorts only**, ranked by matched pairs then stable rule/agent identifiers. Ten cohorts per page, offset 0..5000. A selected cohort is opened with the **existing** `-list-saved-wazuh-pair-offset 0 -pair-rule-id RULE -pair-agent-id AGENT` route, followed by the existing reverse/first Entry selectors and separately reauthorized private report.

## Scope and security
This is a synthetic/offline local DFIR workflow; correlation is a temporal lead (<=60 seconds), not event identity or causation. Both source owners must pass existing principal/tenant/Case/source/SHA checks before any result is written. The catalog never includes raw source previews, private dossiers, listeners or credentials. Mixed catalog/private/selectors are rejected before source access, errors produce no partial stdout. Counts reuse the event-time-sorted first-party Wazuh observations and the same 60-second pairing predicate as existing pagination.

## Dependency
Stacked on immutable PR #121 at `d3a930d1643e2db0b63aafbb0602379b8144b38d`. Parent PR and all upstream PRs remain draft/unmerged, including repository `main` which currently contains only README. No database, external service, dependency or paid AI API.

## Acceptance
- Real two-authorized-Case saved archive → catalog → rule/agent filtered pair page
- Stable ranking/pagination and no unmatched cohort disclosure
- Untrusted principal, wrong Case, changed SHA and mixed private selector fail closed with no output
- Existing GitHub Actions investigation-timeline-runtime + engineering-kernel CI on exact head
- One context-separated read-only risk-based L2 for this product feature, never automatic merge
