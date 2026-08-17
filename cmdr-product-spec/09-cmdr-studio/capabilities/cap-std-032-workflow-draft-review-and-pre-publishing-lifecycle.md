---
id: CAP-STD-032
title: Workflow Draft, Review and Pre-Publishing Lifecycle
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
# CAP-STD-032 — Workflow Draft, Review and Pre-Publishing Lifecycle
## 1. Définition
Lifecycle documentaire/fonctionnel STD-2 d’un Workflow avant publishing détaillé : draft, editing, review-ready, under-review, changes-requested, approved-for-publishing-candidate, deprecated and superseded. Approved-for-publishing-candidate ≠ published/deployed/active runtime.
## 2. Problème utilisateur
Sans frontière pre-publish, une validation ou review peut être interprétée comme déploiement ou activation automatique.
## 3. Objectifs
Formaliser passage draft→review, requested changes, review outcome candidate and history, while deferring assurance/publishing/deployment to STD-4.
## 4. Non-objectifs
Aucune promotion/deployment, runtime activation, canary, final Assurance approval, Govern authorization, release pipeline or detailed publishing mechanics.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Workflow/version candidate | CMDR Studio | draft/version | oui | pinned | no review |
| Validation Assessment | CAP-STD-030 | readiness findings | oui | current for version | review not ready |
| review participants | Settings principals + Studio role context | review identities | oui | current | review blocked |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow | CMDR Studio | lifecycle/version | read/write owner path |
| Version | CMDR Studio | candidate/history | read |
| Evaluation/Simulation/Deployment | CMDR Studio | future STD-4 refs only | read/reference only |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Workflow pre-publish state | transition | CMDR Studio | does not publish/deploy |
| review record concept | create/update | CMDR Studio | human review provenance, not Govern Approval |
## 11. Fonctionnalités
Review readiness; under-review; changes requested; candidate acceptance; deprecation/supersession context; handoff to future Assurance/publishing.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Submit for review | Automation Designer | Workflow | 2 | valid-for-review | under-review | OPEN-013 |
| Review Workflow | Studio Reviewer | Workflow | 0 | review access | review assessment | non |
| Request changes | Studio Reviewer | Workflow | 2 | under-review | changes-requested | OPEN-013 |
| Mark approved-for-publishing-candidate | Studio Reviewer | Workflow | 2 | review satisfied | candidate only | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| summarize changes/findings | oui | oui | oui | oui | deterministic diff + reviewer notes |
| recommend review issues | oui | oui | oui | oui | manual checklist |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
draft, editing, review-ready, under-review, changes-requested, approved-for-publishing-candidate, deprecated, superseded.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| pre-publish lifecycle state | Workflow | Library/future STD-4 | candidate state only |
| review record/provenance | review context | Audit/Quality | not Govern Approval |
| handoff candidate | reference package | future STD-4 Assurance | version + validation + review refs |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Builder | submit for review | Workflow review | version + validation + provenance | Builder |
| review | request changes | Builder | findings/change refs | review |
| review candidate | handoff | future Assurance/Publishing | Workflow/version/assessment | Workflow |
## 18. Dépendances
CAP-STD-018/030/031; Workflow object; Studio reviewers; future STD-4 Assurance/Versions & Deployment. Govern Approval remains separate.
## 19. Source de vérité
Studio owns Workflow pre-publish lifecycle. Future STD-4 owns detailed assurance/publishing/deployment; Govern owns response authority, not Studio content review.
## 20. Provenance et audit
Record state transition, actor/reviewer, Workflow/version, validation snapshot, review findings, change requests, candidate rationale and correlation.
## 21. Permissions fonctionnelles
Workflow read/edit/submit-review/review functional needs; pre-publish candidate marking. No deployment/publish atomic permission finalized. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Stale validation, version changed during review, reviewer denied, unresolved blockers, concurrent edit or superseded candidate returns to explicit blocked/changes state.
## 23. Métriques
Review cycle counts, stale-review invalidations, changes-requested rate concept, candidate age, accidental publish/deploy claims target zero.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** a Workflow changes after review starts, **When** review status is checked, **Then** the prior assessment is stale for the new version and cannot silently approve it.

**Given** a Workflow reaches approved-for-publishing-candidate, **When** status is displayed, **Then** it is explicitly not published, deployed or active runtime.

**Given** AI is disabled, **When** review occurs, **Then** deterministic validation/diff and human review notes provide the complete path.
## 26. Questions ouvertes
OPEN-013 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Workflow/Builder/Library, CAP-STD-030/031, future STD-4 Assurance/Publishing, Security, Quality, Roadmap.
