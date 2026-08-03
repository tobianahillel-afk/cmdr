---
id: 07-investigate-readme
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-OBJ-002
  - REQ-OBJ-004
---
# Investigate

## Mission

Rechercher, collecter, tester des hypothèses, analyser et produire des Evidence et Findings traçables.

## Possède

- Case, Hypothesis, Artifact, Evidence et Finding ;
- investigation lifecycle, Hunt, collection, workbench, Detection Engineering et Intelligence.

## Consomme

Incident, Endpoint/Fleet projections, Decisions et Results Govern, Studio automation, recherche et capabilities partagées.

## Exclusions

- pas d'autorité indépendante pour une réponse à risque élevé ;
- pas d'administration de flotte ;
- pas de seconde Work Queue générale.

## Transitions principales

Incident → Case ; Case → Evidence ; Evidence → Finding ; Finding → Action Request ; Result → validation et amélioration.

## Place de l'IA

L'IA peut proposer hypothèses, requêtes et résumés, mais ne confirme pas seule un Finding et ne contourne pas Govern. Les règles déterministes restent de première classe.

## Sources

- [`../01-product-vision/product-boundaries.md`](../01-product-vision/product-boundaries.md)
- [`../00-governance/ownership-register.md`](../00-governance/ownership-register.md)
- `information-architecture.md`
- `product-definition.md`

## Critère d'acceptation

Un module de ce produit ne peut revendiquer un objet ou une capability exclue sans mise à jour des frontières, du registre de propriété et d'une ADR lorsqu'elle est transversale.
