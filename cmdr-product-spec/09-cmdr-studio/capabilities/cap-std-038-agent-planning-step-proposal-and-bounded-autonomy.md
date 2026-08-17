---
id: CAP-STD-038
title: Agent Planning, Step Proposal and Bounded Autonomy
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
# CAP-STD-038 — Agent Planning, Step Proposal and Bounded Autonomy

## 1. Définition
Contrat fonctionnel de planification agentique bornée : proposer next step, Tool/Skill/Workflow candidate, dependency, missing input, human review, stop or escalation under objective, permissions, eligibility, Workflow constraints, runtime limits, policy context and Human Gates. Agent plan ≠ Workflow Definition or authorization.

## 2. Problème utilisateur
Une proposition agentique peut être prise à tort pour un plan canonique, une action autorisée ou une extension silencieuse du scope.

## 3. Objectifs
Permettre des propositions inspectables et bornées ; vérifier chaque candidate step contre context/permissions/eligibility ; préserver uncertainty/provenance ; arrêter/escalader lorsque contraintes ou information manquent.

## 4. Non-objectifs
Ne crée pas de Workflow Definition, ne modifie pas authority, ne persiste pas hidden chain-of-thought, ne choisit pas framework/model, ne démarre pas Tool Call/Run et ne bypass pas Human Gate.

## 5. Propriétaire
CMDR Studio Product Lead. Studio owns the functional semantics described here. Referenced objects, authority, administrative configuration and generic infrastructure remain owned by their canonical products.

## 6. Utilisateurs
Studio Operator and Automation Designer; human supervisor/Reviewer consumes proposals and escalation context.

## 7. Conditions d’entrée
Agent/version + bounded execution context available, allowed asset references resolved, permission/policy projections available; planning request is no-effect.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Agent/version + objective | CAP-STD-034/035 | planning subject/bounds | oui | pinned/current | no plan |
| allowed Tool/Skill/Workflow refs | CAP-STD-036 + STD-1/2 | candidate universe | oui | current eligibility | candidate excluded |
| current observations/context | source products + Run context | attributed facts/refs | conditionnel | freshness visible | uncertainty/missing input |
| Human Gate/Govern constraints | Studio/Govern | review/authority boundaries | conditionnel | current | proposal limited/escalated |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Automation Agent | CMDR Studio | goal/constraints | read |
| Workflow / Tool / Skill | CMDR Studio | candidate refs/contracts | read/reference |
| Policy / Decision | Govern | authority constraints only | read |

## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Agent Plan Proposal | create/replace conceptual no-effect proposal | CMDR Studio | not Workflow Definition or authority |
| Agent Step Proposal | create/review/reject | CMDR Studio | proposal ≠ action |

## 11. Fonctionnalités
Bounded decomposition; next-step candidates; Tool/Skill/Workflow candidate selection; missing-input detection; stop/escalation; uncertainty; deterministic constraint validation; proposal version/provenance.

## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Request bounded plan | Studio Operator | Agent | 1 | context valid | no-effect proposals | non |
| Review/reject proposal | human supervisor | Agent Step Proposal | 0 | read | disposition recorded | non |
| Accept proposal into Run preparation | authorized operator | proposal | 2 | still eligible/within scope | candidate passed to CAP-STD-044 | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| inspect/prepare capability context | oui | oui | oui | oui | manual forms + deterministic checks |
| explain constraints or failure | oui | oui | oui | oui | source-backed status/rule views |

AI is optional. Agentic assistance never grants permission, expands scope, creates Govern Approval/Decision, bypasses a Human Gate, reveals raw secrets, hides errors, retries indefinitely or mutates provenance. Essential operation remains possible through manual controls and deterministic rules.

## 14. États fonctionnels
planning-requested, planning, proposal-ready, missing-input, scope-conflict, permission-blocked, human-review-required, escalation-required, stop-proposed, rejected, superseded.

## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose source availability, freshness, unknown status and masking explicitly. UI intent never fabricates a runtime confirmation or authorization.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Agent Plan Proposal | proposal set | human supervisor/CAP-STD-044 | bounds + candidate refs + uncertainty |
| escalation/stop candidate | control signal | CAP-STD-039 | reason/context/provenance |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| CAP-STD-035/036 | request plan | CAP-STD-038 | objective/constraints/allowed assets | Agent |
| CAP-STD-038 | proposal accepted for preparation | CAP-STD-044 | proposal ref + recheck context | Run |
| proposal unsafe/uncertain | escalate | CAP-STD-039 | reason/proposal/context | Agent |

## 18. Dépendances
CAP-STD-034..037/039/042/044; STD-1 Tool/Skill and STD-2 Workflow contracts; Security; Govern boundaries; AI constraints; OPEN-013/015.

## 19. Source de vérité
Studio owns plan/proposal semantics and provenance. Source facts remain source-owned; Workflow Definition remains STD-2-owned; permissions/authority remain Security/Govern.

## 20. Provenance et audit
Record Agent/version, objective/context, candidate sources, constraints applied, candidate steps, rejected candidates/reasons, missing inputs, uncertainty, human dispositions, model/provider reference if used (not secret), timestamps and correlation. Do not require hidden chain-of-thought.

## 21. Permissions fonctionnelles
Agent planning request; proposal read/review; referenced-asset metadata; no invocation permission implied. Human review and effectful action permissions remain separate.

`perm.studio.*` and `perm.cmdr-studio.*` remain coexisting historical namespaces. STD-3 performs no bulk rename and defines no final RBAC/ABAC matrix.

## 22. Limites et erreurs
Untrusted content/tool output, missing input, permission denial, scope expansion, unavailable asset, stale context, exceeded iteration/time/resource bound or uncertain effect forces explicit limitation/stop/escalation.

## 23. Métriques
Plans blocked by constraints; scope-expansion proposals; human-review rate; stop/escalation candidates; proposals executed without recheck—target zero; deterministic/manual fallback coverage.

## 24. Classification de livraison
`defined / planned`. This documentary contract proves no runtime implementation, scheduler, agent framework, model/provider, API, protocol, physical JSON Schema, deployment or Endpoint capability.

## 25. Critères d’acceptation
**Given** an Agent proposes a next step that expands scope, **When** the proposal is validated, **Then** it is blocked or escalated and the Run scope is not changed.

**Given** AI proposes a Tool as the next step, **When** the proposal is accepted for preparation, **Then** Tool access and runtime permission are rechecked before any Tool Call.

**Given** AI is unavailable, **When** work must continue, **Then** the Workflow, operator-selected next step and deterministic constraint checks provide a non-AI path.

## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.

## 27. Consommateurs documentaires
Automation Agents, Agent Teams, CAP-STD-039/042/044/049, Human Gates, Security, Govern, Quality.
