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

The CLI is the only path selector; the browser cannot upload files, choose ACLs or mutate Case sources. All feeds are bounded, parsed and owner-authorized before generating the ephemeral bearer or binding the listener. Identical alert IDs across separate files get separate SHA-256-derived source/entry identities, using the local trusted feed ordinal in each Wazuh normalized ID (`wazuh-feed-N-rawid`). The existing JSONL importer remains unchanged; its provenance remains `operator-local-jsonl`, while the feed identifier is encoded in the ID without disclosing a local file path. The owner snapshot and every input file remain untouched. Authenticated Case discovery, Timeline pagination and typed source previews work across both feeds, including the original source descriptions; anonymous/foreign/denied requests remain blocked.

## Reuse-first design

The code composes the existing Wazuh 4.x normalizer, `readLocalEventFile`, `importLocalEvents` in-memory owner extension, `NewMemoryCaseTimelineHTTPHandlers` and first-party analyst UI. One optional local feed-number argument to Wazuh normalization namespaces IDs, before the unchanged JSONL importer; no new serializer, queue, database, generic ACL engine, frontend framework, external dependency or paid API. A streaming parser, Wazuh SDK or live remote ingestion would add trust boundaries without enabling this bounded offline use case.

## Compatibility, limits and risk

Single-file Wazuh viewer and single-file Wazuh cross-Case proximity CLI retain their prior output/identity behavior. **Multi-file input is viewer-only**; explicit cross-Case triage still accepts exactly one alert file. File count 1..4, per-file bytes 8 MiB and combined per-Case entries ≤5000; reject hardlinks, symlinks and file reuse; one bad second feed rejects everything before any bearer/listener or stdout. Feed numbers reflect local command order; they are not cryptographically authenticated Wazuh provenance or immutable file fingerprints. This is an offline synthetic principal/owner demonstration, not production IAM, SIEM ingestion, long-term persistence, remote Wazuh integration, signed forensic chain of custody or Shared Timeline Engine. OPEN-013 unresolved.

## Acceptance

End-to-end CI must exercise real disk files, duplicate IDs across two files, two independently authorized typed source previews, unauthorized and cross-tenant denial, multi-file ambiguity and all-or-nothing failure, plus Go race/vet/gofmt on an exact immutable PR SHA. One independent risk-based read-only L2 after CI. No merge implied by green tests.
