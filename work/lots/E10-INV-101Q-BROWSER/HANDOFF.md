# E10-INV-101Q-BROWSER — Host-mounted read-only Case Queue browser

Canonical: CAP-INV-101 Case Queue (§§3,11,15,16,17,21,25; REQ-PROD-014/020), Case object Investigate-owned. Engineering parent: existing E10-INV-101 epic, depends on P-INSPECT IMPLEMENTED. Code owner: native Go product-runtime/case-queue with same-origin static browser assets. Git parent: immutable #147 HEAD b989d5d1284b8e7b8370021019a49b042d346d7e, NOT main, not product parent.

## Demonstrable behavior

The host opts into `NewAuthorizedQueueBrowserHandler(sessions, cases, access, views)` rather than the existing JSON-only constructor. The same authoritative owner HTTP serves GET /case-queue and POST /case-queue/return; the new GET /case-queue/ui, ui.js, ui.css are static, read-only assets only. The browser fetches `/case-queue?limit=25&q=...` and `cursor=...` with same-origin cookies/session, no explicit auth token storage or tenant parameters. Cells are inert textContent only; no business logic/search engine/ACL duplicate. Rejected or stale API responses clear old rows. No polling, assignment/export/mutation/Incident-workflow ownership.

UI renders title/ID, status, owner, freshness and updated UTC with search and first/next page controls. The server remains the sole source of tenant/environment scope and authorized rows; the host must provide real session resolver and object-read ACL. Strict CSP, no-store, no-referrer and no arbitrary static route delegation. A browser with only bearer Authorization but no host-owned same-origin cookie/session will not authenticate by itself. Do NOT embed bearer tokens or claims of production IAM into assets.

**Not implemented:** final Shared Inspector visual component, drawer/focus cycle, full Case Workspace host UI, Saved View management UI, production auth/session/TLS/rate limits/audit, persistent Case data/OPEN-013 mutation. The new UI is a bounded native read-only proof, not full CAP-INV-101 acceptance or deployment.

## Evidence and next step

Go HTTP behavioral tests prove opt-in assets headers, zero private data in assets, host cookie-based same-origin session, denied Case omitted and all candidate ACL checks invoked; bad static routes/methods denied without contacting host. Existing native Case Queue tests must remain green; exact SHA CI/quality gates mandatory. The parent #147 L2 task is already executed (not GitHub APPROVED) and must not be duplicated. This material new host-facing browser needs one separate read-only L2 once published and green.

If source CI fails, preserve draft and correct only actual defects on this child. Ancestor chain #147 → #146 → ... → #36 → #2 remains DRAFT/unmerged; do not claim merge/production. Reconcile cached graph state and CLI evidence from latest GitHub, not prior conversation.