---
id: deprecated-transition-studio-to-govern
domain: 04-experience-architecture
status: deprecated
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: deprecated
replaced-by: ../cross-product-transitions.md
requirements:
  - REQ-PROD-008
  - REQ-UX-006
  - REQ-UX-007
---

# Pointeur déprécié

## Remplaçant

[`../cross-product-transitions.md`](../cross-product-transitions.md)

## Justification

Les transitions ont été consolidées dans une matrice unique couvrant contexte, permissions, erreurs, retour et trace.

## Migration

L’ancien scénario `studio-to-govern` devient une ligne/variante du contrat canonique ; ses parcours détaillés appartiennent à la Phase 5.

## Dépendants

Parcours, écrans et liens profonds citant ce scénario.

## Date de retrait

2026-08-03. Le chemin reste disponible pour l'historique et les liens de migration, mais ne porte plus de règle normative.

## Critère d'acceptation

**Given** un consommateur de cet ancien chemin,  
**When** sa dépendance est mise à jour,  
**Then** il référence le remplaçant, ne copie aucune règle locale et ce document n'est jamais utilisé comme source active.
