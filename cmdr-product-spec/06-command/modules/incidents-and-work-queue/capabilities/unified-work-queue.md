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
Les analystes doivent traiter un flux mixte sans naviguer entre pages clones. Sans cette capacité, vues, permissions et contexte divergent.

## 3. Objectifs
Conserver une route unique, fournir exactement six vues système, préserver filtres/tri/sélection/scroll/Inspector et réévaluer permission/fraîcheur à chaque ouverture.

## 4. Non-objectifs
Ne crée pas six pages, ne définit pas toutes les colonnes, ne réintroduit pas Team Load, ne possède pas Case/Decision/Run et n’exécute pas d’action dangereuse.

## 5. Propriétaire
Command / Incidents and Work Queue / Command Product Lead. Shared Capabilities possède le mécanisme Saved Views ; Command possède le catalogue système et la sémantique métier.

## 6. Utilisateurs
Principal : SOC Analyst L1/L2. Secondaires : Incident Commander, Team Lead et Business Owner en lecture.

## 7. Conditions d’entrée
Tenant/environnement sélectionnés, au moins un type d’objet autorisé et permission de lecture Incident ou Task.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Work items | Incident et Task | objets Command | oui | version courante | état Empty ou Partial explicite |
| Linked projections | Object Linking Service | Case, Decision et Run refs | non | résolue à l’ouverture | lien indisponible sans fuite de données |
| Saved View | Shared Saved Views et catalogue Command | configuration de workspace | oui | schéma courant | fallback `All` avec explication |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident | Command | coordination autorisée | consulter, filtrer, sélectionner et naviguer |
| Task | Command | coordination autorisée | consulter, filtrer, sélectionner et naviguer |
| Case / Decision / Response Run | produits propriétaires | statut et relation | consulter et naviguer en projection |
| Saved View | Shared Capabilities | configuration | appliquer permission-aware |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Saved View application | appliquer vue, filtres et tri | Shared/Command catalog | classe 0, aucune donnée métier copiée |
| Incident / Task | aucune mutation par simple affichage | Command | actions déléguées aux capabilities propriétaires |

## 11. Fonctionnalités
Changer de vue sans route distincte, filtrer/trier/grouper, sélectionner, ouvrir l’Inspector, ouvrir Incident Detail ou un objet lié et afficher fraîcheur/permissions/résultats partiels.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Changer de vue | lecteur | workspace | 0 | vue autorisée | même route, paramètre différent | non |
| Filtrer ou trier | lecteur | projection | 0 | champs autorisés | projection mise à jour | non |
| Ouvrir un item | lecteur | objet | 0 | permission objet | détail ou Inspector | non |
| Lancer une action locale | rôle autorisé | sélection | selon capability | préconditions propres | résultat délégué | selon classe |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Appliquer une vue | oui | oui | oui | suggestion facultative | six vues et filtres manuels |
| Filtrer et trier | oui | oui | oui | non nécessaire | opérateurs et tris déterministes |
| Suggérer une prochaine action | oui | facteurs sourcés | oui | proposition attribuée | capabilities propriétaires et décision humaine |

L’IA ne devient jamais la Work Queue.

## 14. États fonctionnels
`ready`, `no-authorized-items`, `partial-sources`, `stale-view`, `schema-migrated`, `permission-reduced`, `dirty-view`.

## 15. États d’interface
Loading conserve la structure ; Empty explique scope/absence ; Partial nomme les sources ; Error garde les données valides ; Offline est lecture seule ; Permission denied ne divulgue rien ; Stale expose source/date/effet.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Queue projection | état du workspace | opérateur | objets autorisés, fraîcheur et limites visibles |
| Selected context | état de route | Inspector et Incident Detail | vue, filtres, tri, scroll et sélection conservés |
| Delegated action context | événement | capability propriétaire | sélection et permissions réévaluées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Work Queue | ouverture Incident Detail | Command | view, filtres, sélection et scroll | retour exact |
| Work Queue | ouverture objet externe lié | produit propriétaire | tenant, environnement, item et return origin | retour exact à la file |
| Work Queue | action groupée | Bulk Coordination | sélection, versions, vue et filtres | nouvelle projection après résultats |

## 18. Dépendances
CAP-CMD-102, CAP-CMD-104, CAP-CMD-105, CAP-CMD-108, CAP-CMD-109, Saved Views, Data Table, Filters et Inspector.

## 19. Source de vérité
Incident/Task : Command. Saved View : Shared avec catalogue Command. Projections externes : propriétaires sources.

## 20. Provenance et audit
Application de vue, migration de schéma, permission réduite et actions déléguées enregistrent acteur, source/version, scope et résultat.

## 21. Permissions fonctionnelles
`perm.command.read`, `perm.command.incident.read`, future lecture Task, permissions projections et partage de vue distinct ; atomisation reportée.

## 22. Limites et erreurs
Vue incompatible, donnée stale, projection indisponible, colonne interdite, tenant incompatible ou refus conservent route et données valides.

## 23. Métriques
Temps file→owner, taux de retour exact, items masqués et vues migrées ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** les six vues, **When** elles sont sélectionnées, **Then** la route reste identique et aucun écran autonome n’est créé.

**Given** aucun modèle IA, **When** la file est utilisée, **Then** Saved Views, filtres et actions manuelles suffisent.

**Given** une colonne interdite, **When** une vue est ouverte, **Then** la colonne est retirée sans modifier la source ni divulguer de donnée.

## 26. Questions ouvertes
Quelles colonnes finales et quelle virtualisation seront retenues ? — Requirement IDs ci-dessus. Aucun nouvel OPEN.

## 27. Consommateurs documentaires
Work Queue, Incident Detail, Mission Control, parcours Phase 5, workspace Phase 6, objets Phase 7 et permissions/contrats ultérieurs.