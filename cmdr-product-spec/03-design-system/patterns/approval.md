---
id: pattern-approval
domain: 03-design-system
status: draft
owner: Design System Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-UX-002
  - REQ-UX-004
  - REQ-UX-005
  - REQ-AI-001
  - REQ-AI-007
  - REQ-AI-010
---

# Approval

## Problème

Recueillir une expression d’autorité sans la confondre avec decision.

## Contexte

Utiliser ce pattern lorsqu'une activité traverse plusieurs états, objets ou surfaces communes et qu'un composant isolé ne suffit pas.

## Solution

Request, approver, authority, sod, conditions, expiration, result, trace.

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

Approval implicite, agent auto-approbateur ou absence de justification.

## Exemple

L'utilisateur peut interrompre, comprendre l'état, corriger l'entrée ou revenir à la source sans perdre le workspace.

## Critère d’acceptation

**Given** un utilisateur au clavier, un contexte cross-product, des données partielles et une permission limitée,  
**When** il utilise le pattern puis rencontre une erreur,  
**Then** le contexte et le focus sont préservés, l'erreur et la reprise sont explicites, aucune donnée interdite n'est révélée, la provenance reste visible et aucune IA n'est nécessaire.
