# E10-INV-110B — provider-neutral analyst Case discovery to Timeline HTTP integration

## Functional increment

The first-party read-only Go module now exports `NewCaseInvestigationHTTPHandlers`. An existing CMDR host can inject its own authenticated `CaseTimelineSessionResolver`, authoritative `AccessChecker`, Case candidate owner `CaseListCatalog`, Case's Timeline Entry catalog `CaseTimelineCatalog`, Timeline owner `EntryResolver`, and typed source `SourceOwnerResolver`. Both resulting `Cases` and `Timeline` `http.Handler` values reuse the **same session and Case ACL services** without a new security engine or provider implementation.

The existing memory-backed runnable demo calls this same constructor, so the reference Case picker and investigation view exercise a shared wiring path instead of two separately built flows. Real integration tests invoke both routes through `http.ServeMux` with independently injected owner interfaces, run an authenticated Case selection and authorized kind/query/pagination filter, and check missing sessions, foreign Case denial, hidden source confidentiality, cross-tenant scope spoof rejection, and access recheck after Case revocation. This is the functional path `authenticate → discover Cases → select authorized Case → check Timeline Entry/source → read-only forensic chronology`.

## Alternatives and boundary

- **Reuse selected:** existing `NewCaseDiscoveryHTTPHandler`, `NewCaseTimelineHTTPHandler`, memory demo providers and standard-library `net/http`. This is a small composition seam rather than a second implementation of authorization.
- **Not selected:** proprietary or third-party authorization engine, database, indexer, event-store, new generic dependency-injection framework, and a second independent Case owner. Those choices require specific trusted-base/security and persistence architecture decisions.
- **Retained host obligations:** the actual host must supply verified identity, tenant-scoped Case and Timeline owner/catalog, source permission and data providers, TLS, routes, rate limits, audit, revocation consistency, persistence and lifecycle. The test store is intentionally synthetic; no backend is deployed and no real customer data processed.

The read-only scope of B and Shared Timeline Entry ownership remain unchanged, with `OPEN-013` unresolved. The later C-NAV lot is still blocked in the work graph. This milestone is suitable for one independent L2 review after its GitHub Actions checks are confirmed; no extra L2 per cosmetic correction.
