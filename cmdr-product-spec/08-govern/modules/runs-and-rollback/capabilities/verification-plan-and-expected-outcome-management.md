---
id: CAP-GOV-028
title: Verification Plan and Expected Outcome Management
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
# CAP-GOV-028 — Verification Plan and Expected Outcome Management

## 1. Définition
Définir et maintenir un **Verification Plan** pour un Response Run, explicitant l’expected outcome, les critères de success/failure/partial success, les sources d’observation, l’owner de verification, la fenêtre temporelle, les dépendances, les données indisponibles et les questions de residual risk avant de conclure sur l’effet réel.

## 2. Problème utilisateur
Un executor peut indiquer que son action a terminé sans démontrer que l’objectif de réponse a été atteint. Sans expected outcome et critères définis avant ou pendant l’exécution, l’équipe risque d’adapter rétrospectivement les critères au résultat observé ou de confondre technical completion avec succès.

## 3. Objectifs
- lier expected outcome à Decision, Run, exact targets et action ;
- définir critères observables de success/failure/partial success ;
- identifier verification sources/owners/time window ;
- distinguer source technique, Command/Investigate projection et observation indépendante ;
- documenter unavailable data, quality limitations et residual-risk questions ;
- conserver versions/diffs du plan ;
- préparer CAP-GOV-029 sans préjuger du Result.

## 4. Non-objectifs
Ne pas créer Evidence, modifier Finding/Incident, inventer une source de confirmation, définir un score universel de confiance, considérer Verification Plan comme Result, exécuter une action cible, déclencher rollback automatiquement ou choisir un provider/telemetry protocol.

## 5. Propriétaire
Govern owns verification-plan semantics for response execution. Technical owners, Command and Investigate retain ownership of source observations/objects. Verification consumes their authorized projections without requalifying Evidence/Finding.

## 6. Utilisateurs
Principal : Verification Reviewer / Response Operator. Secondaires : Govern Reviewer, Decision Maker, Incident Commander, Investigator, Endpoint/Runtime Operator, Service/Business Owner as contextual source, Auditor.

## 7. Conditions d’entrée
Decision/Handoff/Execution Plan and Response Run available; expected action/targets known; verification requirement from Decision/Playbook/Plan available or explicitly absent; permission to inspect proposed verification sources.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| intended/expected outcome | Decision/Handoff/Plan | response objective | oui | pinned Run lineage | plan incomplete |
| exact targets/action/scope | Response Run/Plan | verification subject | oui | current Run version | blocked |
| verification requirement/conditions | Decision/Playbook | required checks | selon action | pinned versions | explicit none/unknown |
| candidate success/failure observations | source owners/technical metadata | verification evidence sources | oui for verifiable outcome | source/version declared | insufficient-data risk |
| verification owner/time window | Govern/Decision/Policy | accountability/timing | oui where verification required | current | plan incomplete |
| source quality/freshness/restrictions | source owner | interpretation limits | selon source | timestamp/version | unknown/limited |
| residual-risk questions | Decision/Risk/Run/error context | unresolved safety questions | non initially | current context | none invented |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision / Execution Plan / Response Run | Govern | intended outcome/action/targets/conditions | read |
| Playbook | Govern | declared verification support/checkpoints | read |
| Runtime Status / Error records | Govern + technical sources | technical completion/failures | read |
| Incident | Command | operational state projection | restricted read/link |
| Case/Finding/Evidence | Investigate | contextual/source references | restricted read/link; no requalification |
| Endpoint/Integration/Telemetry source | source owner/Settings | verification source metadata | restricted read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Verification Plan | create/update/version/supersede | Govern local concept | expected outcome and criteria explicit |
| Verification Criterion | add/update/remove before assessment | Govern local concept | observable source/condition required |
| verification-source relation | link/version | source owner + Govern projection | reference does not become Evidence |
| Response Run | attach Plan/version/checkpoint refs | Govern | Plan ≠ Result |

## 11. Fonctionnalités
Capture expected outcome; define per-target/aggregate success, partial and failure conditions; choose source projections; identify independent/cross-product checks where appropriate; set owner/time window; document freshness/quality/restrictions; define unavailable-data behavior; identify residual-risk questions; attach step checkpoints; compare Plan versions; prevent retrospective silent criterion changes; hand off to post-execution assessment.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect expected outcome/sources | reviewer | Verification Plan | 0 | read | verification context | non |
| validate criterion/source completeness | reviewer | Plan | 1 | criteria/source metadata | pass/missing list | non |
| create/update criteria/window/owner | Verification Reviewer | Plan | 2 | Run/Decision context | versioned Plan | OPEN-013 |
| request authorized verification observation | reviewer/operator | source query/check | 1/2 if no-effect | source permission | observation request/ref | governed by source |
| modify target/Decision or declare Result | none here | target/Decision/Result | — | outside capability | no effect | other owner/capability |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| map expected outcome to candidate checks | oui | catalog/rules where known | oui | suggestions | verification checklist/catalog |
| validate criterion completeness | oui | oui | oui | explain missing | structured form |
| compare Plan versions | oui | deterministic diff | oui | summary | diff viewer/table |
| summarize source limitations | oui | structured metadata | oui | sourced summary | source matrix |
| decide success/rollback | later governed assessment | no here | no | prohibited | CAP-GOV-029/030 |

