---
id: analysis-workbench-actions
domain: 07-investigate
status: draft
owner: Investigate Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
---
# Actions and permissions

| Besoin | Capabilities | Risque | Classe | Step-up / séparation | Owner / phase |
|---|---|---|---:|---|---|
| Artifact read/preview/raw/sensitive read | 301,304–309 | exposition de contenu | 0 | policy de sensibilité | Investigate/Permissions |
| Analysis Session create/update/close/reopen | 302 | modification de contexte | 2 | OPEN-013; reopen séparé possible | Investigate/Permissions |
| Tool use/version selection/output read | 303,312 | mauvais Tool/version | 0/1/2 | Studio permissions; visibilité obligatoire | Studio/Permissions |
| strings/embedded/archive extraction | 306,308,309 | volume et contenu sensible | 1 | limites et interruption | Investigate/Permissions |
| Derived Artifact create/read/export | 305–312 | transformation mal attribuée | 1/2 | provenance obligatoire | Investigate/Permissions |
| comparison and annotation | 310 | conclusion prématurée | 0/2 | reviewer selon policy | Investigate |
| Evidence candidate / Finding Draft | 313 | promotion prématurée | 2 | qualification séparée | Investigate |
| cross-tenant analysis | tous | fuite de données | bloqué par défaut | autorité dédiée future | Security/Phase 7 |
