---
id: CAP-GOV-013
title: Emergency and Time-Bounded Approval Governance
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
# CAP-GOV-013 — Emergency and Time-Bounded Approval Governance

## 1. Définition
Gouverner un chemin d’autorisation d’urgence explicitement justifié et borné lorsqu’un contexte critique exige une décision plus rapide, en conservant authority, scope, target, duration/expiry, compensating controls, audit and retrospective-review requirement sans transformer urgence en bypass non gouverné.

## 2. Problème utilisateur
Une urgence réelle peut rendre le chemin normal trop lent, mais « urgent » est souvent utilisé pour contourner review, SoD ou justification. Sans limites fortes, l’urgence devient une permission générale et persistante.

## 3. Objectifs
- exiger une emergency justification et la source de criticité ;
- définir limited action, target, scope et duration ;
- vérifier emergency authority et step-up/SoD applicable ;
- documenter shortened-path candidate et compensating controls ;
- enregistrer expiry/revocation ;
- exiger retrospective verification/review comme condition future ;
- ne jamais exécuter la cible dans GOV-1.

## 4. Non-objectifs
Ne pas autoriser un bypass incontrôlé, créer une autorité d’urgence automatiquement, ignorer un Policy block sans exception/authority explicite, exécuter, démarrer Response Run, rollback ou supprimer l’audit.

## 5. Propriétaire
Govern / Approvals & Authorities owns emergency governance context and authority-bearing Approval/Decision inputs. Security owns emergency-access/step-up/SoD policy. Runtime execution remains future GOV-2/target owner.

## 6. Utilisateurs
Principal : Emergency Approver / authorized Decision Maker. Secondaires : Govern Coordinator, Security Reviewer, Incident Commander as requester/context contributor, Authority Reviewer, Auditor.

## 7. Conditions d’entrée
Action Request/version, explicit emergency justification, exact target/scope, time bounds, impact/risk, normal-path limitation or urgency reason, emergency authority source and available audit/notification context.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| emergency request/justification | requester/Incident context | urgency basis | oui | current incident/request | emergency path denied/incomplete |
| limited action/target/scope | CAP-GOV-004 | bounded effect | oui | reviewed current | blocked |
| impact/risk/reversibility | CAP-GOV-005 | consequence context | oui | current | unknown/block unless policy permits known unknown |
| Policy/conflict/exception context | CAP-GOV-007/008 | governance constraints | according to request | current | unresolved |
| emergency authority/eligible approver | CAP-GOV-009/010 | authority path | oui | current/effective | no emergency authorization |
| duration/expiry | Security/policy/request | time bound | oui | explicit | emergency path prohibited |
| compensating controls/post-review requirement | Security/Govern reviewer | guardrail context | as required | current | incomplete |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | action/target/scope/justification | read |
| Incident | Command | urgency/impact context | read/link |
| Risk/Policy/Exception/Authority assessments | Govern | governance inputs | read |
| Principal/Role | Settings | emergency approver identity | minimal read |
| Security emergency/step-up policy | Security | required controls | read/evaluate |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Emergency Authorization Context | create/review/expire/revoke/supersede | Govern local concept | explicit authority + scope + duration |
| Approval Request/Approval | create/record through CAP-GOV-011 | Govern | emergency route does not auto-approve |
| Decision Draft context | attach emergency restrictions | Govern | Decision still explicit CAP-GOV-014/015 |
| target/runtime | no mutation | future GOV-2/owner | prohibited in GOV-1 |

## 11. Fonctionnalités
Validate emergency basis; compare normal vs shortened path; require exact limited scope/target/duration; verify emergency authority and approver eligibility; apply step-up/SoD where policy says; capture compensating controls; request/record Approval explicitly; define expiry/revocation; mark retrospective verification/review requirement; preserve source/notifications/provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect emergency context | reviewer | request/context | 0 | read | sourced urgency view | non |
| validate bounds/authority/expiry | Security/Authority Reviewer | Emergency Context | 1 | rules available | explainable eligibility | non |
| request emergency route | authorized requester | context | 2 | justification + exact scope | review candidate | OPEN-013 |
| record emergency Approval/authorization | eligible emergency approver | Approval/Emergency Context | 3 | authority/SoD/step-up + time bound | authority record, no execution | required |
| finalize emergency Decision | Decision Maker | Decision | 3 | CAP-GOV-014/015 prerequisites | bounded Decision | required |
| execute target/rollback | none in GOV-1 | target | — | prohibited | no effect | future GOV-2 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate emergency required fields | oui | oui | oui | explain missing | checklist |
| calculate expiry/countdown | oui | oui | oui | no need | timestamps |
| compare normal/shortened path requirements | oui | rules/matrix | oui | summary | policy matrix |
| suggest compensating controls | oui | catalog if defined | oui | candidate only | control catalog |
| approve/decide/execute | authorized humans/future runtime | validation only | never autonomous | prohibited | explicit Govern workflow |

