---
id: deprecated-experience-back-and-return-behavior
domain: 04-experience-architecture
status: deprecated
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: deprecated
replaced-by: history-and-back.md
requirements:
  - REQ-UX-001
  - REQ-UX-006
  - REQ-UX-007
---

# Pointeur déprécié

## Remplaçant

[`history-and-back.md`](history-and-back.md)

## Justification

La responsabilité a été consolidée dans une source Phase 3 unique et testable afin de supprimer les architectures concurrentes.

## Migration

Les consommateurs migrent vers `history-and-back.md` et conservent uniquement leur usage local.

## Dépendants

Produits, écrans et composants qui citaient cet ancien document.

## Date de retrait

2026-08-03. Le chemin reste disponible pour l'historique et les liens de migration, mais ne porte plus de règle normative.

## Critère d'acceptation

**Given** un consommateur de cet ancien chemin,  
**When** sa dépendance est mise à jour,  
**Then** il référence le remplaçant, ne copie aucune règle locale et ce document n'est jamais utilisé comme source active.
