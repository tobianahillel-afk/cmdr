# E10-CMD-101K-OWNERFILTER — exact authorized responsible-owner filtering

## Four owners and canonical obligations

- Product: `CAP-CMD-101` canonical `cmdr-product-spec/06-command/modules/incidents-and-work-queue/capabilities/unified-work-queue.md` §§1,3,5,9–12,14–22,25 and Command Incident/Task owners. **Read-only filtering** on the existing system Work Queue; `CAP-CMD-102` retains ownership of assignment mutations, which are NOT implemented here.
- Engineering: Child task `E10-CMD-101K-OWNERFILTER` under existing `E10-CMD-101` epic; dependency `E10-CMD-101J-DIAGNOSTICS` IMPLEMENTED. No pretend `CAP-CMD-101` ACCEPTED or product-wide release.
- Code: Existing first-party `product-runtime/work-queue/{queue.go,http.go,html.go,pagination.go}`, one targeted Go behavioral test. **No new route/provider/schema/authorization engine**. No modification of shared Saved Views, Investigate or mutable Incident/Task objects.
- Git: New child of PR #170 `fix/e10-cmd-101j-filtered-sla-diagnostics`, frozen parent SHA `d7ac9d3c4141c186722862a383dcb3e97e9d9f10`. Parent #169 and foundational #36 remain DRAFT and unmerged.

## Functional delta

Before: operator could choose exactly six fixed Command-owned views, an exact StatusFilter, chronological sort, and revocation-safe cursor pagination; `Owner` displayed but could not be selected as a filter. Canonical CAP-CMD-101 requires simple manual filtering and a unified no-AI workspace.

After: a bounded optional literal `owner` query parameter filters **by responsible-owner string already provided by the Command owner** (`owner=analyst`, not a user login or tenant claim). The server resolves trusted principal/tenant/environment from SessionResolver and still performs **per-item Owner.CanRead for every in-scope candidate BEFORE owner equality**. This also applies to other-tenant/denied data: other-tenant IDs are not submitted to the ACL; denied items cannot affect rows, summaries, diagnostics or selections. Only authorized filtered items may set source freshness partial/state and relevant SLA-policy-unavailable warnings; a filtered-out unknown SLA must not mark the chosen owner cohort partial. An empty owner value removes the filter; exact unassigned items remain available through the existing `unassigned` view.

The same `GET /command/work-queue` route, six views, native `html/template` form, status filter, oldest/newest sort and cursor navigation are reused. Owner string is validated as a bounded, valid UTF-8 literal (<=128 bytes, no controls or wildcard), never interpreted as authorization or a general query language. Native view switching, next/previous links and current-page row selection preserve the owner filter. Editing filter inputs starts at page one. The cursor's full authorized-row digest is additionally **bound to the owner filter**, preventing cross-owner replay even when row sets happen to look equal. Fresh owner authorization and selection occur on each request.

## Evidence / acceptance

- `Project` confirms that an authorized analyst-owned at-risk Incident remains visible and **ready** even when another-owner Task lacks SLA source policy; a relevant owner-specific missing SLA preserves `partial-sources`. Denied/foreign rows remain hidden.
- Full real HTTP + native browser: 75 authorized matching Tasks span two pages; select row 70 on page two with the exact owner/status/sort/cursor preserved; changing owner or revoking the selected item yields generic 404/no partial result. No new route or client authorization.
- Invalid owner tokens, duplicates and forbidden extra scope arguments fail **before any trusted SessionResolver/Owner call**; HTML escapes values and links URL-encode them.
- Existing Linux/macOS/Windows `command-work-queue-runtime` and `engineering-kernel` including mandatory 90% changed security-critical coverage/SAST are required on exact HEAD. A failing or incomplete workflow blocks merge; parent's green CI does not transfer.

## Limits and handoff

This is a synthetic first-party Command owner projection, **not live production IAM, SLA calculation, assignment update, bulk action, Shared Saved View persistence or a complete CAP-CMD-101 implementation**. `OPEN-013` Class-2 mutations remain excluded. The file graph and state cache are execution aids; GitHub SHA/CI, native approvals and canonical specification remain independent truth. One independent read-only risk-based L2 is warranted for filtering/cursor authorization surface only after a stable published SHA, no duplicate review of #170's smaller diagnostic fix.

Next work should examine remaining canonical Work Queue transitions (source-owned linked object navigation, Inspector integration and fully native approved host adapters) or independently ready work, **not automatically add another filter**.
