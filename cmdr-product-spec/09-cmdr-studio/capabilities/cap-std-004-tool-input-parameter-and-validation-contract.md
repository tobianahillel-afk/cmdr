---
id: CAP-STD-004
title: Tool Input, Parameter and Validation Contract
product: cmdr-studio
module: studio-foundations
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009, REQ-SEC-001]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-004 — Tool Input, Parameter and Validation Contract
## 1. Définition
Contrat fonctionnel des inputs/parameters Tool, validation, source/target refs, tenant/environment et Secret References sans JSON Schema final.
## 2. Problème utilisateur
Une invocation ne doit jamais inventer un default, accepter une valeur interdite ou exposer un secret pour compenser un input incomplet.
## 3. Objectifs
Distinguer required/optional, types fonctionnels, allowed/default concepts, refs, validation et restrictions cross-tenant.
## 4. Non-objectifs
Pas de format technique final, parser runtime, secret value ou orchestration.
## 5. Propriétaire
CMDR Studio possède le Tool input contract; Settings possède Secret References et tenant/env administration.
## 6. Utilisateurs
Automation Designer, Studio Operator, Investigate Analyst et consommateurs autorisés.
## 7. Conditions d’entrée
Tool/version exacts, caller scope, tenant/env et permission context résolus.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Tool definition | Studio | exact Tool/version | oui | pinned | bloque validation |
| input bindings | caller | typed values/refs | oui | request | missing/invalid |
| Secret Reference | Settings | opaque ref | conditionnel | auth courant | restricted |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Tool | CMDR Studio | input contract | read |
| Secret Reference | Platform Settings | metadata/ref | metadata read |
| Tenant/Environment | Platform Settings | scope refs | read |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Tool input contract | define/update | CMDR Studio | functional only |
| validation outcome | derive | CMDR Studio | no caller mutation |
## 11. Fonctionnalités
Required/optional inputs, functional types, constraints/default concept, source/target refs, validation and masking.
## 12. Actions utilisateur
Class 0 inspect; Class 1 validate bindings; Class 2 prepare corrected binding. No execution.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| validation | oui | oui | oui | oui | deterministic validator |
| mapping suggestion | oui | oui | oui | oui | manual binding |
## 14. États fonctionnels
unvalidated, valid, invalid, missing-input, restricted-input, cross-scope-blocked.
## 15. États d’interface
Partial/Error/Permission denied expose field-level reason without protected value leakage.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| validation result | assessment | caller/CAP-STD-008 | explicit missing/invalid/restricted |
| validated bindings | refs | Tool Call prep | no raw secret |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Tool | validate bindings | validation result | Tool/version/caller scope | caller |
| validation result | prepare invocation | CAP-STD-008 | validated refs/tenant/env | return-origin |
## 18. Dépendances
CAP-STD-003, Settings Secret References, Security tenant isolation and future runtime owner.
## 19. Source de vérité
Tool input semantics are Studio-owned; values remain caller-owned and secrets Settings-owned.
## 20. Provenance et audit
Record Tool/version, caller, parameter names/types, ref IDs, validation reasons and correlation; never raw secret.
## 21. Permissions fonctionnelles
Input-contract read/update, restricted-reference metadata read and validation; final atomic namespace deferred.
## 22. Limites et erreurs
Missing/invalid/restricted input, stale refs, cross-tenant target and unsupported type fail explicitly.
## 23. Métriques
Validation failures by reason, missing-input rate, restricted-input blocks and stale refs.
## 24. Classification de livraison
`defined / planned`; no runtime schema/validator implementation selected.
## 25. Critères d’acceptation
**Given** required input missing, **When** validation runs, **Then** it fails explicitly and invents no default.  
**Given** Secret Reference, **When** inspected, **Then** raw secret stays masked and Settings ownership visible.  
**Given** IA off, **When** bindings are prepared, **Then** deterministic validation/manual mapping remain usable.
## 26. Questions ouvertes
No new OPEN; physical schema and atomic permissions remain later-phase work.
## 27. Consommateurs documentaires
CAP-STD-003/008, Skills, future Builder/orchestration, Security, Settings and Quality.
