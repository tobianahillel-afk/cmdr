---
id: deprecated-pattern-list-detail
domain: 03-design-system
status: deprecated
owner: Design System Lead
updated: 2026-08-03
source-of-truth: deprecated
replaced-by: selection-and-inspection.md
requirements:
  - REQ-UX-002
  - REQ-UX-004
---

# Pointeur déprécié

## Remplaçant

[`selection-and-inspection.md`](selection-and-inspection.md)

## Justification

Le pattern Phase 0 était générique ou chevauchait un contrat canonique plus précis.

## Migration

Migrer les consommateurs vers `selection-and-inspection.md` et conserver les règles métier dans leur domaine propriétaire.

## Dépendants

Composants et écrans qui citaient cet ancien pattern.

## Date de retrait

2026-08-03. Le chemin reste disponible pour l'historique et les liens de migration, mais ne porte plus de règle normative.

## Critère d'acceptation

**Given** un consommateur de cet ancien chemin,  
**When** sa dépendance est mise à jour,  
**Then** il référence le remplaçant, ne copie aucune règle locale et ce document n'est jamais utilisé comme source active.
