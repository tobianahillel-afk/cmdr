# E10-INV-110G-EXACT — deterministic handoff

## Canonical ownership and first verified discrepancy

- **Owner of functional behavior:** `CAP-INV-110` Investigation Timeline; `CAP-INV-101` Case Queue is the picker consumer; Shared `CAP-INV-002` Search Job and Investigate `CAP-INV-008` ephemeral comparison consume the **selected Case identity**. Existing product specs, tenant and permission contracts are unchanged.
- **Engineering parent:** existing `E10-INV-110`, bounded task `E10-INV-110G-EXACT`, following previously integrated `E10-INV-110F-SELECTOR` which is IMPLEMENTED but not canonically accepted. This is NOT a new canonical CAP nor a mutation operation.
- **Code owner today:** `product-runtime/investigation-timeline/cmd/case-timeline-demo/analyst.js`, first-party synthetic loopback UI ONLY, not production IAM or persistent Case Workspace.
- **Git base:** `feature/e10-inv-008b-110f-owner-integration` PR #134, frozen `5d7de39e53cf993c76f7811448221b7f67c8b20e`, whose exact-HEAD GitHub Actions runtime and engineering kernel were SUCCESS. Its own independent L2 has already been scheduled/executed in another context, result not claimable here; no change to that reviewed SHA.

**Material R2/root cause:** The backend `validHTTPCaseSelector` accepts an opaque authorized ID containing ordinary leading/trailing spaces (it checks UTF-8, nonblank scope, max length and control characters); Case picker returns that exact authorized ID. The browser historically ran `caseID.value.trim()` in 36 distinct selection/scope/authorization/navigation checks and `compareCase.value.trim()` in ten comparison target checks. This silently converted `" case-disk "` into `"case-disk"`, yielding an unauthorized denial or **opening a different authorized Case** rather than the one actually chosen. Server ACL did not grant unauthorized source access; it could not prevent the browser requesting the wrong Case.

**Correction:** Browser now passes the original exact identifier without normalizing it in Timeline GETs, Search Job POSTs, two-Case comparison selectors, source-open context and asynchronous stale-controls checks. The latter continue to invalidate data on Case/token/query changes. No new selector API, input parser, wrapper, owner permission model, storage or dependency.

**Behavioral regression:** A real Case picker mock displays two server-authorized distinct IDs (`case-disk` / ` case-disk `); selecting the latter opens Timeline with precisely its original bytes and issues one authorized completed Search Job with exactly that Case value, then selecting the former invalidates the previous results/provenance. The existing Case comparator must request both differing owner IDs without normalization. All other browser Case/Timeline/Wazuh/source ACL tests and Go cross-platform tests remain required.

**Recovery / no claim inflation:** Graph/node is IMPLEMENTED only after publication; review/safety/merge obligations unchanged. No accepted full CAP-INV-110 or CAP-INV-008, no production Case repository, no new persistence or evidence. Next step: verify exact-HEAD GitHub Actions `investigation-timeline-runtime` Linux Go+Node and macOS/Windows Go, plus `engineering-kernel` mandatory `CHK-SECURITY-TESTS` 90% floor. Do NOT borrow green statuses from the #134 SHA for this new head. Once verified, fix materially failing tests only; this single R2 compatibility fix does NOT automatically trigger another independent full L2 of already reviewed #134 unless new security behavior warrants it. Current `engineering/state/current.json` is an execution cache and must not promote partial CAP or backend state to accepted.


## 2026-10-19 — Usable Case-picker labels for opaque IDs (same R2 root cause)

After the exact-wire-selector correction, the same existing native `<select>` still
rendered `option.textContent = id`. An ID with significant outer or interior
spaces could therefore visually collapse into another Case label even though its
`option.value` and server authorization were different. This is a material
operator selection ambiguity in the SAME CAP-INV-110 / CAP-INV-101 read-only
handoff, not a new capability or new work-unit. In the existing
`case-timeline-demo/analyst.js`, preserve `option.value = id` byte-for-byte
and render only the option's **visual text** as the JSON-escaped original ID
with literal ASCII spaces shown as `\\u0020`. For example `case-disk`
renders as `"case-disk"`; ` case-disk ` as
`"\\u0020case-disk\\u0020"`. Do not decode the label back into an
authorization argument: the browser continues to submit the original value.

