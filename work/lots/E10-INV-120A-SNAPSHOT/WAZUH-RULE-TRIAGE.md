# E10-INV-120A-SNAPSHOT — authorized Wazuh 4.x rule-level triage

**Functional path.** Operator-selected offline Wazuh 4.x `alerts.json` -> existing JSONL normalization -> existing Case owner import -> authenticated Case/Timeline/typed-source page traversal -> deterministic rule level histogram (0..16) in the local triage receipt. The report also counts missing rule levels as unknown. The pre-existing `-triage-import -wazuh-alerts` CLI interface, viewer, two-Case aggregate and time proximity remain intact. Rule level is reported as raw Wazuh metadata, not a new risk score.

**Trust boundary.** We reuse `prepareWazuh4Case`, `importLocalEvents`, `collectDemoCaseTriageWithObserver`, SHA-256 entry identities, and Go standard library. A rule level can be counted only after the matching imported entry has passed existing bearer-authenticated Case and typed-source authorization. If even one imported alert is not visible, the entire output is withheld. The same in-memory bytes used for import are used for rule-level parsing, eliminating a second-disk-read race. The source file and original owner snapshot remain unchanged.

**Alternatives rejected/deferred.** A second search index or custom SOC rules engine adds new storage, source authority and external dependencies without need for a bounded 5,000-entry offline snapshot. A real Wazuh/manager network connector, durable store, central authorization provider and production retention/audit are still absent and not claimed. OPEN-013 stays unresolved.

**Validation.** Exact-SHA Go CI, tests for three distinct levels/unknown, denied principal/Case and unsupported level, reproducible read-only disk source and absence of preview/foreign-tenant leakage. Independent risk-based L2 needed for this new analyst-visible behavior.
