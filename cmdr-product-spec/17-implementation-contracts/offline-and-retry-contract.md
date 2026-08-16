---
id: contract-offline-and-retry-contract
domain: 17-implementation-contracts
status: draft
owner: Platform Architecture Lead
updated: 2026-08-16
source-of-truth: canonical
---
# Offline And Retry Contract

## Objectif

Définir le contrat Offline/Retry sans le transformer en architecture générique de failover ou recovery.

## Périmètre

Le contrat conserve ses primitives existantes. ADR-0009 borne la Resilience initiale à degradation observability, Offline/Retry existant et provider fallback constraints comme configuration.

## Propriétaire fonctionnel

Platform Architecture Lead.

## Fonctionnalités

- Retry class.
- Backoff.
- Expiry.
- Queue.
- Conflict.
- User visibility.

## Distinctions normatives

- retry != fallback ;
- fallback != failover ;
- failover != recovery ;
- recovery != rollback ;
- offline != permission d'exécuter une mutation ;
- retry policy != SLO target.

Le provider fallback configuré ailleurs ne prouve jamais un automatic failover.

## Hors scope initial

Ce contrat ne définit pas :
- generic runtime redundancy ;
- automatic ou manual failover ;
- data-source fallback ;
- generic recovery ;
- DR ;
- RTO/RPO ;
- un executor de résilience.

## UX et interactions

Les états Offline/Retry montrent dernière synchronisation, retry state et limites sans présenter un failover/recovery inexistant.

## Permissions

Aucune permission de failover/recovery n'est créée. Une opération future effectful doit être séparément sourcée et autorisée.

## Dépendances

- `../00-governance/adr/ADR-0009-slo-health-resilience-source-ownership-and-runtime-boundary.md`
- `health-contract.md`
- `../14-security-permissions-and-trust/permission-model.md`

## Critères d’acceptation

- Les primitives retry restent distinctes des actions runtime de résilience.
- Aucun executor ou SLO n'est inventé.
- Unknown/offline ne devient jamais healthy/compliant.

## Questions ouvertes

- Quel format de schéma et quelle version initiale?
- Backpressure générique, failover, recovery, DR, RTO/RPO restent des décisions futures source-auditées.
