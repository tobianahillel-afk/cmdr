---
id: CAP-GOV-031
title: Rollback Execution and Recovery Coordination
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-031 — Rollback Execution and Recovery Coordination

## 1. Définition
Coordonner l’exécution gouvernée d’un Rollback Plan autorisé, suivre son runtime/progress/partial/failure state, vérifier son effet et, lorsque le rollback est insuffisant ou impossible, coordonner une recovery path bornée ou une intervention manuelle sans prétendre que `rolled-back` signifie restauration exacte de l’état original.

## 2. Problème utilisateur
Une tentative de rollback peut échouer partiellement, restaurer seulement certains targets ou introduire un nouvel écart. Si le système traite une primitive inverse réussie comme une récupération complète, il masque le residual impact et peut clôturer prématurément l’incident.

## 3. Objectifs
- exiger Rollback Plan, authority and fresh preconditions avant start ;
- transmettre exact rollback target/action/scope à l’executor owner ;
- distinguer request, start, progress, partial, failure and technical completion ;
- vérifier le rollback avec des critères/source explicites ;
- déclencher recovery fallback/manual intervention quand requis ;
- conserver residual impact et relation au Result ;
- préserver la provenance complète de l’original Run à la recovery.

## 4. Non-objectifs
Ne pas définir de commande de rollback/recovery, garantir l’état original, auto-réessayer l’action initiale, élargir les targets, exécuter sans authority, réécrire Decision/Run original, créer un Result sans CAP-GOV-032 ou choisir un executor/runtime.

## 5. Propriétaire
Govern owns rollback/recovery governance, Response Rollback semantics and reconciliation. Endpoint/Studio/provider owners retain technical reverse/recovery primitives and raw outputs. Shared Recovery remains generic mechanism, not owner of Govern response semantics.

## 6. Utilisateurs
Principal : authorized Response Operator / Rollback Operator. Secondaires : Govern Reviewer, Decision Maker, Security/Authority Reviewer, Endpoint/Runtime/Studio Operator, Incident Commander, Investigator, Auditor.

## 7. Conditions d’entrée
Rollback Eligibility Assessment and Rollback Plan from CAP-GOV-030, explicit authority for effectful rollback, current target/readiness/preconditions, available technical executor or explicit manual-recovery route, and post-rollback Verification Plan/criteria.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Rollback Plan + Eligibility | CAP-GOV-030 | rollback contract | oui | current reviewed version | no rollback start |
| rollback authority/conditions | Govern/Security | effect authority | oui | current/effective | blocked |
| exact rollback targets/scope | Plan/original Run | bounded reverse effect | oui | current resolved targets | blocked/drift review |
| technical rollback/recovery capability | Endpoint/Studio/provider | executor primitive | oui for automated effect | current metadata | manual/recovery alternative |
| parameter/Secret References | Plan/Settings | technical inputs by reference | selon primitive | current refs | blocked |
| runtime status/output refs | technical owner | rollback/recovery observations | during execution | correlated/timestamped | status-unknown |
| post-rollback verification criteria | CAP-GOV-028/030 | expected restored/recovered state | oui | exact Plan/version | verification incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| original Response Run / Decision | Govern | original effect/authority context | read |
| Rollback Plan / Eligibility | Govern | exact rollback bounds/preconditions | read |
| Response Rollback | Govern | governed rollback identity/state | read/manage |
| Workflow/Automation Run/Tool Call | CMDR Studio | technical rollback/recovery refs | read/link only |
| Agent Command/Endpoint state | Endpoint Agent | technical effect/status refs | read/link only |
| Secret Reference / Integration / health | Platform Settings | input and executor metadata | restricted read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Response Rollback | create/transition/version/close | Govern | technical success ≠ verified restoration |
| Recovery Coordination Context | create/update/version | Govern local concept | fallback/manual route explicit |
| original Response Run | attach rollback/recovery relation/state | Govern | original history not overwritten |
| technical execution objects | no local mutation semantics | source owner | linked/reconciled only |

## 11. Fonctionnalités
Revalidate authority/target/preconditions; create governed Response Rollback; dispatch exact rollback/recovery intent via executor handoff; correlate technical run/call/command refs; track requested/started/running/partial/failed/completed states; preserve per-target rollback outcome; run post-rollback verification; detect residual effect/data loss; choose recovery/manual path via explicit review; stop unsafe continuation; link final rollback/recovery outcome to Result preparation.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect rollback/recovery state | authorized user | Response Rollback | 0 | read | current context | non |
| revalidate rollback preconditions | reviewer | rollback context | 1 | current source facts | pass/block reasons | non |
| prepare/annotate recovery fallback | reviewer/operator | Recovery Context | 2 | rollback inadequate/unsupported | versioned fallback | OPEN-013 |
| start/stop governed rollback/recovery | authorized operator | Response Rollback/recovery action | 3/4 per effect | authority + exact Plan/current preconditions | executor handoff | governed |
| auto-retry original action after rollback failure | none | original action | — | forbidden | no effect | new governed review required |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| recheck Plan/authority/current target | oui | oui | oui | explain only | checklist/diff |
| correlate rollback runtime state | oui | ids/status mapping | oui | sourced summary | status table |
| compare expected restored vs observed | oui | criteria rules | oui | explanation | verification matrix |
| propose recovery/manual next step | oui | policy/capability rules | oui | recommendation only | fallback checklist |
| authorize/start/retry effect | accountable governed path | contract-controlled | only when explicitly authorized | no autonomous AI | explicit controls |

