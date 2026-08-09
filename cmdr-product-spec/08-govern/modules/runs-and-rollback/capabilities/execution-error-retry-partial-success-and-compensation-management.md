---
id: CAP-GOV-027
title: Execution Error, Retry, Partial Success and Compensation Management
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-009, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-027 — Execution Error, Retry, Partial Success and Compensation Management

## 1. Définition
Gérer fonctionnellement les erreurs d’exécution, timeouts, pertes d’état, partial success, retry eligibility/bounds et compensations d’un Response Run, en conservant les cibles/steps affectés et l’autorité d’origine sans transformer un retry en réautorisation ni une compensation en rollback.

## 2. Problème utilisateur
Les réponses de production sont rarement binaires : certains targets peuvent réussir, d’autres échouer, un executor peut devenir indisponible ou perdre son état. Un retry automatique ou non borné peut répéter un effet, dépasser l’expiration de la Decision ou agir sur une cible devenue différente.

## 3. Objectifs
- classifier explicitement precondition/executor/target/permission/secret/timeout/step/output/status errors ;
- préserver partial target/step success au lieu de fabriquer un succès global ;
- calculer retry eligibility depuis Decision/Run/step/current state ;
- appliquer retry limit et delay conceptuel sans retry infini ;
- empêcher retry après expiry, target drift ou Playbook change non revu ;
- distinguer retry, compensation, rollback et manual intervention ;
- fournir des inputs de verification et rollback review.

## 4. Non-objectifs
Ne pas choisir un backoff algorithm final, exécuter un retry silencieux, réautoriser une Decision, élargir le scope, substituer target/Playbook, définir une commande de compensation/rollback, masquer un technical error ou produire le canonical Result avant verification.

## 5. Propriétaire
Govern / Runs & Rollback owns error disposition, retry governance and partial-success semantics for Response Runs. Studio/Endpoint/provider owners retain their technical retry/idempotency/error mechanisms. Govern compensation semantics remain distinct from Govern rollback capabilities CAP-GOV-030/031.

## 6. Utilisateurs
Principal : Response Operator. Secondaires : Govern Reviewer, Studio/Endpoint/Runtime Operator, Decision Maker, Security Reviewer, Incident Commander/Investigator as consumers, Auditor.

## 7. Conditions d’entrée
Existing Response Run/Step with a precondition/runtime/error/timeout/partial/unknown condition, correlated technical refs where available and current authority/readiness context sufficient to evaluate next action.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Run/Step current state | CAP-GOV-022/024/026 | failure context | oui | latest confirmed/unknown | no retry disposition |
| technical error/output refs | Studio/Endpoint/provider | source error evidence | selon failure | exact correlated ref/time | mark unavailable/unknown |
| affected targets/steps | Run/Runtime reconciliation | partial-success scope | oui where effect may exist | current known state | manual review/verification required |
| Decision/conditions/expiry | Govern | retry authority bounds | oui for effectful retry | current/effective | retry blocked |
| current readiness/target drift | CAP-GOV-020/021 | current target/authority context | oui for retry | fresh enough | retry blocked/re-review |
| retry limits/delay/idempotency context | Plan/Playbook/executor metadata | retry constraints | selon action | pinned/current | no implicit unlimited retry |
| compensation/rollback relevance | Plan/Step/Playbook/Decision | recovery relation | selon failure | pinned versions | explicit unknown/not-supported |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run / Response Step | Govern | state/targets/attempts | read/manage error disposition |
| Runtime Status / technical refs | Govern + source owner | raw error/status/output references | restricted read |
| Decision / Execution Plan / Playbook | Govern | expiry/scope/retry/compensation bounds | read |
| Automation Run / Tool Call | CMDR Studio | technical failure/retry refs | read/link only |
| Agent Command / Endpoint state | Endpoint Agent | technical failure/idempotency/status | read/link only |
| Secret Reference / integration health | Platform Settings | availability metadata only | restricted read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Execution Error record | create/update/classify/supersede | Govern local concept | source/error refs preserved |
| Retry Record | propose/authorize-request/record attempt/outcome | Govern local concept | retry ≠ new authorization automatically |
| Compensation Record | prepare/coordinate/record | Govern local concept | compensation ≠ rollback |
| Response Run/Step | mark partial/failed/unknown/manual-intervention | Govern | successful and failed targets remain distinct |

## 11. Fonctionnalités
Classify precondition/executor/target/permission/secret/timeout/step/inconsistent-output/lost-status/duplicate-execution errors; preserve per-target/per-step outcomes; detect duplicate execution candidate; compute retry eligibility; enforce attempt limit/conceptual delay/expiry/idempotency need; require readiness/authorization recheck before effectful retry; prepare compensation where procedure allows; require manual intervention when state is unsafe/unknown; route rollback-relevant failures to CAP-GOV-030.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect errors/attempts/affected targets | operator | Error/Retry record | 0 | read | exact failure context | non |
| evaluate retry eligibility | reviewer | Retry Record | 1 | current authority/state | eligible/blocked reason | non |
| propose retry/compensation/manual intervention | operator/reviewer | Run/record | 2 | rationale + bounds | versioned proposal | OPEN-013 |
| execute authorized effectful retry/compensation | authorized Run path | Execution Action | 3/4 per effect | revalidation + retry/compensation authority | executor handoff | governed |
| silently retry/expand target/bypass expiry | none | Run/target | — | forbidden | no action | prohibited |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| classify known error/status codes | oui | mapping/rules | oui | explanation | error catalog/table |
| calculate retry eligibility/remaining attempts | oui | oui | oui | no authority | deterministic rule view |
| summarize partial success | oui | per-target aggregation | oui | sourced summary | target/step matrix |
| propose retry/compensation | oui | policy/rules may suggest | oui | recommendation only | operator checklist |
| execute retry/compensation/rollback | governed explicit path | contract-controlled | only when authorized | **never silently** | explicit controls |

