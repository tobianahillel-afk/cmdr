---
id: value-proposition
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-001
  - REQ-PROD-012
  - REQ-PROD-013
  - REQ-PROD-018
---
# Proposition de valeur

## Valeur plateforme

CMDR propose une continuité opérationnelle et probatoire : les équipes travaillent sur des objets reliés, conservent le contexte, distinguent fait, conclusion, décision et résultat, et peuvent intégrer moteurs natifs ou externes sans organiser leur activité autour des marques fournisseurs.

## Valeur par produit

| Produit | Valeur différenciante | Limite de promesse |
|---|---|---|
| Command | Une situation opérationnelle, une Work Queue et des handovers reliés à impact, décisions et résultats | ne remplace pas le workbench d'analyse |
| Investigate | Recherche, collecte, hypothèses, Evidence et Findings dans un Case cohérent | les moteurs avancés sont classifiés progressivement |
| Govern | Décisions et réponses fondées sur preuve, impact, autorité, rollback et vérification | ne promet pas une gouvernance instantanée sans politique |
| CMDR Studio | Construction, assurance et supervision communes des automatisations optionnelles | ne remplace pas les produits opérationnels |
| Platform Settings | Administration cohérente des tenants, rôles, sources, providers et flotte | n'est pas un espace d'investigation |
| Endpoint Agent | Cible d'un EDR natif relié à investigation et gouvernance | cible `planned`, non déclarée livrée |
| Shared Capabilities | Services communs réutilisables sans duplication | n'est pas un septième produit utilisateur |

## Différenciation

CMDR n'est pas différencié par « plus d'IA » ou « tout-en-un ». Il est différencié par :

- la continuité entre contexte opérationnel, preuve et autorité ;
- la propriété canonique des objets ;
- l'usage conjoint d'humains, règles et automatisations ;
- la capacité à gouverner l'action et vérifier son résultat ;
- la possibilité de remplacer une intégration temporaire sans changer l'activité utilisateur.

## Formulations interdites

Ne pas utiliser `révolutionnaire`, `totalement autonome`, `intelligent par défaut`, `tout-en-un complet` ou `natif` sans classification et preuve.

## Critère d'acceptation

**Given** une présentation de l'Endpoint Agent,  
**When** elle décrit sa valeur,  
**Then** elle indique la cible native, les produits qui administrent, utilisent et gouvernent l'agent, et distingue clairement la vision de l'état livré.
