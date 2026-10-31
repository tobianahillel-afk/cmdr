# E10-CMD-101C-HTML — Bounded read-only Command browser

## Owners / canonical grounding
- **Product owner:** Command CAP-CMD-101 (Unified Work Queue), exact §§1–27; six view keys in `cmdr-product-spec/06-command/modules/incidents-and-work-queue/saved-views.md`. Incidents and Tasks are Command-owned; Shared owns full Saved Views mechanism, filters and Inspector semantics. No active separate Work Queue Screen ID.
- **Engineering parent:** existing `E10-CMD-101` epic. This task C depends on independently implemented `E10-CMD-101B-HTTP` (task B), not simply on its Git branch number. C is `IMPLEMENTED`, not `VERIFIED` or canonically accepted.
- **Code owner:** `product-runtime/work-queue`; reuse first-party `Project`, `NewHTTPHandler`, `SessionResolver`, `Owner.CanRead`, with no second owner policy, data store, framework or third-party product dependency.
- **Git parent:** immutable PR #159 head `81cdd5ab55786cbe93172a0968f3fac427b1104c`, base branch `feature/e10-cmd-101b-authorized-single-http-route`; ancestors #158, #157, #155, #154, #153, #36, #2 remain draft and unmerged. This child never modifies those SHAs.

## Gap / deliverable / explicit exclusions
- Gap: the existing exact-six-view HTTP handler has no browser-consumable workspace. A current authenticated host must wire its session provider and Command Owner to obtain records; there is no bundled production IAM or deployed service.
- Delivered: a first-party Go `html/template` server-rendered read-only document when an authorized GET to **the same canonical route** `/command/work-queue?view=...` requests `Accept: text/html`. Links to six fixed system views navigate the **same route**, invoking host `ResolveSession`, `Project` and per-item `Owner.CanRead` anew on each request. No JavaScript or inline stylesheet. HTML context escaping, CSP `default-src 'none'`, no cache/no store, frame deny, no referrer and no permission-bearing client state. Presents tenant/environment, state/partial diagnostics, up to first 50 authorized Incident/Task rows and explicit empty state. The existing JSON API output for clients without HTML Accept remains unchanged.
- NOT delivered: full CAP-CMD-101 acceptance, production host login/IAM, persisted/shared Saved Views, user filters/sorting, selected-row context/return/scroll, Inspector, linked Case/Decision/Run navigation, SLA authoritative source, full pagination, mutating or governed responses, deployment. Legacy Team Load remains deprecated; six routes and Screen IDs are never introduced.

## Risk acceptance
- Same request boundary as existing first-party host session and owner ACL. No early rendering or partial row disclosure on revoked session, owner failure, cross-tenant object or invalid view; HTML does not bypass Project.
- Escapes server-owner-controlled fields including strings that would form tags, script or event handler markup. Browser never executes owner HTML. JSON clients remain unaffected.
- Security model: only synthetic host+owner in tests. Host still has to provide trustworthy session and authoritative data with audit for production.

## Exact-head evidence and disposition
See PR body and GitHub Actions for the immutable child HEAD, tests and required CI. Mandatory 90% changed-security critical coverage and SAST/performance checks remain unchanged; no fusion until *all* required gates, ancestors and native approvals pass. One independent targeted READ-ONLY L2 for material new HTML/data display boundary on exact child diff; do not duplicate parent #159 L2.

## Deterministic next steps
1. Confirm exact-head command-work-queue-runtime on Linux/macOS/Windows and engineering-kernel, including HTML injection negative tests.
2. Obtain one risk-based independent L2 on the frozen child SHA; do not assume automatic native GitHub approval.
3. Resolve parent stack review/merge debt separately; product owner must later clarify final Shared Saved Views and Inspector contracts, host IAM and persisted object ownership before full CAP-CMD-101 acceptance.
