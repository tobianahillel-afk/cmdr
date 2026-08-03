---
id: 10-platform-settings-readme
domain: 10-platform-settings
status: draft
owner: Platform Settings Product Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-017
  - REQ-AI-011
  - REQ-OBJ-008
---
# Platform Settings

## Mission

Administrer identités, tenants, environnements, sources, intégrations, providers, secrets, rétention, santé et flotte Endpoint Agent.

## Possède

- objets administratifs de plateforme ;
- Endpoint Agent Fleet et Endpoint Policies ;
- configuration des providers et intégrations.

## Consomme

Permission Model, audit, notifications, health et capabilities partagées.

## Exclusions

- pas d'investigation, forensic ou Case ;
- pas de Decision de réponse ;
- pas d'orchestration métier.

## Transitions principales

Configuration → produits consommateurs ; enrollment/policy/version → Endpoint Agent ; audit administratif → Audit.

## Place de l'IA

Settings administre les fournisseurs de modèles et leurs policies ; Studio possède les usages agentiques. L'absence de provider ne bloque pas les fonctions essentielles.

## Sources

- [`../01-product-vision/product-boundaries.md`](../01-product-vision/product-boundaries.md)
- [`../00-governance/ownership-register.md`](../00-governance/ownership-register.md)
- `information-architecture.md`
- `product-definition.md`

## Critère d'acceptation

Un module de ce produit ne peut revendiquer un objet ou une capability exclue sans mise à jour des frontières, du registre de propriété et d'une ADR lorsqu'elle est transversale.
