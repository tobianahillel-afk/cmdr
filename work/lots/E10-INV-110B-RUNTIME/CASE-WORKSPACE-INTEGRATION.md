# E10-INV-110B — Case-to-Timeline functional integration

## Decision: reuse the existing owner and authorization pipeline

The reusable Go first-party boundaries `AccessChecker`, `EntryResolver`, `SourceOwnerResolver` and `ProjectFromOwnerSnapshotsWithTypedSources` already implement the security-aware read-only projection from selected entry IDs. Previously a caller had to provide these IDs: there was no Case -> Timeline discovery path.

Alternatives considered:
- **Selected:** introduce only `CaseTimelineCatalog.ListCaseEntryIDs(ctx, Case-scoped AccessRequest)`, then reuse the existing typed owner projection, with a renewed Case permission check. No new runtime-to-runtime import edge or dependency.
- Build a new Timeline Engine, database/index, general ACL broker, or multi-service orchestrator: rejected for this bounded step; Shared owns those responsibilities and `OPEN-013` has not been decided.
- Expose caller-supplied entry IDs or trust Case Queue handoff/URL/permission flags as an authority: rejected, because tenant, Case and entry IDs can be forged by a caller.

## Functional path implemented

`OpenCaseTimeline(ctx, scope, caseID, checker, catalog, entries, sources)`:
1. Validates context, Case, tenant/environment and required server integration providers.
2. Requests authoritative Case read permission BEFORE calling the trusted catalog (no enumeration on denied Case).
3. Gets the selected Case's Timeline Entry IDs from its owner-side catalog. These IDs are not accepted from a browser or user-provided snapshot.
4. Delegates to the existing typed provider pipeline, which rechecks Case permission (to handle revocation after discovery), checks each Timeline Entry, resolves canonical entry snapshots, resolves four supported typed source families, enforces each source's owner-specific permission, and returns a deterministic immutable Timeline projection.
5. Any provider failure, malformed/oversized catalog, missing authority or cancellation fails closed. Hidden linked source content is not copied into the result.

The new in-module integration test exercises the **entire in-process functional chain** against one deterministic server-owned **in-memory fixture**, including two tenants, two Cases, mixed inferred/observed chronology, version/supersession, hidden evidence, wrong tenant/environment, catalog failure, cancellation and authorization revocation. This fixture is NOT a deployed owner service or durable repository.

## Evidence scope and delivery blockers

- Implemented: usable read-only Go application entry point joining Case selection, owner ID discovery, typed permission checking and projection in one bounded in-process call.
- Not implemented: real principal-authentication integration, persistent Case/Shared Timeline Entry catalog, actual Evidence/Telemetry/Govern/Command owner-provider adapters, transaction/snapshot consistency, HTTP/API gateway, Case Queue UI handoff to Timeline, final navigation/screen `INV-TIM-001`, replay/export or Shared Timeline Engine.
- The `AccessChecker` must be bound to an authenticated server session by the deployment, never to caller-provided flags; `CaseTimelineCatalog` and all resolvers must verify same principal, Case, tenant and environment. The two Case authorization checks are deliberate at separate trust boundaries, not interchangeable.
- No claim of production readiness, p95 server response, tenant isolation from deployed stores or complete CAP-INV-110 delivery. Historical performance measurements remain tied only to their own SHAs.
- Existing work graph `E10-INV-110B-RUNTIME=READY` and `E10-INV-110C-NAV=BLOCKED` are retained until complete provider integration and independent review/DoD evidence exists. Canonical Product Spec and `OPEN-013` unchanged.

## Next functional work

1. Bind this orchestration to real first-party Case/Shared owner repositories and authenticated access checks (do not fabricate a shared schema or widen tenant privileges). Review ACL revocation/version consistency at that interface.
2. When owner integration is available, run real repository-backed end-to-end tests with denied/foreign/missing source cases, trace context and performance measurements; qualify B separately from production deployment.
3. Only after gates truly hold, progress `E10-INV-110C-NAV` to navigation/UI using the verified Case-to-Timeline service. This decision avoids repeating micro-validation changes that do not increase usable functionality.
