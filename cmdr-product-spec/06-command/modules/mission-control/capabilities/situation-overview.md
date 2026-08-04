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

**Situation.** L’Incident Commander doit comprendre ce qui change maintenant sans reconstruire la situation à partir de files, rapports et produits séparés.

**Utilisateurs concernés.** Incident Commander en premier lieu ; SOC Analyst L1/L2, Business Owner, Response Operator.

**Conséquence sans la capacité.** Sans cette capacité, les priorités divergent, les dépendances critiques restent invisibles et les handovers commencent avec une situation incomplète.

## 3. Objectifs

- rendre les changements significatifs visibles avant les métriques secondaires ;
- relier chaque élément à son objet source et à son owner ;
- montrer les données manquantes, conflictuelles ou stale ;
- permettre un passage direct vers l’Incident, le Case, Govern ou le Run concerné.

## 4. Non-objectifs

- remplacer la Work Queue ;
- devenir un dashboard de widgets configurables ;
- effectuer une investigation technique ;
- transformer une suggestion en priorité effective.

## 5. Propriétaire

Command / Mission Control / Command Product Lead. Command possède la composition opérationnelle et les mutations explicites d’Incident ; les projections externes gardent leur owner.

## 6. Utilisateurs

Rôle principal : Incident Commander. Rôles secondaires : SOC Analyst L1/L2, Business Owner, Response Operator. Chacun agit uniquement dans son tenant, environnement, scope et permissions.

## 7. Conditions d’entrée

Tenant et environnement sélectionnés ; au moins une source Command ou projection autorisée ; permission générale de lecture Command.

## 8. Entrées fonctionnelles

| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Situation items | Incident/Task et projections interproduits | objets et événements | non | temps réel ou dernière fraîcheur connue | afficher une situation partielle et nommer la source absente |
| Service impact | Business Service Catalog et Incident | projection métier | non | fraîcheur déclarée | marquer impact inconnu, jamais l’inférer silencieusement |
| Pending governance | Govern | Decision/Run projections | non | état courant | conserver le contexte Incident même si Govern est indisponible |

## 9. Objets lus

| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident | Command | priorité, état, owner, impact, prochaine action | lecture et navigation |
| Task | Command | urgence, owner, échéance, blocage | lecture et navigation |
| Decision / Response Run / Result | Govern | statut, portée, résultat résumé | projection seulement |
| Service | Shared Business Service Catalog | owner, criticité, dépendances | projection seulement |

## 10. Objets créés ou modifiés

| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident | mise à jour de prochaine action ou accusé de situation | Command | classe 2 ; aucune modification automatique par l’agrégateur |

Aucune projection externe ne transfère son ownership à Command.

## 11. Fonctionnalités

- classer les changements par impact opérationnel explicable ;
- afficher incidents critiques, tâches urgentes et blocages ;
- mettre en évidence Decisions en attente et Runs actifs ;
- afficher fraîcheur, source et qualité de chaque projection ;
- ouvrir l’objet source et restaurer Mission Control au retour.

## 12. Actions utilisateur

| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter la situation | tous rôles autorisés | projections | 0 | tenant sélectionné | lecture filtrée | non |
| Ouvrir l’objet source | tous rôles autorisés | objet référencé | 0 | permission destination | transition avec return origin | non |
| Mettre à jour la prochaine action | Incident Commander | Incident | 2 | Incident accessible et version courante | Incident audité | selon OPEN-013 |

Les classes 3 et 4 ne sont jamais exécutées par Command.

## 13. Automatisation et IA

| Fonction | Humain | Règle | Moteur déterministe | Workflow | Agent | Govern | Alternative sans IA |
|---|---|---|---|---|---|---|---|
| Consulter la situation | oui | sélection explicable | agrégation sourcée | possible | résumé attribué | non | lecture et filtres manuels |
| Mettre à jour la prochaine action | oui, décision finale | possible si policy | validation/version | possible | proposition uniquement | OPEN-013 selon effet | mutation manuelle complète |

