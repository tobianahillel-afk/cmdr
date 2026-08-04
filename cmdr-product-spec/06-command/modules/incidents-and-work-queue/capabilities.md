---
id: command-module-incidents-and-work-queue-capabilities
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-OBJ-001
  - REQ-OBJ-012
  - REQ-PROD-003
  - REQ-PROD-004
  - REQ-PROD-005
  - REQ-PROD-006
  - REQ-PROD-008
  - REQ-PROD-009
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-PROD-053
  - REQ-SEC-001
  - REQ-UX-005
  - REQ-UX-008
  - REQ-UX-009
---

# Capability map — Incidents and Work Queue

| ID | Capability | Primary role | Objects | Representative actions | Status | Mode | OPEN |
|---|---|---|---|---|---|---|---|
| CAP-CMD-101 | Unified Work Queue | SOC Analyst L1/L2 | Incident, Task, projections | Changer de vue (C0), filtrer (C0), ouvrir (C0) | defined | planned | — |
| CAP-CMD-102 | Work Assignment | Incident Commander | Incident/Task, Principal, Notification | Affecter (C2), prendre (C2), libérer (C2) | defined | planned | OPEN-013 |
| CAP-CMD-103 | Operational Ownership | Incident Commander | Incident/Task, Team, Audit | owner (C2), contributeur (C2), watcher (C2) | defined | planned | OPEN-013 |
| CAP-CMD-104 | Priority and Severity Coordination | SOC Analyst L2 | work item, source dimensions | inspecter (C0), priority/impact (C2) | defined | planned | OPEN-013 |
| CAP-CMD-105 | SLA Tracking | SOC Analyst L1/L2 | work item, policy, Task | consulter (C0), pause/reprise (C2) | defined | planned | OPEN-006, OPEN-013 |
| CAP-CMD-106 | Incident Coordination | Incident Commander | Signal/Alert, Case/Finding, Govern objects | créer/coordonner (C2), ouvrir Case (C2) | defined | planned | OPEN-013 |
| CAP-CMD-107 | Task Coordination | SOC Analyst | Task, source object | créer/modifier/lier (C2) | defined | planned | OPEN-013 |
| CAP-CMD-108 | Bulk Coordination | Incident Commander | work items, views, export job | preview (C0), mutations limitées (C2) | defined | planned | OPEN-013 |
| CAP-CMD-109 | Work Freshness and Staleness | SOC Analyst L1/L2 | work/projections/health | inspecter/refresh (C0), follow-up (C2) | defined | planned | — |
| CAP-CMD-110 | Escalation | Incident Commander | work item, Case/Finding/Decision | escalader (C2) | defined | planned | OPEN-013 |

Les sources individuelles possèdent les entrées, sorties et critères. `planned` ne prouve aucune disponibilité runtime.
