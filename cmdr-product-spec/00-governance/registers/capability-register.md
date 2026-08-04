---
id: capability-register
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-04
source-of-truth: registry
requirements:
  - REQ-PROD-006
  - REQ-PROD-012
  - REQ-PROD-019
  - REQ-PROD-013
---

# Capability Register

## Objet

Le registre fournit l’index contrôlable des capabilities. Le fichier canonique propriétaire décrit le comportement ; le registre détecte l’absence de document, d’owner, d’utilisateur, d’objet, de Requirement ID, de preuve de delivery ou les owners concurrents.

## Convention

| Préfixe | Produit |
|---|---|
| `CAP-CMD-` | Command |
| `CAP-INV-` | Investigate |
| `CAP-GOV-` | Govern |
| `CAP-STD-` | CMDR Studio |
| `CAP-SET-` | Platform Settings |
| `CAP-EPT-` | Endpoint Agent |
| `CAP-SHR-` | Shared Capabilities |

Un identifiant est unique, immuable, indépendant du chemin et non recyclé. Les capacités des Phases 4B–4D ne reçoivent pas d’ID prématuré dans la Phase 4A.

## Statuts

- `delivery_status` : maturité fonctionnelle (`defined`, `partial`, `planned`, `proposed`, `out-of-scope`) ;
- `delivery_mode` : preuve de livraison actuelle (`native`, `integrated`, `temporary-integration`, `planned`, `out-of-scope`) ;
- `status` : cycle documentaire. Aucun document ci-dessous n’est `validated` ou `implemented`.

## Command — Phase 4A

