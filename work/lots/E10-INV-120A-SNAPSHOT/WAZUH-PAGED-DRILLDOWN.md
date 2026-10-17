# E10-INV-120A — Wazuh authorized paged analyst drilldown

Functional increment: the existing loopback Wazuh 4.x Case -> imported Timeline -> top-agent/rule cohort -> typed-source preview now exposes **all** authorized matched event identities in 25-row pages rather than truncating at the first 25.

- The existing owner-backed `collectDemoCaseTriageWithObserver` is rerun for **every page**, including source ACLs; raw `agent.id` / `rule.id` remain selectors, not identity or permission.
- Fixed page size, canonical nonnegative bounded offset (multiples of 25 up to 5000), stable traversal ordering within the immutable imported owner snapshot, complete count before emission, and generic deny with zero material on invalid selectors or revoked permissions.
- Response adds `offset` and `has_more` while preserving the existing bounded metadata-only drilldown receipt. Browser Next/Previous uses the same bearer in an Authorization header, never URL, and discards stale/cancelled requests. Entry preview still opens through the existing source reauthorization route.
- Integration tests exercise 53 locally imported realistic Wazuh-shaped events over 3 pages and refusal on invalid/cross-Case/cancelled selectors; Node UI tests cover next/previous, bearer handling and denial clearing.
- Reuse choice: stdlib net/http and the existing authorized in-memory owner/loopback UI are sufficient for <= 5000 imported entries; no database, indexing service, new IAM, third-party package or paid API. Pagination reruns the source traversal, a deliberate lab-only performance cost; production scale and concurrent mutation require a persistent authorized snapshot or cursor design.
- This is **not** production IAM, distributed authorization, durable storage, forensic chain of custody or `OPEN-013` resolution. No canonical product spec was changed.

Do not mark the lot verified until exact-SHA CI and an independent risk-focused review of this meaningful functional integration have passed.
