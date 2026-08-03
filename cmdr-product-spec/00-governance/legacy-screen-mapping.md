---
id: legacy-screen-mapping
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: migration
requirements:
  - REQ-UX-008
  - REQ-UX-009
---

# Mapping des anciens écrans

| Ancienne source | Destination | Traitement Phase 3 |
|---|---|---|
| Mission Control variants | Mission Control workspace | migration Phase 6 ; vues/modes, pas nouvelles pages |
| Work Queue `incidents.md` | Work Queue `view=incidents` | deprecated alias |
| Work Queue `tasks.md` | `view=tasks` | deprecated alias |
| Work Queue `unassigned.md` | `view=unassigned` | deprecated alias |
| Work Queue `sla-risk.md` | `view=sla-risk` | deprecated alias |
| Work Queue `team-load.md` | Work Queue sans vue système | deprecated ; custom/shared view possible |
| Incident Detail | écran pilote Command | conservé |
| Case Workspace | écran pilote Investigate | conservé |
| Action Center | futur Decision Workspace | écran pilote conservé, renommage Phase 6 |
| Runs & Rollback | futur Response Run | écran pilote conservé |
| Endpoint Agent Fleet | Settings Shell | écran pilote conservé |
| Studio Builder | Builder Shell | écran pilote conservé |

Une destination dépréciée indique remplaçant, justification, migration et dépendants. Les aliases ne sont pas comptés comme écrans actifs. L'historique Git conserve le contenu ; une seule architecture reste normative.
