---
id: command-incidents-work-queue
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-OBJ-001
  - REQ-OBJ-012
  - REQ-UX-008
  - REQ-UX-009
---

# Incidents & Work Queue

## Décision canonique

La Work Queue est un workspace Command unique utilisant le Queue Shell. Les variantes sont des Saved Views, jamais des Pages ou écrans autonomes.

## Objets

Incident et Task appartiennent à Command. Alert, Case et autres objets sont des projections autorisées. La vue ne modifie aucune propriété.

## Vues système

`All`, `Incidents`, `Tasks`, `Unassigned`, `SLA Risk`, `My Work`.

## Composition

Context Bar, titre/fraîcheur, Saved Views, recherche, filtres, actions groupées sûres, Data Table et Inspector unique. Vue, filtres, tri, colonnes, densité, sélection et scroll survivent à Incident Detail et au retour.

## Team Load

`Team Load` n'est pas une vue système décidée. Il peut être recréé ultérieurement comme vue partagée personnalisée ou visualisation de capacité dans un module propriétaire, sans Screen ID.

## Migration

Les anciens fichiers `incidents.md`, `tasks.md`, `unassigned.md`, `sla-risk.md` et `team-load.md` sont dépréciés. Les anciennes URLs deviennent des aliases vers le workspace avec paramètre `view`, sauf Team Load qui ouvre la Work Queue avec migration expliquée.

## IA

Priorité ou résumé proposés sont secondaires et attribués. Affectation, filtres, recherche, tri et actions manuelles fonctionnent sans modèle.

## Critère d’acceptation

**Given** la vue `My Work`, des filtres actifs et un Incident sélectionné,  
**When** l'utilisateur ouvre Incident Detail puis revient,  
**Then** le workspace, la vue, les filtres, la sélection, le scroll et l'Inspector sont restaurés.
