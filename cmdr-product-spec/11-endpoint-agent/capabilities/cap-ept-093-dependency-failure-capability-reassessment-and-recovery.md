---
id: CAP-EPT-093
title: Dependency Failure, Capability Reassessment and Recovery
product: endpoint-agent
module: resilience
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-005, REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-093 — Dependency Failure, Capability Reassessment and Recovery

## 1. Définition
Définir la relation entre panne/restauration de dépendance locale, capability availability et recovery/revalidation sans assimiler dépendance restaurée à capability saine.

## 2. Problème utilisateur
Un composant ou sensor peut revenir disponible alors que ses capacités restent dégradées, stale ou nécessitent une nouvelle validation.

## 3. Objectifs
Represent dependency unavailable/restored, sensor/component/runtime failure, affected capability, availability reassessment, degraded/restored/revalidation and provenance.

## 4. Non-objectifs
Aucun dependency manager, auto-healing engine, watchdog, support claim, full endpoint recovery guarantee or implementation.

## 5. Propriétaire
Endpoint owns local dependency/capability-state facts; Shared retains generic Recovery and Settings retains administrative configuration/health aggregation.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, SOC consumer, Reliability/Security Reviewer, Auditor.

## 7. Conditions d’entrée
Known dependency/component/sensor and capability refs, failure/restoration observations, Agent identity/version, freshness and applicable policy context.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| dependency/component state | Endpoint local source | technical fact | oui | current | unknown |
| capability declaration/availability | CAP-EPT-011/027/028 | capability context | oui | current | impact unknown |
| health/resource state | CAP-EPT-008/092 | diagnostic context | non | fresh | partial reassessment |
| policy/config projection | Settings | dependency constraints | non | current | no admin inference |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Endpoint Agent | Endpoint | component/version | read |
| Technical Capability Availability | Endpoint | affected capability | read |
| Agent Health Assessment | Endpoint | health context | read |
| Endpoint Policy | Settings | config constraints | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Dependency State | create/refresh | Endpoint | local technical fact |
| Capability Reassessment | derive/version | Endpoint | restored dependency != healthy capability |
| Recovery/Revalidation State | derive | Endpoint | explicit checks/limits |

## 11. Fonctionnalités
Record unavailable/failure; link affected capabilities; reassess availability; track dependency restoration; require capability revalidation; expose degraded/unknown/healthy-for-checks state and recovery provenance.

## 12. Actions utilisateur
Inspect Class 0; reassess/revalidate Class 1; request bounded diagnostics/recovery Class 2/3 depending effect/policy. No unsupported self-healing inference.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| map dependency to capabilities | oui | oui | oui | explain | dependency map |
| reassess availability | oui | oui | oui | summarize | explicit health/capability checks |
| invent healthy recovery | non | interdit | non | interdit | degraded/unknown/manual review |

## 14. États fonctionnels
`available`, `unavailable`, `failed`, `degraded`, `restoring`, `restored-unverified`, `revalidation-pending`, `revalidated`, `capability-degraded`, `capability-unavailable`, `recovery-failed`, `unknown`, `stale`.

## 15. États d’interface
No Screen ID. Dependency restored and capability revalidated remain separate; unknown/stale are visible.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Dependency State | Endpoint fact | health/Settings/Quality | source/freshness explicit |
| Capability Reassessment | Endpoint derived fact | CAP-EPT-087/098/consumers | revalidation required |
| Recovery provenance | technical context | CAP-EPT-097/099 | no full-service claim |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| component/sensor | failure | CAP-EPT-093 | dependency/capability refs | local context |
| dependency restored | restoration | revalidation | restored state + affected capabilities | reassessment |
| CAP-EPT-093 | capability state changed | consumers/Settings | availability/reason/limits | no owner transfer |

## 18. Dépendances
CAP-EPT-008/011/027/028/087/091/092/097/098; Shared Recovery; Settings; OPEN-008.

## 19. Source de vérité
Endpoint SOT for local dependency/capability reassessment facts; Shared and Settings retain generic/admin ownership.

## 20. Provenance et audit
Dependency/component identity, source/version, failure/restoration time, affected capabilities, health/resource context, reassessment/revalidation results and limitations.

## 21. Permissions fonctionnelles
Dependency/capability state read, diagnostic/revalidation request, sensitive provenance read, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Dependency restored != capability healthy automatically; recovery succeeded locally != business service restored; platform support never inferred.

## 23. Métriques
Dependency failure/restoration, capability degraded/unavailable, revalidation latency/state, recovery failure and unknown/stale counts.

## 24. Classification de livraison
`draft / defined / planned`; no self-healing/dependency manager implementation.

## 25. Critères d’acceptation
**Given** a sensor dependency returns, **When** its capability has not been revalidated, **Then** state is restored-unverified rather than healthy.

**Given** a dependency failure affects two capabilities, **When** reassessed, **Then** each affected capability and limitation is explicit.

**Given** AI is unavailable, **When** dependency mapping is evaluated, **Then** deterministic/manual reassessment remains available.

## 26. Questions ouvertes
OPEN-008 remains open; platform-specific dependencies and recovery mechanisms are not selected.

## 27. Consommateurs documentaires
EPT-6 Resilience/Update/Security, Settings, Shared Recovery, Command/Investigate consumers, Quality, registers and Roadmap.