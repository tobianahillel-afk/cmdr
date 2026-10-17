# E10-INV-120A-SNAPSHOT — Wazuh 4.x → two-Case temporal leads

## Functional increment

An analyst can now take a locally acquired Wazuh 4.x `alerts.json` JSONL file directly to an already authorized Case and look for near-in-time events in a **second independently authorized Case**, using one local CLI command. No intermediate normalized event or merged owner snapshot is persisted:

```sh
go run ./cmd/case-timeline-demo \
  -principal=analyst-a -snapshot=/absolute/owner.json \
  -wazuh-alerts=/absolute/alerts.json -into-case=case-a \
  -triage-import -with-case=case-c -compare-proximity \
  -proximity-seconds=30 > /private/wazuh-temporal-leads.json
```

Output: a single sensitive, unsigned `cmdr-local-wazuh4-case-proximity-v1` JSON receipt. It includes a maximum of 25 nearest-time leads among **already source-authorized** timeline rows. Each second Case is independently authorized through the existing Case/Timeline/typed-source HTTP handler path. No output is emitted before both reads and the proximity selection succeed.

## Reuse-first choice and boundaries

Selected: Go standard library and existing `normalizeWazuh4Alerts`, `readLocalEventFile`, `importLocalEvents`, `newDemoHandlerWithData`, `collectAuthorizedCrossCaseSequence` and `selectCaseProximity`. A new correlation engine, search index, external parser SDK, IAM service or persistence system would increase complexity for a bounded offline analyst journey. Wazuh `agent` or `rule` data is never an authority to grant access. Input source descriptions are not exported with leads, though authorized entry IDs and timestamps are sensitive.

The old Wazuh `-serve-import` browser journey remains independently supported. Ambiguous viewer/triage combinations fail closed. Source and owner files remain byte-for-byte unchanged.

**Limitations**: all owner permissions and identities are operator-supplied local lab state; no real production IAM, live Wazuh API, signed forensic provenance, chain of custody, database, durable audit or SIEM stream is implemented. A near timestamp is a **lead only**, not causation or confirmed correlation. `OPEN-013` remains unresolved. No performance SLO is claimed.

**Acceptance evidence required**: exact-SHA Go race and functional test CI; one imported Wazuh event at 10:02:02.125Z should yield a 125ms lead against a separately authorized Case event at 10:02:02.250Z. Foreign/denied Cases, unsupported ECS input, cancellation and ambiguous flags must emit zero output. Run and L2 outcomes must be separately verified before any completion claim.
