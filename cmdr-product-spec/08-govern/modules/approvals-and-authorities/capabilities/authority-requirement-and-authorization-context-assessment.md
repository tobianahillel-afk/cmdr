---
id: CAP-GOV-009
title: Authority Requirement and Authorization Context Assessment
product: govern
module: approvals-and-authorities
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-015, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013]
source-of-truth: canonical
---
# CAP-GOV-009 — Authority Requirement and Authorization Context Assessment

## 1. Définition
Déterminer quelle autorité contextuelle est requise pour gouverner une Action Request donnée et évaluer le contexte d’autorisation disponible selon action type/class, target/scope, tenant/environment, impact, requester, duration, delegation and restrictions, sans confondre Role, CRUD permission ou configuration administrative avec authority.

## 2. Problème utilisateur
Un utilisateur peut posséder une permission technique ou un Role mais ne pas être autorisé à approuver une action donnée sur une cible donnée. Une matrice statique de rôles peut ignorer scope, impact, expiry ou delegation et conduire à une fausse autorisation.

## 3. Objectifs
- calculer/établir un Authority Requirement explicable ;
- lier l’autorité à l’action, classe, target/scope, tenant/env, impact et duration ;
- distinguer source d’autorité configurée et autorité contextuelle évaluée ;
- représenter missing/conflicting/restricted/expired authority ;
- transmettre des exigences précises à l’éligibilité des approvers et à Decision Preparation.

## 4. Non-objectifs
Ne pas administrer Users/Roles/Groups, accorder une permission, inventer une authority, sélectionner automatiquement l’approver final, enregistrer une Approval ou Decision, exécuter l’action ou finaliser le RBAC/ABAC.

## 5. Propriétaire
Govern / Approvals & Authorities owns contextual authority assessment. Platform Settings owns administrative identity/role/group configuration; Security owns permission/Decision Authority policy; Govern consumes those sources and records the request-specific assessment.

## 6. Utilisateurs
Principal : Authority Reviewer. Secondaires : Govern Reviewer, Security Reviewer, Decision Maker, Approver candidate, Platform Administrator as source owner, Auditor.

## 7. Conditions d’entrée
Current Action Request/version, reviewed target/scope, impact/risk/reversibility, Policy Evaluation/Exception context as applicable, requester identity and authorized access to authority source projections.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| action type/class | Action Request + classification | authority driver | oui | request version | authority unknown/blocked |
| target/scope/tenant/environment | CAP-GOV-004 | authority boundary | oui | reviewed current | blocked |
| impact/risk/reversibility | CAP-GOV-005 | authority threshold context | according to policy | current assessment | unknown/partial |
| Policy/Exception context | CAP-GOV-007/008 | additional requirement/restriction | when applicable | current evaluations | unresolved |
| requester identity/relations | Settings identity projection | requester context | oui | current | blocked |
| configured authority sources | Security/Settings | roles, delegation, scope, expiry, restrictions | oui | current/effective | `authority-context-required` |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | action/requester/version | read |
| Context/Risk/Policy reviews | Govern | target/scope/impact/outcomes | read |
| Principal/Role/Tenant/Environment | Settings | identity/admin context | minimal read |
| permission/Decision Authority policy | Security | authority rules/constraints | read/evaluate |
| Delegation | Govern/Settings/Security context | scope/duration/restrictions | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Authority Requirement | create/version/supersede | Govern local concept | request-specific; source rules cited |
| Authority Context | assess/version/mark missing/conflict/expired | Govern | configured authority ≠ contextual authority |
| Review Context | attach authority result | Govern | no permission/Role mutation |
| User/Role/Permission | no mutation | Settings/Security | source projection only |

## 11. Fonctionnalités
Resolve authority rules; derive required authority categories/scope/duration; compare requester/action/target context; preserve restrictions and expiry; evaluate configured/delegated authority candidates; distinguish missing/conflicting authority; show source and rationale; re-evaluate after material request change; hand off to approver eligibility.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect authority source | Authority Reviewer | authority projection | 0 | read | sourced context | non |
| run authority requirement assessment | Authority Reviewer | Authority Requirement | 1 | request/context available | explainable requirement | non |
| annotate/challenge authority context | reviewer | Authority Context | 2 | rationale | versioned review | OPEN-013 |
| request missing authority information | reviewer | Information Request | 2 | missing source named | clarification/escalation | OPEN-013 |
| grant role/permission/authority | none in GOV-1 | Settings/Security config | — | forbidden | no mutation | source owner only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| determine required authority | oui | rule/matrix when defined | oui | candidate/explanation | authority matrix + checklist |
| compare scope/expiry/restrictions | oui | yes | oui | summary | deterministic comparison |
| identify missing/conflicting authority | oui | yes | oui | suggestion | rule diagnostics |
| propose approver candidate classes | oui | catalog/rules | oui | suggestion | authority viewer |
| invent/grant authority | no | no | no | prohibited | Settings/Security admin process |