## 14. États fonctionnels
`requested`, `incomplete`, `under-review`, `authority-pending`, `approval-pending`, `authorized-context`, `rejected`, `expired`, `revoked`, `superseded`, `retrospective-review-required`, `retrospective-review-pending`. Post-action review execution is future GOV-2/GOV-3; GOV-1 only records the requirement.

## 15. États d’interface
Loading preserves countdown/context ; Empty no emergency request ; Partial names missing justification/authority/control ; Error preserves request and blocks authority write ; Offline cannot create authoritative emergency Approval/Decision ; Permission denied masks protected context ; Stale/expired becomes visibly unusable.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Emergency Authorization Context | time-bound governance record | CAP-GOV-011/014/015 | exact authority/scope/duration/restrictions |
| emergency Approval | Approval record | Decision Preparation | Approval ≠ Decision/execution |
| retrospective review requirement | Decision condition/handoff requirement | CAP-GOV-015/016/future GOV-2/3 | requirement preserved through handoff |
| expiry/revocation event | lifecycle record | Response Inbox/Decision history | authority cannot silently persist |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-010/012 | normal route insufficient and emergency justified | CAP-GOV-013 | request, urgency, failed/slow route, exact scope | authority route |
| CAP-GOV-013 | eligible emergency approver resolved | CAP-GOV-011 | requirement, authority, expiry, restrictions | emergency context |
| CAP-GOV-013/011 | Approval recorded | CAP-GOV-014 | emergency basis/Approval/conditions | Approvals & Authorities |
| CAP-GOV-015 | emergency Decision recorded | CAP-GOV-016 | exact time bounds/conditions/review requirement | Decision Register |
| expiry/revocation | authority no longer effective | Response Inbox/Decision history | event/reason/current status | same request |

## 18. Dépendances
CAP-GOV-003..012/014/015/016, Security emergency-access/step-up/SoD, Command Incident context, Settings identity, Shared Notifications/Trace/Versioning, OPEN-007/013.

## 19. Source de vérité
Govern owns emergency governance/Approval/Decision records; Security owns emergency policy; Command owns Incident urgency context; future GOV-2 owns execution. Emergency does not create a parallel authority system outside Govern.

## 20. Provenance et audit
Record request/version, emergency requester, explicit justification/source, normal-path limitation, target/scope, risk/reversibility, Policies/exceptions, emergency authority, approver/step-up/SoD, duration/expiry, compensating controls, Approval/Decision refs, retrospective-review requirement, expiry/revocation and notifications.

## 21. Permissions fonctionnelles
Emergency request/review; authority context read; approve/reject; Decision review/record through owner capabilities; restricted context; provenance export. High-risk emergency records are class 3 and likely step-up/SoD controlled.

## 22. Limites et erreurs
Missing justification, broad/unbounded scope, absent expiry, unavailable emergency authority, SoD failure, expired request, conflicting Policy without governed resolution, target ambiguity or offline audit path blocks emergency authority. Urgency alone never authorizes.

## 23. Métriques
Emergency requests by reason/disposition; average authorized duration; expired/revoked contexts; retrospective-review-required records; normal-path failure causes; uncontrolled bypass/auto-Approval/auto-Decision — target zero.

## 24. Classification de livraison
`defined` / `planned`; no final break-glass implementation, authentication protocol, response runtime or retrospective-review engine selected.

## 25. Critères d’acceptation
**Given** an urgent Incident but no emergency authority, **When** emergency review occurs, **Then** urgency is visible but the request is not authorized and no shortened path bypass is invented.

**Given** a valid emergency authority with a 30-minute bounded context, **When** it expires, **Then** the authority context becomes expired, remains in history and cannot be reused silently.

**Given** no AI, **When** emergency governance is performed, **Then** explicit forms, authority matrices, time bounds, checklists and human approvals provide full functionality.

## 26. Questions ouvertes
OPEN-007 and OPEN-013 remain open. Final break-glass policy, step-up strength and retrospective-review implementation remain future Security/GOV-2/GOV-3 work.

## 27. Consommateurs documentaires
Approvals & Authorities, Action Center Decision Preparation, Decision Register/Execution Handoff, Security emergency design, Command incident handoff, future Permissions/GOV-2/GOV-3 and quality gates.
