# Investigation Timeline — bounded read-only projection core

## Read-only HTTP Case entry point (bounded integration)

`NewCaseTimelineHTTPHandler` adapts `OpenCaseTimeline` to `net/http`. A server-provided `CaseTimelineSessionResolver` must authenticate each request and attach the authenticated principal to the returned request-derived context. The only accepted user selector is `case_id` on a bodyless GET; tenant, environment, Case, Timeline Entry and linked-source permissions are derived/checked server-side using the same Case/catalog/owner/typed-source chain. Results are non-cacheable JSON; authentication failures are generic 401, malformed selectors 400, and Case/owner denials all receive the same generic 404. No frontend-supplied access flags, cookies or `X-Tenant` headers become authority.

The route is deliberately **not assigned to a public URL in this module**; the host API gateway owns routing, TLS, authentication middleware, audit, rate limiting and deployment. The injected in-memory two-tenant owner fixture in `case_http_test.go` exercises a request-to-owner-to-ACL-to-response chain, including denied/hidden sources, revocation and metadata confidentiality. These tests are NOT a deployed authenticated principal/owner repository, a real persistent Timeline catalog, an actual browser integration, a production endpoint, or evidence of production SLOs. See `work/lots/E10-INV-110B-RUNTIME/HTTP-CASE-INTEGRATION.md`.


## Case-to-Timeline integration (bounded functional increment)

`OpenCaseTimeline` now exposes a read-only in-process service from an authorized Case selection to trusted catalog-discovered Timeline Entry IDs, first-party owner snapshots and typed source ACL checks. It reuses existing `ProjectFromOwnerSnapshotsWithTypedSources`; the server `CaseTimelineCatalog` and all other providers remain injected interfaces. Integration tests exercise a deterministic two-tenant in-memory owner fixture, **not a real deployed authenticated/persistent provider**. See `work/lots/E10-INV-110B-RUNTIME/CASE-WORKSPACE-INTEGRATION.md` for scope, acceptance evidence and outstanding gates.

This module now contains an **initial first-party Go projection** for CAP-INV-110, not a complete Timeline Engine, source authorization service, backend provider or production-ready UI.

`Project` accepts caller-resolved Case/Shared Timeline Entry snapshots, enforces explicit tenant/environment + Case and Timeline Entry permission outcomes, masks denied/cross-tenant sources, preserves observed/inferred kind, source references, provenance, version and supersession, and sorts known timestamps deterministically. Unknown or clock-uncertain time does not become an observed chronology. Inputs are not mutated. The caller MUST obtain authoritative server-side permission decisions separately; caller-provided permission flags are **not** proof of production authorization.

Executable tests cover the 18 predeclared synthetic contract scenarios and focused negative cases. Limits (5,000 entries and 8 MiB string input budget) are provisional engineering protection, not a canonical limit or production SLO. No persistence, Shared Timeline Entry mutation, replay/export, external runtime dependencies, network or final UI exists here. `OPEN-013` is unresolved.

**Projection performance evidence exists only for synthetic on-demand runs**. GitHub Actions run `38842748892` at exact source SHA `dae5fe2c64e1eb1befb0cc0130cc42040ba5bfdf` returned `pass`, measuring a 4.70007566 ms p95 across five 200-operation sample means and 5,147,088 bytes of periodically sampled Go heap over 1,000 iterations. Neither value represents request-level p95, peak RSS, external permission-provider cost or a production SLO. Navigation benchmark and navigation resource handlers remain fail-closed preimplementation. Security coverage and module-specific race/vet are exercised by CI, but no production deployment or full CAP-INV-110 completion is claimed.

Security hardening: `Supersedes` is exposed only when the referenced prior version is also visible and authorized in the bounded snapshot; unavailable/denied predecessor IDs are redacted and status becomes `partial`. A timezone's explicit RFC3339 offset must match the timestamp offset before chronology is considered known. Unknown/named zones conservatively remain uncertain. These are projection policies, not authorization-provider proofs.

Security hardening: trim scope before wildcard checks; uniqueness is enforced among accessible source entries only, preventing denied/cross-tenant ID collisions from denying legitimate projections. Visible supersession links require strictly increasing versions (so cycles cannot be represented); inaccessible predecessor IDs remain redacted. This bounded, caller-supplied projection still is **not** an authoritative authorization provider.

Visibility-first validation: the 8 MiB string-data budget and duplicate/nonempty-ID checks apply to **authorized visible source entries only**. Hidden/foreign/missing source fields must not influence denial; the global 5,000-entry count bound still applies to all supplied entries for bounded runtime. These tests establish in-process noninterference under caller-supplied permission outcomes, **not** independent authority to decide access.

## Runnable Case owner integration (memory-backed reference)

`NewMemoryCaseTimelineHTTPHandler(authenticate, data)` now wires a real, read-only Go HTTP handler from a **trusted, static in-memory owner dataset**. The host supplies a verified authentication callback returning a principal identifier; the adapter derives tenant and environment from its own principal registry. It copies data at construction and implements the existing Case catalog, authoritative-in-this-fixture Case/Timeline/source permission checks, Timeline Entry owner resolver and typed source owner resolver. This is a functional **process-local reference backend**, not a durable or independently verified production authentication/ACL provider.

