---
id: OBJ-ALERT
type: object
domain: 05-domain-model
owner: Command
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Alert

## Objectif

Représenter une unité de triage nécessitant une disposition explicite.

## Propriétaire

Command.

## Position dans la chaîne

- Amont: Signal
- Aval: Incident

## Champs canoniques

- `id`
- `tenant-id`
- `signals`
- `severity`
- `status`
- `assignee`
- `disposition`
- `sla`

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.

## États

- new
- acknowledged
- triaged
- linked
- closed

## Invariants

- Le tenant est obligatoire et immuable.
- La provenance vers l’amont est conservée.
- Les mutations et transitions sont auditées.
- Les références sont stables et ne transfèrent pas les permissions.

## Permissions

- `perm.command.alert.read`
- `perm.command.alert.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
