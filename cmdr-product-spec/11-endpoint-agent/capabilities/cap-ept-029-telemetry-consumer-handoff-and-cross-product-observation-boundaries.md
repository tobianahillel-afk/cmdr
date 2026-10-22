---
id: CAP-EPT-029
title: Telemetry Consumer Handoff and Cross-Product Observation Boundaries
product: endpoint-agent
module: telemetry
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-OBJ-008]
open_decisions: [OPEN-008, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-029 — Telemetry Consumer Handoff and Cross-Product Observation Boundaries

## 1. Définition
Define how Endpoint telemetry/observation/capability projections are consumed by Investigate, Command, Studio, Govern, Settings and Shared without transferring canonical ownership.

## 2. Problème utilisateur
Cross-product consumers need stable source-backed handoffs without automatically promoting technical observations into Evidence, Finding, Result, Tool or administrative source objects.

## 3. Objectifs
- Define consumer-specific handoff context, source references, freshness/limitations/provenance and return links.
- Preserve Investigate Evidence/Finding/Case, Studio Tool, Govern Result/Response Run, Settings administration and Shared generic mechanisms.

## 4. Non-objectifs
- No Evidence/Finding/Result/Tool creation; no source administration; no product-screen design; no execution handoff.

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
| capability/source facts | Shared telemetry-event | source facts | oui | source-owned | unknown/degraded |
| policy/tenant projection | Platform Settings | restriction context | non | source version | no authorization inference |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| telemetry-event | Shared | read/reference only | source permission |
| Evidence | Investigate | read/reference only | source permission |
| Finding | Investigate | read/reference only | source permission |
| Tool | Studio | read/reference only | source permission |
| Automation Run | Studio | read/reference only | source permission |
| Decision | Govern | read/reference only | source permission |
| Response Run | Govern | read/reference only | source permission |
| Result | Govern | read/reference only | source permission |
| Endpoint Agent Fleet | Platform Settings | read/reference only | source permission |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Telemetry Consumer Handoff | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |
| Observation Context Package | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |
| Capability Projection Handoff | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |

## 11. Fonctionnalités
Define how Endpoint telemetry/observation/capability projections are consumed by Investigate, Command, Studio, Govern, Settings and Shared without transferring canonical ownership. Missing conditions remain explicit rather than inferred.

## 12. Actions utilisateur
Inspect source-backed facts (class 0), deterministically assess/normalize them (class 1), and request only bounded status refresh where supported (class 2). No response action or Settings-owned administration is introduced.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| inspect | oui | oui | oui | explain only | structured facts |
| assess/derive | oui | oui | oui | explain only | deterministic rules |
| invent source/support/authority | non | validation only | non | interdit | unknown/denied state |

## 14. États fonctionnels
`prepared`, `consumed`, `restricted`, `stale`, `partial`, `unavailable`.

## 15. États d’interface
No Endpoint Screen ID. Consumers must preserve Partial/Stale/Permission denied/Unavailable and source limitations.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Investigate observation context | Endpoint projection | authorized consumers | source/reason/provenance explicit |
| Settings capability/source projection | Endpoint projection | authorized consumers | source/reason/provenance explicit |
| Studio/Govern/Command technical projection | Endpoint projection | authorized consumers | source/reason/provenance explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| source facts | reassessment | Endpoint projection | source/state/reason/time | provenance retained |
| Endpoint projection | consumer handoff | consumer context | stable refs/limitations | no ownership transfer |
| dependency/context change | re-evaluation | changed state | previous/new reason | history retained |

## 18. Dépendances
- Shared telemetry-event
- Investigate
- Settings
- Studio
- Govern
- Command
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
**Given** Investigate consumes a telemetry event, **When** an analyst opens it, **Then** the event remains telemetry and does not automatically become Evidence or a Finding.

**Given** Studio consumes an Endpoint capability declaration, **When** a workflow references it, **Then** the capability remains distinct from a Studio Tool.

**Given** Govern consumes a technical observation, **When** a Response Run is reviewed, **Then** the observation remains technical context and does not automatically become a Govern Result.

## 26. Questions ouvertes
- OPEN-008
- OPEN-015
`OPEN-008` remains open; no platform/source family is declared globally delivered or supported.

## 27. Consommateurs documentaires
Endpoint EPT-2; later Endpoint lots as consumers only; Shared; Investigate; Settings; Studio; Govern; Command; Security/Trust; Quality; Roadmap.
