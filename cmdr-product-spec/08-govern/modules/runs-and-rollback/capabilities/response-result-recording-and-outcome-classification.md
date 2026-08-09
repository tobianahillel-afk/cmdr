---
id: CAP-GOV-032
title: Response Result Recording and Outcome Classification
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-004, REQ-PROD-005, REQ-PROD-006, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-032 — Response Result Recording and Outcome Classification

## 1. Définition
Créer, revoir et finaliser le **canonical Result** Govern d’un Response Run à partir de l’intended outcome, execution/runtime outcomes, verification, rollback/recovery outcomes, failures, partial success, residual risk, unintended effects and follow-up needs, sans confondre Result avec raw executor output, Evidence, Finding ou Decision.

## 2. Problème utilisateur
Les executors produisent leurs propres résultats techniques, parfois contradictoires ou incomplets. Les consommateurs opérationnels ont besoin d’un outcome canonique qui exprime ce qui a réellement été tenté, observé et vérifié, tout en conservant les sources et incertitudes.

## 3. Objectifs
- relier Result à exact Response Run et Decision ;
- représenter exact action/targets and intended outcome ;
- agréger execution outcome sans écraser raw technical refs ;
- intégrer verification, rollback/recovery and residual risk ;
- conserver partial success/failures/unintended effects ;
- classifier l’outcome sans fausse certitude ;
- transmettre follow-up needs et provenance aux consommateurs.

## 4. Non-objectifs
Ne pas requalifier Evidence/Finding, modifier Incident automatiquement, réécrire Decision, copier un raw executor payload comme Result, considérer runtime success comme success vérifié, supprimer les échecs/partials, finaliser un Report ou créer GOV-3 metrics/audit capabilities.

## 5. Propriétaire
Govern owns canonical Result semantics and lifecycle. Technical owners retain raw outputs/results. Command/Investigate consume Result projections while retaining their objects. Shared Reporting renders/exports under its own ownership.

## 6. Utilisateurs
Principal : Response Operator / Govern Reviewer. Secondaires : Decision Maker, Incident Commander, Investigator, Runtime/Endpoint/Studio Operator as source contributors, Business/Service Owner as contextual consumer, Auditor.

## 7. Conditions d’entrée
Response Run has a terminal, stopped/cancelled, rollback/recovery or verification disposition sufficient to explain the current outcome; technical refs and verification/rollback contexts are available or explicitly unavailable; permission to create/review Result.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Response Run + Decision | Govern | canonical lineage/intent | oui | exact pinned versions | no Result |
| exact action/targets/intended outcome | Run/Plan/Decision | result subject | oui | historical pinned context | no Result |
| execution/runtime outcomes | CAP-GOV-026/027 | technical execution summary + refs | oui where execution occurred | correlated terminal/known state | incomplete/inconclusive |
| Verification Assessment | CAP-GOV-029 | verified outcome/residual risk | expected where verification required | exact version | verification-pending/inconclusive |
| Rollback/Recovery outcomes | CAP-GOV-030/031 | reverse/recovery disposition | when applicable | current finalized/known | unknown/partial explicit |
| errors/partial/unintended effects | CAP-GOV-027/029/031 | outcome limitations | when present | current records | none invented |
| follow-up needs | reviewer/assessment/handoff context | next-action context | non | current | no follow-up invented |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run / Decision / Execution Plan | Govern | lineage, action, targets, intended outcome | read |
| Runtime Status / Error / Retry / Compensation | Govern + technical sources | execution facts and raw refs | read |
| Verification/Residual Risk Assessment | Govern | verified/unknown outcomes | read |
| Response Rollback / Recovery Context | Govern | rollback/recovery outcome | read |
| Incident | Command | downstream operational context only | read/link |
| Case/Finding/Evidence | Investigate | downstream/source context only | read/link, no mutation/requalification |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Result | create/update-review/finalize/dispute/supersede | Govern | source refs/uncertainty/partial states preserved |
| outcome classification | set/review/supersede | Govern | cannot contradict verification evidence silently |
| follow-up need/handoff relation | add/update/close | Govern local context | no destination-object mutation |
| Response Run | link Result/closure context | Govern | Run history preserved |

## 11. Fonctionnalités
Assemble Result draft; pin Run/Decision/action/targets/intended outcome; summarize per-target/per-step execution; retain raw technical refs; include verification outcome and residual risk; include rollback/recovery state; capture failures/partial/unintended effects; classify outcome; compare draft versions; dispute/annotate; final review; supersede with history; create downstream handoff candidates without source-object mutation.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect Result/source refs | authorized user | Result | 0 | read | exact outcome context | non |
| validate classification vs source assessments | reviewer | Result | 1 | source refs available | mismatch/missing list | non |
| create/update Result draft/follow-up needs | Response Operator | Result | 2 | Run/Decision lineage | versioned draft | OPEN-013 |
| finalize/supersede Result | authorized Govern Reviewer | Result | 2/3 depending policy | sufficient source/verification context | canonical outcome record | governed |
| mutate Finding/Evidence/Decision/Incident silently | none | source objects | — | forbidden | no effect | destination owner only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate per-target/step outcomes | oui | structured aggregation | oui | sourced summary | outcome matrix |
| validate classification against verification | oui | explicit consistency rules | oui | explanation | verification/result checklist |
| draft narrative/follow-up needs | oui | templates | oui | draft/recommendation | structured form |
| detect missing/contradictory source refs | oui | oui | oui | summary | source completeness table |
| declare success without evidence/rewrite source object | accountable reviewer only within rules | no autonomous | no silent | prohibited | explicit review |

