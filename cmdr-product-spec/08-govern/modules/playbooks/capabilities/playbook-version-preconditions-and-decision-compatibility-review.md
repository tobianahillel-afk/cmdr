---
id: CAP-GOV-018
title: Playbook Version, Preconditions and Decision Compatibility Review
product: govern
module: playbooks
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-018 — Playbook Version, Preconditions and Decision Compatibility Review

## 1. Définition
Vérifier qu’un Playbook candidat et sa **version exacte** sont compatibles avec la Decision, l’Execution Handoff Package, les target types, le scope, les conditions, l’environnement, les dépendances et les exigences de rollback/verification avant toute préparation d’Execution Plan.

## 2. Problème utilisateur
Un Playbook peut rester catalogué alors que sa version, son Workflow sous-jacent, ses dépendances ou ses préconditions ont changé. Utiliser silencieusement une version plus récente que celle réellement revue peut invalider scope, conditions ou autorisation sans que l’utilisateur le voie.

## 3. Objectifs
- comparer exact Playbook version à exact Decision/handoff version ;
- vérifier action, target type, allowed/prohibited scope, conditions, expiry et required approvals ;
- examiner tenant/environment, executor owner, runtime dependencies et limitations ;
- vérifier rollback et verification capabilities déclarées ;
- détecter deprecated/incompatible/missing-runtime states ;
- exiger re-review lorsqu’une version change après Decision ou handoff.

## 4. Non-objectifs
Ne pas substituer une version, modifier Decision/Approval, finaliser une policy d’autorisation, résoudre un target courant, exécuter un precondition check avec effet, créer un Run, appeler Workflow/Endpoint ou choisir un runtime/provider global.

## 5. Propriétaire
Govern / Playbooks owns Decision-to-Playbook compatibility review. Govern owns Playbook semantics; Studio owns Workflow/version; Settings/Endpoint/runtime owners own dependency and execution-capability truth.

## 6. Utilisateurs
Principal : Govern Reviewer / Response Operator. Secondaires : Decision Maker, Playbook Owner, Studio Operator, Runtime/Endpoint Operator, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Candidate from CAP-GOV-017, exact Decision and Execution Handoff Package, readable Playbook/version metadata, decision conditions/expiry and authorized dependency projections.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| selected Playbook/version | CAP-GOV-017 / Govern | compatibility subject | oui | exact selected version | review impossible |
| Decision/version/disposition | CAP-GOV-015 | authority constraints | oui | current/unexpired for planning | blocked/expired |
| Handoff target/scope/conditions | CAP-GOV-016 | exact execution bounds | oui | current package | blocked |
| required approvals/exceptions | GOV-1 records | authority prerequisites | selon Decision | validity visible | compatibility incomplete |
| runtime/executor dependency metadata | Settings/Endpoint/Studio | capability/health projection | selon Playbook | latest known timestamp | missing-runtime/unknown |
| rollback/verification support | Playbook/version + executor metadata | safety compatibility | selon Decision | exact versions | incompatible/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision | Govern | action/scope/conditions/expiry | read |
| Execution Handoff Package | Govern | targets/prohibited scope/requirements | read |
| Playbook | Govern | exact version/preconditions/limits | read/compare |
| Workflow / Version | CMDR Studio | referenced version/availability | read/link |
| Integration/Environment/Endpoint capability | Settings/Endpoint | dependency/readiness metadata | restricted read |
| Approval/Exception context | Govern | required validity/bounds | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Playbook Compatibility Review | create/version/recheck/supersede | Govern local concept | exact versions and reasons required |
| Selection Context | mark compatible/incompatible/re-review | Govern | no silent version replacement |
| Playbook/Decision/Workflow | no source mutation | respective owner | read/compare only |
| dependency objects | no mutation | Settings/Endpoint/Studio | projection only |

## 11. Fonctionnalités
Diff selected version against Decision requirements; validate intended action/target type; check included/prohibited scope; inspect Decision conditions/start window/expiry; check required Approval/Exception context as input; compare tenant/environment; inspect executor support/dependencies; check rollback/verification support; expose limitations/deprecation; re-run after version/dependency change; mark re-decision/review requirement where authorization impact is possible.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect versions/preconditions | reviewer | Playbook/Decision | 0 | read | exact version view | non |
| run compatibility check | reviewer | Compatibility Review | 1 | exact inputs | explainable compatibility result | non |
| annotate/request alternative | reviewer | Review | 2 | mismatch/rationale | versioned review action | OPEN-013 |
| accept compatible candidate for planning | authorized reviewer | Selection Context | 2 | all required checks pass/known limits accepted | planning candidate | OPEN-013 |
| silently substitute version/start | none | Playbook/Run | — | forbidden | no effect | later capabilities only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| diff Playbook versions | oui | oui | oui | sourced summary | deterministic diff |
| check Decision/action/scope compatibility | oui | explicit rules | oui | explanation | compatibility matrix |
| identify dependency mismatch | oui | metadata comparison | oui | suggestion | dependency checklist |
| suggest alternative version | oui | catalog matching | oui | candidate only | CAP-GOV-017 search |
| declare execution authorized/start | governed later | no autonomous | no | prohibited | CAP-GOV-021..023 |

