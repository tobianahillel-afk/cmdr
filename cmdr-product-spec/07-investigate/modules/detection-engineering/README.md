---
id: investigate-detection-engineering
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-06
source-of-truth: canonical
requirements:
  - REQ-INV-006
  - REQ-PROD-014
  - REQ-PROD-019
  - REQ-PROD-020
  - REQ-AI-002
open_decisions:
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Detection Engineering — Foundations, Authoring and Validation

## Mission
Fournir dans Investigate un espace fonctionnel pour transformer un besoin sourcé en Detection Content Draft documenté, validé, testé, rejoué historiquement, revu et cartographié, puis préparer un package vers la future revue et promotion.

## Capability range
CAP-INV-401..417 — dix-sept capabilities `defined` et `planned`.

## Ownership
Investigate possède les Project, Detection Hypothesis, Detection Content Draft, tests, Expected Outcomes, Validation/Replay/Match Review, Coverage/Gap et Review Package comme concepts fonctionnels. Command conserve runtime Detection, Signal, Alert et Incident. Settings conserve Data Sources, parsers, schemas, health, retention et environments. Endpoint Agent conserve la télémétrie et l’évaluation locale. Studio conserve Tools, Tool Calls, Workflows, Datasets/Evaluations génériques et Automation Runs. Govern conserve Decision, Approval et l’autorité future de production. Shared conserve Query, Telemetry Event, Jobs, Trace, Versioning, Search, Export, Reporting et collaboration.

## Safety and phase boundary
Aucun moteur, langage, syntaxe propriétaire, API, protocole, parser, compilateur, modèle ML, code, deployment, activation, deactivation, rollback, shadow/canary, production exception, Signal/Alert mutation ou Intelligence object n’est créé.

## Delivery
4B.3A.1 définit authoring et validation. 4B.3A.2 et 4B.3B restent non commencées.
