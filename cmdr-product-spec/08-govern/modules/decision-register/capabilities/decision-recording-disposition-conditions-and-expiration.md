---
id: CAP-GOV-015
title: Decision Recording, Disposition, Conditions and Expiration
product: govern
module: decision-register
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-015, REQ-PROD-020, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-015 — Decision Recording, Disposition, Conditions and Expiration

## 1. Définition
Enregistrer une Decision Govern explicite pour une Action Request/version, avec Decision Maker, authority context, disposition, justification, source facts, uncertainty, conditions, permitted/prohibited scope, target, time bounds, rollback/verification requirements, Approvals, dissent, timestamp, version, expiration and supersession, sans exécuter l’action.

## 2. Problème utilisateur
Une Approval, recommendation ou Policy outcome peut être prise à tort pour la décision finale. Une Decision sans scope/conditions/expiry ou version exacte peut aussi être exécutée au-delà de l’autorité réellement accordée.

## 3. Objectifs
- enregistrer une disposition autoritative distincte des inputs ;
- supporter `approve`, `approve-with-conditions`, `reject`, `defer`, `request-more-information`, `cancel`, `supersede` ;
- préserver Decision Maker/authority/Approvals/rationale/dissent ;
- borner permitted/prohibited scope, target, time and start constraints ;
- conserver rollback/verification requirements ;
- gérer expiration/supersession sans deletion ;
- produire une Decision exploitable par CAP-GOV-016 et les produits source.

## 4. Non-objectifs
Ne pas exécuter, démarrer Response Run, produire Result, effectuer rollback, modifier Evidence/Finding, transformer Approval en Decision, auto-décider ou finaliser la machine d’état physique.

## 5. Propriétaire
Govern / Decision Register owns Decision and its authority-bearing disposition. Security defines authority policy; source facts remain with source owners; future GOV-2 owns Response Run execution/verification/rollback.

## 6. Utilisateurs
Principal : authorized Decision Maker. Secondaires : Decision Reviewer, Govern Reviewer, Approvers, Security/Authority Reviewer, requester/source owners in authorized read, Auditor.

## 7. Conditions d’entrée
Decision Draft and `decision-ready` context, current Action Request/version, contextual authority, required Approvals/SoD/step-up disposition, applicable Policy/exception results, exact target/scope, known uncertainty and permission to record the Decision.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Decision Draft + readiness | CAP-GOV-014 | reviewed decision package | oui | exact current version | no Decision record |
| Action Request/version | Govern | decision subject | oui | current | no Decision |
| authority/Decision Maker context | CAP-GOV-009/010/Security | authority | oui | effective at decision time | blocked |
| required Approvals/SoD/step-up | CAP-GOV-011..013 | authority prerequisites | according to action | current/unexpired | blocked or explicit not-required basis |
| Policy/conflict/exception results | CAP-GOV-007/008 | governance constraints | according to applicability | current | unresolved/blocked |
| target/scope/risk/reversibility | CAP-GOV-004/005 | effect bounds | oui for effectful action | current reviewed | blocked/unknown explicit |
| source facts/uncertainty/dissent | CAP-GOV-006/014/source | rationale context | oui as applicable | referenced versions | not hidden |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Action Request | Govern | exact version/action/target/scope | read |
| Decision Draft | Govern | options/conditions/readiness | read |
| Policy Evaluation/Exception | Govern | outcomes/authority refs | read |
| Authority/Approval/Emergency context | Govern/Security | validity/scope/conditions | read/re-evaluate |
| Incident/Case/Finding/Evidence | source owners | source facts/refs/restrictions | read/link, no requalification |
| Workflow/Human Gate/Automation Run | Studio | provenance only | read/link |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Decision | record/version/supersede | Govern | immutable historical version; authority context required |
| Decision Condition | record/read/supersede with Decision | Govern local concept | condition ≠ Policy; exact scope/expiry |
| Decision Expiration | record/derive from explicit bound | Govern local concept | expiry ≠ deletion |
| Action Request processing state | mark decided/deferred/info-required/cancelled as applicable | Govern | request history retained |
| target/source objects | no mutation | source/runtime owners | Decision is authority only |

## 11. Fonctionnalités
Validate current inputs; record Decision Maker/authority; choose explicit disposition; capture rationale/source facts/uncertainty; record conditions, permitted/prohibited scope, target, time/start constraints, expiry, rollback/verification requirements, Approvals and dissent; return information/reject/defer outcomes to source; supersede old Decision with links; expose expiration; prevent stale Decision from silently authorizing future handoff.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect Decision package/history | reviewer/auditor | Decision Draft/Decision | 0 | read | sourced context | non |
| run final deterministic precondition check | Decision Maker | draft package | 1 | all inputs available | missing/stale list | non |
| update pre-final conditions | authorized reviewer/Decision Maker | Decision Draft | 2 | before finalization + re-eval as needed | versioned condition diff | OPEN-013 |
| record reject/defer/info-request/cancel | Decision Maker | Decision | 3 when authority-bearing governance act | authority + current package | authoritative disposition, no execution | yes |
| record approve/approve-with-conditions | Decision Maker | Decision | 3 | authority/Approvals/SoD/Policy current | authoritative Decision only | yes |
| supersede Decision | authorized Decision Maker | Decision | 2/3 | replacement authority and lineage | old remains resolvable | governed |
| execute target | none in GOV-1 | target | — | prohibited | no effect | future GOV-2 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| validate freshness/prerequisites | oui | oui | oui | explain only | final checklist |
| compare Decision conditions/options | oui | deterministic diff | oui | summary/suggestion | comparison table |
| draft rationale | oui | template | oui | draft | manual rationale template |
| calculate explicit expiry from entered bound | oui | oui | oui | no need | timestamp calculator |
| select disposition/finalize | authorized human | validation only | **never autonomous** | prohibited | explicit Decision Maker action |

