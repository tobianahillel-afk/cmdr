---
id: CAP-EPT-028
title: Technical Capability Availability, Degradation and Dependency Projection
product: endpoint-agent
module: telemetry
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-018, REQ-PROD-019]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-028 — Technical Capability Availability, Degradation and Dependency Projection

## 1. Définition
Define dynamic capability availability/degradation/dependency states and transitions using declared conditions, source/sensor health, platform/version/context and freshness facts.

## 2. Problème utilisateur
A declared capability can become unavailable or degraded for many reasons that must be explicit without being conflated with authorization or global support.

## 3. Objectifs
- Represent available/degraded/unavailable/unsupported/unknown, dependency unavailable, source disabled, version incompatible, platform unknown, permission/context restriction and freshness.
- Project reasoned transitions to CAP-EPT-011 summary and authorized consumers.

## 4. Non-objectifs
- No authorization grant; no automatic remediation; no support commitment; no Settings administrative mutation.

## 5. Propriétaire
Endpoint Agent owns only endpoint-local technical semantics. Shared, Settings, Investigate, Studio, Govern and Command retain their canonical ownership.

## 6. Utilisateurs
Endpoint Operator; SOC/Investigate analyst; Platform Administrator as consumer; Security/Privacy reviewer; Auditor; authorized Command/Govern/Studio consumers where relevant.

## 7. Conditions d’entrée
Known Agent/tenant scope, source-backed facts and explicit unknown/unsupported states. `OPEN-008` remains open.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Agent/platform/version context | CAP-EPT-001..011 | technical context | oui | source freshness | unknown |
| capability/source facts | CAP-EPT-008/010/011/022/027 | source facts | oui | source-owned | unknown/degraded |
| policy/tenant projection | Platform Settings | restriction context | non | source version | no authorization inference |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Capability Declaration | Endpoint Agent | read/reference only | source permission |
| Sensor Health Observation | Endpoint Agent | read/reference only | source permission |
| Agent Health | Endpoint Agent | read/reference only | source permission |
| Endpoint Policy | Canonical owner | read/reference only | source permission |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Capability Availability | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |
| Capability Dependency State | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |
| Degradation Reason | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |

## 11. Fonctionnalités
Define dynamic capability availability/degradation/dependency states and transitions using declared conditions, source/sensor health, platform/version/context and freshness facts. Missing conditions remain explicit rather than inferred.

## 12. Actions utilisateur
Inspect source-backed facts (class 0), deterministically assess/normalize them (class 1), and request only bounded status refresh where supported (class 2). No response action or Settings-owned administration is introduced.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| inspect | oui | oui | oui | explain only | structured facts |
| assess/derive | oui | oui | oui | explain only | deterministic rules |
| invent source/support/authority | non | validation only | non | interdit | unknown/denied state |

## 14. États fonctionnels
`available`, `degraded`, `unavailable`, `unsupported`, `unknown`, `stale`, `dependency-unavailable`, `context-restricted`.

## 15. États d’interface
No Endpoint Screen ID. Consumers must preserve Partial/Stale/Permission denied/Unavailable and source limitations.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| availability/degradation projection | Endpoint projection | authorized consumers | source/reason/provenance explicit |
| dependency-state projection | Endpoint projection | authorized consumers | source/reason/provenance explicit |
| CAP-EPT-011 summary input | Endpoint projection | authorized consumers | source/reason/provenance explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| source facts | reassessment | Endpoint projection | source/state/reason/time | provenance retained |
| Endpoint projection | consumer handoff | consumer context | stable refs/limitations | no ownership transfer |
| dependency/context change | re-evaluation | changed state | previous/new reason | history retained |

## 18. Dépendances
- CAP-EPT-008/010/011/022/027
- Settings capability inventory
- OPEN-008

## 19. Source de vérité
Endpoint owns local technical facts. Shared owns generic telemetry-event/normalization mechanisms; Settings owns administration; Investigate owns Evidence/Finding/Case; Studio owns Tool/Tool Call/Automation Run; Govern owns Decision/Response Run/Result.

## 20. Provenance et audit
Preserve Agent/version/platform, exact source/sensor refs, time semantics, transformations, limitations, state reasons and consumer handoff. Local audit remains distinct from Shared Trace.

## 21. Permissions fonctionnelles
Identify read needs for capability/source/restriction/provenance facts and deny cross-tenant access. No final RBAC/ABAC or arbitrary atomic permission is selected.

## 22. Limites et erreurs
Declared/observed technical facts are not authorization, Detection, Finding, Evidence, Tool or Govern Result. No API, protocol, port, physical schema, storage/event-bus/SIEM implementation, product code or EPT-3+ scope.

## 23. Métriques
State distribution, reason/dependency coverage, freshness, limitation/provenance completeness and consumer-handoff completeness; no operational SLO commitment.

## 24. Classification de livraison
`draft / defined / planned`; documentary only, not implemented/native/integrated/deployed.

## 25. Critères d’acceptation
**Given** a declared capability loses a required sensor, **When** availability is recalculated, **Then** it becomes degraded/unavailable with dependency reason.

**Given** a capability is available locally, **When** the platform is not globally approved, **Then** local availability remains distinct from globally supported product delivery.

**Given** a capability is available, **When** the user lacks permission, **Then** availability is not converted into authorization and restricted details remain protected.

## 26. Questions ouvertes
- OPEN-008
`OPEN-008` remains open; no platform/source family is declared globally delivered or supported.

## 27. Consommateurs documentaires
Endpoint EPT-2; later Endpoint lots as consumers only; Shared; Investigate; Settings; Studio; Govern; Command; Security/Trust; Quality; Roadmap.