The HTTP behavior runs without a new database, external API or frontend trust flag. A full two-tenant, two-Case integration test covers successful Case navigation, observed and inferred entries, provenance and supersession, restricted evidence, shared/foreign tenant isolation, incorrect/missing authentication and invalid catalog/ownership construction. This intentionally does NOT add a persistence schema or change Shared-owned Timeline Entry contracts. See `work/lots/E10-INV-110B-RUNTIME/MEMORY-OWNER-INTEGRATION.md`.

For a real deployment, the host must supply its identity verifier, actual Case/Shared/Evidence/Telemetry/Govern/Command repositories and provider-side transaction/revocation semantics. No remote provider, actual production Case repository, installed route, database migration, production SLO or final frontend is claimed.

## Functional analyst viewer (local demonstration only)

`go run ./cmd/case-timeline-demo -principal=analyst-a` now serves **`http://127.0.0.1:<ephemeral-port>/analyst`** alongside the existing `/case-timeline` read-only API. The viewer is a small first-party HTML/CSS/JavaScript page embedded into the same Go binary: enter the terminal-provided **ephemeral** bearer manually, select a Case ID (`case-a` for analyst-a or `case-b` for analyst-b), optionally filter by observed/inferred kind, and page through the **authorized** timeline. Entry identity, timestamps, uncertainty, source provenance, version and accessible supersession are rendered as text, not HTML. The browser uses same-origin `fetch` with an `Authorization` header and no credential persistence or third-party assets. Exact static routes enforce no-store, a restrictive CSP, referrer/frame/MIME protections; the API continues to enforce all Case, Timeline and typed-source permissions before filters or pagination. Anonymous browser access to the static UI never grants data access.

This is **an interactive, loopback-only synthetic product demonstration**, not the final CMDR screen, a deployed UI, remote identity system, owner repository, source ingestion pipeline, audit backend, production storage, TLS ingress or production SLA. Page counts cover visible entries only; sequential pages are not snapshot-isolated under concurrent owner changes. `OPEN-013` is unchanged. Relevant tests: `go test ./cmd/case-timeline-demo` and the existing package HTTP/owner integration suite.

## Authenticated Case discovery for the local analyst demo

The loopback-only `/demo-cases` GET route exposes only IDs of Cases readable by the bearer-authenticated demo principal according to the same immutable synthetic principal/tenant/Case reader data used to build the Timeline owner handler. It rejects arbitrary selectors, write methods and missing credentials, returns no source or Timeline Entry metadata, and never trusts browser-supplied tenant identifiers. The analyst can discover a permitted Case in the UI before opening its timeline, retaining server-side Case/Timeline/source checks when reading events.

This is deliberately a **fixture-only Case picker** owned by `cmd/case-timeline-demo`, not a general Case-catalog API, production Case search or production identity provider. In a deployment, a real host Case-owner catalogue would implement authenticated Case discovery using its own authorization model; the current module must not expose synthetic fixture data as real Cases. No persistence, schemas or ownership decisions were changed.

## Offline Case snapshot import — functional local investigation

The existing first-party Go Case/Timeline owner adapter can now read an **explicit local JSON owner snapshot** at startup, instead of requiring the hard-coded sample. This enables an analyst to open, discover, search, time-filter and page Cases and Timeline Entries from an existing offline file using the same Case + Timeline Entry + typed-source access checks and same-origin UI. No new database, generic importer, frontend flags, authorization middleware or dependencies were introduced.

From `product-runtime/investigation-timeline`:

```sh
go run ./cmd/case-timeline-demo -principal=investigator-42 -dataset=cmd/case-timeline-demo/testdata/example-owners.json
```

Open the **loopback-only** `/analyst` address printed by the program and enter its one-run bearer token. The JSON uses the existing `MemoryTimelineData` shape: `Principals`, `Cases`, `Entries`, and `Sources` (see `testdata/example-owners.json`). The selected principal and permitted Case list now derive from the loaded local owner snapshot instead of fixed demo Case IDs. The file is limited to 16 MiB, decoded with unknown-field rejection and no trailing JSON payload, and passed through the existing owner/ACL constructor. It is read once; changing a file does not silently update authorization state.

**Trust boundary:** the file and its ACL readers are supplied by a trusted local operator, not uploaded by a remote user. If that operator supplies untrustworthy ACL data, it has no independently trustworthy identity/permission authority. Protect any private source references, use only access-controlled local files, and do not commit real incident evidence or secrets. The session token is still intentionally printed to the launching terminal and should be used only on the analyst's machine. This is an offline reference integration, **not** a production Case repository, source ingestion, database-backed snapshot isolation, real SSO, durable audit, real-time revocation or production network service. Canonical CAP-INV-110 and OPEN-013 remain unchanged.
