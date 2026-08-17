---
id: CAP-STD-024
title: Subworkflow Composition and Reusable Workflow References
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
# CAP-STD-024 — Subworkflow Composition and Reusable Workflow References
## 1. Définition
Contrat de référence d’un Workflow versionné comme step/dépendance d’un autre Workflow avec I/O, compatibility, recursion/cycle constraints, failure propagation, return mapping, permissions and provenance. Subworkflow reference ≠ copied Workflow.
## 2. Problème utilisateur
Copier silencieusement un subworkflow ou suivre automatiquement sa dernière version détruit reproductibilité et ownership.
## 3. Objectifs
Permettre composition réutilisable par référence exacte/contrainte explicite, valider recursion/compatibility and preserve separate lifecycle.
## 4. Non-objectifs
Aucune copie automatique, recursive runtime engine, deployment coupling, permission inheritance, final package format or scheduler.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| parent Workflow/version | CMDR Studio | composition context | oui | pinned | no composition |
| referenced Workflow/version | CMDR Studio | exact/constraint ref | oui | version visible | missing/incompatible |
| I/O contract | parent/subworkflow | binding contract | oui | versions pinned | mapping incomplete |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow parent | CMDR Studio | graph/context | read/write |
| Workflow referenced | CMDR Studio | version/I-O/lifecycle | read/reference |
| Version | CMDR Studio | compatibility/supersession | read |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Subworkflow Reference | create/update conceptual | CMDR Studio | reference only, no copy |
| compatibility/cycle assessment | derive | CMDR Studio | no runtime execution |
## 11. Fonctionnalités
Exact/constraint version ref; parent/subworkflow I-O mapping; dependency compatibility; recursion/cycle checks; failure/return semantics; permissions and provenance.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect subworkflow | Automation Designer | Subworkflow Reference | 0 | read | ref/version visible | non |
| Validate compatibility/recursion | Studio Reviewer | reference | 1 | both contracts accessible | assessment | non |
| Add/change subworkflow ref | Automation Designer | Workflow draft | 2 | manage | reference changed | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest reusable Workflow | oui | oui | oui | oui | manual Library selection |
| compatibility/cycle validation | oui | oui | oui | oui | deterministic graph/version checks |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
unresolved, compatible, incompatible, recursion-invalid, deprecated-reference, valid-for-review.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| subworkflow reference | Workflow dependency | parent graph | exact ref/version + mappings |
| compatibility assessment | assessment | CAP-STD-030 | no auto upgrade/copy |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Library/Workflow Detail | select subworkflow | parent Workflow | referenced id/version/return-origin | Builder |
| parent validation | inspect dependency | referenced Workflow | version/I-O/permissions | parent |
| future runtime | enter subworkflow | future Automation Run context | exact ref + correlation | STD-3 only |
## 18. Dépendances
CAP-STD-017/019/020/023/030/031; Shared Linking/Versioning; Security.
## 19. Source de vérité
Each Workflow keeps its own Studio source/lifecycle; Subworkflow Reference is part of parent Workflow configuration only.
## 20. Provenance et audit
Record parent/ref Workflow ids and versions, mapping refs, compatibility result, recursion check, actor/change reason and correlation.
## 21. Permissions fonctionnelles
Parent Workflow edit; referenced Workflow read/restricted-read; mapping validation. No execution permission inheritance. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Missing/deprecated/incompatible ref, invalid recursion/cycle, denied reference, I/O mismatch or stale version are explicit.
## 23. Métriques
Subworkflow reuse, incompatible/deprecated refs, recursion candidates, silent-copy violations, exact-version provenance.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** a referenced subworkflow version is incompatible, **When** parent validation runs, **Then** readiness is blocked and no alternate version is substituted.

**Given** a subworkflow is updated, **When** the parent still references the old version, **Then** the old exact reference remains reproducible and change is visible.

**Given** AI is disabled, **When** a subworkflow is selected, **Then** manual selection plus deterministic compatibility checks remain available.
## 26. Questions ouvertes
No new OPEN decision. Existing programme OPEN decisions remain unchanged.
## 27. Consommateurs documentaires
Workflow/Builder/Library, CAP-STD-017/030/031/033, Security, future STD-3 runtime, Quality.
