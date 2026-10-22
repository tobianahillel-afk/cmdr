---
id: product-vision-readme
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-001
  - REQ-PROD-020
  - REQ-AI-001
---
# Vision produit CMDR

## Mission du domaine

`01-product-vision/` définit pourquoi CMDR existe, la destination recherchée, les personnes servies, les problèmes traités, la valeur, les frontières, le modèle opérationnel, les risques et les horizons futurs. Il ne décrit ni l'UI détaillée ni l'architecture technique.

## Propriétaire et statut

- **Propriétaire :** Head of Product.
- **Statut :** Draft, prêt pour une revue produit structurée mais non validé.
- **Consommateurs :** Product, UX, Design, Security, Engineering, Content, QA et tous les domaines produits.

## Requirement IDs principaux

`REQ-PROD-001` à `REQ-PROD-062`, `REQ-AI-001` à `REQ-AI-011`, `REQ-OBJ-001` à `REQ-OBJ-012`, `REQ-SEC-001` à `REQ-SEC-005`.

## Ordre de lecture

1. [`product-vision.md`](product-vision.md)
2. [`product-mission.md`](product-mission.md)
3. [`product-principles.md`](product-principles.md)
4. [`target-users.md`](target-users.md)
5. [`user-problems.md`](user-problems.md)
6. [`value-proposition.md`](value-proposition.md)
7. [`product-boundaries.md`](product-boundaries.md)
8. [`product-pillars.md`](product-pillars.md)
9. [`platform-overview.md`](platform-overview.md)
10. [`operating-model.md`](operating-model.md)
11. [`success-metrics.md`](success-metrics.md)
12. [`non-goals.md`](non-goals.md)
13. [`product-risks.md`](product-risks.md)
14. [`future-vision.md`](future-vision.md)

## Documents complémentaires

- [`north-star.md`](north-star.md) — résultat d'expérience synthétique.
- [`market-positioning.md`](market-positioning.md) — différenciation sans promesse de livraison.
- [`capability-map.md`](capability-map.md) — propriétaire, cible et statut documentaire des capabilities.
- [`portfolio-map.md`](portfolio-map.md) — relations entre produits.
- [`product-scope.md`](product-scope.md) — périmètre de programme et horizons.
- `personas.md` — alias déprécié vers `target-users.md`.

## Frontières du domaine

Ce domaine décide la mission, les principes et la propriété produit. Il ne décide pas :

- des palettes Investigate, Govern ou Studio ;
- des écrans détaillés ;
- des schémas d'objets ;
- des permissions détaillées ;
- des protocoles Endpoint Agent ;
- des APIs, services ou infrastructures ;
- de l'état réel d'implémentation sans preuve.

## Dépendances

- `00-governance/source-material/` fournit les décisions sources.
- `00-governance/ownership-register.md` définit la propriété canonique.
- les phases 2 à 8 appliquent la vision à la marque, l'UX, les fonctionnalités et la technique.

## Décisions ouvertes

Les décisions `OPEN-001` à `OPEN-015` restent enregistrées dans `00-governance/source-material/unresolved-decisions.md`. Les documents de ce domaine citent les questions qui affectent leur périmètre sans les fermer.

## Contribution

Toute modification doit :

- citer les Requirement IDs ;
- distinguer target, planned, proposed et état prouvé ;
- mettre à jour les dépendants et la matrice ;
- éviter les slogans non mesurables ;
- conserver l'usage sans IA ;
- ne pas déplacer une frontière sans ADR.

## Critère d'acceptation

**Given** un nouveau lecteur,  
**When** il suit l'ordre de lecture,  
**Then** il peut expliquer ce qu'est CMDR, pour qui, pourquoi, comment les six produits coopèrent, ce qu'ils ne possèdent pas, comment l'IA reste optionnelle et quelles décisions restent ouvertes, sans supposer une implémentation déjà livrée.
