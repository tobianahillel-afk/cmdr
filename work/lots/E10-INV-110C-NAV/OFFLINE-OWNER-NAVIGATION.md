# E10-INV-110C-NAV — imported Case owner snapshot → source navigation (functional integration)

## Why this is a functional increment
The active development branch already exposes authenticated Case discovery, paginated Timeline reads, typed source-owner navigation and a safe read-only return selection, but `case-timeline-demo` still only launches built-in synthetic fixtures. Historical draft PR #38 implements a useful bounded local JSON owner snapshot loader but is based on an outdated branch; draft PR #39 duplicates the navigation draft already in the current base and cannot be merged cleanly. This focused PR **reuses PR #38's existing bounded decoder/tests/example**, ports only the CLI integration and wires it into the **current** `NewMemoryCaseTimelineHTTPHandlers`, preserving the current C-NAV `Navigate` route and shared authority.

## One working local analyst journey
1. The operator prepares a **trusted local** Case/Timeline/source owner snapshot with explicit reader lists. A synthetic example is included under `cmd/case-timeline-demo/testdata/example-owners.json`. Never use live private evidence in GitHub tests.
2. From `product-runtime/investigation-timeline`, run:
   `go run ./cmd/case-timeline-demo -principal=investigator-42 -dataset=cmd/case-timeline-demo/testdata/example-owners.json`.
3. The CLI fully loads and validates the snapshot **before** printing a bearer/listener, selects the first readable Case and starts a listener bound to IPv4 loopback only.
4. The analyst opens `/analyst`, enters the one-run bearer locally, discovers the Cases they can read and searches a chronological Timeline. Case, Timeline Entry and canonical source read checks are evaluated by the existing in-memory first-party owner.
5. Choosing an authorized entry uses `GET /case-source-navigation`. The current C-NAV implementation rechecks Case/Entry/source-owner permissions and returns only typed object identity and an inert return selection. Return navigation without a new bearer does not grant read.
6. Denied sources, foreign Cases, manipulated tenant/object parameters and unauthenticated requests do not disclose private source identities. The server never redirects to a caller-supplied URL.

## Reuse and alternatives
Selected: Go standard-library bounded `encoding/json` reader and existing `MemoryTimelineData`, `NewMemoryCaseTimelineHTTPHandlers`, first-party session/checker/catalog/owner contracts, and current analyst viewer. Reuses historical PR #38 blobs, **not** a parallel auth or source navigation subsystem.

Not selected: another dynamic routing engine, SQLite, a source-specific parser zoo, external auth library, generic data-ingest service, production database or paid API. Those solve a materially larger and as-yet-unspecified persistence and identity problem.

## Exact scope and risk
- Input comes only from an **operator-selected local file**, not HTTP. Max 16 MiB, regular file, exactly one JSON document, strict known JSON fields and nonempty owner sections. Existing owner constructor verifies tenant/environment and reader bindings, foreign pointers and canonical source links.
- Read-once snapshot: no live revocation, durable audit trail, encryption-at-rest, real SSO, trusted producer signature, transactional multi-source consistency, remote listener, UI object contents or production SLO.
- Source IDs and Case search metadata may be sensitive even within loopback; never treat this local demo as production software. Do not merge directly into main until integration base and required review gates are satisfied.
- The C-NAV `Navigate` reference is an authorized-at-read **draft**, never the final owner's screen, URL, bearer or permission token. Destination owners must recheck access on actual open. OPEN-013 remains unresolved.

## Verification and acceptance
- Ported first-party tests from PR #38: malformed/extra/trailing/oversize/nonregular JSON, cross-tenant links, real loopback startup/shutdown, Case picker and offline example.
- New end-to-end test: exact imported example → Case discovery → Timeline read → source navigation draft and safe return → unauthorized return denied → authorized return allowed → tenant/selector tampering denied → source permission revoked denied.
- Existing Go `-race`, vet, gofmt, engineering-kernel and Timeline workflow must pass on an **exact published commit SHA**. Passing synthetic integration tests does not establish a production IdP/provider or finished CAP-INV-110.

## PR/branch strategy
The canonical engineering branch and PR #36 are left untouched; this feature is a small stacked PR targeting that branch. Existing #37/#38/#39 remain draft and are not force-pushed or silently merged; their duplicate/overlapping changes require explicit later reconciliation. No canonical Product Spec file is edited.
