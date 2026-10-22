---
id: CAP-EPT-020
title: User, Session and Authentication Activity Observation
product: endpoint-agent
module: telemetry
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-012, REQ-PROD-018, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-020 — User, Session and Authentication Activity Observation

## 1. Définition
Define user, session and authentication activity observations where source/platform support exists, including identity/device context, privilege-context projection, time and privacy.

## 2. Problème utilisateur
Consumers need logon/session facts while ensuring authentication activity is not automatically interpreted as compromise.

## 3. Objectifs
- Represent logon/session activity, user/session/device context, privilege context when available, timestamps, source, platform limitations and provenance.
- Never capture raw secrets and preserve restricted identity metadata.

## 4. Non-objectifs
- No compromise verdict; no credential capture; no identity-provider administration; no session termination.

## 5. Propriétaire
Endpoint Agent owns the endpoint-local semantics defined here. Shared, Platform Settings, Investigate, Studio, Govern and Command retain their canonical objects and generic/admin/interpretive mechanisms. This capability transfers no ownership.

## 6. Utilisateurs
Endpoint Operator; SOC/Investigate analyst; Platform Administrator as an authorized consumer; Security/Privacy reviewer; Govern/Command consumer where relevant; Auditor. Studio may consume technical capability references but receives no Tool ownership from this capability.

## 7. Conditions d’entrée
A known Endpoint Agent/tenant scope, source-backed facts appropriate to this capability, applicable policy/permission projections, and explicit unknown/unsupported states when source/platform/version facts are absent. `OPEN-008` remains open.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Endpoint Agent identity/context | CAP-EPT-001..010 | local technical context | oui | source freshness | state/context unknown |
| source/capability facts | authentication-telemetry | functional facts | oui | source-owned | unknown/degraded rather than inferred |
| tenant/environment/policy projection | Platform Settings | scope/restriction context | non | source version | no authorization/support inference |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| endpoint-agent | Endpoint Agent | source-owned projection only | read/reference within permission |
| telemetry-event | Shared Capabilities | source-owned projection only | read/reference within permission |
| Endpoint Policy | Platform Settings | source-owned projection only | read/reference within permission |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| User Session Observation | derive/create/update documentary functional state | Endpoint Agent | no foreign ownership transfer; source/provenance explicit |
| Authentication Observation | derive/create/update documentary functional state | Endpoint Agent | no foreign ownership transfer; source/provenance explicit |
| Privilege Context Projection | derive/create/update documentary functional state | Endpoint Agent | no foreign ownership transfer; source/provenance explicit |

## 11. Fonctionnalités
Define user, session and authentication activity observations where source/platform support exists, including identity/device context, privilege-context projection, time and privacy.
The contract uses explicit source/platform/version/permission conditions, preserves time and provenance, and exposes partial/degraded/unsupported/unknown states instead of silently filling missing facts.

## 12. Actions utilisateur
| Action | Rôle | Classe | Précondition | Résultat |
|---|---|---:|---|---|
| inspect | authorized Endpoint/consumer role | 0 | read permission and tenant scope | source-backed state only |
| assess/normalize | Endpoint Operator or deterministic service | 1 | required source facts | derived state with provenance |
| request bounded refresh/status | authorized operator | 2 | source supports bounded refresh | refresh/status only; no response action |

No class-3/4 response action is introduced. Any administrative source/policy mutation remains with its canonical owner, especially Platform Settings.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| inspect source-backed state | oui | oui | oui | summary/explanation only | structured fields and reasons |
| derive normalized/quality state | oui | oui | oui | optional explanation | deterministic rules |
| invent event/support/authority | non | validation only | non | interdit | explicit unknown/unsupported/denied state |

Essential functions must work without AI. AI cannot invent events, sources, platform support, maliciousness, Evidence, Finding, Result, authorization or missing provenance.

## 14. États fonctionnels
`observed`, `partial`, `restricted`, `unavailable`, `unsupported`, `unknown`. State names are functional/documentary semantics; they do not define a physical state machine or runtime implementation.

