---
id: 07-investigate-product-definition
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-002
  - REQ-PROD-005
  - REQ-PROD-014
  - REQ-PROD-046
---

# Définition fonctionnelle du produit Investigate

## Objet

Investigate transforme des Signals et données autorisés en travail analytique reproductible, puis en Evidence et Findings revus. Il prépare le contexte nécessaire à Command, Govern, Reporting et aux capacités futures sans posséder leurs objets.

## Questions auxquelles le produit répond

- Que montrent réellement les données, de quelles sources et avec quelles limites ?
- Quelle Hypothesis est testée, et quels éléments la soutiennent ou la contredisent ?
- Quel Case organise le travail, et quels Incidents lui sont liés ?
- Quel Artifact est simplement disponible et lequel a été qualifié comme Evidence ?
- Quelle conclusion constitue un Finding revu ?
- Quel package doit être transmis à Govern sans créer de Decision ?

## Frontières

Investigate ne devient jamais une seconde Work Queue générale, une autorité autonome de réponse, un moteur d’administration Endpoint, un catalogue visible d’outils tiers, une interface centrée sur un chatbot ou un propriétaire de Decision, Response Run, Result, Workflow ou Automation Run.

## Résultat attendu de Phase 4B.1

Signals and Hunt et Cases and Evidence disposent de capabilities à 27 sections, d’actions classées, d’alternatives sans IA, de transitions et de sources de vérité. Collection and Live Response, Analysis Workbench, Detection Engineering et Intelligence restent hors périmètre.

## Invariants

1. Les sources sont visibles.
2. Les transformations sont traçables.
3. Les relations ne transfèrent ni permission ni ownership.
4. Une proposition automatisée reste une proposition.
5. Une Evidence conserve sa source et sa version.
6. Un Finding référence ses Evidence et contradictions.
7. Une Action Request demeure un objet Govern.
8. Une capability `planned` ne prouve aucune implémentation.
