---
id: CAP-STD-007
title: Tool Permission, Risk and Execution Eligibility
product: cmdr-studio
module: studio-foundations
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-SEC-001, REQ-SEC-002]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-STD-007 — Tool Permission, Risk and Execution Eligibility
## 1. Définition
Évalue permission needs, risk/side-effect/action class, runtime availability, compatibility and Govern dependency sans accorder l'autorité.
## 2. Problème utilisateur
Tool visible/available/compatible ne signifie jamais authorized ou executable.
## 3. Objectifs
Distinguer metadata read, invoke, sensitive invoke, cross-tenant refusal, Secret Reference use need, side effects, step-up/SoD and runtime eligibility.
## 4. Non-objectifs
Pas de final RBAC/ABAC, Decision/Approval, provider/runtime choice ou production-action semantics nouvelles.
## 5. Propriétaire
Studio possède Tool eligibility assessment; Security possède permission model; Govern possède response authority; runtime owner exécute.
## 6. Utilisateurs
Automation Designer, Studio Operator, Security Reviewer et consuming operator.
## 7. Conditions d’entrée
Exact Tool/version, tenant/env, permission context, risk metadata, runtime health and authority refs if required.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Tool/version | Studio | risk/permission metadata | oui | exact | block |
| authorization | Security | permission/scope | oui | current | permission-blocked |
| runtime health | Settings/runtime | availability | oui | freshness visible | unavailable |
| Govern context | Govern | authority ref | conditionnel | valid/current | governance-required |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Tool | Studio | risk/permission needs | read |
| Decision | Govern | authority projection | read only |
| Secret Reference | Settings | metadata | read metadata |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| execution eligibility assessment | derive | CMDR Studio | assessment != permission/Decision |
## 11. Fonctionnalités
Risk/action classification, permission need mapping, runtime/compatibility check, governance requirement and consumer eligibility.
## 12. Actions utilisateur
Class 0 inspect; Class 1 assess; Class 2 prepare eligible low-impact invocation. Class 3/4 effect requires external authority/runtime execution.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| eligibility check | oui | oui | oui | oui | permission/policy checks |
| risk explanation | oui | oui | oui | oui | risk matrix |
## 14. États fonctionnels
unknown, evaluating, eligible, permission-blocked, governance-required, runtime-unavailable, incompatible.
## 15. États d’interface
Permission denied masks protected data; Partial/Stale show which eligibility source is unavailable.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| eligibility result | assessment | CAP-STD-008/caller | no implicit grant |
| governance requirement | refs needed | caller/Govern | no auto Decision/Approval |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Tool | check eligibility | assessment | Tool/version/tenant/env/risk | Tool |
| eligible Tool | prepare call | CAP-STD-008 | auth refs/return-origin | caller |
## 18. Dépendances
Security permission model, OPEN-013, Settings runtime health, Govern authority and CAP-STD-003/004.
## 19. Source de vérité
Security/Govern/runtime owners remain authoritative; Studio derives only eligibility for the Tool context.
## 20. Provenance et audit
Record assessment inputs, exact permission/risk references, freshness, authority refs and reasons for eligible/blocked.
## 21. Permissions fonctionnelles
Metadata read, definition inspect, invoke need, sensitive/cross-tenant invoke distinction and output access. Both historical `perm.studio.*` and `perm.cmdr-studio.*` remain unresolved; no bulk rename.
## 22. Limites et erreurs
Missing permission, stale runtime, incompatible version, expired authority and tenant mismatch block eligibility explicitly.
## 23. Métriques
Eligibility outcomes by reason, denied/high-risk requests, stale runtime assessments and permission-namespace ambiguity occurrences.
## 24. Classification de livraison
`defined / planned`; no authorization/runtime engine implemented.
## 25. Critères d’acceptation
**Given** Tool visible but invoke denied, **When** eligibility evaluated, **Then** metadata may remain visible and invoke blocked.  
**Given** runtime unavailable, **When** checked, **Then** Tool remains defined but execution unavailable.  
**Given** Govern authority required, **When** eligibility evaluated, **Then** Studio creates no Approval/Decision and requires external reference.
## 26. Questions ouvertes
OPEN-013 remains open. Permission namespace ambiguity is documented but does not block functional needs and creates no new OPEN.
## 27. Consommateurs documentaires
CAP-STD-008, future orchestration/runtime lots, Security, Govern, Settings, Command/Investigate consumers, Quality.
