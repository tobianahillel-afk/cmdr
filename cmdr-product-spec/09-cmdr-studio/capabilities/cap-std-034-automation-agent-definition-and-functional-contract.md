---
id: CAP-STD-034
title: Automation Agent Definition and Functional Contract
product: cmdr-studio
module: automation-agents
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-AI-001, REQ-AI-002, REQ-AI-004, REQ-AI-007, REQ-OBJ-009, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-STD-034 — Automation Agent Definition and Functional Contract

## 1. Définition
Contrat fonctionnel canonique d’un Automation Agent Studio : identity, owner, intended purpose/users, supported tasks, allowed environments, Tool/Skill/Workflow references, constraints, AI dependency, version/lifecycle reference, risk, permission needs and provenance. Automation Agent ≠ Automation Run, Workflow, Skill, Tool or human user.

## 2. Problème utilisateur
Sans définition stable, un agent peut être confondu avec son exécution, un Workflow ou une identité humaine, et ses capacités peuvent être interprétées comme une autorisation implicite.

## 3. Objectifs
Définir une identité/version et un périmètre fonctionnel réutilisable ; rendre explicites Tools/Skills/Workflows autorisables, limites, environnements et dépendances ; préserver une lecture et une gestion déterministes même si l’IA est indisponible.

## 4. Non-objectifs
Ne définit ni Automation Run, framework agent, model/provider, reasoning storage, scheduler, deployment/promotion, final object schema, API/protocol, permission grant or Govern authority.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Automation Designer principal ; Studio Reviewer, Studio Operator, Security Reviewer and Auditor as secondary users; operational products consume references under permission.

## 7. Conditions d’entrée
Tenant/environment known, actor authenticated, Agent draft/version identifiable or creation authorized, referenced Studio assets resolvable and functional permission context available.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Agent intent | Automation Designer | purpose/users/supported-task declaration | oui | draft current | creation incomplete |
| Tool/Skill/Workflow refs | Studio catalogs | typed exact-version references | conditionnel | version/freshness visible | capability unavailable |
| environment/risk context | Settings + Studio | allowed environment/risk projection | oui | current metadata | Agent not eligible for review |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Agent | CMDR Studio | identity/version/lifecycle/purpose | read/manage |
| Tool / Skill / Workflow | CMDR Studio | exact references + eligibility metadata | read/reference |
| Environment / Principal | Platform Settings | environment and owner/actor projection | read only |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Automation Agent functional definition | create/update draft/version metadata | CMDR Studio | definition ≠ execution |
| Agent dependency links | link/unlink typed references | CMDR Studio | links transfer neither ownership nor permission |

## 11. Fonctionnalités
Identity and immutable-reference concept; purpose/users/tasks; allowed environments; Tool/Skill/Workflow refs; constraints; AI dependency; non-AI boundary; version/lifecycle reference; risk/permission needs; provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Inspect Agent definition | Automation Designer | Automation Agent | 0 | read permission | definition + refs visible | non |
| Validate Agent completeness | Studio Reviewer | Automation Agent | 1 | references available | no-effect assessment | non |
| Create/update Agent draft | Automation Designer | Automation Agent | 2 | manage + tenant/env | versioned draft mutation | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
draft-definition, incomplete, valid-for-review, restricted, incompatible-reference, superseded-reference, retired-reference. These are capability work states and do not replace the canonical object lifecycle.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Agent functional contract | Automation Agent | Builder/Library/Control Room | owner/version/purpose/limits explicit |
| Agent dependency/risk summary | assessment | CAP-STD-035/036/030-style consumers | no authorization or execution claim |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Library/Builder | open Agent | Automation Agent | agent ref/version/tenant/return-origin | source surface |
| Automation Agent | reference Tool/Skill/Workflow | source asset | exact ref/version | Automation Agent |
| Automation Agent | future run request | CAP-STD-042 | agent/version + bounded context | Agent detail |

## 18. Dépendances
STD-1 Tool/Skill contracts `CAP-STD-003..016`; STD-2 Workflow contracts `CAP-STD-017..033`; Automation Agent canonical object; Settings environment/identity projections; Security; OPEN-013; future STD-4 assurance/deployment.

## 19. Source de vérité
CMDR Studio is source for Automation Agent semantics and identity. Tools, Skills, Workflows remain their own Studio objects; Settings remains source for environment/identity/provider metadata.

## 20. Provenance et audit
Record Agent id/version, owner, tenant/env scope, purpose/task changes, exact Tool/Skill/Workflow refs, constraint/risk changes, editor/reviewer, reason and correlation. No hidden reasoning or raw secret is required for provenance.

## 21. Permissions fonctionnelles
Functional needs: Agent read; Agent create/update; restricted config read; dependency-reference read; submit-for-review. Tool/Skill invocation permissions remain independent from Agent configuration.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Missing/incompatible reference, disallowed environment, cross-tenant link, stale dependency, missing owner, permission denial or unsupported AI dependency remains explicit; no substitute asset/provider is selected silently.

## 23. Métriques
Definition completeness; unresolved dependency count; restricted-reference count; version reproducibility; no-AI inspect/manage availability.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** an Automation Agent definition references a Workflow, **When** it is inspected, **Then** the Agent and Workflow remain distinct objects with separate versions and ownership.

**Given** an Agent definition is complete but no runtime exists, **When** it is reviewed, **Then** it can be `defined / planned` without claiming an executable Automation Run.

**Given** AI/model service is unavailable, **When** an authorized user inspects or edits Agent configuration, **Then** forms, reference selectors and deterministic validation remain usable.

## 26. Questions ouvertes
OPEN-013 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Automation Agents, Builder, Library, Agent Teams, Control Room, CAP-STD-035..051, Security, Settings, Quality, Roadmap, future STD-4.
