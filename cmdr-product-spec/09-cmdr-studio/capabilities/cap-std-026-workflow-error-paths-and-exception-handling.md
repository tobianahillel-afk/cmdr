---
id: CAP-STD-026
title: Workflow Error Paths and Exception Handling
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
# CAP-STD-026 — Workflow Error Paths and Exception Handling
## 1. Définition
Contrat fonctionnel des error paths : validation error, Tool error, Skill dependency failure, condition failure, missing input, runtime-unavailable projection, timeout concept, unsupported/permission-blocked, partial output, explicit error route, terminate and continue-with-warning. Studio exception ≠ Govern Policy Exception.
## 2. Problème utilisateur
Sans modèle d’erreur explicite, un Workflow peut masquer un échec, poursuivre sur données invalides ou confondre exception technique et exception de gouvernance.
## 3. Objectifs
Classifier causes, définir routes/termination/warning conceptuels, préserver partial data and provenance; fournir base aux retries/compensation.
## 4. Non-objectifs
Aucune Policy Exception Govern, runtime exception engine, timeout implementation, automatic Incident, silent continue or generic Job error lifecycle.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| step/validation outcome | Workflow/Tool/Skill future status | error/warning/partial metadata | oui | event/assessment time | unknown error |
| error-route configuration | Workflow draft | route/action concept | oui | version pinned | unhandled |
| context/provenance | Studio/Shared refs | step/caller/correlation | oui | current | unattributed |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow | CMDR Studio | error route definitions | read/write |
| Tool Call | CMDR Studio | technical error ref only | read/reference |
| Policy/Exception Candidate | Govern | distinct governance ref | read only |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Workflow error-path definition | create/update | CMDR Studio | orchestration semantics only |
| error assessment | derive | CMDR Studio | does not mutate source error |
## 11. Fonctionnalités
Error taxonomy; explicit routes; terminate; continue-with-warning; partial-output handling; timeout/unsupported/permission-blocked projections; provenance.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect error path | Automation Designer | Workflow | 0 | read | route/reason | non |
| Validate error coverage | Studio Reviewer | Workflow | 1 | graph/contracts | assessment | non |
| Configure error route | Automation Designer | Workflow draft | 2 | manage | draft changed | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest error route | oui | oui | oui | oui | manual configuration |
| classify known error | oui | oui | oui | oui | deterministic taxonomy |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
unhandled, handled-route, terminate, continue-with-warning, partial, permission-blocked, timeout-concept, unsupported.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| error-path contract | Workflow definition | CAP-STD-027/028/030 | explicit route semantics |
| error assessment | assessment | Builder/future runtime | no hidden failure |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| step outcome | error | Workflow error path | step/error/provenance | Workflow |
| error path | retry candidate | CAP-STD-027 | error class + step contract | Workflow |
| error path | compensation candidate | CAP-STD-028 | completed/failed context | Workflow |
## 18. Dépendances
CAP-STD-005/009/019/022/027/028/030; Govern Policy Exception non-equivalence; Shared Trace.
## 19. Source de vérité
Source Tool/Skill/runtime owns raw technical error; Workflow owns orchestration response definition; Govern owns policy exceptions/authority.
## 20. Provenance et audit
Record Workflow/version, step ref, raw source error ref, classification, route choice/config, actor, correlation, partial outputs and warnings.
## 21. Permissions fonctionnelles
Workflow read/edit/validate; technical error metadata read; restricted output masking. No permission bypass by error route. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Unknown error, missing route, permission denial, partial output, timeout concept, unsupported operation or inaccessible source remain explicit.
## 23. Métriques
Unhandled paths, continue-with-warning configurations, hidden-error violations, timeout candidates, partial-output routes.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** a Tool step returns an error, **When** Workflow handling is evaluated, **Then** the configured error route/termination is explicit and no error is hidden.

**Given** a permission-blocked step, **When** error handling is evaluated, **Then** the Workflow cannot route around authorization as a successful step.

**Given** AI is disabled, **When** error paths are authored, **Then** deterministic taxonomy and manual route configuration remain available.
## 26. Questions ouvertes
OPEN-013 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Workflow, Builder, CAP-STD-027/028/030, Tool/Skill owners, Govern boundary, Shared Trace, Quality.
