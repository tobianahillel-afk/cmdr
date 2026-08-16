---
id: settings-health
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-16
source-of-truth: canonical
---
# Health

## Objectif

Voir la santé, la fraîcheur, les dégradations et les projections SLO sourcées de la plateforme sans transférer à Settings la propriété des probes, mesures, calculs ou runtimes sources.

## Périmètre

Module du produit 10-platform-settings. Les objets, permissions, composants et transitions partagés sont référencés et non redéfinis. `ADR-0009` fixe la frontière SLO / Health / Resilience.

## Propriétaire fonctionnel

Platform Settings Product Lead pour l'UI/projection Health et les dérivations de présentation bornées. Platform Architecture reste owner du contrat neutre Health/Metrics. Les source/runtime owners restent owners des probes, mesures et calculs faisant autorité.

## Objets concernés

- integration
- endpoint-agent-fleet
- data-source

`SLO`, `Health Observation`, `Threshold`, `Service Level`, `Availability`, `Reliability`, `Resilience`, `RTO` et `RPO` ne deviennent pas des objets canoniques par ce module.

## Fonctionnalités

- SLO source-attributed projection.
- Freshness.
- Degradation.
- Maintenance visibility.
- Tenant-local Health/SLO.
- MSSP read-only aggregation dans l'Authorized Tenant Set.

Une projection SLO conserve au minimum Tenant, typed subject reference, authoritative source reference, target version/effective period et provenance; Environment est conservé lorsqu'il est sourcé.

## Frontière runtime

Platform Health est une couche de projection/presentation déterministe bornée. Settings ne possède pas :
- les probes ;
- l'acquisition des mesures ;
- le monitoring de service ;
- un calculateur SLO central ;
- le failover ou recovery générique.

Les états missing, stale, conflicting, partial ou unsupported restent explicites; aucun état healthy/compliant n'est inféré sans preuve source.

## Resilience

La portée initiale se limite aux dégradations sourcées, au contrat Offline/Retry existant et aux contraintes de fallback provider existantes comme configuration. Retry, fallback, failover, recovery et rollback restent distincts. Generic failover/recovery/DR/RTO/RPO restent hors scope initial.

## UX et interactions

- Conserver le contexte Tenant, vue et objet.
- Utiliser l’Inspector canonique.
- Afficher les six états obligatoires.
- Afficher source, target version, freshness et provenance pour toute projection SLO/Health.
- Préserver navigation clavier et liens profonds.
- `Acknowledge maintenance` reste désactivé/non exécutable tant qu'un owner, une action class et une permission explicite ne sont pas sourcés.

## Permissions

Voir `../../14-security-permissions-and-trust/permission-model.md` et le registre des permissions.

`perm.settings.health.read` reste strictement read-only. Il n'autorise ni configuration SLO, ni export, ni monitoring runtime, ni failover/recovery. Aucune permission Health/SLO manage n'est créée par ADR-0009.

## MSSP / Customer

La visibilité initiale Phase-6 MVP est Tenant-local plus agrégation MSSP read-only dans l'Authorized Tenant Set. Chaque projection conserve son Tenant. Customer reste une projection externe. La publication externe/customer reste bloquée par `OPEN-019`. Search reste single-selected-Tenant; Report et Export restent single-Tenant et Shared-owned.

## Handoffs

Un SLO breach sourcé peut être projeté et notifier via Shared. Il ne crée pas automatiquement Incident, Task, changement de priorité, Govern flow ou automation. Tout handoff Command initial est explicite et tenant-local.

## États

Les états métier viennent de leurs sources canoniques; la page ajoute uniquement Loading, Empty, Partial, Error, Offline et Permission denied.

## Dépendances

- `../../00-governance/adr/ADR-0009-slo-health-resilience-source-ownership-and-runtime-boundary.md`
- `../../12-shared-capabilities/metrics-engine.md`
- `../../14-security-permissions-and-trust/permission-model.md`
- `../../17-implementation-contracts/health-contract.md`
- `../../17-implementation-contracts/metrics-contract.md`
- `../../17-implementation-contracts/offline-and-retry-contract.md`

## Critères d’acceptation

- Aucune définition d’objet ou de permission locale.
- Aucune valeur SLO/Health n'est présentée comme authoritative sans source, target version, window/calculation provenance et freshness requis.
- Les états inconnus/partiels restent visibles.
- Aucune mutation ou exécution runtime n'est rendue possible par `perm.settings.health.read`.
- Aucun cross-tenant mutation/admin/response/export widening.
- Tous les écrans du module conservent leur contrat de sections et IDs existants.

## Questions ouvertes

- `OPEN-008` conserve la décision d'implémentation sur disponibilité/support des sources.
- `OPEN-013` conserve l'autorité par défaut des futures mutations Class 2.
- `OPEN-015` conserve le bridge Automation Run / Response Run.
- `OPEN-019` conserve dissemination/releasability/sharing/access externe.
