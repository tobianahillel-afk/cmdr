---
id: dynamic-sandbox-shared-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-014
---
# Shared capabilities consumed

| Shared Capability | Usage Dynamic Sandbox | Données locales | Source canonique |
|---|---|---|---|
| Background Jobs | queue/progress/cancel/partial | Run context | `12-shared-capabilities/background-jobs.md` |
| Notifications | fin, erreur, unsafe environment | status/deep link | `12-shared-capabilities/notification-center.md` |
| Trace / Activity / Audit Hooks | attribution et reconstruction | business events | Shared/Trust |
| Inspector / Context Bar | selection and preserved context | local configuration | Design System |
| Timeline | observed/inferred/annotation ordering | event semantics | `12-shared-capabilities/timeline-engine.md` |
| Object Linking | Case/Artifact/Run/Observation links | relation semantics | `12-shared-capabilities/object-linking-service.md` |
| Versioning | environment/profile/Tool/input versions | references | Shared/Studio |
| Export | permission-aware results package | scope/classification | `12-shared-capabilities/export-engine.md` |
| Collaboration | annotations/review | business content | `12-shared-capabilities/collaboration-service.md` |
