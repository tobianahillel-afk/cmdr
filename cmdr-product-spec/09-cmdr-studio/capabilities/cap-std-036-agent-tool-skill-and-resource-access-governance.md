---
id: CAP-STD-036
title: Agent Tool, Skill and Resource Access Governance
product: cmdr-studio
module: automation-agents
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-036 — Agent Tool, Skill and Resource Access Governance

## 1. Définition
Contrat fonctionnel des Tool/Skill/resource references accessibles à un Agent : exact versions, risk/eligibility, input/data restrictions, Secret Reference handling, tenant boundaries, runtime availability, denied/temporary availability and permission needs. Agent access ≠ ownership or permanent permission grant.

## 2. Problème utilisateur
Une liste de Tools sur un Agent peut être interprétée comme une autorisation d’invocation, ignorer les versions/risques ou exposer des ressources sensibles hors tenant.

## 3. Objectifs
Séparer configuration d’accès, permission effective et eligibility ; pinner versions ; contrôler scope/sensitive refs ; revalider availability/permission au moment pertinent ; tracer denied/temporary access.

## 4. Non-objectifs
Ne gère pas provider credentials, ne révèle pas secrets, ne modifie pas Tool/Skill ownership, ne crée pas permission durable, ne sélectionne pas runtime/provider et n’autorise aucun effet de classe 3/4.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Automation Designer and Studio Operator; Security Reviewer, Studio Reviewer and Auditor inspect access decisions.

## 7. Conditions d’entrée
Agent/version and bounded context resolved; Tool/Skill catalog references available; current permission/risk/tenant context queryable.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Agent/version + context | CAP-STD-034/035 | subject and bounded scope | oui | pinned/current | no access assessment |
| Tool/Skill refs | STD-1 | exact versions + risk/input contracts | oui when used | version current | denied/unavailable |
| permission/eligibility context | Security + Tool owner | functional access/eligibility | oui | action-time | denied |
| Secret/Runtime refs | Platform Settings | opaque availability metadata | conditionnel | current | restricted/unavailable |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Agent | CMDR Studio | configured allowed refs | read/manage |
| Tool / Skill | CMDR Studio | version/risk/I-O/eligibility metadata | read/reference |
| Secret Reference / Runtime Reference | Platform Settings | opaque metadata/availability only | restricted read |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Agent Access Policy projection | configure/validate | CMDR Studio | does not mutate Tool/Skill or Security permission |
| Access Decision record | derive/record conceptual | CMDR Studio + Security evidence | request-specific, revocable |

## 11. Fonctionnalités
Allow/deny exact Tool/Skill versions; risk/eligibility checks; input/resource restrictions; sensitive context masking; cross-tenant block; temporary availability; revocation awareness; per-call permission recheck.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect Agent access | Automation Designer | Access projection | 0 | read | allowed/denied refs visible | non |
| Validate Tool/Skill eligibility | Studio Reviewer | Agent + asset | 1 | current context | eligible/blocked reason | non |
| Configure allowed Tool/Skill refs | Automation Designer | Automation Agent | 2 | manage + no permission escalation | draft access changed | OPEN-013 |
| Invoke effectful Tool | none in this capability | Tool | 3 | separate runtime + authority | not executed here | Govern/source authority |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
unconfigured, allowed-reference, denied, permission-missing, risk-blocked, tenant-blocked, version-incompatible, runtime-unavailable, secret-reference-unavailable, temporarily-available, revoked.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Agent access contract | configured projection | CAP-STD-038/042/044 | exact refs + restrictions, no permission grant |
| Access/eligibility assessment | assessment | operator/Control Room | source-backed allow/deny reasons |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Agent plan | requests candidate Tool/Skill | CAP-STD-036 | agent/context/asset ref | plan |
| CAP-STD-036 | eligible candidate | CAP-STD-044 | exact asset/version + restrictions | Run |
| denied/revoked | escalate or choose allowed alternative | CAP-STD-039 | reason/asset/context | Agent |

## 18. Dépendances
CAP-STD-003..016/034/035/038/042/044; Security permission model; Settings providers/secrets/runtime health; Govern authority for effects; OPEN-013/015.

## 19. Source de vérité
Studio owns Agent access configuration and run-local access assessment semantics; Tool/Skill remain Studio source assets, Security owns permission policy, Settings owns secret/provider/runtime administration.

## 20. Provenance et audit
Record Agent/version, asset/version, requested operation, risk class, permission/eligibility sources, tenant/env, decision/reason, masked Secret Reference id, revocation/freshness and correlation. Never raw secret values.

## 21. Permissions fonctionnelles
Agent read/manage; Agent Tool access configure; Agent Skill access configure; Tool/Skill metadata read; restricted resource metadata; invocation permission remains separate and rechecked per call.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Unauthorized Tool, out-of-scope Tool input, revoked asset, unsupported version, permission denial, cross-tenant attempt, Secret Reference unavailable, runtime unavailable or stale eligibility remains blocked and explicit.

## 23. Métriques
Denied requests by reason; unauthorized requests blocked; revoked refs in active configurations; access rechecks; cross-tenant blocks; raw-secret exposure—target zero.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** an Agent requests a Tool not present in its allowed configuration, **When** access is evaluated, **Then** the Tool is denied and no Tool Call is prepared.

**Given** an allowed Tool is requested with an input outside the Agent permitted scope, **When** eligibility is evaluated, **Then** invocation is blocked without expanding scope.

**Given** a required Secret Reference is unavailable or unauthorized, **When** an Agent step is prepared, **Then** the step is blocked without revealing or substituting a secret.

**Given** a cross-tenant Tool/resource reference is requested without permission, **When** access is evaluated, **Then** access is denied and the source tenant remains protected.

## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Automation Agents, Agent Teams, CAP-STD-038/042/044/047/048, Security, Settings, Govern, Quality.
