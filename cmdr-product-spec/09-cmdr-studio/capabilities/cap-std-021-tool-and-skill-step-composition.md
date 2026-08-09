---
id: CAP-STD-021
title: Tool and Skill Step Composition
product: cmdr-studio
module: workflows
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-SEC-001, REQ-AI-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-STD-021 — Tool and Skill Step Composition
## 1. Définition
Contrat de composition de Workflow steps qui référencent un Tool ou une Skill avec version constraint, I/O mappings, dependencies and permission/risk needs. Tool step ≠ Tool Call; Skill step ≠ Skill.
## 2. Problème utilisateur
Une référence mal définie peut hériter silencieusement de permissions, appeler une version inattendue ou confondre la définition d’un step avec l’invocation runtime.
## 3. Objectifs
Composer Tools/Skills existants sans les redéfinir; pin/contraindre versions; valider mappings, dependency and eligibility needs; préserver permissions indépendantes.
## 4. Non-objectifs
Ne crée aucun Tool Call, Automation Run, Agent step runtime, permission grant ou Tool/Skill lifecycle.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Workflow step | Workflow draft | step type + reference | oui | draft | incomplete |
| Tool/Skill contract | CAP-STD-003/010 | definition/version/I-O | oui | pinned/constraint | missing/incompatible |
| eligibility metadata | CAP-STD-007 + Security | permission/risk needs | oui | current assessment basis | not executable |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Tool | CMDR Studio | definition/version/I-O/risk | read/reference |
| Skill | CMDR Studio | definition/version/deps/I-O | read/reference |
| Workflow | CMDR Studio | step context/mappings | read/write owner path |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Tool Step definition | compose in Workflow | CMDR Studio | not Tool Call |
| Skill Step definition | compose in Workflow | CMDR Studio | not Skill object mutation |
## 11. Fonctionnalités
Tool/Skill version constraints; input/output mapping refs; required dependencies; permission/risk metadata; deprecation/compatibility visibility; no permission inheritance.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect step dependency | Automation Designer | Tool/Skill step | 0 | read | contract visible | non |
| Validate step compatibility | Studio Reviewer | step | 1 | dependency accessible | assessment | non |
| Configure Tool/Skill step | Automation Designer | Workflow draft | 2 | Workflow manage | step definition changed | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest Tool/Skill step | oui | oui | oui | oui | manual Library selection |
| validate mapping/version | oui | oui | oui | oui | deterministic contract checks |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
unconfigured, configured, missing-dependency, deprecated-dependency, incompatible, permission-requirement-unresolved, valid-for-review.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Tool/Skill step definition | Workflow step | graph/validation | exact reference/version constraint |
| step dependency assessment | assessment | CAP-STD-030 | visibility ≠ authorization |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Library | select Tool/Skill | Workflow step | exact asset/version/return-origin | Builder |
| Workflow validation | check step | Tool/Skill contracts | mapping/version/permission needs | Workflow |
| future runtime | invoke Tool step | CAP-STD-008 | future run context + exact step ref | STD-3 only |
## 18. Dépendances
CAP-STD-003..014; CAP-STD-019/020/023/030; Security permissions; no Tool Call/Automation Run creation in STD-2.
## 19. Source de vérité
Tool/Skill definitions remain STD-1 sources; Workflow owns only the step reference/configuration. Runtime invocation later creates distinct runtime objects.
## 20. Provenance et audit
Record Workflow/version, step ref, Tool/Skill id/version constraint, mappings, dependency status, permission/risk requirements, actor and validation result.
## 21. Permissions fonctionnelles
Workflow edit; Tool/Skill metadata read; restricted metadata; configuration. Tool invoke/Skill dependent Tool permissions are not inherited. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Missing Tool, deprecated Skill, incompatible version, restricted dependency, invalid mapping or permission requirement unresolved blocks readiness explicitly.
## 23. Métriques
Missing/deprecated dependencies, incompatible version constraints, step-mapping failures, unauthorized-visible dependency cases.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** a Workflow references a missing Tool, **When** step validation runs, **Then** the exact Tool reference is marked missing and readiness is blocked.

**Given** a Skill version is deprecated, **When** the step is reviewed, **Then** deprecation is visible and no silent version substitution occurs.

**Given** a Tool step is visible but the user lacks invoke permission, **When** eligibility is inspected, **Then** configuration visibility remains distinct from authorization.
## 26. Questions ouvertes
OPEN-013 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
CAP-STD-017/020/023/030, Builder, Tool/Skill owners, Security, future STD-3 runtime, Quality.
