---
id: CAP-GOV-011
title: Approval Request and Approval Record Management
product: govern
module: approvals-and-authorities
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013]
source-of-truth: canonical
---
# CAP-GOV-011 — Approval Request and Approval Record Management

## 1. Définition
Créer, router, suivre et enregistrer des Approval Requests et Approvals liés à une Action Request/version et à une Authority Requirement, en conservant approver, deadline/expiry, conditions, justification, authority context, version, révocation/supersession éventuelle et historique, sans confondre Approval avec Decision ou execution.

## 2. Problème utilisateur
Une approbation donnée hors contexte, sur une ancienne request version, par une personne non éligible ou sans justification peut être réutilisée comme autorité indue. Inversement, une Approval valide ne doit pas être prise pour une Decision finale ni déclencher automatiquement une action.

## 3. Objectifs
- créer une Approval Request explicite à partir d’un requirement et d’un approver eligible ;
- lier exact Action Request/version/authority context ;
- gérer deadline, expiry, information request, approve/reject/abstain quand applicable ;
- supporter delegation seulement via CAP-GOV-012 ;
- enregistrer rationale/conditions/timestamp/version ;
- permettre revoke-before-execution where policy permits et supersession sans suppression ;
- fournir un input vérifiable à Decision Preparation.

## 4. Non-objectifs
Ne pas prendre la Decision, exécuter l’action, démarrer Response Run, considérer Human Gate comme Approval, auto-approuver, créer l’autorité manquante ou finaliser la signature/protocole technique d’Approval.

## 5. Propriétaire
Govern / Approvals & Authorities owns Approval Request and Approval lifecycle. Studio owns Human Gate; Settings owns identity/admin config; Security supplies authority/SoD/step-up constraints.

## 6. Utilisateurs
Principal : authorized Approver. Secondaires : Approval Coordinator, Authority Reviewer, Security Reviewer, Decision Maker, requester in read-limited state, Auditor.

## 7. Conditions d’entrée
Current Action Request/version, explicit Approval Requirement, Authority Context, eligible approver result, SoD outcome and permission/step-up context. A candidate or eligible approver is not yet the actual approver until routed explicitly.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Action Request/version | Govern | approval subject | oui | exact reviewed version | no Approval Request |
| Approval Requirement | CAP-GOV-009/Policy context | authority need | oui | current | blocked |
| eligible approver | CAP-GOV-010 | candidate eligibility | oui | current assessment | no route |
| SoD/step-up context | Security/CAP-GOV-010 | authorization control | yes when applicable | current | pending/blocked |
| Policy/Exception context | CAP-GOV-007/008 | conditions/approval need | according to request | current | pending/unknown |
| deadline/expiry | requirement/policy/request | temporal bound | according to policy | current | explicit unknown/not-applicable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | version/action/target/scope | read |
| Authority Requirement/Context | Govern | authority/scope/restrictions | read |
| Eligibility Assessment | Govern | eligible/ineligible reasons | read |
| Policy Evaluation/Exception Candidate | Govern | conditions and context | read |
| Principal | Settings | approver identity projection | minimal read |
| Human Gate/Automation Run | Studio | provenance/automation context only | read/link, never treat as Approval |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Approval Request | create/update/cancel/expire/supersede | Govern local concept | bound to request/version/requirement |
| Approval | create by authorized disposition/read/revoke where policy permits/supersede | Govern | authority context + approver + rationale required |
| information request | create/resolve | Govern local review concept | does not equal rejection |
| Action Request | attach approval refs/stage | Govern | Approval does not finalize Decision |

## 11. Fonctionnalités
Create Approval Request; route to eligible approver; present exact request/target/scope/authority context; request information; approve/reject/abstain if applicable; capture rationale/conditions; apply deadline/expiry; record step-up/SoD result references; support permitted delegation route; revoke before future execution if policy permits; supersede after material request change; preserve history and feed Decision Preparation.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| read Approval Request/context | approver/reviewer | Approval Request | 0 | authorized read | exact context view | non |
| create/cancel Approval Request | coordinator/reviewer | Approval Request | 2 | eligible route + requirement | versioned request | OPEN-013 |
| request more information | approver | info request | 2 | named question | pending response | OPEN-013 |
| approve/reject/abstain | eligible authorized approver | Approval | 3 where authority-bearing | current request/version, SoD/step-up satisfied | Approval record only | yes authority-bearing, no execution |
| revoke before execution where permitted | authorized approver/policy owner | Approval | 2/3 | policy allows + execution not started | revocation record | governed |
| execute target | none in GOV-1 | target | — | prohibited | no effect | future GOV-2 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| create/rout Approval Request | oui | route validation | oui | no autonomous approver choice | eligible-set selection |
| validate request/version/eligibility | oui | oui | oui | explain only | deterministic checks |
| summarize context | oui | structured aggregation | oui | sourced summary | source tables/diffs |
| draft rationale/conditions | oui | templates | oui | draft | manual template |
| record Approval disposition | authorized human | validation only | notification/workflow routing | **never autonomous** | explicit approver action |

