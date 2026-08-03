---
id: 06-command-readme
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-013
  - REQ-PROD-010
  - REQ-OBJ-001
---
# Command

## Mission

Prioriser, coordonner, distribuer et superviser le travail de cybersécurité en maintenant la situation opérationnelle et l'impact métier.

## Possède

- Incident et coordination opérationnelle ;
- Work Queue et Tasks opérationnelles ;
- priorité, SLA, ownership, situation, handover, readiness et customer delivery lorsque applicable.

## Consomme

Signals, Alerts, projections de Case/Evidence/Finding, Decisions, Response Runs, Results, Reporting Engine et Saved Views génériques.

## Exclusions

- pas de forensic, reverse engineering ou terminal endpoint ;
- pas de policy engine détaillé ;
- pas de builder d'agents ;
- pas de cycle de vie concurrent pour Case ou Evidence.

## Transitions principales

Signal/Alert → Incident ; Incident → Case ; Incident/context → Govern ; Result → Incident et situation.

## Place de l'IA

L'IA peut résumer ou proposer une priorité. Command reste propriétaire de l'Incident, de l'ownership et du handover. Les fonctions essentielles restent disponibles sans modèle.

## Sources

- [`../01-product-vision/product-boundaries.md`](../01-product-vision/product-boundaries.md)
- [`../00-governance/ownership-register.md`](../00-governance/ownership-register.md)
- `information-architecture.md`
- `product-definition.md`

## Critère d'acceptation

Un module de ce produit ne peut revendiquer un objet ou une capability exclue sans mise à jour des frontières, du registre de propriété et d'une ADR lorsqu'elle est transversale.
