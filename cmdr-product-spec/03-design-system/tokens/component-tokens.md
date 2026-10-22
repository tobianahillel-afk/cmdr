---
id: deprecated-token-component-tokens
domain: 03-design-system
status: deprecated
owner: Design System Lead
updated: 2026-08-03
source-of-truth: deprecated
replaced-by: ../foundations/tokens.md
requirements:
  - REQ-UX-003
  - REQ-BRAND-003
  - REQ-BRAND-007
---

# Pointeur déprécié

## Remplaçant

[`../foundations/tokens.md`](../foundations/tokens.md)

## Justification

La Phase 3 consolide primitives, sémantiques, composants et slots non résolus dans une architecture unique.

## Migration

Migrer les références vers `../foundations/tokens.md` ; ne résoudre aucun slot Investigate, Govern, Studio ou typographique.

## Dépendants

Composants, thèmes et documents de marque qui citaient ce chemin.

## Date de retrait

2026-08-03. Le chemin reste disponible pour l'historique et les liens de migration, mais ne porte plus de règle normative.

## Critère d'acceptation

**Given** un consommateur de cet ancien chemin,  
**When** sa dépendance est mise à jour,  
**Then** il référence le remplaçant, ne copie aucune règle locale et ce document n'est jamais utilisé comme source active.
