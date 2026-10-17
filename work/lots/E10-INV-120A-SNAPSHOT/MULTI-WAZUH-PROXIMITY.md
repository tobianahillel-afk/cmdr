# E10-INV-120A — correlated multi-feed Wazuh investigations

## Operator journey

A local analyst can now choose **2–4 distinct, bounded Wazuh 4.x alerts.json JSONL files**, import them into one Case using the pre-existing in-memory owner snapshot, and request **nearest certain-time leads against a second independently authorized Case**. The command composes existing Go services, with no intermediate owner-data export:

```sh
go run ./cmd/case-timeline-demo \
  -principal=analyst-a -snapshot=/absolute/path/owner.json \
  -wazuh-alerts=/absolute/path/sensor-a-alerts.json \
  -wazuh-alerts=/absolute/path/sensor-b-alerts.json \
  -into-case=case-a -triage-import -with-case=case-c \
  -compare-proximity -proximity-seconds=1
```

Result: a bounded `cmdr-local-multi-wazuh4-case-proximity-v1` JSON receipt with the number of imported events and up to 25 **temporal leads**. A lead is **not a causal correlation** or a verified incident. Each file contributes a distinct prefix to normalized event IDs, even when raw Wazuh alert IDs repeat across independently exported sources.

## Architecture / reuse decision

Reuse unchanged `normalizeWazuh4Alerts`, `readLocalEventFile`, `namespaceWazuhEvents`, `importLocalEvents`, `newDemoHandlerWithData`, `collectAuthorizedCrossCaseSequence`, and `selectCaseProximity`. One shared multi-file import path feeds either existing loopback Case viewer or the new cross-Case analyst receipt. This is smaller and safer than duplicating an importer, creating a Wazuh SDK, database, index, event bus, correlation engine, or external IAM provider. No paid dependency.

The second Case is independently traversed through the existing HTTP owner ACL for Case / entries / typed source before any receipt is written. One invalid Wazuh file, duplicate path/hardlink, unauthorized principal/Case, unavailable source, malformed Wazuh JSON, invalid timestamp, unacceptable interval or cancelled context fails closed without stdout or bearer/listener. Inputs and owner graph on disk remain untouched. Legacy single-file Wazuh proximity retains the original format.

## Explicit deployment limits

Bounded 2–4 local operator-selected files (8 MiB maximum per file, up to 5000 entries per Case), local synthetic principal and owner ACLs, no production IAM, no live ingestion, no persistent Case database, no digital signature/custody chain, no encrypted import, no remote Wazuh connector or production SLO. OPEN-013 is unresolved and product canonical specs are unchanged. Multi-feed Case timeline viewer remains operational as a separate mode.

## Tests and acceptance

Functional test traverses two Wazuh files containing the **same raw alert ID** through normalized in-memory imports, independently authorized two-Case HTTP traversal and chronological proximity selection. It verifies exact feed-specific entry hashes, observed/inferred kinds, known 120ms/125ms separations, no output of source descriptions or owner ACL, and unchanged operator files. Negative tests cover malformed second file, source-file reuse, owner-file reuse, 2..4 bound, foreign/denied/missing Cases and principals, cancelled request, incompatible CLI modes and invalid interval.

CI requirements: Go race/vet, strict gofmt, engineering security coverage, exact-sha workflows and one risk-based **independent, read-only L2** of this stable functional integration before merging. The previous PR #48 HEAD `a0fc37707c1c08cbc95b468849bbbeb86ba779e2` had green Timeline runtime CI but engineering-kernel failure (changed critical coverage **86.76%**). A new commit must not be considered delivered until **both** independent CI paths pass. Keep PR draft.
