# E10-CMD-106A-CASELINKS — authorized existing Investigate Case links

## Canonical ownership and four distinct parents

**Canonical:** CAP-CMD-106 Incident Coordination §§1–27, Command / incidents-and-work-queue. §9 permits *reading* independently owned Investigate Case status and linked object relationships; §8 requires source-authoritative cross-product links; §15 requires honest Partial/Permission-denied states. The incident remains Command-owned, the Case remains Investigate-owned, and the link remains Shared Object Linking Service owned (canonical objects Incident and Case). CAP-CMD-101 provides the selected current-page Incident in the six-system-view Work Queue. No new status lifecycle or synthetic Incident creation is claimed.

**Engineering:** epic E10-CMD-106 and executable task E10-CMD-106A-CASELINKS depending on already IMPLEMENTED E10-CMD-101I-SELECTION. Work graph formerly had no READY entries; this is a newly decomposed, bounded Class-0 obligation, **not** a claim that an old READY item existed. **Code:** native Go `product-runtime/work-queue`. **Git parent:** frozen PR #176 SHA `415aba77758540673a4c78fd8f7694e300ac1a04` (solely history; **not** the engineering parent). All #2/#36/Command parent PRs remain DRAFT/unmerged.

## Canonical obligation → existing behavior → new behavior → remaining

| Requirement | Before | This bounded lot | Still absent |
|---|---|---|---|
| Show Case links on authorized Incident (CAP-CMD-106 §9/16) | Work Queue already selects and authorizes one Incident, but does not display external Case link summaries | Optional first-party `LinkedCasesOwner.ReadAuthorizedLinkedCases` reads up to 16 source-authorized *existing* Case IDs/status/source/version for one selected current-page Incident; selected JSON and escaped native HTML | Actual host Shared Linking/Investigate provider, authentic Case navigation and persisted bidirectional links |
| Source and Case authorization (Incident/Case invariants, §21) | Command Incident CanRead only | Separate host contract MUST authorize **each** link relation and **each** Investigate-owned Case on exact principal/tenant/environment/Incident ID; never call for unselected/foreign/denied/Task | Real production IAM and owner/permission adapters |
| Degraded state (§15/22) | Missing links not represented | Nil provider/projection means unknown source; authorized empty list means no visible links; malformed claims and refusal return generic HTTP 503 before selected row or links serialized | Full UI partial-source diagnostics and navigation errors |
| Mutating/creating/associating Case (§12/17/25) | None in this projection | **Explicitly out of scope** | Class-2 mutation and Incident→Investigate transition are blocked until accepted decisions and actual source providers; OPEN-013 not resolved |

## Security and implementation

- ONE established `GET /command/work-queue` handler. Normal full owner list, per-object `CanRead`, status/owner filters, pagination and current-page selection execute before the optional external source call.
- Only selected `Incident` triggers **at most one** link adapter call. A Task, any denied/nonexistent/off-page Incident, and an ordinary unselected list never invokes it. The client cannot choose a Case source, tenant or principal.
- This external source is **not** a fabricated local authorization engine. The host implementation is responsible for checking source link visibility and separately Investigate Case read permission for every returned ref. The CMDR consumer validates bounded distinct IDs, canonical Case status, non-empty link provenance/version syntax and UTF-8; invalid or rejected claims fail closed generically without partial selected-row data.
- No raw Case contents, Finding, Evidence, tenant identity, direct Case URLs or unverified source status are loaded or emitted. Native HTML uses Go `html/template` escaping; link metadata cannot execute HTML/script. A missing link provider does NOT imply that no Cases exist.
- Source status values are from canonical Case object states, not inferred from Incident or locally advanced. No new runtime dependency, API fee, route, background job, datastore or Class-2 action.

**Tests:** selected Incident current scope once, JSON multi-Case link summary, unselected/Task/denied/foreign no call, missing owner unknown, HTML escapes source text, denied/duplicate/invalid source generic 503 with no partial Incident/Case details. Existing Work Queue, SLA, ownership, six views and cursor tests must stay green. Three-OS Go `command-work-queue-runtime` and `engineering-kernel` mandatory gates including changed-critical >=90% coverage must pass exact HEAD. One independent read-only L2 for material cross-product boundary, not per commit. A native GitHub approval, source adapter, acceptance and merge have **not** occurred.
