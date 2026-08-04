---
id: 06-command-navigation
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-013
  - REQ-UX-001
  - REQ-UX-007
  - REQ-UX-008
---
# Navigation — Command

## Navigation primaire

Command est une destination produit du Global Header. Investigate, Govern et Studio sont des transitions explicites, jamais des sous-menus Command.

## Navigation locale

| Destination | Type | Condition |
|---|---|---|
| Mission Control | module/workspace | toujours si `perm.command.read` |
| Work Queue | workspace | toujours si au moins Incident ou Task autorisé |
| Risk and Coverage | module | projections disponibles ou Empty explicable |
| Readiness and Operations | module | lecture readiness autorisée |
| Customers and Delivery | module proposé | deployment profile + permission |

## Liens d’objet

Incident Detail ouvre depuis Work Queue, Mission Control, Search, Notification ou lien profond. Case, Decision, Run, Result et Workflow ouvrent dans leur produit propriétaire.

## Retour

Le bouton Back restaure l’origine réelle. Une erreur de permission ou de transition ne remplace pas le workspace source.

## Command Palette et Search

La Command Palette lance navigation et actions autorisées mais pas une action dangereuse directe. Global Search renvoie uniquement les objets autorisés et leur owner produit.

## Critère

**Given** un Incident ouvert depuis `view=sla-risk`, **When** l’utilisateur ouvre un Case puis revient deux fois, **Then** Incident Detail puis la Work Queue `SLA Risk` sont restaurés avec le même contexte.
