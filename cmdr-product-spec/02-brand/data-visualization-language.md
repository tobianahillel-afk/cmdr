---
id: brand-data-visualization-language
domain: 02-brand
status: draft
owner: Brand Design Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-BRAND-001
  - REQ-BRAND-008
  - REQ-PROD-002
  - REQ-PROD-005
---
# Langage de visualisation des données

## Principe

Une visualisation CMDR est un instrument de comparaison, d'explication ou de décision. Elle n'est jamais un décor de dashboard.

## Métadonnées obligatoires

- titre décrivant la question ;
- unité et période ;
- source et fraîcheur ;
- population ou périmètre ;
- définition des mesures ;
- légende ;
- état de données manquantes ;
- incertitude ou limite ;
- alternative tabulaire lorsque nécessaire.

## Hiérarchie

1. conclusion ou question ;
2. série principale ;
3. contexte et comparaison ;
4. annotations ;
5. source et méthode.

La série principale peut utiliser l'accent produit. Les autres séries utilisent des neutres ou une palette qualitative accessible définie en Phase 3.

## Sémantique

- une couleur de produit n'est pas une gravité ;
- une couleur sémantique conserve la même signification entre produits ;
- une échelle quantitative possède un ordre perceptible ;
- les catégories ne sont pas distinguées uniquement par couleur ;
- le rouge n'est pas utilisé pour « rendre le graphique cyber ».

## Types privilégiés

- lignes et small multiples pour évolution ;
- barres pour comparaison ;
- tables et matrices pour précision ;
- timelines pour séquence et causalité ;
- graphes de relation quand les relations sont la question ;
- cartes uniquement lorsque la géographie est décisionnelle.

## Types à éviter

- jauges circulaires décoratives ;
- donut charts pour de nombreuses catégories ;
- 3D ;
- radar charts décoratifs ;
- heatmaps sans échelle ni valeur accessible ;
- score unique sans facteurs ni intervalle ;
- animation continue.

## Produit

- Command privilégie temps, charge, priorité et impact.
- Investigate privilégie relation, timeline, comparaison et provenance.
- Govern privilégie conditions, décisions, exécutions, vérification et rollback.
- Studio privilégie version, run, qualité, coûts et erreurs.

## Critère d'acceptation

**Given** une visualisation sans couleur,  
**When** elle est lue par ses labels, formes et ordre,  
**Then** la question, les séries, la conclusion et les limites restent compréhensibles.
