---
id: investigate-threat-intelligence-command-detection-boundaries
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
---
# Command and Detection Engineering boundaries

| Threat Intelligence concept | Owner projection | Boundary |
|---|---|---|
| Indicator Candidate | Detection Content | candidate ≠ rule; handoff creates no Draft automatically |
| Sighting | Signal / Alert / Incident | observation ≠ operational object or compromise |
| Threat Entity Candidate | Entity | candidate ≠ canonical identity or attributed actor |
| Campaign/Cluster candidate | Incident | analytic grouping ≠ operational situation |
| TTP Mapping | Detection Hypothesis/Coverage | mapping ≠ detection logic or effective coverage |
| Analysis Handoff Package | Continuous Improvement Package | packages remain owner-specific and non-effective |
| Intelligence priority | Command operational priority | different scope and authority |

Command retains runtime Detection, Signal, Alert, Incident, triage and operational priority. Detection Engineering retains Detection Content, Hypothesis, Coverage, Gap and lifecycle. Threat Intelligence can prepare sourced candidate handoffs only.
