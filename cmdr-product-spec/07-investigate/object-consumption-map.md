---
id: investigate-object-consumption-map
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
requirements:
  - REQ-PROD-006
  - REQ-PROD-014
  - REQ-PROD-061
  - REQ-PROD-062
open_decisions:
  - OPEN-005
  - OPEN-008
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Object consumption map — Investigate through Phase 4B.2B.3A

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case / Hypothesis / Artifact / Evidence / Finding | Investigate | context, source, reasoning and handoff | read/link/create under owner capability | final states/cardinality | Objects/Trust |
| Memory Image | Investigate canonical object | acquired source for analysis | read/link/review status; source immutable | detailed schema/states incomplete | Objects |
| Collection Request / Collection Job | Investigate / concept | acquisition context | read/link only | Job object absent | Objects |
| Endpoint / Endpoint Agent | shared / Endpoint Agent | acquisition source and capability | read projection | platform support OPEN-008 | Endpoint/Objects |
| Memory Forensics Session | Investigate concept | durable forensic context | create/update/pause/close/reopen/supersede | object absent | Objects |
| Platform Candidate / Analysis Profile | Investigate concepts using Settings/Studio | analysis interpretation | propose/select/confirm/dispute | objects absent | Objects |
| Process / Thread Observation | Investigate concepts | reconstructed system state | create/read/annotate/link | objects absent | Objects |
| Memory Region / Mapping | Investigate concepts | memory layout observations | read/compare/annotate/extract relation | objects absent | Objects |
| Module / Driver Observation | Investigate concepts | loaded component observations | read/compare/annotate/extract relation | objects absent | Objects |
| Handle / System Object Observation / IPC Relation | Investigate concepts | reconstructed technical relations | read/filter/annotate/link | distinct from canonical CMDR objects | Objects |
| Network State Observation | Investigate concept | memory-resident connection candidate | read/annotate/link/handoff | not full Network Forensics | Objects / 4B.2B.3B |
| Sensitive Material Candidate | Investigate concept under Security policy | masked exposure assessment | detect/mask/review/restrict; reveal/copy/export gated | object/permission model absent | Objects/Permissions |
| Memory Anomaly / Kernel State Observation | Investigate concepts | candidate anomaly and kernel interpretation | create/classify/dispute/link | objects absent | Objects |
| Memory Timeline / Timeline Entry | Investigate projection + Shared object | temporal correlation | read/filter/annotate/link | Memory Timeline not object | Shared/Objects |
| Derived Artifact / Extraction Result | Investigate concept | extracted content and result | create/read/export/withdraw | Derived object absent | Objects |
| Reproducibility Assessment / Provenance Record | Investigate assessment + Shared records | trace and disposition | read/create/dispute/supersede | contracts absent | Objects/Trust |
| Analysis / Reverse / Debugger Sessions and Runtime Snapshot | Investigate concepts | cross-analysis context | read/link only where relevant | objects absent | Objects |
| Tool / Tool Call / Workflow / Automation Run | CMDR Studio | execution and provenance | select/invoke/read/link | final objects absent | Studio/Objects |
| Fleet / Policy / Storage / Retention / Environment | Platform Settings | administrative projections | read only | admin remains Settings | Settings/Objects |
| Background Job / Notification / Trace / Activity / Export | Shared mechanisms | progress, audit and recovery | consume/emit semantics | contracts future | Shared/Trust |
| Action Request / Decision / Response Run / Result | Govern | real-target authority | read/link/prepare only | future contracts | Govern |
| Attachment | open | documentary content | reference only | OPEN-014 | Objects |

No schema, JSON Schema, final cardinality, object state machine, memory format, low-level field model or atomic permission is defined.
