---
id: CAP-GOV-029
title: Post-Execution Verification and Residual Risk Assessment
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-002, REQ-PROD-004, REQ-PROD-005, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-029 — Post-Execution Verification and Residual Risk Assessment

## 1. Définition
Évaluer après ou pendant l’exécution l’**observed outcome** d’un Response Run face à son Verification Plan et à l’intended outcome, en séparant technical outcome, target state, verification evidence sources, residual risk, unintended effects et gaps pour produire une assessment vérifiable avant canonical Result ou rollback review.

## 2. Problème utilisateur
Un executor peut annoncer `success` alors que l’objectif n’est pas atteint, qu’une exposition subsiste ou qu’un effet indésirable apparaît. À l’inverse, des données de verification peuvent être incomplètes. Sans assessment séparée, runtime success devient abusivement business success.

## 3. Objectifs
- comparer intended vs observed outcome par target/criterion ;
- intégrer technical outcome sans lui donner autorité de Result ;
- vérifier current target state et sources cross-product autorisées ;
- exposer gaps, unavailable data, anomalies and unintended effects ;
- qualifier residual risk sans score universel ;
- déterminer besoin de rollback review/further investigation ;
- produire un outcome de verification traçable pour CAP-GOV-030/032.

## 4. Non-objectifs
Ne pas modifier Incident/Case/Finding/Evidence, fabriquer une Evidence, conclure automatiquement au rollback, réécrire Decision, transformer un raw output en Result, exécuter une action corrective ou imposer une formule universelle de risque/confidence.

## 5. Propriétaire
Govern owns post-execution verification and residual-risk assessment for the Response Run. Command/Investigate/technical owners retain their source objects and observations. Rollback governance is CAP-GOV-030/031; canonical Result is CAP-GOV-032.

## 6. Utilisateurs
Principal : Verification Reviewer. Secondaires : Response Operator, Govern Reviewer, Decision Maker, Incident Commander, Investigator, Endpoint/Runtime Operator, Business/Service Owner as contextual contributor, Auditor.

## 7. Conditions d’entrée
Verification Plan, Response Run and current runtime/error/partial state available; required verification observations are present, pending or explicitly unavailable; target/source access is authorized.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Verification Plan/version | CAP-GOV-028 | expected outcome/criteria | oui | exact active version | assessment blocked |
| Run/Step technical outcomes | CAP-GOV-026/027 | execution observations | oui where emitted | correlated/current | insufficient technical context |
| target state observations | target/Endpoint/provider source | observed state | selon criterion | timestamp visible | insufficient-data/inconclusive |
| Incident state projection | Command | operational context | non selon action | current/source time | omit, never infer |
| Finding/Case context | Investigate | source/follow-up context | non selon action | current/versioned | omit/restricted |
| verification source observations | source owners | criterion observations | oui per required criterion | within defined window | criterion insufficient-data |
| unintended-effect/residual-risk context | Run/errors/reviewers/source | safety assessment input | selon observation | current | no unsupported claim |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Verification Plan / Response Run | Govern | expected outcome/criteria/targets | read |
| Runtime Status / Error/Compensation refs | Govern + technical owner | observed technical state | read |
| target/Endpoint/integration state | source owner/Settings | observed target state | restricted read |
| Incident | Command | operational impact/state | read/link, no mutation |
| Case/Finding/Evidence | Investigate | contextual/restricted references | read/link, no requalification |
| Decision | Govern | original authorized objective/conditions | read, never rewrite |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Verification Assessment | create/update/version/supersede | Govern local concept | criterion/source/outcome explicit |
| Residual Risk Assessment | create/update/version | Govern local concept | dimensions/unknowns visible, no universal score |
| Unintended Effect record | create/link/disposition | Govern local concept | source and affected target required |
| Response Run | attach verification outcome/rollback-review need | Govern | technical success not overwritten; no Result yet |

## 11. Fonctionnalités
Collect observations within verification window; compare each criterion expected vs observed; preserve technical source status separately; identify verified/failed/partial/unknown criteria; assess target state and remaining exposure; record new anomalies/unintended effects; expose missing/restricted data; document residual risk/questions; compare assessment versions; determine rollback-review-required and further-investigation need without automatic execution.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect observations/criteria | reviewer | Verification Assessment | 0 | read | source-backed comparison | non |
| run deterministic criterion comparison | reviewer | assessment | 1 | criterion + observation | criterion outcome | non |
| annotate/dispute/mark unavailable data | reviewer | assessment | 2 | rationale/source | versioned assessment | OPEN-013 |
| declare assessment outcome | authorized Verification Reviewer | assessment | 2 | criteria evaluated/unknowns explicit | verification outcome | OPEN-013 |
| execute rollback/change Finding/Decision | none here | target/source objects | — | outside capability | no effect | CAP-GOV-030/031 or source owner |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare expected vs observed | oui | yes where criteria deterministic | oui | explanation | criterion table |
| aggregate per-target outcomes | oui | structured aggregation | oui | sourced summary | target matrix |
| identify residual-risk questions | oui | rules/checklists | oui | suggestion | risk checklist |
| propose rollback/further investigation | oui | policy/condition rules | oui | recommendation only | explicit review matrix |
| auto-declare success/rollback | accountable human/rules where explicitly authorized later | no model authority | no silent | prohibited | deterministic + human governance |

