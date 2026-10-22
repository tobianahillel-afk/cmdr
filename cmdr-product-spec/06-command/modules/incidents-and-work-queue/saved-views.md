---
id: command-work-queue-saved-views
domain: 06-command
status: draft
owner: Command Product Lead
updated: 2026-08-04
source-of-truth: canonical
requirements:
  - REQ-OBJ-012
  - REQ-UX-008
  - REQ-UX-009
capability_ids:
  - CAP-CMD-101
---

# Saved Views de la Work Queue

## Propriété

Shared Capabilities possède stockage/versioning/partage/migration/application permission-aware. Command possède le catalogue et la signification. Le Design System possède l’interaction. CAP-CMD-101 possède le comportement du workspace.

## Vues système exactes

| Key | Nom | Intention | Critère fonctionnel initial |
|---|---|---|---|
| `all` | All | tout travail autorisé | aucun filtre type/owner imposé |
| `incidents` | Incidents | coordonner les Incidents | `object.type=Incident` |
| `tasks` | Tasks | coordonner les Tasks | `object.type=Task` |
| `unassigned` | Unassigned | travail sans owner principal | owner absent |
| `sla-risk` | SLA Risk | échéance proche ou violation | projection de policy SLA |
| `my-work` | My Work | travail assigné ou suivi par le principal | owner/contributor/watcher autorisé |

## Invariants

- une seule route/workspace ;
- la vue est un paramètre, pas une Page ;
- six vues versionnées et non supprimables ;
- permissions/colonnes/objets réévalués à l’ouverture ;
- une Saved View n’enregistre ni données, ni permission, ni secret ;
- `Team Load` n’est pas une vue système ;
- aucun ancien Screen ID n’est réactivé.

## Migration

`CMD-IWQ-001..004` ouvrent la vue correspondante. `CMD-IWQ-005` ouvre la Work Queue avec notice de migration ; il ne devient pas `team-load`.

## Critères d’acceptation

**Given** `view=unassigned`, **When** elle est ouverte, **Then** la route reste identique, les permissions sont réévaluées et aucun écran autonome n’est créé.

**Given** une colonne devenue interdite dans une vue partagée, **When** elle est ouverte, **Then** la colonne est retirée sans modifier la source ni divulguer de donnée.

**Given** un ancien lien Team Load, **When** il est résolu, **Then** aucune vue système `team-load` n’existe et l’ancien ID reste deprecated.
