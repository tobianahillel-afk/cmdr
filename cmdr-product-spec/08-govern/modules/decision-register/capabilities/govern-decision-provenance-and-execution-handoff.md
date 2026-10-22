---
id: CAP-GOV-016
title: Govern Decision Provenance and Execution Handoff
product: govern
module: decision-register
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-016 — Govern Decision Provenance and Execution Handoff

## 1. Définition
Retracer l’intégralité de la chaîne GOV-1 jusqu’à une Decision et préparer un **Execution Handoff Package** exact, versionné et sans effet contenant l’action autorisée, targets/scope, conditions, limits, expiry, Approval/Decision context, rollback/verification requirements and provenance for future GOV-2.

## 2. Problème utilisateur
Même une Decision valide peut devenir dangereuse si l’exécution reçoit une cible différente, perd une condition/expiry, ignore une Approval révoquée ou ne sait plus quel Finding/Incident/Policy/authority a conduit à l’autorisation. Un handoff flou peut aussi être pris pour un Response Run déjà lancé.

## 3. Objectifs
- conserver la chaîne source→request→reviews→Approval→Decision ;
- inclure exact approved action/scope/targets and prohibited scope ;
- inclure conditions/time bounds/expiry/start constraints ;
- inclure rollback/verification requirements sans les exécuter ;
- inclure Tool/Tool Call/Automation Run/Human Gate provenance sans transfert d’ownership ;
- produire/versionner/withdraw/supersede un package ;
- empêcher tout handoff d’une Decision rejetée/deferred/expired/stale.

## 4. Non-objectifs
Ne pas créer Response Run/Result, sélectionner une execution primitive, appeler Endpoint/Cloud/identity/network command, effectuer rollback, vérifier un résultat, supprimer trace, finaliser l’API/protocole du package ou commencer GOV-2/GOV-3.

## 5. Propriétaire
Govern / Decision Register owns Decision provenance and the GOV-1 Execution Handoff Package. Studio owns Tool/Tool Call/Automation Run/Human Gate; future GOV-2 owns Response Run/execution/verification/rollback; source products own their evidence/context.

## 6. Utilisateurs
Principal : Govern Reviewer / Decision Maker. Secondaires : future Response Operator, Incident Commander/Investigation Lead as read-limited source/consumer, Studio Operator for Run provenance, Auditor.

## 7. Conditions d’entrée
Current non-expired Decision with disposition permitting future action (`approve` or `approve-with-conditions`), exact request/version, target/scope, authority/Approval context, conditions and required provenance accessible. Any material stale/missing authority blocks a valid handoff.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Decision/version/disposition | CAP-GOV-015 | authoritative source | oui | current and unexpired | no handoff |
| Action Request/version | CAP-GOV-003 | requested action lineage | oui | exact decided version | no handoff |
| exact targets/scope/limits | CAP-GOV-004/Decision | execution bounds | oui | as decided | no target inference |
| risk/reversibility | CAP-GOV-005 | rollback/recovery context | yes for effectful action | decided review version | package incomplete/block |
| Policy/conflict/exception | CAP-GOV-007/008 | governance provenance | according to decision | current at decision | package marks missing/stale |
| authority/Approvals/emergency | CAP-GOV-009..013 | authority provenance | according to Decision | current/decision snapshot | block if required evidence absent |
| Tool/Tool Call/Automation Run/Human Gate refs | Studio | preparation provenance | when used | exact run/call refs | package records absent/not-used, never fabricates |
| rollback/verification requirements | Decision/GOV-1 review | future execution constraints | according to decision | exact Decision | explicit none/unknown/required |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision | Govern | disposition/scope/conditions/expiry/authority | read |
| Action Request + GOV-1 reviews | Govern | lineage/reasons/versions | read |
| Incident/Case/Finding/Evidence | Command/Investigate | source refs/return origin/restrictions | read/link only |
| Tool/Tool Call/Workflow/Human Gate/Automation Run | Studio | provenance refs/results used in governance | read/link |
| future Response Run/Result | Govern future GOV-2 | destination/reference type only | no creation/mutation in GOV-1 |
| Trace/Activity/Report refs | Shared | provenance/export presentation | consume |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Execution Handoff Package | create/version/withdraw/supersede | Govern local concept | package ≠ Response Run; no target effect |
| provenance links | add/version | source owner + Shared linking | stable refs, permissions preserved |
| Decision handoff state/projection | mark prepared/stale/withdrawn as local projection | Govern | Decision itself remains authoritative history |
| Response Run/Result/target | no creation or mutation | future GOV-2/runtime owner | strict GOV-1 stop line |

