---
id: CAP-CMD-105
title: SLA Tracking
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-005
  - REQ-PROD-013
  - REQ-PROD-053
  - REQ-OBJ-012
open_decisions:
  - OPEN-006
  - OPEN-013
source-of-truth: canonical
---

# CAP-CMD-105 — SLA Tracking

## 1. Définition
Projette le SLA applicable à un Incident/Task, le temps restant, le risque, les pauses autorisées et la provenance tenant/contractuelle, sans définir la policy.

## 2. Problème utilisateur
Les équipes doivent connaître une échéance et la validité d’une pause, pas seulement un timestamp. Sans capacité, violations et responsabilités sont découvertes tard.

## 3. Objectifs
Montrer source/policy/échéance/temps/état, distinguer les états SLA, auditer pause/reprise et créer une action de suivi.

## 4. Non-objectifs
Ne définit pas toutes les policies, ne calcule pas sans source, ne modifie pas un contrat et ne supprime pas l’historique pendant une pause.

## 5. Propriétaire
Command possède la projection et les mutations autorisées du work item ; Settings ou la source deployment-specific possède policy/engagement.

## 6. Utilisateurs
Principal : SOC Analyst L1/L2. Secondaires : Incident Commander, Service Delivery Manager et Customer Success selon `OPEN-006`.

## 7. Conditions d’entrée
Work item accessible, policy/contrat référencé si applicable et horloge/timezone connues.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Applicable SLA | tenant policy ou engagement | règle référencée | non | version effective | `unknown` ou `not-applicable` expliqué |
| Work timestamps | Incident/Task | création, transitions et pauses | oui | audit courant | calcul impossible visible |
| Pause authorization | policy ou authority source | condition et raison | non | évaluée à chaque pause | pause refusée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | timestamps, owner et état | consulter et modifier de façon limitée |
| SLA policy / engagement | Settings ou source deployment | durée, calendrier et règles de pause | consulter et appliquer en projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | lier SLA, pause/reprise et état dérivé | Command | classe 2, policy vérifiée |
| Task | créer un suivi SLA | Command | classe 2, source et échéance liées |
| SLA policy / engagement | aucune mutation | propriétaire source | projection en lecture seule |

## 11. Fonctionnalités
Afficher état/temps restant, expliquer source/calendrier/pause, montrer at-risk/breached, pause/reprise autorisée et Task/escalade de suivi.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter SLA | lecteur | Incident/Task | 0 | objet accessible | état expliqué | non |
| Mettre en pause | rôle autorisé | tracking | 2 | policy et raison | paused audité | OPEN-013 |
| Reprendre | rôle autorisé | tracking | 2 | pause active | calcul repris | OPEN-013 |
| Créer suivi | coordinateur | Task | 2 | risk ou breach | Task liée | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Calculer l’état SLA | vérification humaine | oui | oui | résumé facultatif | policy versionnée et moteur déterministe |
| Détecter at-risk/breach | oui | oui | oui | non nécessaire | seuils et calendrier de la policy |
| Mettre en pause/reprendre | oui | validation de condition | workflow possible | proposition seulement | action humaine avec raison |
| Créer un suivi | oui | déduplication possible | oui | brouillon Task | création manuelle de Task |

## 14. États fonctionnels
`healthy`, `approaching`, `at-risk`, `breached`, `paused`, `not-applicable`, `unknown`.

## 15. États d’interface
Unknown explique la policy absente ; Partial garde les timestamps ; Offline bloque pause/reprise ; Permission denied masque l’engagement ; Stale montre version/date.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| SLA projection | état dérivé | Queue, Mission Control et Reporting | source, version et calcul visibles |
| SLA event | événement Incident/Task | Audit et Notifications | pause/reprise/breach avec acteur et raison |
| SLA follow-up | Task Command | Work Queue | attribuée, datée et liée au work item |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| SLA at-risk sans owner | action utilisateur ou règle | Work Assignment / Escalation | item, délai et policy | vue SLA Risk restaurée |
| SLA contractuel | ouverture du contexte client | Customers and Delivery | tenant, client, engagement et item | uniquement si module activé |
| SLA breach | création de suivi | Task Coordination | item, breach time et owner | Work Queue SLA Risk |

## 18. Dépendances
CAP-CMD-101, CAP-CMD-102, CAP-CMD-110, Notification Center, Metrics Engine et `OPEN-006` pour le contractuel.

## 19. Source de vérité
Policy/engagement : source propriétaire. Timestamps et mutations work item : Command. Calcul : policy version/calendrier/pauses visibles.

## 20. Provenance et audit
Pause/reprise/breach : acteur, policy/version, before/after, raison, tenant, timezone et correlation ID.

## 21. Permissions fonctionnelles
Lecture Command, coordinate/manage Incident/Task et lecture policy ; permissions atomiques de pause et `OPEN-013` reportées.

## 22. Limites et erreurs
Policy absente, timezone incohérente, horloge stale, pause interdite, conflit ou refus ne produisent pas un SLA certain.

## 23. Métriques
Items calculables, durée at-risk sans suivi et pauses sans raison complète ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; SLA contractuel deployment-dependent sous `OPEN-006`.

## 25. Critères d’acceptation
**Given** policy et timestamps valides, **When** SLA est consulté, **Then** source, version, échéance, temps et état sont explicables.

**Given** une pause non permise, **When** elle est demandée, **Then** la mutation est refusée sans altérer l’historique.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** policy, moteur et actions manuelles suffisent.

## 26. Questions ouvertes
Quelles sources sont universelles et quelles pauses requièrent Govern ? — Requirement IDs ci-dessus ; `OPEN-006` et `OPEN-013` restent ouvertes.

## 27. Consommateurs documentaires
Work Queue SLA Risk, Mission Control, Customers proposal, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.