The already existing two-ID browser test now asserts exact option values,
unambiguous visual labels, lossless JSON decoding of the displayed string,
and retains the original real Case-picker → Timeline → Search Job request
and Case context invalidation scenario. Code/test changes are confined to
the **pre-existing authorized paths and PR #135**, no new PR/lot/schema/API,
product spec, backend auth, persistent Case or dependency. Historical parent
PR #134 SHA is unchanged.

**Recovery and validation:** PR #135 is still DRAFT and the E10-INV-110G-EXACT
work node remains **IMPLEMENTED only** until all required exact-head Actions
and any independent/native reviews are resolved. The earlier HEAD
`269598e441fc0cf3e1f7ef8c4fc91b1819e4c7a9` had runtime Linux/Node/macOS/
Windows SUCCESS and mandatory `CHK-SECURITY-TESTS` FAILURE 89.30% <90%;
those results MUST NOT be transferred to this additional code/test diff.
The full-scope Go security coverage gate remains mandatory even for browser
changes. No artificial tests or reduced thresholds should be added solely
to force green.


## 2026-10-19 — Exact Case identity in return links (same R2 root cause)

After correcting the browser picker and comparison, the first-party **return-from-Inspector Timeline URL** still rejected valid owner-issued Case IDs with significant surrounding spaces in `restoreReturnSelection` (`analyst.js` around line 1783). The old `id.trim() !== id` presentation check meant an authorized Case called `" case-disk "` was accepted by the backend and picker, but could **never** have its Timeline/time-range restored from a source-navigation return link. The navigation contract requires a return to the same selected Case, not silent alteration or an invalid link.

Minimal correction **within the already authorized E10-INV-110G-EXACT paths and SAME PR #135**: require the URL `case_id` to be nonempty, length <=1024, nonblank after a read-only whitespace check `!id.trim()`, and not wildcard, but **do not normalize** its identity. Keep existing strict allowed-query-key and single-value validation, page/time/other-filter bounds, no automatic bearer or HTTP fetch on URL load, and fresh server-side reauthorization after explicit submission. No new parser, router, IAM, state model or security rule.

An existing Node browser test harness in `analyst_summary.test.cjs` now takes an optional mocked URL query. New behavioral test checks a return URL with `case_id=%20case-disk%20`, verifies zero initial requests, exact unchanged restored Case value, and an explicit Timeline submit with the same original encoded Case identity and freshly entered bearer. A whitespace-only Case return link is rejected with no request. Existing native Case picker, SearchJob, comparison and inspector tests remain unchanged.

**Source changes recorded:** 2 real code/test commits on the same PR #135 after the prior frozen historical HEAD `44957cad74eeb6a391ede967a08f0dfdf95dd284`: `0d5a653e3ce6374676b5b95bf0a7aa6ba8578c3a` (browser return selector), then `78fa3353f94331d6c49111ac84f31fd125cd68de` (real browser regression). Both changes are in the task manifest's allowlist; no SHA in parent PR #134 was changed. The historical `44957cad...` runtime SUCCESS and security-test 89.30% FAIL are **not evidence for new HEAD**. Exact-head GitHub Actions for `78fa335...` must be read, and any subsequent documentation/ledger commit also mandates checking its **own** SHA. The mandatory 90.00% security-coverage gate is unchanged; if the fallback full-Go scope remains 89.30%, merge stays BLOCKED even when browser/runtime tests succeed. No repetitive tests solely to increase a coverage score; solving the underlying full-scope security coverage and PR stack is independent engineering work outside this small read-only Case selector task and **must not silently widen its allowed paths**.

Official node status remains **IMPLEMENTED**, never VERIFIED/ACCEPTED solely because this local synthetic UI fix works; parent PR #134, PR #36 foundation, approved live IAM and production Case persistence have separate acceptance and merge gates. L2 is not presumed from an existing PR #134 review.
