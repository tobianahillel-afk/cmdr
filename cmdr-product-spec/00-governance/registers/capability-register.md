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
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-019
---

# Capability Register

## Objet

Le registre fournit l’index contrôlable des capabilities. Le fichier canonique propriétaire décrit les comportements ; le registre permet de détecter l’absence de document, owner, utilisateur, objet, entrée, sortie, Requirement ID, preuve de delivery ou les owners concurrents.

## Convention et statuts

| Préfixe | Produit |
|---|---|
| `CAP-CMD-` | Command |
| `CAP-INV-` | Investigate |
| `CAP-GOV-` | Govern |
| `CAP-STD-` | CMDR Studio |
| `CAP-SET-` | Platform Settings |
| `CAP-EPT-` | Endpoint Agent |
| `CAP-SHR-` | Shared Capabilities |

Un ID est unique, immuable, indépendant du chemin et non recyclé. `status` décrit le cycle documentaire ; `delivery_status` la maturité fonctionnelle ; `delivery_mode` la preuve de livraison actuelle. Aucun document ci-dessous n’est `validated` ou `implemented`.

## Command — Phase 4A

| capability_id | name | owner_product | owner_module | status | delivery_status | delivery_mode | canonical_file | primary_roles | primary_objects | consumers | requirement_ids | open_decisions | dependencies | supersedes | last_reviewed |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| CAP-CMD-001 | Situation Overview | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/situation-overview.md` | Incident Commander, SOC Analyst | Incident, Task, Govern projections | Mission Control | REQ-PROD-008,010,013,021 | — | CAP-CMD-003,005,006; Global Search | — | 2026-08-04 |
| CAP-CMD-002 | Priority Management | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/priority-management.md` | Incident Commander, SOC L2 | Incident, Task, SLA context | Mission Control, Work Queue | REQ-PROD-003,010,013,021 | OPEN-013 | CAP-CMD-104,105,204,205 | — | 2026-08-04 |
| CAP-CMD-003 | Situation Timeline | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/situation-timeline.md` | Incident Commander, Auditor | Timeline Entry, Incident, Govern projections | Mission Control, Incident Detail | REQ-PROD-005,008,013; REQ-UX-007 | — | Timeline Engine, Trace | — | 2026-08-04 |
| CAP-CMD-004 | Handover | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/handover.md` | Incident Commander, Team Lead | Incident, Task, Case/Govern projections | Handover, Incident Detail | REQ-PROD-008,010,013,021 | — | CAP-CMD-001,005,006; Notifications | — | 2026-08-04 |
| CAP-CMD-005 | Operational Blockers | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/operational-blockers.md` | Incident Commander, Task owner | Incident, Task, dependency relation | Mission Control, Work Queue | REQ-PROD-006,009,013,021 | OPEN-013 | CAP-CMD-107,110; Object Linking | — | 2026-08-04 |
| CAP-CMD-006 | Recent Results and Outcomes | Command | mission-control | draft | defined | planned | `06-command/modules/mission-control/capabilities/recent-results-and-outcomes.md` | Incident Commander, Response Operator | Result, Run, Decision, Incident | Mission Control, Incident Detail | REQ-PROD-005,008,013,021 | OPEN-013 | CAP-CMD-001,003,106 | — | 2026-08-04 |
| CAP-CMD-101 | Unified Work Queue | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/unified-work-queue.md` | SOC L1/L2, Incident Commander | Incident, Task, external projections | Work Queue | REQ-OBJ-001,012; REQ-PROD-013; REQ-UX-008,009 | — | CAP-CMD-102,104,105,108; Saved Views | — | 2026-08-04 |
| CAP-CMD-102 | Work Assignment | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/work-assignment.md` | Incident Commander, Team Lead | Incident, Task, Principal | Work Queue, Incident Detail | REQ-PROD-003,013,021; REQ-SEC-001 | OPEN-013 | CAP-CMD-103; Identity, Notifications | — | 2026-08-04 |
| CAP-CMD-103 | Operational Ownership | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/operational-ownership.md` | Incident Commander, Analyst | Incident, Task, owner/team | Work Queue, Mission Control | REQ-PROD-006,009,013; REQ-OBJ-001 | OPEN-013 | CAP-CMD-102; Ownership Register | — | 2026-08-04 |
| CAP-CMD-104 | Priority and Severity Coordination | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/priority-and-severity-coordination.md` | SOC L2, Incident Commander | Incident, Task, Signal/Alert context | Work Queue, Priorities | REQ-PROD-005,013,021; REQ-UX-005 | OPEN-013 | CAP-CMD-002,105,204,205 | — | 2026-08-04 |
| CAP-CMD-105 | SLA Tracking | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/sla-tracking.md` | Analyst, Incident Commander, Delivery Manager | Incident, Task, SLA policy | Work Queue SLA Risk, Handover | REQ-PROD-005,013,053; REQ-OBJ-012 | OPEN-006,013 | CAP-CMD-101,102,110; Notifications | — | 2026-08-04 |
| CAP-CMD-106 | Incident Coordination | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/incident-coordination.md` | Incident Commander, SOC Analyst | Incident, Task, Case/Govern links | Incident Detail, Mission Control | REQ-OBJ-001; REQ-PROD-003,008,013; REQ-SEC-001 | OPEN-013 | CAP-CMD-102,104,107,110 | — | 2026-08-04 |
| CAP-CMD-107 | Task Coordination | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/task-coordination.md` | Analyst, Incident Commander | Task, source relation | Work Queue, Readiness | REQ-PROD-009,013,021; REQ-OBJ-012 | OPEN-013 | CAP-CMD-005,102,103; Object Linking | — | 2026-08-04 |
| CAP-CMD-108 | Bulk Coordination | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/bulk-coordination.md` | SOC L2, Team Lead | Incident/Task selection, job | Work Queue, audit | REQ-PROD-004,009,013; REQ-SEC-001 | OPEN-013 | CAP-CMD-101,102,104; Jobs | — | 2026-08-04 |
| CAP-CMD-109 | Work Freshness and Staleness | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/work-freshness-and-staleness.md` | Analyst, Incident Commander, Auditor | Incident/Task and projections | Work Queue, Handover | REQ-PROD-005,008,013; REQ-UX-005 | — | Data Quality, Jobs, Notifications | — | 2026-08-04 |
| CAP-CMD-110 | Escalation | Command | incidents-and-work-queue | draft | defined | planned | `06-command/modules/incidents-and-work-queue/capabilities/escalation.md` | Incident Commander, SOC L2 | Incident/Task, Case/Finding/Action Request links | Incident Detail, Work Queue | REQ-PROD-003,004,008,013; REQ-SEC-001 | OPEN-013 | CAP-CMD-005,106; Notifications | — | 2026-08-04 |
| CAP-CMD-201 | Service Context | Command | risk-and-coverage | draft | defined | planned | `06-command/modules/risk-and-coverage/capabilities/service-context.md` | Incident Commander, Business Owner | Service projection, Incident, Exposure/Coverage | Risk and Coverage | REQ-PROD-005,013,021,032 | OPEN-013 | Service Catalog, CAP-CMD-204 | — | 2026-08-04 |
| CAP-CMD-202 | Exposure Overview | Command | risk-and-coverage | draft | defined | planned | `06-command/modules/risk-and-coverage/capabilities/exposure-overview.md` | SOC L2, Business Owner | Exposure, Service, Incident | Risk and Coverage | REQ-PROD-006,013,032,037 | OPEN-013 | Exposure sources, Data Quality | — | 2026-08-04 |
| CAP-CMD-203 | Coverage Overview | Command | risk-and-coverage | draft | defined | planned | `06-command/modules/risk-and-coverage/capabilities/coverage-overview.md` | Incident Commander, Detection Engineer viewer | Coverage projections, Service, Task | Risk, Readiness | REQ-PROD-005,013,032,037 | OPEN-013 | CAP-CMD-201,301,303; Metrics | — | 2026-08-04 |
| CAP-CMD-204 | Business Impact Context | Command | risk-and-coverage | draft | defined | planned | `06-command/modules/risk-and-coverage/capabilities/business-impact-context.md` | Incident Commander, Business Owner | Incident, Service, Finding/Result refs | Incident Detail, Priorities | REQ-PROD-003,005,013,021 | OPEN-013 | CAP-CMD-201,205,106 | — | 2026-08-04 |
| CAP-CMD-205 | Risk Prioritization Context | Command | risk-and-coverage | draft | defined | planned | `06-command/modules/risk-and-coverage/capabilities/risk-prioritization-context.md` | Incident Commander, SOC L2 | impact, urgency, service, exposure, SLA | Mission Control, Work Queue | REQ-PROD-003,010,013,021 | OPEN-013 | CAP-CMD-002,104,105,201 | — | 2026-08-04 |
| CAP-CMD-301 | Readiness Overview | Command | readiness-and-operations | draft | defined | planned | `06-command/modules/readiness-and-operations/capabilities/readiness-overview.md` | Readiness Coordinator, Incident Commander | capability projection, Tasks, plans/exercises | Readiness, Mission Control | REQ-PROD-005,013,021,057 | OPEN-010,013 | CAP-CMD-203,303,304,305 | — | 2026-08-04 |
| CAP-CMD-302 | Exercise Coordination | Command | readiness-and-operations | draft | defined | planned | `06-command/modules/readiness-and-operations/capabilities/exercise-coordination.md` | Readiness Coordinator, Team Lead | exercise record, plan, Task | Readiness, reporting | REQ-PROD-005,013,021,057 | OPEN-010,013 | CAP-CMD-301,303,304; Collaboration | — | 2026-08-04 |
| CAP-CMD-303 | Improvement Actions | Command | readiness-and-operations | draft | defined | planned | `06-command/modules/readiness-and-operations/capabilities/improvement-actions.md` | Readiness Coordinator, Task owner | Result/exercise source, Task | Readiness, Work Queue | REQ-PROD-009,013,021; REQ-OBJ-012 | OPEN-013 | CAP-CMD-107,301,302,203 | — | 2026-08-04 |
| CAP-CMD-304 | Operational Plans | Command | readiness-and-operations | draft | defined | planned | `06-command/modules/readiness-and-operations/capabilities/operational-plans.md` | Readiness Coordinator, Incident Commander | plan record, roles, service refs | Readiness, Exercises | REQ-PROD-005,013,021,057 | OPEN-010,013 | CAP-CMD-301,302,305 | — | 2026-08-04 |
| CAP-CMD-305 | Capability Readiness | Command | readiness-and-operations | draft | defined | planned | `06-command/modules/readiness-and-operations/capabilities/capability-readiness.md` | Readiness Coordinator, Product Owner | capability, health/assurance, assessment | Readiness, Mission Control | REQ-PROD-012,013,019,057 | OPEN-010,013 | Capability Register, Platform Health | — | 2026-08-04 |
| CAP-CMD-401 | Customers and Delivery Context | Command | customers-and-delivery | draft | proposed | planned | `06-command/modules/customers-and-delivery/capabilities/customers-and-delivery-context.md` | Delivery Manager, Customer Success | Incident/Task, Result, Report/customer context | optional delivery journeys | REQ-PROD-013,019,033,053 | OPEN-006,013 | Reporting, Metrics, Export | — | 2026-08-04 |

