---
id: CAP-STD-033
title: Workflow Provenance and Cross-Product Orchestration Contracts
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
# CAP-STD-033 — Workflow Provenance and Cross-Product Orchestration Contracts
## 1. Définition
Contrat STD-2 de provenance et frontières cross-product couvrant Workflow/Version/Builder Session/Tools/Skills/subworkflows/I-O/variables/conditions/mappings/retries/compensation/Human Gates/validation/owners/changes. Workflow ≠ Playbook; Automation Run ≠ Response Run; Workflow ≠ Endpoint primitive.
## 2. Problème utilisateur
Une orchestration réutilisée entre produits peut perdre le caller, l’owner, la version, l’autorité et le return-origin ou fusionner objets Studio/Govern/Endpoint.
## 3. Objectifs
Rendre tout handoff et reference chain attribuable/reproductible; préserver ownership et qualification chez le consumer; préparer future runtime bridge sans le définir.
## 4. Non-objectifs
Aucun Automation Run lifecycle détaillé, Response Run bridge final, Endpoint capability, API/protocol, runtime engine, output qualification auto or cross-tenant permission transfer.
## 5. Propriétaire
CMDR Studio Product Lead. Studio owns only the Workflow/Builder orchestration semantics described here; referenced products retain their canonical objects and authority.
## 6. Utilisateurs
Automation Designer principal; Studio Reviewer et, selon le handoff, Studio Operator, Command/Investigate analyst, Response Operator ou Auditor autorisés.
## 7. Conditions d’entrée
Tenant et environnement résolus, Workflow/version ou draft identifiable, actor authentifié, permission context disponible et références requises explicitement résolues.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Workflow provenance | CMDR Studio | version/graph/steps/change refs | oui | pinned | provenance Partial |
| caller/consumer context | Command/Investigate/Govern/Studio | source object + return-origin | conditionnel | current | unscoped handoff rejected |
| Govern authority refs | Govern | Decision/Response Run/Approval when applicable | conditionnel | exact/current | no authority inferred |
| Shared trace refs | Shared | Trace/Activity/Job refs | conditionnel | available | provenance gap explicit |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Workflow / Version | CMDR Studio | definition/version/provenance | read |
| Tool / Tool Call / Skill / Human Gate | CMDR Studio | refs and technical provenance | read/reference |
| Decision / Approval / Playbook / Response Run / Result | Govern | distinct refs | read only |
| Secret Reference / Environment | Platform Settings | opaque/config refs | restricted read |
| Trace / Activity / Job | Shared | generic provenance/background refs | consume |
| Endpoint capability | Endpoint Agent | future technical ref only | read projection future |
## 10. Objets créés ou modifiés
| Objet/concept | Opération | Propriétaire | Règle |
|---|---|---|---|
| Workflow provenance projection | append/reference | CMDR Studio | source history preserved |
| cross-product orchestration context | prepare/link | CMDR Studio | ownership/permissions unchanged |
| consumer output reference | link only | consumer owner | no Evidence/Finding/Result auto-qualification |
## 11. Fonctionnalités
End-to-end reference chain; caller/return-origin; exact versions; changes; step/mapping/retry/compensation/gate/validation refs; Govern/Settings/Shared/Endpoint boundaries.
## 12. Actions utilisateur
| Action | Rôle | Objet/concept | Classe | Précondition | Résultat | Govern |
|---|---|---|---|---|---|---|
| Inspect provenance | Automation Designer/Auditor | Workflow | 0 | read | chain visible | non |
| Validate owner/return-origin | Studio Reviewer | handoff context | 1 | refs resolved | assessment | non |
| Prepare cross-product reference context | Authorized operator | handoff context | 2 | source permission | reference package only | OPEN-013/015 |
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| summarize provenance | oui | oui | oui | oui | raw refs/timeline |
| validate ownership | oui | oui | oui | oui | registry/ref checks |

AI is optional. No essential capability in this contract requires a chatbot or model provider; AI suggestions remain reviewable, attributable and non-authorizing.
## 14. États fonctionnels
complete, partial, stale, owner-conflict, unresolvable-ref, authority-required, cross-scope-blocked.
## 15. États d’interface
Loading, Empty, Partial, Error, Offline, Permission denied and Stale expose missing source/freshness/permission explicitly. UI state never changes functional ownership or authorization.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| STD-2 provenance chain | projection | Auditor/consumers | exact Workflow/version/dependency/change refs |
| cross-product context package | references | Command/Investigate/Govern/future Endpoint | caller/owner/return preserved |
| output reference | reference only | consumer | qualification remains consumer-owned |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Command/Investigate | open/reference Workflow | Studio | source object/tenant/return-origin | source product |
| Govern Response Run/Step | handoff future execution | Studio future runtime | Decision/Run/scope/Workflow version refs | Govern via CAP-GOV-025 |
| Workflow | reference future Endpoint capability | Endpoint owner | technical capability ref only | Studio remains non-owner |
| Studio | consume Settings/Shared refs | Settings/Shared | opaque/config/trace refs | Studio |
## 18. Dépendances
CAP-STD-001..032; CAP-GOV-025; OPEN-007/013/015; Settings; Shared; Endpoint future; Command/Investigate ownership contracts.
## 19. Source de vérité
Each product keeps canonical objects and permissions. Studio owns only Workflow-side definitions/context/provenance; consumer owns conclusions/outcomes; Govern owns response authority; Endpoint owns technical primitive.
## 20. Provenance et audit
Record caller product/object, actor/system, tenant/env, Workflow/version, Builder base/change refs, every dependency version, mapping/condition/retry/compensation/gate/validation refs, Govern authority refs, Shared trace/correlation and return-origin.
## 21. Permissions fonctionnelles
Workflow/provenance read, restricted cross-product refs, reference-package preparation, provenance export preparation. No authority or permission transfers through links. `perm.cmdr-studio.*` and `perm.studio.*` remain coexisting historical namespaces; STD-2 performs no bulk rename and final RBAC/ABAC remains future.
## 22. Limites et erreurs
Unresolvable/stale ref, owner conflict, cross-tenant mismatch, missing return-origin, absent required authority or provenance gap blocks/marks Partial without inventing data.
## 23. Métriques
Provenance completeness, owner conflicts, missing return-origin, stale refs, cross-product handoff failures, auto-qualification violations target zero conceptually.
## 24. Classification de livraison
`defined / planned`. Documentary definition does not prove implementation, publishing, deployment, runtime availability or production execution.
## 25. Critères d’acceptation
**Given** Investigate references a Workflow and receives a technical output later, **When** provenance is inspected, **Then** Investigate retains qualification of Evidence/Finding and Studio only supplies attributed refs.

**Given** Govern requires a Response Run handoff, **When** a Workflow is referenced, **Then** Workflow remains distinct from Playbook/Response Run and CAP-GOV-025 authority/correlation is preserved.

**Given** AI is disabled, **When** cross-product provenance is reconstructed, **Then** exact refs, versions, owners, trace and return-origin provide the complete path.
## 26. Questions ouvertes
OPEN-007; OPEN-013; OPEN-015 remain open and are not resolved by this capability.
## 27. Consommateurs documentaires
All Studio modules, Command, Investigate, Govern, Settings, Shared, future Endpoint/STD-3, Security, Quality, Roadmap.
