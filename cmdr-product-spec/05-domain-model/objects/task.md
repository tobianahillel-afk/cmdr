---
id: OBJ-TASK
type: object
domain: 05-domain-model
owner: Shared Capabilities
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Task

## Objectif

Représenter un travail assigné sans modifier automatiquement l’objet source.

## Propriétaire

Shared Capabilities.

## Position dans la chaîne

- Amont: aucun objet canonique direct
- Aval: aucun objet canonique direct

## Champs canoniques

- `id immuable`
- `tenant-id`
- `created-at`
- `updated-at`
- `version`

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.

## États

- open
- in-progress
- blocked
- done
- cancelled

## Invariants

- Le tenant est obligatoire.
- Les mutations sont auditées.
- Les références utilisent des identifiants stables.

## Permissions

- `perm.shared-capabilities.task.read`
- `perm.shared-capabilities.task.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
