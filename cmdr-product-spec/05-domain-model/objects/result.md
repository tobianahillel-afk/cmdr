---
id: OBJ-RESULT
type: object
domain: 05-domain-model
owner: Govern
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Result

## Objectif

Représenter le résultat vérifié d’un Response Run et son effet opérationnel.

## Propriétaire

Govern.

## Position dans la chaîne

- Amont: Response Run
- Aval: aucun objet canonique direct

## Champs canoniques

- `id`
- `tenant-id`
- `run-id`
- `outcome`
- `target-results`
- `verification`
- `residual-risk`
- `created-at`

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.

## États

- draft
- verified
- disputed
- superseded

## Invariants

- Le tenant est obligatoire et immuable.
- La provenance vers l’amont est conservée.
- Les mutations et transitions sont auditées.
- Les références sont stables et ne transfèrent pas les permissions.

## Permissions

- `perm.govern.result.read`
- `perm.govern.result.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
