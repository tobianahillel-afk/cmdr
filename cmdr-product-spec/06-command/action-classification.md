---
id: command-action-classification
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-004
  - REQ-PROD-013
  - REQ-SEC-001
  - REQ-SEC-002
open_decisions:
  - OPEN-013
---
# Action classification — Command

| Classe | Command peut | Exemples | Govern |
|---:|---|---|---|
| 0 | observer, filtrer, rechercher, naviguer, inspecter | lire Incident/Case/Decision/Result, changer de vue | non |
| 1 | demander une collecte | préparer contexte vers Investigate | destination Investigate; pas d’exécution Command |
| 2 | modifier réversiblement un objet Command | assignment, priority, Task, relation, handover, SLA pause autorisée | OPEN-013 reste ouverte |
| 3 | demander ou suivre containment | préparer Action Request, afficher Run | Govern requis par défaut |
| 4 | fournir contexte pour destructif/irréversible | package action, impact, alternatives | Govern obligatoire |

## Rules

- action class describes effect, not button appearance;
- a Command request for a class-3 effect is a class-2 preparation and transition;
- bulk supports class 0 and restricted class 2 only;
- no class 3 or 4 execution endpoint is defined;
- each capability action table is normative for Phase 4A;
- OPEN-013 is not closed.
