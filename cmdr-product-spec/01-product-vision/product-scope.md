---
id: product-scope
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-001
  - REQ-PROD-012
  - REQ-PROD-019
---
# Périmètre produit et programme

## Dans le périmètre de la vision

- coordination, investigation, preuve, décision, réponse, rollback, audit et amélioration ;
- automatisation et IA gouvernées ;
- administration de plateforme et Endpoint Agent cible ;
- capabilities natives et intégrées classifiées ;
- expérience continue entre produits.

## Dans le périmètre de Phase 1

- mission, principes, utilisateurs, problèmes, valeur ;
- frontières et propriété ;
- modèle opérationnel ;
- non-objectifs, risques, métriques conceptuelles ;
- taxonomie de delivery ;
- questions ouvertes.

## Phases ultérieures

| Sujet | Phase |
|---|---|
| marque et identités | 2 |
| UX, shells et Design System | 3 |
| fonctionnalités détaillées | 4 |
| parcours | 5 |
| écrans pilotes et autres écrans | 6 |
| objets, permissions et contrats fonctionnels | 7 |
| architecture technique | 8 |
| implémentation | 9 |

## Hors périmètre actuel

Protocoles, certificats, algorithmes, microservices, bases de données, files de messages, frameworks, drivers, budgets CPU, sizing et déploiement final.

## Règle de portée

Une capability cible peut être dans le scope produit tout en restant `planned`. Elle ne devient `implemented` qu'avec preuve.

## Critère d'acceptation

**Given** un choix de protocole Endpoint,  
**When** il est proposé en Phase 1,  
**Then** il est enregistré comme dépendance ultérieure et ne remplace pas le travail de frontière et d'expérience produit.
