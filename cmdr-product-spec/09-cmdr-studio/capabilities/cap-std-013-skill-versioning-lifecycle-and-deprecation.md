---
id: CAP-STD-013
title: Skill Versioning, Lifecycle and Deprecation
product: cmdr-studio
module: skills
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-013 — Skill Versioning, Lifecycle and Deprecation
## 1. Définition
Version/lifecycle Skill: draft, review-ready, available, deprecated, disabled, superseded et retired-reference lorsqu'applicable; publishing/deployment détaillé reste STD-4.
## 2. Problème utilisateur
Les consumers doivent conserver exact Skill version/dependencies et voir deprecation/supersession sans silent replacement.
## 3. Objectifs
Version identity, compatibility, dependency/consumer compatibility, supersession, migration et provenance.
## 4. Non-objectifs
Pas de promotion/canary/deployment runtime, assurance gates détaillés ou package manager.
## 5. Propriétaire
CMDR Studio / Skills possède lifecycle/version semantics; generic Versioning mechanism reste Shared si consommé.
## 6. Utilisateurs
Automation Designer, Studio Reviewer, consumer maintainers.
## 7. Conditions d’entrée
Skill current version, candidate version, dependency versions and consumers resolved.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Skill | Studio | identity | oui | current | block |
| candidate version | designer | contract snapshot | oui | candidate | reject |
| dependency/consumer refs | Studio | compatibility context | oui | visible | unknown compatibility |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Skill | CMDR Studio | lifecycle | read |
| Version | CMDR Studio | version/supersession | read |
| Tool/Skill deps | CMDR Studio | pinned refs | read |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Skill version/lifecycle | create/deprecate/disable/supersede | CMDR Studio | deployment detail deferred STD-4 |
## 11. Fonctionnalités
Create version, compare compatibility, mark lifecycle, supersede with migration notice and preserve reproducibility.
## 12. Actions utilisateur
Class 0 inspect/compare; Class 1 compatibility validate; Class 2 create version/change lifecycle metadata.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| compatibility validation | oui | oui | oui | oui | dependency checks |
| migration suggestion | oui | oui | oui | oui | manual review |
## 14. États fonctionnels
draft, review-ready, available, deprecated, disabled, superseded, retired-reference.
## 15. États d’interface
Deprecated/superseded shows exact current/historical refs and never hides provenance.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| versioned Skill ref | Version | Library/consumers | exact version preserved |
| lifecycle notice | metadata | consumers | no silent replacement |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Skill draft | create version | versioned Skill | contract/deps/provenance | Skill |
| deprecated Skill | review replacement | consumer | old/new refs/compatibility | no auto replacement |
## 18. Dépendances
CAP-STD-010/011/012, Shared Versioning and future STD-4 for promotion/deployment.
## 19. Source de vérité
Studio Skill/Version sources; consumer refs remain consumer-owned.
## 20. Provenance et audit
Actor, old/new version, dependency snapshot, compatibility result, lifecycle reason and supersession retained.
## 21. Permissions fonctionnelles
Skill/version read, version create, deprecate/disable/supersede; no final namespace.
## 22. Limites et erreurs
Unknown/incompatible dependency, stale consumer ref, forbidden lifecycle mutation and missing migration link explicit.
## 23. Métriques
Deprecated usage, incompatible consumers, orphaned supersession, lifecycle/provenance completeness.
## 24. Classification de livraison
`defined / planned`; no deployment/publishing engine.
## 25. Critères d’acceptation
**Given** deprecated Skill, **When** opened, **Then** exact version/history remain and no silent replacement.  
**Given** dependency incompatible, **When** candidate version checked, **Then** compatibility blocked/unknown explicitly.  
**Given** IA off, **When** versions compared, **Then** deterministic checks remain available.
## 26. Questions ouvertes
No new OPEN; detailed publishing/deployment belongs STD-4.
## 27. Consommateurs documentaires
Library, Skills, future STD-2/3/4, Command/Investigate consumers, Quality.