## 14. États fonctionnels
`not-assessed`, `requirements-known`, `authority-context-required`, `authority-found`, `authority-partial`, `authority-conflict`, `authority-expired`, `authority-restricted`, `delegated-context`, `reviewed`, `superseded`, `blocked`.

## 15. États d’interface
Loading preserves request/source versions ; Empty distinguishes no authority source vs not reviewed ; Partial names missing dimensions ; Error keeps prior valid assessment ; Offline no authoritative reassessment ; Permission denied masks protected authority details ; Stale forces re-evaluation after identity/policy/request changes.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Authority Requirement | review record | CAP-GOV-010/011/014/015 | action/target/scope/source/rationale visible |
| Authority Context Assessment | review result | approver eligibility/Decision prep | configured role not treated as authority by itself |
| missing/conflict condition | blocker/context | Response Inbox/Action Center | no implicit fallback authority |
| escalation need | governance context | CAP-GOV-012 | reason/scope/deadline preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Action Request | authority review required | CAP-GOV-009 | request/context/risk/policy/requester | Action Center |
| CAP-GOV-009 | requirement resolved | CAP-GOV-010 | required authority, scope, restrictions, requester relation | authority review |
| CAP-GOV-009 | authority unavailable/conflicted | CAP-GOV-012 | missing authority, reason, deadline | same request |
| CAP-GOV-009 | decision prep input ready | CAP-GOV-014 | authority requirement/context/unknowns | Action Center |

## 18. Dépendances
CAP-GOV-003..008/010/012/014/015, Security permission/Decision Authority/step-up, Settings users/roles/tenant/env, Shared Trace/Linking, OPEN-007 and OPEN-013.

## 19. Source de vérité
Settings is source for administrative identity/role/group configuration. Security is source for permission/authority policy. Govern is source of the request-specific Authority Requirement/Context Assessment. Role/permission/configuration never silently becomes contextual authority.

## 20. Provenance et audit
Record request/version, action class, target/scope, impact, Policy/Exception refs, requester, authority source/rule/version, delegation source, scope/duration/expiry/restrictions, assessment method, reviewer, conflicts/unknowns, automation provenance and timestamps.

## 21. Permissions fonctionnelles
Authority Context read; authority assessment; restricted identity/authority projection; eligibility handoff; escalation; automated recommendation request; provenance export. User/Role/permission administration remains outside Govern.

## 22. Limites et erreurs
Missing identity, stale Role, expired delegation, mismatched tenant/env, ambiguous target, unknown authority rule, inaccessible policy, conflicting sources or step-up failure leaves authority unknown/blocked and cannot be replaced by a guessed role.

## 23. Métriques
Requests by required authority type; missing/conflict/expired authority; reassessments after scope changes; time to resolve authority gaps; zero invented or automatically granted authority.

## 24. Classification de livraison
`defined` / `planned`; no final authority engine, algorithm, RBAC/ABAC mapping or admin API selected.

## 25. Critères d’acceptation
**Given** a user with a platform admin Role but no contextual response authority, **When** authority assessment runs, **Then** the Role is visible as administrative context but the user is not treated as authorized approver/Decision Maker.

**Given** a delegated authority expired before review, **When** assessed, **Then** it is marked expired and cannot satisfy the requirement.

**Given** no AI, **When** authority is assessed, **Then** explicit matrices, scope comparisons and human review provide the full function.

## 26. Questions ouvertes
OPEN-007 remains open for Human Gate/Govern relationship. OPEN-013 remains open for default class-2 governance. Final authority policy/permission atomization remains future; no new OPEN is created.

## 27. Consommateurs documentaires
Approvals & Authorities, Decision Preparation/Register, Response Inbox blockers, Security/Settings integrations, future Objects/Permissions/Screens/Journeys/GOV-2/Technique and quality gates.
