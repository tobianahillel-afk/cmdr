---
id: deprecated-layout-three-pane-layout
domain: 03-design-system
status: deprecated
owner: Design System Lead
updated: 2026-08-03
source-of-truth: deprecated
replaced-by: technical-workbench-shell.md
requirements:
  - REQ-UX-001
  - REQ-UX-002
---

# Pointeur déprécié

## Remplaçant

[`technical-workbench-shell.md`](technical-workbench-shell.md)

## Justification

La Phase 3 remplace les layouts génériques par huit shells d’activité avec anatomie, clavier, responsive et critères testables.

## Migration

Les consommateurs choisissent le shell `technical-workbench-shell.md` puis appliquent uniquement les adaptations prévues.

## Dépendants

Écrans et modules qui citaient cet ancien layout.

## Date de retrait

2026-08-03. Le chemin reste disponible pour l'historique et les liens de migration, mais ne porte plus de règle normative.

## Critère d'acceptation

**Given** un consommateur de cet ancien chemin,  
**When** sa dépendance est mise à jour,  
**Then** il référence le remplaçant, ne copie aucune règle locale et ce document n'est jamais utilisé comme source active.
