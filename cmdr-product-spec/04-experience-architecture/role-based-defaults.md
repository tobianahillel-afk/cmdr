---
id: role-based-defaults
domain: 04-experience-architecture
status: draft
owner: UX Architecture Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-PROD-007
  - REQ-PROD-021
  - REQ-PROD-036
  - REQ-PROD-057
---

# Valeurs par défaut selon le rôle

Les rôles modifient la page d'arrivée, la densité proposée, les modules favoris et l'ordre initial des sections de l'Inspector ; ils ne créent ni objet, ni permission, ni statut différent.

| Activité | Défaut | Exemples de rôles |
|---|---|---|
| triage et Queue | compact, My Work | L1, Incident Commander |
| Case | standard, Evidence/Timeline | L2, Senior Analyst |
| Technical Workbench | compact, Inspector technique | DFIR, Malware, Reverse |
| Decision | standard, impact/autorité | Approver, Business Owner |
| Run | standard, progression/rollback | Response Operator |
| Builder | standard, validation/version | Automation Designer |
| Settings | standard, formulaires | Platform Administrator |

Un utilisateur peut surcharger densité, favoris et ordre non critique. Permission, provenance, autorité, Human Gates et audit ne sont jamais personnalisables. `OPEN-010` reste ouverte jusqu'aux études utilisateurs.
