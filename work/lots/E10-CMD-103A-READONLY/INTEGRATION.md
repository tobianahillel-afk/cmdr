# E10-CMD-103A-READONLY — canonical product owner in selected Command Work Queue

## Source of authority and traceable parent

**Product:** CAP-CMD-103 Operational Ownership §§1–27 under Command Incidents and Work Queue. The canon distinguishes operational Incident/Task responsibility, canonical product/Ownership Register ownership, and Platform Settings principal/team identity; CAP-CMD-101 supplies the already authorized selected Work Queue item, CAP-CMD-102 owns assignment mutations but OPEN-013 leaves Class-2 update authority unresolved. CAP-CMD-103 acceptance also needs ownership history, audit, mutations and identity resolution: these are **NOT** implemented here.

**Engineering:** new epic E10-CMD-103 (not another child of CAP-CMD-105 SLA), bounded Class-0 task E10-CMD-103A-READONLY depending on already IMPLEMENTED Work Queue E10-CMD-101M-MYWORK; source manifests/graph/current record the correct parent. **Code:** first-party `product-runtime/work-queue`, same GET /command/work-queue single current-page selected Incident/Task JSON or native HTML. **Git:** child branch of frozen PR #175 HEAD `391e4b13ddde572d47a756590bd09bc036495e0b`; this is only the Git parent, not the canonical or engineering parent.

## Canonical gap and actual behavior

| Canonical obligation | Existing code | This lot | Still missing |
|---|---|---|---|
| Owner operational vs canonical | Row.Owner is an owner-supplied operational principal; contributors/watchers supported | Optional, separately authorized `CanonicalOwnershipOwner.ReadAuthorizedCanonicalOwnership` returns canonical product owner + registry source/version + Settings identity state for the ONE currently authorized selected row. JSON property `selected_canonical_ownership` and escaped HTML show separate labels. | Real Ownership Register/Settings adapter, owner-period audit, relation changes |
| Unknown/unresolved identity | No canonical register lookup; cannot derive product owner from Incident/Task type | Nil register stays unknown. `identity-unresolved` explicit, never optimistically resolved or promoted into a right. | Real identity resolver/tenancy policy |
| Restrict linked metadata | Session + full owner authorization before selected-row Factors and SLA | Independent source authorization required before canonical product owner lookup; denied/foreign/out-of-page/unselected => **zero** register calls, invalid claims => generic 503 without selected row or provider message | Production IAM and source-specific permissions |

This is a *host source* contract, NOT a new IAM or complete authoritative register implementation. The host must independently authorize every returned Ownership Register/Settings field in the active principal, tenant, environment and exact Incident/Task ID. Client can only select a currently visible authorized item, never choose product owner or identity source. Source data is bounded and rejects control/newline, invalid Unicode, unsupported identity statuses and claims lacking register provenance. Native `html/template` autoescapes text with XSS-shaped input. No new routes, provider calls for unselected views, third-party dependencies, AI proposals, assignment writes, audit events or source ownership changes.

**Functional tests:** selected JSON retains operational owner distinctly from register product owner, independent current-session request args and one source lookup, unselected/denied/foreign item no call, malformed/denied register generic 503 no partial data, safe HTML rendering, missing adapter explicitly unknown. Other six views, selected factors, SLA, cursor authorization and Change-0 flows remain unchanged.

**Delivery:** DRAFT child PR, not merged/accepted/production. Exact SHA three-OS Go command-work-queue-runtime and mandatory engineering-kernel including unchanged changed-security-critical >=90% gate required; no forced merge. This code/runtime is provisional source adapter placement; accepted ownership and durable audit remain under canonical owners. A single independent read-only L2 on this new source trust boundary after stable CI, not a L2 for every commit. Existing #175 L2 must not be duplicated.
