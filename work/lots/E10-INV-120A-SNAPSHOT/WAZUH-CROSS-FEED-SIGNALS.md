# E10-INV-120A — authorized cross-feed agent/rule co-occurrence

## Functional increment on PR #49

Three or four operator-selected Wazuh 4.x offline alert feeds can share alert IDs, agents and rule identifiers. Previously the local analyst showed authorized top agents and top rules separately, without distinguishing a repeated rule on the *same* agent across independent feeds from repeats within one feed.

The existing multi-feed loader now attaches a trusted **1..4 feed ordinal** to its validated, feed-namespaced metadata. After the **existing Case → Timeline → typed-source ACL pagination** completes, the shared Wazuh triage evaluator groups authorized imported entries by validated numeric **(agent ID, rule ID)**. A signal appears only when that same pair occurs in **2+ distinct feeds**. Output is bounded to ten deterministic signals sorted by distinct feeds then alert count and IDs, with optional maximum Wazuh-reported severity. The loopback UI displays the cooccurrence as escaped text via `textContent`, not HTML. It retains the existing per-agent/per-rule drilldown and source-preview reauthorization; there is no independent permission path or external data access.

**Semantics:** same agent/rule in two source files is merely a repeated *offline detection observation*. It does not prove a causally related attack, disjoint physical sensors, unique alerts or trustworthy agent identity. It is never merged into Shared-owned Timeline Entry, a production finding, or a persistence record.

**Reuse-first:** reuse Go's `map`, `sort.Slice`, existing Wazuh 4.x normalization, in-memory owner authorization and current authenticated analyst page. No index database, graph engine, paid API, background workers or new generic condition framework. First-party memory cost proportional to visible imported entries (max 5000). The owner snapshot and original alert files remain read-only.

**Acceptance:** three-feed import into one authorized Case yields one pair (agent 001, rule 5710) with three alerts in exactly two feeds and severity 12, while identical rule on another agent does not match. Single-feed mode produces no multi-feed signal. Unauthorized principal returns zero bytes and source descriptions, alert IDs and bearer tokens do not appear. Exact-SHA Linux/macOS/Windows/runtime and engineering-kernel CI must pass before considering merge.

**Boundaries:** no live Wazuh connector, deployed IAM, durable Case storage, external tenant entitlement provider, custody guarantee or production telemetry correlation. Product Spec, decision `OPEN-013`, main and PR #36 are unchanged. Original PR #49 L2 on its earlier functional snapshot is not an approval of this material feature and a risk-based review should examine this increment after CI.
