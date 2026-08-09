---
id: CAP-STD-027
title: Workflow Retry, Idempotency and Duplicate-Execution Protection
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
# CAP-STD-027 — Workflow Retry, Idempotency and Duplicate-Execution Protection
## 1. Définition
Contrat fonctionnel des retry semantics : eligibility, conceptual max attempts/delay/backoff, idempotency expectation, duplicate-execution candidates, previous-attempt reference, retry-safe/unsafe, permission/condition recheck and cancellation. Retry ≠ authorization renewal; idempotency ≠ exactly-once guarantee.
## 2. Problème utilisateur
Répéter un step sans bornes ni contexte peut provoquer double effet, réutiliser une autorisation expirée ou transformer un échec en boucle infinie.
## 3. Objectifs
Définir retry eligibility and bounded policy semantics, recheck prerequisites, preserve attempt lineage and duplicate-risk signals without selecting runtime algorithm.
## 4. Non-objectifs
Aucun scheduler, concrete backoff algorithm, exactly-once claim, automatic authorization renewal, Automation Run lifecycle or retry command.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| step/error context | CAP-STD-021/026 | step/error/retry metadata | oui | current assessment | retry blocked |
| idempotency expectation | Tool/Skill/runtime contract | safe/unsafe/key/scope concept | oui | exact dependency version | unknown = unsafe |
| permission/condition context | Security/Govern/Workflow | current refs | oui | future recheck | retry blocked |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow | CMDR Studio | retry config | read/write |
| Tool/Tool Call | CMDR Studio | idempotency/retry metadata + attempt refs | read/reference |
| Decision | Govern | authority expiry/context if required | read only |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Retry Policy concept | create/update in Workflow | CMDR Studio | not runtime scheduler object |
| duplicate-risk assessment | derive | CMDR Studio | candidate, not exactly-once proof |
## 11. Fonctionnalités
Retry eligibility; bounded attempts; delay/backoff concept; previous-attempt lineage; safe/unsafe classification; duplicate detection candidate; permission/condition recheck; cancellation.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect retry policy | Automation Designer | Workflow step | 0 | read | policy visible | non |
| Assess retry safety | Studio Reviewer | retry config | 1 | contracts available | safe/unsafe/unknown | non |
| Configure bounded retry | Automation Designer | Workflow draft | 2 | manage | policy changed | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest retry configuration | oui | oui | oui | oui | manual configuration |
| idempotency/duplicate assessment | oui | oui | oui | oui | deterministic metadata checks |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
not-configured, eligible, unsafe, unknown, bounded, cancelled, duplicate-candidate, authority-recheck-required.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| retry policy contract | Workflow step config | future STD-3 runtime | bounded/recheck semantics |
| duplicate-execution assessment | assessment | Reviewer/future runtime | candidate only, no exactly-once claim |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| error path | consider retry | retry assessment | step/error/attempt refs | Workflow |
| retry candidate | recheck permission/condition | Security/Govern/Workflow condition | current refs | candidate |
| future retry attempt | link previous attempt | Tool Call/Automation Run future | attempt lineage | STD-3 |
## 18. Dépendances
CAP-STD-007/008/009/021/022/026/030; implementation idempotency contract; Security/Govern; future STD-3 runtime.
## 19. Source de vérité
Workflow owns retry configuration semantics; Tool/Tool Call/runtime own operation/attempt facts; Security/Govern remain authoritative at each future retry.
## 20. Provenance et audit
Record Workflow/version/step, policy version, attempt refs, eligibility reasons, idempotency expectation/scope, permission/condition recheck results and cancellation.
## 21. Permissions fonctionnelles
Workflow retry configure/validate; attempt metadata read; permission/Decision read if applicable. Retry never grants missing permission. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Unknown idempotency, non-retry-safe operation, expired authority, missing condition, max-attempt concept exceeded, duplicate candidate or cancellation blocks further retry.
## 23. Métriques
Unsafe/unknown retry configs, duplicate candidates, bounded-policy coverage, permission-recheck failures, retry lineage completeness.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** an operation is marked non-retry-safe, **When** retry readiness is checked, **Then** automatic retry is blocked and manual/governed handling is required.

**Given** a duplicate-execution candidate exists, **When** another retry is considered, **Then** previous-attempt lineage is surfaced and exactly-once is not claimed.

**Given** prior authorization existed, **When** a retry is considered later, **Then** permission/condition/authority are rechecked rather than renewed implicitly.
## 26. Questions ouvertes
OPEN-013; OPEN-015 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Workflow/Builder, future STD-3 runtime, Tool Call contracts, Govern/Security, Quality.
