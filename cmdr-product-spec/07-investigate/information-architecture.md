---
id: 07-investigate-information-architecture
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-008
  - REQ-PROD-014
  - REQ-UX-001
  - REQ-UX-006
---

# Architecture d’information Investigate

## Organisation fonctionnelle

| Module | Objectif principal | Phase |
|---|---|---|
| Signals and Hunt | partir d’un Signal ou d’une question, rechercher, inspecter, pivoter et organiser un Hunt | 4B.1 |
| Cases and Evidence | organiser le Case, tester des Hypotheses, qualifier Evidence et produire Findings | 4B.1 |
| Collection and Live Response | demander et suivre acquisitions, sessions et collecte | 4B.2 |
| Analysis Workbench | analyser Artifacts dans des workbenches techniques | 4B.2 |
| Detection Engineering | construire, tester, déployer et mesurer des détections | 4B.3 |
| Intelligence | gérer indicators, campaigns, actors, watchlists et enrichissements | 4B.3 |

## Workspaces

- **Signal Triage** : qualification et pivot.
- **Event Search** : recherche et inspection.
- **Hunt** : activité exploratoire organisée, sans statut d’objet final décidé.
- **Case Queue** : liste spécialisée des Cases, pas Work Queue.
- **Case Workspace** : travail durable sur Case, Hypotheses, Evidence, Findings et collaboration.
- Evidence, Reasoning, Timeline et Reporting sont des sections, vues ou workspaces spécialisés selon l’objectif ; leur découpage d’écran final reste Phase 6.

## Règles page, vue, mode, filtre

Une page ou un workspace correspond à un objectif réellement distinct. Raw/rendered, graph/table et compare sont des modes. Les filtres et Saved Views ne créent pas de pages. Event Inspector peut être Inspector, panneau ou workspace selon profondeur, sans décision d’écran finale. Chaque capability n’implique pas un écran.

## Ancienne architecture

Les anciens modules basés sur outils ou écrans sont conservés comme sources de besoins ou pointeurs de migration. Static Analysis, Sandbox, Reverse, Debugger, Memory et Disk restent inchangés et appartiennent à 4B.2.

## Context preservation

Tenant, environnement, Incident, Case, Signal, période, Query/version, Hunt, Artifact, sélection, filtres, tabs, scroll et return origin sont transmis seulement lorsqu’ils sont pertinents et autorisés.
