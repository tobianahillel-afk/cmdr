---
id: quality-validation-status
domain: 16-quality-and-validation
status: draft
owner: Quality Lead
updated: 2026-08-03
source-of-truth: canonical
---
# Validation Status

## Objectif

Définir validation status pour la documentation et le produit CMDR.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Quality Lead.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Per-control pass/fail.
- Date and commit.
- No unsupported PASS.
- Failures with owner.

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

- Quelle preuve automatisée ou manuelle est requise?
- Qui signe le contrôle?
