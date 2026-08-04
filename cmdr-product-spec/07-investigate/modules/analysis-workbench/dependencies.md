---
id: analysis-workbench-dependencies
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-014
open_decisions:
  - OPEN-005
  - OPEN-013
  - OPEN-014
  - OPEN-015
---
# Dependencies

| Dépendance | Usage | Owner | État |
|---|---|---|---|
| CAP-INV-105 Artifact Management | source, versions et lineage | Investigate | active |
| CAP-INV-107/108 | Evidence candidate et qualification | Investigate | active |
| CAP-INV-109 | Finding Draft et lifecycle | Investigate | active |
| Studio Tool/Tool Call/Automation Run | exécution et provenance | Studio | partial / Phase 7 |
| Settings environments/providers/secrets | disponibilité administrative | Settings | projection seulement |
| Shared Jobs/Trace/Activity/Linking/Versioning/Export | mécanismes transversaux | Shared | consommés |
| OPEN-005 | futurs moteurs | Investigate | open |
| OPEN-013 | classe 2 | Security | open |
| OPEN-014 | Artifact/Attachment | Investigate | open |
| OPEN-015 | Automation Run bridge | Studio/Govern | open |
| Phases 4B.2B.2 et 4B.2B.3 | destinations futures | Investigate | non commencées |
