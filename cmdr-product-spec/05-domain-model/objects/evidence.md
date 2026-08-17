---
id: OBJ-EVIDENCE
type: object
domain: 05-domain-model
owner: Investigate
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Evidence

## Objectif

Représenter une preuve ou un artefact dérivé avec provenance, intégrité et chaîne de garde.

## Propriétaire

Investigate.

## Position dans la chaîne

- Amont: Case
- Aval: Finding

## Champs canoniques

- `id`
- `tenant-id`
- `case-id`
- `type`
- `source`
- `collector`
- `acquired-at`
- `hashes`
- `integrity-state`
- `classification`
- `parent-evidence`

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.

## États

- registered
- collecting
- available
- verifying
- verified
- rejected
- superseded

## Invariants

- Le tenant est obligatoire et immuable.
- La provenance vers l’amont est conservée.
- Les mutations et transitions sont auditées.
- Les références sont stables et ne transfèrent pas les permissions.

## Permissions

- `perm.investigate.evidence.read`
- `perm.investigate.evidence.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
