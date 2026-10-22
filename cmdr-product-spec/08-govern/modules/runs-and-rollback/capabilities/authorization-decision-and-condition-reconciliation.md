---
id: CAP-GOV-021
title: Authorization, Decision and Condition Reconciliation
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-021 — Authorization, Decision and Condition Reconciliation

## 1. Définition
Réconcilier immédiatement avant création/démarrage de Run la Decision finalisée, son exact version, disposition, action, targets, allowed/prohibited scope, conditions, start window/expiry, Approval/Exception/emergency validity, Playbook compatibility et target readiness afin de déterminer si la préparation du Response Run est encore autorisée.

## 2. Problème utilisateur
Une Decision valide au moment de sa création peut être expirée ou dépassée par un changement de target, Playbook, Approval ou exception au moment de l’exécution. Sans revalidation, une autorisation historique peut être utilisée hors fenêtre ou hors scope.

## 3. Objectifs
- vérifier existence/finalisation/disposition de la Decision ;
- pinner exact Decision/request/handoff/plan/Playbook versions ;
- comparer exact targets, allowed/prohibited scope et conditions ;
- revalider time bounds/expiry et required Approvals/Exceptions/emergency context ;
- intégrer target drift/readiness et Playbook compatibility ;
- produire une disposition explicable `authorized-for-run-preparation` ou blocker.

## 4. Non-objectifs
Ne pas prendre une nouvelle Decision, renouveler une Approval/Exception, élargir le scope, modifier Playbook/target, démarrer le Run, appeler l’executor, auto-réautoriser après expiry ou finaliser le moteur d’autorisation.

## 5. Propriétaire
Govern owns execution-time reconciliation of its Decision/Approval/Exception records. Security owns authority/SoD/step-up policy. Source owners retain current target/runtime facts; GOV-1 remains owner of the historical Decision semantics.

## 6. Utilisateurs
Principal : Response Operator / Govern Reviewer. Secondaires : Decision Maker, Security/Authority Reviewer, Approver, Endpoint/Runtime Operator, Auditor.

## 7. Conditions d’entrée
Current Decision/Handoff, compatible Playbook/version, Execution Plan, current Readiness Assessment, access to relevant Approval/Exception/emergency validity and permission to prepare a Response Run.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| finalized Decision/version | CAP-GOV-015 | execution authority | oui | current/effective | blocked |
| Execution Handoff/Plan | CAP-GOV-016/019 | exact bounds/intention | oui | current versions | blocked |
| Playbook compatibility | CAP-GOV-018 | procedure compatibility | oui | current version | re-review-required |
| Target readiness/drift | CAP-GOV-020 | current execution context | oui | fresh enough for start policy | blocked/limited |
| Approval/authority/SoD | CAP-GOV-009..013 | authority prerequisites | selon Decision | currently valid where required | approval-invalid/blocked |
| Policy/Exception/emergency | CAP-GOV-007/008/013 | validity/conditions | selon Decision | effective/current | exception-expired/blocked |
| current time/start window | trusted time context | temporal condition | oui for bounded Decision | current | expired/condition-unsatisfied |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision / Conditions / Expiration | Govern | authority/bounds | read/reconcile |
| Approval / Exception / Emergency Context | Govern | validity/scope/conditions | read/revalidate |
| Execution Plan / Handoff | Govern | exact action/target/scope | read |
| Playbook Compatibility Review | Govern | version compatibility | read |
| Readiness Assessment | Govern/source projections | current target/runtime state | read |
| Security authority policy | Security | step-up/SoD/validity rules | read/evaluate |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Authorization Reconciliation | create/rerun/version/supersede | Govern local concept | exact source versions and time required |
| Run-preparation disposition | record/block/invalidate | Govern | not a new Decision |
| Execution Plan | attach reconciliation ref | Govern | no silent plan mutation |
| Decision/Approval/Exception | no rewriting | Govern GOV-1 records | history/effective state preserved |

## 11. Fonctionnalités
Validate Decision existence/final state/disposition; compare requested vs approved action; compare Plan targets/scope to allowed/prohibited scope; evaluate each Decision condition; check start window/expiry; verify Approval/Exception/emergency validity; verify Playbook compatibility and target readiness; detect target drift; expose unresolved conditions; determine whether a new Decision/review is required.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect authority chain | operator/reviewer | reconciliation inputs | 0 | read | exact authority context | non |
| run reconciliation | Govern Reviewer | reconciliation | 1 | source versions accessible | deterministic disposition | non |
| annotate/return for review | reviewer | reconciliation | 2 | mismatch/reason | return/review context | OPEN-013 |
| accept authorized-for-run-preparation | authorized reviewer | disposition | 2 | all mandatory constraints satisfied | Run creation may proceed | OPEN-013 |
| renew/override authority/start | none here | Decision/Run | — | forbidden | no effect | GOV-1/new authority or CAP-GOV-023 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare exact ids/versions/scope | oui | oui | oui | explain only | deterministic diff |
| evaluate expiry/start window | oui | oui | oui | no need | time checks |
| validate Approval/Exception references | oui | explicit validity rules | oui | summary | authority checklist |
| explain mismatches | oui | reason codes | oui | sourced explanation | mismatch table |
| authorize new scope/Decision/start | authorized separate path | no autonomous | no | prohibited | explicit GOV-1/Run control |

