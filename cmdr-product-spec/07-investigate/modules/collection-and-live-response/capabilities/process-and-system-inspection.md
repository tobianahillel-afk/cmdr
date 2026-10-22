---
id: CAP-INV-206
title: Process and System Inspection
product: investigate
module: collection-and-live-response
owner: Investigate Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-014
  - REQ-PROD-018
  - REQ-SEC-001
open_decisions:
  - OPEN-008
  - OPEN-013
source-of-truth: canonical
---
# CAP-INV-206 — Process and System Inspection

## 1. Définition
Inspecter en lecture les processus, relations, sessions, services, connexions et informations système autorisées, puis déclencher une collecte complémentaire ou préparer une action gouvernée.

## 2. Problème utilisateur
Une inspection live peut être confondue avec une mutation ou un containment. L’analyste doit connaître l’instantané observé, sa fraîcheur et ses limites avant de tirer une conclusion ou de préparer une action.

## 3. Objectifs
- afficher un snapshot horodaté des processus et relations parent/enfant
- consulter sessions, services/composants, connexions et informations système autorisées
- rechercher, filtrer et sélectionner sans modifier l’Endpoint
- produire un résultat d’inspection ou un Artifact et préparer séparément toute action risquée

## 4. Non-objectifs
Ne pas terminer, suspendre, démarrer, arrêter ou isoler directement un composant ; ne pas définir commandes, moteurs, schémas de télémétrie ou couverture plateforme ; ne pas présenter l’inspection comme containment.

## 5. Propriétaire
Investigate / Collection and Live Response / Investigate Product Lead possède le contexte métier, les drafts et les relations au Case. Platform Settings reste propriétaire de Fleet et Endpoint Policies ; Endpoint Agent exécute et rapporte localement ; Govern conserve l’autorité, Decision, Response Run et Result.

## 6. Utilisateurs
Principal : Case Analyst ou Response Operator. Secondaires : DFIR Analyst, Investigation Lead et Evidence Reviewer.

## 7. Conditions d’entrée
Case et Endpoint accessibles, capability d’inspection déclarée, catégories demandées, snapshot/fraîcheur disponibles, policy et permission de lecture validées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Case et question d’investigation | Investigate | contexte analytique | oui | version courante | inspection standalone interdite sauf raison |
| Endpoint/Agent et snapshot time | Endpoint Agent | cible et observation | oui | horodatage visible | offline/stale |
| Catégories d’inspection | analyste | processus, sessions, services, connexions, système | oui | validées à la requête | incomplete |
| Policy et permission | Settings / Security | droit de lecture et restrictions | oui | snapshot à l’action | denied/restricted |
| Contexte historique éventuel | Telemetry/Case | comparaison | non | période/source visibles | snapshot courant seulement |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case / Hypothesis / Finding | Investigate | question, scope et conclusions existantes | consulter/lier |
| Endpoint / Endpoint Agent | partagé / Endpoint Agent | état, snapshot et capacités | consulter |
| Endpoint Policy | Platform Settings | catégories et données autorisées | consulter uniquement |
| Agent Command | Endpoint Agent | demande d’inspection et statut | consulter |
| Operation Result / Artifact | concept futur / Investigate | snapshot reçu et materialisation | consulter/lier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Inspection request/result conceptuel | créer et enregistrer le snapshot | Investigate, modèle futur | read-only, horodaté et Case-scoped |
| Artifact | créer uniquement si le snapshot est matérialisé | Investigate | source, catégories et acquisition requises |
| Collection Request draft | préparer depuis une sélection | Investigate | scope borné et classe recalculée |
| Endpoint/process/service | aucune mutation | Endpoint Agent / système cible | inspection ne modifie rien |

## 11. Fonctionnalités
- afficher processus, ancestry/descendants, sessions, services/composants et connexions
- rechercher et filtrer en conservant le timestamp du snapshot
- ouvrir les métadonnées autorisées d’un élément sélectionné
- demander une acquisition ciblée depuis un processus, fichier ou connexion
- préparer une Action Request pour toute modification/containment au lieu de l’exécuter

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Lire/filtrer snapshot | Case Analyst | inspection result | 0 | read permission | vue horodatée | non |
| Comparer deux snapshots | DFIR Analyst | inspection results | 0 | snapshots accessibles | diff sourcé | non |
| Créer un Artifact du snapshot | analyste autorisé | Artifact | 1 | résultat matériel et provenance | Artifact lié | non |
| Préparer collecte complémentaire | Case Analyst | Collection Request draft | 1 | sélection et scope borné | passage CAP-INV-202 | selon impact |
| Préparer une modification/containment | Investigation Lead | Action Request | 2/3/4 | Finding/Evidence/impact/rollback | passage CAP-INV-215/113 | selon classe/obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| Filtrer et grouper | oui | règles/champs | oui | résumé | filtres manuels |
| Signaler un snapshot stale | oui | règle de fraîcheur | oui | explication | timestamp visible |
| Proposer une collecte complémentaire | oui | mappings explicites | oui | suggestion | sélection manuelle |
| Résumer les différences | oui | diff déterministe | oui | résumé attribué | diff brut |
| Modifier un processus/service | non dans cette capability | interdit | non | interdit | Action Request/Govern |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run et Tool Calls lorsqu’ils existent, sources, paramètres fonctionnels, timestamp, statut, incertitude, owner humain, acceptation/modification/rejet et trace.

