---
id: status-lifecycle
domain: 00-governance
status: draft
owner: Documentation Governance Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-012
---
# Cycle de vie des documents

## Principes

Le statut décrit le niveau de décision, de preuve et d'adoption. Une réécriture importante ne justifie pas automatiquement `validated`. Les états métier d'un objet sont indépendants du statut documentaire.

## Draft

Conditions minimales :

- contenu non approuvé ;
- propriétaire, date et source indiqués ;
- Requirement IDs applicables ;
- dépendances et consommateurs connus ;
- décisions ouvertes autorisées si elles sont structurées ;
- aucune promesse d'implémentation ;
- propositions distinguées des décisions.

Sortie possible : `in-review` lorsque le contenu est substantiel et que les gates de revue sont satisfaites.

## In Review

Conditions d'entrée :

- contenu spécifique, non limité au template ;
- aucune anomalie P0 connue dans le périmètre ;
- reviewers désignés ;
- critères fonctionnels observables ;
- contradictions résolues ou décision formellement demandée ;
- dépendants et impacts identifiés ;
- liens, noms et Requirement IDs vérifiés.

Résultats possibles :

- retour à `draft` avec demandes de modification ;
- maintien `in-review` avec actions ouvertes non bloquantes ;
- passage à `validated` après approbations.

## Validated

Conditions :

- approbations requises obtenues et enregistrées ;
- aucune décision bloquante ouverte ;
- traçabilité complète ;
- périmètre validé explicite ;
- version et date ;
- preuves de revue accessibles ;
- aucune source concurrente active connue.

`validated` signifie que la spécification est approuvée, pas que le produit est implémenté.

## Implemented

Conditions :

- implémentation ou configuration liée ;
- tests liés aux critères d'acceptation ;
- version produit ou release ;
- preuves d'acceptation ;
- écarts connus et décisions de traitement ;
- statut de la spécification toujours cohérent avec le comportement livré.

## Deprecated

Conditions :

- remplaçant ou décision d'abandon ;
- justification ;
- plan de migration ;
- date de retrait ;
- consommateurs identifiés ;
- liens normatifs supprimés ;
- historique conservé.

## Transitions autorisées

```text
draft → in-review → validated → implemented → deprecated
draft → deprecated
in-review → draft
validated → draft        seulement après découverte d'une contradiction majeure
implemented → draft      lorsqu'un écart invalide la spécification, avec incident de gouvernance
```

Toute transition en arrière enregistre raison, auteur, date, impact et dépendants.

## Preuves de statut

| Statut | Preuve minimale |
|---|---|
| draft | propriétaire, source, exigences |
| in-review | reviewers, checklist, contenu substantiel |
| validated | approbations, matrice, absence de blocage |
| implemented | lien code/config, tests, release |
| deprecated | remplaçant, migration, date |

## Critère d'acceptation

**Given** un document portant `validated`,  
**When** sa preuve est vérifiée,  
**Then** les approbateurs, Requirement IDs, périmètre, date, version, décisions ouvertes non bloquantes et dépendants sont accessibles ; sinon le document revient à `draft`.
