# E10-INV-120A — reconciled offline multi-Wazuh Case integration

## Purpose and foundation
New scoped branch from `engineering/autonomous-development-foundation` at exact SHA `ea1770b48a4bc1b2e850a7ffab15d5839f9e9e49`. Existing historical draft PR #49 is blocked by conflicts with the evolving foundation; this port does **not** rewrite it, merge it, force-push it, or discard its history. The foundation's 25-entry Wazuh pagination and single-feed source-authorization chain remain authoritative. Change scope: new offline 2..4 feed import, cross-feed cooccurrence signals, two-Case comparison/proximity and already-authorized source drilldown.

## Reuse versus alternatives
Reuses CMDR's existing Wazuh 4.x normalizer, JSONL importer, immutable operator-selected owner snapshot, Case/Timeline/typed-source ACL, loopback bearer, Go stdlib, and read-only aggregation. A new OpenSearch index, event bus, live agent connector, state store, generic rule engine or provider IAM would create new trust boundaries without being necessary for the bounded local workflow. No third-party product dependency added.

## Security and behavior
Multiple operator-selected, distinct regular Wazuh files are parsed at startup and each event receives a trustworthy feed prefix; IDs can collide in original Wazuh exports. No data is returned until Case, Timeline and source traversals succeed. User metadata/rule/agent is **not** an authorization decision. Drilldown uses full owner reauthorization, 25-row server-side pagination and a separate fresh source-open. The existing single-feed path is unchanged. Cross-feed agent/rule cooccurrence is an observational lead, never a causal security finding.

## Limits and acceptance
Lab-only, read-only local inputs; synthetic in-memory owner ACL and ephemeral bearer loopback, **not** production IAM, persistent Case store, live Wazuh collection, or custody certification. OPEN-013 remains unresolved. Exact-SHA CI with `go test -race`, `go vet`, `gofmt`, Node browser tests and mandatory engineering-kernel checks must pass before considering this increment verified. PR should stay draft pending a risk-based independent read-only L2.
