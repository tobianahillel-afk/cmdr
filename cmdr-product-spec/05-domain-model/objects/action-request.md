---
id: OBJ-ACTION_REQUEST
type: object
domain: 05-domain-model
owner: Govern
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Action Request

## Objectif

Formaliser une demande d’action gouvernée fondée sur un Finding ou une base d’urgence.

## Propriétaire

Govern.

## Position dans la chaîne

- Amont: Finding
- Aval: Decision

## Champs canoniques

- `id`
- `tenant-id`
- `requester`
- `targets`
- `action-type`
- `justification`
- `finding-links`
- `risk-of-action`
- `risk-of-inaction`
- `rollback-plan`
- `status`

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.

## États

- draft
- submitted
- policy-check
- awaiting-approval
- decided
- cancelled

## Invariants

- Le tenant est obligatoire et immuable.
- La provenance vers l’amont est conservée.
- Les mutations et transitions sont auditées.
- Les références sont stables et ne transfèrent pas les permissions.

## Permissions

- `perm.govern.action-request.read`
- `perm.govern.action-request.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
