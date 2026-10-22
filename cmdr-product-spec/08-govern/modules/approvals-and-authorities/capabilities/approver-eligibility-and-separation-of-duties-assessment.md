---
id: CAP-GOV-010
title: Approver Eligibility and Separation-of-Duties Assessment
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
# CAP-GOV-010 — Approver Eligibility and Separation-of-Duties Assessment

## 1. Définition
Évaluer quels approver candidates sont contextuellement éligibles à satisfaire un Authority Requirement, en tenant compte de l’autorité, du scope, du tenant/environnement, de la relation au requester/action owner, des conflits d’intérêt, de SoD, delegation, disponibilité et expiration, sans sélectionner silencieusement l’approver réel.

## 2. Problème utilisateur
Une liste de personnes possédant un Role peut inclure le requester, un owner en conflit, une délégation expirée ou une autorité hors scope. Sans éligibilité explicable, un système peut créer de l’auto-approbation ou une escalade de privilège implicite.

## 3. Objectifs
- établir la population de candidates à partir des sources autorisées ;
- tester chaque candidate contre authority/scope/expiry/restrictions ;
- évaluer SoD et conflits d’intérêt ;
- distinguer candidate, eligible et actual approver ;
- exposer ineligibility reason et alternatives ;
- représenter N-of-M uniquement comme exigence fonctionnelle, pas algorithme final.

## 4. Non-objectifs
Ne pas administrer identités/rôles, auto-assigner une Approval, considérer availability comme authority, finaliser la matrice SoD, approuver, décider ou exécuter.

## 5. Propriétaire
Govern / Approvals & Authorities owns request-specific eligibility assessment. Settings owns identities/roles/groups. Security owns SoD/authority policy. Govern does not mutate those sources.

## 6. Utilisateurs
Principal : Authority Reviewer. Secondaires : Govern Coordinator, Security Reviewer, Approver candidate, Decision Maker, requester as read-limited consumer, Auditor.

## 7. Conditions d’entrée
Authority Requirement/Context from CAP-GOV-009, Action Request/requester, target/scope, tenant/environment, candidate identity projections and SoD/eligibility rules accessible to the reviewer.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Authority Requirement | CAP-GOV-009 | required authority/scope | oui | current request | no eligibility assessment |
| candidate Principals/groups | Settings identity source | candidate population | oui | current | no candidate invented |
| requester/action owner relations | Action Request/source context | SoD relation | oui | request version | unknown/block |
| SoD/conflict rules | Security | eligibility restrictions | yes when applicable | effective version | unknown, not pass |
| delegation/substitution context | CAP-GOV-012 / source | scoped authority path | non | current/effective | ignored only if not needed |
| availability/expiry | Settings/Govern context | practical eligibility | according to route | current | candidate unavailable/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Authority Requirement/Context | Govern | scope/restrictions | read |
| Action Request | Govern | requester/action/target/version | read |
| Principal/Role/Group | Settings | identity/role/group/availability projection | minimal read |
| SoD/Decision Authority policy | Security | conflict/eligibility rules | read/evaluate |
| Delegation | Govern/Settings/Security | scope/duration/restrictions | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Approver Candidate projection | list/assess eligible/ineligible | Govern local concept | reason/source required |
| Eligibility Assessment | create/version/supersede | Govern local concept | bound to request + authority versions |
| Approval Request routing context | prepare candidates only | Govern | actual approver selection is explicit later |
| Principal/Role/Group | no mutation | Settings | no privilege grant |

## 11. Fonctionnalités
Build candidate set; compare authority scope/tenant/env/action; detect requester==candidate and other conflicts; evaluate SoD; evaluate delegation/expiry/availability; show eligible/ineligible reasons; support alternate candidates/escalation need; model N-of-M as requirement metadata; re-evaluate after authority/request/identity changes.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect candidates/authority | Authority Reviewer | candidate projections | 0 | read | sourced list | non |
| run eligibility/SoD check | Authority Reviewer | Eligibility Assessment | 1 | rules + authority + identities | explainable eligible/ineligible results | non |
| annotate/challenge eligibility | reviewer | assessment | 2 | rationale | versioned review | OPEN-013 |
| select candidate for Approval Request | coordinator/reviewer | routing context | 2 | candidate eligible | explicit route, not Approval | OPEN-013 |
| grant role/override SoD | none in GOV-1 | source config | — | forbidden | no mutation | source/future authority path |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| find candidate population | oui | directory/filter rules | oui | candidate suggestion | authority/identity filters |
| evaluate authority scope | oui | oui | oui | explanation | deterministic matrix |
| evaluate SoD/conflict | oui | explicit rules | oui | conflict suggestion | SoD matrix |
| propose alternative approver | oui | eligible-set rules | oui | suggestion | sorted eligible list |
| approve/override | authorized human later | validation only | no autonomous | prohibited | CAP-GOV-011/015 |

