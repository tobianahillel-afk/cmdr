---
id: product-risks
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
---
# Risques produit

## Objectif

Centraliser les risques produit avec propriétaire et traitement.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Head of Product.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Risque de surcharge des écrans — Design Lead.
- Risque de sources concurrentes — Product Architecture.
- Risque d’automatisation non gouvernée — Govern et Studio.
- Risque d’isolation tenant insuffisante — Security Architecture.
- Risque de capacités endpoint trop faibles — Endpoint Agent Owner.
- Risque d’outils forensics incomplets — Investigate Owner.

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

- Quels éléments nécessitent une validation utilisateur ou marché?
- Quelles mesures deviennent des objectifs contractuels?
