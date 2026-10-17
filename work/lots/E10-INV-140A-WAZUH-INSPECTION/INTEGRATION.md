# E10-INV-140A — Wazuh alert → Event Inspection → Event Search PivotDraft

This is a **functional offline integration** on top of the already validated multi-Wazuh Case viewer (PR #50). It reuses the existing validated Wazuh 4.x JSONL normalizer, namespaced two-to-four-feed Case importer, real synthetic Case/Timeline/typed-source owner HTTP handlers, the first-party `event-inspection.Project` API and `event-inspection.PreparePivot` API. It adds no new ACL service, parser, third-party dependency or storage.

## Operator flow

```sh
cd product-runtime/investigation-timeline
go run ./cmd/case-timeline-demo \
  -snapshot=/absolute/path/to/operator-owner.json \
  -principal=analyst-a \
  -into-case=forensic-case-on-disk \
  -wazuh-alerts=/absolute/path/to/alerts1.jsonl \
  -wazuh-alerts=/absolute/path/to/alerts2.jsonl \
  -inspect-wazuh-entry=local-jsonl-entry-<server-derived-sha256> \
  -pivot-wazuh-field=wazuh.rule.id \
  -pivot-start=2026-10-17T09:00:00Z \
  -pivot-end=2026-10-17T11:00:00Z
```

The entry identifier is the one returned by the already-authorized `/demo-wazuh-drilldown` journey. Omit all three pivot flags for inspection alone. Output is JSON containing the Event Inspection projection and an optional *draft*, never an executed Search Job.

**Security behavior:** selected Wazuh identity must belong to the imported feed, complete owner-authorized Case Timeline pagination must match all imported alerts, and the selected typed source is freshly opened before projection. Only then are bounded `wazuh.rule.id` and `wazuh.agent.id` metadata supplied as normalized fields. The source owner preview is expressly disclosed as **rendered content** to this selected sensitive inspection; raw payload access is always denied. Source metadata never determines authorization. Cross-tenant selection, denied sources, unknown entries, bad time bounds, incompatible modes and cancellations produce no partial response.

**Limits:** this is synthetic local operator-controlled IAM, not production authorization, live ingestion, a persistent Case/Job, chain of custody, browser UI or deployed search connector. Pivot drafts use the existing CAP-INV-004→CAP-INV-002 contract but are **not executed**. The next meaningful integration is to hand off a draft to a real authorized Event Search provider, after the blocked/reviewing PR #52 can be safely reconciled. Performance and security claims remain bounded to existing test fixtures.

**Acceptance:** compile and exercise disk→Wazuh normalization→Case→Timeline→typed source→Event Inspection→optional PivotDraft, including two colliding Wazuh alert IDs and cross-tenant negatives, via GitHub Actions Linux/macOS/Windows + engineering kernel. Remain draft until CI and risk-based independent L2; no merge before dependency PR #50 and branch protections allow it.
