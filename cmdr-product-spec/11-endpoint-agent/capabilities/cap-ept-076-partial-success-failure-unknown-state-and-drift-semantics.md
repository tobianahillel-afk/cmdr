---
id: CAP-EPT-076
title: Partial Success, Failure, Unknown State and Drift Semantics
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-076 — Partial Success, Failure, Unknown State and Drift Semantics

## 1. Définition
Définir les semantics Endpoint pour execution partial, verification mismatch/unavailable, contradictory/stale/unknown target state et post-action drift, afin d’empêcher toute réduction binaire non sourcée.

## 2. Problème utilisateur
Les actions endpoint peuvent perdre le statut, n’affecter qu’une cible ou dériver après succès initial. Forcer success/failure masque le risque et conduit à de mauvaises décisions Govern.

## 3. Objectifs
Represent partial execution, contradictory/stale/unknown, lost status, target-changed-during-response, effect drift/post-action regression, retry eligibility, escalation need, verification unavailable and provenance.

## 4. Non-objectifs
Décider business failure/success, auto-retry, auto-rollback, infer cause, define final retry engine, or collapse unknown into failed.

## 5. Propriétaire
Endpoint owns technical uncertainty/drift semantics. Govern owns Response Run error/retry/verification/rollback decisions and Result classification.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Verification Reviewer, Govern Reviewer, Auditor.

## 7. Conditions d’entrée
Technical Outcome and/or Target State Observation exists, with affected target/scope and timestamps; missing source is explicitly representable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| technical outcome | CAP-EPT-073 | execution state | oui if execution | exact attempt | unknown execution |
| target-state observation | CAP-EPT-075 | verification fact | selon verification | timestamped | verification unavailable |
| expected state | CAP-EPT-066/074 | comparator | oui for drift | pinned | drift-unassessable |
| source health/freshness | EPT-1/2 | reliability | selon source | current/last seen | stale/unknown |
| Run/Step | Govern | correlation | oui | stable | orphaned technical condition |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Technical Outcome | Endpoint | per-target state/errors | read |
| Target State Observation | Endpoint | expected/observed/freshness | read |
| Response Run/Step | Govern | intended scope | read/ref |
| Agent health/telemetry state | Endpoint | source reliability | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Technical Uncertainty State | create/update | Endpoint | source-backed reason |
| Response Drift | create/update/version | Endpoint | observation, not causal verdict |
| Retry Eligibility | derive | Endpoint | eligibility only, no retry execution |
| Govern Run | aucune mutation directe | Govern | receives handoff |

## 11. Fonctionnalités
Detect partial per-target; identify contradictory source states; mark stale/unknown after timeout/disconnect; detect target identity/state change; compare later observation to achieved state for drift/regression; calculate retry eligibility under immutable scope; request escalation without executing retry/rollback.

## 12. Actions utilisateur
Inspect Class 0; deterministic uncertainty/drift evaluation Class 1; annotate/refresh Class 1/2 if no effect. Retry/reversal requires separate Govern-authorized path.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| detect partial/stale/contradiction | oui | oui | oui | explain | state rules |
| detect drift | oui | comparison | oui | summary | before/after matrix |
| suggest escalation/retry candidate | oui | eligibility rules | oui | recommendation | operator review |
| auto-retry/rollback/classify Result | non | interdit | non | interdit | Govern |

## 14. États fonctionnels
`complete-known`, `partial`, `failed-known`, `unknown`, `status-lost`, `contradictory`, `stale`, `verification-unavailable`, `verification-mismatch`, `target-changed`, `drift-detected`, `post-action-regression`, `retry-eligible`, `retry-ineligible`, `escalation-required`, `superseded`.

## 15. États d’interface
No Screen ID. Unknown, failure, partial and stale are separate visual semantics in future UI; no arbitrary fallback color/outcome is specified.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Technical Uncertainty State | Endpoint condition | Govern/Quality | reason/source/time |
| Response Drift | technical observation | CAP-EPT-077/079 | no failure causality assumed |
| Retry Eligibility | local assessment | Govern | retry ≠ execution/rollback |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-EPT-073/075 | partial/mismatch/unknown | CAP-EPT-076 | technical states/sources | Run |
| CAP-EPT-076 | reversal candidate | CAP-EPT-077 | prior/desired/current state + limitations | Govern approval retained |
| CAP-EPT-076 | reconciliation | CAP-EPT-079 | uncertainty/drift/retry eligibility | Response Run |

## 18. Dépendances
CAP-EPT-009/028/053/062/073..075/077..080, CAP-GOV-026..032, OPEN-008/013/015.

## 19. Source de vérité
Endpoint SOT of its technical uncertainty/drift observations. Govern SOT of what those conditions mean for Response Run/rollback/Result.

## 20. Provenance et audit
Execution/target/Run refs, per-target statuses, source timestamps, last confirmed state, contradiction refs, expected/observed comparison, drift timestamps, retry eligibility reason, escalation ref and correlation.

## 21. Permissions fonctionnelles
Uncertainty/drift read, no-effect refresh/verification, sensitive status metadata, provenance export, cross-tenant deny. Retry/effect authority separate.

## 22. Limites et erreurs
Unknown ≠ fail; partial ≠ complete; drift ≠ execution failure automatically; retry ≠ rollback; verification unavailable ≠ action failed; cancellation ≠ rollback.

## 23. Métriques
Partial/unknown/status-lost/contradictory/stale durations, drift/regression, retry eligibility, incorrect binary coercions target zero.

## 24. Classification de livraison
`draft / defined / planned`; no retry/drift engine, API, physical schema or implementation.

## 25. Critères d’acceptation
**Given** desired state is reached then later drifts, **When** a new observation is compared, **Then** `drift-detected` is recorded without rewriting the original execution outcome.

**Given** target state becomes unknown after timeout, **When** reconciled, **Then** unknown remains unknown and is not fabricated as success/failure.

**Given** technical execution partially succeeds, **When** escalation is prepared, **Then** affected targets and known/unknown states remain explicit and no full success is claimed.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; final retry policy and cross-product state precedence are not selected.

## 27. Consommateurs documentaires
CAP-EPT-077/079/080, Govern Verification/Rollback/Result, Security, Quality, Audit.