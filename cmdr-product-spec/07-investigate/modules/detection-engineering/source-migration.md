---
id: investigate-detection-engineering-source-migration
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Source migration

## Audit result
No active standalone Detection Engineering authoring module or canonical Rule Builder functional source existed before this phase. Requirements were distributed across Event Search, Hunt, Case/Evidence/Findings, technical handoffs, Command runtime objects, Endpoint Agent detection, Settings sources/parsers, Shared query/correlation/normalization, Studio datasets/evaluations and user journeys.

## Disposition
- Competing functional sources deprecated: **0**.
- Active screens deprecated: **0**.
- Endpoint Agent runtime detection/update/suppression sources remain active for their owner and future 4B.3A.2.
- Runtime Detection and Detection-to-Signal sources remain active under Command/Product Architecture.
- Event Search and Query Authoring remain active and distinct from Detection Content.
- Studio datasets/evaluations remain active as generic capabilities.
- New canonical replacement for authoring: `07-investigate/modules/detection-engineering/`.

## Acceptance
The canonical module absorbs authoring needs without changing runtime ownership, choosing a language/engine or starting deployment and Intelligence scope.
