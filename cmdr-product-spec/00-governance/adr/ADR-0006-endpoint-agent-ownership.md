---
id: ADR-0006-endpoint-agent-ownership
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-018
  - REQ-OBJ-008
  - REQ-SEC-004
  - REQ-SEC-005
---
# ADR-0006 — Propriété et rôle de l'Endpoint Agent

## 1. Identifiant

`ADR-0006`

## 2. Titre

Propriété et rôle de l'Endpoint Agent

## 3. Statut

Draft. Cette ADR n'est pas approuvée et ne remplace aucune décision source au-delà de ce qu'elle applique explicitement.

## 4. Date

2026-08-03

## 5. Propriétaire

Product Architecture.

## 6. Décideurs attendus

Head of Product et propriétaires des domaines affectés ; Security, UX ou Engineering selon les effets décrits.

## 7. Requirement IDs

`REQ-PROD-018`, `REQ-OBJ-008`, `REQ-SEC-004`, `REQ-SEC-005`

## 8. Contexte

CMDR vise un Endpoint Agent EDR natif complet, mais la flotte, l'investigation et l'autorité d'action appartiennent à des responsabilités différentes.

## 9. Problème

Confondre composant, flotte et workflow métier crée une propriété impossible et peut contourner Govern.

## 10. Forces en présence

Administration, expérience d'investigation, sécurité locale, déploiement progressif et statut de capacité.

## 11. Options étudiées

1. Endpoint Agent propriétaire de tout le workflow.
2. Platform Settings propriétaire du composant et des actions.
3. Composant distinct ; Settings administre, Investigate utilise, Govern autorise le risque.

## 12. Décision

Endpoint Agent est un composant produit distinct et une cible native planifiée. Platform Settings possède la flotte et les policies ; Investigate consomme inspection et collecte ; Govern contrôle les actions risquées.

## 13. Justification

Cette séparation distingue administration, activité analytique et autorité sans réduire l'agent à une intégration.

## 14. Conséquences positives

- Responsabilités lisibles.
- Actions risquées gouvernées.
- UX de flotte séparée du workbench.

## 15. Conséquences négatives

- Transitions et projections à spécifier.
- Plusieurs produits consomment le même endpoint.

## 16. Risques

- Promesse native prématurée.
- Double contrôle d'une action.

## 17. Effets sur la navigation

Settings expose la flotte ; Investigate expose l'activité endpoint ; Govern expose la décision et le run.

## 18. Effets sur les objets

Endpoint partagé, Endpoint Agent et Endpoint Agent Fleet restent distincts.

## 19. Effets sur les permissions

Les permissions et classes d'action sont réévaluées à chaque produit ; aucune autorisation implicite depuis l'agent.

## 20. Effets sur les parcours

Incident to Containment et Endpoint Alert to Result appliquent cette séparation.

## 21. Effets sur les autres documents

README Endpoint, Settings, Investigate, Govern, ownership register et capability map.

## 22. Migration

Remplacer « EDR complet livré » par « cible native planifiée » tant qu'aucune preuve d'implémentation n'existe.

## 23. Critères de réévaluation

Réévaluer après décision du support initial et spécification de l'expérience endpoint.

## 24. Questions encore ouvertes

- OPEN-008 — plateformes initiales.
- Les protocoles et la PKI restent Phase 8.
