# E10-INV-120D-TECHNAV — Wazuh ATT&CK-to-Event navigation

## Canonical functional need
CAP-INV-110 defines source-authorized Case/Timeline navigation and investigation read filtering. Previously, authorized Wazuh source metadata `rule.mitre.id` was presented only as top-ten counts in the local analyst browser; unlike agent/rule cohorts, clicking a technique could not open constituent authorized Events. This work adds that read-only navigation as **a fourth selector to the exact same server-side drilldown**, not another engine or new endpoint. Wazuh labels remain *source claims*, not verified observed ATT&CK behavior or entitlement.

## Four distinct owners
- **Canonical product:** Investigate CAP-INV-110 controls Case selection, Shared controls Timeline Entry, and typed Wazuh-derived source owner controls read access to Event details. This is not a decision in OPEN-013 Class-2.
- **Engineering:** existing epic E10-INV-110; new task E10-INV-120D-TECHNAV after E10-INV-120C-MITREFIX `IMPLEMENTED`, not canonically `ACCEPTED`.
- **Code:** original Go `wazuh_drilldown.go`, original source parser `wazuhRuleLevels` and existing local `analyst.js`. Browser already uses bearer & separate `openSource` Case/typed-source read, no new provider/ACL/storage or dependency.
- **Git:** this branch starts **exactly** from frozen PR #218 SHA `930f41987d00f74f31f7706f91d43d5a144581c2`; Git parent #218 is NOT engineering parent. Never alter the parent undergoing L2.

## Behavior and risk boundary
- Existing fixed GET `/demo-wazuh-drilldown?case_id=...&technique_id=T1059.001&offset=0` validates only strict `Tdddd` or `Tdddd.ddd`. Technique cannot be combined with agent or rule, wrong scope, duplicate or arbitrary query fields. Its selector is a non-authoritative filter.
- Same `wazuhRuleLevels(raw, caseID)` source parse (including strict `null` denial from PR #218), then same original random-bearer handler + `collectDemoCaseTriageWithObserver` Case/Timeline/typed-source ACL traversal. Each authorized source Event is checked before matching against its MITRE source labels. The **whole** imported set must be visible and consistent, even for an empty or final page, before writing a receipt. No raw alert ID, preview/full_log, hostname or secret is output.
- Response uses the existing `cmdr-local-wazuh4-drilldown-v1` limited 25 entries/page, optional offset and `has_more`, and source ID/kind/time only. Browser button uses `textContent`, the existing same-origin bearer request and existing `openSource` which rechecks destination read permission; it does NOT trust the tag or response as permission.
- Existing agent, rule and agent∩rule drilldowns remain unchanged. Offline Wazuh import and synthetic Case owner are still test-only; no production IAM, live SIEM connector, durable audit, Evidence or technique verification.

## Tests and merge limits
Go test exercises disk-backed 53 alerts, 27 two-page matching tag occurrences (duplicate labels on same alert count once), empty/forged scope, bearer absent, mixed agent/rule, duplicate ID, invalid label, bad offset, deny before leaking raw. Node tests exercise source-claimed button, exact same-origin `technique_id` query, bearer never in URL, authorized Event identities and denial when returned selector differs. The workflow is `investigation-timeline-runtime` plus all engineering mandatory checks on **the final SHA only**.

This candidate is an additional Draft PR stacked *only because* #218 is frozen for independent L2; it does not make parent #218 or its long chain merged or production-ready. Do not start another task on this SHA before reviews/CI stabilize; promote integration and native independent approvals through the existing parent chain before merging. A ChatGPT L2 technical verdict is not native GitHub APPROVED.


## Independent L2 R2 correction — strict transport parsing (separate descendant)
Independent read-only review of immutable PR #219 HEAD `164019c252227f65048d292b33e699b2a275c18b` reported `CHANGES_REQUIRED`: `r.URL.Query()` silently drops invalid URL pairs, so an apparently valid selector with an additional malformed pair can pass the exact-key checks. Corrective branch `fix/e10-inv-120d-strict-wazuh-drilldown` starts from that frozen SHA without altering PR #219, its historical checks or product contracts.

In `wazuh_drilldown.go`, the existing endpoint now uses standard-library `url.ParseQuery(r.URL.RawQuery)` and rejects parsing errors before interpreting any remaining selector. Tests cover invalid percent-escape in unknown fields, offset, duplicate Case and semicolon pairs; no source data is emitted.

The nonblocking reviewer R3 observed that browser-side `offset`/`has_more` were optional. The same client now requires exact integer offset equal to the request and boolean has_more consistent with total count before rendering source Event IDs. New Node tests cover missing, mismatched and malformed metadata; existing valid mock is made protocol-compliant. Agent/rule/intersection/technique behavior otherwise unchanged, and source open reauthorizes Case/Timeline/typed sources.

This is a targeted correction of the existing E10-INV-120D-TECHNAV task, not a new product capability, data owner, route, auth engine or dependency. New exact-HEAD CI and targeted independent L2 delta review are required. No native GitHub APPROVED, merge to main or production IAM/remote Wazuh provider is claimed.
