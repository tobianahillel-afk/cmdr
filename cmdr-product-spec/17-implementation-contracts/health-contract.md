---
id: contract-health-contract
domain: 17-implementation-contracts
status: draft
owner: Platform Architecture Lead
updated: 2026-08-16
source-of-truth: canonical
---
# Health Contract

## Objectif

Définir l'enveloppe neutre Health/SLO utilisée par les projections CMDR sans sélectionner de runtime, store ou format physique.

## Périmètre

Le contrat est implementation-agnostic. ADR-0009 définit les frontières SLO/Health/Resilience. Platform Architecture possède le contrat; les sources/runtime owners gardent probes, acquisition et calculs faisant autorité.

## Propriétaire fonctionnel

Platform Architecture Lead.

## Identité de projection

Une projection SLO/Health conserve, selon disponibilité source :
- Tenant ;
- typed subject reference ;
- Environment lorsqu'il est sourcé ;
- authoritative source reference ;
- status/observation source ;
- target value/unit et target version/effective period lorsqu'un SLO target existe ;
- measurement window et calculation provenance lorsqu'un state/breach est présenté ;
- freshness ;
- details scope et limites/unknowns.

SLO n'est pas un objet canonique.

## Fonctionnalités

- Component/typed subject reference.
- Status.
- Capabilities affected.
- Freshness.
- Source-attributed SLO target/state/breach.
- Details scope.
- Unknown/partial/stale/conflicting/unsupported semantics.

## Règles

- Aucun state `healthy` ou `compliant` n'est inféré depuis l'absence de donnée.
- Un SLO state/breach n'est authoritative que si la source, target version, window, calculation provenance et freshness requises sont disponibles.
- Settings peut dériver age/stale/unknown/partial de métadonnées explicites, jamais une mesure ou compliance.
- SLO breach n'est ni Incident, ni Task, ni contractual breach automatiquement.
- Tenant identity et source provenance sont obligatoires pour l'agrégation MSSP read-only.

## Permissions et scope

La lecture Health utilise les permissions existantes et le scope Tenant. `perm.settings.health.read` n'autorise ni configuration, ni export, ni exécution. Search/Report/Export conservent leurs contrats Shared et leurs scopes single-Tenant approuvés.

## Resilience boundary

Ce contrat peut représenter sourced degradation et sourced failover/recovery state, mais ne définit aucun executor. Retry/fallback/failover/recovery/rollback restent distincts.

## Dépendances

- `../00-governance/adr/ADR-0009-slo-health-resilience-source-ownership-and-runtime-boundary.md`
- `metrics-contract.md`
- `offline-and-retry-contract.md`
- `../12-shared-capabilities/metrics-engine.md`
- `../14-security-permissions-and-trust/permission-model.md`

## Critères d’acceptation

- Source et Tenant ne sont jamais perdus.
- Missing/stale/conflicting restent explicites.
- Aucune mesure ou compliance n'est fabriquée.
- Aucun runtime ou permission n'est implicite.
- Cross-tenant aggregation reste read-only.

## Questions ouvertes

- Quel format de schéma et quelle version initiale?
- `OPEN-008` conserve la disponibilité/support effectifs des sources et runtimes.
