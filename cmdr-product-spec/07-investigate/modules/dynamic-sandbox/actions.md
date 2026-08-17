---
id: dynamic-sandbox-actions
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
| Dynamic Analysis read/create/update/close | 314,316 | contexte incorrect | 0/2 | reopen séparé possible | Investigate/Permissions |
| Sandbox Environment view/select | 315 | mauvais environnement | 0/2 | admin reste Settings | Settings |
| Run prepare/start/stop/cancel/retry | 317,326 | exécution/ressource | 1/2 | policy et OPEN-013 | Investigate/Settings |
| network/interaction profile select | 315,323 | réseau ou interaction non voulus | 2 | confirmation explicite | Settings/Studio |
| sensitive Artifact execute-in-sandbox | 314,317 | contenu dangereux | 1/2 | permission et environnement autorisé | Security |
| behavioral/process/network results read | 318–322 | contenu sensible | 0 | policy/redaction | Investigate |
| Runtime Artifact create/read/export | 324 | propagation/fuite | 1/2 | restrictions et provenance | Investigate |
| comparison/reproducibility | 325,327 | conclusion trompeuse | 0/1/2 | reviewer selon policy | Investigate |
| Evidence candidate / Finding Draft | 328 | promotion prématurée | 2 | qualification séparée | Investigate |
| automated analysis request | 314–317 | lancement implicite | 2 | proposition visible; humain/gate | Studio/Investigate |
| restricted/cross-tenant environment | tous | fuite/autorité | bloqué par défaut | autorité future | Security/Phase 7 |
