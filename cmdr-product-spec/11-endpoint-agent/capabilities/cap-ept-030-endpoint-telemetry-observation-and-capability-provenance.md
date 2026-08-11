---
id: CAP-EPT-030
title: Endpoint Telemetry, Observation and Capability Provenance
product: endpoint-agent
module: telemetry
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-OBJ-008, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-030 — Endpoint Telemetry, Observation and Capability Provenance

## 1. Définition
Define end-to-end provenance from Endpoint Agent through source/sensor, local observation, normalization, Shared telemetry-event projection, quality state, capability availability and consumer handoff.

## 2. Problème utilisateur
Consumers and auditors need to reconstruct exactly how Endpoint technical facts were produced and transformed without inventing missing source data.

## 3. Objectifs
- Preserve exact Agent/version/platform facts, source/sensor refs, timestamp semantics, normalization/masking transformations, ordering/freshness/loss state, capability conditions and consumer handoff.
- Link Endpoint local audit where appropriate while keeping it distinct from Shared Trace.

## 4. Non-objectifs
- No immutable-ledger implementation; no cryptographic protocol selection; no Evidence custody claim; no causal inference; no EPT-3+ outcome provenance.

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
| capability/source facts | CAP-EPT-014/015/016/023/024/026/027/028/029 | source facts | oui | source-owned | unknown/degraded |
| policy/tenant projection | Platform Settings | restriction context | non | source version | no authorization inference |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | read/reference only | source permission |
| Local Observation | Endpoint Agent | read/reference only | source permission |
| telemetry-event | Shared | read/reference only | source permission |
| local-audit-event | Endpoint Agent | read/reference only | source permission |
| Trace | Shared | read/reference only | source permission |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Endpoint Telemetry Provenance | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |
| Transformation Record | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |
| Consumer Handoff Provenance | derive/refresh | Endpoint Agent | source-backed; no foreign ownership |

## 11. Fonctionnalités
Define end-to-end provenance from Endpoint Agent through source/sensor, local observation, normalization, Shared telemetry-event projection, quality state, capability availability and consumer handoff. Missing conditions remain explicit rather than inferred.

## 12. Actions utilisateur
Inspect source-backed facts (class 0), deterministically assess/normalize them (class 1), and request only bounded status refresh where supported (class 2). No response action or Settings-owned administration is introduced.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| inspect | oui | oui | oui | explain only | structured facts |
| assess/derive | oui | oui | oui | explain only | deterministic rules |
| invent source/support/authority | non | validation only | non | interdit | unknown/denied state |

## 14. États fonctionnels
`complete`, `partial`, `source-missing`, `transformation-limited`, `restricted`, `unknown`.

## 15. États d’interface
No Endpoint Screen ID. Consumers must preserve Partial/Stale/Permission denied/Unavailable and source limitations.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| provenance chain | Endpoint projection | authorized consumers | source/reason/provenance explicit |
| transformation/masking lineage | Endpoint projection | authorized consumers | source/reason/provenance explicit |
| consumer handoff lineage | Endpoint projection | authorized consumers | source/reason/provenance explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| source facts | reassessment | Endpoint projection | source/state/reason/time | provenance retained |
| Endpoint projection | consumer handoff | consumer context | stable refs/limitations | no ownership transfer |
| dependency/context change | re-evaluation | changed state | previous/new reason | history retained |

## 18. Dépendances
- CAP-EPT-014/015/016/023/024/026/027/028/029
- local-audit.md
- Shared Trace
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
**Given** an observation is normalized and masked, **When** provenance is inspected, **Then** source, mapping version/state and masking transformation remain attributable.

**Given** source metadata is missing, **When** provenance is assembled, **Then** the chain is marked partial/source-missing rather than invented.

**Given** a local audit event is linked, **When** a consumer follows provenance, **Then** the local audit reference remains distinct from Shared Trace and from Evidence custody.

## 26. Questions ouvertes
- OPEN-008
- OPEN-015
`OPEN-008` remains open; no platform/source family is declared globally delivered or supported.

## 27. Consommateurs documentaires
Endpoint EPT-2; later Endpoint lots as consumers only; Shared; Investigate; Settings; Studio; Govern; Command; Security/Trust; Quality; Roadmap.