## 14. États fonctionnels
Approval Request: `draft`, `pending`, `information-requested`, `information-received`, `ready`, `expired`, `cancelled`, `superseded`, `closed`.
Approval: `approved`, `rejected`, `abstained` where policy permits, `revoked` where permitted, `expired`, `superseded`. Final object state machine remains future.

## 15. États d’interface
Loading keeps request/version ; Empty means no Approval required/requested ; Partial names missing authority/context ; Error preserves pending request ; Offline never records an authoritative Approval without guaranteed write ; Permission denied masks restricted facts ; Stale forces revalidation if request/authority/Policy changed.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Approval Request | versioned governance request | Approver/CAP-GOV-014 | exact request/version/requirement/authority refs |
| Approval | authority record | CAP-GOV-014/015 | approver, authority, rationale, conditions, time and version preserved |
| information request/response | review event | source/reviewer | same Approval Request lineage |
| revocation/expiry/supersession | lifecycle event | Decision prep/history | no deletion or silent replacement |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-010 | eligible candidate selected | CAP-GOV-011 | request/version/requirement/authority/SoD | eligibility review |
| Approval Request | approver requests info | source/Govern review | questions/request/version | same Approval Request |
| Approval Request | authorized disposition | Approval | exact context + approver/authority | Approvals & Authorities |
| Approval | record available | CAP-GOV-014 | disposition/rationale/conditions/expiry | Approval context |
| material request change | reapproval required | CAP-GOV-010/011 | old Approval + new request version | current request |

## 18. Dépendances
CAP-GOV-003/007..010/012/013/014/015, Approval object, Security Approval/Decision Authority/SoD/step-up, Settings identity, Studio Human Gate/Automation provenance, Shared Notifications/Trace/Versioning, OPEN-007, OPEN-013.

## 19. Source de vérité
Govern owns Approval Request/Approval. Settings owns Principal identity/config; Security owns authority policy; Studio Human Gate remains distinct. Approval can satisfy an authority requirement but is not itself a Decision or execution.

## 20. Provenance et audit
Record request/version, requirement, eligible-set assessment, actual approver, authority source/scope/expiry, SoD/step-up context, questions/responses, disposition, rationale, conditions, timestamp, revocation/expiry/supersession, automation provenance and correlation ids.

## 21. Permissions fonctionnelles
Approval Request create/update/cancel; Approval read; approve/reject/abstain; delegation route; revoke where permitted; restricted context; provenance export. High-risk Approval can require step-up/SoD and class 3 authority.

## 22. Limites et erreurs
Stale request, ineligible approver, expired authority/delegation, failed step-up, SoD conflict, missing Policy context, deadline expiry, unavailable source or concurrent disposition must refuse authoritative write and preserve explainable pending/error history.

## 23. Métriques
Approval Requests by outcome; information loops; expiry/revocation; reapproval after material changes; SoD/step-up blocks; time from request to disposition; autonomous Approvals and requester self-approval under SoD — target zero.

## 24. Classification de livraison
`defined` / `planned`; no final e-signature, quorum implementation, API/protocol or notification runtime selected.

## 25. Critères d’acceptation
**Given** an Approval Request routed to an eligible approver, **When** the approver approves, **Then** an Approval record is created with authority/rationale/version, but no Decision or execution is created automatically.

**Given** requester == approver and SoD prohibits it, **When** approval is attempted, **Then** the authoritative write is refused and the conflict is audited.

**Given** no AI, **When** Approval is requested/reviewed, **Then** eligibility data, forms, context tables and explicit human disposition provide full functionality.

## 26. Questions ouvertes
OPEN-007 remains open for Human Gate/Approval relation; OPEN-013 remains open for class-2 defaults. Final quorum/N-of-M, step-up and revocation policies remain future.

## 27. Consommateurs documentaires
Approvals & Authorities, Action Center Decision Preparation, Decision Register, source-product projections, Security/Settings/Studio boundaries, future Objects/Permissions/Screens/GOV-2 and quality gates.
