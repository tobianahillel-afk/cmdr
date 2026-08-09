---
id: CAP-GOV-022
title: Response Run Creation and Lifecycle Management
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-008, REQ-PROD-009, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-022 — Response Run Creation and Lifecycle Management

## 1. Définition
Créer et gérer le **canonical Response Run** Govern comme enveloppe versionnée d’une exécution gouvernée sous une Decision, reliant Handoff, Execution Plan, Playbook/version, targets/scope, owner, executor projections, conditions, runtime states, verification, rollback and future Result sans confondre Run créé, Run démarré et Automation Run.

## 2. Problème utilisateur
Sans objet Run canonique, une exécution peut être représentée uniquement par un Job, Tool Call ou Automation Run et perdre sa Decision, son scope, ses conditions ou sa vérification. À l’inverse, créer un Run ne doit pas déclencher automatiquement l’effet.

## 3. Objectifs
- créer un Response Run seulement après reconciliation autorisant la préparation ;
- pinner Decision/Handoff/Plan/Playbook versions et target set ;
- exposer owner, conditions, executor refs and timestamps ;
- définir un lifecycle fonctionnel couvrant préparation, runtime, verification, rollback/recovery and closure ;
- préserver relations vers future Result et technical runs sans fusion ;
- rendre chaque transition attribuée et auditée.

## 4. Non-objectifs
Ne pas finaliser la machine d’état physique, implémenter une queue/runtime, démarrer implicitement l’exécution, remplacer Automation Run/Job, produire Result avant verification, modifier Decision/Evidence ou choisir un executor/protocol global.

## 5. Propriétaire
Govern / Runs & Rollback owns Response Run semantics/lifecycle. Studio owns Automation Run/Workflow/Tool Call; Shared owns generic Job; Endpoint/provider owners execute technical effects. Govern later owns canonical Result.

## 6. Utilisateurs
Principal : Response Operator. Secondaires : Govern Reviewer, Decision Maker, Runtime/Endpoint Operator, Studio Operator, Incident Commander as consumer, Investigator as consumer, Auditor.

## 7. Conditions d’entrée
Authorization Reconciliation `authorized-for-run-preparation`, current Execution Plan and Playbook version, resolved targets/readiness context, tenant/environment and permission to create a Govern Response Run.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| authorization reconciliation | CAP-GOV-021 | run-preparation authority | oui | current versions | no Run |
| Decision/Handoff/Plan | Govern | lineage/bounds | oui | exact pinned versions | no Run |
| Playbook/version | CAP-GOV-018/019 | procedure | oui | exact compatible version | blocked |
| resolved targets/readiness | CAP-GOV-020 | current target context | oui | assessed | Run may remain preparing/block |
| owner/executor projections | Govern + Studio/Endpoint/Settings | responsibility/routing refs | oui for effectful run | current metadata | preparing/blocked |
| verification/rollback requirements | Decision/Plan | downstream constraints | selon action | exact pinned versions | explicit none/unknown/required |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Decision / Handoff / Execution Plan | Govern | exact authority/intention | read |
| Playbook | Govern | exact version/response semantics | read |
| Resolved Targets / Readiness | Govern + source owners | current execution context | read |
| Workflow/Automation Run/Tool Call | CMDR Studio | future executor/provenance refs | read/link only |
| Job | Shared | optional background-work projection | read/link only |
| Endpoint/runtime capability | source owner | technical executor projection | read |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Response Run | create/transition/version/supersede/close | Govern | Run creation ≠ start |
| Response Step relation | initialize/link | Govern | functional step refs, execution later |
| verification/rollback relations | initialize/link | Govern | requirements only until performed |
| Automation Run/Job/technical target | no source mutation | source owners | references do not transfer ownership |

## 11. Fonctionnalités
Create Run id/lineage; pin Decision/Plan/Playbook/targets/scope; assign owner; record executor candidates/projections; carry conditions/time bounds; initialize functional step/status relations; expose current state and timestamps; link future technical runs; preserve verification/rollback requirements; allow pause/stop/cancel/verification/rollback/recovery states through owner capabilities; close/supersede with history.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect Run/lineage | authorized user | Response Run | 0 | read | current state/context | non |
| validate lifecycle transition eligibility | operator/reviewer | Run | 1 | source state/rules | allowed/blocked reason | non |
| create Run / update preparation metadata | Response Operator | Run | 2 | authorized-for-run-preparation | Run in preparing | OPEN-013 |
| annotate/assign owner | operator/coordinator | Run | 2 | permission | versioned local mutation | OPEN-013 |
| start effectful Run | authorized operator/path | Run | 3/4 depending action | CAP-GOV-023 rechecks | separate action | governed |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| create Run shell from pinned inputs | oui | oui | oui | no need | structured form |
| validate lifecycle transition | oui | explicit rules | oui | explanation | state checklist |
| summarize Run context | oui | aggregation | oui | sourced summary | lineage/status tables |
| flag stale pinned dependency | oui | version checks | oui | explanation | deterministic validity checks |
| start/authorize/close success | accountable governed path | validation only | only when contract authorizes | no autonomous authority | explicit controls |

