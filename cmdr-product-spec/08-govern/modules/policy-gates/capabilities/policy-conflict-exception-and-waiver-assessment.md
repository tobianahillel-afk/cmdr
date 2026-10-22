---
id: CAP-GOV-008
title: Policy Conflict, Exception and Waiver Assessment
product: govern
module: policy-gates
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-GOV-008 — Policy Conflict, Exception and Waiver Assessment

## 1. Définition
Analyser les conflits de Policies et préparer/revoir des **Exception Candidates** ou waivers fonctionnels limités par action, target, scope et durée, sans supprimer une Policy, activer silencieusement une exception ou prendre automatiquement une Decision.

## 2. Problème utilisateur
Deux Policies ou versions peuvent imposer des exigences incompatibles, ou une request peut nécessiter une dérogation temporaire. Sans workflow explicite, une équipe peut ignorer un block, considérer une exception proposée comme active ou effacer la Policy conflictuelle.

## 3. Objectifs
- conserver Policies/versions/outcomes en conflit ;
- expliquer conflict reason, scope and priority context if source provides it ;
- créer une Exception Candidate avec justification, target/scope/duration/expiry ;
- documenter compensating-control candidates et approver requirement ;
- supporter review, withdrawal, supersession, expiration/revocation context ;
- transmettre le résultat à authority/Decision review sans auto-bypass.

## 4. Non-objectifs
Ne pas choisir un modèle universel de priorité, supprimer/modifier une Policy, activer une exception sans autorité, masquer un block, exécuter l’action ou considérer urgence comme waiver automatique.

## 5. Propriétaire
Govern / Policy Gates owns conflict and exception-governance review. Security/Policy sources retain rule/authority semantics; Approvals & Authorities owns required authority/approver assessment.

## 6. Utilisateurs
Principal : Policy Reviewer. Secondaires : Authority Reviewer, Decision Maker, Security Reviewer, request owner, Policy Owner, Auditor.

## 7. Conditions d’entrée
At least one Policy Evaluation with conflict/block/warn/unknown or an explicit exception request; exact request/version, Policy versions, target/scope and reviewer permission.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Policy Evaluations | CAP-GOV-007 | outcomes/reasons/versions | oui | current request version | no conflict disposition |
| conflict facts | Policy versions/scopes | incompatibility context | yes for conflict | versioned | conflict stays unresolved |
| exception justification | requester/reviewer | rationale | yes for candidate | request version | candidate incomplete |
| target/scope/time bounds | CAP-GOV-004/request | bounded exception | yes | reviewed current | candidate blocked |
| impact/risk/reversibility | CAP-GOV-005 | consequence context | yes for material exception | current assessment | unknown/incomplete |
| compensating-control candidates | source/Policy/Security | mitigation proposal | according to policy | source/version visible | none claimed |
| authority/approver requirement | CAP-GOV-009/010 | activation authority | before active exception | current | candidate only |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Policy / Policy Evaluation | Govern | versions/outcomes/scope/reasons | read/review |
| Action Request | Govern | action/target/scope/version | read |
| Context/Risk Assessments | Govern | impact/reversibility | read |
| Principal/Role/Authority projection | Settings/Security/Govern | requester/approver context | minimal read |
| Decision/Approval | Govern | resulting authority records | read/link when present |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Policy Conflict | create/review/version/resolve-by-disposition | Govern local concept | does not delete either Policy |
| Policy Exception Candidate | propose/update/withdraw/supersede/expire/revoke-context | Govern | candidate ≠ active exception |
| compensating-control candidate | attach/review | source/Govern context | candidate not proof of control |
| Policy | no mutation by conflict resolution in GOV-1 | Govern future admin source | source version retained |

## 11. Fonctionnalités
Compare Policy versions/scopes/outcomes; classify conflict reason; preserve priority if provided without inventing it; propose Exception Candidate; bound target/scope/duration/expiry; attach justification and compensating controls; identify required authority; review/withdraw/supersede; record active-exception authorization only through later Approval/Decision context; preserve provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect/compare conflicting Policies | reviewer | Policy/Evaluations | 0 | read | sourced conflict view | non |
| run bounded conflict comparison | Policy Reviewer | Policy Conflict | 1 | versions available | explainable differences | non |
| create/update Exception Candidate | Policy/Security Reviewer | candidate | 2 | justification/scope/time bound | versioned candidate | OPEN-013 |
| withdraw/supersede candidate | authorized reviewer | candidate | 2 | no hidden history removal | disposition retained | OPEN-013 |
| authorize high-impact active exception | authorized approver/Decision Maker | exception authority | 3 | required authority/SoD/Approval | authority record only | required |
| bypass Policy/execution | none | target/policy | 4/3 | forbidden | no effect | future governed path only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| compare Policies/versions | oui | diff/rules | oui | summary | deterministic diff |
| identify conflict candidate | oui | condition comparison | oui | suggestion | conflict matrix |
| draft Exception Candidate | oui | template | oui | draft | structured form |
| suggest compensating controls | oui | catalog if available | oui | candidates only | control catalog/human review |
| activate exception | authorized human later | validation only | no autonomous | prohibited | Approval/Decision path |