## 14. États fonctionnels
`rollback-preparing`, `rollback-ready`, `rollback-start-requested`, `rolling-back`, `rollback-partial`, `rollback-failed`, `rollback-technically-completed`, `rollback-verification-pending`, `rolled-back-verified`, `recovery-required`, `recovery-preparing`, `recovering`, `recovery-partial`, `recovery-failed`, `recovered-verified`, `manual-intervention-required`, `status-unknown`, `closed`, `superseded`. `rolled-back`/`recovered` never guarantee exact original state without explicit verification.

## 15. États d’interface
Loading preserves original Run/Rollback context ; Empty no active rollback ; Partial exposes target-specific progress/gaps ; Error preserves last confirmed runtime/verification refs ; Offline prohibits new effects absent future guarantees ; Permission denied masks sensitive executor/output data ; Stale exposes last confirmed rollback/recovery state.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Response Rollback | governed reverse-action record | CAP-GOV-032/033 | original Run/Plan/authority/targets/provenance linked |
| rollback verification outcome | assessment | CAP-GOV-032 | technical completion and verified restoration distinguished |
| Recovery Coordination Context/outcome | recovery record | CAP-GOV-032/033 | expected vs observed recovery + limitations explicit |
| residual/manual-intervention condition | governance event | CAP-GOV-032/033/Command | no false recovery claim |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-030 | plan ready/authorized | CAP-GOV-031 | Rollback Plan/Eligibility/authority/current targets | rollback review |
| CAP-GOV-031 | technical handoff | Studio/Endpoint/provider | exact rollback/recovery intent + refs/limits | same Response Rollback |
| executor/runtime | progress/output/error | CAP-GOV-026/031 | correlated technical refs/status | source runtime |
| CAP-GOV-031 | verification/residual outcome | CAP-GOV-032 | rollback/recovery states/verification/limitations | Runs & Rollback |
| CAP-GOV-031 | further investigation/manual need | CAP-GOV-033 | residual impact/gaps/affected scope | same Run |

## 18. Dépendances
CAP-GOV-020..030/032/033, Response Rollback, Studio/Endpoint/provider technical rollback/recovery primitives, Settings Secret References/health, Shared Trace/Recovery/Notifications, Security authority, OPEN-007/008/013/015.

## 19. Source de vérité
Govern Response Rollback is source of governed rollback identity and lifecycle; technical owner remains source of raw execution facts; verification sources remain authoritative for observed restoration. Recovery Coordination Context is Govern-owned but does not claim exact restoration without evidence.

## 20. Provenance et audit
Record original Run/Decision/Plan, Rollback Plan/Eligibility/version, authority, target/scope/preconditions, executor owner/capability/version, Secret References only, dispatch/correlation, technical refs/status/output/error, per-target outcomes, verification criteria/observations, residual effects/data loss, recovery/manual path, human decisions, AI recommendations/dispositions and timestamps.

## 21. Permissions fonctionnelles
Rollback authorize/start/stop under required authority, Response Rollback read/manage, technical output inspect, recovery coordinate, verification execute/read, restricted parameter metadata, manual-intervention disposition and provenance export. No final RBAC/ABAC is defined.

## 22. Limites et erreurs
Rollback executor unavailable, target drift, partial effect knowledge, permission/secret denial, timeout, partial rollback, verification failure, data loss or recovery failure remain explicit. Rollback failure never automatically retries the original action; manual/re-decision paths remain possible.

## 23. Métriques
Rollback attempts/success/partial/failure; verification mismatch after technical rollback success; recovery/manual intervention frequency; targets not restored; rollback failures followed by unsafe original retry — target zero.

## 24. Classification de livraison
`defined` / `planned`; no rollback/recovery command, API, protocol, provider/runtime or execution implementation selected.

## 25. Critères d’acceptation
**Given** rollback completes technically on two of three targets, **When** status and verification are reconciled, **Then** partial rollback remains explicit and no full restoration claim is made.

**Given** rollback verification succeeds for the changed control but the original full state cannot be proven, **When** recovery outcome is recorded, **Then** `rolled-back` does not claim exact original-state restoration and limitations remain visible.

**Given** rollback fails, **When** next actions are considered, **Then** the original action is not automatically retried; recovery/manual/re-decision paths require explicit governance.

**Given** no AI, **When** rollback/recovery is coordinated, **Then** deterministic precondition/status/verification checks and explicit human controls provide full functionality.

## 26. Questions ouvertes
OPEN-007/008/013/015 remain open. Final technical rollback/recovery implementation and exact effect authority remain future Security/Technique work; no new OPEN.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-032/033, Command/Investigate handoffs, Studio/Endpoint/Settings executors, future Objects/Permissions/Screens/Technique, GOV-2 report and GOV-3 audit/metrics.
