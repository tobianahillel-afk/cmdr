---
id: CAP-GOV-024
title: Response Step and Action Execution Coordination
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-008, REQ-PROD-009, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-024 — Response Step and Action Execution Coordination

## 1. Définition
Coordonner fonctionnellement les **Response Steps** et action intents d’un Response Run, leur ordre ou parallélisme autorisé, dépendances, target, executor owner, Workflow/Tool reference, input bindings, conditions, timeout, retry eligibility, compensation relation, verification checkpoint et rollback relevance sans documenter de commande réelle.

## 2. Problème utilisateur
Un Run multi-step peut produire des effets incohérents si des steps partent avant leurs dépendances, sur de mauvaises cibles ou malgré une condition devenue fausse. L’orchestration métier Govern doit rester distincte de l’orchestration technique Studio/Endpoint tout en conservant un ordre et des limites explicables.

## 3. Objectifs
- matérialiser les Response Steps depuis Playbook/Execution Plan ;
- distinguer ordered et parallel-eligible steps ;
- lier target/action intent/executor owner/input bindings ;
- porter preconditions, stop conditions, timeout/retry/compensation/rollback relevance ;
- créer des execution handoffs corrélables ;
- préserver step state sans transformer technical result en Result canonique.

## 4. Non-objectifs
Ne pas définir une commande offensive, endpoint/API/protocole, implémenter un orchestrator, modifier un Workflow, exécuter directement une primitive, auto-retry, confondre compensation/rollback, ou conclure au succès global depuis un step.

## 5. Propriétaire
Govern owns Response Step semantics and cross-step governance. Studio/Endpoint/provider owners retain technical step execution semantics. Playbook remains Govern-owned; Workflow remains Studio-owned.

## 6. Utilisateurs
Principal : Response Operator. Secondaires : Govern Reviewer, Studio/Endpoint Operator, Decision Maker as authority observer, Incident Commander/Investigator as consumers, Auditor.

## 7. Conditions d’entrée
Response Run in a state permitting execution coordination, current control authorization, pinned Execution Plan/Playbook, resolved targets, valid bindings and executor capability references.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Response Run + pinned Plan | CAP-GOV-022/019 | execution envelope | oui | current Run version | no step dispatch |
| step definitions/order/dependencies | Playbook/Execution Plan | functional orchestration | oui | pinned versions | blocked |
| exact target/action intent | Plan/Resolved Target | bounded effect | oui | current readiness context | blocked |
| parameter/Secret Reference bindings | CAP-GOV-019 | inputs by reference | selon step | current Plan | incomplete |
| executor owner/capability ref | Studio/Endpoint/provider owner | technical destination | oui for effectful step | current projection | executor-unavailable |
| timeout/retry/compensation/verification/rollback constraints | Plan/Playbook/Decision | safety controls | selon step | pinned/current | explicit unknown/block |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run / Response Step | Govern | current state/dependencies | read/manage local semantics |
| Execution Plan / Playbook | Govern | order, bindings, constraints | read |
| Decision/conditions | Govern | effect bounds | read |
| Workflow/Tool | CMDR Studio | executor implementation reference | read/link only |
| Endpoint/runtime capability | technical owner | supported action/current availability | restricted read |
| Secret Reference | Platform Settings | input reference only | restricted read, no value |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Response Step | create from pinned plan/transition/annotate | Govern | step success ≠ Run success |
| Execution Action intent | create/version/link | Govern local concept | exact target/scope/executor owner |
| Executor Handoff relation | prepare/link | Govern | technical object remains source-owned |
| Workflow/Tool/target | no local definition mutation | source owners | request only through CAP-GOV-025 |

## 11. Fonctionnalités
Instantiate ordered steps; compute dependency eligibility; mark parallel-eligible sets; validate step target/action/bindings/conditions; create executor handoff intent; track requested/accepted/running/terminal technical refs; apply timeout eligibility; flag retry/compensation relations; create verification checkpoints; stop downstream dependencies after failed required step; preserve per-target/per-step state.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect steps/dependencies | operator | Response Steps | 0 | read | execution graph/table | non |
| validate dependency/condition eligibility | reviewer | Step | 1 | current refs | eligible/blocked reason | non |
| annotate/reorder only if Plan permits new version | operator/reviewer | Plan/Step | 2 | no authority expansion | versioned change/re-review | OPEN-013 |
| dispatch approved effectful step | authorized operator/control path | Execution Action | 3/4 per underlying effect | Run/start authority + exact bounds | executor handoff | governed |
| issue raw command directly | none via Govern spec | technical target | — | forbidden | no action | executor owner only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| determine dependency eligibility | oui | oui | oui | explanation | dependency graph/rules |
| validate step inputs/conditions | oui | oui | oui | summarize blocker | checklist |
| propose safe ordering/parallelism | oui | explicit constraints | oui | suggestion | plan graph/table |
| summarize step failures | oui | structured outcomes | oui | sourced summary | status table |
| dispatch/authorize new effect | governed explicit path | contract-controlled | only with authority | no autonomous authority | operator + deterministic controls |

