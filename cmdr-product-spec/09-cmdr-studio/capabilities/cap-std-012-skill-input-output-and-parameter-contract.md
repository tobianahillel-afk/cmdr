---
id: CAP-STD-012
title: Skill Input, Output and Parameter Contract
product: cmdr-studio
module: skills
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009, REQ-SEC-001]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-012 — Skill Input, Output and Parameter Contract
## 1. Définition
Contrat fonctionnel Skill I/O: object refs, parameter refs, required/optional, Tool input mapping, consumer output mapping, validation et sensitive/Secret Reference boundaries.
## 2. Problème utilisateur
Un consumer doit mapper I/O explicitement sans imposer un format technique final ni faire transiter un secret brut.
## 3. Objectifs
Définir inputs/outputs conceptuels, validation, mappings, sensitive refs et error semantics.
## 4. Non-objectifs
Aucun JSON Schema final, runtime serialization, raw secret, Tool invocation ou Workflow binding automatique.
## 5. Propriétaire
CMDR Studio / Skills possède le Skill I/O contract; Tool contracts et Settings Secret References restent chez leurs owners.
## 6. Utilisateurs
Automation Designer, Studio Reviewer, consumer product analyst.
## 7. Conditions d’entrée
Skill/version, consumer context et referenced Tool/Secret contracts disponibles.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Skill/version | Studio | exact ref | oui | pinned | block |
| input/parameter refs | consumer | functional bindings | oui | request | invalid/missing |
| Secret Reference | Settings | opaque ref | conditionnel | auth current | restricted |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Skill | CMDR Studio | I/O contract | read |
| Tool | CMDR Studio | input mapping refs | read |
| Secret Reference | Platform Settings | metadata/ref | metadata read |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Skill I/O contract | define/update | CMDR Studio | functional format only |
## 11. Fonctionnalités
Input/output concepts, refs, required/optional, Tool mapping concept, consumer mapping, validation and sensitive handling.
## 12. Actions utilisateur
Class 0 inspect I/O; Class 1 validate/match; Class 2 edit Skill draft mapping. No Tool execution.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| I/O validation | oui | oui | oui | oui | deterministic validation |
| mapping suggestion | oui | oui | oui | oui | manual mapping |
## 14. États fonctionnels
unvalidated, valid, invalid, missing-input, restricted-input, incompatible-output.
## 15. États d’interface
Error/Partial identify field/mapping issue; Secret values never displayed/logged.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| validated Skill bindings | refs | Skill consumer | required/optional explicit; no raw secret |
| output mapping | functional map | consumer | Skill semantics preserved |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| consumer | bind Skill | Studio Skill | refs/tenant/env | consumer |
| Skill | map Tool input | Tool contract | validated refs only | Skill |
## 18. Dépendances
CAP-STD-004/005/010/011/015, Settings Secret References and consumer object owners.
## 19. Source de vérité
Skill I/O semantics are Studio-owned; source values remain consumer-owned; secret values stay Settings-controlled.
## 20. Provenance et audit
Record Skill/version, mapping refs, validation, consumer, tenant/env and correlation without raw sensitive values.
## 21. Permissions fonctionnelles
Skill I/O read/update, restricted ref metadata read and mapping validation; final atomic permissions deferred.
## 22. Limites et erreurs
Missing/invalid/restricted input, incompatible output, stale ref or cross-tenant binding fails explicitly.
## 23. Métriques
Validation failures, missing inputs, mapping incompatibilities, restricted bindings and no-AI coverage.
## 24. Classification de livraison
`defined / planned`; no serialization/runtime implementation selected.
## 25. Critères d’acceptation
**Given** required Skill input absent, **When** validation runs, **Then** failure explicit and no default invented.  
**Given** Secret Reference, **When** mapping displayed, **Then** opaque ref only and Settings owner visible.  
**Given** IA off, **When** I/O mapped, **Then** deterministic/manual path remains functional.
## 26. Questions ouvertes
No new OPEN; technical formats deferred.
## 27. Consommateurs documentaires
Skills, Tools, Library, future Workflow/Agents, Command/Investigate consumers, Security/Settings, Quality.
