---
id: CAP-STD-042
title: Automation Run Creation and Execution Context
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
# CAP-STD-042 — Automation Run Creation and Execution Context

## 1. Définition
Contrat fonctionnel de création d’une Automation Run Studio à partir d’un Workflow/version, Agent/version ou authorized caller context, avec initiator/caller, tenant/environment, inputs, Tool/Skill refs, permission/runtime refs, correlation, reason, return-origin, optional parent Run and provenance. Run created ≠ started.

## 2. Problème utilisateur
Sans Run identity distincte, une exécution peut être confondue avec son Workflow, un Tool Call, un Shared Job ou un Govern Response Run et perdre son contexte d’appel.

## 3. Objectifs
Créer une Run shell traçable sans effet ; pinner source versions/context ; vérifier tenant/input/permissions readiness ; préserver parent/correlation/return-origin and cross-product refs.

## 4. Non-objectifs
Ne démarre aucun effet, ne crée pas final object schema, scheduler/queue, Tool Call, Response Run, endpoint command, API/protocol or deployment.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Studio Operator principal; Automation Designer, source-product caller, human supervisor, Govern/Response Operator and Auditor secondary.

## 7. Conditions d’entrée
Authorized request/caller, Workflow or Agent exact version, tenant/environment, inputs/context and functional permission/runtime references available enough for preparation.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| source Workflow/version | STD-2 | exact orchestration definition | conditionnel | pinned | no Workflow-based Run |
| source Agent/version | CAP-STD-034 | exact Agent definition | conditionnel | pinned | no Agent-based Run |
| initiator/caller/reason | user/source product | identity + purpose | oui | request-time | no Run |
| tenant/environment/inputs | Settings + caller | scope + bound values/refs | oui | current/pinned | Run blocked |
| permission/runtime refs | Security + Settings/runtime | preparation/readiness projections | oui for execution | current | Run preparing/blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow / Automation Agent | CMDR Studio | exact source versions | read/reference |
| Tool / Skill | CMDR Studio | referenced dependencies | read/reference |
| Principal / Tenant / Environment | Platform Settings | caller/scope projections | read |
| Response Run / Decision | Govern | optional bounded caller authority refs | read/link |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Automation Run | create functional identity/shell | CMDR Studio | Run created ≠ started; physical object file deferred Phase 7 |
| Run Execution Context | bind/version temporary context | CMDR Studio | temporary context ≠ source object |

## 11. Fonctionnalités
Run identity; source version refs; caller/initiator/reason; tenant/env; inputs; Tool/Skill refs; permission/runtime refs; correlation/return-origin; parent Run; pre-start validation and provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect prepared Run | Studio Operator | Automation Run | 0 | read | Run shell/context | non |
| Validate creation context | Studio Reviewer | Run context | 1 | source refs + permissions | pass/block findings | non |
| Create Run shell | authorized caller/operator | Automation Run | 2 | create permission + bounded context | preparing Run | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
preparing, validation-pending, ready, blocked, superseded-before-start, cancelled-before-start. Later runtime states are owned by CAP-STD-043.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Automation Run identity/context | Automation Run | CAP-STD-043..051/source caller | pinned lineage + no-start guarantee |
| creation validation findings | assessment | Studio Operator | missing/stale/permission blockers explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Workflow/Agent/caller | request automation | CAP-STD-042 | source versions/context/return-origin | caller |
| CAP-STD-042 | Run ready | CAP-STD-043/045 | Run id + pinned context | source |
| Govern Response Run | authorized Studio handoff | CAP-STD-042 | Response Run/Decision/handoff refs | Govern |

## 18. Dépendances
CAP-STD-017..041; Tool/Skill STD-1; ownership register Automation Run Phase7; Settings/Security; Shared correlation/Jobs; Govern CAP-GOV-025 and OPEN-015; Endpoint future only by ref.

## 19. Source de vérité
Studio is source for Automation Run identity and functional creation context. Workflow/Agent/Tool/Skill stay source-owned; Govern Response Run remains distinct; Shared Job remains generic infrastructure.

## 20. Provenance et audit
Record Run id, source Workflow/Agent versions, caller/initiator, tenant/env, inputs refs, dependency versions, permission/runtime readiness refs, reason, parent Run, return-origin, correlation and creation actor/time. No raw secrets.

## 21. Permissions fonctionnelles
Automation Run read/create; source asset read; restricted Run context; caller/source permissions remain independent. Run creation never grants start/Tool invocation authority.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Missing source version, invalid tenant/env, denied caller, invalid input, unavailable runtime/Secret Reference, duplicate Run candidate or stale authority leaves Run preparing/blocked; no auto-start.

## 23. Métriques
Runs created vs started; blocked creation reasons; duplicate candidates; runs lacking exact source version—target zero; cross-product correlation completeness.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** an Automation Run is created successfully, **When** its state is inspected immediately, **Then** it is preparing/ready but no execution is considered started.

**Given** a Govern Response Run hands off to Studio, **When** an Automation Run is created, **Then** both Run identities remain distinct and correlated through explicit refs.

**Given** AI is unavailable, **When** a Run shell is created, **Then** structured request/context validation provides the complete preparation path.

## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Control Room, Workflows, Automation Agents, Agent Teams, Human Gates, CAP-STD-043..051, source products, Govern, Shared, Settings, Security, Quality.
