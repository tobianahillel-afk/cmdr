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
# Object consumption map — Investigate through Phase 4B.2B.3B.1

| Objet ou concept | Owner actuel | Usage local | Opérations locales | Lacune | Phase propriétaire |
|---|---|---|---|---|---|
| Case / Hypothesis / Artifact / Evidence / Finding | Investigate | context, source, reasoning and handoff | read/link/create under owner capability | final states/cardinality | Objects/Trust |
| Disk Image | Investigate canonical object | acquired immutable source | read/link/review; source never modified | detailed schema/states incomplete | Objects |
| Memory Image | Investigate canonical object | cross-source correlation | read/link only in Disk module | detailed schema/states incomplete | Objects |
| Collection Request / Collection Job | Investigate / concept | acquisition context and errors | read/link only | Job object absent | Objects |
| Endpoint / Endpoint Agent | shared / Endpoint Agent | acquisition source/capability | read projection | platform support OPEN-008 | Endpoint/Objects |
| Disk Forensics Session | Investigate concept | durable Disk analytical context | create/update/pause/close/reopen/supersede | object absent | Objects |
| Image Integrity Review | Investigate concept | source trust and exploitability | create/review/dispute/link | object absent | Objects/Trust |
| Partition / Volume / Filesystem Candidate | Investigate concepts | structural interpretation | propose/select/confirm/dispute | objects/support model absent | Objects |
| Filesystem Entry / File Observation | Investigate concepts | navigation, metadata, content and identity | read/annotate/link/bookmark | path is not stable identity | Objects |
| Deleted Entry Candidate / Unallocated Region | Investigate concepts | deleted/residual analysis | read/compare/annotate/select | attribution model absent | Objects |
| Journal Record | Investigate observation | change-history projection | read/correlate/annotate/link | not a certain user action | Objects |
| System Artifact Observation | Investigate concept | persistent system/configuration state | read/compare/annotate/link | persistent ≠ current | Objects |
| User Activity / Application Artifact | Investigate concepts under Privacy | persistent user/application context | read/minimize/annotate/correlate | authorship/intention unresolved | Objects/Privacy |
| Persistence Candidate | Investigate concept | startup/execution candidate | classify/dispute/link/handoff | not a confirmed Finding | Objects |
| Disk Timeline Entry | Investigate projection + Shared mechanism | temporal reconstruction | read/filter/correlate/annotate | Disk Timeline not object | Shared/Objects |
| Recovery Result | Investigate concept | carving/recovery output | create/read/export/withdraw | original identity not guaranteed | Objects |
| Restricted Content Record | Investigate/Security concept | encrypted/compressed/restricted handling | request/read/annotate; access gated | object/permission model absent | Objects/Permissions |
| Comparison Result | Investigate concept | multi-image/volume/snapshot differences | create/read/export/link | object absent | Objects |
| Derived Artifact | Investigate concept | extracted/recovered content | create/read/export/withdraw | object absent | Objects |
| Reproducibility Assessment / Provenance Record | Investigate assessment + Shared records | trace and disposition | read/create/dispute/supersede | contracts absent | Objects/Trust |
| Tool / Tool Call / Workflow / Automation Run | CMDR Studio | execution and provenance | select/invoke/read/link | final objects absent | Studio/Objects |
| Fleet / Policy / Storage / Retention / Environment | Platform Settings | administrative projections | read only | admin remains Settings | Settings/Objects |
| Background Job / Notification / Trace / Activity / Timeline / Export | Shared mechanisms | progress, audit, correlation and recovery | consume/emit semantics | contracts future | Shared/Trust |
| Action Request / Decision / Response Run / Result | Govern | real-target authority | read/link/prepare only | future contracts | Govern |
| Attachment | open | documentary content | reference only | OPEN-014 | Objects |

No schema, JSON Schema, final cardinality, object state machine, disk/filesystem format, low-level field model or atomic permission is defined.
