# E10-INV-120C-MITREFIX — Correct source-owned ATT&CK cohort consistency

**Review source:** independent, read-only L2 review of PR #187 `E10-INV-120B-MITRE` at frozen commit `53544c2a65f66324d004d48427b07493f7ac2d9e`; verdict `ACCEPTABLE_AVEC_RÉSERVES`, with reproducible R2-01 `analyst.js:102–119` and R2-02 `wazuh_levels.go:96–130`. This is a **new downstream candidate**, NOT a rewrite of the reviewed files or an approval for #187.

## Four distinct owners and source-to-product contract

- Product: CAP-INV-110 Investigation Timeline under Investigate Cases and Evidence. Case remains owned by Investigate; Timeline Entry and the read-only mechanism are Shared. Local Wazuh `rule.mitre.id` is **source-provided metadata** only, not a verified attack, Finding, Evidence, canonical permission, or new taxonomy owner.
- Engineering: existing epic `E10-INV-110`; child `E10-INV-120C-MITREFIX`, explicitly after `E10-INV-120B-MITRE IMPLEMENTED`. Does NOT upgrade the separate stale `E10-INV-120A-SNAPSHOT READY` or claim complete CAP-INV-110 acceptance.
- Code: existing Go first-party source parser `wazuh_levels.go` and local analyst `analyst.js`, with their existing Go/Node tests. No independent IAM, client-side enrichment engine, new route or storage.
- Git: branch `fix/e10-inv-120c-wazuh-mitre-consistency` starts from reviewed PR #187 HEAD `53544c2a65f66324d004d48427b07493f7ac2d9e`. Git parent is not engineering/product parent. Only this new branch may move; PR #187 and other audited SHAs remain immutable.

## Behavior and adversarial proof

**R2-01:** A malformed local HTTP response could claim one imported alert, contradictory unmapped counts, duplicate `T1059` entries, inconsistent total or an incorrectly ordered top-ten, and the UI would render it. The changed browser checks that all source IDs are distinct and lexical tie-sorted after descending count; `top.length = min(10, distinct IDs)`, at most 16 distinct tags per mapped imported alert, each tag count is <= mapped alerts, and for **non-truncated** top lists all mapped alerts appear under some tag. The sum of technique counts can exceed imported alert count because tags on one alert overlap. All validation occurs before changing DOM. Legacy receipts with **all** additive fields absent remain supported. Test the malformed response and the legitimate overlapping-count case in Node 24.

**R2-02:** JSON `rule.mitre.id:null` is not the same as the optional field being missing; decoding `null` into Go `[]string` silently produced `nil` and a misleading unmapped alert. The Go parser now first reads the source `mitre` and `id` as raw JSON, rejects explicit `null`, unexpected arrays/scalars/objects for the source container/ID list, then parses valid string-label arrays under the unchanged 16-label cap. Missing optional metadata remains genuinely unmapped; the parser fails before emitting any receipt for malformed values. Go tests cover both denied and valid cases.

**Bounded scale proof:** The Go test parses 5,000 distinct Wazuh-shaped alert metadata records carrying 16 source labels each (80,000 label occurrences) under the existing 8 MiB source budget. It is a deterministic regression at the documented bound, **not** a network throughput benchmark, live connector proof, or SLO. No time threshold invented.

## Preserved guarantees, exclusions, verification

The earlier full Case/Timeline/source-owner authorization and local bearer loopback semantics are untouched; so are input file caps, source classification, aggregation, old rule/agent UI and drilldowns, tenant and environment boundaries, source-denial behavior, and all Class-2 mutations under OPEN-013. JSON responses remain the same schema for valid sources; only invalid source data is rejected more strictly. Production Wazuh ingestion, real IAM, Shared timeline persistence, provenance custody and complete CAP acceptance remain absent.

**CI status must be filled from GitHub exact candidate HEAD:** timeline Go/Node/Go-race/vet across Linux/Windows/macOS and full engineering-kernel required checks. No earlier #187 green workflow validates changes here. Before merge, reconcile source-approved PR #187 and this PR in one dependency order, obtain required second-identity native reviews and gate results. Do not create another L2 for a minor regression-only correction unless material trust-boundary change is found.
