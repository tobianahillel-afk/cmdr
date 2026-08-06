---
id: investigate-screen-capability-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements: [REQ-UX-001, REQ-UX-010, REQ-INV-006]
---
# Screen capability map — Threat Intelligence complete functional scope

| Existing screen or surface | Consumption | Change in 4B.3B.2 |
|---|---|---|
| Entity Graph / Graph / Timeline | candidates, assessments, relationships, Sightings and versions | links/map only; owners unchanged |
| Event Search / Hunt Workspace / Case Workspace | origins, pivots, handoffs, feedback and return | no detailed rewrite |
| Evidence Board / Hypotheses and Findings | Evidence/Findings and contradictions | read/link only; Intelligence Hypothesis remains distinct |
| Detection Engineering / Technical Workbench | CAP-INV-531/534/535/537 handoffs | no rule/runtime UI defined |
| Analysis Workbench surfaces | technical inputs to CAP-INV-519/524 | no rewrite |
| Platform Settings Sources/Integrations/Tenants/Health | source/access/destination projections and requests | read/request only |
| Studio Control Room / Builder / Library | Tool/Run/Workflow provenance and monitoring requests | ownership unchanged |
| Govern Decisions / Approvals / Action Center | external-release/revocation Action Requests | no sharing execution UI defined |
| Shared Reporting / Export / Notifications | rendering, authorized export and delivery mechanisms | consumed, not redefined |
| Inspector / Context Bar | selected object, restrictions, audience and return origin | consumption only |
| Threat Intelligence Workspace / Intelligence Library / Product Builder / Indicator Explorer | future detailed surfaces | capability links recorded; **no Screen ID** |

Required surfaces read: **18**. Screen specs modified: **0**. Detailed rewrites: **0**. New Screen IDs, wireframes, final buttons/columns/filters/animations/shortcuts and external protocol UI: **0**.
