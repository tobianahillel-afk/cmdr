---
id: CAP-STD-049
title: Control Room Monitoring and Runtime Intervention
product: cmdr-studio
module: control-room
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-049 — Control Room Monitoring and Runtime Intervention

## 1. Définition
Définit la sémantique fonctionnelle du Studio Control Room pour observer active/queued/scheduled/waiting-human/error/stalled Runs, runtime availability, Agent/Workflow/current step/Tool Call/Human Gate refs, authorized context, control requests, intervention, escalation and pivots to source/provenance. Control Room ≠ Command Work Queue, Govern Runs & Rollback or Endpoint management console.

## 2. Problème utilisateur
Les opérateurs ont besoin d’une vue runtime cohérente sans que Studio devienne owner des Incidents, Response Runs, Endpoint fleet ou moteurs génériques.

## 3. Objectifs
Fournir monitoring source-backed; expose freshness/unknown states; enable authorized pause/resume/stop/cancel/intervention requests; preserve pivots/return-origin; route authority/administration to owners.

## 4. Non-objectifs
Ne réécrit pas STD-CTL-001, ne définit pas final columns/filters/buttons, ne possède pas deployment lifecycle STD-4, Command queue, Govern Response Run/Result, Endpoint fleet or generic observability engine.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Studio Operator/human supervisor principal; Automation Designer, source-product operator, Govern/Runtime Operator and Auditor secondary.

## 7. Conditions d’entrée
Permission to Control Room/read Runs, tenant/environment selected, runtime sources/Run projections available; restricted context separately authorized.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Automation Run projections | CAP-STD-043..048 | states/steps/errors/context | oui | source timestamps | Partial/unknown |
| Agent/Workflow refs | CMDR Studio | source configuration/definition refs | oui | pinned/current | limited pivot |
| runtime health/availability | Settings/runtime owner | health/freshness projection | oui | current/last seen | stalled/unknown |
| Human Gate/Tool Call refs | CMDR Studio | waiting/active dependencies | conditionnel | current | partial |
| Govern/source refs | source products | authority/return-origin | conditionnel | current | handoff limited |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Run / Agent / Team / Workflow | CMDR Studio | runtime/source refs | read/control request |
| Tool Call / Human Gate | CMDR Studio | active/waiting refs | read |
| Response Run / Decision | Govern | correlation/authority refs | read only |
| Endpoint/Fleet health | Endpoint/Settings | technical/admin projections only | read only |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Control Room Runtime View projection | derive/refresh | CMDR Studio | no new canonical business object |
| Runtime Intervention | prepare/request/record | CMDR Studio | does not modify Decision/source object |

## 11. Fonctionnalités
Active/queued/scheduled/waiting/error/stalled views conceptually; freshness/unknown; current step/Tool Calls/Human Gates; context masking; controls; interventions/escalation; pivots to source/provenance; no detailed UI spec.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect Run/refs/provenance | Studio Operator | Control Room | 0 | read | source-backed status | non |
| Normalize/detect stalled state | Studio Operator | Run status | 1 | timestamps/health | assessment | non |
| Request pause/resume/stop/cancel/intervention | authorized operator | Automation Run | 2 | control permission + rechecks | control/intervention request | OPEN-013 |
| Escalate authority need | Studio Operator | Govern handoff | 2 | effect/authority required | bounded handoff | Govern |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
healthy-view, loading-source, partial, stale, runtime-unavailable, status-unknown, attention-required, waiting-human, stalled-candidate, intervention-pending, escalation-pending.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Control Room runtime projection | view/assessment | Studio Operator | source refs + freshness + masking |
| Runtime Intervention/Control request | Studio event | CAP-STD-046/Run | actor/reason/bounds/correlation |
| pivot context | navigation handoff | source product/Govern | return-origin and permissions preserved |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| runtime sources | status update | CAP-STD-049 | Run/step/status/source time | source |
| CAP-STD-049 | control request | CAP-STD-046 | Run/control/current refs | Control Room |
| CAP-STD-049 | authority/admin pivot | Govern/Settings/source | bounded context/return-origin | Control Room |

## 18. Dépendances
CAP-STD-039/042..048/050/051; STD-CTL-001 existing screen; Settings Health; Shared Jobs/Trace/Notifications; Govern; Command/Investigate; Endpoint boundary; OPEN-007/013/015.

## 19. Source de vérité
Studio owns Control Room runtime supervision/intervention semantics. Each source owns its object/status; Settings owns runtime/admin health; Shared generic mechanisms; Govern authority; Endpoint technical/admin owners stay separate.

## 20. Provenance et audit
Record viewed source versions/freshness, Run/control/intervention ids, operator, reason, current state/context snapshot refs, pivots/escalations, masking decisions and correlation. Monitoring never rewrites source history.

## 21. Permissions fonctionnelles
Control Room read; Automation Run read/control/intervention; sensitive context read separately; Tool Call/Human Gate read; provenance export preparation; cross-tenant Run read separate. Final RBAC not defined.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Lost runtime status, stale health, permission denial, source contradiction, unsupported control, failed intervention, cross-tenant request or unavailable provenance leaves explicit Partial/unknown/blocked states.

## 23. Métriques
Active/queued/waiting/error/stalled counts conceptually; stale/unknown duration; interventions/escalations; control-denials; source-object mutations from Control Room—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** Control Room detects a stalled Run candidate, **When** an operator intervenes, **Then** the intervention is attributed and uses the Run control contract without modifying source business state directly.

**Given** runtime status is stale, **When** Control Room renders the Run, **Then** the last confirmed status and age are shown and success is not inferred.

**Given** AI is unavailable, **When** operators supervise Runs, **Then** source-backed status, deterministic stale/error rules and explicit human controls remain fully usable.

## 26. Questions ouvertes
OPEN-007; OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
STD-CTL-001, Automation Runs, Agents/Teams, Workflows, Human Gates, Tool Calls, Govern, Settings, Shared, Command/Investigate, Security, Quality.
