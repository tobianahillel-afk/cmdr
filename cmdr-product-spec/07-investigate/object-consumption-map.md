---
id: investigate-object-consumption-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-014
  - REQ-PROD-061
  - REQ-PROD-062
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Object consumption map — Investigate through Phase 4B.2B.1

| Objet ou concept | Owner actuel | Usage Investigate | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case | Investigate | workspace/static context | create/update/link | final states/cardinality | Objets |
| Hypothesis | Investigate | testable question | create/review/link | confidence | Objets |
| Artifact | Investigate | analyzable source | register/version/link/derive | storage/dedup | Objets/Technique |
| Derived Artifact | concept Investigate | static transformation | create/read/export/link | object absent | Objets |
| Attachment | open | documentary file | reference | OPEN-014 | Objets |
| Evidence / Finding | Investigate | explicit qualification/conclusion | candidate/draft handoff | trust/review | Objets/Trust |
| Tool / Tool Call / Workflow / Automation Run | Studio | execution/provenance | select/invoke/read/link | final objects absent | Studio/Objets |
| Analysis Session / Analysis Result | concepts Investigate | static context/output | create/update/review/link | objects absent | Objets |
| Background Job / Provenance Record | Shared mechanisms | progress/reconstruction | consume/emit semantics | technical/object contracts | Shared/Trust |
| Report | Shared Reporting Engine | prepare only | no local lifecycle | object absent | Shared/Objets |

No schema, cardinality, final state machine or atomic permission is defined.
