---
id: CAP-STD-020
title: Workflow Steps, Nodes and Dependency Graph
product: cmdr-studio
module: workflows
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-SEC-001, REQ-AI-002]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-020 — Workflow Steps, Nodes and Dependency Graph
## 1. Définition
Sémantique fonctionnelle des nodes/steps et dependency edges d’un Workflow : step type, predecessor/successor, required/optional paths, disabled steps, unreachable candidates, cycle detection, conceptual start/terminal paths, annotations, grouping and ordering semantics.
## 2. Problème utilisateur
Sans graph contract, une représentation visuelle peut être prise pour l’ordre d’exécution, cacher des cycles ou transformer un node en objet Tool/Skill.
## 3. Objectifs
Définir un graph déterministe analysable et versionnable, indépendant du rendu UI et du storage/runtime scheduler.
## 4. Non-objectifs
Aucun graph storage schema, visual layout contract, runtime scheduler, final node execution object ou UI geometry.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Workflow draft | CMDR Studio | current graph | oui | current draft | no graph |
| step definitions | Studio | typed nodes + references | oui | draft/version pinned | incomplete |
| dependency edges | Studio | precedence/required semantics | oui | draft | unreachable/invalid |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow | CMDR Studio | graph/version | read/write owner path |
| Tool / Skill / Human Gate refs | CMDR Studio | referenced step payload metadata | read/reference |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Workflow graph projection | create/update | CMDR Studio | part of Workflow definition, not independent runtime object |
| Validation Assessment | derive conceptual | CMDR Studio | candidate findings only |
## 11. Fonctionnalités
Typed nodes; predecessors/successors; dependencies; required/optional paths; disabled nodes; start/terminal concepts; cycle/unreachable detection; annotations/grouping; execution-order semantics.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect graph | Automation Designer | Workflow | 0 | read | nodes/edges | non |
| Validate reachability/cycles | Studio Reviewer | Workflow graph | 1 | graph available | assessment | non |
| Add/remove/reorder step | Automation Designer | Workflow draft | 2 | manage | draft graph changed | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest graph structure | oui | oui | oui | oui | manual graph editing |
| cycle/reachability check | oui | oui | oui | oui | deterministic graph algorithms |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
empty, incomplete, structurally-valid, cycle-invalid, unreachable-candidate, disabled-path, superseded.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| normalized functional graph | Workflow projection | CAP-STD-022/025/030 | typed refs + edges, no storage schema |
| graph validation findings | assessment | Builder/Reviewer | exact node/edge refs |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Builder Session | edit graph | Workflow draft | nodes/edges/base version | Builder |
| Workflow graph | validate | Validation Assessment | graph/version | Workflow |
| node reference | open owner | Tool/Skill/Human Gate surface | exact ref/return-origin | Workflow |
## 18. Dépendances
CAP-STD-017/018/021/022/024/025/030; generic graph/layout primitives remain Shared/Design System; no runtime engine.
## 19. Source de vérité
Workflow owns graph semantics; referenced assets own their definitions; generic graph rendering/layout remains external to Studio business semantics.
## 20. Provenance et audit
Record node stable reference concept, type, referenced asset/version, edge changes, actor, base Workflow version, validation findings and correlation.
## 21. Permissions fonctionnelles
Workflow read/edit; graph validation; referenced metadata read. Layout/accessibility permissions remain screen/component concerns. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Cycle, missing endpoint, orphan node, unreachable candidate, forbidden reference, tenant mismatch or conflicting ordering are explicit.
## 23. Métriques
Invalid-cycle count, unreachable candidates, orphan refs, graph validation duration concept, graph-change provenance completeness.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** a graph contains an invalid cycle, **When** static graph validation runs, **Then** exact cycle nodes/edges are reported and readiness is blocked.

**Given** a node references a Tool, **When** the graph is inspected, **Then** the node remains a step definition and does not become the Tool object.

**Given** AI is disabled, **When** a graph is built, **Then** manual editing plus deterministic cycle/reachability checks remain available.
## 26. Questions ouvertes
No new OPEN decision. Existing programme OPEN decisions remain unchanged.
## 27. Consommateurs documentaires
Builder, Workflow Detail, CAP-STD-021..025/030/033, Shared graph primitives, Design System, Quality.
