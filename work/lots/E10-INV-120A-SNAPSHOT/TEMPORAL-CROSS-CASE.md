# E10-INV-120A — temporal triage of shared source references

## Functional increment

The offline analyst viewer already compares **exact source-reference strings** returned by two independently authenticated, Case/Timeline/typed-source authorized HTTP reads. It formerly displayed only per-Case counts. This increment adds the first and last **certain, visible** instants per Case for each shared reference, the number of uncertain observations, and the absolute difference between the first certain instants. Timestamps with explicit RFC3339 UTC offsets are normalized by JavaScript `Date.parse`; uncertain entries never become invented dates. If either Case has no certain event for a reference, the temporal gap is **indeterminate**.

The per-case 5,000-entry ceiling, 200-entry pagination, existing source-reference sort and five-result display limit remain unchanged. Every comparison still completes authorization against both Cases; failure/revocation/inconsistent pages clear results instead of publishing partial information. Each existing pivot button **re-queries** the selected Case through the original server-side permission chain; it does not acquire source-owner authority from correlation.

This is a useful DFIR triage aid, **not** a cross-tenant search index, evidence export engine, proof of causality, transactionally consistent query, live SIEM feed or durable chain-of-custody artifact. The local bearer/ACL snapshot remains synthetic; `OPEN-013` stays open.

## Reuse decision and verification

Reuse the current Go Case/Timeline/source HTTP services, the existing read-only analyst.js comparison and plain JavaScript Maps/Date. Do **not** add an indexer, database, third-party temporal analysis library, new endpoint or permission logic. The per-page work stays linear in the number of authorized visible entries; compare sorting stays on bounded shared reference groups.

Behavioral Node tests include paginated Cases, normalization of explicit offsets, ignored uncertain timestamps, indeterminate time, malformed asserted-certain timestamp, and denied Case. The existing `investigation-timeline-runtime` GitHub workflow is expanded to execute both analyst return and analyst summary/correlation tests, rather than leaving the correlation tests only in the repository. Required CI conclusions must be checked for the exact published HEAD. A single risk-based independent read-only L2 applies to this material functional integration when CI is green. No merge is authorized.
