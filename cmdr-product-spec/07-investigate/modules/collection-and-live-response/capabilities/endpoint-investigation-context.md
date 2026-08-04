---
id: CAP-INV-201
title: Endpoint Investigation Context
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
  - REQ-OBJ-008
open_decisions:
  - OPEN-008
source-of-truth: canonical
---
# CAP-INV-201 — Endpoint Investigation Context

## 1. Définition
Ouvrir un Endpoint dans le contexte d’un Case et comprendre son identité, sa fraîcheur, ses capacités et ses limites sans administrer la Fleet.

## 2. Problème utilisateur
Sans ce contexte, les états offline, partiels ou unsupported peuvent être pris pour un succès et les frontières Fleet, Agent, Case et Govern deviennent ambiguës.

## 3. Objectifs
- relier Case, Endpoint, Agent, Fleet, Policy, Incidents, Artifacts et projections Govern ;
- exposer statut, dernière communication, capacités, limitations et fraîcheur ;
- préparer une collecte ou Live Session en conservant le return origin.

## 4. Non-objectifs
Ne pas administrer Fleet/Policies/enrollment, définir protocole/API/commande/moteur/plateforme, créer automatiquement Evidence/Finding/Decision/Run/Result ni commencer Analysis Workbench.

## 5. Propriétaire
Investigate possède le contexte métier et les relations au Case. Platform Settings administre Fleet/Policies ; Endpoint Agent expose son état et exécute localement ; Govern possède l’autorité risquée.

## 6. Utilisateurs
Principal : Case Analyst. Secondaires : Investigation Lead, Evidence Reviewer, Incident Commander, Response Operator et Platform Administrator en consultation.

## 7. Conditions d’entrée
Tenant/environnement conservés, Case accessible, Endpoint résolu, projections de fraîcheur/capacités/policy disponibles ou lacunes explicites, permission de lecture.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Case et objectif | Investigate | contexte métier | oui | version courante | contexte standalone interdit sauf raison |
| Endpoint et état Agent | Platform Settings / Endpoint Agent | cible et disponibilité | oui | dernière communication visible | offline/unknown, aucune exécution présentée |
| Fleet et Policy | Platform Settings | projection administrative | non | version et affectation visibles | partially-available |
| Incidents, Artifacts et Runs liés | owners respectifs | contexte relationnel | non | version/statut visibles | section partielle |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Case | Investigate | objectif, scope, owner et relations | consulter et naviguer |
| Endpoint | concept partagé à formaliser | identité fonctionnelle et environnement | consulter et sélectionner |
| Endpoint Agent | Endpoint Agent | état, fraîcheur et capacités déclarées | consulter uniquement |
| Endpoint Agent Fleet / Endpoint Policy | Platform Settings | posture, affectation et restrictions | consulter uniquement |
| Incident / Task | Command | contexte et coordination | consulter et lier |
| Artifact / Decision / Response Run / Result | Investigate / Govern | résultats et historique liés | consulter et naviguer |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint context projection | assembler sans persister comme source concurrente | Investigate | chaque champ conserve owner, version et fraîcheur |
| Case–Endpoint relation | créer ou supersede | Investigate / Object Linking | relation sourcée, aucun transfert d’ownership |
| Événement d’ouverture/pivot | émettre | Shared mechanisms, sémantique Investigate | acteur, Case, Endpoint et return origin |
| Fleet/Policy/Agent | aucune mutation | owner externe | lecture seule absolue |

## 11. Fonctionnalités
Rechercher ou ouvrir un Endpoint depuis un Case, afficher tenant/environnement/identité/Agent/fraîcheur/capacités/limitations/policies/opérations récentes/Artifacts/Incidents/Cases/Decisions/Runs, puis préparer collecte ou session.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter ou filtrer | utilisateur autorisé | Endpoint context | 0 | Endpoint read | projection sourcée | non |
| Lier au Case | Case Analyst | relation | 2 | Case et Endpoint accessibles | lien versionné | OPEN-013 selon policy |
| Préparer collecte | Analyst | Collection Request draft | 1 | capacité/policy visibles | passage CAP-INV-202 | selon impact |
| Préparer Live Session | Response Operator | session request | 2 | disponibilité/permission | passage CAP-INV-209 | OPEN-007/013 selon policy |
| Préparer containment | Investigation Lead | Action Request | 3 | Finding/Evidence/impact | passage CAP-INV-215/113 | obligatoire |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Assembler les projections | oui | oui, IDs et relations | oui | non nécessaire | linking et context bar |
| Résumer l’état | oui | agrégation factuelle | oui | résumé attribué | lecture des champs |
| Signaler staleness/capacité absente | oui | oui | oui | explication | règles de fraîcheur |
| Proposer prochaine action | oui | procédures | oui | suggestion | actions et checklists |

