---
id: investigate-analysis-workbench
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-014
  - REQ-PROD-020
  - REQ-INV-002
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Analysis Workbench

## Mission
Fournir le contexte fonctionnel dans lequel un analyste ouvre un Artifact, comprend son identité et sa provenance, choisit une analyse statique autorisée, conserve paramètres et résultats, produit des Derived Artifacts et prépare un handoff explicite vers Evidence ou Finding.

## Périmètre de Phase 4B.2B.1
Le module couvre l’intake, les Analysis Sessions, le contexte du Workbench, la preview, les métadonnées, les structures, les chaînes et contenus, l’analyse statique de binaires, scripts et documents, les archives, la comparaison, les Derived Artifacts, la provenance, la reproductibilité et le handoff.

## Hors périmètre
Dynamic Sandbox, Reverse Engineering, Debugger, Memory/Disk/Network Forensics, cloud/mobile analysis, moteurs finaux, APIs, protocoles, commandes et code produit.

## Ownership
Investigate possède Case, Artifact, Evidence, Finding et le contexte analytique. Studio possède Tool, Tool Call, Workflow et Automation Run. Settings administre providers, secrets et environnements. Shared possède Jobs, Trace, Activity, Linking, Versioning, Notifications et Export.

## Technical Workbench
Un canvas principal, un Inspector droit, deux panneaux auxiliaires maximum, six tabs visibles maximum, console basse optionnelle, explorer contextuel et Automation Tray fermé par défaut. Le module décrit les besoins, pas le wireframe final.

## Classification
Les treize capabilities sont `defined`, `draft` et `planned`. Aucun moteur, support d’implémentation ou statut natif livré n’est revendiqué.
