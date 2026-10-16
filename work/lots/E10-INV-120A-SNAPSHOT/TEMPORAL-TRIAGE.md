# E10-INV-120A-SNAPSHOT — analyst chronology gap triage (read-only)

## Functional increment

The local operator-selected JSON owner snapshot already supports Case discovery, authenticated filtered Timeline pages, per-source reauthorization, a full bounded analyst summary and reference grouping. A remaining investigative gap was the inability to inspect **the reliable time span and the largest gap between visible events across page boundaries**. The operator can now explicitly select **Analyze authorized events** to see these two temporal indicators alongside counts and sources.

The browser computes results exclusively from the same `/case-timeline` pages already admitted by server-side Case + Timeline Entry + source ACLs, without an extra endpoint or stored results. A `time_uncertain` event does not participate in gap calculation; a declared-certain time must be a finite RFC3339 timestamp. All instants normalize to UTC milliseconds. The 5000-entry ceiling remains bounded, with `O(n log n)` local sorting. Less than two certain events yields **no calculable gap**. A malformed certain time on any page discards the entire summary; permission loss/Case/token/filter change retains existing abort and invalidation semantics.

## Reuse-first comparison and security

- **Chosen:** native `Date.parse`, `Array.sort`, already authorized paginated HTTP and existing explicit analyst action. No additional authorization logic in JS.
- **Deferred:** new server aggregate endpoint, external analytics/time-series service, full SIEM correlation index, database or timeline event ownership changes. None is required for a bounded ephemeral local demo.
- **Limits:** a large gap between *visible, timestamp-certain entries* NEVER proves an absence of activity. Hidden/denied records, omitted sources and uncertain clocks may fill it. Millisecond browser resolution is not source-native nanosecond forensic precision. The source snapshots and principals remain host-controlled synthetic offline data, not real IAM, live ingestion, durable audit or production persistence. Canonical `OPEN-013` remains unresolved.

## Verification

- Node behavioral tests cover a 205-event/2-page Case with an uncertain middle instant and the exact 2s largest gap; an all-but-one-uncertain Case; rejection of a malformed *declared-certain* timestamp in the second page without partial display; and timezone-offset normalization across page boundaries.
- Existing Go HTTP integration proves disk JSON -> Case ACL -> authorized page and source preview; security remains server-side. Exact-SHA GitHub Actions and independent risk-based L2 must be verified before any positive CI/L2 assertion.

This is a scoped functional enrichment of E10-INV-120A, not a declaration of production-ready CAP-INV-110.
