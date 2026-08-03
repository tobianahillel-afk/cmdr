---
id: 09-cmdr-studio-readme
domain: 09-cmdr-studio
status: draft
owner: CMDR Studio Product Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-016
  - REQ-AI-001
  - REQ-AI-002
  - REQ-OBJ-009
---
# CMDR Studio

## Mission

Concevoir, versionner, évaluer, déployer et superviser les automatisations déterministes et agentiques.

## Possède

- Skill, Tool, Tool Call, Automation Agent, Agent Team, Workflow, Human Gate et Automation Run ;
- Library, Builder, Assurance et Control Room.

## Consomme

Objets et contextes des produits sous permission, policies Govern, providers et secrets administrés par Settings.

## Exclusions

- ne possède pas Incident, Case, Finding, Decision ou Response Run ;
- ne remplace pas les produits opérationnels ;
- n'est pas l'interface obligatoire.

## Transitions principales

Trigger produit → Automation Run → Tool Calls → output vers produit source ; action risquée → Action Request/Govern.

## Place de l'IA

Studio possède les capacités agentiques, leur assurance, leurs versions, coûts et traces. Les workflows opérationnels restent utilisables sans IA.

## Delivery classification

Les capabilities Studio sont des cibles produit `planned` tant que leur implémentation et leurs contrats détaillés ne sont pas prouvés. Tool Call et Automation Run seront formalisés en Phase 7.

## Sources

- [`../01-product-vision/product-boundaries.md`](../01-product-vision/product-boundaries.md)
- [`../00-governance/ownership-register.md`](../00-governance/ownership-register.md)
- `information-architecture.md`
- `product-definition.md`

## Critère d'acceptation

Un module de ce produit ne peut revendiquer un objet ou une capability exclue sans mise à jour des frontières, du registre de propriété et d'une ADR lorsqu'elle est transversale.
