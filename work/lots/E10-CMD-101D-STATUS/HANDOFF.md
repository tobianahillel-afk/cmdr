# E10-CMD-101D-STATUS — Bounded Command Work Queue status filtering

## Four owners and canonical boundaries

- **Canonical:** `CAP-CMD-101` Unified Work Queue, Command module Incidents and Work Queue; `REQ-OBJ-001/012`, `REQ-PROD-013`, `REQ-UX-008/009`. Full canonical §§1–27 (`cmdr-product-spec/06-command/modules/incidents-and-work-queue/capabilities/unified-work-queue.md`) requires a single route, six system views, manual filters/sort/selection/Inspector. Incident and Task status are Command-owned; Shared owns generic Saved Views persistence/permissions, and the Design System owns final Inspector. Neither is implemented or cloned in this lot.
- **Engineering:** new one-behavior task `E10-CMD-101D-STATUS`, under existing epic `E10-CMD-101`, depends on `E10-CMD-101C-HTML` **IMPLEMENTED**. Complexity: one Command context, no migrations, no alternative authorization engine, no external package.
- **Code:** `product-runtime/work-queue` uses existing Go stdlib `Project`, `NewHTTPHandler`, owner `List/CanRead`, trusted host `SessionResolver`, and `html/template`. Adds a literal, bounded status selection; does NOT invent a filter DSL, Shared Saved View or new path.
- **Git:** a new child of PR #160 EXACT FROZEN HEAD `55d67812747f1893ee4f485247b167d20319fd67` (not `main`), all previously audited parent SHAs unchanged. Parent is a GIT relation, not canonical ownership.

## Requirement → delivered behavior → exclusions

`CAP-CMD-101` §§3/11/12/25 requires manual filter and sort as a no-AI alternative. The first-party read-only Command Work Queue already renders six views on one GET path with host scope+object ACL and server-rendered HTML, but had no manual filter at all. The new bounded capability lets an analyst enter **one literal status** from the existing authorized Incident/Task `Item.Status` field (one query parameter `status`, exact case-sensitive equality; empty means no filter; max 64 UTF-8 bytes, nonempty cannot be blank/control/wildcard). `Project` invokes current `Owner.CanRead` for every in-scope Incident/Task BEFORE evaluating both view and status. The raw status query is **not a permission** and cannot select tenant/environment/principal, Saved View or object owner. Owner errors and revocation remain fail closed; order and 5000-candidate/50-output bounds unchanged.

The HTML adds a semantic GET form on the SAME `/command/work-queue` route; all six canonical navigation links preserve the submitted status via Go `url.QueryEscape` and template contextual HTML escaping. CSP now permits only `form-action 'self'` with `default-src 'none'`; no JavaScript, external libraries/styles, mutable or POST action. JSON consumers without status receive the same original fields; filtered JSON adds `status_filter` and only authorized matching rows. This is ONE material Class-0 step, **not final multi-field filtering, configurable sorting, grouping, Shared Saved Views, persisted user data or full CAP-101 acceptance**.

### Operator-shaped behavior

`GET /command/work-queue?view=tasks&status=open` with a trustworthy host session and authorized Command owner returns only Tasks with exact owner-supplied status `open`. A browser can filter and then change to `incidents` without changing the path or forgetting the current status; a new HTTP GET re-resolves session+ACL before producing any content. `status=secret` cannot expose a denied incident even if one exists. A permission revocation still fails/reduces projection under current first-party rules.

## Acceptance, evidence and limitations

Tests call the real `Project` and real `http.Handler`: positive cases All/Tasks/Incidents, no-match vs denied-only distinction, original no-filter JSON, all six links and GET form, context encoding of HTML-like malicious status, rejection of duplicated/oversized/control/invalid UTF-8 `status`, no client tenant scope, denied host with no output. No additional owner.List or CanRead; only literal string comparison after existing ACL.

Canonical Incident status and Task status have different life-cycle semantics, but this opt-in exact literal equality intentionally does NOT unify them into a new shared machine. Final Task field/state schema and custom Shared Saved View fields are outside the stage; no invented acceptance of OPEN decisions. One partial page is visible when more than 50 authorized matches; no cursor/total count implied. No production IAM, deployed server, host cookies/TLS, persistent backend, per-row Inspector, actions, sorting control, linkage Case/Decision/Run, or Signed Evidence.

**Required GitHub proof:** exact-head `command-work-queue-runtime` three OS jobs + mandatory `engineering-kernel` including `CHK-SECURITY-TESTS` changed-security critical >=90 unchanged. One independent risk-based L2 for new status query/HTML filter boundary only after stable CI; no duplicate #160 L2. Child remains DRAFT/unmerged pending upstream #160/#159/#36/#2 native approvals and merges. The recorded state must never be read as full canonical capability acceptance or production integration.

## Next work after this bounded stage

Only select another ready unit after reconciling work graph against GitHub and the full CAP-CMD-101 obligations. Candidates include a separate bounded user sorting control **only if** source schema/Shared contract allows, or a separate native integration with Shared Saved Views and Inspector once their providers and contracts are ready. The major current delivery risk is the long nonmerged Git ancestry, not the count of these small proof units.
