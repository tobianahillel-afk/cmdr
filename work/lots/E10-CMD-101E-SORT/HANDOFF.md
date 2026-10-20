# E10-CMD-101E-SORT — Command Work Queue deterministic operator sorting

## Canonical and engineering ownership

**Canonical:** `CAP-CMD-101` §§1–27, Command / Incidents and Work Queue; `REQ-OBJ-001/012`, `REQ-PROD-013`, `REQ-UX-008/009`. Canonical objects are Command-owned Incident and Task. Shared Capabilities owns Saved Views persistence, sharing and policy-aware application. The Design System owns the final Inspector, not this lot. Final columns/virtualization are OPEN, so this unit only sorts the already-owned and already-authorized `UpdatedAt` timestamp.

**Engineering:** existing epic `E10-CMD-101`, child task `E10-CMD-101E-SORT`, depends on `E10-CMD-101D-STATUS` IMPLEMENTED (exact parent HEAD `2bc24c784b61ffc6c8b97fb81588ec02a84b4c59` with GitHub command runtime and kernel both success, NOT merged or accepted). Bounded to 1 Command context, 3 production Go files, one behavior, no migration. Full work unit and allowed paths in `manifest.json`.

**Code:** existing `product-runtime/work-queue/queue.go`, `http.go`, `html.go`. No new engine/provider, storage, wrapper, route, dependency or authorization model. Same Go first-party `Owner.CanRead` per in-scope Incident/Task BEFORE filtering and sorting, host-resolved principal/tenant/environment every request.

**Git:** child branch `feature/e10-cmd-101e-updated-at-sort` targets parent `feature/e10-cmd-101d-status-filter`, not main; parent SHA frozen. This is a Git stack edge, NOT a business parent.

## Functional delta

Before: six authorized system views, optional exact status filter, default newest-first stable UpdatedAt timestamp sorting, <=5000 candidates and <=50 displayed rows, native read-only HTML and JSON on one GET route. No user sorting option.

Now: Class-0 manual `sort=newest|oldest` on the SAME `GET /command/work-queue?view=tasks&status=open&sort=oldest`. Exactly one existing `Project` and one sorting pass, ascending or descending on the canonical owner `UpdatedAt` with unchanged deterministic Kind/ID tie-break on equal time. No explicit `sort` preserves original newest-first JSON shape and semantics. Any unknown/duplicate sort input, unsupported column, injected scope or malformed query fails before calling the host owner. No client authorizes an item by selecting sort. In HTML, the same accessible GET form offers oldest/newest and six view navigation links retain permitted status and sort with `url.QueryEscape` + `html/template` contextual escaping; no JavaScript or Saved View persistent state.

Tests (`product-runtime/work-queue/sort_test.go`) exercise authorized three-row ordering and status/view intersection, owner `CanRead` count and excluded foreign items, original JSON compatibility, all six HTML views and form retained navigation, invalid/duplicate arbitrary sort rejection before owner, revoked owner no rows, deterministic bounded 50-row first page. Existing runtime and Node/Go checks are preserved.

## Explicit exclusions and merge readiness

Only manual sorting on already-owned `UpdatedAt`, not arbitrary columns, priority policy, final sorting DSL, grouping, pagination cursors, saved personal views, Incident Detail, Inspector, actual production host IAM, durable Incident/Task source, Case/Decision/Run source linking or complete canonical acceptance. No paid API.

GitHub verification must use the **exact final HEAD**, validate all three OS Command Go jobs and unchanged mandatory `engineering-kernel` security/architecture/complexity checks. At the authoring of this handoff, the final CI result is not embedded: see the PR's head SHA and matching GitHub Actions runs as the source of truth. Do not claim VERIFIED, ACCEPTED, merged or production without evidence. A single independent read-only L2 is warranted for the new query and HTML form trust boundary when published and stable; do not duplicate parent reviews.

The existing upstream ancestry from foundational PR #2 and #36 through Work Queue PR #153–161 remains in draft/unmerged. Do not merge on green child CI alone. Main still documentary; this is a tested *branch runtime*, not a deployed product.

## Next bounded actions

1. Read exact-head Action results and fix only material failures on a new HEAD if needed. Confirm cross-OS, security floor and L2.
2. Reconcile long PR stack, native review, parent checks and accepted product obligations before merging or selecting another lot; never infer canonical acceptance from an IMPLEMENTED graph state.
3. Once an independently proven Shared Saved Views or Inspector provider exists, integrate through its owner instead of duplicating those mechanisms in Command. If not ready, record BLOCKED and choose an unrelated canonical obligation.
