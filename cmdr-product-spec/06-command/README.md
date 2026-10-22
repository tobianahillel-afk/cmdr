---
id: 06-command-readme
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-008
  - REQ-PROD-010
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-OBJ-001
  - REQ-OBJ-012
  - REQ-UX-008
  - REQ-UX-009
open_decisions:
  - OPEN-010
  - OPEN-013
---
# Command

## Mission

Command coordonne la situation opérationnelle, Incidents, Tasks, priorité, ownership, affectation, SLA, impact, handovers, risques et readiness. Il transforme des projections sourcées en travail coordonné ; il ne devient ni espace d’investigation, ni autorité de réponse, ni builder d’automatisation.

## Possède

- Incident et Task opérationnelle ;
- Work Queue unique et ses six vues système ;
- coordination, assignment, priority, operational ownership, SLA et next action ;
- situation, handover et blocages portés par les objets Command ;
- contexte opérationnel de risque, impact et readiness ;
- modules et Capability IDs listés dans `capability-map.md`.

## Consomme sans posséder

- Case, Evidence et Finding — Investigate ;
- Action Request, Decision, Response Run et Result — Govern ;
- Workflow, Automation Agent et Automation Run — CMDR Studio ;
- Endpoint Agent Fleet et configuration — Platform Settings ;
- moteurs et composants partagés — Shared Capabilities et Design System.

## Modules canoniques

| Module | Mission | Source |
|---|---|---|
| Mission Control | conscience partagée, priorité, timeline, handover, blockers et outcomes | `modules/mission-control/README.md` |
| Incidents and Work Queue | file unique, assignment, SLA et coordination Incident/Task | `modules/incidents-and-work-queue/README.md` |
| Risk and Coverage | service, exposure, coverage, impact et contexte de risque | `modules/risk-and-coverage/README.md` |
| Readiness and Operations | readiness, exercices, plans et improvement Tasks | `modules/readiness-and-operations/README.md` |
| Customers and Delivery | contexte deployment-dependent | `modules/customers-and-delivery/README.md` |

Les anciens modules `exposure-and-coverage`, `risk-and-business-impact` et `customer-and-reports` deviennent des points de migration documentaires. Leurs écrans actifs restent inchangés jusqu’à la phase écrans et sont reliés par `screen-capability-map.md`.

## Work Queue

Une seule route/workspace. Vues système exactes : `All`, `Incidents`, `Tasks`, `Unassigned`, `SLA Risk`, `My Work`. `Team Load` n’est ni page ni vue système.

## IA

L’IA est facultative. Elle peut résumer, proposer une priorité, une affectation, une prochaine action, un handover ou un package de transition. Elle ne modifie rien silencieusement, ne confirme pas un Finding, ne crée pas de Decision et ne contourne pas Govern.

## Delivery

Les 27 capabilities Command sont `defined` au niveau product-spec et leur mode courant est `planned` tant qu’aucun logiciel, API ou runtime n’est prouvé. Customers and Delivery est deployment-dependent ; sa définition documentaire n’implique aucune disponibilité runtime.

## Sources de Phase 4A

- `product-definition.md`
- `user-goals.md`
- `capability-map.md`
- `functional-dependency-map.md`
- `object-consumption-map.md`
- `automation-and-ai-model.md`
- `action-classification.md`
- `cross-product-links.md`
- `screen-capability-map.md`
- Capability Register global.

## Critère d’acceptation

**Given** une fonction Command, **When** son owner, ses objets, actions ou transitions sont recherchés, **Then** une capability canonique unique les décrit, les objets externes restent des projections, les actions possèdent une classe et une alternative sans IA existe.
