---
id: shared-business-service-catalog
domain: 12-shared-capabilities
status: draft
owner: Shared Capabilities Product Lead
updated: 2026-08-16
source-of-truth: canonical
---
# Business Service Catalog

## Objectif

Définir la capacité partagée Business Service Catalog et fournir une identité de contexte réutilisable par les projections Health/SLO sans créer un nouvel objet Service.

## Périmètre

Le catalogue reste Shared-owned. ADR-0009 permet à une projection SLO de référencer une Business Service catalog identity comme typed subject reference; cela ne transforme ni Service ni SLO en nouvel objet canonique.

## Propriétaire fonctionnel

Shared Capabilities Product Lead.

## Fonctionnalités

- Service identity.
- Owners.
- Criticality.
- Dependencies.
- Tenant scope.

## SLO / Health relationship

Une projection SLO/Health peut référencer l'identité stable fournie par ce catalogue, avec son Tenant et sa source autoritative. Le catalogue :
- ne possède pas le SLO target ;
- ne calcule pas le SLO state/breach ;
- ne publie pas un Customer commitment ;
- ne devient pas un runtime de monitoring.

Les Customer/contract projections restent externes et la publication external/customer reste gouvernée séparément.

## UX et interactions

- Navigation par liens stables.
- Tenant et identité Service restent visibles dans les consumers.
- Aucune duplication des définitions externes.

## Permissions

Les permissions du consumer et du contexte Tenant continuent de s'appliquer. Une relation Service↔SLO n'accorde aucun droit.

## Dépendances

- `../00-governance/adr/ADR-0009-slo-health-resilience-source-ownership-and-runtime-boundary.md`
- `metrics-engine.md`
- `../17-implementation-contracts/health-contract.md`

## Critères d’acceptation

- Service identity peut être référencée sans nouvelle définition d'objet.
- Ownership Service, SLO target et runtime restent séparés.
- Tenant scope est préservé.

## Questions ouvertes

- Les capacités natives/intégrées et volumes restent source/deployment-dependent.
- `OPEN-019` reste applicable à toute diffusion external/customer.
