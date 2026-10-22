---
id: CAP-EPT-027
title: Technical Capability Declaration and Support Conditions
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
# CAP-EPT-027 — Technical Capability Declaration and Support Conditions

## 1. Définition
Define the detailed technical capability declaration contract beneath CAP-EPT-011, including semantic family, dependencies, platform/version/context conditions, supported observation classes, limitations, state/reason and provenance.

## 2. Problème utilisateur
Consumers need evidence-backed declaration details without confusing declared capability with authorization, runtime availability or globally supported product delivery.

## 3. Objectifs
- Represent capability semantic identifier/family, source/sensor dependency, platform/version conditions, required privilege/context, observation classes, limitations, declared state/reason and provenance.
- Keep CAP-EPT-011 as the foundation summary rather than duplicating it.

## 4. Non-objectifs
- No capability authorization; no global support matrix closure; no Tool definition; no provider/runtime implementation; no EPT-3+ execution.

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
| capability/source facts | CAP-EPT-004/005/011/015 | source facts | oui | source-owned | unknown/degraded |
| policy/tenant projection | Platform Settings | restriction context | non | source version | no authorization inference |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Capability Advertisement | Endpoint Agent | read/reference only | source permission |
| Platform Fact | Endpoint Agent | read/reference only | source permission |
| Agent Version | Endpoint Agent | read/reference only | source permission |
| Observation Source | Endpoint Agent | read/reference only | source permission |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical Capability Declaration | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |
| Capability Dependency | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |
| Support Condition | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |

## 11. Fonctionnalités
Define the detailed technical capability declaration contract beneath CAP-EPT-011, including semantic family, dependencies, platform/version/context conditions, supported observation classes, limitations, state/reason and provenance. Missing conditions remain explicit rather than inferred.

## 12. Actions utilisateur
Inspect source-backed facts (class 0), deterministically assess/normalize them (class 1), and request only bounded status refresh where supported (class 2). No response action or Settings-owned administration is introduced.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| inspect | oui | oui | oui | explain only | structured facts |
| assess/derive | oui | oui | oui | explain only | deterministic rules |
| invent source/support/authority | non | validation only | non | interdit | unknown/denied state |

## 14. États fonctionnels
`declared`, `not-declared`, `conditional`, `unsupported`, `unknown`.

## 15. États d’interface
No Endpoint Screen ID. Consumers must preserve Partial/Stale/Permission denied/Unavailable and source limitations.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| detailed capability declaration | Endpoint projection | authorized consumers | source/reason/provenance explicit |
| support-condition projection | Endpoint projection | authorized consumers | source/reason/provenance explicit |
| dependency references | Endpoint projection | authorized consumers | source/reason/provenance explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| source facts | reassessment | Endpoint projection | source/state/reason/time | provenance retained |
| Endpoint projection | consumer handoff | consumer context | stable refs/limitations | no ownership transfer |
| dependency/context change | re-evaluation | changed state | previous/new reason | history retained |

## 18. Dépendances
- CAP-EPT-004/005/011/015
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
**Given** a capability is declared for a specific source and version condition, **When** consumers inspect it, **Then** the declaration exposes those conditions without asserting global product support.

**Given** a capability is declared, **When** a user lacks execution authority, **Then** declaration remains distinct from authorization.

**Given** CAP-EPT-011 summarizes availability, **When** declaration details change, **Then** CAP-EPT-027 provides detailed conditions while CAP-EPT-011 remains the foundation summary.

## 26. Questions ouvertes
- OPEN-008
`OPEN-008` remains open; no platform/source family is declared globally delivered or supported.

## 27. Consommateurs documentaires
Endpoint EPT-2; later Endpoint lots as consumers only; Shared; Investigate; Settings; Studio; Govern; Command; Security/Trust; Quality; Roadmap.
