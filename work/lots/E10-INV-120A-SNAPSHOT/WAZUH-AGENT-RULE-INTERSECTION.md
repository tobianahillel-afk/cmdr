# E10-INV-120A — Wazuh agent × rule authorized intersection

## New functional capability

The existing offline Wazuh 4.x analyst journey already imports operator-provided alerts into an owner-authorized Case, shows the top agent and rule cohorts, pages through source-authorized Timeline identities and opens permitted typed-source previews. This increment lets an analyst **combine an agent ID AND a rule ID** without exporting raw owner data or invoking a second indexing service.

- In the authenticated loopback analyst viewer, the user picks an agent and a rule from the **already authorized** cohort summaries, then requests their intersection.
- The existing `GET /demo-wazuh-drilldown` accepts `case_id`, both `agent_id` and `rule_id`, and optional bounded `offset` (25 records per page). Legacy single-agent and single-rule requests remain unchanged.
- The server validates one unique value per selector, rejects unknown query keys/invalid IDs/offsets, authenticates the ephemeral bearer and performs the **complete pre-existing Case/Timeline/typed-source authorization traversal on EVERY page** before any response bytes.
- Filtering combines both Wazuh metadata predicates over the same imported source bytes; neither the request query nor Wazuh agent/rule metadata grants access. Only allowed entry ID, source timestamp, kind and uncertainty are returned, never raw descriptions, host names, protected owner graph or tokens.
- The existing `openSource(entry.id)` path independently reauthorizes the selected source. Client-side cohort status and paging are discarded after Case/token changes or denied responses.

## Reuse-first decision

**Selected:** one extra conjunction inside the existing bounded Go observer (up to 5,000 entries) and existing JS cohort panel, HTTP path and source navigation. These primitives already preserve the trust boundary, support cancellation and have tests. No new dependency, SDK, index, query DSL, persistence or asynchronous worker.

**Deferred:** a database-backed facet/index or dedicated correlation engine might reduce repeated page traversals at larger scale, but would introduce synchronization, ACL-index leakage, retention and authorization complexity without a measured need in this 5,000-entry offline lab.

## Tests and limitations

Go integration tests import real-shaped disk Wazuh JSONL, check single- and dual-key filtering, 25-entry pagination, no duplicate rows, foreign Case/anonymous bearer/duplicate selectors/unknown query rejection and absence of raw source previews. Node analyst tests verify selectable agent+rule controls, request identity, page navigation, unchanged bearer-in-header discipline, revocation clearing and rejection of mismatched server metadata.

This remains an **operator-local synthetic owner/ACL loopback prototype**. It is not production IAM, forensic chain of custody, live Wazuh ingestion or causal/incident correlation. `OPEN-013` remains unresolved. Do not report CI PASS, L2 approval or merge readiness until the exact SHA is checked.

## Resume

After CI, read any actually accessible risk-based L2 report for this material feature, address only evidenced R0-R2 failures, leave R3/R4 as notes, then progress the next operator-visible integration. Historical L2 of paged single-cohort drilldown does not attest this new two-cohort intersection.
