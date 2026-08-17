---
id: contract-metrics-contract
domain: 17-implementation-contracts
status: draft
owner: Platform Architecture Lead
updated: 2026-08-16
source-of-truth: canonical
---
# Metrics Contract

## Objectif

Définir l'enveloppe neutre des métriques et des références nécessaires aux projections Health/SLO.

## Périmètre

Platform Architecture possède le contrat neutre. Shared conserve les mécanismes génériques de metric definition/version, aggregation, reconciliation et freshness. Les source/runtime owners conservent acquisition et calculs source-specific faisant autorité. ADR-0009 ne crée aucun calculateur SLO central.

## Propriétaire fonctionnel

Platform Architecture Lead.

## Fonctionnalités

- Definition/version.
- Dimensions.
- Aggregation/window metadata.
- Freshness.
- Privacy threshold.
- Source and calculation provenance.
- Target/version reference when a metric participates in an SLO projection.

## SLO binding

Une valeur utilisée pour projeter un SLO state/breach doit pouvoir référencer :
- authoritative source ;
- metric definition/version ;
- dimensions applicables ;
- measurement/aggregation window ;
- calculation provenance ;
- target version/effective period ;
- freshness.

Le contrat ne transforme pas une metric definition en SLO target et ne transforme pas un target en contractual SLA.

## Calculation boundary

Shared peut fournir les mécanismes génériques de définition/version/aggregation/reconciliation/freshness. L'authoritative source/runtime owner reste responsable de tout calcul source-specific availability/reliability/SLO. Aucun calculateur SLO universel n'est sélectionné.

## UX et interactions

Les consommateurs exposent source, freshness et provenance suffisantes pour expliquer une valeur. Unknown/partial/stale restent visibles.

## Permissions

Le contrat ne crée aucune permission. Lire une projection métrique n'implique ni export, ni configuration, ni execution.

## Dépendances

- `../00-governance/adr/ADR-0009-slo-health-resilience-source-ownership-and-runtime-boundary.md`
- `health-contract.md`
- `../12-shared-capabilities/metrics-engine.md`

## Critères d’acceptation

- Definition/version et provenance sont traçables.
- Une projection SLO ne peut pas perdre target version, window, source ou freshness nécessaires.
- Aucune architecture de store/calcul runtime n'est imposée.

## Questions ouvertes

- Quel format de schéma et quelle version initiale?
- Les volumes/support effectifs restent dépendants des sources et de `OPEN-008`.
