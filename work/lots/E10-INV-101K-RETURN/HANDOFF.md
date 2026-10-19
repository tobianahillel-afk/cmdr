# E10-INV-101K-RETURN — Native authorized Case Queue HTTP return

## Ownership and canonical obligation

- Product: CAP-INV-101 Case Queue §25 and §§12/17 (REQ-PROD-014, REQ-PROD-020), consumer CAP-INV-110 Case Workspace, screen INV-CAS-001. Only Class-0 read-only return is in scope; Saved View Shared ownership and CAP-INV-102 Case mutations remain outside.
- Engineering: existing epic E10-INV-101; new bounded task E10-INV-101K-RETURN depending on E10-INV-101J-HANDOFF **IMPLEMENTED**. Canonical work owner is not its Git parent.
- Code: first-party native `product-runtime/case-queue` with the **existing** `RestoreCaseQueue` engine as the single truth for token and view-state validation; no second parser for domain filters or ACL.
- Git: base frozen PR #141 `dcabc46620e06e0452a0e5d9c6827b4aff2a006c`; do not edit parent/L2 SHA.

## Before / after

The native authenticated GET `/case-queue?open_case=...` already generated a bounded `WorkspaceHandoff` and an opaque base64 return token. The original Go `RestoreCaseQueue` library could restore exactly the caller's prior selected Case, q/filters/sort, Saved View/version and scroll after reprojecting the **current** Case graph, but the native HTTP boundary exposed only open, not return. The canonical Case Queue §25 round trip was not host-consumable.

**Now:** call `POST /case-queue/return` with `Content-Type: application/json` and an exact, single-key JSON body `{"return_token":"..."}`. This does not transmit the token in a GET query URL, browser history or referrer. There is a 16KiB-ish body cap, no query parameters, duplicate keys, unknown keys, malformed/trailing JSON and forged tenant/permission fields rejected **before** provider/session. GET at the return route and POST at the Queue listing route are refused. The host must still mount the handler and enforce TLS, real sessions, rate limiting, audit and CSRF security in its production host; this alone is not a deployed Case Queue.

The existing host-injected `QueueSessionResolver`, `QueueCaseProvider` and `QueueReadAuthorizer` determine current tenant, environment and read access; all Case, Incident, and Finding candidates are independently authorized **before** `RestoreCaseQueue` inspects the untrusted token. Denied nested data is removed. `RestoreCaseQueue` itself checks token version, tenant/environment, valid Saved View/version, Case visibility, current filters, stale/denied Case and exact original selection, and returns only an authorized `ReturnState` JSON object — **not** a Case payload or new permission. Generic error, no partial output, no-store/no-referrer/Vary Authorization preserved. Existing GET list/open cases and first-party projector unchanged.

## Evidence / scenarios

- Real GET with all existing q/status/owner/freshness/sort + exact opaque whitespace-sensitive Case ID and scroll produces WorkspaceHandoff. Independently authorized POST on the **same owner** restores the selected Case, query, filters, sort and scroll; denied matching Case is absent, both requests re-authorize all provider candidates.
- Revoked Case access, Case staleness or changed tenant after issuing the original token reject the return, with no selected Case disclosure. A token is never a signed authorization grant.
- Bad/missing/multiple/oversized keys, malformed JSON, wrong method/path/media type and any URL token/tenant controls reject before the host session resolver or Case provider.
- Existing GET behavior and all Go Case Queue tests remain required, including security-critical changed-code coverage minimum 90% on exact SHA (never waive). Other ancestor PRs have independent L2/merge obligations.

## Explicit exclusions

No production identity provider, persistent Case DB/queue, full UX, Saved View mutation, permission issuance, Case assignment/status/OPEN-013, export, signed return token, socket listener or new third-party package is implemented. The host must route and protect POST and avoid logging request-body token; base64 presentation state is neither authenticated nor secret storage. A read-only server projection is not canonical product ACCEPTED or production.
