---
id: CAP-EPT-092
title: Resource Pressure, Backpressure and Degraded-Operation Semantics
product: endpoint-agent
module: resilience
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008]
source-of-truth: canonical
---
# CAP-EPT-092 — Resource Pressure, Backpressure and Degraded-Operation Semantics

## 1. Définition
Définir les faits locaux de pression CPU/mémoire/stockage/réseau, backpressure, réduction de service et capability degradation sans définir resource manager ni perte silencieuse.

## 2. Problème utilisateur
Sous pression, l’Agent peut réduire certaines fonctions sans être totalement en panne; cette réduction ne doit pas masquer sampling, drops ou capacités affectées.

## 3. Objectifs
Represent resource pressure/unavailable, queue pressure, source degradation, reduced operation, load-shedding boundary, capability degradation, recovery state and consumer projection.

## 4. Non-objectifs
Aucun scheduler de ressources, seuil universel, storage guarantee, tuning OS, kill implementation, SLO final ou implementation.

## 5. Propriétaire
Endpoint owns local resource/degradation facts; Settings may provide policy limits and Shared retains generic backpressure/job mechanisms.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, SOC consumer, Reliability/Security Reviewer, Auditor.

## 7. Conditions d’entrée
Agent/resource observations, relevant queue/source/capability state, applicable policy constraints and freshness.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| CPU/memory/storage/network facts | Endpoint local monitor | pressure facts | si sourcé | current | pressure unknown |
| queue/backpressure | CAP-EPT-025/089 | delivery pressure | non | current | no queue inference |
| capability availability | CAP-EPT-011 | affected capability context | oui for degradation | current | impact unknown |
| policy limits | Settings | boundary projection | non | current | source defaults/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Agent Health Assessment | Endpoint | health context | read |
| Local Queue State | Endpoint | pressure context | read |
| Technical Capability Availability | Endpoint | affected capabilities | read |
| Endpoint Policy | Settings | resource constraints | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Resource Pressure State | derive/refresh | Endpoint | source-backed |
| Degraded Operation State | derive/transition | Endpoint | complete failure not implied |
| Load-Shedding/Reduction Fact | record | Endpoint | no silent loss |

## 11. Fonctionnalités
Observe pressure; correlate queue/source/capability impact; represent reduced operation/load shedding; expose drop/sampling implications separately; reassess after relief; project limitations to consumers.

## 12. Actions utilisateur
Inspect Class 0; calculate/reassess Class 1; request bounded diagnostics/status Class 2. Effectful resource changes remain policy/platform-specific and out of final classification here.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| assess pressure/degradation | oui | oui if thresholds sourced | oui | explain | explicit rules/facts |
| summarize affected capabilities | oui | oui | oui | oui | capability map |
| hide loss/invent recovery | non | interdit | non | interdit | explicit limitation/state |

## 14. États fonctionnels
`normal`, `pressure-observed`, `resource-unavailable`, `backpressured`, `reduced-operation`, `load-shedding`, `capability-degraded`, `recovering`, `recovered-observed`, `state-unknown`, `stale`.

## 15. États d’interface
No Screen ID. Degraded remains distinct from complete failure; any reduction/drop indication is visible and source-attributed.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Resource Pressure State | Endpoint fact | health/Settings/Quality | source/limits explicit |
| Degraded Operation Projection | Endpoint fact | Command/Investigate consumers | affected capability named |
| reduction/drop context | technical limitation | CAP-EPT-089/093 | no silent-loss inference |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| local monitor | pressure | CAP-EPT-092 | resource/source/queue facts | Agent context |
| CAP-EPT-092 | capability affected | CAP-EPT-093 | capability/dependency/pressure refs | reassessment |
| CAP-EPT-092 | pressure relieved | health/capability reassessment | prior/new state | provenance retained |

## 18. Dépendances
CAP-EPT-008/011/025/089/093/097; resilience resource-guardrails; Settings; Shared; OPEN-008.

## 19. Source de vérité
Endpoint SOT for observed local pressure/degradation; Settings owns admin policy, Shared generic mechanisms.

## 20. Provenance et audit
Agent/resource/source refs, observations/times, threshold/policy refs if any, affected capabilities, reduction/drop indicators, prior/new state and recovery observation.

## 21. Permissions fonctionnelles
Resource/degradation read, diagnostics request, sensitive telemetry read, provenance read, cross-tenant deny; no final RBAC.

## 22. Limites et erreurs
Degraded != complete failure; load shedding != arbitrary/silent data loss; backpressure != durability guarantee; recovered pressure != capability healthy automatically.

## 23. Métriques
Pressure/degradation duration, affected capabilities, load-shedding/drop indications, recovery transitions and unknown/stale states.

## 24. Classification de livraison
`draft / defined / planned`; no resource-control implementation.

## 25. Critères d’acceptation
**Given** memory pressure reduces one sensor, **When** state is assessed, **Then** affected capability is degraded without declaring the whole Agent failed.

**Given** load shedding occurs, **When** output quality is reviewed, **Then** reduction/drop indication is explicit.

**Given** AI is unavailable, **When** sourced thresholds are evaluated, **Then** deterministic/manual assessment remains available.

## 26. Questions ouvertes
OPEN-008 remains open; platform-specific budgets and controls are not selected.

## 27. Consommateurs documentaires
EPT-6 Resilience, EPT-2 telemetry, Settings, Shared, Security, Quality, registers and Roadmap.