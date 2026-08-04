---
id: dynamic-sandbox-dependencies
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
| CAP-INV-301..313 | source statique, Derived Artifact, provenance et handoff | Investigate | active |
| CAP-INV-105/107/108/109 | Artifact, Evidence et Finding lifecycles | Investigate | active |
| CAP-INV-213/214 | custody et provenance de collecte | Investigate | active |
| Sandbox Environment | disponibilité, isolation déclarée, health, reset | Platform Settings | projection only |
| Tool/Tool Call/Workflow/Automation Run | exécution et attribution | Studio | partial / Phase 7 |
| Jobs/Trace/Activity/Linking/Export/Notifications | mécanismes partagés | Shared | consumed |
| Govern | cibles réelles et actions risquées hors sandbox | Govern | boundary |
| OPEN-005 | moteurs futurs | Investigate | open |
| OPEN-013 | classe 2 | Security | open |
| OPEN-014 | Artifact/Attachment | Investigate | open |
| OPEN-015 | Automation/Response Run | Studio/Govern | open |
| 4B.2B.2B / 4B.2B.3 | Reverse/Debugger/Forensics | Investigate | non commencées |
