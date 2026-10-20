# E10-INV-004G-HOST — CAP-INV-004 read-only Search → Inspection HTTP authority

This engineering work unit adds one provider-neutral, first-party Go HTTP
adapter to the **existing** `eventinspection.Project` engine.

## Authority flow

1. The trusted embedding application authenticates `SessionResolver`.
   It derives principal, Tenant and environment from the host, never JSON.
2. `SearchOwner.CanInspectSearchResult` independently re-authorizes that
   principal's selected Event membership in the exact Job/run/context.
3. `TelemetryOwner.ReadAuthorizedInspectionEvent` independently re-authorizes
   the immutable Shared-owned Telemetry Event, raw/rendered rights, the
   environment and all returned source/normalized/enriched field values.
   This owner **must pre-redact** every field the principal cannot view.
4. `eventinspection.Project` checks Tenant, permission and distinct raw /
   rendered decisions. The host marshals only its allowed projection as a
   JSON `kind:inspection` object that the browser journey understands.

POST `/investigate/event-inspection` with JSON only containing:
`tenantRef`, `environmentRef`, `eventRef`, `jobId`, `runId`,
`correlationId`, `sourceHint`. The server checks client scope against the
host session but **does not trust** it. Every access denied and upstream
error returns a generic response without any partial projection. Source,
parser and enrichment fields use lowerCamel JSON for the first-party
`product-runtime/event-search-frontend/inspection.mjs` contract on #183.

## Evidence and strict exclusions

- Go `httptest` exercises the actual handler, session/Search/Telemetry
  owner order, tenant and Event mismatches, raw/read access independence,
  classified raw redaction, provenance, malformed requests and revocation.
- Tested providers are synthetic host stubs; this does **not** supply
  production authentication, data source adapters, IAM, search index,
  actual telemetry, TLS, audit sink, CSRF control or browser routing.
- Real production rights must be resolved at source by independent owner
  services; a frontend Event reference never transfers permissions.
- Raw/rendered HTTP transport currently permits bounded **UTF-8 text only**
  to preserve bytes; binary forensic payload transport is **not provided**.
- No mutation, Case link, Artifact or Evidence creation, parser execution,
  final query dialect or automatic PivotDraft execution while OPEN-013/014
  remain unresolved.
- Independent GitHub PR #181/#182 Go extraction/corrections and #183 browser
  journey remain **unchanged** while under review. This branch is
  based on the earlier source-identical Go foundation #36. It does not
  import the fix from #182; no pivot endpoint is exposed.
- The older engineering foundation cache marks E10-INV-120A READY even though
  later stacked Git branches implemented it. This parallel work adds only
  E10-INV-004G-HOST and retains historical statuses rather than rewriting
  an unrelated graph or claiming all parents are integrated.

Actual GitHub Actions CI must be verified on the exact PR HEAD before
calling this work technically proven. The L2 reviewer must separately
inspect Search/Telemetry trust, per-field data classification and lack
of an actual host provider; a green CI is not production acceptance.

## Architecture-safe relocation

The first implementation was staged as an in-tree subpackage. The enforced CHK-EVENT-INSPECTION-E2E correctly rejected touching a runtime whose pivot evidence was already frozen. The effective final diff therefore removes the proposed subpackage entirely and registers **a distinct first-party HTTP host architecture boundary** (`event-inspection-http-host-runtime`) that depends only on the unchanged `event-inspection-runtime` Go module; the historical module retains exactly its verified root. The new on-demand performance target deliberately measures the existing upstream projector only, NOT HTTP or remote owner cost; the separate `BenchmarkAuthorizedInspectionHTTP` reports local synthetic HTTP mean ns/op without a production SLO. Original E10-INV-004B/C evidence and source blobs remain unchanged. No check threshold or accepted SHA was rewritten.

## Mandatory measured security scope

The new runtime boundary is registered in `engineering/security/runtime-security-policy.json` with server-side authorization and Tenant isolation mandatory. `tools/cmdr-dev/security_tests.go` uses the new Go module to measure coverage against the **unchanged >=90% changed-security-critical floor**, and runs `TestDeniedSearchNeverCallsTelemetryAndDenialsAreGeneric` plus `TestSessionScopeIsNeverTakenFromClient` as separate negative gates. The module's own exact-HEAD Linux CI reported 91.8% host coverage, but only a successful global engineering-kernel run proves the registered multi-runtime gate. No gate was lowered or disabled.
