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
Coordonne un exercice opérationnel avec scénario, objectifs, participants, observations, résultats et actions d’amélioration, sans simulation technique Studio ou Sandbox.

## 2. Problème utilisateur
Les exercices hors plateforme produisent rarement des observations et actions traçables. Sans coordination, les leçons n’ont ni owner, ni résultat attendu.

## 3. Objectifs
Planifier scope/participants/objectifs, enregistrer observations/results, créer Tasks d’amélioration et distinguer exercice opérationnel et simulation technique.

## 4. Non-objectifs
N’orchestre pas Sandbox, ne teste pas Agent/règle, n’exécute pas Response Run et ne simule pas un Incident réel.

## 5. Propriétaire
Command possède la coordination et le record opérationnel ; Studio/Sandbox possèdent la simulation technique.

## 6. Utilisateurs
Principal : Readiness Coordinator. Secondaires : Incident Commander, Team Lead, Business Owner et Observers.

## 7. Conditions d’entrée
Scénario/objectifs, participants, tenant, autorisations et procédures liées.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Exercise definition | Readiness Coordinator | scénario, objectifs et scope | oui | version de planification | exercice non planifiable |
| Participants and roles | Platform Settings identity projection | principals, équipes et rôles | oui | résolution avant démarrage | état `partial`, démarrage bloqué si critique |
| Observations and outcomes | participants et observers | observations attribuées et résultat résumé | oui pour completion | au fil de l’exercice | état `observing`, completion bloquée sans résultat |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Operational Plan | Command Readiness | scénario, rôles et procédures | consulter et relier |
| Capability Readiness | Command | statut et gaps | consulter et comparer |
| Workflow / Simulation | Studio | référence et statut autorisé | consulter en projection |
| Principal / Team | Platform Settings | identité et rôle | consulter et sélectionner |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Exercise record | créer, planifier, démarrer, compléter ou superséder | Command Readiness | classe 2, version et audit |
| Observation | ajouter ou corriger | Command Readiness | auteur et timestamp obligatoires |
| Task | créer une amélioration | Command | classe 2, source exercice liée |
| Simulation / Workflow | aucune mutation | Studio | projection en lecture seule |

## 11. Fonctionnalités
Définir scénario/objectifs/participants/schedule, enregistrer observations attribuées, résumer outcome/gaps, créer Tasks et comparer Readiness avant/après conceptuellement.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer ou planifier | coordinateur | Exercise | 2 | scénario/objectifs | draft/planned | OPEN-013 |
| Démarrer ou terminer | coordinateur | Exercise | 2 | participants/scope | running/completed | OPEN-013 |
| Ajouter observation | participant | Observation | 2 | exercice actif | observation attribuée | non |
| Créer Task | coordinateur | Task | 2 | gap | Task liée | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Vérifier la complétude | oui | oui | oui | explication facultative | checklist versionnée |
| Notifier les participants | oui | règles de calendrier | oui | non nécessaire | invitations et notifications déterministes |
| Agréger les observations | oui | oui | oui | synthèse attribuée | lecture et synthèse manuelles |
| Créer les actions | oui | déduplication/validation | oui | brouillon Task | création manuelle de Task |

## 14. États fonctionnels
`draft`, `planned`, `ready`, `running`, `observing`, `completed`, `cancelled`, `superseded`.

## 15. États d’interface
Partial montre participants/sources unresolved ; Offline bloque start/complete ; Permission denied masque les observations protégées ; Stale signale plan/procédure périmé.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Exercise record | record Command | participants et Reporting | versionné, attribué et tenant-scoped |
| Observation set | collection d’événements | Readiness Coordinator | auteur, timestamp et objectif liés |
| Improvement Task | Task Command | Work Queue | source exercice et validation attendue conservées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Exercise Coordination | simulation technique nécessaire | Studio ou Sandbox owner | scénario, objectifs, scope et return origin | exercice restauré |
| Exercise Coordination | gap confirmé | Improvement Actions | exercice, observation, owner et résultat attendu | exercice restauré |
| Exercise Coordination | publication du résultat | Reporting Engine | snapshot, observations autorisées et audience | exercice inchangé |

## 18. Dépendances
CAP-CMD-301, CAP-CMD-303, CAP-CMD-304, Collaboration Service, Notification Center et Reporting Engine.

## 19. Source de vérité
Exercise record/observations : Command. Simulations/workflows : Studio. Identités : Settings.

## 20. Provenance et audit
Auteur et timestamp de chaque observation, version du scénario, participants, transitions, résultat et correlation ID.

## 21. Permissions fonctionnelles
Command read/coordinate et future permission de contribution participant ; `OPEN-010/013` restent ouvertes.

## 22. Limites et erreurs
Participant unresolved, procédure stale, exercice interrompu, conflit ou refus n’autorisent pas un état `completed` trompeur.

## 23. Métriques
Exercises avec result summary, observations→Tasks et délai de vérification des améliorations ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; aucune Sandbox ou runtime défini.

## 25. Critères d’acceptation
**Given** scénario/objectifs/participants, **When** l’exercice termine, **Then** observations, résultat et gaps sont attribués.

**Given** un gap, **When** une Task est créée, **Then** source exercice et validation attendue sont liées.

**Given** aucun modèle, **When** l’exercice est coordonné, **Then** toutes les étapes restent manuelles/déterministes.

## 26. Questions ouvertes
Quel objet porte Exercise et qui modifie observations versus results ? — Requirement IDs ci-dessus ; `OPEN-010/013` restent ouvertes.

## 27. Consommateurs documentaires
Readiness, reports, Improvement Actions, parcours Phase 5, écrans Phase 6, phase Objets et permissions ultérieures.