# E10-INV-120A — actionable source correlation in the local analyst journey

## Functional improvement

The authenticated local Case analyst can already load a host-supplied JSON owner snapshot, discover authorized Cases, view a bounded Timeline, navigate reauthorized typed source previews, search and compute a complete 200-per-page summary. This increment **makes that synthesis useful for evidence triage**: it groups the same authorized entries by the exact existing `source_ref`, counts occurrences and their observed/inferred classification, and shows the three most frequent source references. This is a grouping/co-occurrence aid, **not** causal inference or cross-case entity resolution.

## Reuse-first engineering decision

Selected: extend the already tested browser-native `analyzeAuthorizedCase` loop and use `Map` plus a small deterministic sort. Alternative deferred: database-backed joins, Elasticsearch, separate correlation service, cross-tenant index, additional API and new frontend library. The first-party API already filters via the Case/Timeline Entry/typed source ACL providers before returning each page; a new server-side authority or separate HTTP endpoint would be redundant for this bounded offline demo.

Security and behavior: the same `Authorization: Bearer` header and no-store, same-origin read-only requests; the bearer never enters query strings or storage. Duplicate event IDs, changing pagination totals, malformed/missing authorized source references, revocation or user selection changes fail closed with **no partial aggregate**. Only `textContent` is used to render source references, and names over 120 UTF-16 units are visibly abbreviated. Counts cover the current Case and its active query/kind/date filters, not hidden or cross-tenant events. Up to 5,000 entries, up to 25 authorized HTTP pages; `Map` memory bound and local sorting; no production latency baseline, durable storage or export.

Existing tests: Node browser adapter now validates 205-event/2-page grouping, source kind totals, malformed last page, label escaping, and cancellation; Go HTTP integration separately checks disk JSON → Case/Tenant/source authorized pagination with hidden source exclusion. CI must be checked at exact SHA before claiming pass.

The operator-supplied snapshot remains **synthetic local input**, not live forensic ingestion, a production IAM/ACL provider, or the Shared Timeline Engine. OPEN-013 and tenant boundary contracts stay unchanged. This standalone integration milestone is the appropriate risk-based L2 scope, not every future small test edit.
