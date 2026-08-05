---
id: investigate-detection-engineering-command-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Command boundaries

| Investigate authoring concept | Command runtime concept | Boundary |
|---|---|---|
| Detection Content Draft | Detection | draft is not runtime output |
| Proposed severity/priority | Operational severity/priority | proposal has no runtime effect |
| Replay match | Detection / Signal | historical result creates neither |
| Match Review | Signal/Alert triage | candidate quality review does not triage runtime work |
| Review Package | Production decision/deployment | package is neither Approval nor deployment |
| Coverage Assessment | Command readiness/coverage projection | analytical mapping is not operational effectiveness |

Command retains Detection, Signal, Alert, Incident, triage, prioritization and coordination. No runtime object is fabricated for an undeployed draft.
