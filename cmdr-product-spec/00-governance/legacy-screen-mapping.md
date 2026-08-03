---
id: legacy-screen-mapping
domain: 00-governance
status: draft
owner: Product Architecture
updated: 2026-08-03
source-of-truth: migration
---
# Mapping des anciens écrans

| Ancien chemin condensé | Chemin canonique |
|---|---|
| `10-command-center/MISSION_CONTROL.md` | `06-command/modules/mission-control/screens/mission-control.md` |
| `10-command-center/INCIDENT_DETAIL.md` | `06-command/modules/incidents-and-work-queue/screens/incident-detail.md` |
| `10-command-center/WORK_QUEUE.md` | `06-command/modules/incidents-and-work-queue/screens/work-queue.md` |
| `20-investigation-lab/CASE_WORKSPACE.md` | `07-investigate/modules/case-workspace/screens/case-workspace.md` |
| `20-investigation-lab/EVENT_SEARCH.md` | `07-investigate/modules/event-search/screens/event-search.md` |
| `20-investigation-lab/REVERSE_ENGINEERING.md` | `07-investigate/modules/reverse-engineering/screens/reverse-engineering.md` |
| `30-response-governance/ACTION_CENTER.md` | `08-govern/modules/action-center/screens/action-center.md` |
| `30-response-governance/RUNS_AND_ROLLBACK.md` | `08-govern/modules/runs-and-rollback/screens/runs-and-rollback.md` |
| `40-platform/*` | `10-platform-settings/`, `12-shared-capabilities/`, `14-security-permissions-and-trust/` ou `17-implementation-contracts/` selon propriété |
| `50-quality/*` | `16-quality-and-validation/`, registres propriétaires et questions ouvertes locales |

Les anciens chemins sont supprimés de l’espace actif. Leur contenu utile est migré; Git conserve l’historique.
