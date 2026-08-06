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
# Screen capability map — Detection Engineering

| Existing screen or surface | Primary owner | Detection Engineering consumption | Change |
|---|---|---|---|
| Event Search | Investigate/Shared Query | authoring source, replay and source-event pivot | map and links only |
| Hunt Workspace | Investigate | intake, Hypothesis, gaps and new-cycle return | map and links only |
| Case Workspace / Evidence Board | Investigate | source need, Evidence, Findings and return origin | map and links only |
| Technical Workbench / Builder Shell | Design System / Investigate | CAP-INV-401..435 temporary structured work surface | no detailed composition |
| Command Signal Queue | Command | CAP-INV-427/428 feedback projection | read/link only |
| Command Alert Queue | Command | alert dispositions and operational value | read/link only |
| Command Incident Workspace | Command | incident outcomes and future detection needs | read/link only |
| Platform Settings Sources & Parsers | Settings | readiness, schema and drift projection | read/request only |
| Platform Settings Tenants & Environments | Settings | target/environment projection | read/select proposal only |
| Platform Settings Health | Settings | runtime/source health projection | read/link only |
| Studio Control Room / Assurance | Studio | Automation Runs, Tools, evaluations and operational trace | ownership unchanged |
| Govern Action Center / Approvals | Govern | Change Request handoff, Decision and Approval projection | read/link only |
| Govern Runs & Rollback / Run Shell | Govern | Response Run, Result, partial target and rollback observation | read/link only |
| Detection Engineering autonomous screen | none | CAP-INV-401..435 | absence recorded; no Screen ID |
| Rule Builder autonomous screen | none | CAP-INV-406..413 | absence recorded; no Screen ID |
| Lifecycle/Release screen | none | CAP-INV-418..435 | composition deferred; no Screen ID |

Required surfaces examined: **16**. Detailed screen rewrites: **0**. Screen specifications modified: **0**. New Screen IDs: **0**. Wireframes, final buttons, columns, filters, animations, shortcuts and vendor syntax: **0**.