## 15. États d’interface
No Endpoint Screen ID is created. Any future consumer surface must represent Loading/Empty/Partial/Stale/Permission denied/Unavailable without hiding source limitations, and must preserve source ownership and tenant scope.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| authentication/session observation | Endpoint technical fact/projection | authorized Endpoint/Shared/Investigate/Settings/Studio/Govern/Command consumers as applicable | source, scope, freshness, limitations and provenance explicit |
| identity context projection | Endpoint technical fact/projection | authorized Endpoint/Shared/Investigate/Settings/Studio/Govern/Command consumers as applicable | source, scope, freshness, limitations and provenance explicit |
| privacy/limitation state | Endpoint technical fact/projection | authorized Endpoint/Shared/Investigate/Settings/Studio/Govern/Command consumers as applicable | source, scope, freshness, limitations and provenance explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| source/local facts | new observation or reassessment | Endpoint functional state/projection | Agent/source/scope/time/reason/provenance | source context retained |
| Endpoint projection | consumer handoff | authorized consumer context | stable refs/freshness/limitations | no ownership or permission transfer |
| source condition changes | re-evaluation | degraded/unavailable/unknown or refreshed state | previous/new state and reason | history/provenance retained |

## 18. Dépendances
- authentication-telemetry
- CAP-EPT-015/016/026
- OPEN-008
- `CAP-EPT-001..014` remain intact and are consumed, not rewritten.
- Shared generic event/Trace/Activity/Jobs/Search/Reporting/Export mechanisms remain Shared-owned.
- `OPEN-008` remains unresolved for actual platform/source delivery.

## 19. Source de vérité
Endpoint is source of truth only for endpoint-local technical facts and the semantics defined here. Shared remains source of truth for `telemetry-event` generic envelope/normalization mechanisms where applicable; Settings remains source of truth for administrative configuration; Investigate for Evidence/Finding/Case; Studio for Tool/Tool Call/Automation Run; Govern for Decision/Response Run/Result.

## 20. Provenance et audit
Record or preserve Agent identity/version, tenant/environment scope, exact source/sensor reference where available, source and observation timestamps, transformation/normalization/masking state, limitations, previous/new state when derived, correlation reference and consumer handoff. `local-audit-event` remains Endpoint-owned and is not Shared Trace.

## 21. Permissions fonctionnelles
Functional needs include source/observation read, restricted telemetry read where applicable, normalization/quality metadata read, capability declaration/availability read, provenance read and explicit cross-tenant denial. EPT-2 does not finalize RBAC/ABAC and does not create arbitrary atomic permission namespaces.

## 22. Limites et erreurs
Telemetry/Observation is not Detection, Finding, Evidence, Alert, Incident or Govern Result. Missing/late/duplicate/gap/degraded states remain explicit. Platform/source/version support is never inferred. No API, protocol, port, physical event schema, storage engine, event bus, SIEM, final query language, product code or EPT-3+ behavior is defined.

## 23. Métriques
Conceptual metrics: state distribution; source/freshness coverage; unknown/unsupported/degraded reasons; late/duplicate/gap/loss candidates where applicable; normalization/masking limitations; consumer handoff completeness. Metrics describe documentary semantics, not operational SLO commitments.

## 24. Classification de livraison
`draft / defined / planned`. This capability is documentary functional specification only. It is not implemented/native/integrated/deployed and closes no platform support decision.

## 25. Critères d’acceptation
**Given** a successful logon is observed, **When** it is projected, **Then** the event remains authentication activity and not a compromise conclusion.

**Given** raw secret material could be present in a source, **When** projection is prepared, **Then** raw secret capture is prohibited and sensitive fields are excluded/masked.

**Given** authentication telemetry is unsupported on a platform, **When** capability is assessed, **Then** state is unsupported/platform-dependent without a global support claim.

## 26. Questions ouvertes
- OPEN-008
`OPEN-008` remains open. No Windows/Linux/macOS, cloud/container/mobile, sensor family or observation class is declared globally delivered/supported by this capability.

## 27. Consommateurs documentaires
Endpoint EPT-2 maps/registers/quality evidence; Endpoint later lots as consumers only; Shared telemetry/event mechanisms; Investigate; Platform Settings; Studio; Govern; Command; Security/Trust; Requirements/Open-decision traceability; Delivery Roadmap Phase 5.
