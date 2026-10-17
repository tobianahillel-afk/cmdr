# E10-INV-110B — runnable in-memory Case Timeline owner integration

## Functional deliverable

The existing `OpenCaseTimeline` and `NewCaseTimelineHTTPHandler` already implement the read-only flow, but previously required a custom fixture or independently deployed Case/Timeline/Source providers. `NewMemoryCaseTimelineHTTPHandler` is a **concrete, reusable process-local provider** to exercise the complete HTTP Case-to-Timeline journey.

One immutable constructor-owned dataset ties verified principal IDs to tenant/environment, Cases, Timeline Entry snapshots and typed source owner objects. A host-supplied authentication function returns only verified principal IDs; HTTP headers and frontend-provided permission flags never become authorization evidence. The adapter handles Case membership, Timeline Entry ACL and source-object ACL separately and delegates chronology/provenance/supersession to the existing standard-library Go projection engine.

### Reuse-first comparison

- **Selected:** implement minimal first-party, read-only in-memory adapters for existing `CaseTimelineSessionResolver`, `AccessChecker`, `CaseTimelineCatalog`, `EntryResolver` and `SourceOwnerResolver`. Uses no third-party packages or new runtime dependency; preserves B scope.
- **Deferred:** PostgreSQL/SQLite durable Case/Timeline owner storage, external IAM/ACL providers and atomic revocation transactions — require real service ownership and deployment architecture, and Shared Timeline data model remains unsettled.
- **Rejected:** a second Timeline Engine, generic RBAC library, client-supplied permission/tenant fields or a separate bespoke HTTP authorization pathway.

### Evidence acceptance and limitations

Tests run HTTP to real in-memory owner records for two tenants and Cases, verified principal-to-scope mapping, source-filtered evidence and complete chronological display; validate cross-tenant 404, absent session 401, spoofed selectors 400 and invalid trusted records fail before serving. Constructor snapshots are copied so upstream mutations cannot enlarge permissions. No mutation, replay, export or persistence is introduced.

This **does improve real process-local usability** but is NOT production-ready: authentication still comes from an injected trusted host callback, the backing data is not persistent, external source-owner providers are not integrated and transaction/revocation semantics are not guaranteed. Work graph B remains READY until the true owner services, closure criteria, exact-SHA CI and independent L2 are verified. C-NAV remains BLOCKED. Do not count this as complete CAP-INV-110 or OPEN-013 resolution.
