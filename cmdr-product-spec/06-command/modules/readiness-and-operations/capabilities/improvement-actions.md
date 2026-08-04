---
id: CAP-CMD-303
title: Improvement Actions
product: command
module: readiness-and-operations
owner: Command Product Lead
status: draft
delivery_status: defined
delivery_mode: planned
last_updated: 2026-08-04
requirement_ids: [REQ-PROD-009, REQ-PROD-013, REQ-PROD-021, REQ-OBJ-012]
open_decisions: [OPEN-013]
source-of-truth: canonical
---

# CAP-CMD-303 — Improvement Actions

## 1. Définition
Transforme Result, exercice, Incident review ou gap de Coverage en Task Command avec cause, owner, priorité, résultat attendu et vérification.

## 2. Problème utilisateur
Les leçons restent des recommandations sans owner ni échéance. Sans Task réutilisée, un objet concurrent apparaît ou la clôture n’est pas vérifiée.

## 3. Objectifs
Réutiliser Task, lier cause/résultat attendu, gérer owner/priority/blockers/due et exiger vérification ou exception de clôture.

## 4. Non-objectifs
Ne crée pas un backlog général, ne modifie pas la source, n’évalue pas un agent et ne finalise pas le schéma Task.

## 5. Propriétaire
Command possède la Task ; Result, Exercise, Readiness et Coverage restent propriétaires de leurs preuves sources.

## 6. Utilisateurs
Principal : Readiness Coordinator. Secondaires : Incident Commander, Task owner, Service Owner.

## 7. Conditions d’entrée
Source ou justification, action attendue, owner/team ou raison `unassigned`.

## 8. Entrées fonctionnelles
| Entrée | Source | Type fonctionnel | Requise | Fraîcheur | Si absente |
|---|---|---|---:|---|---|
| Improvement source | Result, Exercise, Coverage ou review | gap et cause | oui ou justification autonome | source datée | création refusée si cause non qualifiée |
| Expected outcome | Readiness Coordinator ou source | résultat vérifiable | oui | au moment de la création | Task refusée comme trop vague |
| Owner, due and priority | Command coordination | responsabilité et échéance | oui ou raison `unassigned` | version courante | Task `unassigned`, lacune visible |

## 9. Objets lus
| Objet | Propriétaire | Projection utilisée | Droit local |
|---|---|---|---|
| Result | Govern | outcome, verification et residual risk | consulter et relier en projection |
| Exercise / Readiness / Coverage | propriétaires fonctionnels | gap, observation et source | consulter et relier |
| Task | Command | owner, état, due et résultat | consulter et modifier |

## 10. Objets créés ou modifiés
| Objet | Opération | Propriétaire | Règle |
|---|---|---|---|
| Task | créer, assigner, prioriser, suivre, clore ou rouvrir | Command | classe 2, cause et résultat attendu obligatoires |
| Source relation | ajouter une référence | Object Linking Service | aucune mutation de la source |
| Result / Exercise / Coverage | aucune mutation | propriétaires sources | projections en lecture seule |

## 11. Fonctionnalités
Créer une Task depuis gap, assigner/prioriser/suivre, lier une preuve de validation, clore avec vérification/exception et rouvrir si invalidation.

## 12. Actions utilisateur
| Action | Rôle | Objet | Classe | Précondition | Résultat | Govern |
|---|---|---|---:|---|---|---|
| Créer | coordinateur | Task | 2 | source et expected outcome | open | OPEN-013 |
| Mettre à jour/assigner | owner/coordinateur | Task | 2 | version courante | updated | OPEN-013 |
| Soumettre vérification | owner | Task | 2 | outcome | verification | non |
| Clore ou rouvrir | verifier/coordinateur | Task | 2 | preuve ou raison | done/reopened | OPEN-013 |

## 13. Automatisation et IA
| Fonction | Manuel | Déterministe | Automatisable | IA possible | Alternative sans IA |
|---|---:|---:|---:|---:|---|
| Détecter un gap | oui | règles possibles | oui | suggestion attribuée | revue des Results/exercices/coverage |
| Dédupliquer les actions | oui | oui | oui | explication possible | recherche de Tasks et comparaison manuelle |
| Créer une Task | oui | validation | oui | brouillon attribué | formulaire manuel |
| Vérifier la clôture | oui | règles de preuve | workflow possible | résumé facultatif | vérification humaine et références sources |

## 14. États fonctionnels
`open`, `assigned`, `in-progress`, `blocked`, `verification`, `closed`, `reopened`, `cancelled`.

## 15. États d’interface
Partial montre preuve/source manquante ; Conflict conserve version ; Offline bloque clôture ; Permission denied masque la source protégée ; Stale invalide la vérification.

## 16. Sorties
| Sortie | Objet ou événement | Consommateur | Garantie |
|---|---|---|---|
| Improvement Task | Task Command | Work Queue et Readiness | cause, owner, due et expected outcome visibles |
| Verification outcome | événement Task | Readiness assessment et source | attribué, sourcé et sans mutation source |
| Reopen event | événement | owner et coordinateur | raison et preuve d’échec conservées |

## 17. Transitions
| Source | Déclencheur | Destination | Contexte transmis | Retour |
|---|---|---|---|---|
| Result, Exercise ou Coverage | créer une action | Task Coordination | source, cause et résultat attendu | source restaurée |
| Improvement Task | validation technique nécessaire | produit owner, Studio ou Settings | Task, source et critères de vérification | Task restaurée |
| Verification | échec | Task Coordination | Task, preuve et motif | Task reopened |

## 18. Dépendances
CAP-CMD-107, CAP-CMD-301, CAP-CMD-302, CAP-CMD-203 et Object Linking Service.

## 19. Source de vérité
Task et coordination : Command. Cause/evidence : propriétaire source.

## 20. Provenance et audit
Source, cause, résultat attendu, owner, changements, evidence, verifier et correlation ID.

## 21. Permissions fonctionnelles
Task manage, source read et future permission de vérification ; `OPEN-013` reportée.

## 22. Limites et erreurs
Source absente, expected outcome vague, evidence inaccessible, conflit ou refus empêchent une clôture confirmée.

## 23. Métriques
Gap→Task avec owner, clôture avec vérification et reopen après échec ; aucune cible définitive.

## 24. Classification de livraison
`defined` / `planned`, cible native ; preuve documentaire.

## 25. Critères d’acceptation
**Given** un gap Coverage, **When** une Task est créée, **Then** cause, owner, due et résultat attendu sont liés.

**Given** une vérification échouée, **When** le résultat est enregistré, **Then** la Task est reopened avec raison/evidence.

**Given** aucun modèle, **When** l’action est créée ou fermée, **Then** la voie manuelle complète reste disponible.

## 26. Questions ouvertes
Qui confirme la vérification et quelles exceptions permettent la clôture sans preuve ? — Requirement IDs ci-dessus ; `OPEN-013` reste ouverte.

## 27. Consommateurs documentaires
Readiness, Work Queue, Result/exercise reviews, parcours Phase 5, écrans Phase 6, objet Task Phase 7 et permissions ultérieures.