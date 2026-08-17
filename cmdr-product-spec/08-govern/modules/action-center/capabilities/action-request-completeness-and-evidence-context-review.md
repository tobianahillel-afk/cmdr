---
id: CAP-GOV-006
title: Action Request Completeness and Evidence Context Review
product: govern
module: action-center
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-014, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001]
open_decisions: [OPEN-013, OPEN-014]
source-of-truth: canonical
---
# CAP-GOV-006 — Action Request Completeness and Evidence Context Review

## 1. Définition
Vérifier qu’une Action Request contient les informations fonctionnelles nécessaires à la review Govern et que ses références Evidence/Finding/Incident/Case sont disponibles ou explicitement limitées, sans requalifier Evidence/Finding ni substituer une absence par une conclusion.

## 2. Problème utilisateur
Une demande peut être formellement remplie mais manquer de justification, source Finding, Evidence accessible, rollback context ou time bounds. Une checklist vague peut faire passer une demande incomplète en autorisée.

## 3. Objectifs
Vérifier requester, action, target, justification, source Finding/Evidence, Incident/Case context, expected outcome, impact/risk, rollback when required, time bounds, authority/Approval context and contradictions; generate explicit missing-information reasons; preserve restricted evidence semantics.

## 4. Non-objectifs
Ne pas confirmer un Finding, qualifier Evidence, produire une conclusion d’investigation, exiger Evidence quand une emergency basis autorisée est prévue, approuver ou décider.

## 5. Propriétaire
Govern / Action Center owns completeness review. Investigate remains owner of Evidence/Finding and their qualification. Command remains owner of Incident. Source permissions remain in force.

## 6. Utilisateurs
Principal : Govern Reviewer. Secondaires : Decision Reviewer, Policy/Authority Reviewer, source requester/contributor, Evidence Reviewer as source owner, Auditor.

## 7. Conditions d’entrée
Action Request/version, Context Review and permission-aware source links. Requiredness depends on action/policy/context and must be explainable rather than a single universal form.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| requester/action/target/justification | Action Request | core request | oui | current request | incomplete |
| Incident/Case/Finding context | source products | reason/context | according to source path | source version | information-required or justified alternate basis |
| Evidence references | Investigate | support/provenance | where assertions require evidence | qualification/version visible | missing/restricted, never fabricated |
| expected outcome + impact/risk | requester + CAP-GOV-005 | decision input | yes for effectful action | current assessment | incomplete |
| rollback/recovery context | requester/CAP-GOV-005 | reversibility | policy/action dependent | current | unknown/incomplete |
| time bounds/expiry | request/policy/authority | temporal constraint | where authority/effect bounded | current | information-required/unknown |
| authority/approval context | CAP-GOV-009..011 | future review prerequisite | not necessarily at early completeness | current | mark pending, not complete authorization |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | all review fields/version | review |
| Context/Risk assessments | Govern | reviewed scope/impact/reversibility | read |
| Finding/Evidence/Case | Investigate | refs, status, restrictions, provenance | read restricted; no requalification |
| Incident | Command | context/impact/urgency | read/link |
| Policy/Authority/Approval projections | Govern/Security | requiredness and pending state | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Completeness Review | create/version | Govern local concept | exact request version + checklist reasons |
| Information Request | create/close | Govern local concept | question names missing/contradictory input |
| Action Request processing state | mark incomplete/ready for next review | Govern | completeness ≠ authorized |
| Evidence/Finding | no mutation | Investigate | never requalified by Govern |

## 11. Fonctionnalités
Dynamic required-field checks by action/policy; source-link accessibility review; contradiction listing; restricted-source handling; missing-information request generation; alternate emergency basis capture where allowed; version diff; readiness disposition for subsequent policy/authority/Decision review.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect completeness/source refs | reviewer | request/references | 0 | read | explicit checklist | non |
| run deterministic completeness check | reviewer | Completeness Review | 1 | rule set known | missing/complete dimensions | non |
| request more information | reviewer | Information Request | 2 | missing item named | source follow-up | OPEN-013 |
| accept justified restricted/alternate context for review | reviewer | Completeness Review | 2 | rationale/policy permits | review disposition, not Decision | OPEN-013 |
| requalify Evidence/Finding | none | Evidence/Finding | — | forbidden | no mutation | source owner only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| field requiredness check | oui | oui | oui | explain missing | checklist/rules |
| source-link accessibility check | oui | yes | yes | no need | link status table |
| identify contradictions | oui | deterministic comparisons | oui | suggestion/explanation | diff/matrix |
| draft information questions | oui | templates | oui | draft | manual template |
| qualify Evidence/Finding | source-owner human workflow | no | no | prohibited | Investigate workflow |

