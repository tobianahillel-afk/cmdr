---
id: review-approval-policy
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# Politique de revue et d’approbation

## Objectif

Définir les revues nécessaires avant qu’un document devienne normatif.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Product Architecture.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Faire relire les frontières produit par tous les propriétaires concernés.
- Faire relire les permissions par Security Architecture.
- Faire relire les écrans par Design et Accessibilité.
- Faire échouer la revue en présence de liens cassés ou doublons.

## UX et interactions

- Navigation par liens stables.
- Contenu lisible en thème clair et sombre.
- Aucune duplication des définitions externes.

## Permissions

Les modifications suivent le modèle défini dans `../14-security-permissions-and-trust/permission-model.md` lorsque le document décrit une capacité exécutable.

## États

Le statut documentaire suit `00-governance/document-status-model.md`; les états métier restent dans leurs sources canoniques.

## Dépendances

- 00-governance/source-of-truth-policy.md

## Critères d’acceptation

- Le document a un propriétaire unique.
- Les liens locaux sont valides.
- Les décisions non tranchées sont attribuées.

## Questions ouvertes

- À compléter — décision source non fournie dans le brief canonique.