## 14. États fonctionnels
Assessment: `not-assessed`, `candidates-found`, `no-candidate`, `partial`, `sod-review`, `eligible-set-ready`, `conflicted`, `stale`, `superseded`.
Candidate: `candidate`, `eligible`, `ineligible-authority`, `ineligible-scope`, `ineligible-sod`, `unavailable`, `delegation-expired`, `step-up-required`.

## 15. États d’interface
Loading keeps requirement/request ; Empty distinguishes no candidate from no access ; Partial lists unverified eligibility dimensions ; Error keeps previous assessment ; Offline cannot finalize new eligibility ; Permission denied masks sensitive identity/authority data ; Stale requires reevaluation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Eligibility Assessment | review result | CAP-GOV-011/014/015 | candidate reasons/authority/SoD sources visible |
| eligible candidate set | projection | Approval routing | candidate ≠ actual approver |
| no-eligible-approver blocker | review context | Response Inbox/CAP-GOV-012 | explicit reason, no fallback privilege |
| step-up/delegation need | governance context | CAP-GOV-012/013 | requirement only, no authority created |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-009 | authority requirement ready | CAP-GOV-010 | authority scope/restrictions/requester | authority review |
| CAP-GOV-010 | eligible set ready | CAP-GOV-011 | candidate set + reasons + SoD result | eligibility review |
| CAP-GOV-010 | no eligible candidate | CAP-GOV-012 | missing route/conflict/deadline | same request |
| CAP-GOV-010 | emergency path required | CAP-GOV-013 | normal-path failure + scope/urgency | eligibility review |
| CAP-GOV-010 | decision prep | CAP-GOV-014 | eligibility/SoD disposition | Action Center |

## 18. Dépendances
CAP-GOV-003/004/009/011/012/013/014/015, Settings Principal/Role/Group projections, Security SoD/Decision Authority/step-up, Shared Trace/Linking/Notifications, OPEN-007 and OPEN-013.

## 19. Source de vérité
Settings owns identity/Role/Group configuration; Security owns SoD/authority policy; Govern owns the request-specific Eligibility Assessment. Candidate/eligible status does not grant permission or authority beyond the evaluated request context.

## 20. Provenance et audit
Record request/authority versions, candidate-source query/context, each candidate identity ref, authority scope/expiry, requester relation, SoD rule/version, delegation/availability, eligibility result/reason, reviewer, overrides attempted/denied, automation provenance and timestamps.

## 21. Permissions fonctionnelles
Approver candidates read, eligibility review, restricted identity/authority context, routing preparation, escalation, automated recommendation request and provenance export. No user/role administration is added.

## 22. Limites et erreurs
Identity unavailable, stale role, unclear requester relation, conflicting SoD rules, candidate outside tenant, delegation expiry, no eligible candidate or permission denial produce explicit no-candidate/unknown/blocker states; never auto-select an ineligible fallback.

## 23. Métriques
Candidate-to-eligible ratios; no-eligible-approver cases; SoD conflicts; expired delegation/unavailable candidates; reassessment after identity/request change; zero requester self-approval under applicable SoD.

## 24. Classification de livraison
`defined` / `planned`; no final N-of-M algorithm, SoD matrix, directory integration or approver-routing engine selected.

## 25. Critères d’acceptation
**Given** requester and candidate are the same Principal and applicable SoD prohibits self-approval, **When** eligibility is evaluated, **Then** the candidate is ineligible with explicit reason and cannot be silently selected.

**Given** a candidate has the correct Role but authority scope excludes the target environment, **When** evaluated, **Then** Role is not treated as contextual authority and the candidate is ineligible.

**Given** no AI, **When** candidates are evaluated, **Then** authority matrices, identity projections and deterministic SoD checks provide full functionality.

## 26. Questions ouvertes
OPEN-007 remains open for Human Gate/Govern semantics. OPEN-013 remains open for class-2 routing governance. Final N-of-M and SoD policy remain future; no new OPEN is created.

## 27. Consommateurs documentaires
Approvals & Authorities, Approval Management, Delegation/Emergency, Decision Preparation/Register, Security/Settings, future Objects/Permissions/Screens/Journeys/GOV-2 and conformance report.
