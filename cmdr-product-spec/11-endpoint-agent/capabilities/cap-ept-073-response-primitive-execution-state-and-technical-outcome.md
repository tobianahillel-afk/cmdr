---
id: CAP-EPT-073
title: Response Primitive Execution State and Technical Outcome
product: endpoint-agent
module: containment
owner: Endpoint Agent Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-11
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-016, REQ-PROD-018, REQ-PROD-019, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002, REQ-SEC-004]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-EPT-073 — Response Primitive Execution State and Technical Outcome

## 1. Définition
Normaliser fonctionnellement le lifecycle target-side des primitives EPT-5 et leurs Technical Outcomes sans convertir ces états en Response Run success ou Govern Result.

## 2. Problème utilisateur
Chaque primitive peut exposer des états différents. Sans vocabulaire borné, `accepted`, `started`, `technically succeeded`, `verified` et `Result success` peuvent être confondus.

## 3. Objectifs
Normaliser requested/validating/blocked/ready/executing/partial/technically-succeeded/failed/unknown/cancelled/reversal/verification-pending ; préserver raw primitive state ; per-target outcomes ; timestamps/freshness ; output/error refs ; no business inference.

## 4. Non-objectifs
Remplacer CAP-EPT-062, créer Response Run/Result, définir final state machine/storage, auto-retry/rollback, masquer raw source state, or decide verification outcome.

## 5. Propriétaire
Endpoint owns local Technical Execution State/Outcome. Govern owns normalized Response Run reconciliation and canonical Result. Shared Job state remains generic and distinct.

## 6. Utilisateurs
Response Operator, Endpoint Operator, Govern Reviewer, Verification Reviewer, Auditor, Incident/Investigation consumers.

## 7. Conditions d’entrée
A valid EPT-5 primitive execution ref, target, Run/Step correlation, raw local state/output/error observation and permission to expose technical status.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| primitive execution | CAP-EPT-067..072/081 | execution subject | oui | current attempt | no outcome |
| Run/Step correlation | Govern | lineage | oui | pinned | orphaned technical state |
| raw state/progress | primitive owner | technical state | oui when emitted | source timestamp | unknown |
| output/error refs | EPT-4/primitive | technical result refs | terminal/progress dependent | exact source ref | partial/incomplete |
| cancellation/reversal refs | EPT-4/EPT-5 | control context | when applicable | current | explicit absent |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| primitive-specific execution | Endpoint | raw state/target/op | read |
| Technical Execution Output | Endpoint EPT-4 | output/error/timeout | read |
| Response Run/Step | Govern | correlation only | read/ref |
| Job/Trace | Shared | optional transport/audit refs | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Response Primitive Execution State | create/update/version | Endpoint | raw state retained |
| Technical Outcome | create/update/finalize technical | Endpoint | not Govern Result |
| per-target outcome | derive/update | Endpoint | partiality retained |
| Response Run/Result | aucune mutation | Govern | handoff only |

## 11. Fonctionnalités
Map primitive raw states to bounded technical states; preserve raw labels; track per-target progress; capture accepted/start/terminal and cancellation; retain timeout/unknown; distinguish reversal-pending/reversed and verification-pending; output deterministic reason codes and raw refs.

## 12. Actions utilisateur
Inspect Class 0; reconcile local state Class 1; request refresh Class 1/2 if no effect; starting/retrying/reversing remains with original effect capability and Govern authority.

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| map raw to bounded states | oui | oui | oui | explanation only | mapping table |
| aggregate per-target outcomes | oui | oui | oui | sourced summary | target matrix |
| detect stale/unknown/contradiction | oui | oui | oui | explain | freshness rules |
| declare Response Run/Result success | non | interdit | non | interdit | Govern reconciliation |

## 14. États fonctionnels
`requested`, `validating`, `blocked`, `authority-missing`, `ready`, `executing`, `partial`, `technically-succeeded`, `technically-failed`, `timed-out`, `unknown`, `cancel-requested`, `cancelled`, `reversal-pending`, `reversing`, `reversed`, `verification-pending`, `stale`, `superseded`.

## 15. États d’interface
No Screen ID. Requested/confirmed, raw/normalized, technical/verified and Endpoint/Govern states remain separately attributable.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| normalized technical state | Endpoint projection | Govern CAP-GOV-026 | raw source/time retained |
| Technical Outcome | Endpoint record | CAP-EPT-074/076/079 | not Result |
| partial/unknown condition | technical condition | CAP-EPT-076 | affected target explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| primitive family | state/output event | CAP-EPT-073 | primitive/target/raw status/time | primitive |
| CAP-EPT-073 | terminal/partial | CAP-EPT-074/076 | technical outcome/known state | Run |
| CAP-EPT-073 | Govern handoff | CAP-EPT-079 | technical status/output refs | Response Run |

## 18. Dépendances
CAP-EPT-062/067..072/081/074..080, CAP-GOV-022/026/032, Shared Jobs/Trace, OPEN-008/013/015.

## 19. Source de vérité
Endpoint is SOT of local technical state/outcome. Govern is SOT of Response Run reconciliation and Result. Raw primitive facts are never overwritten by normalized projection.

## 20. Provenance et audit
Primitive/execution/target ids, Run/Step refs, raw/normalized status, output/error refs, start/end/source/receipt times, operator, Agent/version, cancellation/reversal/verification refs, mapping version and correlation.

## 21. Permissions fonctionnelles
Technical outcome read, status refresh, sensitive output metadata, provenance read/export preparation, cross-tenant deny. Effect permissions remain primitive-specific.

## 22. Limites et erreurs
Technical success ≠ Response Run success; technical failure ≠ response failure automatically; technical output ≠ Result; cancel ≠ rollback; stop request ≠ confirmed stop; timeout/unknown never maps silently to success/fail.

## 23. Métriques
Time in state, partial/unknown/stale, raw→normalized mismatches, terminal technical success later mismatching verification, direct Result promotions target zero.

## 24. Classification de livraison
`draft / defined / planned`; no runtime/state engine, protocol, schema or implementation.

## 25. Critères d’acceptation
**Given** technical execution partially succeeds, **When** normalized, **Then** per-target partial outcome remains visible and overall state is not complete success.

**Given** status is lost after timeout, **When** reconciled, **Then** target state becomes unknown rather than failed or succeeded arbitrarily.

**Given** technical state reports success, **When** sent to Govern, **Then** it remains an input to verification/reconciliation and not a canonical Result.

## 26. Questions ouvertes
OPEN-008/013/015 remain open; final runtime status ontology/bridge is not selected.

## 27. Consommateurs documentaires
EPT-5 verification/drift/reversal/reconciliation, Govern Runs, Shared Trace, Quality and Audit.