---
id: settings-health-platform-health-model
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-16
source-of-truth: canonical
---
# Platform Health Model

## Objectif

Définir la projection Platform Health alignée sur `ADR-0009`, sans devenir runtime de probe, acquisition, monitoring ou calcul SLO.

## Périmètre

Platform Settings possède la présentation/projection Health et les dérivations déterministes bornées de présentation. Platform Architecture possède le contrat neutre Health/Metrics. Les sources/runtime owners conservent les mesures et calculs faisant autorité.

## Propriétaire fonctionnel

Platform Settings Product Lead pour la surface Health. Les ownerships source/runtime restent inchangés.

## Identité et source

SLO reste une projection/configuration source-attributed non canonique. La projection référence :
- Tenant ;
- typed subject reference ;
- authoritative source reference ;
- target value/unit ;
- target version et effective period ;
- measurement window ;
- calculation provenance lorsqu'un état/breach est présenté ;
- freshness ;
- Environment seulement lorsqu'il est sourcé.

Aucun objet canonique SLO/Health Observation/Threshold/Service Level/Availability/Reliability/Resilience/RTO/RPO n'est créé.

## Fonctionnalités

- Service health projection.
- Data freshness.
- Source-attributed SLO target/state/breach visibility.
- Impact mapping.
- Tenant-local visibility.
- MSSP read-only aggregation.
- Customer visibility uniquement lorsqu'une politique future autorise la publication externe.

## États et calcul

Settings peut dériver des faits de présentation bornés comme age, stale, unknown ou partial depuis les métadonnées sources explicites. Settings ne synthétise ni mesure, ni SLO compliance, ni breach lorsque la preuve autoritative manque.

Missing, stale, conflicting et unsupported ne sont jamais normalisés en healthy.

## Resilience boundary

La portée initiale comprend degradation observability, Offline/Retry existant et provider fallback constraints comme configuration. Elle n'inclut pas generic redundancy, automatic/manual failover, data-source fallback, recovery, DR, RTO ou RPO.

## MSSP et Tenant

L'agrégation MSSP read-only est bornée par l'Authorized Tenant Set Security. Le Tenant de chaque projection reste visible. Toute mutation/admin/response reste single-selected-Tenant. Search reste single-selected-Tenant; Report/Export restent single-Tenant.

## Breach handoff

Chaîne initiale :
`authoritative source/runtime → sourced SLO state/breach → Platform Health projection → optional Shared Notification → explicit selected-Tenant human handoff → existing Command action`.

Aucun Incident, Task, changement de priorité, Govern flow ou automation n'est créé automatiquement.

## Permissions

`perm.settings.health.read` couvre uniquement la lecture des projections autorisées. Aucune permission de configuration Health/SLO n'est créée. `Acknowledge maintenance` reste désactivé/non exécutable tant qu'une permission et une autorité ne sont pas séparément sourcées.

## UX et interactions

- Source, freshness, Tenant et target version restent visibles.
- Unknown/partial/stale sont distingués.
- Les actions interdites restent désactivées avec raison sans révéler de données protégées.
- Le contexte Tenant est explicite avant tout handoff tenant-local.

## Dépendances

- `../../00-governance/adr/ADR-0009-slo-health-resilience-source-ownership-and-runtime-boundary.md`
- `../../12-shared-capabilities/business-service-catalog.md`
- `../../12-shared-capabilities/metrics-engine.md`
- `../../14-security-permissions-and-trust/permission-model.md`
- `../../17-implementation-contracts/health-contract.md`
- `../../17-implementation-contracts/metrics-contract.md`

## Critères d’acceptation

- Le document conserve un owner unique de projection.
- Aucun runtime owner n'est transféré à Settings.
- Toute projection SLO/Health est source-attributed.
- Aucun cross-tenant effect ni export widening.
- Les OPEN-008/013/015/019 restent ouvertes.

## Questions ouvertes

Les décisions de source availability/support, futures mutations, automation/response bridge et external dissemination restent respectivement dans `OPEN-008`, `OPEN-013`, `OPEN-015` et `OPEN-019`.
