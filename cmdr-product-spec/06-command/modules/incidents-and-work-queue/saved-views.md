---
id: command-work-queue-saved-views
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Work Queue saved views

## Objectif

Définir la source unique des vues enregistrées de la Work Queue.

## Modèle

Une vue enregistrée référence la page `Incidents & Work Queue`, une vue locale, des filtres, un tri, des colonnes, une densité et une portée personnelle ou partagée. Elle n’enregistre jamais de données d’objet ni de permission.

## Vues système

- `Incidents`
- `Tasks`
- `Unassigned`
- `SLA Risk`
- `Team Load`

## Règles

- Les vues système sont versionnées et non supprimables.
- Une vue personnelle est modifiable uniquement par son propriétaire.
- Une vue partagée exige `perm.command.saved-view.share`.
- L’ouverture réévalue les permissions et retire les colonnes non autorisées.
- Les autres documents pointent ici; ils ne redéfinissent pas le schéma des saved views.

## Critères d’acceptation

- URL, filtres et colonnes sont reproductibles.
- Aucune vue ne traverse les tenants.
- Un changement du schéma produit une migration ou un avertissement.
