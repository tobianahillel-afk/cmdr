---
id: CAP-STD-037
title: Agent Team Composition, Roles and Coordination
product: cmdr-studio
module: agent-teams
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-STD-037 — Agent Team Composition, Roles and Coordination

## 1. Définition
Définit une Agent Team Studio : purpose, member Agent/version references, functional roles, coordinator/specialists, allowed interactions, shared-context bounds, handoffs, conflict handling, unavailable members and version compatibility. Agent Team ≠ human Team or Workflow.

## 2. Problème utilisateur
Sans contrat d’équipe, rôles d’agents, coordination et délégation peuvent être confondus avec permissions, graph Workflow ou équipes humaines et masquer les conflits entre propositions.

## 3. Objectifs
Rendre composition/roles explicites ; limiter shared context ; définir coordinator/handoffs/conflicts/unavailable member ; préserver version compatibility and human oversight without selecting a multi-agent framework.

## 4. Non-objectifs
Ne définit pas topology/runtime protocol, shared-memory engine, user Team, permission delegation, Workflow graph, model/provider, deployment or autonomous authority.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Automation Designer principal; Studio Operator/Reviewer and human supervisor secondary; Auditor consumes team provenance.

## 7. Conditions d’entrée
Agent Team or draft identifiable, member Agent/version refs resolvable, team purpose and tenant/environment known, permission context available.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Team purpose | Automation Designer | purpose/tasks/coordination intent | oui | draft current | team incomplete |
| member Agent refs | CMDR Studio | exact Agent/version + role | oui | compatibility current | member unavailable/incompatible |
| shared-context bounds | Team config + Security | allowed shared references/data | oui | current | coordination blocked |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Agent Team | CMDR Studio | identity/version/purpose/roles | read/manage |
| Automation Agent | CMDR Studio | member versions/capabilities | read/reference |
| Principal | Platform Settings | human owner/supervisor projection | read |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Agent Team definition | create/update draft/version refs | CMDR Studio | team ≠ human Team/Workflow |
| Agent role/handoff rules | configure functional roles | CMDR Studio | role ≠ authorization |

## 11. Fonctionnalités
Member/role composition; coordinator/specialist roles; allowed interactions; bounded shared context; task/handoff contracts; conflict handling; unavailable member behavior; compatibility and stop/escalation.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect team composition | Automation Designer | Agent Team | 0 | read | members/roles/context visible | non |
| Validate member/version compatibility | Studio Reviewer | Agent Team | 1 | member refs available | compatibility findings | non |
| Edit team roles/members | Automation Designer | Agent Team | 2 | manage | versioned draft mutation | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
draft, incomplete, compatible, member-unavailable, member-incompatible, conflict-detected, coordination-blocked, escalation-required, paused-reference, superseded.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Agent Team contract | Agent Team | CAP-STD-038/042/049 | exact members/roles/context bounds |
| coordination/conflict assessment | assessment | human supervisor | conflicts/unavailable members explicit |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Builder/Library | open Team | Agent Team | team/version/tenant | source |
| Agent Team | delegate bounded task | member Agent | task/context bounds/role | Team |
| conflicting proposals | escalate | CAP-STD-039 | proposal refs/reasons | Team |

## 18. Dépendances
CAP-STD-034..036/038/039/042; canonical Agent Team/Automation Agent objects; Security; Settings Principal; Shared communication primitives if any; OPEN-013.

## 19. Source de vérité
Studio owns Agent Team purpose/composition/coordination semantics. Member Agents remain distinct Studio objects; human identities remain Settings-owned; generic messaging infrastructure remains Shared/technical.

## 20. Provenance et audit
Record Team/version, purpose, member Agent versions, functional roles, coordinator, shared-context bounds, handoffs, conflicts, member availability, edits, actor and correlation.

## 21. Permissions fonctionnelles
Agent Team read/configure; member Agent metadata read; shared-context restricted read; no permission delegation through team membership or coordinator role.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Member unavailable/incompatible, circular delegation, conflicting proposals, missing coordinator where required, context-overexposure, cross-tenant member or permission denial remains explicit.

## 23. Métriques
Unavailable-member events; proposal conflicts; incompatible versions; escalations; context-bound violations blocked; role-to-permission conflations—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** an Agent Team member becomes unavailable, **When** the team prepares work, **Then** the member is marked unavailable and work is not silently reassigned outside configured roles/scope.

**Given** two members produce conflicting proposals, **When** coordination evaluates them, **Then** both proposals and provenance remain visible and the configured conflict/escalation path is used.

**Given** a coordinator delegates a task, **When** the member receives it, **Then** delegation transfers task context only and does not transfer the coordinator's permissions.

## 26. Questions ouvertes
OPEN-013 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Agent Teams, Automation Agents, CAP-STD-038/039/042/049/051, Security, Settings, Control Room, Quality.
