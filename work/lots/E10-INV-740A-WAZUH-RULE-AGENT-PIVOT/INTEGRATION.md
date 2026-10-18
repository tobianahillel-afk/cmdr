# E10-INV-740A — Rule/agent-scoped Wazuh pair pivot

## Why this is a product capability, not another correlation engine

The existing first-party two-Case Wazuh 4.x investigator already supports source-authorized, SHA-256-verified temporal pairs (<=60s), paginated metadata-only discovery beyond the first ten, direct discovery by exact first-Case Entry ID, and a separate explicitly requested private pair dossier with real source reauthorization. Without a cohort pivot, an analyst investigating one concrete rule/agent combination must traverse *all* eligible pairs over several pages and manually discard unrelated alerts.

This increment **reuses exactly the existing owner validators, source-derived Wazuh cohort collection, nearest-peer deterministic matching and page engine**. It adds two **optional CLI flags** available ONLY with the existing metadata-only `-list-saved-wazuh-pair-offset`:

```sh
go run ./cmd/case-timeline-demo \
  -principal analyst-a \
  -snapshot /absolute/first-authorized-wazuh-case.json \
  -compare-saved-case CASE_A -expect-snapshot-sha256 FIRST_64HEX \
  -compare-saved-with-snapshot /absolute/second-authorized-wazuh-case.json \
  -compare-saved-with-case CASE_B -expect-other-snapshot-sha256 SECOND_64HEX \
  -list-saved-wazuh-pair-offset 0 \
  -pair-rule-id 5710 -pair-agent-id 001
```

Both independently authorized owners are fully verified before selecting an exact `{rule_id, agent_id}` cohort from already source-authorized Wazuh event maps. Existing page semantics remain **ten authorized matching pairs plus `has_more`**, with offset relative to *matching* results, not global events. The JSON receipt retains format `cmdr-local-authorized-wazuh-pair-page-v1` and adds optional `rule_id`/`agent_id` echo (the operator's exact filter). Unfiltered pages remain byte-compatible except normal JSON formatting, with no filter fields. No new authorization system, heuristic, index, daemon, provider, package, paid API or generic query DSL. Match is exact/case-sensitive. If a cohort is absent, the result is empty without disclosing whether the other Case contains it.

The returned Entry IDs are navigation selectors only, never access grants. The operator can use them in the already existing separate exact-ID private source drilldown; typed Source rights on BOTH Cases are verified again before the 0600 no-replace dossier is written. Temporal co-occurrence does not prove causation.

## Trust boundary and acceptance

- **Both** `-pair-rule-id` and `-pair-agent-id` are required when either is supplied; identifiers are nonempty and bounded to 128 bytes. A pivot without page selection, or combined with exact-ID discovery, sensitive/private dossier, Wazuh event previews or unrelated CLI modes, fails before metadata or file publication.
- The filter happens only **after** independent SHA-256 input verification, principal/tenant/environment Case/Timeline/typed-source authorization and Wazuh metadata extraction. No client-supplied rule/agent values become authorization grants. 5000-offset limit and fail-closed owner checks stay unchanged.
- End-to-end test creates two synthetic imported Wazuh Case owners with **two distinct rule+agent cohorts**; unfiltered output contains four pairs, filtered output two each, mismatched cohort zero. Typed-source drilldown from a filtered pair opens the correct private descriptions without emitting them to stdout.
- Exact negative tests: missing filter half, oversized value, wrong second snapshot fingerprint, foreign principal, foreign Case, conflicting direct-ID discovery, private inspection or generic-preview mode, private artifact mode; modified second JSON still parseable but SHA mismatch yields zero first-Case metadata.

## Integration and release decision

This is an incremental function on the immutable approved-for-L2 (not merged) PR #117 SHA `276b9b4211563873d98bc8c54ba57cea2fe4edf0`. Its sibling PR #116 is an independently green alternative representation of the same #114/#115 integration and must not be blindly merged alongside #117; select one canonical path after targeted diff reconciliation. All ancestor L2-frozen SHAs are unchanged. Run exact-HEAD Go/Node on Linux, Go on Windows/macOS and unchanged engineering kernel, especially mandatory 90% changed security-critical coverage. This PR and its entire parent stack remain DRAFT/UNMERGED until all checks, independent risk-based L2 and required native reviews allow landing. This is a synthetic offline lab capability, not live ingestion, production tenant IAM, signed forensic custody, or a claim that `main` already contains the runtime.
