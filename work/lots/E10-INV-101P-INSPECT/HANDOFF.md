# E10-INV-101P-INSPECT — Authorized Case Queue → Inspector read projection

## Ownership and canonical mapping

- **Canonical owner:** CAP-INV-101 §3, §9, §12, §17, §21, §25 (Case Queue → Inspector), Investigate/Cases and Evidence. `Case` remains Investigate-owned, `Incident` Command-owned, `Finding` Investigate-owned; Shared Inspector interaction/keyboard and Shared Saved View are not reimplemented.
- **Engineering parent:** existing `E10-INV-101` epic, next child after `E10-INV-101O-INCIDENT`. No new business capability invented.
- **Code owner:** `product-runtime/case-queue` native Go read-only owner; `http.go` reuses the original `ProjectWithSavedViewFallback`, host session, server Case/Incident/Finding read ACL, cursor pagination and Shared Saved View selected by the server.
- **Git:** child of frozen PR #146 `6a7a265b473def439d19f2a43bc381feac6256df`, not `main`. Parent and ancestor L2/audited commits not modified.

## Gap and behavior

| Canonical obligation | Previous runtime | This bounded increment | Outstanding |
|---|---|---|---|
| CAP-INV-101 §17 select in Inspector | Authorized Queue listing and full Workspace open, but no distinct Inspector read selection | Optional GET `/case-queue?inspect_case=OPAQUE_CASE_ID` returns exactly one source-authorized, currently filtered Queue `Row`, retaining authorized Case, Incident/Finding, stale/partial flags | Final Shared Inspector UI component, keyboard/focus lifecycle, responsive and accessibility, native frontend |
| §17/25 filtered selection | Authorized query/filter/Shared View; paginated listing, Workspace handoff and return | Inspector selection uses same source-checked projection and can only select a Case in currently authorized page and filters. Stale readable Case may be inspected as stale without changing the stricter full Workspace handoff | Production Case provider/session and persisted Shared view |
| REQ-PROD-020 authorization | Reauthorize every Case/Incident/Finding at HTTP boundary | No separate source resolution, owner bypass or new metadata returned. Revoked Case or stale cursor => generic not found. Denied Incident/Finding metadata omitted. | Production IAM, hosting controls TLS, CSRF, audit and rate limits |

## Acceptance and security

- Go `httptest` with Case, Incident and Finding independently allowed/denied by the injected host: inspect selected readable Case, preserve partial indicators, omit denied nested content and denied Case metadata, with no-cache headers; unrelated search/filter denied.
- Authorized paginated Queue `limit=1` → page-2 Inspector access; same Case from page 1 rejected, and ACL revocation invalidates the old cursor before revealing the Case.
- Invalid, duplicate, oversized or injection selectors and conflicting `open_case`/`scroll` controls rejected **before** session/provider; `inspect_case` is not an authorization grant. Existing Queue/Shared View/Workspace/return tests unchanged.
- No new business object, Case mutation, Command relationship engine, extra read owner, HTTP authorization model, storage or paid API. Source data and API are synthetic host-injected read-only integration only; no product-wide ACCEPTED or production-ready claim.

## Evidence and recovery

- Authoritative requirement: `cmdr-product-spec/07-investigate/modules/cases-and-evidence/capabilities/case-queue.md` (full §§1–27), `cmdr-product-spec/03-design-system/components/inspector.md`, canonical Case/Incident objects and tenant/permission contracts.
- `manifest.json` contains allowed files and bounded acceptance. `work/graph.json` node `E10-INV-101P-INSPECT` is `IMPLEMENTED`, never `VERIFIED` until all required evidence/reviews. `engineering/state/current.json` is only a reconciliation cache.
- Check exact HEAD GitHub engineering-kernel `CHK-MANIFESTS`, `CHK-COMPLEXITY`, `CHK-ARCHITECTURE`, `CHK-RUNTIME-DEPS`, `CHK-SECURITY-TESTS`; check targeted native Case Queue Go tests, not just synthetic projection summaries. Do NOT claim parent PR #146 is merged or approved because its engineering CI was green.
- Risk-based independent L2 on the exact material PR diff; read-only without GitHub write. Native GitHub review and merges separately required.
