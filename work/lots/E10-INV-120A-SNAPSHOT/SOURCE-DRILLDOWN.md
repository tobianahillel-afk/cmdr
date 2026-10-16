# E10-INV-120A — authorized source drilldown from analyst correlation

## Delivered functional path

The bounded, operator-local synthetic Case analyst already supports authenticated Case discovery, pagination, cross-page correlation and typed source previews. Source correlation previously produced a **text-only summary**. This increment makes the three leading authorized reference groups actionable: each control identifies a representative **authorized Timeline Entry** from the corresponding group, then invokes the **existing** source-navigation and owner-preview routes. Those routes independently recheck current Case, Timeline Entry and typed-source permissions. A revoked source fails closed; no source identifier is accepted as a browser authority.

**Reuse-first choice:** retain the current first-party `Map` grouping, browser DOM text APIs and `openSource` HTTP authorization flow. Do not create a second source-authorization API, database, search index, frontend framework or unreviewed owner interface. Rejected alternatives: passing `source_ref` as authority in a URL, using the full-text `q` parameter as an inaccurate exact-source filter, and assuming an old summary is proof that the entry remains authorized.

## Acceptance and limitations

- Buttons appear **only after all authorized pages have passed completeness, identity and revocation checks**. Source group count reflects the selected Case, filters and best-effort nontransactional pages; grouping is not causal proof.
- A representative entry ID is passed only to the existing server-revalidated navigation and source-preview routes; no bearer token or owner reference is placed in any URL. Text uses `textContent`; reference labels are capped at 120 characters.
- Any Case, filter or token change invalidates and clears the controls. Source revocation or navigation failure hides detail and disables the return link. Node behavioral tests cover a 205-event, two-page path, a source with a representative entry on the **second page**, server denial, and invalidation.
- This is an offline **synthetic loopback** demonstration with an operator-selected read-only JSON file, not real SIEM ingestion, durable persistence, production IAM, audit backend, p95 SLO, or full CAP-INV-110. OPEN-013 remains open. Production integration requires externally trusted principal/Case/Shared/source providers and ownership decisions.

CI and independent risk-based L2 must be checked against the published immutable SHA. No native GitHub review is inferred from green tests.
