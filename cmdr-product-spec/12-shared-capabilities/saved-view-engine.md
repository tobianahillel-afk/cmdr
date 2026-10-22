---
id: deprecated-shared-saved-view-engine
domain: 12-shared-capabilities
status: deprecated
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: deprecated
replaced-by: saved-views.md
requirements:
  - REQ-OBJ-011
  - REQ-UX-009
---

# Pointeur déprécié

## Remplaçant

[`saved-views.md`](saved-views.md)

## Justification

Le nom Engine suggérait un contrat technique et concurrençait la capacité fonctionnelle demandée.

## Migration

Les consommateurs utilisent `saved-views.md`; les futurs contrats techniques restent en Phase 8.

## Dépendants

Tous les produits, le registre de propriété et les contrats Saved Views.

## Date de retrait

2026-08-03. Le chemin reste disponible pour l'historique et les liens de migration, mais ne porte plus de règle normative.

## Critère d'acceptation

**Given** un consommateur de cet ancien chemin,  
**When** sa dépendance est mise à jour,  
**Then** il référence le remplaçant, ne copie aucune règle locale et ce document n'est jamais utilisé comme source active.
