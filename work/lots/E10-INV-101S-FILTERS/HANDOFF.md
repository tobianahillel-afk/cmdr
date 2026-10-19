# E10-INV-101S-FILTERS — Browser use of first-party authorized Case Queue filters

## Authority and ownership

Canonical **CAP-INV-101 Case Queue**, `cmdr-product-spec/07-investigate/modules/cases-and-evidence/capabilities/case-queue.md` §§3, 8, 11, 12, 17 and 25; `REQ-PROD-014/020`. Case data/authorization remains Investigate, Incidents Command, Saved Views Shared. Engineering parent **E10-INV-101** epic; predecessor **E10-INV-101R-SELECT** IMPLEMENTED. Code owner `product-runtime/case-queue` native browser; underlying Go queue projection, authorization and sorting already exist at E10-INV-101I-FILTERS/M-PAGE/P-INSPECT. Git base is **immutable draft PR #149** SHA `c73f05b16a72af97c9845cbf4e4624baf9052c9f` (NOT `main`).

## Material behavior

Previously the native authenticated opt-in browser exposed text search, 25-row cursor pagination and independently reauthorized Case selection, while the existing first-party server already supported the `status`, `owner`, `freshness` filters and allowed `sort` fields. This child provides an accessible native UI for the **existing** server query contract (no second filter/search engine). Values are optional; literal text status/owner are bounded to 128, freshness is an explicitly selected first-party candidate value, and sort is restricted to known server-supported field+direction pairs.

On Apply, the browser captures ONE **last applied** search/filter/sort projection snapshot. The original read-only `GET /case-queue`, next page cursor and **new independent** `GET /case-queue?inspect_case=EXACT_OPAQUE_ID` all use that same applied snapshot. Merely editing an unsubmitted form DOES NOT silently change the current queue page or selected row's authorization context. Applying/resetting clears old cursor, search selection and current owner previews; page changes recheck current ACL and wipe old inaccessible results. Reset resets all controls, page and selection. The server alone selects tenant/environment/current host Saved View and reauthorizes per-object Case/Incident/Finding sources before projection. No client-provided permission/tenant/scope, no raw source preview, no Case mutation, no export, and no Shared Inspector visual clone.

### Boundaries and exclusions

This is a bounded **Class-0 browser integration**, not production full `CAP-INV-101` acceptance. Remaining canonical elements include actual host-backed Case Workspace navigation + exact return/scroll in the final frontend, complete Shared Inspector component integration and production host IAM/providing authoritative Case/Incident/Finding projections. Mutations remain outside scope (OPEN-013). No new third-party package, persistent Saved View adapter, database, full UI framework or paid API.

## Behavioral tests and exact evidence

Existing Go HTTP tests already exercise strict server-side filter/sort allowlists, per-request Case/Incident/Finding ACL, Shared Saved View permission-filtering and revocation-safe pagination. New Node built-in DOM tests run the **actual browser.js** in the Go-selected Node suite. They verify:

1. Applying status/owner/freshness/sort drives the known first-party GET query, with no tenant/permission input.
2. The same applied filter and sort persist across cursor page changes and the second, independently authorized Inspector GET.
3. Unsubmitted filter changes NEVER alter the security context used for an Inspector selection.
4. Clearing the controls and page immediately purges old protected details; original denial, stale cursor, inert text rendering and focus restoration remain tested.

**Source/CI truth:** the manifest/graph status is `IMPLEMENTED` only; exact SHA GitHub Actions `engineering-kernel`, changed-security-critical 90% gate, one independent read-only L2 if stable, native approvals and ancestors' checks must still be verified. Do not call a draft accepted/merged, and do not mutate any reviewed parent SHA. The previous `current.json` statement about R awaiting CI was stale: #149 has exact-head run `39318491676` SUCCESS. Full CAP-INV-101 is not closed by this UI adapter.
