# E10-INV-120A — Authorized Wazuh rule-ID cohorts

## Functional capability

Wazuh 4.x offline `alerts.json` (local, operator-selected file) → existing strict normalization → imported Case → existing authenticated owner Case/Timeline/typed-source traversal → actionable aggregate grouped by Wazuh **`rule.id`**. This builds on the existing `rule.level` histogram and 5-minute time windows.

The receipt remains `cmdr-local-wazuh4-case-triage-v1`, additive fields only:
- `distinct_rule_ids`: number of distinct known rule IDs.
- `unknown_rule_ids`: number of imported authorized Wazuh alerts missing rule ID.
- `top_rules`: up to ten rule IDs ranked by authorized count descending, ties by decimal rule ID ascending; includes count, highest *known* `rule.level` (nullable), first/last *certain* UTC alert times (omitted if entirely time-uncertain).
- Alert IDs, agent names, `full_log`, source previews and tokens do not appear in this aggregate. A rule ID is an analyst pivot, NOT a source authorization, risk/incident score, automatic detection or causal relationship.

This extension reads `rule.id` from the SAME bounded source bytes already imported, pairs metadata with its deterministic per-Case Timeline Entry ID, and counts it ONLY inside the observer of `collectDemoCaseTriageWithObserver` after existing Case/Entry/typed-source access checks. The Wazuh `rule.id` must be an ASCII decimal string up to 16 characters when supplied. No new ACL, parser framework, database, webhook, provider, network listener, SDK or paid API is introduced. The 13-group tie test asserts a deterministic top-ten bound, and a fractional-second test ensures earliest/latest ordering is correct independently of input order.

## Security and limitations

Input contains potentially sensitive SOC data; no real customer secrets should be passed to this **synthetic loopback/offline lab**. The operator-supplied owner snapshot is synthetic authorization, not a real IAM or deployed SIEM connector. Zero output on invalid metadata or denied Case; no partial receipt before complete reauthorization and aggregation; no input file is modified. No forensic custody, durable storage, production security logs, service account, cross-Case rule search, or decision of `OPEN-013`. Original histogram/time-window JSON fields remain backward compatible.

Risk-based L2 should focus on actual authorization boundaries, Wazuh 4.x rule-ID formatting and consistent counts with `rule.level`, temporal correctness, no source preview disclosure, and preservation of prior CLI/HTTP journey. R3/R4 style matters should not hold the functional lot.