## 14. États fonctionnels
`error-observed`, `precondition-failed`, `executor-unavailable`, `target-unavailable`, `permission-denied`, `secret-unavailable`, `timed-out`, `step-failed`, `partial-success`, `inconsistent-output`, `status-lost`, `duplicate-execution-candidate`, `retry-eligible`, `retry-blocked`, `retry-scheduled`, `retrying`, `retry-limit-reached`, `compensation-candidate`, `compensating`, `compensated`, `manual-intervention-required`, `rollback-review-required`, `superseded`.

## 15. États d’interface
Loading preserves failed target/attempt context ; Empty means no current failure record ; Partial distinguishes unavailable outputs from failed targets ; Error in the UI never erases runtime evidence ; Offline disables new effectful retry unless future contract guarantees it ; Permission denied masks secret/output details ; Stale requires fresh readiness/authority before retry.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Execution Error record | classified failure context | operator/CAP-GOV-028..030 | source/target/step/attempt preserved |
| Retry Record/disposition | retry governance context | CAP-GOV-023..026 | bound to same authorized scope and attempt limits |
| partial-success matrix | per-target/per-step outcome | CAP-GOV-029/032 | successes/failures remain distinct |
| Compensation Record | bounded compensating-action context | CAP-GOV-029/032 | explicitly not rollback |
| rollback/manual-intervention need | governance event | CAP-GOV-030/operator | reason and affected scope explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-026 | error/timeout/partial/unknown | CAP-GOV-027 | Run/step/targets/raw refs/current state | runtime reconciliation |
| CAP-GOV-027 | effectful retry proposed | CAP-GOV-020/021/023 | retry attempt + current target/authority context | same Run |
| CAP-GOV-027 | partial/compensated state | CAP-GOV-028/029 | affected targets/steps/attempts/outputs | Runs & Rollback |
| CAP-GOV-027 | rollback review required | CAP-GOV-030 | failure/partial/adverse state + scope | same Run |

## 18. Dépendances
CAP-GOV-020..026/028..032, Studio retry/idempotency/error handling, Endpoint queueing/retry/idempotency, Settings secrets/health, Shared Jobs/Trace/Notifications, Security authority, OPEN-008/013/015.

## 19. Source de vérité
Technical owners remain source for raw errors/attempt execution. Govern owns the error classification relative to the Response Run, retry eligibility/disposition, partial-success semantics and compensation record. A retry never changes the source Decision.

## 20. Provenance et audit
Record Run/step/target, error class, raw source/ref/code/message metadata under classification, known effect state, retry attempt number/limit/delay rationale, idempotency context, Decision/readiness/reconciliation versions, target drift checks, compensation proposal/execution refs, human decision, AI recommendation/provenance and timestamps. Raw secrets remain absent.

## 21. Permissions fonctionnelles
Error/read technical output, retry request/proposal, effectful retry through governed execution permission, compensation prepare/coordinate, manual-intervention disposition, restricted output metadata and provenance export. No implicit authority or RBAC finalization.

## 22. Limites et erreurs
Unknown effect state, lost executor response, expired Decision, target/Playbook drift, exceeded retry limit, unavailable secret/executor, duplicate candidate or non-idempotent uncertainty must block automatic retry and require explicit review. Partial success is never collapsed into success.

## 23. Métriques
Errors by class; retries per Run/step; retry-limit reached; partial-success incidence; duplicate candidates blocked; retries prevented by expiry/drift; manual interventions; silent/infinite retries — target zero.

## 24. Classification de livraison
`defined` / `planned`; no retry engine, backoff algorithm, executor command, compensation implementation or provider/runtime selected.

## 25. Critères d’acceptation
**Given** a step times out after one allowed retry, **When** retry is considered, **Then** current Decision/readiness is rechecked, the retry remains in the same scope and no infinite retry occurs.

**Given** three targets where two succeeded and one failed, **When** Run state is reconciled, **Then** the two successes and one failure remain distinct and overall success is not fabricated.

**Given** AI recommends retry, **When** the recommendation is shown, **Then** reason/provenance is visible and no retry starts without deterministic eligibility and required authority.

**Given** no AI, **When** failures are handled, **Then** error mappings, retry counters, checklists, target matrices and human controls provide the complete function.

## 26. Questions ouvertes
OPEN-008 remains for actual executor support, OPEN-013 for default class-2 retry/compensation preparation, OPEN-015 for technical run bridge. Final retry/idempotency/compensation implementation remains future; no new OPEN.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-028..033, Studio/Endpoint/Settings source owners, Command/Investigate consumers, future Objects/Permissions/Screens/Journeys/Technique, GOV-2 report and GOV-3 audit/metrics.
