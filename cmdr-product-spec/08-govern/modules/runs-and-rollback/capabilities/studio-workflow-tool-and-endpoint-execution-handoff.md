---
id: CAP-GOV-025
title: Studio Workflow, Tool and Endpoint Execution Handoff
product: govern
module: runs-and-rollback
owner: Govern Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-004, REQ-PROD-006, REQ-PROD-008, REQ-PROD-009, REQ-PROD-015, REQ-PROD-016, REQ-PROD-020, REQ-AI-002, REQ-OBJ-007, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-008, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-GOV-025 — Studio Workflow, Tool and Endpoint Execution Handoff

## 1. Définition
Transmettre un execution request/handoff borné depuis un Govern Response Run/Response Step vers un owner technique autorisé — Studio Workflow/Tool, Endpoint Agent ou autre provider/runtime configuré — et recevoir des références de retour corrélables sans fusionner Response Run avec Automation Run, Tool Call, Job ou technical result.

## 2. Problème utilisateur
Une exécution cross-product peut perdre ses limites d’autorité lorsqu’elle traverse Govern vers Studio/Endpoint, ou au contraire importer les objets runtime techniques comme s’ils devenaient le Run canonique. Le handoff doit transmettre exactement l’intention autorisée et revenir au Run d’origine.

## 3. Objectifs
- transmettre Run ID, Decision/Plan lineage, exact target/action/scope/conditions ;
- transmettre parameter/Secret References sans raw secret ;
- indiquer executor owner, deadline/expiry et idempotency need conceptually ;
- recevoir accepted/rejected/technical run/call references ;
- préserver return origin/correlation and provenance ;
- garantir Workflow≠Playbook, Automation Run/Tool Call≠Response Run, technical output≠Result.

## 4. Non-objectifs
Ne pas définir API/protocole/command/signing format, administrer Workflow/Tool/Endpoint, résoudre raw secrets dans Govern, choisir provider/runtime global, créer un Automation Run soi-même, réécrire technical result ou contourner executor permissions.

## 5. Propriétaire
Govern owns the bounded execution handoff context. CMDR Studio owns Workflow/Tool/Tool Call/Automation Run; Endpoint owns Agent Command/technical execution; Settings owns integrations/secrets/runtime configuration; technical provider owns its own execution object/output.

## 6. Utilisateurs
Principal : Response Operator. Secondaires : Studio Operator, Endpoint/Runtime Operator, Govern Reviewer, Security Reviewer, Auditor.

## 7. Conditions d’entrée
Response Step eligible for dispatch, Response Run effectful start/control authorized, exact executor selected from declared Playbook/Plan dependencies, current target/readiness and required input/Secret References bound.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Run/Step identity and state | CAP-GOV-022/024 | execution subject | oui | current | no handoff |
| Decision/Plan/Playbook lineage | Govern | authority/procedure bounds | oui | pinned | no handoff |
| exact target/action/scope/conditions | Plan/Step | effect intent | oui | current approved context | blocked |
| parameter + Secret References | CAP-GOV-019 | input references | selon action | current Plan | blocked/incomplete |
| executor owner/capability/version | Studio/Endpoint/Settings/provider | technical destination | oui | current projection | unsupported/unavailable |
| deadline/expiry/idempotency need | Decision/Plan/executor contract | execution guardrail | selon action | current | explicit unknown/block |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Response Run / Step | Govern | context/action/target/state | read/handoff |
| Execution Plan / Decision / Playbook | Govern | bounds/version/conditions | read |
| Workflow / Tool / Tool Call / Automation Run | CMDR Studio | destination and returned refs | read/link, source-owned |
| Agent Command / Endpoint Agent | Endpoint Agent | capability/technical response refs | read/link, source-owned |
| Integration / Secret Reference | Platform Settings | configured destination/input refs | restricted read |
| Job | Shared | optional background correlation | read/link, never Run identity |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Executor Handoff | create/version/dispatch/cancel-if-not-accepted | Govern local concept | exact bounded intent |
| Response Step | attach handoff/technical refs | Govern | no raw output rewrite |
| Automation Run/Tool Call/Agent Command | no local creation semantics | source owner | created/executed by source contract |
| technical output ref | link/annotate | source owner + Govern projection | input to reconciliation only |

