---
id: CAP-STD-022
title: Conditions, Branches and Deterministic Decision Logic
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
# CAP-STD-022 — Conditions, Branches and Deterministic Decision Logic
## 1. Définition
Sémantique des conditions déterministes et branches de Workflow : source value, comparison, boolean/multi-way outcomes, unknown/missing handling, default route concept, explicit error path and provenance. Condition ≠ Govern Policy; Branch ≠ Govern Decision.
## 2. Problème utilisateur
Un contrôle de flux ambigu peut être confondu avec une décision d’autorité ou masquer le traitement de valeurs absentes.
## 3. Objectifs
Rendre chaque condition explicable/rejouable, définir outcomes/routes et traitement unknown/missing, sans imposer langage d’expression.
## 4. Non-objectifs
Aucune Policy Govern, Decision Govern, AI obligatoire, expression language final, runtime evaluator ou autorisation.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| condition definition | Workflow draft | source/comparison/routes | oui | draft/version | invalid condition |
| source value definition | Workflow data context | typed/ref input | oui | binding-time/runtime future | unknown |
| branch targets | Workflow graph | destination nodes | oui | draft | invalid/unreachable |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow | CMDR Studio | condition/branch graph context | read/write owner path |
| Policy / Decision | Govern | distinct references only | read only when linked |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Condition definition | create/update conceptual | CMDR Studio | deterministic control only |
| Branch edge | create/update | CMDR Studio | no authority semantics |
| evaluation preview | derive no-effect | CMDR Studio | not runtime/business decision |
## 11. Fonctionnalités
Boolean/multi-way branch; deterministic comparison; unknown/missing; explicit/default route; error route; evaluation preview; provenance.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect branch | Automation Designer | Condition | 0 | read | routes visible | non |
| Preview condition | Studio Reviewer | Condition | 1 | sample/context value | no-effect result | non |
| Edit condition/routes | Automation Designer | Workflow draft | 2 | manage | draft graph changed | OPEN-013 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| suggest condition/branch | oui | oui | oui | oui | manual expression/config |
| evaluate preview | oui | oui | oui | non requis | deterministic evaluator concept |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
unconfigured, valid, invalid, source-missing, outcome-unknown, route-missing, valid-for-review.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| condition/branch contract | Workflow graph | CAP-STD-020/030 | deterministic and versioned |
| preview result | no-effect assessment | Builder/Reviewer | not Govern Decision/Policy outcome |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Workflow data context | evaluate preview | Condition | source value/type | Builder |
| Condition | select route concept | Workflow graph | evaluation result/provenance | Workflow |
| Condition error | route | error path | reason/source ref | Workflow |
## 18. Dépendances
CAP-STD-019/020/023/026/030; Govern Policy/Decision as non-equivalent; future runtime engine unspecified.
## 19. Source de vérité
Studio owns Workflow control-flow definitions; Govern remains source for Policy/Decision authority. Preview results are derived Studio assessments only.
## 20. Provenance et audit
Record condition version, source ref/type, comparison/config, result incl. unknown, selected preview route, actor and correlation.
## 21. Permissions fonctionnelles
Workflow read/edit; condition configure; source metadata read. Source-value access and Govern authority permissions remain external. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Missing input, incompatible type, ambiguous/unknown result, route missing, inaccessible source or invalid condition remain explicit.
## 23. Métriques
Invalid conditions, unknown outcomes, missing default/error routes, AI-vs-manual authoring parity, provenance completeness.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** a branch input is absent, **When** condition preview runs, **Then** outcome is unknown/missing and the configured missing/error behavior is shown rather than guessed.

**Given** AI proposes a branch, **When** it is accepted into a draft, **Then** deterministic condition semantics and human-visible provenance are preserved.

**Given** a branch is valid, **When** it is reviewed, **Then** it is never represented as a Govern Decision or Policy.
## 26. Questions ouvertes
OPEN-013 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
Workflow graph, Builder, CAP-STD-019/020/023/026/030, Govern boundaries, Quality.
