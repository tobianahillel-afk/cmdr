---
id: command-cross-product-links
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-PROD-008
  - REQ-PROD-013
  - REQ-JRN-001
  - REQ-JRN-003
  - REQ-JRN-007
---
# Cross-product links and transitions — Command

| Transition | Déclencheur | Owner destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Signal/Alert → Incident | promotion ou liaison après déduplication conceptuelle | Command | source IDs, entities, severity, confidence, disposition | triage/Incident |
| Incident → Case | besoin d’investigation; create/open idempotent | Investigate | tenant, env, Incident, time range, entities, Alerts | Incident Detail exact |
| Incident/Finding → Action Request | action proposée, impact et targets | Govern | Incident, Case/Finding refs, evidence summary, impact, urgency, requester, alternatives | Incident Detail |
| Action Request → Govern | submit | Govern | package complet + return origin | Decision/Run projection |
| Result → Incident | Result publié/autorisé | Command consomme; Govern reste owner | Result, Run, Decision, residual risk | Incident/Mission Control |
| Command → Studio | launch deployed workflow ou inspect Automation Run | Studio | context object, deployed version, permissions, return origin | source Command |

## Erreurs

- Case inaccessible : relation visible sans données protégées; option create selon permission ;
- plusieurs Cases : liste de relations, aucun Case “principal” implicite ;
- destination indisponible : source conservée, retry explicite ;
- permission refusée : aucun contexte sensible dans l’URL ;
- retour invalide : fallback au workspace parent avec explication.

## Invariant

Command ne crée jamais directement Decision, Response Run, Result, Automation Run ou objet Endpoint Agent.
