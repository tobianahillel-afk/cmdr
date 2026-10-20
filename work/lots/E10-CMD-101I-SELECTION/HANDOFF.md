# E10-CMD-101I-SELECTION — Recoverable handoff

## Four distinct owners and authoritative contract

- **Product:** Command / CAP-CMD-101 (Unified Work Queue). Canonical capability `cmdr-product-spec/06-command/modules/incidents-and-work-queue/capabilities/unified-work-queue.md`, especially §§3, 11, 16, 17 and 25; work module `saved-views.md`; Command owns Incident and Task objects, Shared owns persisted Saved Views, final Inspector belongs outside this task.
- **Engineering:** bounded task `E10-CMD-101I-SELECTION` under existing epic `E10-CMD-101`; actual prerequisite `E10-CMD-101H-RETURN` recorded IMPLEMENTED with two green GitHub Actions on PR #166.
- **Code:** `product-runtime/work-queue/http.go`, `html.go` and behavioral `pagination_selection_test.go`, no new owner/auth/interface/service or dependency.
- **Git:** draft child of PR #166 branch `feature/e10-cmd-101h-exact-prior-page-return`, IMMUTABLE base `47ecdc1809769624ea1f37a4710336f48ba96301`; the Git parent is NOT the engineering/product parent. Parent frozen for independent review. No force-push, direct main write or merge.

## Concrete canonical gap (reproduced from code, before change)

The existing browser selection button on pages 2+ belonged to the manual filter form, which omitted `cursor`. A click on `Task:task-100` from the third of three authorized pages submitted `view=tasks&status=open&sort=oldest&selected=Task:task-100`. The existing server correctly restricts `selected` to a visible row on **the current page** and returned 404 because the form reset to the first 50 rows. This is a significant, reproducible failure of the canonical selected-context retention workflow, not a missing authorization bypass.

## Functional code delta and invariants

The manual filter/chronological sort GET form is now closed **before** the table; it deliberately contains no page cursor, so changing these conditions always resets to page 1. Every already-owner-authorized visible row has its own standard HTML GET selection form. This form carries the **current page's opaque cursor** plus the current `view`, literal `status`, `sort` and its own selected `Kind:ID`. The single authoritative `Project` first verifies the complete current owner and the revision/digest of the requested page, then writes `Result.CurrentCursor` (private `json:"-"` presentation field) only on successful authorization. The HTTP handler passes **only that verified Result** into `html/template`, never the raw query token; no session, identity or permission comes from HTML. No JavaScript, client-side query assembly or new route. All navigation goes through unchanged `GET /command/work-queue`.

On selection the server still obtains a fresh authenticated host session, calls one `Owner.List` and checks `Owner.CanRead` for every same-scope object, recomputes the bounded filtered/sorted authorized page SHA, denies stale/foreign/revoked page cursors generically, and selects only one of the 50 currently visible rows. No other Case/Incident/Task rights may be inherited from another page.

The HTML table caption now says *page courante* rather than falsely claiming later pages are the first page. The unselected/selected JSON projection schema and the pagination owner are not changed.

## Actual executable behavioral evidence and non-goals

`pagination_selection_test.go` tests three pages (123 authorized Task records, an owner-denied Task and a foreign-tenant Task) and reconstructs the actual `html/template` row GET form. It selects page-three `Task:task-100` and page-two `Task:task-050`, asserts the original cursor, view/status/sort and correct 23/50-row authorized slice and selected context, confirms fresh host and 124 in-scope owner ACL queries on each page, deliberate first-page reset on edited manual filters/sort, off-page 404, revoked old page cursor 404, no foreign/denied source leakage. Prior first-party 6-view, HTML XSS/CSP, filter/sort, forward/back pagination, runtime tests remain part of the same Go module.

GitHub run proofs **must be read at the latest exact HEAD**, not assumed from this handoff. Initial code SHA `ea6f1e4fce31ec2575578dee663cb3cd86e5a9be`: `command-work-queue-runtime` run 39377715483 COMPLETED SUCCESS Ubuntu/macOS/Windows; `engineering-kernel` run 39377715480 FAILED **only due to invalid initial work-unit ID** `E10-CMD-101I-PAGED-SELECTION`. Both initial invalid work path and graph/cache references have since been replaced with the valid `E10-CMD-101I-SELECTION`. The next run pair on corrected cache/graph SHA `c2d35e96c1ad9b61225657e905636a8cfb1002dc` is `command-work-queue-runtime` 39377975899 (3 jobs SUCCESS) and `engineering-kernel` 39377975861 (**must recheck final result**). No coverage thresholds, SAST rules, third-party policy or canonical specifications were changed.

Remaining: check `engineering-kernel` full exact-head result; only if all mandatory checks green, update graph/status from READY to IMPLEMENTED with evidence (no ACCEPTED), produce final exact-head GitHub validation and request ONE independent risk-based READ-ONLY L2 for material page selection/auth boundary. L2 ChatGPT is not a native GitHub APPROVED review. PR chain #166 through #36/#2 still draft and unmerged, `main` remains documentary; production IAM, Shared Saved Views, native Inspector and final canonical CAP acceptance absent.

## Security static-analysis correction on the same bounded cause

The initial implementation passed the raw `pageCursor` local from `http.go` directly to `renderQueuePage`. Although the result had already passed owner authorization, `engineering-kernel` exact-HEAD run **39377975861** correctly failed `CHK-SAST-GO` with `G705/MEDIUM` at `product-runtime/work-queue/http.go:171`: a direct request-to-HTML taint path was still visible. This failure must never be reported as a successful check. The correction does **not** suppress G705 or waive SAST. Instead, authorization and current-page digest validation in `authorizedQueuePage` now determine whether `Project` may populate `Result.CurrentCursor`, a presentation-only `json:"-"` field. The native renderer receives that validated owner projection only; the request token no longer crosses directly into the HTML boundary. Recheck SAST and all tests on the final commit, do not substitute the old kernel failure or old runtime success for its results.
