---
id: investigate-detection-engineering-shared-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Shared capability consumption

| Shared Capability | Usage Detection Engineering | Données locales | Source canonique |
|---|---|---|---|
| Background Jobs | validation, replay, shadow assessment and bounded analysis progress | job refs, target scope and disposition | Shared |
| Notifications | review, completion, error, permission, approval and runtime-observation events | product event context only | Shared |
| Trace / Activity / Audit Hooks | end-to-end authoring, authority, runtime and improvement lineage | business relations and source refs | Shared |
| Versioning / Comparison / Diff | Drafts, candidates, expected/observed runtime versions, target results and proposals | product versions and comparison selections | Shared |
| Linking / Search | source, Govern, Command, Settings, Endpoint and Studio navigation | linked refs and safe filters | Shared |
| Metrics / Health projections | health, latency, throughput, errors, target coverage and freshness | assessment context and limitations | Shared/Settings owners |
| Export / Reporting | authorized review, readiness, health, performance and lifecycle outputs | package content, masking and restrictions | Shared |
| Collaboration / Comments / Assignments | authoring, formal review, ownership, investigation and follow-up | product roles and discussion context | Shared |
| Inspector / Context Bar / Timeline | selected candidate, version, target, Run, Result and return origin | authorized projections | Design System / Experience |
| Recovery | Drafts, plans, assessments, conflicts and partial target results | recovery checkpoints and prior valid state | Shared |

No Shared capability is redefined. Shared objects remain canonical and generic; Detection Engineering stores only its product-specific relations and dispositions.