## 14. États fonctionnels
Response Step functional states may include `pending`, `blocked-dependency`, `eligible`, `dispatch-requested`, `accepted`, `running`, `succeeded-technical`, `failed-technical`, `timed-out`, `cancel-requested`, `cancelled`, `skipped`, `partial`, `status-unknown`, `verification-checkpoint`, `compensation-candidate`, `rollback-relevant`, `superseded`. These are not the final persisted technical machine.

## 15. États d’interface
Loading preserves active Run/step graph ; Empty means no instantiated steps ; Partial names unavailable technical states ; Error preserves last confirmed step state ; Offline forbids new effectful dispatch unless future guarantees allow ; Permission denied masks protected bindings/output ; Stale marks executor status age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Response Step state | Govern run projection | CAP-GOV-025..029 | exact Run/step/target/intent refs |
| Executor Handoff intent | bounded request context | CAP-GOV-025 | no target/procedure expansion |
| timeout/failure/partial condition | runtime event | CAP-GOV-027 | technical status not canonical Result |
| verification checkpoint | Run condition | CAP-GOV-028/029 | expected checkpoint/source retained |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-023 | start/control confirmed for execution | CAP-GOV-024 | Run, Plan, exact bounds | Runs & Rollback |
| CAP-GOV-024 | step eligible | CAP-GOV-025 | Run/step/target/action/bindings/conditions | step state |
| CAP-GOV-025/026 | technical status/output | CAP-GOV-024/027 | correlation + raw ref + normalized state | same Run/step |
| CAP-GOV-024 | verification checkpoint reached | CAP-GOV-028/029 | step/target/expected state refs | Runs & Rollback |

## 18. Dépendances
CAP-GOV-019..023/025..031, Response Step object, Studio Workflow/Tool/Automation Run, Endpoint/provider primitives, Settings Secret References/runtime health, Shared Trace/Jobs/Versioning, Security permissions, OPEN-008/013/015.

## 19. Source de vérité
Govern owns functional Response Step identity, dependency and governance state. Technical executor remains source for actual execution state/output. Studio Workflow steps are never silently converted into Response Steps; mappings are explicit and versioned.

## 20. Provenance et audit
Record Run/step ids, Plan/Playbook versions, dependency/order, exact target/action, executor owner/ref, binding refs, condition evaluation, dispatch actor/authority, handoff correlation, technical status refs, timeouts, retry/compensation eligibility, verification checkpoint, rollback relevance and timestamps.

## 21. Permissions fonctionnelles
Per-step inspect; execution-action dispatch through authorized Run; step annotate; restricted input/output metadata; technical output inspect; provenance export. Effectful dispatch inherits Run/Decision authority and can require step-up/SoD.

## 22. Limites et erreurs
Missing dependency, stale readiness, invalid binding, unavailable executor, lost status, timeout, contradictory technical states or scope mismatch blocks/marks the step; no dependent effect is assumed safe and no success is inferred from request acceptance.

## 23. Métriques
Steps blocked by dependencies; dispatch→confirmation latency; step timeouts; technical partial/failure rates; dependency violations prevented; steps marked successful without technical evidence — target zero.

## 24. Classification de livraison
`defined` / `planned`; no command, executor protocol, orchestration runtime or API selected.

## 25. Critères d’acceptation
**Given** a required prior step failed, **When** a dependent step becomes due, **Then** it remains blocked unless an explicit governed alternate path exists; no automatic continuation is fabricated.

**Given** one step reports technical success, **When** the Run is inspected, **Then** the step may be technically successful but Run and Result success remain undecided pending other steps and verification.

**Given** no AI, **When** steps are coordinated, **Then** deterministic dependency/condition tables and explicit operator controls provide full functionality.

## 26. Questions ouvertes
OPEN-008 covers actual executor support, OPEN-013 class-2 step/plan mutations and OPEN-015 Studio/Response Run bridge. Final mapping of Workflow steps to Response Steps remains future technical/object work.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-025..033, Studio/Endpoint execution owners, Command/Investigate projections, future Objects/Permissions/Screens/Technique, GOV-2 report and GOV-3 audit/metrics.
