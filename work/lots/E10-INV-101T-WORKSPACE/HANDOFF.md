# E10-INV-101T-WORKSPACE — Browser Case Queue → authorized Workspace handoff → reauthorized return

## Canonical authority, ownership and state

**Canonical CAP-INV-101** `cmdr-product-spec/07-investigate/modules/cases-and-evidence/capabilities/case-queue.md` §§11, 12, 17 and 25: filtered queue → open authorized Case Workspace context → exact view/filter/scroll/selection return. Case remains Investigate, Incident Command, Saved Views Shared. Screen destination `INV-CAS-001`; `perm.investigate.case.read` only. **Engineering epic E10-INV-101**, child task E10-INV-101T-WORKSPACE following browser filters E10-INV-101S-FILTERS. Code owner `product-runtime/case-queue`. **Git parent** PR #150 HEAD `a5034c3dfb6ae23a9fa61725e7d9057277db72fe` is NOT the engineering parent and NOT main.

## What was missing

The native opt-in Queue browser already listed authorized Cases with applied q/status/owner/freshness/sort, 25-row cursor pagination and independently reauthorized Inspector selection. First-party Go server `GET /case-queue?open_case=EXACT_CASE_ID&scroll=...` and `POST /case-queue/return` independently existed, tested in E10-INV-101J-HANDOFF/K-RETURN/N-RETURN, but had **no browser consumer**. This means no actual end-user handoff/return journey between the existing components. Building a second Case navigation engine or storing a token client-side would violate owner constraints.

## Implemented bounded browser journey

- Each current authorized Queue row adds **Open workspace handoff**. It invokes the EXISTING same-origin credentialed/no-store `GET /case-queue` with the exact opaque Case ID, current **last-applied** filters/sort and page cursor, plus bounded integer scroll. The server reauthorizes current session, Case, Incident/Finding source owners, Saved View, filters, page and freshness.
- Only a valid server `WorkspaceHandoff` with exact Case ID and `screen_id=INV-CAS-001` may enter a **transient, read-only navigation context**. Previously rendered sensitive Queue rows are wiped. The unsigned return token stays **only in volatile JavaScript memory**, never URL, storage, output or a new owner.
- Explicit **Return to authorized Queue** submits that volatile token exclusively via `POST /case-queue/return` JSON body with same-origin credentials/no-store and existing trusted host session. The first-party Go `RestoreCaseQueue` revalidates tenant/environment, Saved View version and permission, exact selected Case, freshness and original page cursor/filter projection. The client also rejects mismatched returned projection/Case.
- A successful POST alone is insufficient to redisplay an old browser row: the browser clears the transient handoff and fetches the selected Queue page again, so *current* independent Case/Incident/Finding read authorization still gates the UI. Scroll is restored only after this fresh successful GET. Denied/open mismatch/return mismatch or revoked/stale page clears sensitive details. Filters, pagination, clear and pagehide invalidate the transient state.

The web page explicitly labels this **navigation context only**. The actual **full Case Workspace UI is NOT hosted or implemented here**; no fake production page, no Case mutation or assignment/export, no Shared Inspector duplication, no real IAM/CSRF/TLS/audit provider, Saved View storage, new third-party runtime dependency or paid API. OPEN-013 stays untouched. Transient unsigned tokens are **not authorization credentials**.

## Behavioral acceptance and CI

Node built-in `browser.test.cjs` executes real browser.js in a stub DOM, covering successful current-page server handoff+return POST+fresh GET, exact Case ID and applied filter preservation, scroll, token never in URL, rejected/mismatched handoff and revoked/stale return with no private rows. Existing Node filtering, paging, Inspector reauthorization/XSS/focus tests plus Go first-party HTTP J/K/N server round-trip remain in existing kernel suite.

The **required** `engineering-kernel` exact-HEAD CI must include CHK-CASE-QUEUE-CONTRACT, CHK-MANIFESTS, CHK-SECURITY-TESTS with existing >=90% changed-security-critical threshold; NEVER bypass. A draft PR, successful child CI, or ChatGPT read-only L2 does not mean native GitHub approval, acceptance, merge or production. A significant browser→server return trust boundary deserves ONE independent, read-only L2 on fixed base/head, if published and stable. Ancestor draft PR #150 and earlier audited SHA never modified.

## Remaining work

Real host-mounted Workspace screen/nav router, full Shared Inspector and production Case/Saved View providers are separate canonical obligations. The current handoff demonstrates the *real transport and reauthorization loop*, not the final multi-module application deployment. The stacked draft PR lineage and #36 engineering base must be reconciled before any full product delivery; `main` remains documentary.