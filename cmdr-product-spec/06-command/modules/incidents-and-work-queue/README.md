---
id: command-module-incidents-and-work-queue
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

# Incidents and Work Queue

## Mission

Recevoir, organiser, affecter, prioriser et suivre Incidents et Tasks dans une Work Queue unique.

## Utilisateurs

- SOC Analyst L1/L2
- Incident Commander
- Team Lead
- Business Owner en lecture

## Ownership

**Le module possède :** Incident, Task opérationnelle, six vues Work Queue, assignment, ownership opérationnel, coordination SLA/priorité/état.

**Le module ne possède pas :** Case/Evidence/Finding, Decision/Response Run, Team Load comme vue système ou containment direct.

## Capabilities

| ID | Capability | Delivery status | Delivery mode | Source |
|---|---|---|---|---|
| CAP-CMD-101 | Unified Work Queue | defined | planned | capabilities/unified-work-queue.md |
| CAP-CMD-102 | Work Assignment | defined | planned | capabilities/work-assignment.md |
| CAP-CMD-103 | Operational Ownership | defined | planned | capabilities/operational-ownership.md |
| CAP-CMD-104 | Priority and Severity Coordination | defined | planned | capabilities/priority-and-severity-coordination.md |
| CAP-CMD-105 | SLA Tracking | defined | planned | capabilities/sla-tracking.md |
| CAP-CMD-106 | Incident Coordination | defined | planned | capabilities/incident-coordination.md |
| CAP-CMD-107 | Task Coordination | defined | planned | capabilities/task-coordination.md |
| CAP-CMD-108 | Bulk Coordination | defined | planned | capabilities/bulk-coordination.md |
| CAP-CMD-109 | Work Freshness and Staleness | defined | planned | capabilities/work-freshness-and-staleness.md |
| CAP-CMD-110 | Escalation | defined | planned | capabilities/escalation.md |

## Shared Capabilities consommées

| Shared Capability | Usage local | Source canonique |
|---|---|---|
| Saved Views | appliquer/partager une configuration; catalogue système Command | `../../../12-shared-capabilities/saved-views.md` |
| Search / Command Palette | retrouver et ouvrir objets/actions autorisés | `../../../12-shared-capabilities/global-search.md` / `../../../04-experience-architecture/command-palette.md` |
| Notifications / Jobs | assignment, SLA, escalade, export et progression | `../../../12-shared-capabilities/notification-center.md` / `background-jobs.md` |
| Activity / Trace / Timeline | mutations, provenance et chronologie | Design System / `timeline-engine.md` |
| Reporting / Export | citations, snapshots, redaction et audit | `reporting-engine.md` / `export-engine.md` |
| Inspector / Context Bar | sélection et contexte | Design System sources canoniques |
| Object Links / Collaboration / Presence | références, comments, mentions et conflits | Shared sources canoniques |
| Versioning / Localization / Audit Hooks | concurrence, formats et trace | Domain/Shared sources canoniques |

Aucune de ces capacités partagées n’est redéfinie localement.

## États et erreurs

Les états métier sont propres aux capabilities. Loading, Empty, Partial, Error, Offline, Permission denied et Stale suivent le Design System et ne masquent jamais une source manquante.

## IA et automatisation

Toutes les fonctions essentielles ont une voie manuelle ou déterministe. Une suggestion reste une proposition attribuée ; elle n’est jamais une priorité, un Finding ou une Decision effective.

## Critère d’acceptation

**Given** le module et un utilisateur autorisé, **When** une capability est utilisée, **Then** owner, objets, entrées, sorties, classes, états, dépendances et alternative sans IA sont résolus depuis sa source canonique.
