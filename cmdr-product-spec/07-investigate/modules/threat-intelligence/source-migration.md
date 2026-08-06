---
id: investigate-threat-intelligence-source-migration
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Source migration

## Audit result
No active standalone Investigate Threat Intelligence functional module, Intelligence Library, Indicator Explorer or CAP-INV-5xx existed at the starting SHA. Historical needs were distributed across Shared Threat Intelligence Enrichment, Entity/Graph, Event Search/Hunt, Cases/Evidence, technical handoffs, Detection Engineering, Settings sources/providers, Studio workflows, Govern release controls and future dissemination concepts.

## Disposition
- Competing Investigate functional sources deprecated: **0**.
- Active screen sources deprecated: **0**.
- `12-shared-capabilities/threat-intelligence-enrichment.md` remains active under Shared; it is not a competing Intelligence workspace.
- Entity/Graph/Search/Timeline remain active under Shared.
- Source/provider/feed administration remains active under Settings.
- Runtime Detection/Signal/Alert/Incident remain active under Command.
- Detection Engineering and all technical analysis modules remain active owners of their outputs.
- Future Report, dissemination, watchlist and external-sharing sources remain future 4B.3B.2 boundaries.
- Canonical replacement for Investigate foundations: `07-investigate/modules/threat-intelligence/`.

## Acceptance
Only one Investigate functional architecture remains active; no screen or owner source is deprecated merely because it is consumed by Threat Intelligence.
