---
id: deprecated-component-confidence-indicator
domain: 03-design-system
status: deprecated
owner: Design System Lead
updated: 2026-08-03
source-of-truth: deprecated
replaced-by: badges.md
requirements:
  - REQ-UX-002
  - REQ-UX-003
  - REQ-UX-004
---

# Pointeur déprécié

## Remplaçant

[`badges.md`](badges.md)

## Justification

Le contrat Phase 0 était Template-level ou doublonnait une responsabilité désormais possédée par un composant canonique de Phase 3.

## Migration

Migrer variants, états et usages vers `badges.md` ; conserver les données et permissions chez leur propriétaire fonctionnel.

## Dépendants

Écrans, patterns et modules référençant cet ancien nom.

## Date de retrait

2026-08-03. Le chemin reste disponible pour l'historique et les liens de migration, mais ne porte plus de règle normative.

## Critère d'acceptation

**Given** un consommateur de cet ancien chemin,  
**When** sa dépendance est mise à jour,  
**Then** il référence le remplaçant, ne copie aucune règle locale et ce document n'est jamais utilisé comme source active.
