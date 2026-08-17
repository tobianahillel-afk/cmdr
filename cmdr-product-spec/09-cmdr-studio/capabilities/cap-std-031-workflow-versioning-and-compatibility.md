---
id: CAP-STD-031
title: Workflow Versioning and Compatibility
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
# CAP-STD-031 — Workflow Versioning and Compatibility
## 1. Définition
Contrat fonctionnel de Workflow Version : parent version, immutable reference concept, change summary, breaking/non-breaking concept, Tool/Skill/subworkflow dependency versions, consumer compatibility, migration requirement, deprecation, supersession and provenance.
## 2. Problème utilisateur
Sans versioning explicite, un consumer ou subworkflow peut suivre une modification breaking sans le savoir et rendre un run futur non reproductible.
## 3. Objectifs
Créer des versions candidates traçables, comparer compatibilité, préserver historiques/dependencies and expose migration/deprecation without promotion/deployment semantics.
## 4. Non-objectifs
Aucune promotion, deployment, canary, runtime migration, package format, semantic-version algorithm imposed or silent upgrade.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| current Workflow/version | CMDR Studio | parent snapshot | oui | pinned | no version candidate |
| candidate draft | Builder/Workflow | changed contract/graph | oui | candidate | no change to version |
| dependency/consumer refs | Studio + consumers | versions/compatibility context | oui | visible | compatibility unknown |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow | CMDR Studio | identity/lifecycle | read |
| Version | CMDR Studio | version/supersession | read |
| Tool/Skill/subworkflow refs | CMDR Studio | exact dependency versions | read |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Workflow Version metadata | create candidate/deprecate/supersede semantics | CMDR Studio | no deployment/promotion |
| Compatibility Assessment | derive | CMDR Studio | consumer/dependency view only |
## 11. Fonctionnalités
Parent/new version; change summary; breaking concept; dependency snapshot; compatibility; migration requirement; deprecation/supersession; reproducibility.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Compare versions | Automation Designer | Workflow Version | 0 | read | diff/refs | non |
| Assess compatibility | Studio Reviewer | Workflow Version | 1 | dependencies/consumers | assessment | non |
| Create candidate version | Automation Designer | Workflow | 2 | manage + valid draft | version candidate | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| draft change summary | oui | oui | oui | oui | deterministic diff + manual summary |
| compatibility assessment | oui | oui | oui | oui | dependency/contract diff |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
candidate, compatible, breaking-candidate, incompatible, deprecated, superseded, historical-reference.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Workflow Version reference | Version | Library/consumers/subworkflows | exact dependency snapshot |
| compatibility assessment | assessment | CAP-STD-030/032 | no deployability claim |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Workflow draft | create version candidate | Workflow Version | base/new contract + provenance | Workflow |
| consumer/subworkflow | check compatibility | Workflow Version | consumer/dependency refs | consumer |
| deprecated version | inspect replacement | superseding version | old/new refs | no auto migration |
## 18. Dépendances
CAP-STD-006/013/017/024/030/032; canonical Workflow/Version objects; Shared Versioning infrastructure; future STD-4 promotion/deployment.
## 19. Source de vérité
Studio Workflow/Version is source for version identity and dependency snapshots; consumer objects remain consumer-owned; Shared may provide generic versioning mechanism only.
## 20. Provenance et audit
Record parent/new version, actor, change summary, dependency versions, compatibility findings, migration flag, deprecation/supersession reason and correlation.
## 21. Permissions fonctionnelles
Workflow/version read/create/compare/deprecate/supersede functional needs; no deployment/promotion permission. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Unknown consumer, incompatible dependency, missing parent, stale dependency, denied ref or supersession gap remains explicit; no silent migration.
## 23. Métriques
Breaking candidates, incompatible consumers, deprecated usage, migration-required refs, exact-version reproducibility coverage.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** a Workflow version changes, **When** a consumer compares versions, **Then** changed graph/I-O/dependency refs and compatibility impact are visible.

**Given** a subworkflow dependency is incompatible with a new version, **When** compatibility is assessed, **Then** the candidate is incompatible for that consumer and no substitution occurs.

**Given** AI is disabled, **When** versions are compared, **Then** deterministic diff/compatibility checks and manual summary remain available.
## 26. Questions ouvertes
OPEN-013 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Workflow/Builder/Library, CAP-STD-024/030/032/033, future STD-4 Versions & Deployment, consumers, Quality.
