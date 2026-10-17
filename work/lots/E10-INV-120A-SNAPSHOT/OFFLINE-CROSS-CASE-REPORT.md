# E10-INV-120A-SNAPSHOT — Offline comparison of two authorized Cases

## Functional increment
`case-timeline-demo -principal=analyst-a -snapshot=/absolute/path/owner.json -compare-case=case-a -with-case=case-c`

Produce a **single unsigned local JSON comparison** without starting a listener. A single process-selected operator snapshot and synthetic principal are passed to the existing `NewMemoryCaseTimelineHTTPHandlers` owner chain. `collectDemoCaseTriage` is executed through authenticated Case discovery and paginated `/case-timeline` reads independently for **both** Cases. The receipt is emitted only after both sides succeed and the context remains active. No alternate authorization engine, HTTP route, Case object or permission claim is introduced.

The receipt includes both already-authorized aggregate triage summaries, differences in visible entries and observed/inferred counts, and the intersection of each Case's **known UTC timestamp bounds** when defined. It does **not** claim causality, matching events, exact incident overlap or a true temporal join. Unknown clock entries do not contribute to asserted time overlap. Per-Case `top_sources` is limited to five authorized references, and must not be interpreted as a complete cross-Case source intersection.

## Reuse-first design
- **Selected:** existing Go standard-library `encoding/json` and `time`, existing owner/typed-source ACL HTTP handlers, existing bounded page collector and signed-in demo-local bearer.
- **Deferred:** an additional provider authorization layer, a new export endpoint, external search/indexing/database, a new persistence schema, a general correlation engine, third-party dependencies and new product decisions.
- Security constraints: two distinct Case IDs (up to 1024 bytes), each separately discovered and authorized, at most 5,000 timeline entries each via the existing page size and source ACLs, no output after a denied/failed second Case, local stdout only, no bearer or raw Timeline rows in the receipt, no writes to source data.

## Validation and limits
Go tests use an operator-supplied **on-disk** dataset, two same-tenant readable Cases, a denied linked source, a separate foreign tenant, and **207 entries across two authenticated pages** in Case A. Assertions cover actual permitted observations, observed/inferred differences, exact certain-time interval intersection, source-reference aggregation, no raw/hidden entry IDs or source preview, unauthorized second Case, same/missing Case, invalid path, cancelled context and failed output writer.

This remains a local synthetic loopback integration example. The synthetic principal/owner dataset are NOT enterprise IAM or live SIEM source adapters. Sequential reads of two Cases are **not transactionally snapshot-isolated**. This JSON is not signed, a chain-of-custody document, evidence export, or production investigation record. `OPEN-013` remains unresolved. No production SLO is claimed.

**CI and L2:** after publishing an immutable code SHA, record its exact workflow run IDs and request ONE risk-based independent read-only L2 covering this functional increment; do not treat historical passing checks as validating this change. Do not mark `E10-INV-120A-SNAPSHOT` VERIFIED until its definition of done is met.
