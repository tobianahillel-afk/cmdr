---
id: CAP-GOV-023
title: Response Run Scheduling, Start, Pause, Stop and Cancellation Control
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-023 — Response Run Scheduling, Start, Pause, Stop and Cancellation Control

## 1. Définition
Gouverner les contrôles temporels d’un Response Run — immediate/scheduled start, pause, resume, stop and cancellation requests — avec revalidation de Decision/conditions/readiness au moment pertinent, distinction request/confirmation et respect des concurrency/max-target/maintenance/expiry bounds.

## 2. Problème utilisateur
Un Run peut être créé longtemps avant son démarrage, tandis que Decision, target ou executor changent. De même, un bouton Stop/Cancel ne prouve pas que l’executor a cessé d’agir. Sans semantics explicites, l’UI et les opérateurs peuvent confondre intention et état confirmé.

## 3. Objectifs
- supporter immediate et scheduled start avec start window/expiry ;
- réconcilier authority/readiness avant un effectful start ;
- distinguer `start requested` et `start confirmed` ;
- encadrer pause/resume avec revalidation ;
- enregistrer stop/cancel requests et confirmations distinctement ;
- appliquer concurrency, maximum target, maintenance and availability constraints ;
- préserver no-effect cancellation vs partial-effect stop/rollback review.

## 4. Non-objectifs
Ne pas implémenter un scheduler, envoyer une commande propriétaire, considérer cancel comme rollback, auto-réautoriser une Decision expirée, étendre les targets, ignorer stop confirmation, finaliser RBAC/ABAC ou choisir un runtime.

## 5. Propriétaire
Govern / Runs & Rollback owns Run control intent and functional lifecycle. Technical executor owns actual start/pause/stop/cancel primitive and confirmation. Shared may provide scheduling/job mechanisms but does not own Response Run semantics.

## 6. Utilisateurs
Principal : authorized Response Operator. Secondaires : Govern Reviewer, Decision Maker, Runtime/Endpoint Operator, Studio Operator, Incident Commander as observer, Auditor.

## 7. Conditions d’entrée
Existing Response Run, current authorization reconciliation, fresh-enough readiness for start/resume, exact bounds/conditions, schedule/start window where applicable and effectful-operation permission/authority.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Response Run/state | CAP-GOV-022 | control subject | oui | current Run version | action denied |
| authorization reconciliation | CAP-GOV-021 | current authority | oui for start/resume/effectful retry | rechecked at control time | blocked |
| readiness/target state | CAP-GOV-020 | current execution readiness | oui for start/resume | policy-defined freshness | stale/blocked |
| schedule/start window/expiry | Decision/Run | temporal bound | selon Run | current time | blocked/unknown |
| concurrency/max-target limits | Decision/Plan/Playbook/Policy | execution bounds | selon action | current pinned versions | constraint-unsatisfied |
| maintenance/executor availability | Settings/Endpoint/Studio/provider | runtime restriction | selon executor | current projection | start denied/pending |
| executor confirmations | technical owner | actual control response | after request | correlated to request | requested state remains unconfirmed |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run | Govern | state/bounds/targets/owner | read/control |
| Decision/Conditions | Govern | start window/expiry/limits | read/reconcile |
| Readiness/Authorization Reconciliation | Govern | current preconditions | read/rerun |
| Execution Plan/Playbook | Govern | concurrency/stop constraints | read |
| Workflow/Automation Run | Studio | control/confirmation projection | read/link |
| Endpoint/runtime state | technical owner | availability/control response | restricted read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Run control request | create/version/cancel where safe | Govern local event | intent ≠ confirmed state |
| Response Run | transition after validated request/confirmation | Govern | confirmations correlated |
| schedule context | create/update/cancel | Govern local concept | schedule ≠ start |
| technical execution object | no local mutation except authorized handoff | source owner | confirmation consumed by reconciliation |

## 11. Fonctionnalités
Schedule start; validate start window/expiry; check max targets/concurrency/maintenance; rerun authorization/readiness before effectful start; create start request; consume acceptance/started confirmation; request pause/resume with relevant revalidation; request stop/cancel; distinguish pre-effect cancellation from partial-effect stop; capture denial/reason; invalidate scheduled run when Decision expires; preserve operator/executor provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect schedule/control state | user | Run | 0 | read | status | non |
| validate start/resume preconditions | operator/reviewer | reconciliation/readiness | 1 | sources current | pass/block reasons | non |
| schedule/cancel schedule/pause request | authorized operator | Run control | 2 | within Decision conditions | versioned request | OPEN-013 |
| start/resume/stop/cancel effectful Run | authorized operator | Run | 3/4 depending underlying action | reconciliation + authority + safe bounds | technical handoff request | governed |
| expand scope or bypass expiry | none | Run/Decision | — | forbidden | no effect | new governance required |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| calculate schedule/expiry eligibility | oui | oui | oui | no need | timestamp rules |
| rerun authorization/readiness | oui | oui | oui | explain only | deterministic checks |
| propose safe pause/stop timing | oui | explicit step constraints | oui | recommendation | step/status table |
| detect stale schedule | oui | version/time rules | oui | explanation | schedule validator |
| start/stop/rollback authority | accountable governed path | contract-controlled | only where approved policy applies | no autonomous authority | explicit operator/Decision path |