## 14. États fonctionnels
`not-reviewed`, `incomplete`, `information-required`, `information-received`, `restricted-context`, `contradicted`, `complete-for-review`, `complete-with-known-unknowns`, `superseded`.

`complete-for-review` does not mean `authorized`, `approved` or `decision-ready`.

## 15. États d’interface
Loading retains checklist progress ; Empty requires a request ; Partial names inaccessible/missing fields ; Error preserves prior checks ; Offline disables authoritative completeness transition ; Permission denied shows safe “restricted” state ; Stale source versions require recheck.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Completeness Review | review result | CAP-GOV-007/009/014 | required/available/missing/restricted/contradictory dimensions visible |
| Information Request | event/context | source product | exact questions, request/version, return origin |
| complete-for-review disposition | processing state | Response Inbox/Action Center | no authority or Approval implied |
| contradiction set | review context | risk/policy/decision reviewers | sources/version preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Action Request | completeness review | CAP-GOV-006 | request/version/context/risk/source refs | Action Center |
| CAP-GOV-006 | missing info | source product | questions/restricted refs/return origin | same request lineage |
| source product | response submitted | CAP-GOV-003/006 | linked response/new version | source product |
| CAP-GOV-006 | complete for review | CAP-GOV-007/009/014 | completeness disposition/known unknowns | Action Center |

## 18. Dépendances
CAP-GOV-003/004/005/007/009/014, Investigate Evidence/Finding/Case, Command Incident, Shared Linking/Trace, OPEN-013, OPEN-014 for final Artifact/Attachment/material relations.

## 19. Source de vérité
Govern owns completeness disposition. Evidence/Finding qualification remains Investigate; Incident remains Command. A restricted/inaccessible source is represented as restricted/unknown, never copied to bypass permission.

## 20. Provenance et audit
Record checklist version, requiredness rules/sources, request version, each source ref/version/access status, contradictions, reviewer, questions/responses, restricted handling, automation provenance and final completeness disposition.

## 21. Permissions fonctionnelles
Action Request review; restricted context read; request information; context update; source reference read; automated recommendation request; provenance export. No Evidence/Finding mutation permission is introduced.

## 22. Limites et erreurs
Missing Evidence, disputed Finding, unavailable Incident, inaccessible Case, stale source, contradictory requester data, absent rollback context or permission denial must remain explicit. Completeness does not certify truth, safety or authority.

## 23. Métriques
Information requests per request/category; requests complete on first review; restricted/contradicted source count; stale-source rechecks; zero Evidence/Finding requalifications or auto-Decisions.

## 24. Classification de livraison
`defined` / `planned`; no final validation engine, object schema or required-field standard is selected.

## 25. Critères d’acceptation
**Given** a request referencing Evidence the reviewer cannot access, **When** completeness is checked, **Then** the Evidence remains `restricted`, permission is not inherited and the request cannot claim that Evidence was reviewed.

**Given** a request missing required rollback context, **When** the check completes, **Then** the missing item is explicit and `complete-for-review` is not granted unless an applicable documented exception basis permits it.

**Given** no AI, **When** completeness is reviewed, **Then** deterministic checklists, source-link status and human review suffice.

## 26. Questions ouvertes
OPEN-013 remains open; OPEN-014 remains open for final Artifact/Attachment/material relations and retention. No new decision is created.

## 27. Consommateurs documentaires
Action Center, Policy/Authority/Decision capabilities, source-product Action Request preparation, future Objects/Permissions/Screens/Journeys/GOV-2 and quality gates.
