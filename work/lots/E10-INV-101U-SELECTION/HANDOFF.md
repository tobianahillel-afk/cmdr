# E10-INV-101U-SELECTION — Restore authorized selected Case on Workspace return

## Canonical obligation and ownership

- **Product** CAP-INV-101 §17 and especially **§25**: after returning to a filtered Case Queue, recover **view, filters, scroll, and selected Case**. Case owned by Investigate; Saved View and final Inspector owned by Shared; Incident owned by Command.
- **Engineering** existing `E10-INV-101` epic; task `E10-INV-101U-SELECTION` depends on prior `E10-INV-101T-WORKSPACE` IMPLEMENTED only.
- **Code** `product-runtime/case-queue/browser.js` and native DOM tests. No new Go provider/authorization; no access to production Case persistence.
- **Git parent** frozen PR #151 `e3e35e9d8e0fb1f3ba8378c9c24d3b2ee1fc9e99`, NOT an engineering/business parent and NOT `main`. Never rewrite frozen SHA.

## Precise product gap

PR #151 connected first-party `GET /case-queue?open_case` and `POST /case-queue/return` with correct filter/page+scroll restoration and a fresh authorized Queue GET. Yet it then **left the selected Case Inspector closed**, violating the last independent item of CAP-INV-101 §25. The transient return token and old selected row are NOT authority to redisplay a Case.

## Actual bounded behavior

After the existing first-party Workspace handoff, `POST /case-queue/return` revalidates tenant, environment, Case, Incident/Finding, Shared Saved View and page cursor. The browser invalidates its transient unsigned token and then:

1. requests the **current** filtered Queue page again via the existing server-authorized `GET /case-queue`;
2. checks the Case ID in the accepted POST is present **byte-exact** in that freshly authorized page;
3. issues a second **independent** first-party `GET /case-queue?inspect_case=<exact-id>` under the same last-applied filters/sort and cursor, requiring another host permission check, then renders only its currently authorized bounded Case summary;
4. restores original scroll only once both source checks succeed.

A missing Case on the freshly authorized page, a denied/mismatched Inspector read, a stale filter/page or canceled navigation purges old private content; it never reuses the old row, grants a permission from return token, or emits raw Incident/Finding/Evidence/Source data. Existing manual Inspect behavior, Queue paging/filters, no-store and same-origin transport are reused, not duplicated. This is a **read-only Case selection context**; the full `INV-CAS-001` Workspace and Shared Inspector remain **NOT delivered**.

## Falsifiable tests

Existing real Node DOM integration harness runs `browser.js` and already covers backend-normalized filter values, fresh GET + scroll, failed handoff/return, source permission errors, and inert rendering. Extended scenarios assert:

- success: open exact current Case → host-authorized POST return → new Queue GET → separate current Inspector GET (Case ID exact, same applied filter/sort/cursor) → restored selection and scroll;
- old Case disappears between POST and fresh GET: no Inspector request, no old rows or selected Case visible;
- Case still on new page but independent Inspector permission revoked: all older sensitive Queue/selection content cleared and no successful restoration.

No new API, third-party dependency, paid service, business mutation, new IAM or persistence design. The existing host must still provide real IAM, TLS, CSRF, audit, providers and production deployments.

## Evidence / stop-line

Read GitHub Actions on **exact HEAD**, specifically CHK-CASE-QUEUE-CONTRACT, CHK-SECURITY-TESTS mandatory 90% changed-critical floor and CHK-MANIFESTS; do not substitute parent PR #151's success for child CI. Set `IMPLEMENTED` only when code is published; `VERIFIED`/product acceptance/merge require the separate remaining gates and independent review. Request ONE risk-based read-only L2 for this new server-revalidated selection behavior, without duplicating the prior parent #151 L2. Do not merge while draft ancestors #151/#36/#2 and approvals remain unresolved.
