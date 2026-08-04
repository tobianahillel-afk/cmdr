---
id: CAP-CMD-107
title: Task Coordination
product: command
module: incidents-and-work-queue
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids:
  - REQ-PROD-009
  - REQ-PROD-013
  - REQ-PROD-021
  - REQ-OBJ-012
open_decisions:
  - OPEN-013
source-of-truth: canonical
---
# CAP-CMD-107 — Task Coordination

## 1. Définition
Crée et coordonne des Tasks opérationnelles liées à des objets sources avec owner, échéance, priorité, dépendances, blocage et résultat, sans modifier automatiquement la source.

## 2. Problème utilisateur
Les prochaines actions restent dans des notes sans owner ni suivi. Principal : SOC Analyst ; secondaires : Incident Commander, Readiness Coordinator, Business Owner. Sans capacité, le travail est oublié.

## 3. Objectifs
Réutiliser Task ; lier cause/résultat attendu ; gérer owner/échéance/priorité/dépendances/blocage ; éviter un objet Improvement Action concurrent.

## 4. Non-objectifs
Ne pas transformer toute notification en Task, modifier la source à la clôture, posséder les Tasks spécialisées externes ou finaliser le schéma.

## 5. Propriétaire
Command possède la Task opérationnelle ; Shared Task Inbox reste un consommateur/projection.

## 6. Utilisateurs
Principal : SOC Analyst. Secondaires : Incident Commander, Readiness Coordinator, Business Owner en lecture.

## 7. Conditions d’entrée
Contexte ou source, résultat attendu, owner/équipe ou raison unassigned, permission de gestion.

## 8. Entrées fonctionnelles
| Entrée | Source | Type | Requise | Fraîcheur | Si absente |
|---|---|---|---|---|---|
| Task intent | utilisateur/Incident/Result/exercice | travail attendu | oui | création | refuser sans résultat attendu |
| Source object | Object Linking Service | référence | non | courante | Task autonome justifiée |
| Dependencies | Tasks/objects | relations | non | version courante | unknown/blocked explicite |

## 9. Objets lus
| Objet | Owner | Projection | Droit local |
|---|---|---|---|
| Task | Command | owner, state, due, priority, dependencies | lecture/modification |
| Source object | produit propriétaire | résumé/relation | projection |

## 10. Objets créés ou modifiés
| Objet | Opération | Owner | Règle |
|---|---|---|---|
| Task | créer, assigner, prioriser, transitionner, fermer | Command | classe 2, audit |
| Source relation | ajouter lien | Object Linking Service | aucune mutation source |

## 11. Fonctionnalités
Créer avec expected result/source ; gérer owner/due/priority/state ; lier dépendances/blocages ; enregistrer résultat/motif ; réutiliser pour improvement actions.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer Task | autorisé | Task | 2 | intent/résultat | open | OPEN-013 |
| Modifier due/priority | owner/coordinateur | Task | 2 | version courante | mise à jour | OPEN-013 |
| Lier dépendance | owner | Task | 2 | objet accessible | relation | OPEN-013 |
| Clore/réouvrir | owner/coordinateur | Task | 2 | résultat/motif | done/reopened | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Humain | Règle | Moteur | Workflow | Agent | Govern | Sans IA |
|---|---|---|---|---|---|---|---|
| Créer/proposer Task | validation humaine | oui | dedup/validation | possible | brouillon | OPEN-013 | formulaire manuel |
| Suivre échéance/blocage | oui | oui | calcul déterministe | possible | résumé | non | règles/notifications |

## 14. États fonctionnels
`open`, `assigned`, `in-progress`, `blocked`, `pending-external`, `pending-decision`, `verification`, `done`, `cancelled`, `reopened`.

## 15. États d’interface
Partial montre source/dependency manquante ; Offline bloque mutation ; conflit garde les versions ; Permission denied ne révèle pas la source ; Stale indique version. Rendu DS.

## 16. Sorties
| Sortie | Objet/événement | Consommateur | Garantie |
|---|---|---|---|
| Task | objet Command | Queue, Mission Control, Readiness | source, owner, due, expected result |
| Task outcome | event/relation | objet source/audit | résultat/motif sans mutation source |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte | Retour |
|---|---|---|---|---|
| Result/exercise/blocker | improvement Task | Task Coordination | source, cause, expected result | source |
| Task | investigation/decision nécessaire | Investigate/Govern | source, contexte, next action | Task reste Command |

## 18. Dépendances
CAP-CMD-102, 103, 005, Object Linking Service, Notification Center, Task Inbox projection.

## 19. Source de vérité
Task et mutations : Command. Objet source : produit propriétaire. Relation : service de liaison.

## 20. Provenance et audit
Création, affectation, transition, relation et clôture enregistrent acteur, source, version, justification, résultat et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.task.manage`, permission source en lecture, ABAC tenant/owner. Lecture Task atomique à formaliser.

## 22. Limites et erreurs
Intent sans résultat, source inaccessible, dépendance cyclique/inconnue, conflit de version, tenant incompatible ou refus empêchent une clôture/mutation trompeuse.

## 23. Métriques
Tasks avec owner/expected result ; âge blocked ; clôtures avec résultat/motif. Aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** un Result nécessitant amélioration, **When** une Task est créée, **Then** source, cause, owner et expected result sont conservés.

**Given** clôture d’une Task liée, **When** elle est validée, **Then** la source n’est pas modifiée automatiquement et l’outcome est audité.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** formulaires, règles et calculs déterministes suffisent.

## 26. Questions ouvertes
Quels attributs sont communs aux spécialisations et quel lien source est obligatoire ? — REQ-PROD-009, REQ-PROD-013, REQ-PROD-021, REQ-OBJ-012. `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Queue, Incident Detail, Readiness, Mission Control, parcours Phase 5, écrans Phase 6, objet Task Phase 7 et permissions ultérieures.
