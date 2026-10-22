---
id: CAP-EPT-008
title: Agent Health and Self-Check Assessment
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
# CAP-EPT-008 — Agent Health and Self-Check Assessment

## 1. Définition
Define a bounded Endpoint Agent health assessment from self-checks, component/sensor health and explicit local dependency/resource conditions, without treating health as security posture or a successful self-check as proof of endpoint integrity.

## 2. Problème utilisateur
An Agent can answer a heartbeat while sensors or dependencies are degraded, and a self-check can pass without proving the host is secure or untampered.

## 3. Objectifs
Represent Agent health, self-checks, component/sensor health projection, degraded capabilities, unavailable dependencies, privilege/resource issues, timestamp, freshness, reason and provenance.

## 4. Non-objectifs
No anti-tamper proof, endpoint security score, detection verdict, remediation, recovery implementation, resource-tuning engine, response action or final health SLO.

## 5. Propriétaire
Endpoint Agent owns local technical health/self-check semantics. Settings may aggregate administrative health projections but does not own the individual local health fact.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, SOC/Investigate analyst, Command/Govern consumer and Auditor.

## 7. Conditions d’entrée
Agent identity, tenant/environment binding, sourced self-check/component observations and timestamps. Missing dependencies remain explicit.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| self-check observations | Endpoint local health monitor | health facts | oui | explicit timestamp | health partial/unknown |
| component/sensor status | Endpoint components | technical projection | selon platform/source | explicit timestamp | affected capability unknown |
| dependency/resource conditions | Endpoint local state | context facts | non | observation freshness | no failure inferred |
| heartbeat context | CAP-EPT-009 | connectivity context | non | heartbeat freshness | health still independently assessed |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | identity/state/version | existing read |
| environment | Platform Settings | scope/context | read projection |
| endpoint-policy | Platform Settings | effective-state reference only if relevant | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Agent Health Assessment | derive/refresh/supersede | Endpoint Agent | local technical fact, not security posture |
| local-audit-event | append health transition/check provenance | Endpoint Agent | append-only trace |

## 11. Fonctionnalités
Self-check assessment; sensor/component health; degraded-capability linkage; dependency unavailable; privilege/resource issue; health reason; timestamp/freshness; explicit unknown/partial state.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect health | authorized user | Agent Health Assessment | 0 | read | health/reason/source | non |
| normalize health facts | Endpoint Operator | health facts | 1 | sourced observations | healthy/degraded/unknown assessment | non |
| request bounded self-check | authorized operator | endpoint-agent | 2 | source supports + manage need | self-check request only | no response authority |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/normalize health | oui | oui | oui | sourced explanation | deterministic state rules |
| summarize degraded causes | oui | oui for sourced causes | oui | oui, attributed | structured reason list |
| claim security/integrity proof | non | non | non | interdit | explicit bounded health semantics |

## 14. États fonctionnels
`healthy-for-declared-checks`, `degraded`, `health-unknown`, `dependency-unavailable`, `self-check-failed`, `health-stale`.

## 15. États d’interface
Partial lists missing checks; Offline does not automatically equal unhealthy; Stale shows last assessment age; Permission denied masks restricted causes; Error preserves prior valid health with timestamp. No Screen ID.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Agent Health Projection | Endpoint technical fact | Settings Fleet, Command, Investigate, Govern | source/check scope/reason explicit |
| Degraded Capability Context | derived linkage | capability availability consumers | degradation does not imply offline |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| self-check/components | assessment | Agent Health Assessment | check results/source/time | local context retained |
| health assessment | degraded dependency | capability availability context | affected capability/reason | health context retained |
| health assessment | consumer projection | Settings/products | status/freshness/limits | no ownership transfer |

## 18. Dépendances
Endpoint health-monitoring, CAP-EPT-001/003/005/009/011 context, Settings projections, Security/privacy, OPEN-008.

## 19. Source de vérité
Endpoint local health sources are authoritative only for declared technical checks; health is not security posture, compliance, integrity proof or business readiness.

## 20. Provenance et audit
Record check identifiers/versions where known, Agent/scope, observations, source/time, normalization reason, affected capabilities, requester and correlation id.

## 21. Permissions fonctionnelles
Endpoint Agent read for health; self-check request is a functional manage need if supported. Sensitive causes remain masked per Security; no final atomic permission is created.

## 22. Limites et erreurs
Heartbeat received != healthy; health != security posture; self-check PASS != full endpoint integrity proof; missing checks yield partial/unknown, not healthy.

## 23. Métriques
Health-state distribution; failed/missing/stale checks; degraded-capability count; dependency/resource issue counts; self-check request outcomes.

## 24. Classification de livraison
`draft / defined / planned`; no health engine implementation, security posture product, remediation, API/protocol or platform support claim.

## 25. Critères d’acceptation
**Given** heartbeat is received while one sensor self-check fails, **When** health is assessed, **Then** the Agent is degraded rather than automatically healthy.

**Given** every declared self-check passes, **When** status is displayed, **Then** PASS is bounded to those checks and is not an endpoint integrity proof.

**Given** AI is unavailable, **When** health facts are normalized, **Then** deterministic/manual rules still produce the assessment.

## 26. Questions ouvertes
OPEN-008 remains open; platform-specific checks, launch coverage and future tenant policy parameters are not selected.

## 27. Consommateurs documentaires
Endpoint operational state/capability availability, Settings Fleet/health, Investigate/Command/Govern, registers, Quality, Roadmap and future resilience/security lots.