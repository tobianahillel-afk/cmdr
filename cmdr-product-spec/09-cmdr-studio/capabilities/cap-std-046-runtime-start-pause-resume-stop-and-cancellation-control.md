---
id: CAP-STD-046
title: Runtime Start, Pause, Resume, Stop and Cancellation Control
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
# CAP-STD-046 — Runtime Start, Pause, Resume, Stop and Cancellation Control

## 1. Définition
Définit les control intents et confirmations d’une Automation Run : start request/accepted/confirmed, pause request/confirmation, resume request with revalidation, stop request/confirmation, cancel request/cancellation, unsupported pause, too-late-to-cancel and human intervention. pause ≠ stop; stop ≠ cancel; cancel ≠ Govern rollback.

## 2. Problème utilisateur
Un bouton ou une request ne prouve pas l’état effectif de l’executor ; resume peut agir sur un contexte devenu obsolète et cancel après effet peut être confondu avec rollback.

## 3. Objectifs
Séparer intent/acceptance/confirmation; revalider start/resume; encadrer pause/stop/cancel; conserver unknown/partial effect; router authority to Govern when required.

## 4. Non-objectifs
Ne définit pas executor commands/API, scheduler, rollback, Decision modification, final state machine, provider/runtime, UI controls or destructive authority.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Studio Operator principal; human supervisor, source caller, Govern/Runtime Operator and Auditor secondary.

## 7. Conditions d’entrée
Automation Run existing in eligible state; current context/permission/authority and runtime capability available; control-specific safety preconditions known.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Automation Run/state | CAP-STD-043 | control subject | oui | last confirmed | control denied |
| current context/readiness | CAP-STD-035/036/048 + runtime | scope/permission/availability | oui for start/resume | fresh | blocked |
| control capability | runtime owner | pause/stop/cancel support | oui for requested control | current | unsupported |
| Govern authority ref | Govern | required authority for effectful action | conditionnel | valid/current | effect blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Run | CMDR Studio | state/control history | read/control |
| Tool Call / Run Step | CMDR Studio | active work/interruptibility refs | read |
| Decision / Response Run | Govern | authority/correlation refs | read only |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Run Control Request | create/record/confirm | CMDR Studio | intent ≠ confirmed executor state |
| Automation Run lifecycle | transition after attributable evidence | CMDR Studio | no fabricated confirmation |

## 11. Fonctionnalités
Start/pause/resume/stop/cancel requests; acceptance/confirmation; resume revalidation; unsupported/too-late conditions; active Tool Call handling; partial-effect context; human intervention and authority bridge.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Request start/pause/resume/stop/cancel | authorized Studio Operator | Automation Run | 2 | control permission + preconditions | control request | OPEN-013 |
| Validate start/resume context | Studio Operator | Run context | 1 | current refs | pass/block | non |
| Effectful start/resume | authorized path | Automation Run | 3 | required authority + runtime acceptance | executor request | Govern when required |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
start-requested, start-accepted, starting, running, pause-requested, paused, resume-requested, resume-blocked, stop-requested, stopping, stopped, cancel-requested, cancelled-before-effect, too-late-to-cancel, unsupported-control, control-unknown.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Control request/confirmation | Run lifecycle event | CAP-STD-043/049/source caller | intent and confirmation separately timestamped |
| partial-effect/authority review condition | runtime condition | CAP-STD-047/Govern | affected step/Tool Call/scope explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-STD-045 | start eligible | CAP-STD-046 | Run/current context/authority refs | Run |
| CAP-STD-046 | control request | runtime owner | Run/control/correlation/bounds | Control Room |
| runtime owner | accepted/confirmed/rejected | CAP-STD-043/046 | raw status/source time | same Run |

## 18. Dépendances
CAP-STD-039/042..045/047/048/049; Tool Call interruptibility; Security; Govern CAP-GOV-025/open decisions; Settings runtime availability; OPEN-013/015.

## 19. Source de vérité
Studio owns Automation Run control intent and its source-backed functional lifecycle; runtime owner is source of actual technical control response; Govern owns production authority/rollback.

## 20. Provenance et audit
Record Run/control request ids, prior state, actor, authority refs, context/readiness versions, active steps/Tool Calls, request/accept/confirm/reject timestamps, runtime source and correlation.

## 21. Permissions fonctionnelles
Run start/pause/resume/stop/cancel; state read; restricted runtime status; step-up/SoD may be required by Security/Govern. No control permission inherited from Agent role.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Unsupported pause, too-late cancel, lost acknowledgement, active non-interruptible Tool Call, stale context, expired authority, runtime unavailable, duplicate request or conflicting control leaves explicit pending/blocked/unknown state.

## 23. Métriques
Control request→confirmation latency; resume blocks after context drift; unsupported controls; stop during Tool Call outcomes; cancelled-before-effect; control-unknown duration.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** a start request is rejected, **When** the runtime response is received, **Then** the Run does not become running and the rejection remains attributable.

**Given** a Run is paused and its permission/context changes, **When** resume is requested, **Then** the relevant context is revalidated and resume can be blocked.

**Given** stop is requested during an active Tool Call, **When** the runtime cannot stop it immediately, **Then** `stop-requested` remains distinct from `stopped` and the active call remains visible.

**Given** cancel occurs before any side effect and is confirmed, **When** history is reviewed, **Then** the Run is cancelled-before-effect and no rollback is claimed.

## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Control Room, Automation Runs, Tool Calls, CAP-STD-043/047/049/051, Govern, Security, runtime/Settings, Quality.
