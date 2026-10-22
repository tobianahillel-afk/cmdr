---
id: investigate-detection-engineering-command-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Command boundaries

| Investigate concept | Command runtime concept | Boundary |
|---|---|---|
| Detection Content Draft / Release Candidate | runtime Detection | draft/candidate is not runtime output |
| Proposed severity/priority | operational severity/priority | proposal has no runtime effect |
| Historical replay or shadow match | Detection / Signal | neither creates an operational Signal |
| Runtime Version Observation | runtime Detection state | projection does not overwrite Command |
| Detection Health Assessment | runtime Detection | health is technical and distinct from quality/effectiveness |
| Runtime Quality Assessment / Production Match Review | Signal/Alert/Incident disposition | Command feedback is evidence, not absolute ground truth |
| Tuning/Suppression/Exception Proposal | active runtime change | proposal creates no active change |
| Retirement Proposal | inactive/retired runtime state | proposal is not deactivation or deletion |
| Continuous Improvement Package | new authoring cycle | package changes no Command object |

Command retains runtime Detection, Signal, Alert, Incident, Work Queue, triage, prioritization, dispositions and operational coordination. Investigate consumes projections and feedback without silently modifying or deleting them.
