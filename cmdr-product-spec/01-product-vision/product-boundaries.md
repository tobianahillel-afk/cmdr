---
id: product-boundaries
domain: 01-product-vision
status: draft
owner: Head of Product
updated: 2026-08-03
source-of-truth: canonical
---
# Frontières produit

## Objectif

Rendre explicites les décisions qui appartiennent à chaque produit.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Head of Product.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Command décide quoi traiter et coordonne.
- Investigate décide ce qui est établi par les preuves.
- Govern décide ce qui peut être exécuté.
- Studio décide comment construire et assurer l’automatisation.
- Platform Settings administre l’environnement.
- Endpoint Agent collecte, détecte et exécute localement selon contrats signés.

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
