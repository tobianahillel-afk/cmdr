---
id: relationship-model
domain: 05-domain-model
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# Modèle de relations

## Objectif

Définir modèle de relations.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Product Architecture.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Relations typées et directionnelles.
- Provenance obligatoire.
- Tenant isolation.
- Suppression logique sans rupture des références.

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