L’absence de modèle ne bloque aucune fonction essentielle.

## 14. États fonctionnels

`coherent`, `partial`, `stale`, `conflicting`, `critical-change`, `no-active-work`. Ils décrivent la projection de situation, pas les machines d’état des objets.

## 15. États d’interface

Loading conserve le contexte ; Empty explique l’absence ; Partial nomme les sources manquantes ; Error conserve les données valides ; Offline bloque les mutations non garanties ; Permission denied ne divulgue rien ; Stale affiche source, date et conséquence. Le Design System reste propriétaire du rendu.

## 16. Sorties

| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Situation snapshot | projection datée | Mission Control et handover | chaque item conserve source, fraîcheur et owner |
| Navigation context | return origin | produit destination | tenant, environnement et objet source conservés |

## 17. Transitions

| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Mission Control | ouverture d’un Incident ou Case | Command / Investigate | tenant, environnement, objet, période | retour à la même position |
| Mission Control | ouverture d’une Decision ou d’un Run | Govern | Incident et contexte métier | retour à la situation |
| Result | projection disponible | Mission Control | Result, Run, Decision et Incident liés | l’objet source reste Govern |

## 18. Dépendances

CAP-CMD-003, CAP-CMD-005, CAP-CMD-006, Shared Global Search, Object Linking Service, Timeline Engine et Business Service Catalog. Les dépendances n’impliquent aucun transfert d’ownership.

## 19. Source de vérité

Incident reste propriétaire Command. Decision, Response Run, Result et Service sont des projections par référence stable et permission-aware. Chaque donnée conserve source, version ou timestamp, fraîcheur et classification ; les calculs exposent facteurs et limites.

## 20. Provenance et audit

Toute mutation ou proposition enregistre acteur/producteur, source, version/run, objet, avant/après, justification, résultat, tenant, environnement et correlation ID.

## 21. Permissions fonctionnelles

`perm.command.read`, permissions de lecture des projections Govern/Investigate et ABAC tenant/environnement. L’atomisation et les règles finales sont reportées.

## 22. Limites et erreurs

Données absentes/non autorisées, projection stale, dépendance indisponible, conflit de version, tenant/environnement incompatible ou refus d’autorisation empêchent tout résultat présenté comme complet. Une transition échouée conserve le workspace source.

## 23. Métriques

- temps jusqu’à identification de l’owner et de la prochaine action ;
- part des items avec fraîcheur et source visibles ;
- taux de transitions restaurant exactement le contexte.

Aucune cible définitive n’est fixée.

## 24. Classification de livraison

`delivery_status: defined`, `delivery_mode: planned`, cible native. La preuve actuelle est documentaire uniquement ; la promotion exige objets, permissions, parcours, écrans, contrats, implémentation et validation.

## 25. Critères d’acceptation

**Given** un Incident Commander autorisé et des sources courantes, **When** il consulte la situation, **Then** chaque élément expose source, fraîcheur, owner et prochaine action sans dupliquer la Work Queue.

**Given** aucun fournisseur de modèle, **When** Mission Control est ouvert, **Then** règles, agrégations déterministes, filtres et actions manuelles fournissent le résultat essentiel.

**Given** une source indisponible ou non autorisée, **When** la situation est calculée, **Then** la lacune et son effet sont visibles, aucune donnée protégée n’est révélée et les autres données restent utilisables.

## 26. Questions ouvertes

- Quels changements méritent une mise en avant transversale sans score opaque ? — REQ-PROD-013, REQ-PROD-008, REQ-PROD-010, REQ-PROD-021.
- Quelle fenêtre de fraîcheur s’applique par source ? — mêmes Requirement IDs.

Aucune nouvelle décision ouverte n’est créée.

## 27. Consommateurs documentaires

Mission Control Now/Situation, handover, futurs écrans exécutifs, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions/contrats ultérieurs.
