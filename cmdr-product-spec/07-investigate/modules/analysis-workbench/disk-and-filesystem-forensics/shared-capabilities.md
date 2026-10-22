---
id: investigate-disk-filesystem-shared-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Shared Capabilities consumption

| Shared Capability | Disk/Filesystem use | Local data | Canonical source |
|---|---|---|---|
| Background Jobs | bounded analyses, comparison and recovery progress/cancel/partial | scope and business status | `12-shared-capabilities/background-jobs.md` |
| Notifications | completion, failure, blocked/restricted states | safe summary and correlation ID | `12-shared-capabilities/notification-center.md` |
| Trace / Activity | Tool Calls, selections, annotations, access and handoffs | business semantics | Design System Trace / Shared services |
| Timeline | ordering and correlation | Disk timestamp quality and source | `12-shared-capabilities/timeline-engine.md` |
| Object Linking | Case, Artifact, file observations, Entity and candidates | typed relations | `12-shared-capabilities/object-linking-service.md` |
| Versioning | sessions, interpretations, annotations and comparisons | local versions/supersession | Shared/versioning sources |
| Export | authorized Artifact/result/provenance export | restrictions and redaction | `12-shared-capabilities/export-engine.md` |
| Reporting | local analysis package consumption | selected findings/candidates | `12-shared-capabilities/reporting-engine.md` |
| Collaboration | contributors, comments and review | session membership | `12-shared-capabilities/collaboration-service.md` |
| Inspector / Context Bar | contextual inspection and preserved source | selected object projection | Design System canonical components |
| Search | permission-aware path/metadata/content search | query/filter context | `12-shared-capabilities/global-search.md` |
| Audit Hooks | sensitive access and disposition | actor/reason/result | Trust/Audit sources |
| Recovery | resume session and bounded jobs | checkpoints and errors | Design System/Shared recovery |
