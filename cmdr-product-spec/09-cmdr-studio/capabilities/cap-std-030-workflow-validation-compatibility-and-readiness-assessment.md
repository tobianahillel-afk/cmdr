---
id: CAP-STD-030
title: Workflow Validation, Compatibility and Readiness Assessment
product: cmdr-studio
module: workflows
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-SEC-001, REQ-AI-002]
open_decisions: [OPEN-007, OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-030 — Workflow Validation, Compatibility and Readiness Assessment
## 1. Définition
No-effect assessment de graph validity, unreachable/invalid cycles, inputs, Tool/Skill/subworkflow versions, runtime/Secret refs, permission requirements, Human Gate consistency, mappings, unsafe retries, compensation gaps and deprecations. Valid ≠ published; compatible ≠ authorized; graph valid ≠ deployable.
## 2. Problème utilisateur
Un Workflow structurellement valide peut rester incomplet, incompatible, non autorisé ou non exécutable ; ces dimensions doivent être séparées.
## 3. Objectifs
Produire une readiness assessment explicable avec findings/warnings/blockers sans exécuter le Workflow ni simuler le business.
## 4. Non-objectifs
Aucune runtime success prediction, deployment approval, Assurance evaluation complète, production simulation, API/runtime selection or authority grant.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Workflow/version candidate | CMDR Studio | full STD-2 contract set | oui | pinned | cannot assess |
| dependency contracts | Studio/Settings/Shared/Govern refs | versions/health/authority needs | oui | freshness visible | Partial/blocked |
| permission context | Security | functional needs/current actor scope | oui | current | permission gaps |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow / Version | CMDR Studio | graph/I-O/dependencies/lifecycle | read |
| Tool/Skill/Human Gate | CMDR Studio | compatibility/deprecation/permission needs | read |
| Secret Reference / Environment | Platform Settings | metadata/availability | restricted read |
| Decision/Approval | Govern | required-authority projection only | read only |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Validation Assessment | create/refresh conceptual | CMDR Studio | no-effect assessment, no runtime object |
| readiness findings | derive | CMDR Studio | warning/blocker/superseded classifications |
## 11. Fonctionnalités
Graph/reachability/cycle; I/O/mapping; dependency/version; runtime refs; Secret refs; permission needs; Human Gate; retry/compensation; deprecation checks; readiness outcomes.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Run static validation | Automation Designer | Workflow | 1 | read + contracts | assessment | non |
| Review findings | Studio Reviewer | Validation Assessment | 0 | read | findings | non |
| Apply draft fix | Automation Designer | Workflow draft | 2 | manage | new draft state | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest fix/explain finding | oui | oui | oui | oui | manual fix + deterministic finding |
| validation execution | oui | oui | oui | non requis | deterministic check suite |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
not-run, running-concept, valid-for-review, incomplete, incompatible, warning, blocked, superseded.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Validation Assessment | assessment | Builder/Workflow Review | findings with exact refs |
| readiness outcome | classification | CAP-STD-032/future STD-4 | not publish/deploy/execute authorization |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Builder/Workflow | validate | Validation Assessment | Workflow/version/dependency snapshot | Workflow |
| assessment | fix finding | Builder | finding refs/return-origin | assessment refresh |
| valid-for-review | submit | CAP-STD-032 | Workflow/version/assessment | Workflow |
## 18. Dépendances
CAP-STD-017..029/031/032; Security; Settings; Shared; Govern boundaries; future STD-3 runtime and STD-4 assurance.
## 19. Source de vérité
Studio owns the no-effect validation/readiness assessment. Each dependency owner remains source for its facts; final assurance/deployment belongs later lot.
## 20. Provenance et audit
Record exact Workflow/version, validator/check versions conceptually, dependency snapshot, findings, skipped checks, actor, timestamp/freshness and correlation.
## 21. Permissions fonctionnelles
Workflow/assessment read; validation run; restricted dependency metadata; draft fix. Validation never grants invoke/publish/deploy authority. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Unavailable source, stale dependency, denied metadata, unsupported check, validation timeout concept or contradictory findings yields Partial/blocked and never fabricated PASS.
## 23. Métriques
Finding counts by category, unresolved blockers, deprecated dependencies, unknown runtime refs, no-AI validation coverage, stale-assessment usage.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** static validation passes but the future runtime is unavailable, **When** readiness is assessed, **Then** structural validation remains PASS while runtime availability remains unavailable and deployability is not inferred.

**Given** a Secret Reference exists but current permission is absent, **When** readiness is checked, **Then** the sensitive binding is blocked without revealing secret content.

**Given** AI is disabled, **When** validation runs, **Then** deterministic checks produce the complete readiness assessment and findings.
## 26. Questions ouvertes
OPEN-007; OPEN-013; OPEN-015 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Builder, Workflow Detail, CAP-STD-032, future STD-3 runtime/STD-4 Assurance, Security, Settings, Govern, Quality.
