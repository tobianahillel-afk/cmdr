---
id: command-module-customers-and-delivery
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-14
source-of-truth: canonical
requirements:
  - REQ-PROD-053
  - REQ-PROD-013
  - REQ-PROD-019
  - REQ-PROD-033
open_decisions:
  - OPEN-013
  - OPEN-019
---
# Customers and Delivery

## État

`defined`, `deployment-dependent`, `delivery_mode: planned`. `OPEN-006` est résolue par `../../00-governance/adr/ADR-0008-customers-mssp-delivery-deployment-and-cross-tenant-architecture.md`. Cette définition documentaire ne revendique aucune implémentation runtime.

## Modes de déploiement

- **Internal** : reporting et suivi de service sans concept Customer obligatoire ;
- **Enterprise multi-tenant** : Tenants indépendants et contexte Customer/engagement externe lorsque configuré ;
- **MSP/MSSP** : overview read-only sur un Authorized Tenant Set Security, puis sélection explicite d'un Tenant pour toute action tenant-local.

Les capacités MSSP sont deployment-dependent.

## Customer semantics

Customer reste une projection `external/deployment/customer/contract`, jamais un objet CMDR canonique et jamais un alias de Tenant. Command ne possède ni Customer lifecycle, ni contrat, ni billing, ni CRM.

## MSSP boundaries

Le MVP :
- agrège en lecture seule les Tenants explicitement autorisés ;
- conserve l'identité Tenant de chaque objet ;
- permet le context switching explicite ;
- interdit cross-tenant mutation, administration, response, delegated administration et automatic export widening.

Aucun `ManagedTenant`, `TenantGroup`, `Portfolio`, `CustomerTenant` ou équivalent canonique n'est créé.

## Search / Reporting / Export

Initialement : Search, Report et Export sont single-selected-Tenant. Shared conserve leur ownership. Multi-tenant Search/Reporting/Export sont différés.

## Response

Une réponse MSSP suit : sélection Tenant → Security re-evaluation → Govern Decision Authority dans ce Tenant. Aucun droit read-only agrégé ne confère une autorité de réponse.

## Réutilisable sans Customer

Reporting Engine, métriques sourcées, export single-Tenant, Tasks de suivi, Services et outcomes vérifiés restent utilisables dans les déploiements Internal.

## Capability

`CAP-CMD-401 — Customers and Delivery Context`, source `capabilities/customers-and-delivery-context.md`, conserve le même ID et le même owner.

## Explicitement hors scope

- canonical Customer lifecycle/admin ;
- Tenant hierarchy ;
- cross-tenant mutation/admin/response ;
- delegated administration ;
- initial multi-tenant Search/Report/Export ;
- CRM ;
- billing ;
- customer portal ;
- contract mutation.

## OPEN préservées

- `OPEN-013` : default governance des mutations Class 2 ;
- `OPEN-019` : dissemination/releasability/client or external sharing.

## Risques

Fuite inter-tenant, confusion Tenant/Customer, mélange SLA opérationnel/contractuel, duplication du Reporting Engine, export trop large ou autorité MSSP implicite.

## Critères

**Given** un déploiement Internal sans Customer, **When** Command demande un report, **Then** Reporting Engine reste utilisable avec un Tenant sélectionné et aucun Customer canonique n'est requis.

**Given** un MSSP autorisé sur plusieurs Tenants, **When** l'overview s'ouvre, **Then** il est read-only, chaque objet garde son Tenant et toute mutation exige un Tenant unique.
