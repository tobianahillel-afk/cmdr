---
id: CAP-EPT-010
title: Degraded, Offline and Stale State Semantics
product: endpoint-agent
module: foundations
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-OBJ-008, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-010 — Degraded, Offline and Stale State Semantics

## 1. Définition
Define distinct Endpoint operational-state semantics for online, degraded, offline, stale, unknown, revoked and enrollment-pending using explicit evidence and no silent state upgrade.

## 2. Problème utilisateur
Health, connectivity, freshness, revocation and enrollment signals can diverge; a single overloaded status can hide partial capability or incorrectly treat stale data as current.

## 3. Objectifs
Explicit state definitions; conceptual triggers/evidence; transitions; recovery observation; consumer projection; no silent promotion; preserve revoked/enrollment-pending distinctions.

## 4. Non-objectifs
No recovery engine, retry/buffering implementation, response rollback, containment, availability SLA, universal timeout or screen-state design.

## 5. Propriétaire
Endpoint Agent owns individual local operational state. Settings owns Fleet administration; Govern owns response authority/rollback.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, SOC/Investigate analyst, Command/Govern consumer and Auditor.

## 7. Conditions d’entrée
Agent identity plus available health, heartbeat/connectivity, freshness, registration/enrollment and revocation evidence. Missing evidence yields unknown/partial rather than inference.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| health assessment | CAP-EPT-008 | technical state | non | health freshness | no health conclusion |
| heartbeat/connectivity | CAP-EPT-009 | connectivity facts | non | observation freshness | connectivity unknown |
| inventory freshness | CAP-EPT-007 | freshness context | non | snapshot freshness | no stale inventory conclusion |
| registration/enrollment/revocation | CAP-EPT-001/002 | lifecycle context | oui | current refs | operational state constrained/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | canonical lifecycle + local facts | existing read |
| endpoint-agent-fleet | Platform Settings | aggregation context | Settings read projection |
| local-audit-event | Endpoint Agent | transition provenance | read where authorized |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint Operational State assessment | derive/transition | Endpoint Agent | evidence-backed, no Fleet lifecycle mutation |
| local-audit-event | append state transition | Endpoint Agent | prior/new evidence retained |

## 11. Fonctionnalités
Online/degraded/offline/stale/unknown/revoked/enrollment-pending definitions; evidence and reason; transition/recovery observation; subset capability context; no-silent-state-upgrade.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect operational state | authorized user | endpoint-agent | 0 | read | state/evidence/reason | non |
| derive/reconcile state | Endpoint Operator | source facts | 1 | source observations | deterministic state or unknown | non |
| request bounded refresh/self-check | authorized operator | endpoint-agent | 2 | source supports | request only | no response action |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect state evidence | oui | oui | oui | sourced explanation | source matrix |
| derive state | oui | oui under explicit rules | oui | explanation only | deterministic rule set |
| hide/promote state without evidence | non | non | non | interdit | explicit unknown/stale/degraded |

## 14. États fonctionnels
`online`, `degraded`, `offline`, `stale`, `unknown`, `revoked`, `enrollment-pending`.

## 15. États d’interface
Interface Offline/Stale/Partial states must not overwrite canonical/functional semantics; they show evidence/freshness and preserve return origin. No Endpoint Screen ID or detailed UX is created.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Endpoint Operational State | derived Endpoint fact | Settings Fleet, Command, Investigate, Govern | state/evidence/freshness explicit |
| State Transition Provenance | local-audit-event | Audit | prior/new state and reason retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| health/connectivity/freshness | evidence changes | operational state | source refs/reasons/times | source context retained |
| degraded/offline/stale | recovery evidence | newer state | recovery observation + prior state | history retained |
| any state | revocation | revoked | revocation ref/time | no deletion |

## 18. Dépendances
CAP-EPT-001/002/007/008/009/011, Endpoint canonical object states, Settings Fleet projection, OPEN-008.

## 19. Source de vérité
Endpoint owns individual operational-state assessment from named source facts; Settings Fleet aggregation and UI state remain distinct.

## 20. Provenance et audit
Record source fact refs/times, rule/context, prior/new state, reason, missing evidence, actor/engine and correlation id; never rewrite history.

## 21. Permissions fonctionnelles
Endpoint read for state/evidence; restricted underlying facts remain masked. State projection does not grant access to source data or cross-tenant scope.

## 22. Limites et erreurs
Degraded != offline; offline != revoked; stale != offline automatically; online != fully capable; enrollment-pending != online; unknown cannot be silently promoted.

## 23. Métriques
Operational-state distribution; transition counts; duration as descriptive observation; unknown/stale cases; subset-capability degradation and recovery observations.

## 24. Classification de livraison
`draft / defined / planned`; no recovery/queue implementation, timeout/SLO, response action, protocol or screen design.

## 25. Critères d’acceptation
**Given** one required component is degraded but connectivity persists, **When** state is derived, **Then** Agent may be degraded rather than offline.

**Given** inventory is stale while current heartbeat exists, **When** state is read, **Then** stale data is exposed without automatically marking the Agent offline.

**Given** an offline Agent is revoked, **When** revocation evidence arrives, **Then** revoked remains distinct and history is retained.

## 26. Questions ouvertes
OPEN-008 remains open; platform-specific rules/timing and later resilience/recovery behavior are not selected.

## 27. Consommateurs documentaires
Endpoint health/capabilities, Settings Fleet, Investigate/Command/Govern, registers, Quality, Roadmap and future EPT-2/EPT-6 work.