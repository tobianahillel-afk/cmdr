---
id: OBJ-DECISION
type: object
domain: 05-domain-model
owner: Govern
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Decision

## Objectif

Enregistrer l’autorisation, les conditions ou le refus d’une Action Request.

## Propriétaire

Govern.

## Position dans la chaîne

- Amont: Action Request
- Aval: Response Run

## Champs canoniques

- `id`
- `tenant-id`
- `request-id`
- `outcome`
- `authority`
- `approvers`
- `rationale`
- `conditions`
- `policy-snapshot`
- `decided-at`
- `supersedes`

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.

## États

- pending
- approved
- approved-with-conditions
- more-information-required
- refused
- superseded

## Invariants

- Le tenant est obligatoire et immuable.
- La provenance vers l’amont est conservée.
- Les mutations et transitions sont auditées.
- Les références sont stables et ne transfèrent pas les permissions.

## Permissions

- `perm.govern.decision.read`
- `perm.govern.decision.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
