---
id: CAP-GOV-012
title: Delegation, Substitution and Escalation Governance
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
# CAP-GOV-012 — Delegation, Substitution and Escalation Governance

## 1. Définition
Gouverner l’usage contextuel de delegation, substitution et escalation lorsque l’autorité ou l’approver normal est indisponible ou inadéquat, en imposant scope, duration, restrictions, source d’autorité, revocation/expiry et audit, sans transformer delegation en permanent authority ni escalation en Approval.

## 2. Problème utilisateur
Une absence d’approver peut bloquer une action urgente, mais une délégation informelle peut étendre l’autorité au-delà de son scope ou persister après sa nécessité. Une escalation peut aussi être confondue avec une validation.

## 3. Objectifs
- représenter delegator/delegate et source d’autorité ;
- borner action types, target/scope, tenant/env et duration ;
- distinguer delegation, substitution and escalation ;
- gérer unavailable approver, destination, reason, deadline ;
- permettre revocation/expiry et re-evaluation ;
- préserver SoD et préparer une nouvelle eligibility/Approval route.

## 4. Non-objectifs
Ne pas créer un Role permanent, contourner SoD, conférer une permission technique, auto-approuver, considérer escalation comme Approval, gérer emergency path (CAP-GOV-013), décider ou exécuter.

## 5. Propriétaire
Govern / Approvals & Authorities owns contextual delegation/substitution/escalation governance. Settings owns identity/admin role configuration; Security owns authority/SoD policy.

## 6. Utilisateurs
Principal : Authority Reviewer. Secondaires : authorized Delegator, Govern Coordinator, Security Reviewer, Decision Maker, delegate/approver candidate, Auditor.

## 7. Conditions d’entrée
Action Request/version, Authority Requirement, eligibility result or explicit unavailability/conflict, identity projections, applicable delegation/SoD policy and reason/deadline.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| authority requirement | CAP-GOV-009 | scope/restrictions | oui | current | no delegation route |
| eligibility/unavailability | CAP-GOV-010 | normal-route state | oui for substitution/escalation | current | reason unknown/block |
| delegator/delegate identities | Settings | identity refs | yes for delegation | current | blocked |
| delegation source/policy | Security/Govern config | permissible scope/duration | yes | effective version | no delegation inferred |
| SoD context | Security/CAP-GOV-010 | conflict controls | when applicable | current | unknown/block |
| escalation reason/destination/deadline | reviewer/policy | routing context | for escalation | current | incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | action/target/scope/version | read |
| Authority/Eligibility assessments | Govern | requirement/candidates/reasons | read |
| Principal/Role/Group | Settings | delegator/delegate identity projection | minimal read |
| SoD/authority/delegation policy | Security | constraints | read/evaluate |
| Approval Request | Govern | route requiring substitute/escalation | read/update route context |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Delegation context | create/review/revoke/expire | Govern governance concept | scoped/time-bound; no Role mutation |
| Substitution route | create/close | Govern local concept | substitute re-evaluated for eligibility |
| Escalation | create/update/close | Govern local concept | destination/reason/deadline; escalation ≠ Approval |
| Approval Request routing | update route/version | Govern | original requirement/history preserved |

## 11. Fonctionnalités
Validate delegator authority; define delegate/scope/duration/restrictions/action types; detect SoD conflicts; activate only an authorized contextual route where policy permits; revoke/expire; substitute unavailable approver; escalate to an authority destination; preserve original route; re-run eligibility; show deadlines/status.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect delegation/escalation context | reviewer | governance context | 0 | read | sourced context | non |
| validate scope/expiry/SoD | reviewer | assessment | 1 | rules available | explainable result | non |
| create/revoke permitted delegation | authorized delegator/reviewer | Delegation context | 2/3 depending authority | policy + authority + scope | time-bound route | governed |
| create substitution/escalation | coordinator/reviewer | route | 2 | reason/target/deadline | new reviewed route | OPEN-013 |
| grant permanent Role/permission | none | Settings source | — | forbidden | no mutation | Settings only |
| treat escalation as Approval | none | Approval | — | forbidden | no record | CAP-GOV-011 only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate delegation scope/expiry | oui | oui | oui | explain only | deterministic matrix |
| find alternative eligible candidates | oui | eligible-set rules | oui | suggestion | candidate table |
| propose escalation destination | oui | routing rules | oui | suggestion | authority directory/matrix |
| monitor expiry/deadline | oui | time rules | oui | no need | timestamp/status |
| grant authority/approve | authorized human governed path | validation only | no autonomous | prohibited | explicit authority/Approval |

