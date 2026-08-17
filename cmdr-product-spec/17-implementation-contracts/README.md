---
id: implementation-contracts-readme
domain: 17-implementation-contracts
status: draft
owner: Platform Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Implementation Contracts

## Objectif

Définir les contrats d’API, événements, erreurs, permissions, contexte, objets, jobs, Endpoint Agent, Studio et Govern.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Platform Architecture Lead.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Contracts are product-facing, implementation-agnostic.
- Canonical objects remain in domain model.
- Every mutation is idempotent where needed and audited.
- Tenant scope is mandatory.

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
