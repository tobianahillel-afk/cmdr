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
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Object consumption map — Investigate through Phase 4B.2B.2B

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case / Hypothesis / Artifact / Evidence / Finding | Investigate | context, source, reasoning and handoff | read/link/create under owner capability | final states/cardinality | Objects/Trust |
| Derived Artifact / Runtime Artifact | Investigate concepts | static transformation or dynamic/debug output | create/read/export/link | objects absent | Objects |
| Analysis Session / Analysis Result | Investigate concepts | static context/result | create/update/review/link | objects absent | Objects |
| Dynamic Analysis Session / Sandbox Run | Investigate concepts | dynamic context/execution | create/update/start/stop/read | objects absent | Objects |
| Reverse Analysis Session | Investigate concept | durable Reverse context | create/update/pause/close/reopen/supersede | object absent | Objects |
| Debugger Session | Investigate concept | isolated debug context | prepare/open/update/close/reopen | object absent | Objects |
| Code Location | Investigate concept | address/offset/location navigation | read/bookmark/annotate/link | final representation absent | Objects |
| Function / Symbol / Cross Reference | Investigate concepts | interpretation and navigation | read/confirm/rename/annotate/link | objects and confidence model absent | Objects |
| Control Flow Graph / Call Graph | Investigate + Shared mechanism | static graph projection | read/filter/compare/annotate | final graph record absent | Objects/Shared |
| Type Definition / Structure Definition | Investigate concepts | local analytical overlay | create/update/version/apply/revert | objects absent | Objects |
| Annotation / Rename / Bookmark | Investigate concepts | versioned analyst knowledge | create/update/restore/share | object boundary absent | Objects |
| Breakpoint | Investigate concept | session-scoped execution control | create/update/enable/disable/remove | object and state machine absent | Objects |
| Runtime State Snapshot / Thread Observation / Stack Frame | Investigate concepts | debugger state inspection | capture/read/compare/annotate | objects absent | Objects |
| Memory Region / Module Observation | Investigate concepts | debugger-scoped inspection only | read/compare/annotate/extract | not Memory Forensics | Objects / Phase 4B.2B.3 boundary |
| Exception Event / Debug Trace | Investigate + Shared trace mechanism | event analysis and context | read/group/annotate/link | event model absent | Objects/Shared |
| Patch Hypothesis | Investigate concept | isolated reversible experiment | create/review/apply-to-copy/revert/supersede | object/state absent | Objects |
| Reproducibility Assessment / Provenance Record | Investigate assessment + Shared records | linked provenance and disposition | read/create/dispute/supersede | contracts absent | Objects/Trust |
| Tool / Tool Call / Workflow / Automation Run | CMDR Studio | execution and provenance | select/invoke/read/link | final objects absent | Studio/Objects |
| Execution / Sandbox Environment | Platform Settings | authorized environment projection | read/select/request alternative | admin remains Settings | Settings/Objects |
| Background Job / Notification / Trace / Activity | Shared mechanisms | progress, events and recovery | consume/emit semantics | contracts future | Shared/Trust |
| Action Request / Decision / Response Run / Result | Govern | real-target authority and return | read/link/prepare only | future contracts | Govern |
| Attachment | open | documentary content | reference only | OPEN-014 | Objects |

No schema, JSON Schema, final cardinality, object state machine, internal address format or atomic permission is defined.
