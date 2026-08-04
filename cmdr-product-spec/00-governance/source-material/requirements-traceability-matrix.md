---
id: requirements-traceability-matrix
domain: 00-governance
status: draft
owner: QA and Traceability Lead
updated: 2026-08-04
source-of-truth: source-material
requirements:
  - REQ-PROD-001
  - REQ-PROD-062
  - REQ-UX-001
  - REQ-UX-010
  - REQ-OBJ-001
  - REQ-OBJ-012
---

# Requirements Traceability Matrix

## Interprétation

Les 122 Requirement IDs sources sont conservés. `conform` signifie qu’une phase propriétaire actuelle fournit des règles substantielles et sans contradiction active ; cela ne prouve ni implémentation ni livraison. La Phase 4A ajoute les preuves fonctionnelles Command sans promouvoir artificiellement les exigences globales dépendantes des Phases 4B–4E, Objets, Permissions, Parcours, Écrans ou Technique.

## Couverture

| État | Après Phase 3 | Après Phase 4A |
|---|---:|---:|
| conform | 99 | 99 |
| partial | 20 | 20 |
| absent | 3 | 3 |
| contradictory | 0 | 0 |
| total | 122 | 122 |

## Traçabilité Command Phase 4A

| Capability | Module | Nom | Functional status | Delivery mode | Requirement IDs | Objets | Classes | OPEN | Preuve canonique | Consommateurs |
|---|---|---|---|---|---|---|---|---|---|---|
| CAP-CMD-001 | mission-control | Situation Overview | defined | planned | REQ-PROD-013, REQ-PROD-008, REQ-PROD-010, REQ-PROD-021 | Incident, Task, Decision / Response Run / Result, Service | 0, 2 | — | `06-command/modules/mission-control/capabilities/situation-overview.md` | Mission Control — Now, Mission Control — Situation |
| CAP-CMD-002 | mission-control | Priority Management | defined | planned | REQ-PROD-003, REQ-PROD-010, REQ-PROD-013, REQ-PROD-021 | Incident / Task, Service, SLA context, Audit trail | 0, 2 | OPEN-013 | `06-command/modules/mission-control/capabilities/priority-management.md` | Mission Control — Priorities, Unified Work Queue |
| CAP-CMD-003 | mission-control | Situation Timeline | defined | planned | REQ-PROD-005, REQ-PROD-008, REQ-PROD-013, REQ-UX-007 | Timeline Entry, Incident, Decision / Run / Result | 0 | — | `06-command/modules/mission-control/capabilities/situation-timeline.md` | Mission Control — Situation, Incident Detail |
| CAP-CMD-004 | mission-control | Handover | defined | planned | REQ-PROD-008, REQ-PROD-013, REQ-PROD-021, REQ-PROD-010 | Incident / Task, Decision / Response Run / Result, Case / Finding, Handover record, Incident / Task ownership | 2 | — | `06-command/modules/mission-control/capabilities/handover.md` | Mission Control — Handover, parcours de relève |
| CAP-CMD-005 | mission-control | Operational Blockers | defined | planned | REQ-PROD-006, REQ-PROD-009, REQ-PROD-013, REQ-PROD-021 | Incident / Task, Decision / Case / external dependency, Task | 2 | OPEN-013 | `06-command/modules/mission-control/capabilities/operational-blockers.md` | Mission Control, Unified Work Queue |
| CAP-CMD-006 | mission-control | Recent Results and Outcomes | defined | planned | REQ-PROD-005, REQ-PROD-008, REQ-PROD-013, REQ-PROD-021 | Result / Response Run / Decision, Incident, Timeline | 0, 2 | OPEN-013 | `06-command/modules/mission-control/capabilities/recent-results-and-outcomes.md` | Mission Control — Now, Mission Control — Situation |
| CAP-CMD-101 | incidents-and-work-queue | Unified Work Queue | defined | planned | REQ-OBJ-001, REQ-OBJ-012, REQ-PROD-013, REQ-UX-008, REQ-UX-009 | Incident, Task, Case / Decision / Response Run, Saved View application, Incident / Task | 0 | — | `06-command/modules/incidents-and-work-queue/capabilities/unified-work-queue.md` | Work Queue, Incident Detail |
| CAP-CMD-102 | incidents-and-work-queue | Work Assignment | defined | planned | REQ-PROD-003, REQ-PROD-013, REQ-PROD-021, REQ-SEC-001 | Incident / Task, Principal / Role, Notification | 2 | OPEN-013 | `06-command/modules/incidents-and-work-queue/capabilities/work-assignment.md` | Unified Work Queue, Incident Detail |
| CAP-CMD-103 | incidents-and-work-queue | Operational Ownership | defined | planned | REQ-PROD-006, REQ-PROD-009, REQ-PROD-013, REQ-OBJ-001 | Incident / Task, Case / Decision / Run, Principal / Team, Audit | 0, 2 | OPEN-013 | `06-command/modules/incidents-and-work-queue/capabilities/operational-ownership.md` | Work Queue, Incident Detail |
| CAP-CMD-104 | incidents-and-work-queue | Priority and Severity Coordination | defined | planned | REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-UX-005 | Incident / Task, Signal / Alert, Service / SLA, Recommendation disposition | 0, 2 | OPEN-013 | `06-command/modules/incidents-and-work-queue/capabilities/priority-and-severity-coordination.md` | Unified Work Queue, Mission Control — Priorities |
| CAP-CMD-105 | incidents-and-work-queue | SLA Tracking | defined | planned | REQ-PROD-005, REQ-PROD-013, REQ-PROD-053, REQ-OBJ-012 | Incident / Task, SLA policy / engagement, Task | 0, 2 | OPEN-006, OPEN-013 | `06-command/modules/incidents-and-work-queue/capabilities/sla-tracking.md` | Unified Work Queue — SLA Risk, Mission Control |
| CAP-CMD-106 | incidents-and-work-queue | Incident Coordination | defined | planned | REQ-OBJ-001, REQ-PROD-003, REQ-PROD-008, REQ-PROD-013, REQ-SEC-001 | Signal / Alert, Case / Finding, Decision / Response Run / Result, Incident / Task, Incident, Task, Action Request | 2 | OPEN-013 | `06-command/modules/incidents-and-work-queue/capabilities/incident-coordination.md` | Incident Detail, Unified Work Queue |
| CAP-CMD-107 | incidents-and-work-queue | Task Coordination | defined | planned | REQ-PROD-009, REQ-PROD-013, REQ-PROD-021, REQ-OBJ-012 | Task, Source object, Source relation | 2 | OPEN-013 | `06-command/modules/incidents-and-work-queue/capabilities/task-coordination.md` | Unified Work Queue, Incident Detail |
| CAP-CMD-108 | incidents-and-work-queue | Bulk Coordination | defined | planned | REQ-PROD-004, REQ-PROD-009, REQ-PROD-013, REQ-SEC-001 | Incident / Task, Saved View / filters, Export job | 0, 2 | OPEN-013 | `06-command/modules/incidents-and-work-queue/capabilities/bulk-coordination.md` | Unified Work Queue, audit |
| CAP-CMD-109 | incidents-and-work-queue | Work Freshness and Staleness | defined | planned | REQ-PROD-005, REQ-PROD-008, REQ-PROD-013, REQ-UX-005 | Incident / Task, Case / Decision / Run projections, Health, Task, Acknowledgement | 0, 2 | — | `06-command/modules/incidents-and-work-queue/capabilities/work-freshness-and-staleness.md` | Unified Work Queue, Mission Control |
| CAP-CMD-110 | incidents-and-work-queue | Escalation | defined | planned | REQ-PROD-003, REQ-PROD-004, REQ-PROD-008, REQ-PROD-013, REQ-SEC-001 | Incident / Task, Case / Finding / Decision, Escalation record/state, Case or Action Request | 2 | OPEN-013 | `06-command/modules/incidents-and-work-queue/capabilities/escalation.md` | Incident Detail, Unified Work Queue |
| CAP-CMD-201 | risk-and-coverage | Service Context | defined | planned | REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-032 | Service, Incident / Task, Exposure / coverage projections, Data quality follow-up | 0, 2 | OPEN-013 | `06-command/modules/risk-and-coverage/capabilities/service-context.md` | Risk and Coverage, Incident Detail |
| CAP-CMD-202 | risk-and-coverage | Exposure Overview | defined | planned | REQ-PROD-006, REQ-PROD-013, REQ-PROD-032, REQ-PROD-037 | Exposure projection, Service, Incident, Incident / Task, Exposure source | 0, 2 | OPEN-013 | `06-command/modules/risk-and-coverage/capabilities/exposure-overview.md` | Risk and Coverage, Mission Control |
| CAP-CMD-203 | risk-and-coverage | Coverage Overview | defined | planned | REQ-PROD-005, REQ-PROD-013, REQ-PROD-032, REQ-PROD-037 | Coverage projections, Service, Task, Coverage source | 0, 2 | OPEN-013 | `06-command/modules/risk-and-coverage/capabilities/coverage-overview.md` | Risk and Coverage, Readiness Overview |
| CAP-CMD-204 | risk-and-coverage | Business Impact Context | defined | planned | REQ-PROD-003, REQ-PROD-005, REQ-PROD-013, REQ-PROD-021 | Incident, Service, Finding/Result, Action Request context | 2 | OPEN-013 | `06-command/modules/risk-and-coverage/capabilities/business-impact-context.md` | Incident Detail, Mission Control |
| CAP-CMD-205 | risk-and-coverage | Risk Prioritization Context | defined | planned | REQ-PROD-003, REQ-PROD-010, REQ-PROD-013, REQ-PROD-021 | Incident / Task, Service / Exposure / Coverage, Confidence / severity, Priority recommendation, Incident/Task | 0, 2 | OPEN-013 | `06-command/modules/risk-and-coverage/capabilities/risk-prioritization-context.md` | Mission Control — Priorities, Unified Work Queue |
| CAP-CMD-301 | readiness-and-operations | Readiness Overview | defined | planned | REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-057 | Capability projection, Task, Plan/Exercise records, Readiness assessment | 0, 2 | OPEN-010, OPEN-013 | `06-command/modules/readiness-and-operations/capabilities/readiness-overview.md` | Readiness & Operations screen, Mission Control |
| CAP-CMD-302 | readiness-and-operations | Exercise Coordination | defined | planned | REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-057 | Operational Plan, Capability Readiness, Workflow/Simulation, Exercise record, Task | 2 | OPEN-010, OPEN-013 | `06-command/modules/readiness-and-operations/capabilities/exercise-coordination.md` | Readiness & Operations, reporting |
| CAP-CMD-303 | readiness-and-operations | Improvement Actions | defined | planned | REQ-PROD-009, REQ-PROD-013, REQ-PROD-021, REQ-OBJ-012 | Result, Exercise/readiness/coverage, Task, Source relation | 2 | OPEN-013 | `06-command/modules/readiness-and-operations/capabilities/improvement-actions.md` | Readiness Overview, Unified Work Queue |
| CAP-CMD-304 | readiness-and-operations | Operational Plans | defined | planned | REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-057 | Plan/Playbook references, Capability Readiness, Service/Incident, Operational Plan record, Task | 2 | OPEN-010, OPEN-013 | `06-command/modules/readiness-and-operations/capabilities/operational-plans.md` | Readiness Overview, Exercise Coordination |
| CAP-CMD-305 | readiness-and-operations | Capability Readiness | defined | planned | REQ-PROD-012, REQ-PROD-013, REQ-PROD-019, REQ-PROD-057 | Capability, Health/Assurance/Exercise, Task, Readiness assessment | 0, 2 | OPEN-010, OPEN-013 | `06-command/modules/readiness-and-operations/capabilities/capability-readiness.md` | Readiness Overview, Operational Plans |
| CAP-CMD-401 | customers-and-delivery | Customers and Delivery Context | proposed | planned | REQ-PROD-053, REQ-PROD-013, REQ-PROD-019, REQ-PROD-033 | Incident / Task, Result, Report, Customer/engagement context, Report request/draft context, Task, Customer/contract source | 0, 2 | OPEN-006, OPEN-013 | `06-command/modules/customers-and-delivery/capabilities/customers-and-delivery-context.md` | Customer & Reports legacy screen, Reporting Engine |

