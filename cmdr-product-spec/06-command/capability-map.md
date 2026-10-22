---
id: command-capability-map
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-09
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
# Capability map — Command

## Capability Specification phase namespace

This map is the active Command capability source for `Capability Specification Phase 4A — Command`.

- namespace: **Capability Specification**;
- canonical execution label: `Capability Specification Phase 4A — Command`;
- status: **PASS**; unchanged by the phase-numbering reconciliation;
- Capability IDs: unchanged;
- historical reports and SHAs: unchanged.

This numbering is not the Delivery Roadmap numbering. `Capability Specification Phase 4A — Command` is not `Delivery Roadmap Phase 2 — Command`, and it does not create a numeric predecessor relationship to `Delivery Roadmap Phase 4 — Govern`. In particular, 4A/4B do not imply a future `4C Govern`.

See `../18-roadmap-and-releases/phase-numbering-and-namespace-convention.md`.

| ID | Capability | Module | Status | Mode | Source |
|---|---|---|---|---|---|
| CAP-CMD-001 | Situation Overview | Mission Control | defined | planned | `modules/mission-control/capabilities/situation-overview.md` |
| CAP-CMD-002 | Priority Management | Mission Control | defined | planned | `modules/mission-control/capabilities/priority-management.md` |
| CAP-CMD-003 | Situation Timeline | Mission Control | defined | planned | `modules/mission-control/capabilities/situation-timeline.md` |
| CAP-CMD-004 | Handover | Mission Control | defined | planned | `modules/mission-control/capabilities/handover.md` |
| CAP-CMD-005 | Operational Blockers | Mission Control | defined | planned | `modules/mission-control/capabilities/operational-blockers.md` |
| CAP-CMD-006 | Recent Results and Outcomes | Mission Control | defined | planned | `modules/mission-control/capabilities/recent-results-and-outcomes.md` |
| CAP-CMD-101 | Unified Work Queue | Incidents and Work Queue | defined | planned | `modules/incidents-and-work-queue/capabilities/unified-work-queue.md` |
| CAP-CMD-102 | Work Assignment | Incidents and Work Queue | defined | planned | `modules/incidents-and-work-queue/capabilities/work-assignment.md` |
| CAP-CMD-103 | Operational Ownership | Incidents and Work Queue | defined | planned | `modules/incidents-and-work-queue/capabilities/operational-ownership.md` |
| CAP-CMD-104 | Priority and Severity Coordination | Incidents and Work Queue | defined | planned | `modules/incidents-and-work-queue/capabilities/priority-and-severity-coordination.md` |
| CAP-CMD-105 | SLA Tracking | Incidents and Work Queue | defined | planned | `modules/incidents-and-work-queue/capabilities/sla-tracking.md` |
| CAP-CMD-106 | Incident Coordination | Incidents and Work Queue | defined | planned | `modules/incidents-and-work-queue/capabilities/incident-coordination.md` |
| CAP-CMD-107 | Task Coordination | Incidents and Work Queue | defined | planned | `modules/incidents-and-work-queue/capabilities/task-coordination.md` |
| CAP-CMD-108 | Bulk Coordination | Incidents and Work Queue | defined | planned | `modules/incidents-and-work-queue/capabilities/bulk-coordination.md` |
| CAP-CMD-109 | Work Freshness and Staleness | Incidents and Work Queue | defined | planned | `modules/incidents-and-work-queue/capabilities/work-freshness-and-staleness.md` |
| CAP-CMD-110 | Escalation | Incidents and Work Queue | defined | planned | `modules/incidents-and-work-queue/capabilities/escalation.md` |
| CAP-CMD-201 | Service Context | Risk and Coverage | defined | planned | `modules/risk-and-coverage/capabilities/service-context.md` |
| CAP-CMD-202 | Exposure Overview | Risk and Coverage | defined | planned | `modules/risk-and-coverage/capabilities/exposure-overview.md` |
| CAP-CMD-203 | Coverage Overview | Risk and Coverage | defined | planned | `modules/risk-and-coverage/capabilities/coverage-overview.md` |
| CAP-CMD-204 | Business Impact Context | Risk and Coverage | defined | planned | `modules/risk-and-coverage/capabilities/business-impact-context.md` |
| CAP-CMD-205 | Risk Prioritization Context | Risk and Coverage | defined | planned | `modules/risk-and-coverage/capabilities/risk-prioritization-context.md` |
| CAP-CMD-301 | Readiness Overview | Readiness and Operations | defined | planned | `modules/readiness-and-operations/capabilities/readiness-overview.md` |
| CAP-CMD-302 | Exercise Coordination | Readiness and Operations | defined | planned | `modules/readiness-and-operations/capabilities/exercise-coordination.md` |
| CAP-CMD-303 | Improvement Actions | Readiness and Operations | defined | planned | `modules/readiness-and-operations/capabilities/improvement-actions.md` |
| CAP-CMD-304 | Operational Plans | Readiness and Operations | defined | planned | `modules/readiness-and-operations/capabilities/operational-plans.md` |
| CAP-CMD-305 | Capability Readiness | Readiness and Operations | defined | planned | `modules/readiness-and-operations/capabilities/capability-readiness.md` |
| CAP-CMD-401 | Customers and Delivery Context | Customers and Delivery | defined | planned | `modules/customers-and-delivery/capabilities/customers-and-delivery-context.md` |

## Comptage

- 27 Capability IDs canoniques Command ;
- 27 `defined` et 0 `proposed` ;
- 27 `delivery_mode: planned` ;
- cible native pour les 26 capabilities cœur ;
- Customers and Delivery deployment-dependent.

## Identifiants

`CAP-CMD-001..099` Mission Control; `100..199` Incidents and Work Queue; `200..299` Risk and Coverage; `300..399` Readiness and Operations; `400..499` Customers and Delivery. Les IDs ne sont jamais recyclés.