## 14. États fonctionnels
Functional lifecycle may include `preparing`, `readiness-check`, `ready`, `scheduled`, `starting`, `running`, `paused`, `stop-requested`, `cancel-requested`, `partially-completed`, `completed`, `failed`, `verification-pending`, `rollback-required`, `rolling-back`, `rolled-back`, `recovery-required`, `recovered`, `closed`, `superseded`. This is not a final persisted state machine.

## 15. États d’interface
Loading preserves active Run context ; Empty means no Run selected ; Partial exposes missing executor/status projections ; Error retains last valid state and correlation ; Offline cannot assert authoritative runtime transitions without confirmed writes ; Permission denied masks restricted target/output ; Stale marks last update and unknown runtime state.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Response Run | canonical execution-governance record | CAP-GOV-023..033/Command/Investigate | exact pinned lineage and current functional state |
| Run lifecycle event | attributed transition | Shared Trace/GOV-3 future | actor/source/time/reason preserved |
| Response Step relation set | functional execution structure | CAP-GOV-024 | Playbook/Plan lineage preserved |
| verification/rollback requirements | Run constraints | CAP-GOV-028..031 | inherited source and state explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-021 | authorized for preparation | CAP-GOV-022 | Decision/Plan/Playbook/readiness authority context | reconciliation |
| CAP-GOV-022 | Run created/preparing | CAP-GOV-023 | Run id, pinned bounds, conditions | Runs & Rollback |
| CAP-GOV-022/023 | started | CAP-GOV-024/025 | Run/steps/executor intents | same Run |
| runtime/verification/rollback capabilities | status transitions | CAP-GOV-022 | attributed state/outcome refs | Run history |
| CAP-GOV-032 | Result finalized | CAP-GOV-022 | Result ref/outcome/closure context | Result |

## 18. Dépendances
CAP-GOV-018..021/023..033, Response Run/Step/Rollback/Result objects, Studio Automation Run/Workflow/Tool, Endpoint/provider executors, Shared Jobs/Trace/Activity/Versioning/Notifications, Security permissions, OPEN-013/015.

## 19. Source de vérité
Govern Response Run is source of governed-run identity, pinned authority/procedure context and functional run state. It is not source of raw executor truth; source technical runs/results remain with Studio/Endpoint/provider and are reconciled by CAP-GOV-026.

## 20. Provenance et audit
Record Run id, tenant/env, Decision/Handoff/Plan/Playbook versions, target set, scope/conditions, owner, executor refs, creation actor/time, every transition/request/confirmation, step relations, technical run refs, verification/rollback/recovery relations, Result ref, errors, retries and correlation ids.

## 21. Permissions fonctionnelles
Response Run create/read/manage; assignment; lifecycle preparation; schedule/start/pause/resume/stop/cancel through relevant capabilities; restricted context; provenance export. Effectful start/rollback require separate higher-risk authority and possible step-up/SoD.

## 22. Limites et erreurs
Stale reconciliation/readiness, duplicate Run candidate, conflicting Run version, lost executor state, unavailable trace write, tenant mismatch or invalid transition yields explicit blocked/unknown state. A Run cannot claim started/success solely from local intent.

## 23. Métriques
Runs created vs started; time preparing→start; invalid transitions blocked; stale dependency blocks; orphaned technical executions; Run states unknown; Runs without Decision/Plan/Playbook pinned — target zero.

## 24. Classification de livraison
`defined` / `planned`; no final state-machine implementation, queue, protocol, database schema or execution engine selected.

## 25. Critères d’acceptation
**Given** an authorized reconciliation, **When** a Response Run is created, **Then** it enters preparation with exact Decision/Plan/Playbook/targets and does not start execution automatically.

**Given** a Studio Automation Run is later linked, **When** the Govern Run is inspected, **Then** both identities remain distinct and ownership is not merged.

**Given** no AI, **When** Run lifecycle is managed, **Then** deterministic transitions, status tables and explicit operator controls provide full functionality.

## 26. Questions ouvertes
OPEN-013 remains open for class-2 preparation/control defaults and OPEN-015 for Automation Run/Response Run bridge. Final persisted state machine and runtime contract remain future; no new OPEN.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-023..033, Command/Investigate projections, Studio/Endpoint/Shared execution sources, future Objects/Permissions/Screens/Journeys/Technique, GOV-2 report and GOV-3 audit/metrics.
