---
id: CAP-STD-008
title: Tool Call Request and Invocation Context
product: cmdr-studio
module: studio-foundations
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009, REQ-SEC-001]
open_decisions: [OPEN-015]
source-of-truth: canonical
---
# CAP-STD-008 — Tool Call Request and Invocation Context
## 1. Définition
Contrat fonctionnel de préparation d'un Tool Call: Tool/version, initiator, caller product/object, user, tenant/env, bindings, Secret References, targets, reason, permissions, runtime/provider refs, timestamp, correlation et return-origin.
## 2. Problème utilisateur
Une invocation non attribuée ou non bornée rend l'exécution non traçable et peut casser les frontières de permission/tenant.
## 3. Objectifs
Préparer un request explicite, version-pinned, scope-bound et auditable avant runtime dispatch.
## 4. Non-objectifs
Aucun Automation Run orchestration, Response Run, Job générique, API/protocol ou commande réelle.
## 5. Propriétaire
CMDR Studio possède Tool Call semantics; final Tool Call object schema reste différé Phase 7; runtime owner exécute.
## 6. Utilisateurs
Studio Operator, Automation Designer, Investigate Analyst et consuming product operators.
## 7. Conditions d’entrée
Tool/version exact, bindings validés CAP-STD-004, eligibility CAP-STD-007, caller/return-origin et tenant/env résolus.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Tool/version | Studio | exact contract ref | oui | pinned | reject |
| validated bindings | CAP-STD-004 | inputs/targets/secret refs | oui | request | validation-pending |
| caller context | calling product | initiator/object/reason/return | oui | request | reject unattributed |
| eligibility | CAP-STD-007 | authorization/risk/runtime assessment | oui | current | blocked |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Tool | CMDR Studio | exact version/contract refs | read |
| Secret Reference | Platform Settings | opaque ref | metadata read |
| Decision | Govern | authority ref if required | read only |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Tool Call functional request | prepare/submit semantics | CMDR Studio | final schema deferred; not Automation/Response Run |
## 11. Fonctionnalités
Pin Tool/version, caller/calling object, user, tenant/env, validated inputs, targets, Secret refs, reason, auth context, runtime/provider, timestamp, correlation and return-origin.
## 12. Actions utilisateur
Class 0 inspect request; Class 1 validate request context; Class 2 prepare/submit explicitly eligible Tool Call. Side effects still obey runtime/Govern.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| request assembly | oui | oui | oui | oui | form/contract |
| context suggestion | oui | oui | oui | oui | manual binding |
## 14. États fonctionnels
prepared, validation-pending, ready, rejected, permission-blocked, runtime-unavailable; execution lifecycle belongs CAP-STD-009.
## 15. États d’interface
Partial/Denied identify missing/forbidden context without leaking secret/protected object data.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| attributed Tool Call request | functional request | runtime/CAP-STD-009 | exact Tool/version/caller/scope |
| return-origin context | context ref | caller | navigation preserved |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| calling product | prepare call | Studio | caller/object/tenant/env/return | caller |
| prepared call | submit eligible | runtime owner | Tool/version/bindings/auth/correlation | status to Studio |
## 18. Dépendances
CAP-STD-003/004/007/015/016, Security, Settings refs, Govern authority where required and runtime owner.
## 19. Source de vérité
Tool Call functional request semantics are Studio-owned; caller objects, Settings refs and Govern authority stay external.
## 20. Provenance et audit
Always retain initiator, calling product/object, exact Tool/version, bindings refs, targets, reason, auth context, tenant/env, runtime/provider refs, timestamp, correlation and return-origin.
## 21. Permissions fonctionnelles
Tool Call prepare/submit/read, sensitive input use need, target scope and cross-tenant refusal; no atomic permission family selected.
## 22. Limites et erreurs
Missing caller, invalid binding, expired authority, unavailable runtime, unsupported target or tenant mismatch blocks request explicitly.
## 23. Métriques
Prepared/ready/blocked requests, attribution completeness, stale eligibility, missing return-origin and sensitive-binding failures.
## 24. Classification de livraison
`defined / planned`; no Tool Call transport/runtime implementation.
## 25. Critères d’acceptation
**Given** eligible Tool + valid bindings, **When** Tool Call prepared, **Then** exact Tool/version/caller/tenant/env/correlation/return are pinned.  
**Given** Secret Reference, **When** request assembled, **Then** only opaque ref enters context.  
**Given** IA off, **When** request prepared, **Then** deterministic form/validation supports same path.
## 26. Questions ouvertes
OPEN-015 remains open; Tool Call != Automation Run != Response Run is fixed, bridge semantics are not.
## 27. Consommateurs documentaires
CAP-STD-009, future STD-2/3, Investigate/Command callers, Govern handoff, Settings/Security, Quality.
