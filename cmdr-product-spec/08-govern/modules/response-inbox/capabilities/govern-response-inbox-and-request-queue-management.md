---
id: CAP-GOV-002
title: Govern Response Inbox and Request Queue Management
product: govern
module: response-inbox
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-019, REQ-UX-008, REQ-UX-009]
open_decisions: [OPEN-010, OPEN-013]
source-of-truth: canonical
---
# CAP-GOV-002 — Govern Response Inbox and Request Queue Management

## 1. Définition
Fournir la queue **Govern-only** des Action Requests reçues, avec état de review, source, declared criticality, deadlines/expiration, blockers, assignment et navigation vers Action Center, sans devenir la Work Queue générale de Command.

## 2. Problème utilisateur
Les reviewers Govern doivent distinguer les demandes à gouverner du travail opérationnel général. Une queue mélangée avec Incidents/Tasks masquerait la nature d’autorité, la version de request et les délais de décision.

## 3. Objectifs
- présenter seulement le travail Govern pertinent ;
- exposer source product, request/version, action type, état, deadline/expiry, blockers et requirements ;
- permettre assignment/reassignment Govern avec audit ;
- conserver le contexte en ouvrant Action Center et au retour ;
- supporter filtres/vues locales sans créer un deuxième moteur Shared.

## 4. Non-objectifs
Ne pas remplacer Command Work Queue, créer/posséder Incident ou Task, approuver depuis la queue, exécuter une action, définir les colonnes/filters finaux ni réécrire l’écran GOV-INB-001.

## 5. Propriétaire
Govern / Response Inbox / Govern Product Lead possède les **semantics de queue Govern**. Shared possède Search/filter/assignment/Inspector mechanisms ; Command conserve sa Work Queue générale.

## 6. Utilisateurs
Principal : Govern Coordinator. Secondaires : Govern Reviewer, Policy Reviewer, Authority Reviewer, Decision Maker, Auditor en lecture.

## 7. Conditions d’entrée
Au moins une Action Request reçue ou historisée, tenant/scope courant, permission Inbox read et projection d’intake disponible. Une request inaccessible ne doit pas fuiter via les compteurs ou filtres.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Intake disposition | CAP-GOV-001 | state, blockers, source | oui | current request version | item `partial`/unavailable |
| Action Request summary | Govern | action, target/scope, requester, status | oui | current version | no actionable row |
| deadlines/expiration | request/policy/authority context | temporal review context | non selon request | current calculation/source | display unknown/not-applicable |
| assignment context | Govern + Settings identity projection | reviewer/team ownership | non | current | `unassigned` |
| source product / return origin | request provenance | navigation | oui | immutable submission origin | route blocked, item retained |
| approval/policy requirements summary | GOV-1 review projections | triage context | non at intake | latest review | pending/unknown |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | summary/version/state/deadline | read/queue |
| Govern Review Context | Govern | blockers/assignment/review stage | read/update assignment |
| Principal/Team | Settings | identity/availability projection | read minimal |
| Incident/Case/Finding | source owners | source label/return ref only | restricted projection |
| Decision/Approval | Govern | status projection when present | read only from queue |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Govern Queue Item projection | create/update from request state | Govern | no independent lifecycle from Action Request/Review Context |
| Govern assignment | assign/reassign/release | Govern | class 2, audited, does not grant authority |
| saved local view/filter state | configure | Shared mechanism / user | never becomes a new page/product owner |
| Action Request/Decision | no substantive decision mutation from queue | Govern | open Action Center for review |

## 11. Fonctionnalités
Filter/sort/search Govern requests; show source/action/declared criticality/review state/deadline/expiry/blocker/info-request/assignment; distinguish stale/partial; assign/reassign; navigate to Action Center preserving selection; bulk read/triage only where non-authoritative and reversible.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| search/filter/sort | reviewer | queue projection | 0 | read scope | scoped view | non |
| inspect request summary | reviewer | Action Request | 0 | read | Inspector/source refs | non |
| assign/reassign review | coordinator | Review Context | 2 | eligible reviewer identity | audited assignment | OPEN-013 |
| open Action Center | reviewer | request/version | 0 | request read | exact review context | non |
| bulk mark viewed/triage tag | coordinator | queue projection | 2 | non-authoritative selection | partial-success audited | OPEN-013 |
| approve/reject/decide | any | authority records | — | forbidden in Inbox | route to owner capability | required later |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| filter/sort/state grouping | oui | oui | oui | unnecessary | deterministic queue controls |
| propose reviewer assignment | oui | rules/availability | oui | suggestion only | manual assignment + authority viewer |
| summarize blockers | oui | state aggregation | oui | sourced summary | explicit fields/filters |
| flag deadline risk | oui | deterministic time rules | oui | explanation | visible timestamps/status |
| approve/decide | no from queue | no | no | prohibited | Action Center/Decision capabilities |

