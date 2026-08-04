---
id: CAP-CMD-101
title: Unified Work Queue
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-OBJ-001
  - REQ-OBJ-012
  - REQ-PROD-013
  - REQ-UX-008
  - REQ-UX-009
open_decisions:
  - none
source-of-truth: canonical
---
# CAP-CMD-101 — Unified Work Queue

## 1. Définition
Fournit un workspace unique de coordination pour Incidents et Tasks avec ownership, priorité, état, SLA, âge, service, blocage, prochaine action et liens vers Case, Decision et Response Run.

## 2. Problème utilisateur
Les analystes doivent traiter un flux mixte sans naviguer entre pages clones. Rôles : SOC Analyst L1/L2, Incident Commander, Team Lead. Sans la capacité, vues et permissions divergent et le contexte se perd.

## 3. Objectifs
- conserver une seule route ;
- fournir exactement All, Incidents, Tasks, Unassigned, SLA Risk et My Work ;
- préserver filtres, tri, sélection, scroll et Inspector ;
- réévaluer permissions et fraîcheur à chaque ouverture.

## 4. Non-objectifs
Ne pas créer six pages, définir toutes les colonnes, réintroduire Team Load, posséder Case/Decision/Run ou exécuter une action dangereuse.

## 5. Propriétaire
Command / Incidents and Work Queue / Command Product Lead. Shared Capabilities possède le mécanisme Saved Views ; Command possède le catalogue système et le contenu métier.

## 6. Utilisateurs
Principal : SOC Analyst L1/L2. Secondaires : Incident Commander, Team Lead, Business Owner en lecture. Tenant, environnement, scope et permissions sont obligatoires.

## 7. Conditions d’entrée
Tenant/environnement sélectionnés, au moins un type d’objet autorisé et permission de lecture Incident ou Task.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Work items | Incident/Task | objets Command | oui | version courante | Empty ou Partial explicite |
| Linked projections | Object Linking Service | Case/Decision/Run refs | non | résolue à l’ouverture | lien indisponible sans fuite |
| Saved View | Shared + catalogue Command | configuration | oui | schema courant | fallback All avec explication |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident | Command | coordination complète autorisée | lecture/navigation |
| Task | Command | coordination complète autorisée | lecture/navigation |
| Case / Decision / Response Run | produits propriétaires | statut et relation | projection |
| Saved View | Shared | configuration | application permission-aware |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Saved View application | appliquer vue/filtre/tri | Shared/Command catalog | classe 0 |
| Incident / Task | aucune mutation par simple affichage | Command | actions déléguées aux capabilities propriétaires |

## 11. Fonctionnalités
Changer de vue sans route distincte ; filtrer/trier/grouper ; sélectionner un item ; ouvrir l’Inspector ; ouvrir Incident Detail ou un objet lié ; afficher freshness, permission et résultats partiels.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Changer de vue | lecteur | workspace | 0 | vue autorisée | même route, paramètre différent | non |
| Filtrer/trier | lecteur | projection | 0 | champs autorisés | projection mise à jour | non |
| Ouvrir un item | lecteur | objet | 0 | permission objet | détail/Inspector | non |
| Lancer action locale | rôle autorisé | sélection | selon capability | préconditions propres | résultat délégué | selon classe |

## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur déterministe | Workflow | Agent | Govern | Alternative sans IA |
|---|---|---|---|---|---|---|---|
| Appliquer une vue | oui | oui | filtrage/tri | possible | suggestion facultative | non | vue et filtres manuels |
| Suggérer prochaine action | décision humaine | possible | facteurs sourcés | possible | proposition attribuée | selon effet | capacités propriétaires |
L’IA ne devient jamais la Work Queue.

## 14. États fonctionnels
`ready`, `no-authorized-items`, `partial-sources`, `stale-view`, `schema-migrated`, `permission-reduced`, `dirty-view`.

## 15. États d’interface
Loading conserve la structure ; Empty explique scope ou absence ; Partial nomme sources ; Error garde données valides ; Offline est lecture seule ; Permission denied ne divulgue rien ; Stale affiche source/date/effet. Rendu : Design System.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Queue projection | workspace state | opérateur | objets autorisés et fraîcheur visibles |
| Selected context | route state | Inspector/Incident Detail | vue, filtres, tri, scroll et sélection |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Work Queue | Incident Detail | Command | view, filters, selection, scroll | retour exact |
| Work Queue | objet externe lié | produit propriétaire | tenant, env, item, return origin | retour exact à la file |

## 18. Dépendances
CAP-CMD-102, 104, 105, 108, 109 ; Saved Views, Data Table, Filters et Inspector. Aucune dépendance ne transfère l’ownership.

## 19. Source de vérité
Incident/Task restent Command ; Saved View reste Shared avec catalogue Command ; projections externes restent permission-aware. La file ne copie pas leurs cycles de vie.

## 20. Provenance et audit
Application de vue, migration de schéma, permission réduite et actions déléguées enregistrent acteur, source/version, scope, résultat et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.read`, `perm.command.incident.read`, future `perm.command.task.read`, permissions des projections et partage de vue distinct. Atomisation reportée.

## 22. Limites et erreurs
Vue incompatible, données stale, projection indisponible, colonne interdite, conflit de tenant ou permission refusée conservent la route et les données valides sans divulgation.

## 23. Métriques
Temps file→owner ; taux de retour exact ; items masqués pour permission ; vues migrées ; travail sans next action. Aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native. Preuve documentaire seulement ; promotion après objets, permissions, parcours, écran/workspace, contrats, implémentation et validation.

## 25. Critères d’acceptation
**Given** trois Incidents, quatre Tasks, un unassigned, un SLA Risk et deux My Work, **When** les six vues sont sélectionnées, **Then** la route reste identique, le paramètre change, les objets autorisés sont affichés et aucune page n’est créée.

**Given** aucun modèle IA, **When** la file est utilisée, **Then** Saved Views, filtres et actions manuelles fournissent tout le résultat essentiel.

**Given** une colonne devenue interdite, **When** une vue partagée est ouverte, **Then** la colonne est retirée, la vue source reste inchangée et aucune donnée n’est divulguée.

## 26. Questions ouvertes
Quelles colonnes seront essentielles par vue et quelle stratégie de virtualisation sera retenue ? — REQ-OBJ-001, REQ-OBJ-012, REQ-PROD-013, REQ-UX-008, REQ-UX-009. Aucun nouvel OPEN.

## 27. Consommateurs documentaires
Work Queue, Incident Detail, Mission Control, parcours triage Phase 5, workspace Phase 6, objets Phase 7 et permissions/contrats ultérieurs.
