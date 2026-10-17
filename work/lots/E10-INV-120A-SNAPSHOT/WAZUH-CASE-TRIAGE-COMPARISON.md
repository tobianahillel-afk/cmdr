# E10-INV-120A — Wazuh 4.x single-Case triage and two-Case comparison

## Analyst journey delivered in this functional increment

Existing Wazuh 4.x `alerts.json` ingestion already supported the authenticated loopback analyst viewer and near-time cross-Case leads, but **not** a standalone first-Case aggregate triage or an independent two-Case aggregate comparison. The new explicit CLI modes compose the existing owner-authorized Case traversal with the already tested Wazuh normalization.

```sh
# No listener and no merged owner graph written. Sensitive local receipt.
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/owner.json \
  -wazuh-alerts=/absolute/alerts.json -into-case=case-a -triage-import > /private/wazuh-case-triage.json

# Comparison of two separately authorized Cases, with 5-minute windows and counts.
go run ./cmd/case-timeline-demo -principal=analyst-a -snapshot=/absolute/owner.json \
  -wazuh-alerts=/absolute/alerts.json -into-case=case-a -triage-import \
  -with-case=case-c > /private/wazuh-case-comparison.json
```

New receipt formats are `cmdr-local-wazuh4-case-triage-v1` and `cmdr-local-wazuh4-case-comparison-v1`. The latter uses the existing `collectAuthorizedCaseComparison`, which performs **independent authenticated Case and typed-source access checks** before writing one aggregate JSON receipt. The first uses `collectDemoAuthorizedCaseTriage`. The existing `-compare-proximity` mode still produces nearest-time leads; `-serve-import` still exposes the authorized loopback UI. Paths/files and Windows/macOS/Linux compatibility are inherited from the first-party Go file reader and offline importer.

## Reuse-first and security

**Chosen:** existing local Wazuh parser, `importLocalEvents`, `collectDemoAuthorizedCaseTriage`, `collectAuthorizedCaseComparison`, `selectCaseProximity`; one shared internal Wazuh-to-Case import step for each analysis mode. This is simpler than creating a new rule-correlation engine, database, ingest service or IAM layer, and does not change Shared Timeline Entry ownership or `OPEN-013`.

Tests use real on-disk two-Case owner and alert files, 208 authorized first-Case entries spanning two HTTP pages, first-Case triage, separate second-Case aggregate with confirmed time overlap, zero leaked preview/hidden tenant fields, original input file immutability, foreign/denied Case failure, invalid Wazuh failure and cancellation. No browser, bearer or merged owner snapshot is printed.

## Limitations

The owner principal and ACL catalog are **operator-supplied local synthetic authority**, not production IAM; there is no live Wazuh API, persistent Case database, verifiable forensic chain of custody, SIEM streaming or producer-authenticated telemetry. Case comparison reports are unsigned, potentially sensitive analyst receipts. The five-minute overlap is an observational hint, never evidence of causality. Performance measures from earlier synthetic SHA are not transferred to this change.

Lot status remains READY until the source HEAD passes exact CI and material independent L2, and real provider/persistence choices are resolved.
