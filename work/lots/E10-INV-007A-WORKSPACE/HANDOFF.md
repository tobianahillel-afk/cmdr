# CMDR — CAP-INV-007 Class-0 authorized Search Result Organization: implementation handoff

## Authority and ownership
- **Product owner:** `CAP-INV-007`, Investigate / Signals and Hunt, `cmdr-product-spec/07-investigate/modules/signals-and-hunt/capabilities/search-result-organization.md`, full canonical specification. Dependencies: `CAP-INV-002` Event Search (`Search Job` Shared-owned) and `CAP-INV-004` Event Inspector (`Telemetry Event` Shared-owned). Requirement refs `REQ-PROD-014`, `REQ-PROD-020`.
- **Engineering owner:** independent bounded root lot `E10-INV-007A-WORKSPACE` in `work/graph.json`, dependent on VERIFIED `E10-INV-002G` and `E10-INV-110E`; manifest `work/lots/E10-INV-007A-WORKSPACE/manifest.json`. **Not** a child of the historical `E10-INV-120A-SNAPSHOT` or transversally scoped DFIR release epic. Product CAP ownership distinct from Git ancestry.
- **Temporary code placement:** `product-runtime/investigation-timeline/cmd/case-timeline-demo/analyst.{html,js}` because the existing authenticated Case/typed-source/first-party SearchJob viewer resides here. This is NOT a new production module nor definitive architecture ownership. No new API, IAM, persistence service or external package.
- **Git parent:** independent child branch of the immutable PR #127 final SHA `7a36424d5140f784542e30f290f45462b7a28a65` (PR #127 base itself PR #120, not `main`). Never edit frozen #127 or ancestors. All still subject to native reviews and ancestor CI. PR remains DRAFT/unmerged.

## Requirement → implemented behavior → remaining obligation
| Canonical obligation | Actual implementable slice | Still missing / blocked |
|---|---|---|
| Selection of authorized Search results | Client-only checkboxes for up to 200 references of ONE completed authenticated owner-authorized Search Job; counts and selected state | Persisted/versioned workspace state, reconciliation after run expiry |
| Grouped results | Group counts into actual Wazuh vs generic JSONL refs using only membership in source-linked `search_hits` that belong to the SAME completed job | More general grouping, enrichment/Entity comparisons and Saved View integration |
| Pinning | Explicit transient pinning of selected refs, current-job only, no mutation/export | Persistence of pin reasons, annotations, audit trails |
| Manual inspection | Single-selected authorized Entry ID launches EXISTING `openSource` and independent server owner Case/Timeline/typed-source checks | Case relation/linking and selected package, pending class-2 authority |
| Security/partial state | No result workspace until full completed Search Job and correlation leads validated. Duplicate/foreign Wazuh metadata fails closed. New search/filter clears pins/selection; stale controls inert on Case, bearer or criteria mismatch | Production tenant IAM, expiry/stale endpoint, signed evidence custody |

**Explicit exclusions:** `OPEN-013` class-2 mutation approval not decided. No annotation, persisted exclusion, Case linking, Shared Export Job, Evidence qualification, raw content export, new canonical Collection or claim `CAP-INV-007` accepted/complete. All current results are synthetic operator-local sources. No local storage/cookie/token saving and no new access authorization.

## Verification and deterministic next action
Browser integration tests in `analyst_summary.test.cjs` exercise real completed SearchJob envelope and Wazuh+JSONL source grouping, selection, pin lifecycle, source-intent controls, stale clearing and foreign grouping failure. Existing Go+Node multi-OS actions must run on *exact final SHA*, with `engineering-kernel` unchanged mandatory changed-security coverage floor 90%. The current graph status is `IMPLEMENTED`, **not `VERIFIED` or product accepted**. Do not upgrade status until exact CI and independent risk-based L2 plus applicable native reviews.

After this bounded class-0 slice, the next owner should choose a canonical work unit from the refreshed graph and capability register. Do not keep expanding the demo with class-2 exclusions/export without accepted permission/storage authority. For CAP-INV-007 the likely future architectural task is versioned Investigate-owned workspace plus Shared-owned Search Job reference and permission revalidation, but only after canonical OPEN decisions and a separately approved design.
