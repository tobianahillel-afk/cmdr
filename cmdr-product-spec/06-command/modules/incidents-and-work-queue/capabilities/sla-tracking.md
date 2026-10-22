---
id: CAP-CMD-105
title: SLA Tracking
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-14
requirement_ids:
  - REQ-PROD-005
  - REQ-PROD-013
  - REQ-PROD-053
  - REQ-OBJ-012
open_decisions:
  - OPEN-013
source-of-truth: canonical
---

# CAP-CMD-105 — SLA Tracking

## 1. Définition
Projette le SLA applicable à un Incident/Task, le temps restant, le risque, les pauses autorisées et la provenance tenant/contractuelle, sans définir la policy ou muter le contrat.

## 2. Problème utilisateur
Les équipes doivent connaître une échéance et la validité d’une pause, pas seulement un timestamp. Sans capacité, violations et responsabilités sont découvertes tard.

## 3. Objectifs
Montrer source/policy/échéance/temps/état, distinguer les états SLA, auditer pause/reprise et créer une action de suivi.

## 4. Non-objectifs
Ne définit pas toutes les policies, ne calcule pas sans source, ne modifie pas un contrat, ne crée pas Customer et ne supprime pas l’historique pendant une pause.

## 5. Propriétaire
Command possède la projection et les mutations autorisées du work item ; Settings ou la source deployment/customer/contract possède policy/engagement. ADR-0008 fixe seulement la sémantique de déploiement et de projection Customer.

## 6. Utilisateurs
Principal : SOC Analyst L1/L2. Secondaires : Incident Commander, Service Delivery Manager et Customer Success lorsque le deployment et les permissions applicables l'autorisent selon ADR-0008.

## 7. Conditions d’entrée
Work item accessible, Tenant sélectionné, policy/contrat référencé si applicable et horloge/timezone connues. Dans un contexte MSSP agrégé read-only, pause/reprise/création de suivi sont indisponibles tant qu'un Tenant unique n'est pas sélectionné.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Applicable SLA | tenant policy ou engagement externe | règle référencée | non | version effective | `unknown` ou `not-applicable` expliqué |
| Work timestamps | Incident/Task | création, transitions et pauses | oui | audit courant | calcul impossible visible |
| Pause authorization | policy ou authority source | condition et raison | non | évaluée à chaque pause | pause refusée |
| Tenant context | Platform Settings/Security | Tenant sélectionné et autorisé | oui pour mutation | courant | mutation bloquée |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Incident / Task | Command | timestamps, owner et état | consulter et modifier de façon limitée |
| SLA policy / engagement | Settings ou source deployment/customer/contract | durée, calendrier et règles de pause | consulter/appliquer en projection |
| Customer context | source externe | référence d'engagement si présente | read-only, jamais objet canonique |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Incident / Task | lier SLA, pause/reprise et état dérivé | Command | Class 2, Tenant sélectionné, policy vérifiée |
| Task | créer un suivi SLA | Command | Class 2, Tenant sélectionné, source et échéance liées |
| SLA policy / engagement | aucune mutation | propriétaire source | projection en lecture seule |
| Customer / contract source | aucune mutation | source externe | projection en lecture seule |

## 11. Fonctionnalités
Afficher état/temps restant, expliquer source/calendrier/pause, montrer at-risk/breached, pause/reprise autorisée et Task/escalade de suivi. Le contractual SLA reste deployment-dependent et peut être absent en mode Internal.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter SLA | lecteur | Incident/Task | 0 | objet accessible | état expliqué | non |
| Mettre en pause | rôle autorisé | tracking | 2 | Tenant sélectionné, policy et raison | paused audité | OPEN-013 |
| Reprendre | rôle autorisé | tracking | 2 | Tenant sélectionné, pause active | calcul repris | OPEN-013 |
| Créer suivi | coordinateur | Task | 2 | Tenant sélectionné, risk ou breach | Task liée | OPEN-013 selon default C2 |

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
Unknown explique la policy absente ; Partial garde les timestamps ; Offline bloque pause/reprise ; Permission denied masque l’engagement ; Stale montre version/date. Un overview MSSP agrégé est read-only.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| SLA projection | état dérivé | Queue, Mission Control, Customers & Delivery et Reporting | source, version, Tenant et calcul visibles |
| SLA event | événement Incident/Task | Audit et Notifications | pause/reprise/breach avec acteur et raison |
| SLA follow-up | Task Command | Work Queue | attribuée, datée, tenant-scoped et liée au work item |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| SLA at-risk sans owner | action utilisateur ou règle | Work Assignment / Escalation | item, délai et policy | vue SLA Risk restaurée |
| SLA contractuel | ouverture du contexte delivery | Customers and Delivery | Tenant, customer projection, engagement et item | uniquement si module activé |
| SLA breach | création de suivi | Task Coordination | Tenant, item, breach time et owner | Work Queue SLA Risk |

## 18. Dépendances
CAP-CMD-101, CAP-CMD-102, CAP-CMD-110, CAP-CMD-401, ADR-0008, Notification Center et Metrics Engine.

## 19. Source de vérité
Policy/engagement : source propriétaire. Customer : projection externe seulement. Timestamps et mutations work item : Command. Calcul : policy version/calendrier/pauses visibles.

## 20. Provenance et audit
Pause/reprise/breach : acteur, Tenant, policy/version, before/after, raison, timezone et correlation ID.

## 21. Permissions fonctionnelles
Lecture Command, coordinate/manage Incident/Task et lecture policy ; permissions atomiques de pause et `OPEN-013` restent reportées. Aucun droit multi-tenant read-only ne confère une mutation SLA.

## 22. Limites et erreurs
Policy absente, timezone incohérente, horloge stale, Tenant non sélectionné, pause interdite, conflit ou refus ne produisent pas un SLA certain ou une mutation implicite.

## 23. Métriques
Items calculables, durée at-risk sans suivi, pauses sans raison complète et refus de mutation depuis contexte agrégé ; aucune cible définitive.

## 24. Classification de livraison
`defined / planned`, cible native ; SLA contractuel deployment-dependent selon ADR-0008. Aucune source contractuelle ou implémentation n'est revendiquée universellement.

## 25. Critères d’acceptation
**Given** policy et timestamps valides, **When** SLA est consulté, **Then** source, version, Tenant, échéance, temps et état sont explicables.

**Given** une pause non permise ou aucun Tenant unique sélectionné, **When** elle est demandée, **Then** la mutation est refusée sans altérer l’historique.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** policy, moteur et actions manuelles suffisent.

## 26. Questions ouvertes
Quelles pauses requièrent Govern et quelle politique C2 s'applique par défaut ? — `OPEN-013` reste ouverte. `OPEN-006` est résolue par ADR-0008.

## 27. Consommateurs documentaires
Work Queue SLA Risk, Mission Control, Customers & Delivery, Reporting, parcours, écrans et phases d'implémentation ultérieures.
