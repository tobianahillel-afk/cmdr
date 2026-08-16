---
id: shared-metrics-engine
domain: 12-shared-capabilities
status: draft
owner: Shared Capabilities Product Lead
updated: 2026-08-16
source-of-truth: canonical
---
# Metrics Engine

## Objectif

Définir la capacité partagée Metrics Engine comme mécanisme générique de métriques, sans devenir owner universel des targets ou calculs SLO.

## Périmètre

Shared conserve metric definitions/version, reconciliation, aggregation/freshness mechanisms et privacy thresholds. ADR-0009 attribue aux authoritative source/runtime owners l'acquisition et les calculs source-specific availability/reliability/SLO faisant autorité.

## Propriétaire fonctionnel

Shared Capabilities Product Lead.

## Fonctionnalités

- Metric definitions/version.
- Aggregation mechanisms.
- Reconciliation.
- Privacy thresholds.
- Freshness.
- Snapshot/provenance references pour les consumers autorisés.

## SLO boundary

Metrics Engine peut fournir les primitives génériques référencées par une projection SLO, mais :
- Metric != SLO target ;
- Metric != contractual SLA ;
- aucun generic SLO target store n'est créé ;
- aucun central CMDR SLO calculator n'est créé ;
- une projected compliance/breach reste dépendante de la source/calculation provenance autoritative.

## UX et interactions

Les consumers peuvent expliquer definition/version, window, freshness et provenance lorsqu'elles sont disponibles. Missing/partial/stale restent explicites.

## Permissions

Une permission de lecture métrique n'implique ni export, ni SLO configuration, ni runtime action. Reporting/Export restent des mécanismes Shared distincts.

## Dépendances

- `../00-governance/adr/ADR-0009-slo-health-resilience-source-ownership-and-runtime-boundary.md`
- `../17-implementation-contracts/metrics-contract.md`
- `../17-implementation-contracts/health-contract.md`

## Critères d’acceptation

- Shared conserve les mécanismes génériques sans absorber le runtime source.
- Les valeurs chiffrées restent traçables à leur definition/version et provenance.
- Aucun calculateur SLO universel n'est revendiqué.

## Questions ouvertes

- Volumes et support effectifs restent implementation/source-dependent et liés à `OPEN-008`.
