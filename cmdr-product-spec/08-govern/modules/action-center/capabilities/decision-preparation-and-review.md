---
id: CAP-GOV-014
title: Decision Preparation and Review
product: govern
module: action-center
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-AI-002, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-014 — Decision Preparation and Review

## 1. Définition
Assembler, comparer et revoir un Decision Draft pour une Action Request/version à partir du contexte source, impact/risk/reversibility, Policy outcomes, conflicts/exceptions, authority, Approvals, alternatives, unresolved questions and recommendations, sans que le draft ou une recommandation devienne une Decision.

## 2. Problème utilisateur
Les éléments nécessaires à une décision viennent de plusieurs owners et versions. Sans préparation explicite, un Decision Maker peut voir une synthèse non sourcée, ignorer un dissent, une Approval expirée ou un Policy conflict et prendre un draft/AI recommendation pour une autorité finale.

## 3. Objectifs
- agréger les inputs GOV-1 avec références/versions ;
- exposer unresolved questions, contradictions, restrictions and known unknowns ;
- comparer alternatives, Decision options and conditions ;
- permettre annotate/challenge/request changes/return ;
- distinguer recommendation, Decision Draft and Decision ;
- déclarer `decision-ready` seulement lorsque les prerequisites documentaires sont satisfaites.

## 4. Non-objectifs
Ne pas finaliser la Decision, exécuter l’action, démarrer Response Run, modifier source Evidence/Finding, activer une exception, auto-approuver ou autoriser une IA à décider.

## 5. Propriétaire
Govern / Action Center owns Decision Draft and preparation review. Source objects remain at their owners. `CAP-GOV-015` owns the final Decision record.

## 6. Utilisateurs
Principal : Decision Reviewer. Secondaires : Decision Maker, Govern Reviewer, Policy/Risk/Authority Reviewers, authorized Approvers, source requester/contributors, Auditor.

## 7. Conditions d’entrée
Current Action Request/version, relevant context/scope/risk/completeness results, Policy Evaluations/conflict/exception dispositions, Authority Requirement, eligibility and required Approvals or explicit pending/unknown dispositions.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Action Request + requester | CAP-GOV-003 | decision subject | oui | exact current version | not decision-ready |
| target/scope/context | CAP-GOV-004 | authorized-scope candidate | oui | reviewed current | not decision-ready |
| impact/risk/reversibility | CAP-GOV-005 | consequence assessment | oui for effectful action | current | unknown/information-required |
| completeness/Evidence context | CAP-GOV-006 | input quality | oui | current | not decision-ready |
| Policy outcomes/conflicts/exceptions | CAP-GOV-007/008 | governance inputs | according to applicability | current | pending/unknown visible |
| authority/eligibility/Approvals | CAP-GOV-009..013 | authority inputs | according to action/policy | current/effective | pending/not decision-ready |
| alternatives/recommendations | requester/reviewers/automation | options, not authority | non | source/version visible | no option invented |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | exact version and lineage | read |
| Incident | Command | operational impact/urgency | restricted read/link |
| Case/Finding/Evidence | Investigate | source context/restrictions | read/link, no requalification |
| Policy Evaluation/Exception | Govern | outcomes/reasons/authority refs | read |
| Authority/Approval | Govern | validity/scope/conditions | read |
| Workflow/Human Gate/Automation Run/Tool Call | Studio | automation provenance only | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Decision Draft | create/update/version/compare/return/supersede | Govern local concept | draft ≠ Decision |
| Decision option/condition candidate | add/update/withdraw | Govern local concept | source/rationale required |
| Review Context | mark unresolved/decision-ready | Govern | readiness ≠ authorization itself |
| source objects/Approval/Policy | no owner mutation | respective owner | consume exact versions only |

## 11. Fonctionnalités
Build Decision Draft; assemble source facts/versions; show requester/target/scope; summarize impact/risk/reversibility; list Policy outcomes/conflicts/exceptions; validate authority/Approval freshness; show dissent/unresolved questions; compare options/conditions; annotate/challenge/request changes/return; diff draft versions; declare decision-ready; preserve provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect/compare inputs | reviewer | review package | 0 | authorized read | sourced review | non |
| run deterministic readiness check | Decision Reviewer | Decision Draft | 1 | prerequisites defined | missing/pending list | non |
| create/update Decision Draft | Decision Reviewer | Decision Draft | 2 | request/version current | versioned draft | OPEN-013 |
| annotate/challenge/request changes | reviewer | Draft/Review Context | 2 | rationale | attributed review change | OPEN-013 |
| declare decision-ready | authorized reviewer | Review Context | 2 | required inputs valid | ready state, no Decision yet | OPEN-013 |
| record final Decision | Decision Maker | Decision | 3 when authority-bearing | CAP-GOV-015 prerequisites | handled by CAP-GOV-015 | yes |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| aggregate review inputs | oui | structured refs | oui | sourced summary | tables/links |
| detect missing/stale prerequisites | oui | oui | oui | explanation | readiness checklist |
| compare options/conditions | oui | deterministic diff | oui | suggestion | comparison table |
| draft rationale/conditions | oui | templates | oui | draft only | Decision template |
| recommend option | oui | factors/rules may rank transparently | oui | proposal only | human comparison |
| finalize Decision | authorized human | validation only | no autonomous | prohibited | CAP-GOV-015 explicit action |

