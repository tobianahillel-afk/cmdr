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
Projette le SLA applicable à un Incident/Task, temps restant, risque, pauses autorisées et provenance contractuelle/tenant, sans définir la policy.

## 2. Problème utilisateur
Les équipes doivent connaître une échéance et la validité d’une pause, pas seulement un timestamp. Rôles : SOC Analyst L1/L2, Incident Commander, Service Delivery Manager, Customer Success selon déploiement. Sans capacité, violations et responsabilités sont découvertes tard.

## 3. Objectifs
Montrer source/policy/échéance/temps/état ; distinguer les états SLA ; auditer pause/reprise ; créer une action de suivi sans inventer la policy.

## 4. Non-objectifs
Ne pas définir toutes les policies, calculer sans source, modifier un contrat ou supprimer l’historique pendant une pause.

## 5. Propriétaire
Command possède la projection et les mutations autorisées du work item ; Settings ou la source deployment-specific possède policy/engagement.

## 6. Utilisateurs
Principal : SOC Analyst L1/L2. Secondaires : Incident Commander, Service Delivery Manager, Customer Success selon `OPEN-006`.

## 7. Conditions d’entrée
Work item accessible ; policy/contrat référencé si applicable ; horloge/timezone connues.

## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Applicable SLA | tenant policy/engagement | règle référencée | non | version effective | unknown/not-applicable expliqué |
| Work timestamps | Incident/Task | création/transitions/pauses | oui | audit courant | calcul impossible visible |
| Pause authorization | policy/authority | condition/raison | non | chaque pause | pause interdite |

## 9. Objets lus
| Objet | Owner | Projection | Droit local |
|---|---|---|---|
| Incident / Task | Command | timestamps, owner, state | lecture/modification limitée |
| SLA policy / engagement | Settings/deployment source | durée, calendrier, pause | projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Owner | Règle |
|---|---|---|---|
| Incident / Task | référence SLA, pause/reprise, follow-up | Command | classe 2, policy vérifiée |
| Task | action de suivi SLA | Command | classe 2 |

## 11. Fonctionnalités
Afficher état/temps restant ; expliquer source/calendrier/pause ; montrer at-risk/breached dans la même file ; pause/reprise autorisée ; Task/escalade de suivi.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter SLA | lecteur | Incident/Task | 0 | objet accessible | état expliqué | non |
| Mettre en pause | rôle autorisé | tracking | 2 | policy et raison | paused audité | OPEN-013 |
| Reprendre | rôle autorisé | tracking | 2 | pause active | calcul repris | OPEN-013 |
| Créer suivi | coordinateur | Task | 2 | risk/breach | Task liée | non |

## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur | Workflow | Agent | Govern | Sans IA |
|---|---|---|---|---|---|---|---|
| Calculer SLA | validation | policy | calcul déterministe | possible | résumé seulement | non | policy+moteur |
| Pause/reprise | décision humaine | condition | validation | possible | proposition | OPEN-013 | action manuelle |

## 14. États fonctionnels
`healthy`, `approaching`, `at-risk`, `breached`, `paused`, `not-applicable`, `unknown`.

## 15. États d’interface
Unknown explique la policy absente ; Partial garde timestamps valides ; Offline bloque pause/reprise ; Permission denied ne révèle pas l’engagement ; Stale montre version/date. Rendu DS.

## 16. Sorties
| Sortie | Objet/événement | Consommateur | Garantie |
|---|---|---|---|
| SLA projection | derived status | Queue, Mission Control, reporting | source/version/calcul visibles |
| SLA event | Incident/Task event | audit/notifications | pause/reprise/breach avec raison |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte | Retour |
|---|---|---|---|---|
| SLA at risk sans owner | affectation/escalade | Command | item, délai, policy | vue SLA Risk |
| contractual SLA | contexte client | Customers and Delivery | tenant/client | seulement si activé |

## 18. Dépendances
CAP-CMD-101, 102, 110, Notification Center, Metrics Engine, `OPEN-006` pour contractual SLA.

## 19. Source de vérité
Policy/engagement vient de sa source ; timestamps et mutations work item sont Command ; calcul expose policy version, calendrier et pauses.

## 20. Provenance et audit
Pause/reprise/breach enregistrent acteur, policy/version, avant/après, raison, tenant, timezone et correlation ID.

## 21. Permissions fonctionnelles
Lecture Command, coordinate/manage Incident/Task, lecture policy. Permissions atomiques de pause et `OPEN-013` reportées.

## 22. Limites et erreurs
Policy absente, timezone incohérente, horloge stale, pause non autorisée, conflit de version ou refus ne doivent pas produire un SLA certain.

## 23. Métriques
Items avec source/état calculable ; durée at-risk sans suivi ; pauses sans raison complète. Aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; contractual SLA deployment-dependent (`OPEN-006`).

## 25. Critères d’acceptation
**Given** une policy et des timestamps valides, **When** SLA est consulté, **Then** source, version, échéance, temps et état sont explicables.

**Given** une pause non permise, **When** elle est demandée, **Then** la mutation est refusée avec raison sans altérer l’historique.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** policy, moteur déterministe et actions manuelles suffisent.

## 26. Questions ouvertes
Quelles sources sont universelles vs deployment-dependent ? Quelles pauses requièrent Govern ? — REQ-PROD-005, REQ-PROD-013, REQ-PROD-053, REQ-OBJ-012. `OPEN-006` et `OPEN-013` restent ouvertes.

## 27. Consommateurs documentaires
Work Queue SLA Risk, Mission Control, Customers proposal, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.
