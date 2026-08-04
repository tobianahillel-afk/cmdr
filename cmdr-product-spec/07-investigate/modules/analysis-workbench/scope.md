---
id: analysis-workbench-scope
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-014
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Scope

## Inclus
Artifact routing; Analysis Session; Tool selection; safe preview; metadata/format/structure/strings; static binary, script, document and archive analysis; comparison; Derived Artifact; provenance; reproducibility; Evidence/Finding handoff.

## Exclu
Toute exécution de l’Artifact, comportement runtime, désassemblage, décompilation, debugger, memory/disk/network forensics et choix d’outil final.

## Invariants
- Artifact, Evidence et Finding restent distincts.
- Un Derived Artifact ne remplace jamais sa source.
- Analysis Session, Tool Call et Automation Run restent distincts.
- Une sortie automatisée reste une proposition attribuée.
- Toute fonction essentielle possède une voie déterministe ou manuelle.
