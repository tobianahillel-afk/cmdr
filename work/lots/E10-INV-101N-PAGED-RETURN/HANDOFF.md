# E10-INV-101N-PAGED-RETURN — Case Queue page → Workspace → restored page

**Canonical owner:** CAP-INV-101 (Case Queue) §§17/25, REQ-PROD-014/020, screen INV-CAS-001. Case Workspace is a consumer, not a source of Queue permissions; Shared remains Saved Views owner. **Engineering parent** E10-INV-101; dependency E10-INV-101M-PAGE (IMPLEMENTED only); **Git parent** frozen PR #144 HEAD `5051c7e1b9444fc5e35e818c9fabffe4f9ccda2c`. Code owner `product-runtime/case-queue`.

## Canonical obligation → previous code → change

The product requires Queue filtered state, scroll and selected Case to be restored when returning from a Case Workspace. Existing first-party `OpenCaseWorkspace` / `RestoreCaseQueue` implement this for an unpaged Queue, and PR #144 enables owner-authorized `limit/cursor` listing, but it intentionally rejects `open_case` together with paging. Thus a Case selected on page two cannot use the native HTTP workspace handoff **with its originating page context**.

The native `GET /case-queue?limit=N&cursor=OPAQUE&sort=title:asc&open_case=ID&scroll=91` now accepts page selection only if the selected Case belongs to the **fully reauthorized current page**, using the exact `paginateQueueProjection` and `ProjectWithSavedViewFallback` engines. The same original handoff encodes bounded `page_limit` and `page_cursor` in transient, unsigned `ReturnState` fields, preserving previous token compatibility and old unpaged behavior. `POST /case-queue/return` independently resolves session, Case/Incident/Finding authorization and current host Shared Saved View, then `RestoreCaseQueue` repeats current projection and page cursor validation BEFORE returning the original selected Case, filters, sort, scroll and page state. Revoked/hidden/changed page or Case is denied with no partial data. **The cursor and token are neither credentials nor durable snapshots**.

## End-to-end acceptance

- Real Go `httptest`: page 1 returns A/B, page 2 C; open C on page 2, then POST return, restoring `page_limit=2`, exact `page_cursor`, `selected_case_id=C`, scroll 91 and sort. Both requests reauthorize all four candidate Cases including a denied source.
- Select C on page 1 or A via page 2 cursor: reject without token even though either Case is otherwise accessible. Revoke Case C before POST return, or change Case metadata after open: deny stale return.
- Reject missing `limit`, duplicated controls, invalid bounds, tenant injection and malformed cursor before host session where applicable; previous unpaged GET/open/POST, case paging, Shared Saved View and source ACL tests remain active.
- No API paid/external runtime package, new Case Queue projection/route, IAM, data store, mutation, export, production UI, TLS/CSRF or live provider. The host still owns those boundaries.

## Git & evidence

Child PR on frozen #144, not `main`; never rewrite audited ancestor SHA. Record exact HEAD and engineering-kernel CI (including mandatory 90% changed-security coverage) after workflow completion. Status `IMPLEMENTED` only pending one independent read-only L2 and native approvals; no full CAP-INV-101 acceptance or production statement. Graph/state reflect **this** code, not a fabricated merge.
