---
id: CAP-STD-017
title: Workflow Definition and Functional Contract
product: cmdr-studio
module: workflows
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-SEC-001, REQ-AI-002]
open_decisions: [OPEN-013, OPEN-007, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-017 — Workflow Definition and Functional Contract
## 1. Définition
Contrat fonctionnel canonique d’un Workflow Studio : purpose, owner, consumers, trigger/reference concept, entrées/sorties, graph, steps, conditions, branches, dépendances, références Tool/Skill/Human Gate/subworkflow, risque, version, lifecycle, contraintes et provenance. Workflow ≠ Govern Playbook.
## 2. Problème utilisateur
Sans contrat Workflow unique, un auteur ou un consommateur peut confondre composition déterministe, procédure Govern, runtime Automation Run et simple séquence de Tools.
## 3. Objectifs
Définir la sémantique stable du Workflow avant runtime; rendre explicites structure, dépendances, versions, limites et handoffs; conserver une voie déterministe/no-AI.
## 4. Non-objectifs
Ne définit ni Automation Run, scheduler, exécution production, Agent runtime, Govern Playbook, final graph schema, langage d’orchestration, API ou protocole.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Workflow intent | Automation Designer | purpose/owner/consumer | oui | draft courant | création bloquée |
| dependency refs | Studio + Settings projections | Tool/Skill/subworkflow/Human Gate refs | oui | version/freshness visibles | Partial/incomplete |
| consumer context | Command/Investigate/Studio | source object + tenant/env + return-origin | conditionnel | courant | handoff non borné |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow | CMDR Studio | identity/version/lifecycle/tenant | read/write owner path |
| Tool / Skill | CMDR Studio | exact references + contract metadata | read/reference |
| Human Gate | CMDR Studio | reference + context requirements | read/reference |
| Playbook / Decision | Govern | distinct authority/procedure refs | read only |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Workflow functional definition | create/update draft semantics | CMDR Studio | pas de runtime ni de Playbook implicite |
| Workflow dependency projection | compose typed refs | CMDR Studio | références conservent owner/permission |
## 11. Fonctionnalités
Purpose/owner/consumers; conceptual trigger; graph and steps; I/O refs; Tool/Skill/Human Gate/subworkflow refs; conditions/branches; dependency/risk/version/lifecycle/constraints/provenance.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect Workflow contract | Automation Designer | Workflow | 0 | read permission | contract/provenance visible | non |
| Validate definition completeness | Studio Reviewer | Workflow | 1 | draft résolu | assessment déterministe | non |
| Create or edit Workflow draft | Automation Designer | Workflow | 2 | manage + tenant/env | versioned draft mutation | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| author structure | oui | oui | oui | oui | manual graph/config authoring |
| validate references | oui | oui | oui | oui | deterministic dependency checks |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
draft, incomplete, valid-for-review, incompatible, superseded-reference; these are STD-2 work states and do not redefine runtime states.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Workflow contract | Workflow definition | Builder/Library/consumers | exact owner/version/dependencies and non-goals |
| dependency/risk summary | assessment | CAP-STD-030 | no authorization or runtime-success claim |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Builder Session | save valid draft | Workflow | graph/I-O/dependencies/provenance | Builder return-origin |
| Workflow | reference dependency | Tool/Skill/Human Gate/subworkflow owner | exact ref/version/tenant | Workflow |
| Consumer product | open Workflow ref | Workflow surface | source context/return-origin | consumer |
## 18. Dépendances
CAP-STD-003..016; Workflow object; Shared Linking/Versioning; Security permission model; Settings refs; Govern boundaries; future STD-3 runtime.
## 19. Source de vérité
CMDR Studio owns Workflow semantics and the canonical Workflow object; referenced Tools/Skills/Human Gates/Settings/Govern objects remain source-owned.
## 20. Provenance et audit
Record creator/editor, tenant/env, Workflow identity/version, graph hash/reference concept, exact dependency versions, change reason, validation state, correlation and return-origin.
## 21. Permissions fonctionnelles
Workflow read/restricted-read/create/edit/version-reference and validation needs; `perm.studio.*` and `perm.cmdr-studio.*` remain historical/non-final. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Missing dependency, inaccessible source, invalid graph, stale ref, tenant mismatch, unsupported step type or permission denial remain explicit; no silent fallback.
## 23. Métriques
Definition completeness, unresolved dependencies, incompatible references, provenance completeness, no-AI authorability; no numerical SLO target.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** a Workflow references a missing Tool, **When** definition validation runs, **Then** the Workflow is incomplete and no substitute Tool is chosen.

**Given** a Workflow is valid statically, **When** a future runtime is unavailable, **Then** the definition remains defined while executable availability remains false/unknown.

**Given** AI is unavailable, **When** an Automation Designer authors a Workflow, **Then** manual graph/config authoring and deterministic validation remain fully usable.
## 26. Questions ouvertes
OPEN-013; OPEN-007; OPEN-015 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Builder, Workflows, Library, CAP-STD-018..033, future STD-3/4, Govern, Command, Investigate, Security, Quality, Roadmap.