## 14. États fonctionnels
Assessment outcomes include exactly `verified-success`, `verified-partial`, `verification-failed`, `inconclusive`, `insufficient-data`, `adverse-effect-observed`, `rollback-review-required`; supporting states include `verification-pending`, `data-pending`, `disputed`, `stale`, `superseded`. Runtime success is not verified success.

## 15. États d’interface
Loading preserves criterion results ; Empty means verification not started ; Partial lists unevaluated/blocked criteria ; Error preserves prior observations ; Offline never fabricates current target state ; Permission denied masks restricted source data ; Stale marks observation outside required freshness/window.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Verification Assessment | versioned expected-vs-observed result | CAP-GOV-030/032 | per-criterion/target sources and unknowns explicit |
| Residual Risk Assessment | governance assessment | CAP-GOV-030/032/Command | risk/uncertainty source-backed, not execution failure by itself |
| unintended-effect record | review input | CAP-GOV-030/032/Investigate handoff | affected target/source visible |
| rollback-review-required | governance condition | CAP-GOV-030 | no automatic rollback implied |
| further-investigation need | handoff candidate | CAP-GOV-033/Investigate | no Finding/Evidence mutation |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-026/027/028 | terminal/partial observations + Plan | CAP-GOV-029 | Run/criteria/targets/source refs/errors | Runs & Rollback |
| CAP-GOV-029 | rollback review required | CAP-GOV-030 | assessment/residual risk/affected scope | verification |
| CAP-GOV-029 | assessment sufficient for outcome | CAP-GOV-032 | verification outcome/residual risk/technical refs | verification |
| CAP-GOV-029 | further investigation needed | CAP-GOV-033/Investigate | anomalies/gaps/source refs/return origin | same Run |

## 18. Dépendances
CAP-GOV-022/026..028/030/032/033, Command Incident, Investigate Case/Finding/Evidence, Endpoint/Settings/technical observations, Shared Trace/Linking/Reporting, Security access, OPEN-008/013/015.

## 19. Source de vérité
Govern Verification Assessment is source of response-outcome verification semantics; source systems remain authoritative for observations. Decision remains historical authority. Evidence/Finding remain Investigate-owned. Verification failure alone does not execute rollback.

## 20. Provenance et audit
Record Run/Decision/Verification Plan versions, each criterion, expected/observed value/condition, source owner/ref/version/timestamps/restrictions, technical outcome refs, target state, Incident/Finding context refs, unavailable data, residual-risk reasoning/unknowns, unintended effects, reviewer/disputes, automation/AI provenance and rollback/follow-up disposition.

## 21. Permissions fonctionnelles
Verification execute/read, restricted source observation read, residual-risk assessment create/update, rollback-review request, further-investigation handoff preparation, automated recommendation request, provenance export. No source-object mutation or automatic rollback.

## 22. Limites et erreurs
Conflicting observations, late/stale data, source unavailability, inaccessible Evidence/Finding, changed target after execution or circular executor-only evidence may yield inconclusive/insufficient-data. A verification failure is not automatically a rollback unless explicit policy/authority says so.

## 23. Métriques
Verified-success/partial/failure/inconclusive distribution; source-unavailable frequency; runtime-success→verification-failure mismatches; unintended effects; residual-risk follow-ups; automatic success without evidence — target zero.

## 24. Classification de livraison
`defined` / `planned`; no verification engine, risk score, telemetry API/provider/protocol or rollback implementation selected.

## 25. Critères d’acceptation
**Given** runtime reports success but a required verification criterion fails, **When** assessment completes, **Then** the verification outcome cannot be `verified-success`, the mismatch is visible and residual risk is recorded.

**Given** rollback is unavailable after a verification failure, **When** residual risk is assessed, **Then** no recovery claim is fabricated and manual/escalation/further-investigation needs remain explicit.

**Given** verification data is incomplete, **When** assessment is performed, **Then** the outcome is `insufficient-data` or `inconclusive` as supported; missing data is not converted to success.

**Given** no AI, **When** post-execution verification is performed, **Then** deterministic criterion comparisons, source/target tables and human review provide full functionality.

## 26. Questions ouvertes
OPEN-008 covers source/runtime availability, OPEN-013 assessment mutations/default governance and OPEN-015 run/source provenance. Final residual-risk model and rollback-trigger policy remain future; no new OPEN.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-030..033, Command/Investigate handoffs, Studio/Endpoint/Settings sources, future Objects/Permissions/Screens/Technique, GOV-2 report and GOV-3 audit/metrics.
