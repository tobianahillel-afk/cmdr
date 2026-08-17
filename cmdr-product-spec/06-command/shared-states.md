---
id: command-shared-functional-states
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-013
  - REQ-UX-005
  - REQ-UX-010
---
# Shared functional states — Command

## Distinctions obligatoires

| Famille | Exemples | Owner |
|---|---|---|
| état d’objet | Incident `active`, Decision `approved` | produit propriétaire |
| état du travail | assigned, blocked, pending-decision, monitoring | capability Command |
| état de donnée | current, partial, stale, source-unavailable | source/capability |
| readiness | available, not-tested, degraded, planned | Readiness capability |
| delivery status/mode | defined / planned | Capability Register |
| statut documentaire | draft, in-review, validated, deprecated | gouvernance documentation |

## Familles Command

- assignment : unassigned, assigned, accepted, declined, reassignment-pending, unavailable ;
- SLA : healthy, approaching, at-risk, breached, paused, not-applicable, unknown ;
- handover : draft, ready, sent, acknowledged, rejected-for-correction, superseded ;
- readiness : ready/available, partial, degraded, not-tested, unavailable, planned ;
- bulk : prepared, validated, executing, partial-success, completed, cancelled.

## Règle

Une capability peut utiliser un état de travail sans l’ajouter à la machine d’état canonique de l’objet. La Phase Objets décidera les transitions définitives.
