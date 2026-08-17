---
id: CAP-STD-006
title: Tool Versioning, Compatibility and Deprecation
product: cmdr-studio
module: studio-foundations
owner: CMDR Studio Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-09
requirement_ids: [REQ-PROD-006, REQ-PROD-016, REQ-OBJ-009]
open_decisions: []
source-of-truth: canonical
---
# CAP-STD-006 — Tool Versioning, Compatibility and Deprecation
## 1. Définition
Version, compatibility, breaking-change, deprecation, supersession et reproducibility semantics pour Tool. Deprecated ≠ disabled ≠ deleted ≠ superseded.
## 2. Problème utilisateur
Un consumer doit pouvoir reproduire une invocation ancienne et voir toute incompatibilité sans silent upgrade.
## 3. Objectifs
Pin exact version, assess consumer/provider/runtime compatibility, expose migration/deprecation and preserve history.
## 4. Non-objectifs
Pas de final package manager, deployment pipeline, provider selection ou runtime implementation.
## 5. Propriétaire
Studio possède Tool version semantics; generic Versioning mechanisms restent Shared lorsque consommés.
## 6. Utilisateurs
Automation Designer, Studio Reviewer, Studio Operator et consumer maintainer.
## 7. Conditions d’entrée
Tool identity, candidate version, dependency/consumer refs and provenance available.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Tool | Studio | identity | oui | current | block |
| candidate version | designer | contract snapshot | oui | candidate | reject |
| dependency/consumer refs | Studio/Settings | compatibility context | oui | visible | unknown compatibility |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Version | CMDR Studio | version/supersession | read |
| Provider/Integration | Settings | compatibility/availability | read |
| Tool Call | Studio | historical exact-version ref | provenance read |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Tool version metadata | create/deprecate/supersede | CMDR Studio | no silent substitution |
## 11. Fonctionnalités
Immutable version reference, compatibility assessment, breaking-change indication, deprecation/supersession, migration notice and reproducibility impact.
## 12. Actions utilisateur
Class 0 inspect/compare; Class 1 compatibility check; Class 2 create version/deprecate/supersede metadata.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| compatibility check | oui | oui | oui | oui | compatibility matrix |
| migration suggestion | oui | oui | oui | oui | manual review |
## 14. États fonctionnels
draft/current/deprecated/disabled/superseded/retired-reference; historical refs remain resolvable.
## 15. États d’interface
Stale/Partial state exposes unresolved dependency compatibility; deprecated state never hides old version.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| exact version ref | Version | Library/Skills/Calls | reproducible/pinned |
| deprecation notice | metadata | consumers | history preserved |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Tool draft | create version | versioned Tool | contract/provenance | Tool |
| deprecated Tool | inspect replacement | consumer review | old/new refs | no auto upgrade |
## 18. Dépendances
CAP-STD-003, Shared Versioning, Settings dependency metadata and consumer compatibility inputs.
## 19. Source de vérité
Studio Tool/Version sources; external dependency versions remain owner-controlled.
## 20. Provenance et audit
Record old/new version, actor, reason, compatibility result, dependency snapshot and supersession.
## 21. Permissions fonctionnelles
Version read/create, compatibility inspect, deprecate/supersede; no final atomic namespace.
## 22. Limites et erreurs
Unknown compatibility, missing dependency, stale provider state, prohibited migration and unresolved consumer remain explicit.
## 23. Métriques
Deprecated usage, incompatible consumers, unknown compatibility, silent-upgrade violations target zero conceptually.
## 24. Classification de livraison
`defined / planned`; no package/runtime implementation.
## 25. Critères d’acceptation
**Given** deprecated Tool version, **When** opened, **Then** exact old version remains visible and no silent upgrade.  
**Given** provider compatibility unknown, **When** assessed, **Then** compatibility stays unknown/Partial.  
**Given** IA off, **When** versions compared, **Then** deterministic matrix remains usable.
## 26. Questions ouvertes
No new OPEN; detailed deployment/promotion belongs STD-4.
## 27. Consommateurs documentaires
Library, Skills, Tool Calls, future STD lots, Settings projections, Quality.
