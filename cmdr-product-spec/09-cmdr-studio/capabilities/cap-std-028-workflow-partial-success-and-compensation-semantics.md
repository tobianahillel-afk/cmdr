---
id: CAP-STD-028
title: Workflow Partial Success and Compensation Semantics
product: cmdr-studio
module: workflows
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-10
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-PROD-019, REQ-OBJ-009, REQ-SEC-001, REQ-AI-002]
open_decisions: [OPEN-013, OPEN-015]
source-of-truth: canonical
---
# CAP-STD-028 — Workflow Partial Success and Compensation Semantics
## 1. Définition
Contrat fonctionnel de partial success et compensation : completed/failed/skipped steps, blocked downstream paths, compensation candidates/order/limits, manual intervention and resulting Workflow disposition candidate. Compensation ≠ Govern rollback et n’est jamais une inverse operation garantie.
## 2. Problème utilisateur
Une orchestration partiellement réussie peut sinon être présentée comme succès ou déclencher une pseudo-rollback non autorisée.
## 3. Objectifs
Rendre l’état partiel, les effets déjà produits, les candidats de compensation et leurs limites explicitement traçables avant tout runtime.
## 4. Non-objectifs
Aucune commande de compensation, Govern rollback, production mutation, inverse guarantee, Response Run state or runtime compensator.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| step outcomes | future runtime/Tool Calls | completed/failed/skipped refs | conditionnel | event time future | unknown disposition |
| compensation definitions | Workflow draft | candidate step/order/limitations | conditionnel | version pinned | manual only |
| authority/risk metadata | Security/Govern/Tool | effect constraints | conditionnel | future recheck | compensation blocked |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow | CMDR Studio | compensation config | read/write |
| Tool/Tool Call | CMDR Studio | technical outcome/compensation capability refs | read/reference |
| Response Rollback / Response Run | Govern | distinct governance refs | read only |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Compensation Step concept | configure in Workflow | CMDR Studio | definition, no command |
| partial-success assessment | derive conceptual | CMDR Studio | not canonical Result |
## 11. Fonctionnalités
Partial path accounting; completed/failed/skipped distinction; downstream blocking; compensation candidates/order/limitations; manual intervention; disposition candidate.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect compensation config | Automation Designer | Workflow | 0 | read | config visible | non |
| Validate compensation coverage | Studio Reviewer | Workflow | 1 | effects/config known | assessment | non |
| Configure compensation candidate | Automation Designer | Workflow draft | 2 | manage | draft changed | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest compensation candidate | oui | oui | oui | oui | manual selection |
| assess coverage/order | oui | oui | oui | oui | deterministic dependency analysis |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
no-partial-context, potential-partial, compensation-configured, compensation-gap, manual-intervention-required, valid-for-review.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| compensation contract | Workflow config | CAP-STD-030/future STD-3 | candidate/order/limits only |
| partial-success assessment | assessment | future Control Room/Govern consumer | not success/Result |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| error path | partial state | compensation assessment | step outcomes/effects | Workflow |
| compensation candidate | authority check future | Govern/Security/runtime owner | effect/risk refs | STD-3 |
| manual intervention | handoff | operator/Govern | remaining effects/provenance | Workflow context |
## 18. Dépendances
CAP-STD-005/009/020/026/027/030; Govern rollback/Result boundaries; future STD-3 runtime.
## 19. Source de vérité
Studio owns Workflow compensation definitions and partial-state orchestration semantics; Govern owns rollback/Response Run/Result; technical owners own actual effects.
## 20. Provenance et audit
Record completed/failed/skipped refs, effect/output refs, compensation config/order, gaps, manual intervention reason, Workflow/version and correlation.
## 21. Permissions fonctionnelles
Workflow compensation configure/read/validate; effect metadata read; any future effectful compensation requires separate permissions/authority. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Unknown effect, unavailable compensation, unsafe/incompatible compensation, partial outputs, authority absent or ordering conflict remain explicit.
## 23. Métriques
Compensation coverage/gaps, manual-intervention candidates, partial-path configurations, false-success violations, provenance completeness.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** some steps complete and a later step fails, **When** disposition is assessed, **Then** the Workflow is partial rather than success and completed effects remain visible.

**Given** no compensation is available for a completed effect, **When** coverage is checked, **Then** the gap and manual-intervention need are explicit.

**Given** a compensation candidate exists, **When** it is reviewed, **Then** it is not represented as a Govern rollback or guaranteed inverse.
## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Workflow/Builder, CAP-STD-026/027/030/033, Govern, future STD-3 runtime, Quality.