Toute sortie automatisée expose initiateur, producteur/version, Automation Run/Tool Calls s’ils existent, sources, timestamp, limites, owner humain, disposition et trace.

## 14. États fonctionnels
`available`, `partially-available`, `degraded`, `offline`, `stale`, `unsupported`, `restricted`, `policy-blocked`, `agent-not-installed`, `unknown`. Ce ne sont pas des états finaux d’objet.

## 15. États d’interface
Loading conserve Case/Endpoint ; Empty distingue Agent absent et permission ; Partial nomme les projections manquantes ; Error conserve les données valides ; Offline interdit toute fausse exécution ; Permission denied ne fuit rien ; Stale affiche la dernière synchronisation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Endpoint context snapshot | projection | Case Workspace | owners, fraîcheur et limitations visibles |
| Collection/session draft context | package de transition | CAP-INV-202/209 | Case, Endpoint, capacités, policy, permission et return origin |
| Case–Endpoint relation | relation | Case et Command projection | sourcée et tenant-scoped |
| Context activity event | événement | Timeline/Trace/Replay | acteur, source et timestamp |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Case Workspace | ouvrir Endpoint | Endpoint Investigation Context | tenant, environnement, Case, Incident, Endpoint, objectif, return origin | même Case |
| Endpoint Context | préparer collecte | CAP-INV-202 | target, Case, capacités, policy, permission, classe | Endpoint Context |
| Endpoint Context | ouvrir session | CAP-INV-209 | target, Case, raison, participants, expiration, policy | Endpoint Context |
| Endpoint Context | ouvrir Artifact/Run | owner destination | objet, source et return origin | même contexte Endpoint |

## 18. Dépendances
CAP-INV-102/105/110/112, Platform Settings Fleet/Policies/Health, Endpoint Agent health/capabilities, Shared Linking/Context Bar/Inspector/Timeline/Trace, Govern projections et OPEN-008/013.

## 19. Source de vérité
Case et relations locales restent Investigate ; Fleet/Policy restent Settings ; état/capacités Agent restent Endpoint Agent ; Incident reste Command ; Decision/Run/Result restent Govern.

## 20. Provenance et audit
Case, Endpoint, Agent, Fleet, Policy/version, fraîcheur, utilisateur, source d’entrée, liens ouverts, actions préparées et correlation ID.

## 21. Permissions fonctionnelles
Endpoint read, Agent capability read, Fleet/Policy projection read, sensitive context read, cross-tenant/environment restrictions, Case link et préparation de collecte/session. Matrice atomique reportée.

## 22. Limites et erreurs
Endpoint ambigu, Agent absent/revoked/degraded, heartbeat stale, platform unsupported, policy conflictuelle, données protégées, tenant mismatch ou relation cassée produisent un état explicite sans mutation externe.

## 23. Métriques
Contextes avec fraîcheur/capacités complètes, endpoints unsupported/offline, transitions vers collecte/session, liens cassés et retours Case réussis.

## 24. Classification de livraison
`defined` / `planned`. Cible native via Endpoint Agent, mais aucune plateforme, API, protocole, moteur ou release n’est prouvée. OPEN-008 reste ouverte.

## 25. Critères d’acceptation
**Given** un Case actif et un Endpoint disponible **When** le contexte est ouvert **Then** identité, tenant, environnement, fraîcheur, capacités, policy et liens sont visibles sans administration.

**Given** un Endpoint offline **When** une collecte est préparée **Then** aucune exécution n’est montrée, l’attente éventuelle et l’annulation sont explicites et le Case reste actif.

**Given** aucun modèle IA **When** le contexte est utilisé **Then** toutes les projections, filtres et transitions essentielles restent disponibles déterministiquement.

## 26. Questions ouvertes
OPEN-008 conserve le support plateforme. L’objet Endpoint, les cardinalités et permissions fines restent aux phases Objets/Permissions ; OPEN-013 reste ouverte pour les mutations réversibles.

## 27. Consommateurs documentaires
Collection and Live Response, Case Workspace, Screen Capability Map, Platform Settings Fleet/Policies/Health, Endpoint Agent, parcours Endpoint investigation/offline recovery et phases Objets/Permissions.