---
id: CAP-CMD-006
title: Recent Results and Outcomes
product: command
module: mission-control
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-005
  - REQ-PROD-008
  - REQ-PROD-013
  - REQ-PROD-021
open_decisions:
  - OPEN-013
source-of-truth: canonical
---

# CAP-CMD-006 — Recent Results and Outcomes

## 1. Définition
Consomme les Results Govern récents et les relie aux Response Runs, Decisions et Incidents afin d’actualiser situation, risque résiduel et prochaine action sans modifier le Result.

## 2. Problème utilisateur
Une exécution peut réussir techniquement sans vérification ni effet métier compris dans Command. L’Incident risque alors d’être fermé trop tôt ou le travail relancé.

## 3. Objectifs
Distinguer outcome technique/vérification/impact métier, relier Result/Run/Decision/Incident, réinjecter le résultat dans la situation et permettre contestation sans réécriture.

## 4. Non-objectifs
Ne crée ni ne valide un Result, ne détermine pas la réussite d’un Run, ne modifie pas une Decision et ne masque pas un résultat disputé.

## 5. Propriétaire
Command / Mission Control / Command Product Lead pour la consommation opérationnelle ; Govern reste propriétaire des objets de réponse.

## 6. Utilisateurs
Principal : Incident Commander. Secondaires : Response Operator, SOC Analyst L2, Business Owner.

## 7. Conditions d’entrée
Result autorisé ou projection partielle, lien Run/Decision et Incident cible identifiable ou association à confirmer.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Result projection | Govern | outcome, verification et residual risk | oui | état courant | `pending-verification` ou association manuelle |
| Incident context | Command | statut, impact et prochaine action | oui | version courante | aucune mutation si inaccessible |
| Business validation | Business Owner ou Incident Commander | confirmation opérationnelle | non | datée | technique et métier restent distincts |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Result / Response Run / Decision | Govern | résultat, vérification, portée et conditions | consulter et relier en lecture seule |
| Incident | Command | état, impact et prochaine action | consulter et modifier si autorisé |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident | lier Result, ajuster prochaine action ou état de travail | Command | classe 2, jamais automatique par défaut |
| Timeline Entry | produire événement de consommation | Shared Timeline / contenu Command | référence stable vers Result source |
| Result | aucune mutation | Govern | projection en lecture seule |

## 11. Fonctionnalités
Présenter Results et fraîcheur, différencier succeeded/verified/business-accepted, associer à un Incident, marquer contestation/vérification manquante et créer un suivi.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter Result | rôle autorisé | Result | 0 | permission Govern | projection visible | non |
| Associer à Incident | coordinateur | Incident | 2 | objets compatibles | relation auditée | OPEN-013 |
| Modifier prochaine action | Incident Commander | Incident | 2 | effet compris | Incident mis à jour | OPEN-013 |
| Demander vérification | Incident Commander | Task ou retour Govern | 2 | résultat absent/contesté | suivi attribué | selon destination |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Classer les Results | oui | oui | oui | résumé attribué | filtres et règles de statut |
| Proposer une association Incident | oui | matching sourcé | oui | suggestion seulement | recherche et liaison manuelles |
| Mettre à jour la situation | oui | validation/version | workflow possible | proposition seulement | mutation Incident manuelle |
| Préparer un suivi | oui | complétude/déduplication | oui | brouillon Task | création manuelle de Task |

## 14. États fonctionnels
`received`, `pending-verification`, `verified`, `disputed`, `applied-to-situation`, `superseded`.

## 15. États d’interface
Partial montre projections manquantes ; Offline bloque mutation ; Permission denied ne révèle pas Result ; Stale affiche source/date ; Error conserve Incident local.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Outcome projection | Result lié | Mission Control, Incident et Reporting | source Govern et niveau de vérification visibles |
| Operational follow-up | Incident update ou Task | Work Queue et Readiness | attribué, audité et sans mutation du Result |
| Contestation context | événement de transition | Govern ou Investigate | motif, Incident et Result conservés |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Govern Result | publication autorisée | Command Incident | Result, Run, Decision et residual risk | retour Govern possible |
| Command Incident | contestation | Govern | Result, motif, Incident et demande de vérification | Incident restauré |
| Command Incident | investigation complémentaire | Investigate | Incident, Result et question technique | Incident restauré |

## 18. Dépendances
CAP-CMD-001, CAP-CMD-003, CAP-CMD-106, Object Linking Service, Timeline Engine et Reporting Engine.

## 19. Source de vérité
Govern possède Result/Run/Decision ; Command possède seulement les relations et mutations Incident explicites.

## 20. Provenance et audit
Association, contestation, mutation ou proposition enregistrent acteur/producteur, source, before/after, justification, tenant et correlation ID.

## 21. Permissions fonctionnelles
`perm.govern.result.read`, `perm.command.incident.manage`, `perm.command.task.manage` ; atomisation et `OPEN-013` reportés.

## 22. Limites et erreurs
Result absent, stale, inaccessible, non associé ou conflit d’Incident ne doivent jamais être présentés comme outcome complet.

## 23. Métriques
Results associés, distinction vérification technique/métier et délai Result→prochaine action ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire uniquement.

## 25. Critères d’acceptation
**Given** un Result vérifié, **When** Command le consomme, **Then** Run, Decision, vérification et risque résiduel restent visibles.

**Given** un Result disputé, **When** une vérification est demandée, **Then** Command crée seulement le suivi et ne modifie pas le Result.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** filtres, règles et actions manuelles fournissent le résultat essentiel.

## 26. Questions ouvertes
Quelle validation métier précède `applied-to-situation` et quand repartir vers Govern ou Investigate ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Mission Control, Incident Detail, Readiness, reports, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions/contrats ultérieurs.