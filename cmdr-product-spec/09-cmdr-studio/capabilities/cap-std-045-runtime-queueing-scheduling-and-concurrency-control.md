---
id: CAP-STD-045
title: Runtime Queueing, Scheduling and Concurrency Control
product: cmdr-studio
module: control-room
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-045 — Runtime Queueing, Scheduling and Concurrency Control

## 1. Définition
Définit la sémantique Automation Run de queueing/scheduling/concurrency : immediate/delayed/scheduled start candidate, priority projection, resource availability, tenant/runtime concurrency bounds, conflicting/duplicate execution, stale schedule, expiry and blocked scheduling. Shared may own generic queue/scheduling mechanisms; Studio owns Run meaning.

## 2. Problème utilisateur
Une Run en queue ou planifiée peut être prise pour une exécution démarrée, et un scheduler générique peut être confondu avec la sémantique métier de l’Automation Run.

## 3. Objectifs
Représenter queue/schedule intent; vérifier availability/concurrency/expiry/duplicates; séparer scheduling from start; preserve tenant/runtime bounds; integrate Shared mechanisms without ownership transfer.

## 4. Non-objectifs
Ne sélectionne/implémente aucun scheduler, queue backend, algorithm, SLO, provider/runtime, deployment schedule or Govern Response Run scheduling semantics.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Studio Operator principal; human supervisor, source caller, platform/runtime operator and Auditor secondary.

## 7. Conditions d’entrée
Automation Run ready, current execution context/bounds, requested timing, runtime availability projection and scheduling permission available.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Automation Run | CAP-STD-042/043 | Run/timing subject | oui | current | no schedule |
| requested start timing | caller/operator | immediate/delayed/time window | oui | current | unscheduled |
| runtime/resource availability | Settings/runtime owner | availability/capacity projection | oui for scheduling | fresh enough | queued/blocked |
| concurrency/tenant limits | Studio policy context + Settings/Security | functional bounds | oui | current | blocked |
| duplicate/idempotency context | Run/STD-2 | duplicate-risk references | oui | current | duplicate-candidate |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Run | CMDR Studio | state/context/schedule intent | read/manage |
| Job/queue mechanism | Shared Capabilities | generic queue/progress ref | read/link |
| Runtime/Environment | Platform Settings/runtime owner | availability/maintenance | read |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Run Queue/Schedule projection | create/update/cancel | CMDR Studio | semantic intent; generic mechanism remains Shared |
| Shared Job reference | link only | Shared Capabilities | Job ≠ Automation Run |

## 11. Fonctionnalités
Immediate/delayed/scheduled candidate; queue status; priority projection; capacity/tenant/runtime concurrency; conflict/duplicate detection; schedule staleness/expiry; blocked scheduling; no-start guarantee.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect queue/schedule | Studio Operator | Automation Run | 0 | read | timing/position/blocks visible | non |
| Validate concurrency/availability | Studio Operator | Run schedule | 1 | current bounds | eligible/blocked | non |
| Queue/schedule/cancel-before-start | authorized operator | Automation Run | 2 | schedule permission | versioned schedule intent | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
ready, queue-requested, queued, schedule-requested, scheduled, delayed, capacity-blocked, concurrency-blocked, duplicate-candidate, runtime-unavailable, stale-schedule, expired-before-start, cancelled-before-start.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Run scheduling state | Automation Run projection | CAP-STD-043/046/049 | queue/schedule != start |
| generic Job/queue ref | Shared reference | Control Room/Audit | identity remains Shared |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-STD-043 | Run ready | CAP-STD-045 | Run/bounds/requested timing | Run |
| CAP-STD-045 | eligible start time | CAP-STD-046 | Run + fresh scheduling/availability context | schedule |
| Shared queue/runtime | status/capacity update | CAP-STD-045 | source ref/status/time | same Run |

## 18. Dépendances
CAP-STD-027/042/043/046/049; Shared Background Jobs/generic queue/scheduling primitives; Settings health/environment; Security; Govern boundary; OPEN-013/015.

## 19. Source de vérité
Studio owns Automation Run queue/schedule/concurrency semantics. Shared owns generic Job/queue mechanisms; runtime/Settings own actual capacity/availability. A schedule does not assert start.

## 20. Provenance et audit
Record Run id, requested timing/window, queue/job refs, priority projection, concurrency/tenant/runtime bounds, duplicate checks, availability sources/timestamps, cancellations/expiry and actor/correlation.

## 21. Permissions fonctionnelles
Run read/schedule; queue request; cancel-before-effect; runtime availability read; cross-tenant scheduling requires explicit permission. No scheduler-admin permission implied.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Runtime unavailable, capacity exceeded, conflicting Run, duplicate candidate, stale schedule, expired window, missing queue acknowledgement or permission denial keeps Run queued/blocked/unknown.

## 23. Métriques
Queue wait; schedules expiring; concurrency blocks; duplicate candidates; runtime-unavailable queues; queued→running without start confirmation—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** a Run is queued while the runtime is unavailable, **When** queue state is viewed, **Then** it remains queued/blocked and is not marked started.

**Given** a scheduled Run reaches an expired start window, **When** scheduling evaluates it, **Then** it becomes expired-before-start and no silent re-schedule or start occurs.

**Given** two Runs conflict with a per-runtime concurrency bound, **When** eligibility is checked, **Then** at least one remains blocked/queued according to explicit rules without inventing capacity.

## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Control Room, Automation Runs, Shared Background Jobs, Settings Health, CAP-STD-046/049/051, Security, Quality.
