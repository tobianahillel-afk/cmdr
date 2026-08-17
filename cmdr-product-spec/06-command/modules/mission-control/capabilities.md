---
id: command-module-mission-control-capabilities
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-003
  - REQ-PROD-005
  - REQ-PROD-006
  - REQ-PROD-008
  - REQ-PROD-009
  - REQ-PROD-010
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-UX-007
---

# Capability map — Mission Control

| ID | Capability | Primary role | Objects | Representative actions | Status | Mode | OPEN |
|---|---|---|---|---|---|---|---|
| CAP-CMD-001 | Situation Overview | Incident Commander | Incident, Task, Decision / Response Run / Result, Service | Consulter la situation (C0), Ouvrir l’objet source (C0), Mettre à jour la prochaine action (C2) | defined | planned | — |
| CAP-CMD-002 | Priority Management | Incident Commander | Incident / Task, Service, SLA context, Audit trail | Inspecter les facteurs (C0), Modifier la priorité (C2), Accepter une recommandation (C2) | defined | planned | OPEN-013 |
| CAP-CMD-003 | Situation Timeline | Incident Commander | Timeline Entry, Incident, Decision / Run / Result | Filtrer la timeline (C0), Inspecter la provenance (C0), Exporter la sélection (C0) | defined | planned | — |
| CAP-CMD-004 | Handover | Incident Commander | Incident / Task, Decision / Response Run / Result, Case / Finding, Handover record, Incident / Task ownership | Créer ou modifier le brouillon (C2), Marquer prêt (C2), Envoyer (C2) | defined | planned | — |
| CAP-CMD-005 | Operational Blockers | Incident Commander | Incident / Task, Decision / Case / external dependency, Task | Déclarer bloqué (C2), Créer une Task de déblocage (C2), Escalader (C2) | defined | planned | OPEN-013 |
| CAP-CMD-006 | Recent Results and Outcomes | Incident Commander | Result / Response Run / Decision, Incident, Timeline | Consulter le Result (C0), Associer à un Incident (C2), Mettre à jour la prochaine action (C2) | defined | planned | OPEN-013 |

## Règles

- un Capability ID ne change pas avec son chemin ;
- la source individuelle possède les entrées, sorties et critères ;
- les agrégats de ce fichier ne remplacent pas les documents canoniques ;
- un delivery mode `planned` ne prouve aucune disponibilité runtime.
