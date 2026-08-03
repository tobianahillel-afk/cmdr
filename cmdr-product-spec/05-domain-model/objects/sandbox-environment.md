---
id: OBJ-SANDBOX_ENVIRONMENT
type: object
domain: 05-domain-model
owner: Platform Settings
status: draft
updated: 2026-08-03
source-of-truth: canonical
---
# Sandbox Environment

## Objectif

Représenter un environnement d’exécution isolé versionné.

## Propriétaire

Platform Settings.

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

- provisioning
- ready
- busy
- degraded
- retired

## Invariants

- Le tenant est obligatoire.
- Les mutations sont auditées.
- Les références utilisent des identifiants stables.

## Permissions

- `perm.platform-settings.sandbox-environment.read`
- `perm.platform-settings.sandbox-environment.manage`

## Audit et provenance

Toute création, transition, relation et suppression logique enregistre acteur, tenant, justification, version et identifiant de corrélation.

## Critères d’acceptation

- Le schéma ne duplique aucun autre objet.
- Le propriétaire correspond au registre de propriété.
- Les transitions invalides sont refusées côté serveur.
- Les références restent résolubles après versionnement.

## Questions ouvertes

- À compléter — contenu source non fourni dans le brief canonique.