## 14. États fonctionnels
Policy Conflict: `detected`, `under-review`, `unresolved`, `resolved-by-policy-context`, `exception-considered`, `superseded`.
Exception Candidate: `draft`, `incomplete`, `review-ready`, `authority-pending`, `approval-pending`, `approved-context`, `rejected`, `withdrawn`, `expired`, `revoked`, `superseded`. `approved-context` still does not execute the target action.

## 15. États d’interface
Loading keeps both Policy versions ; Empty means no conflict/candidate ; Partial names missing policy/authority inputs ; Error preserves prior review ; Offline read-only ; Permission denied masks protected Policy/rationale ; Stale requires reevaluation after Policy/request change.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Policy Conflict review | review record | CAP-GOV-009/014 | both Policy versions/outcomes/reasons preserved |
| Exception Candidate | versioned candidate | CAP-GOV-009..015 | candidate ≠ active exception |
| compensating-control context | candidate set | Authority/Decision review | candidate/source/uncertainty visible |
| withdrawal/revocation/supersession context | lifecycle event | Decision/history | previous authorization/history retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-007 | conflict detected | CAP-GOV-008 | Policies/versions/outcomes/reasons | Policy Gates |
| CAP-GOV-008 | exception requires authority | CAP-GOV-009/010 | candidate, target/scope/duration/risk | Policy Gates |
| CAP-GOV-008 | Approval required | CAP-GOV-011 | requirement/candidate/authority context | candidate review |
| CAP-GOV-008 | review complete | CAP-GOV-014 | conflict/exception disposition + unresolved issues | Action Center |
| request/Policy version changes | reevaluate | CAP-GOV-007/008 | previous records + new versions | Policy Gates |

## 18. Dépendances
CAP-GOV-004/005/007/009..015, Policy object, Security SoD/authority/step-up, Shared Versioning/Trace/Linking/Notifications, OPEN-013.

## 19. Source de vérité
Govern owns conflict and exception-candidate review records. Policy source remains authoritative for its own version. Only explicit authority/Approval/Decision records can make an exception effective for a governed action; candidate status never does.

## 20. Provenance et audit
Record Policies/versions/outcomes, conflict reason, request/version, target/scope, justification, duration/expiry, compensating controls and sources, requester/reviewer, authority/Approval refs, candidate lifecycle, revocation/supersession, automation provenance and timestamps.

## 21. Permissions fonctionnelles
Policy read/Evaluation read; Exception propose/review/revoke; authority/Approval context read; restricted context; automated recommendation request; provenance export. High-impact activation may be class 3 and step-up/SoD governed.

## 22. Limites et erreurs
Unknown Policy precedence, stale evaluation, missing Policy owner, unavailable authority, contradictory scope, expired candidate, source restriction or automation failure leaves conflict unresolved/candidate inactive; no silent bypass or auto-rejection.

## 23. Métriques
Conflicts by reason; unresolved duration; Exception Candidates by disposition/duration; expired/revoked candidates; compensating-control coverage; zero silent Policy bypasses/automatic active exceptions.

## 24. Classification de livraison
`defined` / `planned`; no final precedence engine, exception object schema, Policy authoring system or enforcement runtime chosen.

## 25. Critères d’acceptation
**Given** two applicable Policies with conflicting requirements, **When** conflict review occurs, **Then** both versions and reasons remain visible and no automatic rejection or silent precedence is applied.

**Given** a temporary Exception Candidate, **When** it is prepared, **Then** target, scope, duration, expiry, justification and required authority are explicit and the candidate is not active without the required authority path.

**Given** no AI, **When** an exception is reviewed, **Then** Policy diffs, forms, matrices and human review provide the full function.

## 26. Questions ouvertes
OPEN-013 remains open. Final Policy precedence, exception enforcement and authority thresholds remain future; no new OPEN is created where existing Security/roadmap dependencies cover the uncertainty.

## 27. Consommateurs documentaires
Policy Gates, Approvals & Authorities, Action Center Decision Preparation, Decision Register, future Objects/Permissions/Screens/GOV-2/Technique and quality report.
