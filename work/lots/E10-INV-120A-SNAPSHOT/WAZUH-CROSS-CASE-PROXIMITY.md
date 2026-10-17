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

## Exact GitHub validation and deterministic handoff (2026-10-17)

Functional implementation and tests are frozen at `83061203d79da148d43da9ecc3b7b002002f4f0d`, relative to prior validated source `8ace16423fec2ab2fcee18084a7a9152a4930be5`. Verified GitHub Actions **SUCCESS** at this exact SHA:

- Timeline runtime push: [39013321428](https://github.com/tobianahillel-afk/cmdr/actions/runs/39013321428), Linux runtime plus Windows and macOS subprocess matrix jobs.
- Timeline runtime PR: [39013323995](https://github.com/tobianahillel-afk/cmdr/actions/runs/39013323995), same three jobs.
- Engineering kernel push: [39013321433](https://github.com/tobianahillel-afk/cmdr/actions/runs/39013321433).
- Engineering kernel PR: [39013324014](https://github.com/tobianahillel-afk/cmdr/actions/runs/39013324014).

One **independent read-only** L2 task `CMDR — Revue L2 — E10-INV-120A-WAZUH-PROXIMITY` was created for this exact functional SHA and diff. **L2_EN_ATTENTE**: neither a detailed L2 report nor native GitHub reviewer approval has been verified; do not claim either.

**Next functional step** after reviewing materially evidenced R0/R1/R2: move this bounded offline Wazuh-to-Case analyst bridge toward a real provider-owned IAM and durable ingestion/evidence boundary, only after their canonical authority and trust choices are settled. Avoid more microtest gates for already green local fixtures. The work graph still marks the broader snapshot lot READY, not production-shippable.