## 14. États fonctionnels
For GOV-1 functional semantics: `pending-record`, `approved`, `approved-with-conditions`, `rejected`, `deferred`, `information-required`, `cancelled`, `expired`, `superseded`. These do not replace the future canonical object state-machine review. `expired` preserves the record; `approved` does not mean executed.

## 15. États d’interface
Loading preserves draft/input versions ; Empty no Decision yet ; Partial lists restricted/unknown source facts ; Error never writes a partial authoritative Decision ; Offline prevents authoritative finalization unless future guaranteed control exists ; Permission denied masks protected rationale/context ; Stale invalidates finalization/handoff until reevaluated.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Decision | authority record | CAP-GOV-016/source products/future GOV-2 | exact request/version/authority/disposition/scope/conditions/time preserved |
| Decision disposition event | return context | Command/Investigate/Detection/TI/source | reject/defer/info/approve semantics explicit |
| Decision conditions/expiry | constraint set | CAP-GOV-016/future GOV-2 | condition ≠ Policy; expiry ≠ deletion |
| supersession relation | immutable history link | Decision Register/Audit | old Decision remains resolvable |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-014 | decision-ready | CAP-GOV-015 | exact Draft and all review refs | Action Center |
| CAP-GOV-015 | reject/defer/info-required | source product | Decision ref/rationale/allowed details/next step | Decision Register/source object |
| CAP-GOV-015 | approve/approve-with-conditions | CAP-GOV-016 | Decision + scope/targets/conditions/expiry/requirements | Decision Register |
| CAP-GOV-015 | superseded | Decision Register history | old/new refs/rationale/effective status | current Decision |
| CAP-GOV-015 | expired before handoff/use | Response Inbox/Decision history | expiry reason/status | source/Decision Register |

## 18. Dépendances
CAP-GOV-003..014/016, Decision/Approval/Policy/Action Request objects, Security Decision Authority/SoD/step-up, Shared Versioning/Trace/Linking/Notifications, Studio provenance, OPEN-007/013/015.

## 19. Source de vérité
Govern Decision is the source of decision disposition and its conditions for the exact request version. It does not become source of Evidence/Finding/Incident facts or runtime target state. Approval and Policy Evaluation remain distinct input records.

## 20. Provenance et audit
Record Decision id/version, request/draft versions, Decision Maker, authority source/scope/expiry, SoD/step-up, Approvals, Policies/conflicts/exceptions, source facts/uncertainty, alternatives, dissent, disposition/rationale, conditions/permitted/prohibited scope, targets, time/start constraints, rollback/verification requirements, timestamp, expiration/supersession and automation provenance.

## 21. Permissions fonctionnelles
Decision Draft read/review; Decision record; approve/reject/defer/request-information/cancel; conditions update before finalization; Decision supersede; restricted context; provenance export. Authority and step-up are re-evaluated independently from CRUD permission.

## 22. Limites et erreurs
Expired Approval/authority, request version drift, Policy reevaluation needed, target/scope change, SoD conflict, missing source fact, offline state, concurrent Decision, failed audit write or expired emergency context prevents authoritative finalization/handoff. No silent fallback Decision.

## 23. Métriques
Decisions by disposition; approve-with-conditions rate; defer/info cycles; expired/superseded Decisions; condition violations caught before handoff; Decisions invalidated by stale inputs; automatic Decisions — target zero.

## 24. Classification de livraison
`defined` / `planned`; no final Decision object schema, signing protocol, storage, API, state machine or runtime enforcement implementation selected.

## 25. Critères d’acceptation
**Given** all required Approvals exist but the Decision Maker has not recorded a Decision, **When** the package is viewed, **Then** status remains undecided and no Execution Handoff Package is authorized.

**Given** a Decision is `approve-with-conditions`, **When** recorded, **Then** exact conditions, permitted/prohibited scope and expiry are stored as Decision constraints and future GOV-2 must receive them unchanged.

**Given** a Decision reaches expiry, **When** history is reviewed, **Then** the Decision remains resolvable but cannot be treated as current execution authority.

## 26. Questions ouvertes
OPEN-007, OPEN-013 and OPEN-015 remain open. Final Decision state machine, signing/quorum and execution-time enforcement are future Objects/Permissions/GOV-2 work.

## 27. Consommateurs documentaires
Decision Register, CAP-GOV-016, source-product return flows, future GOV-2 Runs/Verification/Rollback, GOV-3 Audit/Metrics, Objects/Permissions/Screens/Journeys/Technique and conformance report.
