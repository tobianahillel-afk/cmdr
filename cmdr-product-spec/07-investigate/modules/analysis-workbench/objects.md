---
id: analysis-workbench-objects
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-OBJ-002
  - REQ-OBJ-003
  - REQ-OBJ-004
  - REQ-PROD-061
open_decisions:
  - OPEN-014
  - OPEN-015
---
# Objects and concepts

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case | Investigate | contexte durable | read/link | cardinalités finales | Objets |
| Hypothesis | Investigate | question testée | read/link/propose | confiance finale | Objets |
| Artifact | Investigate | source analysée | read/preview/link/version | stockage/dedup | Objets/Technique |
| Attachment | ouvert | contenu documentaire éventuel | read/reference | OPEN-014 | Objets |
| Evidence | Investigate | qualification aval | candidate handoff only | trust final | Objets/Trust |
| Finding | Investigate | conclusion aval | draft handoff only | validation finale | Objets/Permissions |
| Tool / Tool Call | Studio | exécution et provenance | select/invoke/read | objets finaux absents | Studio/Objets |
| Workflow / Automation Run | Studio | assistance optionnelle | invoke/read/link | OPEN-015 | Studio/Objets |
| Background Job | Shared | queue/progress/cancel | consume | contrat technique | Shared/Technique |
| Analysis Session | concept Investigate | contexte d’analyse | create/update/close/reopen | objet absent | Objets |
| Analysis Result | concept Investigate | résultat structuré | receive/review/link | objet absent | Objets |
| Derived Artifact | concept Investigate | transformation sourcée | create/read/export/link | objet absent | Objets |
| Annotation | concept Investigate/Shared | interprétation humaine | create/supersede | ownership final | Objets |
| Indicator Candidate | concept Investigate | valeur extraite non confirmée | propose/link | bridge Intelligence | 4B.3/Objets |
| Report | Shared Reporting Engine | sortie canonique future | prepare only | objet Report absent | Shared/Objets |
| Provenance Record | Shared mechanism, Investigate semantics | reconstruction | emit/review | objet absent | Shared/Trust |
