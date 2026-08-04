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
# Object consumption map — Investigate through Phase 4B.2B.2A

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case / Hypothesis / Artifact / Evidence / Finding | Investigate | context, source, reasoning and handoff | read/link/create under owner capability | final states/cardinality | Objets/Trust |
| Derived Artifact | concept Investigate | static transformation | create/read/export/link | object absent | Objets |
| Runtime Artifact | concept Investigate | dynamic output | capture/read/export/link | object absent | Objets |
| Analysis Session / Analysis Result | concepts Investigate | static context/result | create/update/review/link | objects absent | Objets |
| Dynamic Analysis Session / Sandbox Run | concepts Investigate | dynamic context/execution | create/update/start/stop/read | objects absent | Objets |
| Behavioral / Process / Network / System Observation | concepts Investigate | sourced behavior | read/annotate/link/compare | records absent | Objets |
| Interaction Profile | Studio/Settings definition; Investigate use | selected scenario | read/select/version | ownership detail future | Studio/Settings |
| Sandbox Environment | Platform Settings | authorized environment projection | read/select/request alternative | admin remains Settings | Settings/Objets |
| Tool / Tool Call / Workflow / Automation Run | Studio | execution and provenance | select/invoke/read/link | final objects absent | Studio/Objets |
| Provenance Record / Background Job | Shared mechanisms | trace and progress | consume/emit semantics | contracts future | Shared/Trust |
| Action Request / Decision / Response Run / Result | Govern | real-target authority and return | read/link/prepare only | future contracts | Govern |
| Attachment | open | documentary content | reference only | OPEN-014 | Objets |

No schema, JSON Schema, final cardinality, object state machine or atomic permission is defined.
