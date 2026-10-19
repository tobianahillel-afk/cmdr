# E10-INV-101L-VIEW — First-party native Queue with authorized Shared Saved View

**Canonical product** CAP-INV-101 Case Queue §§3/8/12/17/25, read-only Saved View use; Saved View object and use/read permissions are owned by **Shared Capabilities**, not by Case Queue. Consumer CAP-INV-110 Case Workspace; screen INV-CAS-001. **Engineering parent** E10-INV-101 and predecessor E10-INV-101K-RETURN IMPLEMENTED. **Code owner** product-runtime/case-queue. **Git parent** frozen PR #142 HEAD `a41bbcf6b2ac40f366b17c29f9bd1c573d7ce6ff` (not parent business, not main).

## Canonical gap closed

BEFORE: the Go Case Queue library E10-INV-101C-VIEW applies caller-owned Saved View permission-aware, preserving version/fields/filters/sort and existing E10-INV-101D return validates version; however native HTTP E10-INV-101H/J/K had no way to accept a **host-selected** view, so its GET/open/return always used default.

AFTER: optional trusted `QueueSavedViewProvider.CurrentQueueSavedView(ctx,scope)` obtains the authenticated principal's current view/permission-checked elements from Shared's host adapter **on each request**. `NewAuthorizedQueueHTTPHandlerWithSavedViews` opts in; the old constructor still uses the default. The host MUST supply actual tenant-scoped Shared view read/use decisions AND independently re-evaluate individual fields, filters and sort permission flags. The client cannot send view ID, view data, tenant or permissions. The existing first-party `ApplySavedView` and `ProjectWithSavedViewFallback` sanitize and apply allowed elements, remove denied elements, and return state diagnostics. The existing `OpenCaseWorkspace` and `RestoreCaseQueue` encode/view-check ID+version and re-check the whole current authorized Case and view upon POST return. No new storage, view engine, permission subsystem or HTTP route.

## Acceptance and boundary evidence

- Injected trusted Shared view has an allowed status=review filter, allowed status field, **denied source-like field and denied owner filter**, approved sort=title:asc. Real native GET lists only matching independently owner-authorized Case with exactly one allowed field and no denied field/filter values or filtered Case ID.
- Real native GET selection and HTTP POST `/case-queue/return` restore Case ID and scroll and preserve Shared view ID/version, with a **fresh** Shared resolver call each time. Changing view version or revoking use right after token emission forbids old-token return without exposing Case or token.
- Foreign-tenant, invalid and failed Shared responses deny with generic 404/no Case data; an already denied view falls back to the generic default only on a fresh GET. A client-controlled `view_id` remains invalid before host access; prior non-Shared constructor still works.
- Case/Incident/Finding separate owner checks and authorization boundaries are not replaced. The host's Shared provider is the ONLY trusted source for per-element decisions; if it lies or is compromised, this adapter does not independently implement Shared IAM.
- **No authenticated production provider, TLS server, CSRF defense, persistent views, full UI, Case assignment, export, mutation, new third-party runtime or paid API.** OPEN-013 remains open for Case mutation and irrelevant to this Class-0 view.

## Status

Production code + behavioral tests committed in new child branch. Check the exact published HEAD GitHub Actions `engineering-kernel` required tests/90% changed-critical floor (and any other required workflows); no CI PASS assumed for changed HEAD without evidence. One independent risk-based L2 review for the material Shared owner boundary only; native GitHub approvals separately required. Set work graph `IMPLEMENTED` until reviews/integration, do NOT claim CAP-INV-101 full ACCEPTED. The runbook must identify Git SHA, PR number, checks and true next step without relying on chat history.
