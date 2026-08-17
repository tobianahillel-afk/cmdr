---
id: CAP-EPT-086
title: Update Progress, Deferral, Failure and Retry Semantics
product: endpoint-agent
module: updates
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-12
requirement_ids: [REQ-PROD-006, REQ-PROD-012, REQ-PROD-017, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013]
source-of-truth: canonical
---
# CAP-EPT-086 — Update Progress, Deferral, Failure and Retry Semantics

## 1. Définition
Définir progress, deferral, failure, interruption et retry d’une opération d’update locale sans transformer retry en recovery ni garantir absence de duplication.

## 2. Problème utilisateur
Un update différé, interrompu ou retryable peut être interprété comme échec définitif ou rétablissement réussi.

## 3. Objectifs
Track progress, deferred/retryable/blocked/failed/interrupted/stale-target, retry request/eligibility, attempt history and provenance.

## 4. Non-objectifs
Aucun retry engine, backoff physique, scheduler, exactly-once guarantee, recovery/reversion, package transport ou implementation.

## 5. Propriétaire
Endpoint owns local update attempt/progress facts; Shared retains generic retry mechanisms and Settings retains admin scheduling/assignment.

## 6. Utilisateurs
Endpoint Operator, Platform Administrator, Release Reviewer, Auditor, Support Operator.

## 7. Conditions d’entrée
Identifiable Local Update Operation/attempt, target/package refs, runtime events, current assignment/policy and explicit error/deferral source.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| update operation | CAP-EPT-085 | attempt context | oui | exact | orphaned/unknown |
| progress/error events | Endpoint runtime | local facts | conditionnel | event time | progress unknown |
| current target | CAP-EPT-082 | admin projection | oui for stale check | current | target unknown |
| retry constraints | Policy/Shared contract | eligibility context | non | current | retry unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Local Update Operation | Endpoint | attempt/state | read |
| Local Update Assignment Projection | Endpoint/Settings source | current target | read |
| Endpoint Policy | Settings | retry/defer constraints | read projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Update Attempt State | create/transition | Endpoint | one attempt, explicit history |
| Update Progress Projection | append/refresh | Endpoint | no success invention |
| Retry Eligibility | derive | Endpoint | retry != recovery |

## 11. Fonctionnalités
Record progress; deferred/blocked/interrupted/failed; detect stale target; assess retryability; record retry request/new attempt relation; preserve previous attempts and late status.

## 12. Actions utilisateur
Inspect Class 0; retry/defer eligibility Class 1; retry/defer request Class 2/3 according effect/policy; no autonomous indefinite retry.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| calculate progress/state | oui | oui | oui | summarize | source events |
| assess retryability | oui | oui if rules exist | oui | explain | explicit rules |
| invent success/retry forever | non | interdit | non | interdit | stop/unknown/manual review |

## 14. États fonctionnels
`pending`, `in-progress`, `deferred`, `blocked`, `retryable`, `retry-requested`, `retrying`, `interrupted`, `failed`, `stale-target`, `completed-technical`, `status-unknown`.

## 15. États d’interface
No Screen ID. Deferred and failed remain distinct; late status cannot rewrite attempt history silently.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| attempt/progress state | Endpoint fact | CAP-EPT-087/088/097 | attempt identity preserved |
| retry eligibility | derived fact | operator/Settings | no recovery guarantee |
| failure/deferral reason | technical fact | Quality/Support | source-backed |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-085 | operation starts | CAP-EPT-086 | attempt/target/package | lifecycle retained |
| CAP-EPT-086 | retry allowed | new attempt/CAP-EPT-085 | prior attempt + reason | history retained |
| CAP-EPT-086 | terminal/partial | CAP-EPT-087/088 | outcome/limits | verification/recovery context |

## 18. Dépendances
CAP-EPT-024/025/062/082/085/087/088; Shared offline/retry contract; Settings; OPEN-008/013.

## 19. Source de vérité
Endpoint SOT for local update attempt/progress facts; Shared generic retry and Settings admin policy remain external.

## 20. Provenance et audit
Attempt ids, parent/retry relation, target/package, progress events, deferral/failure reason, requester, timestamps, stale-target evidence and outcome.

## 21. Permissions fonctionnelles
Update progress read, defer/retry request, sensitive error/provenance read, cross-tenant deny; final RBAC remains Security-owned.

## 22. Limites et erreurs
Deferred != failed; retry != recovery; retry != duplicate-free guarantee; stale target may block further effect; unknown remains unknown.

## 23. Métriques
Deferrals, failures, interruptions, retries/attempts, stale-target blocks, unknown status and retry exhaustion candidates.

## 24. Classification de livraison
`draft / defined / planned`; no retry/update engine implementation.

## 25. Critères d’acceptation
**Given** an update is deferred by policy, **When** state is reviewed, **Then** it is deferred and not failed.

**Given** a failed attempt is retryable, **When** retry starts, **Then** a linked new attempt is recorded rather than rewriting the failed attempt.

**Given** AI is unavailable, **When** retryability rules apply, **Then** deterministic/manual assessment remains possible.

## 26. Questions ouvertes
OPEN-008/013 remain open; no universal retry count/backoff or final action class is selected.

## 27. Consommateurs documentaires
EPT-6 Update/Resilience, Settings, Shared Retry, Security, Quality, registers and Roadmap.