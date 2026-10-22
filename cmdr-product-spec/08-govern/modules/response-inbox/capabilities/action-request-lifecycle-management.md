---
id: CAP-GOV-003
title: Action Request Lifecycle Management
product: govern
module: response-inbox
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-003 — Action Request Lifecycle Management

## 1. Définition
Gérer le lifecycle fonctionnel Govern d’une Action Request depuis sa réception/soumission jusqu’à sa disposition décidée, son retrait, expiration ou supersession, en conservant versions, reviewers, information requests et lineage.

## 2. Problème utilisateur
Une Action Request peut évoluer pendant la review. Sans versioning/lifecycle explicites, des reviewers peuvent approuver une ancienne cible, perdre une réponse d’information ou confondre withdrawal, rejection, expiry et deletion.

## 3. Objectifs
Maintenir request/version/status ; attribuer reviewers ; gérer information requests/responses ; préserver deadlines ; permettre withdrawal/cancel/supersession ; relier les stages Policy/authority/Approval/Decision ; conserver l’historique complet.

## 4. Non-objectifs
Ne pas définir la machine d’état physique finale, supprimer l’historique, requalifier Finding/Evidence, exécuter une Decision, modifier les objets source ni créer un deuxième objet de work item.

## 5. Propriétaire
Govern / Response Inbox possède le processing lifecycle de l’Action Request. Le product source possède le contenu source et peut contribuer/retirer selon permission sans reprendre l’ownership du lifecycle post-soumission.

## 6. Utilisateurs
Principal : Govern Reviewer. Secondaires : Govern Coordinator, source requester/contributor, Policy Reviewer, Authority Reviewer, Approver, Decision Maker, Auditor.

## 7. Conditions d’entrée
Action Request identifiée, intake disposition, version courante, requester/source/return origin, tenant et permission de transition. Toute mutation requiert contrôle de concurrence/version.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| current request/version | Action Request | lifecycle anchor | oui | current | mutation interdite |
| intake/review disposition | CAP-GOV-001/002 | processing state | oui | current | remain submitted/blocked |
| information response | source product | request contribution | non | linked version | keep information-requested |
| policy/authority/approval stage results | CAP-GOV-007..013 | review progression | selon path | current source versions | do not advance silently |
| Decision disposition | CAP-GOV-015 | terminal/continuation context | non until decided | final Decision version | remain decision-ready/pending |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | version/status/requester/target/action | full lifecycle per permission |
| Review Context | Govern | reviewer, blockers, information requests | read/update |
| source objects | Command/Investigate/etc. | response/return references | read restricted |
| Policy Evaluation/Approval/Decision | Govern | progression/disposition | read/link |
| Automation Run/Tool Call | Studio | automated contribution provenance | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Action Request | transition/version/withdraw/cancel/supersede | Govern | immutable history, no delete-on-reject |
| Information Request | create/close/reopen local review need | Govern local concept | tied to exact request version |
| Review Context | assign stage/reviewer/blocker | Govern | no authority grant |
| source object | no mutation | source owner | contribution only through explicit response |

## 11. Fonctionnalités
Support lifecycle states; version changes; optimistic conflict handling conceptually; information requests and responses; reviewer assignment; stage progression; withdrawal; cancellation where allowed; expiration; supersession; lineage; Decision linkage; no destructive deletion.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect history/version | reviewer/auditor | Action Request | 0 | read | immutable lineage view | non |
| transition review stage | reviewer | Action Request | 2 | stage prerequisites | audited state transition | OPEN-013 |
| request more information | reviewer | Information Request | 2 | named gap | source response path | OPEN-013 |
| contribute information | source contributor | request contribution | 2 | same lineage + permission | new version/context | Govern re-reviews |
| withdraw request | requester/authorized owner | Action Request | 2 | execution not started; policy future | withdrawn with history | OPEN-013 |
| supersede request | authorized reviewer/requester | Action Request | 2 | replacement linked | old version retained | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate transition prerequisites | oui | oui | oui | explain only | state checklist |
| diff request versions | oui | oui | oui | summarize diff | deterministic diff |
| propose missing information | oui | field/rule checks | oui | suggestion | checklist |
| classify lifecycle stage | oui | explicit rules | oui | no authority | stage/status table |
| withdraw/supersede/decide | human-authorized | validation only | workflow routing | never autonomous | explicit action |