## 11. Fonctionnalités
Assemble provenance graph; verify Decision current/approved; copy by reference exact action/targets/scope/conditions/expiry; preserve prohibited scope; include authority/Approvals/Policy/exception/emergency; include source-product refs; include Studio automation/tool provenance; include rollback/verification requirements; compare package to Decision; version/supersede/withdraw; prepare future GOV-2 destination; export provenance under permission.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect provenance | reviewer/auditor | provenance chain | 0 | read | traceable source chain | non |
| run deterministic Decision↔package consistency check | reviewer | package | 1 | Decision current | mismatch list | non |
| prepare/update package | authorized reviewer | Execution Handoff Package | 2 | approved Decision + exact bounds | versioned no-effect package | OPEN-013 |
| withdraw/supersede stale package | reviewer/Decision Maker | package | 2 | reason + lineage | old package retained | OPEN-013 |
| start Response Run/execute/rollback | none in GOV-1 | future run/target | — | prohibited | no effect | future GOV-2 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate provenance refs | oui | structured linking | oui | sourced summary | trace/link table |
| compare package to Decision | oui | exact diff | oui | explanation | deterministic diff |
| flag stale/expired authority/Decision | oui | timestamp/version checks | oui | explanation | validity checklist |
| draft handoff summary | oui | template | oui | draft | structured template |
| start execution/rollback | future human/runtime | no in GOV-1 | no | prohibited | GOV-2 only |

## 14. États fonctionnels
Execution Handoff Package: `not-prepared`, `draft`, `incomplete`, `consistency-review`, `ready-for-future-execution`, `stale`, `blocked`, `expired-with-decision`, `withdrawn`, `superseded`. `ready-for-future-execution` means documentary package readiness only and is **not** a Response Run or execution state.

## 15. États d’interface
Loading preserves Decision/package versions ; Empty no package ; Partial names missing/restricted provenance ; Error preserves prior package and blocks release ; Offline can inspect but not declare current handoff readiness ; Permission denied masks restricted refs ; Stale/expired visibly blocks future consumption.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Execution Handoff Package | no-effect governance package | future GOV-2 | exact approved action/scope/targets/conditions/expiry/authority requirements |
| provenance chain | Trace/link set | Audit/source products/future GOV-3 | source ownership/permissions/version preserved |
| package consistency result | review result | Decision Register/future GOV-2 | Decision mismatch/staleness explicit |
| withdrawal/supersession event | lifecycle event | future GOV-2/history | old package remains resolvable and unusable as current authority |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-015 | approve/approve-with-conditions current Decision | CAP-GOV-016 | Decision/request/reviews/authority/conditions | Decision Register |
| CAP-GOV-016 | package prepared and consistent | future GOV-2 | exact package; no Response Run created | Decision Register/GOV-1 |
| Decision expires/supersedes | validity change | CAP-GOV-016 | current vs previous Decision/package refs | package becomes stale/blocked |
| source/authority material change requiring new Decision | invalidation | CAP-GOV-014/015 | changed source/why handoff invalid | Decision preparation |

## 18. Dépendances
CAP-GOV-003..015, Shared Linking/Trace/Activity/Reporting/Export/Versioning, Studio Tool/Tool Call/Workflow/Human Gate/Automation Run, source objects, Security permissions, future GOV-2/GOV-3, OPEN-007/013/015.

## 19. Source de vérité
Govern owns Decision and Execution Handoff Package. Source-product objects remain authoritative for original facts; Studio for automation/tool provenance; future GOV-2 for Response Run and Result. Package content references authority; it is not runtime truth.

## 20. Provenance et audit
Record source product, Incident/Case/Finding/Evidence refs, Action Request/versions, reviewer assessments, target/scope, risk, Policy Evaluations/conflicts/exceptions, authority/eligibility, Approval Requests/Approvals, delegation/emergency, Decision Draft/Decision/conditions/expiry, Tool/Tool Calls/Automation Runs/Human Gates, errors/info requests, human dispositions, package versions/consistency checks/withdrawal/supersession and correlation ids.

## 21. Permissions fonctionnelles
Decision/Approval/provenance read, Execution Handoff prepare/withdraw/supersede, restricted context read, provenance export, automated recommendation request. No Response Run create/start, target mutation or rollback permission is added by GOV-1.

## 22. Limites et erreurs
Rejected/deferred/info-required/cancelled/expired/superseded Decision, stale Approval/authority, mismatch between package and Decision, inaccessible target ref, missing rollback requirement, restricted provenance or destination unavailable keeps package blocked/stale and performs no execution.

## 23. Métriques
Packages by state; consistency mismatches; stale/expired before future use; source/provenance gaps; time Decision→package readiness; packages incorrectly creating Response Runs/target effects — target zero.

## 24. Classification de livraison
`defined` / `planned`; package is a functional contract only. No API, serialization, protocol, queue, execution engine or rollback implementation is selected.

## 25. Critères d’acceptation
**Given** an `approve-with-conditions` Decision, **When** the handoff is prepared, **Then** every condition, prohibited scope, expiry and rollback/verification requirement is present and a missing item blocks readiness.

**Given** a Decision expires before future execution, **When** the package is inspected, **Then** it becomes expired/stale and cannot create or authorize a Response Run.

**Given** a Studio Human Gate and Automation Run contributed to preparation, **When** provenance is assembled, **Then** both remain Studio-owned references; neither is converted into Approval, Decision or Response Run.

## 26. Questions ouvertes
OPEN-007 remains open for Human Gate/Govern relation; OPEN-013 for class-2 preparation governance; OPEN-015 for final Automation Run/Response Run bridge. GOV-1 closes none.

## 27. Consommateurs documentaires
Decision Register, future GOV-2 Playbooks/Runs/Verification/Rollback, future GOV-3 Audit/Metrics, source-product return flows, Studio/Shared provenance, Objects/Permissions/Screens/Journeys/Technique and GOV-1 conformance report.
