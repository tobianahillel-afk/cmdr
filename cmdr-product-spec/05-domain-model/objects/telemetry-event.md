---
id: OBJ-TELEMETRY_EVENT
type: object
domain: 05-domain-model
owner: Shared Capabilities
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Telemetry Event

## Objectif

Représenter une unité immuable de télémétrie normalisée.

## Propriétaire

Shared Capabilities.

## Position dans la chaîne

- Amont: aucun objet canonique direct
- Aval: Detection

## Champs canoniques

- `id`
- `tenant-id`
- `event-time`
- `ingestion-time`
- `source`
- `raw-reference`
- `normalized-fields`
- `schema-version`
- `integrity`

## Relations

- Les relations sont typées, tenant-scoped, bidirectionnellement navigables et sourcées.
- Une relation ne transfère ni propriété ni permission.

## États

- immutable

## Invariants

- Le tenant est obligatoire et immuable.
- La provenance vers l’amont est conservée.
- Les mutations et transitions sont auditées.
- Les références sont stables et ne transfèrent pas les permissions.

## Permissions

- `perm.shared-capabilities.telemetry-event.read`
- `perm.shared-capabilities.telemetry-event.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
