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
# Screen capability map — Detection Engineering Authoring

| Existing screen or surface | Primary capability | Detection Authoring consumption | Change |
|---|---|---|---|
| Event Search | CAP-INV-002 | CAP-INV-401,404,405,414,415 | map and links only; Search ≠ Detection Content |
| Hunt Workspace / signals-and-hunt surface | CAP-INV-005 | CAP-INV-401,403,412,416 | map and links only |
| Case Workspace | CAP-INV-102 | CAP-INV-401..403,410,417 | map and links only |
| Evidence Board / Hypotheses and Findings | CAP-INV-103/107..109 | CAP-INV-401,403,415,417 | map and links only |
| Technical Workbench / Builder Shell | layout/surface | CAP-INV-401..417 | temporary authoring surface pending screen phase |
| Platform Settings Sources and Parsers | Settings surface | CAP-INV-404,405,409,416 | read/request transition only |
| Studio Builder / Assurance / Evaluation / Control Room | Studio surfaces | Tools, Datasets, Evaluations and Runs | ownership unchanged |
| Command Incident / future Detection surface | Command surface | CAP-INV-401 source and future 417 destination only | no runtime Detection fabricated |
| Detection Engineering autonomous screen | none | CAP-INV-401..417 | absence recorded; no Screen ID created |
| Rule Builder autonomous screen | none | CAP-INV-406..413 | absence recorded; no Screen ID created |

Detailed screen rewrites: **0**. Screen specifications modified: **0**. New Screen IDs: **0**. Wireframes, final buttons, columns, filters, dimensions, animations, shortcuts and rule syntax: **0**.
