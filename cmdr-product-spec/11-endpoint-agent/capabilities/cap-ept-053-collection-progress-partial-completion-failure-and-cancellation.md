---
id: CAP-EPT-053
title: Collection Progress, Partial Completion, Failure and Cancellation
product: endpoint-agent
module: collection
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-014, REQ-PROD-018, REQ-PROD-019, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-053 — Collection Progress, Partial Completion, Failure and Cancellation

## 1. Définition
Définir l’état technique transversal d’une opération de collecte et de ses items : preparation, queue/local wait, running, progress, partial, blocked, failed, timed-out, cancel requested/cancelled, completed, transfer pending, offline/unavailable et resume/retry eligibility conceptuels.

## 2. Problème utilisateur
Queue, partial success, offline et cancel peuvent être interprétés à tort comme démarrage, succès ou rollback.

## 3. Objectifs
Séparer requested/confirmed states ; agréger item progress sans masquer failures ; exposer offline queue/expiry ; définir retry/resume eligibility ; préserver last confirmed target-side state.

## 4. Non-objectifs
Aucune state machine technique finale, scheduler, Shared Background Job ownership, exactly-once guarantee, rollback, transport or automatic retry policy.

## 5. Propriétaire
Endpoint owns local operation progress/state. Shared owns generic Background Jobs; Investigate owns Collection Job business state.

## 6. Utilisateurs
DFIR Analyst, Endpoint Operator, Case Analyst, Response Operator, Auditor.

## 7. Conditions d’entrée
Collection operation/attempt reference, planned items, local status events or explicit lack thereof, time/expiry and cancellation permissions.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| operation/items | CAP-EPT-049..052 | execution refs | oui | current | no status |
| target connectivity | EPT-1 | online/offline | oui | last seen | status unknown |
| local queue/progress events | Endpoint | state/progress | conditionnel | event time | stale/unknown |
| cancel/retry/resume request | caller/authority | control intent | non | request time | no control |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Plan/Attempts | Endpoint | items/status | read |
| Collection Job | Investigate | business correlation | read ref |
| Background Job | Shared | generic queue/progress | read ref only |
| Response Run/Automation Run | Govern/Studio | correlation only | read ref |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Collection Operation State | reconcile/version | Endpoint | request != confirmation |
| Item Progress State | reconcile | Endpoint | per item/category |
| Retry/Resume Eligibility | derive | Endpoint | eligibility only, not auto execution |

## 11. Fonctionnalités
Reconcile queue/offline/running/progress/per-item terminal outcomes, compute partial vs complete, expose expiry/timeout, accept cancel/stop intent without fabricating confirmation, mark retry/resume eligibility and link subsequent attempt if explicitly initiated.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect progress | Analyst | operation state | 0 | read | detailed status | non |
| compute partial/completeness | service | items | 1 | item states | derived aggregate | non |
| request cancel | Operator | operation | 2 | cancellable | cancel-requested | OPEN-013 |
| retry/resume | authorized caller | failed/interrupted items | 2 | eligible + gates | new linked attempt/request | according impact |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate progress | oui | oui | oui | summarize | raw states |
| detect partial/timeout | oui | oui | oui | explain | state rules |
| recommend retry | oui | rules | oui | suggestion | manual selection |
| silently retry/resume | non | interdit | non | interdit | explicit control |

## 14. États fonctionnels
`preparing`, `queued-local`, `waiting-endpoint`, `running`, `partial`, `blocked`, `failed`, `timed-out`, `cancel-requested`, `cancelled`, `completed`, `transfer-pending`, `expired`, `unavailable`, `status-unknown`.

## 15. États d’interface
No Screen ID. Last confirmed state/time shown; `cancel-requested` is never rendered as `cancelled`; offline never implies failed or running.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| operation/item progress | Endpoint projection | Investigate/CAP-EPT-054/055 | requested/confirmed distinction |
| aggregate completion | derived state | Investigate | partial != complete |
| retry/resume eligibility | diagnostic | caller | no automatic retry |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-049..052 | status event | CAP-EPT-053 | attempt/item/state/time | attempt |
| CAP-EPT-053 | item terminal | CAP-EPT-054 | item/completeness/status | operation |
| CAP-EPT-053 | transfer pending | CAP-EPT-055 | output refs/status | operation |

## 18. Dépendances
Endpoint collection-queue/resilience boundary, CAP-EPT-049..055, Investigate CAP-INV-203, Shared Background Jobs, OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT for local confirmed operation state; Shared SOT generic job mechanism; Investigate SOT business Collection Job.

## 20. Provenance et audit
Operation/attempt/item IDs, source events, requested vs confirmed controls, progress, timeout/expiry, cancel/retry/resume refs, actor, timestamps, correlation.

## 21. Permissions fonctionnelles
Progress read, cancel/retry/resume request, sensitive status/output refs, provenance, cross-tenant deny.

## 22. Limites et erreurs
Partial ≠ complete; cancel ≠ rollback; timeout ≠ confirmed termination; retry ≠ duplicate-free guarantee; resume ≠ new collection automatically; no exactly-once claim.

## 23. Métriques
Queue/wait duration, partial/failure/timeout/cancel, unknown-state duration, retry/resume rates, false completion target zero.

## 24. Classification de livraison
`draft / defined / planned`; no scheduler, queue engine or state-machine implementation.

## 25. Critères d’acceptation
**Given** a multi-item collection partially succeeds, **When** status is aggregated, **Then** result is `partial` and failed items remain explicit.

**Given** cancel is requested, **When** no target acknowledgement arrives, **Then** state remains `cancel-requested/status-unknown`, not cancelled.

**Given** an offline queued operation reconnects, **When** resume eligibility is computed, **Then** no resume starts without an explicit authorized control path.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; final retry/idempotency bridge deferred.

## 27. Consommateurs documentaires
EPT-4 Collection, Investigate Collection Job, Shared Jobs, Govern/Studio provenance, Quality.
