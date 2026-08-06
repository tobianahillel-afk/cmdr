---
id: investigate-screen-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-010
  - REQ-INV-006
---
# Screen capability map — Detection Engineering and Threat Intelligence Foundations

| Existing screen or surface | Primary owner | Threat Intelligence consumption | Change |
|---|---|---|---|
| Entity Graph | Investigate / Shared Graph | CAP-INV-508,510,511,514,516 relationship and entity-candidate projections | map/links only; Entity/Graph ownership unchanged |
| Event Search | Investigate / Shared Query | CAP-INV-501,506,507,510,513 source-event and Sighting pivots | map/links only |
| Hunt Workspace | Investigate | CAP-INV-501..503,513,518 origin and return | map/links only |
| Case Workspace | Investigate | CAP-INV-501..503,508,513,518 context and handoff | map/links only |
| Evidence Board / Hypotheses and Findings | Investigate | source Evidence/Findings and contradictions | read/link only; Material ≠ Evidence |
| Detection Engineering / Technical Workbench | Investigate / Design System | CAP-INV-501,507,512,518 sourced handoffs | no detailed composition |
| Analysis Workbench technical surfaces | Investigate | CAP-INV-501,506,509,510,512,518 technical handoffs | no rewrite |
| Platform Settings Integrations / Sources & Parsers | Settings | CAP-INV-504..506 source/access/parser/health projections | read/request only |
| Platform Settings Tenants / Environments / Health | Settings | tenant, scope, retention, access and health | read only |
| Studio Control Room / Builder / Library | Studio | Tool, Tool Call, Workflow and Automation Run provenance | ownership unchanged |
| Govern Decisions / Approvals | Govern | future source access or external-sharing authority | read/link only; no sharing now |
| Global Search | Shared | authorized project/source/knowledge retrieval | map/links only |
| Inspector / Context Bar | Design System / Experience | selected source/candidate/project and restrictions | consumption only |
| Graph / Timeline / Comparison components | Shared / Design System | relationship, Sighting, version and contradiction views | alternative tabular view required |
| Intelligence Library autonomous screen | none | CAP-INV-501..518 | absence recorded; no Screen ID |
| Threat Intelligence Workspace autonomous screen | none | CAP-INV-501..518 | absence recorded; no Screen ID |
| Indicator Explorer / Watchlist screen | none | out of scope 4B.3B.1 | no Screen ID |
| Intelligence Reports / Dissemination surface | none | future 4B.3B.2 | not started |

Required surfaces examined: **18**. Detailed screen rewrites: **0**. Screen specifications modified: **0**. New Screen IDs: **0**. Wireframes, final buttons, columns, filters, animations, shortcuts, graph layout, exchange syntax and vendor branding: **0**.
