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
Les prochaines actions restent souvent dans des notes sans owner ni suivi. Sans Task coordonnée, le travail est oublié ou sa clôture modifie implicitement l’objet source.

## 3. Objectifs
Réutiliser Task, lier cause/résultat attendu, gérer owner/échéance/priorité/dépendances/blocage et éviter un objet Improvement Action concurrent.

## 4. Non-objectifs
Ne transforme pas toute notification en Task, ne modifie pas la source à la clôture, ne possède pas des Tasks spécialisées externes et ne finalise pas le schéma.

## 5. Propriétaire
Command possède la Task opérationnelle ; Shared Task Inbox reste consommateur/projection.

## 6. Utilisateurs
Principal : SOC Analyst. Secondaires : Incident Commander, Readiness Coordinator, Business Owner en lecture.

## 7. Conditions d’entrée
Contexte ou source, résultat attendu, owner/équipe ou raison `unassigned` et permission de gestion.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Task intent | utilisateur, Incident, Result ou exercice | travail et résultat attendu | oui | au moment de la création | création refusée sans résultat attendu |
| Source object | Object Linking Service | référence stable | non | résolution courante | Task autonome avec justification explicite |
| Dependencies | Tasks ou objets liés | relations | non | version courante | `unknown` ou `blocked` explicite |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Task | Command | owner, état, due, priority et dependencies | consulter et modifier |
| Source object | produit propriétaire | résumé et relation | consulter, relier et naviguer |
| Related Task | Command | état et dépendance | consulter et relier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Task | créer, assigner, prioriser, transitionner, fermer ou rouvrir | Command | classe 2, audit et résultat attendu |
| Source relation | ajouter ou retirer un lien | Object Linking Service | aucune mutation de l’objet source |
| Source object | aucune mutation | produit propriétaire | résultat de Task seulement référencé |

## 11. Fonctionnalités
Créer avec expected result/source, gérer owner/due/priority/state, lier dépendances/blocages, enregistrer résultat/motif et réutiliser pour improvement actions.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer Task | rôle autorisé | Task | 2 | intent et résultat | open | OPEN-013 |
| Modifier due/priority | owner/coordinateur | Task | 2 | version courante | mise à jour | OPEN-013 |
| Lier dépendance | owner | Task | 2 | objet accessible | relation | OPEN-013 |
| Clore ou rouvrir | owner/coordinateur | Task | 2 | résultat ou motif | done/reopened | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Créer une Task | oui | validation/déduplication | oui | brouillon attribué | formulaire manuel |
| Affecter et prioriser | oui | facteurs/règles possibles | oui | proposition seulement | sélection et priorité manuelles |
| Suivre échéance/blocage | oui | oui | oui | résumé facultatif | règles et notifications |
| Clore ou rouvrir | oui | contrôle de version | workflow possible | jamais silencieusement | action humaine avec résultat/motif |

## 14. États fonctionnels
`open`, `assigned`, `in-progress`, `blocked`, `pending-external`, `pending-decision`, `verification`, `done`, `cancelled`, `reopened`.

## 15. États d’interface
Partial montre source/dépendance manquante ; Offline bloque mutation ; Conflict conserve versions ; Permission denied masque la source ; Stale expose la version.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Task | objet Command | Queue, Mission Control et Readiness | source, owner, due et résultat attendu visibles |
| Task outcome | événement ou relation | objet source et Audit | résultat/motif sans mutation automatique de la source |
| Dependency update | relation typée | Work Queue | versionnée et permission-aware |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Result, exercice ou blocker | créer improvement Task | Task Coordination | source, cause et résultat attendu | objet source restauré |
| Task | investigation nécessaire | Investigate | source, contexte et prochaine action | Task reste Command |
| Task | décision nécessaire | Govern | Task, Incident, impact et action proposée | Task reste Command |

## 18. Dépendances
CAP-CMD-102, CAP-CMD-103, CAP-CMD-005, Object Linking Service, Notification Center et Task Inbox projection.

## 19. Source de vérité
Task et mutations : Command. Objet source : produit propriétaire. Relation : Object Linking Service.

## 20. Provenance et audit
Création, affectation, transition, relation et clôture enregistrent acteur, source, version, justification, résultat et correlation ID.

## 21. Permissions fonctionnelles
`perm.command.task.manage`, permission source en lecture et ABAC tenant/owner ; lecture Task atomique reportée.

## 22. Limites et erreurs
Intent sans résultat, source inaccessible, dépendance cyclique/inconnue, conflit, tenant incompatible ou refus empêchent une clôture trompeuse.

## 23. Métriques
Tasks avec owner/résultat attendu, âge blocked et clôtures avec résultat/motif ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire seulement.

## 25. Critères d’acceptation
**Given** un Result nécessitant amélioration, **When** une Task est créée, **Then** source, cause, owner et résultat attendu sont conservés.

**Given** une Task liée, **When** elle est clôturée, **Then** la source n’est pas modifiée automatiquement et l’outcome est audité.

**Given** aucun modèle IA, **When** la capability est utilisée, **Then** formulaires, règles et calculs déterministes suffisent.

## 26. Questions ouvertes
Quels attributs sont communs aux spécialisations et quel lien source est obligatoire ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Queue, Incident Detail, Readiness, Mission Control, parcours Phase 5, écrans Phase 6, objet Task Phase 7 et permissions ultérieures.