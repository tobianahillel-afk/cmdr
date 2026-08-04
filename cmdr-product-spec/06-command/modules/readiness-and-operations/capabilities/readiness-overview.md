---
id: CAP-CMD-301
title: Readiness Overview
product: command
module: readiness-and-operations
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-005, REQ-PROD-013, REQ-PROD-021, REQ-PROD-057]
open_decisions: [OPEN-010, OPEN-013]
source-of-truth: canonical
---
# CAP-CMD-301 — Readiness Overview
## 1. Définition
Présente la préparation opérationnelle par capability, scénario, équipe, plan, exercice et action d’amélioration, avec owner, échéance, validation et lacunes.
## 2. Problème utilisateur
La présence d’une capability documentée ne prouve ni configuration, ni test, ni disponibilité tenant. Principal : Readiness Coordinator ; secondaires : Incident Commander, Team Lead, Business Owner. Sans capacité, les lacunes sont découvertes pendant l’incident.
## 3. Objectifs
Distinguer readiness et delivery ; montrer sources/date/scope ; relier gaps à Tasks ; ouvrir l’owner technique sans dupliquer Assurance/Health.
## 4. Non-objectifs
Ne pas certifier un agent, administrer la plateforme, promouvoir delivery mode ou confondre readiness/document status.
## 5. Propriétaire
Command possède l’assessment opérationnel ; Studio/Settings/products possèdent assurance, configuration et health sources.
## 6. Utilisateurs
Readiness Coordinator principal ; Incident Commander, Team Lead, Business Owner secondaires.
## 7. Conditions d’entrée
Scope/scenario, projections de capabilities/plans/exercices/health/tasks et permission.
## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Capability inventory | Capability Register | owner/delivery/target | oui | last reviewed | unknown |
| Plans/exercises/tasks | Command | readiness evidence | non | courante | gap visible |
| Health/assurance | Settings/Studio/products | evidence refs | non | timestamp/version | not-tested/unavailable |
## 9. Objets lus
Capability projection, Operational Plan, Exercise, Task, health/assurance sources.
## 10. Objets créés ou modifiés
Readiness assessment (Command, C2) ; improvement Task (Command, C2).
## 11. Fonctionnalités
Afficher readiness par scope ; séparer delivery/document status ; identifier gaps/dependencies ; lier action/owner ; ouvrir Studio Assurance ou Settings Health.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Consulter | lecteur | assessment | 0 | scope | status/sources | non |
| Créer improvement Task | coordinateur | Task | 2 | gap/expected result | Task | OPEN-013 |
| Update assessment | coordinateur | assessment | 2 | source/rationale | status updated | OPEN-013 |
| Ouvrir owner source | lecteur | capability source | 0 | permission | transition | non |
## 13. Automatisation et IA
Moteurs agrègent health/tests ; règles calculent gaps ; IA résume seulement. Humain valide assessment et actions. Sans IA : sources, règles, plans et Tasks.
## 14. États fonctionnels
`ready`, `partial`, `degraded`, `not-tested`, `unavailable`, `planned`, `unknown`.
## 15. États d’interface
Partial/Unknown montrent sources manquantes ; Offline conserve dernière assessment datée ; Permission denied masque l’evidence protégée. Rendu DS.
## 16. Sorties
Readiness assessment vers Mission Control/plans/reporting ; Gap Task vers owner/team avec expected validation/due.
## 17. Transitions
Détail technique vers Studio Assurance ; health vers Settings ; retours restaurent readiness.
## 18. Dépendances
CAP-CMD-203/303/304/305, Studio Assurance, Platform Health, Metrics Engine.
## 19. Source de vérité
Assessment Command ; evidence reste propriétaire de sa source ; delivery mode reste Capability Register.
## 20. Provenance et audit
Scope, source/evidence, date, assessor, rationale, before/after, tenant et correlation ID.
## 21. Permissions fonctionnelles
Command read/coordinate/task manage + source reads. Densité finale `OPEN-010`, mutations `OPEN-013`.
## 22. Limites et erreurs
Evidence stale/absente, capability unknown, tenant mismatch, health down ou permission refusée donnent partial/unknown, jamais ready inventé.
## 23. Métriques
Capabilities avec source/date ; gaps sans owner/due ; délai gap→closure validée.
## 24. Classification de livraison
`defined` / `planned`, cible native ; assessment n’est pas preuve de livraison.
## 25. Critères d’acceptation
**Given** une capability planned non testée, **When** readiness est consultée, **Then** delivery mode et readiness restent distincts et not-tested est sourcé.

**Given** un gap, **When** une Task est créée, **Then** owner, expected validation et source sont liés.

**Given** aucun modèle, **When** la vue est utilisée, **Then** sources/règles/Tasks suffisent.
## 26. Questions ouvertes
Quel record canonique porte assessment ? Quels scénarios obligatoires par tenant ? — requirements ci-dessus. `OPEN-010` et `OPEN-013` restent ouvertes.
## 27. Consommateurs documentaires
Readiness screen, Mission Control, reports, parcours Phase 5, écrans Phase 6, objets Phase 7 et permissions ultérieures.
