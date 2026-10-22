---
id: settings-sources-and-parsers-source-management
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-14
source-of-truth: canonical
open_decisions: [OPEN-008, OPEN-013]
---
# Source Management

## Objectif

Définir l'administration des Data Sources dans Platform Settings sans absorber leur runtime d'acquisition ou d'ingestion.

## Périmètre

Le document encadre la configuration administrative, la portée, la fraîcheur, le cycle de vie, la désactivation sûre, la validation strictement locale et la projection de santé du `Data Source` canonique.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objet canonique

`OBJ-DATA_SOURCE` reste l'unique objet canonique de source pour ce périmètre. Son owner objet reste Platform Settings.

États canoniques exacts:
- `configured`
- `active`
- `degraded`
- `disabled`

Le tenant est obligatoire. Environment reste contextuel/source-dependent lorsqu'une autre source canonique l'établit; aucune relation Environment obligatoire n'est inventée ici.

## Fonctionnalités

- Connection metadata/reference where sourced.
- Scope.
- Freshness.
- Disable safely.
- Strictly local deterministic validation.
- Sourced health projection.
- Audit and provenance.

## Integration et Secret Reference

`Integration` reste distinct de `Data Source`; une configuration de source peut seulement référencer une Integration lorsque la source canonique de cette configuration le permet. Une `Secret Reference` reste reference-only et aucune valeur de secret ne peut être lue, copiée, journalisée ou injectée dans ce document.

## Test Source

Un test local qui vérifie structure, état, tenant, version ou préconditions sans effet externe peut être une validation déterministe locale. Si `Test Source` nécessite un service externe, connecteur, collection, ingestion ou probe, Platform Settings ne revendique pas l'exécuteur: le contrat Settings se limite à la demande, aux préconditions, au handoff et à la projection sourcée du résultat observé.

## Health

La santé affichée par Settings est une projection sourcée avec provenance, fraîcheur et état connus. Une projection de santé n'est jamais une preuve que Settings a exécuté un probe externe.

## Hors périmètre runtime

Ce document ne définit ni acquisition, collection runtime, ingestion engine, connector runtime, event-processing engine, health-probe executor, parser runtime, storage engine ni support effectif d'un fournisseur/source.

## UX et interactions

- Navigation par liens stables.
- Préserver tenant, environnement contextuel, filtres sûrs et objet actif.
- Afficher état, fraîcheur, provenance et limites connues sans fabriquer les données manquantes.

## Permissions

Réutiliser `perm.platform-settings.data-source.read` et `perm.platform-settings.data-source.manage`, ainsi que les aliases UI existants du screen Sources & Parsers. Security reste propriétaire de l'autorisation. Aucune permission nouvelle n'est créée.

## Audit et provenance

Toute mutation conserve acteur, tenant, objet, action, résultat, justification, version et identifiant de corrélation. Les projections runtime conservent leur source et leur fraîcheur.

## Dépendances

- `../../05-domain-model/objects/data-source.md`
- `../../05-domain-model/objects/integration.md`
- `../../05-domain-model/objects/secret-reference.md`
- `../../17-implementation-contracts/health-contract.md`
- `../../14-security-permissions-and-trust/permission-model.md`

## Critères d’acceptation

- Les quatre états canoniques ne sont pas étendus localement.
- Tenant reste obligatoire.
- Data Source n'est pas confondu avec Integration, Connection ou Collector.
- Aucun runtime externe n'est attribué à Settings sans source canonique.
- OPEN-008 et OPEN-013 restent ouverts.

## Questions ouvertes

Le support effectif des sources, SLO/limites et l'autorité par défaut des mutations réversibles restent non résolus dans leurs sources canoniques existantes.
