---
id: CAP-STD-044
title: Automation Run Steps, Attempts and Tool Call Coordination
product: cmdr-studio
module: control-room
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-044 — Automation Run Steps, Attempts and Tool Call Coordination

## 1. Définition
Définit la structure runtime fonctionnelle des Run Steps/Attempts : Workflow/Agent/Tool/Skill step references, exact input bindings, output refs, state/timestamps/errors, retry relation, parent Run, ordering/dependencies and Human Gate relation. Tool Call reste un objet distinct de Automation Run and attempt.

## 2. Problème utilisateur
Sans distinction Run/step/attempt/Tool Call, un retry peut être pris pour une nouvelle Run ou un Tool Call successful pour un Run successful.

## 3. Objectifs
Corréler chaque runtime step et attempt à la définition source ; créer/observer Tool Calls par contrat STD-1 ; préserver ordering/dependency/Human Gate and output provenance; distinguer retries.

## 4. Non-objectifs
Ne redéfinit pas Workflow step semantics, Tool Call lifecycle, retry policy STD-2, scheduler, executor payload schema or final Run Step object.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Studio Operator; Automation Designer/Agent supervisor; source caller and Auditor.

## 7. Conditions d’entrée
Automation Run prepared/active; source Workflow/Agent plan available; exact Tool/Skill refs and current access/permission context for actionable steps.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Run + source version | CAP-STD-042/043 | parent execution context | oui | pinned/current | step not attributable |
| source step/proposal | STD-2 or CAP-STD-038 | definition/proposal ref | oui per step | exact version | step blocked |
| input bindings/context | CAP-STD-048 | run/step-local values/refs | oui | current | step incomplete |
| Tool/Skill access | CAP-STD-036 + STD-1 | eligibility/permission refs | conditionnel | action-time | step denied |
| Human Gate relation | CAP-STD-040/041 | wait/review ref | conditionnel | current | step waiting/blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Run | CMDR Studio | parent/state/context | read |
| Workflow / Agent Plan | CMDR Studio | source step definition/proposal | read |
| Tool Call / Tool / Skill | CMDR Studio | exact execution/dependency refs | read/link |
| Human Gate | CMDR Studio | step gate relation | read/link |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Run Step projection | create/update runtime projection | CMDR Studio | step ≠ source definition object |
| Run Attempt | create/link outcome concept | CMDR Studio | attempt ≠ Run |
| Tool Call | request via STD-1 contract | CMDR Studio | distinct identity/lifecycle |

## 11. Fonctionnalités
Instantiate step projections; bind inputs; preserve source definition/proposal; track attempts; coordinate Tool Calls; states/timestamps/errors/output refs; ordering/dependencies; Human Gate wait; retry linkage.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect step/attempt | Studio Operator | Run Step/Attempt | 0 | read | source/attempt/outputs visible | non |
| Validate step readiness/access | Studio Operator | Run Step | 1 | bindings/access/current context | pass/block | non |
| Prepare next attempt | authorized runtime/operator | Run Attempt | 2 | retry/permission eligible | attempt candidate | OPEN-013 |
| Request effectful Tool Call | authorized Run path | Tool Call | 3 | separate eligibility + authority | Tool Call request | Govern/runtime when required |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
pending, ready, queued, running, waiting-human, succeeded, failed, timed-out, skipped, cancelled, retry-eligible, retrying, partial-output, status-unknown, superseded.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Run Step/Attempt status | runtime projection | CAP-STD-043/047/049/050 | source step + attempt + raw refs preserved |
| Tool Call reference/output ref | Studio source object/ref | Run/consumer | Tool Call distinct; output unqualified |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Run | instantiate source step | CAP-STD-044 | source ref/bindings/context | Run |
| CAP-STD-044 | Tool step eligible | Tool Call contract | Run/step/attempt/tool/input refs | Run Step |
| step requires review | request/wait | CAP-STD-040/041 | Run/step/reason/context | Run Step |

## 18. Dépendances
STD-1 CAP-STD-003..009; STD-2 CAP-STD-017..029; CAP-STD-036/038/040..043/047/048; Shared Trace; Govern OPEN-015; Security.

## 19. Source de vérité
Studio owns Run Step/Attempt coordination semantics and Tool Call identity/lifecycle under STD-1. Workflow definition/proposal remains its source; source outputs remain source-attributed.

## 20. Provenance et audit
Record Run/step/attempt ids, source definition/proposal/version, bindings refs, Tool/Skill versions, Tool Call ids, Human Gate refs, states/timestamps/errors/retry relation/output refs and correlation.

## 21. Permissions fonctionnelles
Run/step/attempt read; prepare attempt; Tool Call request/read according to Tool/access permissions; Human Gate refs read. No permission inheritance from Agent/Run.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Missing source step, invalid binding, denied Tool, unavailable Skill/runtime, lost Tool Call status, duplicate attempt, ambiguous retry, Human Gate expiry or out-of-order result remains explicit.

## 23. Métriques
Attempts per step; duplicate-attempt prevention; Tool Call correlation; step status unknown; Tool Call success with failed Run; missing source-version links—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** a Tool Call succeeds but a later required step fails, **When** the Automation Run is summarized, **Then** Tool Call success remains local and overall Run success is not inferred.

**Given** a retry creates a second attempt for the same step, **When** history is reviewed, **Then** both attempts remain under the same Automation Run and the retry does not create a duplicate Run.

**Given** a Human Gate is attached to a step, **When** execution reaches it, **Then** the step waits on the distinct Human Gate lifecycle rather than fabricating an approval.

## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Control Room, Automation Runs, Tool Calls, Workflows, Human Gates, CAP-STD-045..051, Govern, Shared Trace, Quality.
