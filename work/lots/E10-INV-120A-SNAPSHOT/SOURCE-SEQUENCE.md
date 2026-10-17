# E10-INV-120A-SNAPSHOT — authorized source-sequence investigation

## Functional increment

The local analyst already loads a trusted operator-selected owner JSON snapshot and receives authorized Case Timeline entries from the existing Case/Timeline/typed-source HTTP chain. Previously the summary grouped references and opened only ONE representative event from each top source. This increment adds an actionable **25-entry paginated source sequence** from the same bounded, fully authorized Case summary: the analyst selects one of the top three source groups, browses its already returned entry IDs, observed/inferred kind and declared timestamps, then opens a particular entry in the existing source owner preview route.

This is a concrete analyst journey from on-disk owner data → authenticated Case → authorized Timeline pages → grouped source → chosen entry → reauthorized source owner preview. No new API, source permission decision, persistence, live ingestion, provider, index, UI library or backend abstraction is introduced.

## Reuse-first decision

Chosen: reuse the existing authenticated paginated `/case-timeline` summary fetch, its `Map` keyed by exact `source_ref`, and existing `openSource(entryID)` which calls the Case/source navigation and typed owner preview endpoints. Local sequence pagination is 25 events with max 5,000 previously authorized entries. Avoid extra data fetches, extra server authorizers, parallel index, SQL/full-text search, Elasticsearch and a new Shared-owned Timeline service: none is justified for the 8 MiB/5,000-entry demo.

Rows retain the original server-provided sequence order; the UI **does not invent chronology** from uncertain clocks or claim causal correlation. Each entry shows its original declared timestamp and uncertainty marker. Event ID and source labels are rendered via `textContent`, with long IDs abbreviated for display only; the unmodified ID goes to the already authenticated source route, never into an authority-bearing client flag.

## Acceptance and limits

- A complete 65-entry same-source summary pages 1–25, 26–50 and 51–65 **without additional fetches**.
- Opening an event from the third page uses the existing same-origin/no-store bearer-authenticated source navigation and owner preview; only `entry_id` is a selector. Reauthorize Case/Timeline/typed source at click time.
- Changing Case, token or filters erases the sequence and renders stale action closures inert. Incomplete/denied/inconsistent paginated source summaries never publish controls.
- Everything remains loopback-only, operator-loaded **synthetic** Case data. The local summary is nontransactional; previously visible entries can be revoked before click. This does not establish live SIEM ingestion, production auth/ACL, persistent Case store, final UI, chain-of-custody verification, or production performance SLO. OPEN-013 stays open and Shared Timeline ownership is unchanged.

Evidence policy: CI on the published SHA and independent risk-based L2 review must be checked before claiming milestone VERIFIED. See `analyst_summary.test.cjs` for behavior and authorization boundary expectations.