## 14. États fonctionnels
Queue item views include `unassigned`, `assigned`, `in-review`, `information-required`, `policy-review`, `authority-review`, `approval-pending`, `decision-ready`, `blocked`, `expired`, `withdrawn`, `superseded`, `decided`. They project source/request state; they are not a second object lifecycle.

## 15. États d’interface
Loading preserves filters/selection ; Empty distinguishes no matching requests from no permission ; Partial exposes missing projections ; Error keeps valid items ; Offline is read-only/stale ; Permission denied hides protected item existence where required ; Stale exposes request/version age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Govern queue projection | view state | Response Inbox screen | only authorized Govern work |
| assignment change | Review Context event | reviewers/Notifications | no authority grant |
| Action Center navigation context | deep link context | CAP-GOV-004..014 | exact request/version/return origin |
| deadline/blocker signal | queue context | coordinator/reviewer | source and calculation visible |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-001 | intake ready/partial | CAP-GOV-002 | request/version/intake state | intake |
| Response Inbox | open item | Action Center | request/version/selection/filter/return origin | exact queue selection |
| Response Inbox | source follow-up needed | CAP-GOV-003/006 | request + missing-info reason | same item |
| Decision recorded | projection update | Response Inbox history/Decision Register | Decision ref/disposition | source/queue |

## 18. Dépendances
CAP-GOV-001/003/004/006/009/011/015, Shared Search/Notifications/Assignments/Inspector/Versioning, Settings identities, Action Request object, OPEN-010 density and OPEN-013 class-2 governance.

## 19. Source de vérité
Action Request and Govern Review Context drive queue semantics. Shared stores generic view/filter/assignment mechanisms. Command Work Queue remains separate and authoritative for operational Incidents/Tasks.

## 20. Provenance et audit
Record queue view changes only where policy requires; always audit assignment/reassignment/release, bulk mutation, request/version opened, source/return refs, deadline source, blockers and automation suggestions/dispositions.

## 21. Permissions fonctionnelles
Govern Inbox read, Action Request read, Govern assignment/reassignment, restricted context read, cross-tenant review if explicitly authorized. Approve/Decision permissions are not exercised from Response Inbox.

## 22. Limites et erreurs
Stale version, request superseded, reviewer unavailable, assignment conflict, dependency outage, hidden source, tenant switch or expired request must preserve the item’s safe state and prevent accidental review of the wrong version.

## 23. Métriques
Queue age by review stage; unassigned count; reassignment rate; information-required/blocked/expired counts; context restoration success; zero Approval/Decision recorded directly from queue triage.

## 24. Classification de livraison
`defined` / `planned`; target Govern-native queue semantics consuming Shared mechanisms. No queue backend or final screen implementation is claimed.

## 25. Critères d’acceptation
**Given** an Action Request in `information-required`, **When** a reviewer filters that state and opens it, **Then** the same request/version and return origin open in Action Center without creating a Decision.

**Given** Command Work Queue contains an Incident, **When** Govern Response Inbox is opened, **Then** the Incident appears only through an Action Request/source projection, never as an independent Govern work item.

**Given** no AI, **When** requests are prioritized for review, **Then** deterministic filters, deadlines, states and human assignment remain sufficient.

## 26. Questions ouvertes
OPEN-010 remains open for final density by role/activity. OPEN-013 remains open for default class-2 assignment/bulk governance. No queue-specific OPEN is created.

## 27. Consommateurs documentaires
GOV-INB-001, Action Center, Decision Register, Govern navigation/maps, Shared assignment/search designs, future Permissions/Screens/Journeys/Technique and quality validation.
