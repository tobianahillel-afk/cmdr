---
id: 07-investigate-readme
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-002
  - REQ-PROD-005
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-PROD-046
---

# Investigate

## Mission

Investigate recherche, trie, comprend, teste des Hypotheses, organise des Cases, qualifie des Artifacts et Evidence, produit des Findings traçables et prépare des demandes d’action. Il soutient Command, Govern, Reporting et les futurs domaines Detection Engineering et Intelligence sans absorber leur ownership.

## Modules canoniques

1. **Signals and Hunt** — Signal Triage, Event Search, Query assets, Hunt et provenance.
2. **Cases and Evidence** — Case lifecycle, Hypotheses, Entities, Artifacts, Evidence, Findings, collaboration, timeline, replay, Action Request et reporting.
3. **Collection and Live Response** — Phase 4B.2, non traitée.
4. **Analysis Workbench** — Phase 4B.2, non traitée.
5. **Detection Engineering** — Phase 4B.3, non traitée.
6. **Intelligence** — Phase 4B.3, non traitée.

Les anciens regroupements `intake-and-triage`, `event-search`, `case-workspace`, `case-timeline`, `entity-graph`, `evidence`, `hypotheses-and-findings` et `investigation-report` restent des chemins de migration ou d’écrans. Leurs règles fonctionnelles 4B.1 sont consolidées dans les deux modules canoniques.

## Ownership

Investigate possède Case, Hypothesis, Artifact, Evidence, Finding, les relations et événements analytiques Case-scoped, ainsi que les Saved Searches et Query Assets de domaine.

Investigate consomme sans redéfinir Signal, Alert, Detection, Incident et Task — Command ; Telemetry Event, Query, Search Job, Entity, Timeline, Saved Views et Reporting Engine — Shared ; Action Request, Decision, Approval, Response Run et Result — Govern ; Tool, Workflow et Automation Run — CMDR Studio ; Endpoint Agent Fleet — Platform Settings.

`Action Request` possède un lifecycle Govern ; Investigate est producteur. `Attachment` n’est pas un objet canonique actuel et reste sous `OPEN-014`.

## Principes

- donnée brute, Event, Artifact, Evidence, Hypothesis, Finding, Decision et Result sont distincts ;
- une donnée ou un Artifact ne devient jamais automatiquement Evidence ;
- une Hypothesis ne devient jamais automatiquement Finding ;
- un Finding proposé par IA ne devient jamais confirmé sans revue humaine autorisée ;
- toute activité essentielle dispose d’une voie manuelle ou déterministe ;
- une transition conserve contexte et ownership ;
- aucune capability `planned` n’est présentée comme livrée.

## Phase 4B.1

Cette sous-phase définit 22 Capability IDs : `CAP-INV-001` à `CAP-INV-008`, puis `CAP-INV-101` à `CAP-INV-114`. Les 22 delivery modes restent `planned`. `CAP-INV-106` reste `proposed` sous `OPEN-014`; les 21 autres sont fonctionnellement `defined`.

## Limites

La syntaxe de requête, le moteur de recherche, l’index, les stockages, les algorithmes d’intégrité, les moteurs forensics, les protocoles Endpoint, les APIs, les machines d’état finales, les permissions atomiques, les écrans détaillés et le logiciel restent hors Phase 4B.1.
