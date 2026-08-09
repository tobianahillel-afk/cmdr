---
id: CAP-STD-011
title: Skill Composition, Dependencies and Preconditions
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
# CAP-STD-011 — Skill Composition, Dependencies and Preconditions
## 1. Définition
Contrat fonctionnel des child Skill refs, Tool refs, dependency types/versions, required/optional status, preconditions, constraints, incompatibilities et cycle candidates. Skill composition ≠ Workflow orchestration.
## 2. Problème utilisateur
Une Skill réutilisable doit rendre ses dépendances explicites sans substituer silencieusement une version, un Tool ou une orchestration.
## 3. Objectifs
Construire un dependency graph typé/versionné et valider missing/incompatible/cyclic conditions.
## 4. Non-objectifs
Pas de Workflow graph/branches/retries, runtime scheduler, Agent composition ou Tool permission inheritance.
## 5. Propriétaire
CMDR Studio / Skills possède dependency semantics; Tools et child Skills conservent leurs propres contracts/permissions.
## 6. Utilisateurs
Automation Designer, Studio Reviewer et consumer product analyst.
## 7. Conditions d’entrée
Parent Skill/version, dependency refs/versions et consumer context disponibles.
## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| parent Skill | Studio | Skill/version | oui | pinned | block |
| dependencies | Studio | child Skill/Tool refs | oui | version visible | Partial |
| compatibility context | Studio/Settings | versions/runtime constraints | oui | freshness visible | unknown/incompatible |
## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Skill | CMDR Studio | parent/child refs | read |
| Tool | CMDR Studio | dependency contract | read |
| Version | CMDR Studio | dependency version | read |
## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Skill dependency projection | add/update typed refs | CMDR Studio | composition != orchestration |
## 11. Fonctionnalités
Required/optional deps, dependency type/version, preconditions, incompatibilities, cycle candidate detection and consumer visibility.
## 12. Actions utilisateur
Class 0 inspect graph; Class 1 validate deps/cycles; Class 2 edit dependency refs in a Skill draft.
## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---|---|---|---|---|
| dependency validation | oui | oui | oui | oui | graph/cycle checks |
| dependency suggestion | oui | oui | oui | oui | manual selection |
## 14. États fonctionnels
complete, partial, dependency-missing, incompatible, cycle-candidate, blocked.
## 15. États d’interface
Partial identifies exact missing/incompatible dependency; denied dependency metadata remains masked.
## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| dependency graph | projection | Skill/Library/consumer | typed/versioned refs |
| compatibility assessment | assessment | consumer | missing/incompatible explicit |
## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Skill | validate deps | assessment | parent+dependency versions | Skill |
| Library | inspect dependency | Tool/Skill surface | ref/return-origin | Library |
## 18. Dépendances
CAP-STD-003/006/010, Settings availability projections, Shared Linking/Versioning.
## 19. Source de vérité
Parent Skill owns its dependency declaration; referenced Tool/Skill owners own their contracts/lifecycles.
## 20. Provenance et audit
Record dependency add/remove/change, actor, parent version, child ref/version, required flag and validation result.
## 21. Permissions fonctionnelles
Skill read/update, restricted dependency metadata read; Tool invoke permission remains separate.
## 22. Limites et erreurs
Missing/denied dependency, incompatible version, unresolved ref, cycle candidate or stale runtime info remain explicit.
## 23. Métriques
Dependency completeness, missing/incompatible/cycle candidates and silent-substitution violations.
## 24. Classification de livraison
`defined / planned`; no orchestration runtime.
## 25. Critères d’acceptation
**Given** required dependency missing, **When** validation runs, **Then** Skill becomes Partial/unavailable and no substitution occurs.  
**Given** Skill references Tool, **When** consumer opens Skill, **Then** Tool permission remains distinct.  
**Given** composition graph, **When** displayed, **Then** it does not become implicit Workflow orchestration.
## 26. Questions ouvertes
No new OPEN; orchestration belongs STD-2.
## 27. Consommateurs documentaires
Skills, Library, future Workflow/Agents, Command/Investigate consumers, Quality.
