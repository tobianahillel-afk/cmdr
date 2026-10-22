---
id: OBJ-SIGNAL
type: object
domain: 05-domain-model
owner: Command
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Signal

## Objectif

Représenter un indice opérationnel corrélé ou enrichi qui mérite évaluation.

## Propriétaire

Command.

## Position dans la chaîne

- Amont: Detection
- Aval: Alert

## Champs canoniques

- `id`
- `tenant-id`
- `source-detections`
- `entities`
- `confidence`
- `priority`
- `status`

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.

## États

- new
- correlated
- promoted
- dismissed

## Invariants

- Le tenant est obligatoire et immuable.
- La provenance vers l’amont est conservée.
- Les mutations et transitions sont auditées.
- Les références sont stables et ne transfèrent pas les permissions.

## Permissions

- `perm.command.signal.read`
- `perm.command.signal.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
