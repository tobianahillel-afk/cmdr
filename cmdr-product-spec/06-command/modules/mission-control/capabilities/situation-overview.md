---
id: CAP-CMD-001
title: Situation Overview
product: command
module: mission-control
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-013
  - REQ-PROD-008
  - REQ-PROD-010
  - REQ-PROD-021
open_decisions:
  - none
source-of-truth: canonical
---

# CAP-CMD-001 — Situation Overview

## 1. Définition
Agrège une lecture opérationnelle priorisée des Incidents, Tasks urgentes, services affectés, blocages, Decisions en attente, Response Runs actifs et Results récents, avec fraîcheur, owner et prochaine action.

## 2. Problème utilisateur
L’Incident Commander doit comprendre les changements importants sans reconstruire la situation depuis plusieurs produits. Sans cette capacité, priorités, dépendances et handovers divergent.

## 3. Objectifs
- montrer les changements significatifs et leurs sources ;
- relier chaque item à son owner et sa prochaine action ;
- exposer données manquantes, conflictuelles ou stale ;
- préserver le contexte lors des transitions.

## 4. Non-objectifs
Ne remplace ni la Work Queue, ni l’investigation, ni Govern ; ne devient pas un dashboard de widgets ou un chatbot.

## 5. Propriétaire
Command / Mission Control / Command Product Lead. Command possède la composition opérationnelle ; les objets externes restent propriétaires de leurs données.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : SOC Analyst L1/L2, Business Owner, Response Operator.

## 7. Conditions d’entrée
Tenant et environnement sélectionnés, au moins une source autorisée et permission de lecture Command.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Situation items | Incident, Task et projections interproduits | objets et événements | non | temps réel ou dernière fraîcheur connue | situation `partial` et source absente nommée |
| Service impact | Business Service Catalog et Incident | projection métier | non | timestamp de la source | impact `unknown`, sans inférence silencieuse |
| Pending governance | Govern | Decision et Response Run projections | non | état courant | contexte Incident conservé, Govern signalé indisponible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident | Command | priorité, état, owner, impact, prochaine action | consulter et naviguer |
| Task | Command | urgence, owner, échéance, blocage | consulter et naviguer |
| Decision / Response Run / Result | Govern | statut, portée et résultat résumé | consulter et relier en lecture seule |
| Service | Shared Business Service Catalog | owner, criticité et dépendances | consulter et filtrer |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident | modifier la prochaine action ou accuser la situation | Command | classe 2, action explicite et auditée |
| Projections externes | aucune mutation | produits propriétaires | lecture seule ; aucun transfert d’ownership |

## 11. Fonctionnalités
Classer les changements par facteurs visibles, montrer Incidents/Tasks/blocages/Decisions/Runs, exposer qualité et fraîcheur, ouvrir l’objet source et restaurer Mission Control.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter la situation | rôle autorisé | projections | 0 | tenant sélectionné | lecture filtrée | non |
| Ouvrir l’objet source | rôle autorisé | objet référencé | 0 | permission destination | transition avec return origin | non |
| Modifier la prochaine action | Incident Commander | Incident | 2 | version courante | Incident audité | selon OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Composer la situation | oui | oui | oui | résumé attribué | agrégation sourcée, filtres et sélection humaine |
| Prioriser l’affichage | oui | oui | oui | suggestion seulement | règles explicables et ordre corrigible manuellement |
| Modifier la prochaine action | oui | validation/version | workflow possible | proposition seulement | mutation humaine complète |

## 14. États fonctionnels
`coherent`, `partial`, `stale`, `conflicting`, `critical-change`, `no-active-work`.

## 15. États d’interface
Loading conserve le contexte ; Empty explique l’absence ; Partial nomme les sources ; Error conserve les données valides ; Offline bloque les mutations ; Permission denied ne divulgue rien ; Stale expose source, date et conséquence.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Situation snapshot | projection datée | Mission Control et Handover | source, fraîcheur et owner conservés |
| Navigation context | événement de transition | produit destination | tenant, environnement et return origin préservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Mission Control | ouverture Incident ou Case | Command ou Investigate | tenant, environnement, objet et période | même position de situation |
| Mission Control | ouverture Decision ou Run | Govern | Incident, impact et return origin | situation restaurée |
| Result | projection disponible | Mission Control | Result, Run, Decision et Incident liés | Govern reste propriétaire |

## 18. Dépendances
CAP-CMD-003, CAP-CMD-005, CAP-CMD-006, Global Search, Object Linking Service, Timeline Engine et Business Service Catalog.

## 19. Source de vérité
Incident et Task restent Command ; objets Govern et Service sont des projections permission-aware ; calculs et classements exposent leurs facteurs.

## 20. Provenance et audit
Acteur ou producteur, source, version/run, avant/après, justification, résultat, tenant, environnement et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.read`, permissions des projections et ABAC tenant/environnement ; atomisation reportée.

## 22. Limites et erreurs
Source absente, stale, interdite, conflit de version ou destination indisponible rendent la situation partielle et conservent le workspace source.

## 23. Métriques
Temps jusqu’à owner/prochaine action, part des items sourcés et taux de retour exact ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire uniquement, promotion après objets, permissions, parcours, écrans, contrats, implémentation et validation.

## 25. Critères d’acceptation
**Given** des sources courantes, **When** la situation est ouverte, **Then** chaque item expose source, fraîcheur, owner et prochaine action sans dupliquer la Work Queue.

**Given** aucun fournisseur de modèle, **When** Mission Control est utilisé, **Then** règles, agrégations et actions manuelles fournissent le résultat essentiel.

**Given** une source indisponible, **When** la situation est calculée, **Then** la lacune est visible et les autres données restent utilisables.

## 26. Questions ouvertes
Quels changements sont transversaux et quelles fenêtres de fraîcheur s’appliquent ? — REQ-PROD-013, REQ-PROD-008, REQ-PROD-010, REQ-PROD-021. Aucun nouvel OPEN.

## 27. Consommateurs documentaires
Mission Control Now/Situation, Handover, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions/contrats ultérieurs.