## 14. États fonctionnels
`draft-received`, `submitted`, `incomplete`, `under-review`, `information-requested`, `information-received`, `policy-review`, `authority-review`, `approval-pending`, `decision-ready`, `decided`, `withdrawn`, `expired`, `superseded`, `cancelled`.

These are GOV-1 functional states, not a final persisted machine. `rejected` is a Decision disposition, not deletion of the request.

## 15. États d’interface
Loading preserves version ; Empty no selected request ; Partial names unavailable review inputs ; Error preserves current valid version ; Offline disables authoritative transitions ; Permission denied masks protected source ; Stale/version-conflict requires refresh/diff before mutation.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| lifecycle transition | Action Request event | Response Inbox/Action Center | actor, old/new state, version, reason |
| information request/response link | review event | source product/Govern | same lineage and return origin |
| supersession relation | stable relation | Decision Register/history | predecessor remains resolvable |
| decision-ready state | processing state | CAP-GOV-014 | prerequisites explicit, not a Decision |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-001 | request accepted | CAP-GOV-003 | request/version/intake state | intake |
| Action Request | reviewer requests info | source product | question/version/restrictions | same request |
| source product | information response | CAP-GOV-003/006 | response refs/new version | source workspace |
| Action Request | stage prerequisites met | policy/authority/approval/decision capability | exact request version | Response Inbox/Action Center |
| replacement request | explicit supersede | historical request | new request ref/reason | current replacement |

## 18. Dépendances
CAP-GOV-001/002/004..015, Shared Versioning/Linking/Trace/Notifications, source-product handoffs, Action Request object, OPEN-013, OPEN-015.

## 19. Source de vérité
Govern is source of Action Request processing lifecycle and its version/disposition links. Source-product facts remain source-owned and referenced. Decision/Approval remain distinct Govern objects.

## 20. Provenance et audit
Every state/version change records actor, role/context, timestamp, request/version, before/after, justification, information questions/responses, source refs, automation producer if any, reviewer assignment, expiry/withdrawal/supersession and correlation ids.

## 21. Permissions fonctionnelles
Action Request read/review/context update; request information; withdraw/cancel/supersede; assignment; restricted context; automated recommendation request. Final permission granularity/step-up/SoD remain future.

## 22. Limites et erreurs
Concurrent edit, stale request, invalid transition, missing stage result, requester withdrawal after future execution start, tenant/env change, inaccessible source or Decision mismatch must fail safely with diff/reason and no history loss.

## 23. Métriques
Time by lifecycle stage; information-request loops; version conflicts; withdrawn/expired/superseded ratios; request-to-decision readiness; zero deleted history or silent stage advance.

## 24. Classification de livraison
`defined` / `planned`. No persistence, event schema, workflow engine or API is selected.

## 25. Critères d’acceptation
**Given** a request under review with missing rollback information, **When** the reviewer requests information, **Then** state becomes `information-requested`, the source receives the question and no rejection/Decision is fabricated.

**Given** a submitted request is superseded, **When** a replacement is linked, **Then** the old request remains resolvable and cannot silently become the current execution candidate.

**Given** no AI, **When** lifecycle progresses, **Then** explicit state rules, diffs and human review provide the full workflow.

## 26. Questions ouvertes
OPEN-013 remains open for class-2 mutation governance; OPEN-015 remains open for future Automation Run/Response Run bridge. Final Action Request object states are deferred to Objects.

## 27. Consommateurs documentaires
Response Inbox, Action Center, source-product handoffs, Policy/Authority/Approval/Decision capabilities, Decision Register history, future Objects/Permissions/Screens/Journeys/GOV-2/Technique and quality gates.
