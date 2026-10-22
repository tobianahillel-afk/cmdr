---
id: unresolved-decisions-govern-gov3-audit
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-09
source-of-truth: decision-audit-addendum
---
# Open Decision Audit — Govern GOV-3

The canonical `unresolved-decisions.md` remains unchanged. The repository still contains **18 OPEN decisions** and `OPEN-009` remains the only historically resolved item. GOV-3 creates **0** new OPEN and closes **0**.

## GOV-3 relevant decisions
| OPEN | Relevance | GOV-3 treatment | Status after GOV-3 |
|---|---|---|---|
| OPEN-007 | Human Gate vs Govern Approval/Decision | audit/metrics preserve non-equivalence and source ownership | OPEN |
| OPEN-008 | platform/source availability | audit/metrics distinguish unavailable/partial sources; no provider selected | OPEN |
| OPEN-010 | final role/activity density | no detailed Audit Trail/Response Metrics screen rewrite | OPEN |
| OPEN-013 | default governance for reversible class-2 mutations | GOV-3 annotations/packages/assessments remain classified without final default policy | OPEN |
| OPEN-014 | Artifact/Attachment/material relationships | Audit Evidence Package remains distinct from canonical Evidence | OPEN |
| OPEN-015 | Automation Run / Response Run bridge | audit reconstruction preserves correlation without merging objects | OPEN |
| OPEN-019 | dissemination/releasability/sharing | audit/report/export preparation never grants external sharing authority | OPEN |

## Not resolved by GOV-3
GOV-3 deliberately does not select an audit engine, event format, log pipeline, metrics engine, warehouse, storage schema, immutable-ledger implementation, cryptographic proof mechanism, statistical/anomaly method, universal thresholds, final retention policy, final RBAC/ABAC, external-disclosure policy or implementation architecture.

Existing implementation/delivery OPEN decisions may remain open while the provider-neutral Govern capability specification closes, provided all mandatory behavior is specified and no blocking active contradiction remains.