## 14. États fonctionnels
Decision Draft: `not-started`, `draft`, `incomplete`, `information-required`, `policy-unresolved`, `authority-unresolved`, `approval-pending`, `under-review`, `changes-requested`, `review-ready`, `decision-ready`, `returned`, `superseded`, `cancelled`.

`decision-ready` is not `approved` and creates no target effect.

## 15. États d’interface
Loading preserves request/draft version ; Empty indicates no draft yet ; Partial lists inaccessible/pending inputs ; Error preserves latest valid draft ; Offline permits safe read but no authoritative readiness transition ; Permission denied masks protected context ; Stale shows which input invalidated readiness.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Decision Draft | versioned preparation | CAP-GOV-015 | all key inputs/unknowns/sources attributable |
| decision-ready disposition | Review Context event | Decision Maker/CAP-GOV-015 | readiness only, not Decision |
| request-changes/info need | review event | CAP-GOV-003/source/reviewer | exact issue and return origin |
| option/condition comparison | review artifact | Decision Maker | recommendation ≠ Decision |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-004..013 | inputs sufficient/pending | CAP-GOV-014 | exact review records/versions/unknowns | owner capability |
| CAP-GOV-014 | more information needed | CAP-GOV-003/source product | questions/request version/restricted refs | same draft lineage |
| CAP-GOV-014 | changes requested | relevant GOV-1 review capability | issue + prior result/draft | Decision Preparation |
| CAP-GOV-014 | decision-ready | CAP-GOV-015 | Decision Draft, authority/Approval/Policy context | Action Center |

## 18. Dépendances
CAP-GOV-003..013/015/016, Command/Investigate source projections, Studio provenance, Shared Linking/Versioning/Trace/Collaboration/Reporting, Security permissions, OPEN-007/013/015.

## 19. Source de vérité
Govern owns Decision Draft/preparation. Final Decision is CAP-GOV-015. Source facts, Evidence/Finding, Incident and Studio Runs remain source-owned. A recommendation or policy-engine/AI output is never source of Decision authority.

## 20. Provenance et audit
Record request/draft versions, every input ref/version/status, reviewers, missing/stale flags, options/alternatives, recommendations and producer, conditions, dissent/challenges, changes requested, readiness checks, automation/model/tool provenance and timestamps.

## 21. Permissions fonctionnelles
Decision Draft create/update/read; Decision review; request information; compare/annotate/challenge; conditions update before finalization; automated recommendation request; restricted context; provenance export. Final Decision permission is exercised in CAP-GOV-015.

## 22. Limites et erreurs
Expired Approval, changed request scope, stale Policy Evaluation, unresolved conflict, missing authority, restricted Evidence, failed automation, concurrent draft edit or tenant change invalidates readiness as appropriate and cannot be hidden by a summary.

## 23. Métriques
Drafts by readiness blocker; number of review cycles/changes; stale-input invalidations; unresolved questions at first review; AI suggestions accepted/rejected with attribution; automatic Decisions — target zero.

## 24. Classification de livraison
`defined` / `planned`; no final Decision workflow engine, model provider, API, screen composition or object schema selected.

## 25. Critères d’acceptation
**Given** an Approval exists but no Decision has been recorded, **When** Decision Preparation opens, **Then** Approval is shown as an input and the request remains undecided until CAP-GOV-015 records a Decision.

**Given** an AI recommends `approve`, **When** the Decision Maker reviews the draft, **Then** the recommendation is attributed and source-backed, and no disposition changes automatically.

**Given** no AI provider, **When** a Decision is prepared, **Then** tables, diffs, checklists, templates and human review provide the complete function.

## 26. Questions ouvertes
OPEN-007, OPEN-013 and OPEN-015 remain open. Final Decision object/state/permission semantics remain future Objects/Permissions work; no new OPEN is created.

## 27. Consommateurs documentaires
Action Center, Decision Register, CAP-GOV-015/016, source-product return flows, future Objects/Permissions/Screens/Journeys/GOV-2/GOV-3/Technique and GOV-1 conformance report.
