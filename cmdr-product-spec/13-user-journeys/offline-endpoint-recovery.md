---
id: journey-offline-endpoint-recovery
domain: 13-user-journeys
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: canonical
---
# Offline Endpoint Recovery

## Objectif

Définir le parcours Offline Endpoint Recovery.

## Périmètre

Document canonique du domaine. Il définit uniquement son sujet et renvoie vers les autres sources de vérité pour les concepts partagés.

## Propriétaire fonctionnel

Product Architecture.

## Objets concernés

- Concepts du document
- Références canoniques liées

## Fonctionnalités

- Agent offline detection.
- Local safe mode.
- Queued telemetry/commands.
- Reconnect.
- Deduplicate and verify.

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

- Quels rôles exacts exécutent chaque étape?
- Quels délais et SLO sont requis?