## Résultats de la sous-phase

- `REQ-PROD-013`, `REQ-OBJ-001`, `REQ-OBJ-012`, `REQ-UX-008` et `REQ-UX-009` gagnent des preuves fonctionnelles détaillées sans changer leur état déjà conforme.
- `REQ-PROD-008` et `REQ-PROD-010` restent partiels globalement : Command est détaillé, mais les produits et parcours des Phases 4B–4E et 5 ne le sont pas encore.
- `REQ-PROD-006` reste partiel globalement : le périmètre Command traité ne contient plus de document générique actif hors écrans reportés, mais le référentiel complet conserve des domaines Template-level.
- `REQ-SEC-003` à `REQ-SEC-005` restent partiels : les besoins fonctionnels sont identifiés, la matrice atomique et les contrats de confiance restent ultérieurs.
- `REQ-UX-010` reste partiel : la matrice capability→écran est prête, mais aucun écran Command n’est réécrit en détail.
- `REQ-JRN-002`, `REQ-JRN-005` et `REQ-JRN-008` restent absents et appartiennent à la phase Parcours.
- aucun Requirement ID nouveau ;
- aucune contradiction active ;
- aucune exigence planned présentée comme livrée.

## Inventaire complet

### conform — 99

`REQ-AI-001`, `REQ-AI-002`, `REQ-AI-003`, `REQ-AI-004`, `REQ-AI-005`, `REQ-AI-006`, `REQ-AI-011`, `REQ-BRAND-001`, `REQ-BRAND-002`, `REQ-BRAND-003`, `REQ-BRAND-004`, `REQ-BRAND-005`, `REQ-BRAND-006`, `REQ-BRAND-008`, `REQ-INV-001`, `REQ-INV-002`, `REQ-INV-003`, `REQ-INV-004`, `REQ-INV-005`, `REQ-INV-006`, `REQ-OBJ-001`, `REQ-OBJ-002`, `REQ-OBJ-003`, `REQ-OBJ-004`, `REQ-OBJ-005`, `REQ-OBJ-006`, `REQ-OBJ-007`, `REQ-OBJ-008`, `REQ-OBJ-010`, `REQ-OBJ-011`, `REQ-OBJ-012`, `REQ-PROD-001`, `REQ-PROD-002`, `REQ-PROD-003`, `REQ-PROD-004`, `REQ-PROD-005`, `REQ-PROD-009`, `REQ-PROD-011`, `REQ-PROD-012`, `REQ-PROD-013`, `REQ-PROD-014`, `REQ-PROD-015`, `REQ-PROD-016`, `REQ-PROD-017`, `REQ-PROD-018`, `REQ-PROD-019`, `REQ-PROD-021`, `REQ-PROD-022`, `REQ-PROD-023`, `REQ-PROD-024`, `REQ-PROD-025`, `REQ-PROD-026`, `REQ-PROD-027`, `REQ-PROD-028`, `REQ-PROD-029`, `REQ-PROD-030`, `REQ-PROD-031`, `REQ-PROD-032`, `REQ-PROD-033`, `REQ-PROD-034`, `REQ-PROD-035`, `REQ-PROD-036`, `REQ-PROD-037`, `REQ-PROD-038`, `REQ-PROD-039`, `REQ-PROD-040`, `REQ-PROD-041`, `REQ-PROD-042`, `REQ-PROD-043`, `REQ-PROD-044`, `REQ-PROD-045`, `REQ-PROD-046`, `REQ-PROD-047`, `REQ-PROD-048`, `REQ-PROD-049`, `REQ-PROD-050`, `REQ-PROD-051`, `REQ-PROD-052`, `REQ-PROD-053`, `REQ-PROD-054`, `REQ-PROD-055`, `REQ-PROD-056`, `REQ-PROD-057`, `REQ-PROD-058`, `REQ-PROD-059`, `REQ-PROD-060`, `REQ-PROD-061`, `REQ-PROD-062`, `REQ-SEC-001`, `REQ-SEC-002`, `REQ-UX-001`, `REQ-UX-002`, `REQ-UX-003`, `REQ-UX-004`, `REQ-UX-005`, `REQ-UX-006`, `REQ-UX-007`, `REQ-UX-008`, `REQ-UX-009`

### partial — 20

`REQ-AI-007`, `REQ-AI-008`, `REQ-AI-009`, `REQ-AI-010`, `REQ-BRAND-007`, `REQ-JRN-001`, `REQ-JRN-003`, `REQ-JRN-004`, `REQ-JRN-006`, `REQ-JRN-007`, `REQ-OBJ-009`, `REQ-PROD-006`, `REQ-PROD-007`, `REQ-PROD-008`, `REQ-PROD-010`, `REQ-PROD-020`, `REQ-SEC-003`, `REQ-SEC-004`, `REQ-SEC-005`, `REQ-UX-010`

### absent — 3

`REQ-JRN-002`, `REQ-JRN-005`, `REQ-JRN-008`

### contradictory — 0

Aucun Requirement ID n’est classé contradictoire après l’alignement de la Task opérationnelle et la consolidation des modules Command.

## Règle de maintenance

Une capability peut être fonctionnellement `defined` et techniquement `planned`. Seule une preuve approuvée de moteur, contrat et release permet une promotion de delivery mode.
