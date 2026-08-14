---
id: ADR-0001-product-separation
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-14
source-of-truth: canonical
requirements:
  - REQ-PROD-013
  - REQ-PROD-014
  - REQ-PROD-015
  - REQ-PROD-016
  - REQ-PROD-017
  - REQ-PROD-018
---
# ADR-0001 — Séparation des produits CMDR

## 1. Identifiant

`ADR-0001`

## 2. Titre

Séparation des produits CMDR

## 3. Statut

Draft. Cette ADR n'est pas approuvée et ne remplace aucune décision source au-delà de ce qu'elle applique explicitement.

## 4. Date

2026-08-03

## 5. Propriétaire

Product Architecture.

## 6. Décideurs attendus

Head of Product et propriétaires des domaines affectés ; Security, UX ou Engineering selon les effets décrits.

## 7. Requirement IDs

`REQ-PROD-013`, `REQ-PROD-014`, `REQ-PROD-015`, `REQ-PROD-016`, `REQ-PROD-017`, `REQ-PROD-018`

## 8. Contexte

CMDR doit fournir une expérience continue sans transformer chaque responsabilité en silo ni fusionner toutes les activités dans un dashboard universel.

## 9. Problème

Sans frontières explicites, Command devient un workbench, Investigate une autorité de réponse, Govern une file générale et Studio l'interface principale.

## 10. Forces en présence

Cohérence interproduit, expertise des rôles, propriété des objets, continuité du contexte, évolutivité et simplicité de navigation.

## 11. Options étudiées

1. Un produit universel unique.
2. Des produits indépendants partageant seulement des liens.
3. Six produits coordonnés, un composant Endpoint Agent et des capabilities partagées.

## 12. Décision

Retenir Command, Investigate, Govern, CMDR Studio, Platform Settings et Endpoint Agent, complétés par Shared Capabilities. Chaque produit possède une mission et des exclusions.

## 13. Justification

Cette option sépare les responsabilités sans casser le workflow de bout en bout.

## 14. Conséquences positives

- Propriété lisible.
- UX orientée activité.
- Réduction des objets concurrents.
- Handoffs explicites.

## 15. Conséquences négatives

- Transitions interproduits à concevoir.
- Risque de silos si le contexte n'est pas conservé.

## 16. Risques

- Frontières trop rigides.
- Noms d'interface interprétés comme nouveaux produits.

## 17. Effets sur la navigation

Le Global Shell expose les produits ; les transitions conservent le contexte. Shared Capabilities n'apparaît pas comme produit parallèle.

## 18. Effets sur les objets

Les objets restent chez leur propriétaire et sont projetés ailleurs.

## 19. Effets sur les permissions

Chaque produit référence le Permission Model commun ; il ne crée pas de namespace par simple convenance.

## 20. Effets sur les parcours

Command → Investigate → Govern → Command devient le parcours de référence.

## 21. Effets sur les autres documents

`01-product-vision/product-boundaries.md`, READMEs produits, ownership register et parcours.

## 22. Migration

Les anciennes consoles condensées sont réparties selon propriété ; aucune seconde architecture active n'est conservée.

## 23. Critères de réévaluation

Réévaluer si un workflow majeur ne peut être attribué sans double propriété ou si la navigation exige des duplications.

## 24. Décision Customers and Delivery

`OPEN-006` a été résolue le 2026-08-14 par `ADR-0008 — Customers, MSSP, Delivery Deployment and Cross-Tenant Architecture`. Cette résolution précise le modèle de déploiement sans modifier la séparation des produits définie par ADR-0001 : Command consomme le contexte Customer/engagement, Security conserve l'autorisation, Shared conserve Search/Reporting/Export, Settings conserve Tenant et Govern conserve l'autorité de réponse.
