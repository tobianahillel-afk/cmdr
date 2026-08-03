---
id: pattern-recovery
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-002
  - REQ-UX-004
  - REQ-UX-005
---

# Récupération

## Problème

Préserver les données valides et proposer une reprise sûre.

## Contexte

Utiliser ce pattern lorsqu'une activité traverse plusieurs états, objets ou surfaces communes et qu'un composant isolé ne suffit pas.

## Solution

Error boundary, correlation id, retry idempotent, partial results, restore draft.

## Anatomy

Déclencheur explicite, contexte et source, état courant, contenu/action, feedback persistant, sortie/retour et trace lorsqu'une mutation existe.

## Comportement

Le pattern préserve focus, sélection, filtres et travail valide. Toute mutation précise portée, permission, résultat et reprise. Les mises à jour live ne réordonnent pas silencieusement le travail.

## Composants

Utiliser les composants canoniques de `../components/` et les tokens sémantiques ; ne pas recréer d'Inspector, modal, statut ou provenance.

## Erreurs

Loading, Partial, Error, Offline et Permission denied conservent le contexte utile. Retry est idempotent ou explique le risque. Correlation ID accessible lorsque disponible.

## Accessibilité

Ordre clavier documenté, focus restauré, annonces proportionnées, alternative à couleur/drag/canvas, libellés et conséquences compréhensibles.

## Anti-patterns

Refresh complet, perte silencieuse ou succès global trompeur.

## Exemple

L'utilisateur peut interrompre, comprendre l'état, corriger l'entrée ou revenir à la source sans perdre le workspace.

## Critère d’acceptation

**Given** un utilisateur au clavier, un contexte cross-product, des données partielles et une permission limitée,  
**When** il utilise le pattern puis rencontre une erreur,  
**Then** le contexte et le focus sont préservés, l'erreur et la reprise sont explicites, aucune donnée interdite n'est révélée, la provenance reste visible et aucune IA n'est nécessaire.
