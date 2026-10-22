---
id: CAP-STD-047
title: Runtime Error, Timeout, Retry and Partial Completion Handling
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
# CAP-STD-047 — Runtime Error, Timeout, Retry and Partial Completion Handling

## 1. Définition
Définit la matérialisation runtime fonctionnelle des erreurs/timeouts/retry attempts/partial completion sous les retry/idempotency semantics STD-2 : Tool/Skill/Workflow/Agent/provider/runtime/permission failures, lost status, partial outputs/branches, eligibility, exhausted attempts, manual intervention and terminal failure.

## 2. Problème utilisateur
Une erreur runtime peut déclencher un retry silencieux ou dupliquer un effet ; un Tool Call success peut masquer l’échec du Run global et une perte de statut peut être interprétée comme completion.

## 3. Objectifs
Classer runtime failures; conserver raw refs/known effect; appliquer STD-2 retry bounds et rechecks; distinguer attempt/Run; préserver partial completion; exiger intervention lorsque retry unsafe.

## 4. Non-objectifs
Ne redéfinit pas retry policy STD-2, ne promet pas exactly-once, n’implémente backoff/retry engine, ne crée pas Govern retry authority/rollback or canonical Result.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Studio Operator principal; human supervisor, Automation Designer, source caller, runtime operator and Auditor secondary.

## 7. Conditions d’entrée
Automation Run/step/attempt exists with failure/timeout/unknown/partial observation; retry contract and current permission/authority/context available enough to decide.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Run/step/attempt state | CAP-STD-043/044 | failure context | oui | last confirmed | disposition unknown |
| raw error/status/output refs | Tool Call/runtime/source | technical evidence | conditionnel | source timestamp | lost/unknown |
| STD-2 retry/idempotency contract | CAP-STD-027/028 | eligibility/bounds/compensation semantics | oui for retry | pinned | retry blocked |
| current permission/authority/context | Security/Govern/CAP-STD-048 | safe retry context | oui for effectful retry | fresh | retry blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Run / Step / Attempt | CMDR Studio | state/error/attempt history | read/manage disposition |
| Tool Call / Tool / Skill | CMDR Studio | raw technical failure/output refs | read/link |
| Response Run / Decision | Govern | distinct authority/correlation | read only |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Runtime Error disposition | create/classify/supersede | CMDR Studio | raw source preserved |
| Retry Attempt relation | prepare/record outcome | CMDR Studio | attempt ≠ Run; retry ≠ new Run |
| partial completion projection | record per step/branch | CMDR Studio | partial ≠ success |

## 11. Fonctionnalités
Error taxonomy; timeout/lost status; partial output/branch completion; retry eligibility and attempts; exhausted attempts; safe/unsafe retry; intervention; terminal failure; STD-2 compensation link; unknown status handling.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect failure/attempts | Studio Operator | Run Error | 0 | read | raw/normalized context | non |
| Evaluate retry eligibility | Studio Operator | Run Step | 1 | STD-2 bounds + current context | eligible/blocked | non |
| Prepare retry/manual intervention | authorized operator | Run Attempt | 2 | eligible + permission | new attempt or intervention | OPEN-013 |
| Effectful retry | authorized runtime path | Run Step | 3 | revalidated authority + idempotency safety | executor request | Govern when required |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
error-observed, timed-out, status-lost, partial-output, partially-completed, retry-eligible, retry-blocked, retry-requested, retrying, retry-exhausted, manual-intervention-required, terminal-failure, superseded.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Runtime Error/Retry disposition | Run projection | CAP-STD-043/049/050 | raw refs + attempt + affected scope |
| partial completion matrix | step/branch outcomes | consumer/human supervisor | success/failure/unknown remain distinct |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-STD-044/runtime | failure/timeout/unknown | CAP-STD-047 | Run/step/attempt/raw refs | Run |
| CAP-STD-047 | safe retry | CAP-STD-044/046 | new attempt + revalidated context | same Run |
| CAP-STD-047 | unsafe/exhausted | CAP-STD-039/049/050 | intervention/failure context | Run |

## 18. Dépendances
STD-2 CAP-STD-026..028; STD-1 Tool Call; CAP-STD-039/043/044/046/048..050; Shared idempotency/Jobs/Trace; Govern authority; Settings runtime health; OPEN-013/015.

## 19. Source de vérité
Technical source owns raw error/status. Studio owns their relationship to Automation Run and runtime attempt disposition. STD-2 remains source for retry/idempotency/compensation definition.

## 20. Provenance et audit
Record Run/step/attempt, raw error/status/output refs, known effect state, retry eligibility rule/version, attempt count, idempotency context, current permission/authority/context rechecks, intervention and timestamps/correlation.

## 21. Permissions fonctionnelles
Run/error read; retry prepare/request; restricted technical output; effectful retry permission/authority separately required. No silent retry permission.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Non-idempotent/unknown effect, expired authority, permission denial, runtime/provider unavailable, lost status, exhausted attempts, duplicate candidate or partial branch ambiguity blocks automatic retry and remains explicit.

## 23. Métriques
Errors/timeouts by class; retries per step; unsafe retries blocked; exhausted attempts; partial completion; status-lost duration; silent/infinite retries—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** a Tool Call times out, **When** runtime handling evaluates retry, **Then** the Tool Call/attempt remains recorded and retry occurs only if STD-2 eligibility and current authority permit it.

**Given** retry is not safe because effect state is unknown, **When** eligibility is evaluated, **Then** no automatic retry occurs and manual intervention is required.

**Given** some workflow branches complete and another fails, **When** the Run outcome is prepared, **Then** partial completion remains explicit and overall success is not fabricated.

**Given** a Tool Call succeeded but the overall Run later fails, **When** history is reviewed, **Then** the Tool Call success remains visible without making the Run successful.

**Given** runtime status is lost, **When** the Run is monitored, **Then** status becomes unknown/stale rather than silently successful.

## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Control Room, Automation Runs, Tool Calls, CAP-STD-049/050/051, Human Oversight, Govern, Shared, Settings, Quality.
