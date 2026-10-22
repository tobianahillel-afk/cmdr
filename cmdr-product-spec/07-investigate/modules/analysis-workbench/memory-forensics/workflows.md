---
id: investigate-memory-forensics-workflows
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-INV-001
  - REQ-PROD-014
  - REQ-PROD-020
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Workflows

| Source | Déclencheur | Destination | Ownership | Contexte et retour |
|---|---|---|---|---|
| Case / Artifact | ouvrir Memory Forensics | CAP-INV-347 | Investigate | Case, Endpoint, Memory Image, objectif, Hypothesis, return origin |
| Collection Job | résultat disponible | Memory Image / CAP-INV-349 | Collection puis Investigate | Request, Job, acquisition, errors, custody, Artifact relation |
| Integrity Review | image exploitable ou partial acceptée | CAP-INV-348 | Investigate | image, limitations, owner, objective |
| Session | identifier profil | CAP-INV-350 | Investigate consuming Studio/Settings | image, candidates, Tools, prior results |
| Profile | reconstruire | CAP-INV-351..355 | Investigate | profile, confidence, Tool/version, limitations |
| Region/module | examiner anomaly ou extraire | CAP-INV-357/360 | Investigate | source, relations, restrictions, provenance |
| Observation | construire timeline | CAP-INV-359 | Investigate + Shared Timeline | timestamps, quality, sources |
| Derived Artifact | analyser | Static / Reverse | Investigate | parent image, source context, Case, provenance |
| Memory result | préparer Evidence/Finding | CAP-INV-362 | Investigate owner capabilities | uncertainty, contradictions, source links |
| Memory network result | handoff futur | Phase 4B.2B.3B | future owner | candidates, process links, timestamps, limits |
| Knowledge package | handoff futur | Phase 4B.3 | future owner | behavior, conditions, uncertainty, sources |

Toute erreur conserve la source, les résultats valides et le return origin. Toute action réelle sur l’Endpoint est bloquée ou routée vers Collection/Live Response et Govern.
