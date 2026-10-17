# E10-INV-120A — Wazuh cohort → authorized Timeline event drill-down

## Functional delivery

The loopback analyst browser can now select a Wazuh 4.x agent-ID or rule-ID cohort and inspect its first 25 **individually authorized Timeline events**, then open a source through the already existing `/case-source-navigation` permission check. This closes the gap between aggregate-only Wazuh cohorts and an actual investigation pivot. The new optional route is `GET /demo-wazuh-drilldown?case_id=...&agent_id=...` or `rule_id=...`, only in operator-controlled Wazuh-import mode.

Reuse-first: the existing Wazuh 4.x parser, immutable raw bytes from the original local import, SHA-derived normalized entry mapping, authenticated Case/Timeline/source pagination, and simple Go `net/http` plus vanilla browser JavaScript. A new SIEM backend, index, data store, framework, or duplicated authorization service would add complexity without delivering value in this bounded offline laboratory. No paid API.

## Trust and release boundaries

The analyst's cohort selector is a filter **only**, never an ACL, identity proof or tenant claim. The bearer and exact Case are checked before any data access. The Go handler traverses **all imported alerts** with the same owner Case/Timeline/typed-source access checks on every drill-down, requires every imported alert to be visible before emitting a JSON response, and never serializes Wazuh raw alert IDs, hostnames, full logs, or source previews. An event ID itself is returned only if authorized; clicking it uses the separate existing navigation ACL check. Invalid, foreign, anonymous, cancelled, malformed or ambiguous requests fail closed without partial results.

The browser validates returned shape and source IDs, renders only inert text/buttons, aborts/discards stale responses on bearer/Case changes, and provides a maximum of 25 representative events per cohort with a truthful total. This **does not** implement a fully paginated query service, live Wazuh API integration, persistent Case database, production IAM, audit backend or forensic chain of custody. Wazuh agent ID and rule ID are untrusted source metadata; shared labels are investigation leads, not causality or host identity. OPEN-013 is unchanged.

## Evidence

Go end-to-end HTTP tests use real-shaped local Wazuh JSONL, a 2-Case operator snapshot, bearer HTTP handlers, authorized agent/rule/no-match drilldowns, cross-Case/anonymous/POST/malformed denial, cancelled requests, ordinary non-Wazuh route absence, and exclusion of private raw fields. Node browser tests exercise the complete cohort → drilldown → existing source navigation action, and a denied drilldown with zero partial rows. Final CI conclusions and risk-based L2 must be verified on an exact published SHA before closing this lot.
