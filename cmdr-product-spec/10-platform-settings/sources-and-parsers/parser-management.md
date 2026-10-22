---
id: settings-sources-and-parsers-parser-management
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-14
source-of-truth: canonical
open_decisions: [OPEN-008, OPEN-013]
---
# Parser Management

## Objectif

Définir l'administration versionnée du Parser canonique dans Platform Settings sans définir ni revendiquer son moteur d'exécution.

## Périmètre

Le document encadre le cycle de vie administratif, le versionnement, le contrat de transformation, les fixtures, les erreurs/indicateurs de qualité, la validation de configuration et le rollback administratif lorsqu'il est explicitement source-backed.

## Propriétaire fonctionnel

Platform Settings Product Lead.

## Objet canonique

`OBJ-PARSER` reste l'unique objet canonique Parser de ce périmètre. Son owner objet reste Platform Settings.

États canoniques exacts:
- `draft`
- `testing`
- `active`
- `degraded`
- `retired`

Le tenant est obligatoire. Environment reste contextuel/source-dependent lorsqu'une source canonique l'établit.

## Contrat de transformation

Le Parser Contract reste détenu par Platform Architecture et expose seulement les dimensions canoniques actuelles:
- Input dialect.
- Output schema.
- Version.
- Fixtures.
- Errors.
- Quality.

Les questions de format de schéma, version initiale, SLO et limites restent ouvertes. Ce document ne sélectionne ni ECS, ni OCSF, ni CIM, ni OpenTelemetry, ni aucun autre standard non sourcé.

## Fonctionnalités

- Version.
- Fixtures.
- Test definition and local validation.
- Source-backed administrative rollback.
- Error and quality metadata.
- Audit and provenance.

## Aucune relation Source → Parser inventée

Les objets `Data Source` et `Parser` n'ont actuellement aucune relation canonique directe amont/aval. Ce document ne crée donc ni `ParserAssignment`, ni `SourceParserAssignment`, ni objet `ParserCompatibility`, ni route, sélection, fallback, precedence ou relation Source→Parser.

## Test Parser

La validation strictement locale d'une version, d'un schéma référencé, d'une fixture ou de préconditions peut être sans effet. Si un test exige l'exécution d'un moteur Parser, d'un plugin, d'un sandbox, d'un stream processor ou d'un autre runtime, Settings se limite à la demande, aux préconditions, au handoff et au résultat observé sourcé. Une fixture ou un état `testing` n'est pas une preuve de runtime réussi.

## Rollback

Le rollback administratif signifie revenir vers une version/configuration source-backed et auditable lorsque le corpus l'autorise. Il ne signifie pas rollback d'un pipeline runtime, réécriture de données déjà ingérées ou changement automatique de parser.

## Hors périmètre runtime

Ce document ne définit ni parser execution engine, sandbox, plugin runtime, normalization engine, stream processor, schema-registry implementation, automatic parser selection, source-parser assignment, precedence, fallback execution ni storage.

## UX et interactions

- Navigation par liens stables.
- Présenter version, état, erreurs/qualité et provenance sans transformer une projection en fait runtime.
- Préserver tenant, environnement contextuel et objet actif.

## Permissions

Réutiliser `perm.platform-settings.parser.read` et `perm.platform-settings.parser.manage`, ainsi que les aliases UI existants de Sources & Parsers. Aucune permission nouvelle n'est créée.

## Audit et provenance

Toute mutation conserve acteur, tenant, objet, action, résultat, justification, version et identifiant de corrélation. Tout résultat de test/runtime projeté conserve sa source et ne devient pas autorité Settings.

## Dépendances

- `../../05-domain-model/objects/parser.md`
- `../../17-implementation-contracts/parser-contract.md`
- `../../17-implementation-contracts/event-contract.md`
- `../../14-security-permissions-and-trust/permission-model.md`

## Critères d’acceptation

- Les cinq états canoniques ne sont pas étendus localement.
- Tenant reste obligatoire.
- Aucune relation Source→Parser n'est inventée.
- Aucun standard de schéma n'est sélectionné.
- Aucun moteur de parsing n'est attribué à Settings.
- OPEN-008 et OPEN-013 restent ouverts.

## Questions ouvertes

Format/version initiale de schéma, SLO/limites, compatibilité/association Source↔Parser et runtime effectif restent non résolus sauf décision canonique future.
