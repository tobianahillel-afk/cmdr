---
id: analysis-workbench-shared-capabilities
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-014
---
# Shared capabilities consumed

| Shared Capability | Usage Analysis Workbench | Données locales | Source canonique |
|---|---|---|---|
| Background Jobs | queue, progress, cancellation, partial | Analysis Session/Tool context | `12-shared-capabilities/background-jobs.md` |
| Notifications | fin, erreur, action requise | deep link et statut | `12-shared-capabilities/notification-center.md` |
| Trace / Activity | attribution et reconstruction | événements métier | Shared sources |
| Inspector / Context Bar | sélection et contexte | configuration locale seulement | Design System |
| Object Linking | liens Case/Artifact/result | sémantique relationnelle | `12-shared-capabilities/object-linking-service.md` |
| Versioning | Tool/input/result versions | références | Shared/Studio |
| Export | package permission-aware | scope et classification | `12-shared-capabilities/export-engine.md` |
| Collaboration | annotations et revue | contenu métier | `12-shared-capabilities/collaboration-service.md` |
| Saved Views | configuration de vue | filtres/comparaison | `12-shared-capabilities/saved-views.md` |
| Timeline / Audit Hooks | événements ordonnés | événements analytiques | Shared/Trust |
