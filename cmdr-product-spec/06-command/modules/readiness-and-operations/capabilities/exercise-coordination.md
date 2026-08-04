---
id: CAP-CMD-302
title: Exercise Coordination
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
# CAP-CMD-302 — Exercise Coordination
## 1. Définition
Coordonne un exercice opérationnel avec scénario, objectifs, participants, observations, résultats et actions d’amélioration, sans simulation technique Studio/Sandbox.
## 2. Problème utilisateur
Les exercices hors plateforme ne produisent pas d’actions traçables. Principal : Readiness Coordinator ; secondaires : Incident Commander, Team Lead, Business Owner, Observers.
## 3. Objectifs
Planifier scope/participants/objectifs ; enregistrer observations/results ; créer Tasks d’amélioration ; distinguer exercice opérationnel et simulation technique.
## 4. Non-objectifs
Ne pas orchestrer sandbox, tester agent/règle, exécuter Response Run ou simuler un Incident réel.
## 5. Propriétaire
Command possède la coordination/record opérationnel ; Studio/Sandbox possèdent simulation technique.
## 6. Utilisateurs
Readiness Coordinator principal ; participants et observers selon permission.
## 7. Conditions d’entrée
Scénario/objectifs, participants, tenant, autorisations et procédures liées.
## 8. Entrées fonctionnelles
Exercise definition (coordinator), participants/roles (identity projections), observations/results (participants). Un exercice ne peut être completed sans result summary.
## 9. Objets lus
Operational Plan (Command readiness), Capability Readiness, Workflow/Simulation refs (Studio, projection).
## 10. Objets créés ou modifiés
Exercise record (Command readiness, C2) ; improvement Task (Command, C2).
## 11. Fonctionnalités
Définir scenario/goals/participants/schedule ; observations attribuées ; outcome/gaps ; Tasks ; comparaison readiness before/after conceptuelle.
## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer/planifier | coordinateur | exercise | 2 | scenario/objectives | planned | OPEN-013 |
| Démarrer/terminer | coordinateur | exercise | 2 | participants/scope | running/completed | OPEN-013 |
| Ajouter observation | participant | exercise | 2 | active | observation attribuée | non |
| Créer Task | coordinateur | Task | 2 | gap | Task liée | OPEN-013 |
## 13. Automatisation et IA
Règles vérifient complétude ; moteur agrège observations ; workflow notifie ; agent peut préparer synthèse attribuée. Sans IA : formulaire, observations et Tasks manuelles.
## 14. États fonctionnels
`draft`, `planned`, `ready`, `running`, `observing`, `completed`, `cancelled`, `superseded`.
## 15. États d’interface
Partial montre participants/sources unresolved ; Offline bloque start/complete ; Permission denied masque observations protégées. Rendu DS.
## 16. Sorties
Exercise record vers participants/reporting ; improvement Tasks vers Work Queue avec source/expected validation.
## 17. Transitions
Besoin simulation vers Studio/Sandbox owner ; gap vers Improvement Actions ; résultats référencés au retour.
## 18. Dépendances
CAP-CMD-301/303/304, Collaboration, Notifications, Reporting Engine.
## 19. Source de vérité
Exercise record/observations Command ; simulations et workflows restent Studio ; identités restent Settings.
## 20. Provenance et audit
Auteur/temps de chaque observation, version scenario, participants, transitions, résultat et correlation ID.
## 21. Permissions fonctionnelles
Command read/coordinate + future participant contribution. `OPEN-010/013` restent ouvertes.
## 22. Limites et erreurs
Participant unresolved, source/procedure stale, exercice interrompu, conflit ou permission refusée n’autorisent pas completed trompeur.
## 23. Métriques
Exercises avec result summary ; observations→Tasks ; délai vérification improvement.
## 24. Classification de livraison
`defined` / `planned`, cible native ; aucune sandbox/runtime défini.
## 25. Critères d’acceptation
**Given** scénario/objectifs/participants, **When** exercice termine, **Then** observations, result summary et gaps sont attribués.

**Given** gap, **When** Task créée, **Then** source exercice et expected validation sont liés.

**Given** aucun modèle, **When** l’exercice est coordonné, **Then** toutes les étapes restent manuelles/déterministes.
## 26. Questions ouvertes
Quel objet porte Exercise ? Qui modifie observations vs results ? — requirements ci-dessus. `OPEN-010/013` restent ouvertes.
## 27. Consommateurs documentaires
Readiness, reports, improvement actions, parcours Phase 5, écrans Phase 6, objets Phase 7, permissions ultérieures.
