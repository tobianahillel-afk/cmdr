---
id: CAP-STD-043
title: Automation Run Lifecycle and State Management
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
# CAP-STD-043 — Automation Run Lifecycle and State Management

## 1. Définition
Définit le lifecycle fonctionnel d’une Automation Run Studio : preparing, validation-pending, ready, queued, scheduled, starting, running, waiting-human, paused, resume-pending, stop-requested, cancel-requested, partially-completed, completed, failed, cancelled, stopped, timed-out, blocked, superseded. Il ne constitue pas une machine d’état technique finale.

## 2. Problème utilisateur
L’intention locale ou la file d’attente peut être prise pour un état confirmé, et les états Automation Run peuvent être confondus avec ceux d’un Govern Response Run.

## 3. Objectifs
Distinguer request/confirmed states ; gérer waiting-human/pause/stop/cancel/terminal states ; conserver last confirmed runtime observation and unknown/stale conditions; assurer transition provenance.

## 4. Non-objectifs
Ne finalise pas persisted state machine, executor protocol, scheduler, API, Response Run lifecycle, verification/rollback or canonical Result.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Studio Operator principal; human supervisor, Automation Designer, source caller and Auditor secondary.

## 7. Conditions d’entrée
Automation Run exists; current version/context and last confirmed state known; transition permission and relevant readiness/authority context available.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Automation Run/current state | CAP-STD-042 | lifecycle subject | oui | last confirmed | transition blocked |
| runtime control confirmation | runtime owner/Studio | actual accepted/confirmed observation | conditionnel | source timestamp | requested/unknown remains |
| Human Gate state | CAP-STD-041 | waiting-human dependency | conditionnel | current | waiting remains |
| authority/readiness refs | Security/Govern/Settings | transition constraints | conditionnel | current for start/resume | blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Run | CMDR Studio | state/context/lineage | read/manage |
| Human Gate / Tool Call | CMDR Studio | dependent runtime refs | read/link |
| Response Run | Govern | distinct correlated state | read only |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Automation Run lifecycle projection | transition/version | CMDR Studio | request ≠ confirmation |
| Run lifecycle event | append attributed event | CMDR Studio/Shared trace consumed | history preserved |

## 11. Fonctionnalités
Functional transitions and guards; request vs confirmation; waiting-human; pause/resume/stop/cancel; terminal/partial/unknown treatment; supersession and history.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect Run lifecycle | Studio Operator | Automation Run | 0 | read | confirmed/requested states visible | non |
| Validate transition eligibility | Studio Operator | Automation Run | 1 | current context | allowed/blocked reason | non |
| Request lifecycle control | authorized operator | Automation Run | 2 | control permission + guards | request state/event | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
preparing, validation-pending, ready, queued, scheduled, starting, running, waiting-human, paused, resume-pending, stop-requested, cancel-requested, partially-completed, completed, failed, cancelled, stopped, timed-out, blocked, superseded, status-unknown.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Run lifecycle state | Automation Run projection | Control Room/source caller | confirmed/requested distinction + timestamp |
| Run lifecycle event | attributed event | Trace/Audit consumers | actor/source/reason/correlation |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-STD-042 | Run prepared | CAP-STD-043 | Run/context | source |
| CAP-STD-045/046/runtime | queue/control/confirmation | CAP-STD-043 | request/confirmation/source time | same Run |
| Human Gate | response/expiry | CAP-STD-043 | gate outcome/return step | same Run |

## 18. Dépendances
CAP-STD-042/044..051; Shared Jobs/Trace; Security/Settings; Govern OPEN-015 and Response Run non-equivalence; no final runtime implementation.

## 19. Source de vérité
Studio is source for Automation Run functional lifecycle and correlated control state; actual executor observations remain source-attributed; Govern remains source for Response Run state.

## 20. Provenance et audit
Record every lifecycle state/request/confirmation, prior state, actor/source, runtime ref, Human Gate/Tool Call refs, reason, timestamps, freshness, unknown/stale intervals and correlation.

## 21. Permissions fonctionnelles
Automation Run read/manage; lifecycle control requests; restricted runtime-state read. Effectful controls may require separate Govern/runtime authority and source permissions.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Invalid transition, lost acknowledgement, stale runtime, expired authority, blocked resume, unsupported pause, duplicate control, terminal conflict or source contradiction remains explicit; requested state never becomes confirmed silently.

## 23. Métriques
Request→confirmation latency; invalid transitions blocked; unknown-state duration; runs stuck waiting-human; terminal contradictions; false `created/queued=requested` to `running` mappings—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** a Run is queued, **When** its lifecycle is inspected, **Then** it is not represented as started or running.

**Given** a start request is rejected by the runtime, **When** state is reconciled, **Then** the Run remains non-running with the rejection reason and request provenance.

**Given** a Run is paused, **When** its state is shown, **Then** paused remains distinct from stopped and completed.

## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Control Room, CAP-STD-044..051, Workflows, Human Gates, source products, Govern, Shared, Quality.
