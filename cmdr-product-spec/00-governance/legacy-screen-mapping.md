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

## Finalité

Ce document décrit la destination fonctionnelle d'anciens écrans. Il ne valide pas les écrans actuels et ne crée aucune règle UX.

| Ancienne source | Destination fonctionnelle | Traitement |
|---|---|---|
| `10-command-center/MISSION_CONTROL.md` | `06-command/modules/mission-control/README.md` | workspace unique ; Now et Priorities deviennent des vues, Situation un mode ; les fichiers d'écran actuels restent des entrées de migration jusqu'à Phase 6 |
| `10-command-center/INCIDENT_DETAIL.md` | `06-command/modules/incidents-and-work-queue/screens/incident-detail.md` | écran pilote conservé |
| `10-command-center/WORK_QUEUE.md` | `06-command/modules/incidents-and-work-queue/README.md` et `saved-views.md` | workspace unique ; Incidents, Tasks, Unassigned, SLA Risk et My Work sont des vues enregistrées |
| `20-investigation-lab/CASE_WORKSPACE.md` | `07-investigate/modules/case-workspace/screens/case-workspace.md` | écran pilote |
| `20-investigation-lab/EVENT_SEARCH.md` | `07-investigate/modules/event-search/screens/event-search.md` | workspace de recherche structurée |
| `20-investigation-lab/REVERSE_ENGINEERING.md` | `07-investigate/modules/reverse-engineering/screens/reverse-engineering.md` | Technical Workbench |
| `30-response-governance/ACTION_CENTER.md` | futur Decision Workspace dans `08-govern/` | le fichier actuel sert d'entrée de migration, sans créer une Work Queue générale |
| `30-response-governance/RUNS_AND_ROLLBACK.md` | futur Response Run pilote dans `08-govern/` | séparer lifecycle du run et rollback |
| `40-platform/*` | `10-platform-settings/`, `12-shared-capabilities/`, `14-security-permissions-and-trust/`, `17-implementation-contracts/` | répartir selon propriété ; ne pas recréer un produit Platform parallèle |
| `50-quality/*` | `16-quality-and-validation/`, registres et documents propriétaires | les questions restent chez leur propriétaire |

## Règles

- Aucun chemin ancien n'est normatif.
- Une destination de migration peut rester `draft`.
- Un fichier existant n'est pas déclaré écran final parce qu'il possède un identifiant.
- Les consolidations Work Queue et Mission Control sont bloquées jusqu'aux phases UX et écrans.
- L'historique Git conserve la source ; l'espace actif conserve une seule destination.

## Critère d'acceptation

**Given** l'ancien écran Work Queue,  
**When** son mapping est consulté,  
**Then** il pointe vers un workspace et des Saved Views, pas vers un fichier inexistant ni six pages concurrentes.