| capability_id | name | owner_product | owner_module | status | delivery_status | delivery_mode | canonical_file | primary_roles | primary_objects | consumers | requirement_ids | open_decisions | dependencies | supersedes | last_reviewed |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| CAP-CMD-001 | Situation Overview | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/situation-overview.md` | Incident Commander, SOC Analyst L1/L2, Business Owner | Incident, Task, Decision / Response Run / Result, Service | Mission Control — Now, Mission Control — Situation, futur parcours de handover | REQ-PROD-013, REQ-PROD-008, REQ-PROD-010, REQ-PROD-021 | — | CAP-CMD-003, CAP-CMD-005, CAP-CMD-006, Shared Global Search | — | 2026-08-04 |
| CAP-CMD-002 | Priority Management | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/priority-management.md` | Incident Commander, SOC Analyst L2, Business Owner | Incident / Task, Service, SLA context, Audit trail | Mission Control — Priorities, Unified Work Queue, Incident Detail | REQ-PROD-003, REQ-PROD-010, REQ-PROD-013, REQ-PROD-021 | OPEN-013 | CAP-CMD-104, CAP-CMD-105, CAP-CMD-204, CAP-CMD-205 | — | 2026-08-04 |
| CAP-CMD-003 | Situation Timeline | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/situation-timeline.md` | Incident Commander, SOC Analyst L2, Auditor | Timeline Entry, Incident, Decision / Run / Result | Mission Control — Situation, Incident Detail, Handover | REQ-PROD-005, REQ-PROD-008, REQ-PROD-013, REQ-UX-007 | — | Timeline Engine, Object Linking Service, Trace, Export Engine | — | 2026-08-04 |
| CAP-CMD-004 | Handover | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/handover.md` | Incident Commander, SOC Team Lead, SOC Analyst L2 | Incident / Task, Decision / Response Run / Result, Case / Finding, Handover record, Incident / Task ownership | Mission Control — Handover, parcours de relève, Incident Detail | REQ-PROD-008, REQ-PROD-013, REQ-PROD-021, REQ-PROD-010 | — | CAP-CMD-001, CAP-CMD-005, CAP-CMD-006, Notification Center | — | 2026-08-04 |
| CAP-CMD-005 | Operational Blockers | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/operational-blockers.md` | Incident Commander, SOC Analyst L2, Task owner | Incident / Task, Decision / Case / external dependency, Task | Mission Control, Unified Work Queue, Handover | REQ-PROD-006, REQ-PROD-009, REQ-PROD-013, REQ-PROD-021 | OPEN-013 | CAP-CMD-110, CAP-CMD-107, Object Linking Service, Notification Center | — | 2026-08-04 |
| CAP-CMD-006 | Recent Results and Outcomes | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/recent-results-and-outcomes.md` | Incident Commander, Response Operator, SOC Analyst L2 | Result / Response Run / Decision, Incident, Timeline | Mission Control — Now, Mission Control — Situation, Incident Detail | REQ-PROD-005, REQ-PROD-008, REQ-PROD-013, REQ-PROD-021 | OPEN-013 | CAP-CMD-001, CAP-CMD-003, CAP-CMD-106, Object Linking Service | — | 2026-08-04 |
| CAP-CMD-101 | Unified Work Queue | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/unified-work-queue.md` | SOC Analyst L1/L2, Incident Commander, Team Lead | Incident, Task, Case / Decision / Response Run, Saved View application, Incident / Task | Work Queue, Incident Detail, Mission Control | REQ-OBJ-001, REQ-OBJ-012, REQ-PROD-013, REQ-UX-008, REQ-UX-009 | — | CAP-CMD-102, CAP-CMD-104, CAP-CMD-105, CAP-CMD-108 | — | 2026-08-04 |
| CAP-CMD-102 | Work Assignment | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/work-assignment.md` | Incident Commander, SOC Analyst, Team Lead | Incident / Task, Principal / Role, Notification | Unified Work Queue, Incident Detail, Mission Control | REQ-PROD-003, REQ-PROD-013, REQ-PROD-021, REQ-SEC-001 | OPEN-013 | Identity projection, Notification Center, Collaboration Service, CAP-CMD-103 | — | 2026-08-04 |
| CAP-CMD-103 | Operational Ownership | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/operational-ownership.md` | Incident Commander, SOC Analyst, Team Lead | Incident / Task, Case / Decision / Run, Principal / Team, Audit | Work Queue, Incident Detail, Mission Control | REQ-PROD-006, REQ-PROD-009, REQ-PROD-013, REQ-OBJ-001 | OPEN-013 | CAP-CMD-102, Ownership Register, Identity projections, Object Linking Service | — | 2026-08-04 |
| CAP-CMD-104 | Priority and Severity Coordination | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/priority-and-severity-coordination.md` | SOC Analyst L2, Incident Commander, Business Owner | Incident / Task, Signal / Alert, Service / SLA, Recommendation disposition | Unified Work Queue, Mission Control — Priorities, Incident Detail | REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-UX-005 | OPEN-013 | CAP-CMD-002, CAP-CMD-105, CAP-CMD-204, CAP-CMD-205 | — | 2026-08-04 |
| CAP-CMD-105 | SLA Tracking | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/sla-tracking.md` | SOC Analyst L1/L2, Incident Commander, Service Delivery Manager | Incident / Task, SLA policy / engagement, Task | Unified Work Queue — SLA Risk, Mission Control, Handover | REQ-PROD-005, REQ-PROD-013, REQ-PROD-053, REQ-OBJ-012 | OPEN-006, OPEN-013 | CAP-CMD-101, CAP-CMD-102, CAP-CMD-110, Notification Center | — | 2026-08-04 |
| CAP-CMD-106 | Incident Coordination | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/incident-coordination.md` | Incident Commander, SOC Analyst L1/L2, Business Owner | Signal / Alert, Case / Finding, Decision / Response Run / Result, Incident / Task, Incident, Task, Action Request | Incident Detail, Unified Work Queue, Mission Control | REQ-OBJ-001, REQ-PROD-003, REQ-PROD-008, REQ-PROD-013, REQ-SEC-001 | OPEN-013 | CAP-CMD-102, CAP-CMD-104, CAP-CMD-107, CAP-CMD-110 | — | 2026-08-04 |
| CAP-CMD-107 | Task Coordination | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/task-coordination.md` | SOC Analyst, Incident Commander, Readiness Coordinator | Task, Source object, Source relation | Unified Work Queue, Incident Detail, Readiness Improvement Actions | REQ-PROD-009, REQ-PROD-013, REQ-PROD-021, REQ-OBJ-012 | OPEN-013 | CAP-CMD-102, CAP-CMD-103, CAP-CMD-005, Object Linking Service | — | 2026-08-04 |
| CAP-CMD-108 | Bulk Coordination | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/bulk-coordination.md` | Incident Commander, SOC Analyst L2, Team Lead | Incident / Task, Saved View / filters, Export job | Unified Work Queue, audit, operations reporting | REQ-PROD-004, REQ-PROD-009, REQ-PROD-013, REQ-SEC-001 | OPEN-013 | CAP-CMD-101, CAP-CMD-102, CAP-CMD-104, Background Jobs | — | 2026-08-04 |
| CAP-CMD-109 | Work Freshness and Staleness | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/work-freshness-and-staleness.md` | SOC Analyst L1/L2, Incident Commander, Auditor | Incident / Task, Case / Decision / Run projections, Health, Task, Acknowledgement | Unified Work Queue, Mission Control, Handover | REQ-PROD-005, REQ-PROD-008, REQ-PROD-013, REQ-UX-005 | — | Data Quality Service, Background Jobs, Notification Center, CAP-CMD-107 | — | 2026-08-04 |
| CAP-CMD-110 | Escalation | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/escalation.md` | Incident Commander, SOC Analyst L2, Team Lead | Incident / Task, Case / Finding / Decision, Escalation record/state, Case or Action Request | Incident Detail, Unified Work Queue, Mission Control | REQ-PROD-003, REQ-PROD-004, REQ-PROD-008, REQ-PROD-013, REQ-SEC-001 | OPEN-013 | CAP-CMD-005, CAP-CMD-106, Notification Center, Collaboration Service | — | 2026-08-04 |
| CAP-CMD-201 | Service Context | Command | risk-and-coverage | draft | defined | planned | `06-command/modules/risk-and-coverage/capabilities/service-context.md` | Incident Commander, Business Owner, SOC Analyst L2 | Service, Incident / Task, Exposure / coverage projections, Data quality follow-up | Risk and Coverage, Incident Detail, Mission Control | REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-032 | OPEN-013 | Business Service Catalog, Object Linking Service, CAP-CMD-204, CAP-CMD-106 | — | 2026-08-04 |
| CAP-CMD-202 | Exposure Overview | Command | risk-and-coverage | draft | defined | planned | `06-command/modules/risk-and-coverage/capabilities/exposure-overview.md` | Incident Commander, SOC Analyst L2, Business Owner | Exposure projection, Service, Incident, Incident / Task, Exposure source | Risk and Coverage, Mission Control, Incident Detail | REQ-PROD-006, REQ-PROD-013, REQ-PROD-032, REQ-PROD-037 | OPEN-013 | Business Service Catalog, Object Linking Service, Data Quality Service, CAP-CMD-201 | — | 2026-08-04 |
| CAP-CMD-203 | Coverage Overview | Command | risk-and-coverage | draft | defined | planned | `06-command/modules/risk-and-coverage/capabilities/coverage-overview.md` | Incident Commander, Detection Engineer en consultation, Platform Administrator en consultation | Coverage projections, Service, Task, Coverage source | Risk and Coverage, Readiness Overview, Mission Control | REQ-PROD-005, REQ-PROD-013, REQ-PROD-032, REQ-PROD-037 | OPEN-013 | CAP-CMD-201, CAP-CMD-301, CAP-CMD-303, Metrics Engine | — | 2026-08-04 |
| CAP-CMD-204 | Business Impact Context | Command | risk-and-coverage | draft | defined | planned | `06-command/modules/risk-and-coverage/capabilities/business-impact-context.md` | Incident Commander, Business Owner, Service Owner | Incident, Service, Finding/Result, Action Request context | Incident Detail, Mission Control, Priority Management | REQ-PROD-003, REQ-PROD-005, REQ-PROD-013, REQ-PROD-021 | OPEN-013 | CAP-CMD-201, CAP-CMD-205, CAP-CMD-106, Business Service Catalog | — | 2026-08-04 |
| CAP-CMD-205 | Risk Prioritization Context | Command | risk-and-coverage | draft | defined | planned | `06-command/modules/risk-and-coverage/capabilities/risk-prioritization-context.md` | Incident Commander, SOC Analyst L2, Business Owner | Incident / Task, Service / Exposure / Coverage, Confidence / severity, Priority recommendation, Incident/Task | Mission Control — Priorities, Unified Work Queue, Incident Detail | REQ-PROD-003, REQ-PROD-010, REQ-PROD-013, REQ-PROD-021 | OPEN-013 | CAP-CMD-002, CAP-CMD-104, CAP-CMD-105, CAP-CMD-201 | — | 2026-08-04 |
| CAP-CMD-301 | Readiness Overview | Command | readiness-and-operations | draft | defined | planned | `06-command/modules/readiness-and-operations/capabilities/readiness-overview.md` | Readiness Coordinator, Incident Commander, Team Lead | Capability projection, Task, Plan/Exercise records, Readiness assessment | Readiness & Operations screen, Mission Control, reports | REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-057 | OPEN-010, OPEN-013 | CAP-CMD-203, CAP-CMD-303, CAP-CMD-304, CAP-CMD-305 | — | 2026-08-04 |
| CAP-CMD-302 | Exercise Coordination | Command | readiness-and-operations | draft | defined | planned | `06-command/modules/readiness-and-operations/capabilities/exercise-coordination.md` | Readiness Coordinator, Incident Commander, Team Lead | Operational Plan, Capability Readiness, Workflow/Simulation, Exercise record, Task | Readiness & Operations, reporting, future exercise journey | REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-057 | OPEN-010, OPEN-013 | CAP-CMD-301, CAP-CMD-303, CAP-CMD-304, Collaboration Service | — | 2026-08-04 |
| CAP-CMD-303 | Improvement Actions | Command | readiness-and-operations | draft | defined | planned | `06-command/modules/readiness-and-operations/capabilities/improvement-actions.md` | Readiness Coordinator, Incident Commander, Task owner | Result, Exercise/readiness/coverage, Task, Source relation | Readiness Overview, Unified Work Queue, Mission Control blockers | REQ-PROD-009, REQ-PROD-013, REQ-PROD-021, REQ-OBJ-012 | OPEN-013 | CAP-CMD-107, CAP-CMD-301, CAP-CMD-302, CAP-CMD-203 | — | 2026-08-04 |
| CAP-CMD-304 | Operational Plans | Command | readiness-and-operations | draft | defined | planned | `06-command/modules/readiness-and-operations/capabilities/operational-plans.md` | Readiness Coordinator, Incident Commander, Business Owner | Plan/Playbook references, Capability Readiness, Service/Incident, Operational Plan record, Task | Readiness Overview, Exercise Coordination, Incident Detail | REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-057 | OPEN-010, OPEN-013 | CAP-CMD-301, CAP-CMD-302, CAP-CMD-305, Object Linking Service | — | 2026-08-04 |
| CAP-CMD-305 | Capability Readiness | Command | readiness-and-operations | draft | defined | planned | `06-command/modules/readiness-and-operations/capabilities/capability-readiness.md` | Readiness Coordinator, Incident Commander, Product Owner | Capability, Health/Assurance/Exercise, Task, Readiness assessment | Readiness Overview, Operational Plans, Mission Control | REQ-PROD-012, REQ-PROD-013, REQ-PROD-019, REQ-PROD-057 | OPEN-010, OPEN-013 | Capability Register, CAP-CMD-301, CAP-CMD-303, Platform Health | — | 2026-08-04 |
| CAP-CMD-401 | Customers and Delivery Context | Command | customers-and-delivery | draft | proposed | planned | `06-command/modules/customers-and-delivery/capabilities/customers-and-delivery-context.md` | Service Delivery Manager, Customer Success Manager, Incident Commander | Incident / Task, Result, Report, Customer/engagement context, Report request/draft context, Task, Customer/contract source | Customer & Reports legacy screen, Reporting Engine, optional service delivery journeys | REQ-PROD-053, REQ-PROD-013, REQ-PROD-019, REQ-PROD-033 | OPEN-006, OPEN-013 | OPEN-006, Reporting Engine, Metrics Engine, Export Engine | — | 2026-08-04 |

## Contrôles de registre

Une entrée échoue si son ID est dupliqué/recyclé, fichier/owner/rôle/objet/Requirement absent, delivery claim sans preuve, ou owner concurrent.

## Couverture Phase 4A

- IDs Command : **27** ; documents présents : **27/27** ; owners/utilisateurs/objets/Requirement IDs absents : **0** ;
- delivery mode `planned` : **27** ; status `defined` : **26** ; status `proposed` : **1** (`CAP-CMD-401`) ;
- claims current `native`/`integrated`/`temporary-integration` : **0** ; owners concurrents actifs : **0** après alignement Task.

## Phases suivantes

Investigate, Govern, Studio, Settings, Endpoint Agent et Shared recevront leurs IDs substantiels uniquement dans les Phases 4B–4D.
