---
id: deprecated-work-queue-incidents
domain: 06-command
status: deprecated
owner: Command Product Lead
updated: 2026-08-03
source-of-truth: deprecated
replaced-by: ../README.md
requirements:
  - REQ-OBJ-012
  - REQ-UX-008
---

# Incidents — écran déprécié

## Remplaçant

La Work Queue canonique dans [`../README.md`](../README.md), avec migration `view=incidents`.

## Justification

`Incidents` ne possède pas d'objectif autonome : il s'agit d'une configuration du même travail. Le conserver comme écran créerait une architecture concurrente.

## Migration

Les liens legacy sont résolus vers le workspace et la configuration indiquée. Les données, permissions et identifiants d'objets ne changent pas. Aucun contenu métier détaillé n'est réécrit en Phase 3.

## Dépendants

Screen register, liens profonds, navigation Command et tests de migration.

## Date de retrait

2026-08-03.

## Critère

**Given** un ancien lien vers `Incidents`, **When** il est utilisé, **Then** la Work Queue unique s'ouvre, le contexte autorisé est conservé et aucun second écran n'est créé.
