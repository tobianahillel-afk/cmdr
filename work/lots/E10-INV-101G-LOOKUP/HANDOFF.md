# E10-INV-101G-LOOKUP — Recoverable Case Queue lookup integration

## Correct authorities

- Product: `CAP-INV-101` Case Queue, explicitly read-only Class-0 search/filter and Queue → Case Workspace selection; `CAP-INV-110` consumes selected Case in Investigation Timeline. Canonical authority is `cmdr-product-spec/07-investigate/modules/cases-and-evidence/capabilities/case-queue.md`, with `case.md`, server permissions and screen-routing contracts.
- Engineering: new bounded task `E10-INV-101G-LOOKUP`, parent `E10-INV-101`, dependency `E10-INV-101F` recorded VERIFIED. NOT a child of Timeline `E10-INV-110` just because its current browser adapter resides there.
- Code placement: operator-local **synthetic demonstration adapter** `product-runtime/investigation-timeline/cmd/case-timeline-demo/analyst.{html,js}`; frontend test harnesses `analyst_summary.test.cjs` and `analyst_return.test.cjs`. Permanent Case Queue application owner is **not implemented** by this small UI change.
- Git: draft PR #138 `feature/e10-inv-101g-authorized-case-filter`, based on frozen #137 `b98fac7334985fb028651ff946b86f170cb048eb`, not `main`. Parent PR #137 has exact-head runtime and engineering workflow SUCCESS (39229575679 / 39229575692) but remains DRAFT. No ancestor/L2-frozen SHA updated.

## Canonical requirement → implementation → proof → remaining

- **Search/find/filter authorized Cases**: old local Case picker returned a sorted fully owner-authorized list of up to 5000 IDs but had no local lookup. A new search input filters literal substrings **only** of the single existing authorized `GET /demo-cases` response, never triggers another source request, and preserves exact opaque Case IDs including whitespace-distinct values.
- **Safe Case navigation**: lookup alone does NOT modify the active Case, browser selection does so only via original `caseOptions.change` (which fires Case input invalidation), and the unchanged first-party server authenticates the Case, Timeline and typed-source on every request. No new Case authority or frontend permission computation.
- **Revocation/session changes**: the server response is cached only transiently; credential change aborts inflight discovery and clears the allowed list, disables lookup, and late responses cannot restore stale options. Malformed/denied catalogs fail closed without leaking a partial list.
- **Tests**: Node harness covers exact/space-distinct filtering, no extra HTTP, explicit transition, no hidden active-Case pivot, token rotation/racing request, malformed/denied result; existing local return-context harness updated to recognize the new DOM field.
- **Not implemented**: full native Case Queue, Case status/priority/Incident projections, Saved Views, versioned queue state, cross-navigation scroll restoration, Case assignment, exports, production IAM, durable Case storage, source ingestion and canonical acceptance. `OPEN-013` remains unchanged.

## Required evidence and gates

Record exact GitHub PR/HEAD and `investigation-timeline-runtime` Go+Node Linux, macOS/Windows runs; `engineering-kernel` remains mandatory, especially CHK-SECURITY-TESTS >=90% for changed security-critical coverage. Distinguish workflow `SUCCESS` from L2, GitHub review approvals, merge and production. This Case Queue browser adapter cannot be marked full CAP-INV-101 implemented.

Next: reconcile the #2 documentation baseline, #36 very large draft engineering foundation and stacked feature branches; avoid claiming Case Queue delivery on `main`, and do not add more browser-only fragments if provider/production ownership decisions block the next milestone.