## 14. États fonctionnels
Control-related states include `ready`, `scheduled`, `start-window-pending`, `starting`, `running`, `pause-requested`, `paused`, `resume-requested`, `stop-requested`, `cancel-requested`, `cancelled-before-effect`, `stopping`, `stopped`, `start-denied`, `expired-before-start`, `readiness-stale`, `executor-unavailable`, `control-unknown`. These supplement, not finalize, the Response Run machine.

## 15. États d’interface
Loading keeps Run/control intent ; Empty requires a Run ; Partial shows missing executor confirmations ; Error preserves intent and last confirmed state ; Offline forbids new effectful control unless future guarantees support it ; Permission denied masks restricted details ; Stale marks last confirmed runtime update.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| schedule/control request | Govern event | CAP-GOV-024/025/executor | exact Run/version/action/time/authority context |
| confirmed control transition | Run lifecycle event | CAP-GOV-022/026 | request and executor confirmation correlated |
| denied/expired condition | blocker | operator/GOV-1 | no silent start/reauthorization |
| partial-effect stop context | review condition | CAP-GOV-028..031 | affected scope/verification/rollback need explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-022 | Run ready/preparing | CAP-GOV-023 | Run + pinned bounds | Runs & Rollback |
| CAP-GOV-023 | pre-start check required | CAP-GOV-020/021 | Run/target/time/current refs | control context |
| CAP-GOV-023 | start/control authorized | CAP-GOV-024/025 | Run id/control intent/exact bounds | control request |
| executor | accepted/confirmed/denied | CAP-GOV-026/022 | technical status + correlation | same Run |
| stop after effect | confirmed stop | CAP-GOV-028/029/030 | affected targets/steps/status | Runs & Rollback |

## 18. Dépendances
CAP-GOV-020..022/024..031, Decision/Run, Studio/Endpoint/provider control primitives, Settings health, Shared Jobs/Notifications/Trace, Security authority/step-up/SoD, OPEN-007/013/015.

## 19. Source de vérité
Govern owns requested and reconciled Run control state; executor owns whether a technical action actually started/paused/stopped/cancelled. Response Run state changes only from attributable governance and executor evidence, not from UI intent alone.

## 20. Provenance et audit
Record Run/control request ids, actor/authority, Decision/reconciliation/readiness versions, schedule/window/expiry, limits, target set, prior confirmed state, request timestamp, executor handoff/correlation, acceptance/confirmation/denial, retry of control request, unknown state and human interventions.

## 21. Permissions fonctionnelles
Run schedule/start/pause/resume/stop/cancel, state read, authorization/readiness recheck, restricted executor status, provenance export. Effectful controls can require step-up/SoD; final RBAC/ABAC is not defined here.

## 22. Limites et erreurs
Decision expiry, stale readiness, target drift, exceeded target/concurrency bounds, maintenance restriction, unavailable executor, lost acknowledgement, duplicate request or control timeout leaves explicit denied/requested/unknown state. No state is fabricated from request submission.

## 23. Métriques
Scheduled runs expiring before start; start denied reasons; request→confirmed latency; pause/resume revalidation blocks; stop/cancel unknown states; partial-effect stops requiring rollback review; silent starts after expiry — target zero.

## 24. Classification de livraison
`defined` / `planned`; no scheduler/runtime/API/command or provider implementation selected.

## 25. Critères d’acceptation
**Given** a scheduled Run whose Decision expires before the start window, **When** start time arrives, **Then** the Run cannot start, Decision history remains and no silent reauthorization occurs.

**Given** cancel is requested before any effect and confirmed, **When** Run history is reviewed, **Then** it is `cancelled-before-effect`; rollback is not claimed.

**Given** no AI, **When** Run controls are managed, **Then** deterministic time/condition checks, explicit controls and executor confirmations provide full functionality.

## 26. Questions ouvertes
OPEN-007/013/015 remain open. Final effectful control authority, scheduler and executor acknowledgement semantics are future Security/Technique decisions; no new OPEN.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-024..031, Studio/Endpoint executors, Command/Investigate projections, future Objects/Permissions/Screens/Journeys/Technique, GOV-2 report and GOV-3 audit/metrics.
