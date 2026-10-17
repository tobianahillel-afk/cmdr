# E10-INV-120A — four-file Wazuh 4.x offline Case viewer

## Functional operator journey

A local analyst can now open a single authorized Case in the existing loopback HTML viewer with up to **four distinct Wazuh 4.x `alerts.json` files** loaded from disk, even if separate files reuse identical raw Wazuh alert IDs:

```sh
go run ./cmd/case-timeline-demo -principal=analyst-a \
  -snapshot=/absolute/path/owner.json \
  -wazuh-alerts=/absolute/path/alerts-part-1.json \
  -wazuh-alerts=/absolute/path/alerts-part-2.json \
  -into-case=case-a -serve-import
```

The CLI is the only path selector; the browser cannot upload files, choose ACLs or mutate Case sources. All feeds are bounded, parsed and owner-authorized before generating the ephemeral bearer or binding the listener. Identical alert IDs across separate files get separate SHA-256-derived source/entry identities, using the local trusted feed ordinal (`feed-N-wazuh-rawid`) in an isolated adapter after normalization. Both the Wazuh normalizer and JSONL importer remain unchanged; its provenance remains `operator-local-jsonl`, while the feed identifier is encoded in the ID without disclosing a local file path. The owner snapshot and every input file remain untouched. Authenticated Case discovery, Timeline pagination and typed source previews work across both feeds, including the original source descriptions; anonymous/foreign/denied requests remain blocked.

## Reuse-first design

The code composes the existing Wazuh 4.x normalizer, `readLocalEventFile`, `importLocalEvents` in-memory owner extension, `NewMemoryCaseTimelineHTTPHandlers` and first-party analyst UI. One small isolated adapter namespaces already-normalized event IDs before the unchanged JSONL importer; no new serializer, queue, database, generic ACL engine, frontend framework, external dependency or paid API. A streaming parser, Wazuh SDK or live remote ingestion would add trust boundaries without enabling this bounded offline use case.

## Compatibility, limits and risk

Single-file Wazuh viewer and single-file Wazuh cross-Case proximity CLI retain their prior output/identity behavior. **Multi-file input is viewer-only**; explicit cross-Case triage still accepts exactly one alert file. File count 1..4, per-file bytes 8 MiB and combined per-Case entries ≤5000; reject hardlinks, symlinks and file reuse; one bad second feed rejects everything before any bearer/listener or stdout. Feed numbers reflect local command order; they are not cryptographically authenticated Wazuh provenance or immutable file fingerprints. This is an offline synthetic principal/owner demonstration, not production IAM, SIEM ingestion, long-term persistence, remote Wazuh integration, signed forensic chain of custody or Shared Timeline Engine. OPEN-013 unresolved.

## Acceptance

End-to-end CI must exercise real disk files, duplicate IDs across two files, two independently authorized typed source previews, unauthorized and cross-tenant denial, multi-file ambiguity and all-or-nothing failure, plus Go race/vet/gofmt on an exact immutable PR SHA. One independent risk-based read-only L2 after CI. No merge implied by green tests.

## New offline multi-feed authorized Case triage (functional milestone)

An operator can analyze 2..4 Wazuh 4.x `alerts.json` files together, retaining distinct events when raw alert IDs overlap, without starting an HTTP listener or writing an intermediate owner snapshot:

```sh
go run ./cmd/case-timeline-demo -principal=analyst-a \
  -snapshot=/absolute/path/owner.json \
  -wazuh-alerts=/absolute/path/wazuh-a.json \
  -wazuh-alerts=/absolute/path/wazuh-b.json \
  -into-case=case-a -triage-import
```

The bounded JSON receipt `cmdr-local-multi-wazuh4-case-triage-v1` includes the imported-event count and the existing independently owner-authorized Case triage fields (visible entry counts by observed/inferred kind, certain UTC bounds, distinct visible source references and bounded top source counts). The output is assembled only after all feed files have passed validation and the chosen Case has completed authenticated Case/Timeline/typed-source traversal. Unknown, denied, cross-tenant or malformed Case/feeds/cancelled requests fail with no JSON receipt. No Wazuh alert description, full_log, bearer token, owner dataset or file path is in the receipt. This is descriptive local triage, **not** a live SIEM integration, a causal finding, an audit-grade custody record or deployed IAM.

Comparison: simply compose the existing `loadMultipleWazuhFeeds`, `collectDemoAuthorizedCaseTriage` and `localImportedTriage` (Go standard library); no second parser, source ACL implementation, correlation index or external Wazuh SDK. This change is scoped to PR #48 and does not merge or override the engineering foundation/PR #36. Relevant acceptance includes positive two-feed shared alert IDs; denied principal, bad second feed and cancelled context with **zero output**. On-demand new CI and L2 evidence must be verified separately before declaring this draft PR merge-ready.

Previous note that multi-file is viewer-only referred to an older PR snapshot and is superseded by this section; multi-feed viewer and 2-Case temporal proximity remain unchanged.