## 14. États fonctionnels
Delegation: `proposed`, `reviewing`, `active-context`, `restricted`, `revoked`, `expired`, `superseded`.
Substitution/escalation: `needed`, `routing`, `awaiting-destination`, `accepted-for-review`, `returned`, `resolved`, `expired`, `cancelled`.
`active-context` is scoped governance context, not permanent Role/permission.

## 15. États d’interface
Loading retains original authority route ; Empty no delegation/escalation needed ; Partial exposes missing scope/duration ; Error preserves original route ; Offline cannot create authority-bearing delegation ; Permission denied masks protected identity/authority data ; Stale revalidates expiry/identity/request.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Delegation context | scoped authority route | CAP-GOV-010/011/014 | delegator/source/scope/duration/restrictions visible |
| substitution route | routing context | CAP-GOV-010/011 | substitute must be re-evaluated |
| escalation | governance event | authority destination/Response Inbox | reason/deadline/original route preserved |
| revocation/expiry | lifecycle event | Approval/Decision prep/history | no silent continuing authority |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-010 | no eligible/available normal approver | CAP-GOV-012 | requirement/candidates/reasons/deadline | eligibility review |
| CAP-GOV-012 | delegation/substitute candidate ready | CAP-GOV-010 | candidate + delegated scope/expiry/source | delegation review |
| CAP-GOV-012 | new route eligible | CAP-GOV-011 | Approval Request route + authority context | delegation context |
| CAP-GOV-012 | emergency threshold met | CAP-GOV-013 | failed normal route/urgency/scope | escalation context |
| CAP-GOV-012 | unresolved | Response Inbox | blocker/escalation destination | same request |

## 18. Dépendances
CAP-GOV-009/010/011/013/014, Settings identity/role, Security authority/SoD/step-up, Shared Notifications/Trace/Linking, OPEN-007, OPEN-013.

## 19. Source de vérité
Settings remains source for identities/Roles; Security for authority/delegation rules; Govern for request-specific delegation/substitution/escalation context. Delegation never silently becomes permanent administrative authority.

## 20. Provenance et audit
Record request/version, original authority route, delegator/delegate, authority source, action/target/scope, duration/expiry/restrictions, reason, SoD result, alternative candidates, escalation destination/deadline, acceptance/return, revocation/supersession, automation provenance and timestamps.

## 21. Permissions fonctionnelles
Delegation create/revoke; escalation; authority/eligibility read; Approval Request route update; restricted identity context; automated recommendation request; provenance export. Permanent Role administration remains Settings.

## 22. Limites et erreurs
Expired/revoked delegation, unavailable delegate, insufficient delegator authority, SoD conflict, wrong tenant/env, request scope expansion, escalation destination unavailable or failed step-up must not create a valid route; original state/history remains visible.

## 23. Métriques
Delegations by duration/disposition; expired/revoked before use; substitution/escalation rate; no-eligible-route resolution time; SoD conflicts; permanent-authority grants performed by GOV-1 — zero.

## 24. Classification de livraison
`defined` / `planned`; no final delegation object schema, directory integration, escalation engine or privilege system selected.

## 25. Critères d’acceptation
**Given** an eligible approver becomes unavailable, **When** substitution is proposed, **Then** the substitute is re-evaluated against the same authority/SoD requirements and unavailability does not itself grant authority.

**Given** a delegation reaches expiry, **When** an Approval Request is opened, **Then** the delegated context cannot satisfy eligibility and the route is blocked/re-evaluated.

**Given** no AI, **When** delegation/escalation is managed, **Then** authority matrices, eligibility tables, deadlines and human routing provide full functionality.

## 26. Questions ouvertes
OPEN-007 and OPEN-013 remain open. Final delegation/substitution authority policy and step-up requirements remain future Security/Permissions work.

## 27. Consommateurs documentaires
Approvals & Authorities, Response Inbox, Approval/Decision capabilities, Security/Settings, future Objects/Permissions/Screens/GOV-2 and quality gates.
