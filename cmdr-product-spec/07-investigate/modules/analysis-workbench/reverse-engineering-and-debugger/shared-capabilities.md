---
id: reverse-debugger-shared-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-05
source-of-truth: canonical
---
# Shared capability consumption

| Shared Capability | Usage Reverse/Debugger | Local data | Source canonique |
|---|---|---|---|
| Background Jobs | queued Tool work, progress, cancel and partial status | analytical request/result references | `12-shared-capabilities/background-jobs.md` |
| Notifications | completion, failure, revocation and review requests | notification references only | `12-shared-capabilities/notification-center.md` |
| Trace / Activity Stream | actions, Tool Calls, controls, errors and dispositions | linked analytical context | `03-design-system/components/trace.md`; `components/activity-stream.md` |
| Inspector / Context Bar | selected function, location, frame, event or object context | selection and return origin | `03-design-system/components/inspector.md`; `context-bar.md` |
| Timeline Engine | ordered debugger events and gaps | event references and annotations | `12-shared-capabilities/timeline-engine.md` |
| Graph Engine | CFG, call graph and structured alternatives | node/edge selections and notes | `12-shared-capabilities/graph-engine.md` |
| Object Linking | Case, Artifact, Evidence, Finding and session links | link justification | `12-shared-capabilities/object-linking-service.md` |
| Export Engine | policy-aware result and provenance export | export request/reference | `12-shared-capabilities/export-engine.md` |
| Versioning | sessions, annotations, types, hypotheses and supersession | version relations | `05-domain-model/versioning-and-supersession.md` |
| Collaboration | contributors, comments and review | local collaboration refs | `12-shared-capabilities/collaboration-service.md` |
| Audit Hooks | actor, action, result and correlation | audit reference | `14-security-permissions-and-trust/audit-and-immutability.md` |
| permission-aware search | functions, symbols, annotations and sessions | permitted result projections | `12-shared-capabilities/global-search.md` |
| Recovery | crash, timeout, stale state and unsaved work recovery | last known state reference | `03-design-system/patterns/recovery.md` |

No Shared Capability is redefined locally.
