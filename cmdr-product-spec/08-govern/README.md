---
id: 08-govern-readme
domain: 08-govern
status: draft
owner: Govern Product Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-015
  - REQ-PROD-004
  - REQ-OBJ-005
  - REQ-OBJ-007
---
# Govern

## Mission

Évaluer les demandes d'action, appliquer policies et autorités, décider, exécuter ou autoriser, observer, vérifier et rollback.

## Possède

- Decision, Approval, Response Run, Result, policies, exceptions et audit de réponse ;
- lifecycle gouverné des Action Requests.

## Consomme

Incident, Case, Evidence, Finding, Endpoint context et Studio Human Gates.

## Exclusions

- pas de Work Queue générale ;
- pas de Case Workspace ;
- pas de création de Evidence ;
- pas de builder d'agents.

## Transitions principales

Finding/Incident → Action Request ; Action Request → Decision ; Decision → Response Run → Result ; Result → Command/Investigate.

## Place de l'IA

L'IA peut résumer et vérifier la complétude. Une Decision reste une expression d'autorité humaine ou d'une policy explicitement approuvée ; un agent ne s'auto-approuve pas.

## Sources

- [`../01-product-vision/product-boundaries.md`](../01-product-vision/product-boundaries.md)
- [`../00-governance/ownership-register.md`](../00-governance/ownership-register.md)
- `information-architecture.md`
- `product-definition.md`

## Critère d'acceptation

Un module de ce produit ne peut revendiquer un objet ou une capability exclue sans mise à jour des frontières, du registre de propriété et d'une ADR lorsqu'elle est transversale.
