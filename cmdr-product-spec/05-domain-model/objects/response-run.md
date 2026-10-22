---
id: OBJ-RESPONSE_RUN
type: object
domain: 05-domain-model
owner: Govern
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Response Run

## Objectif

Représenter l’exécution gouvernée d’une action ou d’un playbook sous une Decision.

## Propriétaire

Govern.

## Position dans la chaîne

- Amont: Decision
- Aval: Result

## Champs canoniques

- `id`
- `tenant-id`
- `decision-id`
- `playbook-version`
- `targets`
- `status`
- `steps`
- `started-at`
- `ended-at`
- `rollback-status`

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.

## États

- queued
- preparing
- running
- succeeded
- failed
- partially-succeeded
- cancelled
- rollback-pending
- rolling-back
- rolled-back
- rollback-failed

## Invariants

- Le tenant est obligatoire et immuable.
- La provenance vers l’amont est conservée.
- Les mutations et transitions sont auditées.
- Les références sont stables et ne transfèrent pas les permissions.

## Permissions

- `perm.govern.response-run.read`
- `perm.govern.response-run.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