## Investigate — Phase 4B.1

| capability_id | name | owner_product | owner_module | status | delivery_status | delivery_mode | canonical_file | primary_roles | primary_objects | consumers | requirement_ids | open_decisions | dependencies | supersedes | last_reviewed |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| CAP-INV-001 | Signal Triage | Investigate | signals-and-hunt | draft | defined | planned | `07-investigate/modules/signals-and-hunt/capabilities/signal-triage.md` | SOC Analyst, Threat Hunter | Signal/Alert/Detection/Event projections, Case link | Triage Desk, Case intake | REQ-PROD-002,005,008,014,045 | OPEN-013 | CAP-INV-002,004,102; Command; Trace | — | 2026-08-04 |
| CAP-INV-002 | Event Search | Investigate | signals-and-hunt | draft | defined | planned | `07-investigate/modules/signals-and-hunt/capabilities/event-search.md` | Analyst, Threat Hunter | Query, Search Job, Event, Case refs | Event Search, Hunt, Case | REQ-PROD-014,019 | — | CAP-INV-003,004,006,007,008; Shared Query/Jobs | — | 2026-08-04 |
| CAP-INV-003 | Query Authoring and Assistance | Investigate | signals-and-hunt | draft | defined | planned | `07-investigate/modules/signals-and-hunt/capabilities/query-authoring-and-assistance.md` | Analyst, Query Author | Query draft, schema projection | Event Search, Query Assets | REQ-PROD-014; REQ-AI-002 | OPEN-015 | CAP-INV-002,006; Studio optional | — | 2026-08-04 |
| CAP-INV-004 | Event Inspection and Pivot | Investigate | signals-and-hunt | draft | defined | planned | `07-investigate/modules/signals-and-hunt/capabilities/event-inspection-and-pivot.md` | Analyst, Threat Hunter | Event, Entity, Case, candidate refs | Event Search, Case | REQ-PROD-014,019 | — | CAP-INV-002,007,102,105,107; Inspector | — | 2026-08-04 |
| CAP-INV-005 | Hunt Management | Investigate | signals-and-hunt | draft | defined | planned | `07-investigate/modules/signals-and-hunt/capabilities/hunt-management.md` | Threat Hunter, Hunt Lead | Hunt workspace, Query/Search Job, Hypothesis, Case | Hunt, Case, Replay | REQ-PROD-014,020 | — | CAP-INV-002,003,006,008,102,103 | — | 2026-08-04 |
| CAP-INV-006 | Saved Searches and Query Assets | Investigate | signals-and-hunt | draft | defined | planned | `07-investigate/modules/signals-and-hunt/capabilities/saved-searches-and-query-assets.md` | Query Author, Threat Hunter | Query Asset, Query, versions | Search/Hunt | REQ-PROD-014,019 | — | CAP-INV-002,003,005; Versioning | — | 2026-08-04 |
| CAP-INV-007 | Search Result Organization | Investigate | signals-and-hunt | draft | defined | planned | `07-investigate/modules/signals-and-hunt/capabilities/search-result-organization.md` | Analyst, Threat Hunter | result selection, annotations, Case links | Search, Case, Export | REQ-PROD-014,019 | OPEN-013 | CAP-INV-002,004,102; Export | — | 2026-08-04 |
| CAP-INV-008 | Search and Hunt Provenance | Investigate | signals-and-hunt | draft | defined | planned | `07-investigate/modules/signals-and-hunt/capabilities/search-and-hunt-provenance.md` | Analyst, Reviewer, Auditor | Query/Search Job versions, Trace refs | Search, Hunt, Replay | REQ-PROD-014; REQ-AI-002 | OPEN-015 | CAP-INV-002,005,006; Trace/Activity | — | 2026-08-04 |
| CAP-INV-101 | Case Queue | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/case-queue.md` | Case Analyst, Investigation Lead | Case, Incident projection, Finding summary | Case Queue/Workspace | REQ-PROD-014 | — | CAP-INV-102; Search, Saved Views | — | 2026-08-04 |
| CAP-INV-102 | Case Lifecycle and Coordination | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/case-lifecycle-and-coordination.md` | Investigation Lead, Case Analyst | Case, Incident links, contributors | Case Workspace, Command projection | REQ-PROD-014 | OPEN-013 | CAP-INV-101,103,104,105,107,109,110,111 | — | 2026-08-04 |
| CAP-INV-103 | Hypothesis Management | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/hypothesis-management.md` | Analyst, Reviewer | Hypothesis, Evidence/results refs | Case, Findings, Replay | REQ-PROD-014; REQ-AI-002 | OPEN-013,015 | CAP-INV-102,107,108,109; Versioning | — | 2026-08-04 |
| CAP-INV-104 | Entity and Relationship Management | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/entity-and-relationship-management.md` | Analyst, Threat Hunter | Shared Entity, relations, Case/Incident refs | Entity Graph, Search, Case | REQ-PROD-014 | OPEN-013 | CAP-INV-002,004,102; Entity Resolution | — | 2026-08-04 |
| CAP-INV-105 | Artifact Management | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/artifact-management.md` | Analyst, Specialist | Artifact, versions, derivatives, Case | Artifact Detail, future workbench | REQ-PROD-014,061 | OPEN-014 | CAP-INV-102,106,107; Versioning; 4B.2 boundary | — | 2026-08-04 |
| CAP-INV-106 | Attachment Handling | Investigate | cases-and-evidence | draft | proposed | planned | `07-investigate/modules/cases-and-evidence/capabilities/attachment-handling.md` | Contributor, Report Author | Attachment, Note/Comment/Report refs | Case collaboration, reporting | REQ-PROD-061 | OPEN-014 | CAP-INV-105,107,111,114; Shared Attachments | — | 2026-08-04 |
| CAP-INV-107 | Evidence Creation and Management | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/evidence-creation-and-management.md` | Analyst, Evidence Reviewer | Evidence, Artifact/source, Case, Hypothesis/Finding refs | Evidence Board, Finding | REQ-PROD-014,061,062 | OPEN-013 | CAP-INV-102,103,105,108,109; Provenance | — | 2026-08-04 |
| CAP-INV-108 | Evidence Review and Qualification | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/evidence-review-and-qualification.md` | Evidence Reviewer, Analyst | Evidence dimensions, Hypothesis refs | Evidence Board, Finding Review | REQ-PROD-014,062 | OPEN-013 | CAP-INV-103,107,109; future collection request | — | 2026-08-04 |
| CAP-INV-109 | Finding Management | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/finding-management.md` | Analyst, Reviewer | Finding, Evidence, Incident/Request refs | Findings, Govern handoff, Report | REQ-PROD-014,016 | OPEN-013,015 | CAP-INV-103,107,108,113,114 | — | 2026-08-04 |
| CAP-INV-110 | Investigation Timeline | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/investigation-timeline.md` | Case Analyst, Reviewer | Timeline Entry, Case objects, Govern projections | Case Timeline, Replay, Report | REQ-PROD-014,018 | — | CAP-INV-102,107,109,112,114; Timeline Engine | — | 2026-08-04 |
| CAP-INV-111 | Case Collaboration and Investigation Notes | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/case-collaboration-and-investigation-notes.md` | Case Analyst, Contributors | Note, Comment, Attachment, Case/Task refs | Case Workspace | REQ-PROD-014,018 | OPEN-013,014 | CAP-INV-102,105,106,107; Shared Collaboration | — | 2026-08-04 |
| CAP-INV-112 | Case Replay and Investigation Review | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/case-replay-and-investigation-review.md` | Investigation Lead, Reviewer | Case snapshot, Query/Hypothesis/Evidence/Finding versions | Review, future Detection/Readiness | REQ-PROD-014,020 | OPEN-013,015 | CAP-INV-008,102,103,107,109,110 | — | 2026-08-04 |
| CAP-INV-113 | Action Request Preparation | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/action-request-preparation.md` | Investigation Lead, Senior Analyst | Action Request draft, Findings, Evidence, Case | Govern Review Queue, Case | REQ-PROD-014,016 | OPEN-007,013,015 | CAP-INV-107,108,109; Govern; Command impact | — | 2026-08-04 |
| CAP-INV-114 | Case Reporting Preparation | Investigate | cases-and-evidence | draft | defined | planned | `07-investigate/modules/cases-and-evidence/capabilities/case-reporting-preparation.md` | Case Analyst, Report Author, Reviewer | Report draft, citations, Findings/Evidence/Timeline refs | Investigation Report, Reporting Engine | REQ-PROD-014,018 | OPEN-014 | CAP-INV-109,110,111; Reporting/Export | — | 2026-08-04 |

## Contrôles de registre

Une entrée échoue si son ID est dupliqué ou recyclé, si le fichier, owner, utilisateur, objet, Requirement ID, entrée, sortie, action, alternative sans IA ou preuve de delivery manque, ou si un owner concurrent revendique la même capability.

## Couverture enregistrée

| Produit / phase | IDs | Documents | Sans owner/utilisateur/objet | Delivery status | Delivery mode actuel | Owners concurrents |
|---|---:|---:|---:|---|---|---:|
| Command / Phase 4A | 27 | 27/27 | 0 | 26 defined, 1 proposed | 27 planned | 0 |
| Investigate / Phase 4B.1 | 22 | 22/22 | 0 | 21 defined, 1 proposed | 22 planned | 0 |

Aucune capability CAP-INV-2xx, 3xx, 4xx ou 5xx n’est attribuée. Les futures sous-phases reçoivent leurs IDs uniquement lors de leur traitement propriétaire.
