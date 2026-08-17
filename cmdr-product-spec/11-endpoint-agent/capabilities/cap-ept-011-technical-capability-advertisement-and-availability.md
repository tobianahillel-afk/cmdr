---
id: CAP-EPT-011
title: Technical Capability Advertisement and Availability
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
# CAP-EPT-011 — Technical Capability Advertisement and Availability

## 1. Définition
Define how an Endpoint Agent advertises technical capability semantics and current local availability with explicit unsupported, unavailable, degraded and unknown states, without equating advertisement with authorization, global support or a Studio Tool.

## 2. Problème utilisateur
A capability can be known to an Agent yet unavailable because of platform, version, dependency or health conditions, and visibility must never become execution authority.

## 3. Objectifs
Expose capability semantic ref, advertised state, local availability, dependency/platform/version conditions, degradation, permission-requirement projection, freshness and provenance.

## 4. Non-objectifs
No detailed telemetry/detection/Live Response execution capabilities, no Tool definition, no authorization engine, no platform support declaration, no provider implementation or final schema.

## 5. Propriétaire
Endpoint Agent owns local technical capability advertisement/availability. Settings may aggregate capability inventory; Studio owns Tools; Security/Govern retain permissions/authority.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, Investigate/SOC analyst, Govern/Command consumer, Studio integrator as consumer and Auditor.

## 7. Conditions d’entrée
Agent identity, version/platform/health context, a defined capability semantic reference and any known local dependency/permission requirement projections.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| capability semantic reference | Endpoint capability documentation/runtime metadata | semantic ref | oui | versioned | capability unknown |
| platform/version context | CAP-EPT-004/005 | eligibility facts | non | source freshness | condition unknown |
| health/dependency state | CAP-EPT-008 | availability context | non | health freshness | availability unknown |
| permission requirement projection | Security/source owner | functional requirement | non | source version | authorization not inferred |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | identity/version/health state | existing read |
| endpoint-agent-fleet | Platform Settings | aggregation context only | Settings read projection |
| environment | Platform Settings | scope/context | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical Capability Advertisement | derive/refresh | Endpoint Agent | local technical metadata, not Studio Tool |
| local-audit-event | append availability change/provenance | Endpoint Agent | no permission grant |

## 11. Fonctionnalités
Capability ref; advertised/not-advertised; available/unavailable/degraded/unsupported/unknown; dependency/platform/version conditions; permission-requirement projection; freshness/reason/provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect capability advertisement | authorized user | advertisement | 0 | read | state/reason/conditions | non |
| normalize availability | Endpoint Operator | advertisement + facts | 1 | source facts | availability/unknown/degraded | non |
| request bounded refresh | authorized operator | advertisement | 2 | source supports | refresh only, no execution | non |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/normalize availability | oui | oui | oui | sourced explanation | deterministic conditions |
| summarize degradation | oui | oui for source facts | oui | oui, attributed | structured reason list |
| authorize/invent capability | non | validation only | non | interdit | explicit Security/Govern/source path |

## 14. États fonctionnels
`advertised-available`, `advertised-degraded`, `advertised-unavailable`, `unsupported`, `unknown`, `not-advertised`.

## 15. États d’interface
Partial lists unresolved conditions; Stale shows observation age; Permission denied hides restricted details without changing availability; Offline may affect availability but is not identical. No Endpoint screen.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Capability Advertisement | Endpoint technical fact | Settings Fleet, Investigate, Govern, Command | semantic ref/state/source explicit |
| Availability Projection | derived Endpoint fact | authorized consumers | availability != authorization/support |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| capability metadata + local facts | assessment | advertisement/availability | refs/conditions/reasons | source context retained |
| health/platform/version changes | re-evaluation | changed availability | previous/new conditions | history retained |
| advertisement | Settings aggregation | Fleet capability inventory | state/source/freshness | no ownership transfer |

## 18. Dépendances
CAP-EPT-004/005/008/009/010, Settings capability inventory, Security permission model, Studio Tool boundary, OPEN-008.

## 19. Source de vérité
Endpoint owns local advertisement/availability facts. Settings inventory is an aggregation; Security/Govern own permission/authority; Studio Tool is a separate concept.

## 20. Provenance et audit
Record capability ref/version, Agent/scope, conditions, dependency/health refs, state/reason, timestamp, prior state and correlation id.

## 21. Permissions fonctionnelles
Endpoint read for ordinary capability state; permission requirements are projections only. Advertised/available never grants invoke/execute permission.

## 22. Limites et erreurs
Advertised != authorized; advertised != available; available != globally supported; technical capability != Studio Tool; unknown conditions remain unknown.

## 23. Métriques
Availability-state distribution; degraded/unavailable reasons; unknown/unsupported count; freshness; condition-conflict count.

## 24. Classification de livraison
`draft / defined / planned`; no detailed EPT-2+ capability implementation, authorization engine, protocol, code or global support claim.

## 25. Critères d’acceptation
**Given** a capability is advertised but a required dependency is unavailable, **When** assessed, **Then** it is advertised-unavailable/degraded rather than executable.

**Given** a capability is locally available but the user lacks authority, **When** viewed, **Then** availability remains visible only as permitted and no authorization is granted.

**Given** a technical capability exists, **When** Studio consumes its reference, **Then** it does not become a Studio Tool automatically.

## 26. Questions ouvertes
OPEN-008 remains open; exact platform/version capability support is not selected.

## 27. Consommateurs documentaires
Endpoint foundations/future EPT lots, Settings capability inventory/Fleet, Studio integration boundary, Govern/Investigate/Command, registers, Quality and Roadmap.