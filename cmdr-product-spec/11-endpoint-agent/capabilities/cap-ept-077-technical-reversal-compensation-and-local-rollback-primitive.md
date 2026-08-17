---
id: CAP-EPT-077
title: Technical Reversal, Compensation and Local Rollback Primitive
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-077 — Technical Reversal, Compensation and Local Rollback Primitive

## 1. Définition
Définir une primitive Endpoint de reversal/compensation locale lorsque l’effet original est techniquement réversible, en conservant l’état avant reversal, l’intention Govern, l’exécution inverse, l’état après et la vérification, sans devenir le `Response Rollback` canonique.

## 2. Problème utilisateur
La présence d’une action inverse peut être prise pour une autorisation de rollback ou pour une garantie de restauration exacte. Un reversal peut être partiel ou impossible après drift.

## 3. Objectifs
Expose reversal capability/eligibility ; consume Govern Rollback Plan/authority ; pin original effect/target ; recheck current state ; execute bounded inverse primitive ; track partial/failure/unknown ; verify post-reversal target state ; preserve limitations/provenance.

## 4. Non-objectifs
Créer Govern Response Rollback, décider rollback, garantir original state restoration, auto-retry original action, define rollback engine/API/command or start recovery outside authorized scope.

## 5. Propriétaire
Endpoint owns technical inverse primitive and raw reversal facts. Govern owns rollback eligibility/plan, Response Rollback, recovery governance and Result. Shared Recovery remains generic.

## 6. Utilisateurs
Rollback Operator, Response Operator, Endpoint Operator, Govern Reviewer, Verification Reviewer, Auditor.

## 7. Conditions d’entrée
Original technical execution and known/partial effect state, Govern Rollback Plan/authority where required, reversal capability available, current target state/preconditions rechecked.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| original execution/effect | CAP-EPT-067..073/081 | reversal subject | oui | exact attempt | no reversal |
| current target state | CAP-EPT-075/076 | precondition | oui | fresh enough | unknown/block |
| reversal capability | primitive metadata | technical availability | oui | current | unavailable |
| Rollback Plan/authority ref | Govern | governed intent | oui for governed reversal | current | blocked |
| expected post-reversal state | Govern/technical plan | verification anchor | oui | pinned | incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| original Technical Execution/Outcome | Endpoint | effect/target/state | read |
| Target State Observation/Drift | Endpoint | current precondition | read |
| Rollback Plan/Response Rollback | Govern | exact scope/authority/lineage | read/ref |
| Endpoint Policy | Settings | reversal restrictions | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical Reversal | create/transition/version | Endpoint | exact original effect linked |
| Reversal Technical Outcome | create/update | Endpoint | not Response Rollback result |
| post-reversal state observation | create/link | Endpoint | technical fact only |
| original execution/Run | no rewrite | respective owner | history immutable |

## 11. Fonctionnalités
Check reversal support; bind exact original effect/target; compare current state; reject drifted/unknown when unsafe; consume authority; execute inverse primitive; track requested/start/partial/fail/unknown/completed; run technical verification; expose residual differences and recovery/manual needs.

## 12. Actions utilisateur
Inspect support/history Class 0; precondition/eligibility Class 1; reversal effect Class 3/4 according to original effect and Govern authority. No autonomous reversal.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| match effect to inverse primitive | oui | capability rules | oui | suggestion | reversal matrix |
| recheck current state | oui | oui | oui | explain | state diff |
| summarize reversal outcome | oui | structured | oui | oui | outcome table |
| authorize/start/choose rollback | non | Govern-controlled | non autonome | interdit | Govern Rollback path |

## 14. États fonctionnels
`not-assessed`, `not-supported`, `eligible`, `eligible-with-limitations`, `target-drifted`, `authority-required`, `reversal-requested`, `reversing`, `reversal-partial`, `reversal-failed`, `reversal-unknown`, `reversal-technically-completed`, `verification-pending`, `reversed-observed`, `manual-recovery-required`, `superseded`.

## 15. États d’interface
No Screen ID. Reversal capability, Govern rollback state and observed restoration are shown as distinct future concepts; no “undo” guarantee is implied.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Technical Reversal | Endpoint execution | CAP-EPT-073/075/079 | original effect/authority linked |
| Reversal Technical Outcome | technical fact | Govern CAP-GOV-031 | partial/failure explicit |
| post-reversal observation | verification fact | Govern | exact restoration not assumed |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Govern CAP-GOV-030/031 | authorized reversal | CAP-EPT-077 | Rollback Plan/original effect/target/authority | Response Rollback |
| CAP-EPT-077 | technical state | CAP-EPT-073/075 | reversal outcome/observed state | same rollback |
| CAP-EPT-077 | reconciliation | CAP-EPT-079 | outcome/limitations/verification | Govern |

## 18. Dépendances
CAP-EPT-067..076/078..080/081, CAP-GOV-030/031/032, Settings, Shared Recovery/Trace, Security, OPEN-007/008/013/015.

## 19. Source de vérité
Endpoint SOT for technical reversal execution/state. Govern SOT for Response Rollback eligibility/lifecycle/recovery/result. Neither rewrites original execution history.

## 20. Provenance et audit
Original Run/Step/technical execution, effect/target, pre-reversal state, reversal capability/version, Rollback Plan/authority, requested/observed post-state, operator, Agent/version, timestamps, partial/errors, verification/recovery refs.

## 21. Permissions fonctionnelles
Reversal support/read, technical reversal request, target-state read, verification, sensitive context, cross-tenant deny; Govern authority/step-up/SoD explicit.

## 22. Limites et erreurs
Technical reversal ≠ Govern rollback; local compensation ≠ Response Run rollback; rollback requested ≠ complete; technical reversal success ≠ Govern reconciliation; retry ≠ rollback; drift can invalidate reversal.

## 23. Métriques
Reversal available/eligible/blocked, partial/fail/unknown, verification mismatch, manual recovery needs, autonomous reversals target zero.

## 24. Classification de livraison
`draft / defined / planned`; no rollback/recovery command, engine, API/protocol or implementation.

## 25. Critères d’acceptation
**Given** Govern requests technical reversal with current authority, **When** target drift makes the inverse unsafe/invalid, **Then** Endpoint reports `target-drifted` and does not silently execute.

**Given** technical reversal partially fails, **When** status returns, **Then** partial state and affected targets remain explicit and Govern Response Rollback is not marked complete by Endpoint.

**Given** reversal completes technically, **When** post-state is observed, **Then** Endpoint reports facts/limits but does not claim full response rollback reconciliation.

## 26. Questions ouvertes
OPEN-007/008/013/015 remain open; final rollback authority/triggers and runtime bridge are not selected.

## 27. Consommateurs documentaires
Govern Rollback/Result, EPT-5 release/provenance, Shared Recovery, Security, Quality, Audit.