---
id: OBJ-CASE
type: object
domain: 05-domain-model
owner: Investigate
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Case

## Objectif

Fournir le workspace d’investigation lié à un ou plusieurs Incidents.

## Propriétaire

Investigate.

## Position dans la chaîne

- Amont: Incident
- Aval: Evidence

## Champs canoniques

- `id`
- `tenant-id`
- `human-id`
- `title`
- `status`
- `owner`
- `investigators`
- `linked-incidents`
- `scope`
- `time-range`

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.

## États

- intake
- triage
- investigating
- awaiting-input
- review
- completed
- archived

## Invariants

- Le tenant est obligatoire et immuable.
- La provenance vers l’amont est conservée.
- Les mutations et transitions sont auditées.
- Les références sont stables et ne transfèrent pas les permissions.

## Permissions

- `perm.investigate.case.read`
- `perm.investigate.case.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
