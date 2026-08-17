---
id: ADR-0005-page-view-mode-filter-rules
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-001
  - REQ-UX-008
  - REQ-UX-009
---

# ADR-0005 — Page, Workspace, View, Mode et Filter

## Statut

Draft ; applique les décisions source de Phase 3 sans devenir une approbation `validated`.

## Contexte et décision

Une Page correspond à un objectif autonome ; un Workspace à une activité durable ; une View à un sous-ensemble/configuration ; un Mode à une représentation ; un Filter à une restriction temporaire.

La Work Queue est un workspace unique. `All`, `Incidents`, `Tasks`, `Unassigned`, `SLA Risk` et `My Work` sont des vues système. Les cinq anciens fichiers de variantes sont dépréciés comme aliases ou entrée de migration.

## Conséquences

Moins de routes et de cycles concurrents ; retour exact ; Saved Views génériques réutilisables ; migration des anciens URLs. Les vues restent adressables et permission-aware.

## Migration

- `incidents`, `tasks`, `unassigned`, `sla-risk` → `view` keys ;
- `team-load` → Work Queue avec notice de migration, possibilité de vue personnalisée ;
- screen register et liens seront normalisés sans réécrire les contenus métier ;
- aucune nouvelle page de vue n'est autorisée.

## Réévaluation

Après Phase 6, vérifier les aliases et retirer ceux sans usage. L'ADR reste Draft jusqu'au processus d'approbation.
