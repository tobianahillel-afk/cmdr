---
id: command-module-customers-and-delivery
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: proposal
requirements:
  - REQ-PROD-053
  - REQ-PROD-013
  - REQ-PROD-019
  - REQ-PROD-033
open_decisions:
  - OPEN-006
  - OPEN-013
---
# Customers and Delivery — proposal

## État
`proposed`, `deployment-dependent`, `delivery_mode: planned`. Le module n’est pas universellement actif et `OPEN-006` reste ouverte.

## Scénarios
- déploiement interne : reporting et suivi de service sans concept Customer ;
- entreprise multi-tenant : contexte d’entité/engagement selon configuration ;
- MSP/MSSP : portefeuille, engagements, SLA contractuels et audience, sous décision explicite.

## Réutilisable sans modèle client
Reporting Engine, métriques sourcées, export, Tasks de suivi, Services et outcomes vérifiés.

## Dépendant du modèle économique
Multi-client, portail, branding client, billing, obligations contractuelles, contractual SLA et règles de delivery.

## Capability
`CAP-CMD-401 — Customers and Delivery Context`, source `capabilities/customers-and-delivery-context.md`.

## Risques
Fuite inter-tenant, mélange SLA opérationnel/contractuel, duplication du Reporting Engine, promesse universelle ou modèle client implicite.

## Critère
**Given** un déploiement interne sans Customer, **When** Command produit un report, **Then** Reporting Engine reste utilisable et le module client demeure désactivé/proposé.
