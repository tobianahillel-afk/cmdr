---
id: CAP-STD-018
title: Workflow Builder Session and Editing Context
product: cmdr-studio
module: builder
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-SEC-001, REQ-AI-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-STD-018 — Workflow Builder Session and Editing Context
## 1. Définition
Sémantique fonctionnelle d’une session d’édition Builder : Workflow draft ciblé, owner/collaborators, base version, changements non sauvegardés, validation, selection/viewport context conceptuel, dirty state, conflits, recovery, save/discard et return-origin. Builder Session ≠ Workflow.
## 2. Problème utilisateur
Sans contexte d’édition explicite, des modifications concurrentes ou non sauvegardées peuvent être confondues avec une nouvelle Workflow Version ou écraser silencieusement une base.
## 3. Objectifs
Définir une édition collaborative sûre, détecter conflits, conserver base/version et provenance, permettre save/discard/recovery sans imposer l’UI.
## 4. Non-objectifs
Aucun wireframe, bouton/shortcut final, storage format, collaborative protocol, runtime execution, publication ou deployment.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Workflow/base version | CMDR Studio | draft + exact base version | oui | pinned au début de session | session non créée |
| editor/collaborators | Platform Settings principals + Studio | actors/permissions | oui | courant | read-only/denied |
| working changes | Builder | unsaved change set concept | conditionnel | session courante | clean session |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow | CMDR Studio | base identity/version/lifecycle | read |
| Version | CMDR Studio | base version/supersession | read |
| Principal | Platform Settings | editor/collaborator identity projection | read |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Builder Session | open/update/close conceptual editing context | CMDR Studio | ephemeral/work context, not canonical Workflow |
| Workflow draft | apply explicit save | CMDR Studio | conflict check; no silent overwrite |
## 11. Fonctionnalités
Open/edit against base version; dirty state; collaborator attribution; concurrent-update detection; conflict/recovery; validate; save; save-as-version-candidate; discard; preserve return-origin.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Open/edit session | Automation Designer | Builder Session | 0 | Workflow read | editing context | non |
| Run static validation | Automation Designer | draft graph | 1 | session open | validation assessment | non |
| Save draft changes | Automation Designer | Workflow draft | 2 | manage + no unresolved conflict | versioned draft change | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest edit | oui | oui | oui | oui | manual edit |
| conflict explanation | oui | oui | oui | oui | base/current deterministic diff |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
clean, dirty, validating, conflict, recovery-available, save-blocked, saved-to-draft, discarded, closed.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| updated Workflow draft | Workflow | CAP-STD-017/030 | explicit save + actor/base version |
| conflict assessment | assessment | Builder/user | base/current versions explicit |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Workflow Detail | edit | Builder Session | Workflow/base version/tenant/return-origin | Workflow Detail |
| Builder Session | save | Workflow draft | validated change set/base version | same session |
| Builder Session | conflict detected | conflict resolution context | base/current/change refs | Builder Session |
## 18. Dépendances
CAP-STD-017/030/031/032; Workflow/Version objects; Security; Shared Collaboration/Recovery/Versioning concepts; screen STD-BLD-001.
## 19. Source de vérité
Workflow remains canonical source; Builder Session is Studio editing context only. Principals remain Settings-owned; Shared retains generic collaboration/recovery mechanisms.
## 20. Provenance et audit
Record session ref, actor/collaborators, tenant/env, base version, save/discard actions, changed-node references, validation results, conflict/current version and correlation.
## 21. Permissions fonctionnelles
Workflow read/edit; Builder collaborate; validation run; save/discard; restricted metadata visibility. Final atomic namespace remains future. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Base version superseded, concurrent edit, missing permissions, offline mutation, failed save, inaccessible dependency or recovery gap stays explicit.
## 23. Métriques
Conflict frequency, unsaved-change loss prevention, save failures, recovery use, validation-before-save coverage.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** two editors modify the same base version, **When** one saves after the other, **Then** the second session detects a concurrent-edit conflict and does not silently overwrite.

**Given** a session has unsaved changes, **When** discard is selected, **Then** only the session changes are discarded and the canonical Workflow remains unchanged.

**Given** AI is disabled, **When** an editor resolves a conflict, **Then** deterministic base/current diff and manual selection remain available.
## 26. Questions ouvertes
OPEN-013 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Builder, Workflow Detail, CAP-STD-017/030/031/032, future STD-4 assurance/publishing, Security, Shared collaboration/recovery, Quality.