## 14. États fonctionnels
`available`, `loading-snapshot`, `partial`, `stale`, `restricted`, `unsupported`, `failed`. Ces états décrivent la vue d’inspection, pas les objets de l’Endpoint.

## 15. États d’interface
Loading conserve Case/catégories ; Empty distingue aucun élément et permission ; Partial nomme les catégories manquantes ; Error conserve le snapshot valide ; Offline affiche le dernier snapshot comme stale ; Permission denied ne divulgue rien ; Stale affiche l’âge.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Snapshot d’inspection | Operation Result conceptuel | Case Workspace | timestamp, catégories, source et limites visibles |
| Artifact d’inspection | Artifact | CAP-INV-105/107 | création explicite et provenance conservée |
| Collection Request draft | request | CAP-INV-202 | sélection et raison transmises |
| Action context | package | CAP-INV-215/113 | aucune mutation ou containment exécuté |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Endpoint Context | ouvrir inspection | CAP-INV-206 | Case, Endpoint, catégories, policy, snapshot request | Endpoint Context |
| Inspection selection | collecter complément | CAP-INV-202 | objet sélectionné, scope, reason, class | Inspection |
| Inspection result | materialiser | CAP-INV-105 | snapshot, source, timestamp, limitations | Inspection |
| Finding/inspection | proposer modification | CAP-INV-215/113 | Finding, Evidence, target, impact, rollback | Case |

## 18. Dépendances
CAP-INV-201/202/105/107/212/214/215, Endpoint Agent investigation support, Settings Policy/Health, Shared Inspector/Trace et OPEN-008/013.

## 19. Source de vérité
Endpoint Agent reste source du snapshot local. Investigate possède la demande, le contexte analytique, les relations et tout Artifact créé. Settings reste source de la Policy ; Govern reste source de toute autorité de modification.

## 20. Provenance et audit
Case, Endpoint/Agent, catégories, snapshot time, source freshness, policy, initiateur, filtres, sélection, Artifact éventuel, demandes dérivées et correlation IDs.

## 21. Permissions fonctionnelles
Endpoint/system/process/session/service/connection read, sensitive output read, Artifact create, collection prepare et containment request. Aucune permission de mutation n’est conférée par l’inspection.

## 22. Limites et erreurs
Snapshot incomplet ou stale, processus terminé entre lecture et action, catégorie unsupported, données redacted, permission refusée, Endpoint offline ou horloge divergente. Aucune observation n’est présentée comme état durable.

## 23. Métriques
Snapshots par catégorie, âge moyen, taux partial/stale, pivots vers collecte, Artifacts créés et actions risquées correctement redirigées vers Govern.

## 24. Classification de livraison
`defined` / `planned`. Aucun moteur d’inspection, commande, protocole ou support plateforme n’est prouvé.

## 25. Critères d’acceptation
**Given** un Endpoint disponible **When** l’analyste consulte processus et connexions **Then** le snapshot, son timestamp et ses limitations sont visibles et aucune mutation n’a lieu.

**Given** une action de suspension proposée **When** l’analyste la sélectionne **Then** elle n’est pas exécutée comme inspection et un contexte d’Action Request est préparé.

**Given** aucun modèle IA **When** l’inspection est utilisée **Then** filtres, recherche, diff et pivots déterministes restent disponibles.

## 26. Questions ouvertes
OPEN-008 conserve la couverture plateforme ; OPEN-013 la gouvernance d’actions réversibles préparées. La sémantique finale des snapshots et résultats appartient à la phase Objets.

## 27. Consommateurs documentaires
Case Workspace, Collection Request Preparation, Artifact/Evidence, Containment Preparation, Endpoint Agent investigation et phases Objets/Permissions.
