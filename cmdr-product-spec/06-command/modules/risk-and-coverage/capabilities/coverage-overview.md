---
id: CAP-CMD-203
title: Coverage Overview
product: command
module: risk-and-coverage
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-005, REQ-PROD-013, REQ-PROD-032, REQ-PROD-037]
open_decisions: [OPEN-013]
source-of-truth: canonical
---
# CAP-CMD-203 — Coverage Overview
## 1. Définition
Fournit une lecture sourcée des couvertures detection, endpoint, réponse, données et procédures sans posséder les moteurs/règles.
## 2. Problème utilisateur
Un service peut sembler protégé alors qu’une source, un agent, une procédure ou capacité de réponse est absente. Principal : Incident Commander ; secondaires : Detection Engineer, Platform Admin, Readiness Coordinator.
## 3. Objectifs
Séparer les familles ; afficher définition/population/source/fraîcheur ; relier lacunes aux services/work items ; créer improvement Task sans modifier les moteurs.
## 4. Non-objectifs
Gérer règles de détection, flotte Endpoint, agents Studio ou pourcentage global opaque.
## 5. Propriétaire
Command possède la lecture opérationnelle ; produits sources possèdent les couvertures.
## 6. Utilisateurs
Incident Commander principal ; experts sources et Readiness en consultation/action de suivi.
## 7. Conditions d’entrée
Service/scope, projections autorisées et définition de métrique ou unknown.
## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Detection coverage | Investigate | projection | non | version/date | unavailable |
| Endpoint/data/response | Settings/Endpoint/Govern | projections | non | source | gap par famille |
| Procedural coverage | Readiness | plans/tests | non | review/test | not-tested |
## 9. Objets lus
Coverage projections (source products), Service (Shared), Task (Command).
## 10. Objets créés ou modifiés
Task improvement (Command, classe 2) ; aucune mutation des sources de couverture.
## 11. Fonctionnalités
Afficher cinq familles ; montrer définition/numerator/denominator conceptuels/données manquantes ; lier gaps ; ouvrir owner ; créer action.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | lecteur | projection | 0 | source autorisée | limites visibles | non |
| Ouvrir source | lecteur | source product | 0 | permission | transition | non |
| Signaler gap | coordinateur | Task | 2 | gap sourcé | Task | OPEN-013 |
| Attendre validation | Readiness | Task/status | 2 | owner/due | follow-up | non |
## 13. Automatisation et IA
Calculs déterministes par source ; IA peut résumer, jamais inventer coverage. Sans IA : projections, métriques et actions manuelles.
## 14. États fonctionnels
`covered`, `partial`, `gap`, `degraded`, `not-tested`, `unavailable`, `unknown`.
## 15. États d’interface
Chaque famille peut être Partial indépendamment ; source/date/définition restent visibles. Rendu DS.
## 16. Sorties
Coverage context vers Risk/Mission Control/Readiness ; Task improvement avec gap/expected result.
## 17. Transitions
Détail vers Investigate/Settings/Govern/Studio avec return ; gap vers Readiness, Task reste Command.
## 18. Dépendances
CAP-CMD-201/301/303, Metrics Engine, Catalog, Data Quality.
## 19. Source de vérité
Chaque produit source possède sa couverture ; Command ne fait qu’agréger les projections.
## 20. Provenance et audit
Définition/version/source/population/fraîcheur, Task et acteur audités.
## 21. Permissions fonctionnelles
Lecture Command/source ; `perm.command.task.manage`. Atomisation reportée.
## 22. Limites et erreurs
Définitions incompatibles, données manquantes, double comptage, source stale ou refus empêchent toute comparaison opaque.
## 23. Métriques
Familles avec définition/fraîcheur ; gaps sans owner/action ; délai gap→validation. Pas de cible.
## 24. Classification de livraison
`defined` / `planned`, cible native ; sources techniques non définies ici.
## 25. Critères d’acceptation
**Given** detection covered et endpoint unavailable, **When** coverage est consultée, **Then** les familles restent séparées et le total n’est pas présenté comme complet.

**Given** un gap sourcé, **When** une Task est créée, **Then** source, service, owner et expected result sont liés.

**Given** aucun modèle, **When** la capability est utilisée, **Then** projections et métriques déterministes suffisent.
## 26. Questions ouvertes
Quelles définitions sont comparables ? Comment éviter double comptage ? — requirements ci-dessus. `OPEN-013` reste ouverte.
## 27. Consommateurs documentaires
Risk, Mission Control, Readiness, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.