## 14. États fonctionnels
`not-defined`, `draft`, `criteria-incomplete`, `source-incomplete`, `owner-required`, `window-required`, `ready`, `active`, `data-pending`, `source-unavailable`, `insufficient-data-anticipated`, `review-required`, `stale`, `superseded`, `cancelled`.

## 15. États d’interface
Loading preserves criteria/sources ; Empty means no plan yet ; Partial names unavailable sources/criteria ; Error preserves last valid version ; Offline prevents fresh external observations unless guaranteed ; Permission denied masks restricted source data ; Stale exposes changed Run/source versions and requires review.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Verification Plan | versioned verification contract | CAP-GOV-029/032 | expected outcome, criteria, sources, owner, window explicit |
| criterion/source matrix | plan component | Verification Reviewer | source/freshness/limitations visible |
| unavailable-data condition | plan risk | CAP-GOV-029 | no success assumed from missing data |
| residual-risk questions | review input | CAP-GOV-029/032 | unresolved questions attributed |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Decision/Plan/Run | verification requirement available | CAP-GOV-028 | expected outcome/targets/conditions | Run |
| CAP-GOV-024/026/027 | checkpoint/terminal/partial state | CAP-GOV-028 | step/target/runtime/error context | Runs & Rollback |
| CAP-GOV-028 | Plan ready + observations due | CAP-GOV-029 | exact Plan/version/criteria/sources/window | Verification Plan |
| source unavailable | observation impossible | CAP-GOV-029 | missing source/impact/limitations | same Plan |

## 18. Dépendances
CAP-GOV-015..019/022/024/026/027/029/030/032, Command Incident projections, Investigate context, Endpoint/Settings/technical sources, Shared Trace/Linking/Reporting, Security access, OPEN-008/013/015.

## 19. Source de vérité
Govern owns the expected-outcome and verification-plan semantics. Each source owner remains authoritative for its observations. The Verification Plan is a plan, not proof; it never reclassifies Evidence/Finding or overwrites source state.

## 20. Provenance et audit
Record Decision/Run/Plan versions, expected outcome source, criterion ids/versions, success/partial/failure logic description, verification source owner/ref/version/freshness/restrictions, owner/time window, unavailable-data policy, residual-risk questions, Plan diffs, author/reviewer, automation/AI provenance and timestamps.

## 21. Permissions fonctionnelles
Verification Plan create/update/read; verification source metadata read; authorized no-effect verification query request; restricted context read; automated recommendation request; provenance export. Source data access remains independently permissioned.

## 22. Limites et erreurs
Unverifiable objective, ambiguous criterion, circular source that only repeats executor status, unavailable telemetry, stale target state, inaccessible Finding/Evidence or changed scope must remain explicit and may make the assessment inconclusive. No universal confidence score is invented.

## 23. Métriques
Runs with Plan before terminal execution; missing verification sources; criteria changed after observation; independent-check coverage where required; unavailable-data cases; Plans that directly claim Result — target zero.

## 24. Classification de livraison
`defined` / `planned`; no verification engine, provider, query/API/protocol, scoring algorithm or implementation selected.

## 25. Critères d’acceptation
**Given** a Run’s technical executor reports success, **When** the Verification Plan is reviewed, **Then** executor status may be one source but the expected outcome/criteria remain independently defined and no Result is created yet.

**Given** a required verification source is unavailable, **When** verification becomes due, **Then** the Plan exposes the gap and downstream assessment may become `insufficient-data` rather than fabricating success.

**Given** no AI, **When** a Verification Plan is authored, **Then** structured criteria, source matrices, deterministic validation and human review provide full functionality.

## 26. Questions ouvertes
OPEN-008 covers actual source/runtime availability, OPEN-013 class-2 Plan mutations and OPEN-015 runtime provenance bridge. Final verification engines/quality thresholds remain future; no new OPEN.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-029..033, Command/Investigate source consumers, Studio/Endpoint/Settings technical sources, future Objects/Permissions/Screens/Technique, GOV-2 report and GOV-3 audit/metrics.
