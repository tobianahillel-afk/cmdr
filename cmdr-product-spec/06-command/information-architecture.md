---
id: 06-command-information-architecture
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-013
  - REQ-UX-001
  - REQ-UX-008
  - REQ-UX-009
open_decisions:
  - OPEN-006
  - OPEN-010
---
# Information architecture — Command

## Destinations locales

1. **Mission Control**
2. **Work Queue**
3. **Risk and Coverage**
4. **Readiness and Operations**
5. **Customers and Delivery** — seulement si activé pour le déploiement.

Incident Detail est une destination d’objet, pas un module de navigation parallèle.

## Work Queue

La route reste stable. Les six vues sont des configurations Saved Views ; les anciens Screen IDs `CMD-IWQ-001` à `005` restent aliases dépréciés. La vue change par paramètre, pas par page.

## Hiérarchie

Global Header → Context Bar → navigation locale Command → workspace actif → Inspector unique. Une transition interproduit conserve return origin, tenant, environnement, objet, filtres et sélection lorsque autorisé.

## Modules historiques

| Ancien module | Destination canonique | Écrans concernés |
|---|---|---|
| exposure-and-coverage | Risk and Coverage | CMD-EXC-001 |
| risk-and-business-impact | Risk and Coverage | CMD-RBI-001 |
| customer-and-reports | Customers and Delivery proposé | CMD-CRP-001 |

Les fichiers d’écran ne sont pas réécrits en Phase 4A.

## Progressive disclosure

1. owner, priority, impact, SLA, freshness et next action ;
2. relations, timeline, Tasks, services, blockers et projections ;
3. trace, audit, sources, versions et détails techniques dans le produit propriétaire.

## Critère

**Given** un utilisateur qui change de module ou ouvre un objet lié, **When** il revient, **Then** le module, le workspace, la vue, les filtres, la sélection et le scroll sont restaurés sans créer de destination concurrente.