## 14. États fonctionnels
Outcomes: `authorized-for-run-preparation`, `blocked`, `expired`, `condition-unsatisfied`, `scope-mismatch`, `target-mismatch`, `approval-invalid`, `authority-invalid`, `exception-expired`, `emergency-context-expired`, `playbook-re-review-required`, `readiness-stale`, `re-decision-required`, `superseded`.

## 15. États d’interface
Loading preserves all pinned versions ; Empty means missing authority source ; Partial names unevaluated conditions ; Error blocks authoritative progression ; Offline never treats cached authority as current unless policy explicitly guarantees it ; Permission denied masks protected rationale ; Stale requires rerun.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Authorization Reconciliation | versioned result | CAP-GOV-022/023 | exact Decision/Plan/Playbook/readiness/authority versions |
| authorized-for-run-preparation | disposition | CAP-GOV-022 | permission to prepare Run only, not started |
| blocker/re-decision requirement | governance event | GOV-1/Response Operator | exact mismatch/reason and return origin |
| validity snapshot refs | provenance set | Run/Result/Audit | source records not rewritten |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-018..020 | compatible plan/readiness available | CAP-GOV-021 | exact review/plan/readiness refs | planning/readiness |
| CAP-GOV-021 | authorized for prep | CAP-GOV-022 | reconciliation + exact authority bounds | reconciliation |
| CAP-GOV-021 | expired/mismatch | CAP-GOV-014/015 or source review | reason/current vs approved diff | Runs & Rollback |
| CAP-GOV-023 | pre-start recheck | CAP-GOV-021 | latest time/readiness/validity | control action |

## 18. Dépendances
CAP-GOV-007..020/022/023, Security authority/SoD/step-up, trusted time, Shared Trace/Versioning/Linking, source target/runtime facts, OPEN-007/013/015.

## 19. Source de vérité
GOV-1 Decision/Approval/Exception records remain source of authority semantics; Govern GOV-2 owns the execution-time reconciliation result. Reconciliation never extends authority or rewrites historical records.

## 20. Provenance et audit
Record every pinned id/version, disposition, action, targets/scope diff, conditions and evaluations, timestamps/start window/expiry, Approval/Exception/emergency validity, Playbook compatibility, readiness/drift, reviewer, rule/version, automation explanation and return path.

## 21. Permissions fonctionnelles
Authorization reconcile/read; restricted Decision/Approval/Exception context; readiness/compatibility read; return for re-decision; provenance export. Step-up/SoD may be required before effectful Run start, but final RBAC/ABAC remains future.

## 22. Limites et erreurs
Unavailable authority source, stale readiness, clock uncertainty, conflicting Decision versions, expired Approval/Exception, target drift or inaccessible condition input cannot default to authorized. The safe outcome is explicit blocked/unknown/re-review.

## 23. Métriques
Reconciliations blocked by expiry/drift/condition; re-decision frequency; stale readiness at start; invalid Approval/Exception caught; executions later found outside reconciled scope — target zero.

## 24. Classification de livraison
`defined` / `planned`; no authorization engine/API/protocol or runtime enforcement selected.

## 25. Critères d’acceptation
**Given** a Decision expired before start, **When** reconciliation runs, **Then** Run cannot start, the historical Decision remains resolvable and no silent reauthorization occurs.

**Given** target drift expands the current target beyond approved scope, **When** reconciled, **Then** the difference is visible, scope remains unchanged and re-review/re-decision is required.

**Given** no AI, **When** authority is reconciled, **Then** deterministic version/scope/time/validity checks and human review provide the full function.

## 26. Questions ouvertes
OPEN-007, OPEN-013 and OPEN-015 remain open. Final effectful start authority/step-up rules remain future Security/Permissions implementation; GOV-2 creates no new OPEN.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-022/023/025, GOV-1 Decision/Approval pathways, future Objects/Permissions/Screens/Journeys/Technique, GOV-2 report and GOV-3 audit/metrics.
