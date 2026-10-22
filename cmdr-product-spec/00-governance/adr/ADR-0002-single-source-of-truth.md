---
id: ADR-0002-single-source-of-truth
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-009
---
# ADR-0002 — Une source de vérité par concept

## 1. Identifiant

`ADR-0002`

## 2. Titre

Une source de vérité par concept

## 3. Statut

Draft. Cette ADR n'est pas approuvée et ne remplace aucune décision source au-delà de ce qu'elle applique explicitement.

## 4. Date

2026-08-03

## 5. Propriétaire

Product Architecture.

## 6. Décideurs attendus

Head of Product et propriétaires des domaines affectés ; Security, UX ou Engineering selon les effets décrits.

## 7. Requirement IDs

`REQ-PROD-006`, `REQ-PROD-009`

## 8. Contexte

La fondation actuelle contient de nombreux paragraphes répétés et plusieurs consommateurs qui citent des règles propres à un autre produit.

## 9. Problème

Une définition copiée diverge, masque la propriété et rend les validations contradictoires.

## 10. Forces en présence

Traçabilité, autonomie des domaines, coût de maintenance, lisibilité et historique des décisions.

## 11. Options étudiées

1. Autoriser la duplication synchronisée.
2. Centraliser toute la documentation dans un seul fichier.
3. Définir une source propriétaire et des usages locaux référencés.

## 12. Décision

Chaque concept partagé possède une source et un propriétaire. Les consommateurs décrivent uniquement leur usage local.

## 13. Justification

La référence propriétaire réduit la divergence sans créer un document monolithique.

## 14. Conséquences positives

- Dépendants identifiables.
- Corrections appliquées une fois.
- Registres utiles.

## 15. Conséquences négatives

- Les documents consommateurs doivent être plus précis sur leur usage local.
- Une modification canonique peut avoir de nombreux dépendants.

## 16. Risques

- Liens sans application substantive.
- Source centrale trop abstraite.

## 17. Effets sur la navigation

Les écrans référencent composants, objets et capabilities ; ils ne les redéfinissent pas.

## 18. Effets sur les objets

Un objet possède une seule sémantique et machine d'état.

## 19. Effets sur les permissions

Les permissions sont définies sous Security et référencées localement.

## 20. Effets sur les parcours

Les parcours citent les objets et transitions sans recréer leurs contrats.

## 21. Effets sur les autres documents

Tous les domaines actifs, registres et matrice.

## 22. Migration

Remplacer les paragraphes copiés par des références et des configurations locales ; déprécier les sources concurrentes.

## 23. Critères de réévaluation

Réévaluer si une source devient trop large pour avoir un propriétaire responsable ou si deux usages exigent réellement deux concepts.

## 24. Questions encore ouvertes

- La normalisation détaillée des permissions est différée à la Phase 7.
