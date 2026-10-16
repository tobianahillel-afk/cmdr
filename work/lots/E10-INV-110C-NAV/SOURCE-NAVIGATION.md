# E10-INV-110C-NAV — read-only source destination from Case Timeline

**Status:** first bounded functional increment, NOT final C completion, pending exact-SHA CI and independent L2.

## Functional walkthrough
1. Host authenticates an analyst and exposes Case discovery + Timeline using already injected first-party interfaces.
2. Analyst selects a visible Case, filters and views authorized Timeline Entries.
3. Analyst clicks **Vérifier et ouvrir la source** for a specific permitted entry.
4. The source-navigation endpoint discards caller ACL/source claims; it reuses Case catalog, Timeline Entry owner and typed source permission checks, re-resolves canonical source identity, checks its owner ACL, and rechecks Case/Entry access before returning a typed destination reference and Case/entry/time return anchor.
5. Denied, cross-tenant, missing, malformed and revoked object routes do not disclose identities. The UI renders the *authorized draft* as text; it does not forge an external URL or perform source-owner screen navigation.

## Technical choice
- Reuse existing `OpenCaseTimeline`, `SourceOwnerResolver`, `AccessChecker`, `NewCaseInvestigationHTTPHandlers`, `MemoryTimelineData` and Go `net/http`. The first simple implementation reopens the full bounded 5,000-entry Case Timeline, preserving existing ACL/semantics without a duplicate permission engine.
- Do NOT add a vendor database, search/index engine, generic router, alternate RBAC service, unsigned redirect URL or browser-created permission flags. Case Queue's read-only ReturnState design remains the semantic model but its base64 return token is not imported or treated as a credential.
- This runtime owns no Shared Timeline Entry data, no cross-product object persistence, no actual Event/Evidence/Decision/Incident screen, no final UI, and no source content copy. `OPEN-013` remains open.
- Host responsibilities: verified identity, real owner repositories, audit, TLS, rate limiting, destination-owned routes and authorization/revocation checks. Multi-call ACL reads are not an atomic permission transaction: host must supply consistency appropriate to sensitive data.

## Tests and acceptance
- Successful host-level source navigation for permitted telemetry and Evidence, authorized Case return anchor with preserved entry version/kind/time.
- Generic denial for source-hidden, foreign-tenant, missing entry, unauthenticated and revoked Case/source; forbidden query claims and non-read methods rejected.
- Source owner identity is only taken from provider, not URL/header/SourceRef prefix.
- No in-browser storage or executable HTML from source reference.

A separate **risk-based L2** should review this significant cross-owner read navigation milestone after an actual GitHub commit and CI verdict. No production SLO and no external owner connector are claimed.
