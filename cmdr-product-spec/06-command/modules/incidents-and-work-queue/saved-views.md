---
id: command-work-queue-saved-views
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-03
source-of-truth: canonical
requirements:
  - REQ-OBJ-012
  - REQ-UX-008
  - REQ-UX-009
---

# Saved Views de la Work Queue

## Source fonctionnelle générique

`../../../12-shared-capabilities/saved-views.md` possède stockage, versioning, partage et permission re-evaluation. Ce fichier possède uniquement le catalogue et les règles Command.

## Vues système

| Key | Intention | Critère initial |
|---|---|---|
| all | tout travail autorisé | aucun filtre de type/owner |
| incidents | coordination d'Incidents | object.type=Incident |
| tasks | tâches opérationnelles | object.type=Task |
| unassigned | travail sans owner | owner is empty |
| sla-risk | travail proche/en dépassement SLA | SLA policy-derived |
| my-work | travail assigné ou suivi par l'utilisateur | owner/follower=current principal |

Les critères précis dépendent des objets et policies de Phase 4 ; ils ne créent pas de permission.

## Règles

- vues système versionnées, non supprimables et adressables par `view=key` ;
- vue personnelle modifiable par owner ;
- vue partagée exige permission et audience ;
- ouverture réévalue permissions et colonnes ;
- dirty state visible avant update ;
- `Team Load` n'est pas une vue système ;
- aucune vue n'est un fichier d'écran actif.

## Critère

**Given** un lien legacy `/work-queue/unassigned`,  
**When** il est ouvert,  
**Then** il redirige vers le workspace unique avec `view=unassigned`, conserve les filtres sûrs et n'enregistre aucun nouvel écran.