## 14. États fonctionnels
Result lifecycle may include `draft`, `review-required`, `verification-pending`, `ready-for-review`, `finalized`, `disputed`, `superseded`. Functional outcome classifications include `success`, `partial-success`, `failed`, `cancelled`, `stopped`, `rolled-back`, `recovered`, `verification-failed`, `inconclusive`. Classification is separate from lifecycle state.

## 15. États d’interface
Loading preserves Run/source refs ; Empty means no Result yet ; Partial explicitly lists missing verification/runtime/rollback data ; Error preserves latest valid draft ; Offline cannot finalize if required current checks cannot be guaranteed ; Permission denied masks restricted technical/source details ; Stale identifies late source updates requiring review.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| canonical Result | Govern object | CAP-GOV-033/Command/Investigate/Reporting | Run/Decision/action/targets/outcomes/provenance linked |
| outcome classification | Result field/disposition | Command/Investigate/GOV-3 future | verification/partial/residual context preserved |
| follow-up needs | handoff candidates | CAP-GOV-033 | no silent destination mutation |
| Result supersession/dispute event | history relation | Audit/consumers | prior Result remains resolvable |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-026/027/029/031 | execution/verification/rollback context ready | CAP-GOV-032 | Run/source refs/outcomes/residual risk | Runs & Rollback |
| CAP-GOV-032 | Result finalized | CAP-GOV-033 | Result/version/classification/follow-up needs | Result |
| CAP-GOV-032 | verification insufficient/disputed | CAP-GOV-028/029 | missing/contradictory criteria/source refs | same Result draft |
| CAP-GOV-032 | new source outcome arrives | CAP-GOV-032 | late source ref and affected classification | same Result lineage |

## 18. Dépendances
CAP-GOV-015/022/026..031/033, Result object, Command Incident, Investigate Case/Finding/Evidence, Studio/Endpoint/provider technical refs, Shared Trace/Reporting/Linking/Versioning, Security permissions, OPEN-013/015.

## 19. Source de vérité
Govern Result is the canonical response-outcome object. It does not become source of raw executor truth, Incident, Finding or Evidence. Decision remains authority history; Result records what occurred/was verified, never rewrites what was authorized.

## 20. Provenance et audit
Record Result id/version/lifecycle/classification, Run/Decision/Plan/Playbook refs, exact targets/action/intended outcome, per-target/step execution summary, raw technical refs, error/retry/compensation refs, Verification/Residual Risk assessment, Rollback/Recovery refs, failures/partials/unintended effects, uncertainty/missing data, follow-up needs, author/reviewer/dispute, AI draft provenance and timestamps.

## 21. Permissions fonctionnelles
Result create/update/read/review/finalize/dispute/supersede; restricted technical/source context read; cross-product handoff prepare; automated draft request; provenance export. Source-object transitions require destination owner permission and are not implied by Result finalization.

## 22. Limites et erreurs
Missing verification, unknown executor state, conflicting technical outputs, partial rollback, unresolved residual risk or late data can require `partial-success`, `verification-failed` or `inconclusive`. No raw executor payload is promoted wholesale and no failure is erased to simplify classification.

## 23. Métriques
Results by classification; time Run terminal→Result; Results missing verification when required; runtime-success→non-success Result mismatches; disputed/superseded Results; raw technical result promoted as canonical Result — target zero.

## 24. Classification de livraison
`defined` / `planned`; no Result schema/API/reporting implementation or automated outcome engine selected.

## 25. Critères d’acceptation
**Given** an Endpoint technical result says success but verification fails, **When** Result is finalized, **Then** canonical Result cannot be `success`; raw technical success remains linked as a source and verification failure is visible.

**Given** a Run succeeded on some targets and failed on others, **When** Result is recorded, **Then** per-target outcomes remain distinct and overall outcome is `partial-success` or another evidence-supported classification, not fabricated success.

**Given** a rollback is technically complete but recovery verification is incomplete, **When** Result is reviewed, **Then** the classification cannot imply full recovery without supporting evidence.

**Given** no AI, **When** Result is produced, **Then** structured outcome matrices, deterministic consistency checks and human review provide full functionality.

## 26. Questions ouvertes
OPEN-013 remains for final governance of Result mutations/finalization where relevant; OPEN-015 remains for technical-run provenance bridge. Final object schema/state permissions remain future; no new OPEN.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-033, Command, Investigate, Shared Reporting, future Objects/Permissions/Screens/Journeys/Technique, GOV-2 report and GOV-3 Audit Trail/Response Metrics.