## 14. États fonctionnels
`not-reviewed`, `reviewing`, `compatible`, `compatible-with-limitations`, `incompatible-action`, `incompatible-target`, `scope-mismatch`, `condition-mismatch`, `deprecated-version`, `missing-runtime`, `rollback-mismatch`, `verification-mismatch`, `expired-decision`, `re-review-required`, `re-decision-required`, `superseded`.

## 15. États d’interface
Loading preserves exact versions ; Empty means no selected candidate ; Partial lists unavailable dependency facts ; Error preserves last valid review ; Offline cannot finalize a fresh compatibility result ; Permission denied masks protected dependency details ; Stale forces rerun after any material version change.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Compatibility Review | versioned review result | CAP-GOV-019/021 | exact Decision/Playbook/dependency versions + reasons |
| compatible candidate | planning context | CAP-GOV-019 | compatibility ≠ authorization/readiness |
| incompatibility blocker | review event | CAP-GOV-017/source owner | mismatch dimension explicit |
| re-review/re-decision requirement | governance condition | GOV-1 Decision path/CAP-GOV-021 | no silent authorization carry-over |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-017 | candidate selected | CAP-GOV-018 | exact Playbook/version + Decision/handoff | Playbooks |
| CAP-GOV-018 | compatible | CAP-GOV-019 | Playbook/version/preconditions/limitations | compatibility review |
| CAP-GOV-018 | incompatible | CAP-GOV-017 | mismatch + candidate/version | same handoff |
| Playbook/Workflow/dependency changes | material drift | CAP-GOV-018 | old/new refs and affected checks | same review lineage |
| authorization implication detected | re-decision required | CAP-GOV-014/015 | change and scope/condition impact | Decision Register |

## 18. Dépendances
CAP-GOV-015..017/019/020/021, Playbook, Studio Workflow/Version, Settings integration/environment/health, Endpoint capabilities, Security/Approval/Exception validity, Shared Versioning/Trace/Linking, OPEN-008/013/015.

## 19. Source de vérité
Decision is authoritative for approved action/scope/conditions; Playbook/version is authoritative for response procedure semantics; Studio/Settings/Endpoint own referenced runtime facts. Govern owns the compatibility disposition for the exact versions reviewed.

## 20. Provenance et audit
Record Decision/handoff ids and versions, selected Playbook/version, Workflow/tool/dependency refs and versions, precondition list, check outcomes/reasons, limitations/deprecation, rollback/verification support, reviewer, version changes, automation/AI provenance, re-review/re-decision disposition and timestamps.

## 21. Permissions fonctionnelles
Playbook/version read/compare, compatibility run/read, restricted dependency metadata, request alternative, automated recommendation request, provenance export. No source object management or execution permission is implied.

## 22. Limites et erreurs
Unknown runtime, unavailable metadata, changed Workflow version, stale approval/exception, expired Decision, unsupported target, rollback/verification mismatch or permission denial yields explicit incompatible/unknown/re-review state; never silent pass or substitution.

## 23. Métriques
Compatibility outcomes; version-change rechecks; deprecated/missing-runtime blocks; rollback/verification mismatch rates; re-decision requests caused by version drift; silent substitutions — target zero.

## 24. Classification de livraison
`defined` / `planned`; functional compatibility semantics only, no engine/API/provider/runtime selected.

## 25. Critères d’acceptation
**Given** the selected Playbook changes from v3 to v4 after the Decision, **When** planning resumes, **Then** the exact diff is visible, compatibility is rechecked and v4 cannot silently inherit authorization.

**Given** the Decision requires rollback support but the selected Playbook version declares rollback unsupported, **When** compatibility is evaluated, **Then** the candidate is incompatible or requires explicit re-decision; planning cannot hide the mismatch.

**Given** no AI, **When** compatibility is reviewed, **Then** deterministic version diffs, matrices, checklists and human review provide the full function.

## 26. Questions ouvertes
OPEN-008 remains open for actual runtime/platform support, OPEN-013 for class-2 review mutations, OPEN-015 for Workflow/Response Run bridge. Final compatibility engine and authorization impact policy remain future.

## 27. Consommateurs documentaires
Playbooks, CAP-GOV-019/020/021, Runs & Rollback, GOV-1 Decision review, Studio/Endpoint/Settings boundaries, future Objects/Permissions/Screens/Technique and GOV-2 conformance report.
