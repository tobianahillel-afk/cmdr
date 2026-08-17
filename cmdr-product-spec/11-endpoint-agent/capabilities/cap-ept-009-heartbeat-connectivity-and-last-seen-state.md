---
id: CAP-EPT-009
title: Heartbeat, Connectivity and Last-Seen State
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
# CAP-EPT-009 — Heartbeat, Connectivity and Last-Seen State

## 1. Définition
Define source-backed heartbeat observations, connectivity state and last-seen semantics as separate Endpoint facts, including delay, missed-heartbeat and reconnect context without equating any one signal with full health.

## 2. Problème utilisateur
Last seen, a heartbeat event and connectivity are related but not identical; conflating them can incorrectly mark an Agent healthy, offline or current.

## 3. Objectifs
Heartbeat event/time; expected-interval concept when sourced; last seen; connected/disconnected/unknown; delay; missed-heartbeat; reconnect; freshness and provenance.

## 4. Non-objectifs
No transport protocol, persistent connection implementation, port, heartbeat wire format, network telemetry, health proof, SLA target or response action.

## 5. Propriétaire
Endpoint Agent owns local technical heartbeat/connectivity/last-seen semantics. Settings may aggregate projections in Fleet.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, SOC/Investigate analyst, Command/Govern consumer and Auditor.

## 7. Conditions d’entrée
Agent identity and one or more source-backed communication/observation timestamps. Expected interval must come from a source; otherwise no missed-heartbeat threshold is invented.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| heartbeat observation | Endpoint/control-plane observation | heartbeat fact | non | event time | heartbeat unknown |
| last-seen observation | Endpoint/control-plane source | presence timestamp | oui for last-seen | explicit time | last seen unknown |
| connectivity observation | Endpoint/control-plane source | connectivity fact | non | explicit time | connectivity unknown |
| expected interval | source-owned configuration | timing context | non | configured version | no missed threshold inferred |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | identity/current state | existing read |
| endpoint-agent-fleet | Platform Settings | aggregation context only | Settings read projection |
| environment | Platform Settings | scope | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Heartbeat/Connectivity State | derive/refresh | Endpoint Agent | distinct from health and Fleet state |
| local-audit-event | append state/reconnect provenance | Endpoint Agent | no transport payload required |

## 11. Fonctionnalités
Heartbeat timestamp; expected interval when sourced; delay/missed-heartbeat; last seen; connected/disconnected/unknown; reconnect; freshness and reason.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect heartbeat/last seen | authorized user | endpoint-agent | 0 | read | distinct timestamps/states | non |
| derive connectivity/freshness | Endpoint Operator | observations | 1 | source facts | deterministic state or unknown | non |
| request bounded connectivity check | authorized operator | endpoint-agent | 2 | source supports it | check request only | no response execution |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| timestamp inspection | oui | oui | oui | summary only | direct event review |
| missed/reconnect derivation | oui | oui when interval exists | oui | explanation only | deterministic timing rules |
| invent heartbeat/connectivity | non | non | non | interdit | unknown state |

## 14. États fonctionnels
`connected`, `disconnected`, `connectivity-unknown`, `heartbeat-current`, `heartbeat-missed`, `last-seen-known`, `last-seen-unknown`, `reconnected`.

## 15. États d’interface
Offline UI state is not automatically a business/Agent lifecycle transition; Stale displays event age; Partial differentiates missing heartbeat vs missing connectivity; Permission denied masks restricted detail. No screen design.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Heartbeat Fact | Endpoint technical fact | health/state/Settings Fleet | timestamp/source explicit |
| Connectivity/Last-Seen Projection | derived/projection | Settings/Command/Investigate/Govern | signals remain distinct |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| communication observations | event observed | heartbeat/last-seen facts | source/time | origin retained |
| observations + interval | timing evaluation | current/missed/unknown | interval/source/time | no health conclusion |
| disconnected/missed | new valid observation | reconnected/current | previous/current timestamps | history retained |

## 18. Dépendances
CAP-EPT-001/003/008/010, Endpoint health monitoring, Settings Fleet projection and OPEN-008.

## 19. Source de vérité
Individual heartbeat/connectivity/last-seen observations retain their sources; Endpoint owns their local technical interpretation, not a global availability claim.

## 20. Provenance et audit
Record Agent/scope, source, event/receive times, interval source if any, derived reason, previous state, reconnect and correlation id.

## 21. Permissions fonctionnelles
Endpoint read covers ordinary state; restricted connectivity metadata remains permission-gated. No new atomic permission or cross-tenant visibility is implied.

## 22. Limites et erreurs
Last seen != heartbeat; heartbeat received != healthy; no heartbeat != immediately offline without an applicable state rule; unknown interval means missed status cannot be invented.

## 23. Métriques
Heartbeat current/missed/unknown counts; connectivity states; reconnect counts; descriptive delay/last-seen age; stale/permission-denied observations.

## 24. Classification de livraison
`draft / defined / planned`; no transport, port, wire protocol, connectivity implementation or support claim.

## 25. Critères d’acceptation
**Given** last seen is recent but no heartbeat source exists, **When** state is viewed, **Then** last seen is known while heartbeat remains unknown.

**Given** heartbeat is received but health self-check is degraded, **When** projected, **Then** connectivity can be current while health remains degraded.

**Given** no expected interval is configured, **When** heartbeat ages, **Then** the system does not invent a missed threshold.

## 26. Questions ouvertes
OPEN-008 remains open; platform-specific heartbeat sources and intervals are not selected.

## 27. Consommateurs documentaires
Endpoint health/operational state, Settings Fleet, Investigate/Command/Govern, registers, Quality, Roadmap and future resilience/telemetry lots.