## 11. Fonctionnalités
Build bounded handoff; validate executor capability/version and target; include authority correlation and expiry; include parameter refs/Secret References; create idempotency/correlation need; route to source owner; capture accepted/rejected/technical run/call/command id; preserve source timestamps; handle duplicate/late/expired handoff; route raw references to runtime reconciliation without interpreting business success.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| inspect executor capability/handoff | operator | Handoff | 0 | read | context | non |
| validate handoff bounds/idempotency inputs | reviewer | Handoff | 1 | exact refs | pass/block list | non |
| prepare/cancel unaccepted handoff | operator | Handoff | 2 | Run/step state permits | versioned intent | OPEN-013 |
| dispatch approved effectful handoff | authorized operator/path | Handoff | 3/4 per action | Run authority/current constraints | source-owner execution request | governed |
| bypass source owner/permission | none | executor | — | forbidden | no effect | source owner only |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| construct bounded handoff | oui | structured mapping | oui | draft only | template/mapping table |
| validate target/scope/version | oui | oui | oui | explain | exact diff/checklist |
| correlate returned runtime refs | oui | correlation ids | oui | explanation | reference table |
| flag rejected/late/duplicate handoff | oui | explicit status/time rules | oui | summary | status rules |
| authorize/execute secretly | accountable governed/source path | no autonomous | only explicit approved policy | prohibited | human/deterministic path |

## 14. États fonctionnels
`not-prepared`, `prepared`, `validated`, `dispatch-requested`, `accepted`, `rejected`, `expired-before-acceptance`, `duplicate-candidate`, `technical-run-linked`, `technical-call-linked`, `status-unknown`, `cancel-requested`, `cancelled-before-acceptance`, `superseded`.

## 15. États d’interface
Loading preserves handoff/run refs ; Empty no handoff ; Partial names source refs not yet returned ; Error preserves last confirmed source response ; Offline prevents new dispatch unless guaranteed ; Permission denied masks sensitive executor/input metadata ; Stale marks executor status age.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Executor Handoff | bounded request context | Studio/Endpoint/provider | exact Run/step/target/action/scope/refs/expiry |
| accepted/rejected status | technical response ref | CAP-GOV-026/023 | source owner + correlation preserved |
| Automation Run/Tool Call/Agent Command ref | provenance relation | CAP-GOV-026/033 | identity remains source-owned and distinct from Run |
| handoff error/timeout | execution condition | CAP-GOV-027 | no technical/business success inferred |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-GOV-024 | step eligible | CAP-GOV-025 | Run/step/target/action/bindings/conditions | step |
| CAP-GOV-025 | dispatch | Studio/Endpoint/provider | bounded handoff + correlation | same Run/step |
| technical owner | accept/reject/status/output | CAP-GOV-026 | raw ref/state/output metadata/correlation | executor object |
| handoff timeout/failure | no valid response | CAP-GOV-027 | request/ref/expiry/known state | same Run |

## 18. Dépendances
CAP-GOV-019..024/026/027/033, Studio Workflow/Tool/Automation Run/Human Gate, Endpoint Agent/Agent Command, Settings Integration/Secret Reference/health, Shared Jobs/Trace, Security permissions, OPEN-007/008/013/015.

## 19. Source de vérité
Govern is source of authorized handoff intent; the technical owner is source of accepted/executed/raw output facts. The mapping/correlation links the sources but transfers neither ownership nor permission.

## 20. Provenance et audit
Record Run/step/Decision/Plan/Playbook versions, handoff id/version, executor owner/capability/version, exact target/action/scope, input/Secret Reference ids only, conditions/expiry/idempotency need, dispatch actor/authority, correlation id, source acceptance/rejection/ref, timestamps, errors and automation provenance.

## 21. Permissions fonctionnelles
Executor handoff prepare/dispatch/read/cancel-if-safe; Workflow/Tool/Endpoint capability read; restricted parameter metadata; technical output reference read; provenance export. Source-owner execution and secret-use permissions remain independent.

## 22. Limites et erreurs
Missing executor, unsupported target/action, stale capability, invalid Secret Reference, expired Decision/Run window, duplicate idempotency candidate, lost acknowledgement, source permission denial or late output remains explicit and cannot create a success Result.

## 23. Métriques
Handoff acceptance/rejection/timeout; orphaned source executions; duplicate candidates prevented; late responses; source refs correlated; handoffs with raw secrets — target zero; technical objects merged into Response Run identity — target zero.

## 24. Classification de livraison
`defined` / `planned`; no API, protocol, command, provider/runtime or signing implementation selected.

## 25. Critères d’acceptation
**Given** a Studio Automation Run executes the technical procedure, **When** it is linked back, **Then** it remains Studio-owned and distinct from the Govern Response Run while its status/output refs are correlated.

**Given** an Endpoint technical action returns a success status, **When** the handoff completes, **Then** that status is only an input to reconciliation and cannot directly become the canonical Result.

**Given** no AI, **When** execution is handed off, **Then** templates, deterministic mapping/validation and explicit source-owner routing provide full functionality.

## 26. Questions ouvertes
OPEN-007 remains for Human Gate relation, OPEN-008 actual executor support, OPEN-013 control defaults, OPEN-015 final Automation Run/Response Run bridge. No new OPEN.

## 27. Consommateurs documentaires
Runs & Rollback, CAP-GOV-026/027/033, Studio/Endpoint/Settings contracts, future Objects/Permissions/Screens/Journeys/Technique, GOV-2 report and GOV-3 audit/metrics.
